package alkira

import (
	"log"

	"github.com/alkiranet/alkira-client-go/alkira"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceAlkiraCredentialAzureVnet() *schema.Resource {
	return &schema.Resource{
		Description: "Manage AWS credential for authentication.\n\n" +
			"You can provide your credentials via the `AK_AZURE_APPLICATION_ID` " +
			"`AK_AZURE_SUBSCRIPTION_ID`, `AK_AZURE_SECRET_KEY` and " +
			"`AK_AZURE_TENANT_ID` environment variables.",
		Create: resourceCredentialAzureVnet,
		Read:   resourceCredentialAzureVnetRead,
		Update: resourceCredentialAzureVnetUpdate,
		Delete: resourceCredentialAzureVnetDelete,

		Schema: map[string]*schema.Schema{
			"name": &schema.Schema{
				Description: "The name of the credential",
				Type:        schema.TypeString,
				Required:    true,
			},
			"application_id": &schema.Schema{
				Description: "The Application ID",
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc(
					"AK_AZURE_APPLICATION_ID",
					nil),
			},
			"subscription_id": &schema.Schema{
				Description: "The subscription ID",
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc(
					"AK_AZURE_SUBSCRIPTION_ID",
					nil),
			},
			"secret_key": &schema.Schema{
				Description: "The Secret Key",
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc(
					"AK_AZURE_SECRET_KEY",
					nil),
			},
			"tenant_id": &schema.Schema{
				Description: "The Tenant ID",
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc(
					"AK_AZURE_TENANT_ID",
					nil),
			},
		},
	}
}

func resourceCredentialAzureVnet(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*alkira.AlkiraClient)

	c := alkira.CredentialAzureVnet{
		ApplicationId:  d.Get("application_id").(string),
		SecretKey:      d.Get("secret_key").(string),
		SubscriptionId: d.Get("subscription_id").(string),
		TenantId:       d.Get("tenant_id").(string),
	}

	log.Printf("[INFO] Creating Credential (AZURE-VNET)")
	id, err := client.CreateCredential(d.Get("name").(string), alkira.CredentialTypeAzureVnet, c)

	if err != nil {
		return err
	}

	d.SetId(id)
	return resourceCredentialAzureVnetRead(d, meta)
}

func resourceCredentialAzureVnetRead(d *schema.ResourceData, meta interface{}) error {
	return nil
}

func resourceCredentialAzureVnetUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*alkira.AlkiraClient)

	c := alkira.CredentialAzureVnet{
		ApplicationId:  d.Get("application_id").(string),
		SecretKey:      d.Get("secret_key").(string),
		SubscriptionId: d.Get("subscription_id").(string),
		TenantId:       d.Get("tenant_id").(string),
	}

	log.Printf("[INFO] Updating Credential (AZURE-VNET)")
	err := client.UpdateCredential(d.Id(), d.Get("name").(string), alkira.CredentialTypeAzureVnet, c)

	if err != nil {
		return err
	}

	return resourceCredentialAzureVnetRead(d, meta)
}

func resourceCredentialAzureVnetDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*alkira.AlkiraClient)
	id := d.Id()

	log.Printf("[INFO] Deleting Credential (AZURE-VNET %s)\n", id)
	err := client.DeleteCredential(id, alkira.CredentialTypeAzureVnet)

	if err != nil {
		return err
	}

	log.Printf("[INFO] Deleted Credential (AZURE-VNET %s)\n", id)
	d.SetId("")
	return nil
}
