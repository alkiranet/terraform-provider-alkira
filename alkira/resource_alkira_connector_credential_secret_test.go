package alkira

import (
	"context"
	"net/http"
	"strings"
	"testing"

	"github.com/alkiranet/alkira-client-go/alkira"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// connectorCredentialSecretFields are attributes whose value the provider
// stores in the credential service. The API returns only the credential ID.
var connectorCredentialSecretFields = []string{
	"activation_code",
	"username",
	"password",
	"account_key",
	"akamai_bgp_authentication_key",
}

// AK-76636: every connector secret that sits next to a generated
// credential_id must suppress its diff, or an imported connector plans a
// change that the backend never applies.
func TestConnectorCredentialSecretFieldsSuppressDiff(t *testing.T) {
	var walk func(resource, path string, s map[string]*schema.Schema)
	walk = func(resource, path string, s map[string]*schema.Schema) {
		if cred, ok := s["credential_id"]; ok && cred.Computed && !cred.Optional {
			for _, name := range connectorCredentialSecretFields {
				if secret, ok := s[name]; ok {
					assert.NotNil(t, secret.DiffSuppressFunc,
						"%s: %s%s has no DiffSuppressFunc", resource, path, name)
				}
			}
		}
		for name, sch := range s {
			if r, ok := sch.Elem.(*schema.Resource); ok {
				walk(resource, path+name+".", r.Schema)
			}
		}
	}

	checked := 0
	for name, r := range Provider().ResourcesMap {
		if strings.HasPrefix(name, "alkira_connector_") {
			walk(name, "", r.Schema)
			checked++
		}
	}
	require.NotZero(t, checked)
}

type credentialSecretDiffCase struct {
	name     string
	resource string
	state    map[string]string
	config   map[string]interface{}
	key      string
	changed  bool
}

func TestConnectorCredentialSecretDiff(t *testing.T) {
	vmwareConfig := func(code string) map[string]interface{} {
		return map[string]interface{}{"virtual_edge": []interface{}{
			map[string]interface{}{"hostname": "edge1", "activation_code": code},
		}}
	}
	vmwareState := func(code string) map[string]string {
		s := map[string]string{
			"id":                             "4224",
			"virtual_edge.#":                 "1",
			"virtual_edge.0.id":              "2147",
			"virtual_edge.0.hostname":        "edge1",
			"virtual_edge.0.credential_id":   "cred-1",
			"virtual_edge.0.activation_code": code,
		}
		if code == "" {
			delete(s, "virtual_edge.0.activation_code")
		}
		return s
	}

	cases := []credentialSecretDiffCase{
		{
			name:     "vmware import",
			resource: "alkira_connector_vmware_sdwan",
			state:    vmwareState(""),
			config:   vmwareConfig("AAAA-BBBB"),
			key:      "virtual_edge.0.activation_code",
		},
		{
			name:     "vmware edit with credential",
			resource: "alkira_connector_vmware_sdwan",
			state:    vmwareState("AAAA-BBBB"),
			config:   vmwareConfig("CCCC-DDDD"),
			key:      "virtual_edge.0.activation_code",
		},
		{
			name:     "vmware new edge",
			resource: "alkira_connector_vmware_sdwan",
			state:    vmwareState(""),
			config: map[string]interface{}{"virtual_edge": []interface{}{
				map[string]interface{}{"hostname": "edge1", "activation_code": "AAAA-BBBB"},
				map[string]interface{}{"hostname": "edge2", "activation_code": "EEEE-FFFF"},
			}},
			key:     "virtual_edge.1.activation_code",
			changed: true,
		},
		{
			name:     "cisco import",
			resource: "alkira_connector_cisco_sdwan",
			state: map[string]string{
				"id":                    "1",
				"vedge.#":               "1",
				"vedge.0.id":            "11",
				"vedge.0.hostname":      "vedge1",
				"vedge.0.credential_id": "cred-1",
			},
			config: map[string]interface{}{"vedge": []interface{}{
				map[string]interface{}{"hostname": "vedge1", "username": "admin", "password": "secret"},
			}},
			key: "vedge.0.password",
		},
		{
			name:     "cisco import username",
			resource: "alkira_connector_cisco_sdwan",
			state: map[string]string{
				"id":                    "1",
				"vedge.#":               "1",
				"vedge.0.id":            "11",
				"vedge.0.hostname":      "vedge1",
				"vedge.0.credential_id": "cred-1",
			},
			config: map[string]interface{}{"vedge": []interface{}{
				map[string]interface{}{"hostname": "vedge1", "username": "admin", "password": "secret"},
			}},
			key: "vedge.0.username",
		},
		{
			name:     "fortinet import password",
			resource: "alkira_connector_fortinet_sdwan",
			state: map[string]string{
				"id":                       "1",
				"wan_edge.#":               "1",
				"wan_edge.0.id":            "21",
				"wan_edge.0.hostname":      "fgt1",
				"wan_edge.0.credential_id": "cred-1",
			},
			config: map[string]interface{}{"wan_edge": []interface{}{
				map[string]interface{}{"hostname": "fgt1", "password": "secret"},
			}},
			key: "wan_edge.0.password",
		},
		{
			name:     "fortinet import default username",
			resource: "alkira_connector_fortinet_sdwan",
			state: map[string]string{
				"id":                       "1",
				"wan_edge.#":               "1",
				"wan_edge.0.id":            "21",
				"wan_edge.0.hostname":      "fgt1",
				"wan_edge.0.credential_id": "cred-1",
			},
			config: map[string]interface{}{"wan_edge": []interface{}{
				map[string]interface{}{"hostname": "fgt1", "password": "secret"},
			}},
			key: "wan_edge.0.username",
		},
		{
			name:     "aruba import",
			resource: "alkira_connector_aruba_edge",
			state: map[string]string{
				"id":                        "1",
				"instances.#":               "1",
				"instances.0.id":            "31",
				"instances.0.name":          "aruba1",
				"instances.0.credential_id": "cred-1",
			},
			config: map[string]interface{}{"instances": []interface{}{
				map[string]interface{}{"name": "aruba1", "account_key": "key"},
			}},
			key: "instances.0.account_key",
		},
		{
			name:     "akamai import",
			resource: "alkira_connector_akamai_prolexic",
			state:    map[string]string{"id": "1", "credential_id": "cred-1"},
			config:   map[string]interface{}{"akamai_bgp_authentication_key": "key"},
			key:      "akamai_bgp_authentication_key",
		},
		{
			name:     "akamai edit with credential",
			resource: "alkira_connector_akamai_prolexic",
			state: map[string]string{
				"id":                            "1",
				"credential_id":                 "cred-1",
				"akamai_bgp_authentication_key": "old",
			},
			config:  map[string]interface{}{"akamai_bgp_authentication_key": "new"},
			key:     "akamai_bgp_authentication_key",
			changed: true,
		},
		{
			name:     "akamai without credential",
			resource: "alkira_connector_akamai_prolexic",
			state:    map[string]string{"id": "1"},
			config:   map[string]interface{}{"akamai_bgp_authentication_key": "key"},
			key:      "akamai_bgp_authentication_key",
			changed:  true,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			r := Provider().ResourcesMap[tc.resource]
			require.NotNil(t, r)

			state := &terraform.InstanceState{ID: tc.state["id"], Attributes: tc.state}
			diff, err := r.Diff(context.Background(), state,
				terraform.NewResourceConfigRaw(tc.config), &alkira.AlkiraClient{})
			require.NoError(t, err)

			changed := false
			if diff != nil {
				_, changed = diff.Attributes[tc.key]
			}
			assert.Equal(t, tc.changed, changed, "diff on %s", tc.key)
		})
	}
}

// AK-76636: Read sets credential_id so that an imported Akamai Prolexic
// connector can suppress the diff on its authentication key.
func TestConnectorAkamaiProlexicReadSetsCredentialId(t *testing.T) {
	r := Provider().ResourcesMap["alkira_connector_akamai_prolexic"]
	d := schema.TestResourceDataRaw(t, r.Schema, nil)
	d.SetId("42")

	c := mcgwConnectorClient(t, func(req *http.Request) (int, string) {
		switch {
		case req.Method == "GET" && req.URL.Path == "/tenantnetworks/1/akamai-prolexic-connectors/42":
			return 200, `{"id":42,"name":"akamai","cxp":"US-WEST","size":"SMALL",` +
				`"credentialId":"cred-1","segments":["segment1"]}`
		case req.Method == "GET" && strings.HasPrefix(req.URL.Path, "/tenantnetworks/1/segments"):
			return 200, `[{"id":1,"name":"segment1"}]`
		default:
			t.Fatalf("unexpected request: %s %s", req.Method, req.URL.Path)
			return 500, ""
		}
	})

	diags := resourceConnectorAkamaiProlexicRead(context.Background(), d, c)
	require.False(t, diags.HasError(), "%v", diags)
	assert.Equal(t, "cred-1", d.Get("credential_id"))
}

// AK-76636: after an import the authentication key is not in state. An
// update must keep the existing credential and not create one with an empty
// key.
func TestGenerateConnectorAkamaiProlexicRequestKeepsCredentialWithoutKey(t *testing.T) {
	r := Provider().ResourcesMap["alkira_connector_akamai_prolexic"]
	d := schema.TestResourceDataRaw(t, r.Schema, map[string]interface{}{
		"name":       "akamai",
		"segment_id": "1",
	})
	d.SetId("42")
	require.NoError(t, d.Set("credential_id", "cred-1"))

	c := mcgwConnectorClient(t, func(req *http.Request) (int, string) {
		switch {
		case req.Method == "GET" && req.URL.Path == "/tenantnetworks/1/segments/1":
			return 200, `{"id":1,"name":"segment1"}`
		default:
			t.Fatalf("unexpected request: %s %s", req.Method, req.URL.Path)
			return 500, ""
		}
	})

	connector, err := generateConnectorAkamaiProlexicRequest(d, c)
	require.NoError(t, err)
	assert.Equal(t, "cred-1", connector.CredentialId)
}

// AK-76636: on import no instance is in state yet, so Read must add every
// instance that the API returns.
func TestSetVirtualEdgeAddsAllUntrackedInstances(t *testing.T) {
	r := Provider().ResourcesMap["alkira_connector_vmware_sdwan"]
	d := schema.TestResourceDataRaw(t, r.Schema, nil)

	setVirtualEdge(d, &alkira.ConnectorVmwareSdwan{Instances: []alkira.VmwareSdwanInstance{
		{Id: 2147, HostName: "edge1", CredentialId: "cred-1"},
		{Id: 2148, HostName: "edge2", CredentialId: "cred-2"},
	}})

	edges := d.Get("virtual_edge").([]interface{})
	require.Len(t, edges, 2)
	assert.Equal(t, 2148, edges[1].(map[string]interface{})["id"])
	assert.Equal(t, "cred-2", edges[1].(map[string]interface{})["credential_id"])
}

func TestSetVedgeAddsAllUntrackedInstances(t *testing.T) {
	r := Provider().ResourcesMap["alkira_connector_cisco_sdwan"]
	d := schema.TestResourceDataRaw(t, r.Schema, nil)

	setVedge(d, &alkira.ConnectorCiscoSdwan{CiscoEdgeInfo: []alkira.CiscoSdwanEdgeInfo{
		{Id: 11, HostName: "vedge1", CredentialId: "cred-1"},
		{Id: 12, HostName: "vedge2", CredentialId: "cred-2"},
	}})

	vedges := d.Get("vedge").([]interface{})
	require.Len(t, vedges, 2)
	assert.Equal(t, 12, vedges[1].(map[string]interface{})["id"])
	assert.Equal(t, "cred-2", vedges[1].(map[string]interface{})["credential_id"])
}

func TestSetWanEdgeAddsAllUntrackedInstances(t *testing.T) {
	r := Provider().ResourcesMap["alkira_connector_fortinet_sdwan"]
	d := schema.TestResourceDataRaw(t, r.Schema, nil)

	setWanEdge(d, &alkira.ConnectorFortinetSdwan{Instances: []alkira.ConnectorFortinetSdwanInstance{
		{Id: 21, HostName: "fgt1", CredentialId: "cred-1", Version: "7.2.4"},
		{Id: 22, HostName: "fgt2", CredentialId: "cred-2", Version: "7.2.4"},
	}})

	edges := d.Get("wan_edge").([]interface{})
	require.Len(t, edges, 2)
	assert.Equal(t, 22, edges[1].(map[string]interface{})["id"])
	assert.Equal(t, "cred-2", edges[1].(map[string]interface{})["credential_id"])
}

// AK-76636: version is not a secret. Read must set it from the API on
// import, or the plan shows a diff on it.
func TestSetWanEdgeSetsVersionOnUntrackedInstance(t *testing.T) {
	r := Provider().ResourcesMap["alkira_connector_fortinet_sdwan"]
	d := schema.TestResourceDataRaw(t, r.Schema, nil)

	setWanEdge(d, &alkira.ConnectorFortinetSdwan{Instances: []alkira.ConnectorFortinetSdwanInstance{
		{Id: 21, HostName: "fgt1", CredentialId: "cred-1", Version: "7.2.4"},
	}})

	edges := d.Get("wan_edge").([]interface{})
	require.Len(t, edges, 1)
	assert.Equal(t, "7.2.4", edges[0].(map[string]interface{})["version"])
}
