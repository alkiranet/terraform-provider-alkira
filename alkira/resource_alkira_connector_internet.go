package alkira

import (
	"log"
	"strconv"

	"github.com/alkiranet/alkira-client-go/alkira"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func resourceAlkiraConnectorInternet() *schema.Resource {
	return &schema.Resource{
		Description: "Manage Internet Exit.\n\n" +
			"An internet exit is an exit from the CXP to the" +
			"internet and allows and allows the traffic from" +
			"the various Users & Sites or Cloud Connectors to" +
			"flow towards the Internet.",
		Create: resourceConnectorInternetCreate,
		Read:   resourceConnectorInternetRead,
		Update: resourceConnectorInternetUpdate,
		Delete: resourceConnectorInternetDelete,

		Schema: map[string]*schema.Schema{
			"billing_tag_ids": {
				Description: "The list of billing tag Ids.",
				Type:        schema.TypeList,
				Optional:    true,
				Elem:        &schema.Schema{Type: schema.TypeInt},
			},
			"cxp": {
				Description: "The CXP where the connector should be provisioned.",
				Type:        schema.TypeString,
				Required:    true,
			},
			"description": {
				Description: "The description of the connector.",
				Type:        schema.TypeString,
				Optional:    true,
			},
			"group": {
				Description: "The group of the connector.",
				Type:        schema.TypeString,
				Optional:    true,
			},
			"name": {
				Description: "The name of the connector.",
				Type:        schema.TypeString,
				Required:    true,
			},
			"segment": {
				Description: "The segment of the connector belongs to. Currently, only `1` segment is allowed.",
				Type:        schema.TypeString,
				Required:    true,
			},
			"size": {
				Description:  "The size of the connector, one of `SMALL`, `MEDIUM`, or `LARGE`.",
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringInSlice([]string{"SMALL", "MEDIUM", "LARGE"}, false),
			},
		},
	}
}

func resourceConnectorInternetCreate(d *schema.ResourceData, m interface{}) error {
	client := m.(*alkira.AlkiraClient)

	billingTags := convertTypeListToIntList(d.Get("billing_tags").([]interface{}))
	segments := []string{d.Get("segment").(string)}

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
	return resourceConnectorInternetRead(d, m)
}

func resourceConnectorInternetRead(d *schema.ResourceData, m interface{}) error {
	return nil
}

func resourceConnectorInternetUpdate(d *schema.ResourceData, m interface{}) error {
	return resourceConnectorInternetRead(d, m)
}

func resourceConnectorInternetDelete(d *schema.ResourceData, m interface{}) error {
	client := m.(*alkira.AlkiraClient)

	log.Printf("[INFO] Deleting Connector (INET) %s", d.Id())
	err := client.DeleteConnectorInet(d.Get("connector_id").(int))

	if err != nil {
		return err
	}

	return nil
}
