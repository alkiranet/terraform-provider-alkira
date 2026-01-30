package alkira

import (
	"testing"

	"github.com/alkiranet/alkira-client-go/alkira"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/stretchr/testify/assert"
)

func TestGenerateGCPUserInputPrefixes(t *testing.T) {
	tests := []struct {
		name        string
		subnets     *schema.Set
		expected    []alkira.UserInputPrefixes
		expectError bool
		errorMsg    string
	}{
		{
			name:        "nil subnets",
			subnets:     nil,
			expected:    nil,
			expectError: false,
		},
		{
			name:        "empty subnets",
			subnets:     schema.NewSet(schema.HashString, []interface{}{}),
			expected:    nil,
			expectError: false,
		},
		{
			name: "single subnet",
			subnets: schema.NewSet(
				func(i interface{}) int {
					m := i.(map[string]interface{})
					return schema.HashString(m["id"].(string))
				},
				[]interface{}{
					map[string]interface{}{
						"id":   "projects/test/regions/us-central1/subnetworks/subnet-1",
						"cidr": "10.0.1.0/24",
					},
				},
			),
			expected: []alkira.UserInputPrefixes{
				{
					FqId:  "projects/test/regions/us-central1/subnetworks/subnet-1",
					Value: "10.0.1.0/24",
					Type:  "SUBNET",
				},
			},
			expectError: false,
		},
		{
			name: "multiple subnets",
			subnets: schema.NewSet(
				func(i interface{}) int {
					m := i.(map[string]interface{})
					return schema.HashString(m["id"].(string))
				},
				[]interface{}{
					map[string]interface{}{
						"id":   "projects/test/regions/us-central1/subnetworks/subnet-1",
						"cidr": "10.0.1.0/24",
					},
					map[string]interface{}{
						"id":   "projects/test/regions/us-east1/subnetworks/subnet-2",
						"cidr": "10.0.2.0/24",
					},
				},
			),
			expected: []alkira.UserInputPrefixes{
				{
					FqId:  "projects/test/regions/us-central1/subnetworks/subnet-1",
					Value: "10.0.1.0/24",
					Type:  "SUBNET",
				},
				{
					FqId:  "projects/test/regions/us-east1/subnetworks/subnet-2",
					Value: "10.0.2.0/24",
					Type:  "SUBNET",
				},
			},
			expectError: false,
		},
		{
			name: "empty id field - should error",
			subnets: schema.NewSet(
				func(i interface{}) int {
					m := i.(map[string]interface{})
					if id, ok := m["id"].(string); ok {
						return schema.HashString(id)
					}
					return schema.HashString("")
				},
				[]interface{}{
					map[string]interface{}{
						"id":   "",
						"cidr": "10.0.1.0/24",
					},
				},
			),
			expected:    nil,
			expectError: true,
			errorMsg:    "both subnetwork ID",
		},
		{
			name: "empty cidr field - should error",
			subnets: schema.NewSet(
				func(i interface{}) int {
					m := i.(map[string]interface{})
					if id, ok := m["id"].(string); ok {
						return schema.HashString(id)
					}
					return schema.HashString("")
				},
				[]interface{}{
					map[string]interface{}{
						"id":   "projects/test/regions/us-central1/subnetworks/subnet-1",
						"cidr": "",
					},
				},
			),
			expected:    nil,
			expectError: true,
			errorMsg:    "both subnetwork ID",
		},
		{
			name: "empty id field - should error",
			subnets: schema.NewSet(
				func(i interface{}) int {
					m := i.(map[string]interface{})
					if id, ok := m["id"].(string); ok {
						return schema.HashString(id)
					}
					return schema.HashString("")
				},
				[]interface{}{
					map[string]interface{}{
						"id":   "",
						"cidr": "10.0.1.0/24",
					},
				},
			),
			expected:    nil,
			expectError: true,
			errorMsg:    "both subnetwork ID",
		},
		{
			name: "empty cidr field - should error",
			subnets: schema.NewSet(
				func(i interface{}) int {
					m := i.(map[string]interface{})
					if id, ok := m["id"].(string); ok {
						return schema.HashString(id)
					}
					return schema.HashString("")
				},
				[]interface{}{
					map[string]interface{}{
						"id":   "projects/test/regions/us-central1/subnetworks/subnet-1",
						"cidr": "",
					},
				},
			),
			expected:    nil,
			expectError: true,
			errorMsg:    "both subnetwork ID",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := generateGCPUserInputPrefixes(tt.subnets)

			if tt.expectError {
				assert.Error(t, err)
				if tt.errorMsg != "" {
					assert.Contains(t, err.Error(), tt.errorMsg)
				}
			} else {
				assert.NoError(t, err)
				// For sets, check length and content (order may vary)
				assert.Len(t, result, len(tt.expected))
				for _, expectedPrefix := range tt.expected {
					found := false
					for _, resultPrefix := range result {
						if expectedPrefix.FqId == resultPrefix.FqId &&
							expectedPrefix.Value == resultPrefix.Value &&
							expectedPrefix.Type == resultPrefix.Type {
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

func TestSetGcpRoutingOptions(t *testing.T) {
	r := &schema.Resource{
		Schema: map[string]*schema.Schema{
			"gcp_routing": {
				Type: schema.TypeList,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"prefix_list_ids": {Type: schema.TypeList, Elem: &schema.Schema{Type: schema.TypeInt}},
						"custom_prefix":   {Type: schema.TypeString},
					},
				},
			},
		},
	}

	tests := []struct {
		name           string
		gcpRouting     *alkira.ConnectorGcpVpcRouting
		expectEmpty    bool
		expectedPrefix string
		expectedIds    []int
	}{
		{
			name:        "nil routing - should not set",
			gcpRouting:  nil,
			expectEmpty: true,
		},
		{
			name: "default route mode",
			gcpRouting: &alkira.ConnectorGcpVpcRouting{
				ImportOptions: alkira.ConnectorGcpVpcImportOptions{
					RouteImportMode: "ADVERTISE_DEFAULT_ROUTE",
					PrefixListIds:   []int{},
				},
			},
			expectEmpty:    false,
			expectedPrefix: "ADVERTISE_DEFAULT_ROUTE",
			expectedIds:    nil, // convertTypeListToIntList returns nil for empty slices
		},
		{
			name: "custom prefix mode with prefix lists",
			gcpRouting: &alkira.ConnectorGcpVpcRouting{
				ImportOptions: alkira.ConnectorGcpVpcImportOptions{
					RouteImportMode: "ADVERTISE_CUSTOM_PREFIX",
					PrefixListIds:   []int{1, 2, 3},
				},
			},
			expectEmpty:    false,
			expectedPrefix: "ADVERTISE_CUSTOM_PREFIX",
			expectedIds:    []int{1, 2, 3},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := r.TestResourceData()
			setGcpRoutingOptions(tt.gcpRouting, d)

			result := d.Get("gcp_routing").([]interface{})

			if tt.expectEmpty {
				assert.Empty(t, result)
			} else {
				assert.Len(t, result, 1)
				routing := result[0].(map[string]interface{})
				assert.Equal(t, tt.expectedPrefix, routing["custom_prefix"])
				// Convert []interface{} to []int for comparison
				actualIds := convertTypeListToIntList(routing["prefix_list_ids"].([]interface{}))
				assert.Equal(t, tt.expectedIds, actualIds)
			}
		})
	}
}

func TestSetGcpVpcSubnets(t *testing.T) {
	r := &schema.Resource{
		Schema: map[string]*schema.Schema{
			"vpc_subnet": {
				Type: schema.TypeSet,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id":   {Type: schema.TypeString},
						"cidr": {Type: schema.TypeString},
					},
				},
			},
		},
	}

	tests := []struct {
		name          string
		gcpRouting    *alkira.ConnectorGcpVpcRouting
		expectEmpty   bool
		expectedIds   []string
		expectedCidrs []string
	}{
		{
			name:        "nil routing - should not set",
			gcpRouting:  nil,
			expectEmpty: true,
		},
		{
			name: "empty export options",
			gcpRouting: &alkira.ConnectorGcpVpcRouting{
				ExportOptions: alkira.ConnectorGcpVpcExportOptions{
					ExportAllSubnets: true,
					Prefixes:         []alkira.UserInputPrefixes{},
				},
			},
			expectEmpty: true,
		},
		{
			name: "single subnet",
			gcpRouting: &alkira.ConnectorGcpVpcRouting{
				ExportOptions: alkira.ConnectorGcpVpcExportOptions{
					ExportAllSubnets: false,
					Prefixes: []alkira.UserInputPrefixes{
						{
							FqId:  "projects/test/regions/us-central1/subnetworks/subnet-1",
							Value: "10.0.1.0/24",
							Type:  "SUBNET",
						},
					},
				},
			},
			expectEmpty:   false,
			expectedIds:   []string{"projects/test/regions/us-central1/subnetworks/subnet-1"},
			expectedCidrs: []string{"10.0.1.0/24"},
		},
		{
			name: "multiple subnets",
			gcpRouting: &alkira.ConnectorGcpVpcRouting{
				ExportOptions: alkira.ConnectorGcpVpcExportOptions{
					ExportAllSubnets: false,
					Prefixes: []alkira.UserInputPrefixes{
						{
							FqId:  "projects/test/regions/us-central1/subnetworks/subnet-1",
							Value: "10.0.1.0/24",
							Type:  "SUBNET",
						},
						{
							FqId:  "projects/test/regions/us-east1/subnetworks/subnet-2",
							Value: "10.0.2.0/24",
							Type:  "SUBNET",
						},
					},
				},
			},
			expectEmpty:   false,
			expectedIds:   []string{"projects/test/regions/us-central1/subnetworks/subnet-1", "projects/test/regions/us-east1/subnetworks/subnet-2"},
			expectedCidrs: []string{"10.0.1.0/24", "10.0.2.0/24"},
		},
		{
			name: "nil prefixes",
			gcpRouting: &alkira.ConnectorGcpVpcRouting{
				ExportOptions: alkira.ConnectorGcpVpcExportOptions{
					ExportAllSubnets: true,
					Prefixes:         nil,
				},
			},
			expectEmpty: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := r.TestResourceData()
			setGcpVpcSubnets(tt.gcpRouting, d)

			result := d.Get("vpc_subnet").(*schema.Set)

			if tt.expectEmpty {
				assert.Empty(t, result.List())
			} else {
				assert.Len(t, result.List(), len(tt.expectedIds))
				for _, subnet := range result.List() {
					s := subnet.(map[string]interface{})
					assert.Contains(t, tt.expectedIds, s["id"].(string))
					assert.Contains(t, tt.expectedCidrs, s["cidr"].(string))
				}
			}
		})
	}
}

func TestExpandGcpRouting(t *testing.T) {
	tests := []struct {
		name        string
		gcpRouting  []interface{}
		subnets     *schema.Set
		expectError bool
		expected    *alkira.ConnectorGcpVpcRouting
	}{
		{
			name:        "nil routing with nil subnets",
			gcpRouting:  nil,
			subnets:     nil,
			expectError: false,
			expected: &alkira.ConnectorGcpVpcRouting{
				ExportOptions: alkira.ConnectorGcpVpcExportOptions{
					ExportAllSubnets: true,
					Prefixes:         nil,
				},
				ImportOptions: alkira.ConnectorGcpVpcImportOptions{
					RouteImportMode: "ADVERTISE_DEFAULT_ROUTE",
					PrefixListIds:   nil,
				},
			},
		},
		{
			name: "default route mode with subnets",
			gcpRouting: []interface{}{
				map[string]interface{}{
					"custom_prefix":   "ADVERTISE_DEFAULT_ROUTE",
					"prefix_list_ids": []interface{}{},
				},
			},
			subnets: schema.NewSet(
				func(i interface{}) int {
					m := i.(map[string]interface{})
					return schema.HashString(m["id"].(string))
				},
				[]interface{}{
					map[string]interface{}{
						"id":   "projects/test/regions/us-central1/subnetworks/subnet-1",
						"cidr": "10.0.1.0/24",
					},
				},
			),
			expectError: false,
			expected: &alkira.ConnectorGcpVpcRouting{
				ExportOptions: alkira.ConnectorGcpVpcExportOptions{
					ExportAllSubnets: false,
					Prefixes: []alkira.UserInputPrefixes{
						{
							FqId:  "projects/test/regions/us-central1/subnetworks/subnet-1",
							Value: "10.0.1.0/24",
							Type:  "SUBNET",
						},
					},
				},
				ImportOptions: alkira.ConnectorGcpVpcImportOptions{
					RouteImportMode: "ADVERTISE_DEFAULT_ROUTE",
					PrefixListIds:   nil, // convertTypeListToIntList returns nil for empty slices
				},
			},
		},
		{
			name: "custom prefix mode with prefix lists",
			gcpRouting: []interface{}{
				map[string]interface{}{
					"custom_prefix":   "ADVERTISE_CUSTOM_PREFIX",
					"prefix_list_ids": []interface{}{1, 2, 3},
				},
			},
			subnets:     nil,
			expectError: false,
			expected: &alkira.ConnectorGcpVpcRouting{
				ExportOptions: alkira.ConnectorGcpVpcExportOptions{
					ExportAllSubnets: true,
					Prefixes:         nil,
				},
				ImportOptions: alkira.ConnectorGcpVpcImportOptions{
					RouteImportMode: "ADVERTISE_CUSTOM_PREFIX",
					PrefixListIds:   []int{1, 2, 3},
				},
			},
		},
		{
			name: "invalid subnet data - should error",
			gcpRouting: []interface{}{
				map[string]interface{}{
					"custom_prefix":   "ADVERTISE_DEFAULT_ROUTE",
					"prefix_list_ids": []interface{}{},
				},
			},
			subnets: schema.NewSet(
				func(i interface{}) int {
					m := i.(map[string]interface{})
					if id, ok := m["id"].(string); ok {
						return schema.HashString(id)
					}
					return schema.HashString("")
				},
				[]interface{}{
					map[string]interface{}{
						"id":   "", // empty id
						"cidr": "10.0.1.0/24",
					},
				},
			),
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := expandGcpRouting(tt.gcpRouting, tt.subnets)

			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expected.ImportOptions.RouteImportMode, result.ImportOptions.RouteImportMode)
				assert.Equal(t, tt.expected.ImportOptions.PrefixListIds, result.ImportOptions.PrefixListIds)
				assert.Equal(t, tt.expected.ExportOptions.ExportAllSubnets, result.ExportOptions.ExportAllSubnets)
				assert.Equal(t, tt.expected.ExportOptions.Prefixes, result.ExportOptions.Prefixes)
			}
		})
	}
}

func TestGcpVpcDataStructures(t *testing.T) {
	t.Run("test UserInputPrefixes structure", func(t *testing.T) {
		prefix := alkira.UserInputPrefixes{
			FqId:  "projects/test/regions/us-central1/subnetworks/subnet-1",
			Value: "10.0.1.0/24",
			Type:  "SUBNET",
		}

		assert.Equal(t, "projects/test/regions/us-central1/subnetworks/subnet-1", prefix.FqId)
		assert.Equal(t, "10.0.1.0/24", prefix.Value)
		assert.Equal(t, "SUBNET", prefix.Type)
	})

	t.Run("test ConnectorGcpVpcExportOptions structure", func(t *testing.T) {
		options := alkira.ConnectorGcpVpcExportOptions{
			ExportAllSubnets: false,
			Prefixes: []alkira.UserInputPrefixes{
				{
					FqId:  "projects/test/regions/us-central1/subnetworks/subnet-1",
					Value: "10.0.1.0/24",
					Type:  "SUBNET",
				},
			},
		}

		assert.False(t, options.ExportAllSubnets)
		assert.Len(t, options.Prefixes, 1)
		assert.Equal(t, "SUBNET", options.Prefixes[0].Type)
	})

	t.Run("test ConnectorGcpVpcImportOptions structure", func(t *testing.T) {
		options := alkira.ConnectorGcpVpcImportOptions{
			RouteImportMode: "ADVERTISE_CUSTOM_PREFIX",
			PrefixListIds:   []int{1, 2, 3},
		}

		assert.Equal(t, "ADVERTISE_CUSTOM_PREFIX", options.RouteImportMode)
		assert.Equal(t, []int{1, 2, 3}, options.PrefixListIds)
	})

	t.Run("test ConnectorGcpVpcRouting structure", func(t *testing.T) {
		routing := &alkira.ConnectorGcpVpcRouting{
			ExportOptions: alkira.ConnectorGcpVpcExportOptions{
				ExportAllSubnets: true,
				Prefixes:         nil,
			},
			ImportOptions: alkira.ConnectorGcpVpcImportOptions{
				RouteImportMode: "ADVERTISE_DEFAULT_ROUTE",
				PrefixListIds:   nil,
			},
		}

		assert.True(t, routing.ExportOptions.ExportAllSubnets)
		assert.Equal(t, "ADVERTISE_DEFAULT_ROUTE", routing.ImportOptions.RouteImportMode)
		assert.Nil(t, routing.ExportOptions.Prefixes)
		assert.Nil(t, routing.ImportOptions.PrefixListIds)
	})
}
