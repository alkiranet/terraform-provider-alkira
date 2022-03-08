package alkira

import (
	"github.com/alkiranet/alkira-client-go/alkira"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceAlkiraGroupUser() *schema.Resource {
	return &schema.Resource{
		Description: "Use this data source to get an existing user group.",

		Read: dataSourceAlkiraGroupUserRead,

		Schema: map[string]*schema.Schema{
			"name": {
				Description: "The name of the group.",
				Type:        schema.TypeString,
				Required:    true,
			},
		},
	}
}

func dataSourceAlkiraGroupUserRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*alkira.AlkiraClient)

	group, err := client.GetUserGroupByName(d.Get("name").(string))

	if err != nil {
		return err
	}

	d.SetId(group.Id)
	return nil
}
