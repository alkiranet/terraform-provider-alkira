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
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"name": &schema.Schema{
				Description: "The name of the credential.",
				Type:        schema.TypeString,
				Required:    true,
			},
			"auth_key": &schema.Schema{
				Description: "PAN instance auth key.",
				Type:        schema.TypeString,
				Required:    true,
			},
			"auth_code": &schema.Schema{
				Description: "PAN instance auth code.",
				Type:        schema.TypeString,
				Required:    true,
			},
			"license_key": &schema.Schema{
				Description: "PAN license key.",
				Type:        schema.TypeString,
				Required:    true,
			},
			"password": &schema.Schema{
				Description: "PAN password.",
				Type:        schema.TypeString,
				Required:    true,
			},
			"username": &schema.Schema{
				Description: "PAN username.",
				Type:        schema.TypeString,
				Required:    true,
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
	credentialId, err := client.CreateCredential(d.Get("name").(string), alkira.CredentialTypePanInstance, c, 0)

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
	client := meta.(*alkira.AlkiraClient)

	c := alkira.CredentialPanInstance{
		AuthKey:    d.Get("auth_key").(string),
		AuthCode:   d.Get("auth_code").(string),
		LicenseKey: d.Get("license_key").(string),
		Password:   d.Get("password").(string),
		Username:   d.Get("username").(string),
	}

	log.Printf("[INFO] Updating Credential (PAN Instance)")
	err := client.UpdateCredential(d.Id(), d.Get("name").(string), alkira.CredentialTypePanInstance, c, 0)

	if err != nil {
		return err
	}

	return resourceCredentialPanInstanceRead(d, meta)
}

func resourceCredentialPanInstanceDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*alkira.AlkiraClient)
	credentialId := d.Id()

	log.Printf("[INFO] Deleting Credential (PAN Instance %s)\n", credentialId)
	err := client.DeleteCredential(credentialId, alkira.CredentialTypePanInstance)

	if err != nil {
		log.Printf("[INFO] Credential (PAN Instance %s) was already deleted\n", credentialId)
	}

	return nil
}
