package alkira

import (
	"github.com/alkiranet/alkira-client-go/alkira"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceAlkiraConnectorAzureVnet() *schema.Resource {
	return &schema.Resource{
		Description: "Use this data source to get an existing Azure VNET " +
			"connector.",

		Read: dataSourceAlkiraConnectorAzureVnetRead,

		Schema: map[string]*schema.Schema{
			"name": {
				Description: "The name of the connector.",
				Type:        schema.TypeString,
				Required:    true,
			},
			"implicit_group_id": {
				Description: "The implicit group associated with the connector.",
				Type:        schema.TypeInt,
				Computed:    true,
			},
		},
	}
}

func dataSourceAlkiraConnectorAzureVnetRead(d *schema.ResourceData, m interface{}) error {
	api := alkira.NewConnectorAzureVnet(m.(*alkira.AlkiraClient))

	resource, _, err := api.GetByName(d.Get("name").(string))

	if err != nil {
		return err
	}

	d.SetId(string(resource.Id))
	d.Set("implicit_group_id", resource.ImplicitGroupId)

	return nil
}
