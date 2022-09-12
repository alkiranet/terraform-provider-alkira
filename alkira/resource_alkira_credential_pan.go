package alkira

import (
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
		DeprecationMessage: "alkira_credential_pan has been deprecated. Please specify pan_username and pan_password directly in resource service_pan. See documentation for example.",

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
	return nil
}

func resourceCredentialPanRead(d *schema.ResourceData, meta interface{}) error {
	return nil
}

func resourceCredentialPanUpdate(d *schema.ResourceData, meta interface{}) error {
	return nil
}

func resourceCredentialPanDelete(d *schema.ResourceData, meta interface{}) error {
	return nil
}
