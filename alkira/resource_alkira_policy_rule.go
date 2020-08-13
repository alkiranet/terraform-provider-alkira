package alkira

import (
	"fmt"
	"log"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/alkiranet/alkira-client-go/alkira"
)

func resourceAlkiraPolicyRule() *schema.Resource {
	return &schema.Resource{
		Create: resourcePolicyRule,
		Read:   resourcePolicyRuleRead,
		Update: resourcePolicyRuleUpdate,
		Delete: resourcePolicyRuleDelete,

		Schema: map[string]*schema.Schema{
			"name": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				Description: "The name of the policy rule",
			},
			"description": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The description of the policy rule",

			},
			"src_ip": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
			},
			"dst_ip": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
			},
			"dscp": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
			},
			"protocol": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
			},
			"src_port_list": &schema.Schema{
				Type:        schema.TypeList,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Required:    true,
			},
			"dst_port_list": &schema.Schema{
				Type:        schema.TypeList,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Required:    true,
			},
			"action": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
			},
		},
	}
}

func resourcePolicyRule(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*alkira.AlkiraClient)

	srcPortList := expandStringFromList(d.Get("src_port_list").([]interface{}))
	dstPortList := expandStringFromList(d.Get("dst_port_list").([]interface{}))

	rule   := &alkira.PolicyRuleRequest{
		Description: d.Get("description").(string),
		Name:        d.Get("name").(string),
		MatchCondition: alkira.PolicyRuleMatchCondition{
			SrcIp:       d.Get("src_ip").(string),
			DstIp:       d.Get("dst_ip").(string),
			Dscp:        d.Get("dscp").(string),
			Protocol:    d.Get("protocol").(string),
			SrcPortList: srcPortList,
			DstPortList: dstPortList,
		},
		RuleAction: alkira.PolicyRuleAction{
			Action: d.Get("action").(string),
		},
	}

	id, err := client.CreatePolicyRule(rule)
	log.Printf("[INFO] Policy Rule ID: %d", id)

	if err != nil {
		log.Printf("[ERROR] Failed to create policy rule")
	}

	d.SetId(strconv.Itoa(id))
	return resourcePolicyRuleRead(d, meta)
}

func resourcePolicyRuleRead(d *schema.ResourceData, meta interface{}) error {
        return nil
}

func resourcePolicyRuleUpdate(d *schema.ResourceData, meta interface{}) error {
        return resourcePolicyRuleRead(d, meta)
}

func resourcePolicyRuleDelete(d *schema.ResourceData, meta interface{}) error {
	client    := meta.(*alkira.AlkiraClient)
	PolicyRuleId := d.Id()

	log.Printf("[INFO] Deleting PolicyRule %s", PolicyRuleId)
	err := client.DeletePolicyRule(PolicyRuleId)

	if err != nil {
	 	return fmt.Errorf("failed to delete PolicyRule %s", PolicyRuleId)
	}

	return nil
}
