package alkira

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/alkiranet/alkira-client-go/alkira"
	"github.com/hashicorp/go-retryablehttp"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAlkiraServiceCiscoFTDv_resourceSchema(t *testing.T) {
	resource := resourceAlkiraServiceCiscoFTDv()

	// Test required fields
	nameSchema := resource.Schema["name"]
	assert.True(t, nameSchema.Required, "Name should be required")
	assert.Equal(t, schema.TypeString, nameSchema.Type, "Name should be string type")

	cxpSchema := resource.Schema["cxp"]
	assert.True(t, cxpSchema.Required, "CXP should be required")
	assert.Equal(t, schema.TypeString, cxpSchema.Type, "CXP should be string type")

	sizeSchema := resource.Schema["size"]
	assert.True(t, sizeSchema.Required, "Size should be required")
	assert.Equal(t, schema.TypeString, sizeSchema.Type, "Size should be string type")

	globalCidrListIdSchema := resource.Schema["global_cidr_list_id"]
	assert.True(t, globalCidrListIdSchema.Required, "Global CIDR list ID should be required")
	assert.Equal(t, schema.TypeInt, globalCidrListIdSchema.Type, "Global CIDR list ID should be int type")

	segmentsSchema := resource.Schema["segments"]
	assert.True(t, segmentsSchema.Required, "Segments should be required")
	assert.Equal(t, schema.TypeSet, segmentsSchema.Type, "Segments should be set type")

	managementServerIpSchema := resource.Schema["management_server_ip"]
	assert.True(t, managementServerIpSchema.Required, "Management server IP should be required")
	assert.Equal(t, schema.TypeString, managementServerIpSchema.Type, "Management server IP should be string type")

	managementServerSegmentSchema := resource.Schema["management_server_segment"]
	assert.True(t, managementServerSegmentSchema.Required, "Management server segment should be required")
	assert.Equal(t, schema.TypeString, managementServerSegmentSchema.Type, "Management server segment should be string type")

	managementServerSegmentIdSchema := resource.Schema["management_server_segment_id"]
	assert.True(t, managementServerSegmentIdSchema.Required, "Management server segment ID should be required")
	assert.Equal(t, schema.TypeInt, managementServerSegmentIdSchema.Type, "Management server segment ID should be int type")

	// Test optional fields with defaults
	autoScaleSchema := resource.Schema["auto_scale"]
	assert.True(t, autoScaleSchema.Optional, "Auto scale should be optional")
	assert.Equal(t, schema.TypeString, autoScaleSchema.Type, "Auto scale should be string type")
	assert.Equal(t, "OFF", autoScaleSchema.Default, "Auto scale should default to OFF")

	// Test optional fields
	descSchema := resource.Schema["description"]
	assert.True(t, descSchema.Optional, "Description should be optional")
	assert.Equal(t, schema.TypeString, descSchema.Type, "Description should be string type")

	maxInstanceCountSchema := resource.Schema["max_instance_count"]
	assert.True(t, maxInstanceCountSchema.Optional, "Max instance count should be optional")
	assert.Equal(t, schema.TypeInt, maxInstanceCountSchema.Type, "Max instance count should be int type")

	minInstanceCountSchema := resource.Schema["min_instance_count"]
	assert.True(t, minInstanceCountSchema.Optional, "Min instance count should be optional")
	assert.Equal(t, schema.TypeInt, minInstanceCountSchema.Type, "Min instance count should be int type")

	tunnelProtocolSchema := resource.Schema["tunnel_protocol"]
	assert.True(t, tunnelProtocolSchema.Optional, "Tunnel protocol should be optional")
	assert.Equal(t, schema.TypeString, tunnelProtocolSchema.Type, "Tunnel protocol should be string type")

	ipAllowListSchema := resource.Schema["ip_allow_list"]
	assert.True(t, ipAllowListSchema.Optional, "IP allow list should be optional")
	assert.Equal(t, schema.TypeSet, ipAllowListSchema.Type, "IP allow list should be set type")

	billingTagsSchema := resource.Schema["billing_tags"]
	assert.True(t, billingTagsSchema.Optional, "Billing tags should be optional")
	assert.Equal(t, schema.TypeSet, billingTagsSchema.Type, "Billing tags should be set type")

	instancesSchema := resource.Schema["instances"]
	assert.True(t, instancesSchema.Optional, "Instances should be optional")
	assert.Equal(t, schema.TypeSet, instancesSchema.Type, "Instances should be set type")

	// Test computed fields
	provStateSchema := resource.Schema["provision_state"]
	assert.True(t, provStateSchema.Computed, "Provision state should be computed")
	assert.Equal(t, schema.TypeString, provStateSchema.Type, "Provision state should be string type")

	internalNameSchema := resource.Schema["internal_name"]
	assert.True(t, internalNameSchema.Computed, "Internal name should be computed")
	assert.Equal(t, schema.TypeString, internalNameSchema.Type, "Internal name should be string type")

	stateSchema := resource.Schema["state"]
	assert.True(t, stateSchema.Computed, "State should be computed")
	assert.Equal(t, schema.TypeString, stateSchema.Type, "State should be string type")

	credentialIdSchema := resource.Schema["credential_id"]
	assert.True(t, credentialIdSchema.Computed, "Credential ID should be computed")
	assert.Equal(t, schema.TypeString, credentialIdSchema.Type, "Credential ID should be string type")

	// Test that resource has all required CRUD functions
	assert.NotNil(t, resource.CreateContext, "Resource should have CreateContext")
	assert.NotNil(t, resource.ReadContext, "Resource should have ReadContext")
	assert.NotNil(t, resource.UpdateContext, "Resource should have UpdateContext")
	assert.NotNil(t, resource.DeleteContext, "Resource should have DeleteContext")
	assert.NotNil(t, resource.Importer, "Resource should support import")
	assert.NotNil(t, resource.CustomizeDiff, "Resource should have CustomizeDiff")
}

func TestAlkiraServiceCiscoFTDv_validateAutoScale(t *testing.T) {
	tests := []struct {
		name      string
		input     interface{}
		expectErr bool
		errCount  int
	}{
		{
			name:      "Valid ON",
			input:     "ON",
			expectErr: false,
			errCount:  0,
		},
		{
			name:      "Valid OFF",
			input:     "OFF",
			expectErr: false,
			errCount:  0,
		},
		{
			name:      "Invalid auto scale value",
			input:     "INVALID",
			expectErr: true,
			errCount:  1,
		},
		{
			name:      "Empty string",
			input:     "",
			expectErr: true,
			errCount:  1,
		},
		{
			name:      "Non-string input",
			input:     123,
			expectErr: true,
			errCount:  1,
		},
	}

	resource := resourceAlkiraServiceCiscoFTDv()
	autoScaleSchema := resource.Schema["auto_scale"]

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			warnings, errors := autoScaleSchema.ValidateFunc(tt.input, "auto_scale")

			if tt.expectErr {
				assert.Len(t, errors, tt.errCount, "Expected %d errors for input %v", tt.errCount, tt.input)
			} else {
				assert.Empty(t, errors, "Expected no errors for input %v", tt.input)
			}
			assert.Empty(t, warnings, "Expected no warnings")
		})
	}
}

func TestAlkiraServiceCiscoFTDv_validateSize(t *testing.T) {
	tests := []struct {
		name      string
		input     interface{}
		expectErr bool
		errCount  int
	}{
		{
			name:      "Valid SMALL",
			input:     "SMALL",
			expectErr: false,
			errCount:  0,
		},
		{
			name:      "Valid MEDIUM",
			input:     "MEDIUM",
			expectErr: false,
			errCount:  0,
		},
		{
			name:      "Valid LARGE",
			input:     "LARGE",
			expectErr: false,
			errCount:  0,
		},
		{
			name:      "Invalid size",
			input:     "INVALID_SIZE",
			expectErr: true,
			errCount:  1,
		},
		{
			name:      "Empty string",
			input:     "",
			expectErr: true,
			errCount:  1,
		},
	}

	resource := resourceAlkiraServiceCiscoFTDv()
	sizeSchema := resource.Schema["size"]

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			warnings, errors := sizeSchema.ValidateFunc(tt.input, "size")

			if tt.expectErr {
				assert.Len(t, errors, tt.errCount, "Expected %d errors for input %v", tt.errCount, tt.input)
			} else {
				assert.Empty(t, errors, "Expected no errors for input %v", tt.input)
			}
			assert.Empty(t, warnings, "Expected no warnings")
		})
	}
}

// TEST HELPER
func serveServiceCiscoFTDv(t *testing.T, serviceCiscoFTDv *alkira.ServiceCiscoFTDv) *alkira.AlkiraClient {
	server := httptest.NewServer(http.HandlerFunc(
		func(w http.ResponseWriter, req *http.Request) {
			json.NewEncoder(w).Encode(serviceCiscoFTDv)
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
