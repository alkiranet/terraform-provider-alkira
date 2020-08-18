package alkira

import (
	"log"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/alkiranet/alkira-client-go/alkira"
)

func resourceAlkiraConnectorInet() *schema.Resource {
	return &schema.Resource{
		Create: resourceConnectorInetCreate,
		Read:   resourceConnectorInetRead,
		Update: resourceConnectorInetUpdate,
		Delete: resourceConnectorInetDelete,

		Schema: map[string]*schema.Schema{
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
			"segments": {
				Type:     schema.TypeList,
				Required: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"size": {
				Type: schema.TypeString,
				Required: true,
			},
		},
	}
}

func resourceConnectorInetCreate(d *schema.ResourceData, m interface{}) error {
	client := m.(*alkira.AlkiraClient)

	segments := convertTypeListToStringList(d.Get("segments").([]interface{}))

	connector := &alkira.ConnectorInternet{
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
