package alkira

import (
	"context"
	"fmt"
	"time"

	"github.com/alkiranet/alkira-client-go/alkira"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceAlkiraCxpPeeringGateways() *schema.Resource {
	return &schema.Resource{
		Description:   "Manage CXP Peering Gateways.",
		CreateContext: resourceAlkiraCxpPeeringGatewaysCreate,
		ReadContext:   resourceAlkiraCxpPeeringGatewaysRead,
		UpdateContext: resourceAlkiraCxpPeeringGatewaysUpdate,
		DeleteContext: resourceAlkiraCxpPeeringGatewaysDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"name": {
				Description: "The name of the Cxp Peering Gateway.",
				Type:        schema.TypeString,
				Required:    true,
			},
			"description": {
				Description: "Description of the Cxp Peering Gateway.",
				Type:        schema.TypeString,
				Optional:    true,
			},
			"cxp": {
				Description: "The CXP to which the Gateway is attached.",
				Type:        schema.TypeString,
				Required:    true,
			},
			"cloudProvider": {
				Description: "The cloud provider on which the gateway is created",
				Type:        schema.TypeString,
				Required:    true,
			},
			"cloudRegion": {
				Description: "The cloud region on which the ATH will be created",
				Type:        schema.TypeString,
				Required:    true,
			},
			"segment": {
				Description: "Name of the segment in which the gateway is created.",
				Type:        schema.TypeString,
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

func resourceAlkiraCxpPeeringGatewaysCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	// INIT
	api := alkira.NewCxpPeeringGateways(m.(*alkira.AlkiraClient))

	request, err := generateAlkiraCxpPeeringGatewaysRequest(d, m)
	if err != nil {
		return diag.FromErr(err)
	}

	// CREATE
	response, _, err, _ := api.Create(request)

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

	return resourceAlkiraCxpPeeringGatewaysRead(ctx, d, m)
}

func resourceAlkiraCxpPeeringGatewaysRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	// INIT
	api := alkira.NewCxpPeeringGateways(m.(*alkira.AlkiraClient))

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
	d.Set("cxp", resource.CloudProvider)
	d.Set("cloudProvider", resource.CloudProvider)
	d.Set("cloudRegion", resource.CloudRegion)
	d.Set("segment", resource.Segment)
	d.Set("state", resource.State)

	return nil
}

func resourceAlkiraCxpPeeringGatewaysUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	// INIT
	api := alkira.NewCxpPeeringGateways(m.(*alkira.AlkiraClient))

	request, err := generateAlkiraCxpPeeringGatewaysRequest(d, m)
	if err != nil {
		return diag.FromErr(err)
	}

	// UPDATE
	_, err, _ = api.Update(d.Id(), request)

	return nil
}

func resourceAlkiraCxpPeeringGatewaysDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	// INIT
	api := alkira.NewCxpPeeringGateways(m.(*alkira.AlkiraClient))

	// DELETE
	_, err, _ := api.Delete(d.Id())

	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId("")

	return nil
}

// generateAlkiraCxpPeeringGatewaysRequest generate request
func generateAlkiraCxpPeeringGatewaysRequest(d *schema.ResourceData, m interface{}) (*alkira.CxpPeeringGateways, error) {
	request := &alkira.CxpPeeringGateways{
		Name:          d.Get("name").(string),
		Description:   d.Get("description").(string),
		Cxp:           d.Get("cxp").(string),
		CloudRegion:   d.Get("cloudRegion").(string),
		CloudProvider: d.Get("cloudProivder").(string),
		Segment:       d.Get("segment").(string),
	}

	return request, nil
}
