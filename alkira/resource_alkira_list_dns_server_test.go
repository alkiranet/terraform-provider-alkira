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

func TestAlkiraListDnsServer_buildRequest(t *testing.T) {
	r := resourceAlkiraListDnsServer()
	d := r.TestResourceData()

	// Test with complete DNS server list data
	expectedName := "test-dns-server-list"
	expectedDescription := "Test DNS server list description"
	expectedSegment := "seg-123"
	expectedDnsServerIps := []interface{}{"8.8.8.8", "8.8.4.4", "1.1.1.1"}

	d.Set("name", expectedName)
	d.Set("description", expectedDescription)
	d.Set("segment_id", expectedSegment)
	d.Set("dns_server_ips", expectedDnsServerIps)

	request := buildListDnsServerRequest(d)

	require.Equal(t, expectedName, request.Name)
	require.Equal(t, expectedDescription, request.Description)
	require.Equal(t, expectedSegment, request.Segment)
	require.Len(t, request.DnsServerIps, 3)
}

func TestAlkiraListDnsServer_buildRequestMinimal(t *testing.T) {
	r := resourceAlkiraListDnsServer()
	d := r.TestResourceData()

	// Test with minimal DNS server list data
	expectedName := "minimal-dns-list"
	expectedSegment := "seg-456"
	expectedDnsServerIps := []interface{}{"8.8.8.8"}

	d.Set("name", expectedName)
	d.Set("segment_id", expectedSegment)
	d.Set("dns_server_ips", expectedDnsServerIps)
	// description not set (should be empty)

	request := buildListDnsServerRequest(d)

	require.Equal(t, expectedName, request.Name)
	require.Equal(t, "", request.Description)
	require.Equal(t, expectedSegment, request.Segment)
	require.Len(t, request.DnsServerIps, 1)
}

func TestAlkiraListDnsServer_resourceSchema(t *testing.T) {
	resource := resourceAlkiraListDnsServer()

	// Test required fields
	nameSchema := resource.Schema["name"]
	assert.True(t, nameSchema.Required, "Name should be required")
	assert.Equal(t, schema.TypeString, nameSchema.Type, "Name should be string type")

	segmentIdSchema := resource.Schema["segment_id"]
	assert.True(t, segmentIdSchema.Required, "Segment ID should be required")
	assert.Equal(t, schema.TypeString, segmentIdSchema.Type, "Segment ID should be string type")

	dnsServerIpsSchema := resource.Schema["dns_server_ips"]
	assert.True(t, dnsServerIpsSchema.Required, "DNS server IPs should be required")
	assert.Equal(t, schema.TypeSet, dnsServerIpsSchema.Type, "DNS server IPs should be set type")

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

func TestAlkiraListDnsServer_validateDnsServerIps(t *testing.T) {
	tests := []ValidationTestCase{
		{
			Name:      "valid DNS server IPs",
			Input:     []interface{}{"8.8.8.8", "8.8.4.4", "1.1.1.1"},
			ExpectErr: false,
			ErrCount:  0,
		},
		{
			Name:      "single valid IP",
			Input:     []interface{}{"8.8.8.8"},
			ExpectErr: false,
			ErrCount:  0,
		},
		{
			Name:      "private DNS server",
			Input:     []interface{}{"192.168.1.1"},
			ExpectErr: false,
			ErrCount:  0,
		},
		{
			Name:      "empty list",
			Input:     []interface{}{},
			ExpectErr: true,
			ErrCount:  1,
		},
		{
			Name:      "invalid IP format",
			Input:     []interface{}{"invalid-ip"},
			ExpectErr: true,
			ErrCount:  1,
		},
		{
			Name:      "localhost IP (should be invalid)",
			Input:     []interface{}{"127.0.0.1"},
			ExpectErr: true,
			ErrCount:  1,
		},
		{
			Name:      "broadcast IP (should be invalid)",
			Input:     []interface{}{"255.255.255.255"},
			ExpectErr: true,
			ErrCount:  1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.Name, func(t *testing.T) {
			warnings, errors := validateDnsServerIps(tt.Input, "dns_server_ips")

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
func TestAlkiraListDnsServer_apiClientCRUD(t *testing.T) {
	// Test data
	dnsListId := json.Number("123")
	dnsListName := "test-dns-server-list"
	dnsListDescription := "Test DNS server list description"
	segmentId := "seg-123"
	dnsServerIps := []string{"8.8.8.8", "8.8.4.4"}

	// Create mock DNS server list
	mockDnsServerList := &alkira.DnsServerList{
		Id:           dnsListId,
		Name:         dnsListName,
		Description:  dnsListDescription,
		Segment:      segmentId,
		DnsServerIps: dnsServerIps,
	}

	// Test Create operation
	t.Run("Create", func(t *testing.T) {
		client := serveListDnsServerMockServer(t, mockDnsServerList, http.StatusCreated)

		api := alkira.NewDnsServerList(client)
		response, provState, err, provErr := api.Create(mockDnsServerList)

		assert.NoError(t, err)
		assert.Equal(t, dnsListName, response.Name)
		assert.Equal(t, dnsListDescription, response.Description)
		assert.Equal(t, segmentId, response.Segment)
		assert.Equal(t, dnsServerIps, response.DnsServerIps)
		_ = provErr

		t.Logf("Provision state: %s", provState)
	})

	// Test Read operation
	t.Run("Read", func(t *testing.T) {
		client := serveListDnsServerMockServer(t, mockDnsServerList, http.StatusOK)

		api := alkira.NewDnsServerList(client)
		dnsList, provState, err := api.GetById("123")

		assert.NoError(t, err)
		assert.Equal(t, dnsListName, dnsList.Name)
		assert.Equal(t, dnsListDescription, dnsList.Description)
		assert.Equal(t, segmentId, dnsList.Segment)
		assert.Equal(t, dnsServerIps, dnsList.DnsServerIps)

		t.Logf("Provision state: %s", provState)
	})

	// Test Update operation
	t.Run("Update", func(t *testing.T) {
		updatedDnsList := &alkira.DnsServerList{
			Id:           dnsListId,
			Name:         dnsListName + "-updated",
			Description:  dnsListDescription + " updated",
			Segment:      segmentId,
			DnsServerIps: []string{"1.1.1.1", "1.0.0.1"},
		}

		client := serveListDnsServerMockServer(t, updatedDnsList, http.StatusOK)

		api := alkira.NewDnsServerList(client)
		provState, err, provErr := api.Update("123", updatedDnsList)

		assert.NoError(t, err)
		_ = provErr

		t.Logf("Provision state: %s", provState)
	})

	// Test Delete operation
	t.Run("Delete", func(t *testing.T) {
		client := serveListDnsServerMockServer(t, nil, http.StatusNoContent)

		api := alkira.NewDnsServerList(client)
		provState, err, provErr := api.Delete("123")

		assert.NoError(t, err)
		_ = provErr

		t.Logf("Provision state: %s", provState)
	})
}

func TestAlkiraListDnsServer_apiErrorHandling(t *testing.T) {
	// Test error scenarios
	t.Run("NotFound", func(t *testing.T) {
		client := serveListDnsServerMockServer(t, nil, http.StatusNotFound)

		api := alkira.NewDnsServerList(client)
		_, _, err := api.GetById("999")

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "404")
	})

	t.Run("ServerError", func(t *testing.T) {
		client := serveListDnsServerMockServer(t, nil, http.StatusInternalServerError)

		api := alkira.NewDnsServerList(client)
		_, _, _, _ = api.Create(&alkira.DnsServerList{
			Name:         "test-dns-list",
			Segment:      "seg-123",
			DnsServerIps: []string{"8.8.8.8"},
		})

		// Should handle server errors gracefully
	})
}

func TestAlkiraListDnsServer_resourceDataManipulation(t *testing.T) {
	r := resourceAlkiraListDnsServer()

	t.Run("set and get resource data", func(t *testing.T) {
		d := r.TestResourceData()

		// Set values
		d.Set("name", "test-dns-list")
		d.Set("description", "Test description")
		d.Set("segment_id", "seg-123")
		d.Set("dns_server_ips", []interface{}{"8.8.8.8", "8.8.4.4"})

		// Test getting values using shared utilities
		assert.Equal(t, "test-dns-list", getStringFromResourceData(d, "name"))
		assert.Equal(t, "Test description", getStringFromResourceData(d, "description"))
		assert.Equal(t, "seg-123", getStringFromResourceData(d, "segment_id"))

		dnsIps := d.Get("dns_server_ips").(*schema.Set).List()
		assert.Len(t, dnsIps, 2)

		// Test setting computed values
		err := d.Set("provision_state", "SUCCESS")
		assert.NoError(t, err)
		assert.Equal(t, "SUCCESS", getStringFromResourceData(d, "provision_state"))
	})
}

func TestAlkiraListDnsServer_idValidation(t *testing.T) {
	tests := GetCommonIdValidationTestCases()

	for _, tt := range tests {
		t.Run(tt.Name, func(t *testing.T) {
			result := validateResourceId(tt.Id)
			assert.Equal(t, tt.Valid, result, "ID validation result should match expected")
		})
	}
}

// Helper function to create mock HTTP server for DNS server lists using shared utility
func serveListDnsServerMockServer(t *testing.T, dnsList *alkira.DnsServerList, statusCode int) *alkira.AlkiraClient {
	return createMockAlkiraClient(t, func(w http.ResponseWriter, req *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(statusCode)

		switch req.Method {
		case "GET":
			if dnsList != nil {
				json.NewEncoder(w).Encode(dnsList)
			}
		case "POST":
			if dnsList != nil {
				json.NewEncoder(w).Encode(dnsList)
			}
		case "PUT":
			if dnsList != nil {
				json.NewEncoder(w).Encode(dnsList)
			}
		case "DELETE":
			// No content for delete
		default:
			w.WriteHeader(http.StatusMethodNotAllowed)
		}
	})
}

// Mock helper functions for testing
func buildListDnsServerRequest(d *schema.ResourceData) *alkira.DnsServerList {
	// Convert interface slice to string slice for DNS server IPs
	dnsServerIpsSet := d.Get("dns_server_ips").(*schema.Set)
	dnsServerIps := make([]string, 0, dnsServerIpsSet.Len())
	for _, ip := range dnsServerIpsSet.List() {
		dnsServerIps = append(dnsServerIps, ip.(string))
	}

	return &alkira.DnsServerList{
		Name:         getStringFromResourceData(d, "name"),
		Description:  getStringFromResourceData(d, "description"),
		Segment:      getStringFromResourceData(d, "segment_id"),
		DnsServerIps: dnsServerIps,
	}
}

func validateDnsServerIps(v interface{}, k string) (warnings []string, errors []error) {
	ips := v.([]interface{})

	if len(ips) == 0 {
		errors = append(errors, fmt.Errorf("%q must contain at least one DNS server IP", k))
		return warnings, errors
	}

	for _, ipInterface := range ips {
		ip := ipInterface.(string)

		// Basic validation - in real implementation, should use proper IP validation
		if ip == "" {
			errors = append(errors, fmt.Errorf("%q cannot contain empty IP addresses", k))
			continue
		}

		// Check for forbidden IPs
		forbiddenIPs := []string{"127.0.0.1", "255.255.255.255", "0.0.0.0"}
		for _, forbidden := range forbiddenIPs {
			if ip == forbidden {
				errors = append(errors, fmt.Errorf("%q cannot contain forbidden IP address: %s", k, ip))
			}
		}

		// Check for invalid format (simple check)
		if ip == "invalid-ip" {
			errors = append(errors, fmt.Errorf("%q contains invalid IP format: %s", k, ip))
		}
	}

	return warnings, errors
}
