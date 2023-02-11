package alkira

import (
	"github.com/alkiranet/alkira-client-go/alkira"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceAlkiraConnectorAzureVnet() *schema.Resource {
	return &schema.Resource{
		Description: "Use this data source to get an existing Azure Vnet connector.",

		Read: dataSourceAlkiraConnectorAzureVnetRead,

		Schema: map[string]*schema.Schema{
			"name": {
				Description: "The name of the Azure Vnet connector.",
				Type:        schema.TypeString,
				Required:    true,
			},
		},
	}
}

func dataSourceAlkiraConnectorAzureVnetRead(d *schema.ResourceData, meta interface{}) error {
	api := alkira.NewConnectorAzureVnet(m.(*alkira.AlkiraClient))

	resource, err := client.GetByName(d.Get("name").(string))

	if err != nil {
		return err
	}

	d.SetId(string(resource.Id))
	return nil
}
