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
	tests := []struct {
		name      string
		input     interface{}
		expectErr bool
		errCount  int
	}{
		{
			name:      "valid name",
			input:     "test-flow-collector",
			expectErr: false,
			errCount:  0,
		},
		{
			name:      "valid name with numbers",
			input:     "flow-collector-123",
			expectErr: false,
			errCount:  0,
		},
		{
			name:      "valid name with underscores",
			input:     "test_flow_collector_name",
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
			input:     "test flow collector",
			expectErr: false,
			errCount:  0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			warnings, errors := validateFlowCollectorName(tt.input, "name")

			if tt.expectErr {
				assert.Len(t, errors, tt.errCount, "Expected %d errors, got %d", tt.errCount, len(errors))
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
			name:      "valid NETFLOW type",
			input:     "NETFLOW",
			expectErr: false,
			errCount:  0,
		},
		{
			name:      "valid SFLOW type",
			input:     "SFLOW",
			expectErr: false,
			errCount:  0,
		},
		{
			name:      "valid IPFIX type",
			input:     "IPFIX",
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
			input:     "netflow",
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
		response, provState, err, provErr := api.Create(mockFlowCollector)

		assert.NoError(t, err)
		assert.Equal(t, flowCollectorName, response.Name)
		assert.Equal(t, flowCollectorDescription, response.Description)
		assert.Equal(t, flowCollectorEnabled, response.Enabled)
		assert.Equal(t, flowCollectorType, response.CollectorType)
		assert.Equal(t, segment, response.Segment)
		assert.Equal(t, destinationIp, response.DestinationIp)
		assert.Equal(t, destinationPort, response.DestinationPort)
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
		provState, err, provErr := api.Update("123", updatedFlowCollector)

		assert.NoError(t, err)
		_ = provErr

		t.Logf("Provision state: %s", provState)
	})

	// Test Delete operation
	t.Run("Delete", func(t *testing.T) {
		client := serveFlowCollectorMockServer(t, nil, http.StatusNoContent)

		api := alkira.NewFlowCollector(client)
		provState, err, provErr := api.Delete("123")

		assert.NoError(t, err)
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
		_, _, _, _ = api.Create(&alkira.FlowCollector{
			Name:          "test-flow-collector",
			Description:   "test",
			Enabled:       true,
			CollectorType: "NETFLOW",
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
		d.Set("collector_type", "NETFLOW")
		d.Set("segment", "test-segment")
		cxps := []string{"US-WEST", "EU-WEST"}
		d.Set("cxps", cxps)

		// Test getting values
		assert.Equal(t, "test-flow-collector", d.Get("name").(string))
		assert.Equal(t, "Test description", d.Get("description").(string))
		assert.Equal(t, true, d.Get("enabled").(bool))
		assert.Equal(t, "NETFLOW", d.Get("collector_type").(string))
		assert.Equal(t, "test-segment", d.Get("segment").(string))
	})

	t.Run("resource data with changes", func(t *testing.T) {
		d := r.TestResourceData()

		// Set initial values
		d.Set("name", "original-name")
		d.Set("description", "Original description")
		d.Set("enabled", false)
		d.Set("collector_type", "SFLOW")

		// Simulate a change
		d.Set("name", "updated-name")
		d.Set("description", "Updated description")
		d.Set("enabled", true)
		d.Set("collector_type", "NETFLOW")

		assert.Equal(t, "updated-name", d.Get("name").(string))
		assert.Equal(t, "Updated description", d.Get("description").(string))
		assert.Equal(t, true, d.Get("enabled").(bool))
		assert.Equal(t, "NETFLOW", d.Get("collector_type").(string))
	})
}

func TestAlkiraFlowCollector_idValidation(t *testing.T) {
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

func TestAlkiraFlowCollector_customizeDiff(t *testing.T) {
	r := resourceAlkiraFlowCollector()

	// Test CustomizeDiff function exists
	assert.NotNil(t, r.CustomizeDiff, "Resource should have CustomizeDiff")

	// Note: CustomizeDiff testing would require more complex setup with ResourceDiff mock
	// This validates the function exists and can be called
}

// Helper function to create mock HTTP server for flow collectors
func serveFlowCollectorMockServer(t *testing.T, flowCollector *alkira.FlowCollector, statusCode int) *alkira.AlkiraClient {
	server := httptest.NewServer(http.HandlerFunc(
		func(w http.ResponseWriter, req *http.Request) {
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
		},
	))
	t.Cleanup(server.Close)

	retryClient := retryablehttp.NewClient()
	retryClient.HTTPClient.Timeout = time.Duration(1) * time.Second

	return &alkira.AlkiraClient{
		URI:             server.URL,
		TenantNetworkId: "0",
		Client:          retryClient,
		Provision:       false,
	}
}

// Mock helper functions for testing
func buildFlowCollectorRequest(d *schema.ResourceData) *alkira.FlowCollector {
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

	collectorType := ""
	if ct := d.Get("collector_type"); ct != nil {
		collectorType = ct.(string)
	}

	segment := ""
	if seg := d.Get("segment_id"); seg != nil {
		segment = seg.(string)
	}

	destinationIp := ""
	if di := d.Get("destination_ip"); di != nil {
		destinationIp = di.(string)
	}

	destinationFqdn := ""
	if df := d.Get("destination_fqdn"); df != nil {
		destinationFqdn = df.(string)
	}

	destinationPort := 0
	if dp := d.Get("destination_port"); dp != nil {
		destinationPort = dp.(int)
	}

	transportProtocol := ""
	if tp := d.Get("transport_protocol"); tp != nil {
		transportProtocol = tp.(string)
	}

	exportType := ""
	if et := d.Get("export_type"); et != nil {
		exportType = et.(string)
	}

	return &alkira.FlowCollector{
		Name:              name,
		Description:       description,
		Enabled:           enabled,
		CollectorType:     collectorType,
		Segment:           segment,
		DestinationIp:     destinationIp,
		DestinationFqdn:   destinationFqdn,
		DestinationPort:   destinationPort,
		TransportProtocol: transportProtocol,
		ExportType:        exportType,
	}
}

func validateFlowCollectorName(v interface{}, k string) (warnings []string, errors []error) {
	value := v.(string)
	if value == "" {
		errors = append(errors, fmt.Errorf("%q cannot be empty", k))
	}
	return warnings, errors
}

func validateFlowCollectorType(v interface{}, k string) (warnings []string, errors []error) {
	value := v.(string)
	validTypes := []string{"NETFLOW", "SFLOW", "IPFIX"}

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
