package alkira

import (
	"context"
	"fmt"

	"github.com/alkiranet/alkira-client-go/alkira"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceAlkiraByoipPrefix() *schema.Resource {
	return &schema.Resource{
		Description:   "Manage BYOIP Prefix.",
		Create:        resourceByoipPrefix,
		Read:          resourceByoipPrefixRead,
		Update:        resourceByoipPrefixUpdate,
		Delete:        resourceByoipPrefixDelete,
		CustomizeDiff: resourceByoipPrefixCustomizeDiff,

		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
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

func resourceByoipPrefix(d *schema.ResourceData, m interface{}) error {

	// Init
	client := m.(*alkira.AlkiraClient)
	api := alkira.NewByoip(client)

	// Construct request
	request, err := generateByoipRequest(d, m)

	if err != nil {
		return err
	}

	// Send create request
	resource, provisionState, err := api.Create(request)

	if err != nil {
		return err
	}

	d.SetId(string(resource.Id))

	if client.Provision == true {
		d.Set("provision_state", provisionState)
	}

	return resourceByoipPrefixRead(d, m)
}

func resourceByoipPrefixRead(d *schema.ResourceData, m interface{}) error {

	// Init
	client := m.(*alkira.AlkiraClient)
	api := alkira.NewByoip(client)

	// Get the resource
	byoip, err := api.GetById(d.Id())

	if err != nil {
		return err
	}

	d.Set("prefix", byoip.Prefix)
	d.Set("cxp", byoip.Cxp)
	d.Set("description", byoip.Description)
	d.Set("message", byoip.ExtraAttributes.Message)
	d.Set("signature", byoip.ExtraAttributes.Signature)
	d.Set("public_key", byoip.ExtraAttributes.PublicKey)
	d.Set("do_not_advertise", byoip.DoNotAdvertise)

	// Set provision state
	_, provisionState, err := api.GetByName(d.Get("name").(string))

	if client.Provision == true && provisionState != "" {
		d.Set("provision_state", provisionState)
	}

	return nil
}

func resourceByoipPrefixUpdate(d *schema.ResourceData, m interface{}) error {
	return fmt.Errorf("`alkira_byoip_prefix` doesn't support upgrade. Please delete and create new one.")
}

func resourceByoipPrefixDelete(d *schema.ResourceData, m interface{}) error {

	// Init
	client := m.(*alkira.AlkiraClient)
	api := alkira.NewByoip(client)

	// Delete resource
	provisionState, err := api.Delete(d.Id())

	if err != nil {
		return err
	}

	// Check provision state
	if client.Provision == true && provisionState != "SUCCESS" {
		return fmt.Errorf("failed to delete byoip_prefix %s, provision failed", d.Id())
	}

	d.SetId("")
	return err
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
