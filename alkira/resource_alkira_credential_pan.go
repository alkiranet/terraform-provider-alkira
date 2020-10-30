package alkira

import (
	"log"

	"github.com/alkiranet/alkira-client-go/alkira"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceAlkiraCredentialPan() *schema.Resource {
	return &schema.Resource{
		Create: resourceCredentialPan,
		Read:   resourceCredentialPanRead,
		Update: resourceCredentialPanUpdate,
		Delete: resourceCredentialPanDelete,

		Schema: map[string]*schema.Schema{
			"name": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				Description: "The name of the credential",
			},
			"license_key": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				Description: "PAN license key",
			},
			"password": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				Description: "PAN password",
			},
			"username": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				Description: "PAN username",
			},
		},
	}
}

func resourceCredentialPan(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*alkira.AlkiraClient)

	c := alkira.CredentialPan{
		LicenseKey: d.Get("license_key").(string),
		Password:   d.Get("password").(string),
		Username:   d.Get("username").(string),
	}

	log.Printf("[INFO] Createing Credential (PAN)")
	credentialId, err := client.CreateCredential(d.Get("name").(string), "pan", c)

	if err != nil {
		return err
	}

	d.SetId(credentialId)
	return resourceCredentialPanRead(d, meta)
}

func resourceCredentialPanRead(d *schema.ResourceData, meta interface{}) error {
	return nil
}

func resourceCredentialPanUpdate(d *schema.ResourceData, meta interface{}) error {
	return resourceCredentialPanRead(d, meta)
}

func resourceCredentialPanDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*alkira.AlkiraClient)
	credentialId := d.Id()

	log.Printf("[INFO] Deleting Credential (PAN %s)\n", credentialId)
	err := client.DeleteCredential(credentialId, "pan")

	if err != nil {
		log.Printf("[INFO] Credential (PAN %s) was already deleted\n", credentialId)
	}

	return nil
}
