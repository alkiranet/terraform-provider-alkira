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
			"connector_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"connector_type": {
				Type:     schema.TypeString,
				Required: true,
			},
			"fqdn_prefix": {
				Type:     schema.TypeString,
				Required: true,
			},
			"group": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"internet_application_id": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"private_ip": {
				Type:     schema.TypeString,
				Required: true,
			},
			"private_port": {
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
	d.Set("internet_application_id", id)

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
	err := client.DeleteInternetApplication(d.Get("internet_application_id").(int))

	if err != nil {
		return err
	}

	return nil
}
