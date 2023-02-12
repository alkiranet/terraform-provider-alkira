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

func dataSourceAlkiraGroupUserRead(d *schema.ResourceData, m interface{}) error {
	api := alkira.NewUserGroup(m.(*alkira.AlkiraClient))

	group, _, err := api.GetByName(d.Get("name").(string))

	if err != nil {
		return err
	}

	d.SetId(string(group.Id))
	return nil
}
