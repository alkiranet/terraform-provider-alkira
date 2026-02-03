package alkira

import (
	"context"
	"fmt"

	"github.com/alkiranet/alkira-client-go/alkira"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceAlkiraByoipPrefix() *schema.Resource {
	return &schema.Resource{
		Description:   "Manage BYOIP Prefix.",
		CreateContext: resourceByoipPrefix,
		ReadContext:   resourceByoipPrefixRead,
		UpdateContext: resourceByoipPrefixUpdate,
		DeleteContext: resourceByoipPrefixDelete,
		CustomizeDiff: func(ctx context.Context, d *schema.ResourceDiff, m interface{}) error {
			client := m.(*alkira.AlkiraClient)

			old, _ := d.GetChange("provision_state")

			if client.Provision && old == "FAILED" {
				d.SetNew("provision_state", "SUCCESS")
			}

			return nil
		},
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			"prefix": {
				Description: "Prefix for BYOIP.",
				Type:        schema.TypeString,
				Required:    true,
			},
			"cxp": {
				Description: "CXP region.",
				Type:        schema.TypeString,
				Required:    true,
			},
			"description": {
				Description: "Description for the list.",
				Type:        schema.TypeString,
				Optional:    true,
			},
			"do_not_advertise": {
				Description: "Do not advertise.",
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
			},
			"cloud_provider": {
				Description: "Cloud provider for the BYOIP." +
					"This must match CXP's provider.",
				Type:     schema.TypeString,
				Required: true},
			"message": {
				Description: "Message from BYOIP." +
					"For AWS, the format of the message is" +
					" `1|aws|account|cidr|YYYYMMDD|SHA256|RSAPSS`," +
					" where the date is the expiry date of the message." +
					"For AZURE, the format of the message is" +
					" `subscriptionId|cidr|YYYYMMDD`," +
					" where the date is the validity date on the ROA.",
				Type:     schema.TypeString,
				Required: true,
			},
			"provision_state": {
				Description: "The provision state of the resource.",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"signature": {
				Description: "Signature from the BYOIP." +
					"For AZURE, the signature scheme is `SHA256RSA`.",
				Type:     schema.TypeString,
				Required: true,
			},
			"public_key": {
				Description: "The RSA 2048-bit public key from the BYOIP.",
				Type:        schema.TypeString,
				Required:    true,
			},
		},
	}
}

func resourceByoipPrefix(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {

	// INIT
	client := m.(*alkira.AlkiraClient)
	api := alkira.NewByoip(client)

	// Construct request
	request, err := generateByoipRequest(d)

	if err != nil {
		return diag.FromErr(err)
	}

	// Send create request
	resource, provState, err, valErr, provErr := api.Create(request)

	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(string(resource.Id))

	// Handle validation error
	if client.Validate && valErr != nil {
		var diags diag.Diagnostics
		readDiags := resourceByoipPrefixRead(ctx, d, m)
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

	if client.Provision {
		d.Set("provision_state", provState)

		if provErr != nil {
			return diag.Diagnostics{{
				Severity: diag.Warning,
				Summary:  "PROVISION (CREATE) FAILED",
				Detail:   fmt.Sprintf("%s", provErr),
			}}
		}
	}

	return resourceByoipPrefixRead(ctx, d, m)
}

func resourceByoipPrefixRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {

	// INIT
	client := m.(*alkira.AlkiraClient)
	api := alkira.NewByoip(client)

	// Get the resource
	byoip, provState, err := api.GetById(d.Id())

	if err != nil {
		return diag.Diagnostics{{
			Severity: diag.Warning,
			Summary:  "FAILED TO GET RESOURCE",
			Detail:   fmt.Sprintf("%s", err),
		}}
	}

	d.Set("prefix", byoip.Prefix)
	d.Set("cxp", byoip.Cxp)
	d.Set("description", byoip.Description)
	d.Set("message", byoip.ExtraAttributes.Message)
	d.Set("signature", byoip.ExtraAttributes.Signature)
	d.Set("public_key", byoip.ExtraAttributes.PublicKey)
	d.Set("do_not_advertise", byoip.DoNotAdvertise)
	d.Set("cloud_provider", byoip.CloudProvider)

	// Set provision state
	if client.Provision && provState != "" {
		d.Set("provision_state", provState)
	}

	return nil
}

func resourceByoipPrefixUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	return diag.FromErr(fmt.Errorf("`alkira_byoip_prefix` doesn't support upgrade. Please delete and create new one"))
}

func resourceByoipPrefixDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {

	// INIT
	client := m.(*alkira.AlkiraClient)
	api := alkira.NewByoip(client)

	// Delete resource
	provState, err, valErr, provErr := api.Delete(d.Id())

	if err != nil {
		// Terraform may not print "with <resource address>" for destroys of objects
		// that are no longer in configuration, so include identifying context here.
		name, _ := d.GetOk("name")
		if nameStr, ok := name.(string); ok && nameStr != "" {
			return diag.FromErr(fmt.Errorf("%w alkira_byoip_prefix (name=%q id=%s)", err, nameStr, d.Id()))
		}
		return diag.FromErr(fmt.Errorf("%w alkira_byoip_prefix (id=%s)", err, d.Id()))
	}

	d.SetId("")

	// Handle validation error
	if client.Validate && valErr != nil {
		return diag.Diagnostics{{
			Severity: diag.Error,
			Summary:  "VALIDATION (DELETE) FAILED",
			Detail:   fmt.Sprintf("%s", valErr),
		}}
	}

	// Check provision state
	if client.Provision && provState != "SUCCESS" {
		return diag.Diagnostics{{
			Severity: diag.Warning,
			Summary:  "PROVISION (DELETE) FAILED",
			Detail:   fmt.Sprintf("%s", provErr),
		}}
	}

	return nil
}

func generateByoipRequest(d *schema.ResourceData) (*alkira.Byoip, error) {

	attributes := alkira.ByoipExtraAttributes{
		Message:   d.Get("message").(string),
		Signature: d.Get("signature").(string),
		PublicKey: d.Get("public_key").(string),
	}

	request := &alkira.Byoip{
		Prefix:          d.Get("prefix").(string),
		Cxp:             d.Get("cxp").(string),
		Description:     d.Get("description").(string),
		ExtraAttributes: attributes,
		CloudProvider:   d.Get("cloud_provider").(string),
		DoNotAdvertise:  d.Get("do_not_advertise").(bool),
	}

	return request, nil
}
