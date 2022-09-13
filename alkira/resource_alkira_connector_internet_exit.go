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
				Description: "The list of billing tag IDs.",
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
			"enabled": {
				Description: "Is the connector enabled. Default is `true`.",
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     true,
			},
			"group": {
				Description: "The group of the connector.",
				Type:        schema.TypeString,
				Optional:    true,
			},
			"implicit_group_id": {
				Description: "The ID of implicit group automaticaly created with the connector.",
				Type:        schema.TypeInt,
				Computed:    true,
			},
			"name": {
				Description: "The name of the connector.",
				Type:        schema.TypeString,
				Required:    true,
			},
			"public_ip_number": {
				Description: "The number of the public IPs to the connector. Default is `2`.",
				Type:        schema.TypeInt,
				Optional:    true,
				Default:     2,
			},
			"segment_id": {
				Description: "ID of segment associated with the connector.",
				Type:        schema.TypeString,
				Required:    true,
			},
			"size": {
				Description:  "The size of the connector, one of `SMALL`, `MEDIUM`, or `LARGE`.",
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringInSlice([]string{"SMALL", "MEDIUM", "LARGE"}, false),
			},
			"traffic_distribution_algorithm": {
				Description: "The type of the algorithm to be used for traffic distribution." +
					"Currently, only `HASHING` is supported.",
				Type:         schema.TypeString,
				Optional:     true,
				Default:      "HASHING",
				ValidateFunc: validation.StringInSlice([]string{"HASHING"}, false),
			},
			"traffic_distribution_algorithm_attribute": {
				Description:  "The attributes depends on the algorithm. For now, it's either `DEFAULT` or `SRC_IP`.",
				Type:         schema.TypeString,
				Optional:     true,
				Default:      "DEFAULT",
				ValidateFunc: validation.StringInSlice([]string{"DEFAULT", "SRC_IP"}, false),
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
	client := m.(*alkira.AlkiraClient)

	connector, err := client.GetConnectorInternetExitById(d.Id())

	if err != nil {
		return err
	}

	d.Set("billing_tag_ids", connector.BillingTags)
	d.Set("cxp", connector.CXP)
	d.Set("description", connector.Description)
	d.Set("enabled", connector.Enabled)
	d.Set("group", connector.Group)
	d.Set("implicit_group_id", connector.ImplicitGroupId)
	d.Set("name", connector.Name)
	d.Set("public_ip_number", connector.NumOfPublicIPs)

	// Set segment_id
	if len(connector.Segments) > 0 {
		segment, err := client.GetSegmentByName(connector.Segments[0])

		if err != nil {
			return err
		}
		d.Set("segment_id", segment.Id)
	}

	d.Set("size", connector.Size)

	if connector.TrafficDistribution != nil {
		d.Set("traffic_distribution_algorithm", connector.TrafficDistribution.Algorithm)
		d.Set("traffic_distribution_algorithm_attribute", connector.TrafficDistribution.AlgorithmAttributes.Keys)
	}

	return nil
}

func resourceConnectorInternetExitUpdate(d *schema.ResourceData, m interface{}) error {
	client := m.(*alkira.AlkiraClient)

	connector, err := generateConnectorInternetRequest(d, m)

	if err != nil {
		return err
	}

	log.Printf("[INFO] Updating Connector (INTERNET-EXIT) %s", d.Id())
	err = client.UpdateConnectorInternetExit(d.Id(), connector)

	return err
}

func resourceConnectorInternetExitDelete(d *schema.ResourceData, m interface{}) error {
	client := m.(*alkira.AlkiraClient)

	log.Printf("[INFO] Deleting Connector (INTERNET-EXIT) %s", d.Id())
	return client.DeleteConnectorInternetExit(d.Id())
}

// generateConnectorInternetRequest generate request for connector-internet
func generateConnectorInternetRequest(d *schema.ResourceData, m interface{}) (*alkira.ConnectorInternet, error) {
	client := m.(*alkira.AlkiraClient)

	billingTags := convertTypeListToIntList(d.Get("billing_tag_ids").([]interface{}))
	segment, err := client.GetSegmentById(d.Get("segment_id").(string))

	algorithmAttributes := alkira.AlgorithmAttributes{
		Keys: d.Get("traffic_distribution_algorithm_attribute").(string),
	}

	trafficDistribution := alkira.TrafficDistribution{
		Algorithm:           d.Get("traffic_distribution_algorithm").(string),
		AlgorithmAttributes: algorithmAttributes,
	}

	if err != nil {
		return nil, err
	}

	request := &alkira.ConnectorInternet{
		BillingTags:         billingTags,
		CXP:                 d.Get("cxp").(string),
		Description:         d.Get("description").(string),
		Group:               d.Get("group").(string),
		Enabled:             d.Get("enabled").(bool),
		Name:                d.Get("name").(string),
		NumOfPublicIPs:      d.Get("public_ip_number").(int),
		Segments:            []string{segment.Name},
		Size:                d.Get("size").(string),
		TrafficDistribution: &trafficDistribution,
	}

	return request, nil
}
