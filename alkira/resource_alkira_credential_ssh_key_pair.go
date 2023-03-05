package alkira

import (
	"context"
	"log"

	"github.com/alkiranet/alkira-client-go/alkira"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceAlkiraCredentialSshKeyPair() *schema.Resource {
	return &schema.Resource{
		Description:   "Provides SSH Key Pair credential resource.",
		CreateContext: resourceCredentialSshKeyPairCreate,
		ReadContext:   resourceCredentialSshKeyPairRead,
		UpdateContext: resourceCredentialSshKeyPairUpdate,
		DeleteContext: resourceCredentialSshKeyPairDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
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

func resourceCredentialSshKeyPairCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*alkira.AlkiraClient)

	c := alkira.CredentialKeyPair{
		PublicKey: d.Get("public_key").(string),
		Type:      "IMPORTED",
	}

	id, err := client.CreateCredential(d.Get("name").(string), alkira.CredentialTypeKeyPair, c, 0)

	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(id)
	return resourceCredentialSshKeyPairRead(ctx, d, meta)
}

func resourceCredentialSshKeyPairRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	return nil
}

func resourceCredentialSshKeyPairUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*alkira.AlkiraClient)

	c := alkira.CredentialKeyPair{
		PublicKey: d.Get("public_key").(string),
		Type:      "IMPORTED",
	}

	log.Printf("[INFO] Updating Credential (SSH key pair)")
	err := client.UpdateCredential(d.Id(), d.Get("name").(string), alkira.CredentialTypeKeyPair, c, 0)

	if err != nil {
		return diag.FromErr(err)
	}

	return resourceCredentialSshKeyPairRead(ctx, d, meta)
}

func resourceCredentialSshKeyPairDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*alkira.AlkiraClient)

	log.Printf("[INFO] Deleting credential (SSH key pair %s)\n", d.Id())
	err := client.DeleteCredential(d.Id(), alkira.CredentialTypeKeyPair)

	return diag.FromErr(err)
}
