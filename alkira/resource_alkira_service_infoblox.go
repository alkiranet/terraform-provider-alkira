package alkira

import (
	"errors"
	"log"

	"github.com/alkiranet/alkira-client-go/alkira"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func resourceAlkiraInfoblox() *schema.Resource {
	return &schema.Resource{
		Description: "Provide Infoblox service resource (**BETA**).",
		Create:      resourceInfoblox,
		Read:        resourceInfobloxRead,
		Update:      resourceInfobloxUpdate,
		Delete:      resourceInfobloxDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"anycast": {
				Type:        schema.TypeSet,
				Required:    true,
				Description: "Defines the AnyCast policy.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"enabled": {
							Description: "Defines if AnyCast should be enabled. Default value is `false`.",
							Type:        schema.TypeBool,
							Optional:    true,
							Default:     false,
						},
						"ips": {
							Description: "The IPs to be used when AnyCast is enabled. When AnyCast " +
								"is enabled this list cannot be empty. The IPs used for AnyCast MUST " +
								"NOT overlap the CIDR of `alkira_segment` resource associated with " +
								"the service.",
							Type:     schema.TypeList,
							Optional: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
						"backup_cxps": {
							Description: "The `backup_cxps` to be used when the current " +
								"Infoblox service is not available. It also needs to " +
								"have a configured Infoblox service in order to take advantage of " +
								"this feature. It is NOT required that the `backup_cxps` should have " +
								"a configured Infoblox service before it can be designated as a backup.",
							Type:     schema.TypeList,
							Optional: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
					},
				},
			},
			"billing_tag_ids": {
				Description: "Billing tag IDs to associate with the service.",
				Type:        schema.TypeList,
				Optional:    true,
				Elem:        &schema.Schema{Type: schema.TypeInt},
			},
			"cxp": {
				Description: "The CXP where the service should be provisioned.",
				Type:        schema.TypeString,
				Required:    true,
			},
			"description": {
				Description: "The description of the Infoblox service.",
				Type:        schema.TypeString,
				Optional:    true,
			},
			"global_cidr_list_id": {
				Description: "The ID of the global cidr list to be associated with " +
					"the Infoblox service.",
				Type:     schema.TypeInt,
				Required: true,
			},
			// NOTE: for v1 of Infoblox support we are not supporting the creation of a new infoblox
			// service on behalf of our customer. Instead the customer must have a preexsting
			// infoblox service. Future releases will allow for this. At that time this comment
			// should be removed.
			"grid_master": {
				Type:        schema.TypeList,
				Required:    true,
				Description: "Defines the properties of the Infoblox grid master.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"external": {
							Description: "External indicates if a new grid master should be " +
								"created or if an existing grid master should be used. Default " +
								"value is `false`.",
							Type:     schema.TypeBool,
							Optional: true,
							Default:  false,
						},
						"ip": {
							Description: "The IP address of the grid master.",
							Type:        schema.TypeString,
							Computed:    true,
							Optional:    true,
						},
						"name": {
							Description: "Name of the grid master.",
							Type:        schema.TypeString,
							Required:    true,
						},
						"username": {
							Description: "The Grid Master user name.",
							Type:        schema.TypeString,
							Required:    true,
						},
						"password": {
							Description: "The Grid Master password.",
							Type:        schema.TypeString,
							Required:    true,
						},
						"credential_id": {
							Description: "The credential ID of the Grid Master.",
							Type:        schema.TypeString,
							Computed:    true,
						},
					},
				},
			},
			"instance": {
				Type:        schema.TypeList,
				Required:    true,
				Description: "The properties pertaining to each individual instance of the Infoblox service.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"anycast_enabled": {
							Description: " This knob controls whether AnyCast is to be enabled " +
								"for this instance or not. AnyCast can only be enabled on an " +
								"instance if it is also enabled on the service. The default value is `false`.",
							Type:     schema.TypeBool,
							Optional: true,
							Default:  false,
						},
						"id": {
							Description: "The ID of the Infoblox instance.",
							Type:        schema.TypeInt,
							Computed:    true,
						},
						"credential_id": {
							Description: "The credential ID of the Infoblox instance.",
							Type:        schema.TypeString,
							Computed:    true,
						},
						"hostname": {
							Description: "The host name of the instance. The host name MUST always have a suffix `.localdomain`.",
							Type:        schema.TypeString,
							Required:    true,
						},
						"model": {
							Description: "The model of the Infoblox instance.",
							Type:        schema.TypeString,
							Required:    true,
						},
						"password": {
							Description: "The password associated with the infoblox instance.",
							Type:        schema.TypeString,
							Required:    true,
						},
						"type": {
							Description:  "The type of the Infoblox instance that is to be provisioned. The value could be `MASTER`, `MASTER_CANDIDATE` and `MEMBER`.",
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: validation.StringInSlice([]string{"MASTER", "MASTER_CANDIDATE", "MEMBER"}, false),
						},
						"version": {
							Description: "The version of the Infoblox instance to be used.",
							Type:        schema.TypeString,
							Required:    true,
						},
					},
				},
			},
			"license_type": {
				Description:  "Infoblox license type, only `BRING_YOUR_OWN` is supported right now.",
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringInSlice([]string{"BRING_YOUR_OWN"}, false),
			},
			"name": {
				Description: "Name of the Infoblox service.",
				Type:        schema.TypeString,
				Required:    true,
			},
			"segment_ids": {
				Description: "IDs of segments associated with the service.",
				Type:        schema.TypeList,
				Required:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
			"service_group_name": {
				Description: "The name of the service group to be associated with the service. " +
					"A service group represents the service in traffic policies, route policies " +
					"and when configuring segment resource shares.",
				Type:     schema.TypeString,
				Required: true,
			},
			"shared_secret": {
				Description:  "Shared Secret of the InfoBlox grid. This cannot be empty.",
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringIsNotEmpty,
			},
		},
	}
}

func resourceInfoblox(d *schema.ResourceData, m interface{}) error {
	client := m.(*alkira.AlkiraClient)

	request, err := generateInfobloxRequest(d, m)

	if err != nil {
		log.Printf("[ERROR] failed to generate infoblox request")
		return err
	}

	id, err := client.CreateInfoblox(request)

	if err != nil {
		return err
	}

	d.SetId(id)
	return resourceInfobloxRead(d, m)
}

func resourceInfobloxRead(d *schema.ResourceData, m interface{}) error {
	client := m.(*alkira.AlkiraClient)

	infoblox, err := client.GetInfobloxById(d.Id())
	if err != nil {
		log.Printf("[ERROR] failed to get infoblox %s", d.Id())
		return err
	}

	setAllInfobloxResourceFields(d, infoblox)

	return nil
}

func resourceInfobloxUpdate(d *schema.ResourceData, m interface{}) error {
	client := m.(*alkira.AlkiraClient)

	request, err := generateInfobloxRequest(d, m)

	if err != nil {
		return err
	}

	log.Printf("[INFO] Updating infoblox%s", d.Id())
	err = client.UpdateInfoblox(d.Id(), request)

	return err
}

func resourceInfobloxDelete(d *schema.ResourceData, m interface{}) error {
	client := m.(*alkira.AlkiraClient)

	log.Printf("[INFO] Deleting infoblox %s", d.Id())
	return client.DeleteInfoblox(d.Id())
}

func generateInfobloxRequest(d *schema.ResourceData, m interface{}) (*alkira.Infoblox, error) {
	client := m.(*alkira.AlkiraClient)

	//Create Infoblox Service Credential
	name := d.Get("name").(string)
	nameWithSuffix := name + randomNameSuffix()
	shared_secret := d.Get("shared_secret").(string)

	var infobloxCredentialId string
	if shared_secret != "" {
		err := errors.New("")
		infobloxCredentialId, err = client.CreateCredential(
			nameWithSuffix,
			alkira.CredentialTypeInfoblox,
			&alkira.CredentialInfoblox{SharedSecret: shared_secret},
			0,
		)
		if err != nil {
			return nil, err
		}
	}

	//Parse Grid Master
	gmSet := d.Get("grid_master").([]interface{})
	gridMaster, err := expandInfobloxGridMaster(gmSet, infobloxCredentialId, m)
	if err != nil {
		return nil, err
	}

	//Parse Instances
	instanceList := d.Get("instance").([]interface{})
	instances, err := expandInfobloxInstances(instanceList, m)
	if err != nil {
		return nil, err
	}

	//Parse Anycast
	anycast, err := expandInfobloxAnycast(d.Get("anycast").(*schema.Set))
	if err != nil {
		return nil, err
	}

	//segmentIdsToSegmentNames
	ids := convertTypeListToStringList(d.Get("segment_ids").([]interface{}))
	segment_names, err := convertSegmentIdsToSegmentNames(ids, m)
	if err != nil {
		return nil, err
	}

	return &alkira.Infoblox{
		AnyCast:          *anycast,
		BillingTags:      convertTypeListToIntList(d.Get("billing_tag_ids").([]interface{})),
		Cxp:              d.Get("cxp").(string),
		Description:      d.Get("description").(string),
		GlobalCidrListId: d.Get("global_cidr_list_id").(int),
		GridMaster:       *gridMaster,
		Instances:        instances,
		LicenseType:      d.Get("license_type").(string),
		Name:             name,
		Segments:         segment_names,
		ServiceGroupName: d.Get("service_group_name").(string),
	}, nil
}
