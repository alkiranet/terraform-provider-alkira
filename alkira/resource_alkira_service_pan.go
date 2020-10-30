package alkira

import (
	"log"
	"strconv"

	"github.com/alkiranet/alkira-client-go/alkira"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceAlkiraServicePan() *schema.Resource {
	return &schema.Resource{
		Create: resourceServicePanCreate,
		Read:   resourceServicePanRead,
		Update: resourceServicePanUpdate,
		Delete: resourceServicePanDelete,

		Schema: map[string]*schema.Schema{
			"billing_tags": {
				Type:     schema.TypeList,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeInt},
			},
			"credential_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"cxp": {
				Type:     schema.TypeString,
				Required: true,
			},
			"group": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"instance": {
				Type: schema.TypeSet,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:     schema.TypeString,
							Required: true,
						},
						"credential_id": {
							Type:     schema.TypeString,
							Required: true,
						},
					},
				},
				Required: true,
			},
			"license_type": {
				Type:     schema.TypeString,
				Required: true,
			},
			"panorama_enabled": {
				Type:     schema.TypeString,
				Required: true,
			},
			"panorama_device_group": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"panorama_ip_address": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"panorama_template": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"management_segment": {
				Type:     schema.TypeString,
				Required: true,
			},
			"max_instance_count": {
				Type:     schema.TypeInt,
				Optional: true,
				Default:  1,
			},
			"min_instance_count": {
				Type:     schema.TypeInt,
				Optional: true,
				Default:  1,
			},
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"zones_to_groups": {
				Type: schema.TypeSet,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"segment_name": {
							Type:     schema.TypeString,
							Required: true,
						},
						"zone_name": {
							Type:     schema.TypeString,
							Required: true,
						},
						"groups": {
							Type:     schema.TypeList,
							Required: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
					},
				},
				Optional: true,
			},
			"segments": {
				Type:     schema.TypeList,
				Required: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"service_id": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"size": {
				Type:     schema.TypeString,
				Required: true,
			},
			"type": {
				Type:     schema.TypeString,
				Required: true,
			},
			"version": {
				Type:     schema.TypeString,
				Required: true,
			},
		},
	}
}

func resourceServicePanCreate(d *schema.ResourceData, m interface{}) error {
	client := m.(*alkira.AlkiraClient)

	billingTags := convertTypeListToIntList(d.Get("billing_tags").([]interface{}))
	instances := expandPanInstances(d.Get("instance").(*schema.Set))
	segments := convertTypeListToStringList(d.Get("segments").([]interface{}))
	segmentOptions := expandPanSegmentOptions(d.Get("zones_to_groups").(*schema.Set))

	service := &alkira.ServicePanRequest{
		BillingTags:         billingTags,
		CXP:                 d.Get("cxp").(string),
		CredentialId:        d.Get("credential_id").(string),
		Instances:           instances,
		LicenseType:         d.Get("license_type").(string),
		MaxInstanceCount:    d.Get("max_instance_count").(int),
		MinInstanceCount:    d.Get("min_instance_count").(int),
		ManagementSegment:   d.Get("management_segment").(string),
		Name:                d.Get("name").(string),
		PanoramaEnabled:     d.Get("panorama_enabled").(string),
		PanoramaDeviceGroup: d.Get("panorama_device_group").(string),
		PanoramaIpAddress:   d.Get("panorama_ip_address").(string),
		PanoramaTemplate:    d.Get("panorama_template").(string),
		SegmentOptions:      segmentOptions,
		Segments:            segments,
		Size:                d.Get("size").(string),
		Type:                d.Get("type").(string),
		Version:             d.Get("version").(string),
	}

	log.Printf("[INFO] Creating Service (PAN) %s", d.Id())
	id, err := client.CreateServicePan(service)

	if err != nil {
		return err
	}

	d.SetId(strconv.Itoa(id))
	d.Set("service_id", id)

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
	err := client.DeleteServicePan(d.Get("service_id").(int))

	if err != nil {
		return err
	}

	return nil
}
