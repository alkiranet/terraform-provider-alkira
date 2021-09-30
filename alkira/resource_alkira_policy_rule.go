package alkira

import (
	"log"
	"strings"

	"github.com/alkiranet/alkira-client-go/alkira"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func resourceAlkiraPolicyRule() *schema.Resource {
	return &schema.Resource{
		Description: "Manage policy rule.\n\n" +
			"This resource is usually used along with policy resources:" +
			"`policy_prefix_list`, `policy_rule_list` and `policy`" +
			"control the network traffic.",
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
				Description: "The name of the policy rule.",
				Type:        schema.TypeString,
				Required:    true,
			},
			"description": {
				Description: "The description of the policy rule.",
				Type:        schema.TypeString,
				Optional:    true,
			},
			"src_ip": {
				Description:   "A single source IP as The match condition of the rule.",
				Type:          schema.TypeString,
				ConflictsWith: []string{"src_prefix_list_id"},
				Optional:      true,
			},
			"dst_ip": {
				Description:   "A single destination IP as The match condition of the rule.",
				Type:          schema.TypeString,
				ConflictsWith: []string{"dst_prefix_list_id"},
				Optional:      true,
			},
			"src_ports": {
				Description: "Source ports that can take values: `any` or `1` to `65535`.",
				Type:        schema.TypeList,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Optional:    true,
			},
			"dst_ports": {
				Description: "Destination ports that can take values: `any` or `1` to `65535`.",
				Type:        schema.TypeList,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Optional:    true,
			},
			"src_prefix_list_id": {
				Description:   "The ID of prefix list as source associated with the rule.",
				Type:          schema.TypeInt,
				ConflictsWith: []string{"src_ip"},
				Optional:      true,
			},
			"dst_prefix_list_id": {
				Description:   "The ID of prefix list as destination associated with the rule.",
				Type:          schema.TypeInt,
				ConflictsWith: []string{"dst_ip"},
				Optional:      true,
			},
			"dscp": {
				Description: "The dscp value can be `any` or between `0` to `63` inclusive.",
				Type:        schema.TypeString,
				Required:    true,
			},
			"internet_application_id": {
				Description: "The ID of the internet application associated with the rule. When an internet applciation is selected, destination ip and port will be the private ip and port of the application.",
				Type:        schema.TypeInt,
				Optional:    true,
			},
			"protocol": {
				Description:  "The following protocols are supported, `ICMP`, `TCP`, `UDP` or `ANY`.",
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringInSlice([]string{"ICMP", "TCP", "UDP", "ANY"}, false),
			},
			"rule_action": {
				Description:  "The action that is applied on matched traffic, either `ALLOW` or `DROP`. The default value is `ALLOW`.",
				Type:         schema.TypeString,
				Optional:     true,
				Default:      "ALLOW",
				ValidateFunc: validation.StringInSlice([]string{"ALLOW", "DROP"}, false),
			},
			"rule_action_service_types": {
				Description: "Based on the service type, traffic is routed to service of the given type.",
				Type:        schema.TypeList,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Optional:    true,
			},
			"rule_action_service_ids": {
				Description: "Based on the service IDs, traffic is routed to the specified services.",
				Type:        schema.TypeList,
				Elem:        &schema.Schema{Type: schema.TypeInt},
				Optional:    true,
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
	d.Set("Protocol", strings.ToUpper(rule.MatchCondition.Protocol))

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
	d.Set("rule_action_service_ids", rule.RuleAction.ServiceList)

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

	return client.DeletePolicyRule(d.Id())
}

func generatePolicyRuleRequest(d *schema.ResourceData, m interface{}) (*alkira.PolicyRule, error) {

	srcPortList := convertTypeListToStringList(d.Get("src_ports").([]interface{}))
	dstPortList := convertTypeListToStringList(d.Get("dst_ports").([]interface{}))

	applicationList := convertTypeListToIntList(d.Get("application_ids").([]interface{}))
	applicationFamilyList := convertTypeListToIntList(d.Get("application_family_ids").([]interface{}))
	serviceTypeList := convertTypeListToStringList(d.Get("rule_action_service_types").([]interface{}))
	serviceList := convertTypeListToIntList(d.Get("rule_action_service_ids").([]interface{}))

	request := &alkira.PolicyRule{
		Description: d.Get("description").(string),
		Name:        d.Get("name").(string),
		MatchCondition: alkira.PolicyRuleMatchCondition{
			SrcIp:                 d.Get("src_ip").(string),
			DstIp:                 d.Get("dst_ip").(string),
			Dscp:                  d.Get("dscp").(string),
			Protocol:              strings.ToLower(d.Get("protocol").(string)),
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
			ServiceList:     serviceList,
		},
	}

	return request, nil
}
