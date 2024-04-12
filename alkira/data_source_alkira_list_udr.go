package alkira

import (
	"github.com/alkiranet/alkira-client-go/alkira"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceAlkiraListUdr() *schema.Resource {
	return &schema.Resource{
		Description: "Use this data source to get an existing UDR list.",

		Read: dataSourceAlkiraListUdrRead,

		Schema: map[string]*schema.Schema{
			"name": {
				Description: "The name of the list.",
				Type:        schema.TypeString,
				Required:    true,
			},
		},
	}
}

func dataSourceAlkiraListUdrRead(d *schema.ResourceData, m interface{}) error {
	api := alkira.NewUdrList(m.(*alkira.AlkiraClient))

	list, _, err := api.GetByName(d.Get("name").(string))

	if err != nil {
		return err
	}

	d.SetId(string(list.Id))
	return nil
}
