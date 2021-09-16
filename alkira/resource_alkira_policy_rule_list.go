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

func resourcePolicyRuleList(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*alkira.AlkiraClient)

	rules := expandPolicyRuleListRules(d.Get("rules").(*schema.Set))
	ruleList := &alkira.PolicyRuleList{
		Description: d.Get("description").(string),
		Name:        d.Get("name").(string),
		Rules:       rules,
	}

	log.Printf("[INFO] Policy Rule List Creating")
	id, err := client.CreatePolicyRuleList(ruleList)

	if err != nil {
		log.Printf("[ERROR] failed to create rule list")
		return err
	}

	d.SetId(id)
	return resourcePolicyRuleListRead(d, meta)
}

func resourcePolicyRuleListRead(d *schema.ResourceData, meta interface{}) error {
	return nil
}

func resourcePolicyRuleListUpdate(d *schema.ResourceData, meta interface{}) error {
	client := m.(*alkira.AlkiraClient)

	ruleList, err := generatePolicyRuleListRequest(d, m)

	if err != nil {
		return err
	}

	log.Printf("[INFO] Updateing Policy Rule List %s", d.Id())
	err = client.UpdatePolicyRuleList(d.Id(), ruleList)

	return resourcePolicyRuleListRead(d, meta)
}

func resourcePolicyRuleListDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*alkira.AlkiraClient)

	log.Printf("[INFO] Deleting Policy Rule List %s", d.Id())
	return client.DeletePolicyRuleList(d.Id())
}

func generatePolicyRuleListRequest(d *schema.ResourceData, m interface{}) (*alkira.PolicyRuleList, error) {

	rules := expandPolicyRuleListRules(d.Get("rules").(*schema.Set))
	request := &alkira.PolicyRuleListRequest{
		Description: d.Get("description").(string),
		Name:        d.Get("name").(string),
		Rules:       rules,
	}

	return request, nil
}
