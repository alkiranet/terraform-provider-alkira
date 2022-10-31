// Copyright (C) 2022 Alkira Inc. All Rights Reserved.
package alkira

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceAlkiraCredentialCheckpoint() *schema.Resource {
	return &schema.Resource{
		Create: resourceCredentialCheckpoint,
		Read:   resourceCredentialCheckpointRead,
		Update: resourceCredentialCheckpointUpdate,
		Delete: resourceCredentialCheckpointDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		DeprecationMessage: "alkira_credential_checkpoint has been deprecated. " +
			"Please specify name, password, management_server_password and sic_keys " +
			"directly in resource service_checkpoint. See documentation for example.",

		Schema: map[string]*schema.Schema{
			"name": &schema.Schema{
				Description: "The name of the credential.",
				Type:        schema.TypeString,
				Required:    true,
			},
			"password": &schema.Schema{
				Description: "The Checkpoint Firewall service password.",
				Type:        schema.TypeString,
				Required:    true,
			},
			"management_server_password": &schema.Schema{
				Description: "The password for Checkpoint Firewall Management Server. ",
				Type:        schema.TypeString,
				Required:    true,
			},
			"sic_keys": &schema.Schema{
				Description: "The checkpoint instance sic keys.",
				Type:        schema.TypeList,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Required:    true,
			},
		},
	}
}

func resourceCredentialCheckpoint(d *schema.ResourceData, meta interface{}) error {
	return nil
}

func resourceCredentialCheckpointRead(d *schema.ResourceData, meta interface{}) error {
	return nil
}

func resourceCredentialCheckpointUpdate(d *schema.ResourceData, meta interface{}) error {
	return nil
}

func resourceCredentialCheckpointDelete(d *schema.ResourceData, meta interface{}) error {
	return nil
}
