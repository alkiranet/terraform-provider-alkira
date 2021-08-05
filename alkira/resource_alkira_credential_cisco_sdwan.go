package alkira

import (
	"log"

	"github.com/alkiranet/alkira-client-go/alkira"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceAlkiraCredentialCiscoSdwan() *schema.Resource {
	return &schema.Resource{
		Create: resourceCredentialCiscoSdwan,
		Read:   resourceCredentialCiscoSdwanRead,
		Update: resourceCredentialCiscoSdwanUpdate,
		Delete: resourceCredentialCiscoSdwanDelete,

		Schema: map[string]*schema.Schema{
			"name": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				Description: "The name of the credential.",
			},
			"password": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				Description: "Password of Cisco SD-WAN.",
			},
			"username": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				Description: "Username of Cisco SD-WAN.",
			},
		},
	}
}

func resourceCredentialCiscoSdwan(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*alkira.AlkiraClient)

	c := alkira.CredentialCiscoSdwan{
		Password: d.Get("password").(string),
		Username: d.Get("username").(string),
	}

	log.Printf("[INFO] Creating Credential (ciscosdwan)")
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
	client := meta.(*alkira.AlkiraClient)

	c := alkira.CredentialCiscoSdwan{
		Password: d.Get("password").(string),
		Username: d.Get("username").(string),
	}

	log.Printf("[INFO] Updating Credential (ciscosdwan)")
	err := client.UpdateCredential(d.Id(), d.Get("name").(string), "ciscosdwan", c)

	if err != nil {
		return err
	}

	return resourceCredentialCiscoSdwanRead(d, meta)
}

func resourceCredentialCiscoSdwanDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*alkira.AlkiraClient)
	credentialId := d.Id()

	log.Printf("[INFO] Deleting Credential (ciscosdwan %s)\n", credentialId)
	err := client.DeleteCredential(credentialId, "ciscosdwan")

	if err != nil {
		log.Printf("[INFO] Credential (ciscosdwan %s) was already deleted\n", credentialId)
	}

	return nil
}
