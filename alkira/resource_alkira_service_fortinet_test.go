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
			raw := map[string]interface{}{"instances": instances}
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
