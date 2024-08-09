package alkira

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceAlkiraByoip() *schema.Resource {
	return &schema.Resource{
		Description: "Use this data source to get an existing BYOIP Prefix.",

		Read: dataSourceAlkiraByoipRead,

		Schema: map[string]*schema.Schema{
			"prefix": {
				Description: "Prefix for BYOIP.",
				Type:        schema.TypeString,
				Required:    true,
			},
		},
	}
}
