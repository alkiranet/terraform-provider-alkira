package alkira

import (
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/alkiranet/alkira-client-go/alkira"

	"github.com/hashicorp/go-cty/cty"
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
			StateContext: importWithReadValidation(resourceCredentialSshKeyPairRead),
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Description: "The name of the credential.",
				Type:        schema.TypeString,
				Required:    true,
			},
			"public_key": {
				Description: "Public key.",
				Type:        schema.TypeString,
				Optional:    true,
				WriteOnly:   true,
			},
		},
	}
}

func resourceCredentialSshKeyPairCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*alkira.AlkiraClient)

	publicKey, err := getSshKeyPairCredentialValue(d, "public_key", "AK_SSH_PUBLIC_KEY")
	if err != nil {
		return diag.FromErr(err)
	}

	c := alkira.CredentialKeyPair{
		PublicKey: publicKey,
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

	// Note: public_key is NOT returned by the API for security reasons.
	// The getSshKeyPairCredentialValue helper reads from config or AK_SSH_PUBLIC_KEY env var.
	// Private keys are never returned by the API.

	return nil
}

func resourceCredentialSshKeyPairUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*alkira.AlkiraClient)

	publicKey, err := getSshKeyPairCredentialValue(d, "public_key", "AK_SSH_PUBLIC_KEY")
	if err != nil {
		return diag.FromErr(err)
	}

	c := alkira.CredentialKeyPair{
		PublicKey: publicKey,
		Type:      "IMPORTED",
	}

	err = client.UpdateCredential(d.Id(), d.Get("name").(string), alkira.CredentialTypeKeyPair, c, 0)

	if err != nil {
		return diag.FromErr(err)
	}

	return resourceCredentialSshKeyPairRead(ctx, d, meta)
}

func resourceCredentialSshKeyPairDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*alkira.AlkiraClient)

	err := client.DeleteCredential(d.Id(), alkira.CredentialTypeKeyPair)

	if err != nil {
		// Terraform may not print "with <resource address>" for destroys of objects
		// that are no longer in configuration, so include identifying context here.
		name, _ := d.GetOk("name")
		if nameStr, ok := name.(string); ok && nameStr != "" {
			return diag.FromErr(fmt.Errorf("%w alkira_credential_ssh_key_pair (name=%q id=%s)", err, nameStr, d.Id()))
		}
		return diag.FromErr(fmt.Errorf("%w alkira_credential_ssh_key_pair (id=%s)", err, d.Id()))
	}

	d.SetId("")
	return nil
}

// getSshKeyPairCredentialValue gets a value from config or environment variable.
// For WriteOnly fields, reads from raw config since values are not stored in state.
func getSshKeyPairCredentialValue(d *schema.ResourceData, field string, envVar string) (string, error) {
	// First try raw config (for WriteOnly fields)
	attrPath := cty.Path{cty.GetAttrStep{Name: field}}
	val, diags := d.GetRawConfigAt(attrPath)

	if !diags.HasError() && !val.IsNull() && val.IsKnown() && val.Type() == cty.String {
		strVal := val.AsString()
		if strVal != "" {
			return strVal, nil
		}
	}

	// Fall back to environment variable
	envValue := os.Getenv(envVar)
	if envValue != "" {
		return envValue, nil
	}

	// Return empty string if not set (field is Optional)
	return "", nil
}
