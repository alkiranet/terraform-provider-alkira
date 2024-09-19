package alkira

import (
	"context"
	"fmt"
	"log"

	"github.com/alkiranet/alkira-client-go/alkira"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func resourceAlkiraInfoblox() *schema.Resource {
	return &schema.Resource{
		Description:   "Provide Infoblox service resource (**BETA**).",
		CreateContext: resourceInfoblox,
		ReadContext:   resourceInfobloxRead,
		UpdateContext: resourceInfobloxUpdate,
		DeleteContext: resourceInfobloxDelete,
		CustomizeDiff: func(ctx context.Context, d *schema.ResourceDiff, m interface{}) error {
			client := m.(*alkira.AlkiraClient)

			old, _ := d.GetChange("provision_state")

			if client.Provision == true && old == "FAILED" {
				d.SetNew("provision_state", "SUCCESS")
			}

			return nil
		},
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
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
				Description: "IDs of billing tags to be associated with " +
					"the service.",
				Type:     schema.TypeSet,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeInt},
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
				Description: "The ID of the global cidr list to be " +
					"associated with the Infoblox service.",
				Type:     schema.TypeInt,
				Required: true,
			},
			"provision_state": {
				Description: "The provision state of the resource.",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"grid_master": {
				Type:     schema.TypeList,
				Required: true,
				Description: "Defines the properties of the Infoblox grid " +
					"master.",
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
				Type:     schema.TypeList,
				Required: true,
				Description: "The properties pertaining to each individual " +
					"instance of the Infoblox service.",
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
							Description: "The host name of the instance. The " +
								"host name MUST always have a suffix `.localdomain`.",
							Type:     schema.TypeString,
							Required: true,
						},
						"model": {
							Description: "The model of the Infoblox instance.",
							Type:        schema.TypeString,
							Required:    true,
						},
						"password": {
							Description: "The password associated with the " +
								"infoblox instance.",
							Type:     schema.TypeString,
							Required: true,
						},
						"type": {
							Description: "The type of the Infoblox instance that " +
								"is to be provisioned. The value could be `MASTER`, " +
								"`MASTER_CANDIDATE` and `MEMBER`.",
							Type:     schema.TypeString,
							Required: true,
							ValidateFunc: validation.StringInSlice([]string{
								"MASTER", "MASTER_CANDIDATE", "MEMBER"}, false),
						},
						"version": {
							Description: "The version of the Infoblox to be " +
								"used. Please check Alkira Portal for all " +
								"supported versions",
							Type:     schema.TypeString,
							Required: true,
						},
					},
				},
			},
			"license_type": {
				Description: "Infoblox license type, only " +
					"`BRING_YOUR_OWN` is supported right now.",
				Type:     schema.TypeString,
				Required: true,
				ValidateFunc: validation.StringInSlice([]string{
					"BRING_YOUR_OWN"}, false),
			},
			"name": {
				Description: "Name of the Infoblox service.",
				Type:        schema.TypeString,
				Required:    true,
			},
			"segment_ids": {
				Description: "IDs of segments associated with the service.",
				Type:        schema.TypeSet,
				Required:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
			"service_group_name": {
				Description: "The name of the service group to be associated " +
					"with the service. A service group represents the " +
					"service in traffic policies, route policies " +
					"and when configuring segment resource shares.",
				Type:     schema.TypeString,
				Required: true,
			},
			"service_group_id": {
				Description: "The ID of the service group to be associated " +
					"with the service. A service group represents the " +
					"service in traffic policies, route policies " +
					"and when configuring segment resource shares.",
				Type:     schema.TypeInt,
				Computed: true,
			},
			"service_group_implicit_group_id": {
				Description: "The ID of the implicit group to be associated " +
					"with the service.",
				Type:     schema.TypeInt,
				Computed: true,
			},
			"shared_secret": {
				Description: "Shared Secret of the InfoBlox grid. " +
					"This cannot be empty.",
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringIsNotEmpty,
			},
			"allow_list_id": {
				Description: "The ID of the `alkira_policy_prefix_list` to be used to whitelist prefixes for the service.",
				Type:        schema.TypeInt,
				Optional:    true,
			},
		},
	}
}

func resourceInfoblox(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {

	client := m.(*alkira.AlkiraClient)
	api := alkira.NewServiceInfoblox(m.(*alkira.AlkiraClient))

	// Construct request
	request, err := generateInfobloxRequest(d, m)

	if err != nil {
		return diag.FromErr(err)
	}

	// Send create request
	response, provState, err, provErr := api.Create(request)

	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(string(response.Id))

	// Set provision state
	if client.Provision == true {
		d.Set("provision_state", provState)

		if provState == "FAILED" {
			return diag.Diagnostics{{
				Severity: diag.Warning,
				Summary:  "PROVISION (CREATE) FAILED",
				Detail:   fmt.Sprintf("%s", provErr),
			}}
		}
	}

	return resourceInfobloxRead(ctx, d, m)
}

func resourceInfobloxRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {

	client := m.(*alkira.AlkiraClient)
	api := alkira.NewServiceInfoblox(m.(*alkira.AlkiraClient))

	infoblox, provState, err := api.GetById(d.Id())

	if err != nil {
		return diag.Diagnostics{{
			Severity: diag.Warning,
			Summary:  "FAILED TO GET RESOURCE",
			Detail:   fmt.Sprintf("%s", err),
		}}
	}

	setAllInfobloxResourceFields(d, infoblox)

	// Set provision state
	if client.Provision == true && provState != "" {
		d.Set("provision_state", provState)
	}

	return nil
}

func resourceInfobloxUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {

	// INIT
	client := m.(*alkira.AlkiraClient)
	api := alkira.NewServiceInfoblox(m.(*alkira.AlkiraClient))

	oldCredentialId := d.Get("credential_id").(string)
	// Construct request
	request, err := generateInfobloxRequest(d, m)

	if err != nil {
		return diag.FromErr(err)
	}

	// Send update request
	provState, err, provErr := api.Update(d.Id(), request)

	if err != nil {
		return diag.FromErr(err)
	}

	newCredentialId := d.Get("credential_id").(string)

	if oldCredentialId != newCredentialId {
		err := deleteInfobloxCredential(oldCredentialId, client)
		if err != nil {
			log.Printf("[WARN] failed to delete old credential %s", err)

		}
	}
	// Set provision state
	if client.Provision == true {
		d.Set("provision_state", provState)

		if provState == "FAILED" {
			return diag.Diagnostics{{
				Severity: diag.Warning,
				Summary:  "PROVISION (UPDATE) FAILED",
				Detail:   fmt.Sprintf("%s", provErr),
			}}
		}
	}

	return resourceInfobloxRead(ctx, d, m)
}

func resourceInfobloxDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {

	client := m.(*alkira.AlkiraClient)
	api := alkira.NewServiceInfoblox(m.(*alkira.AlkiraClient))

	provState, err, provErr := api.Delete(d.Id())

	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId("")

	if client.Provision == true && provState != "SUCCESS" {
		return diag.Diagnostics{{
			Severity: diag.Warning,
			Summary:  "PROVISION (DELETE) FAILED",
			Detail:   fmt.Sprintf("%s", provErr),
		}}
	}

	return nil
}
func createInfobloxCredentials(d *schema.ResourceData, client *alkira.AlkiraClient) (string, error) {
	log.Printf("[INFO] Creating Infoblox Credential")

	infobloxCredentialName := d.Get("name").(string) + "_" + randomNameSuffix()
	infobloxCredential := alkira.CredentialInfoblox{SharedSecret: d.Get("shared_secret").(string)}
	return client.CreateCredential(
		infobloxCredentialName,
		alkira.CredentialTypeInfoblox,
		infobloxCredential,
		0,
	)

}

func deleteInfobloxCredential(infobloxCredentialId string, client *alkira.AlkiraClient) error {

	log.Printf("[INFO] Deleting Infoblox Credential")

	return client.DeleteCredential(infobloxCredentialId, alkira.CredentialTypeInfoblox)
}

func generateInfobloxRequest(d *schema.ResourceData, m interface{}) (*alkira.ServiceInfoblox, error) {
	client := m.(*alkira.AlkiraClient)
	var infobloxCredentialId string

	if d.Get("shared_secret").(string) != "" || d.HasChange("shared_secret") {
		credentialId, err := createInfobloxCredentials(d, client)
		if err != nil {
			return nil, err
		}
		infobloxCredentialId = credentialId
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
	segmentNames, err := convertSegmentIdsToSegmentNames(d.Get("segment_ids").(*schema.Set), m)

	if err != nil {
		return nil, err
	}

	return &alkira.ServiceInfoblox{
		AnyCast:          *anycast,
		BillingTags:      convertTypeSetToIntList(d.Get("billing_tag_ids").(*schema.Set)),
		Cxp:              d.Get("cxp").(string),
		Description:      d.Get("description").(string),
		GlobalCidrListId: d.Get("global_cidr_list_id").(int),
		GridMaster:       *gridMaster,
		Instances:        instances,
		LicenseType:      d.Get("license_type").(string),
		Name:             name,
		Segments:         segmentNames,
		ServiceGroupName: d.Get("service_group_name").(string),
		AllowListId:      d.Get("allow_list_id").(int),
	}, nil
}
