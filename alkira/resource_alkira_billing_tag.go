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
	response, _, err, _ := api.Create(request)

	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(string(response.Id))

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
	_, err, _ := api.Update(d.Id(), request)

	if err != nil {
		return diag.FromErr(err)
	}

	return resourceBillingTagRead(ctx, d, m)
}

func resourceBillingTagDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {

	api := alkira.NewBillingTag(m.(*alkira.AlkiraClient))

	_, err, _ := api.Delete(d.Id())

	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId("")
	return nil
}
