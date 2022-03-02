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
				Description: "Billing tag IDs to associate with the service.",
				Type:        schema.TypeList,
				Optional:    true,
				Elem:        &schema.Schema{Type: schema.TypeInt},
			},
			"bundle": {
				Description: "The software image bundle that would be used for" +
					"PAN instance deployment. This is applicable for licenseType" +
					"`PAY_AS_YOU_GO` only. If not provided, the default" +
					"`PAN_VM_300_BUNDLE_2` would be used. However `PAN_VM_300_BUNDLE_2`" +
					"is legacy bundle and is not supported on AWS. It is recommended" +
					"to use `VM_SERIES_BUNDLE_1` and `VM_SERIES_BUNDLE_2` (supports " +
					"Global Protect).",
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"VM_SERIES_BUNDLE_1", "VM_SERIES_BUNDLE_2", "PAN_VM_300_BUNDLE_2"}, false),
			},
			"credential_id": {
				Description: "ID of PAN credential managed by credential resource.",
				Type:        schema.TypeString,
				Required:    true,
			},
			"global_protect_enabled": {
				Description: "Enable global protect option or not. Default is `false`",
				Type:        schema.TypeBool,
				Required:    true,
			},
			"global_protect_segment_options": {
				Description: "A mapping of segment_name -> zones_to_groups. The only segment names " +
					"allowed are the segments that are already associated with the service." +
					"options should apply. If global_protect_enabled is set to false, " +
					"global_protect_segment_options shound not be included in your request.",
				Type: schema.TypeSet,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"segment_name": {
							Description: "The name of the segment to which the global protect options should apply",
							Type:        schema.TypeString,
							Required:    true,
						},
						"remote_user_zone_name": {
							Description: "Firewall security zone is created using the zone name for remote user sessions.",
							Type:        schema.TypeString,
							Required:    true,
						},
						"portal_fqdn_prefix": {
							Description: "Prefix for the global protect portal FQDN, this would " +
								"be prepended to customer specific alkira domain For Example: " +
								"if prefix is abc and tenant name is example then the FQDN would " +
								"be abc.example.gpportal.alkira.com",
							Type:     schema.TypeString,
							Required: true,
						},
						"service_group_name": {
							Description: "The name of the service group. A group with the same name will be created.",
							Type:        schema.TypeString,
							Required:    true,
						},
					},
				},
				Optional: true,
			},
			"cxp": {
				Description: "The CXP where the service should be provisioned.",
				Type:        schema.TypeString,
				Required:    true,
			},
			"instance": {
				Type:     schema.TypeSet,
				Required: true,
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
			},
			"license_type": {
				Description:  "PAN license type, either `BRING_YOUR_OWN` or `PAY_AS_YOU_GO`.",
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringInSlice([]string{"BRING_YOUR_OWN", "PAY_AS_YOU_GO"}, false),
			},
			"panorama_enabled": {
				Description: "Enable Panorama or not. Default value is `false`.",
				Type:        schema.TypeBool,
				Optional:    true,
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
				Description: "Management Segment ID.",
				Type:        schema.TypeInt,
				Required:    true,
			},
			"max_instance_count": {
				Description: "Max number of Panorama instances for auto scale.",
				Type:        schema.TypeInt,
				Required:    true,
			},
			"min_instance_count": {
				Description: "Minimal number of Panorama instances for auto scale. Default value is `0`.",
				Type:        schema.TypeInt,
				Optional:    true,
				Default:     0,
			},
			"name": {
				Description: "Name of the PAN service.",
				Type:        schema.TypeString,
				Required:    true,
			},
			"segment_ids": {
				Description: "IDs of segments associated with the service.",
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
							Optional:    true,
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

	request, err := generateServicePanRequest(d, m)

	if err != nil {
		return err
	}

	log.Printf("[INFO] Creating service-pan %s", d.Id())
	id, err := client.CreateServicePan(request)

	if err != nil {
		return err
	}

	d.SetId(id)
	return resourceServicePanRead(d, m)
}

func resourceServicePanRead(d *schema.ResourceData, m interface{}) error {
	client := m.(*alkira.AlkiraClient)

	pan, err := client.GetServicePanById(d.Id())

	if err != nil {
		return err
	}

	d.Set("billing_tag_ids", pan.BillingTagIds)
	d.Set("credential_id", pan.CredentialId)
	d.Set("cxp", pan.CXP)
	d.Set("license_type", pan.LicenseType)
	d.Set("management_segment_id", pan.ManagementSegmentId)
	d.Set("max_instance_count", pan.MaxInstanceCount)
	d.Set("min_instance_count", pan.MinInstanceCount)
	d.Set("name", pan.Name)
	d.Set("panorama_enabled", pan.PanoramaEnabled)
	d.Set("segment_ids", pan.SegmentIds)
	d.Set("size", pan.Size)
	d.Set("tunnel_protocol", pan.TunnelProtocol)
	d.Set("type", pan.Type)
	d.Set("version", pan.Version)

	var instances []map[string]interface{}

	for _, instance := range pan.Instances {
		i := map[string]interface{}{
			"name":          instance.Name,
			"credential_id": instance.CredentialId,
		}
		instances = append(instances, i)
	}

	d.Set("instance", instances)

	return nil
}

func resourceServicePanUpdate(d *schema.ResourceData, m interface{}) error {
	client := m.(*alkira.AlkiraClient)

	request, err := generateServicePanRequest(d, m)

	if err != nil {
		return err
	}

	log.Printf("[INFO] Updateing service-pan %s", d.Id())
	err = client.UpdateServicePan(d.Id(), request)

	return err
}

func resourceServicePanDelete(d *schema.ResourceData, m interface{}) error {
	client := m.(*alkira.AlkiraClient)

	log.Printf("[INFO] Deleting service-pan %s", d.Id())
	return client.DeleteServicePan(d.Id())
}

func generateServicePanRequest(d *schema.ResourceData, m interface{}) (*alkira.ServicePan, error) {

	billingTagIds := convertTypeListToIntList(d.Get("billing_tag_ids").([]interface{}))
	instances := expandPanInstances(d.Get("instance").(*schema.Set))
	segmentIds := convertTypeListToIntList(d.Get("segment_ids").([]interface{}))
	segmentOptions := expandPanSegmentOptions(d.Get("zones_to_groups").(*schema.Set))
	globalProtectSegmentOptions := expandGlobalProtectSegmentOptions(
		d.Get("global_protect_segment_options").(*schema.Set),
	)

	service := &alkira.ServicePan{
		BillingTagIds:               billingTagIds,
		CXP:                         d.Get("cxp").(string),
		CredentialId:                d.Get("credential_id").(string),
		GlobalProtectEnabled:        d.Get("global_protect_enabled").(bool),
		GlobalProtectSegmentOptions: globalProtectSegmentOptions,
		Instances:                   instances,
		LicenseType:                 d.Get("license_type").(string),
		MaxInstanceCount:            d.Get("max_instance_count").(int),
		MinInstanceCount:            d.Get("min_instance_count").(int),
		ManagementSegmentId:         d.Get("management_segment_id").(int),
		Name:                        d.Get("name").(string),
		PanoramaEnabled:             d.Get("panorama_enabled").(bool),
		PanoramaDeviceGroup:         d.Get("panorama_device_group").(string),
		PanoramaIpAddress:           d.Get("panorama_ip_address").(string),
		PanoramaTemplate:            d.Get("panorama_template").(string),
		SegmentOptions:              segmentOptions,
		SegmentIds:                  segmentIds,
		TunnelProtocol:              d.Get("tunnel_protocol").(string),
		Size:                        d.Get("size").(string),
		Type:                        d.Get("type").(string),
		Version:                     d.Get("version").(string),
	}

	return service, nil
}
