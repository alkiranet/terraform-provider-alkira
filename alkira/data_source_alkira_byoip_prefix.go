package alkira

import (
	"github.com/alkiranet/alkira-client-go/alkira"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceAlkiraByoipPrefix() *schema.Resource {
	return &schema.Resource{
		Description: "Use this data source to get an existing BYOIP Prefix.",

		Read: dataSourceAlkiraByoipRead,

		Schema: map[string]*schema.Schema{
			"prefix": {
				Description: "Prefix for BYOIP.",
				Type:        schema.TypeString,
				Required:    true,
			},
		},
	}
}

func dataSourceAlkiraByoipPrefixRead(d *schema.ResourceData, m interface{}) error {
	api := alkira.NewByoip(m.(*alkira.AlkiraClient))

	resource, _, err := api.GetByName(d.Get("prefix").(string))

	if err != nil {
		return err
	}

	d.SetId(string(resource.Id))

	return nil
}
