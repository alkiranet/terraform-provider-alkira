package alkira

import (
	"errors"
	"fmt"
	"log"
	"os"

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
				Description: "Fortinet license key is treated as a file path by default. `license_key` " +
					"can also be literal file contents but `license_key_is_path` must be set to false in " +
					"this instance.",
				Type:     schema.TypeString,
				Optional: true,
			},
			"license_key_is_path": &schema.Schema{
				Description: "Indicates to the Alkira provider if the `license_key` should " +
					"be treated as a file path or as literal file contents. Default behavior is to" +
					"treat `license_key` as a path to an existing `.lic` file. If it makes more sense " +
					"to enter the contents of the file directly you may use either a data source or" +
					"the built in terraform function `file` https://www.terraform.io/language/functions/file" +
					"This field is included as a convenience.",
				Type:     schema.TypeBool,
				Optional: true,
				Default:  true,
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

	licenseKey := d.Get("license_key").(string)
	licenseKeyIsPath := d.Get("license_key_is_path").(bool)

	licenseKey, err := setLicenseKey(licenseKeyIsPath, licenseKey)
	if err != nil {
		return err
	}

	c := alkira.CredentialFortinetInstance{
		LicenseKey:  licenseKey,
		LicenseType: d.Get("license_type").(string),
	}

	log.Printf("[INFO] Creating Credential (Fortinet Instance)")
	credentialId, err := client.CreateCredential(d.Get("name").(string), alkira.CredentialTypeFortinetInstance, c, 0)

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

	licenseKey := d.Get("license_key").(string)
	licenseKeyIsPath := d.Get("license_key_is_path").(bool)

	licenseKey, err := setLicenseKey(licenseKeyIsPath, licenseKey)
	if err != nil {
		return err
	}

	c := alkira.CredentialFortinetInstance{
		LicenseKey:  licenseKey,
		LicenseType: d.Get("license_type").(string),
	}

	log.Printf("[INFO] Updating Credential (Fortinet Instance)")
	err = client.UpdateCredential(d.Id(), d.Get("name").(string), alkira.CredentialTypeFortinetInstance, c, 0)

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

// setLicenseKey takes ifFilePath which indicates if the licenseKey should be treated as a path or as a
// literal file contents. If licenseKey is a file path setLicenseKey will read the file contents
// into a string and return they value otherwise no alteration is made to the string.
func setLicenseKey(ifFilePath bool, licenseKey string) (string, error) {
	if ifFilePath {
		if _, err := os.Stat(licenseKey); errors.Is(err, os.ErrNotExist) {
			return "", fmt.Errorf("file not found at %s: %w", licenseKey, err)
		}

		b, err := os.ReadFile(licenseKey)
		if err != nil {
			return "", err
		}

		return string(b), nil
	}

	return licenseKey, nil
}
