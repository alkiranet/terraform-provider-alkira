package alkira

import (
	"github.com/alkiranet/alkira-client-go/alkira"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceAlkiraPeeringGatewayAwsTgwAttachment() *schema.Resource {
	return &schema.Resource{
		Description: "This data source allows to retrieve an existing " +
			"Peering Gateway AWS TGW Attachment by its name.",

		Read: dataSourceAlkiraPeeringGatewayAwsTgwAttachmentRead,

		Schema: map[string]*schema.Schema{
			"name": {
				Description: "The name of the group.",
				Type:        schema.TypeString,
				Required:    true,
			},
		},
	}
}

func dataSourceAlkiraPeeringGatewayAwsTgwAttachmentRead(d *schema.ResourceData, m interface{}) error {
	api := alkira.NewPeeringGatewayAwsTgwAttachment(m.(*alkira.AlkiraClient))

	resource, _, err := api.GetByName(d.Get("name").(string))

	if err != nil {
		return err
	}

	d.SetId(string(resource.Id))
	return nil
}
