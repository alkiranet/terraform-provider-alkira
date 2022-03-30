// Copyright (C) 2022 Alkira Inc. All Rights Reserved.
package alkira

import (
	"log"

	"github.com/alkiranet/alkira-client-go/alkira"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceAlkiraCredentialCheckpointManagementServer() *schema.Resource {
	return &schema.Resource{
		Create: resourceCredentialCheckpointManagementServer,
		Read:   resourceCredentialCheckpointManagementServerRead,
		Update: resourceCredentialCheckpointManagementServerUpdate,
		Delete: resourceCredentialCheckpointManagementserverDelete,

		Schema: map[string]*schema.Schema{
			"name": &schema.Schema{
				Description: "The name of the credential.",
				Type:        schema.TypeString,
				Required:    true,
			},
			"password": &schema.Schema{
				Description: "Checkpoint management server password.",
				Type:        schema.TypeString,
				Required:    true,
			},
		},
	}
}

func resourceCredentialCheckpointManagementServer(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*alkira.AlkiraClient)

	c := &alkira.CredentialCheckPointFwManagementServer{
		Password: d.Get("password").(string),
	}

	log.Printf("[INFO] Creating Credential (Checkpoint Management Server)")
	credentialId, err := client.CreateCredential(
		d.Get("name").(string),
		alkira.CredentialTypeChkpFwManagement,
		c,
	)

	if err != nil {
		return err
	}

	d.SetId(credentialId)
	return resourceCredentialCheckpointManagementServerRead(d, meta)
}

func resourceCredentialCheckpointManagementServerRead(d *schema.ResourceData, meta interface{}) error {
	return nil
}

func resourceCredentialCheckpointManagementServerUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*alkira.AlkiraClient)

	c := &alkira.CredentialCheckPointFwManagementServer{
		Password: d.Get("password").(string),
	}

	log.Printf("[INFO] Updating Credential (Checkpoint Management Server)")
	err := client.UpdateCredential(
		d.Id(),
		d.Get("name").(string),
		alkira.CredentialTypeChkpFwManagement,
		c,
	)

	if err != nil {
		return err
	}

	return resourceCredentialCheckpointManagementServer(d, meta)
}

func resourceCredentialCheckpointManagementserverDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*alkira.AlkiraClient)
	credentialId := d.Id()

	log.Printf("[INFO] Deleting Credential (Checkpoint Management Server %s)\n", credentialId)
	err := client.DeleteCredential(credentialId, alkira.CredentialTypeChkpFwInstance)

	if err != nil {
		log.Printf("[INFO] Credential (Checkpoint Management Server %s) was already deleted\n", credentialId)
	}

	return nil
}
