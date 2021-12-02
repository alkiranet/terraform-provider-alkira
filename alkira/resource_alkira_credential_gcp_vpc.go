package alkira

import (
	"log"

	"github.com/alkiranet/alkira-client-go/alkira"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceAlkiraCredentialGcpVpc() *schema.Resource {
	return &schema.Resource{
		Create: resourceCredentialGcpVpc,
		Read:   resourceCredentialGcpVpcRead,
		Update: resourceCredentialGcpVpcUpdate,
		Delete: resourceCredentialGcpVpcDelete,

		Schema: map[string]*schema.Schema{
			"name": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				Description: "The name of the credential",
			},
			"auth_provider": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				Description: "Authentication Provider",
			},
			"auth_uri": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				Description: "Authentication URI",
			},
			"client_email": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				Description: "Client email",
			},
			"client_id": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				Description: "Client ID",
			},
			"client_x509_cert_url": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				Description: "Client X509 Cert URL",
			},
			"private_key_id": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				Description: "Private Key ID",
			},
			"private_key": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				Description: "Private Key",
			},
			"project_id": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				Description: "Project ID",
			},
			"token_uri": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				Description: "Token URI",
			},
			"type": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				Description: "Type",
			},
		},
	}
}

func resourceCredentialGcpVpc(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*alkira.AlkiraClient)

	c := alkira.CredentialGcpVpc{
		AuthProvider:      d.Get("auth_provider").(string),
		AuthUri:           d.Get("auth_uri").(string),
		ClientEmail:       d.Get("client_email").(string),
		ClientId:          d.Get("client_id").(string),
		ClientX509CertUrl: d.Get("client_x509_cert_url").(string),
		PrivateKey:        d.Get("private_key").(string),
		PrivateKeyId:      d.Get("private_key_id").(string),
		ProjectId:         d.Get("project_id").(string),
		TokenUri:          d.Get("token_uri").(string),
		Type:              d.Get("type").(string),
	}

	log.Printf("[INFO] Creating Credential (GCP-VPC)")
	credentialId, err := client.CreateCredential(d.Get("name").(string), "gcpvpc", c)

	if err != nil {
		return err
	}

	d.SetId(credentialId)
	return resourceCredentialGcpVpcRead(d, meta)
}

func resourceCredentialGcpVpcRead(d *schema.ResourceData, meta interface{}) error {
	return nil
}

func resourceCredentialGcpVpcUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*alkira.AlkiraClient)

	c := alkira.CredentialGcpVpc{
		AuthProvider:      d.Get("auth_provider").(string),
		AuthUri:           d.Get("auth_uri").(string),
		ClientEmail:       d.Get("client_email").(string),
		ClientId:          d.Get("client_id").(string),
		ClientX509CertUrl: d.Get("client_x509_cert_url").(string),
		PrivateKey:        d.Get("private_key").(string),
		PrivateKeyId:      d.Get("private_key_id").(string),
		ProjectId:         d.Get("project_id").(string),
		TokenUri:          d.Get("token_uri").(string),
		Type:              d.Get("type").(string),
	}

	log.Printf("[INFO] Updating Credential (GCP-VPC)")
	err := client.UpdateCredential(d.Id(), d.Get("name").(string), "gcpvpc", c)

	if err != nil {
		return err
	}

	return resourceCredentialGcpVpcRead(d, meta)
}

func resourceCredentialGcpVpcDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*alkira.AlkiraClient)
	id := d.Id()

	log.Printf("[INFO] Deleting Credential (GCP-VPC %s)\n", id)
	err := client.DeleteCredential(id, "gcpvpc")

	if err != nil {
		return err
	}

	log.Printf("[INFO] Deleted Credential (GCP-VPC %s)\n", id)
	d.SetId("")
	return nil
}
