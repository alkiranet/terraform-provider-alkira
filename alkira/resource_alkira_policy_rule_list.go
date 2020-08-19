package alkira

import (
	"log"
	"strconv"

	"github.com/alkiranet/alkira-client-go/alkira"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
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
			"rule_list_id": {
				Type:     schema.TypeInt,
				Computed: true,
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
	ruleList := &alkira.PolicyRuleListRequest{
		Description: d.Get("description").(string),
		Name:        d.Get("name").(string),
		Rules:       rules,
	}

	log.Printf("[INFO] Policy Rule List Creating")
	id, err := client.CreatePolicyRuleList(ruleList)
	log.Printf("[INFO] PolicyRuleList ID: %d", id)

	if err != nil {
		log.Printf("[ERROR] failed to create rule list")
		return err
	}

	d.SetId(strconv.Itoa(id))
	d.Set("rule_list_id", id)

	return resourcePolicyRuleListRead(d, meta)
}

func resourcePolicyRuleListRead(d *schema.ResourceData, meta interface{}) error {
	return nil
}

func resourcePolicyRuleListUpdate(d *schema.ResourceData, meta interface{}) error {
	return resourcePolicyRuleListRead(d, meta)
}

func resourcePolicyRuleListDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*alkira.AlkiraClient)
	id := d.Get("rule_list_id").(int)

	log.Printf("[INFO] Deleting Policy Rule List %d", id)
	err := client.DeletePolicyRuleList(id)

	if err != nil {
		return err
	}

	return nil
}
