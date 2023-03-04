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
		CustomizeDiff: resourceByoipPrefixCustomizeDiff,

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
			"message": {
				Description: "Message from AWS BYOIP.",
				Type:        schema.TypeString,
				Required:    true,
			},
			"provision_state": {
				Description: "The provision state of the resource.",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"signature": {
				Description: "Signautre from AWS BYOIP.",
				Type:        schema.TypeString,
				Required:    true,
			},
			"public_key": {
				Description: "Public Key from AWS BYOIP.",
				Type:        schema.TypeString,
				Required:    true,
			},
		},
	}
}

func resourceByoipPrefix(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {

	// Init
	client := m.(*alkira.AlkiraClient)
	api := alkira.NewByoip(client)

	// Construct request
	request, err := generateByoipRequest(d, m)

	if err != nil {
		return diag.FromErr(err)
	}

	// Send create request
	resource, provState, err, provErr := api.Create(request)

	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(string(resource.Id))

	if client.Provision == true {
		d.Set("provision_state", provState)

		if provErr != nil {
			return diag.Diagnostics{{
				Severity: diag.Warning,
				Summary:  "PROVISION FAILED",
				Detail:   fmt.Sprintf("%s", provErr),
			}}
		}
	}

	return resourceByoipPrefixRead(ctx, d, m)
}

func resourceByoipPrefixRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {

	// Init
	client := m.(*alkira.AlkiraClient)
	api := alkira.NewByoip(client)

	// Get the resource
	byoip, provState, err := api.GetById(d.Id())

	if err != nil {
		return diag.FromErr(err)
	}

	d.Set("prefix", byoip.Prefix)
	d.Set("cxp", byoip.Cxp)
	d.Set("description", byoip.Description)
	d.Set("message", byoip.ExtraAttributes.Message)
	d.Set("signature", byoip.ExtraAttributes.Signature)
	d.Set("public_key", byoip.ExtraAttributes.PublicKey)
	d.Set("do_not_advertise", byoip.DoNotAdvertise)

	// Set provision state
	if client.Provision == true && provState != "" {
		d.Set("provision_state", provState)
	}

	return nil
}

func resourceByoipPrefixUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	return diag.FromErr(fmt.Errorf("`alkira_byoip_prefix` doesn't support upgrade. Please delete and create new one."))
}

func resourceByoipPrefixDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {

	// Init
	client := m.(*alkira.AlkiraClient)
	api := alkira.NewByoip(client)

	// Delete resource
	provState, err, provErr := api.Delete(d.Id())

	if err != nil {
		return diag.FromErr(err)
	}

	// Check provision state
	if client.Provision == true && provState != "SUCCESS" {
		return diag.FromErr(fmt.Errorf("failed to delete byoip_prefix %s, provision failed, %v", d.Id(), provErr))
	}

	d.SetId("")
	return diag.FromErr(err)
}

func resourceByoipPrefixCustomizeDiff(ctx context.Context, d *schema.ResourceDiff, m interface{}) error {

	client := m.(*alkira.AlkiraClient)

	// Handle provision_state
	old, _ := d.GetChange("provision_state")

	if client.Provision == true && old == "FAILED" {
		d.SetNew("provision_state", "SUCCESS")
	}

	return nil
}

func generateByoipRequest(d *schema.ResourceData, m interface{}) (*alkira.Byoip, error) {

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
		DoNotAdvertise:  d.Get("do_not_advertise").(bool),
	}

	return request, nil
}
