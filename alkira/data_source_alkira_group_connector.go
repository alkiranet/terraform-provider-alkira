package alkira

import (
	"strconv"

	"github.com/alkiranet/alkira-client-go/alkira"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceAlkiraGroupConnector() *schema.Resource {
	return &schema.Resource{
		Description: "Use this data source to get an existing connector group.",

		Read: dataSourceAlkiraGroupConnectorRead,

		Schema: map[string]*schema.Schema{
			"name": {
				Description: "The name of the group.",
				Type:        schema.TypeString,
				Required:    true,
			},
		},
	}
}

func dataSourceAlkiraGroupConnectorRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*alkira.AlkiraClient)

	group, err := client.GetConnectorGroupByName(d.Get("name").(string))

	if err != nil {
		return err
	}

	d.SetId(strconv.Itoa(group.Id))
	return nil
}
