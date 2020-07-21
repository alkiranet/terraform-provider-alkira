package alkira

import (
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/alkiranet/terraform-provider-alkira/alkira/internal"
)

func resourceAlkiraCredentialAzureVnet() *schema.Resource {
	return &schema.Resource{
		Create: resourceCredentialAzureVnet,
		Read:   resourceCredentialAzureVnetRead,
		Update: resourceCredentialAzureVnetUpdate,
		Delete: resourceCredentialAzureVnetDelete,

		Schema: map[string]*schema.Schema{
			"name": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				Description: "The name of the credential",
			},
			"application_id": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				Description: "The Application ID",
			},
			"subscription_id": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				Description: "The subscription ID",
			},
			"secret_key": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				Description: "The Secret Key",
			},
			"tenant_id": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				Description: "The Tenant ID",
			},
		},
	}
}

func resourceCredentialAzureVnet(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*internal.AlkiraClient)

	log.Printf("[INFO] Creating credential-azure-vnet")
	id, err := client.CreateCredentialAzureVnet(
		d.Get("name").(string),
		d.Get("application_id").(string),
		d.Get("secret_key").(string),
		d.Get("subscription_id").(string),
		d.Get("tenant_id").(string))

	if err != nil {
		return err
	}

	d.SetId(id)
	log.Printf("[INFO] Created credential-azure-vnet")
	return resourceCredentialAzureVnetRead(d, meta)
}

func resourceCredentialAzureVnetRead(d *schema.ResourceData, meta interface{}) error {
        return nil
}

func resourceCredentialAzureVnetUpdate(d *schema.ResourceData, meta interface{}) error {
        return resourceCredentialAzureVnetRead(d, meta)
}

func resourceCredentialAzureVnetDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*internal.AlkiraClient)
	id     := d.Id()

	log.Printf("[INFO] Deleting credential-azure-vnet %s\n", id)
	err := client.DeleteCredentialAzureVnet(id)

	if err != nil {
		return err
	}

	log.Printf("[INFO] Deleted credential-azure-vnet %s\n", id)
	d.SetId("")
	return nil
}
