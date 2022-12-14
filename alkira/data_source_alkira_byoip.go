package alkira

import (
	"github.com/alkiranet/alkira-client-go/alkira"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceAlkiraByoip() *schema.Resource {
	return &schema.Resource{
		Description: "Use this data source to get an existing BYOIP Prefix.",

		Read: dataSourceAlkiraByoipRead,

		Schema: map[string]*schema.Schema{
			"name": {
				Description: "Prefix for BYOIP.",
				Type:        schema.TypeString,
				Required:    true,
			},
		},
	}
}

func dataSourceAlkiraByoipRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*alkira.AlkiraClient)

	prefix, err := client.GetByoipPrefixByName(d.Get("name").(string))

	if err != nil {
		return err
	}

	d.SetId(string(prefix.Id))

	return nil
}
