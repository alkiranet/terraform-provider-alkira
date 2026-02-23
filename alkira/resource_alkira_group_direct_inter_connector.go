package alkira

import (
	"context"
	"fmt"

	"github.com/alkiranet/alkira-client-go/alkira"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func resourceAlkiraDirectInterConnectorGroup() *schema.Resource {
	return &schema.Resource{
		Description:   "Provide direct inter-connector group resource.",
		CreateContext: resourceDirectInterConnectorGroup,
		ReadContext:   resourceDirectInterConnectorGroupRead,
		UpdateContext: resourceDirectInterConnectorGroupUpdate,
		DeleteContext: resourceDirectInterConnectorGroupDelete,
		CustomizeDiff: func(ctx context.Context, d *schema.ResourceDiff, m interface{}) error {
			client := m.(*alkira.AlkiraClient)

			old, _ := d.GetChange("provision_state")

			if client.Provision && old == "FAILED" {
				d.SetNew("provision_state", "SUCCESS")
			}

			return nil
		},
		Importer: &schema.ResourceImporter{
			StateContext: importWithReadValidation(resourceDirectInterConnectorGroupRead),
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Description: "The name of the group.",
				Type:        schema.TypeString,
				Required:    true,
			},
			"description": {
				Description: "The description of the group.",
				Type:        schema.TypeString,
				Optional:    true,
			},
			"segment_id": {
				Description: "The segment ID of the group.",
				Type:        schema.TypeString,
				Required:    true,
			},
			"cxp": {
				Description: "The CXP of the group.",
				Type:        schema.TypeString,
				Optional:    true,
			},
			"connector_provider_region": {
				Description: "The region of the connector.",
				Type:        schema.TypeString,
				Optional:    true,
			},
			"connector_type": {
				Description:  "The type of the connector.",
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringInSlice([]string{"AWS_VPC", "AZURE_VNET"}, false),
			},
			"provision_state": {
				Description: "The provisioning state of the resource.",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"azure_network_manager_id": {
				Description: "The Azure Virtual Network Manager's Alkira ID.",
				Type:        schema.TypeInt,
				Optional:    true,
			},
		},
	}
}

func resourceDirectInterConnectorGroup(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {

	// INIT
	client := m.(*alkira.AlkiraClient)
	api := alkira.NewInterConnectorCommunicationGroup(m.(*alkira.AlkiraClient))

	request, err := generateDirectInterConnectorGroupRequest(d, m)

	if err != nil {
		return diag.FromErr(err)
	}

	// Send create request
	response, provState, err, valErr, provErr := api.Create(request)

	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(string(response.Id))

	// Handle validation error
	if client.Validate && valErr != nil {
		var diags diag.Diagnostics
		readDiags := resourceDirectInterConnectorGroupRead(ctx, d, m)
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
		d.Set("provision_state", provState)

		if provErr != nil {
			return diag.Diagnostics{{
				Severity: diag.Warning,
				Summary:  "PROVISION (CREATE) FAILED",
				Detail:   fmt.Sprintf("%s", provErr),
			}}
		}
	}

	return resourceDirectInterConnectorGroupRead(ctx, d, m)
}

func resourceDirectInterConnectorGroupRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {

	// INIT
	client := m.(*alkira.AlkiraClient)
	api := alkira.NewInterConnectorCommunicationGroup(m.(*alkira.AlkiraClient))

	// Get
	group, provState, err := api.GetById(d.Id())

	if err != nil {
		return diag.Diagnostics{{
			Severity: diag.Warning,
			Summary:  "FAILED TO GET RESOURCE",
			Detail:   fmt.Sprintf("%s", err),
		}}
	}

	d.Set("name", group.Name)
	d.Set("description", group.Description)
	d.Set("cxp", group.Cxp)
	d.Set("connector_provider_region", group.ConnectorProviderRegion)
	d.Set("connector_type", group.ConnectorType)
	d.Set("azure_network_manager_id", group.VirtualNetworkManagerAzureId)

	// Get segment
	segmentId, err := getSegmentIdByName(group.Segment, m)

	if err != nil {
		return diag.FromErr(err)
	}
	d.Set("segment_id", segmentId)

	// Set provision state
	if client.Provision && provState != "" {
		d.Set("provision_state", provState)
	}

	return nil
}

func resourceDirectInterConnectorGroupUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {

	// INIT
	client := m.(*alkira.AlkiraClient)
	api := alkira.NewInterConnectorCommunicationGroup(m.(*alkira.AlkiraClient))

	// Construct request
	// Construct request
	request, err := generateDirectInterConnectorGroupRequest(d, m)

	if err != nil {
		return diag.FromErr(err)
	}

	// Send update request
	provState, err, valErr, provErr := api.Update(d.Id(), request)

	if err != nil {
		return diag.FromErr(err)
	}

	// Handle validation error
	if client.Validate && valErr != nil {
		var diags diag.Diagnostics
		readDiags := resourceDirectInterConnectorGroupRead(ctx, d, m)
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

func resourceDirectInterConnectorGroupDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {

	client := m.(*alkira.AlkiraClient)
	api := alkira.NewInterConnectorCommunicationGroup(m.(*alkira.AlkiraClient))

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

	if client.Provision && provState != "SUCCESS" {
		return diag.Diagnostics{{
			Severity: diag.Warning,
			Summary:  "PROVISION (DELETE) FAILED",
			Detail:   fmt.Sprintf("%s", provErr),
		}}
	}

	return nil
}
func generateDirectInterConnectorGroupRequest(d *schema.ResourceData, m interface{}) (*alkira.InterConnectorCommunicationGroup, error) {

	// Segment
	segmentName, err := getSegmentNameById(d.Get("segment_id").(string), m)
	if err != nil {
		return nil, err
	}
	// Construct request
	request := &alkira.InterConnectorCommunicationGroup{
		Name:                         d.Get("name").(string),
		Description:                  d.Get("description").(string),
		Segment:                      segmentName,
		Cxp:                          d.Get("cxp").(string),
		ConnectorProviderRegion:      d.Get("connector_provider_region").(string),
		ConnectorType:                d.Get("connector_type").(string),
		VirtualNetworkManagerAzureId: d.Get("azure_network_manager_id").(int),
	}
	return request, nil
}
