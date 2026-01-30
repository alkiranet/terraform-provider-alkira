package alkira

import (
	"context"
	"fmt"

	"github.com/alkiranet/alkira-client-go/alkira"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func resourceAlkiraBluecat() *schema.Resource {
	return &schema.Resource{
		Description:   "Provide Bluecat service resource (**BETA**).",
		CreateContext: resourceBluecat,
		ReadContext:   resourceBluecatRead,
		UpdateContext: resourceBluecatUpdate,
		DeleteContext: resourceBluecatDelete,
		CustomizeDiff: func(ctx context.Context, d *schema.ResourceDiff, m interface{}) error {
			client := m.(*alkira.AlkiraClient)

			old, _ := d.GetChange("provision_state")

			if client.Provision && old == "FAILED" {
				d.SetNew("provision_state", "SUCCESS")
			}

			return nil
		},
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"bdds_anycast": {
				Type:        schema.TypeSet,
				Optional:    true,
				Description: "Defines the AnyCast configuration for BDDS type instances",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"ips": {
							Description: "The IPs to be used for AnyCast. The IPs used for AnyCast MUST " +
								"NOT overlap the CIDR of `alkira_segment` resource associated with " +
								"the service.",
							Type:     schema.TypeList,
							Optional: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
						"backup_cxps": {
							Description: "The `backup_cxps` to be used when the current " +
								"Bluecat service is not available. It also needs to " +
								"have a configured Bluecat service in order to take advantage of " +
								"this feature. It is NOT required that the `backup_cxps` should have " +
								"a configured Bluecat service before it can be designated as a backup.",
							Type:     schema.TypeList,
							Optional: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
					},
				},
			},
			"edge_anycast": {
				Type:        schema.TypeSet,
				Optional:    true,
				Description: "Defines the AnyCast configuration for EDGE type instances.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"ips": {
							Description: "The IPs to be used for AnyCast. The IPs used for AnyCast MUST " +
								"NOT overlap the CIDR of `alkira_segment` resource associated with " +
								"the service.",
							Type:     schema.TypeList,
							Optional: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
						"backup_cxps": {
							Description: "The `backup_cxps` to be used when the current " +
								"Bluecat service is not available. It also needs to " +
								"have a configured Bluecat service in order to take advantage of " +
								"this feature. It is NOT required that the `backup_cxps` should have " +
								"a configured Bluecat service before it can be designated as a backup.",
							Type:     schema.TypeList,
							Optional: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
					},
				},
			},
			"billing_tag_ids": {
				Description: "Billing tags to be associated with " +
					"the resource. (see resource `alkira_billing_tag`).",
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
				Description: "The description of the Bluecat service.",
				Type:        schema.TypeString,
				Optional:    true,
			},
			"global_cidr_list_id": {
				Description: "The ID of the global cidr list to be " +
					"associated with the Bluecat service.",
				Type:     schema.TypeInt,
				Required: true,
			},
			"provision_state": {
				Description: "The provision state of the resource.",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"instance": {
				Type:     schema.TypeList,
				Required: true,
				Description: "The properties pertaining to each individual " +
					"instance of the Bluecat service.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Description: "The name of the Bluecat instance. This is set to hostname",
							Type:        schema.TypeString,
							Computed:    true,
						},
						"id": {
							Description: "The ID of the Bluecat instance.",
							Type:        schema.TypeInt,
							Computed:    true,
						},
						"bdds_options": {
							Type:        schema.TypeSet,
							Optional:    true,
							Description: "Defines the options required when instance type is BDDS.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"client_id": {
										Description: "The license clientId of the Bluecat BDDS instance.",
										Type:        schema.TypeString,
										Required:    true,
									},
									"activation_key": {
										Description: "The license activationKey of the Bluecat BDDS instance.",
										Type:        schema.TypeString,
										Required:    true,
										Sensitive:   true,
									},
									"license_credential_id": {
										Description: "The license credential ID of the BDDS instance.",
										Type:        schema.TypeString,
										Computed:    true,
									},
									"hostname": {
										Description: "The host name of the instance.",
										Type:        schema.TypeString,
										Required:    true,
									},
									"model": {
										Description: "The model of the Bluecat BDDS instance.",
										Type:        schema.TypeString,
										Required:    true,
									},
									"version": {
										Description: "The version of the Bluecat BDDS instance to be " +
											"used. Please check Alkira Portal for all " +
											"supported versions",
										Type:     schema.TypeString,
										Required: true,
									},
								},
							},
						},
						"edge_options": {
							Type:        schema.TypeSet,
							Optional:    true,
							Description: "Defines the options required when instance type is EDGE.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"config_data": {
										Description: "The Base64 encoded configuration data generated on Bluecat Edge portal.",
										Type:        schema.TypeString,
										Required:    true,
									},
									"credential_id": {
										Description: "The credential ID of the Edge instance.",
										Type:        schema.TypeString,
										Computed:    true,
									},
									"hostname": {
										Description: "The host name of the Edge instance. This " +
											"should match what was configured on the bluecat edge portal.",
										Type:     schema.TypeString,
										Required: true,
									},
									"version": {
										Description: "The version of the Bluecat Edge instance to be " +
											"used. Please check Alkira Portal for all " +
											"supported versions",
										Type:     schema.TypeString,
										Required: true,
									},
								},
							},
						},
						"type": {
							Description: "The type of the Bluecat instance that " +
								"is to be provisioned. The value could be `BDDS`, " +
								"and `EDGE`.",
							Type:     schema.TypeString,
							Required: true,
							ValidateFunc: validation.StringInSlice([]string{
								"BDDS", "EDGE"}, false),
						},
					},
				},
			},
			"license_type": {
				Description: "Bluecat license type, only " +
					"`BRING_YOUR_OWN` is supported right now.",
				Type:     schema.TypeString,
				Required: true,
				ValidateFunc: validation.StringInSlice([]string{
					"BRING_YOUR_OWN"}, false),
			},
			"name": {
				Description: "Name of the Bluecat service.",
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
		},
	}
}

func resourceBluecat(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {

	client := m.(*alkira.AlkiraClient)
	api := alkira.NewServiceBluecat(m.(*alkira.AlkiraClient))

	// Construct request
	request, err := generateBluecatRequest(d, m)

	if err != nil {
		return diag.FromErr(err)
	}

	// Send create request
	response, provState, err, valErr, provErr := api.Create(request)

	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(string(response.Id))

	// Handle validation errors
	if client.Validate && valErr != nil {
		var diags diag.Diagnostics
		// Try to read the resource to preserve any successfully created state
		readDiags := resourceBluecatRead(ctx, d, m)
		if readDiags.HasError() {
			diags = append(diags, readDiags...)
		}

		// Add the validation error
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "VALIDATION (CREATE) FAILED",
			Detail:   fmt.Sprintf("%s", valErr),
		})

		return diags
	}

	// Set provision state
	if client.Provision {
		d.Set("provision_state", provState)

		if provState == "FAILED" {
			return diag.Diagnostics{{
				Severity: diag.Warning,
				Summary:  "PROVISION (CREATE) FAILED",
				Detail:   fmt.Sprintf("%s", provErr),
			}}
		}
	}

	return resourceBluecatRead(ctx, d, m)
}

func resourceBluecatRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {

	client := m.(*alkira.AlkiraClient)
	api := alkira.NewServiceBluecat(m.(*alkira.AlkiraClient))

	bluecat, provState, err := api.GetById(d.Id())

	if err != nil {
		return diag.Diagnostics{{
			Severity: diag.Warning,
			Summary:  "FAILED TO GET RESOURCE",
			Detail:   fmt.Sprintf("%s", err),
		}}
	}

	// Convert segment names to segment IDs
	segmentIds, err := convertSegmentNamesToSegmentIds(bluecat.Segments, m)
	if err != nil {
		return diag.FromErr(err)
	}

	setAllBluecatResourceFields(d, bluecat)
	d.Set("segment_ids", segmentIds)

	// Set provision state
	if client.Provision && provState != "" {
		d.Set("provision_state", provState)
	}

	return nil
}

func resourceBluecatUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {

	// INIT
	client := m.(*alkira.AlkiraClient)
	api := alkira.NewServiceBluecat(m.(*alkira.AlkiraClient))

	// Construct request
	request, err := generateBluecatRequest(d, m)

	if err != nil {
		return diag.FromErr(err)
	}

	// Send update request
	provState, err, valErr, provErr := api.Update(d.Id(), request)

	if err != nil {
		return diag.FromErr(err)
	}

	// Handle validation errors
	if client.Validate && valErr != nil {
		var diags diag.Diagnostics
		// Try to read the resource to preserve current state
		readDiags := resourceBluecatRead(ctx, d, m)
		if readDiags.HasError() {
			diags = append(diags, readDiags...)
		}

		// Add the validation error
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "VALIDATION (UPDATE) FAILED",
			Detail:   fmt.Sprintf("%s", valErr),
		})

		return diags
	}

	// Set provision state
	if client.Provision {
		d.Set("provision_state", provState)

		if provState == "FAILED" {
			return diag.Diagnostics{{
				Severity: diag.Warning,
				Summary:  "PROVISION (UPDATE) FAILED",
				Detail:   fmt.Sprintf("%s", provErr),
			}}
		}
	}

	return resourceBluecatRead(ctx, d, m)
}

func resourceBluecatDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {

	client := m.(*alkira.AlkiraClient)
	api := alkira.NewServiceBluecat(m.(*alkira.AlkiraClient))

	provState, err, valErr, provErr := api.Delete(d.Id())

	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId("")

	// Handle validation errors
	if client.Validate && valErr != nil {
		return diag.Diagnostics{{
			Severity: diag.Error,
			Summary:  "VALIDATION (DELETE) FAILED",
			Detail:   fmt.Sprintf("%s", valErr),
		}}
	}

	if client.Provision && provState != "SUCCESS" {
		return diag.Diagnostics{{
			Severity: diag.Warning,
			Summary:  "PROVISION (DELETE) FAILED",
			Detail:   fmt.Sprintf("%s", provErr),
		}}
	}

	return nil
}

func generateBluecatRequest(d *schema.ResourceData, m interface{}) (*alkira.ServiceBluecat, error) {
	// Parse Instances
	instanceList := d.Get("instance").([]interface{})
	instances, err := expandBluecatInstances(instanceList, m)
	if err != nil {
		return nil, err
	}

	// Parse BDDS Anycast
	bddsAnycast, err := expandBluecatAnycast(d.Get("bdds_anycast").(*schema.Set))
	if err != nil {
		return nil, err
	}

	// Parse Edge Anycast
	edgeAnycast, err := expandBluecatAnycast(d.Get("edge_anycast").(*schema.Set))
	if err != nil {
		return nil, err
	}

	// Convert segment IDs to segment names
	segmentNames, err := convertSegmentIdsToSegmentNames(d.Get("segment_ids").(*schema.Set), m)
	if err != nil {
		return nil, err
	}

	return &alkira.ServiceBluecat{
		BddsAnycast:      *bddsAnycast,
		EdgeAnycast:      *edgeAnycast,
		BillingTags:      convertTypeSetToIntList(d.Get("billing_tag_ids").(*schema.Set)),
		Cxp:              d.Get("cxp").(string),
		Description:      d.Get("description").(string),
		GlobalCidrListId: d.Get("global_cidr_list_id").(int),
		Instances:        instances,
		LicenseType:      d.Get("license_type").(string),
		Name:             d.Get("name").(string),
		Segments:         segmentNames,
		ServiceGroupName: d.Get("service_group_name").(string),
	}, nil
}
