package alkira

import (
	"context"
	"fmt"

	"github.com/alkiranet/alkira-client-go/alkira"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func resourceAlkiraVirtualNetworkManagerAzure() *schema.Resource {
	return &schema.Resource{
		Description:   "Manager Virtual Network Manager for Azure.",
		CreateContext: resourceAlkiraVirtualNetworkManagerAzureCreate,
		UpdateContext: resourceAlkiraVirtualNetworkManagerAzureUpdate,
		DeleteContext: resourceAlkiraVirtualNetworkManagerAzureDelete,
		ReadContext:   resourceAlkiraVirtualNetworkManagerAzureRead,
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
    "name":{
    Description: "Name of the ",
    }

    },

	}

}

func resourceAlkiraVirtualNetworkManagerAzureCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	return nil
}

func resourceAlkiraVirtualNetworkManagerAzureUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	return nil
}

func resourceAlkiraVirtualNetworkManagerAzureDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	return nil
}

func resourceAlkiraVirtualNetworkManagerAzureRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	return nil
}

func generateVirtualNetworkManagerAzureRequest(d *schema.ResourceData, m interface{}) (*alkira.VirtualNetworkManagerAzure, error) {
	return nil, nil
}
