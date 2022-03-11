package alkira

import (
	"log"
	"strconv"

	"github.com/alkiranet/alkira-client-go/alkira"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func resourceAlkiraInternetApplication() *schema.Resource {
	return &schema.Resource{
		Description: "Manage Internet Application.\n\n" +
			"The internet facing applications could be used with both" +
			"Users & Sites or Cloud Connectors.",
		Create: resourceInternetApplicationCreate,
		Read:   resourceInternetApplicationRead,
		Update: resourceInternetApplicationUpdate,
		Delete: resourceInternetApplicationDelete,

		Schema: map[string]*schema.Schema{
			"billing_tag_ids": {
				Description: "IDs of billing tags.",
				Type:        schema.TypeList,
				Optional:    true,
				Elem:        &schema.Schema{Type: schema.TypeInt},
			},
			"connector_id": {
				Description: "Connector ID.",
				Type:        schema.TypeInt,
				Required:    true,
			},
			"connector_type": {
				Description: "Connector Type.",
				Type:        schema.TypeString,
				Required:    true,
			},
			"fqdn_prefix": {
				Description: "User provided FQDN prefix that will be published on route53.",
				Type:        schema.TypeString,
				Required:    true,
			},
			"group_id": {
				Description: "ID of the auto generated system group.",
				Type:        schema.TypeInt,
				Computed:    true,
			},
			"name": {
				Description: "The name of the internet application.",
				Type:        schema.TypeString,
				Required:    true,
			},
			"segment_id": {
				Description: "The ID of segment associated with the internet application.",
				Type:        schema.TypeInt,
				Required:    true,
			},
			"size": {
				Description:  "The size of the internet application, one of `SMALL`, `MEDIUM` and `LARGE`.",
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringInSlice([]string{"SMALL", "MEDIUM", "LARGE"}, false),
			},
		},
	}
}

func resourceInternetApplicationCreate(d *schema.ResourceData, m interface{}) error {
	client := m.(*alkira.AlkiraClient)

	request, err := generateInternetApplicationRequest(d, m)

	if err != nil {
		return err
	}

	id, groupId, err := client.CreateInternetApplication(request)

	if err != nil {
		return err
	}

	d.SetId(id)
	d.Set("group_id", groupId)
	return resourceInternetApplicationRead(d, m)
}

func resourceInternetApplicationRead(d *schema.ResourceData, m interface{}) error {
	client := m.(*alkira.AlkiraClient)

	app, err := client.GetInternetApplication(d.Id())

	if err != nil {
		return err
	}

	d.Set("billing_tag_ids", app.BillingTags)
	d.Set("connector_id", app.ConnectorId)
	d.Set("connector_type", app.ConnectorType)
	d.Set("fqdn_prefix", app.FqdnPrefix)
	d.Set("Name", app.Name)

	segment, err := client.GetSegmentByName(app.SegmentName)

	if err != nil {
		return err
	}
	d.Set("segment_id", segment.Id)
	d.Set("size", app.Size)

	return nil
}

func resourceInternetApplicationUpdate(d *schema.ResourceData, m interface{}) error {
	client := m.(*alkira.AlkiraClient)

	request, err := generateInternetApplicationRequest(d, m)

	if err != nil {
		return err
	}

	err = client.UpdateInternetApplication(d.Id(), request)

	if err != nil {
		return err
	}

	return resourceInternetApplicationRead(d, m)
}

func resourceInternetApplicationDelete(d *schema.ResourceData, m interface{}) error {
	client := m.(*alkira.AlkiraClient)

	return client.DeleteInternetApplication(d.Id())
}

func generateInternetApplicationRequest(d *schema.ResourceData, m interface{}) (*alkira.InternetApplication, error) {
	client := m.(*alkira.AlkiraClient)

	billingTags := convertTypeListToIntList(d.Get("billing_tag_ids").([]interface{}))

	segment, err := client.GetSegmentById(strconv.Itoa(d.Get("segment_id").(int)))

	if err != nil {
		log.Printf("[ERROR] failed to get segment by Id: %d", d.Get("segment_id"))
		return nil, err
	}

	request := &alkira.InternetApplication{
		BillingTags:   billingTags,
		ConnectorId:   d.Get("connector_id").(int),
		ConnectorType: d.Get("connector_type").(string),
		FqdnPrefix:    d.Get("fqdn_prefix").(string),
		Name:          d.Get("name").(string),
		SegmentName:   segment.Name,
		Size:          d.Get("size").(string),
	}

	return request, nil
}
