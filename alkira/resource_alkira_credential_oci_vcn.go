package alkira

import (
	"context"
	"log"

	"github.com/alkiranet/alkira-client-go/alkira"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceAlkiraCredentialOciVcn() *schema.Resource {
	return &schema.Resource{
		Description: "Credential for accessing Oracle Cloud.\n\n" +
			"You can provide your credentials via the following enviromental " +
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
			"name": &schema.Schema{
				Description: "Name of the credential.",
				Type:        schema.TypeString,
				Required:    true,
			},
			"user_ocid": &schema.Schema{
				Description: "OCID of the user.",
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc(
					"AK_OCI_USER_OCID",
					nil),
			},
			"fingerprint": &schema.Schema{
				Description: "Fingerprint of the API key of the user.",
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc(
					"AK_OCI_FINGERPRINT",
					nil),
			},
			"key": &schema.Schema{
				Description: "API key of the user.",
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc(
					"AK_OCI_KEY",
					nil),
			},
			"tenant_ocid": &schema.Schema{
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

	log.Printf("[INFO] Creating Credential (OCI-VCN)")
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

	log.Printf("[INFO] Updating Credential (OCI-VCN)")
	err := client.UpdateCredential(d.Id(), d.Get("name").(string), alkira.CredentialTypeOciVcn, c, 0)

	if err != nil {
		return diag.FromErr(err)
	}

	return resourceCredentialOciVcnRead(ctx, d, meta)
}

func resourceCredentialOciVcnDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*alkira.AlkiraClient)
	id := d.Id()

	log.Printf("[INFO] Deleting Credential (OCI-VCN) %s", id)
	err := client.DeleteCredential(id, alkira.CredentialTypeOciVcn)

	if err != nil {
		return diag.FromErr(err)
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
