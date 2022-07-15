package alkira

import (
	"fmt"
	"log"

	"github.com/alkiranet/alkira-client-go/alkira"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceAlkiraByoipPrefix() *schema.Resource {
	return &schema.Resource{
		Description: "Manage BYOIP Prefix.",
		Create:      resourceByoipPrefix,
		Read:        resourceByoipPrefixRead,
		Update:      resourceByoipPrefixUpdate,
		Delete:      resourceByoipPrefixDelete,
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
	client := m.(*alkira.AlkiraClient)

	request, err := generateByoipRequest(d, m)

	if err != nil {
		log.Printf("[ERROR] failed to generate BYOIP request")
		return err
	}

	id, err := client.CreateByoip(request)

	if err != nil {
		return err
	}

	d.SetId(id)
	return resourceByoipPrefixRead(d, m)
}

func resourceByoipPrefixRead(d *schema.ResourceData, m interface{}) error {
	client := m.(*alkira.AlkiraClient)

	byoip, err := client.GetByoipById(d.Id())

	if err != nil {
		log.Printf("[ERROR] failed to get BYOIP %s", d.Id())
		return err
	}

	d.Set("prefix", byoip.Prefix)
	d.Set("cxp", byoip.Cxp)
	d.Set("description", byoip.Description)
	d.Set("message", byoip.ExtraAttributes.Message)
	d.Set("signature", byoip.ExtraAttributes.Signature)
	d.Set("public_key", byoip.ExtraAttributes.PublicKey)
	d.Set("do_not_advertise", byoip.DoNotAdvertise)

	return nil
}

func resourceByoipPrefixUpdate(d *schema.ResourceData, m interface{}) error {
	return fmt.Errorf("alkira_byoip doesn't support upgrade. Please delete and create new one.")
}

func resourceByoipPrefixDelete(d *schema.ResourceData, m interface{}) error {
	client := m.(*alkira.AlkiraClient)

	log.Printf("[INFO] Deleting BYOIP %s", d.Id())
	return client.DeleteByoip(d.Id())
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
