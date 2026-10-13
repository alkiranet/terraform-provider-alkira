package alkira

import (
	"context"
	"encoding/json"
	"net/http"
	"strings"
	"testing"

	"github.com/alkiranet/alkira-client-go/alkira"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCheckpointRead(t *testing.T) {
	t.Skip("Test skipped: mock server response format doesn't match API client expectations")
	// expectedCxp := "US-WEST"
	// c := &alkira.ServiceCheckpoint{
	// 	Cxp: expectedCxp,
	// }
	// ac := createMockAlkiraClient(t, func(w http.ResponseWriter, req *http.Request) {
	// 	json.NewEncoder(w).Encode(c)
	// 	w.Header().Set("Content-Type", "application/json")
	// })

	// r := resourceAlkiraServiceCheckpoint()
	// d := r.TestResourceData()

	// err := resourceCheckpointRead(nil, d, ac)
	// require.Nil(t, err)

	// require.Equal(t, expectedCxp, getStringFromResourceData(d, "cxp"))
}

// UNUSED: Commented out to suppress linter warnings
// // TEST HELPER
// func serveCheckpoint(t *testing.T, c *alkira.ServiceCheckpoint) *alkira.AlkiraClient {
// 	return createMockAlkiraClient(t, func(w http.ResponseWriter, req *http.Request) {
// 		json.NewEncoder(w).Encode(c)
// 		w.Header().Set("Content-Type", "application/json")
// 	})
// }

// newCheckpointAllowListMockClient returns a client whose mock routes by URL
// path: segment lookups get a segment, everything else (e.g. instance
// credential creation) gets an ID payload.
func newCheckpointAllowListMockClient(t *testing.T, service *alkira.ServiceCheckpoint) *alkira.AlkiraClient {
	return createMockAlkiraClient(t, func(w http.ResponseWriter, req *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		segment := alkira.Segment{Id: json.Number("1"), Name: "test-segment"}

		switch {
		// Lookup by name returns a list; lookup by ID returns one object.
		case strings.Contains(req.URL.Path, "segments") && req.URL.Query().Get("name") != "":
			json.NewEncoder(w).Encode([]alkira.Segment{segment})
		case strings.Contains(req.URL.Path, "segments"):
			json.NewEncoder(w).Encode(segment)
		case service != nil && strings.Contains(req.URL.Path, "chkp-fw-services"):
			json.NewEncoder(w).Encode(service)
		default:
			w.Write([]byte(`{"id": "credential-123"}`))
		}
	})
}

func TestAlkiraServiceCheckpointAllowListSchema(t *testing.T) {
	resourceSchema := resourceAlkiraCheckpoint().Schema

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

func TestAlkiraServiceCheckpointAllowListValidator(t *testing.T) {
	validate := resourceAlkiraCheckpoint().Schema["allow_list"].Elem.(*schema.Schema).ValidateFunc

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

func TestAlkiraServiceCheckpointAllowListExpand(t *testing.T) {
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

// TestAlkiraServiceCheckpointAllowListGenerateRequest drives the real
// generateCheckpointRequest path to confirm allow_list lands on the request
// payload via the Set helper, and that an absent allow_list is omitted.
func TestAlkiraServiceCheckpointAllowListGenerateRequest(t *testing.T) {
	mockClient := newCheckpointAllowListMockClient(t, nil)

	// generateCheckpointRequest requires at least one instance and resolves
	// the segment by ID.
	instance := []interface{}{
		map[string]interface{}{
			"name":    "checkpoint-instance-1",
			"sic_key": "test-sic-key",
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
			raw := map[string]interface{}{
				"instance":   instance,
				"segment_id": "1",
			}
			if tt.allowList != nil {
				raw["allow_list"] = tt.allowList
			}

			r := resourceAlkiraCheckpoint()
			d := schema.TestResourceDataRaw(t, r.Schema, raw)

			service, err := generateCheckpointRequest(d, mockClient)
			require.NoError(t, err, "generateCheckpointRequest should not return error")
			assert.ElementsMatch(t, tt.expected, service.AllowList)
		})
	}
}

// TestAlkiraServiceCheckpointAllowListRead confirms resourceCheckpointRead
// populates allow_list from the API response with no resulting diff.
func TestAlkiraServiceCheckpointAllowListRead(t *testing.T) {
	allowList := []string{"10.0.0.0/24", "10.0.0.5"}

	// resourceCheckpointRead requires exactly one segment on the response.
	client := newCheckpointAllowListMockClient(t, &alkira.ServiceCheckpoint{
		Id:        json.Number("1"),
		Name:      "test-checkpoint-service",
		AllowList: allowList,
		Segments:  []string{"test-segment"},
	})

	r := resourceAlkiraCheckpoint()
	d := r.TestResourceData()
	d.SetId("1")

	// Assert on the full diagnostics: a failed GET surfaces as a *Warning*
	// here, which HasError() would not catch.
	diags := resourceCheckpointRead(context.Background(), d, client)
	require.Empty(t, diags, "resourceCheckpointRead should return no diagnostics")

	got := convertTypeSetToStringList(d.Get("allow_list").(*schema.Set))
	assert.ElementsMatch(t, allowList, got, "allow_list should round-trip from the API response")
}
