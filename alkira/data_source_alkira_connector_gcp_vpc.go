package alkira

import (
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func dataSourceAlkiraConnectorGcpVpc() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAlkiraConnectorGcpVpcRead,

		Schema: map[string]*schema.Schema{
			"connector_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Connector ID",
			},
		},
	}
}

func dataSourceAlkiraConnectorGcpVpcRead(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[INFO] Get Datasource")

	return nil
}
