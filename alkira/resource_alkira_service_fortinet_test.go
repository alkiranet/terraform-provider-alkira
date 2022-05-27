package alkira

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/alkiranet/alkira-client-go/alkira"
	"github.com/stretchr/testify/require"
)

func TestFortinetRead(t *testing.T) {
	expectedCxp := "US-WEST"
	expectedIp := "10.1.1.0"
	expectedSegment := "default"

	f := &alkira.Fortinet{
		Cxp: expectedCxp,
		ManagementServer: &alkira.FortinetManagmentServer{
			IpAddress: expectedIp,
			Segment:   expectedSegment,
		},
	}
	ac := serveFortinet(t, f)

	r := resourceAlkiraServiceFortinet()
	d := r.TestResourceData()

	err := resourceFortinetRead(d, ac)
	require.Nil(t, err)

	require.Equal(t, expectedCxp, d.Get("cxp"))
	require.Equal(t, expectedIp, d.Get("management_server_ip"))
	require.Equal(t, expectedSegment, d.Get("management_server_segment"))
}

func TestFortinetReadAutoScale(t *testing.T) {
	expectedAutoScaleVal := "ON"

	f := &alkira.Fortinet{
		AutoScale: expectedAutoScaleVal,
	}
	ac := serveFortinet(t, f)
	f.ManagementServer = &alkira.FortinetManagmentServer{}

	r := resourceAlkiraServiceFortinet()
	d := r.TestResourceData()

	err := resourceFortinetRead(d, ac)
	require.Nil(t, err)

	require.Equal(t, expectedAutoScaleVal, d.Get("auto_scale"))
}

//
// TEST HELPER
//

func serveFortinet(t *testing.T, f *alkira.Fortinet) *alkira.AlkiraClient {
	server := httptest.NewServer(http.HandlerFunc(
		func(w http.ResponseWriter, req *http.Request) {
			json.NewEncoder(w).Encode(f)
			w.Header().Set("Content-Type", "application/json")
		},
	))
	t.Cleanup(server.Close)

	return &alkira.AlkiraClient{
		URI:             server.URL,
		TenantNetworkId: "0",
		Client:          &http.Client{Timeout: time.Duration(1) * time.Second},
	}
}
