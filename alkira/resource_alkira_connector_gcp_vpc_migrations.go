package alkira

import (
	"context"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// resourceConnectorGcpVpcV0 returns the V0 schema (TypeList for prefix_list_ids)
// used before the migration to TypeSet. This is needed by the state upgrader to
// parse old state formats.
func resourceConnectorGcpVpcV0() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"name": {
				Description: "Connector name.",
				Type:        schema.TypeString,
				Required:    true,
			},
			"description": {
				Description: "Connector description.",
				Type:        schema.TypeString,
				Optional:    true,
			},
			"cxp": {
				Description: "Alkira CXP where the connector will be instantiated.",
				Type:        schema.TypeString,
				Required:    true,
			},
			"size": {
				Description: "Size of the connector.",
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "SMALL",
			},
			"segment_id": {
				Description: "Segment ID to which this connector belongs.",
				Type:        schema.TypeString,
				Required:    true,
			},
			"provision_state": {
				Description: "Provisioning state of the resource.",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"billing_tag_ids": {
				Description: "List of billing tag IDs to be associated with the connector.",
				Type:        schema.TypeSet,
				Optional:    true,
				Elem:        &schema.Schema{Type: schema.TypeInt},
			},
			"gcp_region": {
				Description: "GCP region where the connector is instantiated.",
				Type:        schema.TypeString,
				Required:    true,
			},
			"gcp_project_id": {
				Description: "GCP project ID.",
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
			},
			"credential_id": {
				Description: "Credential ID to use for GCP connector.",
				Type:        schema.TypeString,
				Required:    true,
			},
			"gcp_vpc_name": {
				Description: "GCP VPC network name.",
				Type:        schema.TypeString,
				Required:    true,
			},
			"internal_id": {
				Description: "Internal ID for the connector.",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"vpc_subnet": {
				Description: "List of subnets to be advertised to the CXP.",
				Type:        schema.TypeSet,
				Optional:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Description: "GCP VPC subnet identifier.",
							Type:        schema.TypeString,
							Optional:    true,
						},
						"internal_id": {
							Description: "Internal ID for the subnet (from UI/import).",
							Type:        schema.TypeString,
							Optional:    true,
						},
						"cidr": {
							Description: "CIDR of the subnet.",
							Type:        schema.TypeString,
							Optional:    true,
						},
					},
				},
			},
			"gcp_routing": {
				Type: schema.TypeList,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"prefix_list_ids": {
							Description: "IDs of prefix lists defined on the " +
								"network.",
							Type:     schema.TypeList,
							Required: true,
							Elem:     &schema.Schema{Type: schema.TypeInt},
						},
						"custom_prefix": {
							Description: "Specifies the source of the routes " +
								"that need to be imported. The value could be " +
								"`ADVERTISE_DEFAULT_ROUTE` and " +
								"`ADVERTISE_CUSTOM_PREFIX`.",
							Type:     schema.TypeString,
							Required: true,
						},
						"export_all_subnets": {
							Description: "Whether to export all subnets to " +
								"CXP. When set to true, all subnets in the VPC " +
								"are advertised to the CXP. When set to false, " +
								"only the subnets specified in vpc_subnet are " +
								"advertised.",
							Type:     schema.TypeBool,
							Optional: true,
							Default:  true,
						},
					},
				},
				Optional: true,
				MaxItems: 1,
			},
		},
	}
}

// resourceConnectorGcpVpcStateUpgradeV0 migrates the state from V0
// (TypeList for prefix_list_ids) to V1 (TypeSet). The flat-map
// state keys change from positional indices (e.g. "gcp_routing.0.prefix_list_ids.0")
// to hash-based keys (e.g. "gcp_routing.0.prefix_list_ids.<hash>").
func resourceConnectorGcpVpcStateUpgradeV0(_ context.Context, rawState map[string]interface{}, _ interface{}) (map[string]interface{}, error) {
	log.Printf("[DEBUG] Starting ConnectorGcpVpc state migration from V0 to V1 (TypeList -> TypeSet for prefix_list_ids)")

	// Migrate "gcp_routing.0.prefix_list_ids" from list to set
	if v, ok := rawState["gcp_routing"]; ok {
		if routingList, ok := v.([]interface{}); ok && len(routingList) > 0 {
			if routing, ok := routingList[0].(map[string]interface{}); ok {
				if v, ok := routing["prefix_list_ids"]; ok {
					if prefixListIdsList, ok := v.([]interface{}); ok {
						oldCount := len(prefixListIdsList)
						log.Printf("[DEBUG] Migrating 'gcp_routing.0.prefix_list_ids' field: %d entries from TypeList to TypeSet", oldCount)

						// State upgraders must return plain JSON-serializable types.
						// The SDK handles TypeSet conversion from []interface{}.
						// JSON numbers are float64; convert to int for TypeInt elements.
						migrated := make([]interface{}, 0, len(prefixListIdsList))
						for i, item := range prefixListIdsList {
							switch v := item.(type) {
							case float64:
								migrated = append(migrated, int(v))
							case int:
								migrated = append(migrated, v)
							default:
								log.Printf("[WARN] Unexpected type for prefix_list_ids[%d]: %T, skipping", i, item)
								continue
							}
							log.Printf("[DEBUG] Migrated prefix_list_ids[%d]: %v", i, item)
						}
						routing["prefix_list_ids"] = migrated
						log.Printf("[DEBUG] 'gcp_routing.0.prefix_list_ids' migration complete: %d entries migrated", oldCount)
					} else {
						log.Printf("[DEBUG] 'gcp_routing.0.prefix_list_ids' field exists but is not a TypeList, skipping migration (type=%T)", v)
					}
				} else {
					log.Printf("[DEBUG] 'gcp_routing.0.prefix_list_ids' field not present in state, skipping migration")
				}
				rawState["gcp_routing"] = []interface{}{routing}
			}
		}
	} else {
		log.Printf("[DEBUG] 'gcp_routing' field not present in state, skipping migration")
	}

	log.Printf("[INFO] ConnectorGcpVpc state migration from V0 to V1 completed successfully")
	return rawState, nil
}
