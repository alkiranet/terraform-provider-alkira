package alkira

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/alkiranet/alkira-client-go/alkira"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// setAwsVpcRoutingOptions sets the routing configuration fields from the API response
func setAwsVpcRoutingOptions(connector *alkira.ConnectorAwsVpc, d *schema.ResourceData) {
	if connector.VpcRouting == nil {
		log.Printf("[DEBUG] VpcRouting is nil, skipping routing options")
		return
	}

	// Unmarshal the interface{} to concrete types
	routingJSON, err := json.Marshal(connector.VpcRouting)
	if err != nil {
		log.Printf("[ERROR] Failed to marshal VpcRouting: %v", err)
		return
	}

	var routing alkira.ConnectorAwsVpcRouting
	if err := json.Unmarshal(routingJSON, &routing); err != nil {
		log.Printf("[ERROR] Failed to unmarshal VpcRouting: %v", err)
		return
	}

	// Set export-related fields (vpc_cidr, vpc_subnet, overlay_subnets)
	setAwsVpcExportPrefixes(routing.Export, d)

	// Set import-related fields (vpc_route_table)
	setAwsVpcImportRouteTables(routing.Import, d)
}

// setAwsVpcExportPrefixes sets export-related fields from ExportOptions
func setAwsVpcExportPrefixes(exportOptions interface{}, d *schema.ResourceData) {
	if exportOptions == nil {
		log.Printf("[DEBUG] Export options is nil, skipping export prefixes")
		return
	}

	// Unmarshal the interface{} to ExportOptions
	exportJSON, err := json.Marshal(exportOptions)
	if err != nil {
		log.Printf("[ERROR] Failed to marshal ExportOptions: %v", err)
		return
	}

	var export alkira.ExportOptions
	if err := json.Unmarshal(exportJSON, &export); err != nil {
		log.Printf("[ERROR] Failed to unmarshal ExportOptions: %v", err)
		return
	}

	var cidrList []string
	var subnetList []interface{}
	var overlaySubnets []string

	for _, prefix := range export.Prefixes {
		switch prefix.Type {
		case "CIDR":
			cidrList = append(cidrList, prefix.Value)
		case "SUBNET":
			subnet := map[string]interface{}{
				"id":   prefix.Id,
				"cidr": prefix.Value,
			}
			subnetList = append(subnetList, subnet)
		case "OVERLAY_SUBNETS":
			overlaySubnets = append(overlaySubnets, prefix.Value)
		default:
			log.Printf("[DEBUG] Unknown prefix type: %s", prefix.Type)
		}
	}

	// Only set non-empty fields to preserve config when field is not in config
	if len(cidrList) > 0 {
		d.Set("vpc_cidr", cidrList)
	}
	if len(subnetList) > 0 {
		d.Set("vpc_subnet", subnetList)
	}
	if len(overlaySubnets) > 0 {
		d.Set("overlay_subnets", overlaySubnets)
	}
}

// setAwsVpcImportRouteTables sets vpc_route_table from ImportOptions
func setAwsVpcImportRouteTables(importOptions interface{}, d *schema.ResourceData) {
	if importOptions == nil {
		log.Printf("[DEBUG] Import options is nil, skipping route tables")
		return
	}

	// Unmarshal the interface{} to ImportOptions
	importJSON, err := json.Marshal(importOptions)
	if err != nil {
		log.Printf("[ERROR] Failed to marshal ImportOptions: %v", err)
		return
	}

	var importOpts alkira.ImportOptions
	if err := json.Unmarshal(importJSON, &importOpts); err != nil {
		log.Printf("[ERROR] Failed to unmarshal ImportOptions: %v", err)
		return
	}

	if len(importOpts.RouteTables) == 0 {
		return
	}

	routeTables := make([]interface{}, len(importOpts.RouteTables))
	for i, rt := range importOpts.RouteTables {
		routeTable := map[string]interface{}{
			"id":              rt.Id,
			"options":         rt.Mode,
			"prefix_list_ids": rt.PrefixListIds,
		}
		routeTables[i] = routeTable
	}

	d.Set("vpc_route_table", routeTables)
}

// setTgwAttachment set tgw_attachment blocks
func setTgwAttachment(d *schema.ResourceData, tgwAttachments []alkira.TgwAttachment) {

	var attachments []map[string]interface{}

	for _, each := range d.Get("tgw_attachment").([]interface{}) {
		a := each.(map[string]interface{})

		for _, cfg := range tgwAttachments {
			if a["subnet_id"].(string) == cfg.SubnetId {
				attachments = append(attachments, a)
				break
			}
		}
	}

	//
	// Go through all tgw_attachments from the API response one more
	// time to find any attachement that has not been tracked from
	// Terraform config.
	//
	for _, cfg := range tgwAttachments {
		new := true

		// Check if the gateway already exists in the Terraform config
		for _, tgw := range d.Get("tgw_attachment").([]interface{}) {
			a := tgw.(map[string]interface{})

			if a["subnet_id"].(string) == cfg.SubnetId {
				new = false
				break
			}
		}

		// If the attachment is new, add it to the tail of the list,
		// this will generate a diff
		if new {
			attachment := map[string]interface{}{
				"subnet_id": cfg.SubnetId,
				"az":        cfg.AvailabilityZone,
			}

			attachments = append(attachments, attachment)
		}
	}

	d.Set("tgw_attachment", attachments)

}

// expandAwsVpcRouteTables expand AWS-VPC route tables
func expandAwsVpcRouteTables(in *schema.Set) []alkira.RouteTables {
	if in == nil || in.Len() == 0 {
		log.Printf("[DEBUG] Empty VPC route table input")
		return []alkira.RouteTables{}
	}

	tables := make([]alkira.RouteTables, in.Len())
	for i, table := range in.List() {
		r := alkira.RouteTables{}
		t := table.(map[string]interface{})
		if v, ok := t["id"].(string); ok {
			r.Id = v
		}
		if v, ok := t["options"].(string); ok {
			r.Mode = v
		}

		r.PrefixListIds = convertTypeSetToIntList(t["prefix_list_ids"].(*schema.Set))
		tables[i] = r
	}

	return tables
}

// expandUserInputPrefixes generate UserInputPrefixes used in AWS-VPC connector
func expandUserInputPrefixes(cidr []interface{}, subnets *schema.Set, overlaySubnets []interface{}) ([]alkira.InputPrefixes, error) {

	if len(cidr) == 0 && subnets == nil {
		return nil, fmt.Errorf("ERROR: either \"vpc_subnet\" or \"vpc_cidr\" must be specified")
	}

	// Processing overlay_subnets
	log.Printf("[DEBUG] Processing overlay_subnets %v", overlaySubnets)
	overlaySubnetList := make([]alkira.InputPrefixes, len(overlaySubnets))

	if len(overlaySubnets) > 0 {
		for i, value := range overlaySubnets {
			overlaySubnetList[i].Value = value.(string)
			overlaySubnetList[i].Type = "OVERLAY_SUBNETS"
		}
	}

	// Processing vpc_cidr
	if len(cidr) > 0 {
		log.Printf("[DEBUG] Processing vpc_cidr %v", cidr)
		cidrList := make([]alkira.InputPrefixes, len(cidr))

		for i, value := range cidr {
			cidrList[i].Value = value.(string)
			cidrList[i].Type = "CIDR"
		}

		cidrList = append(cidrList, overlaySubnetList...)
		return cidrList, nil
	}

	// Processing vpc_subnet
	log.Printf("[DEBUG] Processing vpc_subnet")
	if subnets == nil || subnets.Len() == 0 {
		log.Printf("[DEBUG] Empty vpc_subnet")
		return nil, fmt.Errorf("ERROR: Invalid vpc_subnet")
	}

	prefixes := make([]alkira.InputPrefixes, subnets.Len())
	for i, subnet := range subnets.List() {
		r := alkira.InputPrefixes{}
		t := subnet.(map[string]interface{})
		if v, ok := t["id"].(string); ok {
			r.Id = v
		}
		if v, ok := t["cidr"].(string); ok {
			r.Value = v
		}

		r.Type = "SUBNET"
		prefixes[i] = r
	}

	prefixes = append(prefixes, overlaySubnetList...)
	return prefixes, nil
}

// expandAwsVpcTgwAttachments expand tgw_attachment of connector_aws_vpc
func expandAwsVpcTgwAttachments(in []interface{}) []alkira.TgwAttachment {
	if in == nil || len(in) == 0 {
		return []alkira.TgwAttachment{}
	}

	attachments := make([]alkira.TgwAttachment, len(in))
	for i, attachment := range in {
		r := alkira.TgwAttachment{}
		t := attachment.(map[string]interface{})

		if v, ok := t["subnet_id"].(string); ok {
			r.SubnetId = v
		}
		if v, ok := t["az"].(string); ok {
			r.AvailabilityZone = v
		}

		attachments[i] = r
	}

	return attachments
}

// generateConnectorAwsVpcRequest generate request for connector_aws_vpc
func generateConnectorAwsVpcRequest(d *schema.ResourceData, m interface{}) (*alkira.ConnectorAwsVpc, error) {

	// Segment
	segmentName, err := getSegmentNameById(d.Get("segment_id").(string), m)

	if err != nil {
		return nil, err
	}

	inputPrefixes, err := expandUserInputPrefixes(d.Get("vpc_cidr").([]interface{}), d.Get("vpc_subnet").(*schema.Set), d.Get("overlay_subnets").([]interface{}))

	if err != nil {
		return nil, err
	}

	exportOptions := alkira.ExportOptions{
		Mode:     "USER_INPUT_PREFIXES",
		Prefixes: inputPrefixes,
	}

	routeTables := expandAwsVpcRouteTables(d.Get("vpc_route_table").(*schema.Set))
	tgwAttachments := expandAwsVpcTgwAttachments(d.Get("tgw_attachment").([]interface{}))

	vpcRouting := alkira.ConnectorAwsVpcRouting{
		Export: exportOptions,
		Import: alkira.ImportOptions{RouteTables: routeTables},
	}

	request := &alkira.ConnectorAwsVpc{
		BillingTags:                        convertTypeSetToIntList(d.Get("billing_tag_ids").(*schema.Set)),
		CXP:                                d.Get("cxp").(string),
		CredentialId:                       d.Get("credential_id").(string),
		CustomerName:                       m.(*alkira.AlkiraClient).Username,
		CustomerRegion:                     d.Get("aws_region").(string),
		DirectInterVPCCommunicationEnabled: d.Get("direct_inter_vpc_communication_enabled").(bool),
		DirectInterVPCCommunicationGroup:   d.Get("direct_inter_vpc_communication_group").(string),
		Enabled:                            d.Get("enabled").(bool),
		Group:                              d.Get("group").(string),
		Name:                               d.Get("name").(string),
		Segments:                           []string{segmentName},
		SecondaryCXPs:                      convertTypeSetToStringList(d.Get("failover_cxps").(*schema.Set)),
		Size:                               d.Get("size").(string),
		TgwConnectEnabled:                  d.Get("tgw_connect_enabled").(bool),
		TgwAttachments:                     tgwAttachments,
		VpcId:                              d.Get("vpc_id").(string),
		VpcOwnerId:                         d.Get("aws_account_id").(string),
		VpcRouting:                         vpcRouting,
		ScaleGroupId:                       d.Get("scale_group_id").(string),
		Description:                        d.Get("description").(string),
	}

	return request, nil
}
