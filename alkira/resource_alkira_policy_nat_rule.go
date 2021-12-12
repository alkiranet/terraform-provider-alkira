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
				Type:        schema.TypeString,
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
						"dst_prefixes": {
							Description: "The list of prefixes for destination.",
							Type:        schema.TypeList,
							Elem:        &schema.Schema{Type: schema.TypeString},
							Optional:    true,
						},
						"src_prefix_list_ids": {
							Description: "The list of prefix ids as source.",
							Type:        schema.TypeList,
							Elem:        &schema.Schema{Type: schema.TypeInt},
							Optional:    true,
						},
						"dst_prefix_list_ids": {
							Description: "The list of prefix ids as destination.",
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
						"dst_ports": {
							Description: "The list of ports for destination.",
							Type:        schema.TypeList,
							Elem:        &schema.Schema{Type: schema.TypeString},
							Optional:    true,
						},
						"protocol": {
							Description:  "The following protocols are supported, `ICMP`, `TCP`, `UDP` or `ANY`.",
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: validation.StringInSlice([]string{"ICMP", "TCP", "UDP", "ANY"}, false),
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
							Description:  "The translation type are: `STATIC_IP`, `DYNAMIC`, `DYNAMIC_IP_AND_PORT` and `NONE`.",
							Type:         schema.TypeString,
							Optional:     true,
							ValidateFunc: validation.StringInSlice([]string{"STATIC_IP", "DYNAMIC", "DYNAMIC_IP_AND_PORT", "NONE"}, false),
						},
						"src_addr_translation_prefixes": {
							Description: "The list of prefixes.",
							Type:        schema.TypeList,
							Elem:        &schema.Schema{Type: schema.TypeString},
							Optional:    true,
						},
						"src_addr_translation_prefiex_list_ids": {
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
							Description:  "The translation type are: `STATIC_IP`, `DYNAMIC`, `DYNAMIC_IP_AND_PORT` and `NONE`.",
							Type:         schema.TypeString,
							Optional:     true,
							ValidateFunc: validation.StringInSlice([]string{"STATIC_IP", "DYNAMIC", "DYNAMIC_IP_AND_PORT", "NONE"}, false),
						},
						"dst_addr_translation_prefixes": {
							Description: "The list of prefixes.",
							Type:        schema.TypeList,
							Elem:        &schema.Schema{Type: schema.TypeString},
							Optional:    true,
						},
						"dst_addr_translation_prefiex_list_ids": {
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

	rule, err := client.GetNatRule(d.Id())

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

// expandPolicyNatRuleMatch expand "match" section
func expandPolicyNatRuleMatch(in *schema.Set) *alkira.NatRuleMatch {
	if in == nil || in.Len() == 0 || in.Len() > 1 {
		log.Printf("[ERROR] invalid match section (%d)", in.Len())
		return nil
	}

	match := alkira.NatRuleMatch{}

	for _, m := range in.List() {

		matchValue := m.(map[string]interface{})

		if v, ok := matchValue["src_prefixes"].([]string); ok {
			match.SourcePrefixes = v
		}
		if v, ok := matchValue["src_prefix_list_ids"].([]int); ok {
			match.SourcePrefixListIds = v
		}
		if v, ok := matchValue["dst_prefixes"].([]string); ok {
			match.DestPrefixes = v
		}
		if v, ok := matchValue["dst_prefix_list_ids"].([]int); ok {
			match.DestPrefixListIds = v
		}
		if v, ok := matchValue["src_ports"].([]string); ok {
			match.SourcePortList = v
		}
		if v, ok := matchValue["dst_ports"].([]string); ok {
			match.DestPortList = v
		}
		if v, ok := matchValue["protocol"].(string); ok {
			match.Protocol = v
		}
	}

	return &match
}

// expandPolicyNatRuleAction expand "action" section
func expandPolicyNatRuleAction(in *schema.Set) *alkira.NatRuleAction {
	if in == nil || in.Len() == 0 || in.Len() > 1 {
		log.Printf("[ERROR] invalid action section (%d)", in.Len())
		return nil
	}

	st := alkira.NatRuleActionSrcTranslation{}
	dt := alkira.NatRuleActionDstTranslation{}

	for _, m := range in.List() {

		actionValue := m.(map[string]interface{})

		if v, ok := actionValue["src_addr_translation_type"].(string); ok {
			st.TranslationType = v
		}
		if v, ok := actionValue["src_addr_translation_prefixes"].([]string); ok {
			st.TranslatedPrefixes = v
		}
		if v, ok := actionValue["src_addr_translation_prefix_list_ids"].([]int); ok {
			st.TranslatedPrefixListIds = v
		}
		if v, ok := actionValue["src_addr_translation_bidirectional"].(bool); ok {
			st.Bidirectional = v
		}
		if v, ok := actionValue["src_addr_translation_match_and_invalidate"].(bool); ok {
			st.MatchAndInvalidate = v
		}
		if v, ok := actionValue["dst_addr_translation_type"].(string); ok {
			dt.TranslationType = v
		}
		if v, ok := actionValue["dst_addr_translation_prefixes"].([]string); ok {
			dt.TranslatedPrefixes = v
		}
		if v, ok := actionValue["dst_addr_translation_prefix_list_ids"].([]int); ok {
			dt.TranslatedPrefixListIds = v
		}
		if v, ok := actionValue["dst_addr_translation_ports"].([]string); ok {
			dt.TranslatedPortList = v
		}
		if v, ok := actionValue["dst_addr_translation_bidirectional"].(bool); ok {
			dt.Bidirectional = v
		}
		if v, ok := actionValue["dst_addr_translation_advertise_to_connector"].(bool); ok {
			dt.AdvertiseToConnector = v
		}
	}

	action := alkira.NatRuleAction{
		SourceAddressTranslation:      st,
		DestinationAddressTranslation: dt,
	}

	return &action
}
