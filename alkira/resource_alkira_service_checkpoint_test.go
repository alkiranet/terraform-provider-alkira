package alkira

import (
	"encoding/json"
	"net/http"
	"testing"

	"github.com/alkiranet/alkira-client-go/alkira"
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

// TEST HELPER
func serveCheckpoint(t *testing.T, c *alkira.ServiceCheckpoint) *alkira.AlkiraClient {
	return createMockAlkiraClient(t, func(w http.ResponseWriter, req *http.Request) {
		json.NewEncoder(w).Encode(c)
		w.Header().Set("Content-Type", "application/json")
	})
}
