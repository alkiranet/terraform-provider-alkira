package alkira

import (
	"context"
	"fmt"

	"github.com/alkiranet/alkira-client-go/alkira"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func resourceAlkiraConnectorAzureExpressRoute() *schema.Resource {
	return &schema.Resource{
		Description: "Manage Azure ExpressRoute Connector. (**BETA**)",

		CreateContext: resourceConnectorAzureExpressRouteCreate,
		ReadContext:   resourceConnectorAzureExpressRouteRead,
		UpdateContext: resourceConnectorAzureExpressRouteUpdate,
		DeleteContext: resourceConnectorAzureExpressRouteDelete,
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
			"description": {
				Description: "The description of the connector.",
				Type:        schema.TypeString,
				Optional:    true,
			},
			"size": {
				Description: "The size of the connector, one of `SMALL`, `MEDIUM`, `LARGE`," +
					" `2LARGE`, `5LARGE`, `10LARGE`.",
				Type:     schema.TypeString,
				Required: true,
			},
			"enabled": {
				Description: "Is the connector enabled. Default is `true`.",
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     true,
			},
			"vhub_prefix": {
				Description: "IP address prefix for VWAN Hub. This should be a `/23` prefix.",
				Type:        schema.TypeString,
				Required:    true,
			},
			"tunnel_protocol": {
				Description:  "The tunnel protocol. One of `VXLAN`, `VXLAN_GPE`. Default is `VXLAN_GPE`",
				Type:         schema.TypeString,
				Optional:     true,
				Default:      "VXLAN_GPE",
				ValidateFunc: validation.StringInSlice([]string{"VXLAN", "VXLAN_GPE"}, false),
			},
			"cxp": {
				Description: "The CXP where the connector should be provisioned.",
				Type:        schema.TypeString,
				Required:    true,
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
			"billing_tag_ids": {
				Description: "Billing tags to be associated with " +
					"the resource. (see resource `alkira_billing_tag`).",
				Type:     schema.TypeSet,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeInt},
			},
			"instances": {
				Type:     schema.TypeList,
				Required: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Description: "User provided connector instance name.",
							Type:        schema.TypeString,
							Required:    true,
						},
						"id": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"expressroute_circuit_id": {
							Description: "ExpressRoute circuit ID from Azure. " +
								"ExpresRoute Circuit should have a private " +
								"peering connection provisioned, also an valid " +
								"authorization key associated with it.",
							Type:     schema.TypeString,
							Required: true,
						},
						"redundant_router": {
							Description: "Indicates if ExpressRoute Circuit " +
								"terminates on redundant routers on customer side.",
							Type:     schema.TypeBool,
							Optional: true,
							Default:  false,
						},
						"loopback_subnet": {
							Description: "A `/26` subnet from which loopback " +
								"IPs would be used to establish underlay " +
								"VXLAN GPE tunnels.",
							Type:     schema.TypeString,
							Required: true,
						},
						"credential_id": {
							Description: "An opaque identifier generated when " +
								"storing Azure VNET credentials.",
							Type:     schema.TypeString,
							Required: true,
						},
						"gateway_mac_address": {
							Description: "An array containing the mac addresses " +
								"of VXLAN gateways reachable through ExpressRoute " +
								"circuit. The field is only expected if VXLAN " +
								"tunnel protocol is selected, and 2 gateway MAC " +
								"addresses are expected only if `redundant_router` " +
								"is enabled.",
							Type:     schema.TypeList,
							Optional: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
						"virtual_network_interface": {
							Description: "This is an optional field if the " +
								"`tunnel_protocol` is `VXLAN`. If not specified " +
								"Alkira allocates unique VNI from the range " +
								"`[16773023, 16777215]`.",
							Type:     schema.TypeList,
							Optional: true,
							Elem:     &schema.Schema{Type: schema.TypeInt},
						},
						"segment_options": {
							Description: "Instance level segment specific routing and gateway configurations." +
								"Only required when `tunnel_protocol` is `IPSEC`.",
							Type:     schema.TypeList,
							Required: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"segment_name": {
										Description: "The name of an existing segment in the Alkira environment.",
										Type:        schema.TypeString,
										Required:    true,
									},
									"customer_gateways": {
										Description: "Customer gateway configurations for `IPSEC` tunnels. " +
											"Required only if `tunnel_protocol` is `IPSEC`.",
										Type:     schema.TypeList,
										Required: true,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"name": {
													Description: "A unique name for the customer gateway.",
													Type:        schema.TypeString,
													Required:    true,
												},
												"tunnels": {
													Description: "Tunnel configurations for the gateway. " +
														"At least one tunnel is required for `IPSEC`.",
													Type:     schema.TypeList,
													Required: true,
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"id": {
																Description: "The ID of the tunnel.",
																Type:        schema.TypeString,
																Computed:    true,
															},
															"name": {
																Description: "A unique name for the tunnel.",
																Type:        schema.TypeString,
																Required:    true,
															},
															"initiator": {
																Description: "Whether this endpoint initiates the tunnel connection.",
																Type:        schema.TypeBool,
																Optional:    true,
															},
															"profile_id": {
																Description: "The ID of the tunnel profile to use.",
																Type:        schema.TypeInt,
																Optional:    true,
															},
															"ike_version": {
																Description: "The IKE protocol version. Currently, only `IKEv2` is supported.",
																Type:        schema.TypeString,
																Optional:    true,
															},
															"pre_shared_key": {
																Description: "The pre-shared key for tunnel authentication. " +
																	"This field is sensitive and will not be displayed in logs.",
																Type:      schema.TypeString,
																Optional:  true,
																Sensitive: true,
															},
															"remote_auth_type": {
																Description: "The authentication type for the remote endpoint. " +
																	"Only `FQDN` iscurrently supported.",
																Type:     schema.TypeString,
																Optional: true,
															},
															"remote_auth_value": {
																Description: "The authentication value for the remote endpoint. This field is sensitive.",
																Type:        schema.TypeString,
																Optional:    true,
																Sensitive:   true,
															},
														},
													},
												},
											},
										},
									},
								},
							},
						},
					},
				},
			},
			"segment_options": {
				Type:     schema.TypeList,
				Required: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"segment_name": {
							Description: "The name of an existing segment.",
							Type:        schema.TypeString,
							Required:    true,
						},
						"segment_id": {
							Description: "The ID of the segment.",
							Type:        schema.TypeInt,
							Computed:    true,
						},
						"customer_asn": {
							Description: "ASN on the customer premise side.",
							Type:        schema.TypeInt,
							Required:    true,
						},
						"disable_internet_exit": {
							Description: "Enable or disable access to the " +
								"internet when traffic arrives via this " +
								"connector.",
							Type:     schema.TypeBool,
							Optional: true,
							Default:  true,
						},
						"advertise_on_prem_routes": {
							Description: "Allow routes from the branch/premises " +
								"to be advertised to the cloud.",
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

// resourceConnectorAzureExpressRouteCreate create an Azure ExpressRoute connector
func resourceConnectorAzureExpressRouteCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {

	// INIT
	client := m.(*alkira.AlkiraClient)
	api := alkira.NewConnectorAzureExpressRoute(m.(*alkira.AlkiraClient))

	request, err := generateConnectorAzureExpressRouteRequest(d, m)

	if err != nil {
		return diag.FromErr(err)
	}

	// CREATE
	resource, provState, err, provErr := api.Create(request)

	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(string(resource.Id))

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

	return resourceConnectorAzureExpressRouteRead(ctx, d, m)
}

// resourceConnectorAzureExpressRouteRead get and save an Azure ExpressRoute connectors
func resourceConnectorAzureExpressRouteRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {

	// INIT
	client := m.(*alkira.AlkiraClient)
	api := alkira.NewConnectorAzureExpressRoute(m.(*alkira.AlkiraClient))

	// GET
	connector, provState, err := api.GetById(d.Id())

	if err != nil {
		return diag.Diagnostics{{
			Severity: diag.Warning,
			Summary:  "FAILED TO GET RESOURCE",
			Detail:   fmt.Sprintf("%s", err),
		}}
	}

	d.Set("size", connector.Size)
	d.Set("billing_tag_ids", connector.BillingTags)
	d.Set("cxp", connector.Cxp)
	d.Set("group", connector.Group)
	d.Set("enabled", connector.Enabled)
	d.Set("name", connector.Name)
	d.Set("description", connector.Description)
	d.Set("tunnel_protocol", connector.TunnelProtocol)
	d.Set("vhub_prefix", connector.VhubPrefix)

	instances := make([]map[string]interface{}, len(connector.Instances))
	for i, instance := range connector.Instances {
		instances[i] = flattenInstance(instance)
	}

	d.Set("instances", instances)

	segments := make([]map[string]interface{}, len(connector.SegmentOptions))
	for i, seg := range connector.SegmentOptions {
		segments[i] = map[string]interface{}{
			"segment_name":             seg.SegmentName,
			"segment_id":               seg.SegmentId,
			"customer_asn":             seg.CustomerAsn,
			"disable_internet_exit":    seg.DisableInternetExit,
			"advertise_on_prem_routes": seg.AdvertiseOnPremRoutes,
		}
	}
	d.Set("segment_options", segments)

	// Set provision state
	if client.Provision == true && provState != "" {
		d.Set("provision_state", provState)
	}

	return nil
}

// resourceConnectorAzureExpressRouteUpdate update an Azure ExpressRoute connector
func resourceConnectorAzureExpressRouteUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {

	// INIT
	client := m.(*alkira.AlkiraClient)
	api := alkira.NewConnectorAzureExpressRoute(m.(*alkira.AlkiraClient))

	connector, err := generateConnectorAzureExpressRouteRequest(d, m)

	if err != nil {
		return diag.FromErr(err)
	}

	// UPDATE
	provState, err, provErr := api.Update(d.Id(), connector)

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

	return resourceConnectorAzureExpressRouteRead(ctx, d, m)
}

func resourceConnectorAzureExpressRouteDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {

	// INIT
	client := m.(*alkira.AlkiraClient)
	api := alkira.NewConnectorAzureExpressRoute(m.(*alkira.AlkiraClient))

	// DELETE
	provState, err, provErr := api.Delete((d.Id()))

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

// generateConnectorAzureExpressRouteRequest generate a request for Azure ExpressRoute connector
func generateConnectorAzureExpressRouteRequest(d *schema.ResourceData, m interface{}) (*alkira.ConnectorAzureExpressRoute, error) {

	billingTags := convertTypeSetToIntList(d.Get("billing_tag_ids").(*schema.Set))

	instances, err := expandAzureExpressRouteInstances(d.Get("instances").([]interface{}), m)
	if err != nil {
		return nil, err
	}

	segmentOptions, err := expandAzureExpressRouteSegments(d.Get("segment_options").([]interface{}), m)
	if err != nil {
		return nil, err
	}

	request := &alkira.ConnectorAzureExpressRoute{
		Name:           d.Get("name").(string),
		Description:    d.Get("description").(string),
		Size:           d.Get("size").(string),
		BillingTags:    billingTags,
		Enabled:        d.Get("enabled").(bool),
		TunnelProtocol: d.Get("tunnel_protocol").(string),
		Cxp:            d.Get("cxp").(string),
		Group:          d.Get("group").(string),
		VhubPrefix:     d.Get("vhub_prefix").(string),
		Instances:      instances,
		SegmentOptions: segmentOptions,
	}

	return request, nil
}
