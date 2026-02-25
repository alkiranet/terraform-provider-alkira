package alkira

import (
	"fmt"
	"log"

	"github.com/alkiranet/alkira-client-go/alkira"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func expandGcpRouting(in []interface{}, subnets *schema.Set) (*alkira.ConnectorGcpVpcRouting, error) {

	importOptions := alkira.ConnectorGcpVpcImportOptions{
		RouteImportMode: "ADVERTISE_DEFAULT_ROUTE",
	}

	if in != nil && len(in) == 1 {
		for _, option := range in {
			cfg := option.(map[string]interface{})

			if v, ok := cfg["prefix_list_ids"].([]interface{}); ok {
				importOptions.PrefixListIds = convertTypeListToIntList(v)
			}

			if v, ok := cfg["custom_prefix"].(string); ok {
				importOptions.RouteImportMode = v
			}
		}
	}

	exportAllSubnets := true

	prefixes, err := generateGCPUserInputPrefixes(subnets)

	if err != nil {
		return nil, err
	}

	if prefixes != nil && len(prefixes) > 0 {
		exportAllSubnets = false
	}

	exportOptions := alkira.ConnectorGcpVpcExportOptions{
		ExportAllSubnets: exportAllSubnets,
		Prefixes:         prefixes,
	}

	gcp := &alkira.ConnectorGcpVpcRouting{
		ExportOptions: exportOptions,
		ImportOptions: importOptions,
	}
	return gcp, nil
}

// generateUserInputPrefixes generate UserInputPrefixes used in GCP-VPC connector
func generateGCPUserInputPrefixes(subnets *schema.Set) ([]alkira.UserInputPrefixes, error) {

	if subnets != nil && subnets.Len() > 0 {

		prefixes := make([]alkira.UserInputPrefixes, subnets.Len())

		for i, subnet := range subnets.List() {
			r := alkira.UserInputPrefixes{}
			t := subnet.(map[string]interface{})

			internalId := ""
			if v, ok := t["internal_id"].(string); ok && v != "" {
				internalId = v
			}

			if internalId == "" && (t["id"] == "" || t["cidr"] == "") {
				log.Printf("[ERROR] subnet configuration must have either internal_id or both id and cidr")
				return nil, fmt.Errorf("[ERROR] subnet configuration must have either internal_id or both id and cidr")
			}

			// Set the internal Alkira ID
			r.Id = internalId

			if v, ok := t["id"].(string); ok {
				r.FqId = v
			}
			if v, ok := t["cidr"].(string); ok {
				r.Value = v
			}

			r.Type = "SUBNET"

			prefixes[i] = r
		}

		return prefixes, nil
	}

	return []alkira.UserInputPrefixes{}, nil
}

func setGcpRoutingOptions(c *alkira.ConnectorGcpVpcRouting, d *schema.ResourceData) {

	if c == nil {
		return
	}

	in := make(map[string]interface{})

	in["prefix_list_ids"] = c.ImportOptions.PrefixListIds
	in["custom_prefix"] = c.ImportOptions.RouteImportMode

	d.Set("gcp_routing", []interface{}{in})
}

func setGcpVpcSubnets(c *alkira.ConnectorGcpVpcRouting, d *schema.ResourceData) {
	if c == nil {
		return
	}

	prefixes := c.ExportOptions.Prefixes
	if len(prefixes) == 0 {
		return
	}

	subnets := make([]interface{}, len(prefixes))
	for i, prefix := range prefixes {
		subnet := make(map[string]interface{})
		// Store both FqId (user-provided) and Id (internal Alkira ID)
		subnet["id"] = prefix.FqId
		subnet["internal_id"] = prefix.Id
		subnet["cidr"] = prefix.Value
		subnets[i] = subnet
	}

	d.Set("vpc_subnet", subnets)
}

func generateConnectorGcpVpcRequest(d *schema.ResourceData, m interface{}) (*alkira.ConnectorGcpVpc, error) {

	//
	// Routing
	//
	gcpRouting, err := expandGcpRouting(d.Get("gcp_routing").([]interface{}), d.Get("vpc_subnet").(*schema.Set))

	if err != nil {
		log.Printf("[ERROR] failed to convert gcp routing")
		return nil, err
	}

	//
	// Segment
	//
	segmentName, err := getSegmentNameById(d.Get("segment_id").(string), m)

	if err != nil {
		return nil, err
	}

	// Assemble request
	connector := &alkira.ConnectorGcpVpc{
		BillingTags:    convertTypeSetToIntList(d.Get("billing_tag_ids").(*schema.Set)),
		CXP:            d.Get("cxp").(string),
		CredentialId:   d.Get("credential_id").(string),
		GcpRouting:     gcpRouting,
		CustomerRegion: d.Get("gcp_region").(string),
		Enabled:        d.Get("enabled").(bool),
		Group:          d.Get("group").(string),
		Name:           d.Get("name").(string),
		ProjectId:      d.Get("gcp_project_id").(string),
		Segments:       []string{segmentName},
		SecondaryCXPs:  convertTypeSetToStringList(d.Get("failover_cxps").(*schema.Set)),
		Size:           d.Get("size").(string),
		VpcName:        d.Get("gcp_vpc_name").(string),
		CustomerASN:    d.Get("customer_asn").(int),
		ScaleGroupId:   d.Get("scale_group_id").(string),
		Description:    d.Get("description").(string),
	}

	return connector, nil
}
