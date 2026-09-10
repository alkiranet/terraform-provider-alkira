package alkira

import (
	"fmt"
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
			// AK-74515: the server sends `values` and no `value` for a subnet
			// onboarded with several prefixes, and the reverse otherwise. Round-trip
			// whichever it sent so the config the user wrote is what comes back -
			// writing subnet_cidr from an absent value would show a permanent diff.
			vnetSubnet := map[string]interface{}{
				"subnet_id":    prefixes.Id,
				"subnet_cidr":  prefixes.Value,
				"subnet_cidrs": prefixes.Values,
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
					vnetSubnet["native_services"] = serviceRoutes.NativeServiceNames
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
					vnetCidr["native_services"] = serviceRoutes.NativeServiceNames
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

	d.Set("vnet_subnet", vnetSubnets)
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

		// AK-74515: a subnet is one block, whether it has one prefix or several.
		// subnet_cidrs names them all and maps to `values`; subnet_cidr names the
		// single one and maps to `value`. They are mutually exclusive on the wire,
		// so only the one the config used is set - sending both, or sending an
		// empty `value` alongside `values`, is not the shape the API expects.
		subnetCidr, _ := content["subnet_cidr"].(string)
		var subnetCidrs []string
		if raw, ok := content["subnet_cidrs"].(*schema.Set); ok && raw != nil {
			subnetCidrs = convertTypeSetToStringList(raw)
		}

		subnetId, _ := content["subnet_id"].(string)
		if subnetCidr == "" && len(subnetCidrs) == 0 {
			return nil, fmt.Errorf(
				"vnet_subnet %q: one of subnet_cidr or subnet_cidrs must be set", subnetId)
		}
		if subnetCidr != "" && len(subnetCidrs) > 0 {
			return nil, fmt.Errorf(
				"vnet_subnet %q: subnet_cidr and subnet_cidrs are mutually exclusive - "+
					"use subnet_cidrs alone to name every prefix of a multi-prefix subnet", subnetId)
		}
		// The same rule runs at plan time in CustomizeDiff (validateVnetSubnetCidrForms);
		// this stays because a saved plan can be applied without re-running the diff.

		// Processing export options for subnet
		subnetUserInputPrefix := alkira.ConnectorVnetExportOptionUserInputPrefix{}
		subnetUserInputPrefix.Type = "SUBNET"

		if v, ok := content["subnet_id"].(string); ok {
			subnetUserInputPrefix.Id = v
		}

		if len(subnetCidrs) > 0 {
			subnetUserInputPrefix.Values = subnetCidrs
		} else {
			subnetUserInputPrefix.Value = subnetCidr
		}

		exportOptions.UserInputPrefixes = append(exportOptions.UserInputPrefixes, subnetUserInputPrefix)

		// Processing import options for subnet
		if _, ok := content["routing_options"].(string); ok {
			subnetImportOption := alkira.ConnectorVnetImportOptionsSubnet{}

			if v, ok := content["subnet_id"].(string); ok {
				subnetImportOption.Id = v
			}

			if len(subnetCidrs) > 0 {
				subnetImportOption.Values = subnetCidrs
			} else {
				subnetImportOption.Value = subnetCidr
			}

			if v, ok := content["routing_options"].(string); ok {
				subnetImportOption.RouteImportMode = v
			}

			subnetImportOption.PrefixListIds = convertTypeListToIntList(content["prefix_list_ids"].([]interface{}))

			importOptions.Subnets = append(importOptions.Subnets, subnetImportOption)
		}

		// Processing service routes for subnet
		if (content["service_tags"] != nil && content["service_tags"].(*schema.Set).Len() > 0) ||
			(content["native_services"] != nil && content["native_services"].(*schema.Set).Len() > 0) {
			subnetServiceRoute := alkira.ConnectorVnetServiceRouteSubnet{}

			if v, ok := content["subnet_id"].(string); ok {
				subnetServiceRoute.Id = v
			}

			if len(subnetCidrs) > 0 {
				subnetServiceRoute.Values = subnetCidrs
			} else {
				subnetServiceRoute.Value = subnetCidr
			}

			subnetServiceRoute.ServiceTags = convertTypeSetToStringList(content["service_tags"].(*schema.Set))
			subnetServiceRoute.NativeServiceNames = convertTypeSetToStringList(content["native_services"].(*schema.Set))

			serviceRoutes.Subnets = append(serviceRoutes.Subnets, subnetServiceRoute)
		}

		// Processing UDR list for subnet
		if content["udr_list_ids"] != nil && content["udr_list_ids"].(*schema.Set).Len() > 0 {
			subnetUdrList := alkira.ConnectorVnetUdrListSubnet{}

			if v, ok := content["subnet_id"].(string); ok {
				subnetUdrList.Id = v
			}

			if len(subnetCidrs) > 0 {
				subnetUdrList.Values = subnetCidrs
			} else {
				subnetUdrList.Value = subnetCidr
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
		if (content["service_tags"] != nil && content["service_tags"].(*schema.Set).Len() > 0) ||
			(content["native_services"] != nil && content["native_services"].(*schema.Set).Len() > 0) {

			cidrServiceRoute := alkira.ConnectorVnetServiceRoute{}

			if v, ok := content["cidr"].(string); ok {
				cidrServiceRoute.Value = v
			}

			cidrServiceRoute.ServiceTags = convertTypeSetToStringList(content["service_tags"].(*schema.Set))
			cidrServiceRoute.NativeServiceNames = convertTypeSetToStringList(content["native_services"].(*schema.Set))

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

// validateVnetSubnetCidrForms enforces that each vnet_subnet block names its prefixes
// exactly one way. AK-74515: subnet_cidr for a single-prefix subnet, subnet_cidrs for a
// multi-prefix one - the wire format sets one field or the other, never both.
//
// Shared by CustomizeDiff (plan time) and constructVnetRouting (apply time): a plan can
// be applied from a saved file, so the apply-side check stays as the real guard.
func validateVnetSubnetCidrForms(raw interface{}) error {
	set, ok := raw.(*schema.Set)
	if !ok || set == nil {
		return nil
	}

	for _, block := range set.List() {
		content, ok := block.(map[string]interface{})
		if !ok {
			continue
		}

		subnetId, _ := content["subnet_id"].(string)
		subnetCidr, _ := content["subnet_cidr"].(string)

		n := 0
		if cidrs, ok := content["subnet_cidrs"].(*schema.Set); ok && cidrs != nil {
			n = cidrs.Len()
		}

		if subnetCidr == "" && n == 0 {
			return fmt.Errorf(
				"vnet_subnet %q: one of subnet_cidr or subnet_cidrs must be set", subnetId)
		}
		if subnetCidr != "" && n > 0 {
			return fmt.Errorf(
				"vnet_subnet %q: subnet_cidr and subnet_cidrs are mutually exclusive - "+
					"use subnet_cidrs alone to name every prefix of a multi-prefix subnet", subnetId)
		}
	}

	return nil
}
