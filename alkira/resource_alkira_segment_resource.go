package alkira

import (
	"context"
	"fmt"

	"github.com/alkiranet/alkira-client-go/alkira"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceAlkiraSegmentResource() *schema.Resource {
	return &schema.Resource{
		Description:   "Manage segment resource.",
		CreateContext: resourceSegmentResource,
		ReadContext:   resourceSegmentResourceRead,
		UpdateContext: resourceSegmentResourceUpdate,
		DeleteContext: resourceSegmentResourceDelete,
		CustomizeDiff: func(ctx context.Context, d *schema.ResourceDiff, m interface{}) error {
			client := m.(*alkira.AlkiraClient)

			old, _ := d.GetChange("provision_state")

			if client.Provision == true && old == "FAILED" {
				d.SetNew("provision_state", "SUCCESS")
			}

			return nil
		},
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Description: "The name of the segment resource.",
				Type:        schema.TypeString,
				Required:    true,
			},
			"description": {
				Description: "The description of the segment resource.",
				Type:        schema.TypeString,
				Optional:    true,
			},
			"provision_state": {
				Description: "The provision state of the resource.",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"segment_id": {
				Description: "The segment ID.",
				Type:        schema.TypeString,
				Required:    true,
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
							Description: "The connector group ID associated " +
								"with the segment resource.",
							Type:     schema.TypeInt,
							Optional: true,
						},
						"prefix_list_id": {
							Description: "The Prefix List ID associated with " +
								"the segment resource.",
							Type:     schema.TypeInt,
							Optional: true,
						},
					},
				},
			},
		},
	}
}

func resourceSegmentResource(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {

	// INIT
	client := m.(*alkira.AlkiraClient)
	api := alkira.NewSegmentResource(m.(*alkira.AlkiraClient))

	// Construct request
	request, err := generateSegmentResourceRequest(d, m)

	if err != nil {
		return diag.FromErr(err)
	}

	// Send create request
	response, provState, err, provErr := api.Create(request)

	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(string(response.Id))

	// Set provision state
	if client.Provision == true {
		d.Set("provision_state", provState)

		if provErr != nil {
			return diag.Diagnostics{{
				Severity: diag.Warning,
				Summary:  "PROVISION (CREATE) FAILED",
				Detail:   fmt.Sprintf("%s", provErr),
			}}
		}
	}

	return resourceSegmentResourceRead(ctx, d, m)
}

func resourceSegmentResourceRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {

	client := m.(*alkira.AlkiraClient)
	api := alkira.NewSegmentResource(m.(*alkira.AlkiraClient))

	// Get resource
	resource, provState, err := api.GetById(d.Id())

	if err != nil {
		return diag.Diagnostics{{
			Severity: diag.Warning,
			Summary:  "FAILED TO GET RESOURCE",
			Detail:   fmt.Sprintf("%s", err),
		}}
	}

	d.Set("name", resource.Name)
	d.Set("description", resource.Description)
	d.Set("implicit_group_id", resource.GroupId)

	//
	// Get segemnt
	//
	segmentId, err := getSegmentIdByName(resource.Segment, m)

	if err != nil {
		return diag.FromErr(err)
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
	if client.Provision == true && provState != "" {
		d.Set("provision_state", provState)
	}

	return nil
}

func resourceSegmentResourceUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {

	client := m.(*alkira.AlkiraClient)
	api := alkira.NewSegmentResource(m.(*alkira.AlkiraClient))

	// Construct request
	request, err := generateSegmentResourceRequest(d, m)

	if err != nil {
		return diag.FromErr(err)
	}

	// Send update request
	provState, err, provErr := api.Update(d.Id(), request)

	if err != nil {
		return diag.FromErr(err)
	}

	// Set provision state
	if client.Provision == true {
		d.Set("provision_state", provState)

		if provErr != nil {
			return diag.Diagnostics{{
				Severity: diag.Warning,
				Summary:  "PROVISION (UPDATE) FAILED",
				Detail:   fmt.Sprintf("%s", provErr),
			}}
		}
	}

	return nil
}

func resourceSegmentResourceDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {

	client := m.(*alkira.AlkiraClient)
	api := alkira.NewSegmentResource(m.(*alkira.AlkiraClient))

	provState, err, provErr := api.Delete(d.Id())

	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId("")

	if client.Provision == true && provState != "SUCCESS" {
		return diag.Diagnostics{{
			Severity: diag.Warning,
			Summary:  "PROVISION (DELETE) FAILED",
			Detail:   fmt.Sprintf("%s", provErr),
		}}
	}

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
		Description:   d.Get("description").(string),
		Segment:       segmentName,
		GroupPrefixes: groupPrefixes,
	}

	return &resource, nil
}
