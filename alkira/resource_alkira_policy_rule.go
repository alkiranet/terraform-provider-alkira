package alkira

import (
	"context"
	"fmt"

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
			"provision_state": {
				Description: "The provision state of the resource.",
				Type:        schema.TypeString,
				Computed:    true,
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
				Description: "The ID of prefix list as source associated " +
					"with the rule.",
				Type:          schema.TypeInt,
				ConflictsWith: []string{"src_ip"},
				Optional:      true,
			},
			"dst_prefix_list_id": {
				Description: "The ID of prefix list as destination " +
					"associated with the rule.",
				Type:          schema.TypeInt,
				ConflictsWith: []string{"dst_ip"},
				Optional:      true,
			},
			"dscp": {
				Description: "The dscp value can be `any` or between " +
					"`0` to `63` inclusive.",
				Type:     schema.TypeString,
				Required: true,
			},
			"internet_application_id": {
				Description: "The ID of the `internet_application` associated with the " +
					"rule. When an internet applciation is selected, destination IP " +
					"and port will be the private IP and port of the application.",
				Type:     schema.TypeInt,
				Optional: true,
			},
			"protocol": {
				Description:  "The following protocols are supported, `icmp`, `tcp`, `udp` or `any`.",
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringInSlice([]string{"icmp", "tcp", "udp", "any"}, false),
			},
			"rule_action": {
				Description: "The action that is applied on matched traffic, " +
					"either `ALLOW` or `DROP`. The default value is `ALLOW`.",
				Type:         schema.TypeString,
				Optional:     true,
				Default:      "ALLOW",
				ValidateFunc: validation.StringInSlice([]string{"ALLOW", "DROP"}, false),
			},
			"rule_action_service_types": {
				Description: "Based on the service type, traffic is routed to service " +
					"of the given type. For service chaining, both PAN and ZIA service " +
					"types can be selected but must follow order.",
				Type:     schema.TypeList,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Optional: true,
			},
			"rule_action_service_ids": {
				Description: "Based on the service IDs, traffic is routed to the " +
					"specified services. For service chaining, both `service_pan` " +
					"and `service_zscaler`'s IDs can be added here, but ID of " +
					"`service_pan` must be by followed by ID of `service_zscaler`.",
				Type:     schema.TypeList,
				Elem:     &schema.Schema{Type: schema.TypeInt},
				Optional: true,
			},
		},
	}
}

func resourcePolicyRule(d *schema.ResourceData, m interface{}) error {

	client := m.(*alkira.AlkiraClient)
	api := alkira.NewTrafficPolicyRule(m.(*alkira.AlkiraClient))

	// Construct request
	request := generatePolicyRuleRequest(d, m)

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
	return resourcePolicyRuleRead(d, m)
}

func resourcePolicyRuleRead(d *schema.ResourceData, m interface{}) error {

	client := m.(*alkira.AlkiraClient)
	api := alkira.NewTrafficPolicyRule(m.(*alkira.AlkiraClient))

	rule, err := api.GetById(d.Id())

	if err != nil {
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

	d.Set("internet_application_id", rule.MatchCondition.InternetApplicationId)

	d.Set("rule_action", rule.RuleAction.Action)
	d.Set("rule_action_service_types", rule.RuleAction.ServiceTypeList)
	d.Set("rule_action_service_ids", rule.RuleAction.ServiceList)

	// Set provision state
	_, provisionState, err := api.GetByName(d.Get("name").(string))

	if client.Provision == true && provisionState != "" {
		d.Set("provision_state", provisionState)
	}

	return nil
}

func resourcePolicyRuleUpdate(d *schema.ResourceData, m interface{}) error {

	client := m.(*alkira.AlkiraClient)
	api := alkira.NewTrafficPolicyRule(m.(*alkira.AlkiraClient))

	// Construct request
	request := generatePolicyRuleRequest(d, m)

	// Send update request
	provisionState, err := api.Update(d.Id(), request)

	if err != nil {
		return err
	}

	// Set provision state
	if client.Provision == true {
		d.Set("provision_state", provisionState)
	}

	return resourcePolicyRuleRead(d, m)
}

func resourcePolicyRuleDelete(d *schema.ResourceData, m interface{}) error {

	client := m.(*alkira.AlkiraClient)
	api := alkira.NewTrafficPolicyRule(m.(*alkira.AlkiraClient))

	provisionState, err := api.Delete(d.Id())

	if err != nil {
		return err
	}

	if client.Provision == true && provisionState != "SUCCESS" {
		return fmt.Errorf("failed to delete policy_rule %s, provision failed", d.Id())
	}

	d.SetId("")
	return nil
}

func generatePolicyRuleRequest(d *schema.ResourceData, m interface{}) *alkira.TrafficPolicyRule {

	request := &alkira.TrafficPolicyRule{
		Description: d.Get("description").(string),
		Name:        d.Get("name").(string),
		MatchCondition: alkira.PolicyRuleMatchCondition{
			SrcIp:                 d.Get("src_ip").(string),
			DstIp:                 d.Get("dst_ip").(string),
			Dscp:                  d.Get("dscp").(string),
			Protocol:              d.Get("protocol").(string),
			SrcPortList:           convertTypeListToStringList(d.Get("src_ports").([]interface{})),
			DstPortList:           convertTypeListToStringList(d.Get("dst_ports").([]interface{})),
			SrcPrefixListId:       d.Get("src_prefix_list_id").(int),
			DstPrefixListId:       d.Get("dst_prefix_list_id").(int),
			InternetApplicationId: d.Get("internet_application_id").(int),
			ApplicationList:       convertTypeListToIntList(d.Get("application_ids").([]interface{})),
			ApplicationFamilyList: convertTypeListToIntList(d.Get("application_family_ids").([]interface{})),
		},
		RuleAction: alkira.PolicyRuleAction{
			Action:          d.Get("rule_action").(string),
			ServiceTypeList: convertTypeListToStringList(d.Get("rule_action_service_types").([]interface{})),
			ServiceList:     convertTypeListToIntList(d.Get("rule_action_service_ids").([]interface{})),
		},
	}

	return request
}
