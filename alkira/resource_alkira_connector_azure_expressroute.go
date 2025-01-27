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
		Description: "Manages an Azure ExpressRoute Connector. This resource is in **BETA** and may require careful handling.",

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
				Description: "The unique name of the Azure ExpressRoute connector.",
				Type:        schema.TypeString,
				Required:    true,
			},
			"description": {
				Description: "The description of the connector.",
				Type:        schema.TypeString,
				Optional:    true,
			},
			"size": {
				Description: "The size of the connector. Valid values are " +
					"`SMALL`, `MEDIUM`, `LARGE`, `2LARGE`, `5LARGE`, or `10LARGE`.",
				Type:     schema.TypeString,
				Required: true,
			},
			"enabled": {
				Description: "Whether the connector is enabled. Defaults to `true`.",
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     true,
			},
			"vhub_prefix": {
				Description: "The CIDR block (in `/23` notation) reserved for the Azure Virtual WAN Hub.",
				Type:        schema.TypeString,
				Required:    true,
			},
			"tunnel_protocol": {
				Description: "The encapsulation protocol for the tunnels. Valid " +
					"values are `VXLAN` or `VXLAN_GPE`. Defaults to `VXLAN_GPE`.",
				Type:         schema.TypeString,
				Optional:     true,
				Default:      "VXLAN_GPE",
				ValidateFunc: validation.StringInSlice([]string{"VXLAN", "VXLAN_GPE"}, false),
			},
			"cxp": {
				Description: "The Cloud Exchange Point (CXP) where the connector will be provisioned.",
				Type:        schema.TypeString,
				Required:    true,
			},
			"group": {
				Description: "The organizational group to which this connector belongs within the Alkira platform.",
				Type:        schema.TypeString,
				Optional:    true,
			},
			"provision_state": {
				Description: "The current provisioning state of the connector, as reported by the Alkira platform.",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"billing_tag_ids": {
				Description: "A list of billing tag IDs to associate with this connector. Billing" +
					" tags must be created using the `alkira_billing_tag` resource.",
				Type:     schema.TypeSet,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeInt},
			},
			"instances": {
				Description: "Configuration for individual Azure ExpressRoute instances linked to this connector.",
				Type:        schema.TypeList,
				Required:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Description: "A user-defined name for the connector instance.",
							Type:        schema.TypeString,
							Required:    true,
						},
						"id": {
							Description: "The auto-generated ID of the connector instance.",
							Type:        schema.TypeInt,
							Computed:    true,
						},
						"expressroute_circuit_id": {
							Description: "The Azure-assigned ID of the ExpressRoute Circuit." +
								" The circuit must have private peering configured and a valid authorization key.",
							Type:     schema.TypeString,
							Required: true,
						},
						"redundant_router": {
							Description: "Whether the ExpressRoute Circuit connects to redundant" +
								" routers on the customer side. Defaults to `false`.",
							Type:     schema.TypeBool,
							Optional: true,
							Default:  false,
						},
						"loopback_subnet": {
							Description: "A `/26` subnet used to allocate loopback IPs for " +
								"establishing `VXLAN/GPE` underlay tunnels.",
							Type:     schema.TypeString,
							Required: true,
						},
						"credential_id": {
							Description: "The credential ID obtained after storing Azure VNET credentials in the Alkira platform.",
							Type:        schema.TypeString,
							Required:    true,
						},
						"gateway_mac_address": {
							Description: "MAC addresses of VXLAN gateways accessible via the ExpressRoute Circuit." +
								" Required if `tunnel_protocol` is `VXLAN`. Provide two addresses if `redundant_router` is `true`.",
							Type:     schema.TypeList,
							Optional: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
						"virtual_network_interface": {
							Description: "Virtual Network Interface (VNI) identifiers. " +
								"If unspecified, Alkira auto-allocates from the range `16773023-16777215`." +
								" Onlyapplicable for `VXLAN` tunnels.",
							Type:     schema.TypeList,
							Optional: true,
							Elem:     &schema.Schema{Type: schema.TypeInt},
						},
						"segment_options": {
							Description: "Instance level segment specific routing and gateway configurations.",
							Type:        schema.TypeList,
							Required:    true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"segment_name": {
										Description: "The name of an existing segment in the Alkira environment.",
										Type:        schema.TypeString,
										Required:    true,
									},
									"customer_gateways": {
										Description: "Customer gateway configurations for `IPSEC` tunnels. " +
											"Requiredonly if `tunnel_protocol` is `IPSEC`.",
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
				Description: "Global segment routing and policy configurations for the connector.",
				Type:        schema.TypeList,
				Required:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"segment_name": {
							Description: "The name of the segment to associate with this connector.",
							Type:        schema.TypeString,
							Required:    true,
						},
						"segment_id": {
							Description: "The auto-generated ID of the segment within the Alkira platform.",
							Type:        schema.TypeInt,
							Computed:    true,
						},
						"customer_asn": {
							Description: "The Autonomous System Number (ASN) configured on the customer's premises.",
							Type:        schema.TypeInt,
							Required:    true,
						},
						"disable_internet_exit": {
							Description: "When enabled, traffic arriving via this connector cannot access the internet. Defaults to `true`.",
							Type:        schema.TypeBool,
							Optional:    true,
							Default:     true,
						},
						"advertise_on_prem_routes": {
							Description: "When enabled, routes from on-premises networks are advertised to Azure. Defaults to `false`.",
							Type:        schema.TypeBool,
							Optional:    true,
							Default:     false,
						},
					},
				},
			},
		},
	}
}

// resourceConnectorAzureExpressRouteCreate creates an Azure ExpressRoute connector
func resourceConnectorAzureExpressRouteCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*alkira.AlkiraClient)
	api := alkira.NewConnectorAzureExpressRoute(client)

	request, err := generateConnectorAzureExpressRouteRequest(d, m)
	if err != nil {
		return diag.FromErr(err)
	}

	resource, provState, err, provErr := api.Create(request)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(string(resource.Id))

	if client.Provision {
		d.Set("provision_state", provState)
		if provErr != nil {
			return diag.Diagnostics{{
				Severity: diag.Warning,
				Summary:  "Provisioning (Create) Failed",
				Detail:   fmt.Sprintf("%s", provErr),
			}}
		}
	}

	return resourceConnectorAzureExpressRouteRead(ctx, d, m)
}

// resourceConnectorAzureExpressRouteRead retrieves and updates the state of an Azure ExpressRoute connector
func resourceConnectorAzureExpressRouteRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*alkira.AlkiraClient)
	api := alkira.NewConnectorAzureExpressRoute(client)

	connector, provState, err := api.GetById(d.Id())
	if err != nil {
		return diag.Diagnostics{{
			Severity: diag.Warning,
			Summary:  "Failed to Retrieve Resource",
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

	if client.Provision && provState != "" {
		d.Set("provision_state", provState)
	}

	return nil
}

// resourceConnectorAzureExpressRouteUpdate updates an existing Azure ExpressRoute connector
func resourceConnectorAzureExpressRouteUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*alkira.AlkiraClient)
	api := alkira.NewConnectorAzureExpressRoute(client)

	connector, err := generateConnectorAzureExpressRouteRequest(d, m)
	if err != nil {
		return diag.FromErr(err)
	}

	provState, err, provErr := api.Update(d.Id(), connector)
	if err != nil {
		return diag.FromErr(err)
	}

	if client.Provision {
		d.Set("provision_state", provState)
		if provErr != nil {
			return diag.Diagnostics{{
				Severity: diag.Warning,
				Summary:  "Provisioning (Update) Failed",
				Detail:   fmt.Sprintf("%s", provErr),
			}}
		}
	}

	return resourceConnectorAzureExpressRouteRead(ctx, d, m)
}

// resourceConnectorAzureExpressRouteDelete removes an Azure ExpressRoute connector
func resourceConnectorAzureExpressRouteDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*alkira.AlkiraClient)
	api := alkira.NewConnectorAzureExpressRoute(client)

	provState, err, provErr := api.Delete(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId("")

	if client.Provision && provState != "SUCCESS" {
		return diag.Diagnostics{{
			Severity: diag.Warning,
			Summary:  "Provisioning (Delete) Failed",
			Detail:   fmt.Sprintf("%s", provErr),
		}}
	}

	return nil
}

// generateConnectorAzureExpressRouteRequest constructs a request for Azure ExpressRoute connector operations
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
