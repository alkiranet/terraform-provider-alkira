package alkira

import (
	"context"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// resourceBluecatV0 returns the V0 schema (TypeList for the anycast `ips` and
// `backup_cxps` fields) used before the migration to TypeSet. This is needed by
// the state upgrader to parse old state formats. Only the schema structure
// matters here (it feeds CoreConfigSchema().ImpliedType()), so validators and
// set-hash functions are intentionally omitted.
func resourceBluecatV0() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"bdds_anycast": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"ips": {
							Type:     schema.TypeList,
							Optional: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
						"backup_cxps": {
							Type:     schema.TypeList,
							Optional: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
					},
				},
			},
			"edge_anycast": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"ips": {
							Type:     schema.TypeList,
							Optional: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
						"backup_cxps": {
							Type:     schema.TypeList,
							Optional: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
					},
				},
			},
			"billing_tag_ids": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeInt},
			},
			"cxp": {
				Type:     schema.TypeString,
				Required: true,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"global_cidr_list_id": {
				Type:     schema.TypeInt,
				Required: true,
			},
			"provision_state": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"instance": {
				Type:     schema.TypeSet,
				Required: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"id": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"bdds_options": {
							Type:     schema.TypeList,
							Optional: true,
							MaxItems: 1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"client_id": {
										Type:     schema.TypeString,
										Required: true,
									},
									"activation_key": {
										Type:      schema.TypeString,
										Required:  true,
										Sensitive: true,
									},
									"license_credential_id": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"hostname": {
										Type:     schema.TypeString,
										Required: true,
									},
									"model": {
										Type:     schema.TypeString,
										Required: true,
									},
									"version": {
										Type:     schema.TypeString,
										Required: true,
									},
								},
							},
						},
						"edge_options": {
							Type:     schema.TypeList,
							Optional: true,
							MaxItems: 1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"config_data": {
										Type:     schema.TypeString,
										Required: true,
									},
									"credential_id": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"hostname": {
										Type:     schema.TypeString,
										Required: true,
									},
									"version": {
										Type:     schema.TypeString,
										Required: true,
									},
								},
							},
						},
						"type": {
							Type:     schema.TypeString,
							Required: true,
						},
					},
				},
			},
			"license_type": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"segment_ids": {
				Type:     schema.TypeSet,
				Required: true,
				MinItems: 1,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"service_group_name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"service_group_id": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"service_group_implicit_group_id": {
				Type:     schema.TypeInt,
				Computed: true,
			},
		},
	}
}

// resourceBluecatStateUpgradeV0 migrates the state from V0 (TypeList for the
// anycast `ips` and `backup_cxps` fields) to V1 (TypeSet). The values are plain
// strings, so no type conversion is required: the SDK re-keys the flat-map
// entries from positional indices (e.g. "bdds_anycast.<hash>.ips.0") to
// hash-based keys when it re-encodes the returned state under the V1 schema.
func resourceBluecatStateUpgradeV0(_ context.Context, rawState map[string]interface{}, _ interface{}) (map[string]interface{}, error) {
	log.Printf("[DEBUG] Starting ServiceBluecat state migration from V0 to V1 (TypeList -> TypeSet for anycast ips/backup_cxps)")
	log.Printf("[INFO] ServiceBluecat state migration from V0 to V1 completed successfully")
	return rawState, nil
}
