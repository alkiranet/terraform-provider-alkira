package alkira

import (
	"github.com/alkiranet/alkira-client-go/alkira"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// setVnetRouting set vnet_cidr and vnet_subnet block values
func setVnetRouting(d *schema.ResourceData, routingOptions *alkira.ConnectorVnetRouting) {

	var vnetSubnets []map[string]interface{}
	var vnetCidrs []map[string]interface{}

	// Set vnet_subnet
	for _, prefixes := range routingOptions.ExportOptions.UserInputPrefixes {
		if prefixes.Type == "SUBNET" {
			vnetSubnet := map[string]interface{}{
				"subnet_id":   prefixes.Id,
				"subnet_cidr": prefixes.Value,
			}

			for _, importOptions := range routingOptions.ImportOptions.Subnets {
				if vnetSubnet["subnet_id"] == importOptions.Id {
					vnetSubnet["routing_options"] = importOptions.RouteImportMode
					vnetSubnet["prefix_list_ids"] = importOptions.PrefixListIds
				}
			}

			for _, serviceRoutes := range routingOptions.ServiceRoutes.Subnets {
				if vnetSubnet["subnet_id"] == serviceRoutes.Id {
					vnetSubnet["service_tags"] = serviceRoutes.ServiceTags
				}
			}

			for _, udrLists := range routingOptions.UdrLists.Subnets {
				if vnetSubnet["subnet_id"] == udrLists.Id {
					vnetSubnet["udr_list_ids"] = udrLists.UdrListIds
				}
			}

			vnetSubnets = append(vnetSubnets, vnetSubnet)
		}
	}

	// Set vnet_cidr
	for _, prefixes := range routingOptions.ExportOptions.UserInputPrefixes {
		if prefixes.Type == "CIDR" {
			vnetCidr := map[string]interface{}{
				"cidr": prefixes.Value,
			}

			for _, importOptions := range routingOptions.ImportOptions.Cidrs {
				vnetCidr["routing_options"] = importOptions.RouteImportMode
				vnetCidr["prefix_list_ids"] = importOptions.PrefixListIds
			}

			for _, serviceRoutes := range routingOptions.ServiceRoutes.Cidrs {
				if vnetCidr["cidr"] == serviceRoutes.Value {
					vnetCidr["service_tags"] = serviceRoutes.ServiceTags
				}
			}

			for _, udrLists := range routingOptions.UdrLists.Cidrs {
				if vnetCidr["cidr"] == udrLists.Value {
					vnetCidr["udr_list_ids"] = udrLists.UdrListIds
				}
			}

			vnetCidrs = append(vnetCidrs, vnetCidr)
		}
	}

	d.Set("vnet_subnets", vnetSubnets)
	d.Set("vnet_cidr", vnetCidrs)
}

// constructVnetRouting construct connector_azure_vnet routing options
func constructVnetRouting(d *schema.ResourceData) (*alkira.ConnectorVnetRouting, error) {

	exportOptions := alkira.ConnectorVnetExportOptions{}

	importOptions := alkira.ConnectorVnetImportOptions{}
	importOptions.RouteImportMode = d.Get("routing_options").(string)
	importOptions.PrefixListIds = convertTypeListToIntList(d.Get("routing_prefix_list_ids").([]interface{}))

	serviceRoutes := alkira.ConnectorVnetServiceRoutes{}
	udrLists := alkira.ConnectorVnetUdrLists{}

	// Processing vnet_subnet blocks
	for _, block := range d.Get("vnet_subnet").(*schema.Set).List() {
		content := block.(map[string]interface{})

		// Processing export options for subnet
		subnetUserInputPrefix := alkira.ConnectorVnetExportOptionUserInputPrefix{}
		subnetUserInputPrefix.Type = "SUBNET"

		if v, ok := content["subnet_id"].(string); ok {
			subnetUserInputPrefix.Id = v
		}

		if v, ok := content["subnet_cidr"].(string); ok {
			subnetUserInputPrefix.Value = v
		}

		exportOptions.UserInputPrefixes = append(exportOptions.UserInputPrefixes, subnetUserInputPrefix)

		// Processing import options for subnet
		if _, ok := content["routing_options"].(string); ok {
			subnetImportOption := alkira.ConnectorVnetImportOptionsSubnet{}

			if v, ok := content["subnet_id"].(string); ok {
				subnetImportOption.Id = v
			}

			if v, ok := content["subnet_cidr"].(string); ok {
				subnetImportOption.Value = v
			}

			if v, ok := content["routing_options"].(string); ok {
				subnetImportOption.RouteImportMode = v
			}

			subnetImportOption.PrefixListIds = convertTypeListToIntList(content["prefix_list_ids"].([]interface{}))

			importOptions.Subnets = append(importOptions.Subnets, subnetImportOption)
		}

		// Processing service routes for subnet
		if content["service_tags"] != nil && content["service_tags"].(*schema.Set).Len() > 0 {
			subnetServiceRoute := alkira.ConnectorVnetServiceRoute{}

			if v, ok := content["subnet_id"].(string); ok {
				subnetServiceRoute.Id = v
			}

			if v, ok := content["subnet_cidr"].(string); ok {
				subnetServiceRoute.Value = v
			}

			subnetServiceRoute.ServiceTags = convertTypeSetToStringList(content["service_tags"].(*schema.Set))

			serviceRoutes.Subnets = append(serviceRoutes.Subnets, subnetServiceRoute)
		}

		// Processing UDR list for subnet
		if content["udr_list_ids"] != nil && content["udr_list_ids"].(*schema.Set).Len() > 0 {
			subnetUdrList := alkira.ConnectorVnetUdrList{}

			if v, ok := content["subnet_id"].(string); ok {
				subnetUdrList.Id = v
			}

			if v, ok := content["subnet_cidr"].(string); ok {
				subnetUdrList.Value = v
			}

			subnetUdrList.UdrListIds = convertTypeSetToIntList(content["udr_list_ids"].(*schema.Set))

			udrLists.Subnets = append(udrLists.Subnets, subnetUdrList)
		}
	}

	// Processing vnet_cidr blocks
	for _, block := range d.Get("vnet_cidr").(*schema.Set).List() {
		content := block.(map[string]interface{})

		// Processing export options for CIDR
		cidrUserInputPrefix := alkira.ConnectorVnetExportOptionUserInputPrefix{}
		cidrUserInputPrefix.Type = "CIDR"

		if v, ok := content["cidr"].(string); ok {
			cidrUserInputPrefix.Value = v
		}

		exportOptions.UserInputPrefixes = append(exportOptions.UserInputPrefixes, cidrUserInputPrefix)

		// Processing import options for CIDR
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

		// Processing service routes for CIDR
		if content["service_tags"].(*schema.Set).Len() > 0 {
			cidrServiceRoute := alkira.ConnectorVnetServiceRoute{}

			if v, ok := content["cidr"].(string); ok {
				cidrServiceRoute.Value = v
			}

			cidrServiceRoute.ServiceTags = convertTypeSetToStringList(content["service_tags"].(*schema.Set))

			serviceRoutes.Cidrs = append(serviceRoutes.Cidrs, cidrServiceRoute)
		}

		// Processing UDR lists for CIDR
		if content["udr_list_ids"].(*schema.Set).Len() > 0 {
			cidrUdrList := alkira.ConnectorVnetUdrList{}

			if v, ok := content["cidr"].(string); ok {
				cidrUdrList.Value = v
			}

			cidrUdrList.UdrListIds = convertTypeSetToIntList(content["udr_list_ids"].(*schema.Set))
			udrLists.Cidrs = append(udrLists.Cidrs, cidrUdrList)
		}
	}

	vnetRouting := alkira.ConnectorVnetRouting{
		ExportOptions: exportOptions,
		ImportOptions: importOptions,
		ServiceRoutes: serviceRoutes,
		UdrLists:      udrLists,
	}

	return &vnetRouting, nil
}
