package alkira

import (
	"fmt"
	"log"
	"strconv"

	"github.com/alkiranet/alkira-client-go/alkira"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func resourceAlkiraServiceFortinet() *schema.Resource {
	return &schema.Resource{
		Description: "Manage Fortinet firewall.",
		Create:      resourceFortinetCreate,
		Read:        resourceFortinetRead,
		Update:      resourceFortinetUpdate,
		Delete:      resourceFortinetDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"auto_scale": {
				Description: "Indicate if auto_scale should be enabled for your Fortinet " +
					"firewall. `ON` and `OFF` are accepted values. `OFF` is the default if " +
					"field is omitted",
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"ON", "OFF"}, false),
			},
			"billing_tag_ids": {
				Description: "Billing tag IDs to associate with the service.",
				Type:        schema.TypeList,
				Optional:    true,
				Elem:        &schema.Schema{Type: schema.TypeInt},
			},
			"credential_id": {
				Description: "ID of Fortinet Firewall credential managed by credential resource.",
				Type:        schema.TypeString,
				Required:    true,
			},
			"cxp": {
				Description: "The CXP where the service should be provisioned.",
				Type:        schema.TypeString,
				Required:    true,
			},
			"instances": {
				Type:     schema.TypeList,
				Required: true,
				Description: "An array containing properties for each Fortinet Firewall instance " +
					"that needs to be deployed. The number of instances should be equal to " +
					"`max_instance_count`.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Description: "The name of the Fortinet Firewall instance.",
							Type:        schema.TypeString,
							Required:    true,
						},
						"serial_number": {
							Description: "The serial_number of the Fortinet Firewall instance. " +
								"Required only when licenseType is BRING_YOUR_OWN.",
							Type:     schema.TypeString,
							Optional: true,
						},
						"credential_id": {
							Description: "The ID of the Fortinet Firewall instance credentials. " +
								"Required only when licenseType is BRING_YOUR_OWN.",
							Type:     schema.TypeString,
							Optional: true,
						},
						"id": {
							Description: "The ID of the Fortinet Firewall instance.",
							Type:        schema.TypeInt,
							Computed:    true,
						},
					},
				},
			},
			"license_type": {
				Description:  "Fortinet license type, either `BRING_YOUR_OWN` or `PAY_AS_YOU_GO`.",
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringInSlice([]string{"BRING_YOUR_OWN", "PAY_AS_YOU_GO"}, false),
			},
			"management_server_ip": {
				Description: "The IP addresses used to access the management server.",
				Type:        schema.TypeString,
				Optional:    true,
			},
			"management_server_segment_id": {
				Description: "The segment ID used to access the management server. This segment " +
					"must be present in the list of segments assigned to this Fortinet Firewall service.",
				Type:     schema.TypeInt,
				Required: true,
			},
			"max_instance_count": {
				Description: "The maximum number of Fortinet Firewall instances that should be " +
					"deployed when auto-scale is enabled. Note that auto-scale is not supported " +
					"with Fortinet at this time. max_instance_count must be greater than or " +
					"equal to min_instance_count.",
				Type:     schema.TypeInt,
				Required: true,
			},
			"min_instance_count": {
				Description: "The minimum number of Fortinet Firewall instances that should be " +
					" deployed at any point in time.",
				Type:     schema.TypeInt,
				Optional: true,
				Default:  0,
			},
			"name": {
				Description: "Name of the Fortinet Firewall service.",
				Type:        schema.TypeString,
				Required:    true,
			},
			"segment_ids": {
				Description: "IDs of segments associated with the service.",
				Type:        schema.TypeList,
				Required:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
			"segment_options": {
				Type:        schema.TypeSet,
				Optional:    true,
				Description: "The segment options as used by your Fortinet firewall.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"segment_id": {
							Description: "The ID of the segment.",
							Type:        schema.TypeInt,
							Required:    true,
						},
						"zone_name": {
							Description: "The name of the associated zone.",
							Type:        schema.TypeString,
							Required:    true,
						},
						"groups": {
							Description: "The list of Groups associated with the zone.",
							Type:        schema.TypeList,
							Required:    true,
							Elem:        &schema.Schema{Type: schema.TypeString},
						},
					},
				},
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
			"version": {
				Description: "The version of the Fortinet Firewall.",
				Type:        schema.TypeString,
				Required:    true,
			},
		},
	}
}

func resourceFortinetCreate(d *schema.ResourceData, m interface{}) error {
	client := m.(*alkira.AlkiraClient)

	request, err := generateFortinetRequest(d, m)

	if err != nil {
		return err
	}

	log.Printf("[INFO] Creating fortinet %s", d.Id())
	fmt.Println("request: request.ManagementServer.Segment: ", request.ManagementServer.Segment)
	fmt.Println("request: request.Segments: ", request.ManagementServer.Segment)
	id, err := client.CreateFortinet(request)

	if err != nil {
		return err
	}

	d.SetId(id)
	return resourceFortinetRead(d, m)

}

func resourceFortinetRead(d *schema.ResourceData, m interface{}) error {
	client := m.(*alkira.AlkiraClient)

	f, err := client.GetFortinetById(d.Id())

	if err != nil {
		return err
	}

	d.Set("auto_scale", f.AutoScale)
	d.Set("billing_tag_ids", f.BillingTags)
	d.Set("credential_id", f.CredentialId)
	d.Set("cxp", f.Cxp)
	d.Set("license_type", f.LicenseType)
	d.Set("management_server_ip", f.ManagementServer.IpAddress)
	d.Set("management_server_segment_id", f.ManagementServer.Segment)
	d.Set("max_instance_count", f.MaxInstanceCount)
	d.Set("min_instance_count", f.MinInstanceCount)
	d.Set("name", f.Name)
	d.Set("segment_ids", f.Segments)
	d.Set("segment_options", deflateSegmentOptions(f.SegmentOptions))
	d.Set("size", f.Size)
	d.Set("tunnel_protocol", f.TunnelProtocol)
	d.Set("version", f.Version)

	var instances []map[string]interface{}

	for _, instance := range f.Instances {
		i := map[string]interface{}{
			"id":            instance.Id,
			"name":          instance.Name,
			"serial_number": instance.SerialNumber,
			"credential_id": instance.CredentialId,
		}
		instances = append(instances, i)
	}

	d.Set("instances", instances)

	return nil
}

func resourceFortinetUpdate(d *schema.ResourceData, m interface{}) error {
	client := m.(*alkira.AlkiraClient)

	request, err := generateFortinetRequest(d, m)

	if err != nil {
		return err
	}

	log.Printf("[INFO] Updating Fortinet%s", d.Id())
	err = client.UpdateFortinet(d.Id(), request)

	return err
}

func resourceFortinetDelete(d *schema.ResourceData, m interface{}) error {
	client := m.(*alkira.AlkiraClient)

	log.Printf("[INFO] Deleting Fortinet %s", d.Id())
	return client.DeleteFortinet(d.Id())
}

func generateFortinetRequest(d *schema.ResourceData, m interface{}) (*alkira.Fortinet, error) {
	billingTagIds := convertTypeListToIntList(d.Get("billing_tag_ids").([]interface{}))

	segmentId := strconv.Itoa(d.Get("management_server_segment_id").(int))
	mgmtSegName, err := convertSegmentIdToSegmentName(segmentId, m)
	if err != nil {
		return nil, err
	}

	managementServer := &alkira.FortinetManagmentServer{
		IpAddress: d.Get("management_server_ip").(string),
		Segment:   mgmtSegName,
	}
	instances := expandFortinetInstances(d.Get("instances").([]interface{}))

	// convert segment ids to segment names
	segmentIds := convertTypeListToStringList(d.Get("segment_ids").([]interface{}))
	segmentNames, err := convertSegmentIdsToSegmentNames(segmentIds, m)
	if err != nil {
		return nil, err
	}

	segmentOptions, err := expandSegmentOptions(d.Get("segment_options").(*schema.Set), m)
	if err != nil {
		return nil, err
	}

	service := &alkira.Fortinet{
		AutoScale:        d.Get("auto_scale").(string),
		BillingTags:      billingTagIds,
		CredentialId:     d.Get("credential_id").(string),
		Cxp:              d.Get("cxp").(string),
		Instances:        instances,
		LicenseType:      d.Get("license_type").(string),
		ManagementServer: managementServer,
		MaxInstanceCount: d.Get("max_instance_count").(int),
		MinInstanceCount: d.Get("min_instance_count").(int),
		Name:             d.Get("name").(string),
		Segments:         segmentNames,
		SegmentOptions:   segmentOptions,
		Size:             d.Get("size").(string),
		TunnelProtocol:   d.Get("tunnel_protocol").(string),
		Version:          d.Get("version").(string),
	}

	return service, nil
}
