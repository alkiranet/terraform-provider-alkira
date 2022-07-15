package alkira

import (
	"log"
	"strconv"

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
			"service_ids": {
				Description: "The list of service IDs.",
				Type:        schema.TypeList,
				Elem:        &schema.Schema{Type: schema.TypeInt},
				Required:    true,
			},
			"designated_segment_id": {
				Description: "The designated segment ID.",
				Type:        schema.TypeInt,
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
	client := m.(*alkira.AlkiraClient)

	share, err := generateSegmentResourceShareRequest(d, m)

	if err != nil {
		return err
	}

	id, err := client.CreateSegmentResourceShare(share)

	if err != nil {
		return err
	}

	d.SetId(id)
	return resourceSegmentResourceShareRead(d, m)
}

func resourceSegmentResourceShareRead(d *schema.ResourceData, m interface{}) error {
	client := m.(*alkira.AlkiraClient)

	share, err := client.GetSegmentResourceShareById(d.Id())

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
	client := m.(*alkira.AlkiraClient)

	share, err := generateSegmentResourceShareRequest(d, m)

	if err != nil {
		return err
	}

	return client.UpdateSegmentResourceShare(d.Id(), share)
}

func resourceSegmentResourceShareDelete(d *schema.ResourceData, m interface{}) error {
	client := m.(*alkira.AlkiraClient)

	log.Printf("[INFO] Deleting segment_resource_share: %s", d.Id())
	err := client.DeleteSegmentResourceShare(d.Id())

	return err
}

// generateSegmentResourceShareRequest generate request for segment resource shared
func generateSegmentResourceShareRequest(d *schema.ResourceData, m interface{}) (*alkira.SegmentResourceShare, error) {
	client := m.(*alkira.AlkiraClient)

	serviceList := convertTypeListToIntList(d.Get("service_ids").([]interface{}))
	endAResources := convertTypeListToIntList(d.Get("end_a_segment_resource_ids").([]interface{}))
	endBResources := convertTypeListToIntList(d.Get("end_b_segment_resource_ids").([]interface{}))

	segmentId := d.Get("designated_segment_id").(int)
	segment, err := client.GetSegmentById(strconv.Itoa(segmentId))

	if err != nil {
		log.Printf("[ERROR] failed to get segment by ID: %d", segmentId)
		return nil, err
	}

	request := &alkira.SegmentResourceShare{
		Name:              d.Get("name").(string),
		ServiceList:       serviceList,
		DesignatedSegment: segment.Name,
		EndAResources:     endAResources,
		EndBResources:     endBResources,
		EndARouteLimit:    d.Get("end_a_route_limit").(int),
		EndBRouteLimit:    d.Get("end_b_route_limit").(int),
		Direction:         d.Get("traffic_direction").(string),
	}

	return request, nil
}
