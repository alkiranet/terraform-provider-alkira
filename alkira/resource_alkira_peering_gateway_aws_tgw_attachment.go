package alkira

import (
	"context"
	"fmt"
	"time"

	"github.com/alkiranet/alkira-client-go/alkira"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceAlkiraPeeringGatewayAwsTgwAttachment() *schema.Resource {
	return &schema.Resource{
		Description:   "Manage Peering Gateway AWS TGW Attachment.",
		CreateContext: resourcePeeringGatewayAwsTgwAttachmentCreate,
		ReadContext:   resourcePeeringGatewayAwsTgwAttachmentRead,
		UpdateContext: resourcePeeringGatewayAwsTgwAttachmentUpdate,
		DeleteContext: resourcePeeringGatewayAwsTgwAttachmentDelete,
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
			"requestor": {
				Description: "Initiator of transit gateway attachment.",
				Type:        schema.TypeString,
				Required:    true,
			},
			"peer_aws_region": {
				Description: "The AWS region of the peer TGW.",
				Type:        schema.TypeString,
				Required:    true,
			},
			"peer_aws_tgw_id": {
				Description: "The ID of AWS TGW.",
				Type:        schema.TypeString,
				Required:    true,
			},
			"peer_aws_account_id": {
				Description: "The AWS account ID of TGW.",
				Type:        schema.TypeString,
				Required:    true,
			},
			"peering_gateway_aws_tgw_id": {
				Description: "The ID of Peering Gateway AWS-TGW.",
				Type:        schema.TypeInt,
				Required:    true,
			},
			"state": {
				Description: "The state of the resource.",
				Type:        schema.TypeString,
				Computed:    true,
			},
		},
	}
}

func resourcePeeringGatewayAwsTgwAttachmentCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {

	// INIT
	api := alkira.NewPeeringGatewayAwsTgwAttachment(m.(*alkira.AlkiraClient))

	request, err := generatePeeringGatewayAwsTgwAttachmentRequest(d, m)

	if err != nil {
		return diag.FromErr(err)
	}

	// CREATE
	response, _, err, _, _ := api.Create(request)

	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(string(response.Id))

	// WAIT FOR STATE
	state := response.State

	for state != "PENDING_ACCEPTANCE" {
		resource, _, err := api.GetById(d.Id())

		if err != nil {
			return diag.Diagnostics{{
				Severity: diag.Warning,
				Summary:  "FAILED TO GET RESOURCE",
				Detail:   fmt.Sprintf("%s", err),
			}}
		}

		state = resource.State
		time.Sleep(5 * time.Second)
	}

	return resourcePeeringGatewayAwsTgwAttachmentRead(ctx, d, m)
}

func resourcePeeringGatewayAwsTgwAttachmentRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {

	// INIT
	api := alkira.NewPeeringGatewayAwsTgwAttachment(m.(*alkira.AlkiraClient))

	attachment, _, err := api.GetById(d.Id())

	if err != nil {
		return diag.Diagnostics{{
			Severity: diag.Warning,
			Summary:  "FAILED TO GET RESOURCE",
			Detail:   fmt.Sprintf("%s", err),
		}}
	}

	d.Set("name", attachment.Name)
	d.Set("description", attachment.Description)
	d.Set("requestor", attachment.Requestor)
	d.Set("peer_aws_region", attachment.PeerAwsRegion)
	d.Set("peer_aws_tgw_id", attachment.PeerAwsTgwId)
	d.Set("peer_aws_account_id", attachment.PeerAwsAccountId)
	d.Set("peering_gateway_aws_tgw_id", attachment.AwsTgwId)
	d.Set("state", attachment.State)

	return nil
}

func resourcePeeringGatewayAwsTgwAttachmentUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {

	// INIT
	api := alkira.NewPeeringGatewayAwsTgwAttachment(m.(*alkira.AlkiraClient))

	request, err := generatePeeringGatewayAwsTgwAttachmentRequest(d, m)

	if err != nil {
		return diag.FromErr(err)
	}

	// UPDATE
	_, err, _, _ = api.Update(d.Id(), request)

	return nil
}

func resourcePeeringGatewayAwsTgwAttachmentDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {

	// INIT
	api := alkira.NewPeeringGatewayAwsTgwAttachment(m.(*alkira.AlkiraClient))

	// DELETE
	_, err, _, _ := api.Delete(d.Id())

	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId("")

	return nil
}

// generatePeeringGatewayAwsTgwAttachmentRequest generate request
func generatePeeringGatewayAwsTgwAttachmentRequest(d *schema.ResourceData, m interface{}) (*alkira.PeeringGatewayAwsTgwAttachment, error) {

	request := &alkira.PeeringGatewayAwsTgwAttachment{
		Name:             d.Get("name").(string),
		Description:      d.Get("description").(string),
		Requestor:        d.Get("requestor").(string),
		PeerAwsRegion:    d.Get("peer_aws_region").(string),
		PeerAwsTgwId:     d.Get("peer_aws_tgw_id").(string),
		PeerAwsAccountId: d.Get("peer_aws_account_id").(string),
		AwsTgwId:         d.Get("peering_gateway_aws_tgw_id").(int),
	}

	return request, nil
}
