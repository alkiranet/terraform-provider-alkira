package alkira

import (
	"log"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/alkiranet/alkira-client-go/alkira"
)

func resourceAlkiraInternetApplication() *schema.Resource {
	return &schema.Resource{
		Create: resourceInternetApplicationCreate,
		Read:   resourceInternetApplicationRead,
		Update: resourceInternetApplicationUpdate,
		Delete: resourceInternetApplicationDelete,

		Schema: map[string]*schema.Schema{
			"connector_id": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				Description: "ID of the connector",
			},
			"connector_type": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				Description: "Type of the connector",
			},
			"fqdn_prefix": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				Description: "FQDN Prefix",
			},
			"group": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Description: "Group",
			},
			"name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				Description: "Name of the internet application",
			},
			"private_ip": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				Description: "Private IP",
			},
			"private_port": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				Description: "Private Port",
			},
			"segment": {
				Type: schema.TypeString,
				Required: true,
				Description: "Name of the segment",
			},
			"size": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				Description: "Size",
			},
		},
	}
}

func resourceInternetApplicationCreate(d *schema.ResourceData, m interface{}) error {
	client := m.(*alkira.AlkiraClient)

	connector := &alkira.InternetApplicationRequest{
		ConnectorId:    d.Get("connector_id").(string),
		ConnectorType:  d.Get("connector_type").(string),
		FqdnPrefix:     d.Get("fqdn_prefix").(string),
		Group:          d.Get("group").(string),
		Name:           d.Get("name").(string),
        PrivateIp:      d.Get("private_ip").(string),
        PrivatePort:    d.Get("private_port").(string),
        SegmentName:    d.Get("segment").(string),
        Size:           d.Get("size").(string),
	}

	id, err := client.CreateInternetApplication(connector)

	if err != nil {
		return err
	}

	d.SetId(strconv.Itoa(id))
	return resourceInternetApplicationRead(d, m)
}

func resourceInternetApplicationRead(d *schema.ResourceData, m interface{}) error {
        return nil
}

func resourceInternetApplicationUpdate(d *schema.ResourceData, m interface{}) error {
        return resourceInternetApplicationRead(d, m)
}

func resourceInternetApplicationDelete(d *schema.ResourceData, m interface{}) error {
	client := m.(*alkira.AlkiraClient)

	log.Printf("[INFO] Deleting Internet Application %s", d.Id())
	err := client.DeleteInternetApplication(d.Id())

	if err != nil {
		return err
	}

	return nil
}
