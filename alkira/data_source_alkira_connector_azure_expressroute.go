package alkira

import (
	"github.com/alkiranet/alkira-client-go/alkira"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceAlkiraConnectorAzureExpressRoute() *schema.Resource {
	return &schema.Resource{
		Description: "Use this data source to get an existing Azure ExpressRoute connector.",

		Read: dataSourceAlkiraConnectorAzureExpressRouteRead,

		Schema: map[string]*schema.Schema{
			"name": {
				Description: "The name of the Azure ExpressRoute connector.",
				Type:        schema.TypeString,
				Required:    true,
			},
		},
	}
}

func dataSourceAlkiraConnectorAzureExpressRouteRead(d *schema.ResourceData, m interface{}) error {
	api := alkira.NewConnectorAzureExpressRoute(m.(*alkira.AlkiraClient))

	resource, _, err := api.GetByName(d.Get("name").(string))

	if err != nil {
		return err
	}

	d.SetId(string(resource.Id))
	return nil
}
