package alkira

import (
	"testing"

	"github.com/alkiranet/alkira-client-go/alkira"
	"github.com/stretchr/testify/assert"
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
			input:    map[string]alkira.ConnectorIPSecSegmentOptions{},
			expected: nil,
		},
		{
			name: "valid segment options with all fields",
			input: map[string]alkira.ConnectorIPSecSegmentOptions{
				"segment1": {
					DisableInternetExit:   boolPtr(false),
					AdvertiseOnPremRoutes: boolPtr(true),
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
			input: map[string]alkira.ConnectorIPSecSegmentOptions{
				"segment2": {
					DisableInternetExit:   boolPtr(true),
					AdvertiseOnPremRoutes: boolPtr(false),
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
			name: "segment options with nil pointers (defaults)",
			input: map[string]alkira.ConnectorIPSecSegmentOptions{
				"segment3": {
					DisableInternetExit:   nil,
					AdvertiseOnPremRoutes: nil,
				},
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
			name: "BOTH routing options",
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
