// Copyright (C) 2022 Alkira Inc. All Rights Reserved.
package alkira

import (
	"log"

	"github.com/alkiranet/alkira-client-go/alkira"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceAlkiraCredentialCheckpoint() *schema.Resource {
	return &schema.Resource{
		Create: resourceCredentialCheckpoint,
		Read:   resourceCredentialCheckpointRead,
		Update: resourceCredentialCheckpointUpdate,
		Delete: resourceCredentialCheckpointDelete,

		Schema: map[string]*schema.Schema{
			"name": &schema.Schema{
				Description: "The name of the credential.",
				Type:        schema.TypeString,
				Required:    true,
			},
			"password": &schema.Schema{
				Description: "The checkpoint credential password.",
				Type:        schema.TypeString,
				Required:    true,
			},
		},
	}
}

func resourceCredentialCheckpoint(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*alkira.AlkiraClient)

	c := &alkira.CredentialCheckPointFwService{
		AdminPassword: d.Get("password").(string),
	}

	log.Printf("[INFO] Creating Credential (Checkpoint)")
	credentialId, err := client.CreateCredential(
		d.Get("name").(string),
		alkira.CredentialTypeChkpFw,
		c,
	)

	if err != nil {
		return err
	}

	d.SetId(credentialId)
	return resourceCredentialCheckpointRead(d, meta)
}

func resourceCredentialCheckpointRead(d *schema.ResourceData, meta interface{}) error {
	return nil
}

func resourceCredentialCheckpointUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*alkira.AlkiraClient)

	c := &alkira.CredentialCheckPointFwService{
		AdminPassword: d.Get("password").(string),
	}

	log.Printf("[INFO] Updating Credential (Checkpoint)")
	err := client.UpdateCredential(
		d.Id(),
		d.Get("name").(string),
		alkira.CredentialTypeChkpFw,
		c,
	)

	if err != nil {
		return err
	}

	return resourceCredentialCheckpointRead(d, meta)
}

func resourceCredentialCheckpointDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*alkira.AlkiraClient)
	credentialId := d.Id()

	log.Printf("[INFO] Deleting Credential (Checkpoint %s)\n", credentialId)
	err := client.DeleteCredential(credentialId, alkira.CredentialTypeChkpFw)

	if err != nil {
		log.Printf("[INFO] Credential (Checkpoint %s) was already deleted\n", credentialId)
	}

	return nil
}
