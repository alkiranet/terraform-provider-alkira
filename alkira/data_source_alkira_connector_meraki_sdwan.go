package alkira

import (
	"context"

	"github.com/alkiranet/alkira-client-go/alkira"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceAlkiraConnectorMerakiSDWAN() *schema.Resource {
	return &schema.Resource{
		Description: "Use this data source to get an existing Cisco Meraki SD-WAN connector.",

		ReadContext: dataSourceAlkiraConnectorMerakiSDWANRead,

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

func dataSourceAlkiraConnectorMerakiSDWANRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	api := alkira.NewConnectorMerakiSDWAN(m.(*alkira.AlkiraClient))

	connector, _, err := api.GetByName(d.Get("name").(string))

	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(string(connector.Id))
	d.Set("implicit_group_id", connector.ImplicitGroupId)

	return nil
}
