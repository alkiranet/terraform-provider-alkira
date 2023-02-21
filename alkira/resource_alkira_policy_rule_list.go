package alkira

import (
	"context"
	"fmt"

	"github.com/alkiranet/alkira-client-go/alkira"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceAlkiraPolicyRuleList() *schema.Resource {
	return &schema.Resource{
		Create: resourcePolicyRuleList,
		Read:   resourcePolicyRuleListRead,
		Update: resourcePolicyRuleListUpdate,
		Delete: resourcePolicyRuleListDelete,
		CustomizeDiff: func(ctx context.Context, d *schema.ResourceDiff, m interface{}) error {
			client := m.(*alkira.AlkiraClient)

			old, _ := d.GetChange("provision_state")

			if client.Provision == true && old == "FAILED" {
				d.SetNew("provision_state", "SUCCESS")
			}

			return nil
		},
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
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

func resourcePolicyRuleList(d *schema.ResourceData, m interface{}) error {

	client := m.(*alkira.AlkiraClient)
	api := alkira.NewPolicyRuleList(m.(*alkira.AlkiraClient))

	// Construct request
	request, err := generatePolicyRuleListRequest(d, m)

	if err != nil {
		return err
	}

	// Send create request
	response, provisionState, err := api.Create(request)

	if err != nil {
		return err
	}

	// Set provision state
	if client.Provision == true {
		d.Set("provision_state", provisionState)
	}

	d.SetId(string(response.Id))
	return resourcePolicyRuleListRead(d, m)
}

func resourcePolicyRuleListRead(d *schema.ResourceData, m interface{}) error {

	client := m.(*alkira.AlkiraClient)
	api := alkira.NewPolicyRuleList(m.(*alkira.AlkiraClient))

	ruleList, err := api.GetById(d.Id())

	if err != nil {
		return err
	}

	d.Set("name", ruleList.Name)
	d.Set("description", ruleList.Description)
	d.Set("rules", ruleList.Rules)

	// Set provision state
	_, provisionState, err := api.GetByName(d.Get("name").(string))

	if client.Provision == true && provisionState != "" {
		d.Set("provision_state", provisionState)
	}

	return nil
}

func resourcePolicyRuleListUpdate(d *schema.ResourceData, m interface{}) error {

	client := m.(*alkira.AlkiraClient)
	api := alkira.NewPolicyRuleList(m.(*alkira.AlkiraClient))

	// Construct request
	request, err := generatePolicyRuleListRequest(d, m)

	if err != nil {
		return err
	}

	// Send update request
	provisionState, err := api.Update(d.Id(), request)

	if err != nil {
		return err
	}

	// Set provision state
	if client.Provision == true {
		d.Set("provision_state", provisionState)
	}

	return resourcePolicyRuleListRead(d, m)
}

func resourcePolicyRuleListDelete(d *schema.ResourceData, m interface{}) error {

	client := m.(*alkira.AlkiraClient)
	api := alkira.NewPolicyRuleList(m.(*alkira.AlkiraClient))

	provisionState, err := api.Delete(d.Id())

	if err != nil {
		return err
	}

	if client.Provision == true && provisionState != "SUCCESS" {
		return fmt.Errorf("failed to delete policy_rule_list %s, provision failed", d.Id())
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
