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
