package alkira

import (
	"testing"

	"github.com/alkiranet/alkira-client-go/alkira"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/stretchr/testify/assert"
)

func TestExpandConnectorAdvIPSecAdvancedOptions(t *testing.T) {
	tests := []struct {
		name        string
		input       []interface{}
		expected    *alkira.ConnectorAdvIPSecAdvanced
		expectError bool
	}{
		{
			name:        "nil input",
			input:       nil,
			expected:    nil,
			expectError: false,
		},
		{
			name:        "empty input",
			input:       []interface{}{},
			expected:    nil,
			expectError: false,
		},
		{
			name: "valid advanced options",
			input: []interface{}{
				map[string]interface{}{
					"ike_version":       "2",
					"initiator":         true,
					"remote_auth_type":  "PSK",
					"remote_auth_value": "secret123",
				},
			},
			expected: &alkira.ConnectorAdvIPSecAdvanced{
				IkeVersion:      "2",
				Initiator:       true,
				RemoteAuthType:  "PSK",
				RemoteAuthValue: "secret123",
			},
			expectError: false,
		},
		{
			name: "partial advanced options",
			input: []interface{}{
				map[string]interface{}{
					"ike_version": "2",
					"initiator":   false,
				},
			},
			expected: &alkira.ConnectorAdvIPSecAdvanced{
				IkeVersion: "2",
				Initiator:  false,
			},
			expectError: false,
		},
		{
			name: "multiple inputs - should return nil",
			input: []interface{}{
				map[string]interface{}{"ike_version": "2"},
				map[string]interface{}{"ike_version": "1"},
			},
			expected:    nil,
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := expandConnectorAdvIPSecAdvancedOptions(tt.input)

			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expected, result)
			}
		})
	}
}

func TestExpandConnectorAdvIPSecTunnel(t *testing.T) {
	tests := []struct {
		name          string
		input         []interface{}
		expectedCount int
		checkTunnelId string
		checkPSK      string
	}{
		{
			name:          "nil input",
			input:         nil,
			expectedCount: 0,
		},
		{
			name:          "empty input",
			input:         []interface{}{},
			expectedCount: 0,
		},
		{
			name: "valid tunnel configuration",
			input: []interface{}{
				map[string]interface{}{
					"customer_end_overlay_ip_reservation_id": "overlay-ip-1",
					"cxp_end_overlay_ip_reservation_id":      "cxp-overlay-ip-1",
					"cxp_end_public_ip_reservation_id":       "public-ip-1",
					"id":                                     "tunnel-1",
					"preshared_key":                          "secret123",
					"profile_id":                             1,
					"number":                                 1,
					"advanced_options": []interface{}{
						map[string]interface{}{
							"ike_version": "2",
							"initiator":   true,
						},
					},
				},
			},
			expectedCount: 1,
			checkTunnelId: "tunnel-1",
			checkPSK:      "secret123",
		},
		{
			name: "multiple tunnels",
			input: []interface{}{
				map[string]interface{}{
					"customer_end_overlay_ip_reservation_id": "overlay-ip-1",
					"cxp_end_overlay_ip_reservation_id":      "cxp-overlay-ip-1",
					"cxp_end_public_ip_reservation_id":       "public-ip-1",
					"id":                                     "tunnel-1",
					"preshared_key":                          "secret123",
					"profile_id":                             1,
					"number":                                 1,
				},
				map[string]interface{}{
					"customer_end_overlay_ip_reservation_id": "overlay-ip-2",
					"cxp_end_overlay_ip_reservation_id":      "cxp-overlay-ip-2",
					"cxp_end_public_ip_reservation_id":       "public-ip-2",
					"id":                                     "tunnel-2",
					"preshared_key":                          "secret456",
					"profile_id":                             2,
					"number":                                 2,
				},
			},
			expectedCount: 2,
			checkTunnelId: "tunnel-1",
			checkPSK:      "secret123",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := expandConnectorAdvIPSecTunnel(tt.input)

			if tt.expectedCount == 0 {
				assert.Nil(t, result)
			} else {
				assert.Len(t, result, tt.expectedCount)

				if tt.checkTunnelId != "" {
					found := false
					for _, tunnel := range result {
						if tunnel.Id == tt.checkTunnelId && tunnel.PresharedKey == tt.checkPSK {
							found = true
							break
						}
					}
					assert.True(t, found, "Expected tunnel with ID %s and PSK %s not found", tt.checkTunnelId, tt.checkPSK)
				}
			}
		})
	}
}

func TestExpandConnectorAdvIPSecGateway(t *testing.T) {
	tests := []struct {
		name              string
		input             []interface{}
		expectedCount     int
		checkGatewayName  string
		checkCustomerGwIp string
	}{
		{
			name:          "nil input",
			input:         nil,
			expectedCount: 0,
		},
		{
			name:          "empty input",
			input:         []interface{}{},
			expectedCount: 0,
		},
		{
			name: "valid gateway configuration",
			input: []interface{}{
				map[string]interface{}{
					"name":                "test-gateway",
					"customer_gateway_ip": "192.168.1.1",
					"ha_mode":             "SINGLE",
					"id":                  1,
					"tunnel": []interface{}{
						map[string]interface{}{
							"customer_end_overlay_ip_reservation_id": "overlay-ip-1",
							"cxp_end_overlay_ip_reservation_id":      "cxp-overlay-ip-1",
							"cxp_end_public_ip_reservation_id":       "public-ip-1",
							"id":                                     "tunnel-1",
							"preshared_key":                          "secret123",
							"profile_id":                             1,
							"number":                                 1,
						},
					},
				},
			},
			expectedCount:     1,
			checkGatewayName:  "test-gateway",
			checkCustomerGwIp: "192.168.1.1",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := expandConnectorAdvIPSecGateway(tt.input)

			if tt.expectedCount == 0 {
				assert.Nil(t, result)
			} else {
				assert.Len(t, result, tt.expectedCount)

				if tt.checkGatewayName != "" {
					found := false
					for _, gateway := range result {
						if gateway.Name == tt.checkGatewayName && gateway.CustomerGwIp == tt.checkCustomerGwIp {
							found = true
							// Also check that tunnels were processed
							assert.NotNil(t, gateway.Tunnels)
							break
						}
					}
					assert.True(t, found, "Expected gateway with name %s and IP %s not found", tt.checkGatewayName, tt.checkCustomerGwIp)
				}
			}
		})
	}
}

func TestExpandConnectorAdvIPSecPolicyOptions(t *testing.T) {
	tests := []struct {
		name        string
		input       *schema.Set
		expected    *alkira.ConnectorAdvIPSecPolicyOptions
		expectError bool
	}{
		{
			name:        "nil input",
			input:       nil,
			expected:    &alkira.ConnectorAdvIPSecPolicyOptions{},
			expectError: false,
		},
		{
			name:        "empty input",
			input:       schema.NewSet(schema.HashString, []interface{}{}),
			expected:    &alkira.ConnectorAdvIPSecPolicyOptions{},
			expectError: false,
		},
		{
			name: "valid policy options",
			input: schema.NewSet(
				func(i interface{}) int {
					return schema.HashString("test")
				},
				[]interface{}{
					map[string]interface{}{
						"on_prem_prefix_list_ids": schema.NewSet(schema.HashInt, []interface{}{1, 2, 3}),
						"cxp_prefix_list_ids":     schema.NewSet(schema.HashInt, []interface{}{4, 5, 6}),
					},
				},
			),
			expected: &alkira.ConnectorAdvIPSecPolicyOptions{
				BranchTSPrefixListIds: []int{1, 2, 3},
				CxpTSPrefixListIds:    []int{4, 5, 6},
			},
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := expandConnectorAdvIPSecPolicyOptions(tt.input)

			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				if tt.expected.BranchTSPrefixListIds != nil || tt.expected.CxpTSPrefixListIds != nil {
					assert.ElementsMatch(t, tt.expected.BranchTSPrefixListIds, result.BranchTSPrefixListIds)
					assert.ElementsMatch(t, tt.expected.CxpTSPrefixListIds, result.CxpTSPrefixListIds)
				} else {
					assert.Equal(t, tt.expected, result)
				}
			}
		})
	}
}

func TestExpandConnectorAdvIPSecRoutingOptions(t *testing.T) {
	tests := []struct {
		name        string
		input       *schema.Set
		expected    *alkira.ConnectorAdvIPSecRoutingOptions
		expectError bool
		errorMsg    string
	}{
		{
			name:        "nil input",
			input:       nil,
			expected:    &alkira.ConnectorAdvIPSecRoutingOptions{},
			expectError: false,
		},
		{
			name:        "empty input",
			input:       schema.NewSet(schema.HashString, []interface{}{}),
			expected:    &alkira.ConnectorAdvIPSecRoutingOptions{},
			expectError: false,
		},
		{
			name: "valid static routing options",
			input: schema.NewSet(
				func(i interface{}) int {
					return schema.HashString("test")
				},
				[]interface{}{
					map[string]interface{}{
						"type":           "STATIC",
						"prefix_list_id": 123,
						"availability":   "HIGH",
					},
				},
			),
			expected: &alkira.ConnectorAdvIPSecRoutingOptions{
				StaticRouting: &alkira.ConnectorAdvIPSecStaticRouting{
					PrefixListId: 123,
					Availability: "HIGH",
				},
			},
			expectError: false,
		},
		{
			name: "valid dynamic routing options",
			input: schema.NewSet(
				func(i interface{}) int {
					return schema.HashString("test")
				},
				[]interface{}{
					map[string]interface{}{
						"type":                 "DYNAMIC",
						"customer_gateway_asn": "65001",
						"availability":         "HIGH",
						"bgp_auth_key":         "bgp-secret",
					},
				},
			),
			expected: &alkira.ConnectorAdvIPSecRoutingOptions{
				DynamicRouting: &alkira.ConnectorAdvIPSecDynamicRouting{
					CustomerGwAsn:    "65001",
					Availability:     "HIGH",
					BgpAuthKeyAlkira: "bgp-secret",
				},
			},
			expectError: false,
		},
		{
			name: "static routing without prefix_list_id - should error",
			input: schema.NewSet(
				func(i interface{}) int {
					return schema.HashString("test")
				},
				[]interface{}{
					map[string]interface{}{
						"type":         "STATIC",
						"availability": "HIGH",
					},
				},
			),
			expected:    nil,
			expectError: true,
			errorMsg:    "if STATIC routing type is specified, prefix_list_id is required",
		},
		{
			name: "dynamic routing without customer_gateway_asn - should error",
			input: schema.NewSet(
				func(i interface{}) int {
					return schema.HashString("test")
				},
				[]interface{}{
					map[string]interface{}{
						"type":         "DYNAMIC",
						"availability": "HIGH",
					},
				},
			),
			expected:    nil,
			expectError: true,
			errorMsg:    "if DYNAMIC routing type is specified, customer_gateway_asn is required",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := expandConnectorAdvIPSecRoutingOptions(tt.input)

			if tt.expectError {
				assert.Error(t, err)
				if tt.errorMsg != "" {
					assert.Contains(t, err.Error(), tt.errorMsg)
				}
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expected, result)
			}
		})
	}
}

func TestDeflateConnectorAdvIPSecTunnel(t *testing.T) {
	tests := []struct {
		name              string
		input             *alkira.ConnectorAdvIPSecTunnel
		expectNil         bool
		checkFields       map[string]interface{}
		hasAdvancedOption bool
	}{
		{
			name:      "nil input",
			input:     nil,
			expectNil: true, // Function returns nil for nil input
		},
		{
			name: "valid tunnel configuration with advanced options",
			input: &alkira.ConnectorAdvIPSecTunnel{
				Id:           "tunnel-1",
				PresharedKey: "secret123",
				ProfileId:    1,
				TunnelNo:     1,
				Advanced: &alkira.ConnectorAdvIPSecAdvanced{
					IkeVersion:      "2",
					Initiator:       true,
					RemoteAuthType:  "PSK",
					RemoteAuthValue: "secret123",
				},
			},
			expectNil: false,
			checkFields: map[string]interface{}{
				"id":            "tunnel-1",
				"preshared_key": "secret123",
				"profile_id":    1,
				"number":        1,
			},
			hasAdvancedOption: true,
		},
		{
			name: "tunnel without advanced options",
			input: &alkira.ConnectorAdvIPSecTunnel{
				Id:           "tunnel-2",
				PresharedKey: "secret456",
				ProfileId:    2,
				TunnelNo:     2,
			},
			expectNil: false,
			checkFields: map[string]interface{}{
				"id":            "tunnel-2",
				"preshared_key": "secret456",
				"profile_id":    2,
				"number":        2,
			},
			hasAdvancedOption: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := deflateConnectorAdvIPSecTunnel(tt.input)

			if tt.expectNil {
				assert.Nil(t, result)
				return
			}

			// Check basic fields
			for key, expectedValue := range tt.checkFields {
				assert.Equal(t, expectedValue, result[key], "Field %s should match", key)
			}

			// Check if advanced options are present when expected
			if tt.hasAdvancedOption {
				assert.Contains(t, result, "advanced_options")
				advOptions := result["advanced_options"]
				assert.NotNil(t, advOptions)
			}
		})
	}
}

func TestInputValidation(t *testing.T) {
	t.Run("test type assertions", func(t *testing.T) {
		// Test that the functions can handle various input formats without panicking
		validInput := []interface{}{
			map[string]interface{}{
				"ike_version": "2",
				"initiator":   true,
			},
		}

		result, err := expandConnectorAdvIPSecAdvancedOptions(validInput)
		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, "2", result.IkeVersion)
		assert.True(t, result.Initiator)
	})

	t.Run("test empty advanced options", func(t *testing.T) {
		tunnelInput := []interface{}{
			map[string]interface{}{
				"customer_end_overlay_ip_reservation_id": "overlay-ip-1",
				"cxp_end_overlay_ip_reservation_id":      "cxp-overlay-ip-1",
				"cxp_end_public_ip_reservation_id":       "public-ip-1",
				"id":                                     "tunnel-1",
				"preshared_key":                          "secret123",
				"profile_id":                             1,
				"number":                                 1,
				"advanced_options":                       []interface{}{},
			},
		}

		result := expandConnectorAdvIPSecTunnel(tunnelInput)
		assert.Len(t, result, 1)
		assert.Nil(t, result[0].Advanced)
	})
}
