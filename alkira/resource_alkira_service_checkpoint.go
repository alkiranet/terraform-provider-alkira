package alkira

import (
	"log"
	"strconv"

	"github.com/alkiranet/alkira-client-go/alkira"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func resourceAlkiraCheckpoint() *schema.Resource {
	return &schema.Resource{
		Description: "Manage checkpoint services",
		Create:      resourceCheckpoint,
		Read:        resourceCheckpointRead,
		Update:      resourceCheckpointUpdate,
		Delete:      resourceCheckpointDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"auto_scale": {
				Description: "Indicate if `auto_scale` should be enabled for your checkpoint" +
					"firewall. `ON` and `OFF` are accepted values. `OFF` is the default if " +
					"field is omitted",
				Type:         schema.TypeString,
				Default:      "OFF",
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"ON", "OFF"}, false),
			},
			"billing_tag_ids": {
				Description: "Billing tag IDs to associate with the service.",
				Type:        schema.TypeList,
				Optional:    true,
				Elem:        &schema.Schema{Type: schema.TypeInt},
			},
			"cxp": {
				Description: "CXP region.",
				Type:        schema.TypeString,
				Required:    true,
			},
			"password": {
				Description: "The Checkpoint Firewall service password.",
				Type:        schema.TypeString,
				Required:    true,
			},
			"credential_id": {
				Description: "ID of Checkpoint Firewall credential.",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"description": {
				Description: "The description of the checkpoint service.",
				Type:        schema.TypeString,
				Optional:    true,
			},
			"instance": {
				Type:     schema.TypeList,
				Required: true,
				Description: "An array containing properties for each Checkpoint Firewall instance " +
					"that needs to be deployed. The number of instances should be equal to " +
					"`max_instance_count`.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Description: "The name of the checkpoint instance.",
							Type:        schema.TypeString,
							Required:    true,
						},
						"id": {
							Description: "The ID of the checkpoint instance.",
							Type:        schema.TypeInt,
							Computed:    true,
						},
						"credential_id": {
							Description: "ID of Checkpoint Firewall Instance credential.",
							Type:        schema.TypeString,
							Computed:    true,
						},
						"sic_key": {
							Description: "The checkpoint instance sic keys.",
							Type:        schema.TypeString,
							Required:    true,
						},
					},
				},
			},
			"license_type": {
				Description:  "Checkpoint license type, either `BRING_YOUR_OWN` or `PAY_AS_YOU_GO`.",
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringInSlice([]string{"BRING_YOUR_OWN", "PAY_AS_YOU_GO"}, false),
			},
			"management_server": {
				Type:        schema.TypeSet,
				Required:    true,
				Description: "",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"configuration_mode": {
							Description: "The configuration_mode specifies whether the firewall is " +
								"to be automatically configured by Alkira or not. To automatically " +
								"configure the firewall Alkira needs access to the CheckPoint " +
								"management server. If you choose to use manual configuration " +
								"Alkira will provide the customer information about the Checkpoint " +
								"instances so that you can manually configure the firewall.",
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: validation.StringInSlice([]string{"MANUAL", "AUTOMATED"}, false),
						},
						"domain": {
							Description: "Management server domain.",
							Type:        schema.TypeString,
							Optional:    true,
						},
						"credential_id": {
							Description: "ID of Checkpoint Firewall Managment server credential.",
							Type:        schema.TypeString,
							Computed:    true,
						},
						"global_cidr_list_id": {
							Description: "The ID of the global cidr list to be associated with " +
								"the management server.",
							Type:     schema.TypeInt,
							Required: true,
						},
						"ips": {
							Description: "Management server IPs.",
							Type:        schema.TypeList,
							Required:    true,
							Elem:        &schema.Schema{Type: schema.TypeString},
						},
						"reachability": {
							Description: "Specifies whether the management server " +
								"is publicly reachable or not. If the reachability is " +
								"private then you need to provide the segment to be " +
								"used to access the management server. Default value " +
								"is `PUBLIC`.",
							Type:         schema.TypeString,
							Default:      "PUBLIC",
							Optional:     true,
							ValidateFunc: validation.StringInSlice([]string{"PRIVATE", "PUBLIC"}, false),
						},
						"segment_id": {
							Description: "The ID of the segment to be used to access the management server.",
							Type:        schema.TypeInt,
							Optional:    true,
						},
						"type": {
							Description:  "The type of the management server. either `SMS` or `MDS`.",
							Type:         schema.TypeString,
							Optional:     true,
							ValidateFunc: validation.StringInSlice([]string{"SMS", "MDS"}, false),
						},
						"user_name": {
							Description: "The username of the management server.",
							Type:        schema.TypeString,
							Optional:    true,
						},
						"management_server_password": {
							Description: "The password of the management server.",
							Type:        schema.TypeString,
							Optional:    true,
						},
					},
				},
			},
			"max_instance_count": {
				Description: "The maximum number of Checkpoint Firewall instances that should be " +
					"deployed when auto-scale is enabled. Note that auto-scale is not supported " +
					"with Checkpoint at this time. `max_instance_count` must be greater than or " +
					"equal to `min_instance_count`. (**BETA**)",
				Type:     schema.TypeInt,
				Required: true,
			},
			"min_instance_count": {
				Description: "The minimum number of Checkpoint Firewall instances that should be " +
					"deployed at any point in time. If auto-scale is OFF, min_instance_count must " +
					"equal max_instance_count.",
				Type:     schema.TypeInt,
				Optional: true,
				Default:  0,
			},
			"name": {
				Description: "Name of the Checkpoint Firewall service.",
				Type:        schema.TypeString,
				Required:    true,
			},
			"pdp_ips": {
				Description: "The IPs of the PDP Brokers.",
				Type:        schema.TypeList,
				Optional:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
			"segment_id": {
				Description: "The ID of the segments associated with the service.",
				Type:        schema.TypeInt,
				Required:    true,
			},
			"segment_options": {
				Description: "The segment options as used by your Checkpoint firewall. No more than one " +
					"segment option will be accepted.",
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"segment_id": {
							Description: "The ID of the segment.",
							Type:        schema.TypeInt,
							Required:    true,
						},
						"zone_name": {
							Description: "The name of the associated zone. Default value is `DEFAULT`.",
							Type:        schema.TypeString,
							Optional:    true,
							Default:     "DEFAULT",
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
				Description: "The version of the Checkpoint Firewall.",
				Type:        schema.TypeString,
				Required:    true,
			},
		},
	}
}

func resourceCheckpoint(d *schema.ResourceData, m interface{}) error {
	client := m.(*alkira.AlkiraClient)

	request, err := generateCheckpointRequest(d, m)

	if err != nil {
		log.Printf("[ERROR] failed to generate checkpoint request")
		return err
	}

	id, err := client.CreateCheckpoint(request)

	if err != nil {
		return err
	}

	d.SetId(id)
	return resourceCheckpointRead(d, m)
}

func resourceCheckpointRead(d *schema.ResourceData, m interface{}) error {
	client := m.(*alkira.AlkiraClient)

	checkpoint, err := client.GetCheckpointById(d.Id())
	if err != nil {
		log.Printf("[ERROR] failed to get checkpoint %s", d.Id())
		return err
	}

	segmentIds, err := convertCheckpointSegmentNameToSegmentId(checkpoint.Segments, m)
	if err != nil {
		return err
	}

	d.Set("auto_scale", checkpoint.AutoScale)
	d.Set("billing_tag_ids", checkpoint.BillingTags)
	d.Set("credential_id", checkpoint.CredentialId)
	d.Set("cxp", checkpoint.Cxp)
	d.Set("description", checkpoint.Description)
	d.Set("instance", setCheckpointInstances(d, checkpoint.Instances))
	d.Set("license_type", checkpoint.LicenseType)
	d.Set("management_server", deflateCheckpointManagementServer(*checkpoint.ManagementServer))
	d.Set("max_instance_count", checkpoint.MaxInstanceCount)
	d.Set("min_instance_count", checkpoint.MinInstanceCount)
	d.Set("name", checkpoint.Name)
	d.Set("pdp_ips", checkpoint.PdpIps)
	d.Set("segment_id", segmentIds)
	d.Set("size", checkpoint.Size)
	d.Set("segment_options", deflateSegmentOptions(checkpoint.SegmentOptions))
	d.Set("tunnel_protocol", checkpoint.TunnelProtocol)
	d.Set("version", checkpoint.Version)

	return nil
}

func resourceCheckpointUpdate(d *schema.ResourceData, m interface{}) error {
	client := m.(*alkira.AlkiraClient)

	request, err := generateCheckpointRequest(d, m)

	if err != nil {
		return err
	}

	log.Printf("[INFO] Updating Checkpoint %s", d.Id())
	err = client.UpdateCheckpoint(d.Id(), request)

	return err
}

func resourceCheckpointDelete(d *schema.ResourceData, m interface{}) error {
	client := m.(*alkira.AlkiraClient)

	log.Printf("[INFO] Deleting Checkpoint %s", d.Id())
	return client.DeleteCheckpoint(d.Id())
}

func generateCheckpointRequest(d *schema.ResourceData, m interface{}) (*alkira.Checkpoint, error) {
	client := m.(*alkira.AlkiraClient)

	chpfwCredId := d.Get("credential_id").(string)
	if 0 == len(chpfwCredId) {
		log.Printf("[INFO] Creating Checkpoint Firewall Service Credentials")
		chkpfwName := d.Get("name").(string) + "-" + randomNameSuffix()
		c := alkira.CredentialCheckPointFwService{AdminPassword: d.Get("password").(string)}
		credentialId, err := client.CreateCredential(chkpfwName, alkira.CredentialTypeChkpFw, c, 0)
		if err != nil {
			return nil, err
		}
		d.Set("credential_id", credentialId)

	}

	managementServer, err := expandCheckpointManagementServer(d.Get("name").(string), d.Get("management_server").(*schema.Set), m)
	if err != nil {
		return nil, err
	}

	instances, err := expandCheckpointInstances(d.Get("instance").([]interface{}), m)
	if err != nil {
		return nil, err
	}

	segmentIds := []string{strconv.Itoa(d.Get("segment_id").(int))}
	segmentNames, err := convertSegmentIdsToSegmentNames(segmentIds, m)
	if err != nil {
		return nil, err
	}

	segmentOptions, err := expandCheckpointSegmentOptions(segmentNames[0], d.Get("segment_options").(*schema.Set), m)
	if err != nil {
		return nil, err
	}

	billingTagIds := convertTypeListToIntList(d.Get("billing_tag_ids").([]interface{}))

	return &alkira.Checkpoint{
		AutoScale:        d.Get("auto_scale").(string),
		BillingTags:      billingTagIds,
		CredentialId:     d.Get("credential_id").(string),
		Cxp:              d.Get("cxp").(string),
		Description:      d.Get("description").(string),
		Instances:        instances,
		LicenseType:      d.Get("license_type").(string),
		ManagementServer: managementServer,
		MinInstanceCount: d.Get("min_instance_count").(int),
		MaxInstanceCount: d.Get("max_instance_count").(int),
		Name:             d.Get("name").(string),
		PdpIps:           convertTypeListToStringList(d.Get("pdp_ips").([]interface{})),
		Segments:         segmentNames,
		SegmentOptions:   segmentOptions,
		Size:             d.Get("size").(string),
		TunnelProtocol:   d.Get("tunnel_protocol").(string),
		Version:          d.Get("version").(string),
	}, nil
}
