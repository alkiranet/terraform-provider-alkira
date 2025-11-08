package alkira

import (
	"context"

	"github.com/alkiranet/alkira-client-go/alkira"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceAlkiraCredentialOciVcn() *schema.Resource {
	return &schema.Resource{
		Description: "Credential for accessing Oracle Cloud.\n\n" +
			"You can provide your credentials via the following environmental " +
			"variables:\n\n * AK_OCI_USER_OCID\n " +
			"* AK_OCI_FINGERPRINT\n * AK_OCI_KEY\n " +
			"* AK_OCI_TENANT_OCID\n",
		CreateContext: resourceCredentialOciVcn,
		ReadContext:   resourceCredentialOciVcnRead,
		UpdateContext: resourceCredentialOciVcnUpdate,
		DeleteContext: resourceCredentialOciVcnDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Description: "Name of the credential.",
				Type:        schema.TypeString,
				Required:    true,
			},
			"user_ocid": {
				Description: "OCID of the user.",
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc(
					"AK_OCI_USER_OCID",
					nil),
			},
			"fingerprint": {
				Description: "Fingerprint of the API key of the user.",
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc(
					"AK_OCI_FINGERPRINT",
					nil),
			},
			"key": {
				Description: "API key of the user.",
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc(
					"AK_OCI_KEY",
					nil),
			},
			"tenant_ocid": {
				Description: "OCID of the tenant.",
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc(
					"AK_OCI_TENANT_OCID",
					nil),
			},
		},
	}
}

func resourceCredentialOciVcn(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*alkira.AlkiraClient)

	c := generateCredentialOciVcnRequest(d)

	credentialId, err := client.CreateCredential(d.Get("name").(string), alkira.CredentialTypeOciVcn, c, 0)

	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(credentialId)
	return resourceCredentialOciVcnRead(ctx, d, meta)
}

func resourceCredentialOciVcnRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	return nil
}

func resourceCredentialOciVcnUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*alkira.AlkiraClient)

	c := generateCredentialOciVcnRequest(d)

	err := client.UpdateCredential(d.Id(), d.Get("name").(string), alkira.CredentialTypeOciVcn, c, 0)

	if err != nil {
		return diag.FromErr(err)
	}

	return resourceCredentialOciVcnRead(ctx, d, meta)
}

func resourceCredentialOciVcnDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*alkira.AlkiraClient)

	err := client.DeleteCredential(d.Id(), alkira.CredentialTypeOciVcn)

	if err != nil {
		return diag.FromErr(err)
	}

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
