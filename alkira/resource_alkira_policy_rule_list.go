package alkira

import (
	"log"

	"github.com/alkiranet/alkira-client-go/alkira"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceAlkiraPolicyRuleList() *schema.Resource {
	return &schema.Resource{
		Create: resourcePolicyRuleList,
		Read:   resourcePolicyRuleListRead,
		Update: resourcePolicyRuleListUpdate,
		Delete: resourcePolicyRuleListDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
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

	api := alkira.NewPolicyRuleList(m.(*alkira.AlkiraClient))

	ruleList, err := generatePolicyRuleListRequest(d, m)

	if err != nil {
		return err
	}

	response, _, err := api.Create(ruleList)

	if err != nil {
		return err
	}

	d.SetId(string(response.Id))
	return resourcePolicyRuleListRead(d, m)
}

func resourcePolicyRuleListRead(d *schema.ResourceData, m interface{}) error {

	api := alkira.NewPolicyRuleList(m.(*alkira.AlkiraClient))

	ruleList, err := api.GetById(d.Id())

	if err != nil {
		return err
	}

	d.Set("name", ruleList.Name)
	d.Set("description", ruleList.Description)
	d.Set("rules", ruleList.Rules)

	return nil
}

func resourcePolicyRuleListUpdate(d *schema.ResourceData, m interface{}) error {

	api := alkira.NewPolicyRuleList(m.(*alkira.AlkiraClient))

	// Construct request
	ruleList, err := generatePolicyRuleListRequest(d, m)

	if err != nil {
		return err
	}

	// Send update request
	_, err = api.Update(d.Id(), ruleList)

	if err != nil {
		return err
	}

	return resourcePolicyRuleListRead(d, m)
}

func resourcePolicyRuleListDelete(d *schema.ResourceData, m interface{}) error {

	api := alkira.NewPolicyRuleList(m.(*alkira.AlkiraClient))

	_, err := api.Delete(d.Id())

	if err != nil {
		return err
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
