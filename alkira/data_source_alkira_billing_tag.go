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
		},
	}
}

func dataSourceAlkiraBillingTagRead(data *schema.ResourceData, meta interface{}) error {
	api := alkira.NewBillingTag(meta.(*alkira.AlkiraClient))

	billingTag, err := api.GetByName(d.Get("name").(string))

	if err != nil {
		return err
	}

	d.SetId(strconv.Itoa(billingTag.Id))

	return nil
}
