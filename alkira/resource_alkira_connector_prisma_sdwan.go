package alkira

import (
	"context"
	"fmt"

	"github.com/alkiranet/alkira-client-go/alkira"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceAlkiraConnectorPrismaSDWAN() *schema.Resource {
	return &schema.Resource{
		Description:   "Manage Prisma SD-WAN Connector.",
		CreateContext: resourceConnectorPrismaSDWANCreate,
		ReadContext:   resourceConnectorPrismaSDWANRead,
		UpdateContext: resourceConnectorPrismaSDWANUpdate,
		DeleteContext: resourceConnectorPrismaSDWANDelete,
		CustomizeDiff: func(ctx context.Context, d *schema.ResourceDiff, m interface{}) error {
			client := m.(*alkira.AlkiraClient)

			old, _ := d.GetChange("provision_state")

			if client.Provision && old == "FAILED" {
				d.SetNew("provision_state", "SUCCESS")
			}

			return nil
		},
		Importer: &schema.ResourceImporter{
			StateContext: importWithReadValidation(resourceConnectorPrismaSDWANRead),
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
			"billing_tag_ids": {
				Description: "Billing tags to be associated with " +
					"the resource. (see resource `alkira_billing_tag`).",
				Type:     schema.TypeSet,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeInt},
			},
			"cxp": {
				Description: "The CXP where the connector should be " +
					"provisioned.",
				Type:     schema.TypeString,
				Required: true,
			},
			"group": {
				Description: "The group of the connector.",
				Type:        schema.TypeString,
				Optional:    true,
			},
			"provision_state": {
				Description: "The provision state of the connector.",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"implicit_group_id": {
				Description: "The ID of implicit group automaticaly created " +
					"with the connector.",
				Type:     schema.TypeInt,
				Computed: true,
			},
			"size": {
				Description: "The size of the connector, one of `SMALL`, " +
					"`MEDIUM`, `LARGE`, `2LARGE`.",
				Type:     schema.TypeString,
				Required: true,
			},
			"tunnel_protocol": {
				Description: "The tunnel protocol. It could be either `IPSEC`" +
					"or `GRE`. Default value is `IPSEC`.",
				Type:     schema.TypeString,
				Optional: true,
				Default:  "IPSEC",
			},
			"enabled": {
				Description: "Is the connector enabled. Default is `true`.",
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     true,
			},
			"scale_group_id": {
				Description: "The ID of the scale group associated " +
					"with the connector.",
				Type:     schema.TypeString,
				Optional: true,
			},
			"instances": {
				Description: "Prisma SD-WAN connector instances.",
				Type:        schema.TypeList,
				Required:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Description: "The ID of the instance.",
							Type:        schema.TypeInt,
							Computed:    true,
						},
						"host_name": {
							Description: "The host name of the Prisma SD-WAN instance.",
							Type:        schema.TypeString,
							Required:    true,
						},
						"credential_id": {
							Description: "The credential ID for the instance.",
							Type:        schema.TypeString,
							Required:    true,
						},
						"ion_model": {
							Description: "The ION model of the Prisma SD-WAN " +
								"instance. Please check Alkira Portal for " +
								"all supported models.",
							Type:     schema.TypeString,
							Required: true,
						},
						"version": {
							Description: "The version of the Prisma SD-WAN " +
								"instance. Please check Alkira Portal for " +
								"all supported versions.",
							Type:     schema.TypeString,
							Required: true,
						},
					},
				},
			},
			"target_segment": {
				Description: "Specify target segment and VRF mapping.",
				Type:        schema.TypeSet,
				Required:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"segment_id": {
							Description: "Alkira Segment ID.",
							Type:        schema.TypeInt,
							Required:    true,
						},
						"gateway_bgp_asn": {
							Description: "BGP ASN on the customer premise side.",
							Type:        schema.TypeInt,
							Required:    true,
						},
						"vrf_name": {
							Description: "The VRF name for correlating with " +
								"the Alkira segment.",
							Type:     schema.TypeString,
							Required: true,
						},
						"advertise_on_prem_routes": {
							Description: "Whether advertising On Prem Routes. " +
								"Default value is `false`.",
							Type:     schema.TypeBool,
							Optional: true,
							Default:  false,
						},
						"advertise_default_route": {
							Description: "Whether advertise default route of " +
								"internet connector. Default value is `false`.",
							Type:     schema.TypeBool,
							Optional: true,
							Default:  false,
						},
					},
				},
			},
		},
	}
}

func resourceConnectorPrismaSDWANCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {

	// INIT
	client := m.(*alkira.AlkiraClient)
	api := alkira.NewConnectorPrismaSDWAN(m.(*alkira.AlkiraClient))

	// Construct request
	request, err := generateConnectorPrismaSDWANRequest(d, m)

	if err != nil {
		return diag.FromErr(err)
	}

	// Send create request
	response, provState, err, valErr, provErr := api.Create(request)

	if err != nil {
		return diag.FromErr(err)
	}

	// Set states
	d.SetId(string(response.Id))

	if client.Validate && valErr != nil {
		var diags diag.Diagnostics
		readDiags := resourceConnectorPrismaSDWANRead(ctx, d, m)
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

	if client.Provision {
		d.Set("provision_state", provState)

		if provState == "FAILED" {
			return diag.Diagnostics{{
				Severity: diag.Warning,
				Summary:  "PROVISION (CREATE) FAILED",
				Detail:   fmt.Sprintf("%s", provErr),
			}}
		}
	}

	return resourceConnectorPrismaSDWANRead(ctx, d, m)
}

func resourceConnectorPrismaSDWANRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {

	// INIT
	client := m.(*alkira.AlkiraClient)
	api := alkira.NewConnectorPrismaSDWAN(m.(*alkira.AlkiraClient))

	// GET
	connector, provState, err := api.GetById(d.Id())

	if err != nil {
		return diag.Diagnostics{{
			Severity: diag.Warning,
			Summary:  "FAILED TO GET RESOURCE",
			Detail:   fmt.Sprintf("%s", err),
		}}
	}

	d.Set("billing_tag_ids", connector.BillingTags)
	d.Set("cxp", connector.Cxp)
	d.Set("description", connector.Description)
	d.Set("enabled", connector.Enabled)
	d.Set("group", connector.Group)
	d.Set("implicit_group_id", connector.ImplicitGroupId)
	d.Set("name", connector.Name)
	d.Set("size", connector.Size)
	d.Set("tunnel_protocol", connector.TunnelProtocol)
	d.Set("scale_group_id", connector.ScaleGroupId)

	// Set instances
	setPrismaSDWANInstances(d, connector)

	// Set VRF mapping
	var mappings []map[string]interface{}

	for _, m := range connector.PrismaSDWANVRFMappings {
		mapping := map[string]interface{}{
			"advertise_on_prem_routes": m.AdvertiseOnPremRoutes,
			"advertise_default_route":  !m.DisableInternetExit,
			"gateway_bgp_asn":          m.GatewayBgpAsn,
			"segment_id":               m.SegmentId,
			"vrf_name":                 m.VrfName,
		}
		mappings = append(mappings, mapping)
	}

	d.Set("target_segment", mappings)

	// Set provision state
	if client.Provision && provState != "" {
		d.Set("provision_state", provState)
	}

	return nil
}

func resourceConnectorPrismaSDWANUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {

	// INIT
	client := m.(*alkira.AlkiraClient)
	api := alkira.NewConnectorPrismaSDWAN(m.(*alkira.AlkiraClient))

	request, err := generateConnectorPrismaSDWANRequest(d, m)

	if err != nil {
		return diag.FromErr(err)
	}

	// UPDATE
	provState, err, valErr, provErr := api.Update(d.Id(), request)

	if err != nil {
		return diag.FromErr(err)
	}

	if client.Validate && valErr != nil {
		var diags diag.Diagnostics
		readDiags := resourceConnectorPrismaSDWANRead(ctx, d, m)
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
		if provState == "FAILED" {
			return diag.Diagnostics{{
				Severity: diag.Warning,
				Summary:  "PROVISION (UPDATE) FAILED",
				Detail:   fmt.Sprintf("%s", provErr),
			}}
		}
	}

	return resourceConnectorPrismaSDWANRead(ctx, d, m)
}

func resourceConnectorPrismaSDWANDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {

	// INIT
	client := m.(*alkira.AlkiraClient)
	api := alkira.NewConnectorPrismaSDWAN(m.(*alkira.AlkiraClient))

	// DELETE
	provState, err, valErr, provErr := api.Delete(d.Id())

	if err != nil {
		// Terraform may not print "with <resource address>" for destroys of objects
		// that are no longer in configuration, so include identifying context here.
		name, _ := d.GetOk("name")
		if nameStr, ok := name.(string); ok && nameStr != "" {
			return diag.FromErr(fmt.Errorf("%w alkira_connector_prisma_sdwan (name=%q id=%s)", err, nameStr, d.Id()))
		}
		return diag.FromErr(fmt.Errorf("%w alkira_connector_prisma_sdwan (id=%s)", err, d.Id()))
	}

	d.SetId("")

	if client.Validate && valErr != nil {
		return diag.Diagnostics{{
			Severity: diag.Error,
			Summary:  "VALIDATION (DELETE) FAILED",
			Detail:   fmt.Sprintf("%s", valErr),
		}}
	}

	// Check provision state
	if client.Provision && provState != "SUCCESS" {
		return diag.Diagnostics{{
			Severity: diag.Warning,
			Summary:  "PROVISION (DELETE) FAILED",
			Detail:   fmt.Sprintf("%s", provErr),
		}}
	}

	return nil
}

// generateConnectorPrismaSDWANRequest generate request for Prisma SD-WAN connector
func generateConnectorPrismaSDWANRequest(d *schema.ResourceData, m interface{}) (*alkira.ConnectorPrismaSDWAN, error) {

	//
	// Expand instances block
	//
	instances := expandPrismaSDWANInstances(d.Get("instances").([]interface{}))

	// Construct the request payload
	connector := &alkira.ConnectorPrismaSDWAN{
		BillingTags:            convertTypeSetToIntList(d.Get("billing_tag_ids").(*schema.Set)),
		Cxp:                    d.Get("cxp").(string),
		Description:            d.Get("description").(string),
		Enabled:                d.Get("enabled").(bool),
		Group:                  d.Get("group").(string),
		Instances:              instances,
		Name:                   d.Get("name").(string),
		PrismaSDWANVRFMappings: expandPrismaSDWANVrfMappings(d.Get("target_segment").(*schema.Set)),
		Size:                   d.Get("size").(string),
		TunnelProtocol:         d.Get("tunnel_protocol").(string),
		ScaleGroupId:           d.Get("scale_group_id").(string),
	}

	return connector, nil
}
