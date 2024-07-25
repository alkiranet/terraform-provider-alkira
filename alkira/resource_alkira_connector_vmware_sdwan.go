package alkira

import (
	"context"
	"fmt"

	"github.com/alkiranet/alkira-client-go/alkira"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func resourceAlkiraConnectorVmwareSdwan() *schema.Resource {
	return &schema.Resource{
		Description:   "Manage VMWARE SD-WAN Connector.",
		CreateContext: resourceConnectorVmwareSdwanCreate,
		ReadContext:   resourceConnectorVmwareSdwanRead,
		UpdateContext: resourceConnectorVmwareSdwanUpdate,
		DeleteContext: resourceConnectorVmwareSdwanDelete,
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
				Description: "The name of the connector.",
				Type:        schema.TypeString,
				Required:    true,
			},
			"billing_tag_ids": {
				Description: "IDs of Billing Tags to be associated " +
					"with the connector.",
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
			"orchestrator_host": {
				Description: "VMWare (Velo) Orchestrator portal host address.",
				Type:        schema.TypeString,
				Required:    true,
			},
			"implicit_group_id": {
				Description: "The ID of implicit group automaticaly created " +
					"with the connector.",
				Type:     schema.TypeInt,
				Computed: true,
			},
			"size": &schema.Schema{
				Description: "The size of the connector, one of `SMALL`, " +
					"`MEDIUM`, `LARGE`, `2LARGE`.",
				Type:     schema.TypeString,
				Required: true,
				ValidateFunc: validation.StringInSlice([]string{
					"SMALL", "MEDIUM",
					"LARGE", "2LARGE"}, false),
			},
			"tunnel_protocol": {
				Description: "Only supported tunnel protocol is `IPSEC` for now.",
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "IPSEC",
			},
			"virtual_edge": &schema.Schema{
				Description: "Virtual Edge",
				Type:        schema.TypeList,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"credential_id": {
							Description: "The generated credential ID.",
							Type:        schema.TypeString,
							Computed:    true,
						},
						"hostname": {
							Description: "The hostname of the virtual edge.",
							Type:        schema.TypeString,
							Required:    true,
						},
						"id": {
							Description: "The ID of the virtual edge.",
							Type:        schema.TypeInt,
							Computed:    true,
						},
						"activation_code": &schema.Schema{
							Description: "Activation code generated in " +
								"VMWare orchestrator account.",
							Type:     schema.TypeString,
							Required: true,
							DefaultFunc: schema.EnvDefaultFunc(
								"AK_VMWARE_SDWAN_ACTIVATION_CODE",
								nil),
						},
					},
				},
				Required: true,
			},
			"version": {
				Description: "The version of VMWARE SD-WAN.",
				Type:        schema.TypeString,
				Required:    true,
			},
			"target_segment": {
				Description: "Specify target segment.",
				Type:        schema.TypeSet,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
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
						"gateway_bgp_asn": {
							Description: "BGP ASN on the customer premise side. " +
								"A typical value for 2 byte segment " +
								"is `64523` and `4200064523` for 4 byte segment.",
							Type:     schema.TypeInt,
							Optional: true,
							Default:  65000,
						},
						"segment_id": {
							Description: "Alkira Segment ID.",
							Type:        schema.TypeInt,
							Required:    true,
						},
						"vmware_sdwan_segment_name": {
							Description: "VMWare SD-WAN Segment name for " +
								"correlating with Alkria segment.",
							Type:     schema.TypeString,
							Required: true,
						},
					},
				},
				Required: true,
			},
			"enabled": {
				Description: "Is the connector enabled. Default is `true`.",
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     true,
			},
		},
	}
}

func resourceConnectorVmwareSdwanCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {

	// INIT
	client := m.(*alkira.AlkiraClient)
	api := alkira.NewConnectorVmwareSdwan(m.(*alkira.AlkiraClient))

	// Construct request
	request, err := generateConnectorVmwareSdwanRequest(d, m)

	if err != nil {
		return diag.FromErr(err)
	}

	// Send create request
	response, provState, err, provErr := api.Create(request)

	if err != nil {
		return diag.FromErr(err)
	}

	// Set states
	d.SetId(string(response.Id))

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

	return resourceConnectorVmwareSdwanRead(ctx, d, m)
}

func resourceConnectorVmwareSdwanRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {

	// INIT
	client := m.(*alkira.AlkiraClient)
	api := alkira.NewConnectorVmwareSdwan(m.(*alkira.AlkiraClient))

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
	d.Set("group", connector.Group)
	d.Set("implicit_group_id", connector.ImplicitGroupId)
	d.Set("name", connector.Name)
	d.Set("orchestrator_host", connector.OrchestratorHostAddress)
	d.Set("size", connector.Size)
	d.Set("tunnel_protocol", connector.TunnelProtocol)
	d.Set("enabled", connector.Enabled)

	// Set virtual edge
	setVirtualEdge(d, connector)

	// Set vrf_segment_mapping
	var mappings []map[string]interface{}

	for _, m := range connector.VmWareSdWanVRFMappings {
		mapping := map[string]interface{}{
			"advertise_on_prem_routes":   m.AdvertiseOnPremRoutes,
			"advertise_default_route":    !m.DisableInternetExit,
			"gateway_bgp_asn":            m.GatewayBgpAsn,
			"segment_id":                 m.SegmentId,
			"vmware_sdwang_segment_name": m.VmWareSdWanSegmentName,
		}
		mappings = append(mappings, mapping)
	}

	d.Set("target_segment", mappings)
	d.Set("version", connector.Version)

	// Set provision state
	if client.Provision == true && provState != "" {
		d.Set("provision_state", provState)
	}

	return nil
}

func resourceConnectorVmwareSdwanUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {

	// INIT
	client := m.(*alkira.AlkiraClient)
	api := alkira.NewConnectorVmwareSdwan(m.(*alkira.AlkiraClient))

	request, err := generateConnectorVmwareSdwanRequest(d, m)

	if err != nil {
		return diag.FromErr(err)
	}

	// UPDATE
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

	return resourceConnectorVmwareSdwanRead(ctx, d, m)
}

func resourceConnectorVmwareSdwanDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {

	// INIT
	client := m.(*alkira.AlkiraClient)
	api := alkira.NewConnectorVmwareSdwan(m.(*alkira.AlkiraClient))

	// DELETE
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

// generateConnectorVmwareSdwanRequest generate request for VMWARE SD-WAN connector
func generateConnectorVmwareSdwanRequest(d *schema.ResourceData, m interface{}) (*alkira.ConnectorVmwareSdwan, error) {

	//
	// Expand virtual_edge block
	//
	virtualEdges, err := expandVmwareSdwanVirtualEdges(m.(*alkira.AlkiraClient), d.Get("virtual_edge").([]interface{}))

	if err != nil {
		return nil, err
	}

	// Construct the request payload
	connector := &alkira.ConnectorVmwareSdwan{
		BillingTags:             convertTypeSetToIntList(d.Get("billing_tag_ids").(*schema.Set)),
		Instances:               virtualEdges,
		VmWareSdWanVRFMappings:  expandVmwareSdwanVrfMappings(d.Get("target_segment").(*schema.Set)),
		Cxp:                     d.Get("cxp").(string),
		Group:                   d.Get("group").(string),
		OrchestratorHostAddress: d.Get("orchestrator_host").(string),
		Name:                    d.Get("name").(string),
		Size:                    d.Get("size").(string),
		TunnelProtocol:          d.Get("tunnel_protocol").(string),
		Version:                 d.Get("version").(string),
		Enabled:                 d.Get("enabled").(bool),
	}

	return connector, nil
}
