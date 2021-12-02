package alkira

import (
	"github.com/alkiranet/alkira-client-go/alkira"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceAlkiraPolicyPrefixList() *schema.Resource {
	return &schema.Resource{
		Description: "Use this data source to get an existing policy prefix list.",

		Read: dataSourceAlkiraPolicyPrefixListRead,

		Schema: map[string]*schema.Schema{
			"name": {
				Description: "The name of the policy prefix list.",
				Type:        schema.TypeString,
				Required:    true,
			},
			"prefixes": {
				Description: "Prefixes in the prefix list.",
				Type:        schema.TypeString,
				Computed:    true,
				Optional:    true,
			},
		},
	}
}

func dataSourceAlkiraPolicyPrefixListRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*alkira.AlkiraClient)

	prefixList, err := client.GetPolicyPrefixListByName(d.Get("name").(string))

	if err != nil {
		return err
	}

	d.SetId(string(prefixList.Id))
	d.Set("prefixes", prefixList.Prefixes)
	return nil
}
