package alkira

import (
	"log"

	"github.com/alkiranet/alkira-client-go/alkira"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceAlkiraCredentialCiscoSdwan() *schema.Resource {
	return &schema.Resource{
		Description: "Manage Cisco SD-WAN credential for authentication.\n\n" +
			"The following methods are supported:\n\n" +
			" - Static credentials\n" +
			" - Environment variables\n\n" +
			"### Static Credentials:\n\n" +
			"Static credentials can be provided by adding an `username`" +
			"and `password` in-line in the CISCO SD-WAN block.\n\n" +
			"### Environment Variables:\n\n" +
			"You can provide your credentials via the `AK_CISCO_SDWAN_USERNAME` and" +
			"`AK_CISCO_SDWAN_PASSWORD`, environment variables, representing your" +
			"Cisco SD-WAN username and password, respectively.",
		Create: resourceCredentialCiscoSdwanCreate,
		Read:   resourceCredentialCiscoSdwanRead,
		Update: resourceCredentialCiscoSdwanUpdate,
		Delete: resourceCredentialCiscoSdwanDelete,

		Schema: map[string]*schema.Schema{
			"name": &schema.Schema{
				Description: "The name of the credential.",
				Type:        schema.TypeString,
				Required:    true,
			},
			"username": &schema.Schema{
				Description: "Cisco SD-WAN username.",
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc(
					"AK_CISCO_SDWAN_USERNAME",
					nil),
			},
			"password": &schema.Schema{
				Description: "Cisco SD-WAN password.",
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc(
					"AK_CISCO_SDWAN_PASSWORD",
					nil),
			},
		},
	}
}

func resourceCredentialCiscoSdwanCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*alkira.AlkiraClient)

	c := alkira.CredentialCiscoSdwan{
		Username: d.Get("username").(string),
		Password: d.Get("password").(string),
	}

	id, err := client.CreateCredential(d.Get("name").(string), alkira.CredentialTypeCiscoSdwan, c, 0)

	if err != nil {
		return err
	}

	d.SetId(id)
	return resourceCredentialCiscoSdwanRead(d, meta)
}

func resourceCredentialCiscoSdwanRead(d *schema.ResourceData, meta interface{}) error {
	return nil
}

func resourceCredentialCiscoSdwanUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*alkira.AlkiraClient)

	c := alkira.CredentialCiscoSdwan{
		Password: d.Get("password").(string),
		Username: d.Get("username").(string),
	}

	log.Printf("[INFO] Updating Credential (ciscosdwan)")
	err := client.UpdateCredential(d.Id(), d.Get("name").(string), alkira.CredentialTypeCiscoSdwan, c, 0)

	if err != nil {
		return err
	}

	return resourceCredentialCiscoSdwanRead(d, meta)
}

func resourceCredentialCiscoSdwanDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*alkira.AlkiraClient)

	log.Printf("[INFO] Deleting credential (Cisco SD-WAN %s)\n", d.Id())
	err := client.DeleteCredential(d.Id(), alkira.CredentialTypeCiscoSdwan)

	return err
}
