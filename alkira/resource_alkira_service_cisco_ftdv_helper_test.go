package alkira

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/require"
)

// TestExpandCiscoFtdvManagementServerPropagatesSegmentIdError pins that a
// segment name in the management server's segment_id fails the expand
// instead of being silently dropped, and that the failure happens before
// any API call.
func TestExpandCiscoFtdvManagementServerPropagatesSegmentIdError(t *testing.T) {
	client := createMockAlkiraClient(t, http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		t.Errorf("segment API must not be called for a rejected segment_id")
	}))

	in := []interface{}{
		map[string]interface{}{
			"server_ip":     "10.0.0.1",
			"username":      "admin",
			"password":      "secret",
			"credential_id": "existing-cred-id",
			"segment_id":    "ak74335-seg-a",
			"ip_allow_list": []interface{}{},
		},
	}

	_, _, _, err := expandCiscoFtdvManagementServer(in, client)

	require.Error(t, err)
}
