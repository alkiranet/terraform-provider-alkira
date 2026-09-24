package alkira

import (
	"testing"

	"github.com/alkiranet/alkira-client-go/alkira"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestFlattenConnectorIPSecSegmentOptions(t *testing.T) {
	tests := []struct {
		name     string
		input    interface{}
		expected []map[string]interface{}
	}{
		{
			name:     "nil input",
			input:    nil,
			expected: nil,
		},
		{
			name:     "empty map",
			input:    map[string]interface{}{},
			expected: nil,
		},
		{
			name: "valid segment options with all fields",
			input: map[string]interface{}{
				"segment1": map[string]interface{}{
					"disableInternetExit":   false,
					"advertiseOnPremRoutes": true,
				},
			},
			expected: []map[string]interface{}{
				{
					"name":                     "segment1",
					"advertise_default_route":  true,
					"advertise_on_prem_routes": true,
				},
			},
		},
		{
			name: "segment options with disable_internet_exit=true",
			input: map[string]interface{}{
				"segment2": map[string]interface{}{
					"disableInternetExit":   true,
					"advertiseOnPremRoutes": false,
				},
			},
			expected: []map[string]interface{}{
				{
					"name":                     "segment2",
					"advertise_default_route":  false,
					"advertise_on_prem_routes": false,
				},
			},
		},
		{
			name: "segment options with missing fields (defaults)",
			input: map[string]interface{}{
				"segment3": map[string]interface{}{},
			},
			expected: []map[string]interface{}{
				{
					"name":                     "segment3",
					"advertise_default_route":  false,
					"advertise_on_prem_routes": false,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := flattenConnectorIPSecSegmentOptions(tt.input)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestFlattenConnectorIPSecRoutingOptions(t *testing.T) {
	tests := []struct {
		name     string
		input    *alkira.ConnectorIPSecRoutingOptions
		expected []map[string]interface{}
	}{
		{
			name:     "nil input",
			input:    nil,
			expected: nil,
		},
		{
			name:     "empty routing options",
			input:    &alkira.ConnectorIPSecRoutingOptions{},
			expected: nil,
		},
		{
			name: "STATIC routing options",
			input: &alkira.ConnectorIPSecRoutingOptions{
				StaticRouting: &alkira.ConnectorIPSecStaticRouting{
					PrefixListId: 123,
					Availability: "IKE_STATUS",
				},
			},
			expected: []map[string]interface{}{
				{
					"type":           "STATIC",
					"prefix_list_id": 123,
					"availability":   "IKE_STATUS",
				},
			},
		},
		{
			name: "DYNAMIC routing options",
			input: &alkira.ConnectorIPSecRoutingOptions{
				DynamicRouting: &alkira.ConnectorIPSecDynamicRouting{
					CustomerGwAsn:    "65001",
					Availability:     "IPSEC_INTERFACE_PING",
					BgpAuthKeyAlkira: "bgp-secret",
				},
			},
			expected: []map[string]interface{}{
				{
					"type":                 "DYNAMIC",
					"customer_gateway_asn": "65001",
					"availability":         "IPSEC_INTERFACE_PING",
					"bgp_auth_key":         "bgp-secret",
				},
			},
		},
		{
			// PING is a legacy backend value that the API may still return for
			// older connectors, even though it is no longer accepted in
			// configuration (AK-73307). Flatten must pass it through unchanged
			// so the drift stays visible instead of being silently normalized.
			name: "BOTH routing options with legacy PING availability from API",
			input: &alkira.ConnectorIPSecRoutingOptions{
				StaticRouting: &alkira.ConnectorIPSecStaticRouting{
					PrefixListId: 456,
					Availability: "PING",
				},
				DynamicRouting: &alkira.ConnectorIPSecDynamicRouting{
					CustomerGwAsn:    "65002",
					Availability:     "PING",
					BgpAuthKeyAlkira: "bgp-secret-2",
				},
			},
			expected: []map[string]interface{}{
				{
					"type":                 "BOTH",
					"prefix_list_id":       456,
					"customer_gateway_asn": "65002",
					"availability":         "PING",
					"bgp_auth_key":         "bgp-secret-2",
				},
			},
		},
		{
			name: "STATIC routing with default availability",
			input: &alkira.ConnectorIPSecRoutingOptions{
				StaticRouting: &alkira.ConnectorIPSecStaticRouting{
					PrefixListId: 789,
				},
			},
			expected: []map[string]interface{}{
				{
					"type":           "STATIC",
					"prefix_list_id": 789,
					"availability":   "IPSEC_INTERFACE_PING",
				},
			},
		},
		{
			name: "DYNAMIC routing without bgp_auth_key",
			input: &alkira.ConnectorIPSecRoutingOptions{
				DynamicRouting: &alkira.ConnectorIPSecDynamicRouting{
					CustomerGwAsn: "65003",
				},
			},
			expected: []map[string]interface{}{
				{
					"type":                 "DYNAMIC",
					"customer_gateway_asn": "65003",
					"availability":         "IPSEC_INTERFACE_PING",
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := flattenConnectorIPSecRoutingOptions(tt.input)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestFlattenConnectorIPSecPolicyOptions(t *testing.T) {
	tests := []struct {
		name     string
		input    *alkira.ConnectorIPSecPolicyOptions
		expected []map[string]interface{}
	}{
		{
			name:     "nil input",
			input:    nil,
			expected: nil,
		},
		{
			name: "valid policy options",
			input: &alkira.ConnectorIPSecPolicyOptions{
				BranchTSPrefixListIds: []int{1, 2, 3},
				CxpTSPrefixListIds:    []int{4, 5, 6},
			},
			expected: []map[string]interface{}{
				{
					"on_prem_prefix_list_ids": []int{1, 2, 3},
					"cxp_prefix_list_ids":     []int{4, 5, 6},
				},
			},
		},
		{
			name: "empty policy options",
			input: &alkira.ConnectorIPSecPolicyOptions{
				BranchTSPrefixListIds: []int{},
				CxpTSPrefixListIds:    []int{},
			},
			expected: []map[string]interface{}{
				{
					"on_prem_prefix_list_ids": []int{},
					"cxp_prefix_list_ids":     []int{},
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := flattenConnectorIPSecPolicyOptions(tt.input)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestAlkiraConnectorIPSec_validateAvailability(t *testing.T) {
	tests := []ValidationTestCase{
		{
			Name:      "Valid IKE_STATUS availability",
			Input:     "IKE_STATUS",
			ExpectErr: false,
			ErrCount:  0,
		},
		{
			Name:      "Valid IPSEC_INTERFACE_PING availability",
			Input:     "IPSEC_INTERFACE_PING",
			ExpectErr: false,
			ErrCount:  0,
		},
		{
			Name:      "Legacy PING availability is no longer accepted (AK-73307)",
			Input:     "PING",
			ExpectErr: true,
			ErrCount:  1,
		},
		{
			Name:      "Wrong case is rejected",
			Input:     "ike_status",
			ExpectErr: true,
			ErrCount:  1,
		},
		{
			Name:      "Empty string",
			Input:     "",
			ExpectErr: true,
			ErrCount:  1,
		},
	}

	resource := resourceAlkiraConnectorIPSec()

	routingOptionsSchema, exists := resource.Schema["routing_options"]
	require.True(t, exists, "routing_options schema field not found")

	routingOptionsElem, ok := routingOptionsSchema.Elem.(*schema.Resource)
	require.True(t, ok, "routing_options schema element is not a resource")

	availabilitySchema, exists := routingOptionsElem.Schema["availability"]
	require.True(t, exists, "availability schema field not found in routing_options")
	require.NotNil(t, availabilitySchema.ValidateFunc, "availability has no ValidateFunc")

	for _, tt := range tests {
		t.Run(tt.Name, func(t *testing.T) {
			warnings, errors := availabilitySchema.ValidateFunc(tt.Input, "availability")

			if tt.ExpectErr {
				assert.Len(t, errors, tt.ErrCount, "Expected %d errors for input %v", tt.ErrCount, tt.Input)
			} else {
				assert.Empty(t, errors, "Expected no errors for input %v", tt.Input)
			}
			assert.Empty(t, warnings, "Expected no warnings")
		})
	}
}
