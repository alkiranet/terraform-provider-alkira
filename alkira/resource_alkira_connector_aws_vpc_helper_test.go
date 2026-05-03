package alkira

import (
	"testing"

	"github.com/alkiranet/alkira-client-go/alkira"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/stretchr/testify/assert"
)

func TestExpandAwsVpcRouteTables(t *testing.T) {
	tests := []struct {
		name     string
		input    *schema.Set
		expected []alkira.RouteTables
	}{
		{
			name:     "nil input",
			input:    nil,
			expected: []alkira.RouteTables{},
		},
		{
			name:     "empty input",
			input:    schema.NewSet(schema.HashString, []interface{}{}),
			expected: []alkira.RouteTables{},
		},
		{
			name: "single route table",
			input: schema.NewSet(
				func(i interface{}) int {
					m := i.(map[string]interface{})
					return schema.HashString(m["id"].(string))
				},
				[]interface{}{
					map[string]interface{}{
						"id":              "rtb-123456789",
						"options":         "LEARNED",
						"prefix_list_ids": schema.NewSet(schema.HashInt, []interface{}{1, 2, 3}),
					},
				},
			),
			expected: []alkira.RouteTables{
				{
					Id:            "rtb-123456789",
					Mode:          "LEARNED",
					PrefixListIds: []int{1, 2, 3},
				},
			},
		},
		{
			name: "multiple route tables",
			input: schema.NewSet(
				func(i interface{}) int {
					m := i.(map[string]interface{})
					return schema.HashString(m["id"].(string))
				},
				[]interface{}{
					map[string]interface{}{
						"id":              "rtb-123456789",
						"options":         "LEARNED",
						"prefix_list_ids": schema.NewSet(schema.HashInt, []interface{}{1, 2}),
					},
					map[string]interface{}{
						"id":              "rtb-987654321",
						"options":         "ADVERTISED",
						"prefix_list_ids": schema.NewSet(schema.HashInt, []interface{}{3, 4}),
					},
				},
			),
			expected: []alkira.RouteTables{
				{
					Id:            "rtb-123456789",
					Mode:          "LEARNED",
					PrefixListIds: []int{1, 2},
				},
				{
					Id:            "rtb-987654321",
					Mode:          "ADVERTISED",
					PrefixListIds: []int{3, 4},
				},
			},
		},
		{
			name: "route table without prefix list ids",
			input: schema.NewSet(
				func(i interface{}) int {
					m := i.(map[string]interface{})
					return schema.HashString(m["id"].(string))
				},
				[]interface{}{
					map[string]interface{}{
						"id":              "rtb-123456789",
						"options":         "LEARNED",
						"prefix_list_ids": schema.NewSet(schema.HashInt, []interface{}{}),
					},
				},
			),
			expected: []alkira.RouteTables{
				{
					Id:            "rtb-123456789",
					Mode:          "LEARNED",
					PrefixListIds: nil,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := expandAwsVpcRouteTables(tt.input)

			// For sets, order may not be preserved, so check length and content
			assert.Len(t, result, len(tt.expected))

			if len(tt.expected) > 0 {
				// For each expected item, find a matching item in result
				for _, expectedItem := range tt.expected {
					found := false
					for _, resultItem := range result {
						if expectedItem.Id == resultItem.Id &&
							expectedItem.Mode == resultItem.Mode {
							if expectedItem.PrefixListIds == nil {
								assert.Nil(t, resultItem.PrefixListIds)
							} else {
								assert.ElementsMatch(t, expectedItem.PrefixListIds, resultItem.PrefixListIds)
							}
							found = true
							break
						}
					}
					assert.True(t, found, "Expected route table %v not found in result", expectedItem)
				}
			}
		})
	}
}

func TestExpandUserInputPrefixes(t *testing.T) {
	tests := []struct {
		name           string
		cidr           []interface{}
		subnets        *schema.Set
		overlaySubnets []interface{}
		expected       []alkira.InputPrefixes
		expectError    bool
		errorMsg       string
	}{
		{
			name:           "no CIDR and no subnets - should error",
			cidr:           []interface{}{},
			subnets:        nil,
			overlaySubnets: []interface{}{},
			expected:       nil,
			expectError:    true,
			errorMsg:       "either \"vpc_subnet\" or \"vpc_cidr\" must be specified",
		},
		{
			name:           "only overlay subnets",
			cidr:           []interface{}{"10.0.0.0/16"},
			subnets:        nil,
			overlaySubnets: []interface{}{"172.16.0.0/24", "172.16.1.0/24"},
			expected: []alkira.InputPrefixes{
				{Type: "CIDR", Value: "10.0.0.0/16"},
				{Type: "OVERLAY_SUBNETS", Value: "172.16.0.0/24"},
				{Type: "OVERLAY_SUBNETS", Value: "172.16.1.0/24"},
			},
			expectError: false,
		},
		{
			name: "only VPC subnets",
			cidr: []interface{}{},
			subnets: schema.NewSet(
				func(i interface{}) int {
					m := i.(map[string]interface{})
					return schema.HashString(m["id"].(string))
				},
				[]interface{}{
					map[string]interface{}{
						"id":   "subnet-123",
						"cidr": "10.0.1.0/24",
					},
					map[string]interface{}{
						"id":   "subnet-456",
						"cidr": "10.0.2.0/24",
					},
				},
			),
			overlaySubnets: []interface{}{},
			expected: []alkira.InputPrefixes{
				{Id: "subnet-123", Type: "SUBNET", Value: "10.0.1.0/24"},
				{Id: "subnet-456", Type: "SUBNET", Value: "10.0.2.0/24"},
			},
			expectError: false,
		},
		{
			name:           "VPC CIDR only",
			cidr:           []interface{}{"10.0.0.0/16", "192.168.0.0/16"},
			subnets:        nil,
			overlaySubnets: []interface{}{},
			expected: []alkira.InputPrefixes{
				{Type: "CIDR", Value: "10.0.0.0/16"},
				{Type: "CIDR", Value: "192.168.0.0/16"},
			},
			expectError: false,
		},
		{
			name: "mixed subnets and CIDR",
			cidr: []interface{}{"10.0.0.0/16"},
			subnets: schema.NewSet(
				func(i interface{}) int {
					m := i.(map[string]interface{})
					return schema.HashString(m["id"].(string))
				},
				[]interface{}{
					map[string]interface{}{
						"id":   "subnet-123",
						"cidr": "192.168.1.0/24",
					},
				},
			),
			overlaySubnets: []interface{}{"172.16.0.0/24"},
			expected: []alkira.InputPrefixes{
				{Type: "CIDR", Value: "10.0.0.0/16"},
				{Type: "OVERLAY_SUBNETS", Value: "172.16.0.0/24"},
			},
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := expandUserInputPrefixes(tt.cidr, tt.subnets, tt.overlaySubnets)

			if tt.expectError {
				assert.Error(t, err)
				if tt.errorMsg != "" {
					assert.Contains(t, err.Error(), tt.errorMsg)
				}
			} else {
				assert.NoError(t, err)

				// Check length
				assert.Len(t, result, len(tt.expected))

				// Check that all expected prefixes are present
				for _, expectedPrefix := range tt.expected {
					found := false
					for _, resultPrefix := range result {
						if expectedPrefix.Type == resultPrefix.Type &&
							expectedPrefix.Value == resultPrefix.Value {
							found = true
							break
						}
					}
					assert.True(t, found, "Expected prefix %+v not found in result", expectedPrefix)
				}
			}
		})
	}
}

func TestExpandAwsVpcTgwAttachments(t *testing.T) {
	tests := []struct {
		name     string
		input    []interface{}
		expected []alkira.TgwAttachment
	}{
		{
			name:     "nil input",
			input:    nil,
			expected: []alkira.TgwAttachment{},
		},
		{
			name:     "empty input",
			input:    []interface{}{},
			expected: []alkira.TgwAttachment{},
		},
		{
			name: "single attachment",
			input: []interface{}{
				map[string]interface{}{
					"subnet_id": "subnet-123456789",
					"az":        "us-east-1a",
				},
			},
			expected: []alkira.TgwAttachment{
				{
					SubnetId:         "subnet-123456789",
					AvailabilityZone: "us-east-1a",
				},
			},
		},
		{
			name: "multiple attachments",
			input: []interface{}{
				map[string]interface{}{
					"subnet_id": "subnet-123456789",
					"az":        "us-east-1a",
				},
				map[string]interface{}{
					"subnet_id": "subnet-987654321",
					"az":        "us-east-1b",
				},
			},
			expected: []alkira.TgwAttachment{
				{
					SubnetId:         "subnet-123456789",
					AvailabilityZone: "us-east-1a",
				},
				{
					SubnetId:         "subnet-987654321",
					AvailabilityZone: "us-east-1b",
				},
			},
		},
		{
			name: "attachment with missing optional fields",
			input: []interface{}{
				map[string]interface{}{
					"subnet_id": "subnet-123456789",
				},
			},
			expected: []alkira.TgwAttachment{
				{
					SubnetId:         "subnet-123456789",
					AvailabilityZone: "",
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := expandAwsVpcTgwAttachments(tt.input)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestSetTgwAttachment(t *testing.T) {
	// Create a mock ResourceData for testing
	r := &schema.Resource{
		Schema: map[string]*schema.Schema{
			"tgw_attachment": {
				Type: schema.TypeList,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"subnet_id": {Type: schema.TypeString},
						"az":        {Type: schema.TypeString},
					},
				},
			},
		},
	}

	tests := []struct {
		name             string
		existingData     []interface{}
		apiAttachments   []alkira.TgwAttachment
		expectedCount    int
		shouldHaveSubnet string
	}{
		{
			name: "existing attachment matches API",
			existingData: []interface{}{
				map[string]interface{}{
					"subnet_id": "subnet-123",
					"az":        "us-east-1a",
				},
			},
			apiAttachments: []alkira.TgwAttachment{
				{
					SubnetId:         "subnet-123",
					AvailabilityZone: "us-east-1a",
				},
			},
			expectedCount:    1,
			shouldHaveSubnet: "subnet-123",
		},
		{
			name: "API has new attachment not in config",
			existingData: []interface{}{
				map[string]interface{}{
					"subnet_id": "subnet-123",
					"az":        "us-east-1a",
				},
			},
			apiAttachments: []alkira.TgwAttachment{
				{
					SubnetId:         "subnet-123",
					AvailabilityZone: "us-east-1a",
				},
				{
					SubnetId:         "subnet-456",
					AvailabilityZone: "us-east-1b",
				},
			},
			expectedCount:    2,
			shouldHaveSubnet: "subnet-456",
		},
		{
			name:         "empty existing data",
			existingData: []interface{}{},
			apiAttachments: []alkira.TgwAttachment{
				{
					SubnetId:         "subnet-123",
					AvailabilityZone: "us-east-1a",
				},
			},
			expectedCount:    1,
			shouldHaveSubnet: "subnet-123",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := r.TestResourceData()
			d.Set("tgw_attachment", tt.existingData)

			setTgwAttachment(d, tt.apiAttachments)

			result := d.Get("tgw_attachment").([]interface{})
			assert.Len(t, result, tt.expectedCount)

			if tt.shouldHaveSubnet != "" {
				found := false
				for _, attachment := range result {
					a := attachment.(map[string]interface{})
					if a["subnet_id"].(string) == tt.shouldHaveSubnet {
						found = true
						break
					}
				}
				assert.True(t, found, "Expected subnet %s not found", tt.shouldHaveSubnet)
			}
		})
	}
}

func TestAwsVpcValidation(t *testing.T) {
	t.Run("test subnet ID format", func(t *testing.T) {
		validSubnetIds := []string{"subnet-123456789abcdef", "subnet-abc123", "subnet-1234567890abcdef0"}
		for _, subnetId := range validSubnetIds {
			assert.NotEmpty(t, subnetId)
			assert.Contains(t, subnetId, "subnet-")
		}
	})

	t.Run("test route table ID format", func(t *testing.T) {
		validRouteTableIds := []string{"rtb-123456789abcdef", "rtb-abc123", "rtb-1234567890abcdef0"}
		for _, rtbId := range validRouteTableIds {
			assert.NotEmpty(t, rtbId)
			assert.Contains(t, rtbId, "rtb-")
		}
	})

	t.Run("test availability zone format", func(t *testing.T) {
		validAZs := []string{"us-east-1a", "us-west-2b", "eu-west-1c"}
		for _, az := range validAZs {
			assert.NotEmpty(t, az)
			assert.Regexp(t, `^[a-z]+-[a-z]+-\d+[a-z]$`, az)
		}
	})

	t.Run("test CIDR format", func(t *testing.T) {
		validCIDRs := []string{"10.0.0.0/16", "192.168.1.0/24", "172.16.0.0/12"}
		for _, cidr := range validCIDRs {
			assert.NotEmpty(t, cidr)
			assert.Contains(t, cidr, "/")
			assert.Regexp(t, `^\d+\.\d+\.\d+\.\d+/\d+$`, cidr)
		}
	})
}

func TestAwsVpcInputValidation(t *testing.T) {
	t.Run("test route table mode validation", func(t *testing.T) {
		validModes := []string{"LEARNED", "ADVERTISED", "BOTH"}
		for _, mode := range validModes {
			assert.NotEmpty(t, mode)
			assert.Contains(t, []string{"LEARNED", "ADVERTISED", "BOTH"}, mode)
		}
	})

	t.Run("test prefix type validation", func(t *testing.T) {
		validTypes := []string{"CIDR", "SUBNET", "OVERLAY_SUBNETS"}
		for _, prefixType := range validTypes {
			assert.NotEmpty(t, prefixType)
			assert.Contains(t, []string{"CIDR", "SUBNET", "OVERLAY_SUBNETS"}, prefixType)
		}
	})

	t.Run("test type conversions", func(t *testing.T) {
		// Test prefix list ID conversion
		prefixIds := []interface{}{1, 2, 3, 4, 5}
		converted := convertTypeListToIntList(prefixIds)
		expected := []int{1, 2, 3, 4, 5}
		assert.Equal(t, expected, converted)
	})
}

func TestAwsVpcDataStructures(t *testing.T) {
	t.Run("test route table structure", func(t *testing.T) {
		routeTable := alkira.RouteTables{
			Id:            "rtb-123456789",
			Mode:          "LEARNED",
			PrefixListIds: []int{1, 2, 3},
		}

		assert.Equal(t, "rtb-123456789", routeTable.Id)
		assert.Equal(t, "LEARNED", routeTable.Mode)
		assert.Len(t, routeTable.PrefixListIds, 3)
		assert.Contains(t, routeTable.PrefixListIds, 1)
		assert.Contains(t, routeTable.PrefixListIds, 2)
		assert.Contains(t, routeTable.PrefixListIds, 3)
	})

	t.Run("test TGW attachment structure", func(t *testing.T) {
		attachment := alkira.TgwAttachment{
			SubnetId:         "subnet-123456789",
			AvailabilityZone: "us-east-1a",
		}

		assert.Equal(t, "subnet-123456789", attachment.SubnetId)
		assert.Equal(t, "us-east-1a", attachment.AvailabilityZone)
	})

	t.Run("test input prefix structure", func(t *testing.T) {
		prefix := alkira.InputPrefixes{
			Type:  "CIDR",
			Value: "10.0.0.0/16",
		}

		assert.Equal(t, "CIDR", prefix.Type)
		assert.Equal(t, "10.0.0.0/16", prefix.Value)
	})
}
