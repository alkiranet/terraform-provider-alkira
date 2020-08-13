package alkira

import (
	"log"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/alkiranet/alkira-client-go/alkira"
)

func resourceAlkiraConnectorIPSec() *schema.Resource {
	return &schema.Resource{
		Create: resourceConnectorIPSecCreate,
		Read:   resourceConnectorIPSecRead,
		Update: resourceConnectorIPSecUpdate,
		Delete: resourceConnectorIPSecDelete,

		Schema: map[string]*schema.Schema{
			"cxp": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				Description: "The CXP to be used for the connector",
			},
			"group": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Description: "A user group that the connector belongs to",
			},
			"name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				Description: "The name of the connector",
			},
			"segment_options": {
				Type: schema.TypeSet,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:     schema.TypeString,
							Required: true,
						},
						"disable_internet_exit": {
							Type:      schema.TypeString,
							Optional:  true,
						},

						"disable_advertise_on_prem_routes": {
							Type:     schema.TypeString,
							Optional: true,
						},
					},
				},
				Optional: true,
			},
			"segments": {
				Type: schema.TypeString,
				Required: true,
				Description: "A segment associated with the connector",
			},
			"sites": &schema.Schema{
				Type: schema.TypeSet,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type: schema.TypeString,
							Required: true,
						},
						"customer_gateway_asn": {
							Type: schema.TypeString,
							Optional:  true,
						},
						"customer_gateway_ip": {
							Type: schema.TypeString,
							Optional: true,
						},
						"preshared_keys": {
							Type: schema.TypeList,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
							Optional: true,
						},
					},
				},
				Optional: true,
			},
			"size": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
		},
	}
}

func resourceConnectorIPSecCreate(d *schema.ResourceData, m interface{}) error {
	client := m.(*alkira.AlkiraClient)

	segments := []string{d.Get("segments").(string)}
	sites := expandIPSecSites(d.Get("sites").(*schema.Set))

	connector := &alkira.ConnectorIPSecRequest{
		CXP:            d.Get("cxp").(string),
		Group:          d.Get("group").(string),
		Name:           d.Get("name").(string),
//		SegmentOptions: d.Get("segment_options").(*schema.Set).List(),
        Segments:       segments,
		Sites:          sites,
        Size:           d.Get("size").(string),
	}

	log.Printf("[INFO] Creating Connector (IPSec) %s", d.Id())
	id, err := client.CreateConnectorIPSec(connector)

	if err != nil {
		return err
	}

	d.SetId(strconv.Itoa(id))
	return resourceConnectorIPSecRead(d, m)
}

func resourceConnectorIPSecRead(d *schema.ResourceData, m interface{}) error {
        return nil
}

func resourceConnectorIPSecUpdate(d *schema.ResourceData, m interface{}) error {
        return resourceConnectorIPSecRead(d, m)
}

func resourceConnectorIPSecDelete(d *schema.ResourceData, m interface{}) error {
	client := m.(*alkira.AlkiraClient)

	log.Printf("[INFO] Deleting Connector (IPSec) %s", d.Id())
	err := client.DeleteConnectorIPSec(d.Id())

	if err != nil {
		return err
	}

	return nil
}
