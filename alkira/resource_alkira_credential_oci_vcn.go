package alkira

import (
	"log"

	"github.com/alkiranet/alkira-client-go/alkira"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceAlkiraCredentialOciVcn() *schema.Resource {
	return &schema.Resource{
		Description: "Credential for accessing Oracle Cloud.",
		Create:      resourceCredentialOciVcn,
		Read:        resourceCredentialOciVcnRead,
		Update:      resourceCredentialOciVcnUpdate,
		Delete:      resourceCredentialOciVcnDelete,

		Schema: map[string]*schema.Schema{
			"name": &schema.Schema{
				Description: "The name of the credential.",
				Type:        schema.TypeString,
				Required:    true,
			},
			"user_ocid": &schema.Schema{
				Description: "OCID of the user.",
				Type:        schema.TypeString,
				Required:    true,
			},
			"fingerprint": &schema.Schema{
				Description: "Fingerprint of the API key of the user.",
				Type:        schema.TypeString,
				Required:    true,
			},
			"key": &schema.Schema{
				Description: "API key of the user.",
				Type:        schema.TypeString,
				Required:    true,
			},
			"tenant_ocid": &schema.Schema{
				Description: "OCID of the tenant.",
				Type:        schema.TypeString,
				Required:    true,
			},
		},
	}
}

func resourceCredentialOciVcn(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*alkira.AlkiraClient)

	c := generateCredentialOciVcnRequest(d)

	log.Printf("[INFO] Creating Credential (OCI-VCN)")
	credentialId, err := client.CreateCredential(d.Get("name").(string), alkira.CredentialTypeOciVcn, c)

	if err != nil {
		return err
	}

	d.SetId(credentialId)
	return resourceCredentialOciVcnRead(d, meta)
}

func resourceCredentialOciVcnRead(d *schema.ResourceData, meta interface{}) error {
	return nil
}

func resourceCredentialOciVcnUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*alkira.AlkiraClient)

	c := generateCredentialOciVcnRequest(d)

	log.Printf("[INFO] Updating Credential (OCI-VCN)")
	err := client.UpdateCredential(d.Id(), d.Get("name").(string), alkira.CredentialTypeOciVcn, c)

	if err != nil {
		return err
	}

	return resourceCredentialOciVcnRead(d, meta)
}

func resourceCredentialOciVcnDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*alkira.AlkiraClient)
	id := d.Id()

	log.Printf("[INFO] Deleting Credential (OCI-VCN) %s", id)
	err := client.DeleteCredential(id, alkira.CredentialTypeOciVcn)

	if err != nil {
		return err
	}

	log.Printf("[INFO] Deleted Credential (OCI-VCN) %s", id)
	d.SetId("")
	return nil
}

func generateCredentialOciVcnRequest(d *schema.ResourceData) alkira.CredentialOciVcn {
	c := alkira.CredentialOciVcn{
		UserId:      d.Get("user_ocid").(string),
		FingerPrint: d.Get("fingerprint").(string),
		Key:         d.Get("key").(string),
		TenantId:    d.Get("tenant_ocid").(string),
	}

	return c
}
