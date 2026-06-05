package alkira

import (
	"github.com/alkiranet/alkira-client-go/alkira"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceAlkiraPeeringGatewayAzureVnetThirdPartyConnectorAttachment() *schema.Resource {
	return &schema.Resource{
		Description: "This data source allows to retrieve an existing Azure VNET Third Party Connector Attachment by its name.",
		Read:        dataSourceAlkiraPeeringGatewayAzureVnetThirdPartyConnectorAttachmentRead,

		Schema: map[string]*schema.Schema{
			"name": {
				Description: "The name of the attachment.",
				Type:        schema.TypeString,
				Required:    true,
			},
		},
	}
}

func dataSourceAlkiraPeeringGatewayAzureVnetThirdPartyConnectorAttachmentRead(d *schema.ResourceData, m interface{}) error {
	api := alkira.NewAzureVnetThirdPartyConnectorAttachment(m.(*alkira.AlkiraClient))

	resource, _, err := api.GetByName(d.Get("name").(string))
	if err != nil {
		return err
	}

	d.SetId(string(resource.Id))
	return nil
}
