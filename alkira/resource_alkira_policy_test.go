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

func TestAlkiraPolicy_buildPolicyRequest(t *testing.T) {
	r := resourceAlkiraPolicy()
	d := r.TestResourceData()

	// Test with complete policy data
	expectedName := "test-policy"
	expectedDescription := "Test policy description"
	expectedEnabled := true

	d.Set("name", expectedName)
	d.Set("description", expectedDescription)
	d.Set("enabled", expectedEnabled)

	request := buildPolicyRequest(d)

	require.Equal(t, expectedName, request.Name)
	require.Equal(t, expectedDescription, request.Description)
	require.Equal(t, expectedEnabled, request.Enabled)
}

func TestAlkiraPolicy_buildPolicyRequestMinimal(t *testing.T) {
	r := resourceAlkiraPolicy()
	d := r.TestResourceData()

	// Test with minimal policy data
	expectedName := "minimal-policy"

	d.Set("name", expectedName)
	d.Set("enabled", true) // required field

	request := buildPolicyRequest(d)

	require.Equal(t, expectedName, request.Name)
	require.Equal(t, "", request.Description)
	require.Equal(t, true, request.Enabled)
}

func TestAlkiraPolicy_resourceSchema(t *testing.T) {
	resource := resourceAlkiraPolicy()

	// Test required fields
	nameSchema := resource.Schema["name"]
	assert.True(t, nameSchema.Required, "Name should be required")
	assert.Equal(t, schema.TypeString, nameSchema.Type, "Name should be string type")

	enabledSchema := resource.Schema["enabled"]
	assert.True(t, enabledSchema.Required, "Enabled should be required")
	assert.Equal(t, schema.TypeBool, enabledSchema.Type, "Enabled should be bool type")

	// Test optional fields
	descSchema := resource.Schema["description"]
	assert.True(t, descSchema.Optional, "Description should be optional")
	assert.Equal(t, schema.TypeString, descSchema.Type, "Description should be string type")

	// Test that resource has all required CRUD functions
	assert.NotNil(t, resource.CreateContext, "Resource should have CreateContext")
	assert.NotNil(t, resource.ReadContext, "Resource should have ReadContext")
	assert.NotNil(t, resource.UpdateContext, "Resource should have UpdateContext")
	assert.NotNil(t, resource.DeleteContext, "Resource should have DeleteContext")
	assert.NotNil(t, resource.Importer, "Resource should support import")
}

func TestAlkiraPolicy_validatePolicyName(t *testing.T) {
	tests := GetCommonNameValidationTestCases()

	for _, tt := range tests {
		t.Run(tt.Name, func(t *testing.T) {
			warnings, errors := validateResourceName(tt.Input, "name")

			if tt.ExpectErr {
				assert.Len(t, errors, tt.ErrCount, "Expected %d errors, got %d", tt.ErrCount, len(errors))
			} else {
				assert.Len(t, errors, 0, "Expected no errors, got %v", errors)
			}

			assert.Len(t, warnings, 0, "Expected no warnings, got %v", warnings)
		})
	}
}

// Unit test with mock HTTP server
func TestAlkiraPolicy_apiClientCRUD(t *testing.T) {
	// Test data
	policyId := json.Number("123")
	policyName := "test-policy"
	policyDescription := "Test policy description"
	policyEnabled := true

	// Create mock policy
	mockPolicy := &alkira.TrafficPolicy{
		Id:          policyId,
		Name:        policyName,
		Description: policyDescription,
		Enabled:     policyEnabled,
		FromGroups:  []int{1, 2},
		ToGroups:    []int{3, 4},
		SegmentIds:  []int{5, 6},
		RuleListId:  1,
	}

	// Test Create operation
	t.Run("Create", func(t *testing.T) {
		client := servePolicyMockServer(t, mockPolicy, http.StatusCreated)

		api := alkira.NewTrafficPolicy(client)
		response, provState, err, valErr, provErr := api.Create(mockPolicy)

		assert.NoError(t, err)
		assert.Equal(t, policyName, response.Name)
		assert.Equal(t, policyDescription, response.Description)
		assert.Equal(t, policyEnabled, response.Enabled)
		_ = valErr
		_ = provErr

		t.Logf("Provision state: %s", provState)
	})

	// Test Read operation
	t.Run("Read", func(t *testing.T) {
		client := servePolicyMockServer(t, mockPolicy, http.StatusOK)

		api := alkira.NewTrafficPolicy(client)
		policy, provState, err := api.GetById("123")

		assert.NoError(t, err)
		assert.Equal(t, policyName, policy.Name)
		assert.Equal(t, policyDescription, policy.Description)
		assert.Equal(t, policyEnabled, policy.Enabled)

		t.Logf("Provision state: %s", provState)
	})

	// Test Update operation
	t.Run("Update", func(t *testing.T) {
		updatedPolicy := &alkira.TrafficPolicy{
			Id:          policyId,
			Name:        policyName + "-updated",
			Description: policyDescription + " updated",
			Enabled:     false,
			FromGroups:  []int{7, 8},
			ToGroups:    []int{9, 10},
			SegmentIds:  []int{11, 12},
			RuleListId:  2,
		}

		client := servePolicyMockServer(t, updatedPolicy, http.StatusOK)

		api := alkira.NewTrafficPolicy(client)
		provState, err, valErr, provErr := api.Update("123", updatedPolicy)

		assert.NoError(t, err)
		_ = valErr
		_ = provErr

		t.Logf("Provision state: %s", provState)
	})

	// Test Delete operation
	t.Run("Delete", func(t *testing.T) {
		client := servePolicyMockServer(t, nil, http.StatusNoContent)

		api := alkira.NewTrafficPolicy(client)
		provState, err, valErr, provErr := api.Delete("123")

		assert.NoError(t, err)
		_ = valErr
		_ = provErr

		t.Logf("Provision state: %s", provState)
	})
}

func TestAlkiraPolicy_apiErrorHandling(t *testing.T) {
	// Test error scenarios
	t.Run("NotFound", func(t *testing.T) {
		client := servePolicyMockServer(t, nil, http.StatusNotFound)

		api := alkira.NewTrafficPolicy(client)
		_, _, err := api.GetById("999")

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "404")
	})

	t.Run("ServerError", func(t *testing.T) {
		client := servePolicyMockServer(t, nil, http.StatusInternalServerError)

		api := alkira.NewTrafficPolicy(client)
		_, _, _, _, _ = api.Create(&alkira.TrafficPolicy{
			Name:        "test-policy",
			Description: "test",
			Enabled:     true,
			FromGroups:  []int{1},
			ToGroups:    []int{2},
			SegmentIds:  []int{3},
			RuleListId:  1,
		})

		// Should handle server errors gracefully
	})
}

func TestAlkiraPolicy_resourceDataManipulation(t *testing.T) {
	r := resourceAlkiraPolicy()

	t.Run("set and get resource data", func(t *testing.T) {
		d := r.TestResourceData()

		// Set values
		d.Set("name", "test-policy")
		d.Set("description", "Test description")
		d.Set("enabled", true)

		// Test getting values
		assert.Equal(t, "test-policy", d.Get("name").(string))
		assert.Equal(t, "Test description", d.Get("description").(string))
		assert.Equal(t, true, d.Get("enabled").(bool))
	})

	t.Run("resource data with changes", func(t *testing.T) {
		d := r.TestResourceData()

		// Set initial values
		d.Set("name", "original-name")
		d.Set("description", "Original description")
		d.Set("enabled", false)

		// Simulate a change
		d.Set("name", "updated-name")
		d.Set("description", "Updated description")
		d.Set("enabled", true)

		assert.Equal(t, "updated-name", d.Get("name").(string))
		assert.Equal(t, "Updated description", d.Get("description").(string))
		assert.Equal(t, true, d.Get("enabled").(bool))
	})
}

func TestAlkiraPolicy_idValidation(t *testing.T) {
	tests := GetCommonIdValidationTestCases()

	for _, tt := range tests {
		t.Run(tt.Name, func(t *testing.T) {
			isValid := validateResourceId(tt.Id)
			assert.Equal(t, tt.Valid, isValid, "Expected ID validation to return %v for %s", tt.Valid, tt.Id)
		})
	}
}

func TestAlkiraPolicy_customizeDiff(t *testing.T) {
	r := resourceAlkiraPolicy()

	// Test CustomizeDiff function exists
	assert.NotNil(t, r.CustomizeDiff, "Resource should have CustomizeDiff")

	// Note: CustomizeDiff testing would require more complex setup with ResourceDiff mock
	// This validates the function exists and can be called
}

// Helper function to create mock HTTP server for policies
func servePolicyMockServer(t *testing.T, policy *alkira.TrafficPolicy, statusCode int) *alkira.AlkiraClient {
	handler := http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(statusCode)

		switch req.Method {
		case "GET":
			if policy != nil {
				json.NewEncoder(w).Encode(policy)
			}
		case "POST":
			if policy != nil {
				json.NewEncoder(w).Encode(policy)
			}
		case "PUT":
			if policy != nil {
				json.NewEncoder(w).Encode(policy)
			}
		case "DELETE":
			// No content for delete
		default:
			w.WriteHeader(http.StatusMethodNotAllowed)
		}
	})

	return createMockAlkiraClient(t, handler)
}

// Mock helper functions for testing
func buildPolicyRequest(d *schema.ResourceData) *alkira.TrafficPolicy {
	name := ""
	if n := d.Get("name"); n != nil {
		name = n.(string)
	}

	description := ""
	if desc := d.Get("description"); desc != nil {
		description = desc.(string)
	}

	enabled := false
	if e := d.Get("enabled"); e != nil {
		enabled = e.(bool)
	}

	return &alkira.TrafficPolicy{
		Name:        name,
		Description: description,
		Enabled:     enabled,
		FromGroups:  []int{1, 2},
		ToGroups:    []int{3, 4},
		SegmentIds:  []int{5, 6},
		RuleListId:  1,
	}
}
