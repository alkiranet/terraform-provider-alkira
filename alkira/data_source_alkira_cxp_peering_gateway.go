package alkira

import (
	"github.com/alkiranet/alkira-client-go/alkira"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceAlkiraCxpPeeringGateway() *schema.Resource {
	return &schema.Resource{
		Description: "This data source allows to retrieve an existing " +
			"Cxp Peering Gateway by its name.",

		Read: dataSourceAlkiraCxpPeeringGatewayRead,

		Schema: map[string]*schema.Schema{
			"name": {
				Description: "The name of the resource.",
				Type:        schema.TypeString,
				Required:    true,
			},
		},
	}
}

func dataSourceAlkiraCxpPeeringGatewayRead(d *schema.ResourceData, m interface{}) error {
	api := alkira.NewCxpPeeringGateway(m.(*alkira.AlkiraClient))

	resource, _, err := api.GetByName(d.Get("name").(string))
	if err != nil {
		return err
	}

	d.SetId(string(resource.Id))

	return nil
}
