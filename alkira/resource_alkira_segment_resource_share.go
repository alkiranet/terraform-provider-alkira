package alkira

import (
	"github.com/alkiranet/alkira-client-go/alkira"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func resourceAlkiraSegmentResourceShare() *schema.Resource {
	return &schema.Resource{
		Description: "Manages segment resource share.\n\n" +
			"Select resources to share between Resource End-A " +
			"in a segment and Resource End-B in another segment.",
		Create: resourceSegmentResourceShare,
		Read:   resourceSegmentResourceShareRead,
		Update: resourceSegmentResourceShareUpdate,
		Delete: resourceSegmentResourceShareDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Description: "The name of the segment resource share.",
				Type:        schema.TypeString,
				Required:    true,
			},
			"provision_state": {
				Description: "The provision state of the segment resource share.",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"service_ids": {
				Description: "The list of service IDs.",
				Type:        schema.TypeList,
				Elem:        &schema.Schema{Type: schema.TypeInt},
				Required:    true,
			},
			"designated_segment_id": {
				Description: "The designated segment ID.",
				Type:        schema.TypeString,
				Required:    true,
			},
			"end_a_segment_resource_ids": {
				Description: "The End-A segment resource IDs. All segment resources must be on the same segment.",
				Type:        schema.TypeList,
				Elem:        &schema.Schema{Type: schema.TypeInt},
				Required:    true,
			},
			"end_b_segment_resource_ids": {
				Description: "The End-B segment resource IDs. All segment resources must be on the same segment.",
				Type:        schema.TypeList,
				Elem:        &schema.Schema{Type: schema.TypeInt},
				Required:    true,
			},
			"end_a_route_limit": {
				Description: "The End-A route limit. The default value is `100`.",
				Type:        schema.TypeInt,
				Optional:    true,
				Default:     100,
			},
			"end_b_route_limit": {
				Description: "The End-B route limit. The default value is `100`.",
				Type:        schema.TypeInt,
				Optional:    true,
				Default:     100,
			},
			"traffic_direction": {
				Description: "Specify the direction in which traffic is orignated at " +
					"both Resource End-A and Resource End-B. The default value is `BIDIRECTIONAL`.",
				Type:         schema.TypeString,
				Optional:     true,
				Default:      "BIDIRECTIONAL",
				ValidateFunc: validation.StringInSlice([]string{"UNIDIRECTIONAL", "BIDIRECTIONAL"}, false),
			},
		},
	}
}

func resourceSegmentResourceShare(d *schema.ResourceData, m interface{}) error {

	api := alkira.NewSegmentResourceShare(m.(*alkira.AlkiraClient))

	// Construct request
	request, err := generateSegmentResourceShareRequest(d, m)

	if err != nil {
		return err
	}

	// Send create request
	response, provisionState, err := api.Create(request)

	if err != nil {
		return err
	}

	d.SetId(string(response.Id))
	d.Set("provision_state", provisionState)

	return resourceSegmentResourceShareRead(d, m)
}

func resourceSegmentResourceShareRead(d *schema.ResourceData, m interface{}) error {

	api := alkira.NewSegmentResourceShare(m.(*alkira.AlkiraClient))

	share, err := api.GetById(d.Id())

	if err != nil {
		return err
	}

	d.Set("name", share.Name)
	d.Set("service_ids", share.ServiceList)
	d.Set("designated_segement_id", share.DesignatedSegment)
	d.Set("end_a_segment_resource_ids", share.EndAResources)
	d.Set("end_b_segment_resource_ids", share.EndBResources)
	d.Set("end_a_route_limit", share.EndARouteLimit)
	d.Set("end_b_route_limit", share.EndBRouteLimit)
	d.Set("traffic_direction", share.Direction)

	return nil
}

func resourceSegmentResourceShareUpdate(d *schema.ResourceData, m interface{}) error {

	api := alkira.NewSegmentResourceShare(m.(*alkira.AlkiraClient))

	// Construct request
	share, err := generateSegmentResourceShareRequest(d, m)

	if err != nil {
		return err
	}

	// Send update request
	provisionState, err := api.Update(d.Id(), share)

	if err != nil {
		return err
	}

	d.Set("provision_state", provisionState)
	return resourceSegmentResourceShareRead(d, m)
}

func resourceSegmentResourceShareDelete(d *schema.ResourceData, m interface{}) error {

	api := alkira.NewSegmentResourceShare(m.(*alkira.AlkiraClient))

	provisionState, err := api.Delete(d.Id())

	if err != nil {
		return err
	}

	if provisionState != "SUCCESS" {
	}

	d.SetId("")
	return nil
}

// generateSegmentResourceShareRequest generate request for segment resource shared
func generateSegmentResourceShareRequest(d *schema.ResourceData, m interface{}) (*alkira.SegmentResourceShare, error) {

	// Get segment name
	segmentName, err := getSegmentNameById(d.Get("designated_segment_id").(string), m)

	if err != nil {
		return nil, err
	}

	request := &alkira.SegmentResourceShare{
		Name:              d.Get("name").(string),
		ServiceList:       convertTypeListToIntList(d.Get("service_ids").([]interface{})),
		DesignatedSegment: segmentName,
		EndAResources:     convertTypeListToIntList(d.Get("end_a_segment_resource_ids").([]interface{})),
		EndBResources:     convertTypeListToIntList(d.Get("end_b_segment_resource_ids").([]interface{})),
		EndARouteLimit:    d.Get("end_a_route_limit").(int),
		EndBRouteLimit:    d.Get("end_b_route_limit").(int),
		Direction:         d.Get("traffic_direction").(string),
	}

	return request, nil
}
