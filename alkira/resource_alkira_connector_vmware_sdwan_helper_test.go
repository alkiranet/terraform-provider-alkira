package alkira

import (
	"encoding/json"
	"io"
	"net/http"
	"testing"

	"github.com/alkiranet/alkira-client-go/alkira"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type vmwareSdwanCredentialCall struct {
	Method string
	Path   string
	Body   map[string]interface{}
}

// createVmwareSdwanCredentialMockClient returns a mock client that serves the
// credentials API and records every request it receives
func createVmwareSdwanCredentialMockClient(t *testing.T, calls *[]vmwareSdwanCredentialCall) *alkira.AlkiraClient {
	return createMockAlkiraClient(t, http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		call := vmwareSdwanCredentialCall{Method: req.Method, Path: req.URL.Path}

		if data, _ := io.ReadAll(req.Body); len(data) > 0 {
			json.Unmarshal(data, &call.Body)
		}
		*calls = append(*calls, call)

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		switch {
		case req.Method == http.MethodGet && req.URL.Path == "/api/credentials/":
			json.NewEncoder(w).Encode([]alkira.CredentialResponseDetail{
				{Id: "cred-1", Type: "VMWARE_SDWAN_CONNECTOR_INSTANCE", Name: "edge1abcdef"},
			})
		case req.Method == http.MethodPost:
			json.NewEncoder(w).Encode(alkira.CredentialResponse{Id: "cred-new"})
		}
	}))
}

func TestVmwareSdwanExpandVirtualEdgesCreatesCredentialForNewEdge(t *testing.T) {
	var calls []vmwareSdwanCredentialCall
	client := createVmwareSdwanCredentialMockClient(t, &calls)

	in := []interface{}{
		map[string]interface{}{
			"hostname":        "edge1",
			"activation_code": "AAAA-AAAA-AAAA-AAAA",
			"credential_id":   "",
			"id":              0,
		},
	}

	result, err := expandVmwareSdwanVirtualEdges(client, in, nil)

	require.NoError(t, err)
	require.Len(t, result, 1)
	assert.Equal(t, "cred-new", result[0].CredentialId)
	require.Len(t, calls, 1)
	assert.Equal(t, http.MethodPost, calls[0].Method)
	assert.Equal(t, "/api/credentials/vmware-sdwan-connector-instance", calls[0].Path)
}

func TestVmwareSdwanExpandVirtualEdgesUpdatesCredentialOnActivationCodeChange(t *testing.T) {
	var calls []vmwareSdwanCredentialCall
	client := createVmwareSdwanCredentialMockClient(t, &calls)

	old := []interface{}{
		map[string]interface{}{
			"hostname":        "edge1",
			"activation_code": "AAAA-AAAA-AAAA-AAAA",
			"credential_id":   "cred-1",
			"id":              2134,
		},
	}
	in := []interface{}{
		map[string]interface{}{
			"hostname":        "edge2",
			"activation_code": "BBBB-BBBB-BBBB-BBBB",
			"credential_id":   "cred-1",
			"id":              2134,
		},
	}

	result, err := expandVmwareSdwanVirtualEdges(client, in, old)

	require.NoError(t, err)
	require.Len(t, result, 1)
	assert.Equal(t, "cred-1", result[0].CredentialId)
	assert.Equal(t, "edge2", result[0].HostName)
	assert.Equal(t, 2134, result[0].Id)

	// The existing credential is looked up for its name, then updated in place
	require.Len(t, calls, 2)
	assert.Equal(t, http.MethodGet, calls[0].Method)
	assert.Equal(t, "/api/credentials/vmware-sdwan-connector-instance/cred-1", calls[1].Path)
	assert.NotEqual(t, http.MethodPost, calls[1].Method)
	assert.Equal(t, "edge1abcdef", calls[1].Body["name"])

	credentials, ok := calls[1].Body["credentials"].(map[string]interface{})
	require.True(t, ok)
	assert.Equal(t, "BBBB-BBBB-BBBB-BBBB", credentials["activationCode"])
}

func TestVmwareSdwanExpandVirtualEdgesSkipsCredentialWhenActivationCodeUnchanged(t *testing.T) {
	var calls []vmwareSdwanCredentialCall
	client := createVmwareSdwanCredentialMockClient(t, &calls)

	old := []interface{}{
		map[string]interface{}{
			"hostname":        "edge1",
			"activation_code": "AAAA-AAAA-AAAA-AAAA",
			"credential_id":   "cred-1",
			"id":              2134,
		},
	}
	in := []interface{}{
		map[string]interface{}{
			"hostname":        "edge2",
			"activation_code": "AAAA-AAAA-AAAA-AAAA",
			"credential_id":   "cred-1",
			"id":              2134,
		},
	}

	result, err := expandVmwareSdwanVirtualEdges(client, in, old)

	require.NoError(t, err)
	require.Len(t, result, 1)
	assert.Equal(t, "cred-1", result[0].CredentialId)
	assert.Empty(t, calls, "an unchanged activation code must not touch the credential")
}

func TestVmwareSdwanExpandVirtualEdgesFailsWhenCredentialMissing(t *testing.T) {
	var calls []vmwareSdwanCredentialCall
	client := createVmwareSdwanCredentialMockClient(t, &calls)

	old := []interface{}{
		map[string]interface{}{
			"hostname":        "edge1",
			"activation_code": "AAAA-AAAA-AAAA-AAAA",
			"credential_id":   "cred-gone",
			"id":              2134,
		},
	}
	in := []interface{}{
		map[string]interface{}{
			"hostname":        "edge1",
			"activation_code": "BBBB-BBBB-BBBB-BBBB",
			"credential_id":   "cred-gone",
			"id":              2134,
		},
	}

	result, err := expandVmwareSdwanVirtualEdges(client, in, old)

	assert.Error(t, err)
	assert.Nil(t, result)
}
