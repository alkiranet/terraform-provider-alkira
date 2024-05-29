package alkira

import (
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
		},
	}
}

func dataSourceAlkiraGroupRead(d *schema.ResourceData, m interface{}) error {
	api := alkira.NewGroup(m.(*alkira.AlkiraClient))

	group, _, err := api.GetByName(d.Get("name").(string))

	if err != nil {
		return err
	}

	d.SetId(string(group.Id))
	return nil
}
