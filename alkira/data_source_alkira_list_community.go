package alkira

import (
	"github.com/alkiranet/alkira-client-go/alkira"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceAlkiraListCommunity() *schema.Resource {
	return &schema.Resource{
		Description: "Use this data source to get an existing community list.",

		Read: dataSourceAlkiraListCommunityRead,

		Schema: map[string]*schema.Schema{
			"name": {
				Description: "The name of the policy prefix list.",
				Type:        schema.TypeString,
				Required:    true,
			},
			"values": {
				Description: "The value of the list.",
				Type:        schema.TypeString,
				Computed:    true,
				Optional:    true,
			},
		},
	}
}

func dataSourceAlkiraListCommunityRead(d *schema.ResourceData, m interface{}) error {
	api := alkira.NewListCommunity(m.(*alkira.AlkiraClient))

	list, err := api.GetByName(d.Get("name").(string))

	if err != nil {
		return err
	}

	d.SetId(string(list.Id))
	d.Set("values", list.Values)
	return nil
}
