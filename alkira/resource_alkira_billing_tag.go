package alkira

import (
	"log"
	"strconv"

	"github.com/alkiranet/alkira-client-go/alkira"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceAlkiraBillingTag() *schema.Resource {
	return &schema.Resource{
		Create: resourceBillingTag,
		Read:   resourceBillingTagRead,
		Update: resourceBillingTagUpdate,
		Delete: resourceBillingTagDelete,

		Schema: map[string]*schema.Schema{
			"tag_id": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
		},
	}
}

func resourceBillingTag(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*alkira.AlkiraClient)
	name := d.Get("name").(string)

	log.Printf("[INFO] Billing Tag Creating")
	id, err := client.CreateBillingTag(name)
	log.Printf("[INFO] Billing Tag ID: %d", id)

	if err != nil {
		return err
	}

	d.SetId(strconv.Itoa(id))
	d.Set("tag_id", id)

	return resourceBillingTagRead(d, meta)
}

func resourceBillingTagRead(d *schema.ResourceData, meta interface{}) error {
	return nil
}

func resourceBillingTagUpdate(d *schema.ResourceData, meta interface{}) error {
	return resourceBillingTagRead(d, meta)
}

func resourceBillingTagDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*alkira.AlkiraClient)
	id := d.Get("tag_id").(int)

	log.Printf("[INFO] Deleting Billing Tag %d", id)
	err := client.DeleteBillingTag(id)

	if err != nil {
		return err
	}

	return nil
}