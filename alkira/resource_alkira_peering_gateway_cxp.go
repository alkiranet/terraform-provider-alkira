package alkira

import (
	"context"
	"fmt"
	"time"

	"github.com/alkiranet/alkira-client-go/alkira"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceAlkiraPeeringGatewayCxp() *schema.Resource {
	return &schema.Resource{
		Description:   "Manage CXP Peering Gateways.",
		CreateContext: resourceAlkiraPeeringGatewayCxpCreate,
		ReadContext:   resourceAlkiraPeeringGatewayCxpRead,
		UpdateContext: resourceAlkiraPeeringGatewayCxpUpdate,
		DeleteContext: resourceAlkiraPeeringGatewayCxpDelete,
		Importer: &schema.ResourceImporter{
			StateContext: importWithReadValidation(resourceAlkiraPeeringGatewayCxpRead),
		},
		Schema: map[string]*schema.Schema{
			"name": {
				Description: "The name of the Peering Gateway.",
				Type:        schema.TypeString,
				Required:    true,
			},
			"description": {
				Description: "Description of the resource.",
				Type:        schema.TypeString,
				Optional:    true,
			},
			"cxp": {
				Description: "The CXP to which the Gateway is attached.",
				Type:        schema.TypeString,
				Required:    true,
			},

			"cloud_provider": {
				Description: "The cloud provider where this resource will be " +
					"created. The default value is `AZURE` and only `AZURE` " +
					"is supported for now.",
				Type:     schema.TypeString,
				Optional: true,
				Default:  "AZURE",
			},
			"cloud_region": {
				Description: "The region of the specified cloud provider on " +
					"which the resource should be created. E.g. if " +
					"`cloud_provider` is `AZURE`, the region could be like " +
					"`eastus`.",
				Type:     schema.TypeString,
				Required: true,
			},
			"segment_id": {
				Description: "The ID of the segment that is associated with " +
					"the resource.",
				Type:     schema.TypeString,
				Required: true,
			},
			"state": {
				Description: "The state of the resource.",
				Type:        schema.TypeString,
				Computed:    true,
			},
		},
	}
}

func resourceAlkiraPeeringGatewayCxpCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	// INIT
	api := alkira.NewPeeringGatewayCxp(m.(*alkira.AlkiraClient))

	request, err := generateAlkiraCxpPeeringGatewayRequest(d, m)
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

	return resourceAlkiraPeeringGatewayCxpRead(ctx, d, m)
}

func resourceAlkiraPeeringGatewayCxpRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	// INIT
	api := alkira.NewPeeringGatewayCxp(m.(*alkira.AlkiraClient))

	resource, _, err := api.GetById(d.Id())
	if err != nil {
		return diag.Diagnostics{{
			Severity: diag.Warning,
			Summary:  "FAILED TO GET RESOURCE",
			Detail:   fmt.Sprintf("%s", err),
		}}
	}
	segmentId, err := getSegmentIdByName(resource.Segment, m)
	if err != nil {
		return diag.FromErr(err)
	}

	d.Set("name", resource.Name)
	d.Set("description", resource.Description)
	d.Set("cxp", resource.Cxp)
	d.Set("cloud_provider", resource.CloudProvider)
	d.Set("cloud_region", resource.CloudRegion)
	d.Set("segment_id", segmentId)
	d.Set("state", resource.State)

	return nil
}

func resourceAlkiraPeeringGatewayCxpUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	// checking if only description field is changed
	if d.HasChanges("name", "cloud_provider", "cloud_region", "segment_id") {
		// return error if any field except description is changed.
		return diag.Diagnostics{{
			Severity: diag.Error,
			Summary:  "Invalid Update!",
			Detail:   "Only the description field can be updated.",
		}}
	}

	// INIT
	api := alkira.NewPeeringGatewayCxp(m.(*alkira.AlkiraClient))

	request, err := generateAlkiraCxpPeeringGatewayRequest(d, m)
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

func resourceAlkiraPeeringGatewayCxpDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	// INIT
	api := alkira.NewPeeringGatewayCxp(m.(*alkira.AlkiraClient))

	// DELETE
	_, err, _, _ := api.Delete(d.Id())

	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId("")

	return nil
}

// generateAlkiraCxpPeeringGatewayRequest generate request
func generateAlkiraCxpPeeringGatewayRequest(d *schema.ResourceData, m interface{}) (*alkira.PeeringGatewayCxp, error) {
	segmentName, err := getSegmentNameById(d.Get("segment_id").(string), m)
	if err != nil {
		return nil, err
	}
	request := &alkira.PeeringGatewayCxp{
		Name:          d.Get("name").(string),
		Description:   d.Get("description").(string),
		Cxp:           d.Get("cxp").(string),
		CloudRegion:   d.Get("cloud_region").(string),
		CloudProvider: d.Get("cloud_provider").(string),
		Segment:       segmentName,
	}

	return request, nil
}
