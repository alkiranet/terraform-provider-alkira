package alkira

import (
	"log"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/alkiranet/alkira-client-go/alkira"
)

func resourceAlkiraServicePan() *schema.Resource {
	return &schema.Resource{
		Create: resourceServicePanCreate,
		Read:   resourceServicePanRead,
		Update: resourceServicePanUpdate,
		Delete: resourceServicePanDelete,

		Schema: map[string]*schema.Schema{
			"credential_id": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Description: "A user group that the connector belongs to",
			},
			"cxp": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				Description: "The CXP to be used for the connector",
			},
			"group": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Description: "A user group that the connector belongs to",
			},
			"instance": {
				Type:     schema.TypeSet,
				Elem:     &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:      schema.TypeString,
							Required:  true,
						},
						"credential_id": {
							Type:      schema.TypeString,
							Required:  true,
						},
					},
				},
				Required: true,
			},
			"license_type": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				Description: "PAN license type",
			},
			"panorama_enabled": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				Description: "PAN license type",
			},
			"management_segment": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				Description: "The management segment",
			},
			"max_instance_count": &schema.Schema{
				Type:     schema.TypeInt,
				Optional: true,
				Default:  1,
			},
			"min_instance_count": &schema.Schema{
				Type:     schema.TypeInt,
				Optional: true,
				Default:  1,
			},
			"name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				Description: "The name of the PAN service",
			},
			"zones_to_groups": {
				Type:     schema.TypeSet,
				Elem:     &schema.Resource{
					Schema: map[string]*schema.Schema{
						"segment_name": {
							Type:      schema.TypeString,
							Required:  true,
						},
						"zone_name": {
							Type:      schema.TypeString,
							Required:  true,
						},
						"groups": {
							Type:      schema.TypeList,
							Required:  true,
							Elem:      &schema.Schema{Type: schema.TypeString},
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
			"service_id": &schema.Schema{
				Type:     schema.TypeInt,
				Computed: true,
			},
			"size": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"type": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"version": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
		},
	}
}

func resourceServicePanCreate(d *schema.ResourceData, m interface{}) error {
	client := m.(*alkira.AlkiraClient)

	instances      := expandPanInstances(d.Get("instance").(*schema.Set))
	segments       := convertTypeListToStringList(d.Get("segments").([]interface{}))
	segmentOptions := expandPanSegmentOptions(d.Get("zones_to_groups").(*schema.Set))

	service := &alkira.ServicePanRequest{
		CXP:              d.Get("cxp").(string),
		CredentialId:     d.Get("credential_id").(string),
		Instances:        instances,
		LicenseType:      d.Get("license_type").(string),
		MaxInstanceCount: d.Get("max_instance_count").(int),
		MinInstanceCount: d.Get("min_instance_count").(int),
		ManagementSegment:d.Get("management_segment").(string),
		Name:             d.Get("name").(string),
		PanoramaEnabled:  d.Get("panorama_enabled").(string),
		SegmentOptions:   segmentOptions,
        Segments:         segments,
        Size:             d.Get("size").(string),
        Type:             d.Get("type").(string),
        Version:          d.Get("version").(string),
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
