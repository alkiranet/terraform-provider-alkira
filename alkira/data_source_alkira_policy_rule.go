package alkira

import (
	"github.com/alkiranet/alkira-client-go/alkira"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceAlkiraPolicyRule() *schema.Resource {
	return &schema.Resource{
		Description: "Use this data source to get an existing policy rule.",

		Read: dataSourceAlkiraPolicyRuleRead,

		Schema: map[string]*schema.Schema{
			"name": {
				Description: "The name of the policy rule.",
				Type:        schema.TypeString,
				Required:    true,
			},
		},
	}
}

func dataSourceAlkiraPolicyRuleRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*alkira.AlkiraClient)

	rule, err := client.GetPolicyRuleByName(d.Get("name").(string))

	if err != nil {
		return err
	}

	d.SetId(string(rule.Id))
	return nil
}
