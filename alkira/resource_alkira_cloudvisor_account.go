package alkira

import (
	"context"
	"fmt"

	"github.com/alkiranet/alkira-client-go/alkira"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func resourceAlkiraCloudVisorAccount() *schema.Resource {
	return &schema.Resource{
		Description:   "Manage CloudVisor Accounts",
		CreateContext: resourceCloudVisorAccount,
		ReadContext:   resourceCloudVisorAccountRead,
		UpdateContext: resourceCloudVisorAccountUpdate,
		DeleteContext: resourceCloudVisorAccountDelete,

		Schema: map[string]*schema.Schema{
			"name": {
				Description: "The name of the account.",
				Type:        schema.TypeString,
				Required:    true,
			},
			"credential_id": {
				Description: "Credential Id to be used for the account.",
				Type:        schema.TypeString,
				Required:    true,
			},
			"cloud_provider": {
				Description: "Cloud provider of the account, currently, " +
					"`AWS` and `AZURE` are supported.",
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringInSlice([]string{"AWS", "AZURE"}, false),
			},
			"auto_sync": {
				Description: "The interval at which the account should be auto " +
					"synced. The value could be `NONE`, `DAILY`, `WEEKLY` and `MONTHLY`.",
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringInSlice([]string{"NONE", "DAILY", "WEEKLY", "MONTHLY"}, false),
			},
			"native_id": {
				Description: "The native cloud provider account Id.",
				Type:        schema.TypeString,
				Optional:    true,
			},
		},
	}
}

func resourceCloudVisorAccount(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {

	api := alkira.NewCloudProviderAccounts(m.(*alkira.AlkiraClient))

	// Construct request
	request := generateCloudVisorAccountRequest(d)

	// Send create request
	resource, _, err, valErr, _ := api.Create(request)

	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(resource.Id)

	// Handle validation error
	client := m.(*alkira.AlkiraClient)
	if client.Validate && valErr != nil {
		var diags diag.Diagnostics
		readDiags := resourceCloudVisorAccountRead(ctx, d, m)
		if readDiags.HasError() {
			diags = append(diags, readDiags...)
		}

		// Add the validation error
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "VALIDATION (CREATE) FAILED",
			Detail:   fmt.Sprintf("%s", valErr),
		})

		return diags
	}

	return resourceCloudVisorAccountRead(ctx, d, m)
}

func resourceCloudVisorAccountRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {

	api := alkira.NewCloudProviderAccounts(m.(*alkira.AlkiraClient))

	// Get resource
	account, _, err := api.GetById(d.Id())

	if err != nil {
		return diag.Diagnostics{{
			Severity: diag.Warning,
			Summary:  "FAILED TO GET RESOURCE",
			Detail:   fmt.Sprintf("%s", err),
		}}
	}

	d.Set("name", account.Name)
	d.Set("credential_id", account.CredentialId)
	d.Set("cloud_provider", account.CloudProvider)
	d.Set("auto_sync", account.AutoSync)
	d.Set("native_id", account.NativeId)

	return nil
}

func resourceCloudVisorAccountUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {

	api := alkira.NewCloudProviderAccounts(m.(*alkira.AlkiraClient))

	// Construct request
	request := generateCloudVisorAccountRequest(d)

	// Send update request
	_, err, valErr, _ := api.Update(d.Id(), request)

	if err != nil {
		return diag.FromErr(err)
	}

	// Handle validation error
	client := m.(*alkira.AlkiraClient)
	if client.Validate && valErr != nil {
		var diags diag.Diagnostics
		readDiags := resourceCloudVisorAccountRead(ctx, d, m)
		if readDiags.HasError() {
			diags = append(diags, readDiags...)
		}

		// Add the validation error
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "VALIDATION (UPDATE) FAILED",
			Detail:   fmt.Sprintf("%s", valErr),
		})

		return diags
	}

	return resourceCloudVisorAccountRead(ctx, d, m)
}

func resourceCloudVisorAccountDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {

	api := alkira.NewCloudProviderAccounts(m.(*alkira.AlkiraClient))

	_, err, valErr, _ := api.Delete(d.Id())

	if err != nil {
		// Terraform may not print "with <resource address>" for destroys of objects
		// that are no longer in configuration, so include identifying context here.
		name, _ := d.GetOk("name")
		if nameStr, ok := name.(string); ok && nameStr != "" {
			return diag.FromErr(fmt.Errorf("%w alkira_cloudvisor_account (name=%q id=%s)", err, nameStr, d.Id()))
		}
		return diag.FromErr(fmt.Errorf("%w alkira_cloudvisor_account (id=%s)", err, d.Id()))
	}

	d.SetId("")

	// Handle validation error
	client := m.(*alkira.AlkiraClient)
	if client.Validate && valErr != nil {
		return diag.Diagnostics{{
			Severity: diag.Error,
			Summary:  "VALIDATION (DELETE) FAILED",
			Detail:   fmt.Sprintf("%s", valErr),
		}}
	}

	return nil
}

func generateCloudVisorAccountRequest(d *schema.ResourceData) *alkira.CloudProviderAccount {
	return &alkira.CloudProviderAccount{
		Name:          d.Get("name").(string),
		CredentialId:  d.Get("credential_id").(string),
		CloudProvider: d.Get("cloud_provider").(string),
		AutoSync:      d.Get("auto_sync").(string),
		NativeId:      d.Get("native_id").(string),
	}
}
