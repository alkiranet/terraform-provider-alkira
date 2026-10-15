package alkira

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"testing"
	"time"

	client "github.com/alkiranet/alkira-client-go/alkira"
	"github.com/hashicorp/go-retryablehttp"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func mcgwResource(t *testing.T) *schema.Resource {
	t.Helper()
	r := Provider().ResourcesMap["alkira_multi_cloud_gateway"]
	if r == nil {
		t.Fatal("alkira_multi_cloud_gateway resource is not registered")
	}
	return r
}

type mcgwTransport func(*http.Request) (*http.Response, error)

func (f mcgwTransport) RoundTrip(r *http.Request) (*http.Response, error) { return f(r) }
func mcgwClient(t *testing.T, handler func(*http.Request) (int, string)) *client.AlkiraClient {
	t.Helper()
	c := retryablehttp.NewClient()
	c.Logger = nil
	c.RetryMax = 0
	c.HTTPClient.Transport = mcgwTransport(func(r *http.Request) (*http.Response, error) {
		if !strings.HasPrefix(r.URL.Path, "/tenantnetworks/1/multi-cloud-gateways") {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		code, body := handler(r)
		return &http.Response{StatusCode: code, Header: make(http.Header), Body: io.NopCloser(strings.NewReader(body)), Request: r}, nil
	})
	return &client.AlkiraClient{Client: c, URI: "https://example.invalid", TenantNetworkId: "1"}
}
func mcgwValues() map[string]interface{} {
	return map[string]interface{}{"name": "gateway", "tier": "10GBPS", "description": "original"}
}
func mcgwJSON(tier, description, state string) string {
	return fmt.Sprintf(`{"id":42,"name":"gateway","tier":%q,"description":%q,"nativeId":"native-42","state":%q}`, tier, description, state)
}
func TestOnlyNameForceNew(t *testing.T) {
	r := mcgwResource(t)
	for name, s := range r.Schema {
		if s.ForceNew != (name == "name") {
			t.Errorf("%s ForceNew = %v", name, s.ForceNew)
		}
	}
}
func TestTierRequired(t *testing.T) {
	r := mcgwResource(t)
	if !r.Schema["tier"].Required {
		t.Fatal("tier must be required")
	}
	errors := r.Validate(terraform.NewResourceConfigRaw(map[string]interface{}{"name": "gateway"}))
	if len(errors) == 0 {
		t.Fatal("missing tier must fail validation")
	}
	for _, tier := range []string{"10GBPS", "50GBPS", "UNLIMITED"} {
		errors = r.Validate(terraform.NewResourceConfigRaw(map[string]interface{}{"name": "gateway", "tier": tier}))
		if len(errors) != 0 {
			t.Fatal(errors)
		}
	}
}
func TestTierAndDescriptionUpdateInPlace(t *testing.T) {
	for _, field := range []string{"tier", "description", "clear_description"} {
		t.Run(field, func(t *testing.T) {
			r := mcgwResource(t)
			old := schema.TestResourceDataRaw(t, r.Schema, mcgwValues())
			old.SetId("42")
			values := mcgwValues()
			switch field {
			case "tier":
				values[field] = "50GBPS"
			case "clear_description":
				field = "description"
				values[field] = ""
			default:
				values[field] = "changed"
			}
			diff, err := r.Diff(context.Background(), old.State(), terraform.NewResourceConfigRaw(values), nil)
			if err != nil {
				t.Fatal(err)
			}
			if diff == nil || diff.Attributes[field] == nil || diff.RequiresNew() {
				t.Fatalf("expected in-place %s diff: %#v", field, diff)
			}
			puts := 0
			c := mcgwClient(t, func(req *http.Request) (int, string) {
				if req.Method == "PUT" {
					puts++
					if req.URL.Path != "/tenantnetworks/1/multi-cloud-gateways/42" {
						t.Error(req.URL.Path)
					}
					var body map[string]interface{}
					if err := json.NewDecoder(req.Body).Decode(&body); err != nil {
						t.Fatal(err)
					}
					if len(body) != 3 || body["tier"] != values["tier"] || body["description"] != values["description"] || body["name"] != "gateway" {
						t.Errorf("unexpected PUT payload: %v", body)
					}
				} else if req.Method != "GET" {
					t.Errorf("unexpected method %s", req.Method)
				}
				return 200, mcgwJSON(values["tier"].(string), values["description"].(string), "PROVISIONED")
			})
			_, diags := r.Apply(context.Background(), old.State(), diff, c)
			if diags.HasError() {
				t.Fatal(diags)
			}
			if puts != 1 {
				t.Fatalf("PUT count = %d", puts)
			}
		})
	}
}
func TestReadSetsTierAndDescription(t *testing.T) {
	r := mcgwResource(t)
	d := schema.TestResourceDataRaw(t, r.Schema, mcgwValues())
	d.SetId("42")
	c := mcgwClient(t, func(req *http.Request) (int, string) { return 200, mcgwJSON("50GBPS", "external", "PROVISIONED") })
	if diags := r.ReadContext(context.Background(), d, c); diags.HasError() {
		t.Fatal(diags)
	}
	for k, v := range map[string]string{"tier": "50GBPS", "description": "external", "native_id": "native-42", "state": "PROVISIONED"} {
		if d.Get(k) != v {
			t.Errorf("%s = %v", k, d.Get(k))
		}
	}
	diff, err := r.Diff(context.Background(), d.State(), terraform.NewResourceConfigRaw(mcgwValues()), nil)
	if err != nil {
		t.Fatal(err)
	}
	if diff == nil || diff.RequiresNew() || diff.Attributes["tier"] == nil || diff.Attributes["description"] == nil {
		t.Fatalf("expected in-place drift diff: %#v", diff)
	}
}
func TestCreateWaitsForProvisioned(t *testing.T) {
	for _, tc := range []struct {
		name      string
		states    []string
		timeout   time.Duration
		wantError bool
	}{
		{"provisioned", []string{"PROVISIONING", "PROVISIONED"}, 10 * time.Second, false},
		{"failed", []string{"FAILED"}, time.Second, true},
		{"timeout", []string{"PROVISIONING"}, 50 * time.Millisecond, true},
		{"get_error", []string{"ERROR"}, time.Second, true},
	} {
		t.Run(tc.name, func(t *testing.T) {
			r := mcgwResource(t)
			d := schema.TestResourceDataRaw(t, r.Schema, mcgwValues())
			r.Timeouts = &schema.ResourceTimeout{Create: &tc.timeout}
			d = r.Data(d.State())
			gets := 0
			c := mcgwClient(t, func(req *http.Request) (int, string) {
				if req.Method == "POST" {
					return 201, mcgwJSON("10GBPS", "original", "PROVISIONING")
				}
				if req.Method != "GET" {
					t.Errorf("unexpected method %s", req.Method)
				}
				state := tc.states[min(gets, len(tc.states)-1)]
				gets++
				if state == "ERROR" {
					return 500, `{"message":"unavailable"}`
				}
				return 200, mcgwJSON("10GBPS", "original", state)
			})
			start := time.Now()
			diags := r.CreateContext(context.Background(), d, c)
			if diags.HasError() != tc.wantError {
				t.Fatalf("diagnostics: %v", diags)
			}
			if d.Id() != "42" {
				t.Fatalf("created ID must be retained: %s", d.Id())
			}
			if !tc.wantError && (gets < 2 || d.Get("state") != "PROVISIONED") {
				t.Fatalf("returned before provisioned: gets=%d state=%v", gets, d.Get("state"))
			}
			if tc.name == "failed" && !strings.Contains(fmt.Sprint(diags), "FAILED") {
				t.Fatal(diags)
			}
			if tc.name == "timeout" && time.Since(start) > time.Second {
				t.Fatal("create timeout did not bound wait")
			}
		})
	}
}

func TestMultiCloudGatewayImportAndDelete(t *testing.T) {
	r := mcgwResource(t)
	for _, code := range []int{200, 404} {
		t.Run(fmt.Sprint(code), func(t *testing.T) {
			d := schema.TestResourceDataRaw(t, r.Schema, nil)
			d.SetId("42")
			c := mcgwClient(t, func(req *http.Request) (int, string) { return code, mcgwJSON("50GBPS", "imported", "PROVISIONED") })
			states, err := r.Importer.StateContext(context.Background(), d, c)
			if code == 404 {
				if err == nil {
					t.Fatal("missing gateway import succeeded")
				}
				return
			}
			if err != nil || len(states) != 1 || d.Get("tier") != "50GBPS" {
				t.Fatalf("import: %v %v", states, err)
			}
		})
	}
	d := schema.TestResourceDataRaw(t, r.Schema, nil)
	d.SetId("../invalid")
	if _, err := r.Importer.StateContext(context.Background(), d, nil); err == nil {
		t.Fatal("invalid import ID accepted")
	}
	d.SetId("42")
	c := mcgwClient(t, func(req *http.Request) (int, string) {
		if req.Method != "DELETE" || req.URL.Path != "/tenantnetworks/1/multi-cloud-gateways/42" {
			t.Errorf("unexpected delete: %s %s", req.Method, req.URL.Path)
		}
		return 204, ""
	})
	if diags := r.DeleteContext(context.Background(), d, c); diags.HasError() || d.Id() != "" {
		t.Fatalf("delete: %v, ID %q", diags, d.Id())
	}
}

func TestMultiCloudGatewayCreateCancelsRequest(t *testing.T) {
	for _, stage := range []string{"POST", "GET"} {
		t.Run(stage, func(t *testing.T) {
			r := mcgwResource(t)
			r.Timeouts = &schema.ResourceTimeout{Create: schema.DefaultTimeout(50 * time.Millisecond)}
			d := r.Data(nil)
			for k, v := range mcgwValues() {
				if err := d.Set(k, v); err != nil {
					t.Fatal(err)
				}
			}
			c := mcgwClient(t, func(req *http.Request) (int, string) { return 201, mcgwJSON("10GBPS", "original", "PROVISIONING") })
			transport := c.Client.HTTPClient.Transport
			c.Client.HTTPClient.Transport = mcgwTransport(func(req *http.Request) (*http.Response, error) {
				if req.Method == stage {
					<-req.Context().Done()
					return nil, req.Context().Err()
				}
				return transport.RoundTrip(req)
			})
			start := time.Now()
			diags := r.CreateContext(context.Background(), d, c)
			if !diags.HasError() || time.Since(start) > time.Second {
				t.Fatalf("request not bounded by timeout: %v", diags)
			}
			if stage == "GET" && d.Id() != "42" {
				t.Fatal("lost created ID")
			}
			if c.Client.RequestLogHook != nil {
				t.Fatal("shared client was modified")
			}
		})
	}
}
