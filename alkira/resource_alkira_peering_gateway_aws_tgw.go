package alkira

import (
	"context"
	"fmt"
	"time"

	"github.com/alkiranet/alkira-client-go/alkira"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceAlkiraPeeringGatewayAwsTgw() *schema.Resource {
	return &schema.Resource{
		Description:   "Manage Peering Gateway AWS-TGW.",
		CreateContext: resourcePeeringGatewayAwsTgwCreate,
		ReadContext:   resourcePeeringGatewayAwsTgwRead,
		UpdateContext: resourcePeeringGatewayAwsTgwUpdate,
		DeleteContext: resourcePeeringGatewayAwsTgwDelete,
		Importer: &schema.ResourceImporter{
			StateContext: importWithReadValidation(resourcePeeringGatewayAwsTgwRead),
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
			"asn": {
				Description: "Initiator of transit gateway attachment.",
				Type:        schema.TypeInt,
				Required:    true,
			},
			"cxp": {
				Description: "The AWS region of the peer TGW.",
				Type:        schema.TypeString,
				Required:    true,
			},
			"aws_region": {
				Description: "AWS region of TGW.",
				Type:        schema.TypeString,
				Required:    true,
			},
			"state": {
				Description: "The state of the resource.",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"aws_tgw_id": {
				Description: "The ID of AWS TGW.",
				Type:        schema.TypeString,
				Computed:    true,
			},
		},
	}
}

func resourcePeeringGatewayAwsTgwCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {

	// INIT
	api := alkira.NewPeeringGatewayAwsTgw(m.(*alkira.AlkiraClient))

	request, err := generatePeeringGatewayAwsTgwRequest(d)

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

	for state != "ACTIVE" {
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

	return resourcePeeringGatewayAwsTgwRead(ctx, d, m)
}

func resourcePeeringGatewayAwsTgwRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {

	// INIT
	api := alkira.NewPeeringGatewayAwsTgw(m.(*alkira.AlkiraClient))

	resource, _, err := api.GetById(d.Id())

	if err != nil {
		return diag.Diagnostics{{
			Severity: diag.Warning,
			Summary:  "FAILED TO GET RESOURCE",
			Detail:   fmt.Sprintf("%s", err),
		}}
	}

	d.Set("name", resource.Name)
	d.Set("description", resource.Description)
	d.Set("asn", resource.Asn)
	d.Set("cxp", resource.Cxp)
	d.Set("aws_region", resource.AwsRegion)
	d.Set("state", resource.State)
	d.Set("aws_tgw_id", resource.ProviderTransitGatewayId)

	return nil
}

func resourcePeeringGatewayAwsTgwUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {

	// INIT
	api := alkira.NewPeeringGatewayAwsTgw(m.(*alkira.AlkiraClient))

	request, err := generatePeeringGatewayAwsTgwRequest(d)

	if err != nil {
		return diag.FromErr(err)
	}

	// UPDATE
	_, err, _, _ = api.Update(d.Id(), request)

	if err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func resourcePeeringGatewayAwsTgwDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {

	// INIT
	api := alkira.NewPeeringGatewayAwsTgw(m.(*alkira.AlkiraClient))

	// DELETE
	_, err, _, _ := api.Delete(d.Id())

	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId("")

	return nil
}

// generatePeeringGatewayAwsTgwRequest generate request
func generatePeeringGatewayAwsTgwRequest(d *schema.ResourceData) (*alkira.PeeringGatewayAwsTgw, error) {

	request := &alkira.PeeringGatewayAwsTgw{
		Name:        d.Get("name").(string),
		Description: d.Get("description").(string),
		Asn:         d.Get("asn").(int),
		Cxp:         d.Get("cxp").(string),
		AwsRegion:   d.Get("aws_region").(string),
	}

	return request, nil
}
