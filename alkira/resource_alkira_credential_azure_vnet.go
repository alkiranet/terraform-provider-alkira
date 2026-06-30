package alkira

import (
	"context"
	"fmt"
	"log"

	"github.com/alkiranet/alkira-client-go/alkira"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceAlkiraCredentialAzureVnet() *schema.Resource {
	return &schema.Resource{
		Description: "Credential for interacting with Azure.\n\n" +
			"You could also provide your credentials via the following " +
			"environmental variables:\n\n * AK_AZURE_APPLICATION_ID\n " +
			"* AK_AZURE_SUBSCRIPTION_ID\n * AK_AZURE_SECRET_KEY\n " +
			"* AK_AZURE_TENANT_ID\n * AK_AZURE_ENVIRONMENT\n ",
		CreateContext: resourceCredentialAzureVnet,
		ReadContext:   resourceCredentialAzureVnetRead,
		UpdateContext: resourceCredentialAzureVnetUpdate,
		DeleteContext: resourceCredentialAzureVnetDelete,
		Importer: &schema.ResourceImporter{
			StateContext: importWithReadValidation(resourceCredentialAzureVnetRead),
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Description: "The name of the credential.",
				Type:        schema.TypeString,
				Required:    true,
			},
			"application_id": {
				Description: "Azure Application ID.",
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc(
					"AK_AZURE_APPLICATION_ID",
					nil),
			},
			"subscription_id": {
				Description: "Azure subscription ID.",
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc(
					"AK_AZURE_SUBSCRIPTION_ID",
					nil),
			},
			"secret_key": {
				Description: "Azure Secret Key.",
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc(
					"AK_AZURE_SECRET_KEY",
					nil),
			},
			"tenant_id": {
				Description: "Azure Tenant ID.",
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc(
					"AK_AZURE_TENANT_ID",
					nil),
			},
			"environment": {
				Description: "Azure environment can be `AZURE`, " +
					"`AZURE_CHINA` or `AZURE_US_GOVERNMENT`. The " +
					"default value is `AZURE`.",
				Type:     schema.TypeString,
				Optional: true,
				DefaultFunc: schema.EnvDefaultFunc(
					"AK_AZURE_ENVIRONMENT",
					nil),
			},
		},
	}
}

func resourceCredentialAzureVnet(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*alkira.AlkiraClient)

	c := alkira.CredentialAzureVnet{
		ApplicationId:  d.Get("application_id").(string),
		SecretKey:      d.Get("secret_key").(string),
		SubscriptionId: d.Get("subscription_id").(string),
		TenantId:       d.Get("tenant_id").(string),
		Environment:    d.Get("environment").(string),
	}

	log.Printf("[INFO] Creating Credential (AZURE-VNET)")
	id, err := client.CreateCredential(d.Get("name").(string), alkira.CredentialTypeAzureVnet, c, 0)

	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(id)
	return resourceCredentialAzureVnetRead(ctx, d, meta)
}

func resourceCredentialAzureVnetRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	return nil
}

func resourceCredentialAzureVnetUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*alkira.AlkiraClient)

	c := alkira.CredentialAzureVnet{
		ApplicationId:  d.Get("application_id").(string),
		SecretKey:      d.Get("secret_key").(string),
		SubscriptionId: d.Get("subscription_id").(string),
		TenantId:       d.Get("tenant_id").(string),
		Environment:    d.Get("environment").(string),
	}

	log.Printf("[INFO] Updating Credential (AZURE-VNET)")
	err := client.UpdateCredential(d.Id(), d.Get("name").(string), alkira.CredentialTypeAzureVnet, c, 0)

	if err != nil {
		return diag.FromErr(err)
	}

	return resourceCredentialAzureVnetRead(ctx, d, meta)
}

func resourceCredentialAzureVnetDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*alkira.AlkiraClient)

	err := client.DeleteCredential(d.Id(), alkira.CredentialTypeAzureVnet)

	if err != nil {
		// Terraform may not print "with <resource address>" for destroys of objects
		// that are no longer in configuration, so include identifying context here.
		name, _ := d.GetOk("name")
		if nameStr, ok := name.(string); ok && nameStr != "" {
			return diag.FromErr(fmt.Errorf("%w alkira_credential_azure_vnet (name=%q id=%s)", err, nameStr, d.Id()))
		}
		return diag.FromErr(fmt.Errorf("%w alkira_credential_azure_vnet (id=%s)", err, d.Id()))
	}

	d.SetId("")
	return nil
}
