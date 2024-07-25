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

func resourceAlkiraServicePan() *schema.Resource {
	return &schema.Resource{
		Description: "Manage Palo Alto Firewall service.\n\n" +
			"When `panorama_enabled` is set to `true`, `pan_username` and " +
			"`pan_password` are required.",
		CreateContext: resourceServicePanCreate,
		ReadContext:   resourceServicePanRead,
		UpdateContext: resourceServicePanUpdate,
		DeleteContext: resourceServicePanDelete,
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
			"billing_tag_ids": {
				Description: "IDs of billing tags to be associated with " +
					"the service.",
				Type:     schema.TypeSet,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeInt},
			},
			"bundle": {
				Description: "The software image bundle that would be used for" +
					"PAN instance deployment. This is applicable for licenseType" +
					"`PAY_AS_YOU_GO` only. If not provided, the default" +
					"`PAN_VM_300_BUNDLE_2` would be used. However `PAN_VM_300_BUNDLE_2`" +
					"is legacy bundle and is not supported on AWS. It is recommended" +
					"to use `VM_SERIES_BUNDLE_1` and `VM_SERIES_BUNDLE_2` (supports " +
					"Global Protect).",
				Type:     schema.TypeString,
				Optional: true,
				ValidateFunc: validation.StringInSlice([]string{
					"VM_SERIES_BUNDLE_1",
					"VM_SERIES_BUNDLE_2",
					"PAN_VM_300_BUNDLE_2"}, false),
			},
			"provision_state": {
				Description: "The provision state of the service.",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"pan_password": {
				Description: "PAN Panorama password.",
				Type:        schema.TypeString,
				Required:    true,
			},
			"pan_username": {
				Description: "PAN Panorama username. For AWS, username should " +
					"be `admin`. For AZURE, it should be `akadmin`.",
				Type:     schema.TypeString,
				Required: true,
			},
			"pan_license_key": {
				Description: "PAN Panorama license Key.",
				Type:        schema.TypeString,
				Optional:    true,
			},
			"credential_id": {
				Description: "ID of PAN credential.",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"global_protect_enabled": {
				Description: "Enable global protect option or not. " +
					"Default is `false`",
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"global_protect_segment_options": {
				Description: "Segment options for segments that are already " +
					"associated with the service. Options should " +
					"apply. If `global_protect_enabled` is set to false, " +
					"`global_protect_segment_options` shound not be included " +
					"in your request.",
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"segment_id": {
							Description: "The name of the segment to which the " +
								"global protect options should apply",
							Type:     schema.TypeString,
							Required: true,
						},
						"remote_user_zone_name": {
							Description: "Firewall security zone is created using " +
								"the zone name for remote user sessions.",
							Type:     schema.TypeString,
							Required: true,
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
							Description: "The name of the service group. A group " +
								"with the same name will be created.",
							Type:     schema.TypeString,
							Required: true,
						},
					},
				},
			},
			"cxp": {
				Description: "The CXP where the service should be provisioned.",
				Type:        schema.TypeString,
				Required:    true,
			},
			"instance": {
				Type:     schema.TypeList,
				Required: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Description: "The name of the PAN instance.",
							Type:        schema.TypeString,
							Default:     "",
							Optional:    true,
						},
						"id": {
							Description: "The ID of the PAN instance.",
							Type:        schema.TypeInt,
							Computed:    true,
						},
						"auth_key": {
							Description: "PAN instance auth key. This is only required " +
								"when `panorama_enabled` is set to `true`.",
							Type:     schema.TypeString,
							Optional: true,
						},
						"auth_code": {
							Description: "PAN instance auth code. Only required " +
								"when `license_type` is `BRING_YOUR_OWN`.",
							Type:     schema.TypeString,
							Optional: true,
						},
						"credential_id": {
							Description: "ID of PAN instance credential.",
							Type:        schema.TypeString,
							Computed:    true,
						},
						"global_protect_segment_options": {
							Description: "These options should be set only when global protect is " +
								"enabled on service. These are set per segment. It is expected that " +
								"on a segment where global protect is enabled at least 1 instance " +
								"should be set with portal_enabled and at least one with " +
								"gateway_enabled. It can be on the same instance or a different " +
								"instance under the segment.",
							Type:     schema.TypeSet,
							Optional: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"segment_id": {
										Description: "The segment ID for Global " +
											"Protect options.",
										Type:     schema.TypeString,
										Required: true,
									},
									"portal_enabled": {
										Description: "indicates if the " +
											"GlobalProtect Portal is enabled " +
											"on this PAN instance",
										Type:     schema.TypeBool,
										Required: true,
									},
									"gateway_enabled": {
										Description: "indicates if the Global " +
											"Protect Gateway is enabled on " +
											"this PAN instance",
										Type:     schema.TypeBool,
										Required: true,
									},
									"prefix_list_id": {
										Description: "Prefix List with " +
											"Client IP Pool.",
										Type:     schema.TypeInt,
										Required: true,
									},
								},
							},
						},
						"enable_traffic": {
							Description: "Enable traffic on the PAN instance." + " Default is `true`",
							Default:     true,
							Type:        schema.TypeBool,
							Optional:    true,
						},
					},
				},
			},
			"license_type": {
				Description: "PAN license type, either `BRING_YOUR_OWN` " +
					"or `PAY_AS_YOU_GO`.",
				Type:     schema.TypeString,
				Required: true,
				ValidateFunc: validation.StringInSlice([]string{
					"BRING_YOUR_OWN", "PAY_AS_YOU_GO"}, false),
			},
			"license_sub_type": {
				Description: "PAN sub license type, either `CREDIT_BASED` " +
					"or `MODEL_BASED`. (BETA)",
				Type:     schema.TypeString,
				Optional: true,
				ValidateFunc: validation.StringInSlice([]string{
					"CREDIT_BASED", "MODEL_BASED"}, false),
			},
			"panorama_enabled": {
				Description: "Enable Panorama or not. Default value " +
					"is `false`.",
				Type:     schema.TypeBool,
				Optional: true,
			},
			"panorama_device_group": {
				Description: "Panorama device group.",
				Type:        schema.TypeString,
				Optional:    true,
			},
			"panorama_ip_addresses": {
				Description: "Panorama IP addresses.",
				Type:        schema.TypeList,
				Optional:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
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
			"master_key": {
				Description: "Master Key for PAN instances.",
				Type:        schema.TypeString,
				Optional:    true,
			},
			"master_key_enabled": {
				Description: "Enable Master Key for PAN instances or not. " +
					"It's default to `false`.",
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"master_key_expiry": {
				Description: "PAN Master Key Expiry. The date should be in " +
					"format of `YYYY-MM-DD`, e.g. `2000-01-01`.",
				Type:     schema.TypeString,
				Optional: true,
			},
			"max_instance_count": {
				Description: "Max number of Panorama instances for auto scale.",
				Type:        schema.TypeInt,
				Required:    true,
			},
			"min_instance_count": {
				Description: "Minimal number of Panorama instances for auto " +
					"scale. Default value is `0`.",
				Type:     schema.TypeInt,
				Optional: true,
				Default:  0,
			},
			"registration_pin_id": {
				Description: "PAN Registration PIN ID.",
				Type:        schema.TypeString,
				Required:    true,
			},
			"registration_pin_value": {
				Description: "PAN Registration PIN Value.",
				Type:        schema.TypeString,
				Required:    true,
			},
			"registration_pin_expiry": {
				Description: "PAN Registration PIN Expiry. The date " +
					"should be in format of `YYYY-MM-DD`, e.g. `2000-01-01`.",
				Type:     schema.TypeString,
				Required: true,
			},
			"name": {
				Description: "Name of the PAN service.",
				Type:        schema.TypeString,
				Required:    true,
			},
			"segment_ids": {
				Description: "IDs of segments associated with the service.",
				Type:        schema.TypeSet,
				Required:    true,
				Elem:        &schema.Schema{Type: schema.TypeInt},
			},
			"size": {
				Description: "The size of the service, one of " +
					"`SMALL`, `MEDIUM`, `LARGE`, `2LARGE`",
				Type:     schema.TypeString,
				Required: true,
				ValidateFunc: validation.StringInSlice([]string{
					"SMALL",
					"MEDIUM",
					"LARGE",
					"2LARGE"}, false),
			},
			"tunnel_protocol": {
				Description: "Tunnel Protocol, default to `IPSEC`, " +
					"could be either `IPSEC` or `GRE`.",
				Type:     schema.TypeString,
				Optional: true,
				Default:  "IPSEC",
				ValidateFunc: validation.StringInSlice([]string{
					"IPSEC", "GRE"}, false),
			},
			"type": {
				Description: "The type of the PAN firewall. Either " +
					"'VM-300', 'VM-500' or 'VM-700'",
				Type:     schema.TypeString,
				Optional: true,
				ValidateFunc: validation.StringInSlice([]string{
					"VM-300", "VM-500", "VM-700", "VM-SIM"}, false),
			},
			"version": {
				Description: "The version of the PAN firewall. Please check " +
					"Alkira Portal for all supported versions.",
				Type:     schema.TypeString,
				Required: true,
			},
			"segment_options": {
				Type:        schema.TypeSet,
				Optional:    true,
				Description: "The segment options as used by your PAN firewall.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"segment_id": {
							Description: "The ID of the segment.",
							Type:        schema.TypeString,
							Required:    true,
						},
						"zone_name": {
							Description: "The name of the associated " +
								"firewall zone.",
							Type:     schema.TypeString,
							Required: true,
						},
						"groups": {
							Description: "The list of groups associated " +
								"with the zone.",
							Type:     schema.TypeList,
							Required: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
					},
				},
			},
		},
	}
}

func resourceServicePanCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {

	// INIT
	client := m.(*alkira.AlkiraClient)
	api := alkira.NewServicePan(client)

	// Construct request
	request, err := generateServicePanRequest(d, m)

	if err != nil {
		return diag.FromErr(err)
	}

	// Send create request
	response, provState, err, provErr := api.Create(request)

	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(string(response.Id))

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

	return resourceServicePanRead(ctx, d, m)
}

func resourceServicePanRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {

	// INIT
	client := m.(*alkira.AlkiraClient)
	api := alkira.NewServicePan(client)

	// Get the service
	pan, provState, err := api.GetById(d.Id())

	if err != nil {
		return diag.Diagnostics{{
			Severity: diag.Warning,
			Summary:  "FAILED TO GET RESOURCE",
			Detail:   fmt.Sprintf("%s", err),
		}}
	}

	d.Set("billing_tag_ids", pan.BillingTagIds)
	d.Set("bundle", pan.Bundle)
	d.Set("credential_id", pan.CredentialId)
	d.Set("cxp", pan.CXP)
	d.Set("instance", setPanInstances(d, pan.Instances))
	d.Set("license_type", pan.LicenseType)
	d.Set("license_sub_type", pan.SubLicenseType)
	d.Set("management_segment_id", pan.ManagementSegmentId)
	d.Set("master_key_enabled", pan.MasterKeyEnabled)
	d.Set("max_instance_count", pan.MaxInstanceCount)
	d.Set("min_instance_count", pan.MinInstanceCount)
	d.Set("name", pan.Name)
	d.Set("panorama_enabled", pan.PanoramaEnabled)
	d.Set("segment_ids", pan.SegmentIds)
	d.Set("segment_options", deflateSegmentOptions(pan.SegmentOptions))
	d.Set("size", pan.Size)
	d.Set("tunnel_protocol", pan.TunnelProtocol)
	d.Set("type", pan.Type)
	d.Set("version", pan.Version)

	if pan.PanoramaDeviceGroup != nil {
		d.Set("panorama_device_group", pan.PanoramaDeviceGroup)
	}

	if pan.PanoramaTemplate != nil {
		d.Set("panorama_template", pan.PanoramaTemplate)
	}

	// Set provision state
	if client.Provision == true && provState != "" {
		d.Set("provision_state", provState)
	}

	return nil
}

func resourceServicePanUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {

	// INIT
	client := m.(*alkira.AlkiraClient)
	api := alkira.NewServicePan(client)

	// Construct request
	request, err := generateServicePanRequest(d, m)

	if err != nil {
		return diag.FromErr(err)
	}

	// Send update request
	provState, err, provErr := api.Update(d.Id(), request)

	if err != nil {
		return diag.FromErr(err)
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

	return resourceServicePanRead(ctx, d, m)
}

func resourceServicePanDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {

	// INIT
	client := m.(*alkira.AlkiraClient)
	api := alkira.NewServicePan(client)

	// DELETE
	provState, err, provErr := api.Delete(d.Id())

	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId("")

	// Check provision state
	if client.Provision == true && provState != "SUCCESS" {
		return diag.Diagnostics{{
			Severity: diag.Warning,
			Summary:  "PROVISION (DELETE) FAILED",
			Detail:   fmt.Sprintf("%s", provErr),
		}}
	}

	return nil
}

func generateServicePanRequest(d *schema.ResourceData, m interface{}) (*alkira.ServicePan, error) {

	client := m.(*alkira.AlkiraClient)

	panoramaDeviceGroup := d.Get("panorama_device_group").(string)
	panoramaIpAddresses := convertTypeListToStringList(d.Get("panorama_ip_addresses").([]interface{}))
	panoramaTemplate := d.Get("panorama_template").(string)

	//
	// Construct credentials
	//
	panCredentialId := d.Get("credential_id").(string)

	if 0 == len(panCredentialId) {
		log.Printf("[INFO] Creating PAN Credential")

		panCredName := d.Get("name").(string) + randomNameSuffix()
		panCredential := alkira.CredentialPan{
			Username:   d.Get("pan_username").(string),
			Password:   d.Get("pan_password").(string),
			LicenseKey: d.Get("pan_license_key").(string),
		}

		credentialId, err := client.CreateCredential(
			panCredName,
			alkira.CredentialTypePan,
			panCredential,
			0,
		)

		if err != nil {
			return nil, err
		}

		d.Set("credential_id", credentialId)
	}

	//
	// Construct instances
	//
	instances, err := expandPanInstances(d.Get("instance").([]interface{}), m)

	if err != nil {
		return nil, err
	}

	//
	// Construct segment options
	//
	segmentOptions, err := expandSegmentOptions(d.Get("segment_options").(*schema.Set), m)

	if err != nil {
		return nil, err
	}

	//
	// Construct global protect
	//
	globalProtectSegmentOptions, err := expandGlobalProtectSegmentOptions(d.Get("global_protect_segment_options").(*schema.Set), m)

	if err != nil {
		return nil, err
	}

	// PAN Registration PIN saved as credential
	regCredentialName := d.Get("name").(string) + randomNameSuffix()
	regCredential := alkira.CredentialPanRegistration{
		RegistrationPinId:    d.Get("registration_pin_id").(string),
		RegistrationPinValue: d.Get("registration_pin_value").(string),
	}

	regCredentialExpiry, err := convertInputTimeToEpoch(d.Get("registration_pin_expiry").(string))

	if err != nil {
		log.Printf("[ERROR] failed to parse 'registration_pin_exiry', %v", err)
		return nil, err
	}

	regCredentialId, err := client.CreateCredential(regCredentialName, alkira.CredentialTypePanRegistration, regCredential, regCredentialExpiry)

	if err != nil {
		log.Printf("[ERROR] failed to process PAN registration pin, %v", err)
		return nil, err
	}

	// PAN Master Key saved as credential
	var masterKeyCredentialId string
	if d.Get("master_key_enabled").(bool) {
		masterKeyCredentialName := d.Get("name").(string) + randomNameSuffix()
		masterKeyCredential := alkira.CredentialPanMasterKey{
			MasterKey: d.Get("master_key").(string),
		}

		masterKeyCredentialExpiry, err := convertInputTimeToEpoch(d.Get("master_key_expiry").(string))

		if err != nil {
			log.Printf("[ERROR] failed to parse 'master_key_expiry', %v", err)
			return nil, err
		}

		if masterKeyCredentialExpiry == 0 {
			log.Printf("[ERROR] argument 'master_key_expiry' is required when master key was enabled.")
			return nil, err
		}

		masterKeyCredentialId, err = client.CreateCredential(masterKeyCredentialName, alkira.CredentialTypePanMasterKey, masterKeyCredential, masterKeyCredentialExpiry)

		if err != nil {
			log.Printf("[ERROR] failed to process PAN master key, %v", err)
			return nil, err
		}
	}

	service := &alkira.ServicePan{
		BillingTagIds:               convertTypeSetToIntList(d.Get("billing_tag_ids").(*schema.Set)),
		Bundle:                      d.Get("bundle").(string),
		CXP:                         d.Get("cxp").(string),
		CredentialId:                d.Get("credential_id").(string),
		GlobalProtectEnabled:        d.Get("global_protect_enabled").(bool),
		GlobalProtectSegmentOptions: globalProtectSegmentOptions,
		Instances:                   instances,
		LicenseType:                 d.Get("license_type").(string),
		SubLicenseType:              d.Get("license_sub_type").(string),
		MasterKeyCredentialId:       masterKeyCredentialId,
		MasterKeyEnabled:            d.Get("master_key_enabled").(bool),
		MaxInstanceCount:            d.Get("max_instance_count").(int),
		MinInstanceCount:            d.Get("min_instance_count").(int),
		ManagementSegmentId:         d.Get("management_segment_id").(int),
		Name:                        d.Get("name").(string),
		PanoramaEnabled:             d.Get("panorama_enabled").(bool),
		PanoramaDeviceGroup:         &panoramaDeviceGroup,
		PanoramaIpAddresses:         panoramaIpAddresses,
		PanoramaTemplate:            &panoramaTemplate,
		RegistrationCredentialId:    regCredentialId,
		SegmentOptions:              segmentOptions,
		SegmentIds:                  convertTypeSetToIntList(d.Get("segment_ids").(*schema.Set)),
		TunnelProtocol:              d.Get("tunnel_protocol").(string),
		Size:                        d.Get("size").(string),
		Type:                        d.Get("type").(string),
		Version:                     d.Get("version").(string),
	}

	return service, nil
}
