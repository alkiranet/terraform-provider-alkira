package alkira

import (
	"context"
	"fmt"
	"log"

	"github.com/alkiranet/alkira-client-go/alkira"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func resourceAlkiraConnectorAzureVnetThirdParty() *schema.Resource {
	return &schema.Resource{
		Description:   "Manage Azure VNET Third Party Connector.",
		CreateContext: resourceConnectorAzureVnetThirdPartyCreate,
		ReadContext:   resourceConnectorAzureVnetThirdPartyRead,
		UpdateContext: resourceConnectorAzureVnetThirdPartyUpdate,
		DeleteContext: resourceConnectorAzureVnetThirdPartyDelete,
		CustomizeDiff: func(ctx context.Context, d *schema.ResourceDiff, m interface{}) error {
			client := m.(*alkira.AlkiraClient)

			old, _ := d.GetChange("provision_state")

			if client.Provision && old == "FAILED" {
				d.SetNew("provision_state", "SUCCESS")
			}

			return nil
		},
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Description: "The name of the connector.",
				Type:        schema.TypeString,
				Required:    true,
			},
			"description": {
				Description: "The description of the connector.",
				Type:        schema.TypeString,
				Optional:    true,
			},
			"cxp": {
				Description: "The CXP where the connector should be provisioned.",
				Type:        schema.TypeString,
				Required:    true,
			},
			"enabled": {
				Description: "Is the connector enabled. Default is `true`.",
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     true,
			},
			"group": {
				Description: "The group of the connector.",
				Type:        schema.TypeString,
				Optional:    true,
			},
			"segment_id": {
				Description: "The ID of the segment associated with the connector.",
				Type:        schema.TypeString,
				Required:    true,
			},
			"size": {
				Description: "The size of the connector, one of `XSMALL`, `SMALL`, `MEDIUM`, `LARGE`, `XLARGE`, `2XLARGE`.",
				Type:        schema.TypeString,
				Required:    true,
				ValidateFunc: validation.StringInSlice([]string{
					"XSMALL", "SMALL", "MEDIUM", "LARGE", "XLARGE", "2XLARGE"}, false),
			},
			"azure_vnet_third_party_connector_attachment_id": {
				Description: "The ID of the Azure VNET Third Party Connector Attachment.",
				Type:        schema.TypeInt,
				Required:    true,
			},
			"billing_tag_ids": {
				Description: "Billing tags to be associated with the resource.",
				Type:        schema.TypeSet,
				Optional:    true,
				Elem:        &schema.Schema{Type: schema.TypeInt},
			},
			"static_route_prefix_list_ids": {
				Description: "Policy Prefix List IDs to be associated with the connector's static routes.",
				Type:        schema.TypeSet,
				Optional:    true,
				Elem:        &schema.Schema{Type: schema.TypeInt},
			},
			"implicit_group_id": {
				Description: "The ID of implicit group automatically created with the connector.",
				Type:        schema.TypeInt,
				Computed:    true,
			},
			"provision_state": {
				Description: "The provision state of the connector.",
				Type:        schema.TypeString,
				Computed:    true,
			},
		},
	}
}

func resourceConnectorAzureVnetThirdPartyCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*alkira.AlkiraClient)
	api := alkira.NewAzureVnetThirdPartyConnector(client)

	request, err := generateConnectorAzureVnetThirdPartyRequest(d, m)
	if err != nil {
		return diag.FromErr(err)
	}

	resource, provState, err, valErr, provErr := api.Create(request)

	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(string(resource.Id))

	// Handle validation error
	if client.Validate && valErr != nil {
		var diags diag.Diagnostics
		readDiags := resourceConnectorAzureVnetThirdPartyRead(ctx, d, m)
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
		d.Set("provision_state", provState)

		if provErr != nil {
			return diag.Diagnostics{{
				Severity: diag.Warning,
				Summary:  "PROVISION (CREATE) FAILED",
				Detail:   fmt.Sprintf("%s", provErr),
			}}
		}
	}

	return resourceConnectorAzureVnetThirdPartyRead(ctx, d, m)
}

func resourceConnectorAzureVnetThirdPartyRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*alkira.AlkiraClient)
	api := alkira.NewAzureVnetThirdPartyConnector(client)

	connector, provState, err := api.GetById(d.Id())

	if err != nil {
		return diag.Diagnostics{{
			Severity: diag.Warning,
			Summary:  "FAILED TO GET RESOURCE",
			Detail:   fmt.Sprintf("%s", err),
		}}
	}

	d.Set("name", connector.Name)
	d.Set("description", connector.Description)
	d.Set("cxp", connector.CXP)
	d.Set("enabled", connector.Enabled)
	d.Set("group", connector.Group)
	d.Set("size", connector.Size)
	d.Set("azure_vnet_third_party_connector_attachment_id", connector.AzureVnetThirdPartyConnectorAttachmentId)
	d.Set("implicit_group_id", connector.ImplicitGroupId)

	numOfSegments := len(connector.Segments)
	if numOfSegments == 1 {
		segmentId, err := getSegmentIdByName(connector.Segments[0], m)
		if err != nil {
			return diag.FromErr(err)
		}
		d.Set("segment_id", segmentId)
	}

	// Set provision state
	if client.Provision == true && provState != "" {
		d.Set("provision_state", provState)
	}

	// Set billing tags and static routes
	d.Set("billing_tag_ids", connector.BillingTags)
	d.Set("static_route_prefix_list_ids", connector.StaticRoutes)

	return nil
}

func resourceConnectorAzureVnetThirdPartyUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*alkira.AlkiraClient)
	api := alkira.NewAzureVnetThirdPartyConnector(client)

	request, err := generateConnectorAzureVnetThirdPartyRequest(d, m)
	if err != nil {
		return diag.FromErr(err)
	}

	provState, err, valErr, provErr := api.Update(d.Id(), request)

	if err != nil {
		return diag.FromErr(err)
	}

	// Handle validation error
	if client.Validate && valErr != nil {
		var diags diag.Diagnostics
		readDiags := resourceConnectorAzureVnetThirdPartyRead(ctx, d, m)
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
		d.Set("provision_state", provState)

		if provErr != nil {
			return diag.Diagnostics{{
				Severity: diag.Warning,
				Summary:  "PROVISION (UPDATE) FAILED",
				Detail:   fmt.Sprintf("%s", provErr),
			}}
		}
	}

	return resourceConnectorAzureVnetThirdPartyRead(ctx, d, m)
}

func resourceConnectorAzureVnetThirdPartyDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*alkira.AlkiraClient)
	api := alkira.NewAzureVnetThirdPartyConnector(client)

	provState, err, valErr, provErr := api.Delete(d.Id())

	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId("")

	// Handle validation error
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

func generateConnectorAzureVnetThirdPartyRequest(d *schema.ResourceData, m interface{}) (*alkira.AzureVnetThirdPartyConnector, error) {
	segmentName, err := getSegmentNameById(d.Get("segment_id").(string), m)
	if err != nil {
		return nil, err
	}

	billingTags := convertTypeSetToIntList(d.Get("billing_tag_ids").(*schema.Set))
	staticRoutes := convertTypeSetToIntList(d.Get("static_route_prefix_list_ids").(*schema.Set))

	request := &alkira.AzureVnetThirdPartyConnector{
		Name:                                     d.Get("name").(string),
		Description:                              d.Get("description").(string),
		CXP:                                      d.Get("cxp").(string),
		Enabled:                                  d.Get("enabled").(bool),
		Group:                                    d.Get("group").(string),
		Segments:                                 []string{segmentName},
		Size:                                     d.Get("size").(string),
		AzureVnetThirdPartyConnectorAttachmentId: d.Get("azure_vnet_third_party_connector_attachment_id").(int),
		BillingTags:                              billingTags,
		StaticRoutes:                             staticRoutes,
	}

	log.Printf("[DEBUG] Azure VNET Third Party Connector Request: %+v", request)

	return request, nil
}
