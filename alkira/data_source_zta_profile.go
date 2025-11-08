package alkira

import (
	"github.com/alkiranet/alkira-client-go/alkira"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceZtaProfile() *schema.Resource {
	return &schema.Resource{
		Description: "The zta profile data source allows a zta profile to be retrieved by its name.",
		Read:        dataSourceZtaProfileRead,
		Schema: map[string]*schema.Schema{
			"name": {
				Description: "The name of the ZTA profile",
				Type:        schema.TypeString,
				Required:    true,
			},
		},
	}
}

func dataSourceZtaProfileRead(d *schema.ResourceData, m interface{}) error {
	api := alkira.NewZtaProfile(m.(*alkira.AlkiraClient))

	ztaProfile, _, err := api.GetByName(d.Get("name").(string))

	if err != nil {
		return err
	}

	d.SetId(ztaProfile.Id)
	return nil
}
