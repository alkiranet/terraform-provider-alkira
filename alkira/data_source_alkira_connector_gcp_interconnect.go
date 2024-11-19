package alkira

import (
	"github.com/alkiranet/alkira-client-go/alkira"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceAlkiraConnectorGcpInterconnect() *schema.Resource {
	return &schema.Resource{
		Description: "Use this data source to get an existing GCP Interconnect connector.",

		Read: dataSourceAlkiraConnectorGcpInterconnectRead,

		Schema: map[string]*schema.Schema{
			"name": {
				Description: "The name of the GCP Interconnect connector.",
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

func dataSourceAlkiraConnectorGcpInterconnectRead(d *schema.ResourceData, m interface{}) error {
	api := alkira.NewConnectorGcpInterconnect(m.(*alkira.AlkiraClient))

	connector, _, err := api.GetByName(d.Get("name").(string))

	if err != nil {
		return err
	}

	d.SetId(string(connector.Id))
	d.Set("implicit_group_id", connector.ImplicitGroupId)

	return nil
}
