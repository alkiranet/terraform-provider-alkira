// Package alkira - Copyright (C) 2023-2025 Alkira Inc. All Rights Reserved.
package alkira

import (
	"context"
	"fmt"
	"strconv"

	"github.com/alkiranet/alkira-client-go/alkira"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func resourceAlkiraNetworkEntityScaleOptions() *schema.Resource {
	return &schema.Resource{
		Description:   "Scale Options are flexible configurations that elevate the capacity and performance characteristics of your network resource (any connector or a service) on Alkira's platform based on your specific needs. For example, you are experiencing traffic congestion with any of your exiting branch connectors, that too only on a particular segment, you can choose to define the scale options to add extra capacity to that connector on that segment. This can be done by specifying additional tunnels or additional nodes to the existing connector. \nScale options are made available only in certain scenarios when the existing connector or service is not meeting the required needs. \nUnderstanding scale options is crucial for planning and optimizing your network architecture on Alkira's platform. Choosing the right scale option ensures that your resources can handle the expected load.",
		CreateContext: resourceNetworkEntityScaleOptionsCreate,
		ReadContext:   resourceNetworkEntityScaleOptionsRead,
		UpdateContext: resourceNetworkEntityScaleOptionsUpdate,
		DeleteContext: resourceNetworkEntityScaleOptionsDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Description: "The name of the network entity scale options.",
				Type:        schema.TypeString,
				Required:    true,
			},
			"description": {
				Description: "The description of the network entity scale options.",
				Type:        schema.TypeString,
				Optional:    true,
			},
			"entity_id": {
				Description: "The entity ID of the connector or service.",
				Type:        schema.TypeString,
				Required:    true,
			},
			"entity_type": {
				Description: "The entity type.",
				Type:        schema.TypeString,
				Required:    true,
				ValidateFunc: validation.StringInSlice([]string{
					"CONNECTOR", "SERVICE"}, false),
			},
			"network_entity_id": {
				Description: "The network entity ID.",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"network_entity_type": {
				Description: "The network entity type.",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"segment_scale_options": {
				Description: "Segment Scale Options.",
				Type:        schema.TypeList,
				Required:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"additional_tunnels_per_node": {
							Description: "Number of additional Tunnels to be added per node. By default there is one tunnel per node. There are 2 nodes per connector or service on an average. Maximum tunnels are based on the limits allocated to a tenant. Either additionalTunnelsPerNode or additionalNodes either must be defined in a scale option.",
							Type:        schema.TypeInt,
							Optional:    true,
						},
						"additional_nodes": {
							Description: "Number of additional nodes to be added. By default there are 2 nodes per connector or service on an average. Maximum nodes are based on the limits allocated to a tenant. Either additionalTunnelsPerNode or additionalNodes either must be defined in a scale option.",
							Type:        schema.TypeInt,
							Optional:    true,
						},
						"segment_id": {
							Description: "Id of the segment for which custom scale is required. Segment Id is mandatory, a segment can occur only ones in connector segment scale options. For Service Scale the segment can repeat for unique segment and zone combination.",
							Type:        schema.TypeInt,
							Required:    true,
						},
						"zone_name": {
							Description: "optional field, if provided only tunnels associated with given zone would be scaled. Not applicable if scale options are defined for a connector.",
							Type:        schema.TypeString,
							Optional:    true,
						},
					},
				},
			},
			"doc_state": {
				Description: "The document state.",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"last_config_updated_at": {
				Description: "The last config updated at timestamp.",
				Type:        schema.TypeInt,
				Computed:    true,
			},
			"state": {
				Description: "The state of the network entity scale options.",
				Type:        schema.TypeString,
				Computed:    true,
			},
		},
	}
}

// convertEntityIdToInt converts entity_id string to int
// Returns 0 if conversion fails (will be caught by API validation)
func convertEntityIdToInt(entityId string) int {
	id, err := strconv.Atoi(entityId)
	if err != nil {
		return 0
	}
	return id
}

func resourceNetworkEntityScaleOptionsCreate(ctx context.Context, d *schema.ResourceData, m any) diag.Diagnostics {
	client := m.(*alkira.AlkiraClient)
	api := alkira.NewControllerScaleOptions(client)

	request, err := generateNetworkEntityScaleOptionsRequest(d)
	if err != nil {
		return diag.FromErr(err)
	}

	response, provState, err, valErr, provErr := api.Create(request)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(string(response.Id))

	// Handle validation errors
	if client.Validate && valErr != nil {
		var diags diag.Diagnostics
		// Try to read the resource to preserve any successfully created state
		readDiags := resourceNetworkEntityScaleOptionsRead(ctx, d, m)
		if readDiags.HasError() {
			diags = append(diags, readDiags...)
		}

		// Add the validation error
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "VALIDATION (CREATE) FAILED",
			Detail:   fmt.Sprintf("%s", valErr),
		})

		return diags
	}

	// Set provision state
	if client.Provision {
		d.Set("state", provState)

		if provState == "FAILED" {
			return diag.Diagnostics{{
				Severity: diag.Warning,
				Summary:  "PROVISION (CREATE) FAILED",
				Detail:   fmt.Sprintf("%s", provErr),
			}}
		}
	}

	return resourceNetworkEntityScaleOptionsRead(ctx, d, m)
}

func resourceNetworkEntityScaleOptionsRead(ctx context.Context, d *schema.ResourceData, m any) diag.Diagnostics {
	client := m.(*alkira.AlkiraClient)
	api := alkira.NewControllerScaleOptions(client)

	networkEntityScaleOptions, provState, err := api.GetById(d.Id())
	if err != nil {
		return diag.Diagnostics{{
			Severity: diag.Warning,
			Summary:  "FAILED TO GET RESOURCE",
			Detail:   fmt.Sprintf("%s", err),
		}}
	}

	d.Set("name", networkEntityScaleOptions.Name)
	d.Set("description", networkEntityScaleOptions.Description)
	d.Set("entity_id", strconv.Itoa(networkEntityScaleOptions.EntityId))
	d.Set("entity_type", networkEntityScaleOptions.EntityType)
	d.Set("network_entity_id", networkEntityScaleOptions.NetworkEntityId)
	d.Set("network_entity_sub_type", networkEntityScaleOptions.NetworkEntitySubType)
	d.Set("network_entity_type", networkEntityScaleOptions.NetworkEntityType)
	d.Set("doc_state", networkEntityScaleOptions.DocState)
	d.Set("last_config_updated_at", networkEntityScaleOptions.LastConfigUpdatedAt)

	var segmentScaleOptions []map[string]any
	for _, sso := range networkEntityScaleOptions.SegmentScaleOptions {
		ssoMap := map[string]any{
			"additional_tunnels_per_node": sso.AdditionalTunnelsPerNode,
			"additional_nodes":            sso.AdditionalNodes,
			"segment_id":                  sso.SegmentId,
			"zone_name":                   sso.ZoneName,
		}
		segmentScaleOptions = append(segmentScaleOptions, ssoMap)
	}
	d.Set("segment_scale_options", segmentScaleOptions)

	if client.Provision && provState != "" {
		d.Set("state", provState)
	}

	return nil
}

func resourceNetworkEntityScaleOptionsUpdate(ctx context.Context, d *schema.ResourceData, m any) diag.Diagnostics {
	client := m.(*alkira.AlkiraClient)
	api := alkira.NewControllerScaleOptions(client)

	request, err := generateNetworkEntityScaleOptionsRequest(d)
	if err != nil {
		return diag.FromErr(err)
	}

	provState, err, valErr, provErr := api.Update(d.Id(), request)
	if err != nil {
		return diag.FromErr(err)
	}

	// Handle validation errors
	if client.Validate && valErr != nil {
		var diags diag.Diagnostics
		// Try to read the resource to preserve current state
		readDiags := resourceNetworkEntityScaleOptionsRead(ctx, d, m)
		if readDiags.HasError() {
			diags = append(diags, readDiags...)
		}

		// Add the validation error
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "VALIDATION (UPDATE) FAILED",
			Detail:   fmt.Sprintf("%s", valErr),
		})

		return diags
	}

	// Set provision state
	if client.Provision {
		d.Set("state", provState)

		if provState == "FAILED" {
			return diag.Diagnostics{{
				Severity: diag.Warning,
				Summary:  "PROVISION (UPDATE) FAILED",
				Detail:   fmt.Sprintf("%s", provErr),
			}}
		}
	}

	return resourceNetworkEntityScaleOptionsRead(ctx, d, m)
}

func resourceNetworkEntityScaleOptionsDelete(ctx context.Context, d *schema.ResourceData, m any) diag.Diagnostics {
	client := m.(*alkira.AlkiraClient)
	api := alkira.NewControllerScaleOptions(client)

	provState, err, valErr, provErr := api.Delete(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId("")

	// Handle validation errors
	if client.Validate && valErr != nil {
		return diag.Diagnostics{{
			Severity: diag.Error,
			Summary:  "VALIDATION (DELETE) FAILED",
			Detail:   fmt.Sprintf("%s", valErr),
		}}
	}

	if client.Provision && provState != "SUCCESS" {
		return diag.Diagnostics{{
			Severity: diag.Warning,
			Summary:  "PROVISION (DELETE) FAILED",
			Detail:   fmt.Sprintf("%s", provErr),
		}}
	}

	return nil
}

func generateNetworkEntityScaleOptionsRequest(d *schema.ResourceData) (*alkira.ControllerScaleOptions, error) {
	var segmentScaleOptions []alkira.SegmentScaleOptions
	if v, ok := d.Get("segment_scale_options").([]any); ok {
		for _, item := range v {
			ssoMap := item.(map[string]any)
			segmentScaleOptions = append(segmentScaleOptions, alkira.SegmentScaleOptions{
				AdditionalTunnelsPerNode: ssoMap["additional_tunnels_per_node"].(int),
				AdditionalNodes:          ssoMap["additional_nodes"].(int),
				SegmentId:                ssoMap["segment_id"].(int),
				ZoneName:                 ssoMap["zone_name"].(string),
			})
		}
	}

	networkEntityScaleOptions := &alkira.ControllerScaleOptions{
		Name:                d.Get("name").(string),
		Description:         d.Get("description").(string),
		EntityId:            convertEntityIdToInt(d.Get("entity_id").(string)),
		EntityType:          d.Get("entity_type").(string),
		NetworkEntityType:   d.Get("network_entity_type").(string),
		SegmentScaleOptions: segmentScaleOptions,
	}

	return networkEntityScaleOptions, nil
}
