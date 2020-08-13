package alkira

import (
	"fmt"
	"log"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/alkiranet/alkira-client-go/alkira"
)

func resourceAlkiraPolicyRuleList() *schema.Resource {
	return &schema.Resource{
		Create: resourcePolicyRuleList,
		Read:   resourcePolicyRuleListRead,
		Update: resourcePolicyRuleListUpdate,
		Delete: resourcePolicyRuleListDelete,

		Schema: map[string]*schema.Schema{
			"name": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				Description: "The name of the policy rule list",
			},
			"description": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The description of the policy rule list",

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
							Type:      schema.TypeInt,
							Optional:  true,
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

	rules    := expandPolicyRuleListRules(d.Get("rules").(*schema.Set))
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
	}

	d.SetId(strconv.Itoa(id))
	return resourcePolicyRuleListRead(d, meta)
}

func resourcePolicyRuleListRead(d *schema.ResourceData, meta interface{}) error {
        return nil
}

func resourcePolicyRuleListUpdate(d *schema.ResourceData, meta interface{}) error {
        return resourcePolicyRuleListRead(d, meta)
}

func resourcePolicyRuleListDelete(d *schema.ResourceData, meta interface{}) error {
	client    := meta.(*alkira.AlkiraClient)
	id        := d.Id()

	log.Printf("[INFO] Deleting Policy Rule List %s", id)
	err := client.DeletePolicyRuleList(id)

	if err != nil {
	 	return fmt.Errorf("failed to delete Policy Rule List %s", id)
	}

	return nil
}
