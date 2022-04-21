package alkira

import (
	"github.com/alkiranet/alkira-client-go/alkira"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func resourceAlkiraConnectorArubaEdge() *schema.Resource {
	return &schema.Resource{
		Description: "Manage Aruba Edge Connector",

		Create: resourceConnectorArubaEdgeCreate,
		Read:   resourceConnectorArubaEdgeRead,
		Update: resourceConnectorArubaEdgeUpdate,
		Delete: resourceConnectorArubaEdgeDelete,

		Schema: map[string]*schema.Schema{
			"aruba_edge_vrf_mapping": {
				Description: "The connector will accept multiple segments as a part of VRF mappings.",
				Type:        schema.TypeSet,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"advertise_on_prem_routes": {
							Description: "Allow routes from the branch/premises to be advertised to the cloud. The default value is False.",
							Type:        schema.TypeBool,
							Optional:    true,
							Default:     false,
						},
						"alkira_segment_id": {
							Description: "The segment id associated with the Aruba Edge connector.",
							Type:        schema.TypeString,
							Required:    true,
						},
						"aruba_edge_connect_segment_name": {
							Description: "The segment name of the Aruba Edge connector.",
							Type:        schema.TypeString,
							Required:    true,
						},
						"disable_internet_exit": {
							Description: "Enables or disables access to the internet when traffic arrives via this connector. The default value is False.",
							Type:        schema.TypeBool,
							Optional:    true,
							Default:     false,
						},
						"gateway_gbp_asn": {
							Description: "The gateway BGP ASN.",
							Type:        schema.TypeInt,
							Required:    true,
						},
					},
				},
				Optional: true,
			},
			"billing_tag_ids": {
				Description: "Tags for billing.",
				Type:        schema.TypeList,
				Optional:    true,
				Elem:        &schema.Schema{Type: schema.TypeInt},
			},
			"boost_mode": {
				Description: "If enabled the Aruba Edge Connect image supporting the boost mode " +
					"for given size(or bandwidth) would be deployed in Alkira CXP.",
				Type:     schema.TypeBool,
				Required: true,
			},
			"cxp": {
				Description: "The CXP where the connector should be provisioned.",
				Type:        schema.TypeString,
				Required:    true,
			},
			"gateway_gbp_asn": {
				Description: "The gateway BGP ASN.",
				Type:        schema.TypeInt,
				Required:    true,
			},
			"group": {
				Description: "The group of the connector.",
				Type:        schema.TypeString,
				Optional:    true,
			},
			"instances": {
				Description: "The Aruba Edge connector instances.",
				Type:        schema.TypeSet,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"account_name": {
							Description: "The account name given in SilverPeak orchestrator registration.",
							Type:        schema.TypeString,
							Required:    true,
						},
						"credential_id": {
							Description: "The credential ID of the Aruba Edge instance.",
							Type:        schema.TypeString,
							Required:    true,
						},
						"host_name": {
							Description: "The host name given to the Aruba SD-WAN appliance that appears in SilverPeak orchestrator.",
							Type:        schema.TypeString,
							Required:    true,
						},
						"name": {
							Description: "The instance name associated with aruba edge connect instance.",
							Type:        schema.TypeString,
							Required:    true,
						},
						"site_tag": {
							Description: "site tag that appears on the SD-WAN appliance on SilverPeak orchestrator",
							Type:        schema.TypeString,
							Required:    true,
						},
					},
				},
				Required: true,
			},
			"name": {
				Description: "The name of the connector.",
				Type:        schema.TypeString,
				Required:    true,
			},
			"segment_names": {
				Description: "The names of the segments associated with the Aruba Edge connector.",
				Type:        schema.TypeList,
				Required:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
			"size": {
				Description: "The size of the connector, one of `SMALL`, `MEDIUM` or `LARGE`.",
				Type:        schema.TypeString,
				Required:    true,
			},
			"tunnel_protocol": {
				Description:  "Tunnel Protocol, default to `IPSEC`, could be either `IPSEC` or `GRE`.",
				Type:         schema.TypeString,
				Optional:     true,
				Default:      "IPSEC",
				ValidateFunc: validation.StringInSlice([]string{"IPSEC", "GRE"}, false),
			},
			"version": {
				Description: "The version of the Aruba Edge connector.",
				Type:        schema.TypeString,
				Required:    true,
			},
		},
	}
}

func resourceConnectorArubaEdgeCreate(d *schema.ResourceData, m interface{}) error {
	client := m.(*alkira.AlkiraClient)

	connector, err := generateConnectorArubaEdgeRequest(d, m)
	if err != nil {
		return err
	}

	id, err := client.CreateConnectorArubaEdge(connector)
	if err != nil {
		return err
	}

	d.SetId(id)
	return resourceConnectorArubaEdgeRead(d, m)
}

func resourceConnectorArubaEdgeRead(d *schema.ResourceData, m interface{}) error {
	client := m.(*alkira.AlkiraClient)

	connector, err := client.GetConnectorArubaEdgeById(d.Id())

	if err != nil {
		return err
	}

	setArubaEdgeResourceFields(connector, d)

	return err
}

func resourceConnectorArubaEdgeUpdate(d *schema.ResourceData, m interface{}) error {
	client := m.(*alkira.AlkiraClient)

	connector, err := generateConnectorArubaEdgeRequest(d, m)
	if err != nil {
		return err
	}

	err = client.UpdateConnectorArubaEdge(d.Id(), connector)
	if err != nil {
		return err
	}

	return resourceConnectorArubaEdgeRead(d, m)
}

func resourceConnectorArubaEdgeDelete(d *schema.ResourceData, m interface{}) error {
	client := m.(*alkira.AlkiraClient)

	return client.DeleteConnectorArubaEdge(d.Id())
}

func generateConnectorArubaEdgeRequest(d *schema.ResourceData, m interface{}) (*alkira.ConnectorArubaEdge, error) {
	return &alkira.ConnectorArubaEdge{
		ArubaEdgeVrfMapping: expandArubeEdgeVrfMapping(d.Get("aruba_edge_vrf_mapping").(*schema.Set)),
		BillingTags:         convertTypeListToIntList(d.Get("billing_tag_ids").([]interface{})),
		BoostMode:           d.Get("boost_mode").(bool),
		Cxp:                 d.Get("cxp").(string),
		GatewayBgpAsn:       d.Get("gateway_gbp_asn").(int),
		Group:               d.Get("group").(string),
		Instances:           expandArubaEdgeInstances(d.Get("instances").(*schema.Set)),
		Name:                d.Get("name").(string),
		Segments:            convertTypeListToStringList(d.Get("segment_names").([]interface{})),
		Size:                d.Get("size").(string),
		TunnelProtocol:      d.Get("tunnel_protocol").(string),
		Version:             d.Get("version").(string),
	}, nil
}
