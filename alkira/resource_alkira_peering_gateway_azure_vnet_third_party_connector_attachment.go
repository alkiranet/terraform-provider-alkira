package alkira

import (
	"context"
	"fmt"

	"github.com/alkiranet/alkira-client-go/alkira"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceAlkiraPeeringGatewayAzureVnetThirdPartyConnectorAttachment() *schema.Resource {
	return &schema.Resource{
		Description:   "Manage Azure VNET Third Party Connector Attachment.",
		CreateContext: resourcePeeringGatewayAzureVnetThirdPartyConnectorAttachmentCreate,
		ReadContext:   resourcePeeringGatewayAzureVnetThirdPartyConnectorAttachmentRead,
		UpdateContext: resourcePeeringGatewayAzureVnetThirdPartyConnectorAttachmentUpdate,
		DeleteContext: resourcePeeringGatewayAzureVnetThirdPartyConnectorAttachmentDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"name": {
				Description: "The name of the attachment.",
				Type:        schema.TypeString,
				Required:    true,
			},
			"description": {
				Description: "Description of the attachment.",
				Type:        schema.TypeString,
				Optional:    true,
			},
			"cxp_peering_gateway_id": {
				Description: "The ID of the CXP peering gateway.",
				Type:        schema.TypeInt,
				Required:    true,
			},
			"azure_vnet_id": {
				Description: "Azure Virtual Network Resource ID. Format: /subscriptions/{subscriptionId}/resourceGroups/{resourceGroup}/providers/Microsoft.Network/virtualNetworks/{vnetName}",
				Type:        schema.TypeString,
				Required:    true,
			},
			"state": {
				Description: "The state of the attachment.",
				Type:        schema.TypeString,
				Computed:    true,
			},
		},
	}
}

func resourcePeeringGatewayAzureVnetThirdPartyConnectorAttachmentCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	api := alkira.NewAzureVnetThirdPartyConnectorAttachment(m.(*alkira.AlkiraClient))

	request, err := generatePeeringGatewayAzureVnetThirdPartyConnectorAttachmentRequest(d)
	if err != nil {
		return diag.FromErr(err)
	}

	response, _, err, _, _ := api.Create(request)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(string(response.Id))
	return resourcePeeringGatewayAzureVnetThirdPartyConnectorAttachmentRead(ctx, d, m)
}

func resourcePeeringGatewayAzureVnetThirdPartyConnectorAttachmentRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	api := alkira.NewAzureVnetThirdPartyConnectorAttachment(m.(*alkira.AlkiraClient))

	attachment, _, err := api.GetById(d.Id())
	if err != nil {
		return diag.Diagnostics{{
			Severity: diag.Warning,
			Summary:  "FAILED TO GET RESOURCE",
			Detail:   fmt.Sprintf("Failed to get attachment: %s", err),
		}}
	}

	d.Set("name", attachment.Name)
	d.Set("description", attachment.Description)
	d.Set("cxp_peering_gateway_id", attachment.CxpPeeringGatewayId)
	d.Set("azure_vnet_id", attachment.VnetId)
	d.Set("state", attachment.State)

	return nil
}

func resourcePeeringGatewayAzureVnetThirdPartyConnectorAttachmentUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	api := alkira.NewAzureVnetThirdPartyConnectorAttachment(m.(*alkira.AlkiraClient))

	request, err := generatePeeringGatewayAzureVnetThirdPartyConnectorAttachmentRequest(d)
	if err != nil {
		return diag.FromErr(err)
	}

	_, err, _, _ = api.Update(d.Id(), request)
	if err != nil {
		return diag.FromErr(err)
	}

	return resourcePeeringGatewayAzureVnetThirdPartyConnectorAttachmentRead(ctx, d, m)
}

func resourcePeeringGatewayAzureVnetThirdPartyConnectorAttachmentDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	api := alkira.NewAzureVnetThirdPartyConnectorAttachment(m.(*alkira.AlkiraClient))

	_, err, _, _ := api.Delete(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId("")
	return nil
}

func generatePeeringGatewayAzureVnetThirdPartyConnectorAttachmentRequest(d *schema.ResourceData) (*alkira.AzureVnetThirdPartyConnectorAttachment, error) {
	request := &alkira.AzureVnetThirdPartyConnectorAttachment{
		Name:                d.Get("name").(string),
		Description:         d.Get("description").(string),
		CxpPeeringGatewayId: d.Get("cxp_peering_gateway_id").(int),
		VnetId:              d.Get("azure_vnet_id").(string),
	}

	return request, nil
}
