package alkira

import (
	"github.com/alkiranet/alkira-client-go/alkira"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceAlkiraBillingTag() *schema.Resource {
	return &schema.Resource{
		Description: "Manage Billing Tag.",
		Create:      resourceBillingTag,
		Read:        resourceBillingTagRead,
		Update:      resourceBillingTagUpdate,
		Delete:      resourceBillingTagDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Description: "Billing Tag Name.",
				Type:        schema.TypeString,
				Required:    true,
			},
			"description": {
				Description: "Billing Tag Description.",
				Type:        schema.TypeString,
				Optional:    true,
			},
		},
	}
}

func resourceBillingTag(d *schema.ResourceData, m interface{}) error {

	api := alkira.NewBillingTag(m.(*alkira.AlkiraClient))

	// Construct request
	request := &alkira.BillingTag{
		Name:        d.Get("name").(string),
		Description: d.Get("description").(string),
	}

	// Send create request
	response, _, err, _ := api.Create(request)

	if err != nil {
		return err
	}

	d.SetId(string(response.Id))

	return resourceBillingTagRead(d, m)
}

func resourceBillingTagRead(d *schema.ResourceData, m interface{}) error {

	api := alkira.NewBillingTag(m.(*alkira.AlkiraClient))

	// Get resource
	tag, _, err := api.GetById(d.Id())

	if err != nil {
		return err
	}

	d.Set("name", tag.Name)
	d.Set("description", tag.Description)

	return nil
}

func resourceBillingTagUpdate(d *schema.ResourceData, m interface{}) error {

	api := alkira.NewBillingTag(m.(*alkira.AlkiraClient))

	// Construct request
	request := &alkira.BillingTag{
		Name:        d.Get("name").(string),
		Description: d.Get("description").(string),
	}

	// Send update request
	_, err, _ := api.Update(d.Id(), request)

	if err != nil {
		return err
	}

	return resourceBillingTagRead(d, m)
}

func resourceBillingTagDelete(d *schema.ResourceData, m interface{}) error {

	api := alkira.NewBillingTag(m.(*alkira.AlkiraClient))

	_, err, _ := api.Delete(d.Id())

	if err != nil {
		return err
	}

	d.SetId("")
	return err
}
