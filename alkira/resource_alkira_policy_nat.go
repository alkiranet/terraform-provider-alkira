package alkira

import (
	"context"
	"fmt"

	"github.com/alkiranet/alkira-client-go/alkira"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func resourceAlkiraPolicyNat() *schema.Resource {
	return &schema.Resource{
		Description:   "Manage NAT policy.",
		CreateContext: resourcePolicyNat,
		ReadContext:   resourcePolicyNatRead,
		UpdateContext: resourcePolicyNatUpdate,
		DeleteContext: resourcePolicyNatDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Description: "The name of the policy.",
				Type:        schema.TypeString,
				Required:    true,
			},
			"description": {
				Description: "The description of the policy.",
				Type:        schema.TypeString,
				Optional:    true,
			},
			"category": {
				Description: "The category of NAT policy. " +
					"The vaule could be `DEFAULT` or " +
					"`INTERNET_CONNECTOR`. Default value is " +
					"`DEFAULT`.",
				Type:     schema.TypeString,
				Optional: true,
				Default:  "DEFAULT",
				ValidateFunc: validation.StringInSlice([]string{
					"DEFAULT", "INTERNET_CONNECTOR"}, false),
			},
			"type": {
				Description: "The type of NAT policy, currently only " +
					"`INTRA_SEGMENT` is supported.",
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringInSlice([]string{"INTRA_SEGMENT"}, false),
			},
			"segment_id": {
				Description: "IDs of the segment that will define the policy" +
					"scope.",
				Type:     schema.TypeString,
				Required: true,
			},
			"included_group_ids": {
				Description: "Defines the scope for the policy. Connectors " +
					"associated with groups defined here is where this policy " +
					"would be applied.",
				Type:     schema.TypeSet,
				Elem:     &schema.Schema{Type: schema.TypeInt},
				Required: true,
			},
			"excluded_group_ids": {
				Description: "Excludes given associated connector from " +
					"`included_groups`. Implicit group of a branch or on-premise " +
					"connector for which a user defined group is used in " +
					"`included_groups` can be used here.",
				Type:     schema.TypeSet,
				Elem:     &schema.Schema{Type: schema.TypeInt},
				Optional: true,
			},
			"nat_rule_ids": {
				Description: "The list of NAT rules to be applied by the policy.",
				Type:        schema.TypeList,
				Elem:        &schema.Schema{Type: schema.TypeInt},
				Required:    true,
			},
			"allow_overlapping_translated_source_addresses": {
				Description: "Allow overlapping translated source address. " +
					"Default value is `false`. (**BETA**)",
				Type:     schema.TypeBool,
				Default:  false,
				Optional: true,
			},
		},
	}
}

func resourcePolicyNat(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {

	// INIT
	client := m.(*alkira.AlkiraClient)
	api := alkira.NewNatPolicy(m.(*alkira.AlkiraClient))

	// Construct request
	request, err := generatePolicyNatRequest(d, m)

	if err != nil {
		return diag.FromErr(err)
	}

	// Send request
	resource, provState, err, valErr, provErr := api.Create(request)

	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(string(resource.Id))

	// Handle validation error
	if client.Validate && valErr != nil {
		var diags diag.Diagnostics
		readDiags := resourcePolicyNatRead(ctx, d, m)
		if readDiags.HasError() {
			diags = append(diags, readDiags...)
		}

		// Add the validation error
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "VALIDATION (CREATE) FAILED",
			Detail:   fmt.Sprintf("%s", valErr),
		})

		return diags
	}

	if client.Provision {
		d.Set("provision_state", provState)

		if provErr != nil {
			return diag.Diagnostics{{
				Severity: diag.Warning,
				Summary:  "PROVISION (CREATE) FAILED",
				Detail:   fmt.Sprintf("%s", provErr),
			}}
		}
	}

	return resourcePolicyNatRead(ctx, d, m)
}

func resourcePolicyNatRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {

	client := m.(*alkira.AlkiraClient)
	api := alkira.NewNatPolicy(m.(*alkira.AlkiraClient))

	policy, provState, err := api.GetById(d.Id())

	if err != nil {
		// Check if resource was deleted outside Terraform (404)
		if handleResourceNotFound(err, d, "NAT Policy") {
			return nil
		}
		// For other errors, return as warning
		return diag.Diagnostics{{
			Severity: diag.Warning,
			Summary:  "FAILED TO GET RESOURCE",
			Detail:   fmt.Sprintf("%s", err),
		}}
	}

	d.Set("name", policy.Name)
	d.Set("description", policy.Description)
	d.Set("type", policy.Type)
	d.Set("included_group_ids", policy.IncludedGroups)
	d.Set("excluded_group_ids", policy.ExcludedGroups)
	d.Set("nat_rule_ids", policy.NatRuleIds)
	d.Set("category", policy.Category)
	d.Set("allow_overlapping_translated_source_address", policy.AllowOverlappingTranslatedPrefixes)

	// Get segment
	segmentId, err := getSegmentIdByName(policy.Segment, m)

	if err != nil {
		return diag.FromErr(err)
	}

	d.Set("segment_id", segmentId)

	// Set provision state
	if client.Provision && provState != "" {
		d.Set("provision_state", provState)
	}

	return nil
}

func resourcePolicyNatUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {

	// INIT
	client := m.(*alkira.AlkiraClient)
	api := alkira.NewNatPolicy(m.(*alkira.AlkiraClient))

	// Construct request
	request, err := generatePolicyNatRequest(d, m)

	if err != nil {
		return diag.FromErr(err)
	}

	// Send update request
	provState, err, valErr, provErr := api.Update(d.Id(), request)

	if err != nil {
		return diag.FromErr(err)
	}

	// Handle validation error
	if client.Validate && valErr != nil {
		var diags diag.Diagnostics
		readDiags := resourcePolicyNatRead(ctx, d, m)
		if readDiags.HasError() {
			diags = append(diags, readDiags...)
		}

		// Add the validation error
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "VALIDATION (UPDATE) FAILED",
			Detail:   fmt.Sprintf("%s", valErr),
		})

		return diags
	}

	if client.Provision {
		d.Set("provision_state", provState)

		if provErr != nil {
			return diag.Diagnostics{{
				Severity: diag.Warning,
				Summary:  "PROVISION (UPDATE) FAILED",
				Detail:   fmt.Sprintf("%s", provErr),
			}}
		}
	}

	return resourcePolicyNatRead(ctx, d, m)
}

func resourcePolicyNatDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {

	client := m.(*alkira.AlkiraClient)
	api := alkira.NewNatPolicy(m.(*alkira.AlkiraClient))

	provState, err, valErr, provErr := api.Delete(d.Id())

	if err != nil {
		return diag.FromErr(err)
	}

	// Handle validation error
	if client.Validate && valErr != nil {
		return diag.Diagnostics{{
			Severity: diag.Error,
			Summary:  "VALIDATION (DELETE) FAILED",
			Detail:   fmt.Sprintf("%s", valErr),
		}}
	}

	d.SetId("")

	// Set provision state
	if client.Provision {
		d.Set("provision_state", provState)

		if provErr != nil {
			return diag.Diagnostics{{
				Severity: diag.Warning,
				Summary:  "PROVISION(DELETE) FAILED",
				Detail:   fmt.Sprintf("%s", provErr),
			}}
		}
	}

	return nil
}

func generatePolicyNatRequest(d *schema.ResourceData, m interface{}) (*alkira.NatPolicy, error) {

	//
	// Segment
	//
	segmentName, err := getSegmentNameById(d.Get("segment_id").(string), m)

	if err != nil {
		return nil, err
	}

	allowOverlappingTranslatedPrefixes := new(bool)

	if d.Get("allow_overlapping_translated_source_addresses") != nil {
		*allowOverlappingTranslatedPrefixes = d.Get("allow_overlapping_translated_source_addresses").(bool)
	}

	// Assemble request
	policy := &alkira.NatPolicy{
		Name:                               d.Get("name").(string),
		Description:                        d.Get("description").(string),
		Type:                               d.Get("type").(string),
		Segment:                            segmentName,
		IncludedGroups:                     convertTypeSetToIntList(d.Get("included_group_ids").(*schema.Set)),
		ExcludedGroups:                     convertTypeSetToIntList(d.Get("excluded_group_ids").(*schema.Set)),
		NatRuleIds:                         convertTypeListToIntList(d.Get("nat_rule_ids").([]interface{})),
		Category:                           d.Get("category").(string),
		AllowOverlappingTranslatedPrefixes: allowOverlappingTranslatedPrefixes,
	}

	return policy, nil
}
