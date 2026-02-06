package alkira

import (
	"context"
	"fmt"

	"github.com/alkiranet/alkira-client-go/alkira"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceAlkiraBillingTag() *schema.Resource {
	return &schema.Resource{
		Description:   "Manage Billing Tag.",
		CreateContext: resourceBillingTag,
		ReadContext:   resourceBillingTagRead,
		UpdateContext: resourceBillingTagUpdate,
		DeleteContext: resourceBillingTagDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Description: "Billing Tag Name.",
				Type:        schema.TypeString,
				Required:    true,
			},
			"description": {
				Description: "Billing Tag Description.",
				Type:        schema.TypeString,
				Optional:    true,
			},
		},
	}
}

func resourceBillingTag(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {

	api := alkira.NewBillingTag(m.(*alkira.AlkiraClient))

	// Construct request
	request := &alkira.BillingTag{
		Name:        d.Get("name").(string),
		Description: d.Get("description").(string),
	}

	// Send create request
	response, _, err, valErr, _ := api.Create(request)

	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(string(response.Id))

	// Handle validation error
	client := m.(*alkira.AlkiraClient)
	if client.Validate && valErr != nil {
		var diags diag.Diagnostics
		readDiags := resourceBillingTagRead(ctx, d, m)
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

	return resourceBillingTagRead(ctx, d, m)
}

func resourceBillingTagRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {

	api := alkira.NewBillingTag(m.(*alkira.AlkiraClient))

	// Get resource
	tag, _, err := api.GetById(d.Id())

	if err != nil {
		return diag.Diagnostics{{
			Severity: diag.Warning,
			Summary:  "FAILED TO GET RESOURCE",
			Detail:   fmt.Sprintf("%s", err),
		}}
	}

	d.Set("name", tag.Name)
	d.Set("description", tag.Description)

	return nil
}

func resourceBillingTagUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {

	api := alkira.NewBillingTag(m.(*alkira.AlkiraClient))

	// Construct request
	request := &alkira.BillingTag{
		Name:        d.Get("name").(string),
		Description: d.Get("description").(string),
	}

	// Send update request
	_, err, valErr, _ := api.Update(d.Id(), request)

	if err != nil {
		return diag.FromErr(err)
	}

	// Handle validation error
	client := m.(*alkira.AlkiraClient)
	if client.Validate && valErr != nil {
		var diags diag.Diagnostics
		readDiags := resourceBillingTagRead(ctx, d, m)
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

	return resourceBillingTagRead(ctx, d, m)
}

func resourceBillingTagDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {

	api := alkira.NewBillingTag(m.(*alkira.AlkiraClient))

	_, err, valErr, _ := api.Delete(d.Id())

	if err != nil {
		// Terraform may not print "with <resource address>" for destroys of objects
		// that are no longer in configuration, so include identifying context here.
		name, _ := d.GetOk("name")
		if nameStr, ok := name.(string); ok && nameStr != "" {
			return diag.FromErr(fmt.Errorf("%w alkira_billing_tag (name=%q id=%s)", err, nameStr, d.Id()))
		}
		return diag.FromErr(fmt.Errorf("%w alkira_billing_tag (id=%s)", err, d.Id()))
	}

	d.SetId("")

	// Handle validation error
	client := m.(*alkira.AlkiraClient)
	if client.Validate && valErr != nil {
		return diag.Diagnostics{{
			Severity: diag.Warning,
			Summary:  "VALIDATION (DELETE) FAILED",
			Detail:   fmt.Sprintf("%s", valErr),
		}}
	}

	return nil
}
