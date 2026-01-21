package alkira

import (
	"encoding/json"
	"net/http"
	"testing"

	"github.com/alkiranet/alkira-client-go/alkira"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAlkiraServicePanResourceSchema(t *testing.T) {
	resource := resourceAlkiraServicePan()

	// Test required fields
	panUsernameSchema := resource.Schema["pan_username"]
	assert.True(t, panUsernameSchema.Required, "PAN username should be required")
	assert.Equal(t, schema.TypeString, panUsernameSchema.Type, "PAN username should be string type")

	panPasswordSchema := resource.Schema["pan_password"]
	assert.True(t, panPasswordSchema.Required, "PAN password should be required")
	assert.Equal(t, schema.TypeString, panPasswordSchema.Type, "PAN password should be string type")

	cxpSchema := resource.Schema["cxp"]
	assert.True(t, cxpSchema.Required, "CXP should be required")
	assert.Equal(t, schema.TypeString, cxpSchema.Type, "CXP should be string type")

	instanceSchema := resource.Schema["instance"]
	assert.True(t, instanceSchema.Required, "Instance should be required")
	assert.Equal(t, schema.TypeList, instanceSchema.Type, "Instance should be list type")

	// Test optional fields
	bundleSchema := resource.Schema["bundle"]
	assert.True(t, bundleSchema.Optional, "Bundle should be optional")
	assert.Equal(t, schema.TypeString, bundleSchema.Type, "Bundle should be string type")

	globalProtectEnabledSchema := resource.Schema["global_protect_enabled"]
	assert.True(t, globalProtectEnabledSchema.Optional, "Global protect enabled should be optional")
	assert.Equal(t, schema.TypeBool, globalProtectEnabledSchema.Type, "Global protect enabled should be bool type")
	assert.Equal(t, false, globalProtectEnabledSchema.Default, "Global protect enabled should default to false")

	billingTagIdsSchema := resource.Schema["billing_tag_ids"]
	assert.True(t, billingTagIdsSchema.Optional, "Billing tag IDs should be optional")
	assert.Equal(t, schema.TypeSet, billingTagIdsSchema.Type, "Billing tag IDs should be set type")

	// Test computed fields
	provStateSchema := resource.Schema["provision_state"]
	assert.True(t, provStateSchema.Computed, "Provision state should be computed")
	assert.Equal(t, schema.TypeString, provStateSchema.Type, "Provision state should be string type")

	panCredentialIdSchema := resource.Schema["pan_credential_id"]
	assert.True(t, panCredentialIdSchema.Computed, "PAN credential ID should be computed")
	assert.Equal(t, schema.TypeString, panCredentialIdSchema.Type, "PAN credential ID should be string type")

	// Test that resource has all required CRUD functions
	assert.NotNil(t, resource.CreateContext, "Resource should have CreateContext")
	assert.NotNil(t, resource.ReadContext, "Resource should have ReadContext")
	assert.NotNil(t, resource.UpdateContext, "Resource should have UpdateContext")
	assert.NotNil(t, resource.DeleteContext, "Resource should have DeleteContext")
	assert.NotNil(t, resource.Importer, "Resource should support import")
	assert.NotNil(t, resource.CustomizeDiff, "Resource should have CustomizeDiff")
}

func TestAlkiraServicePanValidateBundle(t *testing.T) {
	tests := []ValidationTestCase{
		{
			Name:      "Valid VM_SERIES_BUNDLE_1",
			Input:     "VM_SERIES_BUNDLE_1",
			ExpectErr: false,
			ErrCount:  0,
		},
		{
			Name:      "Valid VM_SERIES_BUNDLE_2",
			Input:     "VM_SERIES_BUNDLE_2",
			ExpectErr: false,
			ErrCount:  0,
		},
		{
			Name:      "Valid PAN_VM_300_BUNDLE_2",
			Input:     "PAN_VM_300_BUNDLE_2",
			ExpectErr: false,
			ErrCount:  0,
		},
		{
			Name:      "Invalid bundle",
			Input:     "INVALID_BUNDLE",
			ExpectErr: true,
			ErrCount:  1,
		},
		{
			Name:      "Empty string",
			Input:     "",
			ExpectErr: true,
			ErrCount:  1,
		},
		{
			Name:      "Non-string input",
			Input:     123,
			ExpectErr: true,
			ErrCount:  1,
		},
	}

	resource := resourceAlkiraServicePan()
	bundleSchema := resource.Schema["bundle"]

	for _, tt := range tests {
		t.Run(tt.Name, func(t *testing.T) {
			warnings, errors := bundleSchema.ValidateFunc(tt.Input, "bundle")

			if tt.ExpectErr {
				assert.Len(t, errors, tt.ErrCount, "Expected %d errors for input %v", tt.ErrCount, tt.Input)
			} else {
				assert.Empty(t, errors, "Expected no errors for input %v", tt.Input)
			}
			assert.Empty(t, warnings, "Expected no warnings")
		})
	}
}

func TestAlkiraServicePanExpandPanInstances(t *testing.T) {
	// Test the expandPanInstances helper function
	instances := []interface{}{
		map[string]interface{}{
			"name":      "pan-instance-1",
			"auth_key":  "test-auth-key",
			"auth_code": "test-auth-code",
		},
		map[string]interface{}{
			"name":      "pan-instance-2",
			"auth_key":  "test-auth-key-2",
			"auth_code": "test-auth-code-2",
		},
	}

	// Create a mock client since the function requires it
	mockClient := createMockAlkiraClient(t, func(w http.ResponseWriter, req *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	result, err := expandPanInstances(instances, mockClient)
	require.NoError(t, err, "expandPanInstances should not return error")
	require.Len(t, result, 2, "Should return 2 instances")

	// Check first instance
	assert.Equal(t, "pan-instance-1", result[0].Name)

	// Check second instance
	assert.Equal(t, "pan-instance-2", result[1].Name)
}

func TestAlkiraServicePanSetPanInstances(t *testing.T) {
	// Test the setPanInstances helper function
	instances := []alkira.ServicePanInstance{
		{
			Id:           1,
			Name:         "pan-instance-1",
			CredentialId: "cred-1",
		},
		{
			Id:           2,
			Name:         "pan-instance-2",
			CredentialId: "cred-2",
		},
	}

	r := resourceAlkiraServicePan()
	d := r.TestResourceData()

	// Create a mock client since the function requires it
	mockClient := createMockAlkiraClient(t, func(w http.ResponseWriter, req *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	result := setPanInstances(d, instances, mockClient)
	require.Len(t, result, 2, "Should return 2 instances")

	// Check first instance
	assert.Equal(t, "pan-instance-1", result[0]["name"])
	assert.Equal(t, 1, result[0]["id"])
	assert.Equal(t, "cred-1", result[0]["credential_id"])

	// Check second instance
	assert.Equal(t, "pan-instance-2", result[1]["name"])
	assert.Equal(t, 2, result[1]["id"])
	assert.Equal(t, "cred-2", result[1]["credential_id"])
}

func TestAlkiraServicePanCreatePanCredential(t *testing.T) {
	// Create mock client using shared utility
	client := createMockAlkiraClient(t, func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		w.Write([]byte(`{"id": "credential-123"}`))
	})

	// Test data
	r := resourceAlkiraServicePan()
	d := r.TestResourceData()
	d.Set("name", "test-pan-service")
	d.Set("pan_username", "admin")
	d.Set("pan_password", "test-password")
	d.Set("pan_license_key", "test-license-key")

	credentialId, err := createPanCredential(d, client)
	require.NoError(t, err, "createPanCredential should not return error")
	assert.NotEmpty(t, credentialId, "Credential ID should not be empty")
}

func TestAlkiraServicePanFlattenGlobalProtectSegmentOptions(t *testing.T) {
	tests := []struct {
		name     string
		input    map[string]*alkira.GlobalProtectSegmentName
		expected int // expected number of results
	}{
		{
			name:     "Nil input",
			input:    nil,
			expected: 0,
		},
		{
			name:     "Empty map",
			input:    map[string]*alkira.GlobalProtectSegmentName{},
			expected: 0,
		},
		{
			name: "Single segment option",
			input: map[string]*alkira.GlobalProtectSegmentName{
				"test-segment": {
					RemoteUserZoneName: "remote-zone-1",
					PortalFqdnPrefix:   "portal-prefix-1",
					ServiceGroupName:   "service-group-1",
				},
			},
			expected: 1,
		},
		{
			name: "Multiple segment options",
			input: map[string]*alkira.GlobalProtectSegmentName{
				"segment-1": {
					RemoteUserZoneName: "remote-zone-1",
					PortalFqdnPrefix:   "portal-prefix-1",
					ServiceGroupName:   "service-group-1",
				},
				"segment-2": {
					RemoteUserZoneName: "remote-zone-2",
					PortalFqdnPrefix:   "portal-prefix-2",
					ServiceGroupName:   "service-group-2",
				},
			},
			expected: 2,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create a mock client that returns segment IDs
			mockClient := createMockAlkiraClient(t, func(w http.ResponseWriter, req *http.Request) {
				w.Header().Set("Content-Type", "application/json")
				// Return a segment with the ID based on the requested name
				response := []alkira.Segment{
					{
						Id:   "123",
						Name: "test-segment",
					},
				}
				// Check if it's segment-1 or segment-2
				if req.URL.Query().Get("name") == "segment-1" {
					response[0].Id = "101"
					response[0].Name = "segment-1"
				} else if req.URL.Query().Get("name") == "segment-2" {
					response[0].Id = "102"
					response[0].Name = "segment-2"
				}
				w.WriteHeader(http.StatusOK)
				jsonBytes, _ := json.Marshal(response)
				w.Write(jsonBytes)
			})

			result := flattenGlobalProtectSegmentOptions(tt.input, mockClient)

			if tt.expected == 0 {
				assert.Nil(t, result, "Expected nil result for empty/nil input")
			} else {
				require.Len(t, result, tt.expected, "Expected %d results", tt.expected)

				// Verify the flattened structure has all required fields
				for _, opt := range result {
					assert.Contains(t, opt, "segment_id", "Result should contain segment_id")
					assert.Contains(t, opt, "remote_user_zone_name", "Result should contain remote_user_zone_name")
					assert.Contains(t, opt, "portal_fqdn_prefix", "Result should contain portal_fqdn_prefix")
					assert.Contains(t, opt, "service_group_name", "Result should contain service_group_name")
				}
			}
		})
	}
}

func TestAlkiraServicePanFlattenGlobalProtectSegmentOptionsValues(t *testing.T) {
	// Test that the flattened values match the input
	input := map[string]*alkira.GlobalProtectSegmentName{
		"test-segment": {
			RemoteUserZoneName: "my-remote-zone",
			PortalFqdnPrefix:   "my-portal-prefix",
			ServiceGroupName:   "my-service-group",
		},
	}

	mockClient := createMockAlkiraClient(t, func(w http.ResponseWriter, req *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		response := []alkira.Segment{
			{
				Id:   "999",
				Name: "test-segment",
			},
		}
		w.WriteHeader(http.StatusOK)
		jsonBytes, _ := json.Marshal(response)
		w.Write(jsonBytes)
	})

	result := flattenGlobalProtectSegmentOptions(input, mockClient)

	require.Len(t, result, 1, "Should return 1 option")
	assert.Equal(t, "999", result[0]["segment_id"], "segment_id should be 999")
	assert.Equal(t, "my-remote-zone", result[0]["remote_user_zone_name"], "remote_user_zone_name should match")
	assert.Equal(t, "my-portal-prefix", result[0]["portal_fqdn_prefix"], "portal_fqdn_prefix should match")
	assert.Equal(t, "my-service-group", result[0]["service_group_name"], "service_group_name should match")
}

func TestAlkiraServicePanFlattenGlobalProtectSegmentOptionsInstance(t *testing.T) {
	tests := []struct {
		name     string
		input    map[string]*alkira.GlobalProtectSegmentNameInstance
		expected int // expected number of results
	}{
		{
			name:     "Nil input",
			input:    nil,
			expected: 0,
		},
		{
			name:     "Empty map",
			input:    map[string]*alkira.GlobalProtectSegmentNameInstance{},
			expected: 0,
		},
		{
			name: "Single instance option",
			input: map[string]*alkira.GlobalProtectSegmentNameInstance{
				"test-segment": {
					PortalEnabled:  true,
					GatewayEnabled: false,
					PrefixListId:   100,
				},
			},
			expected: 1,
		},
		{
			name: "Multiple instance options",
			input: map[string]*alkira.GlobalProtectSegmentNameInstance{
				"segment-1": {
					PortalEnabled:  true,
					GatewayEnabled: true,
					PrefixListId:   101,
				},
				"segment-2": {
					PortalEnabled:  false,
					GatewayEnabled: true,
					PrefixListId:   102,
				},
			},
			expected: 2,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockClient := createMockAlkiraClient(t, func(w http.ResponseWriter, req *http.Request) {
				w.Header().Set("Content-Type", "application/json")
				response := []alkira.Segment{
					{
						Id:   "123",
						Name: "test-segment",
					},
				}
				if req.URL.Query().Get("name") == "segment-1" {
					response[0].Id = "201"
					response[0].Name = "segment-1"
				} else if req.URL.Query().Get("name") == "segment-2" {
					response[0].Id = "202"
					response[0].Name = "segment-2"
				}
				w.WriteHeader(http.StatusOK)
				jsonBytes, _ := json.Marshal(response)
				w.Write(jsonBytes)
			})

			result := flattenGlobalProtectSegmentOptionsInstance(tt.input, mockClient)

			if tt.expected == 0 {
				assert.Nil(t, result, "Expected nil result for empty/nil input")
			} else {
				require.Len(t, result, tt.expected, "Expected %d results", tt.expected)

				// Verify the flattened structure has all required fields
				for _, opt := range result {
					assert.Contains(t, opt, "segment_id", "Result should contain segment_id")
					assert.Contains(t, opt, "portal_enabled", "Result should contain portal_enabled")
					assert.Contains(t, opt, "gateway_enabled", "Result should contain gateway_enabled")
					assert.Contains(t, opt, "prefix_list_id", "Result should contain prefix_list_id")
				}
			}
		})
	}
}

func TestAlkiraServicePanFlattenGlobalProtectSegmentOptionsInstanceValues(t *testing.T) {
	// Test that the flattened values match the input
	input := map[string]*alkira.GlobalProtectSegmentNameInstance{
		"test-segment": {
			PortalEnabled:  true,
			GatewayEnabled: false,
			PrefixListId:   555,
		},
	}

	mockClient := createMockAlkiraClient(t, func(w http.ResponseWriter, req *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		response := []alkira.Segment{
			{
				Id:   "888",
				Name: "test-segment",
			},
		}
		w.WriteHeader(http.StatusOK)
		jsonBytes, _ := json.Marshal(response)
		w.Write(jsonBytes)
	})

	result := flattenGlobalProtectSegmentOptionsInstance(input, mockClient)

	require.Len(t, result, 1, "Should return 1 option")
	assert.Equal(t, "888", result[0]["segment_id"], "segment_id should be 888")
	assert.Equal(t, true, result[0]["portal_enabled"], "portal_enabled should be true")
	assert.Equal(t, false, result[0]["gateway_enabled"], "gateway_enabled should be false")
	assert.Equal(t, 555, result[0]["prefix_list_id"], "prefix_list_id should be 555")
}

// UNUSED: Commented out to suppress linter warnings
// // TEST HELPER
// func serveServicePan(t *testing.T, servicePan *alkira.ServicePan) *alkira.AlkiraClient {
// 	return createMockAlkiraClient(t, func(w http.ResponseWriter, req *http.Request) {
// 		json.NewEncoder(w).Encode(servicePan)
// 		w.Header().Set("Content-Type", "application/json")
// 	})
// }
