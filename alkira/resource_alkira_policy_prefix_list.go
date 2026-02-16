package alkira

import (
	"bytes"
	"context"
	"fmt"

	"github.com/alkiranet/alkira-client-go/alkira"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceAlkiraPolicyPrefixList() *schema.Resource {
	return &schema.Resource{
		Description:   "Manage policy prefix list.",
		CreateContext: resourcePolicyPrefixList,
		ReadContext:   resourcePolicyPrefixListRead,
		UpdateContext: resourcePolicyPrefixListUpdate,
		DeleteContext: resourcePolicyPrefixListDelete,
		SchemaVersion: 1,
		StateUpgraders: []schema.StateUpgrader{
			{
				Version: 0,
				Type:    resourcePolicyPrefixListV0().CoreConfigSchema().ImpliedType(),
				Upgrade: resourcePolicyPrefixListStateUpgradeV0,
			},
		},
		CustomizeDiff: func(ctx context.Context, d *schema.ResourceDiff, m interface{}) error {
			client := m.(*alkira.AlkiraClient)

			old, _ := d.GetChange("provision_state")

			if client.Provision && old == "FAILED" {
				d.SetNew("provision_state", "SUCCESS")
			}

			return nil
		},
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Description: "The name of the prefix list.",
				Type:        schema.TypeString,
				Required:    true,
			},
			"description": {
				Description: "The description of the prefix list.",
				Type:        schema.TypeString,
				Optional:    true,
			},
			"provision_state": {
				Description: "The provisioning state of the resource.",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"prefixes": {
				Description: "A list of prefixes. (**DEPRECATED**)",
				Type:        schema.TypeSet,
				Optional:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
			"prefix": {
				Description: "Prefix with description. This new block should " +
					"replace the old `prefixes` field.",
				Type:     schema.TypeSet,
				Optional: true,
				Set:      prefixHash,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"cidr": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "The network prefix in CIDR notation.",
						},
						"description": {
							Type:        schema.TypeString,
							Optional:    true,
							Default:     "",
							Description: "Description for the prefix.",
						},
					},
				},
			},
			"prefix_range": {
				Description: "A valid prefix range that could be used to " +
					"define a prefix of type `ROUTE`.",
				Type:     schema.TypeSet,
				Optional: true,
				Set:      prefixRangeHash,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"prefix": {
							Description: "A valid CIDR as prefix in " +
								"`x.x.x.x/m` format.",
							Type:     schema.TypeString,
							Required: true,
						},
						"description": {
							Type:     schema.TypeString,
							Optional: true,
							Default:  "",
						},
						"le": {
							Description: "Integer less than `32` but " +
								"greater than mask `m` in prefix",
							Type:     schema.TypeInt,
							Optional: true,
							Default:  0,
						},
						"ge": {
							Description: "Integer less than `32` but " +
								"greater than mask `m` in prefix and less " +
								"than `le`.",
							Type:     schema.TypeInt,
							Optional: true,
							Default:  0,
						},
					},
				},
			},
		},
	}
}

func resourcePolicyPrefixList(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {

	client := m.(*alkira.AlkiraClient)
	api := alkira.NewPolicyPrefixList(m.(*alkira.AlkiraClient))

	// Construct request
	request, err := generatePolicyPrefixListRequest(d)

	if err != nil {
		return diag.FromErr(err)
	}

	// Send request
	response, provState, err, valErr, provErr := api.Create(request)

	if err != nil {
		return diag.FromErr(err)
	}

	// Handle validation error
	if client.Validate && valErr != nil {
		var diags diag.Diagnostics
		readDiags := resourcePolicyPrefixListRead(ctx, d, m)
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

	// Set provision state
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

	d.SetId(string(response.Id))
	return resourcePolicyPrefixListRead(ctx, d, m)
}

func resourcePolicyPrefixListRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {

	client := m.(*alkira.AlkiraClient)
	api := alkira.NewPolicyPrefixList(m.(*alkira.AlkiraClient))

	list, provState, err := api.GetById(d.Id())

	if err != nil {
		return diag.Diagnostics{{
			Severity: diag.Warning,
			Summary:  "FAILED TO GET RESOURCE",
			Detail:   fmt.Sprintf("%s", err),
		}}
	}

	d.Set("name", list.Name)
	d.Set("description", list.Description)

	// Set "prefix" block
	setPrefix(d, list.Prefixes, list.PrefixDetails)

	// Set "prefix_ranges" block
	setPrefixRanges(d, list.PrefixRanges)

	// Set provision state
	if client.Provision && provState != "" {
		d.Set("provision_state", provState)
	}

	return nil
}

func resourcePolicyPrefixListUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {

	client := m.(*alkira.AlkiraClient)
	api := alkira.NewPolicyPrefixList(m.(*alkira.AlkiraClient))

	// Construct request
	request, err := generatePolicyPrefixListRequest(d)

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
		readDiags := resourcePolicyPrefixListRead(ctx, d, m)
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

	// Set provision state
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

	return resourcePolicyPrefixListRead(ctx, d, m)
}

func resourcePolicyPrefixListDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {

	client := m.(*alkira.AlkiraClient)
	api := alkira.NewPolicyPrefixList(m.(*alkira.AlkiraClient))

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

	if client.Provision && provState != "SUCCESS" {
		return diag.Diagnostics{{
			Severity: diag.Warning,
			Summary:  "PROVISION (DELETE) FAILED",
			Detail:   fmt.Sprintf("%s", provErr),
		}}
	}

	return nil
}

// prefixHash computes a hash for a prefix block based on its cidr field.
// This allows Terraform to identify prefixes by their CIDR value rather than
// their position in the list, preventing unwanted reordering when prefixes
// are deleted from the middle of the list.
func prefixHash(v interface{}) int {
	var buf bytes.Buffer
	m := v.(map[string]interface{})
	fmt.Fprintf(&buf, "%s-", m["cidr"])
	return schema.HashString(buf.String())
}

// prefixRangeHash computes a hash for a prefix_range block based on its
// prefix, le, and ge fields. This allows Terraform to identify prefix ranges
// by their content rather than their position in the list, preventing unwanted
// reordering when ranges are deleted from the middle of the list.
func prefixRangeHash(v interface{}) int {
	var buf bytes.Buffer
	m := v.(map[string]interface{})
	fmt.Fprintf(&buf, "%s-%d-%d-", m["prefix"], toInt(m["le"]), toInt(m["ge"]))
	return schema.HashString(buf.String())
}
