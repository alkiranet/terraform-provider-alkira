package alkira

import (
	"log"

	"github.com/alkiranet/alkira-client-go/alkira"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func convertGcpRouting(in *schema.Set) *alkira.ConnectorGcpVpcRouting {
	if in == nil || in.Len() > 1 {
		log.Printf("[DEBUG] Only one object allowed in gcp routing options")
		return nil
	}

	if in.Len() < 1 {
		return nil
	}

	gcp := &alkira.ConnectorGcpVpcRouting{
		ImportOptions: alkira.ConnectorGcpVpcImportOptions{},
	}

	for _, option := range in.List() {
		cfg := option.(map[string]interface{})

		if v, ok := cfg["prefix_list_ids"].([]interface{}); ok {
			gcp.ImportOptions.PrefixListIds = convertTypeListToIntList(v)
		}

		if v, ok := cfg["custom_prefix"].(string); ok {
			gcp.ImportOptions.RouteImportMode = v
		}
	}

	return gcp
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
