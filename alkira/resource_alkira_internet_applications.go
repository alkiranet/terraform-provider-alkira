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
			"inbound_connector_id": {
				Description: "Inbound connector ID. When `inbound_connector_type` is `DEFAULT`, " +
					"it could be left empty.",
				Type:     schema.TypeString,
				Optional: true,
			},
			"inbound_connector_type": {
				Description: "The inbound connector type specifies how the internet application " +
					"is to be opened up to the external world. By `DEFAULT` the native cloud " +
					"internet connector is used. In this scenario, Alkira takes care of creating " +
					"this inbound internet connector implicitly. If instead inbound access is via " +
					"the `AKAMAI_PROLEXIC` connector, then you need to create and configure " +
					"that connector and use it with the internet application.",
				Type:         schema.TypeString,
				Optional:     true,
				Default:      "DEFAULT",
				ValidateFunc: validation.StringInSlice([]string{"DEFAULT", "AKAMAI_PROLEXIC"}, false),
			},
			"name": {
				Description: "The name of the internet application.",
				Type:        schema.TypeString,
				Required:    true,
			},
			"public_ips": {
				Description: "This option pertains to the `AKAMAI_PROLEXIC` inbound_connector_type. " +
					"The public IPs are to be used to access the internet application. These public IPs " +
					"must belong to one of the BYOIP ranges configured for the Akamai Prolexic Connector.",
				Type:     schema.TypeList,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
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
			"target": {
				Type: schema.TypeSet,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"type": {
							Description:  "The type of the target, one of `IP` or `ILB_NAME`.",
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: validation.StringInSlice([]string{"IP", "ILB_NAME"}, false),
						},
						"value": {
							Description: "IFA ILB name or private IP.",
							Type:        schema.TypeString,
							Required:    true,
						},
						"ports": {
							Description: "list of internet application ports.",
							Type:        schema.TypeList,
							Elem:        &schema.Schema{Type: schema.TypeInt},
							Required:    true,
						},
					},
				},
				Required: true,
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
	d.Set("name", app.Name)
	d.Set("public_ips", app.PublicIps)
	d.Set("size", app.Size)

	// segment_id
	segment, err := client.GetSegmentByName(app.SegmentName)

	if err != nil {
		return err
	}
	d.Set("segment_id", segment.Id)

	// targets
	var targets []map[string]interface{}

	for _, target := range app.Targets {
		i := map[string]interface{}{
			"type":  target.Type,
			"value": target.Value,
			"ports": target.Ports,
		}
		targets = append(targets, i)
	}

	d.Set("targets", targets)

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
	publicIps := convertTypeListToStringList(d.Get("public_ips").([]interface{}))

	targets := expandInternetApplicationTargets(d.Get("target").(*schema.Set))
	segment, err := client.GetSegmentById(strconv.Itoa(d.Get("segment_id").(int)))

	if err != nil {
		log.Printf("[ERROR] failed to get segment by Id: %d", d.Get("segment_id"))
		return nil, err
	}

	request := &alkira.InternetApplication{
		BillingTags:          billingTags,
		ConnectorId:          d.Get("connector_id").(int),
		ConnectorType:        d.Get("connector_type").(string),
		FqdnPrefix:           d.Get("fqdn_prefix").(string),
		InboundConnectorId:   d.Get("inbound_connector_id").(string),
		InboundConnectorType: d.Get("inbound_connector_type").(string),
		Name:                 d.Get("name").(string),
		PublicIps:            publicIps,
		SegmentName:          segment.Name,
		Size:                 d.Get("size").(string),
		Targets:              targets,
	}

	return request, nil
}
