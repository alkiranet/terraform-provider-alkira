package alkira

import (
	"github.com/alkiranet/alkira-client-go/alkira"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceAlkiraConnectorAwsVpc() *schema.Resource {
	return &schema.Resource{
		Description: "Use this data source to get an existing Aws Vpc connector.",

		Read: dataSourceAlkiraConnectorAwsVpcRead,

		Schema: map[string]*schema.Schema{
			"name": {
				Description: "The name of the Aws Vpc connector.",
				Type:        schema.TypeString,
				Required:    true,
			},
		},
	}
}

func dataSourceAlkiraConnectorAwsVpcRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*alkira.AlkiraClient)

	group, err := client.GetConnectorAwsVpcByName(d.Get("name").(string))

	if err != nil {
		return err
	}

	d.SetId(string(group.Id))
	return nil
}
