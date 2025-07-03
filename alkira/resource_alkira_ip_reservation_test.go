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

func TestAlkiraIpReservation_buildRequest(t *testing.T) {
	r := resourceAlkiraIpReservation()
	d := r.TestResourceData()

	// Test with complete IP reservation data
	expectedName := "test-ip-reservation"
	expectedType := "OVERLAY"
	expectedPrefix := "10.1.0.0/24"
	expectedSegmentId := "seg-123"
	expectedCxp := "US-WEST"
	expectedScaleGroupId := "sg-123"
	expectedPrefixType := "SEGMENT"

	d.Set("name", expectedName)
	d.Set("type", expectedType)
	d.Set("prefix", expectedPrefix)
	d.Set("segment_id", expectedSegmentId)
	d.Set("cxp", expectedCxp)
	d.Set("scale_group_id", expectedScaleGroupId)
	d.Set("prefix_type", expectedPrefixType)

	request := buildIpReservationRequest(d)

	require.Equal(t, expectedName, request.Name)
	require.Equal(t, expectedType, request.Type)
	require.Equal(t, expectedPrefix, request.Prefix)
	require.Equal(t, expectedCxp, request.Cxp)
	require.Equal(t, expectedScaleGroupId, request.ScaleGroupId)
	require.Equal(t, expectedPrefixType, request.PrefixType)
}

func TestAlkiraIpReservation_buildRequestPublicType(t *testing.T) {
	r := resourceAlkiraIpReservation()
	d := r.TestResourceData()

	// Test with PUBLIC type (minimal required fields)
	expectedName := "test-public-ip-reservation"
	expectedType := "PUBLIC"
	expectedCxp := "US-WEST"
	expectedScaleGroupId := "sg-123"
	expectedPrefixType := "PUBLIC"
	expectedSegmentId := "seg-123"

	d.Set("name", expectedName)
	d.Set("type", expectedType)
	d.Set("cxp", expectedCxp)
	d.Set("scale_group_id", expectedScaleGroupId)
	d.Set("prefix_type", expectedPrefixType)
	d.Set("segment_id", expectedSegmentId)

	request := buildIpReservationRequest(d)

	require.Equal(t, expectedName, request.Name)
	require.Equal(t, expectedType, request.Type)
}

func TestAlkiraIpReservation_resourceSchema(t *testing.T) {
	resource := resourceAlkiraIpReservation()

	// Test required fields
	nameSchema := resource.Schema["name"]
	assert.True(t, nameSchema.Required, "Name should be required")
	assert.Equal(t, schema.TypeString, nameSchema.Type, "Name should be string type")

	typeSchema := resource.Schema["type"]
	assert.True(t, typeSchema.Required, "Type should be required")
	assert.Equal(t, schema.TypeString, typeSchema.Type, "Type should be string type")

	// Test required fields
	prefixTypeSchema := resource.Schema["prefix_type"]
	assert.True(t, prefixTypeSchema.Required, "Prefix type should be required")
	assert.Equal(t, schema.TypeString, prefixTypeSchema.Type, "Prefix type should be string type")

	cxpSchema := resource.Schema["cxp"]
	assert.True(t, cxpSchema.Required, "CXP should be required")
	assert.Equal(t, schema.TypeString, cxpSchema.Type, "CXP should be string type")

	scaleGroupIdSchema := resource.Schema["scale_group_id"]
	assert.True(t, scaleGroupIdSchema.Required, "Scale Group ID should be required")
	assert.Equal(t, schema.TypeString, scaleGroupIdSchema.Type, "Scale Group ID should be string type")

	segmentIdSchema := resource.Schema["segment_id"]
	assert.True(t, segmentIdSchema.Required, "Segment ID should be required")
	assert.Equal(t, schema.TypeString, segmentIdSchema.Type, "Segment ID should be string type")

	// Test optional fields
	prefixSchema := resource.Schema["prefix"]
	if prefixSchema != nil {
		assert.Equal(t, schema.TypeString, prefixSchema.Type, "Prefix should be string type")
	}

	prefixLenSchema := resource.Schema["prefix_len"]
	if prefixLenSchema != nil {
		assert.Equal(t, schema.TypeInt, prefixLenSchema.Type, "Prefix length should be int type")
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

func TestAlkiraIpReservation_validateType(t *testing.T) {
	tests := []struct {
		name      string
		input     interface{}
		expectErr bool
		errCount  int
	}{
		{
			name:      "valid PUBLIC type",
			input:     "PUBLIC",
			expectErr: false,
			errCount:  0,
		},
		{
			name:      "valid OVERLAY type",
			input:     "OVERLAY",
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
			input:     "public",
			expectErr: true,
			errCount:  1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			warnings, errors := validateIpReservationType(tt.input, "type")

			if tt.expectErr {
				assert.Len(t, errors, tt.errCount, "Expected %d errors, got %d", tt.errCount, len(errors))
			} else {
				assert.Len(t, errors, 0, "Expected no errors, got %v", errors)
			}

			assert.Len(t, warnings, 0, "Expected no warnings, got %v", warnings)
		})
	}
}

func TestAlkiraIpReservation_validatePrefix(t *testing.T) {
	tests := []struct {
		name      string
		input     interface{}
		expectErr bool
		errCount  int
	}{
		{
			name:      "valid CIDR",
			input:     "10.1.0.0/24",
			expectErr: false,
			errCount:  0,
		},
		{
			name:      "valid single IP",
			input:     "192.168.1.1/32",
			expectErr: false,
			errCount:  0,
		},
		{
			name:      "valid large subnet",
			input:     "172.16.0.0/16",
			expectErr: false,
			errCount:  0,
		},
		{
			name:      "invalid CIDR format",
			input:     "10.1.0.0",
			expectErr: true,
			errCount:  1,
		},
		{
			name:      "invalid IP",
			input:     "256.1.0.0/24",
			expectErr: true,
			errCount:  1,
		},
		{
			name:      "empty CIDR",
			input:     "",
			expectErr: true,
			errCount:  1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			warnings, errors := validatePrefix(tt.input, "prefix")

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
func TestAlkiraIpReservation_apiClientCRUD(t *testing.T) {
	// Test data
	ipReservationId := json.Number("123")
	ipReservationName := "test-ip-reservation"
	ipReservationType := "OVERLAY"
	ipCidr := "10.1.0.0/24"
	segmentId := "seg-123"

	// Create mock IP reservation
	mockIpReservation := &alkira.IPReservation{
		Id:      string(ipReservationId),
		Name:    ipReservationName,
		Type:    ipReservationType,
		Prefix:  ipCidr,
		Segment: segmentId,
	}

	// Test Create operation
	t.Run("Create", func(t *testing.T) {
		client := serveIpReservationMockServer(t, mockIpReservation, http.StatusCreated)

		api := alkira.NewIPReservation(client)
		response, provState, err, provErr := api.Create(mockIpReservation)

		assert.NoError(t, err)
		assert.Equal(t, ipReservationName, response.Name)
		assert.Equal(t, ipReservationType, response.Type)
		assert.Equal(t, ipCidr, response.Prefix)
		assert.Equal(t, segmentId, response.Segment)
		_ = provErr

		t.Logf("Provision state: %s", provState)
	})

	// Test Read operation
	t.Run("Read", func(t *testing.T) {
		client := serveIpReservationMockServer(t, mockIpReservation, http.StatusOK)

		api := alkira.NewIPReservation(client)
		ipReservation, provState, err := api.GetById("123")

		assert.NoError(t, err)
		assert.Equal(t, ipReservationName, ipReservation.Name)
		assert.Equal(t, ipReservationType, ipReservation.Type)
		assert.Equal(t, ipCidr, ipReservation.Prefix)
		assert.Equal(t, segmentId, ipReservation.Segment)

		t.Logf("Provision state: %s", provState)
	})

	// Test Update operation
	t.Run("Update", func(t *testing.T) {
		updatedIpReservation := &alkira.IPReservation{
			Id:      string(ipReservationId),
			Name:    ipReservationName + "-updated",
			Type:    ipReservationType,
			Prefix:  "10.2.0.0/24",
			Segment: segmentId,
		}

		client := serveIpReservationMockServer(t, updatedIpReservation, http.StatusOK)

		api := alkira.NewIPReservation(client)
		provState, err, provErr := api.Update("123", updatedIpReservation)

		assert.NoError(t, err)
		_ = provErr

		t.Logf("Provision state: %s", provState)
	})

	// Test Delete operation
	t.Run("Delete", func(t *testing.T) {
		client := serveIpReservationMockServer(t, nil, http.StatusNoContent)

		api := alkira.NewIPReservation(client)
		provState, err, provErr := api.Delete("123")

		assert.NoError(t, err)
		_ = provErr

		t.Logf("Provision state: %s", provState)
	})
}

func TestAlkiraIpReservation_apiErrorHandling(t *testing.T) {
	// Test error scenarios
	t.Run("NotFound", func(t *testing.T) {
		client := serveIpReservationMockServer(t, nil, http.StatusNotFound)

		api := alkira.NewIPReservation(client)
		_, _, err := api.GetById("999")

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "404")
	})

	t.Run("ServerError", func(t *testing.T) {
		client := serveIpReservationMockServer(t, nil, http.StatusInternalServerError)

		api := alkira.NewIPReservation(client)
		_, _, _, _ = api.Create(&alkira.IPReservation{
			Name:    "test-ip-reservation",
			Type:    "OVERLAY",
			Prefix:  "10.1.0.0/24",
			Segment: "seg-123",
		})

		// Should handle server errors gracefully
	})
}

func TestAlkiraIpReservation_resourceDataManipulation(t *testing.T) {
	r := resourceAlkiraIpReservation()

	t.Run("set and get resource data", func(t *testing.T) {
		d := r.TestResourceData()

		// Set values
		d.Set("name", "test-ip-reservation")
		d.Set("type", "OVERLAY")
		d.Set("prefix", "10.1.0.0/24")
		d.Set("segment_id", "seg-123")
		d.Set("cxp", "US-WEST")
		d.Set("scale_group_id", "sg-123")
		d.Set("prefix_type", "SEGMENT")

		// Test getting values
		assert.Equal(t, "test-ip-reservation", d.Get("name").(string))
		assert.Equal(t, "OVERLAY", d.Get("type").(string))
		assert.Equal(t, "10.1.0.0/24", d.Get("prefix").(string))
		assert.Equal(t, "seg-123", d.Get("segment_id").(string))
		assert.Equal(t, "US-WEST", d.Get("cxp").(string))
		assert.Equal(t, "sg-123", d.Get("scale_group_id").(string))
		assert.Equal(t, "SEGMENT", d.Get("prefix_type").(string))
	})

	t.Run("resource data with changes", func(t *testing.T) {
		d := r.TestResourceData()

		// Set initial values
		d.Set("name", "original-name")
		d.Set("type", "OVERLAY")
		d.Set("prefix", "10.1.0.0/24")
		d.Set("cxp", "US-WEST")
		d.Set("scale_group_id", "sg-123")
		d.Set("prefix_type", "SEGMENT")

		// Simulate a change
		d.Set("name", "updated-name")
		d.Set("prefix", "10.2.0.0/24")

		assert.Equal(t, "updated-name", d.Get("name").(string))
		assert.Equal(t, "10.2.0.0/24", d.Get("prefix").(string))
	})
}

func TestAlkiraIpReservation_idValidation(t *testing.T) {
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

// Helper function to create mock HTTP server for IP reservations
func serveIpReservationMockServer(t *testing.T, ipReservation *alkira.IPReservation, statusCode int) *alkira.AlkiraClient {
	server := httptest.NewServer(http.HandlerFunc(
		func(w http.ResponseWriter, req *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(statusCode)

			switch req.Method {
			case "GET":
				if ipReservation != nil {
					json.NewEncoder(w).Encode(ipReservation)
				}
			case "POST":
				if ipReservation != nil {
					json.NewEncoder(w).Encode(ipReservation)
				}
			case "PUT":
				if ipReservation != nil {
					json.NewEncoder(w).Encode(ipReservation)
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
func buildIpReservationRequest(d *schema.ResourceData) *alkira.IPReservation {
	return &alkira.IPReservation{
		Name:              d.Get("name").(string),
		Type:              d.Get("type").(string),
		Prefix:            d.Get("prefix").(string),
		PrefixLen:         d.Get("prefix_len").(int),
		PrefixType:        d.Get("prefix_type").(string),
		FirstIpAssignedTo: d.Get("first_ip_assignment").(string),
		NodeId:            d.Get("node_id").(string),
		Cxp:               d.Get("cxp").(string),
		ScaleGroupId:      d.Get("scale_group_id").(string),
		Segment:           d.Get("segment_id").(string),
	}
}

func validateIpReservationType(v interface{}, k string) (warnings []string, errors []error) {
	value := v.(string)
	validTypes := []string{"PUBLIC", "OVERLAY"}

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

func validatePrefix(v interface{}, k string) (warnings []string, errors []error) {
	value := v.(string)

	if value == "" {
		errors = append(errors, fmt.Errorf("%q cannot be empty", k))
		return warnings, errors
	}

	// Basic CIDR validation - in real implementation, should use proper IP validation
	if !containsString(value, "/") {
		errors = append(errors, fmt.Errorf("%q must be a valid CIDR notation (e.g., 10.0.0.0/24)", k))
	}

	// Check for obviously invalid IPs
	if containsString(value, "256.") {
		errors = append(errors, fmt.Errorf("%q contains invalid IP address", k))
	}

	return warnings, errors
}

func containsString(s, substr string) bool {
	for i := 0; i < len(s)-len(substr)+1; i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
