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
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"name": &schema.Schema{
				Description: "The name of the credential.",
				Type:        schema.TypeString,
				Required:    true,
			},
			"license_key": &schema.Schema{
				Description: "PAN license key.",
				Type:        schema.TypeString,
				Optional:    true,
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

func resourceCredentialPan(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*alkira.AlkiraClient)

	c := alkira.CredentialPan{
		LicenseKey: d.Get("license_key").(string),
		Password:   d.Get("password").(string),
		Username:   d.Get("username").(string),
	}

	log.Printf("[INFO] Creating Credential (PAN)")
	credentialId, err := client.CreateCredential(d.Get("name").(string), alkira.CredentialTypePan, c, 0)

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
	client := meta.(*alkira.AlkiraClient)

	c := alkira.CredentialPan{
		LicenseKey: d.Get("license_key").(string),
		Password:   d.Get("password").(string),
		Username:   d.Get("username").(string),
	}

	log.Printf("[INFO] Updating Credential (PAN)")
	err := client.UpdateCredential(d.Id(), d.Get("name").(string), alkira.CredentialTypePan, c, 0)

	if err != nil {
		return err
	}

	return resourceCredentialPanRead(d, meta)
}

func resourceCredentialPanDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*alkira.AlkiraClient)
	credentialId := d.Id()

	log.Printf("[INFO] Deleting Credential (PAN %s)\n", credentialId)
	err := client.DeleteCredential(credentialId, alkira.CredentialTypePan)

	if err != nil {
		log.Printf("[INFO] Credential (PAN %s) was already deleted\n", credentialId)
	}

	return nil
}
