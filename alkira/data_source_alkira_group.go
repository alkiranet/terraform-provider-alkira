package alkira

import (
	"strconv"

	"github.com/alkiranet/alkira-client-go/alkira"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceAlkiraGroup() *schema.Resource {
	return &schema.Resource{
		Description: "This data source allows to retrieve an existing group by its name.",

		Read: dataSourceAlkiraGroupRead,

		Schema: map[string]*schema.Schema{
			"name": {
				Description: "The name of the group.",
				Type:        schema.TypeString,
				Required:    true,
			},
			"group_id": {
				Description: "The ID of the group.",
				Type:        schema.TypeInt,
				Computed:    true,
				Optional:    true,
			},
		},
	}
}

func dataSourceAlkiraGroupRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*alkira.AlkiraClient)

	group, err := client.GetConnectorGroupByName(d.Get("name").(string))

	if err != nil {
		return err
	}

	d.SetId(strconv.Itoa(group.Id))
	d.Set("group_id", group.Id)

	return nil
}
