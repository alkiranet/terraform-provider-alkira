package alkira

import (
	"context"
	"fmt"

	"github.com/alkiranet/alkira-client-go/alkira"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func resourceAlkiraConnectorAzureVhub() *schema.Resource {
	return &schema.Resource{
		Description:   "Manage Azure VHUB Connector.",
		CreateContext: resourceConnectorAzureVhubCreate,
		ReadContext:   resourceConnectorAzureVhubRead,
		UpdateContext: warnOnFailedStateUpdate(resourceConnectorAzureVhubUpdate),
		DeleteContext: resourceConnectorAzureVhubDelete,
		CustomizeDiff: func(ctx context.Context, d *schema.ResourceDiff, m interface{}) error {
			client := m.(*alkira.AlkiraClient)

			old, _ := d.GetChange("provision_state")

			if client.Provision && old == "FAILED" {
				d.SetNew("provision_state", "SUCCESS")
			}

			return nil
		},
		Importer: &schema.ResourceImporter{
			StateContext: importWithReadValidation(resourceConnectorAzureVhubRead),
		},

		Schema: map[string]*schema.Schema{
			"asn": {
				Description: "The BGP ASN of the Azure VHUB VPN Gateway. Always 65515.",
				Type:        schema.TypeInt,
				Computed:    true,
			},
			"billing_tag_ids": {
				Description: "Billing tags to be associated with " +
					"the resource. (see resource `alkira_billing_tag`).",
				Type:     schema.TypeSet,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeInt},
			},
			"credential_id": {
				Description: "ID of the Azure credential.",
				Type:        schema.TypeString,
				Required:    true,
			},
			"cxp": {
				Description: "The CXP where the connector should be provisioned.",
				Type:        schema.TypeString,
				Required:    true,
			},
			"description": {
				Description: "The description of the connector.",
				Type:        schema.TypeString,
				Optional:    true,
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
			"implicit_group_id": {
				Description: "The ID of implicit group automatically created " +
					"with the connector.",
				Type:     schema.TypeInt,
				Computed: true,
			},
			"name": {
				Description: "The name of the connector.",
				Type:        schema.TypeString,
				Required:    true,
			},
			"provision_state": {
				Description: "The provision state of the connector.",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"resource_group_name": {
				Description: "The Azure resource group name derived from the Virtual Hub ID.",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"scale_group_id": {
				Description: "The ID of the scale group associated with " +
					"the connector. Can only be set at create time and " +
					"cannot be changed after provisioning.",
				Type:     schema.TypeString,
				Optional: true,
			},
			"segment_id": {
				Description: "The ID of the segment associated with the connector.",
				Type:        schema.TypeString,
				Required:    true,
			},
			"size": {
				Description: "The size of the connector, one of `SMALL`, `MEDIUM`, " +
					"`LARGE`, `2LARGE`, `5LARGE`, `10LARGE`.",
				Type:     schema.TypeString,
				Required: true,
				ValidateFunc: validation.StringInSlice([]string{
					"SMALL", "MEDIUM", "LARGE", "2LARGE", "5LARGE", "10LARGE"}, false),
			},
			"subscription_id": {
				Description: "The Azure subscription ID derived from the Virtual Hub ID.",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"vhub_routing": {
				Description: "Routing options for the Azure VHUB connector.",
				Type:        schema.TypeList,
				Required:    true,
				MaxItems:    1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"route_import_mode": {
							Description: "The route import mode, one of " +
								"`ADVERTISE_DEFAULT_ROUTE`, `ADVERTISE_CUSTOM_PREFIX`.",
							Type:     schema.TypeString,
							Optional: true,
							ValidateFunc: validation.StringInSlice([]string{
								"ADVERTISE_DEFAULT_ROUTE",
								"ADVERTISE_CUSTOM_PREFIX"}, false),
						},
						"prefix_list_ids": {
							Description: "Prefix List IDs. Used when `route_import_mode` " +
								"is `ADVERTISE_CUSTOM_PREFIX`.",
							Type:     schema.TypeList,
							Optional: true,
							Elem:     &schema.Schema{Type: schema.TypeInt},
						},
					},
				},
			},
			"virtual_hub_id": {
				Description: "The ARM resource ID of the Azure Virtual Hub " +
					"(Microsoft.Network/virtualHubs).",
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"virtual_hub_name": {
				Description: "The name of the Azure Virtual Hub, resolved during provisioning.",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"virtual_wan_id": {
				Description: "The ARM resource ID of the parent Azure Virtual WAN, " +
					"resolved during provisioning.",
				Type:     schema.TypeString,
				Computed: true,
			},
			"vpn_gateway_id": {
				Description: "The ARM resource ID of the VPN Gateway attached to the " +
					"Virtual Hub, resolved during provisioning.",
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceConnectorAzureVhubCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {

	// INIT
	client := m.(*alkira.AlkiraClient)
	api := alkira.NewConnectorAzureVhub(m.(*alkira.AlkiraClient))

	request, err := generateConnectorAzureVhubRequest(d, m)

	if err != nil {
		return diag.FromErr(err)
	}

	// CREATE
	response, provState, err, valErr, provErr := api.Create(request)

	if err != nil {
		return diag.FromErr(err)
	}

	// Set states
	d.SetId(string(response.Id))

	// Handle validation error
	if client.Validate && valErr != nil {
		var diags diag.Diagnostics
		readDiags := resourceConnectorAzureVhubRead(ctx, d, m)
		if readDiags.HasError() {
			diags = append(diags, readDiags...)
		}

		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "VALIDATION (CREATE) FAILED",
			Detail:   fmt.Sprintf("%s", valErr),
		})

		return diags
	}

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

	return resourceConnectorAzureVhubRead(ctx, d, m)
}

func resourceConnectorAzureVhubRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {

	client := m.(*alkira.AlkiraClient)
	api := alkira.NewConnectorAzureVhub(m.(*alkira.AlkiraClient))

	connector, provState, err := api.GetById(d.Id())

	if err != nil {
		return diag.Diagnostics{{
			Severity: diag.Warning,
			Summary:  "FAILED TO GET RESOURCE",
			Detail:   fmt.Sprintf("%s", err),
		}}
	}

	d.Set("asn", connector.ASN)
	d.Set("billing_tag_ids", connector.BillingTags)
	d.Set("credential_id", connector.CredentialId)
	d.Set("cxp", connector.CXP)
	d.Set("description", connector.Description)
	d.Set("enabled", connector.Enabled)
	d.Set("group", connector.Group)
	d.Set("implicit_group_id", connector.ImplicitGroupId)
	d.Set("name", connector.Name)
	d.Set("resource_group_name", connector.ResourceGroupName)
	d.Set("scale_group_id", connector.ScaleGroupId)
	d.Set("size", connector.Size)
	d.Set("subscription_id", connector.SubscriptionId)
	d.Set("virtual_hub_id", connector.VirtualHubId)
	d.Set("virtual_hub_name", connector.VirtualHubName)
	d.Set("virtual_wan_id", connector.VirtualWanId)
	d.Set("vpn_gateway_id", connector.VpnGatewayId)

	setVhubRouting(d, connector.VhubRouting)

	// Get segment
	numOfSegments := len(connector.Segments)
	if numOfSegments == 1 {
		segmentId, err := getSegmentIdByName(connector.Segments[0], m)

		if err != nil {
			return diag.FromErr(err)
		}
		d.Set("segment_id", segmentId)
	} else {
		return diag.FromErr(fmt.Errorf("failed to find segment"))
	}

	// Set provision state
	if client.Provision && provState != "" {
		d.Set("provision_state", provState)
	}

	return nil
}

func resourceConnectorAzureVhubUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {

	// INIT
	client := m.(*alkira.AlkiraClient)
	api := alkira.NewConnectorAzureVhub(m.(*alkira.AlkiraClient))

	request, err := generateConnectorAzureVhubRequest(d, m)

	if err != nil {
		return diag.FromErr(err)
	}

	// UPDATE
	provState, err, valErr, provErr := api.Update(d.Id(), request)

	if err != nil {
		return diag.FromErr(err)
	}

	// Handle validation error
	if client.Validate && valErr != nil {
		var diags diag.Diagnostics
		readDiags := resourceConnectorAzureVhubRead(ctx, d, m)
		if readDiags.HasError() {
			diags = append(diags, readDiags...)
		}

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

	return resourceConnectorAzureVhubRead(ctx, d, m)
}

func resourceConnectorAzureVhubDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {

	// INIT
	client := m.(*alkira.AlkiraClient)
	api := alkira.NewConnectorAzureVhub(m.(*alkira.AlkiraClient))

	// DELETE
	_, err, valErr, provErr := api.Delete(d.Id())

	if err != nil {
		name, _ := d.GetOk("name")
		if nameStr, ok := name.(string); ok && nameStr != "" {
			return diag.FromErr(fmt.Errorf("%w alkira_connector_azure_vhub (name=%q id=%s)", err, nameStr, d.Id()))
		}
		return diag.FromErr(fmt.Errorf("%w alkira_connector_azure_vhub (id=%s)", err, d.Id()))
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

	if client.Provision && provErr != nil {
		return diag.Diagnostics{{
			Severity: diag.Warning,
			Summary:  "PROVISION (DELETE) FAILED",
			Detail:   fmt.Sprintf("%s", provErr),
		}}
	}

	return nil
}

// generateConnectorAzureVhubRequest generate request for connector-azure-vhub
func generateConnectorAzureVhubRequest(d *schema.ResourceData, m interface{}) (*alkira.ConnectorAzureVhub, error) {

	//
	// Segment
	//
	segmentName, err := getSegmentNameById(d.Get("segment_id").(string), m)

	if err != nil {
		return nil, err
	}

	//
	// VhubRouting
	//
	vhubRouting := constructVhubRouting(d)

	// Assemble request
	request := &alkira.ConnectorAzureVhub{
		BillingTags:  convertTypeSetToIntList(d.Get("billing_tag_ids").(*schema.Set)),
		CXP:          d.Get("cxp").(string),
		CredentialId: d.Get("credential_id").(string),
		Description:  d.Get("description").(string),
		Enabled:      d.Get("enabled").(bool),
		Group:        d.Get("group").(string),
		Name:         d.Get("name").(string),
		ScaleGroupId: d.Get("scale_group_id").(string),
		Segments:     []string{segmentName},
		Size:         d.Get("size").(string),
		VirtualHubId: d.Get("virtual_hub_id").(string),
		VhubRouting:  vhubRouting,
	}

	return request, nil
}
