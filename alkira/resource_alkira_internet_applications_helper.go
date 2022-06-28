package alkira

import (
	"log"

	"github.com/alkiranet/alkira-client-go/alkira"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func expandInternetApplicationTargets(in *schema.Set) []alkira.InternetApplicationTargets {
	if in == nil || in.Len() == 0 {
		log.Printf("[DEBUG] invalid internet application targets")
		return nil
	}

	targets := make([]alkira.InternetApplicationTargets, in.Len())

	for i, target := range in.List() {
		r := alkira.InternetApplicationTargets{}
		content := target.(map[string]interface{})
		if v, ok := content["type"].(string); ok {
			r.Type = v
		}
		if v, ok := content["value"].(string); ok {
			r.Value = v
		}
		if v, ok := content["port_ranges"].([]interface{}); ok {
			r.PortRanges = convertTypeListToStringList(v)
		}
		targets[i] = r
	}

	return targets
}
