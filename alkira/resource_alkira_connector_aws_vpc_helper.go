package alkira

import (
	"fmt"
	"log"

	"github.com/alkiranet/alkira-client-go/alkira"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// setTgwAttachement set tgw_attachement blocks
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
		for _, tgw := range d.Get("tgw_attachement").([]interface{}) {
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
			break
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
func expandUserInputPrefixes(cidr []interface{}, subnets *schema.Set) ([]alkira.InputPrefixes, error) {

	if len(cidr) == 0 && subnets == nil {
		return nil, fmt.Errorf("ERROR: either \"vpc_subnet\" or \"vpc_cidr\" must be specified.")
	}

	// Processing "vpc_cidr"
	if len(cidr) > 0 {
		log.Printf("[DEBUG] Processing vpc_cidr %v", cidr)
		cidrList := make([]alkira.InputPrefixes, len(cidr))

		for i, value := range cidr {
			cidrList[i].Value = value.(string)
			cidrList[i].Type = "CIDR"
		}

		return cidrList, nil
	}

	// Processing VPC subnets
	log.Printf("[DEBUG] Processing vpc_subnet")
	if subnets == nil || subnets.Len() == 0 {
		log.Printf("[DEBUG] Empty vpc_subnet")
		return nil, fmt.Errorf("ERROR: Invalid vpc_subnet.")
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

	inputPrefixes, err := expandUserInputPrefixes(d.Get("vpc_cidr").([]interface{}), d.Get("vpc_subnet").(*schema.Set))

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
		Import: alkira.ImportOptions{routeTables},
	}

	request := &alkira.ConnectorAwsVpc{
		BillingTags:                        convertTypeSetToIntList(d.Get("billing_tag_ids").(*schema.Set)),
		CXP:                                d.Get("cxp").(string),
		CredentialId:                       d.Get("credential_id").(string),
		CustomerName:                       m.(*alkira.AlkiraClient).Username,
		CustomerRegion:                     d.Get("aws_region").(string),
		DirectInterVPCCommunicationEnabled: d.Get("direct_inter_vpc_communication").(bool),
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
	}

	return request, nil
}
