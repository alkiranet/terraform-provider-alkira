package alkira

import (
	"context"
	"fmt"
	"net/http"
	"regexp"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/stretchr/testify/require"
)

var segmentGetPath = regexp.MustCompile(`/segments/([0-9]+)$`)

// credentialOrderClient answers GET /segments/<id> and answers a POST to
// rejectedCredential, when set, with a 400 and records it in rejected. Any
// other request fails the test, so a credential created out of order shows
// up as a failure naming the credential endpoint.
func credentialOrderClient(t *testing.T, rejectedCredential string, rejected *bool) interface{} {
	return createMockAlkiraClient(t, http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		if match := segmentGetPath.FindStringSubmatch(req.URL.Path); req.Method == http.MethodGet && match != nil {
			w.Header().Set("Content-Type", "application/json")
			fmt.Fprintf(w, `{"id": %s, "name": "seg-%s"}`, match[1], match[1])
			return
		}

		if rejectedCredential != "" && req.Method == http.MethodPost && strings.HasSuffix(req.URL.Path, "/credentials/"+rejectedCredential) {
			*rejected = true
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		t.Errorf("unexpected request before the apply failed: %s %s", req.Method, req.URL.Path)
		w.WriteHeader(http.StatusInternalServerError)
	}))
}

// TestSegmentLookupsPrecedeCredentialCreation pins the order every service
// with implicit credentials follows: segment lookups, then the service
// credential, then management, instance and license credentials. A rejected
// segment_id creates no credential, and a rejected service credential
// creates no other credential, so neither leaves an orphan in the tenant.
func TestSegmentLookupsPrecedeCredentialCreation(t *testing.T) {
	const badSegment = "ak74389-seg-a"

	segmentOptions := func(segmentId string) []interface{} {
		return []interface{}{
			map[string]interface{}{"segment_id": segmentId, "zone_name": "zone-a", "groups": []interface{}{"group-a"}},
		}
	}
	checkpointManagementServer := func(segmentId string) []interface{} {
		return []interface{}{
			map[string]interface{}{"configuration_mode": "AUTOMATED", "credential_id": "", "password": "secret", "segment_id": segmentId},
		}
	}
	checkpointInstance := []interface{}{
		map[string]interface{}{"name": "ak74389-chkp-a", "sic_key": "key", "credential_id": ""},
	}
	fortinetInstance := []interface{}{
		map[string]interface{}{"name": "ak74389-ftnt-a", "serial_number": "sn", "license_key": "key", "credential_id": ""},
	}
	panInstance := func(segmentId string) []interface{} {
		instance := map[string]interface{}{"name": "ak74389-pan-a", "auth_code": "code", "auth_key": "key", "credential_id": ""}
		if segmentId != "" {
			instance["global_protect_segment_options"] = []interface{}{
				map[string]interface{}{"segment_id": segmentId},
			}
		}
		return []interface{}{instance}
	}
	pan := func(extra map[string]interface{}) map[string]interface{} {
		v := map[string]interface{}{
			"name":         "ak74389-pan",
			"pan_username": "admin",
			"pan_password": "secret",
			"instance":     panInstance(""),
		}
		for k, e := range extra {
			v[k] = e
		}
		return v
	}

	cases := []struct {
		name string
		// rejectedCredential names a credential type the mock answers with
		// a 400; empty means the case expects a segment_id rejection.
		rejectedCredential string
		resource           *schema.Resource
		update             bool
		values             map[string]interface{}
		// computed holds Computed attributes an update reads from state.
		computed map[string]interface{}
	}{
		{
			name:     "checkpoint create, segment_id",
			resource: resourceAlkiraCheckpoint(),
			values: map[string]interface{}{
				"name":              "ak74389-chkp",
				"password":          "secret",
				"segment_id":        badSegment,
				"management_server": checkpointManagementServer("1"),
				"instance":          checkpointInstance,
			},
		},
		{
			name:     "checkpoint create, segment_options segment_id",
			resource: resourceAlkiraCheckpoint(),
			values: map[string]interface{}{
				"name":              "ak74389-chkp",
				"password":          "secret",
				"segment_id":        "1",
				"segment_options":   segmentOptions(badSegment),
				"management_server": checkpointManagementServer("1"),
				"instance":          checkpointInstance,
			},
		},
		{
			name:     "checkpoint create, management_server segment_id",
			resource: resourceAlkiraCheckpoint(),
			values: map[string]interface{}{
				"name":              "ak74389-chkp",
				"password":          "secret",
				"segment_id":        "1",
				"management_server": checkpointManagementServer(badSegment),
				"instance":          checkpointInstance,
			},
		},
		{
			name:               "checkpoint create, service credential rejected",
			rejectedCredential: "chkp-fw",
			resource:           resourceAlkiraCheckpoint(),
			values: map[string]interface{}{
				"name":              "ak74389-chkp",
				"password":          "weak",
				"segment_id":        "1",
				"management_server": checkpointManagementServer("1"),
				"instance":          checkpointInstance,
			},
		},
		{
			name:     "checkpoint update with a password change, segment_id",
			resource: resourceAlkiraCheckpoint(),
			update:   true,
			values: map[string]interface{}{
				"name":       "ak74389-chkp",
				"password":   "new-secret",
				"segment_id": badSegment,
				"instance":   checkpointInstance,
			},
		},
		{
			name:               "checkpoint update with a password change, service credential rejected",
			rejectedCredential: "chkp-fw",
			resource:           resourceAlkiraCheckpoint(),
			update:             true,
			values: map[string]interface{}{
				"name":       "ak74389-chkp",
				"password":   "weak",
				"segment_id": "1",
				"instance":   checkpointInstance,
			},
		},
		{
			name:     "fortinet update with a password change, segment_ids",
			resource: resourceAlkiraServiceFortinet(),
			update:   true,
			values: map[string]interface{}{
				"name":                         "ak74389-ftnt",
				"username":                     "admin",
				"password":                     "new-secret",
				"license_type":                 "BRING_YOUR_OWN",
				"management_server_segment_id": "1",
				"segment_ids":                  []interface{}{badSegment},
				"instances":                    fortinetInstance,
			},
			computed: map[string]interface{}{"credential_id": "existing-cred-id", "credential_name": "existing-cred-name"},
		},
		{
			name:     "fortinet create, segment_ids",
			resource: resourceAlkiraServiceFortinet(),
			values: map[string]interface{}{
				"name":                         "ak74389-ftnt",
				"username":                     "admin",
				"password":                     "secret",
				"license_type":                 "BRING_YOUR_OWN",
				"management_server_segment_id": "1",
				"segment_ids":                  []interface{}{badSegment},
				"instances":                    fortinetInstance,
			},
		},
		{
			name:     "fortinet create, segment_options segment_id",
			resource: resourceAlkiraServiceFortinet(),
			values: map[string]interface{}{
				"name":                         "ak74389-ftnt",
				"username":                     "admin",
				"password":                     "secret",
				"license_type":                 "BRING_YOUR_OWN",
				"management_server_segment_id": "1",
				"segment_options":              segmentOptions(badSegment),
				"instances":                    fortinetInstance,
			},
		},
		{
			name:               "fortinet create, service credential rejected",
			rejectedCredential: "ftntfw",
			resource:           resourceAlkiraServiceFortinet(),
			values: map[string]interface{}{
				"name":                         "ak74389-ftnt",
				"license_type":                 "BRING_YOUR_OWN",
				"management_server_segment_id": "1",
				"instances":                    fortinetInstance,
			},
		},
		{
			name:     "pan create, segment_options segment_id",
			resource: resourceAlkiraServicePan(),
			values:   pan(map[string]interface{}{"segment_options": segmentOptions(badSegment)}),
		},
		{
			name:     "pan create, global_protect_segment_options segment_id",
			resource: resourceAlkiraServicePan(),
			values: pan(map[string]interface{}{
				"global_protect_segment_options": []interface{}{map[string]interface{}{"segment_id": badSegment}},
			}),
		},
		{
			name:     "pan create, instance global_protect_segment_options segment_id",
			resource: resourceAlkiraServicePan(),
			values:   pan(map[string]interface{}{"instance": panInstance(badSegment)}),
		},
		{
			name:     "pan update with a password change, global_protect_segment_options segment_id",
			resource: resourceAlkiraServicePan(),
			update:   true,
			values: pan(map[string]interface{}{
				"pan_password":                   "new-secret",
				"global_protect_segment_options": []interface{}{map[string]interface{}{"segment_id": badSegment}},
			}),
			computed: map[string]interface{}{"pan_credential_id": "existing-cred-id", "pan_credential_name": "existing-cred-name"},
		},
		{
			name:               "pan create, service credential rejected",
			rejectedCredential: "pan",
			resource:           resourceAlkiraServicePan(),
			values:             pan(nil),
		},
		{
			name:     "cisco ftdv create, firepower_management_center segment_id",
			resource: resourceAlkiraServiceCiscoFTDv(),
			values: map[string]interface{}{
				"name": "ak74389-ftdv",
				"firepower_management_center": []interface{}{
					map[string]interface{}{"server_ip": "10.0.0.1", "username": "admin", "password": "secret", "credential_id": "", "segment_id": badSegment},
				},
				"instance": []interface{}{
					map[string]interface{}{"hostname": "ak74389-ftdv-a", "admin_password": "secret", "fmc_registration_key": "key", "ftdv_nat_id": "nat", "credential_id": ""},
				},
			},
		},
		{
			name:     "bluecat create, segment_ids",
			resource: resourceAlkiraBluecat(),
			values: map[string]interface{}{
				"name":        "ak74389-bluecat",
				"segment_ids": []interface{}{badSegment},
				"instance": []interface{}{
					map[string]interface{}{
						"type": "BDDS",
						"bdds_options": []interface{}{
							map[string]interface{}{"hostname": "bdds-a", "client_id": "id", "activation_key": "key", "license_credential_id": ""},
						},
					},
				},
			},
		},
		{
			name:     "f5 create, segment_ids",
			resource: resourceAlkiraF5LoadBalancer(),
			values: map[string]interface{}{
				"name":        "ak74389-f5",
				"segment_ids": []interface{}{badSegment},
				"instance": []interface{}{
					map[string]interface{}{"name": "ak74389-f5-a", "license_type": "PAY_AS_YOU_GO", "f5_username": "admin", "f5_password": "secret", "credential_id": ""},
				},
			},
		},
		{
			name:     "f5 create, segment_options segment_id",
			resource: resourceAlkiraF5LoadBalancer(),
			values: map[string]interface{}{
				"name":            "ak74389-f5",
				"segment_options": []interface{}{map[string]interface{}{"segment_id": badSegment, "elb_nic_count": 1}},
				"instance": []interface{}{
					map[string]interface{}{"name": "ak74389-f5-a", "license_type": "PAY_AS_YOU_GO", "f5_username": "admin", "f5_password": "secret", "credential_id": ""},
				},
			},
		},
		{
			name:     "infoblox create, segment_ids",
			resource: resourceAlkiraInfoblox(),
			values: map[string]interface{}{
				"name":          "ak74389-infoblox",
				"shared_secret": "secret",
				"segment_ids":   []interface{}{badSegment},
			},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			var rejected bool
			client := credentialOrderClient(t, tc.rejectedCredential, &rejected)

			// TestResourceDataRaw builds a diff from the config, which
			// GetChange and HasChanges read; d.Set would not reach them.
			d := schema.TestResourceDataRaw(t, tc.resource.Schema, tc.values)
			for k, v := range tc.computed {
				require.NoError(t, d.Set(k, v), k)
			}

			var diags diag.Diagnostics
			if tc.update {
				d.SetId("1")
				diags = tc.resource.UpdateContext(context.Background(), d, client)
			} else {
				diags = tc.resource.CreateContext(context.Background(), d, client)
			}

			require.True(t, diags.HasError(), "expected the apply to fail")

			if tc.rejectedCredential == "" {
				require.Contains(t, diags[0].Summary, badSegment)
			} else {
				require.True(t, rejected, "expected the apply to reach the %s credential", tc.rejectedCredential)
			}
		})
	}
}
