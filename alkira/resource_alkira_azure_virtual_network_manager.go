package alkira

import (
	"context"
	"fmt"
	"time"

	"github.com/alkiranet/alkira-client-go/alkira"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceAzureVirtualNetworkManager() *schema.Resource {
	return &schema.Resource{
		Description:   "Manager Virtual Network Manager for Azure.",
		CreateContext: resourceAzureVirtualNetworkManagerCreate,
		UpdateContext: resourceAzureVirtualNetworkManagerUpdate,
		DeleteContext: resourceAzureVirtualNetworkManagerDelete,
		ReadContext:   resourceAzureVirtualNetworkManagerRead,
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
				Description: "The name of the connector.",
				Type:        schema.TypeString,
				Required:    true,
			},
		},
	}

}

func resourceAzureVirtualNetworkManagerCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {

	// INIT
	api := alkira.NewAzureVirtualNetworkManager(m.(*alkira.AlkiraClient))

	request, err := generateAzureVirtualNetworkManagerRequest(d, m)

	if err != nil {
		return diag.FromErr(err)
	}

	resource, _, err, _ := api.Create(request)

	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId(string(resource.Id))

	// WAITING FOR THE STATE
	state := resource.State
	for state != "ACTIVE" {
		resource, _, err := api.GetById(d.Id())

		if err != nil {
			return diag.Diagnostics{{
				Severity: diag.Warning,
				Summary:  "FAILED TO GET RESOURCE",
				Detail:   fmt.Sprintf("%s", err),
			}}
		}

		state = resource.State
		time.Sleep(5 * time.Second)

	}
	return resourceAzureVirtualNetworkManagerRead(ctx, d, m)
}

func resourceAzureVirtualNetworkManagerUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	return nil
}

func resourceAzureVirtualNetworkManagerDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	return nil
}

func resourceAzureVirtualNetworkManagerRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	return nil
}

func generateAzureVirtualNetworkManagerRequest(d *schema.ResourceData, m interface{}) (*alkira.AzureVirtualNetworkManager, error) {
	return nil, nil
}
