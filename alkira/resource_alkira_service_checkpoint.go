package alkira

import (
	"log"

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
			"credential_id": {
				Description: "ID of Checkpoint Firewall credential managed by credential resource.",
				Type:        schema.TypeString,
				Required:    true,
			},
			"description": {
				Description: "The description of the checkpoint service.",
				Type:        schema.TypeString,
				Optional:    true,
			},
			"instances": {
				Type:     schema.TypeSet,
				Optional: true,
				Description: "An array containing properties for each Checkpoint Firewall instance " +
					"that needs to be deployed. The number of instances should be equal to " +
					"`max_instance_count`.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Description: "The name of the Checkpoint Firewall instance.",
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
								"Alkira will provide the customer information about the checkpoint " +
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
							Description: "This option specifies whether the management server is " +
								"publicly reachable or not. If the reachability is private then you " +
								"need to provide the segment to be used to access the management server. " +
								"Default value is `PUBLIC`.",
							Type:     schema.TypeString,
							Default:  "PUBLIC",
							Optional: true,
						},
						"segment_id": {
							Description: "The ID of the segment to be used to access the management server.",
							Type:        schema.TypeInt,
							Optional:    true,
						},
						"type": {
							Description: "The type of the management server.",
							Type:        schema.TypeString,
							Optional:    true,
						},
						"user_name": {
							Description: "The user_name of the management server.",
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
					"equal to `min_instance_count`.",
				Type:     schema.TypeInt,
				Required: true,
			},
			"min_instance_count": {
				Description: "The minimum number of Checkpoint Firewall instances that should be " +
					"deployed at any point in time.",
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
			"segment_ids": {
				Description: "The IDs of the segments associated with the service.",
				Type:        schema.TypeList,
				Required:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
			"segment_options": {
				Type:        schema.TypeSet,
				Required:    true,
				Description: "The segment options as used by your checkpoint firewall.",
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

	d.Set("auto_scale", checkpoint.AutoScale)
	d.Set("billing_tag_ids", checkpoint.BillingTags)
	d.Set("credential_id", checkpoint.CredentialId)
	d.Set("cxp", checkpoint.Cxp)
	d.Set("description", checkpoint.Description)
	d.Set("instances", deflateCheckpointInstances(checkpoint.Instances))
	d.Set("license_type", checkpoint.LicenseType)
	d.Set("management_server", deflateCheckpointManagementServer(*checkpoint.ManagementServer))
	d.Set("max_instance_count", checkpoint.MaxInstanceCount)
	d.Set("min_instance_count", checkpoint.MinInstanceCount)
	d.Set("name", checkpoint.Name)
	d.Set("pdp_ips", checkpoint.PdpIps)
	d.Set("segment_ids", checkpoint.Segments)
	d.Set("size", checkpoint.Size)
	d.Set("segment_options", deflateCheckpointSegmentOptions(checkpoint.SegmentOptions))
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

	log.Printf("[INFO] Updating Checkpoint%s", d.Id())
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

	allCheckpointResponseDetails, err := getAllCheckpointCredentials(client)
	if err != nil {
		return nil, err
	}

	managementServer, err := expandCheckpointManagementServer(d.Get("management_server").(*schema.Set), client.GetSegmentById)
	if err != nil {
		return nil, err
	}

	managementServerCredential := parseCheckpointCredentialManagementServer(allCheckpointResponseDetails)
	var managementServerCredentialId string
	if managementServerCredential != nil {
		managementServerCredentialId = managementServerCredential.Id
	}
	managementServer.CredentialId = managementServerCredentialId

	instanceRespDetails := parseAllCheckpointCredentialInstances(allCheckpointResponseDetails)
	instances := fromCheckpointCredentialRespDetailsToCheckpointInstance(instanceRespDetails)

	segmentOptions, err := expandCheckpointSegmentOptions(d.Get("segment_options").(*schema.Set), client.GetSegmentById)
	if err != nil {
		return nil, err
	}

	segmentIds := convertTypeListToStringList(d.Get("segment_ids").([]interface{}))
	segmentNames, err := convertSegmentIdsToSegmentNames(client.GetSegmentById, segmentIds)
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
