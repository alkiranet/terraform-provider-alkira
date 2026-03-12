package alkira

import (
	"context"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/alkiranet/alkira-client-go/alkira"

	"github.com/hashicorp/go-cty/cty"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceAlkiraCredentialAzureVnet() *schema.Resource {
	return &schema.Resource{
		Description: "Credential for interacting with Azure.\n\n" +
			"You could also provide your credentials via the following " +
			"environmental variables:\n\n * AK_AZURE_APPLICATION_ID\n " +
			"* AK_AZURE_SUBSCRIPTION_ID\n * AK_AZURE_SECRET_KEY\n " +
			"* AK_AZURE_TENANT_ID\n * AK_AZURE_ENVIRONMENT\n ",
		CreateContext: resourceCredentialAzureVnet,
		ReadContext:   resourceCredentialAzureVnetRead,
		UpdateContext: resourceCredentialAzureVnetUpdate,
		DeleteContext: resourceCredentialAzureVnetDelete,
		Importer: &schema.ResourceImporter{
			StateContext: importWithReadValidation(resourceCredentialAzureVnetRead),
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Description: "The name of the credential.",
				Type:        schema.TypeString,
				Required:    true,
			},
			"application_id": {
				Description: "Azure Application ID.",
				Type:        schema.TypeString,
				Required:    true,
				Sensitive:   true,
				WriteOnly:   true,
			},
			"subscription_id": {
				Description: "Azure subscription ID.",
				Type:        schema.TypeString,
				Optional:    true,
				Sensitive:   true,
				WriteOnly:   true,
			},
			"secret_key": {
				Description: "Azure Secret Key.",
				Type:        schema.TypeString,
				Required:    true,
				Sensitive:   true,
				WriteOnly:   true,
			},
			"tenant_id": {
				Description: "Azure Tenant ID.",
				Type:        schema.TypeString,
				Required:    true,
				Sensitive:   true,
				WriteOnly:   true,
			},
			"environment": {
				Description: "Azure environment can be `AZURE`, " +
					"`AZURE_CHINA` or `AZURE_US_GOVERNMENT`. The " +
					"default value is `AZURE`.",
				Type:     schema.TypeString,
				Optional: true,
				DefaultFunc: schema.EnvDefaultFunc(
					"AK_AZURE_ENVIRONMENT",
					nil),
			},
		},
	}
}

// getAzureVnetCredentialValue gets a value from config or environment variable
// For WriteOnly fields, reads from raw config since values are not stored in state
// Returns the value and an error if required but not found
func getAzureVnetCredentialValue(d *schema.ResourceData, field string, envVar string, required bool) (string, error) {
	// First try raw config (for WriteOnly fields and normal config values)
	attrPath := cty.Path{cty.GetAttrStep{Name: field}}
	val, diags := d.GetRawConfigAt(attrPath)

	if !diags.HasError() && !val.IsNull() && val.IsKnown() && val.Type() == cty.String {
		strVal := val.AsString()
		if strVal != "" {
			return strVal, nil
		}
	}

	// Check environment variable
	envValue := os.Getenv(envVar)
	if envValue != "" {
		return envValue, nil
	}

	if required {
		return "", fmt.Errorf("required field '%s' is not set in configuration and environment variable '%s' is not set", field, envVar)
	}
	return "", nil
}

// buildAzureVnetCredential builds the CredentialAzureVnet struct from ResourceData
func buildAzureVnetCredential(d *schema.ResourceData) (alkira.CredentialAzureVnet, error) {
	// Get values from config or environment variables
	applicationId, err := getAzureVnetCredentialValue(d, "application_id", "AK_AZURE_APPLICATION_ID", true)
	if err != nil {
		return alkira.CredentialAzureVnet{}, err
	}

	secretKey, err := getAzureVnetCredentialValue(d, "secret_key", "AK_AZURE_SECRET_KEY", true)
	if err != nil {
		return alkira.CredentialAzureVnet{}, err
	}

	tenantId, err := getAzureVnetCredentialValue(d, "tenant_id", "AK_AZURE_TENANT_ID", true)
	if err != nil {
		return alkira.CredentialAzureVnet{}, err
	}

	subscriptionId, _ := getAzureVnetCredentialValue(d, "subscription_id", "AK_AZURE_SUBSCRIPTION_ID", false)

	// Environment is not WriteOnly, use d.Get() directly
	environment := d.Get("environment").(string)

	return alkira.CredentialAzureVnet{
		ApplicationId:  applicationId,
		SecretKey:      secretKey,
		SubscriptionId: subscriptionId,
		TenantId:       tenantId,
		Environment:    environment,
	}, nil
}

func resourceCredentialAzureVnet(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*alkira.AlkiraClient)

	c, err := buildAzureVnetCredential(d)
	if err != nil {
		return diag.FromErr(err)
	}

	log.Printf("[INFO] Creating Credential (AZURE-VNET)")
	id, err := client.CreateCredential(d.Get("name").(string), alkira.CredentialTypeAzureVnet, c, 0)

	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(id)
	return resourceCredentialAzureVnetRead(ctx, d, meta)
}

func resourceCredentialAzureVnetRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*alkira.AlkiraClient)

	// Use GetCredentialById which now properly filters from GetCredentials
	credential, err := client.GetCredentialById(d.Id())
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			d.SetId("")
			return nil
		}
		return diag.FromErr(err)
	}

	// Set fields returned by API
	d.Set("name", credential.Name)

	// Set environment from SubType if available
	if credential.SubType != "" {
		d.Set("environment", credential.SubType)
	}

	// Note: Sensitive fields (secret_key, application_id, tenant_id, subscription_id)
	// are NOT returned by the API for security reasons and must be maintained
	// in the user's HCL configuration.

	return nil
}

func resourceCredentialAzureVnetUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*alkira.AlkiraClient)

	c, err := buildAzureVnetCredential(d)
	if err != nil {
		return diag.FromErr(err)
	}

	log.Printf("[INFO] Updating Credential (AZURE-VNET)")
	err = client.UpdateCredential(d.Id(), d.Get("name").(string), alkira.CredentialTypeAzureVnet, c, 0)

	if err != nil {
		return diag.FromErr(err)
	}

	return resourceCredentialAzureVnetRead(ctx, d, meta)
}

func resourceCredentialAzureVnetDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*alkira.AlkiraClient)

	err := client.DeleteCredential(d.Id(), alkira.CredentialTypeAzureVnet)

	if err != nil {
		// Terraform may not print "with <resource address>" for destroys of objects
		// that are no longer in configuration, so include identifying context here.
		name, _ := d.GetOk("name")
		if nameStr, ok := name.(string); ok && nameStr != "" {
			return diag.FromErr(fmt.Errorf("%w alkira_credential_azure_vnet (name=%q id=%s)", err, nameStr, d.Id()))
		}
		return diag.FromErr(fmt.Errorf("%w alkira_credential_azure_vnet (id=%s)", err, d.Id()))
	}

	d.SetId("")
	return nil
}
