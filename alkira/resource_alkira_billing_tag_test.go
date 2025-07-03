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

func TestAlkiraBillingTag_buildBillingTagRequest(t *testing.T) {
	r := resourceAlkiraBillingTag()
	d := r.TestResourceData()

	// Test with complete billing tag data
	expectedName := "test-billing-tag"
	expectedDescription := "Test billing tag description"

	d.Set("name", expectedName)
	d.Set("description", expectedDescription)

	request := buildBillingTagRequest(d)

	require.Equal(t, expectedName, request.Name)
	require.Equal(t, expectedDescription, request.Description)
}

func TestAlkiraBillingTag_buildBillingTagRequestMinimal(t *testing.T) {
	r := resourceAlkiraBillingTag()
	d := r.TestResourceData()

	// Test with minimal billing tag data
	expectedName := "minimal-billing-tag"

	d.Set("name", expectedName)
	// description not set (should be empty)

	request := buildBillingTagRequest(d)

	require.Equal(t, expectedName, request.Name)
	require.Equal(t, "", request.Description)
}

func TestAlkiraBillingTag_resourceSchema(t *testing.T) {
	resource := resourceAlkiraBillingTag()

	// Test required fields
	nameSchema := resource.Schema["name"]
	assert.True(t, nameSchema.Required, "Name should be required")
	assert.Equal(t, schema.TypeString, nameSchema.Type, "Name should be string type")

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

func TestAlkiraBillingTag_validateBillingTagName(t *testing.T) {
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
func TestAlkiraBillingTag_apiClientCRUD(t *testing.T) {
	// Test data
	billingTagId := json.Number("123")
	billingTagName := "test-billing-tag"
	billingTagDescription := "Test billing tag description"

	// Create mock billing tag
	mockBillingTag := &alkira.BillingTag{
		Id:          billingTagId,
		Name:        billingTagName,
		Description: billingTagDescription,
	}

	// Test Create operation
	t.Run("Create", func(t *testing.T) {
		client := serveBillingTagMockServer(t, mockBillingTag, http.StatusCreated)

		api := alkira.NewBillingTag(client)
		response, provState, err, provErr := api.Create(mockBillingTag)

		assert.NoError(t, err)
		assert.Equal(t, billingTagName, response.Name)
		assert.Equal(t, billingTagDescription, response.Description)
		_ = provErr

		t.Logf("Provision state: %s", provState)
	})

	// Test Read operation
	t.Run("Read", func(t *testing.T) {
		client := serveBillingTagMockServer(t, mockBillingTag, http.StatusOK)

		api := alkira.NewBillingTag(client)
		billingTag, provState, err := api.GetById("123")

		assert.NoError(t, err)
		assert.Equal(t, billingTagName, billingTag.Name)
		assert.Equal(t, billingTagDescription, billingTag.Description)

		t.Logf("Provision state: %s", provState)
	})

	// Test Update operation
	t.Run("Update", func(t *testing.T) {
		updatedBillingTag := &alkira.BillingTag{
			Id:          billingTagId,
			Name:        billingTagName + "-updated",
			Description: billingTagDescription + " updated",
		}

		client := serveBillingTagMockServer(t, updatedBillingTag, http.StatusOK)

		api := alkira.NewBillingTag(client)
		provState, err, provErr := api.Update("123", updatedBillingTag)

		assert.NoError(t, err)
		_ = provErr

		t.Logf("Provision state: %s", provState)
	})

	// Test Delete operation
	t.Run("Delete", func(t *testing.T) {
		client := serveBillingTagMockServer(t, nil, http.StatusNoContent)

		api := alkira.NewBillingTag(client)
		provState, err, provErr := api.Delete("123")

		assert.NoError(t, err)
		_ = provErr

		t.Logf("Provision state: %s", provState)
	})
}

func TestAlkiraBillingTag_apiErrorHandling(t *testing.T) {
	// Test error scenarios
	t.Run("NotFound", func(t *testing.T) {
		client := serveBillingTagMockServer(t, nil, http.StatusNotFound)

		api := alkira.NewBillingTag(client)
		_, _, err := api.GetById("999")

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "404")
	})

	t.Run("ServerError", func(t *testing.T) {
		client := serveBillingTagMockServer(t, nil, http.StatusInternalServerError)

		api := alkira.NewBillingTag(client)
		_, _, _, _ = api.Create(&alkira.BillingTag{
			Name:        "test-billing-tag",
			Description: "test",
		})

		// Should handle server errors gracefully
	})
}

func TestAlkiraBillingTag_resourceDataManipulation(t *testing.T) {
	r := resourceAlkiraBillingTag()

	t.Run("set and get resource data", func(t *testing.T) {
		d := r.TestResourceData()

		// Set values
		d.Set("name", "test-billing-tag")
		d.Set("description", "Test description")

		// Test getting values
		assert.Equal(t, "test-billing-tag", d.Get("name").(string))
		assert.Equal(t, "Test description", d.Get("description").(string))
	})

	t.Run("resource data with changes", func(t *testing.T) {
		d := r.TestResourceData()

		// Set initial values
		d.Set("name", "original-name")
		d.Set("description", "Original description")

		// Simulate a change
		d.Set("name", "updated-name")
		d.Set("description", "Updated description")

		assert.Equal(t, "updated-name", d.Get("name").(string))
		assert.Equal(t, "Updated description", d.Get("description").(string))
	})
}

func TestAlkiraBillingTag_idValidation(t *testing.T) {
	tests := GetCommonIdValidationTestCases()

	for _, tt := range tests {
		t.Run(tt.Name, func(t *testing.T) {
			isValid := validateResourceId(tt.Id)
			assert.Equal(t, tt.Valid, isValid, "Expected ID validation to return %v for %s", tt.Valid, tt.Id)
		})
	}
}

// Helper function to create mock HTTP server for billing tags
func serveBillingTagMockServer(t *testing.T, billingTag *alkira.BillingTag, statusCode int) *alkira.AlkiraClient {
	handler := http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(statusCode)

		switch req.Method {
		case "GET":
			if billingTag != nil {
				json.NewEncoder(w).Encode(billingTag)
			}
		case "POST":
			if billingTag != nil {
				json.NewEncoder(w).Encode(billingTag)
			}
		case "PUT":
			if billingTag != nil {
				json.NewEncoder(w).Encode(billingTag)
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
func buildBillingTagRequest(d *schema.ResourceData) *alkira.BillingTag {
	return &alkira.BillingTag{
		Name:        getStringFromResourceData(d, "name"),
		Description: getStringFromResourceData(d, "description"),
	}
}
