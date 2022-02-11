package alkira

import (
	"log"

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

func resourceBillingTag(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*alkira.AlkiraClient)

	log.Printf("[INFO] Billing Tag Creating")
	id, err := client.CreateBillingTag(d.Get("name").(string), d.Get("description").(string))
	log.Printf("[INFO] Billing Tag ID: %s", id)

	if err != nil {
		return err
	}

	d.SetId(id)

	return resourceBillingTagRead(d, meta)
}

func resourceBillingTagRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*alkira.AlkiraClient)

	tag, err := client.GetBillingTagById(d.Id())

	if err != nil {
		return err
	}

	d.Set("name", tag.Name)
	d.Set("description", tag.Description)

	return nil
}

func resourceBillingTagUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*alkira.AlkiraClient)

	log.Printf("[INFO] Billing Tag Updating")
	err := client.UpdateBillingTag(d.Id(), d.Get("name").(string), d.Get("description").(string))

	if err != nil {
		return err
	}

	return resourceBillingTagRead(d, meta)
}

func resourceBillingTagDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*alkira.AlkiraClient)

	log.Printf("[INFO] Deleting Billing Tag %s", d.Id())
	err := client.DeleteBillingTag(d.Id())

	return err
}
