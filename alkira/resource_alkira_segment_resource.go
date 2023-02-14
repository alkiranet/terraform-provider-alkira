package alkira

import (
	"log"
	"strconv"

	"github.com/alkiranet/alkira-client-go/alkira"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceAlkiraSegmentResource() *schema.Resource {
	return &schema.Resource{
		Description: "Manage segment resource\n\n" +
			"To use this resource, you will need to use or create `alkira_group`, `alkira_segment` and " +
			"`alkira_policy_prefix_list`. " +
			"There could be multiple `group_prefix` section defined as needed. The `group_prefix` should " +
			"be defined like this:\n\n" +
			"* ANY -> ANY:  where `group_id` must be `-1` and prefix_list_id must be `-1`. When an " +
			"ANY -> ANY mapping is present then it should be the only mapping in the group_prefix\n\n" +
			"* EXPLICIT Group -> ANY:  where `group_id` must be the ID of group of type EXPLICIT and " +
			"prefix_list_id MUST be `-1`.\n\n" +
			"* IMPLICIT Group -> ANY: where group_id must be the ID of group of type IMPLICIT, this is " +
			"also known as a Connector Group and `prefix_list_id` must be `-1`. If an IMPLICIT group is " +
			"mapped to ANY `prefix_list_id`, then an IMPLICIT Group -> `prefix_list_id` must NOT be present " +
			"in `group_prefix`\n\n" +
			"* IMPLICIT Group -> PrefixList ID: where `group_id` must be the ID of group of type IMPLICIT " +
			"and `prefix_list_id` MUST be the ID of an existing `prefix_list_id`\n\n" +
			"* SERVICE Group -> ANY: where `group_id` must be the ID of group of type SERVICE and `prefix_list_id` " +
			"MUST be -1.",
		Create: resourceSegmentResource,
		Read:   resourceSegmentResourceRead,
		Update: resourceSegmentResourceUpdate,
		Delete: resourceSegmentResourceDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Description: "The name of the segment resource.",
				Type:        schema.TypeString,
				Required:    true,
			},
			"provision_state": {
				Description: "The provision state of the segment resource.",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"segment_id": {
				Description: "The segment ID.",
				Type:        schema.TypeInt,
				Optional:    true,
			},
			"implicit_group_id": {
				Description: "The ID of automatically created implicit group.",
				Type:        schema.TypeInt,
				Computed:    true,
			},
			"group_prefix": {
				Type:     schema.TypeSet,
				Required: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"group_id": {
							Description: "The connector group ID associated with segment resource.",
							Type:        schema.TypeInt,
							Optional:    true,
						},
						"prefix_list_id": {
							Description: "The Prefix List ID associated with segment resource.",
							Type:        schema.TypeInt,
							Optional:    true,
						},
					},
				},
			},
		},
	}
}

func resourceSegmentResource(d *schema.ResourceData, m interface{}) error {

	api := alkira.NewSegmentResource(m.(*alkira.AlkiraClient))

	// Construct request
	request, err := generateSegmentResourceRequest(d, m)

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

	return resourceSegmentResourceRead(d, m)
}

func resourceSegmentResourceRead(d *schema.ResourceData, m interface{}) error {

	api := alkira.NewSegmentResource(m.(*alkira.AlkiraClient))

	resource, err := api.GetById(d.Id())

	if err != nil {
		return err
	}

	d.Set("name", resource.Name)

	segmentApi := alkira.NewSegment(m.(*alkira.AlkiraClient))
	segment, _, err := segmentApi.GetByName(resource.Segment)

	if err != nil {
		return err
	}

	d.Set("segment_id", segment.Id)

	var prefixes []map[string]interface{}

	for _, prefix := range resource.GroupPrefixes {
		i := map[string]interface{}{
			"group_id":       prefix.GroupId,
			"prefix_list_id": prefix.PrefixListId,
		}
		prefixes = append(prefixes, i)
	}

	d.Set("group_prefix", prefixes)
	d.Set("implicit_group_id", resource.GroupId)

	return nil
}

func resourceSegmentResourceUpdate(d *schema.ResourceData, m interface{}) error {

	api := alkira.NewSegmentResource(m.(*alkira.AlkiraClient))

	// Construct request
	resource, err := generateSegmentResourceRequest(d, m)

	if err != nil {
		return err
	}

	// Send update request
	provisionState, err := api.Update(d.Id(), resource)

	if err != nil {
		return err
	}

	d.Set("provision_state", provisionState)
	return nil
}

func resourceSegmentResourceDelete(d *schema.ResourceData, m interface{}) error {

	api := alkira.NewSegmentResource(m.(*alkira.AlkiraClient))

	provisionState, err := api.Delete(d.Id())

	if err != nil {
		return err
	}

	if provisionState != "SUCCESS" {
	}

	d.SetId("")
	return nil
}

func generateSegmentResourceRequest(d *schema.ResourceData, m interface{}) (*alkira.SegmentResource, error) {

	groupPrefixes := expandSegmentResourceGroupPrefix(d.Get("group_prefix").(*schema.Set))

	segmentApi := alkira.NewSegment(m.(*alkira.AlkiraClient))
	segment, err := segmentApi.GetById(strconv.Itoa(d.Get("segment_id").(int)))

	if err != nil {
		log.Printf("[ERROR] failed to get segment by ID: %d", d.Get("segment_id"))
		return nil, err
	}

	resource := alkira.SegmentResource{
		Name:          d.Get("name").(string),
		Segment:       segment.Name,
		GroupPrefixes: groupPrefixes,
	}

	return &resource, nil
}
