// Copyright (C) 2022 Alkira Inc. All Rights Reserved.
package alkira

import (
	"log"

	"github.com/alkiranet/alkira-client-go/alkira"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceAlkiraCredentialCheckpointInstance() *schema.Resource {
	return &schema.Resource{
		Create: resourceCredentialCheckpointInstance,
		Read:   resourceCredentialCheckpointInstanceRead,
		Update: resourceCredentialCheckpointInstanceUpdate,
		Delete: resourceCredentialCheckpointInstanceDelete,

		Schema: map[string]*schema.Schema{
			"name": &schema.Schema{
				Description: "The name of the credential.",
				Type:        schema.TypeString,
				Required:    true,
			},
			"sic_key": &schema.Schema{
				Description: "CheckpointInstance sic key.",
				Type:        schema.TypeString,
				Required:    true,
			},
		},
	}
}

func resourceCredentialCheckpointInstance(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*alkira.AlkiraClient)

	c := &alkira.CredentialCheckPointFwServiceInstance{
		SicKey: d.Get("sic_key").(string),
	}

	log.Printf("[INFO] Creating Credential (CheckpointInstance)")
	credentialId, err := client.CreateCredential(
		d.Get("name").(string),
		alkira.CredentialTypeChkpFwInstance,
		c,
	)

	if err != nil {
		return err
	}

	d.SetId(credentialId)
	return resourceCredentialCheckpointInstanceRead(d, meta)
}

func resourceCredentialCheckpointInstanceRead(d *schema.ResourceData, meta interface{}) error {
	return nil
}

func resourceCredentialCheckpointInstanceUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*alkira.AlkiraClient)

	c := &alkira.CredentialCheckPointFwServiceInstance{
		SicKey: d.Get("sic_key").(string),
	}

	log.Printf("[INFO] Updating Credential (CheckpointInstance)")
	err := client.UpdateCredential(
		d.Id(),
		d.Get("name").(string),
		alkira.CredentialTypeChkpFwInstance,
		c,
	)

	if err != nil {
		return err
	}

	return resourceCredentialCheckpointInstanceRead(d, meta)
}

func resourceCredentialCheckpointInstanceDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*alkira.AlkiraClient)
	credentialId := d.Id()

	log.Printf("[INFO] Deleting Credential (CheckpointInstance %s)\n", credentialId)
	err := client.DeleteCredential(credentialId, alkira.CredentialTypeChkpFwInstance)

	if err != nil {
		log.Printf("[INFO] Credential (CheckpointInstance %s) was already deleted\n", credentialId)
	}

	return nil
}
