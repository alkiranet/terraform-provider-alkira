package alkira

import (
	"context"

	"github.com/alkiranet/alkira-client-go/alkira"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceAlkiraMultiCloudGateway() *schema.Resource {
	return &schema.Resource{
		Description: "Retrieve an existing Multi-Cloud Gateway by its exact name. The lookup searches all pages and fails with the match count if zero or multiple gateways have that name.",
		ReadContext: dataSourceAlkiraMultiCloudGatewayRead,
		Schema: map[string]*schema.Schema{
			"name": {
				Description: "The exact name of the gateway.",
				Type:        schema.TypeString,
				Required:    true,
			},
			"native_id": {
				Description: "Native gateway ID assigned by Lumen.",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"tier": {
				Description: "The gateway bandwidth tier reported by the API.",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"state": {
				Description: "Provisioning state of the gateway.",
				Type:        schema.TypeString,
				Computed:    true,
			},
		},
	}
}

func dataSourceAlkiraMultiCloudGatewayRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	gateway, _, err := alkira.GetMultiCloudGatewayByName(m.(*alkira.AlkiraClient), d.Get("name").(string))
	if err != nil {
		return diag.FromErr(err)
	}
	for key, value := range map[string]string{
		"native_id": gateway.NativeId,
		"tier":      gateway.Tier,
		"state":     gateway.State,
	} {
		if err := d.Set(key, value); err != nil {
			return diag.FromErr(err)
		}
	}
	d.SetId(gateway.Id.String())
	return nil
}
