package alkira

import (
	"github.com/alkiranet/alkira-client-go/alkira"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// constructVnetRouting construct connector_azure_vnet routing options
func constructVnetRouting(d *schema.ResourceData) (*alkira.ConnectorVnetRouting, error) {

	importOptions := alkira.ConnectorVnetImportOptions{}
	importOptions.RouteImportMode = d.Get("routing_options").(string)
	importOptions.PrefixListIds = convertTypeListToIntList(d.Get("routing_prefix_list_ids").([]interface{}))

	serviceRoutes := alkira.ConnectorVnetServiceRoutes{}

	// Processing vnet_subnet blocks
	for _, block := range d.Get("vnet_subnet").(*schema.Set).List() {
		content := block.(map[string]interface{})

		if _, ok := content["routing_options"].(string); ok {
			subnetImportOption := alkira.ConnectorVnetImportOptionsSubnet{}

			if v, ok := content["subnet_id"].(string); ok {
				subnetImportOption.Id = v
			}

			if v, ok := content["subnet_value"].(string); ok {
				subnetImportOption.Value = v
			}

			if v, ok := content["routing_options"].(string); ok {
				subnetImportOption.RouteImportMode = v
			}

			subnetImportOption.PrefixListIds = convertTypeListToIntList(content["prefix_list_ids"].([]interface{}))

			importOptions.Subnets = append(importOptions.Subnets, subnetImportOption)
		}

		if len(content["service_tags"].([]interface{})) > 0 {
			subnetServiceRoute := alkira.ConnectorVnetServiceRoute{}

			if v, ok := content["subnet_id"].(string); ok {
				subnetServiceRoute.Id = v
			}

			if v, ok := content["subnet_value"].(string); ok {
				subnetServiceRoute.Value = v
			}

			subnetServiceRoute.ServiceTags = convertTypeListToStringList(content["service_tags"].([]interface{}))

			serviceRoutes.Subnets = append(serviceRoutes.Subnets, subnetServiceRoute)
		}
	}

	// Processing vnet_cidr blocks
	for _, block := range d.Get("vnet_cidr").(*schema.Set).List() {
		content := block.(map[string]interface{})

		if _, ok := content["routing_options"].(string); ok {
			cidrImportOption := alkira.ConnectorVnetImportOptionsCidr{}

			if v, ok := content["cidr"].(string); ok {
				cidrImportOption.Value = v
			}

			if v, ok := content["routing_options"].(string); ok {
				cidrImportOption.RouteImportMode = v
			}

			cidrImportOption.PrefixListIds = convertTypeListToIntList(content["prefix_list_ids"].([]interface{}))

			importOptions.Cidrs = append(importOptions.Cidrs, cidrImportOption)
		}

		if len(content["service_tags"].([]interface{})) > 0 {
			cidrServiceRoute := alkira.ConnectorVnetServiceRoute{}

			if v, ok := content["cidr"].(string); ok {
				cidrServiceRoute.Value = v
			}

			cidrServiceRoute.ServiceTags = convertTypeListToStringList(content["service_tags"].([]interface{}))

			serviceRoutes.Cidrs = append(serviceRoutes.Cidrs, cidrServiceRoute)
		}
	}

	vnetRouting := alkira.ConnectorVnetRouting{
		ImportOptions: importOptions,
		ServiceRoutes: serviceRoutes,
	}

	return &vnetRouting, nil
}