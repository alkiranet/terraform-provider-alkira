package alkira

import (
	"log"
	"strconv"

	"github.com/alkiranet/alkira-client-go/alkira"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceAlkiraConnectorInet() *schema.Resource {
	return &schema.Resource{
		Create: resourceConnectorInetCreate,
		Read:   resourceConnectorInetRead,
		Update: resourceConnectorInetUpdate,
		Delete: resourceConnectorInetDelete,

		Schema: map[string]*schema.Schema{
			"billing_tags": {
				Type:     schema.TypeList,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"connector_id": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"cxp": {
				Type:     schema.TypeString,
				Required: true,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"group": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"segment": {
				Type:     schema.TypeString,
				Required: true,
			},
			"size": {
				Type:     schema.TypeString,
				Required: true,
			},
		},
	}
}

func resourceConnectorInetCreate(d *schema.ResourceData, m interface{}) error {
	client := m.(*alkira.AlkiraClient)

	billingTags := convertTypeListToStringList(d.Get("billing_tags").([]interface{}))
	segments    := []string{d.Get("segment").(string)}

	connector := &alkira.ConnectorInternet{
		BillingTags: billingTags,
		CXP:         d.Get("cxp").(string),
		Description: d.Get("description").(string),
		Group:       d.Get("group").(string),
		Name:        d.Get("name").(string),
		Segments:    segments,
		Size:        d.Get("size").(string),
	}

	id, err := client.CreateConnectorInternet(connector)

	if err != nil {
		return err
	}

	d.SetId(strconv.Itoa(id))
	d.Set("connector_id", id)
	return resourceConnectorInetRead(d, m)
}

func resourceConnectorInetRead(d *schema.ResourceData, m interface{}) error {
	return nil
}

func resourceConnectorInetUpdate(d *schema.ResourceData, m interface{}) error {
	return resourceConnectorInetRead(d, m)
}

func resourceConnectorInetDelete(d *schema.ResourceData, m interface{}) error {
	client := m.(*alkira.AlkiraClient)

	log.Printf("[INFO] Deleting Connector (INET) %s", d.Id())
	err := client.DeleteConnectorInet(d.Get("connector_id").(int))

	if err != nil {
		return err
	}

	return nil
}
