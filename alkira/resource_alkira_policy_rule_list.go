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
							Optional: true,
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

	ruleList, err := generatePolicyRuleListRequest(d, m)

	if err != nil {
		log.Printf("[ERROR] failed to generate rule list request")
		return err
	}

	log.Printf("[INFO] Creating policy rule list")
	id, err := client.CreatePolicyRuleList(ruleList)

	if err != nil {
		log.Printf("[ERROR] failed to create policy rule list")
		return err
	}

	d.SetId(id)
	return resourcePolicyRuleListRead(d, m)
}

func resourcePolicyRuleListRead(d *schema.ResourceData, m interface{}) error {
	client := m.(*alkira.AlkiraClient)

	ruleList, err := client.GetPolicyRuleList(d.Id())

	if err != nil {
		log.Printf("[ERROR] failed to get policy rule list %s", d.Id())
		return err
	}

	d.Set("name", ruleList.Name)
	d.Set("description", ruleList.Description)
	d.Set("rules", ruleList.Rules)

	return nil
}

func resourcePolicyRuleListUpdate(d *schema.ResourceData, m interface{}) error {
	client := m.(*alkira.AlkiraClient)

	ruleList, err := generatePolicyRuleListRequest(d, m)

	if err != nil {
		return err
	}

	log.Printf("[INFO] Updateing policy rule list %s", d.Id())
	err = client.UpdatePolicyRuleList(d.Id(), ruleList)

	return resourcePolicyRuleListRead(d, m)
}

func resourcePolicyRuleListDelete(d *schema.ResourceData, m interface{}) error {
	client := m.(*alkira.AlkiraClient)

	log.Printf("[INFO] Deleting policy rule list %s", d.Id())
	return client.DeletePolicyRuleList(d.Id())
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
