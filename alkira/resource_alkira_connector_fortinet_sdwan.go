package alkira

import (
	"context"
	"fmt"

	"github.com/alkiranet/alkira-client-go/alkira"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func resourceAlkiraConnectorFortinetSdwan() *schema.Resource {
	return &schema.Resource{
		Description:   "Manage Fortinet SD-WAN Connector. (**BETA**)",
		CreateContext: resourceConnectorFortinetSdwanCreate,
		ReadContext:   resourceConnectorFortinetSdwanRead,
		UpdateContext: resourceConnectorFortinetSdwanUpdate,
		DeleteContext: resourceConnectorFortinetSdwanDelete,
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
				Description: "A list of Billing Tag IDs associated " +
					"with the connector.",
				Type:     schema.TypeList,
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
			"allow_list": {
				Description: "This list allows the IP addresses or subnets to " +
					"be whitelisted so that they can communicate with the " +
					"Fortinet SD-WAN instance. The value could be `/32` IPs or " +
					"can also be a mask.",
				Type:     schema.TypeList,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"implicit_group_id": {
				Description: "The ID of implicit group automaticaly created " +
					"with the connector.",
				Type:     schema.TypeInt,
				Computed: true,
			},
			"size": &schema.Schema{
				Description: "The size of the connector, one of `SMALL`, " +
					"`MEDIUM` and `LARGE`, `2LARGE`, `5LARGE`.",
				Type:     schema.TypeString,
				Required: true,
				ValidateFunc: validation.StringInSlice([]string{
					"SMALL", "MEDIUM",
					"LARGE", "2LARGE",
					"5LARGE"}, false),
			},
			"tunnel_protocol": {
				Description: "The tunnel protocol. It could be either `IPSEC`" +
					"or `GRE`. Default value is `IPSEC`.",
				Type:     schema.TypeString,
				Optional: true,
				Default:  "IPSEC",
			},
			"wan_edge": &schema.Schema{
				Description: "WAN Edge",
				Type:        schema.TypeList,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"credential_id": {
							Description: "The generated credential ID.",
							Type:        schema.TypeString,
							Computed:    true,
						},
						"username": {
							Description: "The username of the WAN Edge instance." +
								"The default value is `admin`.",
							Type:     schema.TypeString,
							Optional: true,
							Default:  "admin",
						},
						"password": {
							Description: "The password of the WAN Edge instance.",
							Type:        schema.TypeString,
							Required:    true,
						},
						"hostname": {
							Description: "The hostname of the WAN Edge.",
							Type:        schema.TypeString,
							Required:    true,
						},
						"id": {
							Description: "The ID of the WAN Edge instance.",
							Type:        schema.TypeInt,
							Computed:    true,
						},
						"license_type": {
							Description: "The type of license. Either `PAY_AS_YOU_GO` " +
								"or `BRING_YOUR_OWN`.",
							Type:     schema.TypeString,
							Required: true,
							ValidateFunc: validation.StringInSlice([]string{
								"PAY_AS_YOU_GO", "BRING_YOUR_OWN"}, false),
						},
						"serial_number": &schema.Schema{
							Description: "Serial Number of the WAN Edge. It's " +
								"only required when `license_type` is " +
								"`BRING_YOUR_OWN`. It could be set by ENV " +
								"variable `AK_FORTINET_SDWAN_SERIAL_NUMBER`.",
							Type:     schema.TypeString,
							Optional: true,
							DefaultFunc: schema.EnvDefaultFunc(
								"AK_FORTINET_SDWAN_SERIAL_NUMBER",
								nil),
						},
						"version": {
							Description: "The version of WAN Edge.",
							Type:        schema.TypeString,
							Required:    true,
						},
					},
				},
				Required: true,
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
						"vrf_id": {
							Description: "Fortinet SD-WAN Segment name for " +
								"correlating with Alkria segment.",
							Type:     schema.TypeInt,
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

func resourceConnectorFortinetSdwanCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {

	// INIT
	client := m.(*alkira.AlkiraClient)
	api := alkira.NewConnectorFortinetSdwan(m.(*alkira.AlkiraClient))

	// Construct request
	request, err := generateConnectorFortinetSdwanRequest(d, m)

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

	return resourceConnectorFortinetSdwanRead(ctx, d, m)
}

func resourceConnectorFortinetSdwanRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {

	// INIT
	client := m.(*alkira.AlkiraClient)
	api := alkira.NewConnectorFortinetSdwan(m.(*alkira.AlkiraClient))

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
	d.Set("allow_list", connector.AllowList)
	d.Set("size", connector.Size)
	d.Set("tunnel_protocol", connector.TunnelProtocol)
	d.Set("enabled", connector.Enabled)

	// Set WAN Edge instances
	setWanEdge(d, connector)

	// Set VRF mapping
	var mappings []map[string]interface{}

	for _, m := range connector.FtntSdWanVRFMappings {
		mapping := map[string]interface{}{
			"advertise_on_prem_routes": m.AdvertiseOnPremRoutes,
			"advertise_default_route":  !m.DisableInternetExit,
			"gateway_bgp_asn":          m.GatewayBgpAsn,
			"segment_id":               m.SegmentId,
			"vrf_id":                   m.Vrf,
		}
		mappings = append(mappings, mapping)
	}

	d.Set("target_segment", mappings)

	// Set provision state
	if client.Provision == true && provState != "" {
		d.Set("provision_state", provState)
	}

	return nil
}

func resourceConnectorFortinetSdwanUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {

	// INIT
	client := m.(*alkira.AlkiraClient)
	api := alkira.NewConnectorFortinetSdwan(m.(*alkira.AlkiraClient))

	request, err := generateConnectorFortinetSdwanRequest(d, m)

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

	return resourceConnectorFortinetSdwanRead(ctx, d, m)
}

func resourceConnectorFortinetSdwanDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {

	// INIT
	client := m.(*alkira.AlkiraClient)
	api := alkira.NewConnectorFortinetSdwan(m.(*alkira.AlkiraClient))

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

// generateConnectorFortinetSdwanRequest generate request for FORTINET SD-WAN connector
func generateConnectorFortinetSdwanRequest(d *schema.ResourceData, m interface{}) (*alkira.ConnectorFortinetSdwan, error) {

	//
	// Expand wan_edge block
	//
	wanEdges, err := expandFortinetSdwanWanEdges(m.(*alkira.AlkiraClient), d.Get("wan_edge").([]interface{}))

	if err != nil {
		return nil, err
	}

	// Construct the request payload
	connector := &alkira.ConnectorFortinetSdwan{
		BillingTags:          convertTypeListToIntList(d.Get("billing_tag_ids").([]interface{})),
		Instances:            wanEdges,
		FtntSdWanVRFMappings: expandFortinetSdwanVrfMappings(d.Get("target_segment").(*schema.Set)),
		Cxp:                  d.Get("cxp").(string),
		Group:                d.Get("group").(string),
		AllowList:            convertTypeListToStringList(d.Get("allow_list").([]interface{})),
		Name:                 d.Get("name").(string),
		Size:                 d.Get("size").(string),
		TunnelProtocol:       d.Get("tunnel_protocol").(string),
		Enabled:              d.Get("enabled").(bool),
	}

	return connector, nil
}
