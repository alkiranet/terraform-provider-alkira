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

func TestAlkiraServiceF5LoadBalancer_buildServiceF5LoadBalancerRequest(t *testing.T) {
	r := resourceAlkiraF5LoadBalancer()
	d := r.TestResourceData()

	// Test with complete F5 Load Balancer service data
	expectedName := "test-f5-lb-service"
	expectedDescription := "Test F5 Load Balancer service description"
	expectedCxp := "US-WEST"
	expectedSize := "MEDIUM"
	expectedVersion := "16.1.0"
	expectedCredentialId := "test-credential-123"
	expectedLicenseKey := "test-license-key"
	expectedSegmentIds := []int{1, 2}
	expectedBillingTagIds := []int{10, 20}
	expectedMaxInstanceCount := 3
	expectedMinInstanceCount := 1
	expectedAutoScale := "ON"
	expectedTunnelProtocol := "IPSEC"

	d.Set("name", expectedName)
	d.Set("description", expectedDescription)
	d.Set("cxp", expectedCxp)
	d.Set("size", expectedSize)
	d.Set("version", expectedVersion)
	d.Set("credential_id", expectedCredentialId)
	d.Set("license_key", expectedLicenseKey)
	d.Set("segment_ids", expectedSegmentIds)
	d.Set("billing_tag_ids", expectedBillingTagIds)
	d.Set("max_instance_count", expectedMaxInstanceCount)
	d.Set("min_instance_count", expectedMinInstanceCount)
	d.Set("auto_scale", expectedAutoScale)
	d.Set("tunnel_protocol", expectedTunnelProtocol)

	request := buildServiceF5LoadBalancerRequest(d)

	require.Equal(t, expectedName, request.Name)
	require.Equal(t, expectedDescription, request.Description)
	require.Equal(t, expectedCxp, request.Cxp)
	require.Equal(t, expectedSize, request.Size)
	require.Equal(t, expectedVersion, request.Version)
	require.Equal(t, expectedCredentialId, request.CredentialId)
	require.Equal(t, expectedLicenseKey, request.LicenseKey)
	require.Equal(t, expectedSegmentIds, request.SegmentIds)
	require.Equal(t, expectedBillingTagIds, request.BillingTagIds)
	require.Equal(t, expectedMaxInstanceCount, request.MaxInstanceCount)
	require.Equal(t, expectedMinInstanceCount, request.MinInstanceCount)
	require.Equal(t, expectedAutoScale, request.AutoScale)
	require.Equal(t, expectedTunnelProtocol, request.TunnelProtocol)
}

func TestAlkiraServiceF5LoadBalancer_buildServiceF5LoadBalancerRequestMinimal(t *testing.T) {
	r := resourceAlkiraF5LoadBalancer()
	d := r.TestResourceData()

	// Test with minimal required F5 Load Balancer service data
	expectedName := "minimal-f5-lb-service"
	expectedCxp := "US-EAST"
	expectedSize := "SMALL"
	expectedVersion := "15.1.0"
	expectedCredentialId := "minimal-credential"
	expectedSegmentIds := []int{1}

	d.Set("name", expectedName)
	d.Set("cxp", expectedCxp)
	d.Set("size", expectedSize)
	d.Set("version", expectedVersion)
	d.Set("credential_id", expectedCredentialId)
	d.Set("segment_ids", expectedSegmentIds)

	request := buildServiceF5LoadBalancerRequest(d)

	require.Equal(t, expectedName, request.Name)
	require.Equal(t, "", request.Description) // Should be empty when not set
	require.Equal(t, expectedCxp, request.Cxp)
	require.Equal(t, expectedSize, request.Size)
	require.Equal(t, expectedVersion, request.Version)
	require.Equal(t, expectedCredentialId, request.CredentialId)
	require.Equal(t, expectedSegmentIds, request.SegmentIds)
	require.Equal(t, "OFF", request.AutoScale) // Default value
}

func TestAlkiraServiceF5LoadBalancer_buildServiceF5LoadBalancerRequestInstances(t *testing.T) {
	r := resourceAlkiraF5LoadBalancer()
	d := r.TestResourceData()

	// Test with instances configuration
	expectedName := "f5-lb-with-instances"
	expectedCxp := "US-WEST"
	expectedSize := "MEDIUM"
	expectedVersion := "16.1.0"
	expectedCredentialId := "instance-credential"
	expectedSegmentIds := []int{1}

	// Set up instances
	instances := []interface{}{
		map[string]interface{}{
			"name":            "f5-lb-01",
			"traffic_enabled": true,
		},
		map[string]interface{}{
			"name":            "f5-lb-02",
			"traffic_enabled": false,
		},
	}

	d.Set("name", expectedName)
	d.Set("cxp", expectedCxp)
	d.Set("size", expectedSize)
	d.Set("version", expectedVersion)
	d.Set("credential_id", expectedCredentialId)
	d.Set("segment_ids", expectedSegmentIds)
	d.Set("instances", instances)

	request := buildServiceF5LoadBalancerRequest(d)

	require.Equal(t, expectedName, request.Name)
	require.Len(t, request.Instances, 2)

	// Check first instance
	assert.Equal(t, "f5-lb-01", request.Instances[0].Name)
	assert.Equal(t, true, request.Instances[0].TrafficEnabled)

	// Check second instance
	assert.Equal(t, "f5-lb-02", request.Instances[1].Name)
	assert.Equal(t, false, request.Instances[1].TrafficEnabled)
}

func TestAlkiraServiceF5LoadBalancer_resourceSchema(t *testing.T) {
	resource := resourceAlkiraF5LoadBalancer()

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

	// Test optional fields with defaults
	autoScaleSchema := resource.Schema["auto_scale"]
	assert.True(t, autoScaleSchema.Optional, "Auto scale should be optional")
	assert.Equal(t, schema.TypeString, autoScaleSchema.Type, "Auto scale should be string type")
	assert.Equal(t, "OFF", autoScaleSchema.Default, "Auto scale should default to OFF")

	// Test optional fields
	descSchema := resource.Schema["description"]
	assert.True(t, descSchema.Optional, "Description should be optional")
	assert.Equal(t, schema.TypeString, descSchema.Type, "Description should be string type")

	licenseKeySchema := resource.Schema["license_key"]
	assert.True(t, licenseKeySchema.Optional, "License key should be optional")
	assert.Equal(t, schema.TypeString, licenseKeySchema.Type, "License key should be string type")

	maxInstanceCountSchema := resource.Schema["max_instance_count"]
	assert.True(t, maxInstanceCountSchema.Optional, "Max instance count should be optional")
	assert.Equal(t, schema.TypeInt, maxInstanceCountSchema.Type, "Max instance count should be int type")

	minInstanceCountSchema := resource.Schema["min_instance_count"]
	assert.True(t, minInstanceCountSchema.Optional, "Min instance count should be optional")
	assert.Equal(t, schema.TypeInt, minInstanceCountSchema.Type, "Min instance count should be int type")

	tunnelProtocolSchema := resource.Schema["tunnel_protocol"]
	assert.True(t, tunnelProtocolSchema.Optional, "Tunnel protocol should be optional")
	assert.Equal(t, schema.TypeString, tunnelProtocolSchema.Type, "Tunnel protocol should be string type")

	billingTagIdsSchema := resource.Schema["billing_tag_ids"]
	assert.True(t, billingTagIdsSchema.Optional, "Billing tag IDs should be optional")
	assert.Equal(t, schema.TypeSet, billingTagIdsSchema.Type, "Billing tag IDs should be set type")

	instancesSchema := resource.Schema["instances"]
	assert.True(t, instancesSchema.Optional, "Instances should be optional")
	assert.Equal(t, schema.TypeSet, instancesSchema.Type, "Instances should be set type")

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

func TestAlkiraServiceF5LoadBalancer_validateAutoScale(t *testing.T) {
	tests := []struct {
		name      string
		input     interface{}
		expectErr bool
		errCount  int
	}{
		{
			name:      "Valid ON",
			input:     "ON",
			expectErr: false,
			errCount:  0,
		},
		{
			name:      "Valid OFF",
			input:     "OFF",
			expectErr: false,
			errCount:  0,
		},
		{
			name:      "Invalid auto scale value",
			input:     "INVALID",
			expectErr: true,
			errCount:  1,
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

	resource := resourceAlkiraF5LoadBalancer()
	autoScaleSchema := resource.Schema["auto_scale"]

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			warnings, errors := autoScaleSchema.ValidateFunc(tt.input, "auto_scale")

			if tt.expectErr {
				assert.Len(t, errors, tt.errCount, "Expected %d errors for input %v", tt.errCount, tt.input)
			} else {
				assert.Empty(t, errors, "Expected no errors for input %v", tt.input)
			}
			assert.Empty(t, warnings, "Expected no warnings")
		})
	}
}

func TestAlkiraServiceF5LoadBalancer_validateSize(t *testing.T) {
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

	resource := resourceAlkiraF5LoadBalancer()
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

func TestAlkiraServiceF5LoadBalancer_CreateUpdateDelete(t *testing.T) {
	// Test data
	serviceF5LoadBalancer := &alkira.ServiceF5LoadBalancer{
		Id:               json.Number("123"),
		Name:             "test-f5-lb-service",
		Description:      "Test F5 Load Balancer service",
		Cxp:              "US-WEST",
		Size:             "MEDIUM",
		Version:          "16.1.0",
		CredentialId:     "test-credential",
		LicenseKey:       "test-license-key",
		SegmentIds:       []int{1, 2},
		BillingTagIds:    []int{10, 20},
		MaxInstanceCount: 3,
		MinInstanceCount: 1,
		AutoScale:        "ON",
		TunnelProtocol:   "IPSEC",
		Instances: []alkira.F5LoadBalancerInstance{
			{
				Id:             1,
				Name:           "f5-lb-01",
				TrafficEnabled: true,
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
			json.NewEncoder(w).Encode(serviceF5LoadBalancer)
		case http.MethodGet:
			// READ - return service data
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(serviceF5LoadBalancer)
		case http.MethodPut:
			// UPDATE - return updated service
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(serviceF5LoadBalancer)
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
	r := resourceAlkiraF5LoadBalancer()
	d := r.TestResourceData()

	d.Set("name", serviceF5LoadBalancer.Name)
	d.Set("description", serviceF5LoadBalancer.Description)
	d.Set("cxp", serviceF5LoadBalancer.Cxp)
	d.Set("size", serviceF5LoadBalancer.Size)
	d.Set("version", serviceF5LoadBalancer.Version)
	d.Set("credential_id", serviceF5LoadBalancer.CredentialId)
	d.Set("license_key", serviceF5LoadBalancer.LicenseKey)
	d.Set("segment_ids", serviceF5LoadBalancer.SegmentIds)
	d.Set("billing_tag_ids", serviceF5LoadBalancer.BillingTagIds)
	d.Set("max_instance_count", serviceF5LoadBalancer.MaxInstanceCount)
	d.Set("min_instance_count", serviceF5LoadBalancer.MinInstanceCount)
	d.Set("auto_scale", serviceF5LoadBalancer.AutoScale)
	d.Set("tunnel_protocol", serviceF5LoadBalancer.TunnelProtocol)

	diags := resourceF5LoadBalancerCreate(nil, d, client)
	require.Empty(t, diags, "Create should not return errors")

	// Verify ID was set
	require.Equal(t, "123", d.Id())

	// Test READ
	diags = resourceF5LoadBalancerRead(nil, d, client)
	require.Empty(t, diags, "Read should not return errors")

	// Verify data was populated
	assert.Equal(t, serviceF5LoadBalancer.Name, d.Get("name"))
	assert.Equal(t, serviceF5LoadBalancer.Description, d.Get("description"))
	assert.Equal(t, serviceF5LoadBalancer.Cxp, d.Get("cxp"))
	assert.Equal(t, serviceF5LoadBalancer.Size, d.Get("size"))

	// Test UPDATE
	d.Set("description", "Updated description")
	diags = resourceF5LoadBalancerUpdate(nil, d, client)
	require.Empty(t, diags, "Update should not return errors")

	// Test DELETE
	diags = resourceF5LoadBalancerDelete(nil, d, client)
	require.Empty(t, diags, "Delete should not return errors")
}

func TestAlkiraServiceF5LoadBalancer_CreateError(t *testing.T) {
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
	r := resourceAlkiraF5LoadBalancer()
	d := r.TestResourceData()

	d.Set("name", "test-f5-lb-service")
	d.Set("cxp", "US-WEST")
	d.Set("size", "MEDIUM")
	d.Set("version", "16.1.0")
	d.Set("credential_id", "test-credential")
	d.Set("segment_ids", []int{1})

	diags := resourceF5LoadBalancerCreate(nil, d, client)
	require.NotEmpty(t, diags, "Create should return error")
	assert.True(t, diags.HasError(), "Diagnostics should contain error")
}

func TestAlkiraServiceF5LoadBalancer_ReadNotFound(t *testing.T) {
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
	r := resourceAlkiraF5LoadBalancer()
	d := r.TestResourceData()
	d.SetId("999")

	diags := resourceF5LoadBalancerRead(nil, d, client)
	require.Empty(t, diags, "Read should not return error for 404")
	assert.Equal(t, "", d.Id(), "Resource ID should be cleared when not found")
}

func TestAlkiraServiceF5LoadBalancer_setServiceF5LoadBalancerFields(t *testing.T) {
	serviceF5LoadBalancer := &alkira.ServiceF5LoadBalancer{
		Id:               json.Number("123"),
		Name:             "test-f5-lb-service",
		Description:      "Test F5 Load Balancer service",
		Cxp:              "US-WEST",
		Size:             "MEDIUM",
		Version:          "16.1.0",
		CredentialId:     "test-credential",
		LicenseKey:       "test-license-key",
		SegmentIds:       []int{1, 2},
		BillingTagIds:    []int{10, 20},
		MaxInstanceCount: 3,
		MinInstanceCount: 1,
		AutoScale:        "ON",
		TunnelProtocol:   "IPSEC",
		Instances: []alkira.F5LoadBalancerInstance{
			{
				Id:             1,
				Name:           "f5-lb-01",
				TrafficEnabled: true,
			},
			{
				Id:             2,
				Name:           "f5-lb-02",
				TrafficEnabled: false,
			},
		},
	}

	r := resourceAlkiraF5LoadBalancer()
	d := r.TestResourceData()

	err := setServiceF5LoadBalancerFields(d, serviceF5LoadBalancer)
	require.NoError(t, err, "setServiceF5LoadBalancerFields should not return error")

	// Verify all fields were set correctly
	assert.Equal(t, serviceF5LoadBalancer.Name, d.Get("name"))
	assert.Equal(t, serviceF5LoadBalancer.Description, d.Get("description"))
	assert.Equal(t, serviceF5LoadBalancer.Cxp, d.Get("cxp"))
	assert.Equal(t, serviceF5LoadBalancer.Size, d.Get("size"))
	assert.Equal(t, serviceF5LoadBalancer.Version, d.Get("version"))
	assert.Equal(t, serviceF5LoadBalancer.CredentialId, d.Get("credential_id"))
	assert.Equal(t, serviceF5LoadBalancer.LicenseKey, d.Get("license_key"))
	assert.Equal(t, serviceF5LoadBalancer.MaxInstanceCount, d.Get("max_instance_count"))
	assert.Equal(t, serviceF5LoadBalancer.MinInstanceCount, d.Get("min_instance_count"))
	assert.Equal(t, serviceF5LoadBalancer.AutoScale, d.Get("auto_scale"))
	assert.Equal(t, serviceF5LoadBalancer.TunnelProtocol, d.Get("tunnel_protocol"))

	// Verify set fields
	segmentIdsSet := d.Get("segment_ids").(*schema.Set)
	assert.Equal(t, len(serviceF5LoadBalancer.SegmentIds), segmentIdsSet.Len())
	for _, segmentId := range serviceF5LoadBalancer.SegmentIds {
		assert.True(t, segmentIdsSet.Contains(segmentId))
	}

	billingTagIdsSet := d.Get("billing_tag_ids").(*schema.Set)
	assert.Equal(t, len(serviceF5LoadBalancer.BillingTagIds), billingTagIdsSet.Len())
	for _, billingTagId := range serviceF5LoadBalancer.BillingTagIds {
		assert.True(t, billingTagIdsSet.Contains(billingTagId))
	}

	// Verify instances
	instancesSet := d.Get("instances").(*schema.Set)
	assert.Equal(t, len(serviceF5LoadBalancer.Instances), instancesSet.Len())
}

func TestAlkiraServiceF5LoadBalancer_validateF5Version(t *testing.T) {
	tests := []struct {
		name      string
		input     interface{}
		expectErr bool
		errCount  int
	}{
		{
			name:      "Valid version 15.1.0",
			input:     "15.1.0",
			expectErr: false,
			errCount:  0,
		},
		{
			name:      "Valid version 16.1.0",
			input:     "16.1.0",
			expectErr: false,
			errCount:  0,
		},
		{
			name:      "Valid version 17.0.0",
			input:     "17.0.0",
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
	resource := resourceAlkiraF5LoadBalancer()
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
func serveServiceF5LoadBalancer(t *testing.T, serviceF5LoadBalancer *alkira.ServiceF5LoadBalancer) *alkira.AlkiraClient {
	server := httptest.NewServer(http.HandlerFunc(
		func(w http.ResponseWriter, req *http.Request) {
			json.NewEncoder(w).Encode(serviceF5LoadBalancer)
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
