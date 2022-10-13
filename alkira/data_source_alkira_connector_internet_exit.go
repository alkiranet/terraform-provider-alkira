package alkira

import (
	"github.com/alkiranet/alkira-client-go/alkira"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceAlkiraConnectorInternetExit() *schema.Resource {
	return &schema.Resource{
		Description: "Use this data source to get an existing Internet Exit connector.",

		Read: dataSourceAlkiraConnectorInternetExitRead,

		Schema: map[string]*schema.Schema{
			"name": {
				Description: "The name of the Internet Exit connector.",
				Type:        schema.TypeString,
				Required:    true,
			},
		},
	}
}

func dataSourceAlkiraConnectorInternetExitRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*alkira.AlkiraClient)

	group, err := client.GetConnectorInternetByName(d.Get("name").(string))

	if err != nil {
		return err
	}

	d.SetId(string(group.Id))
	return nil
}
