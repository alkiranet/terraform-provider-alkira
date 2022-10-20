package alkira

import (
	"github.com/alkiranet/alkira-client-go/alkira"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceAlkiraConnectorGcpVpc() *schema.Resource {
	return &schema.Resource{
		Description: "Use this data source to get an existing Gcp Vpc connector.",

		Read: dataSourceAlkiraConnectorGcpVpcRead,

		Schema: map[string]*schema.Schema{
			"name": {
				Description: "The name of the Gcp Vpc connector.",
				Type:        schema.TypeString,
				Required:    true,
			},
		},
	}
}

func dataSourceAlkiraConnectorGcpVpcRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*alkira.AlkiraClient)

	group, err := client.GetConnectorGcpVpcByName(d.Get("name").(string))

	if err != nil {
		return err
	}

	d.SetId(string(group.Id))
	return nil
}