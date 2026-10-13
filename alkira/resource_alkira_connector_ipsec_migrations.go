package alkira

import (
	"context"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// connectorRoutingOptionsSetUpgraderV0 returns the V0 -> V1 state upgrader
// shared by alkira_connector_ipsec and alkira_connector_ipsec_adv, where V1
// changed routing_options from TypeSet to TypeList (MaxItems 1).
//
// The V0 type is derived from the current schema with routing_options switched
// back to TypeSet, instead of keeping a hand-copied V0 schema. Both types
// serialize to a JSON array, so the upgrade itself needs no transformation.
func connectorRoutingOptionsSetUpgraderV0(r *schema.Resource) schema.StateUpgrader {
	v0 := make(map[string]*schema.Schema, len(r.Schema))
	for k, v := range r.Schema {
		v0[k] = v
	}

	routingOptions := *r.Schema["routing_options"]
	routingOptions.Type = schema.TypeSet
	routingOptions.MaxItems = 0
	v0["routing_options"] = &routingOptions

	return schema.StateUpgrader{
		Version: 0,
		Type:    (&schema.Resource{Schema: v0}).CoreConfigSchema().ImpliedType(),
		Upgrade: connectorRoutingOptionsStateUpgradeV0,
	}
}

func connectorRoutingOptionsStateUpgradeV0(_ context.Context, rawState map[string]interface{}, _ interface{}) (map[string]interface{}, error) {
	log.Printf("[DEBUG] IPSec connector state migration from V0 to V1 (TypeSet -> TypeList for routing_options)")
	return rawState, nil
}
