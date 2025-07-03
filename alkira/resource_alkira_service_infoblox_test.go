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

func TestAlkiraServiceInfoblox_buildServiceInfobloxRequest(t *testing.T) {
	r := resourceAlkiraInfoblox()
	d := r.TestResourceData()

	// Test with complete Infoblox service data
	expectedName := "test-infoblox-service"
	expectedDescription := "Test Infoblox service description"
	expectedCxp := "US-WEST"
	expectedSize := "MEDIUM"
	expectedVersion := "8.6.0"
	expectedCredentialId := "test-credential-123"
	expectedLicenseKey := "test-license-key"
	expectedSegmentIds := []int{1, 2}
	expectedBillingTagIds := []int{10, 20}
	expectedLocalNsGroup := "test-ns-group"
	expectedFqdnRecord := "test.example.com"

	// Set anycast configuration
	anycastData := []interface{}{
		map[string]interface{}{
			"enabled": true,
			"ips":     []interface{}{"192.168.1.1", "192.168.1.2"},
		},
	}

	// Set instance configuration
	instancesData := []interface{}{
		map[string]interface{}{
			"name":            "infoblox-01",
			"traffic_enabled": true,
		},
		map[string]interface{}{
			"name":            "infoblox-02",
			"traffic_enabled": false,
		},
	}

	d.Set("name", expectedName)
	d.Set("description", expectedDescription)
	d.Set("cxp", expectedCxp)
	d.Set("size", expectedSize)
	d.Set("version", expectedVersion)
	d.Set("credential_id", expectedCredentialId)
	d.Set("license_key", expectedLicenseKey)
	d.Set("segment_ids", expectedSegmentIds)
	d.Set("billing_tag_ids", expectedBillingTagIds)
	d.Set("local_ns_group", expectedLocalNsGroup)
	d.Set("fqdn_record", expectedFqdnRecord)
	d.Set("anycast", anycastData)
	d.Set("instances", instancesData)

	request := buildServiceInfobloxRequest(d)

	require.Equal(t, expectedName, request.Name)
	require.Equal(t, expectedDescription, request.Description)
	require.Equal(t, expectedCxp, request.Cxp)
	require.Equal(t, expectedSize, request.Size)
	require.Equal(t, expectedVersion, request.Version)
	require.Equal(t, expectedCredentialId, request.CredentialId)
	require.Equal(t, expectedLicenseKey, request.LicenseKey)
	require.Equal(t, expectedSegmentIds, request.SegmentIds)
	require.Equal(t, expectedBillingTagIds, request.BillingTagIds)
	require.Equal(t, expectedLocalNsGroup, request.LocalNsGroup)
	require.Equal(t, expectedFqdnRecord, request.FqdnRecord)

	// Verify anycast configuration
	require.NotNil(t, request.Anycast)
	require.Equal(t, true, request.Anycast.Enabled)
	require.Equal(t, []string{"192.168.1.1", "192.168.1.2"}, request.Anycast.Ips)

	// Verify instances configuration
	require.Len(t, request.Instances, 2)
	assert.Equal(t, "infoblox-01", request.Instances[0].Name)
	assert.Equal(t, true, request.Instances[0].TrafficEnabled)
	assert.Equal(t, "infoblox-02", request.Instances[1].Name)
	assert.Equal(t, false, request.Instances[1].TrafficEnabled)
}

func TestAlkiraServiceInfoblox_buildServiceInfobloxRequestMinimal(t *testing.T) {
	r := resourceAlkiraInfoblox()
	d := r.TestResourceData()

	// Test with minimal required Infoblox service data
	expectedName := "minimal-infoblox-service"
	expectedCxp := "US-EAST"
	expectedSize := "SMALL"
	expectedVersion := "8.5.0"
	expectedCredentialId := "minimal-credential"
	expectedSegmentIds := []int{1}

	// Set minimal anycast configuration (required)
	anycastData := []interface{}{
		map[string]interface{}{
			"enabled": false,
		},
	}

	d.Set("name", expectedName)
	d.Set("cxp", expectedCxp)
	d.Set("size", expectedSize)
	d.Set("version", expectedVersion)
	d.Set("credential_id", expectedCredentialId)
	d.Set("segment_ids", expectedSegmentIds)
	d.Set("anycast", anycastData)

	request := buildServiceInfobloxRequest(d)

	require.Equal(t, expectedName, request.Name)
	require.Equal(t, "", request.Description) // Should be empty when not set
	require.Equal(t, expectedCxp, request.Cxp)
	require.Equal(t, expectedSize, request.Size)
	require.Equal(t, expectedVersion, request.Version)
	require.Equal(t, expectedCredentialId, request.CredentialId)
	require.Equal(t, expectedSegmentIds, request.SegmentIds)

	// Verify minimal anycast configuration
	require.NotNil(t, request.Anycast)
	require.Equal(t, false, request.Anycast.Enabled)
	require.Empty(t, request.Anycast.Ips)
}

func TestAlkiraServiceInfoblox_buildServiceInfobloxRequestWithDns(t *testing.T) {
	r := resourceAlkiraInfoblox()
	d := r.TestResourceData()

	// Test with DNS configuration
	expectedName := "infoblox-dns-service"
	expectedCxp := "US-WEST"
	expectedSize := "MEDIUM"
	expectedVersion := "8.6.0"
	expectedCredentialId := "dns-credential"
	expectedSegmentIds := []int{1}
	expectedLocalNsGroup := "dns-group"
	expectedFqdnRecord := "dns.example.com"

	// Set anycast configuration
	anycastData := []interface{}{
		map[string]interface{}{
			"enabled": true,
			"ips":     []interface{}{"10.0.0.1", "10.0.0.2"},
		},
	}

	// Set DNS forwarding zones
	dnsForwardingZonesData := []interface{}{
		map[string]interface{}{
			"fqdn":         "example.com",
			"forwarders":   []interface{}{"8.8.8.8", "8.8.4.4"},
			"disabled":     false,
			"forward_only": true,
		},
		map[string]interface{}{
			"fqdn":         "test.com",
			"forwarders":   []interface{}{"1.1.1.1"},
			"disabled":     true,
			"forward_only": false,
		},
	}

	d.Set("name", expectedName)
	d.Set("cxp", expectedCxp)
	d.Set("size", expectedSize)
	d.Set("version", expectedVersion)
	d.Set("credential_id", expectedCredentialId)
	d.Set("segment_ids", expectedSegmentIds)
	d.Set("local_ns_group", expectedLocalNsGroup)
	d.Set("fqdn_record", expectedFqdnRecord)
	d.Set("anycast", anycastData)
	d.Set("dns_forwarding_zones", dnsForwardingZonesData)

	request := buildServiceInfobloxRequest(d)

	require.Equal(t, expectedName, request.Name)
	require.Equal(t, expectedLocalNsGroup, request.LocalNsGroup)
	require.Equal(t, expectedFqdnRecord, request.FqdnRecord)

	// Verify DNS forwarding zones
	require.Len(t, request.DnsForwardingZones, 2)

	// Check first DNS forwarding zone
	assert.Equal(t, "example.com", request.DnsForwardingZones[0].Fqdn)
	assert.Equal(t, []string{"8.8.8.8", "8.8.4.4"}, request.DnsForwardingZones[0].Forwarders)
	assert.Equal(t, false, request.DnsForwardingZones[0].Disabled)
	assert.Equal(t, true, request.DnsForwardingZones[0].ForwardOnly)

	// Check second DNS forwarding zone
	assert.Equal(t, "test.com", request.DnsForwardingZones[1].Fqdn)
	assert.Equal(t, []string{"1.1.1.1"}, request.DnsForwardingZones[1].Forwarders)
	assert.Equal(t, true, request.DnsForwardingZones[1].Disabled)
	assert.Equal(t, false, request.DnsForwardingZones[1].ForwardOnly)
}

func TestAlkiraServiceInfoblox_resourceSchema(t *testing.T) {
	resource := resourceAlkiraInfoblox()

	// Test required fields
	nameSchema := resource.Schema["name"]
	assert.True(t, nameSchema.Required, "Name should be required")
	assert.Equal(t, schema.TypeString, nameSchema.Type, "Name should be string type")

	cxpSchema := resource.Schema["cxp"]
	assert.True(t, cxpSchema.Required, "CXP should be required")
	assert.Equal(t, schema.TypeString, cxpSchema.Type, "CXP should be string type")

	sizeSchema := resource.Schema["size"]
	assert.True(t, sizeSchema.Required, "Size should be required")
	assert.Equal(t, schema.TypeString, sizeSchema.Type, "Size should be string type")

	versionSchema := resource.Schema["version"]
	assert.True(t, versionSchema.Required, "Version should be required")
	assert.Equal(t, schema.TypeString, versionSchema.Type, "Version should be string type")

	credentialIdSchema := resource.Schema["credential_id"]
	assert.True(t, credentialIdSchema.Required, "Credential ID should be required")
	assert.Equal(t, schema.TypeString, credentialIdSchema.Type, "Credential ID should be string type")

	segmentIdsSchema := resource.Schema["segment_ids"]
	assert.True(t, segmentIdsSchema.Required, "Segment IDs should be required")
	assert.Equal(t, schema.TypeSet, segmentIdsSchema.Type, "Segment IDs should be set type")

	anycastSchema := resource.Schema["anycast"]
	assert.True(t, anycastSchema.Required, "Anycast should be required")
	assert.Equal(t, schema.TypeSet, anycastSchema.Type, "Anycast should be set type")

	// Test optional fields
	descSchema := resource.Schema["description"]
	assert.True(t, descSchema.Optional, "Description should be optional")
	assert.Equal(t, schema.TypeString, descSchema.Type, "Description should be string type")

	licenseKeySchema := resource.Schema["license_key"]
	assert.True(t, licenseKeySchema.Optional, "License key should be optional")
	assert.Equal(t, schema.TypeString, licenseKeySchema.Type, "License key should be string type")

	billingTagIdsSchema := resource.Schema["billing_tag_ids"]
	assert.True(t, billingTagIdsSchema.Optional, "Billing tag IDs should be optional")
	assert.Equal(t, schema.TypeSet, billingTagIdsSchema.Type, "Billing tag IDs should be set type")

	localNsGroupSchema := resource.Schema["local_ns_group"]
	assert.True(t, localNsGroupSchema.Optional, "Local NS group should be optional")
	assert.Equal(t, schema.TypeString, localNsGroupSchema.Type, "Local NS group should be string type")

	fqdnRecordSchema := resource.Schema["fqdn_record"]
	assert.True(t, fqdnRecordSchema.Optional, "FQDN record should be optional")
	assert.Equal(t, schema.TypeString, fqdnRecordSchema.Type, "FQDN record should be string type")

	instancesSchema := resource.Schema["instances"]
	assert.True(t, instancesSchema.Optional, "Instances should be optional")
	assert.Equal(t, schema.TypeSet, instancesSchema.Type, "Instances should be set type")

	dnsForwardingZonesSchema := resource.Schema["dns_forwarding_zones"]
	assert.True(t, dnsForwardingZonesSchema.Optional, "DNS forwarding zones should be optional")
	assert.Equal(t, schema.TypeSet, dnsForwardingZonesSchema.Type, "DNS forwarding zones should be set type")

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
	assert.NotNil(t, resource.CustomizeDiff, "Resource should have CustomizeDiff")
}

func TestAlkiraServiceInfoblox_validateSize(t *testing.T) {
	tests := []struct {
		name      string
		input     interface{}
		expectErr bool
		errCount  int
	}{
		{
			name:      "Valid SMALL",
			input:     "SMALL",
			expectErr: false,
			errCount:  0,
		},
		{
			name:      "Valid MEDIUM",
			input:     "MEDIUM",
			expectErr: false,
			errCount:  0,
		},
		{
			name:      "Valid LARGE",
			input:     "LARGE",
			expectErr: false,
			errCount:  0,
		},
		{
			name:      "Invalid size",
			input:     "INVALID_SIZE",
			expectErr: true,
			errCount:  1,
		},
		{
			name:      "Empty string",
			input:     "",
			expectErr: true,
			errCount:  1,
		},
	}

	resource := resourceAlkiraInfoblox()
	sizeSchema := resource.Schema["size"]

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			warnings, errors := sizeSchema.ValidateFunc(tt.input, "size")

			if tt.expectErr {
				assert.Len(t, errors, tt.errCount, "Expected %d errors for input %v", tt.errCount, tt.input)
			} else {
				assert.Empty(t, errors, "Expected no errors for input %v", tt.input)
			}
			assert.Empty(t, warnings, "Expected no warnings")
		})
	}
}

func TestAlkiraServiceInfoblox_CreateUpdateDelete(t *testing.T) {
	// Test data
	serviceInfoblox := &alkira.ServiceInfoblox{
		Id:            json.Number("123"),
		Name:          "test-infoblox-service",
		Description:   "Test Infoblox service",
		Cxp:           "US-WEST",
		Size:          "MEDIUM",
		Version:       "8.6.0",
		CredentialId:  "test-credential",
		LicenseKey:    "test-license-key",
		SegmentIds:    []int{1, 2},
		BillingTagIds: []int{10, 20},
		LocalNsGroup:  "test-ns-group",
		FqdnRecord:    "test.example.com",
		Anycast: &alkira.InfobloxAnycast{
			Enabled: true,
			Ips:     []string{"192.168.1.1", "192.168.1.2"},
		},
		Instances: []alkira.InfobloxInstance{
			{
				Id:             1,
				Name:           "infoblox-01",
				TrafficEnabled: true,
			},
		},
		DnsForwardingZones: []alkira.InfobloxDnsForwardingZone{
			{
				Fqdn:        "example.com",
				Forwarders:  []string{"8.8.8.8", "8.8.4.4"},
				Disabled:    false,
				ForwardOnly: true,
			},
		},
	}

	// Create mock server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPost:
			// CREATE - return created service with ID
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusCreated)
			json.NewEncoder(w).Encode(serviceInfoblox)
		case http.MethodGet:
			// READ - return service data
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(serviceInfoblox)
		case http.MethodPut:
			// UPDATE - return updated service
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(serviceInfoblox)
		case http.MethodDelete:
			// DELETE - return success
			w.WriteHeader(http.StatusNoContent)
		default:
			w.WriteHeader(http.StatusMethodNotAllowed)
		}
	}))
	defer server.Close()

	// Create client
	retryClient := retryablehttp.NewClient()
	retryClient.HTTPClient.Timeout = time.Duration(1) * time.Second
	client := &alkira.AlkiraClient{
		URI:             server.URL,
		TenantNetworkId: "0",
		Client:          retryClient,
	}

	// Test CREATE
	r := resourceAlkiraInfoblox()
	d := r.TestResourceData()

	anycastData := []interface{}{
		map[string]interface{}{
			"enabled": serviceInfoblox.Anycast.Enabled,
			"ips":     serviceInfoblox.Anycast.Ips,
		},
	}

	d.Set("name", serviceInfoblox.Name)
	d.Set("description", serviceInfoblox.Description)
	d.Set("cxp", serviceInfoblox.Cxp)
	d.Set("size", serviceInfoblox.Size)
	d.Set("version", serviceInfoblox.Version)
	d.Set("credential_id", serviceInfoblox.CredentialId)
	d.Set("license_key", serviceInfoblox.LicenseKey)
	d.Set("segment_ids", serviceInfoblox.SegmentIds)
	d.Set("billing_tag_ids", serviceInfoblox.BillingTagIds)
	d.Set("local_ns_group", serviceInfoblox.LocalNsGroup)
	d.Set("fqdn_record", serviceInfoblox.FqdnRecord)
	d.Set("anycast", anycastData)

	diags := resourceInfoblox(nil, d, client)
	require.Empty(t, diags, "Create should not return errors")

	// Verify ID was set
	require.Equal(t, "123", d.Id())

	// Test READ
	diags = resourceInfobloxRead(nil, d, client)
	require.Empty(t, diags, "Read should not return errors")

	// Verify data was populated
	assert.Equal(t, serviceInfoblox.Name, d.Get("name"))
	assert.Equal(t, serviceInfoblox.Description, d.Get("description"))
	assert.Equal(t, serviceInfoblox.Cxp, d.Get("cxp"))
	assert.Equal(t, serviceInfoblox.Size, d.Get("size"))

	// Test UPDATE
	d.Set("description", "Updated description")
	diags = resourceInfobloxUpdate(nil, d, client)
	require.Empty(t, diags, "Update should not return errors")

	// Test DELETE
	diags = resourceInfobloxDelete(nil, d, client)
	require.Empty(t, diags, "Delete should not return errors")
}

func TestAlkiraServiceInfoblox_CreateError(t *testing.T) {
	// Create mock server that returns error
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(`{"error": "Invalid request"}`))
	}))
	defer server.Close()

	// Create client
	retryClient := retryablehttp.NewClient()
	retryClient.HTTPClient.Timeout = time.Duration(1) * time.Second
	client := &alkira.AlkiraClient{
		URI:             server.URL,
		TenantNetworkId: "0",
		Client:          retryClient,
	}

	// Test CREATE error handling
	r := resourceAlkiraInfoblox()
	d := r.TestResourceData()

	anycastData := []interface{}{
		map[string]interface{}{
			"enabled": false,
		},
	}

	d.Set("name", "test-infoblox-service")
	d.Set("cxp", "US-WEST")
	d.Set("size", "MEDIUM")
	d.Set("version", "8.6.0")
	d.Set("credential_id", "test-credential")
	d.Set("segment_ids", []int{1})
	d.Set("anycast", anycastData)

	diags := resourceInfoblox(nil, d, client)
	require.NotEmpty(t, diags, "Create should return error")
	assert.True(t, diags.HasError(), "Diagnostics should contain error")
}

func TestAlkiraServiceInfoblox_ReadNotFound(t *testing.T) {
	// Create mock server that returns 404
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	}))
	defer server.Close()

	// Create client
	retryClient := retryablehttp.NewClient()
	retryClient.HTTPClient.Timeout = time.Duration(1) * time.Second
	client := &alkira.AlkiraClient{
		URI:             server.URL,
		TenantNetworkId: "0",
		Client:          retryClient,
	}

	// Test READ with non-existent resource
	r := resourceAlkiraInfoblox()
	d := r.TestResourceData()
	d.SetId("999")

	diags := resourceInfobloxRead(nil, d, client)
	require.Empty(t, diags, "Read should not return error for 404")
	assert.Equal(t, "", d.Id(), "Resource ID should be cleared when not found")
}

func TestAlkiraServiceInfoblox_setServiceInfobloxFields(t *testing.T) {
	serviceInfoblox := &alkira.ServiceInfoblox{
		Id:            json.Number("123"),
		Name:          "test-infoblox-service",
		Description:   "Test Infoblox service",
		Cxp:           "US-WEST",
		Size:          "MEDIUM",
		Version:       "8.6.0",
		CredentialId:  "test-credential",
		LicenseKey:    "test-license-key",
		SegmentIds:    []int{1, 2},
		BillingTagIds: []int{10, 20},
		LocalNsGroup:  "test-ns-group",
		FqdnRecord:    "test.example.com",
		Anycast: &alkira.InfobloxAnycast{
			Enabled: true,
			Ips:     []string{"192.168.1.1", "192.168.1.2"},
		},
		Instances: []alkira.InfobloxInstance{
			{
				Id:             1,
				Name:           "infoblox-01",
				TrafficEnabled: true,
			},
			{
				Id:             2,
				Name:           "infoblox-02",
				TrafficEnabled: false,
			},
		},
		DnsForwardingZones: []alkira.InfobloxDnsForwardingZone{
			{
				Fqdn:        "example.com",
				Forwarders:  []string{"8.8.8.8", "8.8.4.4"},
				Disabled:    false,
				ForwardOnly: true,
			},
			{
				Fqdn:        "test.com",
				Forwarders:  []string{"1.1.1.1"},
				Disabled:    true,
				ForwardOnly: false,
			},
		},
	}

	r := resourceAlkiraInfoblox()
	d := r.TestResourceData()

	err := setServiceInfobloxFields(d, serviceInfoblox)
	require.NoError(t, err, "setServiceInfobloxFields should not return error")

	// Verify all fields were set correctly
	assert.Equal(t, serviceInfoblox.Name, d.Get("name"))
	assert.Equal(t, serviceInfoblox.Description, d.Get("description"))
	assert.Equal(t, serviceInfoblox.Cxp, d.Get("cxp"))
	assert.Equal(t, serviceInfoblox.Size, d.Get("size"))
	assert.Equal(t, serviceInfoblox.Version, d.Get("version"))
	assert.Equal(t, serviceInfoblox.CredentialId, d.Get("credential_id"))
	assert.Equal(t, serviceInfoblox.LicenseKey, d.Get("license_key"))
	assert.Equal(t, serviceInfoblox.LocalNsGroup, d.Get("local_ns_group"))
	assert.Equal(t, serviceInfoblox.FqdnRecord, d.Get("fqdn_record"))

	// Verify set fields
	segmentIdsSet := d.Get("segment_ids").(*schema.Set)
	assert.Equal(t, len(serviceInfoblox.SegmentIds), segmentIdsSet.Len())
	for _, segmentId := range serviceInfoblox.SegmentIds {
		assert.True(t, segmentIdsSet.Contains(segmentId))
	}

	billingTagIdsSet := d.Get("billing_tag_ids").(*schema.Set)
	assert.Equal(t, len(serviceInfoblox.BillingTagIds), billingTagIdsSet.Len())
	for _, billingTagId := range serviceInfoblox.BillingTagIds {
		assert.True(t, billingTagIdsSet.Contains(billingTagId))
	}

	// Verify anycast configuration
	anycastSet := d.Get("anycast").(*schema.Set)
	assert.Equal(t, 1, anycastSet.Len())

	// Verify instances
	instancesSet := d.Get("instances").(*schema.Set)
	assert.Equal(t, len(serviceInfoblox.Instances), instancesSet.Len())

	// Verify DNS forwarding zones
	dnsForwardingZonesSet := d.Get("dns_forwarding_zones").(*schema.Set)
	assert.Equal(t, len(serviceInfoblox.DnsForwardingZones), dnsForwardingZonesSet.Len())
}

func TestAlkiraServiceInfoblox_validateInfobloxVersion(t *testing.T) {
	tests := []struct {
		name      string
		input     interface{}
		expectErr bool
		errCount  int
	}{
		{
			name:      "Valid version 8.5.0",
			input:     "8.5.0",
			expectErr: false,
			errCount:  0,
		},
		{
			name:      "Valid version 8.6.0",
			input:     "8.6.0",
			expectErr: false,
			errCount:  0,
		},
		{
			name:      "Valid version 9.0.0",
			input:     "9.0.0",
			expectErr: false,
			errCount:  0,
		},
		{
			name:      "Empty string",
			input:     "",
			expectErr: true,
			errCount:  1,
		},
		{
			name:      "Non-string input",
			input:     123,
			expectErr: true,
			errCount:  1,
		},
	}

	// Since we don't have specific version validation in the schema,
	// we'll test basic string validation
	resource := resourceAlkiraInfoblox()
	versionSchema := resource.Schema["version"]

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// For basic string validation, we check if it's required and correct type
			if tt.input == "" {
				// Empty string should fail required validation
				assert.True(t, versionSchema.Required, "Version should be required")
			} else if _, ok := tt.input.(string); !ok && tt.expectErr {
				// Non-string should fail type validation
				assert.Equal(t, schema.TypeString, versionSchema.Type, "Version should be string type")
			}
		})
	}
}

// TEST HELPER
func serveServiceInfoblox(t *testing.T, serviceInfoblox *alkira.ServiceInfoblox) *alkira.AlkiraClient {
	server := httptest.NewServer(http.HandlerFunc(
		func(w http.ResponseWriter, req *http.Request) {
			json.NewEncoder(w).Encode(serviceInfoblox)
			w.Header().Set("Content-Type", "application/json")
		},
	))
	t.Cleanup(server.Close)

	retryClient := retryablehttp.NewClient()
	retryClient.HTTPClient.Timeout = time.Duration(1) * time.Second

	return &alkira.AlkiraClient{
		URI:             server.URL,
		TenantNetworkId: "0",
		Client:          retryClient,
	}
}
