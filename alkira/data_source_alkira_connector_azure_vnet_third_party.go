package alkira

import (
	"github.com/alkiranet/alkira-client-go/alkira"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceAlkiraConnectorAzureVnetThirdParty() *schema.Resource {
	return &schema.Resource{
		Description: "This data source allows to retrieve an existing Azure VNET Third Party Connector by its name.",
		Read:        dataSourceAlkiraConnectorAzureVnetThirdPartyRead,

		Schema: map[string]*schema.Schema{
			"name": {
				Description: "The name of the connector.",
				Type:        schema.TypeString,
				Required:    true,
			},
			"implicit_group_id": {
				Description: "The ID of implicit group automatically created with the connector.",
				Type:        schema.TypeInt,
				Computed:    true,
			},
		},
	}
}

func dataSourceAlkiraConnectorAzureVnetThirdPartyRead(d *schema.ResourceData, m interface{}) error {
	api := alkira.NewAzureVnetThirdPartyConnector(m.(*alkira.AlkiraClient))

	resource, _, err := api.GetByName(d.Get("name").(string))
	if err != nil {
		return err
	}

	d.SetId(string(resource.Id))
	d.Set("implicit_group_id", resource.ImplicitGroupId)

	return nil
}
