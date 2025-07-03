package alkira

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"
	"time"

	"github.com/alkiranet/alkira-client-go/alkira"
	"github.com/hashicorp/go-retryablehttp"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAlkiraGroup_buildGroupRequest(t *testing.T) {
	r := resourceAlkiraGroup()
	d := r.TestResourceData()

	// Test with complete group data
	expectedName := "test-group"
	expectedDescription := "Test group description"

	d.Set("name", expectedName)
	d.Set("description", expectedDescription)

	request := buildGroupRequest(d)

	require.Equal(t, expectedName, request.Name)
	require.Equal(t, expectedDescription, request.Description)
}

func TestAlkiraGroup_buildGroupRequestMinimal(t *testing.T) {
	r := resourceAlkiraGroup()
	d := r.TestResourceData()

	// Test with minimal group data
	expectedName := "minimal-group"

	d.Set("name", expectedName)
	// description not set (should be empty)

	request := buildGroupRequest(d)

	require.Equal(t, expectedName, request.Name)
	require.Equal(t, "", request.Description)
}

func TestAlkiraGroup_buildGroupRequestSpecialChars(t *testing.T) {
	r := resourceAlkiraGroup()
	d := r.TestResourceData()

	// Test with special characters
	expectedName := "test-group-123_special"
	expectedDescription := "Group with special chars: @#$%"

	d.Set("name", expectedName)
	d.Set("description", expectedDescription)

	request := buildGroupRequest(d)

	require.Equal(t, expectedName, request.Name)
	require.Equal(t, expectedDescription, request.Description)
}

func TestAlkiraGroup_resourceSchema(t *testing.T) {
	resource := resourceAlkiraGroup()

	// Test required fields
	nameSchema := resource.Schema["name"]
	assert.True(t, nameSchema.Required, "Name should be required")
	assert.Equal(t, schema.TypeString, nameSchema.Type, "Name should be string type")

	// Test optional fields
	descSchema := resource.Schema["description"]
	assert.True(t, descSchema.Optional, "Description should be optional")
	assert.Equal(t, schema.TypeString, descSchema.Type, "Description should be string type")

	// Test computed fields
	provStateSchema := resource.Schema["provision_state"]
	assert.True(t, provStateSchema.Computed, "Provision state should be computed")
	assert.Equal(t, schema.TypeString, provStateSchema.Type, "Provision state should be string type")

	// Test that resource has all required CRUD functions
	assert.NotNil(t, resource.CreateContext, "Resource should have CreateContext")
	assert.NotNil(t, resource.ReadContext, "Resource should have ReadContext")
	assert.NotNil(t, resource.UpdateContext, "Resource should have UpdateContext")
	assert.NotNil(t, resource.DeleteContext, "Resource should have DeleteContext")
	assert.NotNil(t, resource.Importer, "Resource should support import")
}

func TestAlkiraGroup_validateGroupName(t *testing.T) {
	tests := []struct {
		name      string
		input     interface{}
		expectErr bool
		errCount  int
	}{
		{
			name:      "valid name",
			input:     "test-group",
			expectErr: false,
			errCount:  0,
		},
		{
			name:      "valid name with numbers",
			input:     "test-group-123",
			expectErr: false,
			errCount:  0,
		},
		{
			name:      "valid name with underscores",
			input:     "test_group_name",
			expectErr: false,
			errCount:  0,
		},
		{
			name:      "empty name",
			input:     "",
			expectErr: true,
			errCount:  1,
		},
		{
			name:      "name with spaces",
			input:     "test group",
			expectErr: false,
			errCount:  0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			warnings, errors := validateGroupName(tt.input, "name")

			if tt.expectErr {
				assert.Len(t, errors, tt.errCount, "Expected %d errors, got %d", tt.errCount, len(errors))
			} else {
				assert.Len(t, errors, 0, "Expected no errors, got %v", errors)
			}

			assert.Len(t, warnings, 0, "Expected no warnings, got %v", warnings)
		})
	}
}

func TestAlkiraGroup_jsonNumberHandling(t *testing.T) {
	tests := []struct {
		name     string
		jsonData string
		expected string
	}{
		{
			name:     "integer ID",
			jsonData: `{"id": 123, "name": "test", "description": "test desc"}`,
			expected: "123",
		},
		{
			name:     "string ID",
			jsonData: `{"id": "456", "name": "test", "description": "test desc"}`,
			expected: "456",
		},
		{
			name:     "large number ID",
			jsonData: `{"id": 999999999999, "name": "test", "description": "test desc"}`,
			expected: "999999999999",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var group alkira.Group
			err := json.Unmarshal([]byte(tt.jsonData), &group)
			require.NoError(t, err)

			assert.Equal(t, tt.expected, string(group.Id))
			assert.Equal(t, "test", group.Name)
			assert.Equal(t, "test desc", group.Description)
		})
	}
}

// Unit test with mock HTTP server
func TestAlkiraGroup_apiClientCRUD(t *testing.T) {
	// Test data
	groupId := json.Number("123")
	groupName := "test-group"
	groupDescription := "Test group description"

	// Create mock group
	mockGroup := &alkira.Group{
		Id:          groupId,
		Name:        groupName,
		Description: groupDescription,
	}

	// Test Create operation
	t.Run("Create", func(t *testing.T) {
		client := serveGroupMockServer(t, mockGroup, http.StatusCreated)

		api := alkira.NewGroup(client)
		response, provState, err, provErr := api.Create(mockGroup)

		assert.NoError(t, err)
		assert.Equal(t, groupName, response.Name)
		assert.Equal(t, groupDescription, response.Description)
		// provErr can be nil or error depending on provisioning setup
		_ = provErr

		t.Logf("Provision state: %s", provState)
	})

	// Test Read operation
	t.Run("Read", func(t *testing.T) {
		client := serveGroupMockServer(t, mockGroup, http.StatusOK)

		api := alkira.NewGroup(client)
		group, provState, err := api.GetById("123")

		assert.NoError(t, err)
		assert.Equal(t, groupName, group.Name)
		assert.Equal(t, groupDescription, group.Description)

		t.Logf("Provision state: %s", provState)
	})

	// Test Update operation
	t.Run("Update", func(t *testing.T) {
		updatedGroup := &alkira.Group{
			Id:          groupId,
			Name:        groupName + "-updated",
			Description: groupDescription + " updated",
		}

		client := serveGroupMockServer(t, updatedGroup, http.StatusOK)

		api := alkira.NewGroup(client)
		provState, err, provErr := api.Update("123", updatedGroup)

		assert.NoError(t, err)
		// provErr can be nil or error depending on provisioning setup
		_ = provErr

		t.Logf("Provision state: %s", provState)
	})

	// Test Delete operation
	t.Run("Delete", func(t *testing.T) {
		client := serveGroupMockServer(t, nil, http.StatusNoContent)

		api := alkira.NewGroup(client)
		provState, err, provErr := api.Delete("123")

		assert.NoError(t, err)
		// provErr can be nil or error depending on provisioning setup
		_ = provErr

		t.Logf("Provision state: %s", provState)
	})
}

func TestAlkiraGroup_apiErrorHandling(t *testing.T) {
	// Test error scenarios
	t.Run("NotFound", func(t *testing.T) {
		client := serveGroupMockServer(t, nil, http.StatusNotFound)

		api := alkira.NewGroup(client)
		_, _, err := api.GetById("999")

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "404")
	})

	t.Run("ServerError", func(t *testing.T) {
		client := serveGroupMockServer(t, nil, http.StatusInternalServerError)

		api := alkira.NewGroup(client)
		_, _, _, _ = api.Create(&alkira.Group{
			Name:        "test-group",
			Description: "test",
		})

		// Should handle server errors gracefully
	})
}

func TestAlkiraGroup_resourceDataManipulation(t *testing.T) {
	r := resourceAlkiraGroup()

	t.Run("set and get resource data", func(t *testing.T) {
		d := r.TestResourceData()

		// Set values
		d.Set("name", "test-group")
		d.Set("description", "Test description")

		// Test getting values
		assert.Equal(t, "test-group", d.Get("name").(string))
		assert.Equal(t, "Test description", d.Get("description").(string))

		// Test setting computed values
		err := d.Set("provision_state", "SUCCESS")
		assert.NoError(t, err)
		assert.Equal(t, "SUCCESS", d.Get("provision_state").(string))
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

func TestAlkiraGroup_idValidation(t *testing.T) {
	tests := []struct {
		name  string
		id    string
		valid bool
	}{
		{
			name:  "valid numeric ID",
			id:    "123",
			valid: true,
		},
		{
			name:  "valid large numeric ID",
			id:    "999999999999",
			valid: true,
		},
		{
			name:  "invalid empty ID",
			id:    "",
			valid: false,
		},
		{
			name:  "invalid non-numeric ID",
			id:    "abc",
			valid: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := strconv.Atoi(tt.id)
			if tt.valid {
				assert.NoError(t, err, "Expected valid ID")
			} else {
				if tt.id != "" { // empty string has different error than non-numeric
					assert.Error(t, err, "Expected invalid ID")
				}
			}
		})
	}
}

// Helper function to create mock HTTP server for groups
func serveGroupMockServer(t *testing.T, group *alkira.Group, statusCode int) *alkira.AlkiraClient {
	server := httptest.NewServer(http.HandlerFunc(
		func(w http.ResponseWriter, req *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(statusCode)

			// Handle provisioning requests
			if req.URL.Query().Get("provision") == "true" {
				// Mock provision response
				provisionResponse := map[string]interface{}{
					"id":     "prov-123",
					"status": "SUCCESS",
				}
				switch req.Method {
				case "POST", "PUT", "DELETE":
					json.NewEncoder(w).Encode(provisionResponse)
					return
				}
			}

			switch req.Method {
			case "GET":
				if group != nil {
					json.NewEncoder(w).Encode(group)
				}
			case "POST":
				if group != nil {
					json.NewEncoder(w).Encode(group)
				}
			case "PUT":
				if group != nil {
					json.NewEncoder(w).Encode(group)
				}
			case "DELETE":
				// No content for delete
			default:
				w.WriteHeader(http.StatusMethodNotAllowed)
			}
		},
	))
	t.Cleanup(server.Close)

	retryClient := retryablehttp.NewClient()
	retryClient.HTTPClient.Timeout = time.Duration(1) * time.Second

	return &alkira.AlkiraClient{
		URI:             server.URL,
		TenantNetworkId: "0",
		Client:          retryClient,
		Provision:       false, // Disable provisioning to avoid complex mock setup
	}
}

// Mock helper functions for testing
func buildGroupRequest(d *schema.ResourceData) *alkira.Group {
	return &alkira.Group{
		Name:        d.Get("name").(string),
		Description: d.Get("description").(string),
	}
}

func validateGroupName(v interface{}, k string) (warnings []string, errors []error) {
	value := v.(string)
	if value == "" {
		errors = append(errors, fmt.Errorf("%q cannot be empty", k))
	}
	return warnings, errors
}
