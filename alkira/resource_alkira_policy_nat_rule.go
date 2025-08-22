package alkira

import (
	"context"
	"fmt"

	"github.com/alkiranet/alkira-client-go/alkira"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func resourceAlkiraPolicyNatRule() *schema.Resource {
	return &schema.Resource{
		Description: "Manage policy NAT rule.\n\n" +
			"This resource is usually used along with policy resources:" +
			"`policy_nat_policy`.",
		CreateContext: resourcePolicyNatRule,
		ReadContext:   resourcePolicyNatRuleRead,
		UpdateContext: resourcePolicyNatRuleUpdate,
		DeleteContext: resourcePolicyNatRuleDelete,
		CustomizeDiff: func(ctx context.Context, d *schema.ResourceDiff, m interface{}) error {
			client := m.(*alkira.AlkiraClient)

			old, _ := d.GetChange("provision_state")

			if client.Provision == true && old == "FAILED" {
				d.SetNew("provision_state", "SUCCESS")
			}

			return nil
		},
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
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
			"provision_state": {
				Description: "the provision state of the resource.",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"category": {
				Description: "The category of NAT rule. The value could be " +
					"`DEFAULT` or `INTERNET_CONNECTOR`. Default value is " +
					"`DEFAULT`.",
				Type:     schema.TypeString,
				Optional: true,
				Default:  "DEFAULT",
				ValidateFunc: validation.StringInSlice(
					[]string{"DEFAULT", "INTERNET_CONNECTOR"}, false),
			},
			"direction": {
				Description: "The direction of NAT rule. The value could be `INBOUND` or `OUTBOUND`.",
				Type:        schema.TypeString,
				Optional:    true,
				ValidateFunc: validation.StringInSlice(
					[]string{"INBOUND", "OUTBOUND"}, false),
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
							Description: "The following protocols are supported, " +
								"`icmp`, `tcp`, `udp` or `any`.",
							Type:     schema.TypeString,
							Required: true,
							ValidateFunc: validation.StringInSlice(
								[]string{"icmp", "tcp", "udp", "any"}, false),
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
							Description: "The translation type are: `STATIC_IP`, " +
								"`DYNAMIC_IP_AND_PORT` and `NONE`. Default value " +
								"is `NONE`.",
							Type:     schema.TypeString,
							Optional: true,
							Default:  "NONE",
							ValidateFunc: validation.StringInSlice(
								[]string{"STATIC_IP", "DYNAMIC_IP_AND_PORT", "NONE"}, false),
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
						"src_addr_translation_match_and_invalidate": {
							Description: "Whether the translation match and " +
								"invalidate. Default is `true`.",
							Type:     schema.TypeBool,
							Optional: true,
						},
						"src_addr_translation_routing_track_prefixes": {
							Description: "The list of prefixes to track.",
							Type:        schema.TypeList,
							Elem:        &schema.Schema{Type: schema.TypeString},
							Optional:    true,
						},
						"src_addr_translation_routing_track_prefix_list_ids": {
							Description: "The list of prefix list IDs.",
							Type:        schema.TypeList,
							Elem:        &schema.Schema{Type: schema.TypeInt},
							Optional:    true,
						},
						"src_addr_translation_routing_track_invalidate_prefixes": {
							Description: "Whether to invalidate the track prefixes. " +
								"Default value is `false`.",
							Type:     schema.TypeBool,
							Optional: true,
							Default:  false,
						},
						"dst_addr_translation_type": {
							Description: "The translation type are: `STATIC_IP`, " +
								"`STATIC_IP_AND_PORT` , `STATIC_PORT` and `NONE`. Default " +
								"value is `NONE`.",
							Type:     schema.TypeString,
							Optional: true,
							Default:  "NONE",
							ValidateFunc: validation.StringInSlice(
								[]string{"STATIC_IP", "STATIC_IP_AND_PORT", "STATIC_PORT", "NONE"}, false),
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
							Description: "The port list to translate the " +
								"destination prefixes to.",
							Type:     schema.TypeList,
							Elem:     &schema.Schema{Type: schema.TypeString},
							Optional: true,
						},
						"dst_addr_translation_list_policy_fqdn_id": {
							Description: "The ID of policy FQDN list.",
							Type:        schema.TypeInt,
							Optional:    true,
						},
						"dst_addr_translation_advertise_to_connector": {
							Description: "Whether the destination address " +
								"should be advertised to connector.",
							Type:     schema.TypeBool,
							Optional: true,
							Default:  false,
						},
						"dst_addr_translation_routing_track_prefixes": {
							Description: "The list of prefixes to track.",
							Type:        schema.TypeList,
							Elem:        &schema.Schema{Type: schema.TypeString},
							Optional:    true,
						},
						"dst_addr_translation_routing_track_prefix_list_ids": {
							Description: "The list of prefix list IDs to track.",
							Type:        schema.TypeList,
							Elem:        &schema.Schema{Type: schema.TypeInt},
							Optional:    true,
						},
						"dst_addr_translation_routing_invalidate_prefixes": {
							Description: "Whether to invalidate the track prefixes. " +
								"Default value is `false`.",
							Type:     schema.TypeBool,
							Optional: true,
							Default:  false,
						},
						"egress_type": {
							Description: "The egress type to use with the " +
								"match. Options are are `ALKIRA_PUBLIC_IP` " +
								"or `BYOIP`.",
							Type:     schema.TypeString,
							Optional: true,
							ValidateFunc: validation.StringInSlice(
								[]string{"ALKIRA_PUBLIC_IP", "BYOIP"}, false),
						},
					},
				},
			},
		},
	}
}

func resourcePolicyNatRule(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {

	client := m.(*alkira.AlkiraClient)
	api := alkira.NewNatRule(m.(*alkira.AlkiraClient))

	// Construct request
	request, err := generatePolicyNatRuleRequest(d, m)

	if err != nil {
		return diag.FromErr(err)
	}

	// Send create request
	response, provState, err, valErr, provErr := api.Create(request)

	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(string(response.Id))

	// Handle validation error
	if client.Validate && valErr != nil {
		var diags diag.Diagnostics
		readDiags := resourcePolicyNatRuleRead(ctx, d, m)
		if readDiags.HasError() {
			diags = append(diags, readDiags...)
		}

		// Add the validation error
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Warning,
			Summary:  "VALIDATION (CREATE) FAILED",
			Detail:   fmt.Sprintf("%s", valErr),
		})

		return diags
	}

	// Set provision state
	if client.Provision == true {
		d.Set("provision_state", provState)

		if provErr != nil {
			return diag.Diagnostics{{
				Severity: diag.Warning,
				Summary:  "PROVISION (CREATE) FAILED",
				Detail:   fmt.Sprintf("%s", provErr),
			}}
		}
	}

	return resourcePolicyNatRuleRead(ctx, d, m)
}

func resourcePolicyNatRuleRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {

	client := m.(*alkira.AlkiraClient)
	api := alkira.NewNatRule(m.(*alkira.AlkiraClient))

	rule, provState, err := api.GetById(d.Id())

	if err != nil {
		return diag.Diagnostics{{
			Severity: diag.Warning,
			Summary:  "FAILED TO GET RESOURCE",
			Detail:   fmt.Sprintf("%s", err),
		}}
	}

	d.Set("name", rule.Name)
	d.Set("description", rule.Description)
	d.Set("enabled", rule.Enabled)
	d.Set("category", rule.Category)

	setNatRuleActionOptions(rule.Action, d)
	setNatRuleMatch(rule.Match, d)

	// Set provision state
	if client.Provision == true && provState != "" {
		d.Set("provision_state", provState)
	}

	return nil
}

func resourcePolicyNatRuleUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {

	client := m.(*alkira.AlkiraClient)
	api := alkira.NewNatRule(m.(*alkira.AlkiraClient))

	// Construct request
	request, err := generatePolicyNatRuleRequest(d, m)

	if err != nil {
		return diag.FromErr(err)
	}

	// Send requset
	provState, err, valErr, provErr := api.Update(d.Id(), request)

	if err != nil {
		return diag.FromErr(err)
	}

	// Handle validation error
	if client.Validate && valErr != nil {
		var diags diag.Diagnostics
		readDiags := resourcePolicyNatRuleRead(ctx, d, m)
		if readDiags.HasError() {
			diags = append(diags, readDiags...)
		}

		// Add the validation error
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Warning,
			Summary:  "VALIDATION (UPDATE) FAILED",
			Detail:   fmt.Sprintf("%s", valErr),
		})

		return diags
	}

	// Set provision state
	if client.Provision == true {
		d.Set("provision_state", provState)

		if provErr != nil {
			return diag.Diagnostics{{
				Severity: diag.Warning,
				Summary:  "PROVISION (UPDATE) FAILED",
				Detail:   fmt.Sprintf("%s", provErr),
			}}
		}
	}

	return resourcePolicyNatRuleRead(ctx, d, m)
}

func resourcePolicyNatRuleDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {

	client := m.(*alkira.AlkiraClient)
	api := alkira.NewNatRule(m.(*alkira.AlkiraClient))

	provState, err, valErr, provErr := api.Delete(d.Id())

	if err != nil {
		return diag.FromErr(err)
	}

	// Handle validation error
	if client.Validate && valErr != nil {
		return diag.Diagnostics{{
			Severity: diag.Warning,
			Summary:  "VALIDATION (DELETE) FAILED",
			Detail:   fmt.Sprintf("%s", valErr),
		}}
	}

	d.SetId("")

	if client.Provision == true && provState != "SUCCESS" {
		return diag.Diagnostics{{
			Severity: diag.Warning,
			Summary:  "PROVISION (DELETE) FAILED",
			Detail:   fmt.Sprintf("%s", provErr),
		}}
	}

	return nil
}
