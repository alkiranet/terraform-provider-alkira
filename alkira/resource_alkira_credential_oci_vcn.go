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
			StateContext: importWithReadValidation(resourceCredentialOciVcnRead),
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
				Sensitive:   true,
				WriteOnly:   true,
			},
			"fingerprint": {
				Description: "Fingerprint of the API key of the user.",
				Type:        schema.TypeString,
				Required:    true,
				Sensitive:   true,
				WriteOnly:   true,
			},
			"key": {
				Description: "API key of the user.",
				Type:        schema.TypeString,
				Required:    true,
				Sensitive:   true,
				WriteOnly:   true,
			},
			"tenant_ocid": {
				Description: "OCID of the tenant.",
				Type:        schema.TypeString,
				Required:    true,
				Sensitive:   true,
				WriteOnly:   true,
			},
		},
	}
}

func resourceCredentialOciVcn(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*alkira.AlkiraClient)

	c, err := generateCredentialOciVcnRequest(d)

	if err != nil {
		return diag.FromErr(err)
	}

	credentialId, err := client.CreateCredential(d.Get("name").(string), alkira.CredentialTypeOciVcn, c, 0)

	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(credentialId)
	return resourceCredentialOciVcnRead(ctx, d, meta)
}

func resourceCredentialOciVcnRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
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

	// Sensitive fields (user_ocid, fingerprint, key, tenant_ocid)
	// are NOT returned by the API for security reasons and are
	// maintained in the user's HCL configuration via WriteOnly.

	return nil
}

func resourceCredentialOciVcnUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*alkira.AlkiraClient)

	c, err := generateCredentialOciVcnRequest(d)

	if err != nil {
		return diag.FromErr(err)
	}

	err = client.UpdateCredential(d.Id(), d.Get("name").(string), alkira.CredentialTypeOciVcn, c, 0)

	if err != nil {
		return diag.FromErr(err)
	}

	return resourceCredentialOciVcnRead(ctx, d, meta)
}

func resourceCredentialOciVcnDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*alkira.AlkiraClient)

	err := client.DeleteCredential(d.Id(), alkira.CredentialTypeOciVcn)

	if err != nil {
		// Terraform may not print "with <resource address>" for destroys of objects
		// that are no longer in configuration, so include identifying context here.
		name, _ := d.GetOk("name")
		if nameStr, ok := name.(string); ok && nameStr != "" {
			return diag.FromErr(fmt.Errorf("%w alkira_credential_oci_vcn (name=%q id=%s)", err, nameStr, d.Id()))
		}
		return diag.FromErr(fmt.Errorf("%w alkira_credential_oci_vcn (id=%s)", err, d.Id()))
	}

	d.SetId("")
	return nil
}

// getOciVcnCredentialValue gets a value from config or environment variable.
// For WriteOnly fields, reads from raw config since values are not stored in state.
func getOciVcnCredentialValue(d *schema.ResourceData, field string, envVar string, required bool) (string, error) {
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

	if required {
		return "", fmt.Errorf("required field '%s' is not set in configuration and environment variable '%s' is not set", field, envVar)
	}
	return "", nil
}

func generateCredentialOciVcnRequest(d *schema.ResourceData) (alkira.CredentialOciVcn, error) {
	userOcid, err := getOciVcnCredentialValue(d, "user_ocid", "AK_OCI_USER_OCID", true)
	if err != nil {
		return alkira.CredentialOciVcn{}, err
	}

	fingerprint, err := getOciVcnCredentialValue(d, "fingerprint", "AK_OCI_FINGERPRINT", true)
	if err != nil {
		return alkira.CredentialOciVcn{}, err
	}

	key, err := getOciVcnCredentialValue(d, "key", "AK_OCI_KEY", true)
	if err != nil {
		return alkira.CredentialOciVcn{}, err
	}

	tenantOcid, err := getOciVcnCredentialValue(d, "tenant_ocid", "AK_OCI_TENANT_OCID", true)
	if err != nil {
		return alkira.CredentialOciVcn{}, err
	}

	c := alkira.CredentialOciVcn{
		UserId:      userOcid,
		FingerPrint: fingerprint,
		Key:         key,
		TenantId:    tenantOcid,
	}

	return c, nil
}
