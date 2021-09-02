package alkira

import (
	"log"

	"github.com/alkiranet/alkira-client-go/alkira"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func resourceAlkiraServicePan() *schema.Resource {
	return &schema.Resource{
		Description: "Manage PAN firewall.",
		Create:      resourceServicePanCreate,
		Read:        resourceServicePanRead,
		Update:      resourceServicePanUpdate,
		Delete:      resourceServicePanDelete,

		Schema: map[string]*schema.Schema{
			"billing_tag_ids": {
				Description: "A list of billing tag ids to associate with the service.",
				Type:        schema.TypeList,
				Optional:    true,
				Elem:        &schema.Schema{Type: schema.TypeInt},
			},
			"credential_id": {
				Description: "ID of PAN credential managed by credential resource.",
				Type:        schema.TypeString,
				Optional:    true,
			},
			"cxp": {
				Description: "The CXP where the service should be provisioned.",
				Type:        schema.TypeString,
				Required:    true,
			},
			"instance": {
				Type: schema.TypeSet,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Description: "The name of the PAN instance.",
							Type:        schema.TypeString,
							Required:    true,
						},
						"credential_id": {
							Description: "ID of PAN instance credential managed by credential resource.",
							Type:        schema.TypeString,
							Required:    true,
						},
					},
				},
				Required: true,
			},
			"license_type": {
				Description:  "PAN license type, either `BRING_YOUR_OWN` or `PAY_AS_YOU_GO`.",
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringInSlice([]string{"BRING_YOUR_OWN", "PAY_AS_YOU_GO"}, false),
			},
			"panorama_enabled": {
				Description: "Enable Panorama or not.",
				Type:        schema.TypeBool,
				Required:    true,
			},
			"panorama_device_group": {
				Description: "Panorama device group.",
				Type:        schema.TypeString,
				Optional:    true,
			},
			"panorama_ip_address": {
				Description: "Panorama IP address.",
				Type:        schema.TypeString,
				Optional:    true,
			},
			"panorama_template": {
				Description: "Panorama Template.",
				Type:        schema.TypeString,
				Optional:    true,
			},
			"management_segment_id": {
				Description: "Management Segment Id.",
				Type:        schema.TypeInt,
				Required:    true,
			},
			"max_instance_count": {
				Description: "Max number of Panorama instances for auto scale.",
				Type:        schema.TypeInt,
				Optional:    true,
				Default:     1,
			},
			"min_instance_count": {
				Description: "Minimal number of Panorama instances for auto scale.",
				Type:        schema.TypeInt,
				Optional:    true,
				Default:     1,
			},
			"name": {
				Description: "Name of the PAN service.",
				Type:        schema.TypeString,
				Required:    true,
			},
			"segment_ids": {
				Description: "The list of segment Ids the service belongs to.",
				Type:        schema.TypeList,
				Required:    true,
				Elem:        &schema.Schema{Type: schema.TypeInt},
			},
			"size": {
				Description:  "The size of the service, one of `SMALL`, `MEDIUM`, `LARGE`.",
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringInSlice([]string{"SMALL", "MEDIUM", "LARGE"}, false),
			},
			"tunnel_protocol": {
				Description:  "Tunnel Protocol, default to `IPSEC`, could be either `IPSEC` or `GRE`.",
				Type:         schema.TypeString,
				Optional:     true,
				Default:      "IPSEC",
				ValidateFunc: validation.StringInSlice([]string{"IPSEC", "GRE"}, false),
			},
			"type": {
				Description: "The type of the PAN firewall.",
				Type:        schema.TypeString,
				Required:    true,
			},
			"version": {
				Description: "The version of the PAN firewall.",
				Type:        schema.TypeString,
				Required:    true,
			},
			"zones_to_groups": {
				Type: schema.TypeSet,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"segment_name": {
							Description: "The name of the segment.",
							Type:        schema.TypeString,
							Required:    true,
						},
						"zone_name": {
							Description: "The name of the zone.",
							Type:        schema.TypeString,
							Required:    true,
						},
						"groups": {
							Description: "The name of the group.",
							Type:        schema.TypeList,
							Required:    true,
							Elem:        &schema.Schema{Type: schema.TypeString},
						},
					},
				},
				Optional: true,
			},
		},
	}
}

func resourceServicePanCreate(d *schema.ResourceData, m interface{}) error {
	client := m.(*alkira.AlkiraClient)

	billingTagIds := convertTypeListToIntList(d.Get("billing_tag_ids").([]interface{}))
	instances := expandPanInstances(d.Get("instance").(*schema.Set))
	segmentIds := convertTypeListToIntList(d.Get("segment_ids").([]interface{}))
	segmentOptions := expandPanSegmentOptions(d.Get("zones_to_groups").(*schema.Set))

	service := &alkira.ServicePan{
		BillingTagIds:       billingTagIds,
		CXP:                 d.Get("cxp").(string),
		CredentialId:        d.Get("credential_id").(string),
		Instances:           instances,
		LicenseType:         d.Get("license_type").(string),
		MaxInstanceCount:    d.Get("max_instance_count").(int),
		MinInstanceCount:    d.Get("min_instance_count").(int),
		ManagementSegmentId: d.Get("management_segment_id").(int),
		Name:                d.Get("name").(string),
		PanoramaEnabled:     d.Get("panorama_enabled").(bool),
		PanoramaDeviceGroup: d.Get("panorama_device_group").(string),
		PanoramaIpAddress:   d.Get("panorama_ip_address").(string),
		PanoramaTemplate:    d.Get("panorama_template").(string),
		SegmentOptions:      segmentOptions,
		SegmentIds:          segmentIds,
		Size:                d.Get("size").(string),
		Type:                d.Get("type").(string),
		Version:             d.Get("version").(string),
	}

	log.Printf("[INFO] Creating Service (PAN) %s", d.Id())
	id, err := client.CreateServicePan(service)

	if err != nil {
		return err
	}

	d.SetId(id)
	return resourceServicePanRead(d, m)
}

func resourceServicePanRead(d *schema.ResourceData, m interface{}) error {
	return nil
}

func resourceServicePanUpdate(d *schema.ResourceData, m interface{}) error {
	return resourceServicePanRead(d, m)
}

func resourceServicePanDelete(d *schema.ResourceData, m interface{}) error {
	client := m.(*alkira.AlkiraClient)

	log.Printf("[INFO] Deleting Service (PAN) %s", d.Id())
	return client.DeleteServicePan(d.Id())
}
