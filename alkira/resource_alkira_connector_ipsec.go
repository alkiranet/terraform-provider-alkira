package alkira

import (
	"context"
	"fmt"

	"github.com/alkiranet/alkira-client-go/alkira"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func resourceAlkiraConnectorIPSec() *schema.Resource {
	return &schema.Resource{
		Description:   "Manage IPSec Connector.",
		CreateContext: resourceConnectorIPSecCreate,
		ReadContext:   resourceConnectorIPSecRead,
		UpdateContext: resourceConnectorIPSecUpdate,
		DeleteContext: resourceConnectorIPSecDelete,
		CustomizeDiff: func(ctx context.Context, d *schema.ResourceDiff, m interface{}) error {
			client := m.(*alkira.AlkiraClient)

			old, _ := d.GetChange("provision_state")

			if client.Provision && old == "FAILED" {
				d.SetNew("provision_state", "SUCCESS")
			}

			return nil
		},
		Importer: &schema.ResourceImporter{
			StateContext: importWithReadValidation(resourceConnectorIPSecRead),
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
			"cxp": {
				Description: "The CXP where the connector should be " +
					"provisioned.",
				Type:     schema.TypeString,
				Required: true,
			},
			"enabled": {
				Description: "Is the connector enabled. Default is `true`.",
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     true,
			},
			"endpoint": {
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
						"customer_ip_type": {
							Description: "The type of `customer_gateway_ip`. It " +
								"could be either `STATIC` or `DYNAMIC`. " +
								"Default value is `STATIC`. When it's `DYNAMIC`, " +
								"`customer_gateway_ip` should be set to `0.0.0.0`. " +
								"`remote_auth_type` in `advanced_options` is " +
								"required as well.",
							Type:         schema.TypeString,
							Optional:     true,
							Default:      "STATIC",
							ValidateFunc: validation.StringInSlice([]string{"STATIC", "DYNAMIC"}, false),
						},
						"id": {
							Description: "The ID of the endpoint.",
							Type:        schema.TypeInt,
							Computed:    true,
						},
						"preshared_keys": {
							Description: "An array of preshared keys, one per " +
								"tunnel. The value needs to be provided explicitly.",
							Type: schema.TypeList,
							Elem: &schema.Schema{
								Type:         schema.TypeString,
								ValidateFunc: validation.StringIsNotWhiteSpace,
							},
							Required: true,
						},
						"enable_tunnel_redundancy": {
							Description: "Disable this if all tunnels will not " +
								"be configured or enabled on the on-premise " +
								"device. If it's set to `false`, connector " +
								"health will be shown as `UP` if at least " +
								"one of the tunnels is `UP`. If enabled, " +
								"all tunnels need to be `UP` for the connector" +
								"health to be shown as `UP`.",
							Type:     schema.TypeBool,
							Optional: true,
							Default:  true,
						},
						"billing_tag_ids": {
							Description: "Billing tags to be associated with " +
								"the resource. (see resource `alkira_billing_tag`).",
							Type:     schema.TypeSet,
							Optional: true,
							Elem: &schema.Schema{
								Type: schema.TypeInt,
							},
						},
						"ha_mode": {
							Description: "The value could be `ACTIVE` or `STANDBY`. " +
								"A endpoint in `STANDBY` mode will not be used for " +
								"traffic unless all other endpoints for the " +
								"connector are down. There can only be one " +
								"endpoint in `STANDBY` mode per connector and " +
								"there must be at least one endpoint " +
								"that isn't in `STANDBY` mode per connector.",
							Type:     schema.TypeString,
							Optional: true,
							ValidateFunc: validation.StringInSlice(
								[]string{"ACTIVE", "STANDBY"}, false),
						},
						"advanced_options": {
							Description: "Advanced options for IPSec endpoint.",
							Type:        schema.TypeList,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"esp_encryption_algorithms": {
										Description: "Encryption algorithms to " +
											"use for IPsec SA. Value " +
											"could be `AES256CBC`, `AES192CBC`, " +
											"`AES128CBC`, `AES256GCM16` " +
											"`3DESCBC`, or `NULL`.",
										Type:     schema.TypeList,
										Elem:     &schema.Schema{Type: schema.TypeString},
										Required: true,
									},
									"esp_integrity_algorithms": {
										Description: "Integrity algorithms to " +
											"use for IPsec SA. Value could " +
											"`SHA1`, `SHA256`, `SHA384`, " +
											"`SHA512` or `MD5`.",
										Type:     schema.TypeList,
										Elem:     &schema.Schema{Type: schema.TypeString},
										Required: true,
									},
									"esp_dh_group_numbers": {
										Description: "Diffie Hellman groups to use " +
											"for IPsec SA. Value could " +
											"`MODP1024`, `MODP2048`, `MODP3072`, " +
											"`MODP4096`, `MODP6144`, " +
											"`MODP8192`, `ECP256`, `ECP384`, " +
											"`ECP521`, `CURVE25519` and `NONE`.",
										Type:     schema.TypeList,
										Elem:     &schema.Schema{Type: schema.TypeString},
										Required: true,
									},
									"initiator": {
										Description: "When the value is `false`, " +
											"CXP will initiate the IKE connection " +
											"and in all other cases the customer " +
											"gateway should initiate IKE connection. " +
											"When `gateway_ip_type` is `DYNAMIC`, " +
											"initiator must be `true`.",
										Type:     schema.TypeBool,
										Required: true,
									},
									"ike_version": {
										Description: "IKE version, either `IKEv1` " +
											"or `IKEv2`",
										Type:     schema.TypeString,
										Required: true,
										ValidateFunc: validation.StringInSlice(
											[]string{"IKEv1", "IKEv2"}, false),
									},
									"ike_encryption_algorithms": {
										Description: "Encryption algorithms to " +
											"use for IKE SA, one of " +
											"`AES256CBC`, `AES192CBC`, " +
											"`AES128CBC` and `3DESCBC`.",
										Type:     schema.TypeList,
										Elem:     &schema.Schema{Type: schema.TypeString},
										Required: true,
									},
									"ike_integrity_algorithms": {
										Description: "Integrity algorithms to use " +
											"for IKE SA, one of " +
											"`SHA1`, `SHA256`, `SHA384`, `SHA512`.",
										Type:     schema.TypeList,
										Elem:     &schema.Schema{Type: schema.TypeString},
										Required: true,
									},
									"ike_dh_group_numbers": {
										Description: "Diffie Hellman groups to " +
											"use for IKE SA, one of `MODP1024`, " +
											"`MODP2048`, `MODP3072`, `MODP4096`, " +
											"`MODP6144`, `MODP8192`, `ECP256`, " +
											"`ECP384`, `ECP521`, or `CURVE25519`.",
										Type:     schema.TypeList,
										Elem:     &schema.Schema{Type: schema.TypeString},
										Required: true,
									},
									"remote_auth_type": {
										Description: "IKE identity to use for " +
											"authentication round, one of " +
											"`FQDN`, `USER_FQDN`, " +
											"`KEYID`, or `IP_ADDR`.",
										Type:     schema.TypeString,
										Required: true,
										ValidateFunc: validation.StringInSlice(
											[]string{"FQDN", "USER_FQDN", "KEYID", "IP_ADDR"}, false),
									},
									"remote_auth_value": {
										Description: "Remote-ID value.",
										Type:        schema.TypeString,
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
				Description: "The group of the connector. (see resource " +
					"`alkira_group`)",
				Type:     schema.TypeString,
				Optional: true,
			},
			"implicit_group_id": {
				Description: "The ID of implicit group automaticaly created " +
					"with the connector.",
				Type:     schema.TypeInt,
				Computed: true,
			},
			"provision_state": {
				Description: "The provision state of the connector.",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"policy_options": {
				Description: "Policy options, both `on_prem_prefix_list_ids` " +
					"and `cxp_prefix_list_ids` must be provided if `vpn_mode` " +
					"is `POLICY_BASED`.",
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
							Description: "Routing type, one of `STATIC`, " +
								"`DYNAMIC`, or `BOTH`.",
							Type:     schema.TypeString,
							Required: true,
						},
						"availability": {
							Description: "The method to determine the " +
								"availability of the routes. The value " +
								"could be `IKE_STATUS` or " +
								"`IPSEC_INTERFACE_PING`. Default value is " +
								"`IPSEC_INTERFACE_PING`.",
							Type: schema.TypeString,
							ValidateFunc: validation.StringInSlice(
								[]string{"IKE_STATUS", "IPSEC_INTERFACE_PING", "PING"}, false),
							Optional: true,
							Default:  "IPSEC_INTERFACE_PING",
						},
						"prefix_list_id": {
							Description: "The ID of prefix list to use for " +
								"static route propagation.",
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
			"segment_options": {
				Description: "Additional options for each segment associated " +
					"with the connector.",
				Type: schema.TypeSet,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Description: "Segment Name.",
							Type:        schema.TypeString,
							Required:    true,
						},
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
					},
				},
				Optional: true,
			},
			"scale_group_id": {
				Description: "The ID of the scale group associated with the connector.",
				Type:        schema.TypeString,
				Optional:    true,
			},
			"segment_id": {
				Description: "The ID of the segment associated with the connector.",
				Type:        schema.TypeString,
				Required:    true,
			},
			"size": {
				Description: "The size of the connector, one of `SMALL`, " +
					"`MEDIUM`, `LARGE`, `2LARGE`, `5LARGE`, " +
					"`10LARGE`.",
				Type:     schema.TypeString,
				Required: true,
			},
			"vpn_mode": {
				Description: "The mode can be configured either as `ROUTE_BASED` " +
					"or `POLICY_BASED`.",
				Type:     schema.TypeString,
				Required: true,
				ValidateFunc: validation.StringInSlice(
					[]string{"ROUTE_BASED", "POLICY_BASED"}, false),
			},
		},
	}
}

func resourceConnectorIPSecCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {

	// INIT
	client := m.(*alkira.AlkiraClient)
	api := alkira.NewConnectorIPSec(m.(*alkira.AlkiraClient))

	request, err := generateConnectorIPSecRequest(d, m)

	if err != nil {
		return diag.FromErr(err)
	}

	// CREATE
	response, provState, err, valErr, provErr := api.Create(request)

	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(string(response.Id))

	// Handle validation error
	if client.Validate && valErr != nil {
		var diags diag.Diagnostics
		readDiags := resourceConnectorIPSecRead(ctx, d, m)
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

	// Set state
	if client.Provision {
		d.Set("provision_state", provState)
		if provErr != nil {
			return diag.Diagnostics{{
				Severity: diag.Warning,
				Summary:  "PROVISION (CREATE) FAILED",
				Detail:   fmt.Sprintf("%s", provErr),
			}}
		}
	}

	return resourceConnectorIPSecRead(ctx, d, m)
}

func resourceConnectorIPSecRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {

	// INIT
	client := m.(*alkira.AlkiraClient)
	api := alkira.NewConnectorIPSec(m.(*alkira.AlkiraClient))

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
	d.Set("cxp", connector.CXP)
	d.Set("enabled", connector.Enabled)
	d.Set("group", connector.Group)
	d.Set("implicit_group_id", connector.ImplicitGroupId)
	d.Set("name", connector.Name)
	d.Set("scale_group_id", connector.ScaleGroupId)
	d.Set("size", connector.Size)
	d.Set("vpn_mode", connector.VpnMode)
	d.Set("description", connector.Description)

	// Get segment
	numOfSegments := len(connector.Segments)
	if numOfSegments == 1 {
		segmentId, err := getSegmentIdByName(connector.Segments[0], m)

		if err != nil {
			return diag.FromErr(err)
		}
		d.Set("segment_id", segmentId)
	} else {
		return diag.FromErr(fmt.Errorf("failed to find segment"))
	}

	// Set segment_options from API response
	if connector.SegmentOptions != nil {
		flatSegmentOpts := flattenConnectorIPSecSegmentOptions(connector.SegmentOptions)
		if flatSegmentOpts != nil {
			d.Set("segment_options", flatSegmentOpts)
		}
	}

	// Set routing_options and policy_options based on vpn_mode
	switch connector.VpnMode {
	case "ROUTE_BASED":
		if connector.RoutingOptions != nil {
			flatRoutingOpts := flattenConnectorIPSecRoutingOptions(connector.RoutingOptions)
			if flatRoutingOpts != nil {
				d.Set("routing_options", flatRoutingOpts)
			}
		}
	case "POLICY_BASED":
		if connector.PolicyOptions != nil {
			flatPolicyOpts := flattenConnectorIPSecPolicyOptions(connector.PolicyOptions)
			if flatPolicyOpts != nil {
				d.Set("policy_options", flatPolicyOpts)
			}
		}
	}

	//
	// Go through all endpoints from the config firstly to find a
	// match, either endpoint's ID or endpoint's name should be
	// uniquely identifying an endpoint.
	//
	// On the first read call at the end of the create call, Terraform
	// didn't track any endpoint IDs yet.
	//
	var endpoints []map[string]interface{}

	for _, endpoint := range d.Get("endpoint").([]interface{}) {
		endpointConfig := endpoint.(map[string]interface{})

		for _, site := range connector.Sites {
			if endpointConfig["id"].(int) == site.Id || endpointConfig["name"].(string) == site.Name {
				// Get the configured preshared_keys count from user's config
				configuredKeyCount := 0
				if keys, ok := endpointConfig["preshared_keys"].([]interface{}); ok {
					configuredKeyCount = len(keys)
				}
				endpoint := setConnectorIPSecEndpoint(site, configuredKeyCount)
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
			// New endpoint not in config, pass 0 to disable deduplication
			endpoint := setConnectorIPSecEndpoint(site, 0)
			endpoints = append(endpoints, endpoint)
			break
		}
	}

	// Set provision state
	if client.Provision && provState != "" {
		d.Set("provision_state", provState)
	}

	d.Set("endpoint", endpoints)

	return nil
}

func resourceConnectorIPSecUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {

	// INIT
	client := m.(*alkira.AlkiraClient)
	api := alkira.NewConnectorIPSec(m.(*alkira.AlkiraClient))

	request, err := generateConnectorIPSecRequest(d, m)

	if err != nil {
		return diag.FromErr(err)
	}

	// UPDATE
	provState, err, valErr, provErr := api.Update(d.Id(), request)

	if err != nil {
		return diag.FromErr(err)
	}

	// Handle validation error
	if client.Validate && valErr != nil {
		var diags diag.Diagnostics
		readDiags := resourceConnectorIPSecRead(ctx, d, m)
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

	if client.Provision {
		d.Set("provision_state", provState)
		if provErr != nil {
			return diag.Diagnostics{{
				Severity: diag.Warning,
				Summary:  "PROVISION (UPDATE) FAILED",
				Detail:   fmt.Sprintf("%s", provErr),
			}}
		}
	}

	return resourceConnectorIPSecRead(ctx, d, m)
}

func resourceConnectorIPSecDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {

	// INIT
	client := m.(*alkira.AlkiraClient)
	api := alkira.NewConnectorIPSec(m.(*alkira.AlkiraClient))

	// DELETE
	provState, err, valErr, provErr := api.Delete(d.Id())

	if err != nil {
		// Terraform may not print "with <resource address>" for destroys of objects
		// that are no longer in configuration, so include identifying context here.
		name, _ := d.GetOk("name")
		if nameStr, ok := name.(string); ok && nameStr != "" {
			return diag.FromErr(fmt.Errorf("%w alkira_connector_ipsec (name=%q id=%s)", err, nameStr, d.Id()))
		}
		return diag.FromErr(fmt.Errorf("%w alkira_connector_ipsec (id=%s)", err, d.Id()))
	}

	d.SetId("")

	// Handle validation error
	if client.Validate && valErr != nil {
		return diag.Diagnostics{{
			Severity: diag.Error,
			Summary:  "VALIDATION (DELETE) FAILED",
			Detail:   fmt.Sprintf("%s", valErr),
		}}
	}

	if client.Provision && provState != "SUCCESS" {
		return diag.Diagnostics{{
			Severity: diag.Warning,
			Summary:  "PROVISION (DELETE) FAILED",
			Detail:   fmt.Sprintf("%s", provErr),
		}}
	}

	return nil
}

// generateConnectorIPSecRequest generate request for connector-ipsec
func generateConnectorIPSecRequest(d *schema.ResourceData, m interface{}) (*alkira.ConnectorIPSec, error) {

	sites := expandConnectorIPSecEndpoint(d.Get("endpoint").([]interface{}))

	//
	// Segment
	//
	segmentName, err := getSegmentNameById(d.Get("segment_id").(string), m)

	if err != nil {
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
		Segments:       []string{segmentName},
		Sites:          sites,
		ScaleGroupId:   d.Get("scale_group_id").(string),
		Size:           d.Get("size").(string),
		VpnMode:        vpnMode,
		Description:    d.Get("description").(string),
	}

	return connector, nil
}
