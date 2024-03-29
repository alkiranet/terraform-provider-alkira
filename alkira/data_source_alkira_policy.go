package alkira

import (
	"github.com/alkiranet/alkira-client-go/alkira"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceAlkiraPolicy() *schema.Resource {
	return &schema.Resource{
		Description: "Use this data source to get an existing policy.",

		Read: dataSourceAlkiraPolicyRuleRead,

		Schema: map[string]*schema.Schema{
			"name": {
				Description: "The name of the policy.",
				Type:        schema.TypeString,
				Required:    true,
			},
		},
	}
}

func dataSourceAlkiraPolicyRead(d *schema.ResourceData, m interface{}) error {
	api := alkira.NewTrafficPolicy(m.(*alkira.AlkiraClient))

	policy, _, err := api.GetByName(d.Get("name").(string))

	if err != nil {
		return err
	}

	d.SetId(string(policy.Id))
	return nil
}
