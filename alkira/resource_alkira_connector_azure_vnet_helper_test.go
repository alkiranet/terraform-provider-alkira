package alkira

import (
	"context"
	"testing"

	"github.com/alkiranet/alkira-client-go/alkira"
	"github.com/hashicorp/go-cty/cty"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
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

func TestAzureVnetFieldNameMatch(t *testing.T) {
	t.Run("schema uses vnet_subnet field name", func(t *testing.T) {
		resourceSchema := resourceAlkiraConnectorAzureVnet().Schema

		// Verify "vnet_subnet" field exists in schema
		vnetSubnetField, exists := resourceSchema["vnet_subnet"]
		assert.True(t, exists, "Schema must have 'vnet_subnet' field")
		assert.NotNil(t, vnetSubnetField, "vnet_subnet field must not be nil")

		// Verify "vnet_subnets" (plural - the bug) does NOT exist in schema
		_, wrongFieldExists := resourceSchema["vnet_subnets"]
		assert.False(t, wrongFieldExists, "Schema should NOT have 'vnet_subnets' field (bug was using plural)")
	})
}

func TestAzureVnetCidrFieldNameMatch(t *testing.T) {
	t.Run("schema uses vnet_cidr field name", func(t *testing.T) {
		resourceSchema := resourceAlkiraConnectorAzureVnet().Schema

		// Verify "vnet_cidr" exists (this one was already correct)
		vnetCidrField, exists := resourceSchema["vnet_cidr"]
		assert.True(t, exists, "Schema must have 'vnet_cidr' field")
		assert.NotNil(t, vnetCidrField, "vnet_cidr field must not be nil")
	})
}

// The backend auto-populates `customerAsn` when the user omits it in VGW mode
// (either to the existing Azure VGW's ASN, or to a DEFAULT_ASN constant for a
// fresh VNet). If the provider schema declares this field as Optional only,
// the backend-supplied value lands in TF state and reads as drift versus the
// empty user config, producing a perpetual "remove customer_asn" plan diff
// that apply cannot resolve.
//
// Marking the field Optional + Computed tells Terraform that a backend-supplied
// value is a valid settled outcome — no drift. This test pins that contract so
// the regression cannot be reintroduced silently.
func TestCustomerAsnSchemaIsOptionalAndComputed(t *testing.T) {
	resourceSchema := resourceAlkiraConnectorAzureVnet().Schema

	customerAsn, exists := resourceSchema["customer_asn"]
	assert.True(t, exists, "Schema must have 'customer_asn' field")
	assert.NotNil(t, customerAsn, "customer_asn field must not be nil")

	assert.Equal(t, schema.TypeInt, customerAsn.Type, "customer_asn must be TypeInt")
	assert.True(t, customerAsn.Optional, "customer_asn must be Optional so users can set it explicitly")
	assert.True(t, customerAsn.Computed,
		"customer_asn must be Computed because the backend auto-populates it "+
			"in VGW mode when the user omits it (AK-68129). "+
			"Without Computed, Terraform treats the backend-supplied value as drift.")
}

func TestPeeringGatewayCxpIdSchemaIsOptionalAndComputed(t *testing.T) {
	r := resourceAlkiraConnectorAzureVnet()

	peeringGateway, ok := r.Schema["peering_gateway_cxp_id"]
	require.True(t, ok, "peering_gateway_cxp_id must exist")

	// The backend may assign a CXP peering gateway when the config omits it and
	// returns that value on Read, and rejects any change to it once the connector
	// is provisioned. Without Computed, an omitted value plans as a removal and
	// the update sends no gateway, which the backend rejects. See AK-75338.
	assert.Equal(t, schema.TypeInt, peeringGateway.Type, "peering_gateway_cxp_id must be TypeInt")
	assert.True(t, peeringGateway.Optional, "peering_gateway_cxp_id must stay Optional")
	assert.True(t, peeringGateway.Computed,
		"peering_gateway_cxp_id must be Computed, otherwise an omitted value plans as a removal")

	// Guard the behaviour Computed provides: with the backend value in state and
	// the attribute omitted from config, the plan keeps the state value.
	state := &terraform.InstanceState{
		ID: "1",
		Attributes: map[string]string{
			"id":                     "1",
			"peering_gateway_cxp_id": "7",
		},
	}
	client := &alkira.AlkiraClient{}

	diff, err := r.Diff(context.Background(), state, terraform.NewResourceConfigRaw(map[string]interface{}{}), client)
	require.NoError(t, err)
	if diff != nil {
		_, changed := diff.Attributes["peering_gateway_cxp_id"]
		assert.False(t, changed, "an omitted peering_gateway_cxp_id must keep the state value")
	}

	// An explicit value still plans as a change.
	diff, err = r.Diff(context.Background(), state,
		terraform.NewResourceConfigRaw(map[string]interface{}{"peering_gateway_cxp_id": 8}), client)
	require.NoError(t, err)
	require.NotNil(t, diff)
	attr, changed := diff.Attributes["peering_gateway_cxp_id"]
	require.True(t, changed, "an explicit peering_gateway_cxp_id must plan as a change")
	assert.Equal(t, "7", attr.Old)
	assert.Equal(t, "8", attr.New)
}

// AK-74515: a multi-prefix Azure subnet is ONE vnet_subnet block naming every
// prefix in subnet_cidrs. Two blocks sharing a subnet_id serialise to two entries
// with the same id, which the API rejects - the orchestrator keys its subnet maps
// by id and a duplicate aborts task generation for the whole tenant network.
func TestConstructVnetRoutingMultiPrefixSubnet(t *testing.T) {
	r := resourceAlkiraConnectorAzureVnet()
	d := r.TestResourceData()

	block := map[string]interface{}{
		"subnet_id":       "/subscriptions/s/.../subnets/private-endpoints",
		"subnet_cidrs":    schema.NewSet(schema.HashString, []interface{}{"10.169.142.0/28", "10.169.142.32/27"}),
		"routing_options": "ADVERTISE_DEFAULT_ROUTE",
		"service_tags":    schema.NewSet(schema.HashString, []interface{}{"AzureKeyVault"}),
		"udr_list_ids":    schema.NewSet(schema.HashInt, []interface{}{88}),
	}
	d.Set("vnet_subnet", schema.NewSet(
		func(i interface{}) int { return schema.HashString("test") },
		[]interface{}{block}))

	result, err := constructVnetRouting(d)
	assert.NoError(t, err)
	assert.NotNil(t, result)

	// One entry for the subnet, naming both prefixes - not one entry per prefix.
	assert.Len(t, result.ExportOptions.UserInputPrefixes, 1)
	sel := result.ExportOptions.UserInputPrefixes[0]
	assert.ElementsMatch(t, []string{"10.169.142.0/28", "10.169.142.32/27"}, sel.Values)
	assert.Empty(t, sel.Value, "value must be empty so it is omitted from the payload")

	assert.Len(t, result.ImportOptions.Subnets, 1)
	assert.ElementsMatch(t, []string{"10.169.142.0/28", "10.169.142.32/27"}, result.ImportOptions.Subnets[0].Values)
	assert.Empty(t, result.ImportOptions.Subnets[0].Value)

	assert.Len(t, result.ServiceRoutes.Subnets, 1)
	assert.ElementsMatch(t, []string{"10.169.142.0/28", "10.169.142.32/27"}, result.ServiceRoutes.Subnets[0].Values)

	assert.Len(t, result.UdrLists.Subnets, 1)
	assert.ElementsMatch(t, []string{"10.169.142.0/28", "10.169.142.32/27"}, result.UdrLists.Subnets[0].Values)
}

// The single-prefix block keeps using subnet_cidr and must be byte-identical to
// what the provider sent before this change.
func TestConstructVnetRoutingSinglePrefixUnchanged(t *testing.T) {
	r := resourceAlkiraConnectorAzureVnet()
	d := r.TestResourceData()

	d.Set("vnet_subnet", schema.NewSet(
		func(i interface{}) int { return schema.HashString("test") },
		[]interface{}{map[string]interface{}{
			"subnet_id":       "/subscriptions/s/.../subnets/default",
			"subnet_cidr":     "10.0.0.0/24",
			"routing_options": "ADVERTISE_DEFAULT_ROUTE",
		}}))

	result, err := constructVnetRouting(d)
	assert.NoError(t, err)
	assert.Len(t, result.ExportOptions.UserInputPrefixes, 1)
	assert.Equal(t, "10.0.0.0/24", result.ExportOptions.UserInputPrefixes[0].Value)
	assert.Empty(t, result.ExportOptions.UserInputPrefixes[0].Values)
}

func TestConstructVnetRoutingRejectsBothCidrForms(t *testing.T) {
	r := resourceAlkiraConnectorAzureVnet()
	d := r.TestResourceData()

	d.Set("vnet_subnet", schema.NewSet(
		func(i interface{}) int { return schema.HashString("test") },
		[]interface{}{map[string]interface{}{
			"subnet_id":    "/subscriptions/s/.../subnets/default",
			"subnet_cidr":  "10.0.0.0/24",
			"subnet_cidrs": schema.NewSet(schema.HashString, []interface{}{"10.0.0.0/24", "20.1.0.0/24"}),
		}}))

	_, err := constructVnetRouting(d)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "mutually exclusive")
}

func TestConstructVnetRoutingRejectsNeitherCidrForm(t *testing.T) {
	r := resourceAlkiraConnectorAzureVnet()
	d := r.TestResourceData()

	d.Set("vnet_subnet", schema.NewSet(
		func(i interface{}) int { return schema.HashString("test") },
		[]interface{}{map[string]interface{}{
			"subnet_id": "/subscriptions/s/.../subnets/default",
		}}))

	_, err := constructVnetRouting(d)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "must be set")
}

// AK-74515: the same rule must fail at PLAN time, not only at apply - subnet_cidr moved
// from Required to Optional, so the schema no longer catches a block that sets neither.
//
// Driven with cty values because the validator walks the raw config: on a ResourceDiff
// d.Get flattens an unknown to the zero value, which is the bug the unknown cases below
// pin down.
func TestValidateVnetSubnetCidrForms(t *testing.T) {
	blockType := cty.Object(map[string]cty.Type{
		"subnet_id":    cty.String,
		"subnet_cidr":  cty.String,
		"subnet_cidrs": cty.Set(cty.String),
	})
	cfg := func(blocks ...cty.Value) cty.Value {
		return cty.ObjectVal(map[string]cty.Value{"vnet_subnet": cty.SetVal(blocks)})
	}
	blk := func(id string, cidr cty.Value, cidrs cty.Value) cty.Value {
		return cty.ObjectVal(map[string]cty.Value{
			"subnet_id":    cty.StringVal(id),
			"subnet_cidr":  cidr,
			"subnet_cidrs": cidrs,
		})
	}
	noCidr := cty.StringVal("")
	noCidrs := cty.SetValEmpty(cty.String)
	set := func(vals ...string) cty.Value {
		out := make([]cty.Value, 0, len(vals))
		for _, v := range vals {
			out = append(out, cty.StringVal(v))
		}
		return cty.SetVal(out)
	}

	t.Run("neither form set is rejected", func(t *testing.T) {
		assert.ErrorContains(t, validateVnetSubnetCidrForms(cfg(blk("/sub/default", noCidr, noCidrs))), "must be set")
	})

	t.Run("both forms set is rejected", func(t *testing.T) {
		assert.ErrorContains(t, validateVnetSubnetCidrForms(cfg(blk("/sub/default", cty.StringVal("10.0.0.0/24"), set("10.0.0.0/24")))), "mutually exclusive")
	})

	t.Run("each valid form on its own passes", func(t *testing.T) {
		assert.NoError(t, validateVnetSubnetCidrForms(cfg(blk("/sub/a", cty.StringVal("10.0.0.0/24"), noCidrs))))
		assert.NoError(t, validateVnetSubnetCidrForms(cfg(blk("/sub/b", noCidr, set("10.0.1.0/24", "20.1.1.0/24")))))
		assert.NoError(t, validateVnetSubnetCidrForms(cfg(blk("/sub/c", noCidr, set("10.0.1.0/24")))))
	})

	t.Run("nothing configured at all is not this validators problem", func(t *testing.T) {
		assert.NoError(t, validateVnetSubnetCidrForms(cty.NullVal(cty.EmptyObject)))
		assert.NoError(t, validateVnetSubnetCidrForms(cty.UnknownVal(cty.EmptyObject)))
	})

	// An unknown at plan time must not read as "neither form set":
	// subnet_cidrs = azurerm_subnet.pe.address_prefixes is unknown whenever the subnet
	// is created in the same run, and d.Get would flatten it to an empty set.
	t.Run("unknown subnet_cidrs is not an error", func(t *testing.T) {
		assert.NoError(t, validateVnetSubnetCidrForms(cfg(blk("/sub/a", noCidr, cty.UnknownVal(cty.Set(cty.String))))))
	})

	t.Run("unknown subnet_cidr is not an error", func(t *testing.T) {
		assert.NoError(t, validateVnetSubnetCidrForms(cfg(blk("/sub/a", cty.UnknownVal(cty.String), noCidrs))))
	})

	t.Run("wholly unknown block is not an error", func(t *testing.T) {
		assert.NoError(t, validateVnetSubnetCidrForms(cfg(cty.UnknownVal(blockType))))
	})

	t.Run("unknown whole vnet_subnet set is not an error", func(t *testing.T) {
		assert.NoError(t, validateVnetSubnetCidrForms(cty.ObjectVal(map[string]cty.Value{
			"vnet_subnet": cty.UnknownVal(cty.Set(blockType)),
		})))
	})

	t.Run("cannot judge exclusivity while one side is unknown", func(t *testing.T) {
		assert.NoError(t, validateVnetSubnetCidrForms(cfg(blk("/sub/a", cty.StringVal("10.0.0.0/24"), cty.UnknownVal(cty.Set(cty.String))))))
	})
}

// AK-74515: the server echoes rows as sent, so a multi-prefix row carries `value`
// alongside `values` only when the client that wrote it sent both - but when it does,
// the read path must not write subnet_cidr from it.
// A block written with only subnet_cidrs plans subnet_cidr="" (Optional, not Computed);
// Create ends in Read, so a non-empty subnet_cidr in state makes core reject the apply
// with "planned set element does not correlate with any element in actual", and the
// changed TypeSet hash then shows the block as removed-and-re-added on every plan.
func TestSetVnetRoutingMultiPrefixLeavesSubnetCidrEmpty(t *testing.T) {
	r := resourceAlkiraConnectorAzureVnet()
	d := r.TestResourceData()

	routing := &alkira.ConnectorVnetRouting{
		ExportOptions: alkira.ConnectorVnetExportOptions{
			UserInputPrefixes: []alkira.ConnectorVnetExportOptionUserInputPrefix{
				{
					Type:   "SUBNET",
					Id:     "/subscriptions/s/.../subnets/private-endpoints",
					Value:  "10.169.142.0/28",
					Values: []string{"10.169.142.0/28", "10.169.142.32/27"},
				},
			},
		},
	}

	setVnetRouting(d, routing)

	subnets := d.Get("vnet_subnet").(*schema.Set).List()
	if len(subnets) != 1 {
		t.Fatalf("expected 1 vnet_subnet, got %d", len(subnets))
	}
	block := subnets[0].(map[string]interface{})

	assert.Equal(t, "", block["subnet_cidr"],
		"subnet_cidr must stay empty when the server sent values - otherwise the applied "+
			"element does not correlate with the planned one")

	got := convertTypeSetToStringList(block["subnet_cidrs"].(*schema.Set))
	assert.ElementsMatch(t, []string{"10.169.142.0/28", "10.169.142.32/27"}, got)
}

// AK-74515: a single-element subnet_cidrs must go out as `value`, not a one-element
// `values`. The server returns `value` for a subnet onboarded with one prefix, so
// sending `values` there cannot round-trip: Read would populate subnet_cidr and leave
// subnet_cidrs empty, inverting the config and failing the apply the same way.
// Nothing requires subnet_cidrs to hold two or more entries, so a user writing the
// plural form uniformly is following the documentation.
func TestConstructVnetRoutingSinglePrefixCidrsSendsValue(t *testing.T) {
	r := resourceAlkiraConnectorAzureVnet()
	d := r.TestResourceData()

	block := map[string]interface{}{
		"subnet_id":       "/subscriptions/s/.../subnets/single",
		"subnet_cidrs":    schema.NewSet(schema.HashString, []interface{}{"10.0.0.0/24"}),
		"routing_options": "ADVERTISE_DEFAULT_ROUTE",
	}
	if err := d.Set("vnet_subnet", schema.NewSet(
		func(i interface{}) int { return schema.HashString("test") },
		[]interface{}{block})); err != nil {
		t.Fatalf("set vnet_subnet: %v", err)
	}

	routing, err := constructVnetRouting(d)
	if err != nil {
		t.Fatalf("constructVnetRouting: %v", err)
	}

	prefixes := routing.ExportOptions.UserInputPrefixes
	if len(prefixes) != 1 {
		t.Fatalf("expected 1 userInputPrefix, got %d", len(prefixes))
	}

	assert.Equal(t, "10.0.0.0/24", prefixes[0].Value,
		"a single prefix must travel as value so the server response round-trips")
	assert.Empty(t, prefixes[0].Values,
		"a single prefix must not be sent as a one-element values list")
}

// AK-74515: the real round trip for a one-element subnet_cidrs, which
// TestConstructVnetRoutingSinglePrefixCidrsSendsValue does not cover - that one asserts
// the wire shape only. A single prefix travels as `value` whichever form the config
// used, so the server response cannot tell Read which form to write back. Read must key
// off the form already in state, or the block returns as subnet_cidr and its TypeSet
// element hash stops matching the planned one - a permanent diff.
func TestSetVnetRoutingPreservesSinglePrefixPluralForm(t *testing.T) {
	r := resourceAlkiraConnectorAzureVnet()
	d := r.TestResourceData()

	// The config/state the user wrote: the plural form with one prefix.
	block := map[string]interface{}{
		"subnet_id":    "/subscriptions/s/.../subnets/single",
		"subnet_cidrs": schema.NewSet(schema.HashString, []interface{}{"10.0.0.0/24"}),
	}
	if err := d.Set("vnet_subnet", schema.NewSet(
		func(i interface{}) int { return schema.HashString("test") },
		[]interface{}{block})); err != nil {
		t.Fatalf("set vnet_subnet: %v", err)
	}

	// What the server returns for a subnet onboarded with one prefix: `value` only.
	routing := &alkira.ConnectorVnetRouting{
		ExportOptions: alkira.ConnectorVnetExportOptions{
			UserInputPrefixes: []alkira.ConnectorVnetExportOptionUserInputPrefix{
				{
					Type:  "SUBNET",
					Id:    "/subscriptions/s/.../subnets/single",
					Value: "10.0.0.0/24",
				},
			},
		},
	}

	setVnetRouting(d, routing)

	subnets := d.Get("vnet_subnet").(*schema.Set).List()
	if len(subnets) != 1 {
		t.Fatalf("expected 1 vnet_subnet, got %d", len(subnets))
	}
	got := subnets[0].(map[string]interface{})

	assert.Equal(t, "", got["subnet_cidr"],
		"a block written with subnet_cidrs must not come back as subnet_cidr")
	assert.ElementsMatch(t, []string{"10.0.0.0/24"},
		convertTypeSetToStringList(got["subnet_cidrs"].(*schema.Set)),
		"the plural form the config used must survive the round trip")
}

// The mirror case: a block written with the singular form must stay singular, even
// though the response is byte-identical to the one above.
func TestSetVnetRoutingPreservesSingularForm(t *testing.T) {
	r := resourceAlkiraConnectorAzureVnet()
	d := r.TestResourceData()

	block := map[string]interface{}{
		"subnet_id":   "/subscriptions/s/.../subnets/single",
		"subnet_cidr": "10.0.0.0/24",
	}
	if err := d.Set("vnet_subnet", schema.NewSet(
		func(i interface{}) int { return schema.HashString("test") },
		[]interface{}{block})); err != nil {
		t.Fatalf("set vnet_subnet: %v", err)
	}

	routing := &alkira.ConnectorVnetRouting{
		ExportOptions: alkira.ConnectorVnetExportOptions{
			UserInputPrefixes: []alkira.ConnectorVnetExportOptionUserInputPrefix{
				{Type: "SUBNET", Id: "/subscriptions/s/.../subnets/single", Value: "10.0.0.0/24"},
			},
		},
	}

	setVnetRouting(d, routing)

	got := d.Get("vnet_subnet").(*schema.Set).List()[0].(map[string]interface{})
	assert.Equal(t, "10.0.0.0/24", got["subnet_cidr"])
	assert.Equal(t, 0, got["subnet_cidrs"].(*schema.Set).Len(),
		"a singular block must not gain a subnet_cidrs")
}

// AK-74515: import has no prior state, so subnetIdsUsingPluralForm is empty and the
// response shape is all Read has. A multi-prefix subnet must still land in
// subnet_cidrs - falling back to subnet_cidr would drop every prefix but one.
func TestSetVnetRoutingImportWithoutPriorState(t *testing.T) {
	r := resourceAlkiraConnectorAzureVnet()
	d := r.TestResourceData() // no vnet_subnet set: this is the import case

	routing := &alkira.ConnectorVnetRouting{
		ExportOptions: alkira.ConnectorVnetExportOptions{
			UserInputPrefixes: []alkira.ConnectorVnetExportOptionUserInputPrefix{
				{
					Type:   "SUBNET",
					Id:     "/sub/multi",
					Value:  "10.169.142.0/28",
					Values: []string{"10.169.142.0/28", "10.169.142.32/27"},
				},
				{Type: "SUBNET", Id: "/sub/single", Value: "10.0.0.0/24"},
			},
		},
	}

	setVnetRouting(d, routing)

	byId := map[string]map[string]interface{}{}
	for _, b := range d.Get("vnet_subnet").(*schema.Set).List() {
		m := b.(map[string]interface{})
		byId[m["subnet_id"].(string)] = m
	}

	multi := byId["/sub/multi"]
	assert.Equal(t, "", multi["subnet_cidr"], "multi-prefix import must not collapse to subnet_cidr")
	assert.ElementsMatch(t, []string{"10.169.142.0/28", "10.169.142.32/27"},
		convertTypeSetToStringList(multi["subnet_cidrs"].(*schema.Set)),
		"import must keep every prefix")

	single := byId["/sub/single"]
	assert.Equal(t, "10.0.0.0/24", single["subnet_cidr"])
	assert.Equal(t, 0, single["subnet_cidrs"].(*schema.Set).Len())
}

// AK-74515: a subnet that genuinely grew a prefix out-of-band must surface as a diff,
// not be masked. The config says one prefix, the server now reports two - Read must
// report both so the plan shows the drift.
func TestSetVnetRoutingSurfacesOutOfBandPrefixAddition(t *testing.T) {
	r := resourceAlkiraConnectorAzureVnet()
	d := r.TestResourceData()

	if err := d.Set("vnet_subnet", schema.NewSet(
		func(i interface{}) int { return schema.HashString("test") },
		[]interface{}{map[string]interface{}{
			"subnet_id":    "/sub/grew",
			"subnet_cidrs": schema.NewSet(schema.HashString, []interface{}{"10.0.0.0/24"}),
		}})); err != nil {
		t.Fatalf("set vnet_subnet: %v", err)
	}

	routing := &alkira.ConnectorVnetRouting{
		ExportOptions: alkira.ConnectorVnetExportOptions{
			UserInputPrefixes: []alkira.ConnectorVnetExportOptionUserInputPrefix{
				{
					Type:   "SUBNET",
					Id:     "/sub/grew",
					Value:  "10.0.0.0/24",
					Values: []string{"10.0.0.0/24", "10.0.1.0/24"},
				},
			},
		},
	}

	setVnetRouting(d, routing)

	got := d.Get("vnet_subnet").(*schema.Set).List()[0].(map[string]interface{})
	assert.ElementsMatch(t, []string{"10.0.0.0/24", "10.0.1.0/24"},
		convertTypeSetToStringList(got["subnet_cidrs"].(*schema.Set)),
		"real drift must reach state so the plan shows it")
}

// AK-74515: a subnet_id present in state but absent from the response must not leak the
// old form onto an unrelated block, and an empty response must not panic.
func TestSetVnetRoutingHandlesMissingAndEmptyResponses(t *testing.T) {
	r := resourceAlkiraConnectorAzureVnet()
	d := r.TestResourceData()

	if err := d.Set("vnet_subnet", schema.NewSet(
		func(i interface{}) int { return schema.HashString("test") },
		[]interface{}{map[string]interface{}{
			"subnet_id":    "/sub/gone",
			"subnet_cidrs": schema.NewSet(schema.HashString, []interface{}{"10.0.0.0/24"}),
		}})); err != nil {
		t.Fatalf("set vnet_subnet: %v", err)
	}

	setVnetRouting(d, &alkira.ConnectorVnetRouting{})

	assert.Equal(t, 0, d.Get("vnet_subnet").(*schema.Set).Len(),
		"a subnet no longer returned must disappear from state")
}

// AK-74515: the failure this whole round is about is a TypeSet ELEMENT HASH mismatch -
// core compares planned and actual elements by hash, so asserting field values alone is
// not enough. Both sides are round-tripped through the schema so the comparison is
// like-for-like: a sparse config literal and a state map differ in field COUNT, which
// would make this pass or fail for the wrong reason.
func TestSetVnetRoutingElementHashMatchesConfig(t *testing.T) {
	elem := resourceAlkiraConnectorAzureVnet().Schema["vnet_subnet"].Elem.(*schema.Resource)
	hashElem := schema.HashResource(elem)

	// normalize returns the block as the SDK would store it, so both sides carry every
	// schema field and the hash compares the same shape.
	normalize := func(block map[string]interface{}) map[string]interface{} {
		r := resourceAlkiraConnectorAzureVnet()
		d := r.TestResourceData()
		if err := d.Set("vnet_subnet", schema.NewSet(hashElem, []interface{}{block})); err != nil {
			t.Fatalf("normalize: %v", err)
		}
		return d.Get("vnet_subnet").(*schema.Set).List()[0].(map[string]interface{})
	}

	cases := []struct {
		name     string
		config   map[string]interface{}
		response alkira.ConnectorVnetExportOptionUserInputPrefix
	}{
		{
			name: "single prefix in plural form",
			config: map[string]interface{}{
				"subnet_id":    "/sub/single",
				"subnet_cidrs": schema.NewSet(schema.HashString, []interface{}{"10.0.0.0/24"}),
			},
			// the server answers a one-prefix subnet with value only, whichever
			// form the config used
			response: alkira.ConnectorVnetExportOptionUserInputPrefix{
				Type: "SUBNET", Id: "/sub/single", Value: "10.0.0.0/24",
			},
		},
		{
			name: "single prefix in singular form",
			config: map[string]interface{}{
				"subnet_id":   "/sub/single",
				"subnet_cidr": "10.0.0.0/24",
			},
			response: alkira.ConnectorVnetExportOptionUserInputPrefix{
				Type: "SUBNET", Id: "/sub/single", Value: "10.0.0.0/24",
			},
		},
		{
			name: "multi prefix",
			config: map[string]interface{}{
				"subnet_id": "/sub/multi",
				"subnet_cidrs": schema.NewSet(schema.HashString,
					[]interface{}{"10.169.142.0/28", "10.169.142.32/27"}),
			},
			// a row written by a client that sent both value and values
			response: alkira.ConnectorVnetExportOptionUserInputPrefix{
				Type: "SUBNET", Id: "/sub/multi",
				Value:  "10.169.142.0/28",
				Values: []string{"10.169.142.0/28", "10.169.142.32/27"},
			},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			r := resourceAlkiraConnectorAzureVnet()
			d := r.TestResourceData()
			if err := d.Set("vnet_subnet", schema.NewSet(hashElem, []interface{}{tc.config})); err != nil {
				t.Fatalf("set vnet_subnet: %v", err)
			}

			planned := hashElem(normalize(tc.config))

			setVnetRouting(d, &alkira.ConnectorVnetRouting{
				ExportOptions: alkira.ConnectorVnetExportOptions{
					UserInputPrefixes: []alkira.ConnectorVnetExportOptionUserInputPrefix{tc.response},
				},
			})

			actual := hashElem(d.Get("vnet_subnet").(*schema.Set).List()[0].(map[string]interface{}))

			assert.Equal(t, planned, actual,
				"planned and actual element hashes must match, or core rejects the apply "+
					"with \"planned set element does not correlate with any element in actual\"")
		})
	}
}
