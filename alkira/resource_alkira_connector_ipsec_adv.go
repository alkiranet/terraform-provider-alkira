package alkira

import (
	"context"
	"fmt"

	"github.com/alkiranet/alkira-client-go/alkira"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func resourceAlkiraConnectorIPSecAdv() *schema.Resource {
	return &schema.Resource{
		Description:   "Manage Advanced IPSec Connector. (**BETA**)",
		CreateContext: resourceConnectorIPSecAdvCreate,
		ReadContext:   resourceConnectorIPSecAdvRead,
		UpdateContext: resourceConnectorIPSecAdvUpdate,
		DeleteContext: resourceConnectorIPSecAdvDelete,
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
			"advertise_default_route": {
				Description: "Enable or disable access to the internet " +
					"when traffic arrives via this connector. Default " +
					"is `false`.",
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"advertise_on_prem_routes": {
				Description: "Additional options for each segment " +
					"associated with the connector. Default is `false`.",
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"billing_tag_ids": {
				Description: "A list of IDs of billing tag associated " +
					"with the gateway.",
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeInt,
				},
			},
			"cxp": {
				Description: "The CXP where the connector should be provisioned.",
				Type:        schema.TypeString,
				Required:    true,
			},
			"destination_type": {
				Description: "The destination type of the connector. The value " +
					"could be `IPSEC_ENDPOINT`, `AWS_VPN_CONNECTION`, " +
					"`AZURE_VPN_CONNECTION`. The default value is `IPSEC_ENDPOINT`.",
				Type:     schema.TypeString,
				Optional: true,
				Default:  "IPSEC_ENDPOINT",
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
				Description: "The ID of implicit group automaticaly created " +
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
			"segment_id": {
				Description: "The ID of the segment associated " +
					"with the connector.",
				Type:     schema.TypeString,
				Required: true,
			},
			"size": &schema.Schema{
				Description: "The size of the connector, one of " +
					"`SMALL`, `MEDIUM`, `LARGE`, `2LARGE`, `4LARGE`, " +
					"`5LARGE`, `10LARGE` and `20LARGE`.",
				Type:     schema.TypeString,
				Required: true,
				ValidateFunc: validation.StringInSlice([]string{
					"SMALL", "MEDIUM", "LARGE", "2LARGE", "4LARGE",
					"5LARGE", "10LARGE", "20LARGE"}, false),
			},
			"tunnels_per_gateway": {
				Description: "The number of tunnels per gateway instance. " +
					"Default is `1`.",
				Type:     schema.TypeInt,
				Optional: true,
				Default:  1,
			},
			"vpn_mode": &schema.Schema{
				Description: "The connector can be configured " +
					"either in `ROUTE_BASED` or `POLICY_BASED`.",
				Type:     schema.TypeString,
				Required: true,
				ValidateFunc: validation.StringInSlice([]string{
					"ROUTE_BASED", "POLICY_BASED"}, false),
			},
			"gateway": &schema.Schema{
				Description: "The gateway.",
				Type:        schema.TypeList,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Description: "The name of the endpoint.",
							Type:        schema.TypeString,
							Required:    true,
						},
						"customer_gateway_ip": {
							Description: "The IP address of the customer gateway. " +
								"This should be `0.0.0.0` to signify that this " +
								"is a dynamic gateway.",
							Type:     schema.TypeString,
							Required: true,
						},
						"ha_mode": {
							Description: "The value could be `ACTIVE` or" +
								"`STANDBY`. A gateway in `STANDBY` mode " +
								"will not be used for traffic unless all " +
								"other gateway for the connector are down. " +
								"There can only be one gateway in `STANDBY` " +
								"mode per connector and there must be at " +
								"least one gateway that isn't in `STANDBY` " +
								"mode per connector.",
							Type:     schema.TypeString,
							Optional: true,
							ValidateFunc: validation.StringInSlice([]string{
								"ACTIVE", "STANDBY"}, false),
						},
						"id": {
							Description: "The ID of the gateway.",
							Type:        schema.TypeInt,
							Computed:    true,
						},
						"tunnel": {
							Description: "Tunnel of the gateway. The number " +
								"of the tunnels should be equal to " +
								"`tunnel_per_gateway`.",
							Type: schema.TypeList,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"id": {
										Description: "The ID of the tunnel.",
										Type:        schema.TypeString,
										Computed:    true,
									},
									"number": {
										Description: "The number of the tunnel.",
										Type:        schema.TypeInt,
										Required:    true,
									},
									"preshared_key": {
										Description:  "The preshared key of the tunnel ",
										Type:         schema.TypeString,
										ValidateFunc: validation.StringIsNotWhiteSpace,
										Required:     true,
									},
									"profile_id": {
										Description: "The profile ID of the tunnel.",
										Type:        schema.TypeInt,
										Optional:    true,
									},
									"customer_end_overlay_ip_reservation_id": {
										Description: "The overlay IP reservation " +
											"ID of the customer end of the tunnel.",
										Type:     schema.TypeString,
										Required: true,
									},
									"cxp_end_overlay_ip_reservation_id": {
										Description: "The overlay IP reservation " +
											"ID of the CXP end of the tunnel.",
										Type:     schema.TypeString,
										Required: true,
									},
									"cxp_end_public_ip_reservation_id": {
										Description: "The public IP reservation " +
											"ID of the CXP end of the tunnel.",
										Type:     schema.TypeString,
										Required: true,
									},
									"advanced_options": {
										Description: "Advanced options for the " +
											"IPSec gateway.",
										Type:     schema.TypeList,
										MaxItems: 1,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"dpd_delay": {
													Description: "Interval to check the " +
														"liveness of a peer.",
													Type:     schema.TypeInt,
													Required: true,
												},
												"dpd_timeout": {
													Description: "Timeouts to check the " +
														"liveness of a peer. IKEv1 only.",
													Type:     schema.TypeInt,
													Required: true,
												},
												"esp_rekey_time": {
													Description: "IPsec SA rekey time in " +
														"seconds.",
													Type:     schema.TypeInt,
													Required: true,
												},
												"esp_life_time": {
													Description: "Maximum IPsec ESP lifetime" +
														"if the IPsec ESP does not rekey.",
													Type:     schema.TypeInt,
													Required: true,
												},
												"esp_random_time": {
													Description: "Time range from which to " +
														"choose a random value to subtract " +
														"from rekey times in seconds.",
													Type:     schema.TypeInt,
													Required: true,
												},
												"esp_encryption_algorithms": {
													Description: "Encryption algorithms to use for " +
														"IPsec SA. Value could be `AES256CBC`, " +
														"`AES192CBC`, `AES128CBC`, `AES256GCM16` " +
														"`3DESCBC`, or `NULL`.",
													Type:     schema.TypeList,
													Elem:     &schema.Schema{Type: schema.TypeString},
													Required: true,
												},
												"esp_integrity_algorithms": {
													Description: "Integrity algorithms to use for " +
														"IPsec SA. Value could `SHA1`, `SHA256`, " +
														"`SHA384`, `SHA512` or `MD5`.",
													Type:     schema.TypeList,
													Elem:     &schema.Schema{Type: schema.TypeString},
													Required: true,
												},
												"esp_dh_group_numbers": {
													Description: "Diffie Hellman groups to use " +
														"for IPsec SA. Value could `MODP1024`, " +
														"`MODP2048`, `MODP3072`, `MODP4096`, " +
														"`MODP6144`, `MODP8192`, `ECP256`, " +
														"`ECP384`, `ECP521` and `CURVE25519`.",
													Type:     schema.TypeList,
													Elem:     &schema.Schema{Type: schema.TypeString},
													Required: true,
												},
												"initiator": {
													Description: "When true CXP will initiate " +
														"the IKE connection and if false then " +
														"the customer gateway should initiate " +
														"IKE. When `gateway_ip_type` is `DYNAMIC`," +
														"initiator must be `true`.",
													Type:     schema.TypeBool,
													Required: true,
												},
												"ike_version": {
													Description: "IKE version, either `IKEv1` " +
														"or `IKEv2`",
													Type:     schema.TypeString,
													Required: true,
													ValidateFunc: validation.StringInSlice([]string{
														"IKEv1", "IKEv2"}, false),
												},
												"ike_rekey_time": {
													Description: "IKE tunnel rekey time.",
													Type:        schema.TypeInt,
													Required:    true,
												},
												"ike_over_time": {
													Description: "Maximum IKE SA lifetime if " +
														"the IKE SA does not rekey.",
													Type:     schema.TypeInt,
													Required: true,
												},
												"ike_random_time": {
													Description: "Time range from which to " +
														"choose a random value to " +
														"subtract from rekey times.",
													Type:     schema.TypeInt,
													Required: true,
												},
												"ike_encryption_algorithms": {
													Description: "Encryption algorithms to use " +
														"for IKE SA, one of `AES256CBC`, `AES192CBC`," +
														"`AES128CBC` and `3DESCBC`.",
													Type:     schema.TypeList,
													Elem:     &schema.Schema{Type: schema.TypeString},
													Required: true,
												},
												"ike_integrity_algorithms": {
													Description: "Integrity algorithms to use " +
														"for IKE SA, one of `SHA1`, `SHA256`, " +
														"`SHA384`, `SHA512`.",
													Type:     schema.TypeList,
													Elem:     &schema.Schema{Type: schema.TypeString},
													Required: true,
												},
												"ike_dh_group_numbers": {
													Description: "Diffie Hellman groups to use " +
														"for IKE SA, one of `MODP1024`, " +
														"`MODP2048`, `MODP3072`, `MODP4096`, " +
														"`MODP6144`, `MODP8192`, `ECP256`, " +
														"`ECP384`, `ECP521`, `CURVE25519`.",
													Type:     schema.TypeList,
													Elem:     &schema.Schema{Type: schema.TypeString},
													Required: true,
												},
												"local_auth_type": {
													Description: "Local-ID type - IKE identity " +
														"to use for authentication round, one " +
														"of `FQDN`, `USER_FQDN`, `KEYID`, " +
														"`IP_ADDR`.",
													Type:     schema.TypeString,
													Required: true,
													ValidateFunc: validation.StringInSlice([]string{
														"FQDN", "USER_FQDN", "KEYID", "IP_ADDR"}, false),
												},
												"local_auth_value": {
													Description: "Local-ID value.",
													Type:        schema.TypeString,
													Required:    true,
												},
												"remote_auth_type": {
													Description: "Remote-ID type - IKE " +
														"identity to use for authentication " +
														"round, one of `FQDN`, `USER_FQDN`, " +
														"`KEYID`, `IP_ADDR`.",
													Type:     schema.TypeString,
													Required: true,
													ValidateFunc: validation.StringInSlice([]string{
														"FQDN", "USER_FQDN", "KEYID", "IP_ADDR"}, false),
												},
												"remote_auth_value": {
													Description: "Remote-ID value.",
													Type:        schema.TypeString,
													Required:    true,
												},
												"replay_window_size": {
													Description: "IPsec replay window for the " +
														"IPsec SA.",
													Type:     schema.TypeInt,
													Required: true,
												},
											},
										},
										Optional: true,
									}, // advanced_options
								},
							},
							Required: true,
						}, // tunnel
					},
				},
				Required: true,
			}, // gateway
			"policy_options": {
				Description: "Policy options, both `on_prem_prefix_list_ids` " +
					"and `cxp_prefix_list_ids` must be provided if `vpn_mode` " +
					"is `POLICY_BASED`",
				Type: schema.TypeSet,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"on_prem_prefix_list_ids": {
							Description: "On-Prem Prefix List IDs.",
							Type:        schema.TypeSet,
							Required:    true,
							Elem:        &schema.Schema{Type: schema.TypeInt},
						},
						"cxp_prefix_list_ids": {
							Description: "CXP Prefix List IDs.",
							Type:        schema.TypeSet,
							Required:    true,
							Elem:        &schema.Schema{Type: schema.TypeInt},
						},
					},
				},
				Optional: true,
			},
			"routing_options": {
				Description: "Routing options, type is `STATIC`, `DYNAMIC`, or" +
					"`BOTH` must be provided if `vpn_mode` is `ROUTE_BASED`",
				Type: schema.TypeSet,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"type": {
							Description: "Routing type, one of `STATIC`, " +
								"`DYNAMIC`, or `BOTH`.",
							Type:     schema.TypeString,
							Required: true,
						},
						"availability": {
							Description: "The method to determine the availability " +
								"of the routes. The value could be `IKE_STATUS` " +
								"or `IPSEC_INTERFACE_PING`. Default value is " +
								"`IPSEC_INTERFACE_PING`.",
							Type: schema.TypeString,
							ValidateFunc: validation.StringInSlice([]string{
								"IKE_STATUS", "IPSEC_INTERFACE_PING", "PING"}, false),
							Optional: true,
							Default:  "IPSEC_INTERFACE_PING",
						},
						"prefix_list_id": {
							Description: "The ID of prefix list to use for static " +
								"route propagation.",
							Type:     schema.TypeInt,
							Optional: true,
						},
						"customer_gateway_asn": {
							Description: "The customer gateway ASN to use for " +
								"dynamic route propagation.",
							Type:     schema.TypeString,
							Optional: true,
						},
						"bgp_auth_key": {
							Description: " BGP MD5 auth key for Alkira to " +
								"authenticate Alkira CXP (On Premise Gateway).",
							Type:     schema.TypeString,
							Optional: true,
						},
					},
				},
				Optional: true,
			},
		},
	}
}

func resourceConnectorIPSecAdvCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {

	// INIT
	client := m.(*alkira.AlkiraClient)
	api := alkira.NewConnectorAdvIPSec(m.(*alkira.AlkiraClient))

	request, err := generateConnectorIPSecAdvRequest(d, m)

	if err != nil {
		return diag.FromErr(err)
	}

	// CREATE
	response, provState, err, provErr := api.Create(request)

	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(string(response.Id))

	// Set state
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

	return resourceConnectorIPSecAdvRead(ctx, d, m)
}

func resourceConnectorIPSecAdvRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {

	// INIT
	client := m.(*alkira.AlkiraClient)
	api := alkira.NewConnectorAdvIPSec(m.(*alkira.AlkiraClient))

	// GET
	connector, provState, err := api.GetById(d.Id())

	if err != nil {
		return diag.Diagnostics{{
			Severity: diag.Warning,
			Summary:  "FAILED TO GET RESOURCE",
			Detail:   fmt.Sprintf("%s", err),
		}}
	}

	d.Set("advertise_default_route", connector.AdvertiseDefaultRoute)
	d.Set("advertise_on_prem_routes", connector.AdvertiseOnPremRoutes)
	d.Set("billing_tag_ids", connector.BillingTags)
	d.Set("cxp", connector.CXP)
	d.Set("enabled", connector.Enabled)
	d.Set("destination_type", connector.DestinationType)
	d.Set("group", connector.Group)
	d.Set("implicit_group_id", connector.ImplicitGroupId)
	d.Set("name", connector.Name)
	d.Set("size", connector.Size)
	d.Set("tunnels_per_gateway", connector.TunnelsPerGateway)
	d.Set("vpn_mode", connector.VpnMode)

	d.Set("gateway", deflateConnectorAdvIPSecGateway(connector, d))

	if connector.PolicyOptions != nil {
		d.Set("policy_options", deflateConnectorAdvIPSecPolicyOptions(connector.PolicyOptions))
	}

	// Get segment
	segmentId, err := getSegmentIdByName(connector.Segment, m)

	if err != nil {
		return diag.FromErr(err)
	}
	d.Set("segment_id", segmentId)

	// Set provision state
	if client.Provision == true && provState != "" {
		d.Set("provision_state", provState)
	}

	return nil
}

func resourceConnectorIPSecAdvUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {

	// INIT
	client := m.(*alkira.AlkiraClient)
	api := alkira.NewConnectorAdvIPSec(m.(*alkira.AlkiraClient))

	request, err := generateConnectorIPSecAdvRequest(d, m)

	if err != nil {
		return diag.FromErr(err)
	}

	// UPDATE
	provState, err, provErr := api.Update(d.Id(), request)

	if err != nil {
		return diag.FromErr(err)
	}

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

	return resourceConnectorIPSecAdvRead(ctx, d, m)
}

func resourceConnectorIPSecAdvDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {

	// INIT
	client := m.(*alkira.AlkiraClient)
	api := alkira.NewConnectorAdvIPSec(m.(*alkira.AlkiraClient))

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

// generateConnectorIPSecAdvRequest generate request for connector-ipsec
func generateConnectorIPSecAdvRequest(d *schema.ResourceData, m interface{}) (*alkira.ConnectorAdvIPSec, error) {

	gateways := expandConnectorAdvIPSecGateway(d.Get("gateway").([]interface{}))

	//
	// Segment
	//
	segmentName, err := getSegmentNameById(d.Get("segment_id").(string), m)

	if err != nil {
		return nil, err
	}

	//
	// Construct Policy Options and Routing Options
	//
	// Base on the vpn_mode, switch what options to use
	//
	vpnMode := d.Get("vpn_mode").(string)

	var policyOptions *alkira.ConnectorAdvIPSecPolicyOptions
	var routingOptions *alkira.ConnectorAdvIPSecRoutingOptions

	switch vpnMode := d.Get("vpn_mode").(string); vpnMode {
	case "ROUTE_BASED":
		{
			routingOptions, err = expandConnectorAdvIPSecRoutingOptions(
				d.Get("routing_options").(*schema.Set))

			if err != nil {
				return nil, err
			}
		}
	case "POLICY_BASED":
		{
			policyOptions, err = expandConnectorAdvIPSecPolicyOptions(
				d.Get("policy_options").(*schema.Set))

			if err != nil {
				return nil, err
			}
		}
	}

	//
	// Construct the request
	//
	connector := &alkira.ConnectorAdvIPSec{
		AdvertiseDefaultRoute: d.Get("advertise_default_route").(bool),
		AdvertiseOnPremRoutes: d.Get("advertise_on_prem_routes").(bool),
		BillingTags:           convertTypeSetToIntList(d.Get("billing_tag_ids").(*schema.Set)),
		CXP:                   d.Get("cxp").(string),
		Enabled:               d.Get("enabled").(bool),
		DestinationType:       d.Get("destination_type").(string),
		Group:                 d.Get("group").(string),
		Name:                  d.Get("name").(string),
		Segment:               segmentName,
		Size:                  d.Get("size").(string),
		TunnelsPerGateway:     d.Get("tunnels_per_gateway").(int),
		VpnMode:               vpnMode,
		Gateways:              gateways,
		PolicyOptions:         policyOptions,
		RoutingOptions:        routingOptions,
	}

	return connector, nil
}
