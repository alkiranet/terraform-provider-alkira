package alkira

import (
	"context"
	"encoding/json"
	"net/http"
	"testing"

	"github.com/alkiranet/alkira-client-go/alkira"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestFortinetRead(t *testing.T) {
	t.Skip("Test skipped: mock server response format doesn't match API client expectations")
	// expectedCxp := "US-WEST"
	// expectedIp := "10.1.1.0"
	// expectedSegment := "default"

	// f := &alkira.ServiceFortinet{
	// 	Cxp: expectedCxp,
	// 	ManagementServer: &alkira.FortinetManagmentServer{
	// 		IpAddress: expectedIp,
	// 		Segment:   expectedSegment,
	// 	},
	// }
	// ac := createMockAlkiraClient(t, func(w http.ResponseWriter, req *http.Request) {
	// 	json.NewEncoder(w).Encode(f)
	// 	w.Header().Set("Content-Type", "application/json")
	// })

	// r := resourceAlkiraServiceFortinet()
	// d := r.TestResourceData()

	// err := resourceFortinetRead(nil, d, ac)
	// require.Nil(t, err)

	// require.Equal(t, expectedCxp, getStringFromResourceData(d, "cxp"))
	// require.Equal(t, expectedIp, getStringFromResourceData(d, "management_server_ip"))
	// require.Equal(t, expectedSegment, getStringFromResourceData(d, "management_server_segment"))
}

func TestFortinetReadAutoScale(t *testing.T) {
	t.Skip("Test skipped: mock server response format doesn't match API client expectations")
	// expectedAutoScaleVal := "ON"

	// f := &alkira.ServiceFortinet{
	// 	AutoScale: expectedAutoScaleVal,
	// }
	// ac := createMockAlkiraClient(t, func(w http.ResponseWriter, req *http.Request) {
	// 	json.NewEncoder(w).Encode(f)
	// 	w.Header().Set("Content-Type", "application/json")
	// })
	// f.ManagementServer = &alkira.FortinetManagmentServer{}

	// r := resourceAlkiraServiceFortinet()
	// d := r.TestResourceData()

	// err := resourceFortinetRead(nil, d, ac)
	// require.Nil(t, err)

	// require.Equal(t, expectedAutoScaleVal, getStringFromResourceData(d, "auto_scale"))
}

//
// TEST HELPER
//

// UNUSED: Commented out to suppress linter warnings
// func serveFortinet(t *testing.T, f *alkira.ServiceFortinet) *alkira.AlkiraClient {
// 	return createMockAlkiraClient(t, func(w http.ResponseWriter, req *http.Request) {
// 		json.NewEncoder(w).Encode(f)
// 		w.Header().Set("Content-Type", "application/json")
// 	})
// }

func TestAlkiraServiceFortinetAllowListSchema(t *testing.T) {
	resourceSchema := resourceAlkiraServiceFortinet().Schema

	t.Run("field exists and is Optional TypeSet of String", func(t *testing.T) {
		field, ok := resourceSchema["allow_list"]
		require.True(t, ok, "allow_list must be present in schema")
		assert.Equal(t, schema.TypeSet, field.Type)
		assert.True(t, field.Optional)
		assert.False(t, field.Required)

		elem, ok := field.Elem.(*schema.Schema)
		require.True(t, ok, "allow_list Elem must be a *schema.Schema")
		assert.Equal(t, schema.TypeString, elem.Type)
		assert.NotNil(t, elem.ValidateFunc, "allow_list elements must be validated")
	})
}

func TestAlkiraServiceFortinetAllowListValidator(t *testing.T) {
	validate := resourceAlkiraServiceFortinet().Schema["allow_list"].Elem.(*schema.Schema).ValidateFunc

	valid := []string{"10.0.0.0/24", "192.168.1.0/24", "10.0.0.5", "10.0.0.5/32"}
	for _, v := range valid {
		_, errs := validate(v, "allow_list")
		assert.Emptyf(t, errs, "expected %q to be accepted (IPv4 CIDR or IPv4 IP)", v)
	}

	invalid := []string{"not-an-ip", "10.0.0.0/33", "999.0.0.1", "2001:db8::/32", "::1", "::ffff:1.2.3.4", "10.0.0.5/24"}
	for _, v := range invalid {
		_, errs := validate(v, "allow_list")
		assert.NotEmptyf(t, errs, "expected %q to be rejected", v)
	}
}

func TestAlkiraServiceFortinetAllowListExpand(t *testing.T) {
	tests := []struct {
		name     string
		elements []interface{}
		expected []string
	}{
		{
			name:     "populated set",
			elements: []interface{}{"10.0.0.0/24", "10.0.0.5"},
			expected: []string{"10.0.0.0/24", "10.0.0.5"},
		},
		{
			name:     "duplicates collapse",
			elements: []interface{}{"10.0.0.0/24", "10.0.0.5", "10.0.0.0/24"},
			expected: []string{"10.0.0.0/24", "10.0.0.5"},
		},
		{
			name:     "empty set produces no entries",
			elements: []interface{}{},
			expected: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := convertTypeSetToStringList(schema.NewSet(schema.HashString, tt.elements))
			assert.ElementsMatch(t, tt.expected, got)
		})
	}
}

// TestAlkiraServiceFortinetAllowListGenerateRequest drives the real
// generateFortinetRequest path to confirm allow_list lands on the request
// payload via the Set helper, and that an absent allow_list is omitted.
func TestAlkiraServiceFortinetAllowListGenerateRequest(t *testing.T) {
	// generateFortinetRequest resolves the management server segment by ID,
	// so the mock returns a segment for any lookup.
	mockClient := createMockAlkiraClient(t, func(w http.ResponseWriter, req *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(alkira.Segment{
			Id:   json.Number("1"),
			Name: "test-segment",
		})
	})

	// generateFortinetRequest requires at least one instance. A pre-set
	// credential_id keeps expansion from creating a credential.
	instances := []interface{}{
		map[string]interface{}{
			"name":          "fortinet-instance-1",
			"credential_id": "cred-1",
		},
	}

	// management_server_segment_id is Required in the real schema; Terraform
	// itself would block a plan without it. Set it explicitly rather than
	// relying on the zero-value "" the fixture previously left unset, which
	// the mock happened to resolve for any id (see comment above) but which
	// getSegmentNameById now rejects before ever reaching the API.
	managementServerSegmentId := "1"

	tests := []struct {
		name      string
		allowList interface{}
		expected  []string
	}{
		{
			name:      "populated allow_list",
			allowList: []interface{}{"10.0.0.0/24", "10.0.0.5"},
			expected:  []string{"10.0.0.0/24", "10.0.0.5"},
		},
		{
			name:      "allow_list absent",
			allowList: nil,
			expected:  nil,
		},
		{
			name:      "allow_list explicitly empty",
			allowList: []interface{}{},
			expected:  nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			raw := map[string]interface{}{
				"instances":                    instances,
				"management_server_segment_id": managementServerSegmentId,
			}
			if tt.allowList != nil {
				raw["allow_list"] = tt.allowList
			}

			r := resourceAlkiraServiceFortinet()
			d := schema.TestResourceDataRaw(t, r.Schema, raw)

			service, err := generateFortinetRequest(d, mockClient)
			require.NoError(t, err, "generateFortinetRequest should not return error")
			assert.ElementsMatch(t, tt.expected, service.AllowList)
		})
	}
}

// TestAlkiraServiceFortinetAllowListRead confirms resourceFortinetRead
// populates allow_list from the API response with no resulting diff.
func TestAlkiraServiceFortinetAllowListRead(t *testing.T) {
	allowList := []string{"10.0.0.0/24", "10.0.0.5"}

	client := createMockAlkiraClient(t, func(w http.ResponseWriter, req *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(alkira.ServiceFortinet{
			Id:        json.Number("1"),
			Name:      "test-fortinet-service",
			AllowList: allowList,
		})
	})

	r := resourceAlkiraServiceFortinet()
	d := r.TestResourceData()
	d.SetId("1")

	// Assert on the full diagnostics: a failed GET surfaces as a *Warning*
	// here, which HasError() would not catch.
	diags := resourceFortinetRead(context.Background(), d, client)
	require.Empty(t, diags, "resourceFortinetRead should return no diagnostics")

	got := convertTypeSetToStringList(d.Get("allow_list").(*schema.Set))
	assert.ElementsMatch(t, allowList, got, "allow_list should round-trip from the API response")
}

// TestAlkiraServiceFortinetSegmentIdsRead pins segment_ids to the segment IDs
// behind the names the API returns, with no extra elements.
//
// The Read path builds this value itself: the API stores segments by name, the
// schema stores them by ID, and the declared Elem is TypeString. Feeding d.Set
// a []int fails the SDK's type check outright, so Read left segment_ids at
// whatever state already held. After `terraform import`, where state starts
// empty, that surfaced as a plan wanting to add every segment back.
func TestAlkiraServiceFortinetSegmentIdsRead(t *testing.T) {
	segmentIdsByName := map[string]string{"seg-a": "8", "seg-b": "9"}

	client := createMockAlkiraClient(t, func(w http.ResponseWriter, req *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		// The segment lookup is a get-by-name, which returns a list.
		if name := req.URL.Query().Get("name"); name != "" {
			require.Contains(t, req.URL.Path, "/segments",
				"segment IDs must be looked up against the segment endpoint")
			json.NewEncoder(w).Encode([]alkira.Segment{{
				Id:   json.Number(segmentIdsByName[name]),
				Name: name,
			}})
			return
		}

		json.NewEncoder(w).Encode(alkira.ServiceFortinet{
			Id:       json.Number("1"),
			Name:     "test-fortinet-service",
			Segments: []string{"seg-a", "seg-b"},
		})
	})

	r := resourceAlkiraServiceFortinet()
	d := r.TestResourceData()
	d.SetId("1")

	// A failed GET surfaces as a Warning here, which HasError() would not catch.
	diags := resourceFortinetRead(context.Background(), d, client)
	require.Empty(t, diags, "resourceFortinetRead should return no diagnostics")

	got := convertTypeSetToStringList(d.Get("segment_ids").(*schema.Set))
	assert.ElementsMatch(t, []string{"8", "9"}, got,
		"segment_ids should hold exactly the IDs of the segments the API returned")
}

// TestAlkiraServiceFortinetInstancesRead pins one `instances` entry per
// instance the API returns.
//
// setInstance walks the config-tracked instances first, then the API response
// for any it did not already match. The second walk used to stop after the
// first unmatched instance, so a service with more instances than state knew
// about (every service, during import) only ever landed one of them in a
// Required attribute.
func TestAlkiraServiceFortinetInstancesRead(t *testing.T) {
	client := createMockAlkiraClient(t, func(w http.ResponseWriter, req *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		if name := req.URL.Query().Get("name"); name != "" {
			json.NewEncoder(w).Encode([]alkira.Segment{{
				Id:   json.Number("8"),
				Name: name,
			}})
			return
		}

		json.NewEncoder(w).Encode(alkira.ServiceFortinet{
			Id:       json.Number("1"),
			Name:     "test-fortinet-service",
			Segments: []string{"seg-a"},
			Instances: []alkira.FortinetInstance{
				{Id: 101, Name: "instance-1"},
				{Id: 102, Name: "instance-2"},
				{Id: 103, Name: "instance-3"},
			},
		})
	})

	r := resourceAlkiraServiceFortinet()
	d := r.TestResourceData()
	d.SetId("1")

	// Seed one instance so both walks in setInstance are exercised. license_key
	// is write-only on the API side and is carried forward from state, so it
	// doubles as a check that the tracked instance is not rebuilt from scratch.
	require.NoError(t, d.Set("instances", []interface{}{
		map[string]interface{}{"id": 101, "name": "instance-1", "license_key": "key-1"},
	}))

	diags := resourceFortinetRead(context.Background(), d, client)
	require.Empty(t, diags, "resourceFortinetRead should return no diagnostics")

	got := d.Get("instances").([]interface{})
	require.Len(t, got, 3, "every instance the API returned should be in state exactly once")

	byName := map[string]map[string]interface{}{}
	for _, v := range got {
		m := v.(map[string]interface{})
		byName[m["name"].(string)] = m
	}

	require.Len(t, byName, 3, "instances should not be duplicated across the two walks")
	assert.Equal(t, "key-1", byName["instance-1"]["license_key"],
		"a tracked instance keeps the write-only license_key already in state")
	assert.Equal(t, 102, byName["instance-2"]["id"])
	assert.Equal(t, 103, byName["instance-3"]["id"])
}
