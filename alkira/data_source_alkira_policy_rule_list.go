package alkira

import (
	"github.com/alkiranet/alkira-client-go/alkira"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceAlkiraPolicyRuleList() *schema.Resource {
	return &schema.Resource{
		Description: "Use this data source to get an existing policy rule list.",

		Read: dataSourceAlkiraPolicyRuleListRead,

		Schema: map[string]*schema.Schema{
			"name": {
				Description: "The name of the policy rule list.",
				Type:        schema.TypeString,
				Required:    true,
			},
		},
	}
}

func dataSourceAlkiraPolicyRuleListRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*alkira.AlkiraClient)

	list, err := client.GetPolicyRuleListByName(d.Get("name").(string))

	if err != nil {
		return err
	}

	d.SetId(string(list.Id))
	return nil
}
