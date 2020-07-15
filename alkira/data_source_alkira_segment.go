package alkira

import (
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func dataSourceAlkiraSegment() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAlkiraSegmentRead,

		Schema: map[string]*schema.Schema{
			"connector_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Connector ID",
			},
		},
	}
}

func dataSourceAlkiraSegmentRead(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[INFO] Get Datasource")

	return nil
}
