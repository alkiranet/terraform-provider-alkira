package alkira

import (
	"context"

	"github.com/alkiranet/alkira-client-go/alkira"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceAlkiraConnectorMultiCloudGateway() *schema.Resource {
	return &schema.Resource{
		Description: "Use this data source to get an existing Multi-Cloud Gateway connector.",

		ReadContext: dataSourceAlkiraConnectorMultiCloudGatewayRead,

		Schema: map[string]*schema.Schema{
			"name": {
				Description: "The name of the Multi-Cloud Gateway connector.",
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

func dataSourceAlkiraConnectorMultiCloudGatewayRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	api := alkira.NewConnectorMultiCloudGateway(m.(*alkira.AlkiraClient))

	connector, _, err := api.GetByName(d.Get("name").(string))

	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(string(connector.Id))
	d.Set("implicit_group_id", connector.ImplicitGroupId)

	return nil
}
