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

func TestAlkiraServiceF5vServerEndpoint_buildServiceF5vServerEndpointRequest(t *testing.T) {
	r := resourceAlkiraServiceF5vServerEndpoint()
	d := r.TestResourceData()

	// Test with complete F5 vServer Endpoint service data
	expectedName := "test-f5-vserver-endpoint"
	expectedF5ServiceId := 123
	expectedF5ServiceInstanceIds := []int{1, 2}
	expectedVirtualServerName := "test-virtual-server"
	expectedVirtualServerAddress := "192.168.1.100"
	expectedVirtualServerPort := 80
	expectedProtocol := "HTTP"
	expectedLoadBalancingMethod := "round-robin"
	expectedDescription := "Test F5 vServer Endpoint description"

	// Set up pool members
	poolMembers := []interface{}{
		map[string]interface{}{
			"name":    "member1",
			"address": "10.0.0.10",
			"port":    8080,
			"enabled": true,
		},
		map[string]interface{}{
			"name":    "member2",
			"address": "10.0.0.11",
			"port":    8080,
			"enabled": false,
		},
	}

	// Set up health monitors
	healthMonitors := []interface{}{
		map[string]interface{}{
			"name":     "http-monitor",
			"type":     "HTTP",
			"interval": 30,
			"timeout":  10,
			"uri":      "/health",
		},
	}

	d.Set("name", expectedName)
	d.Set("f5_service_id", expectedF5ServiceId)
	d.Set("f5_service_instance_ids", expectedF5ServiceInstanceIds)
	d.Set("virtual_server_name", expectedVirtualServerName)
	d.Set("virtual_server_address", expectedVirtualServerAddress)
	d.Set("virtual_server_port", expectedVirtualServerPort)
	d.Set("protocol", expectedProtocol)
	d.Set("load_balancing_method", expectedLoadBalancingMethod)
	d.Set("description", expectedDescription)
	d.Set("pool_members", poolMembers)
	d.Set("health_monitors", healthMonitors)

	request := buildServiceF5vServerEndpointRequest(d)

	require.Equal(t, expectedName, request.Name)
	require.Equal(t, expectedF5ServiceId, request.F5ServiceId)
	require.Equal(t, expectedF5ServiceInstanceIds, request.F5ServiceInstanceIds)
	require.Equal(t, expectedVirtualServerName, request.VirtualServerName)
	require.Equal(t, expectedVirtualServerAddress, request.VirtualServerAddress)
	require.Equal(t, expectedVirtualServerPort, request.VirtualServerPort)
	require.Equal(t, expectedProtocol, request.Protocol)
	require.Equal(t, expectedLoadBalancingMethod, request.LoadBalancingMethod)
	require.Equal(t, expectedDescription, request.Description)

	// Verify pool members
	require.Len(t, request.PoolMembers, 2)
	assert.Equal(t, "member1", request.PoolMembers[0].Name)
	assert.Equal(t, "10.0.0.10", request.PoolMembers[0].Address)
	assert.Equal(t, 8080, request.PoolMembers[0].Port)
	assert.Equal(t, true, request.PoolMembers[0].Enabled)

	assert.Equal(t, "member2", request.PoolMembers[1].Name)
	assert.Equal(t, "10.0.0.11", request.PoolMembers[1].Address)
	assert.Equal(t, 8080, request.PoolMembers[1].Port)
	assert.Equal(t, false, request.PoolMembers[1].Enabled)

	// Verify health monitors
	require.Len(t, request.HealthMonitors, 1)
	assert.Equal(t, "http-monitor", request.HealthMonitors[0].Name)
	assert.Equal(t, "HTTP", request.HealthMonitors[0].Type)
	assert.Equal(t, 30, request.HealthMonitors[0].Interval)
	assert.Equal(t, 10, request.HealthMonitors[0].Timeout)
	assert.Equal(t, "/health", request.HealthMonitors[0].Uri)
}

func TestAlkiraServiceF5vServerEndpoint_buildServiceF5vServerEndpointRequestMinimal(t *testing.T) {
	r := resourceAlkiraServiceF5vServerEndpoint()
	d := r.TestResourceData()

	// Test with minimal required F5 vServer Endpoint service data
	expectedName := "minimal-f5-vserver-endpoint"
	expectedF5ServiceId := 456
	expectedF5ServiceInstanceIds := []int{3}
	expectedVirtualServerName := "minimal-virtual-server"
	expectedVirtualServerAddress := "172.16.1.100"
	expectedVirtualServerPort := 443

	d.Set("name", expectedName)
	d.Set("f5_service_id", expectedF5ServiceId)
	d.Set("f5_service_instance_ids", expectedF5ServiceInstanceIds)
	d.Set("virtual_server_name", expectedVirtualServerName)
	d.Set("virtual_server_address", expectedVirtualServerAddress)
	d.Set("virtual_server_port", expectedVirtualServerPort)

	request := buildServiceF5vServerEndpointRequest(d)

	require.Equal(t, expectedName, request.Name)
	require.Equal(t, expectedF5ServiceId, request.F5ServiceId)
	require.Equal(t, expectedF5ServiceInstanceIds, request.F5ServiceInstanceIds)
	require.Equal(t, expectedVirtualServerName, request.VirtualServerName)
	require.Equal(t, expectedVirtualServerAddress, request.VirtualServerAddress)
	require.Equal(t, expectedVirtualServerPort, request.VirtualServerPort)
	require.Equal(t, "", request.Description)                    // Should be empty when not set
	require.Equal(t, "HTTP", request.Protocol)                   // Default value
	require.Equal(t, "round-robin", request.LoadBalancingMethod) // Default value
}

func TestAlkiraServiceF5vServerEndpoint_buildServiceF5vServerEndpointRequestSSL(t *testing.T) {
	r := resourceAlkiraServiceF5vServerEndpoint()
	d := r.TestResourceData()

	// Test with SSL configuration
	expectedName := "ssl-f5-vserver-endpoint"
	expectedF5ServiceId := 789
	expectedF5ServiceInstanceIds := []int{4}
	expectedVirtualServerName := "ssl-virtual-server"
	expectedVirtualServerAddress := "10.1.1.100"
	expectedVirtualServerPort := 443
	expectedProtocol := "HTTPS"
	expectedSslEnabled := true
	expectedSslCertificate := "test-certificate"
	expectedSslKey := "test-key"

	d.Set("name", expectedName)
	d.Set("f5_service_id", expectedF5ServiceId)
	d.Set("f5_service_instance_ids", expectedF5ServiceInstanceIds)
	d.Set("virtual_server_name", expectedVirtualServerName)
	d.Set("virtual_server_address", expectedVirtualServerAddress)
	d.Set("virtual_server_port", expectedVirtualServerPort)
	d.Set("protocol", expectedProtocol)
	d.Set("ssl_enabled", expectedSslEnabled)
	d.Set("ssl_certificate", expectedSslCertificate)
	d.Set("ssl_key", expectedSslKey)

	request := buildServiceF5vServerEndpointRequest(d)

	require.Equal(t, expectedName, request.Name)
	require.Equal(t, expectedProtocol, request.Protocol)
	require.Equal(t, expectedSslEnabled, request.SslEnabled)
	require.Equal(t, expectedSslCertificate, request.SslCertificate)
	require.Equal(t, expectedSslKey, request.SslKey)
}

func TestAlkiraServiceF5vServerEndpoint_resourceSchema(t *testing.T) {
	resource := resourceAlkiraServiceF5vServerEndpoint()

	// Test required fields
	nameSchema := resource.Schema["name"]
	assert.True(t, nameSchema.Required, "Name should be required")
	assert.Equal(t, schema.TypeString, nameSchema.Type, "Name should be string type")

	f5ServiceIdSchema := resource.Schema["f5_service_id"]
	assert.True(t, f5ServiceIdSchema.Required, "F5 service ID should be required")
	assert.Equal(t, schema.TypeInt, f5ServiceIdSchema.Type, "F5 service ID should be int type")

	f5ServiceInstanceIdsSchema := resource.Schema["f5_service_instance_ids"]
	assert.True(t, f5ServiceInstanceIdsSchema.Required, "F5 service instance IDs should be required")
	assert.Equal(t, schema.TypeSet, f5ServiceInstanceIdsSchema.Type, "F5 service instance IDs should be set type")

	virtualServerNameSchema := resource.Schema["virtual_server_name"]
	assert.True(t, virtualServerNameSchema.Required, "Virtual server name should be required")
	assert.Equal(t, schema.TypeString, virtualServerNameSchema.Type, "Virtual server name should be string type")

	virtualServerAddressSchema := resource.Schema["virtual_server_address"]
	assert.True(t, virtualServerAddressSchema.Required, "Virtual server address should be required")
	assert.Equal(t, schema.TypeString, virtualServerAddressSchema.Type, "Virtual server address should be string type")

	virtualServerPortSchema := resource.Schema["virtual_server_port"]
	assert.True(t, virtualServerPortSchema.Required, "Virtual server port should be required")
	assert.Equal(t, schema.TypeInt, virtualServerPortSchema.Type, "Virtual server port should be int type")

	// Test optional fields with defaults
	protocolSchema := resource.Schema["protocol"]
	assert.True(t, protocolSchema.Optional, "Protocol should be optional")
	assert.Equal(t, schema.TypeString, protocolSchema.Type, "Protocol should be string type")
	assert.Equal(t, "HTTP", protocolSchema.Default, "Protocol should default to HTTP")

	loadBalancingMethodSchema := resource.Schema["load_balancing_method"]
	assert.True(t, loadBalancingMethodSchema.Optional, "Load balancing method should be optional")
	assert.Equal(t, schema.TypeString, loadBalancingMethodSchema.Type, "Load balancing method should be string type")
	assert.Equal(t, "round-robin", loadBalancingMethodSchema.Default, "Load balancing method should default to round-robin")

	sslEnabledSchema := resource.Schema["ssl_enabled"]
	assert.True(t, sslEnabledSchema.Optional, "SSL enabled should be optional")
	assert.Equal(t, schema.TypeBool, sslEnabledSchema.Type, "SSL enabled should be bool type")
	assert.Equal(t, false, sslEnabledSchema.Default, "SSL enabled should default to false")

	// Test optional fields
	descSchema := resource.Schema["description"]
	assert.True(t, descSchema.Optional, "Description should be optional")
	assert.Equal(t, schema.TypeString, descSchema.Type, "Description should be string type")

	poolMembersSchema := resource.Schema["pool_members"]
	assert.True(t, poolMembersSchema.Optional, "Pool members should be optional")
	assert.Equal(t, schema.TypeSet, poolMembersSchema.Type, "Pool members should be set type")

	healthMonitorsSchema := resource.Schema["health_monitors"]
	assert.True(t, healthMonitorsSchema.Optional, "Health monitors should be optional")
	assert.Equal(t, schema.TypeSet, healthMonitorsSchema.Type, "Health monitors should be set type")

	sslCertificateSchema := resource.Schema["ssl_certificate"]
	assert.True(t, sslCertificateSchema.Optional, "SSL certificate should be optional")
	assert.Equal(t, schema.TypeString, sslCertificateSchema.Type, "SSL certificate should be string type")

	sslKeySchema := resource.Schema["ssl_key"]
	assert.True(t, sslKeySchema.Optional, "SSL key should be optional")
	assert.Equal(t, schema.TypeString, sslKeySchema.Type, "SSL key should be string type")

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

func TestAlkiraServiceF5vServerEndpoint_validateProtocol(t *testing.T) {
	tests := []struct {
		name      string
		input     interface{}
		expectErr bool
		errCount  int
	}{
		{
			name:      "Valid HTTP",
			input:     "HTTP",
			expectErr: false,
			errCount:  0,
		},
		{
			name:      "Valid HTTPS",
			input:     "HTTPS",
			expectErr: false,
			errCount:  0,
		},
		{
			name:      "Valid TCP",
			input:     "TCP",
			expectErr: false,
			errCount:  0,
		},
		{
			name:      "Valid UDP",
			input:     "UDP",
			expectErr: false,
			errCount:  0,
		},
		{
			name:      "Invalid protocol",
			input:     "INVALID_PROTOCOL",
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

	resource := resourceAlkiraServiceF5vServerEndpoint()
	protocolSchema := resource.Schema["protocol"]

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			warnings, errors := protocolSchema.ValidateFunc(tt.input, "protocol")

			if tt.expectErr {
				assert.Len(t, errors, tt.errCount, "Expected %d errors for input %v", tt.errCount, tt.input)
			} else {
				assert.Empty(t, errors, "Expected no errors for input %v", tt.input)
			}
			assert.Empty(t, warnings, "Expected no warnings")
		})
	}
}

func TestAlkiraServiceF5vServerEndpoint_validateLoadBalancingMethod(t *testing.T) {
	tests := []struct {
		name      string
		input     interface{}
		expectErr bool
		errCount  int
	}{
		{
			name:      "Valid round-robin",
			input:     "round-robin",
			expectErr: false,
			errCount:  0,
		},
		{
			name:      "Valid least-connections",
			input:     "least-connections",
			expectErr: false,
			errCount:  0,
		},
		{
			name:      "Valid weighted-round-robin",
			input:     "weighted-round-robin",
			expectErr: false,
			errCount:  0,
		},
		{
			name:      "Invalid load balancing method",
			input:     "INVALID_METHOD",
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

	resource := resourceAlkiraServiceF5vServerEndpoint()
	loadBalancingMethodSchema := resource.Schema["load_balancing_method"]

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			warnings, errors := loadBalancingMethodSchema.ValidateFunc(tt.input, "load_balancing_method")

			if tt.expectErr {
				assert.Len(t, errors, tt.errCount, "Expected %d errors for input %v", tt.errCount, tt.input)
			} else {
				assert.Empty(t, errors, "Expected no errors for input %v", tt.input)
			}
			assert.Empty(t, warnings, "Expected no warnings")
		})
	}
}

func TestAlkiraServiceF5vServerEndpoint_CreateUpdateDelete(t *testing.T) {
	// Test data
	serviceF5vServerEndpoint := &alkira.ServiceF5vServerEndpoint{
		Id:                   json.Number("123"),
		Name:                 "test-f5-vserver-endpoint",
		Description:          "Test F5 vServer Endpoint",
		F5ServiceId:          456,
		F5ServiceInstanceIds: []int{1, 2},
		VirtualServerName:    "test-virtual-server",
		VirtualServerAddress: "192.168.1.100",
		VirtualServerPort:    80,
		Protocol:             "HTTP",
		LoadBalancingMethod:  "round-robin",
		SslEnabled:           false,
		PoolMembers: []alkira.F5PoolMember{
			{
				Name:    "member1",
				Address: "10.0.0.10",
				Port:    8080,
				Enabled: true,
			},
		},
		HealthMonitors: []alkira.F5HealthMonitor{
			{
				Name:     "http-monitor",
				Type:     "HTTP",
				Interval: 30,
				Timeout:  10,
				Uri:      "/health",
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
			json.NewEncoder(w).Encode(serviceF5vServerEndpoint)
		case http.MethodGet:
			// READ - return service data
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(serviceF5vServerEndpoint)
		case http.MethodPut:
			// UPDATE - return updated service
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(serviceF5vServerEndpoint)
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
	r := resourceAlkiraServiceF5vServerEndpoint()
	d := r.TestResourceData()

	d.Set("name", serviceF5vServerEndpoint.Name)
	d.Set("description", serviceF5vServerEndpoint.Description)
	d.Set("f5_service_id", serviceF5vServerEndpoint.F5ServiceId)
	d.Set("f5_service_instance_ids", serviceF5vServerEndpoint.F5ServiceInstanceIds)
	d.Set("virtual_server_name", serviceF5vServerEndpoint.VirtualServerName)
	d.Set("virtual_server_address", serviceF5vServerEndpoint.VirtualServerAddress)
	d.Set("virtual_server_port", serviceF5vServerEndpoint.VirtualServerPort)
	d.Set("protocol", serviceF5vServerEndpoint.Protocol)
	d.Set("load_balancing_method", serviceF5vServerEndpoint.LoadBalancingMethod)
	d.Set("ssl_enabled", serviceF5vServerEndpoint.SslEnabled)

	diags := resourceF5vServerEndpointCreate(nil, d, client)
	require.Empty(t, diags, "Create should not return errors")

	// Verify ID was set
	require.Equal(t, "123", d.Id())

	// Test READ
	diags = resourceF5vServerEndpointRead(nil, d, client)
	require.Empty(t, diags, "Read should not return errors")

	// Verify data was populated
	assert.Equal(t, serviceF5vServerEndpoint.Name, d.Get("name"))
	assert.Equal(t, serviceF5vServerEndpoint.Description, d.Get("description"))
	assert.Equal(t, serviceF5vServerEndpoint.F5ServiceId, d.Get("f5_service_id"))
	assert.Equal(t, serviceF5vServerEndpoint.VirtualServerName, d.Get("virtual_server_name"))

	// Test UPDATE
	d.Set("description", "Updated description")
	diags = resourceF5vServerEndpointUpdate(nil, d, client)
	require.Empty(t, diags, "Update should not return errors")

	// Test DELETE
	diags = resourceF5vServerEndpointDelete(nil, d, client)
	require.Empty(t, diags, "Delete should not return errors")
}

func TestAlkiraServiceF5vServerEndpoint_CreateError(t *testing.T) {
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
	r := resourceAlkiraServiceF5vServerEndpoint()
	d := r.TestResourceData()

	d.Set("name", "test-f5-vserver-endpoint")
	d.Set("f5_service_id", 123)
	d.Set("f5_service_instance_ids", []int{1})
	d.Set("virtual_server_name", "test-virtual-server")
	d.Set("virtual_server_address", "192.168.1.100")
	d.Set("virtual_server_port", 80)

	diags := resourceF5vServerEndpointCreate(nil, d, client)
	require.NotEmpty(t, diags, "Create should return error")
	assert.True(t, diags.HasError(), "Diagnostics should contain error")
}

func TestAlkiraServiceF5vServerEndpoint_ReadNotFound(t *testing.T) {
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
	r := resourceAlkiraServiceF5vServerEndpoint()
	d := r.TestResourceData()
	d.SetId("999")

	diags := resourceF5vServerEndpointRead(nil, d, client)
	require.Empty(t, diags, "Read should not return error for 404")
	assert.Equal(t, "", d.Id(), "Resource ID should be cleared when not found")
}

func TestAlkiraServiceF5vServerEndpoint_setServiceF5vServerEndpointFields(t *testing.T) {
	serviceF5vServerEndpoint := &alkira.ServiceF5vServerEndpoint{
		Id:                   json.Number("123"),
		Name:                 "test-f5-vserver-endpoint",
		Description:          "Test F5 vServer Endpoint",
		F5ServiceId:          456,
		F5ServiceInstanceIds: []int{1, 2},
		VirtualServerName:    "test-virtual-server",
		VirtualServerAddress: "192.168.1.100",
		VirtualServerPort:    80,
		Protocol:             "HTTP",
		LoadBalancingMethod:  "round-robin",
		SslEnabled:           true,
		SslCertificate:       "test-certificate",
		SslKey:               "test-key",
		PoolMembers: []alkira.F5PoolMember{
			{
				Name:    "member1",
				Address: "10.0.0.10",
				Port:    8080,
				Enabled: true,
			},
			{
				Name:    "member2",
				Address: "10.0.0.11",
				Port:    8080,
				Enabled: false,
			},
		},
		HealthMonitors: []alkira.F5HealthMonitor{
			{
				Name:     "http-monitor",
				Type:     "HTTP",
				Interval: 30,
				Timeout:  10,
				Uri:      "/health",
			},
			{
				Name:     "tcp-monitor",
				Type:     "TCP",
				Interval: 60,
				Timeout:  20,
			},
		},
	}

	r := resourceAlkiraServiceF5vServerEndpoint()
	d := r.TestResourceData()

	err := setServiceF5vServerEndpointFields(d, serviceF5vServerEndpoint)
	require.NoError(t, err, "setServiceF5vServerEndpointFields should not return error")

	// Verify all fields were set correctly
	assert.Equal(t, serviceF5vServerEndpoint.Name, d.Get("name"))
	assert.Equal(t, serviceF5vServerEndpoint.Description, d.Get("description"))
	assert.Equal(t, serviceF5vServerEndpoint.F5ServiceId, d.Get("f5_service_id"))
	assert.Equal(t, serviceF5vServerEndpoint.VirtualServerName, d.Get("virtual_server_name"))
	assert.Equal(t, serviceF5vServerEndpoint.VirtualServerAddress, d.Get("virtual_server_address"))
	assert.Equal(t, serviceF5vServerEndpoint.VirtualServerPort, d.Get("virtual_server_port"))
	assert.Equal(t, serviceF5vServerEndpoint.Protocol, d.Get("protocol"))
	assert.Equal(t, serviceF5vServerEndpoint.LoadBalancingMethod, d.Get("load_balancing_method"))
	assert.Equal(t, serviceF5vServerEndpoint.SslEnabled, d.Get("ssl_enabled"))
	assert.Equal(t, serviceF5vServerEndpoint.SslCertificate, d.Get("ssl_certificate"))
	assert.Equal(t, serviceF5vServerEndpoint.SslKey, d.Get("ssl_key"))

	// Verify set fields
	f5ServiceInstanceIdsSet := d.Get("f5_service_instance_ids").(*schema.Set)
	assert.Equal(t, len(serviceF5vServerEndpoint.F5ServiceInstanceIds), f5ServiceInstanceIdsSet.Len())
	for _, instanceId := range serviceF5vServerEndpoint.F5ServiceInstanceIds {
		assert.True(t, f5ServiceInstanceIdsSet.Contains(instanceId))
	}

	// Verify pool members
	poolMembersSet := d.Get("pool_members").(*schema.Set)
	assert.Equal(t, len(serviceF5vServerEndpoint.PoolMembers), poolMembersSet.Len())

	// Verify health monitors
	healthMonitorsSet := d.Get("health_monitors").(*schema.Set)
	assert.Equal(t, len(serviceF5vServerEndpoint.HealthMonitors), healthMonitorsSet.Len())
}

func TestAlkiraServiceF5vServerEndpoint_validatePortRange(t *testing.T) {
	tests := []struct {
		name      string
		input     interface{}
		expectErr bool
		errCount  int
	}{
		{
			name:      "Valid port 80",
			input:     80,
			expectErr: false,
			errCount:  0,
		},
		{
			name:      "Valid port 443",
			input:     443,
			expectErr: false,
			errCount:  0,
		},
		{
			name:      "Valid port 8080",
			input:     8080,
			expectErr: false,
			errCount:  0,
		},
		{
			name:      "Valid port 65535",
			input:     65535,
			expectErr: false,
			errCount:  0,
		},
		{
			name:      "Invalid port 0",
			input:     0,
			expectErr: true,
			errCount:  1,
		},
		{
			name:      "Invalid port 65536",
			input:     65536,
			expectErr: true,
			errCount:  1,
		},
		{
			name:      "Invalid negative port",
			input:     -1,
			expectErr: true,
			errCount:  1,
		},
	}

	resource := resourceAlkiraServiceF5vServerEndpoint()
	virtualServerPortSchema := resource.Schema["virtual_server_port"]

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			warnings, errors := virtualServerPortSchema.ValidateFunc(tt.input, "virtual_server_port")

			if tt.expectErr {
				assert.Len(t, errors, tt.errCount, "Expected %d errors for input %v", tt.errCount, tt.input)
			} else {
				assert.Empty(t, errors, "Expected no errors for input %v", tt.input)
			}
			assert.Empty(t, warnings, "Expected no warnings")
		})
	}
}

// TEST HELPER
func serveServiceF5vServerEndpoint(t *testing.T, serviceF5vServerEndpoint *alkira.ServiceF5vServerEndpoint) *alkira.AlkiraClient {
	server := httptest.NewServer(http.HandlerFunc(
		func(w http.ResponseWriter, req *http.Request) {
			json.NewEncoder(w).Encode(serviceF5vServerEndpoint)
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
