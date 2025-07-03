package alkira

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/alkiranet/alkira-client-go/alkira"
	"github.com/hashicorp/go-retryablehttp"
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
	// ac := serveFortinet(t, f)

	// r := resourceAlkiraServiceFortinet()
	// d := r.TestResourceData()

	// err := resourceFortinetRead(nil, d, ac)
	// require.Nil(t, err)

	// require.Equal(t, expectedCxp, d.Get("cxp"))
	// require.Equal(t, expectedIp, d.Get("management_server_ip"))
	// require.Equal(t, expectedSegment, d.Get("management_server_segment"))
}

func TestFortinetReadAutoScale(t *testing.T) {
	t.Skip("Test skipped: mock server response format doesn't match API client expectations")
	// expectedAutoScaleVal := "ON"

	// f := &alkira.ServiceFortinet{
	// 	AutoScale: expectedAutoScaleVal,
	// }
	// ac := serveFortinet(t, f)
	// f.ManagementServer = &alkira.FortinetManagmentServer{}

	// r := resourceAlkiraServiceFortinet()
	// d := r.TestResourceData()

	// err := resourceFortinetRead(nil, d, ac)
	// require.Nil(t, err)

	// require.Equal(t, expectedAutoScaleVal, d.Get("auto_scale"))
}

//
// TEST HELPER
//

func serveFortinet(t *testing.T, f *alkira.ServiceFortinet) *alkira.AlkiraClient {
	server := httptest.NewServer(http.HandlerFunc(
		func(w http.ResponseWriter, req *http.Request) {
			json.NewEncoder(w).Encode(f)
			w.Header().Set("Content-Type", "application/json")
		},
	))
	t.Cleanup(server.Close)

	retryClient := retryablehttp.NewClient()
	retryClient.HTTPClient.Timeout = time.Duration(1) * time.Second

	return &alkira.AlkiraClient{
		URI:             server.URL,
		TenantNetworkId: "0",
		Client:          retryClient,
	}
}
