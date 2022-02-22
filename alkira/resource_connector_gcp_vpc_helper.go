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
