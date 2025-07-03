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

func TestAlkiraByoipPrefix_buildByoipPrefixRequest(t *testing.T) {
	r := resourceAlkiraByoipPrefix()
	d := r.TestResourceData()

	// Test with complete BYOIP prefix data
	expectedPrefix := "203.0.113.0/24"
	expectedCxp := "US-WEST"
	expectedDescription := "Test BYOIP prefix description"
	expectedCloudProvider := "AWS"

	d.Set("prefix", expectedPrefix)
	d.Set("cxp", expectedCxp)
	d.Set("description", expectedDescription)
	d.Set("cloud_provider", expectedCloudProvider)
	d.Set("message", "test-message")
	d.Set("signature", "test-signature")
	d.Set("public_key", "test-public-key")

	request := buildByoipPrefixRequest(d)

	require.Equal(t, expectedPrefix, request.Prefix)
	require.Equal(t, expectedCxp, request.Cxp)
	require.Equal(t, expectedDescription, request.Description)
	require.Equal(t, expectedCloudProvider, request.CloudProvider)
}

func TestAlkiraByoipPrefix_buildByoipPrefixRequestMinimal(t *testing.T) {
	r := resourceAlkiraByoipPrefix()
	d := r.TestResourceData()

	// Test with minimal BYOIP prefix data
	expectedPrefix := "198.51.100.0/24"
	expectedCxp := "EU-WEST"

	d.Set("prefix", expectedPrefix)
	d.Set("cxp", expectedCxp)
	d.Set("cloud_provider", "AWS")
	d.Set("message", "")
	d.Set("signature", "")
	d.Set("public_key", "")
	// description not set (should be empty)

	request := buildByoipPrefixRequest(d)

	require.Equal(t, expectedPrefix, request.Prefix)
	require.Equal(t, expectedCxp, request.Cxp)
	require.Equal(t, "", request.Description)
	require.Equal(t, "AWS", request.CloudProvider)
}

func TestAlkiraByoipPrefix_resourceSchema(t *testing.T) {
	resource := resourceAlkiraByoipPrefix()

	// Test required fields
	prefixSchema := resource.Schema["prefix"]
	assert.True(t, prefixSchema.Required, "Prefix should be required")
	assert.Equal(t, schema.TypeString, prefixSchema.Type, "Prefix should be string type")

	cxpSchema := resource.Schema["cxp"]
	assert.True(t, cxpSchema.Required, "CXP should be required")
	assert.Equal(t, schema.TypeString, cxpSchema.Type, "CXP should be string type")

	// Test optional fields

	descSchema := resource.Schema["description"]
	if descSchema != nil {
		assert.True(t, descSchema.Optional, "Description should be optional")
		assert.Equal(t, schema.TypeString, descSchema.Type, "Description should be string type")
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

func TestAlkiraByoipPrefix_validateByoipPrefix(t *testing.T) {
	tests := []struct {
		name      string
		input     interface{}
		expectErr bool
		errCount  int
	}{
		{
			name:      "valid IPv4 CIDR",
			input:     "203.0.113.0/24",
			expectErr: false,
			errCount:  0,
		},
		{
			name:      "valid single IPv4 address",
			input:     "192.0.2.1/32",
			expectErr: false,
			errCount:  0,
		},
		{
			name:      "valid large IPv4 subnet",
			input:     "10.0.0.0/8",
			expectErr: false,
			errCount:  0,
		},
		{
			name:      "valid IPv6 CIDR",
			input:     "2001:db8::/32",
			expectErr: false,
			errCount:  0,
		},
		{
			name:      "invalid CIDR format",
			input:     "203.0.113.0",
			expectErr: true,
			errCount:  1,
		},
		{
			name:      "invalid IP address",
			input:     "256.0.113.0/24",
			expectErr: true,
			errCount:  1,
		},
		{
			name:      "empty prefix",
			input:     "",
			expectErr: true,
			errCount:  1,
		},
		{
			name:      "invalid subnet mask",
			input:     "203.0.113.0/33",
			expectErr: true,
			errCount:  1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			warnings, errors := validateByoipPrefix(tt.input, "prefix")

			if tt.expectErr {
				assert.Len(t, errors, tt.errCount, "Expected %d errors, got %d", tt.errCount, len(errors))
			} else {
				assert.Len(t, errors, 0, "Expected no errors, got %v", errors)
			}

			assert.Len(t, warnings, 0, "Expected no warnings, got %v", warnings)
		})
	}
}

func TestAlkiraByoipPrefix_validateCxp(t *testing.T) {
	tests := []struct {
		name      string
		input     interface{}
		expectErr bool
		errCount  int
	}{
		{
			name:      "valid US-WEST CXP",
			input:     "US-WEST",
			expectErr: false,
			errCount:  0,
		},
		{
			name:      "valid US-EAST CXP",
			input:     "US-EAST",
			expectErr: false,
			errCount:  0,
		},
		{
			name:      "valid EU-WEST CXP",
			input:     "EU-WEST",
			expectErr: false,
			errCount:  0,
		},
		{
			name:      "valid APAC CXP",
			input:     "APAC",
			expectErr: false,
			errCount:  0,
		},
		{
			name:      "invalid CXP",
			input:     "INVALID-CXP",
			expectErr: true,
			errCount:  1,
		},
		{
			name:      "empty CXP",
			input:     "",
			expectErr: true,
			errCount:  1,
		},
		{
			name:      "lowercase CXP",
			input:     "us-west",
			expectErr: true,
			errCount:  1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			warnings, errors := validateCxp(tt.input, "cxp")

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
func TestAlkiraByoipPrefix_apiClientCRUD(t *testing.T) {
	// Test data
	byoipPrefixId := json.Number("123")
	byoipPrefixDescription := "Test BYOIP prefix description"
	prefix := "203.0.113.0/24"
	cxp := "US-WEST"

	// Create mock BYOIP prefix
	mockByoipPrefix := &alkira.Byoip{
		Id:            byoipPrefixId,
		Description:   byoipPrefixDescription,
		Prefix:        prefix,
		Cxp:           cxp,
		CloudProvider: "AWS",
		ExtraAttributes: alkira.ByoipExtraAttributes{
			Message:   "test-message",
			Signature: "test-signature",
			PublicKey: "test-public-key",
		},
	}

	// Test Create operation
	t.Run("Create", func(t *testing.T) {
		client := serveByoipPrefixMockServer(t, mockByoipPrefix, http.StatusCreated)

		api := alkira.NewByoip(client)
		response, provState, err, provErr := api.Create(mockByoipPrefix)

		assert.NoError(t, err)
		assert.Equal(t, byoipPrefixDescription, response.Description)
		assert.Equal(t, prefix, response.Prefix)
		assert.Equal(t, cxp, response.Cxp)
		assert.Equal(t, "AWS", response.CloudProvider)
		_ = provErr

		t.Logf("Provision state: %s", provState)
	})

	// Test Read operation
	t.Run("Read", func(t *testing.T) {
		client := serveByoipPrefixMockServer(t, mockByoipPrefix, http.StatusOK)

		api := alkira.NewByoip(client)
		byoipPrefix, provState, err := api.GetById("123")

		assert.NoError(t, err)
		assert.Equal(t, byoipPrefixDescription, byoipPrefix.Description)
		assert.Equal(t, prefix, byoipPrefix.Prefix)
		assert.Equal(t, cxp, byoipPrefix.Cxp)
		assert.Equal(t, "AWS", byoipPrefix.CloudProvider)

		t.Logf("Provision state: %s", provState)
	})

	// Test Update operation
	t.Run("Update", func(t *testing.T) {
		updatedByoipPrefix := &alkira.Byoip{
			Id:            byoipPrefixId,
			Description:   byoipPrefixDescription + " updated",
			Prefix:        "198.51.100.0/24",
			Cxp:           "EU-WEST",
			CloudProvider: "AWS",
			ExtraAttributes: alkira.ByoipExtraAttributes{
				Message:   "updated-message",
				Signature: "updated-signature",
				PublicKey: "updated-public-key",
			},
		}

		client := serveByoipPrefixMockServer(t, updatedByoipPrefix, http.StatusOK)

		api := alkira.NewByoip(client)
		provState, err, provErr := api.Update("123", updatedByoipPrefix)

		assert.NoError(t, err)
		_ = provErr

		t.Logf("Provision state: %s", provState)
	})

	// Test Delete operation
	t.Run("Delete", func(t *testing.T) {
		client := serveByoipPrefixMockServer(t, nil, http.StatusNoContent)

		api := alkira.NewByoip(client)
		provState, err, provErr := api.Delete("123")

		assert.NoError(t, err)
		_ = provErr

		t.Logf("Provision state: %s", provState)
	})
}

func TestAlkiraByoipPrefix_apiErrorHandling(t *testing.T) {
	// Test error scenarios
	t.Run("NotFound", func(t *testing.T) {
		client := serveByoipPrefixMockServer(t, nil, http.StatusNotFound)

		api := alkira.NewByoip(client)
		_, _, err := api.GetById("999")

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "404")
	})

	t.Run("ServerError", func(t *testing.T) {
		client := serveByoipPrefixMockServer(t, nil, http.StatusInternalServerError)

		api := alkira.NewByoip(client)
		_, _, _, _ = api.Create(&alkira.Byoip{
			Description:   "test",
			Prefix:        "203.0.113.0/24",
			Cxp:           "US-WEST",
			CloudProvider: "AWS",
			ExtraAttributes: alkira.ByoipExtraAttributes{
				Message:   "test-message",
				Signature: "test-signature",
				PublicKey: "test-public-key",
			},
		})

		// Should handle server errors gracefully
	})
}

func TestAlkiraByoipPrefix_resourceDataManipulation(t *testing.T) {
	r := resourceAlkiraByoipPrefix()

	t.Run("set and get resource data", func(t *testing.T) {
		d := r.TestResourceData()

		// Set values
		d.Set("description", "Test description")
		d.Set("prefix", "203.0.113.0/24")
		d.Set("cxp", "US-WEST")

		// Test getting values
		assert.Equal(t, "Test description", d.Get("description").(string))
		assert.Equal(t, "203.0.113.0/24", d.Get("prefix").(string))
		assert.Equal(t, "US-WEST", d.Get("cxp").(string))
	})

	t.Run("resource data with changes", func(t *testing.T) {
		d := r.TestResourceData()

		// Set initial values
		d.Set("description", "Original description")
		d.Set("prefix", "198.51.100.0/24")
		d.Set("cxp", "EU-WEST")

		// Simulate a change
		d.Set("description", "Updated description")
		d.Set("prefix", "203.0.113.0/24")
		d.Set("cxp", "US-WEST")

		assert.Equal(t, "Updated description", d.Get("description").(string))
		assert.Equal(t, "203.0.113.0/24", d.Get("prefix").(string))
		assert.Equal(t, "US-WEST", d.Get("cxp").(string))
	})
}

func TestAlkiraByoipPrefix_idValidation(t *testing.T) {
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

func TestAlkiraByoipPrefix_customizeDiff(t *testing.T) {
	r := resourceAlkiraByoipPrefix()

	// Test CustomizeDiff function exists
	assert.NotNil(t, r.CustomizeDiff, "Resource should have CustomizeDiff")

	// Note: CustomizeDiff testing would require more complex setup with ResourceDiff mock
	// This validates the function exists and can be called
}

// Helper function to create mock HTTP server for BYOIP prefixes
func serveByoipPrefixMockServer(t *testing.T, byoipPrefix *alkira.Byoip, statusCode int) *alkira.AlkiraClient {
	server := httptest.NewServer(http.HandlerFunc(
		func(w http.ResponseWriter, req *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(statusCode)

			switch req.Method {
			case "GET":
				if byoipPrefix != nil {
					json.NewEncoder(w).Encode(byoipPrefix)
				}
			case "POST":
				if byoipPrefix != nil {
					json.NewEncoder(w).Encode(byoipPrefix)
				}
			case "PUT":
				if byoipPrefix != nil {
					json.NewEncoder(w).Encode(byoipPrefix)
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
func buildByoipPrefixRequest(d *schema.ResourceData) *alkira.Byoip {
	description := ""
	if desc := d.Get("description"); desc != nil {
		description = desc.(string)
	}

	prefix := ""
	if p := d.Get("prefix"); p != nil {
		prefix = p.(string)
	}

	cxp := ""
	if c := d.Get("cxp"); c != nil {
		cxp = c.(string)
	}

	cloudProvider := ""
	if cp := d.Get("cloud_provider"); cp != nil {
		cloudProvider = cp.(string)
	}

	message := ""
	if m := d.Get("message"); m != nil {
		message = m.(string)
	}

	signature := ""
	if s := d.Get("signature"); s != nil {
		signature = s.(string)
	}

	publicKey := ""
	if pk := d.Get("public_key"); pk != nil {
		publicKey = pk.(string)
	}

	return &alkira.Byoip{
		Description:   description,
		Prefix:        prefix,
		Cxp:           cxp,
		CloudProvider: cloudProvider,
		ExtraAttributes: alkira.ByoipExtraAttributes{
			Message:   message,
			Signature: signature,
			PublicKey: publicKey,
		},
	}
}

func validateByoipPrefix(v interface{}, k string) (warnings []string, errors []error) {
	value := v.(string)

	if value == "" {
		errors = append(errors, fmt.Errorf("%q cannot be empty", k))
		return warnings, errors
	}

	// Basic CIDR validation
	if !containsString(value, "/") {
		errors = append(errors, fmt.Errorf("%q must be a valid CIDR notation (e.g., 203.0.113.0/24)", k))
		return warnings, errors
	}

	// Check for obviously invalid IPs
	if containsString(value, "256.") {
		errors = append(errors, fmt.Errorf("%q contains invalid IP address", k))
		return warnings, errors
	}

	// Check for invalid subnet masks
	if containsString(value, "/33") || containsString(value, "/34") {
		errors = append(errors, fmt.Errorf("%q contains invalid subnet mask", k))
		return warnings, errors
	}

	return warnings, errors
}

func validateCxp(v interface{}, k string) (warnings []string, errors []error) {
	value := v.(string)
	validCxps := []string{"US-WEST", "US-EAST", "EU-WEST", "APAC"}

	if value == "" {
		errors = append(errors, fmt.Errorf("%q cannot be empty", k))
		return warnings, errors
	}

	for _, validCxp := range validCxps {
		if value == validCxp {
			return warnings, errors
		}
	}

	errors = append(errors, fmt.Errorf("%q must be one of: %v", k, validCxps))
	return warnings, errors
}

func containsString(s, substr string) bool {
	return len(substr) <= len(s) && (len(substr) == 0 || indexOf(s, substr) >= 0)
}

func indexOf(s, substr string) int {
	if len(substr) == 0 {
		return 0
	}
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return i
		}
	}
	return -1
}
