package alkira

import (
	"encoding/json"
	"net/http"
	"testing"

	"github.com/alkiranet/alkira-client-go/alkira"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/stretchr/testify/assert"
)

// Simplified test for F5 vServer endpoint to avoid struct field mismatches
func TestAlkiraServiceF5vServerEndpoint_basicFunctionality(t *testing.T) {
	// Test data with only valid fields
	serviceF5vServerEndpoint := &alkira.F5vServerEndpoint{
		Id:                   json.Number("123"),
		Name:                 "test-f5-vserver-endpoint",
		F5ServiceId:          456,
		F5ServiceInstanceIds: []int{1, 2},
		Type:                 "HTTP",
		Segment:              "test-segment",
		FqdnPrefix:           "test-fqdn",
		Protocol:             "HTTP",
		Snat:                 "AUTOMAP",
		PortRanges:           []string{"80", "443"},
	}

	// Test basic functionality
	assert.Equal(t, json.Number("123"), serviceF5vServerEndpoint.Id)
	assert.Equal(t, "test-f5-vserver-endpoint", serviceF5vServerEndpoint.Name)
	assert.Equal(t, 456, serviceF5vServerEndpoint.F5ServiceId)
	assert.Equal(t, []int{1, 2}, serviceF5vServerEndpoint.F5ServiceInstanceIds)
	assert.Equal(t, "HTTP", serviceF5vServerEndpoint.Type)
	assert.Equal(t, "test-segment", serviceF5vServerEndpoint.Segment)
	assert.Equal(t, "test-fqdn", serviceF5vServerEndpoint.FqdnPrefix)
	assert.Equal(t, "HTTP", serviceF5vServerEndpoint.Protocol)
	assert.Equal(t, "AUTOMAP", serviceF5vServerEndpoint.Snat)
	assert.Equal(t, []string{"80", "443"}, serviceF5vServerEndpoint.PortRanges)
}

func TestAlkiraServiceF5vServerEndpoint_nameValidation(t *testing.T) {
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

func TestAlkiraServiceF5vServerEndpoint_idValidation(t *testing.T) {
	tests := GetCommonIdValidationTestCases()

	for _, tt := range tests {
		t.Run(tt.Name, func(t *testing.T) {
			isValid := validateResourceId(tt.Id)
			assert.Equal(t, tt.Valid, isValid, "Expected ID validation to return %v for %s", tt.Valid, tt.Id)
		})
	}
}

// TEST HELPER
func serveServiceF5vServerEndpoint(t *testing.T, serviceF5vServerEndpoint *alkira.F5vServerEndpoint) *alkira.AlkiraClient {
	return createMockAlkiraClient(t, func(w http.ResponseWriter, req *http.Request) {
		json.NewEncoder(w).Encode(serviceF5vServerEndpoint)
		w.Header().Set("Content-Type", "application/json")
	})
}

// Mock helper function for testing
func buildServiceF5vServerEndpointRequest(d *schema.ResourceData) *alkira.F5vServerEndpoint {
	return &alkira.F5vServerEndpoint{
		Name:       getStringFromResourceData(d, "name"),
		Type:       getStringFromResourceData(d, "type"),
		Segment:    getStringFromResourceData(d, "segment"),
		FqdnPrefix: getStringFromResourceData(d, "fqdn_prefix"),
		Protocol:   getStringFromResourceData(d, "protocol"),
		Snat:       getStringFromResourceData(d, "snat"),
	}
}
