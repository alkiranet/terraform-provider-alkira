package alkira

import (
	"log"

	"github.com/alkiranet/alkira-client-go/alkira"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceAlkiraCredentialGcpVpc() *schema.Resource {
	return &schema.Resource{
		Description: "Manage Credential for GCP.",
		Create:      resourceCredentialGcpVpc,
		Read:        resourceCredentialGcpVpcRead,
		Update:      resourceCredentialGcpVpcUpdate,
		Delete:      resourceCredentialGcpVpcDelete,

		Schema: map[string]*schema.Schema{
			"name": &schema.Schema{
				Description: "The name of the credential",
				Type:        schema.TypeString,
				Required:    true,
			},
			"auth_provider": &schema.Schema{
				Description: "GCP Authentication Provider",
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "https://www.googleapis.com/oauth2/v1/certs",
			},
			"auth_uri": &schema.Schema{
				Description: "GCP Authentication URI",
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "https://accounts.google.com/o/oauth2/auth",
			},
			"client_email": &schema.Schema{
				Description: "GCP Client email",
				Type:        schema.TypeString,
				Required:    true,
			},
			"client_id": &schema.Schema{
				Description: "GCP Client ID",
				Type:        schema.TypeString,
				Required:    true,
			},
			"client_x509_cert_url": &schema.Schema{
				Description: "GCP Client X509 Cert URL",
				Type:        schema.TypeString,
				Required:    true,
			},
			"private_key_id": &schema.Schema{
				Description: "GCP Private Key ID",
				Type:        schema.TypeString,
				Required:    true,
			},
			"private_key": &schema.Schema{
				Description: "GCP Private Key",
				Type:        schema.TypeString,
				Required:    true,
			},
			"project_id": &schema.Schema{
				Description: "GCP Project ID",
				Type:        schema.TypeString,
				Required:    true,
			},
			"token_uri": &schema.Schema{
				Description: "Token URI",
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "https://oauth2.googleapis.com/token",
			},
			"type": &schema.Schema{
				Description: "GCP Auth Type, default value is `service_account`.",
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "service_account",
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
