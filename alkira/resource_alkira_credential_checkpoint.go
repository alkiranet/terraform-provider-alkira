// Copyright (C) 2022 Alkira Inc. All Rights Reserved.
package alkira

import (
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
	client := meta.(*alkira.AlkiraClient)
	name := d.Get("name").(string)
	password := d.Get("password").(string)

	credentialId, err := createCheckpointCredential(name, password, client)
	if err != nil {
		return err
	}
	d.SetId(credentialId)

	sicKeys := convertTypeListToStringList(d.Get("sic_keys").([]interface{}))
	err = createCheckpointCredentialInstances(sicKeys, client)
	if err != nil {
		return err
	}

	err = createCheckpointCredentialManagementServer(name, password, client)
	if err != nil {
		return err
	}

	return resourceCredentialCheckpointRead(d, meta)
}

func resourceCredentialCheckpointRead(d *schema.ResourceData, meta interface{}) error {
	return nil
}

func resourceCredentialCheckpointUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*alkira.AlkiraClient)
	name := d.Get("name").(string)
	password := d.Get("password").(string)

	//update Checkpoint Credential
	err := updateCheckpointCredential(d.Id(), name, password, client)
	if err != nil {
		return err
	}

	//update Checkpoint Credential Management Server
	err = updateCheckpointCredentialManagementServerByName(name, password, client)
	if err != nil {
		return err
	}

	return resourceCredentialCheckpointRead(d, meta)
}

func resourceCredentialCheckpointDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*alkira.AlkiraClient)

	err := deleteCheckpointCredential(d.Id(), client)
	if err != nil {
		return err
	}

	err = deleteCheckpointCredentialInstances(client)
	if err != nil {
		return err
	}

	//NOTE: normally we would check for an error after an attempt to delete but there is somse
	//inconsistency with the API's clean up of the management server credential. When a call to
	//delete the checkpoint service is made it also deletes the management server credential. It
	//does not attempt to delete either the instance credentials or the base checkpoint service
	//credentials. In this case, we simply make a call so we have confidence the resource has been
	//removed. Hopefully we can clean this up in the future.
	deleteCheckpointCredentialManagementServerByName(d.Get("name").(string), client)

	return nil
}
