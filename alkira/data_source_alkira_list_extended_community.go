package alkira

import (
	"github.com/alkiranet/alkira-client-go/alkira"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceAlkiraListExtendedCommunity() *schema.Resource {
	return &schema.Resource{
		Description: "Use this data source to get an existing extended community list.",

		Read: dataSourceAlkiraListExtendedCommunityRead,

		Schema: map[string]*schema.Schema{
			"name": {
				Description: "The name of the extended community list.",
				Type:        schema.TypeString,
				Required:    true,
			},
			"values": {
				Description: "The values of the list.",
				Type:        schema.TypeString,
				Computed:    true,
				Optional:    true,
			},
		},
	}
}

func dataSourceAlkiraListExtendedCommunityRead(d *schema.ResourceData, m interface{}) error {
	api := alkira.NewListExtendedCommunity(m.(*alkira.AlkiraClient))

	list, err := api.GetByName(d.Get("name").(string))

	if err != nil {
		return err
	}

	d.SetId(string(list.Id))
	d.Set("values", list.Values)
	return nil
}
