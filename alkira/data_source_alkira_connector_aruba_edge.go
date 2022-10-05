package alkira

import (
	"github.com/alkiranet/alkira-client-go/alkira"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceAlkiraConnectorArubaEdge() *schema.Resource {
	return &schema.Resource{
		Description: "Use this data source to get an existing Aruba Edge connector.",

		Read: dataSourceAlkiraConnectorArubaEdgeRead,

		Schema: map[string]*schema.Schema{
			"name": {
				Description: "The name of the Aruba Edge connector.",
				Type:        schema.TypeString,
				Required:    true,
			},
		},
	}
}

func dataSourceAlkiraConnectorArubaEdgeRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*alkira.AlkiraClient)

	group, err := client.GetConnectorArubaEdgeByName(d.Get("name").(string))

	if err != nil {
		return err
	}

	d.SetId(string(group.Id))
	return nil
}
