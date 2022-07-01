package alkira

import (
	"github.com/alkiranet/alkira-client-go/alkira"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceAlkiraListGlobalCidr() *schema.Resource {
	return &schema.Resource{
		Description: "Use this data source to get an existing Global CIDR List.",

		Read: dataSourceAlkiraListGlobalCidrRead,

		Schema: map[string]*schema.Schema{
			"name": {
				Description: "The name of the Global CIDR List.",
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

func dataSourceAlkiraListGlobalCidrRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*alkira.AlkiraClient)

	list, err := client.GetGlobalCidrListByName(d.Get("name").(string))

	if err != nil {
		return err
	}

	d.SetId(string(list.Id))
	d.Set("values", list.Values)
	return nil
}
