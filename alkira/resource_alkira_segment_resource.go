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
			"segment_id": {
				Description: "The segment ID.",
				Type:        schema.TypeInt,
				Optional:    true,
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
	client := m.(*alkira.AlkiraClient)

	resource, err := generateSegmentResourceRequest(d, m)

	if err != nil {
		return err
	}

	id, err := client.CreateSegmentResource(resource)

	if err != nil {
		return err
	}

	d.SetId(id)
	return resourceSegmentResourceRead(d, m)
}

func resourceSegmentResourceRead(d *schema.ResourceData, m interface{}) error {
	client := m.(*alkira.AlkiraClient)

	resource, err := client.GetSegmentResourceById(d.Id())

	if err != nil {
		return err
	}

	d.Set("name", resource.Name)

	segment, err := client.GetSegmentByName(resource.Segment)

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

	return nil
}

func resourceSegmentResourceUpdate(d *schema.ResourceData, m interface{}) error {
	client := m.(*alkira.AlkiraClient)

	resource, err := generateSegmentResourceRequest(d, m)

	if err != nil {
		return err
	}

	log.Printf("[INFO] Updating Segment Resource (%s)", d.Id())
	err = client.UpdateSegmentResource(d.Id(), resource)

	if err != nil {
		return err
	}

	return nil
}

func resourceSegmentResourceDelete(d *schema.ResourceData, m interface{}) error {
	client := m.(*alkira.AlkiraClient)

	log.Printf("[INFO] Deleting SegmentResource (%s)", d.Id())
	err := client.DeleteSegmentResource(d.Id())

	if err != nil {
		return err
	}

	log.Printf("[INFO] Deleted SegmentResource (%s)", d.Id())
	d.SetId("")
	return nil
}

func generateSegmentResourceRequest(d *schema.ResourceData, m interface{}) (*alkira.SegmentResource, error) {
	client := m.(*alkira.AlkiraClient)

	groupPrefixes := expandSegmentResourceGroupPrefix(d.Get("group_prefix").(*schema.Set))
	segment, err := client.GetSegmentById(strconv.Itoa(d.Get("segment_id").(int)))

	if err != nil {
		log.Printf("[ERROR] failed to get segment by Id: %d", d.Get("segment_id"))
		return nil, err
	}

	resource := alkira.SegmentResource{
		Name:          d.Get("name").(string),
		Segment:       segment.Name,
		GroupPrefixes: groupPrefixes,
	}

	return &resource, nil
}
