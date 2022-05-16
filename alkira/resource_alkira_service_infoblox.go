package alkira

import (
	"fmt"
	"log"

	"github.com/alkiranet/alkira-client-go/alkira"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func resourceAlkiraInfoblox() *schema.Resource {
	return &schema.Resource{
		Description: "Manage Infoblox services",
		Create:      resourceInfoblox,
		Read:        resourceInfobloxRead,
		Update:      resourceInfobloxUpdate,
		Delete:      resourceInfobloxDelete,

		Schema: map[string]*schema.Schema{
			"anycast": {
				Type:     schema.TypeSet,
				Required: true,
				Description: "Defines the AnyCast policy to be used with the Infoblox Service. " +
					"Based on this AnyCast policy some implicit route policies and prefix lists get " +
					"generated. These route policies and prefix lists will have the prefix " +
					"ALK-SYSTEM-GENERATED-INFOBLOX. These route policies and prefix lists cannot " +
					"be deleted or modified directly their lifecycle is bound by the Infoblox " +
					"services that are configured on the network. AnyCast may be enabled/disabled " +
					"at the instance level as well. For AnyCast to be enabled for an instance it " +
					"MUST be enabled both at the service and the instance level. If AnyCast is " +
					"NOT enabled at the service level it will stay disabled for all instances.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"enabled": {
							Description: "Defines if AnyCast should be enabled. Default is `false`",
							Type:        schema.TypeBool,
							Optional:    true,
							Default:     false,
						},
						"ips": {
							Description: "The IPs to be used when AnyCast is enabled. When AnyCast " +
								"is enabled this list cannot be empty. The ips used for AnyCast MUST " +
								"NOT overlap the cidr used for the segment IP block associated with " +
								"the service",
							Type:     schema.TypeList,
							Required: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
						"backup_cxps": {
							Description: "The backup_cxps to be used when the current " +
								"Infoblox service is not available. The backup_cxps also need to " +
								"have a configured Infoblox service inorder to take advantage of " +
								"this feature. It is NOT required that the backup_cxps should have " +
								"a configured Infoblox service before it can be designated as a backup.",
							Type:     schema.TypeList,
							Required: true,
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
			//NOTE: for v1 of Infoblox support we are not supporting the creation of a new infoblox
			//service on behalf of our customer. Instead the customer must have a preexsting
			//infoblox service. Future releases will allow for this. At that time this comment
			//should be removed.
			"grid_master": {
				Type:     schema.TypeSet,
				Required: true,
				Description: "Defines the properties of the Infoblox grid master. The Infoblox " +
					"grid master needs to exist before other instances of a the grid can be added. " +
					"The grid master can either be provisioned by Alkira or could already be " +
					"provisioned externally. Some of these properties only need to be provided when " +
					"the grid master is external. If the grid master needs to be provisioned " +
					"internally by Alkira then an instance needs to be added to Infoblox service " +
					"configuration with type = MASTER",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"external": {
							Description: "External indicates if a new grid master should be " +
								"created or if an existing grid master should be used. NOTE: " +
								"creation of new external grid masters is not supported at " +
								"this time, but will be supported in future releases.",
							Type:         schema.TypeBool,
							Optional:     true,
							Default:      false,
							ValidateFunc: ExternalMustBeFalse(),
						},
						"grid_master_credential_id": {
							Description: "The grid master credentials.",
							Type:        schema.TypeString,
							Required:    true,
						},
						"ip": {
							Description: "The ip address of the grid master.",
							Type:        schema.TypeString,
							Required:    true,
						},
						"name": {
							Description: "Name of the grid master.",
							Type:        schema.TypeString,
							Required:    true,
						},
						"shared_secret_credential_id": {
							Description: "The shared secret credentials. This is the shared " +
								"secret to be used when InfoBlox instances need to join the grid.",
							Type:     schema.TypeString,
							Required: true,
						},
					},
				},
			},
			"instances": {
				Type:        schema.TypeSet,
				Required:    true,
				Description: "The properties pertaining to each individual instance of the Infoblox service.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"any_cast_enabled": {
							Description: " This knob controls whether AnyCast is to be enabled " +
								"for this instance or not. AnyCast can only be enabled on an " +
								"instance if it is also enabled on the service.",
							Type:     schema.TypeBool,
							Optional: true,
							Default:  false,
						},
						"credential_id": {
							Description: "The credentials for the instance.",
							Type:        schema.TypeString,
							Required:    true,
						},
						"host_name": {
							Description: "The host name of the instance. The host name MUST always have a suffix `.localdomain`.",
							Type:        schema.TypeString,
							Required:    true,
						},
						"model": {
							Description: "The model of the InfoBlox instance.",
							Type:        schema.TypeString,
							Required:    true,
						},
						"type": {
							Description: "The type of the InfoBlox instance that is to be provisioned. " +
								"There can only be one MASTER ever provisioned. When the grid master " +
								"is provisioned by Alkira, provisioning needs to happen in two steps. " +
								"First the InfoBlox service must be provisioned with only 1 instance " +
								"of type MASTER. Subsequently other instances of the grid may be " +
								"added to the instances list and provisioned. When the grid master " +
								"is external (i.e not provisioned by Alkira) then no instances of " +
								"type MASTER should be configured.",
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: validation.StringInSlice([]string{"MASTER", "MASTER_CANDIDATE, MEMBER"}, false),
						},
						"version": {
							Description: "The version of the InfoBlox instance to be used.",
							Type:        schema.TypeString,
							Required:    true,
						},
					},
				},
			},
			"license_type": {
				Description:  "Infoblox license type, either `BRING_YOUR_OWN` or `PAY_AS_YOU_GO`.",
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringInSlice([]string{"BRING_YOUR_OWN", "PAY_AS_YOU_GO"}, false),
			},
			"name": {
				Description: "Name of the Infoblox service.",
				Type:        schema.TypeString,
				Required:    true,
			},
			"segment_names": { //TODO(mac): these needs to be converted from names to ids in the request to the backend
				Description: "Names of segments associated with the service.",
				Type:        schema.TypeList,
				Required:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
			"service_group_name": {
				Description: " The name of the service group to be associated with the service. " +
					"A service group represents the service in traffic policies, route policies " +
					"and when configuring segment resource shares.",
				Type:     schema.TypeString,
				Required: true,
			},
			"size": {
				Description:  "The size of the service, one of `SMALL`, `MEDIUM`, `LARGE`.",
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"SMALL", "MEDIUM", "LARGE"}, false),
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
	gridMaster, err := expandInfobloxGridMaster(d.Get("grid_master").(*schema.Set))
	if err != nil {
		return nil, err
	}

	anycast, err := expandInfobloxAnycast(d.Get("anycast").(*schema.Set))
	if err != nil {
		return nil, err
	}

	billingTagIds := convertTypeListToIntList(d.Get("billing_tag_ids").([]interface{}))
	instances := expandInfobloxInstances(d.Get("instances").(*schema.Set))
	segmentNames := convertTypeListToStringList(d.Get("segment_names").([]interface{}))

	return &alkira.Infoblox{
		AnyCast:          *anycast,
		BillingTags:      billingTagIds,
		Cxp:              d.Get("cxp").(string),
		Description:      d.Get("description").(string),
		GlobalCidrListId: d.Get("global_cidr_list_id").(int),
		GridMaster:       *gridMaster,
		Instances:        instances,
		LicenseType:      d.Get("license_type").(string),
		Name:             d.Get("name").(string),
		Segments:         segmentNames,
		ServiceGroupName: d.Get("service_group_name").(string),
		Size:             d.Get("size").(string),
	}, nil
}

func ExternalMustBeFalse() schema.SchemaValidateFunc {
	return func(i interface{}, b string) (warnings []string, errors []error) {
		v, ok := i.(bool)
		if !ok {
			errors = append(errors, fmt.Errorf("expected type boolean"))
			return warnings, errors
		}

		if !v {
			return warnings, errors
		}

		errors = append(errors, fmt.Errorf("expected value to be false: future software updates will allow for an input value of true."))
		return warnings, errors
	}
}
