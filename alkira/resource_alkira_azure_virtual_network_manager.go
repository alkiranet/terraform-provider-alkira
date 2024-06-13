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
				Description: "The name of the Azure Virtual Network Manager.",
				Type:        schema.TypeString,
				Required:    true,
			},
			"region": {
				Description: "The Azure Region where the Azure Virtual Network" +
					" Manager is created. eg: `eastus`.",
				Type:     schema.TypeString,
				Required: true,
			},
			"subscription_id": {
				Description: "The ID of the Azure Subscription account in which the" +
					" Azure Virtual Network Manager is created.",
				Type:     schema.TypeString,
				Required: true,
			},
			"resource_group": {
				Description: "The name of the Azure Resource Group in which the" +
					" Azure Virtual Network Manager is created.",
				Type:     schema.TypeString,
				Required: true,
			},
			"Description": {
				Description: "The description of the Azure Virtual Network Manager.",
				Type:        schema.TypeString,
				Optional:    true,
			},
			"credential_id": {
				Description: "ID of credential managed by Credential Manager.",
				Type:        schema.TypeString,
				Required:    true,
			},
			"subscriptions_in_scope": {
				Description: "List IDs of Azure Subscription Accounts that will be" +
					" managed by this Azure Virtual Network Manager.",
				Type:     schema.TypeList,
				Required: true,
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
	// CHECK FOR CHANGES IN ANY FIELD EXCEPT DESCRIPTION AND RETURN ERR.
	if d.HasChange("description") {

		// INIT
		api := alkira.NewAzureVirtualNetworkManager(m.(*alkira.AlkiraClient))

		request, err := generateAzureVirtualNetworkManagerRequest(d, m)

		if err != nil {
			return diag.FromErr(err)
		}

		// UPDATE
		_, err, _ = api.Update(d.Id(), request)

		return nil
	}

	// RETURN ERROR IF ANY FIELD EXCEPT DESCRIPTION IS CHANGED.
	return diag.Diagnostics{{
		Severity: diag.Error,
		Summary:  "INVALID UPDATE!",
		Detail:   "ONLY THE DESCRIPTION FIELD CAN BE UPDATED.",
	}}
}

func resourceAzureVirtualNetworkManagerDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {

	// INIT
	api := alkira.NewAzureVirtualNetworkManager(m.(*alkira.AlkiraClient))

	// DELETE
	_, err, _ := api.Delete(d.Id())

	if err != nil {
		return diag.FromErr(err)
	}

	// SETTING THE ID TO EMPTY STRING
	d.SetId("")
	return nil
}

func resourceAzureVirtualNetworkManagerRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {

	// INIT
	api := alkira.NewAzureVirtualNetworkManager(m.(*alkira.AlkiraClient))
	resource, _, err := api.GetById(d.Id())

	if err != nil {
		return diag.Diagnostics{{
			Severity: diag.Warning,
			Summary:  "FAILED TO GET RESOURCE",
			Detail:   fmt.Sprintf("%s", err),
		}}
	}

	d.Set("name", resource.Name)
	d.Set("region", resource.Region)
	d.Set("subscription_id", resource.SubscriptionId)
	d.Set("resource_group", resource.ResourceGroup)
	d.Set("description", resource.Description)
	d.Set("credential_id", resource.CredentialId)
	d.Set("subscriptions_in_scope", resource.SubscriptionsInScope)

	return nil
}

func generateAzureVirtualNetworkManagerRequest(d *schema.ResourceData, m interface{}) (*alkira.AzureVirtualNetworkManager, error) {

	request := &alkira.AzureVirtualNetworkManager{
		Name:                 d.Get("name").(string),
		Region:               d.Get("region").(string),
		SubscriptionId:       d.Get("subscription_id").(string),
		ResourceGroup:        d.Get("resource_group").(string),
		Description:          d.Get("description").(string),
		CredentialId:         d.Get("credential_id").(string),
		SubscriptionsInScope: d.Get("subscriptions_in_scope").([]string),
	}
	return request, nil
}
