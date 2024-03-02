package alkira

import (
	"github.com/alkiranet/alkira-client-go/alkira"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceAlkiraInternetApplication() *schema.Resource {
	return &schema.Resource{
		Description: "This data source allows to retrieve an internet application.",
		Read:        dataSourceAlkiraInternetApplicationRead,
		Schema: map[string]*schema.Schema{
			"name": {
				Description: "The name of the Internet Application.",
				Type:        schema.TypeString,
				Required:    true,
			},
		},
	}
}

func dataSourceAlkiraInternetApplicationRead(d *schema.ResourceData, m interface{}) error {

	// INIT
	api := alkira.NewInternetApplication(m.(*alkira.AlkiraClient))
	app, _, err := api.GetByName(d.Get("name").(string))

	if err != nil {
		return err
	}

	d.SetId(string(app.Id))

	return nil
}
