package alkira

import (
	"fmt"

	"github.com/alkiranet/alkira-client-go/alkira"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func resourceAlkiraConnectorAzureErVnet() *schema.Resource {
	return &schema.Resource{
		Description: "Manage Azure Cloud Express RouterConnector.",

		Create: resourceConnectorAzureErCreate,
		Read:   resourceConnectorAzureErRead,
		Update: resourceConnectorAzureErUpdate,
		Delete: resourceConnectorAzureErDelete,

		Schema: map[string]*schema.Schema{
			"name": {
				Description: "The name of the connector.",
				Type:        schema.TypeString,
				Required:    true,
			},
			"credential_id": {
				Description: "An opaque identifier generated when storing Azure VNET credentials.",
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
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
				Description: "IP address prefix for VWAN Hub. This should be a /23 prefix.",
				Type:        schema.TypeString,
				Required:    true,
			},
			"tunnel_protocol": {
				Description:  "The tunnel protocol. One of `VXLAN`, `VXLAN_GPE`",
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
				Description: "Tags for billing.",
				Type:        schema.TypeList,
				Optional:    true,
				Elem:        &schema.Schema{Type: schema.TypeInt},
			},
			"instance": {
				Type:     schema.TypeList,
				Required: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Description: "User provided connector instance name",
							Type:        schema.TypeString,
							Required:    true,
						},
						"id": {
							Description: "The ID of the connector instance.",
							Type:        schema.TypeInt,
							Computed:    true,
						},
						"express_route_circuit_id": {
							Description: "Express Route circuit id from Azure. ER Circuit should have a private peering connection provisioned, also an valid authorization key associated with it.",
							Type:        schema.TypeString,
							Required:    true,
						},
						"redundant_router": {
							Description: "Indicates if ER Circuit terminates on redundant routers on customer side.",
							Type:        schema.TypeBool,
							Optional:    true,
							Default:     false,
						},
						"loopback_subnet": {
							Description: "A `/26` subnet from which loopback IPs would be used to establish underlay vXLan GPE tunnels.",
							Type:        schema.TypeString,
							Required:    true,
						},
						"credential_id": {
							Description: "An opaque identifier generated when storing Azure VNET credentials.",
							Type:        schema.TypeString,
							Required:    true,
						},
						"gateway_mac_address": {
							Description: "An array containing the mac addresses of VXLAN gateways reachable through Express Route circuit. The gatewayMacAddresses is only expected if VXLAN tunnel protocol is selected and 2 gateway mac addresses are expected only if redundantRouter is enabled.",
							Type:        schema.TypeList,
							Optional:    true,
							Elem:        &schema.Schema{Type: schema.TypeString},
						},
						"vnis": {
							Description: "This is an optional field if the tunnel protocol is VXLAN. If not specified Alkira allocates unique VNI from the range [16773023, 16777215]",
							Type:        schema.TypeList,
							Optional:    true,
							Elem:        &schema.Schema{Type: schema.TypeInt},
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
							Description: "The name of an existing segment",
							Type:        schema.TypeString,
							Required:    true,
						},
						"segment_id": {
							Description: "The ID of the segment",
							Type:        schema.TypeInt,
							Computed:    true,
						},
						"customer_asn": {
							Description: "ASN on the customer premise side",
							Type:        schema.TypeInt,
							Required:    true,
						},
						"disable_internet_exit": {
							Description: "Enable or disable access to the internet when traffic arrives via this connector",
							Type:        schema.TypeBool,
							Optional:    true,
							Default:     false,
						},
						"advertise_on_prem_routes": {
							Description: "Allow routes from the branch/premises to be advertised to the cloud",
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

func resourceConnectorAzureErCreate(d *schema.ResourceData, m interface{}) error {
	client := m.(*alkira.AlkiraClient)
	connector, err := generateConnectorAzureErRequest(d, m)

	if err != nil {
		return err
	}

	id, err := client.CreateConnectorAzureErVnet(connector)

	if err != nil {
		return err
	}

	d.SetId(id)
	return resourceConnectorAzureErRead(d, m)
}

func resourceConnectorAzureErRead(d *schema.ResourceData, m interface{}) error {
	client := m.(*alkira.AlkiraClient)

	connector, err := client.GetConnectorAzureErVnet(d.Id())

	if err != nil {
		return err
	}

	d.Set("size", connector.Size)
	d.Set("billing_tag_ids", connector.BillingTags)
	d.Set("cxp", connector.Cxp)
	d.Set("enabled", connector.Enabled)
	d.Set("name", connector.Name)
	d.Set("tunnel_protocol", connector.TunnelProtocol)
	d.Set("vhub_prefix", connector.VhubPrefix)

	var instances []map[string]interface{}
	for _, instance := range connector.Instances {
		i := map[string]interface{}{
			"credential_id":            instance.CredentialId,
			"express_route_circuit_id": instance.ExpressRouteCircuitId,
			"gateway_mac_address":      instance.GatewayMacAddress,
			"id":                       instance.Id,
			"loopback_subnet":          instance.LoopbackSubnet,
			"name":                     instance.Name,
			"redundant_router":         instance.RedundantRouter,
			"vnis":                     instance.Vnis,
		}
		instances = append(instances, i)
	}

	d.Set("instance", instances)

	var segments []map[string]interface{}

	for _, seg := range connector.SegmentOptions {
		i := map[string]interface{}{
			"name":       seg.SegmentName,
			"segment_id": seg.SegmentId,
		}
		segments = append(segments, i)
	}
	d.Set("segment_options", segments)

	return nil
}

func resourceConnectorAzureErUpdate(d *schema.ResourceData, m interface{}) error {
	client := m.(*alkira.AlkiraClient)
	connector, err := generateConnectorAzureErRequest(d, m)

	if err != nil {
		return fmt.Errorf("UpdateConnectorAzureErVnet: failed to marshal: %v", err)
		// return err
	}

	err = client.UpdateConnectorAzureErVnet(d.Id(), connector)

	if err != nil {
		return err
	}

	return resourceConnectorAzureVnetRead(d, m)
}

func resourceConnectorAzureErDelete(d *schema.ResourceData, m interface{}) error {
	client := m.(*alkira.AlkiraClient)

	err := client.DeleteConnectorAzureErVnet((d.Id()))

	return err
}

// generateConnectorAzureErRequest generate request for connector-azure-vnet
func generateConnectorAzureErRequest(d *schema.ResourceData, m interface{}) (*alkira.ConnectorAzureErVnet, error) {
	// client := m.(*alkira.AlkiraClient)

	billingTags := convertTypeListToIntList(d.Get("billing_tag_ids").([]interface{}))

	instances, err := expandAzureErInstances(d.Get("instance").([]interface{}), m)
	if err != nil {
		return nil, err
	}

	segmentOptions, err := expandAzureErSegments(d.Get("segment_options").([]interface{}), m)
	if err != nil {
		return nil, err
	}

	request := &alkira.ConnectorAzureErVnet{
		Name:           d.Get("name").(string),
		Size:           d.Get("size").(string),
		CredentialId:   d.Get("credential_id").(string),
		BillingTags:    billingTags,
		Enabled:        d.Get("enabled").(bool),
		TunnelProtocol: d.Get("tunnel_protocol").(string),
		Cxp:            d.Get("cxp").(string),
		VhubPrefix:     d.Get("vhub_prefix").(string),
		Instances:      instances,
		SegmentOptions: segmentOptions,
	}

	// return request, fmt.Errorf("helloll: %s", d.Get("instance"))

	return request, nil
}
