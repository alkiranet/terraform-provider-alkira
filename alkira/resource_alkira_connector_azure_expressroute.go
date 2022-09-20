package alkira

import (
	"fmt"

	"github.com/alkiranet/alkira-client-go/alkira"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func resourceAlkiraConnectorAzureExpressRoute() *schema.Resource {
	return &schema.Resource{
		Description: "Manage Azure ExpressRoute Connector.",

		Create: resourceConnectorAzureExpressRouteCreate,
		Read:   resourceConnectorAzureExpressRouteRead,
		Update: resourceConnectorAzureExpressRouteUpdate,
		Delete: resourceConnectorAzureExpressRouteDelete,

		Schema: map[string]*schema.Schema{
			"name": {
				Description: "The name of the connector.",
				Type:        schema.TypeString,
				Required:    true,
			},
			"size": {
				Description:  "The size of the connector, one of `LARGE`, `2LARGE`, `5LARGE`, `10LARGE`.",
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringInSlice([]string{"LARGE", `2LARGE`, `5LARGE`, `10LARGE`}, false),
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
			"billing_tag_ids": {
				Description: "IDs of Billing Tags.",
				Type:        schema.TypeList,
				Optional:    true,
				Elem:        &schema.Schema{Type: schema.TypeInt},
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
								"ExpresRoute Circuit should have a private peering connection provisioned, " +
								"also an valid authorization key associated with it.",
							Type:     schema.TypeString,
							Required: true,
						},
						"redundant_router": {
							Description: "Indicates if ExpressRoute Circuit terminates on redundant routers on customer side.",
							Type:        schema.TypeBool,
							Optional:    true,
							Default:     false,
						},
						"loopback_subnet": {
							Description: "A `/26` subnet from which loopback IPs would be used to establish underlay VXLAN GPE tunnels.",
							Type:        schema.TypeString,
							Required:    true,
						},
						"credential_id": {
							Description: "An opaque identifier generated when storing Azure VNET credentials.",
							Type:        schema.TypeString,
							Required:    true,
						},
						"gateway_mac_address": {
							Description: "An array containing the mac addresses of VXLAN gateways reachable through ExpressRoute circuit. " +
								"The gatewayMacAddresses is only expected if VXLAN tunnel protocol is selected, " +
								"and 2 gateway mac addresses are expected only if redundant_router is enabled.",
							Type:     schema.TypeList,
							Optional: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
						"virtual_network_interface": {
							Description: "This is an optional field if the tunnel protocol is VXLAN. " +
								"If not specified Alkira allocates unique VNI from the range [16773023, 16777215].",
							Type:     schema.TypeList,
							Optional: true,
							Elem:     &schema.Schema{Type: schema.TypeInt},
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
							Description: "Enable or disable access to the internet when traffic arrives via this connector.",
							Type:        schema.TypeBool,
							Optional:    true,
							Default:     true,
						},
						"advertise_on_prem_routes": {
							Description: "Allow routes from the branch/premises to be advertised to the cloud.",
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

// resourceConnectorAzureExpressRouteCreate create an Azure ExpressRoute connector
func resourceConnectorAzureExpressRouteCreate(d *schema.ResourceData, m interface{}) error {
	client := m.(*alkira.AlkiraClient)
	connector, err := generateConnectorAzureExpressRouteRequest(d, m)

	if err != nil {
		return err
	}

	id, err := client.CreateConnectorAzureExpressRoute(connector)

	if err != nil {
		return err
	}

	d.SetId(id)
	return resourceConnectorAzureExpressRouteRead(d, m)
}

// resourceConnectorAzureExpressRouteRead get and save an Azure ExpressRoute connectors
func resourceConnectorAzureExpressRouteRead(d *schema.ResourceData, m interface{}) error {
	client := m.(*alkira.AlkiraClient)

	connector, err := client.GetConnectorAzureExpressRoute(d.Id())

	if err != nil {
		return err
	}

	d.Set("size", connector.Size)
	d.Set("billing_tag_ids", connector.BillingTags)
	d.Set("cxp", connector.Cxp)
	d.Set("group", connector.Group)
	d.Set("enabled", connector.Enabled)
	d.Set("name", connector.Name)
	d.Set("tunnel_protocol", connector.TunnelProtocol)
	d.Set("vhub_prefix", connector.VhubPrefix)

	var instances []map[string]interface{}
	for _, instance := range connector.Instances {
		i := map[string]interface{}{
			"credential_id":             instance.CredentialId,
			"expressroute_circuit_id":   instance.ExpressRouteCircuitId,
			"gateway_mac_address":       instance.GatewayMacAddress,
			"id":                        instance.Id,
			"loopback_subnet":           instance.LoopbackSubnet,
			"name":                      instance.Name,
			"redundant_router":          instance.RedundantRouter,
			"virtual_network_interface": instance.Vnis,
		}
		instances = append(instances, i)
	}

	d.Set("instances", instances)

	var segments []map[string]interface{}

	for _, seg := range connector.SegmentOptions {
		i := map[string]interface{}{
			"segment_name":             seg.SegmentName,
			"segment_id":               seg.SegmentId,
			"customer_asn":             seg.CustomerAsn,
			"disable_internet_exit":    seg.DisableInternetExit,
			"advertise_on_prem_routes": seg.AdvertiseOnPremRoutes,
		}
		segments = append(segments, i)
	}
	d.Set("segment_options", segments)

	return nil
}

// resourceConnectorAzureExpressRouteUpdate update an Azure ExpressRoute connector
func resourceConnectorAzureExpressRouteUpdate(d *schema.ResourceData, m interface{}) error {
	client := m.(*alkira.AlkiraClient)
	connector, err := generateConnectorAzureExpressRouteRequest(d, m)

	if err != nil {
		return fmt.Errorf("UpdateConnectorAzureExpressRoute: failed to marshal: %v", err)
	}

	err = client.UpdateConnectorAzureExpressRoute(d.Id(), connector)

	if err != nil {
		return err
	}

	return resourceConnectorAzureExpressRouteRead(d, m)
}

func resourceConnectorAzureExpressRouteDelete(d *schema.ResourceData, m interface{}) error {
	client := m.(*alkira.AlkiraClient)

	err := client.DeleteConnectorAzureExpressRoute((d.Id()))

	return err
}

// generateConnectorAzureExpressRouteRequest generate a request for Azure ExpressRoute connector
func generateConnectorAzureExpressRouteRequest(d *schema.ResourceData, m interface{}) (*alkira.ConnectorAzureExpressRoute, error) {

	billingTags := convertTypeListToIntList(d.Get("billing_tag_ids").([]interface{}))

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
