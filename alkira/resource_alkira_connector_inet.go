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
			"cxp": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"group": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"segment": {
				Type: schema.TypeString,
				Required: true,
				Description: "A segment associated with the connector AWS-VPC",
			},
		},
	}
}

func resourceConnectorInetCreate(d *schema.ResourceData, m interface{}) error {
	client := m.(*alkira.AlkiraClient)

	segments := []string{d.Get("segment").(string)}

	connectorInet := &alkira.ConnectorInetRequest{
		CXP:            d.Get("cxp").(string),
		Group:          d.Get("group").(string),
		Name:           d.Get("name").(string),
        Segments:       segments,
	}

	id, err := client.CreateConnectorInet(connectorInet)

	if err != nil {
		return err
	}

	d.SetId(strconv.Itoa(id))
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
	err := client.DeleteConnectorInet(d.Id())

	if err != nil {
		return err
	}

	return nil
}
