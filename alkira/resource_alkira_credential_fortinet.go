// Copyright (C) 2022 Alkira Inc. All Rights Reserved.
package alkira

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceAlkiraCredentialFortinet() *schema.Resource {
	return &schema.Resource{
		Create: resourceCredentialFortinet,
		Read:   resourceCredentialFortinetRead,
		Update: resourceCredentialFortinetUpdate,
		Delete: resourceCredentialFortinetDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		DeprecationMessage: "alkira_credential_fortinet has been deprecated. " +
			"Please specify `username` and `password` directly in resource alkira_service_fortinet. " +
			"See documentation for example.",

		Schema: map[string]*schema.Schema{
			"name": &schema.Schema{
				Description: "The name of the credential.",
				Type:        schema.TypeString,
				Required:    true,
			},
			"password": &schema.Schema{
				Description: "Fortinet password.",
				Type:        schema.TypeString,
				Required:    true,
			},
			"username": &schema.Schema{
				Description: "Fortinet username.",
				Type:        schema.TypeString,
				Required:    true,
			},
		},
	}
}

func resourceCredentialFortinet(d *schema.ResourceData, meta interface{}) error {
	return nil
}

func resourceCredentialFortinetRead(d *schema.ResourceData, meta interface{}) error {
	return nil
}

func resourceCredentialFortinetUpdate(d *schema.ResourceData, meta interface{}) error {
	return nil
}

func resourceCredentialFortinetDelete(d *schema.ResourceData, meta interface{}) error {
	return nil
}
