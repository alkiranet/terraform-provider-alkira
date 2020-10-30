package alkira

import (
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceAlkiraTenantNetwork() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAlkiraTenantNetworkRead,

		Schema: map[string]*schema.Schema{
			"connector_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Connector ID",
			},
		},
	}
}

func dataSourceAlkiraTenantNetworkRead(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[INFO] Get Datasource")

	return nil
}
