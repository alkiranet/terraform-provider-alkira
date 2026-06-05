package alkira

import (
	"context"
	"fmt"

	"github.com/alkiranet/alkira-client-go/alkira"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func resourceAlkiraPolicyInterCxpRouting() *schema.Resource {
	return &schema.Resource{
		Description: "Manage Inter-CXP Routing Policy.\n\n" +
			"Configure an inter-CXP route policy to control route redistribution " +
			"between CXP pairs within a segment.",
		CreateContext: resourcePolicyInterCxpRouting,
		ReadContext:   resourcePolicyInterCxpRoutingRead,
		UpdateContext: warnOnFailedStateUpdate(resourcePolicyInterCxpRoutingUpdate),
		DeleteContext: resourcePolicyInterCxpRoutingDelete,
		CustomizeDiff: func(ctx context.Context, d *schema.ResourceDiff, m interface{}) error {
			client := m.(*alkira.AlkiraClient)

			old, _ := d.GetChange("provision_state")

			if client.Provision && old == "FAILED" {
				d.SetNew("provision_state", "SUCCESS")
			}

			// direction is immutable after creation.
			if d.Id() != "" && d.HasChange("direction") {
				return fmt.Errorf("direction cannot be changed after creation")
			}

			// Validate rule-level constraints.
			rules := d.Get("rule").([]interface{})
			if err := validateInterCxpRoutingRules(rules); err != nil {
				return err
			}

			return nil
		},
		Importer: &schema.ResourceImporter{
			StateContext: importWithReadValidation(resourcePolicyInterCxpRoutingRead),
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Description: "The name of the inter-CXP routing policy. Must be unique within the tenant network.",
				Type:        schema.TypeString,
				Required:    true,
			},
			"description": {
				Description: "The description of the inter-CXP routing policy.",
				Type:        schema.TypeString,
				Optional:    true,
			},
			"provision_state": {
				Description: "The provision state of the resource.",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"enabled": {
				Description: "Whether the inter-CXP routing policy is enabled.",
				Type:        schema.TypeBool,
				Required:    true,
			},
			"direction": {
				Description: "The direction of the policy. Only `OUTBOUND` is supported in Phase 1. " +
					"Immutable after creation.",
				Type:     schema.TypeString,
				Required: true,
				ValidateFunc: validation.StringInSlice(
					[]string{"OUTBOUND"}, false),
			},
			"segment_id": {
				Description: "ID of the segment that defines the policy scope. Both source and " +
					"destination CXPs must carry this segment.",
				Type:     schema.TypeString,
				Required: true,
			},
			"source_cxps": {
				Description: "List of source CXP names from which routes are redistributed. " +
					"Exactly one CXP is allowed. The CXP must carry the policy segment.",
				Type:     schema.TypeList,
				Required: true,
				MinItems: 1,
				MaxItems: 1,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"dest_cxps": {
				Description: "Set of destination CXP names to which routes are redistributed. " +
					"Each CXP must carry the policy segment. A source CXP cannot also be a destination.",
				Type:     schema.TypeSet,
				Required: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"rule": {
				Type:     schema.TypeList,
				MinItems: 1,
				Required: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Description: "The name of the rule. Must be unique within the policy.",
							Type:        schema.TypeString,
							Required:    true,
						},
						"sequence_no": {
							Description: "System-assigned sequence number starting at `1000`. " +
								"Defines rule evaluation order (top-down, first match wins).",
							Type:     schema.TypeInt,
							Computed: true,
						},
						"action": {
							Description: "Action for matched routes. `ALLOW` permits redistribution " +
								"(with optional set operations). `DENY` blocks redistribution.",
							Type:     schema.TypeString,
							Required: true,
							ValidateFunc: validation.StringInSlice([]string{
								"ALLOW", "DENY"}, false),
						},
						"match_all": {
							Description: "Match all routes. Exclusive — cannot be combined with any other match condition.",
							Type:        schema.TypeBool,
							Optional:    true,
						},
						"match_prefix_list_ids": {
							Description: "IDs of Prefix Lists to match.",
							Type:        schema.TypeSet,
							MaxItems:    1,
							Elem:        &schema.Schema{Type: schema.TypeInt},
							Optional:    true,
						},
						"match_community_list_ids": {
							Description: "IDs of Community Lists to match.",
							Type:        schema.TypeSet,
							MaxItems:    1,
							Elem:        &schema.Schema{Type: schema.TypeInt},
							Optional:    true,
						},
						"match_extended_community_list_ids": {
							Description: "IDs of Extended Community Lists to match. " +
								"Mutually exclusive with `match_group_ids`.",
							Type:     schema.TypeSet,
							MaxItems: 1,
							Elem:     &schema.Schema{Type: schema.TypeInt},
							Optional: true,
						},
						"match_as_path_list_ids": {
							Description: "IDs of AS Path Lists to match.",
							Type:        schema.TypeSet,
							MaxItems:    1,
							Elem:        &schema.Schema{Type: schema.TypeInt},
							Optional:    true,
						},
						"match_segment_resource_ids": {
							Description: "IDs of segment resources to match. Each resource must be shared " +
								"with the policy segment. Mutually exclusive with `match_group_ids`.",
							Type:     schema.TypeSet,
							Elem:     &schema.Schema{Type: schema.TypeInt},
							Optional: true,
						},
						"match_group_ids": {
							Description: "IDs of connector groups to match. Service groups are not allowed. " +
								"Mutually exclusive with `match_segment_resource_ids` and `match_extended_community_list_ids`.",
							Type:     schema.TypeSet,
							Elem:     &schema.Schema{Type: schema.TypeInt},
							Optional: true,
						},
						"set_as_path_prepend": {
							Description: "Prepend one or more AS numbers to the AS PATH. " +
								"Space-separated values 0–65535. Example: `100 100 100`. " +
								"Valid only when action is `ALLOW`.",
							Type:     schema.TypeString,
							Optional: true,
						},
						"set_community": {
							Description: "Add BGP community attributes. Format: `as-number:community-value` " +
								"(values 0–65535). Example: `65512:20 65512:21`. " +
								"Valid only when action is `ALLOW`.",
							Type:     schema.TypeString,
							Optional: true,
						},
						"set_extended_community": {
							Description: "Add BGP extended community attributes. " +
								"Format: `type:administrator:assigned-number`. Only SOO/RT types supported. " +
								"Example: `soo:65512:21 soo:10.1.1.1:1234`. " +
								"Valid only when action is `ALLOW`.",
							Type:     schema.TypeString,
							Optional: true,
						},
					},
				},
			}, // rule
		},
	}
}

func resourcePolicyInterCxpRouting(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {

	client := m.(*alkira.AlkiraClient)
	api := alkira.NewInterCxpRoutePolicy(m.(*alkira.AlkiraClient))

	request, err := generatePolicyInterCxpRoutingRequest(d, m)

	if err != nil {
		return diag.FromErr(err)
	}

	response, provState, err, valErr, provErr := api.Create(request)

	if err != nil {
		return diag.FromErr(err)
	}

	if client.Provision {
		d.Set("provision_state", provState)

		if provErr != nil {
			return diag.Diagnostics{{
				Severity: diag.Warning,
				Summary:  "PROVISION (CREATE) FAILED",
				Detail:   fmt.Sprintf("%s", provErr),
			}}
		}
	}

	if client.Validate && valErr != nil {
		var diags diag.Diagnostics
		readDiags := resourcePolicyInterCxpRoutingRead(ctx, d, m)
		if readDiags.HasError() {
			diags = append(diags, readDiags...)
		}
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "VALIDATION (CREATE) FAILED",
			Detail:   fmt.Sprintf("%s", valErr),
		})
		return diags
	}

	d.SetId(string(response.Id))
	return resourcePolicyInterCxpRoutingRead(ctx, d, m)
}

func resourcePolicyInterCxpRoutingRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {

	client := m.(*alkira.AlkiraClient)
	api := alkira.NewInterCxpRoutePolicy(m.(*alkira.AlkiraClient))

	policy, provState, err := api.GetById(d.Id())

	if err != nil {
		return diag.Diagnostics{{
			Severity: diag.Warning,
			Summary:  "FAILED TO GET RESOURCE",
			Detail:   fmt.Sprintf("%s", err),
		}}
	}

	// Only overwrite description if the API returned a value; the API
	// may omit or null-out the field on read even though it was sent on
	// create, which would cause a perpetual plan diff (BUG-1).
	if policy.Description != "" {
		d.Set("description", policy.Description)
	}
	d.Set("direction", policy.Direction)
	d.Set("enabled", policy.Enabled)
	d.Set("name", policy.Name)
	d.Set("source_cxps", policy.SourceCxps)
	d.Set("dest_cxps", policy.DestCxps)

	segmentId, err := getSegmentIdByName(policy.Segment, m)

	if err != nil {
		return diag.FromErr(err)
	}
	d.Set("segment_id", segmentId)

	err = setPolicyInterCxpRoutingRules(policy.Rules, d)

	if err != nil {
		return diag.FromErr(err)
	}

	if client.Provision && provState != "" {
		d.Set("provision_state", provState)
	}

	return nil
}

func resourcePolicyInterCxpRoutingUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {

	client := m.(*alkira.AlkiraClient)
	api := alkira.NewInterCxpRoutePolicy(m.(*alkira.AlkiraClient))

	request, err := generatePolicyInterCxpRoutingRequest(d, m)

	if err != nil {
		return diag.FromErr(err)
	}

	provState, err, valErr, provErr := api.Update(d.Id(), request)

	if err != nil {
		return diag.FromErr(err)
	}

	if client.Provision {
		d.Set("provision_state", provState)

		if provErr != nil {
			return diag.Diagnostics{{
				Severity: diag.Warning,
				Summary:  "PROVISION (UPDATE) FAILED",
				Detail:   fmt.Sprintf("%s", provErr),
			}}
		}
	}

	if client.Validate && valErr != nil {
		var diags diag.Diagnostics
		readDiags := resourcePolicyInterCxpRoutingRead(ctx, d, m)
		if readDiags.HasError() {
			diags = append(diags, readDiags...)
		}
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "VALIDATION (UPDATE) FAILED",
			Detail:   fmt.Sprintf("%s", valErr),
		})
		return diags
	}

	return resourcePolicyInterCxpRoutingRead(ctx, d, m)
}

func resourcePolicyInterCxpRoutingDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {

	client := m.(*alkira.AlkiraClient)
	api := alkira.NewInterCxpRoutePolicy(m.(*alkira.AlkiraClient))

	_, err, valErr, provErr := api.Delete(d.Id())

	if err != nil {
		name, _ := d.GetOk("name")
		if nameStr, ok := name.(string); ok && nameStr != "" {
			return diag.FromErr(fmt.Errorf("%w alkira_policy_inter_cxp_routing (name=%q id=%s)", err, nameStr, d.Id()))
		}
		return diag.FromErr(fmt.Errorf("%w alkira_policy_inter_cxp_routing (id=%s)", err, d.Id()))
	}

	if client.Validate && valErr != nil {
		return diag.Diagnostics{{
			Severity: diag.Error,
			Summary:  "VALIDATION (DELETE) FAILED",
			Detail:   fmt.Sprintf("%s", valErr),
		}}
	}

	if client.Provision && provErr != nil {
		return diag.Diagnostics{{
			Severity: diag.Warning,
			Summary:  "PROVISION (DELETE) FAILED",
			Detail:   fmt.Sprintf("%s", provErr),
		}}
	}

	d.SetId("")
	return nil
}

func generatePolicyInterCxpRoutingRequest(d *schema.ResourceData, m interface{}) (*alkira.InterCxpRoutePolicy, error) {

	segmentName, err := getSegmentNameById(d.Get("segment_id").(string), m)

	if err != nil {
		return nil, err
	}

	rules, err := expandPolicyInterCxpRoutingRule(d.Get("rule").([]interface{}))

	if err != nil {
		return nil, err
	}

	policy := &alkira.InterCxpRoutePolicy{
		Name:        d.Get("name").(string),
		Description: d.Get("description").(string),
		Direction:   d.Get("direction").(string),
		Enabled:     d.Get("enabled").(bool),
		Segment:     segmentName,
		SourceCxps:  convertTypeListToStringList(d.Get("source_cxps").([]interface{})),
		DestCxps:    convertTypeSetToStringList(d.Get("dest_cxps").(*schema.Set)),
		Rules:       rules,
	}

	return policy, nil
}
