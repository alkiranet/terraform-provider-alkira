package alkira

import (
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
		Description: "There are two options for providing the required license key for Fortinet " +
			"instance credentials. You can either input the value directly into the `license_key` field " +
			"or provide the file path for the license key file using the `license_key_file_path`. " +
			"Either `license_key` and `license_key_file_path` must have an input. If both are provided, " +
			"the Alkira provider will treat the `license_key` field with precedence. \n\n\n " +
			"You may also use terraform's built in `file` helper function as a literal input for " +
			"`license_key`. Ex: `license_key = file('/path/to/license/file')`.",
		DeprecationMessage: "alkira_credential_fortinet_instance has been deprecated. " +
			"Please specify `license_key` or `license_key_file_path` directly in resource alkira_service_fortinet. " +
			"See documentation for example.",

		Schema: map[string]*schema.Schema{
			"name": &schema.Schema{
				Description: "The name of the credential.",
				Type:        schema.TypeString,
				Required:    true,
			},
			"license_key": &schema.Schema{
				Description: "Fortinet license key. Interpreted by the Alkira provider as a literal input.",
				Type:        schema.TypeString,
				Optional:    true,
			},
			//"license_key_is_path": &schema.Schema{
			"license_key_file_path": &schema.Schema{
				Description: "Fortinet license key file path. The path to the desired license key. " +
					"`license_key_file_path` will be if both `license_key` and `license_key_file_path` " +
					"are provided in your configuration file. ",
				Type:     schema.TypeString,
				Optional: true,
			},
			"license_type": {
				Description:  "Fortinet instance license type, either `BRING_YOUR_OWN` or `PAY_AS_YOU_GO`.",
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringInSlice([]string{"BRING_YOUR_OWN", "PAY_AS_YOU_GO"}, false),
				Deprecated:   "Not supported anymore. Set `license_type` in alkira_service_fortinet",
			},
		},
	}
}

func resourceCredentialFortinetInstance(d *schema.ResourceData, meta interface{}) error {
	return nil
}

func resourceCredentialFortinetInstanceRead(d *schema.ResourceData, meta interface{}) error {
	return nil
}

func resourceCredentialFortinetInstanceUpdate(d *schema.ResourceData, meta interface{}) error {
	return nil
}

func resourceCredentialFortinetInstanceDelete(d *schema.ResourceData, meta interface{}) error {
	return nil
}
