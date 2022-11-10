package alkira

import (
	"fmt"

	"github.com/alkiranet/alkira-client-go/alkira"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func resourceAlkiraConnectorCiscoFTDv() *schema.Resource {
	return &schema.Resource{
		Description: "Manage Cisco FTDv Connector. (**BETA**)",

		Create: resourceConnectorCiscoFTDvCreate,
		Read:   resourceConnectorCiscoFTDvRead,
		Update: resourceConnectorCiscoFTDvUpdate,
		Delete: resourceConnectorCiscoFTDvDelete,

		Schema: map[string]*schema.Schema{
			"name": {
				Description: "The name of the connector.",
				Type:        schema.TypeString,
				Required:    true,
			},
			"auto_scale": {
				Description: "Indicate if `auto_scale` should be enabled for your Cisco FTDv connector." +
					" `ON` and `OFF` are accepted values. `OFF` is the default if " +
					"field is omitted.",
				Type:         schema.TypeString,
				Default:      "OFF",
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"ON", "OFF"}, false),
			},
			"size": {
				Description:  "The size of the connector, one of `SMALL`, `MEDIUM`, `LARGE`.",
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringInSlice([]string{"SMALL", `MEDIUM`, `LARGE`}, false),
			},
			"tunnel_protocol": {
				Description:  "The tunnel protocol. One of `IPSEC`. Default is `IPSEC`",
				Type:         schema.TypeString,
				Optional:     true,
				Default:      "IPSEC",
				ValidateFunc: validation.StringInSlice([]string{"IPSEC"}, false),
			},
			"cxp": {
				Description: "The CXP where the connector should be provisioned.",
				Type:        schema.TypeString,
				Required:    true,
			},
			"global_cidr_list_id": {
				Description: "The ID of the global cidr list to be associated with " +
					"the management server.",
				Type:     schema.TypeInt,
				Required: true,
			},
			"max_instance_count": {
				Description: "The maximum number of Cisco FTDv instances that should be deployed.",
				Type:        schema.TypeInt,
				Required:    true,
			},
			"min_instance_count": {
				Description: "The minimum number of Cisco FTDv instances that should be deployed.",
				Type:        schema.TypeInt,
				Optional:    true,
				Default:     0,
			},
			"ip_allow_list": {
				Description: "",
				Type:        schema.TypeList,
				Optional:    true,
			},
			"group": {
				Description: "The group of the connector.",
				Type:        schema.TypeString,
				Optional:    true,
			},
			"billing_tag_ids": {
				Description: "IDs of Billing Tags.",
				Type:        schema.TypeList,
				Optional:    true,
				Elem:        &schema.Schema{Type: schema.TypeInt},
			},
			"segment_id": {
				Description: "The ID of the segments associated with the service.",
				Type:        schema.TypeInt,
				Required:    true,
			},
			"management_server": {
				Type:        schema.TypeSet,
				Required:    true,
				Description: "",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"ip_address": {
							Description: "",
							Type:        schema.TypeString,
							Computed:    true,
						},
						"segment_id": {
							Description: "The ID of the segment to be used to access the management server.",
							Type:        schema.TypeInt,
							Optional:    true,
						},
					},
				},
			},
			"instances": {
				Type:        schema.TypeSet,
				Required:    true,
				Description: "",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"credential_id": {
							Description: "",
							Type:        schema.TypeString,
							Computed:    true,
						},
						"hostname": {
							Description: "",
							Type:        schema.TypeString,
							Required:    true,
						},
						"version": {
							Description: "",
							Type:        schema.TypeString,
							Required:    true,
						},
						"license_type": {
							Description:  ", either `BRING_YOUR_OWN` or `PAY_AS_YOU_GO`.",
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: validation.StringInSlice([]string{"BRING_YOUR_OWN", "PAY_AS_YOU_GO"}, false),
						},
					},
				},
			},
		},
	}
}

// resourceConnectorCiscoFTDvCreate create a Cisco FTDv connector
func resourceConnectorCiscoFTDvCreate(d *schema.ResourceData, m interface{}) error {
	client := m.(*alkira.AlkiraClient)
	connector, err := generateConnectorCiscoFTDvRequest(d, m)

	if err != nil {
		return err
	}

	id, err := client.CreateConnectorCiscoFTDv(connector)

	if err != nil {
		return err
	}

	d.SetId(id)
	return resourceConnectorCiscoFTDvRead(d, m)
}

// resourceConnectorCiscoFTDvRead get and save a Cisco FTDv connectors
func resourceConnectorCiscoFTDvRead(d *schema.ResourceData, m interface{}) error {
	client := m.(*alkira.AlkiraClient)

	connector, err := client.GetConnectorCiscoFTDv(d.Id())

	if err != nil {
		return err
	}

	d.Set("size", connector.Size)
	d.Set("billing_tag_ids", connector.BillingTags)
	d.Set("cxp", connector.Cxp)
	d.Set("group", connector.Group)
	d.Set("enabled", connector.Enabled)
	d.Set("name", connector.Name)
	d.Set("tunnel_protocol", connector.TunnelProtocol)
	d.Set("vhub_prefix", connector.VhubPrefix)

	var instances []map[string]interface{}
	for _, instance := range connector.Instances {
		i := map[string]interface{}{}
		instances = append(instances, i)
	}

	d.Set("instances", instances)

	var segments []map[string]interface{}

	return nil
}

// resourceConnectorCiscoFTDvUpdate update a Cisco FTDv connector
func resourceConnectorCiscoFTDvUpdate(d *schema.ResourceData, m interface{}) error {
	client := m.(*alkira.AlkiraClient)
	connector, err := generateConnectorCiscoFTDvRequest(d, m)

	if err != nil {
		return fmt.Errorf("UpdateConnectorCiscoFTDv: failed to marshal: %v", err)
	}

	err = client.UpdateConnectorCiscoFTDv(d.Id(), connector)

	if err != nil {
		return err
	}

	return resourceConnectorCiscoFTDvRead(d, m)
}

func resourceConnectorCiscoFTDvDelete(d *schema.ResourceData, m interface{}) error {
	client := m.(*alkira.AlkiraClient)

	err := client.DeleteConnectorCiscoFTDv((d.Id()))

	return err
}

// generateConnectorCiscoFTDvRequest generate a request for Azure ExpressRoute connector
func generateConnectorCiscoFTDvRequest(d *schema.ResourceData, m interface{}) (*alkira.ConnectorCiscoFTDv, error) {

	billingTags := convertTypeListToIntList(d.Get("billing_tag_ids").([]interface{}))

	instances, err := expandCiscoFTDvInstances(d.Get("instances").([]interface{}), m)
	if err != nil {
		return nil, err
	}

	request := &alkira.ConnectorCiscoFTDv{
		Name:             d.Get("name").(string),
		AutoScale:        d.Get("auto_scale").(string),
		GlobalCidrListId: d.Get("global_cidr_list").(int),
		Size:             d.Get("size").(string),
		Cxp:              d.Get("cxp").(string),
		MaxInstanceCount: d.Get("max_instance_count").(int),
		MinInstanceCount: d.Get("min_instance_count").(int),
		TunnelProtocol:   d.Get("tunnel_protocol").(string),
		BillingTags:      billingTags,
		// Instances:      instances,
		// SegmentOptions: segmentOptions,
	}

	return request, nil
}
