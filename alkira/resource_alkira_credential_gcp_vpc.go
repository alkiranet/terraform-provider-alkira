package alkira

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/alkiranet/alkira-client-go/alkira"

	"github.com/hashicorp/go-cty/cty"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceAlkiraCredentialGcpVpc() *schema.Resource {
	return &schema.Resource{
		Description:   "Manage Credential for GCP.",
		CreateContext: resourceCredentialGcpVpcCreate,
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
				Sensitive:   true,
				WriteOnly:   true,
			},
			"client_id": {
				Description: "GCP Client ID",
				Type:        schema.TypeString,
				Required:    true,
				Sensitive:   true,
				WriteOnly:   true,
			},
			"client_x509_cert_url": {
				Description: "GCP Client X509 Cert URL",
				Type:        schema.TypeString,
				Required:    true,
				Sensitive:   true,
				WriteOnly:   true,
			},
			"private_key_id": {
				Description: "GCP Private Key ID",
				Type:        schema.TypeString,
				Required:    true,
				Sensitive:   true,
				WriteOnly:   true,
			},
			"private_key": {
				Description: "GCP Private Key",
				Type:        schema.TypeString,
				Required:    true,
				Sensitive:   true,
				WriteOnly:   true,
			},
			"project_id": {
				Description: "GCP Project ID",
				Type:        schema.TypeString,
				Required:    true,
				Sensitive:   true,
				WriteOnly:   true,
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

func resourceCredentialGcpVpcCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*alkira.AlkiraClient)

	c, err := buildCredentialGcpVpc(d)
	if err != nil {
		return diag.FromErr(err)
	}

	log.Printf("[INFO] Creating Credential (GCP-VPC)")
	credentialId, err := client.CreateCredential(d.Get("name").(string), alkira.CredentialTypeGcpVpc, c, 0)

	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(credentialId)
	return resourceCredentialGcpVpcRead(ctx, d, meta)
}

func resourceCredentialGcpVpcRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*alkira.AlkiraClient)

	credential, err := client.GetCredentialById(d.Id())
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			d.SetId("")
			return nil
		}
		return diag.FromErr(err)
	}

	d.Set("name", credential.Name)

	// Note: Sensitive WriteOnly fields (client_email, client_id, client_x509_cert_url,
	// private_key_id, private_key, project_id) are NOT returned by the API for security
	// reasons and must be maintained in the user's HCL configuration.

	return nil
}

func resourceCredentialGcpVpcUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*alkira.AlkiraClient)

	c, err := buildCredentialGcpVpc(d)
	if err != nil {
		return diag.FromErr(err)
	}

	log.Printf("[INFO] Updating Credential (GCP-VPC)")
	err = client.UpdateCredential(d.Id(), d.Get("name").(string), alkira.CredentialTypeGcpVpc, c, 0)

	if err != nil {
		return diag.FromErr(err)
	}

	return resourceCredentialGcpVpcRead(ctx, d, meta)
}

func resourceCredentialGcpVpcDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*alkira.AlkiraClient)

	err := client.DeleteCredential(d.Id(), alkira.CredentialTypeGcpVpc)

	if err != nil {
		// Terraform may not print "with <resource address>" for destroys of objects
		// that are no longer in configuration, so include identifying context here.
		name, _ := d.GetOk("name")
		if nameStr, ok := name.(string); ok && nameStr != "" {
			return diag.FromErr(fmt.Errorf("%w alkira_credential_gcp_vpc (name=%q id=%s)", err, nameStr, d.Id()))
		}
		return diag.FromErr(fmt.Errorf("%w alkira_credential_gcp_vpc (id=%s)", err, d.Id()))
	}

	d.SetId("")
	return nil
}

// buildCredentialGcpVpc constructs the CredentialGcpVpc struct from ResourceData
// For WriteOnly fields, reads from raw config since values are not stored in state
func buildCredentialGcpVpc(d *schema.ResourceData) (alkira.CredentialGcpVpc, error) {
	// Helper function to get WriteOnly field values from raw config
	getWriteOnlyString := func(field string) (string, error) {
		attrPath := cty.Path{cty.GetAttrStep{Name: field}}
		val, diags := d.GetRawConfigAt(attrPath)

		if diags.HasError() {
			return "", fmt.Errorf("error reading %s from config: %v", field, diags)
		}

		if val.IsNull() || !val.IsKnown() {
			return "", fmt.Errorf("required field '%s' is not set in configuration", field)
		}

		if val.Type() != cty.String {
			return "", fmt.Errorf("field '%s' is not a string", field)
		}

		return val.AsString(), nil
	}

	// Get WriteOnly sensitive fields from raw config
	clientEmail, err := getWriteOnlyString("client_email")
	if err != nil {
		return alkira.CredentialGcpVpc{}, err
	}

	clientId, err := getWriteOnlyString("client_id")
	if err != nil {
		return alkira.CredentialGcpVpc{}, err
	}

	clientX509CertUrl, err := getWriteOnlyString("client_x509_cert_url")
	if err != nil {
		return alkira.CredentialGcpVpc{}, err
	}

	privateKeyId, err := getWriteOnlyString("private_key_id")
	if err != nil {
		return alkira.CredentialGcpVpc{}, err
	}

	privateKey, err := getWriteOnlyString("private_key")
	if err != nil {
		return alkira.CredentialGcpVpc{}, err
	}

	projectId, err := getWriteOnlyString("project_id")
	if err != nil {
		return alkira.CredentialGcpVpc{}, err
	}

	return alkira.CredentialGcpVpc{
		AuthProvider:      d.Get("auth_provider").(string),
		AuthUri:           d.Get("auth_uri").(string),
		ClientEmail:       clientEmail,
		ClientId:          clientId,
		ClientX509CertUrl: clientX509CertUrl,
		PrivateKey:        privateKey,
		PrivateKeyId:      privateKeyId,
		ProjectId:         projectId,
		TokenUri:          d.Get("token_uri").(string),
		Type:              d.Get("type").(string),
	}, nil
}
