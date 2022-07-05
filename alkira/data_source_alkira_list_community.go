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

func dataSourceAlkiraListCommunityRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*alkira.AlkiraClient)

	list, err := client.GetListByName(d.Get("name").(string), alkira.ListTypeCommunity)

	if err != nil {
		return err
	}

	d.SetId(string(list.Id))
	d.Set("values", list.Values)
	return nil
}
