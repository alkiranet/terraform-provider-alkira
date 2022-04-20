// Copyright (C) 2022 Alkira Inc. All Rights Reserved.
package alkira

import (
	"log"

	"github.com/alkiranet/alkira-client-go/alkira"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceAlkiraCredentialFortinet() *schema.Resource {
	return &schema.Resource{
		Create: resourceCredentialFortinet,
		Read:   resourceCredentialFortinetRead,
		Update: resourceCredentialFortinetUpdate,
		Delete: resourceCredentialFortinetDelete,

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
	client := meta.(*alkira.AlkiraClient)

	c := &alkira.CredentialFortinet{
		UserName: d.Get("username").(string),
		Password: d.Get("password").(string),
	}

	log.Printf("[INFO] Creating Credential (Fortinet)")
	credentialId, err := client.CreateCredential(d.Get("name").(string), alkira.CredentialTypeFortinet, c)

	if err != nil {
		return err
	}

	d.SetId(credentialId)
	return resourceCredentialFortinetRead(d, meta)
}

func resourceCredentialFortinetRead(d *schema.ResourceData, meta interface{}) error {
	return nil
}

func resourceCredentialFortinetUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*alkira.AlkiraClient)

	c := &alkira.CredentialFortinet{
		UserName: d.Get("username").(string),
		Password: d.Get("password").(string),
	}

	log.Printf("[INFO] Updating Credential (Fortinet)")
	err := client.UpdateCredential(
		d.Id(),
		d.Get("name").(string),
		alkira.CredentialTypeFortinet,
		c,
	)

	if err != nil {
		return err
	}

	return resourceCredentialFortinetRead(d, meta)
}

func resourceCredentialFortinetDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*alkira.AlkiraClient)
	credentialId := d.Id()

	log.Printf("[INFO] Deleting Credential (Fortinet %s)\n", credentialId)
	err := client.DeleteCredential(credentialId, alkira.CredentialTypeFortinet)

	if err != nil {
		log.Printf("[INFO] Credential (Fortinet %s) was already deleted\n", credentialId)
	}

	return nil
}
