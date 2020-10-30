package alkira

import (
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceAlkiraConnectorAzureVnet() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAlkiraConnectorAzureVnetRead,

		Schema: map[string]*schema.Schema{
			"connector_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Connector ID",
			},
		},
	}
}

func dataSourceAlkiraConnectorAzureVnetRead(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[INFO] Get Datasource")

	return nil
}
