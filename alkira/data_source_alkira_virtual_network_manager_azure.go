package alkira

import (
	"github.com/alkiranet/alkira-client-go/alkira"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceAlkiraVirtualNetworkManagerAzure() *schema.Resource {
	return &schema.Resource{
		Description: "This data source allows to retrieve an existing" +
			"Alkira Azure Virtual Network Manager by its name.",

		Read: dataSourceVirtualNetworkManagerAzureRead,

		Schema: map[string]*schema.Schema{
			"name": {
				Description: "The name of the Azure Virtual Network Manager.",
				Type:        schema.TypeString,
				Required:    true,
			},
		},
	}
}

func dataSourceVirtualNetworkManagerAzureRead(d *schema.ResourceData, m interface{}) error {
	api := alkira.NewVirtualNetworkManagerAzure(m.(*alkira.AlkiraClient))

	resource, _, err := api.GetByName(d.Get("name").(string))

	if err != nil {
		return err
	}

	d.SetId(string(resource.Id))
	return nil
}
