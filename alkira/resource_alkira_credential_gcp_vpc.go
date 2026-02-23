package alkira

import (
	"context"
	"log"

	"github.com/alkiranet/alkira-client-go/alkira"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceAlkiraCredentialGcpVpc() *schema.Resource {
	return &schema.Resource{
		Description:   "Manage Credential for GCP.",
		CreateContext: resourceCredentialGcpVpc,
		ReadContext:   resourceCredentialGcpVpcRead,
		UpdateContext: resourceCredentialGcpVpcUpdate,
		DeleteContext: resourceCredentialGcpVpcDelete,
		Importer: &schema.ResourceImporter{
			StateContext: importWithReadValidation(resourceCredentialGcpVpcRead),
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Description: "The name of the credential",
				Type:        schema.TypeString,
				Required:    true,
			},
			"auth_provider": {
				Description: "GCP Authentication Provider",
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "https://www.googleapis.com/oauth2/v1/certs",
			},
			"auth_uri": {
				Description: "GCP Authentication URI",
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "https://accounts.google.com/o/oauth2/auth",
			},
			"client_email": {
				Description: "GCP Client email",
				Type:        schema.TypeString,
				Required:    true,
			},
			"client_id": {
				Description: "GCP Client ID",
				Type:        schema.TypeString,
				Required:    true,
			},
			"client_x509_cert_url": {
				Description: "GCP Client X509 Cert URL",
				Type:        schema.TypeString,
				Required:    true,
			},
			"private_key_id": {
				Description: "GCP Private Key ID",
				Type:        schema.TypeString,
				Required:    true,
			},
			"private_key": {
				Description: "GCP Private Key",
				Type:        schema.TypeString,
				Required:    true,
			},
			"project_id": {
				Description: "GCP Project ID",
				Type:        schema.TypeString,
				Required:    true,
			},
			"token_uri": {
				Description: "Token URI",
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "https://oauth2.googleapis.com/token",
			},
			"type": {
				Description: "GCP Auth Type, default value is `service_account`.",
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "service_account",
			},
		},
	}
}

func resourceCredentialGcpVpc(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
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
	credentialId, err := client.CreateCredential(d.Get("name").(string), "gcpvpc", c, 0)

	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(credentialId)
	return resourceCredentialGcpVpcRead(ctx, d, meta)
}

func resourceCredentialGcpVpcRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	return nil
}

func resourceCredentialGcpVpcUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
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
	err := client.UpdateCredential(d.Id(), d.Get("name").(string), "gcpvpc", c, 0)

	if err != nil {
		return diag.FromErr(err)
	}

	return resourceCredentialGcpVpcRead(ctx, d, meta)
}

func resourceCredentialGcpVpcDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*alkira.AlkiraClient)

	err := client.DeleteCredential(d.Id(), "gcpvpc")

	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId("")
	return nil
}
