package alkira

import (
	"log"

	"github.com/alkiranet/alkira-client-go/alkira"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func resourceAlkiraCredentialFortinetInstance() *schema.Resource {
	return &schema.Resource{
		Create: resourceCredentialFortinetInstance,
		Read:   resourceCredentialFortinetInstanceRead,
		Update: resourceCredentialFortinetInstanceUpdate,
		Delete: resourceCredentialFortinetInstanceDelete,

		Schema: map[string]*schema.Schema{
			"name": &schema.Schema{
				Description: "The name of the credential.",
				Type:        schema.TypeString,
				Required:    true,
			},
			"license_key": &schema.Schema{
				Description: "Fortinet license key.",
				Type:        schema.TypeString,
				Required:    true,
			},
			"license_type": {
				Description:  "Fortinet instance license type, either `BRING_YOUR_OWN` or `PAY_AS_YOU_GO`.",
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringInSlice([]string{"BRING_YOUR_OWN", "PAY_AS_YOU_GO"}, false),
			},
		},
	}
}

func resourceCredentialFortinetInstance(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*alkira.AlkiraClient)

	c := alkira.CredentialFortinetInstance{
		LicenseKey:  d.Get("license_key").(string),
		LicenseType: d.Get("license_type").(string),
	}

	log.Printf("[INFO] Creating Credential (Fortinet Instance)")
	credentialId, err := client.CreateCredential(d.Get("name").(string), alkira.CredentialTypeFortinetInstance, c)

	if err != nil {
		return err
	}

	d.SetId(credentialId)
	return resourceCredentialFortinetInstanceRead(d, meta)
}

func resourceCredentialFortinetInstanceRead(d *schema.ResourceData, meta interface{}) error {
	return nil
}

func resourceCredentialFortinetInstanceUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*alkira.AlkiraClient)

	c := alkira.CredentialFortinetInstance{
		LicenseKey:  d.Get("license_key").(string),
		LicenseType: d.Get("license_type").(string),
	}

	log.Printf("[INFO] Updating Credential (Fortinet Instance)")
	err := client.UpdateCredential(d.Id(), d.Get("name").(string), alkira.CredentialTypeFortinetInstance, c)

	if err != nil {
		return err
	}

	return resourceCredentialFortinetInstanceRead(d, meta)
}

func resourceCredentialFortinetInstanceDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*alkira.AlkiraClient)
	credentialId := d.Id()

	log.Printf("[INFO] Deleting Credential (Fortinet Instance %s)\n", credentialId)
	err := client.DeleteCredential(credentialId, alkira.CredentialTypeFortinetInstance)

	if err != nil {
		log.Printf("[INFO] Credential (Fortinet Instance %s) was already deleted\n", credentialId)
	}

	return nil
}
