package alkira

import (
	"fmt"
	"github.com/alkiranet/alkira-client-go/alkira"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"log"
)

func convertGcpRouting(in *schema.Set, subnets *schema.Set) (*alkira.ConnectorGcpVpcRouting, error) {
	importOptions := alkira.ConnectorGcpVpcImportOptions{
		RouteImportMode: "ADVERTISE_DEFAULT_ROUTE",
	}

	if in != nil && in.Len() == 1 {
		for _, option := range in.List() {
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

			if (t["id"] == "" && t["fq_id"] == "") || t["cidr"] == "" {
				log.Printf("[ERROR] both id %s and cidr %s must be populated", t["id"], t["cidr"])
				return nil, fmt.Errorf("[ERROR] both id %s and cidr %s must be populated", t["id"], t["cidr"])
			}

			if v, ok := t["id"].(string); ok {
				r.Id = v
			}
			if v, ok := t["fq_id"].(string); ok {
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

	return nil, nil
}

func setGcpRoutingOptions(c *alkira.ConnectorGcpVpcRouting, d *schema.ResourceData) {
	in := make(map[string]interface{})
	in["prefix_list_ids"] = c.ImportOptions.PrefixListIds
	in["custom_prefix"] = c.ImportOptions.RouteImportMode

	r := resourceAlkiraConnectorGcpVpc()
	f := schema.HashResource(r)
	s := schema.NewSet(f, []interface{}{in})
	d.Set("gcp_routing", s)
}
