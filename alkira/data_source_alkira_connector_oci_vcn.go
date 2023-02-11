package alkira

import (
	"github.com/alkiranet/alkira-client-go/alkira"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceAlkiraConnectorOciVcn() *schema.Resource {
	return &schema.Resource{
		Description: "Use this data source to get an existing OCI VCN connector.",

		Read: dataSourceAlkiraConnectorOciVcnRead,

		Schema: map[string]*schema.Schema{
			"name": {
				Description: "The name of the OCI VCN connector.",
				Type:        schema.TypeString,
				Required:    true,
			},
		},
	}
}

func dataSourceAlkiraConnectorOciVcnRead(d *schema.ResourceData, m interface{}) error {
	api := alkira.NewConnectorOciVcn(m.(*alkira.AlkiraClient))

	resource, err := client.GetByName(d.Get("name").(string))

	if err != nil {
		return err
	}

	d.SetId(string(resource.Id))
	return nil
}
