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

func TestAlkiraServicePan_resourceSchema(t *testing.T) {
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

func TestAlkiraServicePan_validateBundle(t *testing.T) {
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

func TestAlkiraServicePan_expandPanInstances(t *testing.T) {
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

func TestAlkiraServicePan_setPanInstances(t *testing.T) {
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

	result := setPanInstances(d, instances)
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

func TestAlkiraServicePan_createPanCredential(t *testing.T) {
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

// TEST HELPER
func serveServicePan(t *testing.T, servicePan *alkira.ServicePan) *alkira.AlkiraClient {
	return createMockAlkiraClient(t, func(w http.ResponseWriter, req *http.Request) {
		json.NewEncoder(w).Encode(servicePan)
		w.Header().Set("Content-Type", "application/json")
	})
}
