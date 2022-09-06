package alkira

import (
	"github.com/alkiranet/alkira-client-go/alkira"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceAlkiraIpsecConnector() *schema.Resource {
	return &schema.Resource{
		Description: "Use this data source to get an existing ipsec connector.",

		Read: dataSourceAlkiraIpsecConnectorRead,

		Schema: map[string]*schema.Schema{
			"name": {
				Description: "The name of the ipsec connector.",
				Type:        schema.TypeString,
				Required:    true,
			},
		},
	}
}

func dataSourceAlkiraIpsecConnectorRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*alkira.AlkiraClient)

	group, err := client.GetConnectorIpsecByName(d.Get("name").(string))

	if err != nil {
		return err
	}

	d.SetId(string(group.Id))
	return nil
}
