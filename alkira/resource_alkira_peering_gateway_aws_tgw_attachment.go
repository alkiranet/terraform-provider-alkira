package alkira

import (
	"context"
	"fmt"
	"time"

	"github.com/alkiranet/alkira-client-go/alkira"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func resourceAlkiraPeeringGatewayAwsTgwAttachment() *schema.Resource {
	return &schema.Resource{
		Description:   "Manage Peering Gateway AWS TGW Attachment.",
		CreateContext: resourcePeeringGatewayAwsTgwAttachmentCreate,
		ReadContext:   resourcePeeringGatewayAwsTgwAttachmentRead,
		UpdateContext: resourcePeeringGatewayAwsTgwAttachmentUpdate,
		DeleteContext: resourcePeeringGatewayAwsTgwAttachmentDelete,
		Importer: &schema.ResourceImporter{
			StateContext: importWithReadValidation(resourcePeeringGatewayAwsTgwAttachmentRead),
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
				Optional:    true,
			},
			"peer_aws_tgw_id": {
				Description: "The ID of AWS TGW.",
				Type:        schema.TypeString,
				Optional:    true,
			},
			"peer_aws_account_id": {
				Description: "The AWS account ID of TGW.",
				Type:        schema.TypeString,
				Required:    true,
			},
			"peer_allowed_prefixes": {
				Description: "List of allowed CIDR prefixes for the peer.",
				Type:        schema.TypeList,
				Elem: &schema.Schema{
					Type:         schema.TypeString,
					ValidateFunc: validation.StringIsNotWhiteSpace,
				},
				Optional: true,
			},
			"peer_direct_connect_gateway_id": {
				Description: "The AWS Direct Connect Gateway ID.",
				Type:        schema.TypeString,
				Optional:    true,
			},
			"type": {
				Description: "The type of attachment. " +
					"Can be one of `AWS_TRANSIT_GATEWAY` and `AWS_DIRECT_CONNECT_GATEWAY`.",
				Type:         schema.TypeString,
				ValidateFunc: validation.StringInSlice([]string{"AWS_TRANSIT_GATEWAY", "AWS_DIRECT_CONNECT_GATEWAY"}, false),
				Optional:     true,
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
			"failure_reason": {
				Description: "Failure reason if there is any failure in creation/deletion",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"direct_connect_gateway_association_proposal_state": {
				Description: "State of latest direct connect gateway association proposal created by AWS_DIRECT_CONNECT_GATEWAY create/update",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"direct_connect_gateway_association_proposal_id": {
				Description: "id of latest direct connect gateway association proposal created by AWS_DIRECT_CONNECT_GATEWAY create/update",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"direct_connect_gateway_association_proposal_created_at": {
				Description: "Timestamp indicating when direct connect gateway association proposal was created",
				Type:        schema.TypeInt,
				Computed:    true,
			},
			"direct_connect_gateway_association_proposal_updated_at": {
				Description: "Timestamp indicating when direct connect gateway association proposal was last updated",
				Type:        schema.TypeInt,
				Computed:    true,
			},
		},
	}
}

func resourcePeeringGatewayAwsTgwAttachmentCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {

	// INIT
	api := alkira.NewPeeringGatewayAwsTgwAttachment(m.(*alkira.AlkiraClient))

	request, err := generateAwsTgwAttachmentRequest(d)

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
	d.Set("failure_reason", attachment.FailureReason)
	if attachment.ProposalDetails != nil {
		d.Set("direct_connect_gateway_association_proposal_state", attachment.ProposalDetails.ProposalState)
		d.Set("direct_connect_gateway_association_proposal_id", attachment.ProposalDetails.ProposalId)
		d.Set("direct_connect_gateway_association_proposal_created_at", attachment.ProposalDetails.CreatedAt)
		d.Set("direct_connect_gateway_association_proposal_updated_at", attachment.ProposalDetails.UpdatedAt)
	}

	return nil
}

func resourcePeeringGatewayAwsTgwAttachmentUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {

	// INIT
	api := alkira.NewPeeringGatewayAwsTgwAttachment(m.(*alkira.AlkiraClient))

	request, err := generateAwsTgwAttachmentRequest(d)

	if err != nil {
		return diag.FromErr(err)
	}

	// UPDATE
	_, err, _, _ = api.Update(d.Id(), request)
	if err != nil {
		return diag.FromErr(err)
	}

	// Only poll for proposal status when type is AWS_DIRECT_CONNECT_GATEWAY
	attachmentType := d.Get("type").(string)
	if attachmentType == "AWS_DIRECT_CONNECT_GATEWAY" {
		maxRetries := 60 // 5 minutes with 5-second intervals
		retryCount := 0

		for retryCount < maxRetries {
			resource, _, err := api.GetById(d.Id())
			if err != nil {
				return diag.Diagnostics{{
					Severity: diag.Warning,
					Summary:  "FAILED TO GET RESOURCE",
					Detail:   fmt.Sprintf("%s", err),
				}}
			}

			if resource.ProposalStatus != "" {
				switch resource.ProposalStatus {
				case "SUCCESS":
					break
				case "FAILED":
					return diag.Diagnostics{{
						Severity: diag.Error,
						Summary:  "Proposal failed",
						Detail:   fmt.Sprintf("ProposalStatus: FAILED, FailureReason: %s", resource.FailureReason),
					}}
				case "PENDING":
					time.Sleep(5 * time.Second)
					retryCount++
					continue
				default:
					time.Sleep(5 * time.Second)
					retryCount++
					continue
				}
				break
			}

			time.Sleep(5 * time.Second)
			retryCount++
		}

		if retryCount >= maxRetries {
			return diag.Diagnostics{{
				Severity: diag.Error,
				Summary:  "Timeout waiting for proposal to complete",
				Detail:   fmt.Sprintf("Timed out after %d retries (%d minutes)", maxRetries, maxRetries*5/60),
			}}
		}
	}

	return resourcePeeringGatewayAwsTgwAttachmentRead(ctx, d, m)
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

func generateAwsTgwAttachmentRequest(d *schema.ResourceData) (*alkira.PeeringGatewayAwsTgwAttachment, error) {
	if attachmentType, ok := d.Get("type").(string); ok && attachmentType == "AWS_DIRECT_CONNECT_GATEWAY" {
		request, err := generateDirectConnectAssociationTransitGatewayAwsAttachmentd(d)
		if err != nil {
			return nil, err
		}
		request.Type = attachmentType
		return request, nil
	} else {
		request, err := generatePeeringGatewayAwsTgwAttachmentRequest(d)
		if err != nil {
			return nil, err
		}
		if ok {
			request.Type = attachmentType
		}
		return request, nil
	}
}

// generatePeeringGatewayAwsTgwAttachmentRequest generate request
func generatePeeringGatewayAwsTgwAttachmentRequest(d *schema.ResourceData) (*alkira.PeeringGatewayAwsTgwAttachment, error) {

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

func generateDirectConnectAssociationTransitGatewayAwsAttachmentd(d *schema.ResourceData) (*alkira.PeeringGatewayAwsTgwAttachment, error) {

	request := &alkira.PeeringGatewayAwsTgwAttachment{
		Name:                       d.Get("name").(string),
		Description:                d.Get("description").(string),
		Requestor:                  d.Get("requestor").(string),
		PeerDirectConnectGatewayId: d.Get("peer_direct_connect_gateway_id").(string),
		PeerAllowedPrefixes:        convertTypeListToStringList(d.Get("peer_allowed_prefixes").([]any)),
		PeerAwsAccountId:           d.Get("peer_aws_account_id").(string),
		AwsTgwId:                   d.Get("peering_gateway_aws_tgw_id").(int),
	}

	return request, nil
}
