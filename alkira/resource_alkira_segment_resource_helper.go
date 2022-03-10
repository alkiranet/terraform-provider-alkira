package alkira

import (
	"log"

	"github.com/alkiranet/alkira-client-go/alkira"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func expandSegmentResourceGroupPrefix(in *schema.Set) []alkira.SegmentResourceGroupPrefix {
	if in == nil || in.Len() == 0 {
		log.Printf("[DEBUG] invalid input for segment resource group prefix.")
		return nil
	}

	prefixes := make([]alkira.SegmentResourceGroupPrefix, in.Len())

	for i, p := range in.List() {
		prefix := alkira.SegmentResourceGroupPrefix{}
		prefixValue := p.(map[string]interface{})

		if v, ok := prefixValue["group_id"].(int); ok {
			prefix.GroupId = v
		}
		if v, ok := prefixValue["prefix_list_id"].(int); ok {
			prefix.PrefixListId = v
		}

		prefixes[i] = prefix
	}

	return prefixes
}
