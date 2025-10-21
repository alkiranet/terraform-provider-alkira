// Package alkira - Copyright (C) 2023-2025 Alkira Inc. All Rights Reserved.
package alkira

import (
	"context"
	"fmt"

	"github.com/alkiranet/alkira-client-go/alkira"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func resourceAlkiraControllerScaleOptions() *schema.Resource {
	return &schema.Resource{
		Description:   "Manage Controller Scale Options.",
		CreateContext: resourceControllerScaleOptionsCreate,
		ReadContext:   resourceControllerScaleOptionsRead,
		UpdateContext: resourceControllerScaleOptionsUpdate,
		DeleteContext: resourceControllerScaleOptionsDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Description: "The name of the controller scale options.",
				Type:        schema.TypeString,
				Required:    true,
			},
			"description": {
				Description: "The description of the controller scale options.",
				Type:        schema.TypeString,
				Optional:    true,
			},
			"entity_id": {
				Description: "The entity ID.",
				Type:        schema.TypeInt,
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
				Required:    true,
			},
			"network_entity_sub_type": {
				Description: "The network entity sub type.",
				Type:        schema.TypeString,
				Required:    true,
			},
			"network_entity_type": {
				Description: "The network entity type.",
				Type:        schema.TypeString,
				Required:    true,
			},
			"segment_scale_options": {
				Description: "Segment Scale Options.",
				Type:        schema.TypeList,
				Required:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"additional_tunnels_per_node": {
							Description: "Additional tunnels per node.",
							Type:        schema.TypeInt,
							Required:    true,
						},
						"segment_id": {
							Description: "Segment ID.",
							Type:        schema.TypeInt,
							Required:    true,
						},
						"zone_name": {
							Description: "Zone name.",
							Type:        schema.TypeString,
							Required:    true,
						},
					},
				},
			},
			"doc_state": {
				Description: "The document state.",
				Type:        schema.TypeString,
				Required:    true,
				ValidateFunc: validation.StringInSlice([]string{
					"ADDED", "MARKED_FOR_DELETION", "DELETED"}, false),
			},
			"last_config_updated_at": {
				Description: "The last config updated at timestamp.",
				Type:        schema.TypeInt,
				Computed:    true,
			},
			"state": {
				Description: "The state of the controller scale options.",
				Type:        schema.TypeString,
				Required:    true,
				ValidateFunc: validation.StringInSlice([]string{
					"PENDING", "IN_PROGRESS",
					"SUCCESS", "FAILED", "SCHEDULED",
					"PARTIAL_SUCCESS", "INITIALIZING"}, false),
			},
		},
	}
}

func resourceControllerScaleOptionsCreate(ctx context.Context, d *schema.ResourceData, m any) diag.Diagnostics {
	client := m.(*alkira.AlkiraClient)
	api := alkira.NewControllerScaleOptions(client)

	request, err := generateControllerScaleOptionsRequest(d)
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
		readDiags := resourceControllerScaleOptionsRead(ctx, d, m)
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
	if client.Provision == true {
		d.Set("state", provState)

		if provState == "FAILED" {
			return diag.Diagnostics{{
				Severity: diag.Warning,
				Summary:  "PROVISION (CREATE) FAILED",
				Detail:   fmt.Sprintf("%s", provErr),
			}}
		}
	}

	return resourceControllerScaleOptionsRead(ctx, d, m)
}

func resourceControllerScaleOptionsRead(ctx context.Context, d *schema.ResourceData, m any) diag.Diagnostics {
	client := m.(*alkira.AlkiraClient)
	api := alkira.NewControllerScaleOptions(client)

	controllerScaleOptions, provState, err := api.GetById(d.Id())
	if err != nil {
		return diag.Diagnostics{{
			Severity: diag.Warning,
			Summary:  "FAILED TO GET RESOURCE",
			Detail:   fmt.Sprintf("%s", err),
		}}
	}

	d.Set("name", controllerScaleOptions.Name)
	d.Set("description", controllerScaleOptions.Description)
	d.Set("entity_id", controllerScaleOptions.EntityId)
	d.Set("entity_type", controllerScaleOptions.EntityType)
	d.Set("network_entity_id", controllerScaleOptions.NetworkEntityId)
	d.Set("network_entity_sub_type", controllerScaleOptions.NetworkEntitySubType)
	d.Set("network_entity_type", controllerScaleOptions.NetworkEntityType)
	d.Set("doc_state", controllerScaleOptions.DocState)
	d.Set("last_config_updated_at", controllerScaleOptions.LastConfigUpdatedAt)

	var segmentScaleOptions []map[string]any
	for _, sso := range controllerScaleOptions.SegmentScaleOptions {
		ssoMap := map[string]any{
			"additional_tunnels_per_node": sso.AdditionalTunnelsPerNode,
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

func resourceControllerScaleOptionsUpdate(ctx context.Context, d *schema.ResourceData, m any) diag.Diagnostics {
	client := m.(*alkira.AlkiraClient)
	api := alkira.NewControllerScaleOptions(client)

	request, err := generateControllerScaleOptionsRequest(d)
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
		readDiags := resourceControllerScaleOptionsRead(ctx, d, m)
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
	if client.Provision == true {
		d.Set("state", provState)

		if provState == "FAILED" {
			return diag.Diagnostics{{
				Severity: diag.Warning,
				Summary:  "PROVISION (UPDATE) FAILED",
				Detail:   fmt.Sprintf("%s", provErr),
			}}
		}
	}

	return resourceControllerScaleOptionsRead(ctx, d, m)
}

func resourceControllerScaleOptionsDelete(ctx context.Context, d *schema.ResourceData, m any) diag.Diagnostics {
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

	if client.Provision == true && provState != "SUCCESS" {
		return diag.Diagnostics{{
			Severity: diag.Warning,
			Summary:  "PROVISION (DELETE) FAILED",
			Detail:   fmt.Sprintf("%s", provErr),
		}}
	}

	return nil
}

func generateControllerScaleOptionsRequest(d *schema.ResourceData) (*alkira.ControllerScaleOptions, error) {
	var segmentScaleOptions []alkira.SegmentScaleOptions
	if v, ok := d.Get("segment_scale_options").([]any); ok {
		for _, item := range v {
			ssoMap := item.(map[string]any)
			segmentScaleOptions = append(segmentScaleOptions, alkira.SegmentScaleOptions{
				AdditionalTunnelsPerNode: int32(ssoMap["additional_tunnels_per_node"].(int)),
				SegmentId:                int64(ssoMap["segment_id"].(int)),
				ZoneName:                 ssoMap["zone_name"].(string),
			})
		}
	}

	controllerScaleOptions := &alkira.ControllerScaleOptions{
		Name:                 d.Get("name").(string),
		Description:          d.Get("description").(string),
		EntityId:             int64(d.Get("entity_id").(int)),
		EntityType:           d.Get("entity_type").(string),
		NetworkEntityId:      d.Get("network_entity_id").(string),
		NetworkEntitySubType: d.Get("network_entity_sub_type").(string),
		NetworkEntityType:    d.Get("network_entity_type").(string),
		State:                d.Get("state").(string),
		DocState:             d.Get("doc_state").(string),
		SegmentScaleOptions:  segmentScaleOptions,
	}

	return controllerScaleOptions, nil
}
