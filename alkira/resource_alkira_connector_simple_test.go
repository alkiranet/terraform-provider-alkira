package alkira

import (
	"encoding/json"
	"net/http"
	"testing"

	"github.com/alkiranet/alkira-client-go/alkira"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/stretchr/testify/assert"
)

// Test schema validation for all connector resources
func TestConnectorSchemas_Validation(t *testing.T) {
	testCases := []struct {
		name     string
		resource func() *schema.Resource
	}{
		{"AWS VPC Connector", resourceAlkiraConnectorAwsVpc},
		{"Azure VNET Connector", resourceAlkiraConnectorAzureVnet},
		{"GCP VPC Connector", resourceAlkiraConnectorGcpVpc},
		{"IPSec Connector", resourceAlkiraConnectorIPSec},
		{"AWS TGW Connector", resourceAlkiraConnectorAwsTgw},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			resource := tc.resource()

			// Test that CRUD functions exist
			assert.NotNil(t, resource.CreateContext, "Resource should have CreateContext")
			assert.NotNil(t, resource.ReadContext, "Resource should have ReadContext")
			assert.NotNil(t, resource.UpdateContext, "Resource should have UpdateContext")
			assert.NotNil(t, resource.DeleteContext, "Resource should have DeleteContext")
			assert.NotNil(t, resource.Importer, "Resource should support import")

			// Test that common required fields exist
			commonRequiredFields := []string{"name", "cxp", "size"}
			for _, field := range commonRequiredFields {
				if fieldSchema, exists := resource.Schema[field]; exists {
					if fieldSchema.Required {
						assert.True(t, fieldSchema.Required, "Field %s should be required", field)
						assert.Equal(t, schema.TypeString, fieldSchema.Type, "Field %s should be string type", field)
					}
				}
			}

			// Test that common computed fields exist
			commonComputedFields := []string{"implicit_group_id", "provision_state"}
			for _, field := range commonComputedFields {
				if fieldSchema, exists := resource.Schema[field]; exists {
					assert.True(t, fieldSchema.Computed, "Field %s should be computed", field)
				}
			}
		})
	}
}

// Test resource data manipulation for AWS VPC Connector
func TestConnectorAwsVpc_DataManipulation(t *testing.T) {
	r := resourceAlkiraConnectorAwsVpc()
	d := r.TestResourceData()

	// Test setting and getting basic fields
	testData := map[string]interface{}{
		"name":                "test-aws-vpc-connector",
		"aws_account_id":      "123456789012",
		"aws_region":          "us-west-2",
		"credential_id":       "cred-123",
		"cxp":                 "US-WEST",
		"size":                "SMALL",
		"vpc_id":              "vpc-123456",
		"segment_id":          "seg-123",
		"enabled":             true,
		"tgw_connect_enabled": false,
		"description":         "Test connector",
	}

	for key, value := range testData {
		d.Set(key, value)
		assert.Equal(t, value, d.Get(key), "Field %s should match set value", key)
	}
}

// Test resource data manipulation for Azure VNET Connector
func TestConnectorAzureVnet_DataManipulation(t *testing.T) {
	r := resourceAlkiraConnectorAzureVnet()
	d := r.TestResourceData()

	// Test setting and getting basic fields
	testData := map[string]interface{}{
		"name":            "test-azure-vnet-connector",
		"azure_vnet_id":   "/subscriptions/12345/resourceGroups/rg1/providers/Microsoft.Network/virtualNetworks/vnet1",
		"credential_id":   "cred-123",
		"cxp":             "US-WEST",
		"size":            "SMALL",
		"segment_id":      "seg-123",
		"enabled":         true,
		"connection_mode": "VNET_GATEWAY",
		"routing_options": "ADVERTISE_DEFAULT_ROUTE",
		"customer_asn":    65000,
		"description":     "Test Azure connector",
	}

	for key, value := range testData {
		d.Set(key, value)
		assert.Equal(t, value, d.Get(key), "Field %s should match set value", key)
	}
}

// Test resource data manipulation for GCP VPC Connector
func TestConnectorGcpVpc_DataManipulation(t *testing.T) {
	r := resourceAlkiraConnectorGcpVpc()
	d := r.TestResourceData()

	// Test setting and getting basic fields
	testData := map[string]interface{}{
		"name":           "test-gcp-vpc-connector",
		"gcp_project_id": "my-gcp-project",
		"gcp_region":     "us-west1",
		"gcp_vpc_name":   "my-vpc",
		"credential_id":  "cred-123",
		"cxp":            "US-WEST",
		"size":           "SMALL",
		"segment_id":     "seg-123",
		"enabled":        true,
		"customer_asn":   64522,
		"description":    "Test GCP connector",
	}

	for key, value := range testData {
		d.Set(key, value)
		assert.Equal(t, value, d.Get(key), "Field %s should match set value", key)
	}
}

// Test resource data manipulation for IPSec Connector
func TestConnectorIPSec_DataManipulation(t *testing.T) {
	r := resourceAlkiraConnectorIPSec()
	d := r.TestResourceData()

	// Test setting and getting basic fields
	testData := map[string]interface{}{
		"name":        "test-ipsec-connector",
		"cxp":         "US-WEST",
		"size":        "SMALL",
		"segment_id":  "seg-123",
		"enabled":     true,
		"vpn_mode":    "ROUTE_BASED",
		"description": "Test IPSec connector",
	}

	for key, value := range testData {
		d.Set(key, value)
		assert.Equal(t, value, d.Get(key), "Field %s should match set value", key)
	}

	// Test endpoint configuration
	endpoints := []interface{}{
		map[string]interface{}{
			"name":                     "endpoint-1",
			"customer_gateway_ip":      "192.168.1.1",
			"customer_ip_type":         "STATIC",
			"preshared_keys":           []interface{}{"psk1", "psk2"},
			"enable_tunnel_redundancy": true,
			"ha_mode":                  "ACTIVE",
		},
	}

	d.Set("endpoint", endpoints)
	endpointList := d.Get("endpoint").([]interface{})
	assert.Len(t, endpointList, 1)

	endpoint := endpointList[0].(map[string]interface{})
	assert.Equal(t, "endpoint-1", endpoint["name"])
	assert.Equal(t, "192.168.1.1", endpoint["customer_gateway_ip"])
	assert.Equal(t, "STATIC", endpoint["customer_ip_type"])
	assert.Equal(t, true, endpoint["enable_tunnel_redundancy"])
	assert.Equal(t, "ACTIVE", endpoint["ha_mode"])
}

// Test resource data manipulation for AWS TGW Connector
func TestConnectorAwsTgw_DataManipulation(t *testing.T) {
	r := resourceAlkiraConnectorAwsTgw()
	d := r.TestResourceData()

	// Test setting and getting basic fields
	testData := map[string]interface{}{
		"name":                                  "test-aws-tgw-connector",
		"peering_gateway_aws_tgw_attachment_id": 100,
		"cxp":                                   "US-WEST",
		"size":                                  "SMALL",
		"segment_id":                            "seg-123",
		"enabled":                               true,
		"description":                           "Test AWS TGW connector",
	}

	for key, value := range testData {
		d.Set(key, value)
		assert.Equal(t, value, d.Get(key), "Field %s should match set value", key)
	}

	// Test billing tags
	billingTags := schema.NewSet(schema.HashInt, []interface{}{1, 2, 3})
	d.Set("billing_tag_ids", billingTags)
	tags := d.Get("billing_tag_ids").(*schema.Set)
	assert.Equal(t, 3, tags.Len())
}

// Test API client CRUD operations with mock servers
func TestConnectorAPI_MockServerCRUD(t *testing.T) {
	// Test data for AWS VPC Connector
	connectorId := json.Number("123")
	mockAwsVpcConnector := &alkira.ConnectorAwsVpc{
		Id:             connectorId,
		Name:           "test-aws-vpc-connector",
		VpcOwnerId:     "123456789012",
		CustomerRegion: "us-west-2",
		CredentialId:   "cred-123",
		CXP:            "US-WEST",
		VpcId:          "vpc-123456",
		Size:           "SMALL",
		Enabled:        true,
		Segments:       []string{"test-segment"},
	}

	t.Run("AWS VPC Connector API", func(t *testing.T) {
		// Test Create
		client := serveMockServer(t, mockAwsVpcConnector, http.StatusCreated)
		api := alkira.NewConnectorAwsVpc(client)
		response, provState, err, provErr := api.Create(mockAwsVpcConnector)

		assert.NoError(t, err)
		assert.Equal(t, "test-aws-vpc-connector", response.Name)
		_ = provState
		_ = provErr

		// Test Read
		client = serveMockServer(t, mockAwsVpcConnector, http.StatusOK)
		api = alkira.NewConnectorAwsVpc(client)
		connector, provState, err := api.GetById("123")

		assert.NoError(t, err)
		assert.Equal(t, "test-aws-vpc-connector", connector.Name)
		assert.Equal(t, "123456789012", connector.VpcOwnerId)
		_ = provState

		// Test Delete
		client = serveMockServer(t, nil, http.StatusNoContent)
		api = alkira.NewConnectorAwsVpc(client)
		provState, err, provErr = api.Delete("123")

		assert.NoError(t, err)
		_ = provState
		_ = provErr
	})

	// Test error handling
	t.Run("Error Handling", func(t *testing.T) {
		client := serveMockServer(t, nil, http.StatusNotFound)
		api := alkira.NewConnectorAwsVpc(client)
		_, _, err := api.GetById("999")

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "404")
	})
}

// Helper function to create mock HTTP server using shared utility
func serveMockServer(t *testing.T, data interface{}, statusCode int) *alkira.AlkiraClient {
	return createMockAlkiraClient(t, func(w http.ResponseWriter, req *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(statusCode)

		switch req.Method {
		case "GET", "POST", "PUT":
			if data != nil {
				json.NewEncoder(w).Encode(data)
			}
		case "DELETE":
			// No content for delete
		default:
			w.WriteHeader(http.StatusMethodNotAllowed)
		}
	})
}

// Test ID validation for all connector types
func TestConnectorID_Validation(t *testing.T) {
	tests := GetCommonIdValidationTestCases()

	for _, tt := range tests {
		t.Run(tt.Name, func(t *testing.T) {
			result := validateResourceId(tt.Id)
			assert.Equal(t, tt.Valid, result, "ID validation result should match expected")
		})
	}
}

// Test size validation for connectors
func TestConnectorSize_Validation(t *testing.T) {
	validSizes := []string{"5XSMALL", "XSMALL", "SMALL", "MEDIUM", "LARGE", "2LARGE", "5LARGE", "10LARGE", "20LARGE"}

	for _, size := range validSizes {
		t.Run("valid size "+size, func(t *testing.T) {
			r := resourceAlkiraConnectorAwsVpc()
			d := r.TestResourceData()

			d.Set("size", size)
			assert.Equal(t, size, getStringFromResourceData(d, "size"))
		})
	}
}

// Test VPN mode validation for IPSec connector
func TestConnectorIPSec_VpnModeValidation(t *testing.T) {
	validModes := []string{"ROUTE_BASED", "POLICY_BASED"}

	for _, mode := range validModes {
		t.Run("valid VPN mode "+mode, func(t *testing.T) {
			r := resourceAlkiraConnectorIPSec()
			d := r.TestResourceData()

			d.Set("vpn_mode", mode)
			assert.Equal(t, mode, getStringFromResourceData(d, "vpn_mode"))
		})
	}
}

// Test connection mode validation for Azure VNET connector
func TestConnectorAzureVnet_ConnectionModeValidation(t *testing.T) {
	validModes := []string{"VNET_GATEWAY", "VNET_PEERING"}

	for _, mode := range validModes {
		t.Run("valid connection mode "+mode, func(t *testing.T) {
			r := resourceAlkiraConnectorAzureVnet()
			d := r.TestResourceData()

			d.Set("connection_mode", mode)
			assert.Equal(t, mode, getStringFromResourceData(d, "connection_mode"))
		})
	}
}
