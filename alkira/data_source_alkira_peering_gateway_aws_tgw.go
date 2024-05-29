package alkira

import (
	"github.com/alkiranet/alkira-client-go/alkira"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceAlkiraPeeringGatewayAwsTgw() *schema.Resource {
	return &schema.Resource{
		Description: "This data source allows to retrieve an existing " +
			"Peering Gateway AWS TGW by its name.",

		Read: dataSourceAlkiraPeeringGatewayAwsTgwRead,

		Schema: map[string]*schema.Schema{
			"name": {
				Description: "The name of the resource.",
				Type:        schema.TypeString,
				Required:    true,
			},
			"aws_tgw_id": {
				Description: "The ID of the associated AWS TGW.",
				Type:        schema.TypeInt,
				Computed:    true,
				Optional:    true,
			},
		},
	}
}

func dataSourceAlkiraPeeringGatewayAwsTgwRead(d *schema.ResourceData, m interface{}) error {
	api := alkira.NewPeeringGatewayAwsTgw(m.(*alkira.AlkiraClient))

	resource, _, err := api.GetByName(d.Get("name").(string))

	if err != nil {
		return err
	}

	d.SetId(string(resource.Id))
	d.Set("aws_tgw_id", resource.AwsTgwId)

	return nil
}
