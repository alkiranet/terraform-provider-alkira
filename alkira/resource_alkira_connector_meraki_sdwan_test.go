package alkira

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"sync"
	"testing"

	"github.com/alkiranet/alkira-client-go/alkira"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestResourceAlkiraConnectorMerakiSDWANSchema(t *testing.T) {
	require.NoError(t, resourceAlkiraConnectorMerakiSDWAN().InternalValidate(nil, true))
	require.NoError(t, dataSourceAlkiraConnectorMerakiSDWAN().InternalValidate(nil, false))
}

func TestValidateMerakiSDWANBgpAsn(t *testing.T) {
	var aboveMax int64 = 4294967296

	_, errs := validateMerakiSDWANBgpAsn(0, "meraki_bgp_asn")
	assert.NotEmpty(t, errs)

	_, errs = validateMerakiSDWANBgpAsn(65500, "meraki_bgp_asn")
	assert.Empty(t, errs)

	if int64(int(aboveMax)) == aboveMax {
		_, errs = validateMerakiSDWANBgpAsn(int(aboveMax), "meraki_bgp_asn")
		assert.NotEmpty(t, errs)
		_, errs = validateMerakiSDWANBgpAsn(int(aboveMax-1), "meraki_bgp_asn")
		assert.Empty(t, errs)
	}
}

// merakiSDWANKey is the state key of a field of the instance with
// hostName and authToken
func merakiSDWANKey(hostName, authToken, field string) string {
	return fmt.Sprintf("instances.%d.%s", merakiSDWANInstanceHash(map[string]interface{}{
		"host_name": hostName, "auth_token": authToken}), field)
}

// merakiSDWANChanges is the keys under prefix that the diff changes
func merakiSDWANChanges(diff *terraform.InstanceDiff, prefix string) []string {
	var keys []string
	for key, a := range diff.Attributes {
		if strings.HasPrefix(key, prefix) && (a.Old != a.New || a.NewComputed || a.NewRemoved) {
			keys = append(keys, key)
		}
	}
	return keys
}

// merakiSDWANState is the state of a connector with two instances. The
// tokens are empty after an import.
func merakiSDWANState(tokenA, tokenB string) *terraform.InstanceState {
	return &terraform.InstanceState{
		ID: "10",
		Attributes: map[string]string{
			"id": "10", "name": "c", "cxp": "AZURE-WESTUS2", "size": "SMALL", "version": "18.2.14",
			"tunnel_protocol": "IPSEC", "enabled": "true", "description": "", "group": "",
			"target_segment.#":                               "1",
			"target_segment.0.segment_id":                    "1",
			"target_segment.0.meraki_bgp_asn":                "65500",
			"target_segment.0.advertise_on_prem_routes":      "false",
			"target_segment.0.advertise_default_route":       "false",
			"instances.#":                                    "2",
			merakiSDWANKey("vmx-a", tokenA, "host_name"):     "vmx-a",
			merakiSDWANKey("vmx-a", tokenA, "credential_id"): "cred-a",
			merakiSDWANKey("vmx-a", tokenA, "id"):            "1",
			merakiSDWANKey("vmx-a", tokenA, "auth_token"):    tokenA,
			merakiSDWANKey("vmx-b", tokenB, "host_name"):     "vmx-b",
			merakiSDWANKey("vmx-b", tokenB, "credential_id"): "cred-b",
			merakiSDWANKey("vmx-b", tokenB, "id"):            "2",
			merakiSDWANKey("vmx-b", tokenB, "auth_token"):    tokenB,
			merakiSDWANKey("vmx-a", tokenA, "name"):          "inst-1",
			merakiSDWANKey("vmx-a", tokenA, "az"):            "1",
			merakiSDWANKey("vmx-a", tokenA, "wan_ip"):        "10.0.0.4",
			merakiSDWANKey("vmx-a", tokenA, "public_ip"):     "20.1.1.1",
			merakiSDWANKey("vmx-b", tokenB, "name"):          "inst-2",
			merakiSDWANKey("vmx-b", tokenB, "az"):            "2",
			merakiSDWANKey("vmx-b", tokenB, "wan_ip"):        "10.0.0.5",
			merakiSDWANKey("vmx-b", tokenB, "public_ip"):     "20.1.1.2",
		},
	}
}

func merakiSDWANConfig(instances ...map[string]interface{}) *terraform.ResourceConfig {
	list := []interface{}{}
	for _, i := range instances {
		list = append(list, i)
	}

	return terraform.NewResourceConfigRaw(map[string]interface{}{
		"name":           "c",
		"cxp":            "AZURE-WESTUS2",
		"size":           "SMALL",
		"target_segment": []interface{}{map[string]interface{}{"segment_id": 1, "meraki_bgp_asn": 65500}},
		"instances":      list,
	})
}

func TestResourceAlkiraConnectorMerakiSDWANPlan(t *testing.T) {
	resource := resourceAlkiraConnectorMerakiSDWAN()

	t.Run("no diff when the tokens are unchanged", func(t *testing.T) {
		diff, err := resource.Diff(context.Background(), merakiSDWANState("ta", "tb"), merakiSDWANConfig(
			map[string]interface{}{"host_name": "vmx-a", "auth_token": "ta"},
			map[string]interface{}{"host_name": "vmx-b", "auth_token": "tb"},
		), &alkira.AlkiraClient{})

		require.NoError(t, err)
		assert.Nil(t, diff)
	})

	t.Run("a changed token is planned", func(t *testing.T) {
		diff, err := resource.Diff(context.Background(), merakiSDWANState("ta", "tb"), merakiSDWANConfig(
			map[string]interface{}{"host_name": "vmx-a", "auth_token": "ta2"},
			map[string]interface{}{"host_name": "vmx-b", "auth_token": "tb"},
		), &alkira.AlkiraClient{})

		require.NoError(t, err)
		require.NotNil(t, diff)
		assert.True(t, diff.Attributes[merakiSDWANKey("vmx-a", "ta", "host_name")].NewRemoved)
		assert.Equal(t, "ta2", diff.Attributes[merakiSDWANKey("vmx-a", "ta2", "auth_token")].New)
		assert.Equal(t, "vmx-a", diff.Attributes[merakiSDWANKey("vmx-a", "ta2", "host_name")].New)
		assert.Empty(t, merakiSDWANChanges(diff, merakiSDWANKey("vmx-b", "tb", "")))
	})

	t.Run("the token is planned after an import", func(t *testing.T) {
		diff, err := resource.Diff(context.Background(), merakiSDWANState("", ""), merakiSDWANConfig(
			map[string]interface{}{"host_name": "vmx-a", "auth_token": "ta"},
			map[string]interface{}{"host_name": "vmx-b", "auth_token": "tb"},
		), &alkira.AlkiraClient{})

		require.NoError(t, err)
		require.NotNil(t, diff)
		assert.Equal(t, "ta", diff.Attributes[merakiSDWANKey("vmx-a", "ta", "auth_token")].New)
		assert.Equal(t, "tb", diff.Attributes[merakiSDWANKey("vmx-b", "tb", "auth_token")].New)
	})

	t.Run("no diff when the order changes", func(t *testing.T) {
		diff, err := resource.Diff(context.Background(), merakiSDWANState("ta", "tb"), merakiSDWANConfig(
			map[string]interface{}{"host_name": "vmx-b", "auth_token": "tb"},
			map[string]interface{}{"host_name": "vmx-a", "auth_token": "ta"},
		), &alkira.AlkiraClient{})

		require.NoError(t, err)
		assert.Nil(t, diff)
	})

	t.Run("removing the first instance changes only that instance", func(t *testing.T) {
		diff, err := resource.Diff(context.Background(), merakiSDWANState("ta", "tb"), merakiSDWANConfig(
			map[string]interface{}{"host_name": "vmx-b", "auth_token": "tb"},
		), &alkira.AlkiraClient{})

		require.NoError(t, err)
		require.NotNil(t, diff)
		assert.Equal(t, "1", diff.Attributes["instances.#"].New)
		assert.True(t, diff.Attributes[merakiSDWANKey("vmx-a", "ta", "host_name")].NewRemoved)
		assert.Empty(t, merakiSDWANChanges(diff, merakiSDWANKey("vmx-b", "tb", "")), "vmx-b must not change")
	})

	t.Run("a new instance gets a new ID and credential", func(t *testing.T) {
		diff, err := resource.Diff(context.Background(), merakiSDWANState("ta", "tb"), merakiSDWANConfig(
			map[string]interface{}{"host_name": "vmx-a", "auth_token": "ta"},
			map[string]interface{}{"host_name": "vmx-b", "auth_token": "tb"},
			map[string]interface{}{"host_name": "vmx-c", "auth_token": "tc"},
		), &alkira.AlkiraClient{})

		require.NoError(t, err)
		require.NotNil(t, diff)
		assert.Equal(t, "3", diff.Attributes["instances.#"].New)
		assert.Equal(t, "vmx-c", diff.Attributes[merakiSDWANKey("vmx-c", "tc", "host_name")].New)
		assert.True(t, diff.Attributes[merakiSDWANKey("vmx-c", "tc", "id")].NewComputed)
		assert.True(t, diff.Attributes[merakiSDWANKey("vmx-c", "tc", "credential_id")].NewComputed)
		assert.Empty(t, merakiSDWANChanges(diff, merakiSDWANKey("vmx-a", "ta", "")), "vmx-a must not change")
	})
}

type merakiSDWANCall struct {
	Method string
	Path   string
	Query  string
	Body   map[string]interface{}
}

var merakiSDWANConnectorResponse = map[string]interface{}{
	"id": 10, "name": "c", "cxp": "AZURE-WESTUS2", "size": "SMALL", "version": "18.2.14",
	"tunnelProtocol": "IPSEC", "enabled": true, "implicitGroupId": 77,
	"instances": []map[string]interface{}{
		{"id": 1, "hostName": "vmx-a", "credentialId": "cred-a"},
		{"id": 2, "hostName": "vmx-b", "credentialId": "cred-b"},
	},
	"segmentOptions": []map[string]interface{}{{"segmentId": 1, "merakiBgpAsn": 65500, "disableInternetExit": true}},
}

// merakiSDWANServer serves connector 10 and its credentials and records
// every call. The connector POST, PUT and DELETE reply with
// connectorStatus; 202 starts a validation that fails, 200 a provision
// request that succeeds. A credential POST returns cred-new-N.
func merakiSDWANServer(t *testing.T, connectorStatus int) (*alkira.AlkiraClient, *[]merakiSDWANCall) {
	var mu sync.Mutex
	calls := []merakiSDWANCall{}
	next := 0

	client := createMockAlkiraClient(t, func(w http.ResponseWriter, req *http.Request) {
		mu.Lock()
		defer mu.Unlock()

		call := merakiSDWANCall{Method: req.Method, Path: req.URL.Path, Query: req.URL.RawQuery}
		if data, _ := io.ReadAll(req.Body); len(data) > 0 {
			assert.NoError(t, json.Unmarshal(data, &call.Body))
		}
		calls = append(calls, call)

		w.Header().Set("Content-Type", "application/json")

		switch {
		case req.Method == "GET" && req.URL.Path == "/api/credentials/":
			json.NewEncoder(w).Encode([]map[string]string{{"credentialId": "cred-a", "name": "cred-a-name"}})
		case req.Method == "POST" && req.URL.Path == "/api/credentials/meraki-sdwan-connector-instance":
			next++
			w.WriteHeader(http.StatusCreated)
			json.NewEncoder(w).Encode(map[string]string{"id": fmt.Sprintf("cred-new-%d", next)})
		case strings.HasPrefix(req.URL.Path, "/api/credentials/meraki-sdwan-connector-instance/"):
			w.WriteHeader(http.StatusOK)
		case req.Method != "GET" && strings.HasPrefix(req.URL.Path, "/tenantnetworks/0/meraki-sdwan-connectors"):
			if connectorStatus == http.StatusAccepted {
				w.Header().Set("x-ak-validation-state-id", "v1")
			}
			if connectorStatus == http.StatusOK {
				w.Header().Set("x-provision-request-id", "p1")
			}
			w.WriteHeader(connectorStatus)
		case req.URL.Path == "/tenantnetworks/0/validation-states/v1":
			json.NewEncoder(w).Encode(map[string]string{"state": "FAILED"})
		case req.URL.Path == "/tenantnetworks/0/provision-requests/p1":
			json.NewEncoder(w).Encode(map[string]string{"state": "SUCCESS"})
		case req.Method == "GET" && req.URL.Path == "/tenantnetworks/0/meraki-sdwan-connectors/10":
			json.NewEncoder(w).Encode(merakiSDWANConnectorResponse)
		case req.Method == "GET" && req.URL.Path == "/tenantnetworks/0/meraki-sdwan-connectors":
			json.NewEncoder(w).Encode([]interface{}{merakiSDWANConnectorResponse})
		default:
			w.WriteHeader(http.StatusNotFound)
		}
	})
	client.Client.RetryMax = 0

	return client, &calls
}

func merakiSDWANCallIndex(calls []merakiSDWANCall, method, path string) int {
	for i, c := range calls {
		if c.Method == method && c.Path == path {
			return i
		}
	}
	return -1
}

func merakiSDWANCallsTo(calls []merakiSDWANCall, method, prefix string) []merakiSDWANCall {
	var result []merakiSDWANCall
	for _, c := range calls {
		if c.Method == method && strings.HasPrefix(c.Path, prefix) {
			result = append(result, c)
		}
	}
	return result
}

func applyMerakiSDWAN(t *testing.T, client *alkira.AlkiraClient, state *terraform.InstanceState, config *terraform.ResourceConfig) (*terraform.InstanceState, diag.Diagnostics) {
	resource := resourceAlkiraConnectorMerakiSDWAN()

	diff, err := resource.Diff(context.Background(), state, config, client)
	require.NoError(t, err)
	require.NotNil(t, diff)

	return resource.Apply(context.Background(), state, diff, client)
}

func applyMerakiSDWANTokenReplace(t *testing.T, client *alkira.AlkiraClient) (*terraform.InstanceState, diag.Diagnostics) {
	return applyMerakiSDWAN(t, client, merakiSDWANState("ta", "tb"), merakiSDWANConfig(
		map[string]interface{}{"host_name": "vmx-a", "auth_token": "ta2"},
		map[string]interface{}{"host_name": "vmx-b", "auth_token": "tb"},
	))
}

func TestResourceAlkiraConnectorMerakiSDWANFailedCreateDeletesCredentials(t *testing.T) {
	client, calls := merakiSDWANServer(t, http.StatusInternalServerError)

	_, diags := applyMerakiSDWAN(t, client, nil, merakiSDWANConfig(
		map[string]interface{}{"host_name": "vmx-a", "auth_token": "ta"},
		map[string]interface{}{"host_name": "vmx-b", "auth_token": "tb"},
	))
	require.True(t, diags.HasError())

	deletes := merakiSDWANCallsTo(*calls, "DELETE", "/api/credentials/")
	require.Len(t, deletes, 2)
	assert.Equal(t, "/api/credentials/meraki-sdwan-connector-instance/cred-new-1", deletes[0].Path)
	assert.Equal(t, "/api/credentials/meraki-sdwan-connector-instance/cred-new-2", deletes[1].Path)
}

func TestResourceAlkiraConnectorMerakiSDWANFailedUpdate(t *testing.T) {
	client, calls := merakiSDWANServer(t, http.StatusInternalServerError)

	// Replace vmx-a's token and add vmx-c in one apply
	newState, diags := applyMerakiSDWAN(t, client, merakiSDWANState("ta", "tb"), merakiSDWANConfig(
		map[string]interface{}{"host_name": "vmx-a", "auth_token": "ta2"},
		map[string]interface{}{"host_name": "vmx-b", "auth_token": "tb"},
		map[string]interface{}{"host_name": "vmx-c", "auth_token": "tc"},
	))
	require.True(t, diags.HasError())

	tokenPuts := merakiSDWANCallsTo(*calls, "PUT", "/api/credentials/")
	require.Len(t, tokenPuts, 1)
	assert.Equal(t, "/api/credentials/meraki-sdwan-connector-instance/cred-a", tokenPuts[0].Path)
	assert.Equal(t, "cred-a-name", tokenPuts[0].Body["name"])

	connectorPuts := merakiSDWANCallsTo(*calls, "PUT", "/tenantnetworks/0/meraki-sdwan-connectors/10")
	require.Len(t, connectorPuts, 1)
	assert.ElementsMatch(t, []interface{}{
		map[string]interface{}{"id": float64(1), "hostName": "vmx-a", "credentialId": "cred-a"},
		map[string]interface{}{"id": float64(2), "hostName": "vmx-b", "credentialId": "cred-b"},
		map[string]interface{}{"hostName": "vmx-c", "credentialId": "cred-new-1"},
	}, connectorPuts[0].Body["instances"])

	deletes := merakiSDWANCallsTo(*calls, "DELETE", "/api/credentials/")
	require.Len(t, deletes, 1, "only the new instance's credential")
	assert.Equal(t, "/api/credentials/meraki-sdwan-connector-instance/cred-new-1", deletes[0].Path)

	assert.Contains(t, diags[0].Summary, "vmx-a (credential cred-a)",
		"the error must say the token was already replaced")
	assert.Equal(t, "ta", newState.Attributes[merakiSDWANKey("vmx-a", "ta", "auth_token")],
		"a failed update must keep the old token, so the next plan retries the replace")
	assert.Equal(t, "2", newState.Attributes["instances.#"])
}

func TestResourceAlkiraConnectorMerakiSDWANRemoveFirstInstance(t *testing.T) {
	client, calls := merakiSDWANServer(t, http.StatusOK)

	_, diags := applyMerakiSDWAN(t, client, merakiSDWANState("ta", "tb"), merakiSDWANConfig(
		map[string]interface{}{"host_name": "vmx-b", "auth_token": "tb"},
	))
	require.False(t, diags.HasError(), "%v", diags)

	connectorPuts := merakiSDWANCallsTo(*calls, "PUT", "/tenantnetworks/0/meraki-sdwan-connectors/10")
	require.Len(t, connectorPuts, 1)
	assert.Equal(t, []interface{}{
		map[string]interface{}{"id": float64(2), "hostName": "vmx-b", "credentialId": "cred-b"},
	}, connectorPuts[0].Body["instances"], "vmx-b keeps its ID and credential")
	assert.Empty(t, merakiSDWANCallsTo(*calls, "PUT", "/api/credentials/"))
	assert.Empty(t, merakiSDWANCallsTo(*calls, "POST", "/api/credentials/"))
}

func TestResourceAlkiraConnectorMerakiSDWANValidationFailureSavesState(t *testing.T) {
	if testing.Short() {
		t.Skip("waits for one validation poll")
	}

	client, _ := merakiSDWANServer(t, http.StatusAccepted)
	client.Validate = true

	newState, diags := applyMerakiSDWANTokenReplace(t, client)
	require.True(t, diags.HasError())

	assert.Equal(t, "ta2", newState.Attributes[merakiSDWANKey("vmx-a", "ta2", "auth_token")],
		"the connector is saved before validation, so the state must hold the new token")
}

func TestResourceAlkiraConnectorMerakiSDWANTokenReplace(t *testing.T) {
	if testing.Short() {
		t.Skip("waits for one provision poll")
	}

	client, calls := merakiSDWANServer(t, http.StatusOK)
	client.Provision = true

	newState, diags := applyMerakiSDWANTokenReplace(t, client)
	require.False(t, diags.HasError(), "%v", diags)

	credentialPut := merakiSDWANCallIndex(*calls, "PUT", "/api/credentials/meraki-sdwan-connector-instance/cred-a")
	connectorPut := merakiSDWANCallIndex(*calls, "PUT", "/tenantnetworks/0/meraki-sdwan-connectors/10")
	require.NotEqual(t, -1, credentialPut)
	require.NotEqual(t, -1, connectorPut)
	assert.Less(t, credentialPut, connectorPut, "the token is replaced before the update that redeploys")
	assert.Contains(t, (*calls)[connectorPut].Query, "provision=true")
	assert.Equal(t, "ta2", newState.Attributes[merakiSDWANKey("vmx-a", "ta2", "auth_token")])
	assert.Equal(t, "SUCCESS", newState.Attributes["provision_state"])
}

func TestResourceAlkiraConnectorMerakiSDWANRead(t *testing.T) {
	client, _ := merakiSDWANServer(t, http.StatusOK)

	d := resourceAlkiraConnectorMerakiSDWAN().TestResourceData()
	d.SetId("10")

	require.False(t, resourceConnectorMerakiSDWANRead(context.Background(), d, client).HasError())

	assert.Equal(t, "c", d.Get("name"))
	assert.Equal(t, "AZURE-WESTUS2", d.Get("cxp"))
	assert.Equal(t, "SMALL", d.Get("size"))
	assert.Equal(t, "18.2.14", d.Get("version"))
	assert.Equal(t, "IPSEC", d.Get("tunnel_protocol"))
	assert.Equal(t, true, d.Get("enabled"))
	assert.Equal(t, 77, d.Get("implicit_group_id"))
	assert.Equal(t, 1, d.Get("target_segment.0.segment_id"))
	assert.Equal(t, 65500, d.Get("target_segment.0.meraki_bgp_asn"))
	assert.Equal(t, false, d.Get("target_segment.0.advertise_on_prem_routes"))
	assert.Equal(t, false, d.Get("target_segment.0.advertise_default_route"), "disableInternetExit true")
}

func TestResourceAlkiraConnectorMerakiSDWANImport(t *testing.T) {
	client, _ := merakiSDWANServer(t, http.StatusOK)

	resource := resourceAlkiraConnectorMerakiSDWAN()
	d := resource.TestResourceData()
	d.SetId("10")

	result, err := resource.Importer.StateContext(context.Background(), d, client)
	require.NoError(t, err)
	require.Len(t, result, 1)

	instances := merakiSDWANInstancesByHostName(result[0])
	require.Len(t, instances, 2)
	assert.Equal(t, "cred-a", instances["vmx-a"]["credential_id"])
	assert.Equal(t, 1, instances["vmx-a"]["id"])
	assert.Equal(t, "", instances["vmx-a"]["auth_token"], "the token is not imported")
}

// merakiSDWANInstancesByHostName is the instances of d, keyed by host_name
func merakiSDWANInstancesByHostName(d *schema.ResourceData) map[string]map[string]interface{} {
	result := map[string]map[string]interface{}{}
	for _, inst := range d.Get("instances").(*schema.Set).List() {
		i := inst.(map[string]interface{})
		result[i["host_name"].(string)] = i
	}
	return result
}

func TestDataSourceAlkiraConnectorMerakiSDWAN(t *testing.T) {
	client, calls := merakiSDWANServer(t, http.StatusOK)

	d := dataSourceAlkiraConnectorMerakiSDWAN().TestResourceData()
	require.NoError(t, d.Set("name", "c"))

	require.False(t, dataSourceAlkiraConnectorMerakiSDWANRead(context.Background(), d, client).HasError())

	assert.Equal(t, "10", d.Id())
	assert.Equal(t, 77, d.Get("implicit_group_id"))
	assert.Contains(t, (*calls)[0].Query, "name=c")
}

func TestResourceAlkiraConnectorMerakiSDWANDelete(t *testing.T) {
	client, calls := merakiSDWANServer(t, http.StatusOK)

	d := resourceAlkiraConnectorMerakiSDWAN().TestResourceData()
	d.SetId("10")

	require.False(t, resourceConnectorMerakiSDWANDelete(context.Background(), d, client).HasError())

	assert.Equal(t, "", d.Id())
	assert.NotEqual(t, -1, merakiSDWANCallIndex(*calls, "DELETE", "/tenantnetworks/0/meraki-sdwan-connectors/10"))
	assert.Empty(t, merakiSDWANCallsTo(*calls, "DELETE", "/api/credentials/"), "Delete does not delete instance credentials")
}

func TestResourceAlkiraConnectorMerakiSDWANRetryFailedProvision(t *testing.T) {
	state := merakiSDWANState("ta", "tb")
	state.Attributes["provision_state"] = "FAILED"

	diff, err := resourceAlkiraConnectorMerakiSDWAN().Diff(context.Background(), state, merakiSDWANConfig(
		map[string]interface{}{"host_name": "vmx-a", "auth_token": "ta"},
		map[string]interface{}{"host_name": "vmx-b", "auth_token": "tb"},
	), &alkira.AlkiraClient{Provision: true})

	require.NoError(t, err)
	require.NotNil(t, diff)
	assert.Equal(t, "SUCCESS", diff.Attributes["provision_state"].New)
}
