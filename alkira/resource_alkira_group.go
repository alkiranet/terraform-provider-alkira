package alkira

import (
	"context"
	"fmt"

	"github.com/alkiranet/alkira-client-go/alkira"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceAlkiraGroup() *schema.Resource {
	return &schema.Resource{
		Description:   "Provide group resource.",
		CreateContext: resourceGroup,
		ReadContext:   resourceGroupRead,
		UpdateContext: resourceGroupUpdate,
		DeleteContext: resourceGroupDelete,
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
				Description: "The name of the group.",
				Type:        schema.TypeString,
				Required:    true,
			},
			"description": {
				Description: "The description of the group.",
				Type:        schema.TypeString,
				Optional:    true,
			},
			"provision_state": {
				Description: "The provisioning state of the resource.",
				Type:        schema.TypeString,
				Computed:    true,
			},
		},
	}
}

func resourceGroup(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {

	client := m.(*alkira.AlkiraClient)
	api := alkira.NewGroup(m.(*alkira.AlkiraClient))

	// Construct request
	request := &alkira.Group{
		Name:        d.Get("name").(string),
		Description: d.Get("description").(string),
	}

	// Send create request
	response, provState, err, provErr := api.Create(request)

	if err != nil {
		return diag.FromErr(err)
	}

	// Set provision state
	if client.Provision == true {
		d.Set("provision_state", provState)

		if provErr != nil {
			return diag.Diagnostics{{
				Severity: diag.Warning,
				Summary:  "PROVISION FAILED",
				Detail:   fmt.Sprintf("%s", provErr),
			}}
		}
	}

	d.SetId(string(response.Id))
	return resourceGroupRead(ctx, d, m)
}

func resourceGroupRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {

	client := m.(*alkira.AlkiraClient)
	api := alkira.NewGroup(m.(*alkira.AlkiraClient))

	group, provState, err := api.GetById(d.Id())

	if err != nil {
		return diag.FromErr(err)
	}

	d.Set("name", group.Name)
	d.Set("description", group.Description)

	// Set provision state
	if client.Provision == true && provState != "" {
		d.Set("provision_state", provState)
	}

	return nil
}

func resourceGroupUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {

	client := m.(*alkira.AlkiraClient)
	api := alkira.NewGroup(m.(*alkira.AlkiraClient))

	// Construct request
	request := &alkira.Group{
		Name:        d.Get("name").(string),
		Description: d.Get("description").(string),
	}

	// Send update request
	provState, err, provErr := api.Update(d.Id(), request)

	if err != nil {
		return diag.FromErr(err)
	}

	// Set provision state
	if client.Provision == true {
		d.Set("provision_state", provState)

		if provErr != nil {
			return diag.Diagnostics{{
				Severity: diag.Warning,
				Summary:  "PROVISION FAILED",
				Detail:   fmt.Sprintf("%s", provErr),
			}}
		}
	}

	return nil
}

func resourceGroupDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {

	client := m.(*alkira.AlkiraClient)
	api := alkira.NewGroup(m.(*alkira.AlkiraClient))

	provState, err, provErr := api.Delete(d.Id())

	if err != nil {
		return diag.FromErr(err)
	}

	if client.Provision == true && provState != "SUCCESS" {
		return diag.FromErr(provErr)
	}

	d.SetId("")
	return nil
}
