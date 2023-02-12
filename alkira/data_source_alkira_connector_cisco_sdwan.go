package alkira

import (
	"github.com/alkiranet/alkira-client-go/alkira"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceAlkiraConnectorCiscoSdwan() *schema.Resource {
	return &schema.Resource{
		Description: "Use this data source to get an existing Cisco Sdwan connector.",

		Read: dataSourceAlkiraConnectorCiscoSdwanRead,

		Schema: map[string]*schema.Schema{
			"name": {
				Description: "The name of the Cisco Sdwan connector.",
				Type:        schema.TypeString,
				Required:    true,
			},
		},
	}
}

func dataSourceAlkiraConnectorCiscoSdwanRead(d *schema.ResourceData, m interface{}) error {
	api := alkira.NewConnectorGcpVpc(m.(*alkira.AlkiraClient))

	resource, _, err := api.GetByName(d.Get("name").(string))

	if err != nil {
		return err
	}

	d.SetId(string(resource.Id))
	return nil
}
