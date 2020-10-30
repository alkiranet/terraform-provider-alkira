package alkira

import (
	"strconv"

	"github.com/alkiranet/alkira-client-go/alkira"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceAlkiraGroup() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAlkiraGroupRead,

		Schema: map[string]*schema.Schema{
			"group_id": {
				Type:     schema.TypeInt,
				Required: true,
			},
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
		},
	}
}

func dataSourceAlkiraGroupRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*alkira.AlkiraClient)

	group, err := client.GetGroupByName(d.Get("name").(string))

	if err != nil {
		return err
	}

	d.SetId(strconv.Itoa(group.Id))
	d.Set("group_id", group.Id)

	return nil
}
