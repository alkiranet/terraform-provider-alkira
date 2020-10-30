package alkira

import (
	"log"

	"github.com/alkiranet/alkira-client-go/alkira"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceAlkiraCredentialPanInstance() *schema.Resource {
	return &schema.Resource{
		Create: resourceCredentialPanInstance,
		Read:   resourceCredentialPanInstanceRead,
		Update: resourceCredentialPanInstanceUpdate,
		Delete: resourceCredentialPanInstanceDelete,

		Schema: map[string]*schema.Schema{
			"name": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				Description: "The name of the credential",
			},
			"auth_key": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				Description: "PAN instance auth key",
			},
			"auth_code": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				Description: "PAN instance auth code",
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

func resourceCredentialPanInstance(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*alkira.AlkiraClient)

	c := alkira.CredentialPanInstance{
		AuthKey:    d.Get("auth_key").(string),
		AuthCode:   d.Get("auth_code").(string),
		LicenseKey: d.Get("license_key").(string),
		Password:   d.Get("password").(string),
		Username:   d.Get("username").(string),
	}

	log.Printf("[INFO] Creating Credential (PAN Instance)")
	credentialId, err := client.CreateCredential(d.Get("name").(string), "paninstance", c)

	if err != nil {
		return err
	}

	d.SetId(credentialId)
	return resourceCredentialPanInstanceRead(d, meta)
}

func resourceCredentialPanInstanceRead(d *schema.ResourceData, meta interface{}) error {
	return nil
}

func resourceCredentialPanInstanceUpdate(d *schema.ResourceData, meta interface{}) error {
	return resourceCredentialPanInstanceRead(d, meta)
}

func resourceCredentialPanInstanceDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*alkira.AlkiraClient)
	credentialId := d.Id()

	log.Printf("[INFO] Deleting Credential (PAN Instance %s)\n", credentialId)
	err := client.DeleteCredential(credentialId, "paninstance")

	if err != nil {
		return err
	}

	return nil
}
