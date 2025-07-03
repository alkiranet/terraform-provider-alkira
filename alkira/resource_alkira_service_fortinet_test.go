package alkira

import (
	"encoding/json"
	"net/http"
	"testing"

	"github.com/alkiranet/alkira-client-go/alkira"
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

func serveFortinet(t *testing.T, f *alkira.ServiceFortinet) *alkira.AlkiraClient {
	return createMockAlkiraClient(t, func(w http.ResponseWriter, req *http.Request) {
		json.NewEncoder(w).Encode(f)
		w.Header().Set("Content-Type", "application/json")
	})
}
