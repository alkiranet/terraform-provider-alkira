package alkira

import (
	"context"
	"fmt"

	"github.com/alkiranet/alkira-client-go/alkira"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceAlkiraPolicyRuleList() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourcePolicyRuleList,
		ReadContext:   resourcePolicyRuleListRead,
		UpdateContext: resourcePolicyRuleListUpdate,
		DeleteContext: resourcePolicyRuleListDelete,
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
			"name": {
				Description: "The name of the policy rule list.",
				Type:        schema.TypeString,
				Required:    true,
			},
			"description": {
				Description: "The description of the policy rule list.",
				Type:        schema.TypeString,
				Optional:    true,
			},
			"provision_state": {
				Description: "The provisioning state of the resource.",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"rules": {
				Type: schema.TypeSet,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"priority": {
							Type:     schema.TypeInt,
							Required: true,
						},
						"rule_id": {
							Type:     schema.TypeInt,
							Required: true,
						},
					},
				},
				Required: true,
			},
		},
	}
}

func resourcePolicyRuleList(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {

	client := m.(*alkira.AlkiraClient)
	api := alkira.NewPolicyRuleList(m.(*alkira.AlkiraClient))

	// Construct request
	request, err := generatePolicyRuleListRequest(d, m)

	if err != nil {
		return diag.FromErr(err)
	}

	// Send create request
	response, provState, err, provErr := api.Create(request)

	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(string(response.Id))

	// Set provision state
	if client.Provision == true {
		d.Set("provision_state", provState)

		if provErr != nil {
			return diag.Diagnostics{{
				Severity: diag.Warning,
				Summary:  "PROVISION FAILED",
				Detail:   fmt.Sprintf("%s", provErr),
			}}
		}
	}

	return resourcePolicyRuleListRead(ctx, d, m)
}

func resourcePolicyRuleListRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {

	client := m.(*alkira.AlkiraClient)
	api := alkira.NewPolicyRuleList(m.(*alkira.AlkiraClient))

	ruleList, provState, err := api.GetById(d.Id())

	if err != nil {
		return diag.FromErr(err)
	}

	d.Set("name", ruleList.Name)
	d.Set("description", ruleList.Description)
	d.Set("rules", ruleList.Rules)

	// Set provision state
	if client.Provision == true && provState != "" {
		d.Set("provision_state", provState)
	}

	return nil
}

func resourcePolicyRuleListUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {

	client := m.(*alkira.AlkiraClient)
	api := alkira.NewPolicyRuleList(m.(*alkira.AlkiraClient))

	// Construct request
	request, err := generatePolicyRuleListRequest(d, m)

	if err != nil {
		return diag.FromErr(err)
	}

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
				Summary:  "PROVISION FAILED",
				Detail:   fmt.Sprintf("%s", provErr),
			}}
		}
	}

	return resourcePolicyRuleListRead(ctx, d, m)
}

func resourcePolicyRuleListDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {

	client := m.(*alkira.AlkiraClient)
	api := alkira.NewPolicyRuleList(m.(*alkira.AlkiraClient))

	provState, err, provErr := api.Delete(d.Id())

	if err != nil {
		return diag.FromErr(err)
	}

	if client.Provision == true && provState != "SUCCESS" {
		return diag.FromErr(provErr)
	}

	d.SetId("")
	return nil
}

func generatePolicyRuleListRequest(d *schema.ResourceData, m interface{}) (*alkira.PolicyRuleList, error) {

	rules := expandPolicyRuleListRules(d.Get("rules").(*schema.Set))

	request := &alkira.PolicyRuleList{
		Name:        d.Get("name").(string),
		Description: d.Get("description").(string),
		Rules:       rules,
	}

	return request, nil
}
