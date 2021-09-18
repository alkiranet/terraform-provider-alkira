package alkira

import (
	"log"

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

		Schema: map[string]*schema.Schema{
			"name": {
				Description: "The name of the connector.",
				Type:        schema.TypeString,
				Required:    true,
			},
			"billing_tags": {
				Description: "A list of Billing Tag by Id associated with the connector.",
				Type:        schema.TypeList,
				Optional:    true,
				Elem:        &schema.Schema{Type: schema.TypeInt},
			},
			"cxp": {
				Description: "The CXP where the connector should be provisioned.",
				Type:        schema.TypeString,
				Required:    true,
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
				Description: "Policy options, both on-prem and cxp prefix list ids must be provided if vpnMode is `POLICY_BASED`",
				Type:        schema.TypeSet,
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
				Description: "Routing options, type is `STATIC`, `DYNAMIC`, or `BOTH` must be provided if `vpn_mode` is `ROUTE_BASED`",
				Type:        schema.TypeSet,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"type": {
							Description: "Routing type, one of `STATIC`, `DYNAMIC`, or `BOTH`.",
							Type:        schema.TypeString,
							Required:    true,
						},
						"availability": {
							Description: "The method to determine the availability of static route. The value could be `IKE_STATUS` or `IPSEC_INTERFACE_PING`.",
							Type:        schema.TypeString,
							ValidateFunc: validation.StringInSlice([]string{"IKE_STATUS", "IPSEC_INTERFACE_PING"}, false),
							Optional:    true,
						},
						"prefix_list_id": {
							Description: "The id of prefix list to use for static route propagation.",
							Type:        schema.TypeInt,
							Optional:    true,
						},
						"customer_gateway_asn": {
							Description: "The customer gateway ASN to use for dynamic route propagation.",
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
						"disable_internet_exit": {
							Description: "Enable or disable access to the internet when traffic arrives via this connector.",
							Type:        schema.TypeBool,
							Optional:    true,
						},

						"disable_advertise_on_prem_routes": {
							Description: "Additional options for each segment associated with the connector.",
							Type:        schema.TypeBool,
							Optional:    true,
						},
					},
				},
				Optional: true,
			},
			"segment": {
				Description: "The name of the segment associated with the connector.",
				Type:        schema.TypeString,
				Required:    true,
			},
			"size": &schema.Schema{
				Description:  "The size of the connector. one of `SMALL`, `MEDIUM` and `LARGE`.",
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringInSlice([]string{"SMALL", "MEDIUM", "LARGE"}, false),
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

	log.Printf("[INFO] Creating Connector (IPSec)")
	id, err := client.CreateConnectorIPSec(connector)

	if err != nil {
		return err
	}

	d.SetId(id)

	return resourceConnectorIPSecRead(d, m)
}

func resourceConnectorIPSecRead(d *schema.ResourceData, m interface{}) error {
	return nil
}

func resourceConnectorIPSecUpdate(d *schema.ResourceData, m interface{}) error {
	client := m.(*alkira.AlkiraClient)

	connector, err := generateConnectorIPSecRequest(d, m)

	if err != nil {
		return err
	}

	log.Printf("[INFO] Updating Connector (IPSec) %s", d.Id())
	err = client.UpdateConnectorIPSec(d.Id(), connector)

	if err != nil {
		return err
	}

	return resourceConnectorIPSecRead(d, m)
}

func resourceConnectorIPSecDelete(d *schema.ResourceData, m interface{}) error {
	client := m.(*alkira.AlkiraClient)
	id := d.Id()

	log.Printf("[INFO] Deleting Connector (IPSec) %s", id)
	err := client.DeleteConnectorIPSec(id)

	return err
}

// generateConnectorIPSecRequest generate request for connector-ipsec
func generateConnectorIPSecRequest(d *schema.ResourceData, m interface{}) (*alkira.ConnectorIPSec, error) {
	billingTags := convertTypeListToIntList(d.Get("billing_tags").([]interface{}))
	sites := expandConnectorIPSecEndpoint(d.Get("endpoint").(*schema.Set))

	// For now, IPSec connector only support single segment
	segments := []string{d.Get("segment").(string)}
	segmentOptions, optErr := expandConnectorIPSecSegmentOptions(d.Get("segment_options").(*schema.Set))

	if optErr != nil {
		return nil, optErr
	}

	// Base on the vpn_mode, switch what options to use
	vpnMode := d.Get("vpn_mode").(string)

	var policyOptions *alkira.ConnectorIPSecPolicyOptions
	var routingOptions *alkira.ConnectorIPSecRoutingOptions
	var err error

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
		BillingTags:    billingTags,
		CXP:            d.Get("cxp").(string),
		Group:          d.Get("group").(string),
		Name:           d.Get("name").(string),
		PolicyOptions:  policyOptions,
		RoutingOptions: routingOptions,
		SegmentOptions: segmentOptions,
		Segments:       segments,
		Sites:          sites,
		Size:           d.Get("size").(string),
		VpnMode:        vpnMode,
	}

	return connector, nil
}
