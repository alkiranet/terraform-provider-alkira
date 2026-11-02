package alkira

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"sync"
	"testing"

	"github.com/alkiranet/alkira-client-go/alkira"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type merakiSDWANCredentialCall struct {
	Method string
	Path   string
	Body   map[string]interface{}
}

// merakiSDWANCredentialServer records credential calls. A POST whose
// authToken equals failToken returns 500. GET /api/credentials/ lists
// each existing credential as "<id>-name".
func merakiSDWANCredentialServer(t *testing.T, failToken string) (*alkira.AlkiraClient, *[]merakiSDWANCredentialCall) {
	var mu sync.Mutex
	calls := []merakiSDWANCredentialCall{}
	next := 0

	client := createMockAlkiraClient(t, func(w http.ResponseWriter, req *http.Request) {
		mu.Lock()
		defer mu.Unlock()

		call := merakiSDWANCredentialCall{Method: req.Method, Path: req.URL.Path}
		if data, _ := io.ReadAll(req.Body); len(data) > 0 {
			assert.NoError(t, json.Unmarshal(data, &call.Body))
		}
		calls = append(calls, call)

		w.Header().Set("Content-Type", "application/json")
		switch {
		case req.Method == "GET" && req.URL.Path == "/api/credentials/":
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode([]map[string]string{
				{"credentialId": "cred-a", "name": "cred-a-name"},
				{"credentialId": "cred-b", "name": "cred-b-name"},
				{"credentialId": "cred-d", "name": "cred-d-name"},
			})
		case req.Method == "POST":
			creds := call.Body["credentials"].(map[string]interface{})
			if creds["authToken"] == failToken {
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
			next++
			w.WriteHeader(http.StatusCreated)
			json.NewEncoder(w).Encode(map[string]string{"id": fmt.Sprintf("cred-new-%d", next)})
		default:
			w.WriteHeader(http.StatusOK)
		}
	})
	client.Client.RetryMax = 0

	return client, &calls
}

func merakiSDWANCredentialWrites(calls []merakiSDWANCredentialCall) []merakiSDWANCredentialCall {
	var writes []merakiSDWANCredentialCall
	for _, c := range calls {
		if c.Method != "GET" {
			writes = append(writes, c)
		}
	}
	return writes
}

func TestExpandMerakiSDWANSegmentOptions(t *testing.T) {
	var maxAsn int64 = 4294967295

	in := []interface{}{
		map[string]interface{}{
			"segment_id":               7,
			"meraki_bgp_asn":           int(maxAsn),
			"advertise_on_prem_routes": true,
			"advertise_default_route":  true,
		},
	}

	result := expandMerakiSDWANSegmentOptions(in)

	require.Len(t, result, 1)
	assert.Equal(t, 7, result[0].SegmentId)
	assert.Equal(t, int64(int(maxAsn)), result[0].MerakiBgpAsn)
	assert.True(t, result[0].AdvertiseOnPremRoutes)
	assert.False(t, result[0].DisableInternetExit)

	in[0].(map[string]interface{})["advertise_default_route"] = false
	assert.True(t, expandMerakiSDWANSegmentOptions(in)[0].DisableInternetExit)
}

func TestMerakiSDWANReplacedTokensError(t *testing.T) {
	cause := assert.AnError

	assert.Equal(t, cause, merakiSDWANReplacedTokensError(nil, cause), "nothing replaced")

	err := merakiSDWANReplacedTokensError([]string{"vmx-a (credential cred-a)"}, cause)
	assert.ErrorIs(t, err, cause)
	assert.Contains(t, err.Error(), "vmx-a (credential cred-a)")
	assert.Contains(t, err.Error(), "apply again")
}

func TestExpandMerakiSDWANInstances(t *testing.T) {
	client, calls := merakiSDWANCredentialServer(t, "")

	old := []interface{}{
		map[string]interface{}{"id": 11, "host_name": "vmx-a", "credential_id": "cred-a", "auth_token": "ta"},
		map[string]interface{}{"id": 12, "host_name": "vmx-b", "credential_id": "cred-b", "auth_token": "tb"},
		// Imported: the token is not known
		map[string]interface{}{"id": 14, "host_name": "vmx-d", "credential_id": "cred-d", "auth_token": ""},
	}
	// The IDs and credentials come from the prior state of the same
	// host_name, not from the planned values
	in := []interface{}{
		map[string]interface{}{"host_name": "vmx-a", "auth_token": "ta"},
		map[string]interface{}{"host_name": "vmx-b", "auth_token": "tb2"},
		map[string]interface{}{"host_name": "vmx-d", "auth_token": "td"},
		map[string]interface{}{"id": 0, "host_name": "vmx-c", "credential_id": "", "auth_token": "tc"},
	}

	result, created, replaced, err := expandMerakiSDWANInstances(client, in, old)
	require.NoError(t, err)
	require.Len(t, result, 4)
	assert.Equal(t, []string{"cred-new-1"}, created)
	assert.Equal(t, []string{"vmx-b (credential cred-b)", "vmx-d (credential cred-d)"}, replaced)

	assert.Equal(t, alkira.ConnectorMerakiSDWANInstance{Id: "11", HostName: "vmx-a", CredentialId: "cred-a"}, result[0])
	assert.Equal(t, alkira.ConnectorMerakiSDWANInstance{Id: "12", HostName: "vmx-b", CredentialId: "cred-b"}, result[1])
	assert.Equal(t, alkira.ConnectorMerakiSDWANInstance{Id: "14", HostName: "vmx-d", CredentialId: "cred-d"}, result[2])
	assert.Equal(t, alkira.ConnectorMerakiSDWANInstance{HostName: "vmx-c", CredentialId: "cred-new-1"}, result[3])

	writes := merakiSDWANCredentialWrites(*calls)
	require.Len(t, writes, 3)

	assert.Equal(t, "PUT", writes[0].Method)
	assert.Equal(t, "/api/credentials/meraki-sdwan-connector-instance/cred-b", writes[0].Path)
	assert.Equal(t, "tb2", writes[0].Body["credentials"].(map[string]interface{})["authToken"])
	assert.Equal(t, "cred-b-name", writes[0].Body["name"], "the credential keeps its name")

	assert.Equal(t, "/api/credentials/meraki-sdwan-connector-instance/cred-d", writes[1].Path)
	assert.Equal(t, "td", writes[1].Body["credentials"].(map[string]interface{})["authToken"])

	assert.Equal(t, "POST", writes[2].Method)
	assert.Equal(t, "/api/credentials/meraki-sdwan-connector-instance", writes[2].Path)
	assert.Equal(t, "tc", writes[2].Body["credentials"].(map[string]interface{})["authToken"])
	assert.True(t, strings.HasPrefix(writes[2].Body["name"].(string), "vmx-c"))

	payload, err := json.Marshal(result[3])
	require.NoError(t, err)
	assert.NotContains(t, string(payload), `"id"`)
}

func TestExpandMerakiSDWANInstancesRejected(t *testing.T) {
	old := []interface{}{
		map[string]interface{}{"id": 11, "host_name": "vmx-a", "credential_id": "cred-a", "auth_token": "ta"},
		map[string]interface{}{"id": 12, "host_name": "vmx-b", "credential_id": "cred-b", "auth_token": "tb"},
	}

	t.Run("a renamed instance keeps the old token", func(t *testing.T) {
		client, calls := merakiSDWANCredentialServer(t, "")

		in := []interface{}{
			// A valid change before the failing one must not write a credential
			map[string]interface{}{"host_name": "vmx-a", "auth_token": "ta2"},
			map[string]interface{}{"host_name": "vmx-x", "auth_token": "tb"},
		}

		_, created, _, err := expandMerakiSDWANInstances(client, in, old)
		require.Error(t, err)
		assert.Contains(t, err.Error(), `"vmx-x"`)
		assert.Contains(t, err.Error(), `"vmx-b"`)
		assert.Empty(t, created)
		assert.Empty(t, merakiSDWANCredentialWrites(*calls), "no credential call before the plan is valid")
	})

	t.Run("two instances have the same token", func(t *testing.T) {
		client, calls := merakiSDWANCredentialServer(t, "")

		in := []interface{}{
			map[string]interface{}{"host_name": "vmx-a", "auth_token": "ta"},
			map[string]interface{}{"host_name": "vmx-c", "auth_token": "tc"},
			map[string]interface{}{"host_name": "vmx-d", "auth_token": "tc"},
		}

		_, _, _, err := expandMerakiSDWANInstances(client, in, old)
		require.Error(t, err)
		assert.Contains(t, err.Error(), `"vmx-c"`)
		assert.Contains(t, err.Error(), `"vmx-d"`)
		assert.Empty(t, merakiSDWANCredentialWrites(*calls))
	})

	t.Run("a host_name is used twice", func(t *testing.T) {
		client, calls := merakiSDWANCredentialServer(t, "")

		in := []interface{}{
			map[string]interface{}{"host_name": "vmx-a", "auth_token": "ta2"},
			map[string]interface{}{"host_name": "vmx-a", "auth_token": "ta3"},
		}

		_, _, _, err := expandMerakiSDWANInstances(client, in, old)
		require.Error(t, err)
		assert.Contains(t, err.Error(), `host_name "vmx-a"`)
		assert.Empty(t, merakiSDWANCredentialWrites(*calls))
	})

	t.Run("a renamed instance with a new token is a new instance", func(t *testing.T) {
		client, calls := merakiSDWANCredentialServer(t, "")

		in := []interface{}{
			map[string]interface{}{"host_name": "vmx-a", "auth_token": "ta"},
			map[string]interface{}{"host_name": "vmx-x", "auth_token": "tx"},
		}

		result, created, replaced, err := expandMerakiSDWANInstances(client, in, old)
		require.NoError(t, err)
		assert.Equal(t, alkira.ConnectorMerakiSDWANInstance{HostName: "vmx-x", CredentialId: "cred-new-1"}, result[1])
		assert.Equal(t, []string{"cred-new-1"}, created)
		assert.Empty(t, replaced)

		writes := merakiSDWANCredentialWrites(*calls)
		require.Len(t, writes, 1)
		assert.Equal(t, "POST", writes[0].Method)
	})
}

func TestExpandMerakiSDWANInstancesCredentialFailure(t *testing.T) {
	client, calls := merakiSDWANCredentialServer(t, "bad")

	in := []interface{}{
		map[string]interface{}{"host_name": "vmx-a", "credential_id": "", "auth_token": "good"},
		map[string]interface{}{"host_name": "vmx-b", "credential_id": "", "auth_token": "bad"},
	}

	_, _, _, err := expandMerakiSDWANInstances(client, in, nil)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "vmx-b")

	var deletes []string
	for _, c := range *calls {
		if c.Method == "DELETE" {
			deletes = append(deletes, c.Path)
		}
	}
	assert.Equal(t, []string{"/api/credentials/meraki-sdwan-connector-instance/cred-new-1"}, deletes)
}

func TestSetMerakiSDWANInstances(t *testing.T) {
	d := resourceAlkiraConnectorMerakiSDWAN().TestResourceData()
	require.NoError(t, d.Set("instances", []interface{}{
		map[string]interface{}{"host_name": "vmx-b", "id": 2, "auth_token": "tb"},
		// Not tracked yet, as on the read at the end of a create
		map[string]interface{}{"host_name": "vmx-a", "auth_token": "ta"},
	}))

	connector := &alkira.ConnectorMerakiSDWAN{
		Instances: []alkira.ConnectorMerakiSDWANInstance{
			{Id: "1", HostName: "vmx-a", CredentialId: "cred-a", Name: "inst-1", Az: 2, WanIp: "10.0.0.4", PublicIp: "20.1.1.1"},
			{Id: "2", HostName: "vmx-b", CredentialId: "cred-b", Name: "inst-2", Az: 1},
			{Id: "3", HostName: "vmx-c", CredentialId: "cred-c"},
		},
	}

	setMerakiSDWANInstances(d, connector)

	instances := merakiSDWANInstancesByHostName(d)
	require.Len(t, instances, 3)

	b := instances["vmx-b"]
	assert.Equal(t, 2, b["id"])
	assert.Equal(t, "cred-b", b["credential_id"])
	assert.Equal(t, "tb", b["auth_token"])
	assert.Equal(t, 1, b["az"])

	a := instances["vmx-a"]
	assert.Equal(t, 1, a["id"])
	assert.Equal(t, "inst-1", a["name"])
	assert.Equal(t, "10.0.0.4", a["wan_ip"])
	assert.Equal(t, "20.1.1.1", a["public_ip"])
	assert.Equal(t, "ta", a["auth_token"])

	assert.Equal(t, "", instances["vmx-c"]["auth_token"], "an instance not in the config has no token, as after an import")
}

func TestSetMerakiSDWANInstancesMatchesByHostName(t *testing.T) {
	d := resourceAlkiraConnectorMerakiSDWAN().TestResourceData()
	require.NoError(t, d.Set("instances", []interface{}{
		// Another ID than the API's, as in a planned value after a change
		map[string]interface{}{"host_name": "vmx-b", "id": 1, "auth_token": "tb"},
	}))

	connector := &alkira.ConnectorMerakiSDWAN{
		Instances: []alkira.ConnectorMerakiSDWANInstance{
			{Id: "2", HostName: "vmx-b", CredentialId: "cred-b"},
		},
	}

	setMerakiSDWANInstances(d, connector)

	instances := merakiSDWANInstancesByHostName(d)
	require.Len(t, instances, 1)
	assert.Equal(t, 2, instances["vmx-b"]["id"])
	assert.Equal(t, "tb", instances["vmx-b"]["auth_token"], "the token is kept by host_name")
}
