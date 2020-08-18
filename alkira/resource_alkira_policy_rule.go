package alkira

import (
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
			"application_list": {
				Type:     schema.TypeList,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Optional: true,
			},
			"application_family_list": {
				Type:     schema.TypeList,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Optional: true,
			},
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The name of the policy rule",
			},
			"description": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The description of the policy rule",

			},
			"src_ip": {
				Type:     schema.TypeString,
				Required: true,
			},
			"dst_ip": {
				Type:     schema.TypeString,
				Required: true,
			},
			"dscp": {
				Type:     schema.TypeString,
				Required: true,
			},
			"internet_application_id": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"protocol": {
				Type:     schema.TypeString,
				Required: true,
			},
			"src_port_list": {
				Type:     schema.TypeList,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Required: true,
			},
			"dst_port_list": {
				Type:     schema.TypeList,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Required: true,
			},
			"rule_id": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"rule_action": {
				Type:     schema.TypeString,
				Optional: true,
			    Default:  "ALLOW",
			},
			"rule_action_service_type_list": {
				Type:     schema.TypeList,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Optional: true,
			},
		},
	}
}

func resourcePolicyRule(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*alkira.AlkiraClient)

	srcPortList := convertTypeListToStringList(d.Get("src_port_list").([]interface{}))
	dstPortList := convertTypeListToStringList(d.Get("dst_port_list").([]interface{}))

	applicationList       := convertTypeListToStringList(d.Get("application_list").([]interface{}))
	applicationFamilyList := convertTypeListToStringList(d.Get("application_family_list").([]interface{}))
	serviceTypeList       := convertTypeListToStringList(d.Get("rule_action_service_type_list").([]interface{}))

	rule   := &alkira.PolicyRuleRequest{
		Description: d.Get("description").(string),
		Name:        d.Get("name").(string),
		MatchCondition: alkira.PolicyRuleMatchCondition{
			SrcIp:                 d.Get("src_ip").(string),
			DstIp:                 d.Get("dst_ip").(string),
			Dscp:                  d.Get("dscp").(string),
			InternetApplicationId: d.Get("internet_application_id").(int),
			Protocol:              d.Get("protocol").(string),
			SrcPortList:           srcPortList,
			DstPortList:           dstPortList,
			ApplicationList:       applicationList,
			ApplicationFamilyList: applicationFamilyList,
		},
		RuleAction: alkira.PolicyRuleAction{
			Action:          d.Get("rule_action").(string),
			ServiceTypeList: serviceTypeList,
		},
	}

	id, err := client.CreatePolicyRule(rule)
	log.Printf("[INFO] Policy Rule ID: %d", id)

	if err != nil {
		log.Printf("[ERROR] Failed to create policy rule")
		return err
	}

	d.SetId(strconv.Itoa(id))
	d.Set("rule_id", id)

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
	PolicyRuleId := d.Get("rule_id").(int)

	log.Printf("[INFO] Deleting PolicyRule %s", PolicyRuleId)
	err := client.DeletePolicyRule(PolicyRuleId)

	if err != nil {
	 	return err
	}

	return nil
}
