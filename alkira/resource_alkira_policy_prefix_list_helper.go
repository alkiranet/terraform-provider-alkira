package alkira

import (
	"log"

	"github.com/alkiranet/alkira-client-go/alkira"
)

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
