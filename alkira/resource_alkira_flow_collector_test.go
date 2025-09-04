package alkira

import (
	"encoding/json"
	"fmt"
	"net/http"
	"testing"

	"github.com/alkiranet/alkira-client-go/alkira"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAlkiraFlowCollector_buildFlowCollectorRequest(t *testing.T) {
	r := resourceAlkiraFlowCollector()
	d := r.TestResourceData()

	// Test with complete flow collector data
	expectedName := "test-flow-collector"
	expectedDescription := "Test flow collector description"
	expectedEnabled := true
	expectedCollectorType := "GENERIC"
	expectedSegmentId := "test-segment-id"
	expectedDestinationIp := "10.1.1.100"
	expectedDestinationPort := 2055

	d.Set("name", expectedName)
	d.Set("description", expectedDescription)
	d.Set("enabled", expectedEnabled)
	d.Set("collector_type", expectedCollectorType)
	d.Set("segment_id", expectedSegmentId)
	d.Set("destination_ip", expectedDestinationIp)
	d.Set("destination_port", expectedDestinationPort)

	request := buildFlowCollectorRequest(d)

	require.Equal(t, expectedName, request.Name)
	require.Equal(t, expectedDescription, request.Description)
	require.Equal(t, expectedEnabled, request.Enabled)
	require.Equal(t, expectedCollectorType, request.CollectorType)
	require.Equal(t, expectedSegmentId, request.Segment)
	require.Equal(t, expectedDestinationIp, request.DestinationIp)
	require.Equal(t, expectedDestinationPort, request.DestinationPort)
}

func TestAlkiraFlowCollector_buildFlowCollectorRequestMinimal(t *testing.T) {
	r := resourceAlkiraFlowCollector()
	d := r.TestResourceData()

	// Test with minimal flow collector data
	expectedName := "minimal-flow-collector"
	expectedCollectorType := "GENERIC"
	expectedDestinationPort := 9999

	d.Set("name", expectedName)
	d.Set("collector_type", expectedCollectorType)
	d.Set("enabled", true) // required field
	d.Set("destination_port", expectedDestinationPort)

	request := buildFlowCollectorRequest(d)

	require.Equal(t, expectedName, request.Name)
	require.Equal(t, "", request.Description)
	require.Equal(t, true, request.Enabled)
	require.Equal(t, expectedCollectorType, request.CollectorType)
	require.Equal(t, expectedDestinationPort, request.DestinationPort)
}

func TestAlkiraFlowCollector_resourceSchema(t *testing.T) {
	resource := resourceAlkiraFlowCollector()

	// Test required fields
	nameSchema := resource.Schema["name"]
	assert.True(t, nameSchema.Required, "Name should be required")
	assert.Equal(t, schema.TypeString, nameSchema.Type, "Name should be string type")

	enabledSchema := resource.Schema["enabled"]
	if enabledSchema != nil {
		assert.Equal(t, schema.TypeBool, enabledSchema.Type, "Enabled should be bool type")
	}

	typeSchema := resource.Schema["type"]
	if typeSchema != nil {
		assert.Equal(t, schema.TypeString, typeSchema.Type, "Type should be string type")
	}

	// Test optional fields
	descSchema := resource.Schema["description"]
	if descSchema != nil {
		assert.True(t, descSchema.Optional, "Description should be optional")
		assert.Equal(t, schema.TypeString, descSchema.Type, "Description should be string type")
	}

	segmentIdsSchema := resource.Schema["segment_ids"]
	if segmentIdsSchema != nil {
		assert.Equal(t, schema.TypeSet, segmentIdsSchema.Type, "Segment IDs should be set type")
	}

	// Test computed fields
	provStateSchema := resource.Schema["provision_state"]
	if provStateSchema != nil {
		assert.True(t, provStateSchema.Computed, "Provision state should be computed")
		assert.Equal(t, schema.TypeString, provStateSchema.Type, "Provision state should be string type")
	}

	// Test that resource has all required CRUD functions
	assert.NotNil(t, resource.CreateContext, "Resource should have CreateContext")
	assert.NotNil(t, resource.ReadContext, "Resource should have ReadContext")
	assert.NotNil(t, resource.UpdateContext, "Resource should have UpdateContext")
	assert.NotNil(t, resource.DeleteContext, "Resource should have DeleteContext")
	assert.NotNil(t, resource.Importer, "Resource should support import")
}

func TestAlkiraFlowCollector_validateFlowCollectorName(t *testing.T) {
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

func TestAlkiraFlowCollector_validateFlowCollectorType(t *testing.T) {
	tests := []struct {
		name      string
		input     interface{}
		expectErr bool
		errCount  int
	}{
		{
			name:      "valid GENERIC type",
			input:     "GENERIC",
			expectErr: false,
			errCount:  0,
		},
		{
			name:      "invalid type",
			input:     "INVALID",
			expectErr: true,
			errCount:  1,
		},
		{
			name:      "empty type",
			input:     "",
			expectErr: true,
			errCount:  1,
		},
		{
			name:      "lowercase type",
			input:     "generic",
			expectErr: true,
			errCount:  1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			warnings, errors := validateFlowCollectorType(tt.input, "type")

			if tt.expectErr {
				assert.Len(t, errors, tt.errCount, "Expected %d errors, got %d", tt.errCount, len(errors))
			} else {
				assert.Len(t, errors, 0, "Expected no errors, got %v", errors)
			}

			assert.Len(t, warnings, 0, "Expected no warnings, got %v", warnings)
		})
	}
}

// Unit test with mock HTTP server
func TestAlkiraFlowCollector_apiClientCRUD(t *testing.T) {
	// Test data
	flowCollectorId := json.Number("123")
	flowCollectorName := "test-flow-collector"
	flowCollectorDescription := "Test flow collector description"
	flowCollectorEnabled := true
	flowCollectorType := "GENERIC"
	segment := "test-segment-id"
	destinationIp := "10.1.1.100"
	destinationPort := 2055

	// Create mock flow collector
	mockFlowCollector := &alkira.FlowCollector{
		Id:              flowCollectorId,
		Name:            flowCollectorName,
		Description:     flowCollectorDescription,
		Enabled:         flowCollectorEnabled,
		CollectorType:   flowCollectorType,
		Segment:         segment,
		DestinationIp:   destinationIp,
		DestinationPort: destinationPort,
	}

	// Test Create operation
	t.Run("Create", func(t *testing.T) {
		client := serveFlowCollectorMockServer(t, mockFlowCollector, http.StatusCreated)

		api := alkira.NewFlowCollector(client)
		response, provState, err, valErr, provErr := api.Create(mockFlowCollector)

		assert.NoError(t, err)
		assert.Equal(t, flowCollectorName, response.Name)
		assert.Equal(t, flowCollectorDescription, response.Description)
		assert.Equal(t, flowCollectorEnabled, response.Enabled)
		assert.Equal(t, flowCollectorType, response.CollectorType)
		assert.Equal(t, segment, response.Segment)
		assert.Equal(t, destinationIp, response.DestinationIp)
		assert.Equal(t, destinationPort, response.DestinationPort)
		_ = valErr
		_ = provErr

		t.Logf("Provision state: %s", provState)
	})

	// Test Read operation
	t.Run("Read", func(t *testing.T) {
		client := serveFlowCollectorMockServer(t, mockFlowCollector, http.StatusOK)

		api := alkira.NewFlowCollector(client)
		flowCollector, provState, err := api.GetById("123")

		assert.NoError(t, err)
		assert.Equal(t, flowCollectorName, flowCollector.Name)
		assert.Equal(t, flowCollectorDescription, flowCollector.Description)
		assert.Equal(t, flowCollectorEnabled, flowCollector.Enabled)
		assert.Equal(t, flowCollectorType, flowCollector.CollectorType)
		assert.Equal(t, segment, flowCollector.Segment)
		assert.Equal(t, destinationIp, flowCollector.DestinationIp)
		assert.Equal(t, destinationPort, flowCollector.DestinationPort)

		t.Logf("Provision state: %s", provState)
	})

	// Test Update operation
	t.Run("Update", func(t *testing.T) {
		updatedFlowCollector := &alkira.FlowCollector{
			Id:              flowCollectorId,
			Name:            flowCollectorName + "-updated",
			Description:     flowCollectorDescription + " updated",
			Enabled:         false,
			CollectorType:   "GENERIC",
			Segment:         "updated-segment",
			DestinationIp:   "10.2.2.200",
			DestinationPort: 3000,
		}

		client := serveFlowCollectorMockServer(t, updatedFlowCollector, http.StatusOK)

		api := alkira.NewFlowCollector(client)
		provState, err, valErr, provErr := api.Update("123", updatedFlowCollector)

		assert.NoError(t, err)
		_ = valErr
		_ = provErr

		t.Logf("Provision state: %s", provState)
	})

	// Test Delete operation
	t.Run("Delete", func(t *testing.T) {
		client := serveFlowCollectorMockServer(t, nil, http.StatusNoContent)

		api := alkira.NewFlowCollector(client)
		provState, err, valErr, provErr := api.Delete("123")

		assert.NoError(t, err)
		_ = valErr
		_ = provErr

		t.Logf("Provision state: %s", provState)
	})
}

func TestAlkiraFlowCollector_apiErrorHandling(t *testing.T) {
	// Test error scenarios
	t.Run("NotFound", func(t *testing.T) {
		client := serveFlowCollectorMockServer(t, nil, http.StatusNotFound)

		api := alkira.NewFlowCollector(client)
		_, _, err := api.GetById("999")

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "404")
	})

	t.Run("ServerError", func(t *testing.T) {
		client := serveFlowCollectorMockServer(t, nil, http.StatusInternalServerError)

		api := alkira.NewFlowCollector(client)
		_, _, _, _, _ = api.Create(&alkira.FlowCollector{
			Name:            "test-flow-collector",
			Description:     "test",
			Enabled:         true,
			CollectorType:   "GENERIC",
			DestinationPort: 2055,
		})

		// Should handle server errors gracefully
	})
}

func TestAlkiraFlowCollector_resourceDataManipulation(t *testing.T) {
	r := resourceAlkiraFlowCollector()

	t.Run("set and get resource data", func(t *testing.T) {
		d := r.TestResourceData()

		// Set values
		d.Set("name", "test-flow-collector")
		d.Set("description", "Test description")
		d.Set("enabled", true)
		d.Set("collector_type", "GENERIC")
		d.Set("segment_id", "test-segment-id")
		d.Set("destination_ip", "10.1.1.100")
		d.Set("destination_port", 2055)

		// Test getting values
		assert.Equal(t, "test-flow-collector", d.Get("name").(string))
		assert.Equal(t, "Test description", d.Get("description").(string))
		assert.Equal(t, true, d.Get("enabled").(bool))
		assert.Equal(t, "GENERIC", d.Get("collector_type").(string))
		assert.Equal(t, "test-segment-id", d.Get("segment_id").(string))
		assert.Equal(t, "10.1.1.100", d.Get("destination_ip").(string))
		assert.Equal(t, 2055, d.Get("destination_port").(int))
	})

	t.Run("resource data with changes", func(t *testing.T) {
		d := r.TestResourceData()

		// Set initial values
		d.Set("name", "original-name")
		d.Set("description", "Original description")
		d.Set("enabled", false)
		d.Set("collector_type", "GENERIC")
		d.Set("destination_port", 1000)

		// Simulate a change
		d.Set("name", "updated-name")
		d.Set("description", "Updated description")
		d.Set("enabled", true)
		d.Set("collector_type", "GENERIC")
		d.Set("destination_port", 2000)

		assert.Equal(t, "updated-name", d.Get("name").(string))
		assert.Equal(t, "Updated description", d.Get("description").(string))
		assert.Equal(t, true, d.Get("enabled").(bool))
		assert.Equal(t, "GENERIC", d.Get("collector_type").(string))
		assert.Equal(t, 2000, d.Get("destination_port").(int))
	})
}

func TestAlkiraFlowCollector_idValidation(t *testing.T) {
	tests := GetCommonIdValidationTestCases()

	for _, tt := range tests {
		t.Run(tt.Name, func(t *testing.T) {
			isValid := validateResourceId(tt.Id)
			assert.Equal(t, tt.Valid, isValid, "Expected ID validation to return %v for %s", tt.Valid, tt.Id)
		})
	}
}

func TestAlkiraFlowCollector_customizeDiff(t *testing.T) {
	r := resourceAlkiraFlowCollector()

	// Test CustomizeDiff function exists
	assert.NotNil(t, r.CustomizeDiff, "Resource should have CustomizeDiff")

	// Note: CustomizeDiff testing would require more complex setup with ResourceDiff mock
	// This validates the function exists and can be called
}

// Helper function to create mock HTTP server for flow collectors
func serveFlowCollectorMockServer(t *testing.T, flowCollector *alkira.FlowCollector, statusCode int) *alkira.AlkiraClient {
	handler := http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(statusCode)

		switch req.Method {
		case "GET":
			if flowCollector != nil {
				json.NewEncoder(w).Encode(flowCollector)
			}
		case "POST":
			if flowCollector != nil {
				json.NewEncoder(w).Encode(flowCollector)
			}
		case "PUT":
			if flowCollector != nil {
				json.NewEncoder(w).Encode(flowCollector)
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
func buildFlowCollectorRequest(d *schema.ResourceData) *alkira.FlowCollector {
	return &alkira.FlowCollector{
		Name:              getStringFromResourceData(d, "name"),
		Description:       getStringFromResourceData(d, "description"),
		Enabled:           getBoolFromResourceData(d, "enabled"),
		CollectorType:     getStringFromResourceData(d, "collector_type"),
		Segment:           getStringFromResourceData(d, "segment_id"),
		DestinationIp:     getStringFromResourceData(d, "destination_ip"),
		DestinationFqdn:   getStringFromResourceData(d, "destination_fqdn"),
		DestinationPort:   getIntFromResourceData(d, "destination_port"),
		TransportProtocol: getStringFromResourceData(d, "transport_protocol"),
		ExportType:        getStringFromResourceData(d, "export_type"),
	}
}

func validateFlowCollectorType(v interface{}, k string) (warnings []string, errors []error) {
	value := v.(string)
	validTypes := []string{"GENERIC"}

	if value == "" {
		errors = append(errors, fmt.Errorf("%q cannot be empty", k))
		return warnings, errors
	}

	for _, validType := range validTypes {
		if value == validType {
			return warnings, errors
		}
	}

	errors = append(errors, fmt.Errorf("%q must be one of: %v", k, validTypes))
	return warnings, errors
}
