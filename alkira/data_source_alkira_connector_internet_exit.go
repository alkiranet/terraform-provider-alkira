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

func dataSourceAlkiraConnectorInternetExitRead(d *schema.ResourceData, m interface{}) error {
	api := alkira.NewConnectorInternet(m.(*alkira.AlkiraClient))

	resource, err := client.GetByName(d.Get("name").(string))

	if err != nil {
		return err
	}

	d.SetId(string(resource.Id))
	return nil
}
