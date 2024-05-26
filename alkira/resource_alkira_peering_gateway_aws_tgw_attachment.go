package alkira

import (
	"context"
	"fmt"

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
		CustomizeDiff: func(ctx context.Context, d *schema.ResourceDiff, m interface{}) error {
			client := m.(*alkira.AlkiraClient)

			old, _ := d.GetChange("provision_state")

			if client.Provision == true && old == "FAILED" {
				d.SetNew("provision_state", "SUCCESS")
			}

			return nil
		},
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
			"alkira_aws_tgw_id": {
				Description: "The ID of Alkira TGW.",
				Type:        schema.TypeInt,
				Required:    true,
			},
			"provision_state": {
				Description: "The provision state of the attachment.",
				Type:        schema.TypeString,
				Computed:    true,
			},
		},
	}
}

func resourcePeeringGatewayAwsTgwAttachmentCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {

	// INIT
	client := m.(*alkira.AlkiraClient)
	api := alkira.NewPeeringGatewayAwsTgwAttachment(m.(*alkira.AlkiraClient))

	request, err := generatePeeringGatewayAwsTgwAttachmentRequest(d, m)

	if err != nil {
		return diag.FromErr(err)
	}

	// CREATE
	response, provState, err, provErr := api.Create(request)

	if err != nil {
		return diag.FromErr(err)
	}

	// Set the state
	d.SetId(string(response.Id))

	if client.Provision == true {
		d.Set("provision_state", provState)

		if provErr != nil {
			return diag.Diagnostics{{
				Severity: diag.Warning,
				Summary:  "PROVISION (CREATE) FAILED",
				Detail:   fmt.Sprintf("%s", provErr),
			}}
		}
	}

	return resourcePeeringGatewayAwsTgwAttachmentRead(ctx, d, m)
}

func resourcePeeringGatewayAwsTgwAttachmentRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {

	// INIT
	client := m.(*alkira.AlkiraClient)
	api := alkira.NewPeeringGatewayAwsTgwAttachment(m.(*alkira.AlkiraClient))

	attachment, provState, err := api.GetById(d.Id())

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
	d.Set("alkira_aws_tgw_id", attachment.AwsTgwId)

	// Set provision state
	if client.Provision == true && provState != "" {
		d.Set("provision_state", provState)
	}

	return nil
}

func resourcePeeringGatewayAwsTgwAttachmentUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {

	// INIT
	client := m.(*alkira.AlkiraClient)
	api := alkira.NewPeeringGatewayAwsTgwAttachment(m.(*alkira.AlkiraClient))

	request, err := generatePeeringGatewayAwsTgwAttachmentRequest(d, m)

	if err != nil {
		return diag.FromErr(err)
	}

	// UPDATE
	provState, err, provErr := api.Update(d.Id(), request)

	// Set provision state
	if client.Provision == true {
		d.Set("provision_state", provState)

		if provErr != nil {
			return diag.Diagnostics{{
				Severity: diag.Warning,
				Summary:  "PROVISION (UPDATE) FAILED",
				Detail:   fmt.Sprintf("%s", provErr),
			}}
		}
	}

	return nil
}

func resourcePeeringGatewayAwsTgwAttachmentDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {

	// INIT
	client := m.(*alkira.AlkiraClient)
	api := alkira.NewPeeringGatewayAwsTgwAttachment(m.(*alkira.AlkiraClient))

	// DELETE
	provState, err, provErr := api.Delete(d.Id())

	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId("")

	if client.Provision == true && provState != "SUCCESS" {
		return diag.Diagnostics{{
			Severity: diag.Warning,
			Summary:  "PROVISION (DELETE) FAILED",
			Detail:   fmt.Sprintf("%s", provErr),
		}}
	}

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
		AwsTgwId:         d.Get("alkira_aws_tgw_id").(int),
	}

	return request, nil
}
