package alkira

import (
	"log"

	"github.com/alkiranet/alkira-client-go/alkira"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceAlkiraPolicyRule() *schema.Resource {
	return &schema.Resource{
		Create: resourcePolicyRule,
		Read:   resourcePolicyRuleRead,
		Update: resourcePolicyRuleUpdate,
		Delete: resourcePolicyRuleDelete,

		Schema: map[string]*schema.Schema{
			"application_ids": {
				Type:     schema.TypeList,
				Elem:     &schema.Schema{Type: schema.TypeInt},
				Optional: true,
			},
			"application_family_ids": {
				Type:     schema.TypeList,
				Elem:     &schema.Schema{Type: schema.TypeInt},
				Optional: true,
			},
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"src_ip": {
				Type:          schema.TypeString,
				ConflictsWith: []string{"src_prefix_list_id"},
				Optional:      true,
			},
			"dst_ip": {
				Type:          schema.TypeString,
				ConflictsWith: []string{"dst_prefix_list_id"},
				Optional:      true,
			},
			"src_ports": {
				Type:     schema.TypeList,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Optional: true,
			},
			"dst_ports": {
				Type:     schema.TypeList,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Optional: true,
			},
			"src_prefix_list_id": {
				Type:          schema.TypeInt,
				ConflictsWith: []string{"src_ip"},
				Optional:      true,
			},
			"dst_prefix_list_id": {
				Type:          schema.TypeInt,
				ConflictsWith: []string{"dst_ip"},
				Optional:      true,
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
			"rule_action": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "ALLOW",
			},
			"rule_action_service_types": {
				Type:     schema.TypeList,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Optional: true,
			},
		},
	}
}

func resourcePolicyRule(d *schema.ResourceData, m interface{}) error {
	client := m.(*alkira.AlkiraClient)

	request, err := generatePolicyRuleRequest(d, m)

	if err != nil {
		log.Printf("[ERROR] Failed to generate policy rule request")
		return err
	}

	log.Printf("[INFO] Creating policy rule")
	id, err := client.CreatePolicyRule(request)

	if err != nil {
		log.Printf("[ERROR] Failed to create policy rule")
		return err
	}

	d.SetId(id)
	return resourcePolicyRuleRead(d, m)
}

func resourcePolicyRuleRead(d *schema.ResourceData, m interface{}) error {
	client := m.(*alkira.AlkiraClient)

	rule, err := client.GetPolicyRule(d.Id())

	if err != nil {
		log.Printf("[ERROR] Failed to get policy rule %s", d.Id())
		return err
	}

	d.Set("name", rule.Name)
	d.Set("description", rule.Description)

	d.Set("dscp", rule.MatchCondition.Dscp)
	d.Set("Protocol", rule.MatchCondition.Protocol)

	d.Set("src_ip", rule.MatchCondition.SrcIp)
	d.Set("dst_ip", rule.MatchCondition.DstIp)

	d.Set("src_prefix_list_id", rule.MatchCondition.SrcPrefixListId)
	d.Set("dst_prefix_list_id", rule.MatchCondition.DstPrefixListId)

	d.Set("src_ports", rule.MatchCondition.SrcPortList)
	d.Set("dst_ports", rule.MatchCondition.DstPortList)

	d.Set("application_ids", rule.MatchCondition.ApplicationList)
	d.Set("application_family_ids", rule.MatchCondition.ApplicationFamilyList)

	d.Set("rule_action", rule.RuleAction.Action)
	d.Set("rule_action_service_types", rule.RuleAction.ServiceTypeList)

	return nil
}

func resourcePolicyRuleUpdate(d *schema.ResourceData, m interface{}) error {
	client := m.(*alkira.AlkiraClient)

	request, err := generatePolicyRuleRequest(d, m)

	if err != nil {
		log.Printf("[ERROR] Failed to generate policy rule request")
		return err
	}

	err = client.UpdatePolicyRule(d.Id(), request)

	if err != nil {
		log.Printf("[ERROR] Failed to update policy rule %s", d.Id())
		return err
	}

	return resourcePolicyRuleRead(d, m)
}

func resourcePolicyRuleDelete(d *schema.ResourceData, m interface{}) error {
	client := m.(*alkira.AlkiraClient)

	log.Printf("[INFO] Deleting policy rule %s", d.Id())
	return client.DeletePolicyRule(d.Id())
}

func generatePolicyRuleRequest(d *schema.ResourceData, m interface{}) (*alkira.PolicyRule, error) {

	srcPortList := convertTypeListToStringList(d.Get("src_ports").([]interface{}))
	dstPortList := convertTypeListToStringList(d.Get("dst_ports").([]interface{}))

	applicationList := convertTypeListToIntList(d.Get("application_ids").([]interface{}))
	applicationFamilyList := convertTypeListToIntList(d.Get("application_family_ids").([]interface{}))
	serviceTypeList := convertTypeListToStringList(d.Get("rule_action_service_types").([]interface{}))

	request := &alkira.PolicyRule{
		Description: d.Get("description").(string),
		Name:        d.Get("name").(string),
		MatchCondition: alkira.PolicyRuleMatchCondition{
			SrcIp:                 d.Get("src_ip").(string),
			DstIp:                 d.Get("dst_ip").(string),
			Dscp:                  d.Get("dscp").(string),
			Protocol:              d.Get("protocol").(string),
			SrcPortList:           srcPortList,
			DstPortList:           dstPortList,
			SrcPrefixListId:       d.Get("src_prefix_list_id").(int),
			DstPrefixListId:       d.Get("dst_prefix_list_id").(int),
			ApplicationList:       applicationList,
			ApplicationFamilyList: applicationFamilyList,
		},
		RuleAction: alkira.PolicyRuleAction{
			Action:          d.Get("rule_action").(string),
			ServiceTypeList: serviceTypeList,
		},
	}

	return request, nil
}
