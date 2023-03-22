package alkira

import (
	"context"
	"fmt"

	"github.com/alkiranet/alkira-client-go/alkira"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceAlkiraPolicy() *schema.Resource {
	return &schema.Resource{
		Description:   "Manage policy.",
		CreateContext: resourcePolicy,
		ReadContext:   resourcePolicyRead,
		UpdateContext: resourcePolicyUpdate,
		DeleteContext: resourcePolicyDelete,
		CustomizeDiff: func(ctx context.Context, d *schema.ResourceDiff, m interface{}) error {
			client := m.(*alkira.AlkiraClient)

			old, _ := d.GetChange("provision_state")

			if client.Provision == true && old == "FAILED" {
				d.SetNew("provision_state", "SUCCESS")
			}

			return nil
		},
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			"description": {
				Description: "The description of the policy.",
				Type:        schema.TypeString,
				Optional:    true,
			},
			"enabled": {
				Description: "Whether the policy is enabled.",
				Type:        schema.TypeBool,
				Required:    true,
			},
			"from_groups": {
				Description: "IDs of groups that will define source in the policy scope",
				Type:        schema.TypeList,
				Elem:        &schema.Schema{Type: schema.TypeInt},
				Required:    true,
			},
			"name": {
				Description: "The name of the policy.",
				Type:        schema.TypeString,
				Required:    true,
			},
			"provision_state": {
				Description: "The provision state of the resource.",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"rule_list_id": {
				Description: "The `rulelist` that will be used by the policy.",
				Type:        schema.TypeInt,
				Required:    true,
			},
			"segment_ids": {
				Description: "IDs of segments that will define the policy scope.",
				Type:        schema.TypeList,
				Elem:        &schema.Schema{Type: schema.TypeInt},
				Required:    true,
			},
			"to_groups": {
				Description: "IDs of groups that will define destination in the policy scope.",
				Type:        schema.TypeList,
				Elem:        &schema.Schema{Type: schema.TypeInt},
				Required:    true,
			},
		},
	}
}

func resourcePolicy(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {

	client := m.(*alkira.AlkiraClient)
	api := alkira.NewTrafficPolicy(m.(*alkira.AlkiraClient))

	// Construct request
	request := generatePolicyRequest(d, m)

	// Send request
	response, provState, err, provErr := api.Create(request)

	if err != nil {
		return diag.FromErr(err)
	}

	// Set provision state
	if client.Provision == true {
		d.Set("provision_state", provState)

		if provErr != nil {
			return diag.Diagnostics{{
				Severity: diag.Warning,
				Summary:  "PROVISION (CREATE) FAILED",
				Detail:   fmt.Sprintf("%s", provErr),
			}}
		}
	}

	d.SetId(string(response.Id))
	return resourcePolicyRead(ctx, d, m)
}

func resourcePolicyRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {

	client := m.(*alkira.AlkiraClient)
	api := alkira.NewTrafficPolicy(m.(*alkira.AlkiraClient))

	policy, provState, err := api.GetById(d.Id())

	if err != nil {
		return diag.Diagnostics{{
			Severity: diag.Warning,
			Summary:  "FAILED TO GET RESOURCE",
			Detail:   fmt.Sprintf("%s", err),
		}}
	}

	d.Set("description", policy.Description)
	d.Set("enabled", policy.Enabled)
	d.Set("name", policy.Name)
	d.Set("rule_list_id", policy.RuleListId)
	d.Set("segment_ids", policy.SegmentIds)
	d.Set("from_groups", policy.FromGroups)
	d.Set("to_groups", policy.ToGroups)

	// Set provision state
	if client.Provision == true && provState != "" {
		d.Set("provision_state", provState)
	}

	return nil
}

func resourcePolicyUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {

	client := m.(*alkira.AlkiraClient)
	api := alkira.NewTrafficPolicy(m.(*alkira.AlkiraClient))

	// Construct request
	request := generatePolicyRequest(d, m)

	// Send update request
	provState, err, provErr := api.Update(d.Id(), request)

	if err != nil {
		return diag.FromErr(err)
	}

	// Set provision state
	if client.Provision == true {
		d.Set("provision_state", provState)

		if provErr != nil {
			return diag.Diagnostics{{
				Severity: diag.Warning,
				Summary:  "PROVISION (UPDATE) FAILED",
				Detail:   fmt.Sprintf("%s", provErr),
			}}
		}
	}

	return resourcePolicyRead(ctx, d, m)
}

func resourcePolicyDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {

	client := m.(*alkira.AlkiraClient)
	api := alkira.NewTrafficPolicy(m.(*alkira.AlkiraClient))

	provState, err, provErr := api.Delete(d.Id())

	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId("")

	if client.Provision == true && provState != "SUCCESS" {
		return diag.Diagnostics{{
			Severity: diag.Warning,
			Summary:  "PROVISION (DELETE) FAILED",
			Detail:   fmt.Sprintf("%s", provErr),
		}}
	}

	return nil
}

func generatePolicyRequest(d *schema.ResourceData, m interface{}) *alkira.TrafficPolicy {

	policy := &alkira.TrafficPolicy{
		Description: d.Get("description").(string),
		Enabled:     d.Get("enabled").(bool),
		FromGroups:  convertTypeListToIntList(d.Get("from_groups").([]interface{})),
		Name:        d.Get("name").(string),
		RuleListId:  d.Get("rule_list_id").(int),
		SegmentIds:  convertTypeListToIntList(d.Get("segment_ids").([]interface{})),
		ToGroups:    convertTypeListToIntList(d.Get("to_groups").([]interface{})),
	}

	return policy
}
