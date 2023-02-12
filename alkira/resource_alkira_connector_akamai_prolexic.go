package alkira

import (
	"log"
	"strconv"

	"github.com/alkiranet/alkira-client-go/alkira"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func resourceAlkiraConnectorAkamaiProlexic() *schema.Resource {
	return &schema.Resource{
		Description: "Manage Connector for Akamai Prolexic. (BETA)\n\n" +
			"This resource is still under active development. It may have further " +
			"changes in the near future. Today, to use this connector, you will need " +
			"to have onboarded a BYOIP with Do Not Advertise set to `true`. Also, the " +
			"segment with public IPs needs to be reported to Akamai Representative.",
		Create: resourceConnectorAkamaiProlexicCreate,
		Read:   resourceConnectorAkamaiProlexicRead,
		Update: resourceConnectorAkamaiProlexicUpdate,
		Delete: resourceConnectorAkamaiProlexicDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Description: "The name of the connector.",
				Type:        schema.TypeString,
				Required:    true,
			},
			"akamai_bgp_asn": {
				Description: "The Akamai BGP ASN.",
				Type:        schema.TypeInt,
				Required:    true,
			},
			"akamai_bgp_authentication_key": {
				Description: "The Akamai BGP Authentication Key.",
				Type:        schema.TypeString,
				Required:    true,
			},
			"byoip_options": &schema.Schema{
				Description: "BYOIP options.",
				Type:        schema.TypeSet,
				Required:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"byoip_prefix_id": {
							Description: "BYOIP prefix ID.",
							Type:        schema.TypeInt,
							Required:    true,
						},
						"enable_route_advertisement": {
							Description: "Whether enabling route advertisement.",
							Type:        schema.TypeBool,
							Required:    true,
						},
					},
				},
			},
			"billing_tag_ids": {
				Description: "A list of Billing Tag by ID associated with the connector.",
				Type:        schema.TypeList,
				Optional:    true,
				Elem:        &schema.Schema{Type: schema.TypeInt},
			},
			"credential_id": {
				Description: "The credential ID for storing Akamai BGP authentication key.",
				Type:        schema.TypeString,
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
			"group": {
				Description: "The group of the connector.",
				Type:        schema.TypeString,
				Optional:    true,
			},
			"provision": {
				Description: "The state of provisioning.",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"implicit_group_id": {
				Description: "The ID of implicit group automaticaly created " +
					"with the connector.",
				Type:     schema.TypeInt,
				Computed: true,
			},
			"size": &schema.Schema{
				Description: "The size of the connector, one of `SMALL`, `MEDIUM`, " +
					"`LARGE`, `2LARGE`, `4LARGE`, `5LARGE`, `10LARGE`, `20LARGE`.",
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringInSlice([]string{"SMALL", "MEDIUM", "LARGE", "2LARGE", "4LARGE", "5LARGE", "10LARGE", "20LARGE"}, false),
			},
			"segment_id": {
				Description: "The ID of segments associated with the connector. " +
					"Currently, only `1` segment is allowed.",
				Type:     schema.TypeInt,
				Required: true,
			},
			"tunnel_configuration": &schema.Schema{
				Description: "Tunnel Configurations.",
				Type:        schema.TypeSet,
				Required:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"alkira_public_ip": {
							Description: "Alkira public IP.",
							Type:        schema.TypeString,
							Required:    true,
						},
						"tunnel_ips": {
							Description: "Tunnel IPs.",
							Type:        schema.TypeSet,
							Required:    true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"ran_tunnel_ip": {
										Description: "The underlay tunnel IP on the Akamai " +
											"side to be used to configure tunnels between " +
											"the Alkira CXP and the Akamai Prolexic service. " +
											"A RAN (Routed Access Network) is the unit of " +
											"availability for the Route GRE 3.0 service.",
										Type:     schema.TypeString,
										Required: true,
									},
									"alkira_overlay_tunnel_ip": {
										Description: "The overlay IP of the GRE tunnel on the Alkira side.",
										Type:        schema.TypeString,
										Required:    true,
									},
									"akamai_overlay_tunnel_ip": {
										Description: "The overlay IP of the GRE tunnel on the Alkira side.",
										Type:        schema.TypeString,
										Required:    true,
									},
								},
							},
						},
					},
				},
			},
		},
	}
}

func resourceConnectorAkamaiProlexicCreate(d *schema.ResourceData, m interface{}) error {
	api := alkira.NewConnectorAkamaiProlexic(m.(*alkira.AlkiraClient))

	log.Printf("[INFO] Creating Connector (Akamai-Prolexic)")
	connector, err := generateConnectorAkamaiProlexicRequest(d, m)

	if err != nil {
		return err
	}

	resource, provision, err := api.Create(connector)

	if err != nil {
		return err
	}

	d.Set("provision", provision)
	d.SetId(string(resource.Id))

	return resourceConnectorAkamaiProlexicRead(d, m)
}

func resourceConnectorAkamaiProlexicRead(d *schema.ResourceData, m interface{}) error {
	api := alkira.NewConnectorAkamaiProlexic(m.(*alkira.AlkiraClient))

	connector, err := api.GetById(d.Id())

	if err != nil {
		return err
	}

	d.Set("akamai_bgp_asn", connector.AkamaiBgpAsn)
	d.Set("billing_tag_ids", connector.BillingTags)
	d.Set("cxp", connector.CXP)
	d.Set("enabled", connector.Enabled)
	d.Set("group", connector.Group)
	d.Set("implicit_group_id", connector.ImplicitGroupId)
	d.Set("name", connector.Name)
	d.Set("size", connector.Size)

	// segment_id
	if len(connector.Segments) > 0 {
		segApi := alkira.NewSegment(m.(*alkira.AlkiraClient))
		segment, _, err := segApi.GetByName(connector.Segments[0])

		if err != nil {
			return err
		}
		d.Set("segment_id", segment.Id)
	}

	// byoip_options
	var options []map[string]interface{}

	for _, option := range connector.ByoipOptions {
		i := map[string]interface{}{
			"byoip_prefix_id":            option.ByoipId,
			"enable_route_advertisement": option.RouteAdvertisementEnabled,
		}
		options = append(options, i)
	}

	d.Set("byoip_options", options)

	return nil
}

func resourceConnectorAkamaiProlexicUpdate(d *schema.ResourceData, m interface{}) error {
	api := alkira.NewConnectorAkamaiProlexic(m.(*alkira.AlkiraClient))

	connector, err := generateConnectorAkamaiProlexicRequest(d, m)

	if err != nil {
		return err
	}

	log.Printf("[INFO] Updating Connector (Akamai-Prolexic): %s", d.Id())
	provisionState, err := api.Update(d.Id(), connector)

	if err != nil {
		return err
	}

	d.Set("provision", provisionState)

	return resourceConnectorAkamaiProlexicRead(d, m)
}

func resourceConnectorAkamaiProlexicDelete(d *schema.ResourceData, m interface{}) error {
	api := alkira.NewConnectorAkamaiProlexic(m.(*alkira.AlkiraClient))

	log.Printf("[INFO] Deleting Connector (Akamai-Prolexic): %s", d.Id())
	provisionState, err := api.Delete(d.Id())

	if provisionState != "SUCCESS" {
		log.Printf("[ERROR] De-provision connector-akamai-prolexic failed: %s", d.Id())
	}

	return err
}

// generateConnectorAkamaiProlexicRequest generate request for the connector
func generateConnectorAkamaiProlexicRequest(d *schema.ResourceData, m interface{}) (*alkira.ConnectorAkamaiProlexic, error) {
	client := m.(*alkira.AlkiraClient)

	billingTags := convertTypeListToIntList(d.Get("billing_tag_ids").([]interface{}))
	byoipOptions := expandConnectorAkamaiByoipOptions(d.Get("byoip_options").(*schema.Set))
	tunnelConfigurations := expandConnectorAkamaiTunnelConfiguration(d.Get("tunnel_configuration").(*schema.Set))

	// Get Segment
	segApi := alkira.NewSegment(m.(*alkira.AlkiraClient))
	segment, err := segApi.GetById(strconv.Itoa(d.Get("segment_id").(int)))

	if err != nil {
		log.Printf("[ERROR] failed to get segment by Id: %d", d.Get("segment_id"))
		return nil, err
	}

	// Create implict akamai-prolexic credential
	c := alkira.CredentialAkamaiProlexic{
		BgpAuthenticationKey: d.Get("akamai_bgp_authentication_key").(string),
	}

	log.Printf("[INFO] Creating Credential (akamai-prolexic)")
	credentialId, err := client.CreateCredential(d.Get("name").(string), alkira.CredentialTypeAkamaiProlexic, c, 0)

	if err != nil {
		return nil, err
	}

	d.Set("credential_id", credentialId)

	connector := &alkira.ConnectorAkamaiProlexic{
		AkamaiBgpAsn:         d.Get("akamai_bgp_asn").(int),
		BillingTags:          billingTags,
		ByoipOptions:         byoipOptions,
		CXP:                  d.Get("cxp").(string),
		CredentialId:         credentialId,
		Group:                d.Get("group").(string),
		Enabled:              d.Get("enabled").(bool),
		Name:                 d.Get("name").(string),
		Segments:             []string{segment.Name},
		Size:                 d.Get("size").(string),
		OverlayConfiguration: tunnelConfigurations,
	}

	return connector, nil
}
