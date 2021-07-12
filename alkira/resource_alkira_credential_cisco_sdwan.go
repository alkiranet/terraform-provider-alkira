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
			"*** Static Credentials:\n\n" +
			"Static credentials can be provided by adding an `username`" +
			"and `password` in-line in the CISCO SD-WAN block.\n\n" +
			"*** Environment Variables:\n\n" +
			"You can provide your credentials via the `CISCO_SDWAN_USERNAME` and" +
			"`CISCO_SDWAN_PASSWORD`, environment variables, representing your" +
			"Cisco SD-WAN username and password, respectively.",
		Create: resourceCredentialCiscoSdwanCreate,
		Read:   resourceCredentialCiscoSdwanRead,
		Update: resourceCredentialCiscoSdwanUpdate,
		Delete: resourceCredentialCiscoSdwanDelete,

		Schema: map[string]*schema.Schema{
			"name": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				Description: "The name of the credential",
			},
			"username": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				DefaultFunc: schema.EnvDefaultFunc(
					"CISCO_SDWAN_USERNAME",
					nil),
				Description: "Cisco SD-WAN username",
			},
			"password": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				DefaultFunc: schema.EnvDefaultFunc(
					"CISCO_SDWAN_PASSWORD",
					nil),
				Description: "Cisco SD-WAN password",
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

	id, err := client.CreateCredential(d.Get("name").(string), "ciscosdwan", c)

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
	return resourceCredentialCiscoSdwanRead(d, meta)
}

func resourceCredentialCiscoSdwanDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*alkira.AlkiraClient)

	log.Printf("[INFO] Deleting credential (Cisco SD-WAN %s)\n", d.Id())
	err := client.DeleteCredential(d.Id(), "ciscosdwan")

	if err != nil {
		return err
	}

	return nil
}
