package alkira

import (
	"context"
	"fmt"
	"time"

	"github.com/alkiranet/alkira-client-go/alkira"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceAlkiraCxpPeeringGateway() *schema.Resource {
	return &schema.Resource{
		Description:   "Manage CXP Peering Gateways.",
		CreateContext: resourceAlkiraCxpPeeringGatewayCreate,
		ReadContext:   resourceAlkiraCxpPeeringGatewayRead,
		UpdateContext: resourceAlkiraCxpPeeringGatewayUpdate,
		DeleteContext: resourceAlkiraCxpPeeringGatewayDelete,
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

			// TODO: change cloud_provider to be a Required value when more cloud providers are added and remove the default value.
			"cloud_provider": {
				Description: "The cloud provider on which the gateway is created",
				Type:        schema.TypeString,
				// Required:    true,
				Optional: true,
				Default:  "AZURE",
			},
			"cloud_region": {
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

func resourceAlkiraCxpPeeringGatewayCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	// INIT
	api := alkira.NewCxpPeeringGateway(m.(*alkira.AlkiraClient))

	request, err := generateAlkiraCxpPeeringGatewayRequest(d, m)
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

	return resourceAlkiraCxpPeeringGatewayRead(ctx, d, m)
}

func resourceAlkiraCxpPeeringGatewayRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	// INIT
	api := alkira.NewCxpPeeringGateway(m.(*alkira.AlkiraClient))

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
	d.Set("cloud_provider", resource.CloudProvider)
	d.Set("cloud_region", resource.CloudRegion)
	d.Set("segment", resource.Segment)
	d.Set("state", resource.State)

	return nil
}

func resourceAlkiraCxpPeeringGatewayUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	// INIT
	api := alkira.NewCxpPeeringGateway(m.(*alkira.AlkiraClient))

	request, err := generateAlkiraCxpPeeringGatewayRequest(d, m)
	if err != nil {
		return diag.FromErr(err)
	}

	// UPDATE
	_, err, _ = api.Update(d.Id(), request)

	return nil
}

func resourceAlkiraCxpPeeringGatewayDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	// INIT
	api := alkira.NewCxpPeeringGateway(m.(*alkira.AlkiraClient))

	// DELETE
	_, err, _ := api.Delete(d.Id())

	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId("")

	return nil
}

// generateAlkiraCxpPeeringGatewayRequest generate request
func generateAlkiraCxpPeeringGatewayRequest(d *schema.ResourceData, m interface{}) (*alkira.CxpPeeringGateway, error) {
	request := &alkira.CxpPeeringGateway{
		Name:          d.Get("name").(string),
		Description:   d.Get("description").(string),
		Cxp:           d.Get("cxp").(string),
		CloudRegion:   d.Get("cloud_region").(string),
		CloudProvider: d.Get("cloud_proivder").(string),
		Segment:       d.Get("segment").(string),
	}

	return request, nil
}
