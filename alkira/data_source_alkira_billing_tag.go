package alkira

import (
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
		},
	}
}

func dataSourceAlkiraBillingTagRead(d *schema.ResourceData, m interface{}) error {
	api := alkira.NewBillingTag(m.(*alkira.AlkiraClient))

	billingTag, _, err := api.GetByName(d.Get("name").(string))

	if err != nil {
		return err
	}

	d.SetId(string(billingTag.Id))

	return nil
}
