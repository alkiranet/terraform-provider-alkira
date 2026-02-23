package alkira

import (
	"fmt"
	"log"

	"github.com/alkiranet/alkira-client-go/alkira"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// setPrefixRange set "prefix_range" blocks from API response
func setPrefixRanges(d *schema.ResourceData, r []alkira.PolicyPrefixListRange) {
	set := schema.NewSet(prefixRangeHash, nil)

	for _, rng := range r {
		prefixRange := map[string]interface{}{
			"prefix": rng.Prefix,
			"le":     rng.Le,
			"ge":     rng.Ge,
		}
		if rng.Description != "" {
			prefixRange["description"] = rng.Description
		}
		set.Add(prefixRange)
	}

	d.Set("prefix_range", set)
}

// expandPrefixListPrefixRanges expand block "prefix_range" to
// construct payload
func expandPrefixListPrefixRanges(in *schema.Set) ([]alkira.PolicyPrefixListRange, error) {

	if in == nil || in.Len() == 0 {
		log.Printf("[DEBUG] prefix_ranges is empty.")
		return nil, nil
	}

	values := in.List()
	prefixListRanges := make([]alkira.PolicyPrefixListRange, len(values))

	for i, value := range values {
		prefixListRange := alkira.PolicyPrefixListRange{}
		valueMap := value.(map[string]interface{})

		if v, ok := valueMap["prefix"].(string); ok {
			prefixListRange.Prefix = v
		}
		if v, ok := valueMap["le"].(int); ok {
			prefixListRange.Le = v
		}
		if v, ok := valueMap["ge"].(int); ok {
			prefixListRange.Ge = v
		}
		if v, ok := valueMap["description"].(string); ok {
			prefixListRange.Description = v
		}

		prefixListRanges[i] = prefixListRange
	}

	return prefixListRanges, nil
}

func extractPrefixes(d *schema.ResourceData) []string {
	var prefixes []string

	if v, ok := d.GetOk("prefix"); ok {
		for _, p := range v.(*schema.Set).List() {
			prefixMap := p.(map[string]interface{})
			prefixes = append(prefixes, prefixMap["cidr"].(string))
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

	if v, ok := d.GetOk("prefix"); ok {
		for _, p := range v.(*schema.Set).List() {

			prefixMap := p.(map[string]interface{})
			prefix := prefixMap["cidr"].(string)

			if desc, ok := prefixMap["description"].(string); ok && desc != "" {
				details[prefix] = &alkira.PolicyPrefixListDetails{Description: desc}
			}
		}
	}
	return details
}

// setPrefix Set prefix block when reading from API
func setPrefix(d *schema.ResourceData, prefixes []string, details map[string]*alkira.PolicyPrefixListDetails) {
	set := schema.NewSet(prefixHash, nil)

	for _, p := range prefixes {
		prefixEntry := map[string]interface{}{
			"cidr": p,
		}
		if details[p] != nil && details[p].Description != "" {
			prefixEntry["description"] = details[p].Description
		}
		set.Add(prefixEntry)
	}

	d.Set("prefix", set)
}

// generatePolicyPrefixListRequest
func generatePolicyPrefixListRequest(d *schema.ResourceData) (*alkira.PolicyPrefixList, error) {

	prefixRanges, err := expandPrefixListPrefixRanges(d.Get("prefix_range").(*schema.Set))

	if err != nil {
		return nil, err
	}

	if d.Get("prefixes").(*schema.Set).Len() > 0 {
		return nil, fmt.Errorf("ERROR: Please use the new 'prefix' block to replace the old 'prefixes' field")
	}

	prefixes, prefixDetailsMap := expandPrefixListPrefixes(d)

	list := &alkira.PolicyPrefixList{
		Description:   d.Get("description").(string),
		Name:          d.Get("name").(string),
		Prefixes:      prefixes,
		PrefixDetails: prefixDetailsMap,
		PrefixRanges:  prefixRanges,
	}

	return list, nil
}
