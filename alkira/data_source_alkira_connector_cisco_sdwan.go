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

func dataSourceAlkiraConnectorCiscoSdwanRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*alkira.AlkiraClient)

	group, err := client.GetConnectorCiscoSdwanByName(d.Get("name").(string))

	if err != nil {
		return err
	}

	d.SetId(string(group.Id))
	return nil
}
