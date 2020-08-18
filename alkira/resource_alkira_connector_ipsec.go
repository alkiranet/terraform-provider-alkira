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
			"connector_id": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"cxp": {
				Type:     schema.TypeString,
				Required: true,
			},
			"group": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"name": {
				Type:     schema.TypeString,
				Required: true,
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
				Type:     schema.TypeList,
				Required: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
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

	segments := convertTypeListToStringList(d.Get("segments").([]interface{}))
	sites    := expandIPSecSites(d.Get("sites").(*schema.Set))

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
	d.Set("connector_id", id)

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
	err := client.DeleteConnectorIPSec(d.Get("connector_id").(int))

	if err != nil {
		return err
	}

	return nil
}
