package alkira

import (
	"strconv"

	"github.com/alkiranet/alkira-client-go/alkira"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceAlkiraBillingTag() *schema.Resource {
	return &schema.Resource{
		Description: "Use this data source to get information on an existing billing tag.",

		Read: dataSourceAlkiraBillingTagRead,

		Schema: map[string]*schema.Schema{
			"name": {
				Description: "The name of the billing tag.",
				Type:        schema.TypeString,
				Required:    true,
			},
			"billing_tag_id": {
				Description: "The ID of the billing tag.",
				Type:        schema.TypeInt,
				Computed:    true,
				Optional:    true,
			},
		},
	}
}

func dataSourceAlkiraBillingTagRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*alkira.AlkiraClient)

	billingTag, err := client.GetBillingTagByName(d.Get("name").(string))

	if err != nil {
		return err
	}

	d.SetId(strconv.Itoa(billingTag.Id))
	d.Set("billing_tag_id", billingTag.Id)

	return nil
}
