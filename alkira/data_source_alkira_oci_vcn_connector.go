package alkira

import (
	"github.com/alkiranet/alkira-client-go/alkira"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceAlkiraOciVcnConnector() *schema.Resource {
	return &schema.Resource{
		Description: "Use this data source to get an existing connector group.",

		Read: dataSourceAlkiraOciVcnConnectorRead,

		Schema: map[string]*schema.Schema{
			"name": {
				Description: "The name of the OCI VCN connector.",
				Type:        schema.TypeString,
				Required:    true,
			},
		},
	}
}

func dataSourceAlkiraOciVcnConnectorRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*alkira.AlkiraClient)

	group, err := client.GetConnectorOciVcnByName(d.Get("name").(string))

	if err != nil {
		return err
	}

	d.SetId(string(group.Id))
	return nil
}
