package alkira

import (
	"log"

	"github.com/alkiranet/alkira-client-go/alkira"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func resourceAlkiraPolicyNatRule() *schema.Resource {
	return &schema.Resource{
		Description: "Manage policy NAT rule.\n\n" +
			"This resource is usually used along with policy resources:" +
			"`policy_nat_policy`.",
		Create: resourcePolicyNatRule,
		Read:   resourcePolicyNatRuleRead,
		Update: resourcePolicyNatRuleUpdate,
		Delete: resourcePolicyNatRuleDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
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
			"enabled": {
				Description: "Enable the rule or not.",
				Type:        schema.TypeBool,
				Required:    true,
			},
			"match": {
				Description: "Match condition for the rule.",
				Type:        schema.TypeSet,
				Required:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"src_prefixes": {
							Description: "The list of prefixes for source.",
							Type:        schema.TypeList,
							Elem:        &schema.Schema{Type: schema.TypeString},
							Optional:    true,
						},
						"src_prefix_list_ids": {
							Description: "The list of prefix IDs as source.",
							Type:        schema.TypeList,
							Elem:        &schema.Schema{Type: schema.TypeInt},
							Optional:    true,
						},
						"src_ports": {
							Description: "The list of ports for source.",
							Type:        schema.TypeList,
							Elem:        &schema.Schema{Type: schema.TypeString},
							Optional:    true,
						},
						"dst_prefixes": {
							Description: "The list of prefixes for destination.",
							Type:        schema.TypeList,
							Elem:        &schema.Schema{Type: schema.TypeString},
							Optional:    true,
						},
						"dst_prefix_list_ids": {
							Description: "The list of prefix IDs as destination.",
							Type:        schema.TypeList,
							Elem:        &schema.Schema{Type: schema.TypeInt},
							Optional:    true,
						},
						"dst_ports": {
							Description: "The list of ports for destination.",
							Type:        schema.TypeList,
							Elem:        &schema.Schema{Type: schema.TypeString},
							Optional:    true,
						},
						"protocol": {
							Description:  "The following protocols are supported, `icmp`, `tcp`, `udp` or `any`.",
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: validation.StringInSlice([]string{"icmp", "tcp", "udp", "any"}, false),
						},
					},
				},
			},
			"action": {
				Description: "The action of the rule.",
				Type:        schema.TypeSet,
				Required:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"src_addr_translation_type": {
							Description:  "The translation type are: `STATIC_IP`, `DYNAMIC_IP_AND_PORT` and `NONE`. Default value is `NONE`.",
							Type:         schema.TypeString,
							Optional:     true,
							Default:      "NONE",
							ValidateFunc: validation.StringInSlice([]string{"STATIC_IP", "DYNAMIC_IP_AND_PORT", "NONE"}, false),
						},
						"src_addr_translation_prefixes": {
							Description: "The list of prefixes.",
							Type:        schema.TypeList,
							Elem:        &schema.Schema{Type: schema.TypeString},
							Optional:    true,
						},
						"src_addr_translation_prefix_list_ids": {
							Description: "The list of prefix list IDs.",
							Type:        schema.TypeList,
							Elem:        &schema.Schema{Type: schema.TypeInt},
							Optional:    true,
						},
						"src_addr_translation_bidirectional": {
							Description: "Is the translation bidirectional.",
							Type:        schema.TypeBool,
							Optional:    true,
						},
						"src_addr_translation_match_and_invalidate": {
							Description: "Whether the translation match and invalidate.",
							Type:        schema.TypeBool,
							Optional:    true,
						},
						"dst_addr_translation_type": {
							Description:  "The translation type are: `STATIC_IP`, `DYNAMIC_IP_AND_PORT` and `NONE`. Default value is `NONE`.",
							Type:         schema.TypeString,
							Optional:     true,
							Default:      "NONE",
							ValidateFunc: validation.StringInSlice([]string{"STATIC_IP", "DYNAMIC_IP_AND_PORT", "NONE"}, false),
						},
						"dst_addr_translation_prefixes": {
							Description: "The list of prefixes.",
							Type:        schema.TypeList,
							Elem:        &schema.Schema{Type: schema.TypeString},
							Optional:    true,
						},
						"dst_addr_translation_prefix_list_ids": {
							Description: "The list of prefix list IDs.",
							Type:        schema.TypeList,
							Elem:        &schema.Schema{Type: schema.TypeInt},
							Optional:    true,
						},
						"dst_addr_translation_ports": {
							Description: "The port list to translate the destination prefixes to.",
							Type:        schema.TypeList,
							Elem:        &schema.Schema{Type: schema.TypeString},
							Optional:    true,
						},
						"dst_addr_translation_bidirectional": {
							Description: "Is the translation bidirectional.",
							Type:        schema.TypeBool,
							Optional:    true,
						},
						"dst_addr_translation_advertise_to_connector": {
							Description: "Whether the destination address should be advertised to connector.",
							Type:        schema.TypeBool,
							Optional:    true,
						},
					},
				},
			},
		},
	}
}

func resourcePolicyNatRule(d *schema.ResourceData, m interface{}) error {
	client := m.(*alkira.AlkiraClient)

	request, err := generatePolicyNatRuleRequest(d, m)

	if err != nil {
		log.Printf("[ERROR] Failed to generate policy rule request")
		return err
	}

	id, err := client.CreateNatRule(request)

	if err != nil {
		log.Printf("[ERROR] Failed to create policy rule")
		return err
	}

	d.SetId(id)
	return resourcePolicyNatRuleRead(d, m)
}

func resourcePolicyNatRuleRead(d *schema.ResourceData, m interface{}) error {
	client := m.(*alkira.AlkiraClient)

	rule, err := client.GetNatRuleById(d.Id())

	if err != nil {
		log.Printf("[ERROR] Failed to get policy rule %s", d.Id())
		return err
	}

	d.Set("name", rule.Name)
	d.Set("description", rule.Description)
	d.Set("enabled", rule.Enabled)

	return nil
}

func resourcePolicyNatRuleUpdate(d *schema.ResourceData, m interface{}) error {
	client := m.(*alkira.AlkiraClient)

	request, err := generatePolicyNatRuleRequest(d, m)

	if err != nil {
		log.Printf("[ERROR] Failed to generate policy rule request")
		return err
	}

	err = client.UpdateNatRule(d.Id(), request)

	if err != nil {
		log.Printf("[ERROR] Failed to update policy rule %s", d.Id())
		return err
	}

	return resourcePolicyNatRuleRead(d, m)
}

func resourcePolicyNatRuleDelete(d *schema.ResourceData, m interface{}) error {
	client := m.(*alkira.AlkiraClient)

	return client.DeleteNatRule(d.Id())
}

func generatePolicyNatRuleRequest(d *schema.ResourceData, m interface{}) (*alkira.NatRule, error) {

	match := expandPolicyNatRuleMatch(d.Get("match").(*schema.Set))
	action := expandPolicyNatRuleAction(d.Get("action").(*schema.Set))

	request := &alkira.NatRule{
		Name:        d.Get("name").(string),
		Description: d.Get("description").(string),
		Enabled:     d.Get("enabled").(bool),
		Match:       *match,
		Action:      *action,
	}

	return request, nil
}
