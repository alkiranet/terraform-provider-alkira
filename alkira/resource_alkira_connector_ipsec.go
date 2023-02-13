package alkira

import (
	"log"
	"strconv"

	"github.com/alkiranet/alkira-client-go/alkira"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func resourceAlkiraConnectorIPSec() *schema.Resource {
	return &schema.Resource{
		Description: "Manage IPSec Connector.",
		Create:      resourceConnectorIPSecCreate,
		Read:        resourceConnectorIPSecRead,
		Update:      resourceConnectorIPSecUpdate,
		Delete:      resourceConnectorIPSecDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Description: "The name of the connector.",
				Type:        schema.TypeString,
				Required:    true,
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
			"endpoint": &schema.Schema{
				Description: "The endpoint.",
				Type:        schema.TypeList,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Description: "The name of the endpoint.",
							Type:        schema.TypeString,
							Required:    true,
						},
						"customer_gateway_ip": {
							Description: "The IP address of the customer gateway.",
							Type:        schema.TypeString,
							Required:    true,
						},
						"id": {
							Description: "The ID of the endpoint.",
							Type:        schema.TypeInt,
							Computed:    true,
						},
						"preshared_keys": {
							Description: "An array of preshared keys, one per " +
								"tunnel. The value needs to be provided explictly " +
								"unlike portal.",
							Type: schema.TypeList,
							Elem: &schema.Schema{
								Type:         schema.TypeString,
								ValidateFunc: validation.StringIsNotWhiteSpace,
							},
							Required: true,
						},
						"enable_tunnel_redundancy": {
							Description: "Disable this if all tunnels will not be configured or enabled " +
								"on the on-premise device. If disabled, connector health will be shown " +
								"as `UP` if at least one of the tunnels is `UP`. If enabled, all tunnels " +
								"need to be `UP` for the connector health to be shown as `UP`.",
							Type:     schema.TypeBool,
							Optional: true,
						},
						"billing_tag_ids": {
							Description: "A list of IDs of billing tag associated with the endpoint.",
							Type:        schema.TypeList,
							Optional:    true,
							Elem: &schema.Schema{
								Type: schema.TypeInt,
							},
						},
						"ha_mode": {
							Description: "The value could be `ACTIVE` or `STANDBY`. A endpoint in `STANDBY` mode will not " +
								"be used for traffic unless all other endpoints for the connector are down. There can only " +
								"be one endpoint in `STANDBY` mode per connector and there must be at least one endpoint " +
								"that isn't in `STANDBY` mode per connector.",
							Type:         schema.TypeString,
							Optional:     true,
							ValidateFunc: validation.StringInSlice([]string{"ACTIVE", "STANDBY"}, false),
						},
						"advanced_options": {
							Description: "Advanced options for IPSec endpoint.",
							Type:        schema.TypeList,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"dpd_delay": {
										Description: "Interval to check the liveness of a peer.",
										Type:        schema.TypeInt,
										Required:    true,
									},
									"dpd_timeout": {
										Description: "Timeouts to check the liveness of a peer. IKEv1 only.",
										Type:        schema.TypeInt,
										Required:    true,
									},
									"esp_rekey_time": {
										Description: "IPsec SA rekey time in seconds.",
										Type:        schema.TypeInt,
										Required:    true,
									},
									"esp_life_time": {
										Description: "Maximum IPsec ESP lifetime if the IPsec " +
											"ESP does not rekey.",
										Type:     schema.TypeInt,
										Required: true,
									},
									"esp_random_time": {
										Description: "Time range from which to choose " +
											"a random value to subtract from rekey times in seconds.",
										Type:     schema.TypeInt,
										Required: true,
									},
									"esp_encryption_algorithms": {
										Description: "Encryption algorithms to use for IPsec SA. Value " +
											"could be `AES256CBC`, `AES192CBC`, `AES128CBC`, `AES256GCM16` " +
											"`3DESCBC`, or `NULL`.",
										Type:     schema.TypeList,
										Elem:     &schema.Schema{Type: schema.TypeString},
										Required: true,
									},
									"esp_integrity_algorithms": {
										Description: "Integrity algorithms to use for IPsec SA. Value could " +
											"`SHA1`, `SHA256`, `SHA384`, `SHA512` or `MD5`.",
										Type:     schema.TypeList,
										Elem:     &schema.Schema{Type: schema.TypeString},
										Required: true,
									},
									"esp_dh_group_numbers": {
										Description: "Diffie Hellman groups to use for IPsec SA. Value could " +
											"`MODP1024`, `MODP2048`, `MODP3072`, `MODP4096`, `MODP6144`, " +
											"`MODP8192`, `ECP256`, `ECP384`, `ECP521` and `CURVE25519`.",
										Type:     schema.TypeList,
										Elem:     &schema.Schema{Type: schema.TypeString},
										Required: true,
									},
									"initiator": {
										Description: "When true CXP will initiate the IKE connection " +
											"and if false then the customer gateway should initiate IKE. " +
											"When `gateway_ip_type` is `DYNAMIC`, initiator must be `true`.",
										Type:     schema.TypeBool,
										Required: true,
									},
									"ike_version": {
										Description:  "IKE version, either `IKEv1` or `IKEv2`",
										Type:         schema.TypeString,
										Required:     true,
										ValidateFunc: validation.StringInSlice([]string{"IKEv1", "IKEv2"}, false),
									},
									"ike_rekey_time": {
										Description: "IKE tunnel rekey time.",
										Type:        schema.TypeInt,
										Required:    true,
									},
									"ike_over_time": {
										Description: "Maximum IKE SA lifetime if the IKE SA does not rekey.",
										Type:        schema.TypeInt,
										Required:    true,
									},
									"ike_random_time": {
										Description: "Time range from which to choose a random value to " +
											"subtract from rekey times.",
										Type:     schema.TypeInt,
										Required: true,
									},
									"ike_encryption_algorithms": {
										Description: "Encryption algorithms to use for IKE SA, one of " +
											"`AES256CBC`, `AES192CBC`, `AES128CBC` and `3DESCBC`.",
										Type:     schema.TypeList,
										Elem:     &schema.Schema{Type: schema.TypeString},
										Required: true,
									},
									"ike_integrity_algorithms": {
										Description: "Integrity algorithms to use for IKE SA, one of " +
											"`SHA1`, `SHA256`, `SHA384`, `SHA512`.",
										Type:     schema.TypeList,
										Elem:     &schema.Schema{Type: schema.TypeString},
										Required: true,
									},
									"ike_dh_group_numbers": {
										Description: "Diffie Hellman groups to use for IKE SA, one of " +
											"`MODP1024`, `MODP2048`, `MODP3072`, `MODP4096`, `MODP6144`, " +
											"`MODP8192`, `ECP256`, `ECP384`, `ECP521`, `CURVE25519`.",
										Type:     schema.TypeList,
										Elem:     &schema.Schema{Type: schema.TypeString},
										Required: true,
									},
									"local_auth_type": {
										Description: "Local-ID type - IKE identity to use for " +
											"authentication round, one of `FQDN`, `USER_FQDN`, " +
											"`KEYID`, `IP_ADDR`.",
										Type:         schema.TypeString,
										Required:     true,
										ValidateFunc: validation.StringInSlice([]string{"FQDN", "USER_FQDN", "KEYID", "IP_ADDR"}, false),
									},
									"local_auth_value": {
										Description: "Local-ID value.",
										Type:        schema.TypeString,
										Required:    true,
									},
									"remote_auth_type": {
										Description: "Remote-ID type - IKE identity to use for " +
											"authentication round, one of `FQDN`, `USER_FQDN`, " +
											"`KEYID`, `IP_ADDR`.",
										Type:         schema.TypeString,
										Required:     true,
										ValidateFunc: validation.StringInSlice([]string{"FQDN", "USER_FQDN", "KEYID", "IP_ADDR"}, false),
									},
									"remote_auth_value": {
										Description: "Remote-ID value.",
										Type:        schema.TypeString,
										Required:    true,
									},
									"replay_window_size": {
										Description: "IPsec replay window for the IPsec SA.",
										Type:        schema.TypeInt,
										Required:    true,
									},
								},
							},
							Optional: true,
						},
					},
				},
				Required: true,
			},
			"group": {
				Description: "The group of the connector.",
				Type:        schema.TypeString,
				Optional:    true,
			},
			"implicit_group_id": {
				Description: "The ID of implicit group automaticaly created with the connector.",
				Type:        schema.TypeInt,
				Computed:    true,
			},
			"provision_state": {
				Description: "The provision state of the connector.",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"policy_options": {
				Description: "Policy options, both on-prem and cxp prefix" +
					"list ids must be provided if vpnMode is `POLICY_BASED`",
				Type: schema.TypeSet,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"on_prem_prefix_list_ids": {
							Description: "On Prem Prefix List IDs.",
							Type:        schema.TypeList,
							Required:    true,
							Elem:        &schema.Schema{Type: schema.TypeInt},
						},
						"cxp_prefix_list_ids": {
							Description: "CXP Prefix List IDs.",
							Type:        schema.TypeList,
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
							Description: "Routing type, one of `STATIC`, `DYNAMIC`, or `BOTH`.",
							Type:        schema.TypeString,
							Required:    true,
						},
						"availability": {
							Description: "The method to determine the availability of the routes. The value could be " +
								"`IKE_STATUS` or `IPSEC_INTERFACE_PING`. Default value is `IPSEC_INTERFACE_PING`. (**BETA**)",
							Type:         schema.TypeString,
							ValidateFunc: validation.StringInSlice([]string{"IKE_STATUS", "IPSEC_INTERFACE_PING", "PING"}, false),
							Optional:     true,
							Default:      "IPSEC_INTERFACE_PING",
						},
						"prefix_list_id": {
							Description: "The ID of prefix list to use for static route propagation.",
							Type:        schema.TypeInt,
							Optional:    true,
						},
						"customer_gateway_asn": {
							Description: "The customer gateway ASN to use for dynamic route propagation.",
							Type:        schema.TypeString,
							Optional:    true,
						},
						"bgp_auth_key": {
							Description: " BGP MD5 auth key for Alkira to authenticate Alkira CXP (On Premise Gateway).",
							Type:        schema.TypeString,
							Optional:    true,
						},
					},
				},
				Optional: true,
			},
			"segment_options": {
				Description: "Additional options for each segment associated with the connector",
				Type:        schema.TypeSet,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Description: "Segment Name.",
							Type:        schema.TypeString,
							Required:    true,
						},
						"allow_nat_exit": {
							Description: "Enable or disable access to the internet when traffic arrives via this connector. Default is `true`.",
							Type:        schema.TypeBool,
							Optional:    true,
							Default:     true,
						},

						"advertise_on_prem_routes": {
							Description: "Additional options for each segment associated with the connector. Default is `false`.",
							Type:        schema.TypeBool,
							Optional:    true,
							Default:     false,
						},
					},
				},
				Optional: true,
			},
			"segment_id": {
				Description: "The ID of the segment associated with the connector.",
				Type:        schema.TypeInt,
				Required:    true,
			},
			"size": &schema.Schema{
				Description:  "The size of the connector, one of `SMALL`, `MEDIUM`, `LARGE`, `2LARGE`, `4LARGE`, `5LARGE`, `10LARGE` and `20LARGE`.",
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringInSlice([]string{"SMALL", "MEDIUM", "LARGE", "2LARGE", "4LARGE", "5LARGE", "10LARGE", "20LARGE"}, false),
			},
			"vpn_mode": &schema.Schema{
				Description:  "The connector can be configured either in `ROUTE_BASED` or `POLICY_BASED`.",
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringInSlice([]string{"ROUTE_BASED", "POLICY_BASED"}, false),
			},
		},
	}
}

func resourceConnectorIPSecCreate(d *schema.ResourceData, m interface{}) error {
	api := alkira.NewConnectorIPSec(m.(*alkira.AlkiraClient))

	// Generate request for creating connector
	connector, err := generateConnectorIPSecRequest(d, m)

	if err != nil {
		return err
	}

	// Send request to create connector
	resource, provisionState, err := api.Create(connector)

	if err != nil {
		return err
	}

	d.Set("provision_state", provisionState)
	d.SetId(string(resource.Id))

	return resourceConnectorIPSecRead(d, m)
}

func resourceConnectorIPSecRead(d *schema.ResourceData, m interface{}) error {
	api := alkira.NewConnectorIPSec(m.(*alkira.AlkiraClient))

	connector, err := api.GetById(d.Id())

	if err != nil {
		return err
	}

	d.Set("billing_tag_ids", connector.BillingTags)
	d.Set("cxp", connector.CXP)
	d.Set("enabled", connector.Enabled)
	d.Set("group", connector.Group)
	d.Set("implicit_group_id", connector.ImplicitGroupId)
	d.Set("name", connector.Name)
	d.Set("size", connector.Size)
	d.Set("vpn_mode", connector.VpnMode)

	if len(connector.Segments) > 0 {
		segmentApi := alkira.NewSegment(m.(*alkira.AlkiraClient))
		segment, _, err := segmentApi.GetByName(connector.Segments[0])

		if err != nil {
			return err
		}
		d.Set("segment_id", segment.Id)
	}

	var endpoints []map[string]interface{}

	//
	// Go through all endpoints from the config firstly to find a
	// match, either endpoint's ID or endpoint's name should be
	// uniquely identifying an endpoint.
	//
	// On the first read call at the end of the create call, Terraform
	// didn't track any endpoint IDs yet.
	//
	for _, endpoint := range d.Get("endpoint").([]interface{}) {
		endpointConfig := endpoint.(map[string]interface{})

		for _, site := range connector.Sites {
			if endpointConfig["id"].(int) == site.Id || endpointConfig["name"].(string) == site.Name {
				endpoint := setConnectorIPSecEndpoint(site)
				endpoints = append(endpoints, endpoint)
				break
			}
		}
	}

	//
	// Go through all endpoints from the API response one more time to
	// find any endpoint that has not been tracked from Terraform
	// config.
	//
	for _, site := range connector.Sites {
		new := true

		// Check if the endpoint already exists in the Terraform config
		for _, endpoint := range d.Get("endpoint").([]interface{}) {
			endpointConfig := endpoint.(map[string]interface{})

			if endpointConfig["id"].(int) == site.Id || endpointConfig["name"].(string) == site.Name {
				new = false
				break
			}
		}

		// If the endpoint is new, add it to the tail of the list,
		// this will generate a diff
		if new {
			endpoint := setConnectorIPSecEndpoint(site)
			endpoints = append(endpoints, endpoint)
			break
		}
	}

	d.Set("endpoint", endpoints)

	return nil
}

func resourceConnectorIPSecUpdate(d *schema.ResourceData, m interface{}) error {
	api := alkira.NewConnectorIPSec(m.(*alkira.AlkiraClient))

	// Generate new request for updating connector
	connector, err := generateConnectorIPSecRequest(d, m)

	if err != nil {
		return err
	}

	// Send request to update connector
	provisionState, err := api.Update(d.Id(), connector)

	if err != nil {
		return err
	}

	d.Set("provision_state", provisionState)
	return resourceConnectorIPSecRead(d, m)
}

func resourceConnectorIPSecDelete(d *schema.ResourceData, m interface{}) error {
	api := alkira.NewConnectorIPSec(m.(*alkira.AlkiraClient))

	provisionState, err := api.Delete(d.Id())

	if err != nil {
		return err
	}

	if provisionState != "SUCCESS" {
	}

	d.SetId("")
	return nil
}

// generateConnectorIPSecRequest generate request for connector-ipsec
func generateConnectorIPSecRequest(d *schema.ResourceData, m interface{}) (*alkira.ConnectorIPSec, error) {

	sites := expandConnectorIPSecEndpoint(d.Get("endpoint").([]interface{}))

	//
	// Construct Segment
	//
	segmentApi := alkira.NewSegment(m.(*alkira.AlkiraClient))
	segment, err := segmentApi.GetById(strconv.Itoa(d.Get("segment_id").(int)))

	if err != nil {
		log.Printf("[ERROR] failed to get segment by Id: %d", d.Get("segment_id"))
		return nil, err
	}

	//
	// Construct segment options
	//
	segmentOptions, optErr := expandConnectorIPSecSegmentOptions(d.Get("segment_options").(*schema.Set))

	if optErr != nil {
		return nil, optErr
	}

	//
	// Construct Policy Options and Routing Options
	//
	// Base on the vpn_mode, switch what options to use
	//
	vpnMode := d.Get("vpn_mode").(string)

	var policyOptions *alkira.ConnectorIPSecPolicyOptions
	var routingOptions *alkira.ConnectorIPSecRoutingOptions

	switch vpnMode := d.Get("vpn_mode").(string); vpnMode {
	case "ROUTE_BASED":
		{
			routingOptions, err = expandConnectorIPSecRoutingOptions(d.Get("routing_options").(*schema.Set))

			if err != nil {
				return nil, err
			}
		}
	case "POLICY_BASED":
		{
			policyOptions, err = expandConnectorIPSecPolicyOptions(d.Get("policy_options").(*schema.Set))

			if err != nil {
				return nil, err
			}
		}
	}

	//
	// Construct the request
	//
	connector := &alkira.ConnectorIPSec{
		CXP:            d.Get("cxp").(string),
		Enabled:        d.Get("enabled").(bool),
		Group:          d.Get("group").(string),
		Name:           d.Get("name").(string),
		PolicyOptions:  policyOptions,
		RoutingOptions: routingOptions,
		SegmentOptions: segmentOptions,
		Segments:       []string{segment.Name},
		Sites:          sites,
		Size:           d.Get("size").(string),
		VpnMode:        vpnMode,
	}

	return connector, nil
}
