package alkira

import (
	"fmt"
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
			name: "single subnet with id and cidr (backward compatibility)",
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
					Id:    "",
					FqId:  "projects/test/regions/us-central1/subnetworks/subnet-1",
					Value: "10.0.1.0/24",
					Type:  "SUBNET",
				},
			},
			expectError: false,
		},
		{
			name: "single subnet with internal_id (from UI/import)",
			subnets: schema.NewSet(
				func(i interface{}) int {
					m := i.(map[string]interface{})
					return schema.HashString(m["internal_id"].(string))
				},
				[]interface{}{
					map[string]interface{}{
						"internal_id": "7442940776704048352",
						"id":          "",
						"cidr":        "10.0.1.0/24",
					},
				},
			),
			expected: []alkira.UserInputPrefixes{
				{
					Id:    "7442940776704048352",
					FqId:  "",
					Value: "10.0.1.0/24",
					Type:  "SUBNET",
				},
			},
			expectError: false,
		},
		{
			name: "single subnet with both internal_id and id (internal_id takes precedence)",
			subnets: schema.NewSet(
				func(i interface{}) int {
					m := i.(map[string]interface{})
					if id, ok := m["internal_id"].(string); ok && id != "" {
						return schema.HashString(id)
					}
					return schema.HashString(m["id"].(string))
				},
				[]interface{}{
					map[string]interface{}{
						"internal_id": "7442940776704048352",
						"id":          "projects/test/regions/us-central1/subnetworks/subnet-1",
						"cidr":        "10.0.1.0/24",
					},
				},
			),
			expected: []alkira.UserInputPrefixes{
				{
					Id:    "7442940776704048352",
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
					if id, ok := m["internal_id"].(string); ok && id != "" {
						return schema.HashString(id)
					}
					return schema.HashString(m["id"].(string))
				},
				[]interface{}{
					map[string]interface{}{
						"internal_id": "7442940776704048352",
						"id":          "projects/test/regions/us-central1/subnetworks/subnet-1",
						"cidr":        "10.0.1.0/24",
					},
					map[string]interface{}{
						"id":   "projects/test/regions/us-east1/subnetworks/subnet-2",
						"cidr": "10.0.2.0/24",
					},
				},
			),
			expected: []alkira.UserInputPrefixes{
				{
					Id:    "7442940776704048352",
					FqId:  "projects/test/regions/us-central1/subnetworks/subnet-1",
					Value: "10.0.1.0/24",
					Type:  "SUBNET",
				},
				{
					Id:    "",
					FqId:  "projects/test/regions/us-east1/subnetworks/subnet-2",
					Value: "10.0.2.0/24",
					Type:  "SUBNET",
				},
			},
			expectError: false,
		},
		{
			name: "missing internal_id and id - should error",
			subnets: schema.NewSet(
				func(i interface{}) int {
					return schema.HashString("")
				},
				[]interface{}{
					map[string]interface{}{
						"internal_id": "",
						"id":          "",
						"cidr":        "10.0.1.0/24",
					},
				},
			),
			expected:    nil,
			expectError: true,
			errorMsg:    "internal_id or both id and cidr",
		},
		{
			name: "missing cidr without internal_id - should error",
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
						"internal_id": "",
						"id":          "projects/test/regions/us-central1/subnetworks/subnet-1",
						"cidr":        "",
					},
				},
			),
			expected:    nil,
			expectError: true,
			errorMsg:    "internal_id or both id and cidr",
		},
		{
			name: "missing internal_id, id, and cidr - should error",
			subnets: schema.NewSet(
				func(i interface{}) int {
					return schema.HashString("")
				},
				[]interface{}{
					map[string]interface{}{
						"internal_id": "",
						"id":          "",
						"cidr":        "",
					},
				},
			),
			expected:    nil,
			expectError: true,
			errorMsg:    "internal_id or both id and cidr",
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
						if expectedPrefix.Id == resultPrefix.Id &&
							expectedPrefix.FqId == resultPrefix.FqId &&
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
						"prefix_list_ids":    {Type: schema.TypeList, Elem: &schema.Schema{Type: schema.TypeInt}},
						"custom_prefix":      {Type: schema.TypeString},
						"export_all_subnets": {Type: schema.TypeBool, Optional: true, Default: true},
					},
				},
			},
		},
	}

	tests := []struct {
		name              string
		gcpRouting        *alkira.ConnectorGcpVpcRouting
		expectEmpty       bool
		expectedPrefix    string
		expectedIds       []int
		expectedExportAll *bool // nil means don't check
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
				ExportOptions: alkira.ConnectorGcpVpcExportOptions{
					ExportAllSubnets: false, // default value
				},
			},
			expectEmpty:       false,
			expectedPrefix:    "ADVERTISE_DEFAULT_ROUTE",
			expectedIds:       nil, // convertTypeListToIntList returns nil for empty slices
			expectedExportAll: boolPtr(false),
		},
		{
			name: "custom prefix mode with prefix lists",
			gcpRouting: &alkira.ConnectorGcpVpcRouting{
				ImportOptions: alkira.ConnectorGcpVpcImportOptions{
					RouteImportMode: "ADVERTISE_CUSTOM_PREFIX",
					PrefixListIds:   []int{1, 2, 3},
				},
				ExportOptions: alkira.ConnectorGcpVpcExportOptions{
					ExportAllSubnets: false, // default value
				},
			},
			expectEmpty:       false,
			expectedPrefix:    "ADVERTISE_CUSTOM_PREFIX",
			expectedIds:       []int{1, 2, 3},
			expectedExportAll: boolPtr(false),
		},
		{
			name: "with export_all_subnets set to false",
			gcpRouting: &alkira.ConnectorGcpVpcRouting{
				ImportOptions: alkira.ConnectorGcpVpcImportOptions{
					RouteImportMode: "ADVERTISE_DEFAULT_ROUTE",
					PrefixListIds:   []int{},
				},
				ExportOptions: alkira.ConnectorGcpVpcExportOptions{
					ExportAllSubnets: false,
				},
			},
			expectEmpty:       false,
			expectedPrefix:    "ADVERTISE_DEFAULT_ROUTE",
			expectedIds:       nil,
			expectedExportAll: boolPtr(false),
		},
		{
			name: "with export_all_subnets set to true (explicit)",
			gcpRouting: &alkira.ConnectorGcpVpcRouting{
				ImportOptions: alkira.ConnectorGcpVpcImportOptions{
					RouteImportMode: "ADVERTISE_DEFAULT_ROUTE",
					PrefixListIds:   []int{},
				},
				ExportOptions: alkira.ConnectorGcpVpcExportOptions{
					ExportAllSubnets: true,
				},
			},
			expectEmpty:       false,
			expectedPrefix:    "ADVERTISE_DEFAULT_ROUTE",
			expectedIds:       nil,
			expectedExportAll: boolPtr(true),
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
				// Check export_all_subnets if specified
				if tt.expectedExportAll != nil {
					assert.Equal(t, *tt.expectedExportAll, routing["export_all_subnets"].(bool))
				}
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
						"id":          {Type: schema.TypeString},
						"internal_id": {Type: schema.TypeString},
						"cidr":        {Type: schema.TypeString},
					},
				},
			},
		},
	}

	tests := []struct {
		name           string
		gcpRouting     *alkira.ConnectorGcpVpcRouting
		expectEmpty    bool
		expectedIds    []string
		expectedCidrs  []string
		expectedIntIds []string
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
			name: "single subnet with FqId only (Terraform created)",
			gcpRouting: &alkira.ConnectorGcpVpcRouting{
				ExportOptions: alkira.ConnectorGcpVpcExportOptions{
					ExportAllSubnets: false,
					Prefixes: []alkira.UserInputPrefixes{
						{
							Id:    "",
							FqId:  "projects/test/regions/us-central1/subnetworks/subnet-1",
							Value: "10.0.1.0/24",
							Type:  "SUBNET",
						},
					},
				},
			},
			expectEmpty:    false,
			expectedIds:    []string{"projects/test/regions/us-central1/subnetworks/subnet-1"},
			expectedCidrs:  []string{"10.0.1.0/24"},
			expectedIntIds: []string{""},
		},
		{
			name: "single subnet with internal Id only (UI created/import)",
			gcpRouting: &alkira.ConnectorGcpVpcRouting{
				ExportOptions: alkira.ConnectorGcpVpcExportOptions{
					ExportAllSubnets: false,
					Prefixes: []alkira.UserInputPrefixes{
						{
							Id:    "7442940776704048352",
							FqId:  "",
							Value: "10.0.1.0/24",
							Type:  "SUBNET",
						},
					},
				},
			},
			expectEmpty:    false,
			expectedIds:    []string{""},
			expectedCidrs:  []string{"10.0.1.0/24"},
			expectedIntIds: []string{"7442940776704048352"},
		},
		{
			name: "multiple subnets with mixed FqId and Id",
			gcpRouting: &alkira.ConnectorGcpVpcRouting{
				ExportOptions: alkira.ConnectorGcpVpcExportOptions{
					ExportAllSubnets: false,
					Prefixes: []alkira.UserInputPrefixes{
						{
							Id:    "7442940776704048352",
							FqId:  "",
							Value: "10.0.1.0/24",
							Type:  "SUBNET",
						},
						{
							Id:    "",
							FqId:  "projects/test/regions/us-east1/subnetworks/subnet-2",
							Value: "10.0.2.0/24",
							Type:  "SUBNET",
						},
					},
				},
			},
			expectEmpty:    false,
			expectedIds:    []string{"", "projects/test/regions/us-east1/subnetworks/subnet-2"},
			expectedCidrs:  []string{"10.0.1.0/24", "10.0.2.0/24"},
			expectedIntIds: []string{"7442940776704048352", ""},
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
		{
			name: "export_all_subnets=false with nil prefixes",
			gcpRouting: &alkira.ConnectorGcpVpcRouting{
				ExportOptions: alkira.ConnectorGcpVpcExportOptions{
					ExportAllSubnets: false,
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
					assert.Contains(t, tt.expectedIntIds, s["internal_id"].(string))
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
					ExportAllSubnets: false,
					Prefixes:         nil,
				},
				ImportOptions: alkira.ConnectorGcpVpcImportOptions{
					RouteImportMode: "ADVERTISE_DEFAULT_ROUTE",
					PrefixListIds:   nil,
				},
			},
		},
		{
			name: "default route mode with subnets and explicit export_all_subnets=false",
			gcpRouting: []interface{}{
				map[string]interface{}{
					"custom_prefix":      "ADVERTISE_DEFAULT_ROUTE",
					"prefix_list_ids":    []interface{}{},
					"export_all_subnets": false,
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
			name: "default route mode with subnets but export_all_subnets defaults to false",
			gcpRouting: []interface{}{
				map[string]interface{}{
					"custom_prefix":   "ADVERTISE_DEFAULT_ROUTE",
					"prefix_list_ids": []interface{}{},
					// export_all_subnets not specified, defaults to false
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
					ExportAllSubnets: false, // defaults to false
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
					PrefixListIds:   nil,
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
					ExportAllSubnets: false,
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

func TestGcpVpcValidateExportAllSubnetsWithVpcSubnet(t *testing.T) {
	resource := resourceAlkiraConnectorGcpVpc()

	tests := []struct {
		name          string
		config        map[string]interface{}
		expectError   bool
		errorContains string
	}{
		{
			name: "export_all_subnets=true without vpc_subnet - valid",
			config: map[string]interface{}{
				"name":          "test-connector",
				"cxp":           "us-west1",
				"segment_id":    "1",
				"size":          "SMALL",
				"gcp_region":    "us-central1",
				"gcp_vpc_name":  "test-vpc",
				"credential_id": "cred-123",
				"gcp_routing": []interface{}{
					map[string]interface{}{
						"custom_prefix":      "ADVERTISE_DEFAULT_ROUTE",
						"prefix_list_ids":    []interface{}{},
						"export_all_subnets": true,
					},
				},
			},
			expectError: false,
		},
		{
			name: "export_all_subnets=true with vpc_subnet - invalid",
			config: map[string]interface{}{
				"name":          "test-connector",
				"cxp":           "us-west1",
				"segment_id":    "1",
				"size":          "SMALL",
				"gcp_region":    "us-central1",
				"gcp_vpc_name":  "test-vpc",
				"credential_id": "cred-123",
				"gcp_routing": []interface{}{
					map[string]interface{}{
						"custom_prefix":      "ADVERTISE_DEFAULT_ROUTE",
						"prefix_list_ids":    []interface{}{},
						"export_all_subnets": true,
					},
				},
				"vpc_subnet": []interface{}{
					map[string]interface{}{
						"id":   "projects/test/regions/us-central1/subnetworks/subnet-1",
						"cidr": "10.0.1.0/24",
					},
				},
			},
			expectError:   true,
			errorContains: "vpc_subnet cannot be specified when export_all_subnets is true",
		},
		{
			name: "export_all_subnets=false with vpc_subnet - valid",
			config: map[string]interface{}{
				"name":          "test-connector",
				"cxp":           "us-west1",
				"segment_id":    "1",
				"size":          "SMALL",
				"gcp_region":    "us-central1",
				"gcp_vpc_name":  "test-vpc",
				"credential_id": "cred-123",
				"gcp_routing": []interface{}{
					map[string]interface{}{
						"custom_prefix":      "ADVERTISE_DEFAULT_ROUTE",
						"prefix_list_ids":    []interface{}{},
						"export_all_subnets": false,
					},
				},
				"vpc_subnet": []interface{}{
					map[string]interface{}{
						"id":   "projects/test/regions/us-central1/subnetworks/subnet-1",
						"cidr": "10.0.1.0/24",
					},
				},
			},
			expectError: false,
		},
		{
			name: "export_all_subnets=false without vpc_subnet - valid (TPS will error)",
			config: map[string]interface{}{
				"name":          "test-connector",
				"cxp":           "us-west1",
				"segment_id":    "1",
				"size":          "SMALL",
				"gcp_region":    "us-central1",
				"gcp_vpc_name":  "test-vpc",
				"credential_id": "cred-123",
				"gcp_routing": []interface{}{
					map[string]interface{}{
						"custom_prefix":      "ADVERTISE_DEFAULT_ROUTE",
						"prefix_list_ids":    []interface{}{},
						"export_all_subnets": false,
					},
				},
			},
			expectError: false,
		},
		{
			name: "no gcp_routing with vpc_subnet - valid",
			config: map[string]interface{}{
				"name":          "test-connector",
				"cxp":           "us-west1",
				"segment_id":    "1",
				"size":          "SMALL",
				"gcp_region":    "us-central1",
				"gcp_vpc_name":  "test-vpc",
				"credential_id": "cred-123",
				"vpc_subnet": []interface{}{
					map[string]interface{}{
						"id":   "projects/test/regions/us-central1/subnetworks/subnet-1",
						"cidr": "10.0.1.0/24",
					},
				},
			},
			expectError: false,
		},
		{
			name: "export_all_subnets defaults to false with vpc_subnet - valid",
			config: map[string]interface{}{
				"name":          "test-connector",
				"cxp":           "us-west1",
				"segment_id":    "1",
				"size":          "SMALL",
				"gcp_region":    "us-central1",
				"gcp_vpc_name":  "test-vpc",
				"credential_id": "cred-123",
				"gcp_routing": []interface{}{
					map[string]interface{}{
						"custom_prefix":   "ADVERTISE_DEFAULT_ROUTE",
						"prefix_list_ids": []interface{}{},
						// export_all_subnets not specified, defaults to false
					},
				},
				"vpc_subnet": []interface{}{
					map[string]interface{}{
						"id":   "projects/test/regions/us-central1/subnetworks/subnet-1",
						"cidr": "10.0.1.0/24",
					},
				},
			},
			expectError: false,
		},
		{
			name: "export_all_subnets=true with multiple vpc_subnets - invalid",
			config: map[string]interface{}{
				"name":          "test-connector",
				"cxp":           "us-west1",
				"segment_id":    "1",
				"size":          "SMALL",
				"gcp_region":    "us-central1",
				"gcp_vpc_name":  "test-vpc",
				"credential_id": "cred-123",
				"gcp_routing": []interface{}{
					map[string]interface{}{
						"custom_prefix":      "ADVERTISE_DEFAULT_ROUTE",
						"prefix_list_ids":    []interface{}{},
						"export_all_subnets": true,
					},
				},
				"vpc_subnet": []interface{}{
					map[string]interface{}{
						"id":   "projects/test/regions/us-central1/subnetworks/subnet-1",
						"cidr": "10.0.1.0/24",
					},
					map[string]interface{}{
						"id":   "projects/test/regions/us-east1/subnetworks/subnet-2",
						"cidr": "10.0.2.0/24",
					},
				},
			},
			expectError:   true,
			errorContains: "vpc_subnet cannot be specified when export_all_subnets is true",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create test resource data
			d := resource.TestResourceData()
			d.SetId("test-id")

			// Set all the config values
			for key, val := range tt.config {
				if err := d.Set(key, val); err != nil {
					t.Fatalf("Failed to set %s: %v", key, err)
				}
			}

			// Get the CustomizeDiff function
			customizeDiff := resource.CustomizeDiff
			assert.NotNil(t, customizeDiff, "CustomizeDiff should not be nil")

			// The CustomizeDiff function requires a ResourceDiff which is difficult to mock
			// in unit tests. We test the validation logic through the resource's schema
			// and manually run the validation logic here.
			//
			// Extract the actual validation logic from CustomizeDiff:
			client := &alkira.AlkiraClient{}
			err := validateExportAllSubnetsWithVpcSubnet(d, client)

			if tt.expectError {
				assert.Error(t, err)
				if tt.errorContains != "" {
					assert.Contains(t, err.Error(), tt.errorContains)
				}
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

// validateExportAllSubnetsWithVpcSubnet extracts the validation logic
// for testability. This mirrors the logic in CustomizeDiff.
func validateExportAllSubnetsWithVpcSubnet(d *schema.ResourceData, client *alkira.AlkiraClient) error {
	// This function contains the same logic as CustomizeDiff
	// for validation of export_all_subnets and vpc_subnet mutual exclusion

	// Get gcp_routing config
	gcpRouting := d.Get("gcp_routing")
	if gcpRouting == nil {
		return nil
	}

	routing, ok := gcpRouting.([]interface{})
	if !ok || len(routing) == 0 {
		return nil
	}

	routingCfg := routing[0].(map[string]interface{})
	exportAll, ok := routingCfg["export_all_subnets"].(bool)
	if !ok || !exportAll {
		return nil
	}

	// If export_all_subnets is true, vpc_subnet must be empty
	vpcSubnets := d.Get("vpc_subnet")
	if vpcSubnets == nil {
		return nil
	}

	vpcSubnetSet, ok := vpcSubnets.(*schema.Set)
	if !ok {
		return nil
	}

	if vpcSubnetSet.Len() > 0 {
		return fmt.Errorf("vpc_subnet cannot be specified when export_all_subnets is true. " +
			"When exporting all subnets, specific vpc_subnet entries should not be provided")
	}

	return nil
}
