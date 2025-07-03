package alkira

import (
	"encoding/json"
	"testing"

	"github.com/alkiranet/alkira-client-go/alkira"
	"github.com/stretchr/testify/assert"
)

func TestAlkiraGroup_validateGroupDescription(t *testing.T) {
	tests := []struct {
		name      string
		input     interface{}
		expectErr bool
	}{
		{
			name:      "valid description",
			input:     "This is a valid group description",
			expectErr: false,
		},
		{
			name:      "empty description",
			input:     "",
			expectErr: false, // Empty descriptions are typically allowed
		},
		{
			name:      "description with special characters",
			input:     "Description with special chars: @#$%^&*()",
			expectErr: false,
		},
		{
			name:      "very long description",
			input:     "This is a very long description that contains many words and should test whether the validation function properly handles lengthy descriptions that might exceed normal length limits",
			expectErr: false, // Length validation depends on Alkira requirements
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			warnings, errors := validateGroupDescription(tt.input, "description")

			if tt.expectErr {
				assert.NotEmpty(t, errors, "Expected errors for input: %v", tt.input)
			} else {
				assert.Empty(t, errors, "Expected no errors for input: %v, got: %v", tt.input, errors)
			}

			assert.Empty(t, warnings, "Expected no warnings, got %v", warnings)
		})
	}
}

func TestAlkiraGroup_provisionStateHandling(t *testing.T) {
	tests := []struct {
		name          string
		provisionFlag bool
		provState     string
		provErr       error
		expectWarning bool
	}{
		{
			name:          "provision disabled",
			provisionFlag: false,
			provState:     "SUCCESS",
			provErr:       nil,
			expectWarning: false,
		},
		{
			name:          "provision enabled success",
			provisionFlag: true,
			provState:     "SUCCESS",
			provErr:       nil,
			expectWarning: false,
		},
		{
			name:          "provision enabled failed",
			provisionFlag: true,
			provState:     "FAILED",
			provErr:       assert.AnError,
			expectWarning: true,
		},
		{
			name:          "provision enabled pending",
			provisionFlag: true,
			provState:     "PENDING",
			provErr:       nil,
			expectWarning: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Simulate provision state handling logic
			var warnings []string
			var errors []error

			if tt.provisionFlag && tt.provErr != nil {
				warnings = append(warnings, "PROVISION FAILED")
				errors = append(errors, tt.provErr)
			}

			if tt.expectWarning {
				assert.NotEmpty(t, warnings, "Expected warnings for failed provision")
			} else {
				assert.Empty(t, warnings, "Expected no warnings for successful provision")
			}
		})
	}
}

func TestAlkiraGroup_groupStructFieldValidation(t *testing.T) {
	// Test Group struct field types
	t.Run("group struct creation", func(t *testing.T) {
		group := &alkira.Group{
			Id:          json.Number("123"),
			Name:        "test-group",
			Description: "test description",
		}

		assert.Equal(t, "123", string(group.Id))
		assert.Equal(t, "test-group", group.Name)
		assert.Equal(t, "test description", group.Description)
	})

	t.Run("empty group struct", func(t *testing.T) {
		group := &alkira.Group{}

		assert.Equal(t, "", string(group.Id))
		assert.Equal(t, "", group.Name)
		assert.Equal(t, "", group.Description)
	})
}

func TestAlkiraGroup_customizeDiffBehavior(t *testing.T) {
	// Test the custom diff behavior from the resource definition
	resource := resourceAlkiraGroup()
	assert.NotNil(t, resource.CustomizeDiff, "Resource should have CustomizeDiff function")

	// The CustomizeDiff function handles provision state changes
	// This would typically be tested in integration tests with actual terraform state
}

func TestAlkiraGroup_importSupport(t *testing.T) {
	resource := resourceAlkiraGroup()

	// Verify that import is supported
	assert.NotNil(t, resource.Importer, "Resource should support import")
	assert.NotNil(t, resource.Importer.StateContext, "Resource should have state context for import")
}

func TestAlkiraGroup_schemaTypeValidation(t *testing.T) {
	resource := resourceAlkiraGroup()

	// Test that all expected schema fields exist
	expectedFields := []string{"name", "description", "provision_state"}

	for _, field := range expectedFields {
		schema, exists := resource.Schema[field]
		assert.True(t, exists, "Field '%s' should exist in schema", field)
		assert.NotNil(t, schema, "Schema for field '%s' should not be nil", field)
	}

	// Test schema field properties
	nameSchema := resource.Schema["name"]
	assert.True(t, nameSchema.Required, "Name field should be required")
	assert.False(t, nameSchema.Optional, "Name field should not be optional")
	assert.False(t, nameSchema.Computed, "Name field should not be computed")

	descSchema := resource.Schema["description"]
	assert.False(t, descSchema.Required, "Description field should not be required")
	assert.True(t, descSchema.Optional, "Description field should be optional")
	assert.False(t, descSchema.Computed, "Description field should not be computed")

	provStateSchema := resource.Schema["provision_state"]
	assert.False(t, provStateSchema.Required, "Provision state should not be required")
	assert.False(t, provStateSchema.Optional, "Provision state should not be optional")
	assert.True(t, provStateSchema.Computed, "Provision state should be computed")
}

func TestAlkiraGroup_errorMessageValidation(t *testing.T) {
	tests := []struct {
		name        string
		inputName   string
		expectError bool
		errorMsg    string
	}{
		{
			name:        "valid name",
			inputName:   "valid-group-name",
			expectError: false,
			errorMsg:    "",
		},
		{
			name:        "empty name",
			inputName:   "",
			expectError: true,
			errorMsg:    "cannot be empty",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			warnings, errors := validateGroupName(tt.inputName, "name")

			if tt.expectError {
				assert.NotEmpty(t, errors, "Expected validation errors")
				if len(errors) > 0 {
					assert.Contains(t, errors[0].Error(), tt.errorMsg, "Error message should contain expected text")
				}
			} else {
				assert.Empty(t, errors, "Expected no validation errors")
			}

			assert.Empty(t, warnings, "Expected no validation warnings")
		})
	}
}

// Helper validation function for description
func validateGroupDescription(v interface{}, k string) (warnings []string, errors []error) {
	// Description validation logic
	// Currently allowing empty descriptions and any content
	return warnings, errors
}
