package alkira

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"strings"
	"testing"

	client "github.com/alkiranet/alkira-client-go/alkira"
	"github.com/hashicorp/go-retryablehttp"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

type mcgwConnectorTransport func(*http.Request) (*http.Response, error)

func (f mcgwConnectorTransport) RoundTrip(r *http.Request) (*http.Response, error) { return f(r) }

func mcgwConnectorClient(t *testing.T, handler func(*http.Request) (int, string)) *client.AlkiraClient {
	t.Helper()
	c := retryablehttp.NewClient()
	c.Logger = nil
	c.RetryMax = 0
	c.HTTPClient.Transport = mcgwConnectorTransport(func(r *http.Request) (*http.Response, error) {
		code, body := handler(r)
		return &http.Response{StatusCode: code, Header: make(http.Header), Body: io.NopCloser(strings.NewReader(body)), Request: r}, nil
	})
	return &client.AlkiraClient{Client: c, URI: "https://example.invalid", TenantNetworkId: "1"}
}

// TestConnectorMultiCloudGatewayReadSetsOnRampRegion: Read sets each instance's on_ramp_region
// from the API response, matched by instance, not by list position.
func TestConnectorMultiCloudGatewayReadSetsOnRampRegion(t *testing.T) {
	r := Provider().ResourcesMap["alkira_connector_multi_cloud_gateway"]
	if r == nil {
		t.Fatal("alkira_connector_multi_cloud_gateway resource is not registered")
	}

	d := schema.TestResourceDataRaw(t, r.Schema, nil)
	d.SetId("42")

	connectorJSON := `{"id":42,"name":"test","cxp":"US-WEST","size":"SMALL","enabled":true,` +
		`"multiCloudGatewayId":1,"loopbackPrefix":"10.0.0.0/26",` +
		`"segmentOptions":[{"segmentName":"segment1"}],` +
		`"instances":[{"id":1,"name":"a","onRampRegion":"us-sanjose-1 (San Jose)"},` +
		`{"id":2,"name":"b","onRampRegion":"us-ashburn-1 (Ashburn)"}]}`

	c := mcgwConnectorClient(t, func(req *http.Request) (int, string) {
		switch {
		case req.Method == "GET" && req.URL.Path == "/tenantnetworks/1/multi-cloud-gateway-connectors/42":
			return 200, connectorJSON
		case req.Method == "GET" && strings.HasPrefix(req.URL.Path, "/tenantnetworks/1/segments"):
			return 200, `[{"id":1,"name":"segment1"}]`
		default:
			t.Fatalf("unexpected request: %s %s", req.Method, req.URL.Path)
			return 500, ""
		}
	})

	if diags := r.ReadContext(context.Background(), d, c); diags.HasError() {
		t.Fatal(diags)
	}

	want := map[string]string{"a": "us-sanjose-1 (San Jose)", "b": "us-ashburn-1 (Ashburn)"}
	instances := d.Get("instance").([]interface{})
	if len(instances) != len(want) {
		t.Fatalf("got %d instances, want %d", len(instances), len(want))
	}
	for _, v := range instances {
		cfg := v.(map[string]interface{})
		name := cfg["name"].(string)
		if cfg["on_ramp_region"] != want[name] {
			t.Errorf("instance %s on_ramp_region = %v, want %v", name, cfg["on_ramp_region"], want[name])
		}
	}
}

// TestConnectorMultiCloudGatewayUpdateRemovesInstance exercises the real Update path
// (generateConnectorMultiCloudGatewayRequest + d.GetChange): state has instances [a:1, b:2, c:3],
// config removes b, and the PUT body sent must carry instances [{id:1,name:a},{id:3,name:c}].
func TestConnectorMultiCloudGatewayUpdateRemovesInstance(t *testing.T) {
	r := Provider().ResourcesMap["alkira_connector_multi_cloud_gateway"]
	if r == nil {
		t.Fatal("alkira_connector_multi_cloud_gateway resource is not registered")
	}

	oldValues := map[string]interface{}{
		"name":                   "test",
		"cxp":                    "US-WEST",
		"size":                   "SMALL",
		"multi_cloud_gateway_id": 1,
		"loopback_prefix":        "10.0.0.0/26",
		"segment_options": []interface{}{
			map[string]interface{}{"segment_id": "1"},
		},
		"instance": []interface{}{
			map[string]interface{}{"name": "a", "id": 1},
			map[string]interface{}{"name": "b", "id": 2},
			map[string]interface{}{"name": "c", "id": 3},
		},
	}
	old := schema.TestResourceDataRaw(t, r.Schema, oldValues)
	old.SetId("42")

	newValues := map[string]interface{}{
		"name":                   "test",
		"cxp":                    "US-WEST",
		"size":                   "SMALL",
		"multi_cloud_gateway_id": 1,
		"loopback_prefix":        "10.0.0.0/26",
		"segment_options": []interface{}{
			map[string]interface{}{"segment_id": "1"},
		},
		"instance": []interface{}{
			map[string]interface{}{"name": "a"},
			map[string]interface{}{"name": "c"},
		},
	}

	connectorJSON := `{"id":42,"name":"test","cxp":"US-WEST","size":"SMALL","enabled":true,` +
		`"multiCloudGatewayId":1,"loopbackPrefix":"10.0.0.0/26",` +
		`"segmentOptions":[{"segmentName":"segment1"}],` +
		`"instances":[{"id":1,"name":"a"},{"id":3,"name":"c"}]}`

	var puts int
	c := mcgwConnectorClient(t, func(req *http.Request) (int, string) {
		switch {
		case req.Method == "PUT" && req.URL.Path == "/tenantnetworks/1/multi-cloud-gateway-connectors/42":
			puts++
			var body map[string]interface{}
			if err := json.NewDecoder(req.Body).Decode(&body); err != nil {
				t.Fatal(err)
			}
			instances, ok := body["instances"].([]interface{})
			if !ok || len(instances) != 2 {
				t.Fatalf("unexpected instances in PUT body: %#v", body["instances"])
			}
			want := []map[string]interface{}{
				{"id": float64(1), "name": "a"},
				{"id": float64(3), "name": "c"},
			}
			for i, w := range want {
				got := instances[i].(map[string]interface{})
				if got["id"] != w["id"] || got["name"] != w["name"] {
					t.Fatalf("instance[%d] = %#v, want %#v", i, got, w)
				}
			}
			return 200, connectorJSON
		case req.Method == "GET" && req.URL.Path == "/tenantnetworks/1/multi-cloud-gateway-connectors/42":
			return 200, connectorJSON
		case req.Method == "GET" && req.URL.Path == "/tenantnetworks/1/segments/1":
			return 200, `{"id":1,"name":"segment1"}`
		case req.Method == "GET" && req.URL.Path == "/tenantnetworks/1/segments":
			return 200, `[{"id":1,"name":"segment1"}]`
		default:
			t.Fatalf("unexpected request: %s %s", req.Method, req.URL.Path)
			return 500, ""
		}
	})

	diff, err := r.Diff(context.Background(), old.State(), terraform.NewResourceConfigRaw(newValues), c)
	if err != nil {
		t.Fatal(err)
	}
	if diff == nil {
		t.Fatal("expected a diff removing instance b")
	}

	if _, diags := r.Apply(context.Background(), old.State(), diff, c); diags.HasError() {
		t.Fatal(diags)
	}
	if puts != 1 {
		t.Fatalf("PUT count = %d", puts)
	}
}
