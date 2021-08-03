package alkira

import (
	"log"

	"github.com/alkiranet/alkira-client-go/alkira"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func resourceAlkiraConnectorInternetExit() *schema.Resource {
	return &schema.Resource{
		Description: "Manage Internet Exit.\n\n" +
			"An internet exit is an exit from the CXP to the" +
			"internet and allows the traffic from" +
			"the various Users & Sites or Cloud Connectors to" +
			"flow towards the Internet.",
		Create: resourceConnectorInternetExitCreate,
		Read:   resourceConnectorInternetExitRead,
		Update: resourceConnectorInternetExitUpdate,
		Delete: resourceConnectorInternetExitDelete,

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

func resourceConnectorInternetExitCreate(d *schema.ResourceData, m interface{}) error {
	client := m.(*alkira.AlkiraClient)

	connector, err := generateConnectorInternetRequest(d, m)

	if err != nil {
		return err
	}

	id, err := client.CreateConnectorInternetExit(connector)

	if err != nil {
		return err
	}

	d.SetId(id)
	return resourceConnectorInternetExitRead(d, m)
}

func resourceConnectorInternetExitRead(d *schema.ResourceData, m interface{}) error {
	return nil
}

func resourceConnectorInternetExitUpdate(d *schema.ResourceData, m interface{}) error {
	client := m.(*alkira.AlkiraClient)

	connector, err := generateConnectorInternetRequest(d, m)

	if err != nil {
		return err
	}

	log.Printf("[INFO] Updateing Connector (INTERNET) %s", d.Id())
	err = client.UpdateConnectorInternetExit(d.Id(), connector)

	return err
}

func resourceConnectorInternetExitDelete(d *schema.ResourceData, m interface{}) error {
	client := m.(*alkira.AlkiraClient)

	log.Printf("[INFO] Deleting Connector (internet-exit) %s", d.Id())
	return client.DeleteConnectorInternetExit(d.Id())
}

// generateConnectorInternetRequest generate request for connector-internet
func generateConnectorInternetRequest(d *schema.ResourceData, m interface{}) (*alkira.ConnectorInternet, error) {
	billingTags := convertTypeListToIntList(d.Get("billing_tag_ids").([]interface{}))
	segments := []string{d.Get("segment").(string)}

	request := &alkira.ConnectorInternet{
		BillingTags: billingTags,
		CXP:         d.Get("cxp").(string),
		Description: d.Get("description").(string),
		Group:       d.Get("group").(string),
		Name:        d.Get("name").(string),
		Segments:    segments,
		Size:        d.Get("size").(string),
	}

	return request, nil
}
