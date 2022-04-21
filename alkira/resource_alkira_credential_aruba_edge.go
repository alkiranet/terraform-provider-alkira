package alkira

import (
	"log"

	"github.com/alkiranet/alkira-client-go/alkira"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceAlkiraCredentialArubaEdge() *schema.Resource {
	return &schema.Resource{
		Description: "Credential Management for the Aruba Edge Connector.",
		Create:      resourceCredentialArubaEdge,
		Read:        resourceCredentialArubaEdgeRead,
		Update:      resourceCredentialArubaEdgeUpdate,
		Delete:      resourceCredentialArubaDelete,

		Schema: map[string]*schema.Schema{
			"name": &schema.Schema{
				Description: "The name of the credential",
				Type:        schema.TypeString,
				Required:    true,
			},
			"account_key": &schema.Schema{
				Description: "The account key generated in SilverPeak orchestrator account.",
				Type:        schema.TypeString,
				Required:    true,
			},
		},
	}
}

func resourceCredentialArubaEdge(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*alkira.AlkiraClient)

	c := alkira.CredentialArubaEdgeConnectInstance{
		AccountKey: d.Get("account_key").(string),
	}

	id, err := client.CreateCredential(d.Get("name").(string), alkira.CredentialTypeArubaEdgeConnectInstance, c)
	if err != nil {
		return err
	}

	d.SetId(id)
	return resourceCredentialArubaEdgeRead(d, meta)
}

func resourceCredentialArubaEdgeRead(d *schema.ResourceData, meta interface{}) error {
	return nil
}

func resourceCredentialArubaEdgeUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*alkira.AlkiraClient)

	c := alkira.CredentialArubaEdgeConnectInstance{
		AccountKey: d.Get("account_key").(string),
	}

	log.Printf("[INFO] Updating credential (Aruba Edge) %s", d.Id())
	err := client.UpdateCredential(d.Id(), d.Get("name").(string), alkira.CredentialTypeArubaEdgeConnectInstance, c)

	if err != nil {
		return err
	}

	return resourceCredentialArubaEdgeRead(d, meta)
}

func resourceCredentialArubaDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*alkira.AlkiraClient)
	credentialId := d.Id()

	log.Printf("[INFO] Deleting credential (Aruba Edge %s)\n", credentialId)
	return client.DeleteCredential(credentialId, alkira.CredentialTypeArubaEdgeConnectInstance)
}
