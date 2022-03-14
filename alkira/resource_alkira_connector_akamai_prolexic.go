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
		Description: "Manage Connector for Akamai Prolexic.",
		Create:      resourceConnectorAkamaiProlexicCreate,
		Read:        resourceConnectorAkamaiProlexicRead,
		Update:      resourceConnectorAkamaiProlexicUpdate,
		Delete:      resourceConnectorAkamaiProlexicDelete,

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
			"byoip_prefix": &schema.Schema{
				Description: "BYOIP prefixe.",
				Type:        schema.TypeSet,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"byoip_prefix_id": {
							Description: "BYOIP prefix ID.",
							Type:        schema.TypeString,
							Required:    true,
						},
						"enable_route_advertisement": {
							Description: "Whether enabling route advertisement.",
							Type:        schema.TypeBool,
							Required:    true,
						},
					},
				},
				Required: true,
			},
			"billing_tag_ids": {
				Description: "A list of Billing Tag by ID associated with the connector.",
				Type:        schema.TypeList,
				Optional:    true,
				Elem:        &schema.Schema{Type: schema.TypeInt},
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
			"size": &schema.Schema{
				Description:  "The size of the connector. one of `SMALL`, `MEDIUM` and `LARGE`.",
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringInSlice([]string{"SMALL", "MEDIUM", "LARGE"}, false),
			},
			"tunnel_configuration": &schema.Schema{
				Description: "Tunnel Configurations.",
				Type:        schema.TypeSet,
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
				Required: true,
			},
		},
	}
}

func resourceConnectorAkamaiProlexicCreate(d *schema.ResourceData, m interface{}) error {
	client := m.(*alkira.AlkiraClient)
	connector, err := generateConnectorAkamaiProlexicRequest(client, d, m)

	if err != nil {
		return err
	}

	log.Printf("[INFO] Creating Connector (Akamai-Prolexic)")
	id, err := client.CreateConnectorAkamaiProlexic(connector)

	if err != nil {
		return err
	}

	d.SetId(id)

	return resourceConnectorAkamaiProlexicRead(d, m)
}

func resourceConnectorAkamaiProlexicRead(d *schema.ResourceData, m interface{}) error {
	client := m.(*alkira.AlkiraClient)

	connector, err := client.GetConnectorAkamaiProlexic(d.Id())

	if err != nil {
		return err
	}

	d.Set("billing_tag_ids", connector.BillingTags)
	d.Set("cxp", connector.CXP)
	d.Set("group", connector.Group)
	d.Set("enabled", connector.Enabled)
	d.Set("name", connector.Name)
	d.Set("size", connector.Size)

	if len(connector.Segments) > 0 {
		segment, err := client.GetSegmentByName(connector.Segments[0])

		if err != nil {
			return err
		}
		d.Set("segment_id", segment.Id)
	}

	return nil
}

func resourceConnectorAkamaiProlexicUpdate(d *schema.ResourceData, m interface{}) error {
	client := m.(*alkira.AlkiraClient)

	connector, err := generateConnectorAkamaiProlexicRequest(client, d, m)

	if err != nil {
		return err
	}

	log.Printf("[INFO] Updating Connector (Akamai-Prolexic): %s", d.Id())
	err = client.UpdateConnectorAkamaiProlexic(d.Id(), connector)

	if err != nil {
		return err
	}

	return resourceConnectorAkamaiProlexicRead(d, m)
}

func resourceConnectorAkamaiProlexicDelete(d *schema.ResourceData, m interface{}) error {
	client := m.(*alkira.AlkiraClient)
	id := d.Id()

	log.Printf("[INFO] Deleting Connector (Akamai-Prolexic): %s", id)
	return client.DeleteConnectorAkamaiProlexic(id)
}

// generateConnectorAkamaiProlexicRequest generate request for the connector
func generateConnectorAkamaiProlexicRequest(ac *alkira.AlkiraClient, d *schema.ResourceData, m interface{}) (*alkira.ConnectorAkamaiProlexic, error) {
	client := m.(*alkira.AlkiraClient)

	billingTags := convertTypeListToIntList(d.Get("billing_tag_ids").([]interface{}))

	segment, err := client.GetSegmentById(strconv.Itoa(d.Get("segment_id").(int)))

	if err != nil {
		log.Printf("[ERROR] failed to get segment by Id: %d", d.Get("segment_id"))
		return nil, err
	}

	// Create hidden akamai-prolexic credential
	c := alkira.CredentialAkamaiProlexic{
		BgpAuthenticationKey: d.Get("bgp_authentication_key").(string),
	}

	log.Printf("[INFO] Creating Credential (akamai-prolexic)")
	credentialId, err := client.CreateCredential(d.Get("name").(string), alkira.CredentialTypeAkamaiProlexic, c)

	if err != nil {
		return nil, err
	}

	connector := &alkira.ConnectorAkamaiProlexic{
		BillingTags:  billingTags,
		CXP:          d.Get("cxp").(string),
		CredentialId: credentialId,
		Group:        d.Get("group").(string),
		Enabled:      d.Get("enabled").(bool),
		Name:         d.Get("name").(string),
		Segments:     []string{segment.Name},
		Size:         d.Get("size").(string),
	}

	return connector, nil
}
