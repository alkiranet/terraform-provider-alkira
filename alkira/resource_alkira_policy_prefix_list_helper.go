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
			"prefix":      rng.Prefix,
			"le":          rng.Le,
			"ge":          rng.Ge,
			"description": rng.Description,
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
		if v, ok := value["description"].(string); ok {
			prefixListRange.Description = v
		}

		prefixListRanges[i] = prefixListRange
	}

	return prefixListRanges, nil
}

func extractPrefixes(d *schema.ResourceData) []string {
	var prefixes []string
	if v, ok := d.GetOk("prefixes"); ok {
		for _, p := range v.([]interface{}) {
			prefixMap := p.(map[string]interface{})
			prefixes = append(prefixes, prefixMap["prefix"].(string))
		}
	}
	return prefixes
}

func expandPrefixListPrefixes(d *schema.ResourceData) ([]string, map[string]*alkira.PolicyPrefixListDetails) {

	prefixes := extractPrefixes(d)
	prefixMap := buildPrefixDetailsMap(d)
	return prefixes, prefixMap

}

func buildPrefixDetailsMap(d *schema.ResourceData) map[string]*alkira.PolicyPrefixListDetails {
	details := make(map[string]*alkira.PolicyPrefixListDetails)
	if v, ok := d.GetOk("prefixes"); ok {
		for _, p := range v.([]interface{}) {
			prefixMap := p.(map[string]interface{})
			prefix := prefixMap["prefix"].(string)
			if desc, ok := prefixMap["description"].(string); ok && desc != "" {
				details[prefix] = &alkira.PolicyPrefixListDetails{Description: desc}
			}
		}
	}
	return details
}

func setPrefixes(d *schema.ResourceData, prefixes []string, details map[string]*alkira.PolicyPrefixListDetails) {
	var prefixList []map[string]interface{}
	for _, p := range prefixes {
		prefixEntry := map[string]interface{}{"prefix": p}
		if details[p] != nil {
			prefixEntry["description"] = details[p].Description
		}
		prefixList = append(prefixList, prefixEntry)
	}
	d.Set("prefixes", prefixList)
}
