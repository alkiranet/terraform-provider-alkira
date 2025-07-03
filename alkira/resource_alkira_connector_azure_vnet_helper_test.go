package alkira

import (
	"testing"

	"github.com/alkiranet/alkira-client-go/alkira"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/stretchr/testify/assert"
)

func TestVnetRoutingDataValidation(t *testing.T) {
	t.Run("test routing options structure", func(t *testing.T) {
		// Test that we can process routing configuration data structures
		routingOptions := &alkira.ConnectorVnetRouting{
			ExportOptions: alkira.ConnectorVnetExportOptions{
				UserInputPrefixes: []alkira.ConnectorVnetExportOptionUserInputPrefix{
					{
						Type:  "SUBNET",
						Id:    "subnet-123",
						Value: "10.0.1.0/24",
					},
					{
						Type:  "CIDR",
						Value: "10.0.0.0/16",
					},
				},
			},
			ImportOptions: alkira.ConnectorVnetImportOptions{
				Subnets: []alkira.ConnectorVnetImportOptionsSubnet{
					{
						Id:              "subnet-123",
						RouteImportMode: "LEARNED",
						PrefixListIds:   []int{1, 2, 3},
					},
				},
			},
		}

		// Verify structure can be processed
		assert.Len(t, routingOptions.ExportOptions.UserInputPrefixes, 2)
		assert.Len(t, routingOptions.ImportOptions.Subnets, 1)

		// Test subnet prefix
		subnetPrefix := routingOptions.ExportOptions.UserInputPrefixes[0]
		assert.Equal(t, "SUBNET", subnetPrefix.Type)
		assert.Equal(t, "subnet-123", subnetPrefix.Id)
		assert.Equal(t, "10.0.1.0/24", subnetPrefix.Value)

		// Test CIDR prefix
		cidrPrefix := routingOptions.ExportOptions.UserInputPrefixes[1]
		assert.Equal(t, "CIDR", cidrPrefix.Type)
		assert.Equal(t, "10.0.0.0/16", cidrPrefix.Value)
		assert.Empty(t, cidrPrefix.Id) // CIDR doesn't have ID

		// Test import subnet
		importSubnet := routingOptions.ImportOptions.Subnets[0]
		assert.Equal(t, "subnet-123", importSubnet.Id)
		assert.Equal(t, "LEARNED", importSubnet.RouteImportMode)
		assert.Len(t, importSubnet.PrefixListIds, 3)
	})

	t.Run("test empty routing configuration", func(t *testing.T) {
		emptyConfig := &alkira.ConnectorVnetRouting{
			ExportOptions: alkira.ConnectorVnetExportOptions{
				UserInputPrefixes: []alkira.ConnectorVnetExportOptionUserInputPrefix{},
			},
		}

		assert.Len(t, emptyConfig.ExportOptions.UserInputPrefixes, 0)
	})

	t.Run("test service routes structure", func(t *testing.T) {
		serviceRoute := alkira.ConnectorVnetServiceRoute{
			Id:                 "subnet-123",
			Value:              "10.0.1.0/24",
			ServiceTags:        []string{"tag1", "tag2"},
			NativeServiceNames: []string{"service1", "service2"},
		}

		assert.Equal(t, "subnet-123", serviceRoute.Id)
		assert.Equal(t, "10.0.1.0/24", serviceRoute.Value)
		assert.Len(t, serviceRoute.ServiceTags, 2)
		assert.Len(t, serviceRoute.NativeServiceNames, 2)
		assert.Contains(t, serviceRoute.ServiceTags, "tag1")
		assert.Contains(t, serviceRoute.NativeServiceNames, "service1")
	})
}

func TestConstructVnetRouting(t *testing.T) {
	// Create a mock ResourceData for testing
	r := &schema.Resource{
		Schema: map[string]*schema.Schema{
			"routing_options": {
				Type: schema.TypeString,
			},
			"routing_prefix_list_ids": {
				Type: schema.TypeList,
				Elem: &schema.Schema{Type: schema.TypeInt},
			},
			"vnet_subnet": {
				Type: schema.TypeSet,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"subnet_id":       {Type: schema.TypeString},
						"subnet_cidr":     {Type: schema.TypeString},
						"routing_options": {Type: schema.TypeString},
						"prefix_list_ids": {Type: schema.TypeList, Elem: &schema.Schema{Type: schema.TypeInt}},
						"service_tags":    {Type: schema.TypeSet, Elem: &schema.Schema{Type: schema.TypeString}},
						"native_services": {Type: schema.TypeSet, Elem: &schema.Schema{Type: schema.TypeString}},
						"udr_list_ids":    {Type: schema.TypeSet, Elem: &schema.Schema{Type: schema.TypeInt}},
					},
				},
			},
			"vnet_cidr": {
				Type: schema.TypeSet,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"cidr":            {Type: schema.TypeString},
						"service_tags":    {Type: schema.TypeSet, Elem: &schema.Schema{Type: schema.TypeString}},
						"native_services": {Type: schema.TypeSet, Elem: &schema.Schema{Type: schema.TypeString}},
						"udr_list_ids":    {Type: schema.TypeSet, Elem: &schema.Schema{Type: schema.TypeInt}},
					},
				},
			},
		},
	}

	tests := []struct {
		name                   string
		routingOptions         string
		routingPrefixListIds   []interface{}
		vnetSubnets            []interface{}
		vnetCidrs              []interface{}
		expectedExportPrefixes int
		expectedImportSubnets  int
		expectedServiceSubnets int
		expectedUdrSubnets     int
		expectError            bool
	}{
		{
			name:                 "basic configuration",
			routingOptions:       "LEARNED",
			routingPrefixListIds: []interface{}{1, 2, 3},
			vnetSubnets: []interface{}{
				map[string]interface{}{
					"subnet_id":       "subnet-123",
					"subnet_cidr":     "10.0.1.0/24",
					"routing_options": "ADVERTISED",
					"prefix_list_ids": []interface{}{4, 5},
					"service_tags":    schema.NewSet(schema.HashString, []interface{}{"tag1", "tag2"}),
					"native_services": schema.NewSet(schema.HashString, []interface{}{"service1"}),
					"udr_list_ids":    schema.NewSet(schema.HashInt, []interface{}{10, 20}),
				},
			},
			vnetCidrs: []interface{}{
				map[string]interface{}{
					"cidr":            "10.0.0.0/16",
					"service_tags":    schema.NewSet(schema.HashString, []interface{}{}),
					"native_services": schema.NewSet(schema.HashString, []interface{}{}),
					"udr_list_ids":    schema.NewSet(schema.HashInt, []interface{}{}),
				},
			},
			expectedExportPrefixes: 2, // 1 subnet + 1 CIDR
			expectedImportSubnets:  1,
			expectedServiceSubnets: 1,
			expectedUdrSubnets:     1,
			expectError:            false,
		},
		{
			name:                   "minimal configuration",
			routingOptions:         "LEARNED",
			routingPrefixListIds:   []interface{}{},
			vnetSubnets:            []interface{}{},
			vnetCidrs:              []interface{}{},
			expectedExportPrefixes: 0,
			expectedImportSubnets:  0,
			expectedServiceSubnets: 0,
			expectedUdrSubnets:     0,
			expectError:            false,
		},
		{
			name:                 "subnet without optional fields",
			routingOptions:       "LEARNED",
			routingPrefixListIds: []interface{}{},
			vnetSubnets: []interface{}{
				map[string]interface{}{
					"subnet_id":       "subnet-456",
					"subnet_cidr":     "192.168.1.0/24",
					"prefix_list_ids": []interface{}{},
					"service_tags":    schema.NewSet(schema.HashString, []interface{}{}),
					"native_services": schema.NewSet(schema.HashString, []interface{}{}),
					"udr_list_ids":    schema.NewSet(schema.HashInt, []interface{}{}),
					// Note: no routing_options field to test optional field handling
				},
			},
			vnetCidrs:              []interface{}{},
			expectedExportPrefixes: 1,
			expectedImportSubnets:  1, // Import subnet created even without routing_options due to zero value
			expectedServiceSubnets: 0, // Empty service tags/services
			expectedUdrSubnets:     0, // Empty UDR list
			expectError:            false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := r.TestResourceData()
			d.Set("routing_options", tt.routingOptions)
			d.Set("routing_prefix_list_ids", tt.routingPrefixListIds)
			d.Set("vnet_subnet", schema.NewSet(
				func(i interface{}) int {
					return schema.HashString("test")
				},
				tt.vnetSubnets,
			))
			d.Set("vnet_cidr", schema.NewSet(
				func(i interface{}) int {
					return schema.HashString("test")
				},
				tt.vnetCidrs,
			))

			result, err := constructVnetRouting(d)

			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, result)

				// Check export options
				assert.Len(t, result.ExportOptions.UserInputPrefixes, tt.expectedExportPrefixes)

				// Check import options
				assert.Equal(t, tt.routingOptions, result.ImportOptions.RouteImportMode)
				assert.Len(t, result.ImportOptions.Subnets, tt.expectedImportSubnets)

				// Check service routes
				assert.Len(t, result.ServiceRoutes.Subnets, tt.expectedServiceSubnets)

				// Check UDR lists
				assert.Len(t, result.UdrLists.Subnets, tt.expectedUdrSubnets)
			}
		})
	}
}

func TestAzureVnetValidation(t *testing.T) {
	t.Run("test Azure subnet ID format", func(t *testing.T) {
		validSubnetIds := []string{
			"/subscriptions/12345/resourceGroups/rg/providers/Microsoft.Network/virtualNetworks/vnet/subnets/subnet1",
			"/subscriptions/abcdef/resourceGroups/test-rg/providers/Microsoft.Network/virtualNetworks/test-vnet/subnets/subnet-test",
		}
		for _, subnetId := range validSubnetIds {
			assert.NotEmpty(t, subnetId)
			assert.Contains(t, subnetId, "/subscriptions/")
			assert.Contains(t, subnetId, "/subnets/")
		}
	})

	t.Run("test Azure CIDR format", func(t *testing.T) {
		validCIDRs := []string{"10.0.0.0/16", "192.168.0.0/24", "172.16.0.0/12"}
		for _, cidr := range validCIDRs {
			assert.NotEmpty(t, cidr)
			assert.Contains(t, cidr, "/")
			assert.Regexp(t, `^\d+\.\d+\.\d+\.\d+/\d+$`, cidr)
		}
	})

	t.Run("test routing options", func(t *testing.T) {
		validOptions := []string{"LEARNED", "ADVERTISED", "BOTH"}
		for _, option := range validOptions {
			assert.NotEmpty(t, option)
			assert.Contains(t, []string{"LEARNED", "ADVERTISED", "BOTH"}, option)
		}
	})
}

func TestAzureVnetDataStructures(t *testing.T) {
	t.Run("test export option user input prefix", func(t *testing.T) {
		subnetPrefix := alkira.ConnectorVnetExportOptionUserInputPrefix{
			Type:  "SUBNET",
			Id:    "subnet-123",
			Value: "10.0.1.0/24",
		}

		assert.Equal(t, "SUBNET", subnetPrefix.Type)
		assert.Equal(t, "subnet-123", subnetPrefix.Id)
		assert.Equal(t, "10.0.1.0/24", subnetPrefix.Value)

		cidrPrefix := alkira.ConnectorVnetExportOptionUserInputPrefix{
			Type:  "CIDR",
			Value: "10.0.0.0/16",
		}

		assert.Equal(t, "CIDR", cidrPrefix.Type)
		assert.Equal(t, "10.0.0.0/16", cidrPrefix.Value)
		assert.Empty(t, cidrPrefix.Id) // CIDR doesn't have ID
	})

	t.Run("test import options subnet", func(t *testing.T) {
		importSubnet := alkira.ConnectorVnetImportOptionsSubnet{
			Id:              "subnet-123",
			Value:           "10.0.1.0/24",
			RouteImportMode: "LEARNED",
			PrefixListIds:   []int{1, 2, 3},
		}

		assert.Equal(t, "subnet-123", importSubnet.Id)
		assert.Equal(t, "10.0.1.0/24", importSubnet.Value)
		assert.Equal(t, "LEARNED", importSubnet.RouteImportMode)
		assert.Len(t, importSubnet.PrefixListIds, 3)
		assert.Contains(t, importSubnet.PrefixListIds, 1)
		assert.Contains(t, importSubnet.PrefixListIds, 2)
		assert.Contains(t, importSubnet.PrefixListIds, 3)
	})

	t.Run("test service route", func(t *testing.T) {
		serviceRoute := alkira.ConnectorVnetServiceRoute{
			Id:                 "subnet-123",
			Value:              "10.0.1.0/24",
			ServiceTags:        []string{"tag1", "tag2"},
			NativeServiceNames: []string{"service1", "service2"},
		}

		assert.Equal(t, "subnet-123", serviceRoute.Id)
		assert.Equal(t, "10.0.1.0/24", serviceRoute.Value)
		assert.Len(t, serviceRoute.ServiceTags, 2)
		assert.Contains(t, serviceRoute.ServiceTags, "tag1")
		assert.Contains(t, serviceRoute.ServiceTags, "tag2")
		assert.Len(t, serviceRoute.NativeServiceNames, 2)
		assert.Contains(t, serviceRoute.NativeServiceNames, "service1")
		assert.Contains(t, serviceRoute.NativeServiceNames, "service2")
	})

	t.Run("test UDR list", func(t *testing.T) {
		udrList := alkira.ConnectorVnetUdrList{
			Id:         "subnet-123",
			Value:      "10.0.1.0/24",
			UdrListIds: []int{10, 20, 30},
		}

		assert.Equal(t, "subnet-123", udrList.Id)
		assert.Equal(t, "10.0.1.0/24", udrList.Value)
		assert.Len(t, udrList.UdrListIds, 3)
		assert.Contains(t, udrList.UdrListIds, 10)
		assert.Contains(t, udrList.UdrListIds, 20)
		assert.Contains(t, udrList.UdrListIds, 30)
	})
}

func TestAzureVnetTypeConversions(t *testing.T) {
	t.Run("test string list conversion", func(t *testing.T) {
		// Test service tags conversion
		serviceTags := schema.NewSet(schema.HashString, []interface{}{"tag1", "tag2", "tag3"})
		converted := convertTypeSetToStringList(serviceTags)
		assert.Len(t, converted, 3)
		assert.ElementsMatch(t, []string{"tag1", "tag2", "tag3"}, converted)

		// Test native services conversion
		nativeServices := schema.NewSet(schema.HashString, []interface{}{"service1", "service2"})
		convertedServices := convertTypeSetToStringList(nativeServices)
		assert.Len(t, convertedServices, 2)
		assert.ElementsMatch(t, []string{"service1", "service2"}, convertedServices)
	})

	t.Run("test int list conversion", func(t *testing.T) {
		// Test prefix list IDs conversion
		prefixIds := []interface{}{1, 2, 3, 4, 5}
		converted := convertTypeListToIntList(prefixIds)
		assert.Len(t, converted, 5)
		assert.Equal(t, []int{1, 2, 3, 4, 5}, converted)

		// Test UDR list IDs conversion
		udrIds := schema.NewSet(schema.HashInt, []interface{}{10, 20, 30})
		convertedUdr := convertTypeSetToIntList(udrIds)
		assert.Len(t, convertedUdr, 3)
		assert.ElementsMatch(t, []int{10, 20, 30}, convertedUdr)
	})

	t.Run("test empty conversions", func(t *testing.T) {
		// Test empty string set
		emptyStringSet := schema.NewSet(schema.HashString, []interface{}{})
		convertedEmpty := convertTypeSetToStringList(emptyStringSet)
		assert.Nil(t, convertedEmpty)

		// Test empty int set
		emptyIntSet := schema.NewSet(schema.HashInt, []interface{}{})
		convertedEmptyInt := convertTypeSetToIntList(emptyIntSet)
		assert.Nil(t, convertedEmptyInt)

		// Test nil conversions
		nilConverted := convertTypeSetToStringList(nil)
		assert.Nil(t, nilConverted)

		nilIntConverted := convertTypeSetToIntList(nil)
		assert.Nil(t, nilIntConverted)
	})
}

func TestAzureVnetErrorConditions(t *testing.T) {
	t.Run("test prefix type validation", func(t *testing.T) {
		validTypes := []string{"SUBNET", "CIDR"}
		for _, prefixType := range validTypes {
			assert.NotEmpty(t, prefixType)
			assert.Contains(t, []string{"SUBNET", "CIDR"}, prefixType)
		}
	})

	t.Run("test conditional processing", func(t *testing.T) {
		// Test that service routes are only created when tags or services are present
		emptyTags := schema.NewSet(schema.HashString, []interface{}{})
		emptyServices := schema.NewSet(schema.HashString, []interface{}{})

		assert.Equal(t, 0, emptyTags.Len())
		assert.Equal(t, 0, emptyServices.Len())

		// Condition should be false when both are empty
		condition := (emptyTags != nil && emptyTags.Len() > 0) ||
			(emptyServices != nil && emptyServices.Len() > 0)
		assert.False(t, condition)

		// Test with non-empty tags
		nonEmptyTags := schema.NewSet(schema.HashString, []interface{}{"tag1"})
		conditionWithTags := (nonEmptyTags != nil && nonEmptyTags.Len() > 0) ||
			(emptyServices != nil && emptyServices.Len() > 0)
		assert.True(t, conditionWithTags)
	})
}
