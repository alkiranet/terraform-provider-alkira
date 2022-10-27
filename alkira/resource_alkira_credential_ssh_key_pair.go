package alkira

import (
	"log"

	"github.com/alkiranet/alkira-client-go/alkira"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceAlkiraCredentialSshKeyPair() *schema.Resource {
	return &schema.Resource{
		Description: "Provides SSH Key Pair credential resource.",
		Create:      resourceCredentialSshKeyPairCreate,
		Read:        resourceCredentialSshKeyPairRead,
		Update:      resourceCredentialSshKeyPairUpdate,
		Delete:      resourceCredentialSshKeyPairDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"name": &schema.Schema{
				Description: "The name of the credential.",
				Type:        schema.TypeString,
				Required:    true,
			},
			"public_key": &schema.Schema{
				Description: "Public key.",
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc(
					"AK_SSH_PUBLIC_KEY",
					nil),
			},
		},
	}
}

func resourceCredentialSshKeyPairCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*alkira.AlkiraClient)

	c := alkira.CredentialKeyPair{
		PublicKey: d.Get("public_key").(string),
		Type:      "IMPORTED",
	}

	id, err := client.CreateCredential(d.Get("name").(string), alkira.CredentialTypeKeyPair, c, 0)

	if err != nil {
		return err
	}

	d.SetId(id)
	return resourceCredentialSshKeyPairRead(d, meta)
}

func resourceCredentialSshKeyPairRead(d *schema.ResourceData, meta interface{}) error {
	return nil
}

func resourceCredentialSshKeyPairUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*alkira.AlkiraClient)

	c := alkira.CredentialKeyPair{
		PublicKey: d.Get("public_key").(string),
		Type:      "IMPORTED",
	}

	log.Printf("[INFO] Updating Credential (SSH key pair)")
	err := client.UpdateCredential(d.Id(), d.Get("name").(string), alkira.CredentialTypeKeyPair, c, 0)

	if err != nil {
		return err
	}

	return resourceCredentialSshKeyPairRead(d, meta)
}

func resourceCredentialSshKeyPairDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*alkira.AlkiraClient)

	log.Printf("[INFO] Deleting credential (SSH key pair %s)\n", d.Id())
	err := client.DeleteCredential(d.Id(), alkira.CredentialTypeKeyPair)

	return err
}
