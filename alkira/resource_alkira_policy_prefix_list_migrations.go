package alkira

import (
	"context"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// resourcePolicyPrefixListV0 returns the V0 schema (TypeList) used before
// the migration to TypeSet. This is needed by the state upgrader to parse
// old state formats.
func resourcePolicyPrefixListV0() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"name":            {Type: schema.TypeString, Required: true},
			"description":     {Type: schema.TypeString, Optional: true},
			"provision_state": {Type: schema.TypeString, Computed: true},
			"prefixes": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"prefix": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"cidr":        {Type: schema.TypeString, Required: true},
						"description": {Type: schema.TypeString, Optional: true},
					},
				},
			},
			"prefix_range": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"prefix":      {Type: schema.TypeString, Required: true},
						"description": {Type: schema.TypeString, Optional: true},
						"le":          {Type: schema.TypeInt, Optional: true},
						"ge":          {Type: schema.TypeInt, Optional: true},
					},
				},
			},
		},
	}
}

// resourcePolicyPrefixListStateUpgradeV0 migrates the state from V0
// (TypeList for prefix and prefix_range) to V1 (TypeSet). The flat-map
// state keys change from positional indices (e.g. "prefix.0.cidr") to
// hash-based keys (e.g. "prefix.<hash>.cidr").
func resourcePolicyPrefixListStateUpgradeV0(_ context.Context, rawState map[string]interface{}, _ interface{}) (map[string]interface{}, error) {
	log.Printf("[DEBUG] Starting PolicyPrefixList state migration from V0 to V1 (TypeList -> TypeSet)")

	// Migrate "prefix" from list to set
	if v, ok := rawState["prefix"]; ok {
		if prefixList, ok := v.([]interface{}); ok {
			oldCount := len(prefixList)
			log.Printf("[DEBUG] Migrating 'prefix' field: %d entries from TypeList to TypeSet", oldCount)

			set := schema.NewSet(prefixHash, nil)
			for i, item := range prefixList {
				if m, ok := item.(map[string]interface{}); ok {
					entry := map[string]interface{}{
						"cidr":        m["cidr"],
						"description": m["description"],
					}
					set.Add(entry)
					log.Printf("[DEBUG] Migrated prefix[%d]: cidr=%v, description=%v", i, m["cidr"], m["description"])
				}
			}
			rawState["prefix"] = set.List()
			log.Printf("[DEBUG] 'prefix' migration complete: %d entries migrated to %d set entries", oldCount, set.Len())
		} else {
			log.Printf("[DEBUG] 'prefix' field exists but is not a TypeList, skipping migration (type=%T)", v)
		}
	} else {
		log.Printf("[DEBUG] 'prefix' field not present in state, skipping migration")
	}

	// Migrate "prefix_range" from list to set
	if v, ok := rawState["prefix_range"]; ok {
		if rangeList, ok := v.([]interface{}); ok {
			oldCount := len(rangeList)
			log.Printf("[DEBUG] Migrating 'prefix_range' field: %d entries from TypeList to TypeSet", oldCount)

			set := schema.NewSet(prefixRangeHash, nil)
			for i, item := range rangeList {
				if m, ok := item.(map[string]interface{}); ok {
					le := 0
					ge := 0
					if leVal, ok := m["le"]; ok {
						le = toInt(leVal)
					}
					if geVal, ok := m["ge"]; ok {
						ge = toInt(geVal)
					}
					entry := map[string]interface{}{
						"prefix":      m["prefix"],
						"description": m["description"],
						"le":          le,
						"ge":          ge,
					}
					set.Add(entry)
					log.Printf("[DEBUG] Migrated prefix_range[%d]: prefix=%v, description=%v, le=%d, ge=%d", i, m["prefix"], m["description"], le, ge)
				}
			}
			rawState["prefix_range"] = set.List()
			log.Printf("[DEBUG] 'prefix_range' migration complete: %d entries migrated to %d set entries", oldCount, set.Len())
		} else {
			log.Printf("[DEBUG] 'prefix_range' field exists but is not a TypeList, skipping migration (type=%T)", v)
		}
	} else {
		log.Printf("[DEBUG] 'prefix_range' field not present in state, skipping migration")
	}

	log.Printf("[INFO] PolicyPrefixList state migration from V0 to V1 completed successfully")
	return rawState, nil
}
