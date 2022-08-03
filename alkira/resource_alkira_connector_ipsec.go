package alkira

import (
	"encoding/json"
	"log"
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
			"id": {
				Description: "The ID of the IPSec connector.",
				Type:        schema.TypeInt,
				Computed:    true,
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
				Type:        schema.TypeSet,
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
						"preshared_keys": {
							Description: "An array of presharedKeys, one per tunnel.",
							Type:        schema.TypeList,
							Elem: &schema.Schema{
								Type: schema.TypeString,
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
							Elem:        &schema.Schema{Type: schema.TypeInt},
						},
						"advanced": {
							Type: schema.TypeSet,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"dpd_timeout": {
										Description: "Timeouts to check the liveness of a peer. IKEv1 only.",
										Type:        schema.TypeString,
										Required:    true,
									},
									"dpd_delay": {
										Description: "Interval to check the liveness of a peer.",
										Type:        schema.TypeString,
										Required:    true,
									},
									"esp_rekey_time": {
										Type:     schema.TypeString,
										Required: true,
									},
									"esp_life_time": {
										Description: "Maximum IPsec ESP lifetime if the IPsec ESP does not rekey.",
										Type:        schema.TypeInt,
										Required:    true,
									},
									"esp_random_time": {
										Type:     schema.TypeString,
										Required: true,
									},
									"esp_encryption_algorithms": {
										Type:     schema.TypeString,
										Required: true,
									},
									"esp_integrity_algorithms": {
										Type:     schema.TypeString,
										Required: true,
									},
									"esp_dh_group_numbers": {
										Type:     schema.TypeString,
										Required: true,
									},
									"initiator": {
										Type:     schema.TypeString,
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
										Description: "Time range from which to choose a random value to subtract from rekey times.",
										Type:        schema.TypeInt,
										Required:    true,
									},
									"ike_encryption_algorithms": {
										Description:  "Encryption algorithms to use for IKE SA, one of `AES256CBC`, `AES192CBC`, `AES128CBC`.",
										Type:         schema.TypeString,
										Required:     true,
										ValidateFunc: validation.StringInSlice([]string{"AES256CBC", "AES192CBC", "AES128CBC"}, false),
									},
									"ike_integrity_algorithms": {
										Description:  "Integrity algorithms to use for IKE SA, one of `SHA1`, `SHA256`, `SHA384`, `SHA512`.",
										Type:         schema.TypeString,
										Required:     true,
										ValidateFunc: validation.StringInSlice([]string{"SHA1", "SHA256", "SHA384", "SHA512"}, false),
									},
									"ike_dh_group_numbers": {
										Description:  "Diffie Hellman groups to use for IKE SA, one of `MODP1024`, `MODP2048`, `MODP3072`, `MODP4096`, `MODP6144`, `MODP8192`, `ECP256`, `ECP384`, `ECP521`, `CURVE25519`",
										Type:         schema.TypeString,
										Required:     true,
										ValidateFunc: validation.StringInSlice([]string{"MODP1024", "MODP2048", "MODP3072", "MODP4096", "MODP6144", "MODP8192", "ECP256", "ECP384", "ECP521", "CURVE25519"}, false),
									},
									"local_auth_type": {
										Description:  "Local-ID type - IKE identity to use for authentication round, one of `FQDN`, `USER_FQDN`, `KEYID`, `IP_ADDR`.",
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
										Description:  "Remote-ID type - IKE identity to use for authentication round, one of `FQDN`, `USER_FQDN`, `KEYID`, `IP_ADDR`.",
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

	if len(connector.Segments) > 0 {
		segment, err := client.GetSegmentByName(connector.Segments[0])

		if err != nil {
			return err
		}
		d.Set("segment_id", segment.Id)
	}

	d.Set("size", connector.Size)
	d.Set("vpn_mode", connector.VpnMode)

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

	sites := expandConnectorIPSecEndpoint(d.Get("endpoint").(*schema.Set))

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
		Id:             d.Get("id").(json.Number),
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
