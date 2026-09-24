package alkira

import (
	"fmt"

	"github.com/alkiranet/alkira-client-go/alkira"
	"github.com/hashicorp/go-cty/cty"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// subnetIdsUsingPluralForm reports which vnet_subnet blocks already express their
// prefixes as subnet_cidrs, so Read can write each block back in the form its config
// uses. AK-74515: a one-element subnet_cidrs is sent as `value`, so the response alone
// cannot distinguish it from a subnet_cidr block - without this, such a block would come
// back in the singular form and show a permanent diff.
func subnetIdsUsingPluralForm(d *schema.ResourceData) map[string]bool {
	plural := map[string]bool{}

	set, ok := d.Get("vnet_subnet").(*schema.Set)
	if !ok || set == nil {
		return plural
	}

	for _, block := range set.List() {
		content, ok := block.(map[string]interface{})
		if !ok {
			continue
		}
		subnetId, _ := content["subnet_id"].(string)
		if subnetId == "" {
			continue
		}
		if cidrs, ok := content["subnet_cidrs"].(*schema.Set); ok && cidrs != nil && cidrs.Len() > 0 {
			plural[subnetId] = true
		}
	}

	return plural
}

// setVnetRouting set vnet_cidr and vnet_subnet block values
func setVnetRouting(d *schema.ResourceData, routingOptions *alkira.ConnectorVnetRouting) {

	var vnetSubnets []map[string]interface{}
	var vnetCidrs []map[string]interface{}

	// AK-74515: which subnet_ids the config/state writes in the plural form. The two
	// forms are interchangeable on the wire - a single prefix travels as `value`
	// whichever one the user wrote - so the response shape cannot tell us which to
	// write back. Keying off what is already there keeps the round trip stable and
	// removes the read path's dependence on the response shape entirely.
	pluralSubnetIds := subnetIdsUsingPluralForm(d)

	// Set vnet_subnet
	for _, prefixes := range routingOptions.ExportOptions.UserInputPrefixes {
		if prefixes.Type == "SUBNET" {
			// A response for a multi-prefix subnet carries BOTH `values` and `value` -
			// the server keeps `value` as the representative prefix for older clients -
			// so a non-empty `value` never by itself means "single prefix". Write the
			// form the config used: leaving state in the other form makes the planned
			// and actual TypeSet element hashes differ, which core rejects after apply
			// ("planned set element does not correlate with any element in actual") and
			// which otherwise shows the block as removed-and-re-added on every plan.
			vnetSubnet := map[string]interface{}{
				"subnet_id": prefixes.Id,
			}

			cidrs := prefixes.Values
			if len(cidrs) == 0 && prefixes.Value != "" {
				cidrs = []string{prefixes.Value}
			}

			if pluralSubnetIds[prefixes.Id] || len(prefixes.Values) > 1 {
				vnetSubnet["subnet_cidrs"] = cidrs
			} else if len(cidrs) > 0 {
				vnetSubnet["subnet_cidr"] = cidrs[0]
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

		// AK-74515: a ONE-element subnet_cidrs is the singular form on the wire. The
		// server returns `value` (not a one-element `values`) for a subnet onboarded
		// with a single prefix, so sending `values` here could not round-trip: Read
		// would populate subnet_cidr and leave subnet_cidrs empty, inverting the config
		// and failing the apply. Nothing requires subnet_cidrs to hold two or more
		// entries - the docs only say exactly one of the two forms is required - so
		// collapsing here, after validation, lets a user write the plural form
		// uniformly. Every branch below reads these two variables, so normalizing once
		// keeps the export, import, service-route and UDR payloads consistent.
		if len(subnetCidrs) == 1 {
			subnetCidr = subnetCidrs[0]
			subnetCidrs = nil
		}

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
// exactly one way, at plan time.
//
// AK-74515: it walks the RAW CONFIG rather than d.Get. On a ResourceDiff, d.Get reads an
// unknown value as the zero value, so a block whose CIDRs are not resolved yet - e.g.
// subnet_cidrs = azurerm_subnet.pe.address_prefixes, where the subnet is created in the
// same run - would read as "neither form set" and fail the plan. Walking the raw config
// lets us tell "not set" apart from "not known yet" and skip the latter; the apply-time
// check in constructVnetRouting still catches a genuinely empty block once values are
// known.
func validateVnetSubnetCidrForms(rawConfig cty.Value) error {
	if rawConfig.IsNull() || !rawConfig.IsKnown() {
		return nil
	}

	subnets := rawConfig.GetAttr("vnet_subnet")
	if subnets.IsNull() || !subnets.IsKnown() {
		return nil
	}

	for it := subnets.ElementIterator(); it.Next(); {
		_, block := it.Element()
		if block.IsNull() || !block.IsKnown() {
			continue
		}

		subnetId := ""
		if id := block.GetAttr("subnet_id"); id.IsKnown() && !id.IsNull() {
			subnetId = id.AsString()
		}

		cidr := block.GetAttr("subnet_cidr")
		cidrs := block.GetAttr("subnet_cidrs")

		// Unknown on either side means the form cannot be judged yet. The
		// mutual-exclusion arm below only under-fires in that case, and apply
		// re-checks both rules.
		cidrKnown := cidr.IsKnown()
		cidrsKnown := cidrs.IsKnown()

		hasCidr := cidrKnown && !cidr.IsNull() && cidr.AsString() != ""
		hasCidrs := cidrsKnown && !cidrs.IsNull() && cidrs.LengthInt() > 0

		if cidrKnown && cidrsKnown && !hasCidr && !hasCidrs {
			return fmt.Errorf(
				"vnet_subnet %q: one of subnet_cidr or subnet_cidrs must be set", subnetId)
		}
		if hasCidr && hasCidrs {
			return fmt.Errorf(
				"vnet_subnet %q: subnet_cidr and subnet_cidrs are mutually exclusive - "+
					"use subnet_cidrs alone to name every prefix of a multi-prefix subnet", subnetId)
		}
	}

	return nil
}
