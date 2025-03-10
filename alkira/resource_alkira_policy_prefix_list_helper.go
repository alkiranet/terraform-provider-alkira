package alkira

import (
	"log"

	"github.com/alkiranet/alkira-client-go/alkira"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// setPrefixRange set "prefix_range" blocks from API response
func setPrefixRanges(d *schema.ResourceData, r []alkira.PolicyPrefixListRange) {
	var prefixRanges []map[string]interface{}

	for _, rng := range r {
		prefixRange := map[string]interface{}{
			"prefix": rng.Prefix,
			"le":     rng.Le,
			"ge":     rng.Ge,
		}
		prefixRanges = append(prefixRanges, prefixRange)
	}

	d.Set("prefix_range", prefixRanges)
}

// expandPrefixListPrefixRanges expand block "prefix_range" to
// construct payload
func expandPrefixListPrefixRanges(in []interface{}) ([]alkira.PolicyPrefixListRange, error) {

	if in == nil || len(in) == 0 {
		log.Printf("[DEBUG] prefix_ranges is empty.")
		return nil, nil
	}

	prefixListRanges := make([]alkira.PolicyPrefixListRange, len(in))

	for i, values := range in {
		prefixListRange := alkira.PolicyPrefixListRange{}
		value := values.(map[string]interface{})

		if v, ok := value["prefix"].(string); ok {
			prefixListRange.Prefix = v
		}
		if v, ok := value["le"].(int); ok {
			prefixListRange.Le = v
		}
		if v, ok := value["ge"].(int); ok {
			prefixListRange.Ge = v
		}

		prefixListRanges[i] = prefixListRange
	}

	return prefixListRanges, nil
}
