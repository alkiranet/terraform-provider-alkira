package alkira

import (
	"log"
	"reflect"
	"strconv"

	"github.com/alkiranet/alkira-client-go/alkira"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func resourceAlkiraConnectorIPSec() *schema.Resource {
	return &schema.Resource{
		Description: "Manage IPSec Connector.\n\n\n\n" +
			"## VPN Mode\n\n" +
			"`vpn_mode` could be either `ROUTE_BASED` or `POLICY_BASED`. When it's " +
			"defined as `ROUTE_BASED`, `routing_options` block is required. When " +
			"it's defined as `POLICY_BASED`, `policy_options` block is required.",
		Create: resourceConnectorIPSecCreate,
		Read:   resourceConnectorIPSecRead,
		Update: resourceConnectorIPSecUpdate,
		Delete: resourceConnectorIPSecDelete,
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
							Optional: true,
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
								Type:         schema.TypeInt,
							},
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
											"`AES256CBC`, `AES192CBC`, `AES128CBC`.",
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
				Optional: true,
			},
			"group": {
				Description: "The group of the connector.",
				Type:        schema.TypeString,
				Optional:    true,
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
							Description: "The method to determine the availability of static route. The value could be " +
								"`IKE_STATUS` or `IPSEC_INTERFACE_PING`. Default value is `IPSEC_INTERFACE_PING`.",
							Type:         schema.TypeString,
							ValidateFunc: validation.StringInSlice([]string{"IKE_STATUS", "IPSEC_INTERFACE_PING"}, false),
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
	client := m.(*alkira.AlkiraClient)
	connector, err := generateConnectorIPSecRequest(d, m)

	if err != nil {
		return err
	}

	id, err := client.CreateConnectorIPSec(connector)

	if err != nil {
		return err
	}

	d.SetId(id)

	return resourceConnectorIPSecRead(d, m)
}

func resourceConnectorIPSecRead(d *schema.ResourceData, m interface{}) error {
	client := m.(*alkira.AlkiraClient)

	connector, err := client.GetConnectorIPSec(d.Id())

	if err != nil {
		return err
	}

	d.Set("billing_tag_ids", connector.BillingTags)
	d.Set("cxp", connector.CXP)
	d.Set("enabled", connector.Enabled)
	d.Set("group", connector.Group)
	d.Set("name", connector.Name)
	d.Set("size", connector.Size)
	d.Set("vpn_mode", connector.VpnMode)

	if len(connector.Segments) > 0 {
		segment, err := client.GetSegmentByName(connector.Segments[0])

		if err != nil {
			return err
		}
		d.Set("segment_id", segment.Id)
	}

	// Set endpoint
	var endpoints []map[string]interface{}

	for _, site := range connector.Sites {

		var advanced []map[string]interface{}

		siteValue := reflect.ValueOf(site).Elem()
		siteAdvanced := siteValue.FieldByName("Advanced")

		if siteAdvanced == (reflect.Value{}) {
			advancedConfig := map[string]interface{}{
				"dpd_delay":                 site.Advanced.DPDDelay,
				"dpd_timeout":               site.Advanced.DPDTimeout,
				"esp_dh_group_numbers":      site.Advanced.EspDHGroupNumbers,
				"esp_encryption_algorithms": site.Advanced.EspEncryptionAlgorithms,
				"esp_integrity_algorithms":  site.Advanced.EspIntegrityAlgorithms,
				"esp_life_time":             site.Advanced.EspLifeTime,
				"esp_random_time":           site.Advanced.EspRandomTime,
				"esp_rekey_time":            site.Advanced.EspRekeyTime,
				"ike_encryption_algorithms": site.Advanced.IkeEncryptionAlgorithms,
				"ike_integrity_algorithms":  site.Advanced.IkeIntegrityAlgorithms,
				"ike_over_time":             site.Advanced.IkeOverTime,
				"ike_random_time":           site.Advanced.IkeRandomTime,
				"ike_rekey_time":            site.Advanced.IkeRekeyTime,
				"ike_version":               site.Advanced.IkeVersion,
				"local_auth_type":           site.Advanced.LocalAuthType,
				"local_auth_value":          site.Advanced.LocalAuthValue,
				"remote_auth_type":          site.Advanced.RemoteAuthType,
				"remote_auth_value":         site.Advanced.RemoteAuthValue,
				"replay_window_size":        site.Advanced.ReplayWindowSize,
			}
			advanced = append(advanced, advancedConfig)
		}

		endpoint := map[string]interface{}{
			"name":                     site.Name,
			"billing_tag_ids":          site.BillingTags,
			"customer_gateway_ip":      site.CustomerGwIp,
			"enable_tunnel_redundancy": site.EnableTunnelRedundancy,
			"preshared_keys":           site.PresharedKeys,
			"id":                       site.Id,
			"advanced_options":         advanced,
		}
		endpoints = append(endpoints, endpoint)
	}

	d.Set("endpoint", endpoints)

	return nil
}

func resourceConnectorIPSecUpdate(d *schema.ResourceData, m interface{}) error {
	client := m.(*alkira.AlkiraClient)
	connector, err := generateConnectorIPSecRequest(d, m)

	if err != nil {
		return err
	}

	err = client.UpdateConnectorIPSec(d.Id(), connector)

	if err != nil {
		return err
	}

	return resourceConnectorIPSecRead(d, m)
}

func resourceConnectorIPSecDelete(d *schema.ResourceData, m interface{}) error {
	client := m.(*alkira.AlkiraClient)

	return client.DeleteConnectorIPSec(d.Id())
}

// generateConnectorIPSecRequest generate request for connector-ipsec
func generateConnectorIPSecRequest(d *schema.ResourceData, m interface{}) (*alkira.ConnectorIPSec, error) {
	client := m.(*alkira.AlkiraClient)

	sites := expandConnectorIPSecEndpoint(d.Get("endpoint").([]interface{}))

	// For now, IPSec connector only support single segment
	segment, err := client.GetSegmentById(strconv.Itoa(d.Get("segment_id").(int)))

	if err != nil {
		log.Printf("[ERROR] failed to get segment by Id: %d", d.Get("segment_id"))
		return nil, err
	}

	segmentOptions, optErr := expandConnectorIPSecSegmentOptions(d.Get("segment_options").(*schema.Set))

	if optErr != nil {
		return nil, optErr
	}

	// Base on the vpn_mode, switch what options to use
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
