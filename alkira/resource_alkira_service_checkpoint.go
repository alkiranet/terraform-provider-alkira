package alkira

import (
	"context"
	"fmt"

	"github.com/alkiranet/alkira-client-go/alkira"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func resourceAlkiraCheckpoint() *schema.Resource {
	return &schema.Resource{
		Description:   "Manage checkpoint services",
		CreateContext: resourceCheckpoint,
		ReadContext:   resourceCheckpointRead,
		UpdateContext: resourceCheckpointUpdate,
		DeleteContext: resourceCheckpointDelete,
		CustomizeDiff: func(ctx context.Context, d *schema.ResourceDiff, m interface{}) error {
			client := m.(*alkira.AlkiraClient)

			old, _ := d.GetChange("provision_state")

			if client.Provision && old == "FAILED" {
				d.SetNew("provision_state", "SUCCESS")
			}

			return nil
		},
		Importer: &schema.ResourceImporter{
			StateContext: importWithReadValidation(resourceCheckpointRead),
		},

		Schema: map[string]*schema.Schema{
			"auto_scale": {
				Description: "Indicate if `auto_scale` should be enabled " +
					"for your checkpoint firewall. `ON` and `OFF` are " +
					"accepted values. `OFF` is the default if field is omitted",
				Type:         schema.TypeString,
				Default:      "OFF",
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"ON", "OFF"}, false),
			},
			"billing_tag_ids": {
				Description: "Billing tags to be associated with " +
					"the resource. (see resource `alkira_billing_tag`).",
				Type:     schema.TypeSet,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeInt},
			},
			"cxp": {
				Description: "CXP region.",
				Type:        schema.TypeString,
				Required:    true,
			},
			"provision_state": {
				Description: "The provision state of the resource.",
				Type:        schema.TypeString,
				Computed:    true,
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
				Description: "An array containing properties for each " +
					"Checkpoint Firewall instance that needs to be " +
					"deployed. The number of instances should be equal to " +
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
						"enable_traffic": {
							Description: "Enable traffic on the checkpoint instance. " +
								"Default value is `true`",
							Default:  true,
							Type:     schema.TypeBool,
							Optional: true,
						},
					},
				},
			},
			"license_type": {
				Description: "Checkpoint license type, either " +
					"`BRING_YOUR_OWN` or `PAY_AS_YOU_GO`.",
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
							Description: "ID of Checkpoint Firewall Management server credential.",
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
							Description: "The IDs of the segment to be used to " +
								"access the management server.",
							Type:     schema.TypeString,
							Optional: true,
						},
						"type": {
							Description:  "The type of the management server. either `SMS` or `MDS`.",
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: validation.StringInSlice([]string{"SMS", "MDS"}, false),
						},
						"username": {
							Description: "The username of the management server.",
							Type:        schema.TypeString,
							Optional:    true,
						},
						"password": {
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
				Type:        schema.TypeString,
				Required:    true,
			},
			"segment_options": {
				Description: "The segment options as used by your Checkpoint " +
					"firewall. No more than one segment option will be " +
					"accepted.",
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"segment_id": {
							Description: "The ID of the segment.",
							Type:        schema.TypeString,
							Required:    true,
						},
						"zone_name": {
							Description: "The name of the associated zone. " +
								"Default value is `DEFAULT`.",
							Type:     schema.TypeString,
							Optional: true,
							Default:  "DEFAULT",
						},
						"groups": {
							Description: "The list of Groups associated with the zone.",
							Type:        schema.TypeList,
							Optional:    true,
							Elem:        &schema.Schema{Type: schema.TypeString},
						},
					},
				},
			},
			"size": {
				Description: "The size of the service, one of `SMALL`, " +
					"`MEDIUM`, `LARGE`.",
				Type:     schema.TypeString,
				Required: true,
			},
			"tunnel_protocol": {
				Description: "Tunnel Protocol, default to `IPSEC`, could be " +
					"either `IPSEC` or `GRE`.",
				Type:     schema.TypeString,
				Optional: true,
				Default:  "IPSEC",
				ValidateFunc: validation.StringInSlice([]string{
					"IPSEC", "GRE"}, false),
			},
			"version": {
				Description: "The version of the Checkpoint Firewall. Please " +
					"check all supported versions from Alkira Portal.",
				Type:     schema.TypeString,
				Required: true,
			},
		},
	}
}

func resourceCheckpoint(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {

	client := m.(*alkira.AlkiraClient)
	api := alkira.NewServiceCheckpoint(m.(*alkira.AlkiraClient))

	// Create checkpoint service credentail
	credentialId, err := createCheckpointCredential(d, client)
	if err != nil {
		return diag.FromErr(err)
	}
	d.Set("credential_id", credentialId)

	// Construct request
	request, err := generateCheckpointRequest(d, m)

	if err != nil {
		return diag.FromErr(err)
	}

	// Send create request
	response, provState, err, valErr, provErr := api.Create(request)

	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(string(response.Id))

	if client.Validate && valErr != nil {
		var diags diag.Diagnostics
		readDiags := resourceCheckpointRead(ctx, d, m)
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

	return resourceCheckpointRead(ctx, d, m)
}

func resourceCheckpointRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {

	client := m.(*alkira.AlkiraClient)
	api := alkira.NewServiceCheckpoint(m.(*alkira.AlkiraClient))

	checkpoint, provState, err := api.GetById(d.Id())

	if err != nil {
		return diag.Diagnostics{{
			Severity: diag.Warning,
			Summary:  "FAILED TO GET RESOURCE",
			Detail:   fmt.Sprintf("%s", err),
		}}
	}

	// Get segment
	numOfSegments := len(checkpoint.Segments)

	if numOfSegments == 1 {
		segmentId, err := getSegmentIdByName(checkpoint.Segments[0], m)

		if err != nil {
			return diag.FromErr(err)
		}
		d.Set("segment_id", segmentId)
	} else {
		return diag.FromErr(fmt.Errorf("failed to find segment"))
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
	d.Set("size", checkpoint.Size)
	d.Set("segment_options", deflateSegmentOptions(checkpoint.SegmentOptions))
	d.Set("tunnel_protocol", checkpoint.TunnelProtocol)
	d.Set("version", checkpoint.Version)

	// Set provision state
	if client.Provision && provState != "" {
		d.Set("provision_state", provState)
	}

	return nil
}

func resourceCheckpointUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {

	client := m.(*alkira.AlkiraClient)
	api := alkira.NewServiceCheckpoint(m.(*alkira.AlkiraClient))

	// Update checkpoint service credential
	err := updateCheckpointCredential(d, client)
	if err != nil {
		return diag.FromErr(err)
	}

	// Construct request
	request, err := generateCheckpointRequest(d, m)

	if err != nil {
		return diag.FromErr(err)
	}

	// Send update request
	provState, err, valErr, provErr := api.Update(d.Id(), request)

	if err != nil {
		return diag.FromErr(err)
	}

	if client.Validate && valErr != nil {
		var diags diag.Diagnostics
		readDiags := resourceCheckpointRead(ctx, d, m)
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

	return nil
}

func resourceCheckpointDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {

	client := m.(*alkira.AlkiraClient)
	api := alkira.NewServiceCheckpoint(m.(*alkira.AlkiraClient))

	provState, err, valErr, provErr := api.Delete(d.Id())

	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId("")

	if client.Validate && valErr != nil {
		return diag.Diagnostics{{
			Severity: diag.Error,
			Summary:  "VALIDATION (DELETE) FAILED",
			Detail:   fmt.Sprintf("%s", valErr),
		}}
	}

	// Check provision state
	if client.Provision && provState != "SUCCESS" {
		return diag.Diagnostics{{
			Severity: diag.Warning,
			Summary:  "PROVISION (DELETE) FAILED",
			Detail:   fmt.Sprintf("%s", provErr),
		}}
	}

	return nil
}
