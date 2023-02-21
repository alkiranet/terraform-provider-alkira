package alkira

import (
	"context"
	"fmt"

	"github.com/alkiranet/alkira-client-go/alkira"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceAlkiraSegmentResource() *schema.Resource {
	return &schema.Resource{
		Description: "Manage segment resource.",
		Create:      resourceSegmentResource,
		Read:        resourceSegmentResourceRead,
		Update:      resourceSegmentResourceUpdate,
		Delete:      resourceSegmentResourceDelete,
		CustomizeDiff: func(ctx context.Context, d *schema.ResourceDiff, m interface{}) error {
			client := m.(*alkira.AlkiraClient)

			old, _ := d.GetChange("provision_state")

			if client.Provision == true && old == "FAILED" {
				d.SetNew("provision_state", "SUCCESS")
			}

			return nil
		},
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
				Description: "The provision state of the resource.",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"segment_id": {
				Description: "The segment ID.",
				Type:        schema.TypeString,
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

	client := m.(*alkira.AlkiraClient)
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

	// Set provision state
	if client.Provision == true {
		d.Set("provision_state", provisionState)
	}

	d.SetId(string(response.Id))
	return resourceSegmentResourceRead(d, m)
}

func resourceSegmentResourceRead(d *schema.ResourceData, m interface{}) error {

	client := m.(*alkira.AlkiraClient)
	api := alkira.NewSegmentResource(m.(*alkira.AlkiraClient))

	// Get resource
	resource, err := api.GetById(d.Id())

	if err != nil {
		return err
	}

	d.Set("name", resource.Name)
	d.Set("implicit_group_id", resource.GroupId)

	//
	// Get segemnt
	//
	segmentId, err := getSegmentIdByName(resource.Segment, m)

	if err != nil {
		return err
	}

	d.Set("segment_id", segmentId)

	//
	// Get Prefixes
	//
	var prefixes []map[string]interface{}

	for _, prefix := range resource.GroupPrefixes {
		i := map[string]interface{}{
			"group_id":       prefix.GroupId,
			"prefix_list_id": prefix.PrefixListId,
		}
		prefixes = append(prefixes, i)
	}

	d.Set("group_prefix", prefixes)

	// Set provision state
	_, provisionState, err := api.GetByName(d.Get("name").(string))

	if client.Provision == true && provisionState != "" {
		d.Set("provision_state", provisionState)
	}

	return nil
}

func resourceSegmentResourceUpdate(d *schema.ResourceData, m interface{}) error {

	client := m.(*alkira.AlkiraClient)
	api := alkira.NewSegmentResource(m.(*alkira.AlkiraClient))

	// Construct request
	request, err := generateSegmentResourceRequest(d, m)

	if err != nil {
		return err
	}

	// Send update request
	provisionState, err := api.Update(d.Id(), request)

	if err != nil {
		return err
	}

	// Set provision state
	if client.Provision == true {
		d.Set("provision_state", provisionState)
	}

	return nil
}

func resourceSegmentResourceDelete(d *schema.ResourceData, m interface{}) error {

	client := m.(*alkira.AlkiraClient)
	api := alkira.NewSegmentResource(m.(*alkira.AlkiraClient))

	provisionState, err := api.Delete(d.Id())

	if err != nil {
		return err
	}

	if client.Provision == true && provisionState != "SUCCESS" {
		return fmt.Errorf("failed to delete segment_resource %s, provision failed", d.Id())
	}

	d.SetId("")
	return nil
}

func generateSegmentResourceRequest(d *schema.ResourceData, m interface{}) (*alkira.SegmentResource, error) {
	//
	// Segment
	//
	segmentName, err := getSegmentNameById(d.Get("segment_id").(string), m)

	if err != nil {
		return nil, err
	}

	//
	// Group Prefix
	//
	groupPrefixes := expandSegmentResourceGroupPrefix(d.Get("group_prefix").(*schema.Set))

	// Assemble request
	resource := alkira.SegmentResource{
		Name:          d.Get("name").(string),
		Segment:       segmentName,
		GroupPrefixes: groupPrefixes,
	}

	return &resource, nil
}
