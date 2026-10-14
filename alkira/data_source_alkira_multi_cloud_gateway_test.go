package alkira

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"reflect"
	"strings"
	"testing"

	client "github.com/alkiranet/alkira-client-go/alkira"
	"github.com/hashicorp/go-retryablehttp"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

type multiCloudGatewayDataSourceTransport func(*http.Request) (*http.Response, error)

func (f multiCloudGatewayDataSourceTransport) RoundTrip(r *http.Request) (*http.Response, error) {
	return f(r)
}

func TestReadByName(t *testing.T) {
	const gateway = `{"id":9223372036854775806,"name":"gateway","nativeId":"native-42","tier":"50GBPS","state":"PROVISIONED"}`
	const other = `{"id":7,"name":"gateway-other"}`
	for _, tc := range []struct {
		name      string
		pages     []string
		offsets   []string
		status    int
		wantError string
	}{
		{name: "single_match", pages: []string{`{"data":[` + gateway + `],"pagination":{"hits":1}}`}, offsets: []string{"0"}},
		{name: "later_page", pages: []string{`{"data":[` + other + `],"pagination":{"hits":2}}`, `{"data":[` + gateway + `],"pagination":{"hits":2}}`}, offsets: []string{"0", "1"}},
		{name: "duplicate_across_pages", pages: []string{`{"data":[` + gateway + `],"pagination":{"hits":2}}`, `{"data":[` + gateway + `],"pagination":{"hits":2}}`}, offsets: []string{"0", "1"}, wantError: "found 2 Multi-Cloud Gateways"},
		{name: "no_exact_match", pages: []string{`{"data":[` + other + `],"pagination":{"hits":1}}`}, offsets: []string{"0"}, wantError: "found 0 Multi-Cloud Gateways"},
		{name: "empty_list", pages: []string{`{"data":[],"pagination":{"hits":0}}`}, offsets: []string{"0"}, wantError: "found 0 Multi-Cloud Gateways"},
		{name: "api_error", pages: []string{`{"message":"unavailable"}`}, offsets: []string{"0"}, status: 500, wantError: "failed to get page"},
		{name: "malformed_response", pages: []string{`{"data":`}, offsets: []string{"0"}, wantError: "failed to unmarshal page"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			resource := Provider().DataSourcesMap["alkira_multi_cloud_gateway"]
			if resource == nil {
				t.Fatal("alkira_multi_cloud_gateway data source is not registered")
			}
			if !resource.Schema["name"].Required {
				t.Fatal("name must be required")
			}
			for _, key := range []string{"native_id", "tier", "state"} {
				field := resource.Schema[key]
				if field == nil || !field.Computed || field.Optional || field.Required {
					t.Fatalf("%s must be computed only", key)
				}
			}
			if diags := resource.Validate(terraform.NewResourceConfigRaw(map[string]interface{}{})); !diags.HasError() {
				t.Fatal("missing name accepted")
			}
			d := schema.TestResourceDataRaw(t, resource.Schema, map[string]interface{}{"name": "gateway"})
			retryClient := retryablehttp.NewClient()
			retryClient.Logger = nil
			retryClient.RetryMax = 0
			var offsets []string
			retryClient.HTTPClient.Transport = multiCloudGatewayDataSourceTransport(func(req *http.Request) (*http.Response, error) {
				if req.Method != "GET" || req.URL.Path != "/tenantnetworks/1/multi-cloud-gateways" {
					t.Errorf("unexpected request: %s %s", req.Method, req.URL.Path)
				}
				query := req.URL.Query()
				if len(query) != 2 || query.Get("limit") != "100" || !query.Has("offset") {
					t.Errorf("expected offset/limit pagination, got %s", req.URL.RawQuery)
				}
				index := len(offsets)
				offsets = append(offsets, query.Get("offset"))
				if index >= len(tc.pages) {
					return nil, fmt.Errorf("unexpected extra page request")
				}
				status := tc.status
				if status == 0 {
					status = 200
				}
				return &http.Response{StatusCode: status, Header: make(http.Header), Body: io.NopCloser(strings.NewReader(tc.pages[index])), Request: req}, nil
			})
			c := &client.AlkiraClient{Client: retryClient, URI: "https://example.invalid", TenantNetworkId: "1"}
			diags := resource.ReadContext(context.Background(), d, c)
			if !reflect.DeepEqual(offsets, tc.offsets) {
				t.Errorf("requested offsets %v, want %v", offsets, tc.offsets)
			}
			if tc.wantError != "" {
				if !diags.HasError() || !strings.Contains(fmt.Sprint(diags), tc.wantError) {
					t.Fatalf("expected %q, got %v", tc.wantError, diags)
				}
				if d.Id() != "" {
					t.Fatalf("failed lookup set ID %q", d.Id())
				}
				return
			}
			if diags.HasError() {
				t.Fatal(diags)
			}
			if d.Id() != "9223372036854775806" {
				t.Fatalf("Alkira int64 ID lost precision: %q", d.Id())
			}
			for key, want := range map[string]string{"name": "gateway", "native_id": "native-42", "tier": "50GBPS", "state": "PROVISIONED"} {
				if got := d.Get(key); got != want {
					t.Errorf("%s = %v, want %s", key, got, want)
				}
			}
		})
	}
}
