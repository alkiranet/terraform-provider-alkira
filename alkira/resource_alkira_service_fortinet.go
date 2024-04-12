package alkira

import (
	"context"
	"fmt"
	"strconv"

	"github.com/alkiranet/alkira-client-go/alkira"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func resourceAlkiraServiceFortinet() *schema.Resource {
	return &schema.Resource{
		Description:   "Manage Fortinet firewall.",
		CreateContext: resourceFortinetCreate,
		ReadContext:   resourceFortinetRead,
		UpdateContext: resourceFortinetUpdate,
		DeleteContext: resourceFortinetDelete,
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
			"auto_scale": {
				Description: "Whether enable auto scale for Fortinet firewall. " +
					"It could be either `ON` and `OFF`. Default value is `OFF`.",
				Type:     schema.TypeString,
				Optional: true,
				Default:  "OFF",
				ValidateFunc: validation.StringInSlice([]string{
					"ON", "OFF"}, false),
			},
			"billing_tag_ids": {
				Description: "IDs of billing tags to associate with " +
					"the service.",
				Type:     schema.TypeSet,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeInt},
			},
			"credential_id": {
				Description: "ID of Fortinet Firewall credential managed " +
					"by credential resource.",
				Type:     schema.TypeString,
				Computed: true,
			},
			"cxp": {
				Description: "The CXP where the service should be provisioned.",
				Type:        schema.TypeString,
				Required:    true,
			},
			"provision_state": {
				Description: "The provision state of the resource.",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"password": {
				Description: "Fortinet password.",
				Type:        schema.TypeString,
				Optional:    true,
			},
			"username": {
				Description: "Fortinet username.",
				Type:        schema.TypeString,
				Optional:    true,
			},
			"instances": {
				Type:     schema.TypeList,
				Required: true,
				Description: "An array containing properties for each Fortinet " +
					"Firewall instance that needs to be deployed. The number of " +
					"instances should be equal to `max_instance_count`.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Description: "The name of the Fortinet Firewall " +
								"instance.",
							Type:     schema.TypeString,
							Optional: true,
						},
						"license_key_file_path": {
							Description: "Fortinet license key file path. The path to the desired " +
								"license key. \n\n\nThere are two options for providing the required " +
								"license key for Fortinet instance credentials. You can either input " +
								"the value directly into the `license_key` field or provide the file " +
								"path for the license key file using the `license_key_file_path`. " +
								"Either `license_key` or `license_key_file_path` must have an input. " +
								"If both are provided, the Alkira provider will treat the `license_key` " +
								"field with precedence.",
							Type:     schema.TypeString,
							Optional: true,
						},
						"license_key": {
							Description: "The Fortinet license key literal. You may copy and " +
								"paste the contents of your license key here. You may also use " +
								"terraform's built in `file` helper function as a literal input " +
								"for `license_key`. Ex: `license_key = file('/path/to/license/file')`" +
								"the `file` helper function will copy the contents of your file " +
								"and place them as literal data into your configuration. \n\n\n" +
								"Instead of using this field you may also use `license_key_file_path`" +
								"to simply place the path to the license key file you'd like to use. ",
							Type:     schema.TypeString,
							Optional: true,
						},

						"serial_number": {
							Description: "The serial_number of the Fortinet Firewall instance. " +
								"Required only when `license_type` is `BRING_YOUR_OWN.",
							Type:     schema.TypeString,
							Optional: true,
						},
						"credential_id": {
							Description: "The ID of the Fortinet Firewall instance credentials. " +
								"Required only when `license_type` is `BRING_YOUR_OWN`.",
							Type:     schema.TypeString,
							Computed: true,
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
				Description: "Fortinet license type, either `BRING_YOUR_OWN`" +
					"or `PAY_AS_YOU_GO`.",
				Type:     schema.TypeString,
				Required: true,
				ValidateFunc: validation.StringInSlice([]string{
					"BRING_YOUR_OWN", "PAY_AS_YOU_GO"},
					false,
				),
			},
			"license_scheme": {
				Description: "The license scheme tells more about BYOL license " +
					"method. `POINT_BASED` scheme refers to FortiFlex license " +
					"whereas `TERM_BASED` refers to regular BYOL.",
				Type:     schema.TypeString,
				Optional: true,
				Default:  "TERM_BASED",
				ValidateFunc: validation.StringInSlice([]string{
					"POINT_BASED", "TERM_BASED"},
					false,
				),
			},
			"management_server_ip": {
				Description: "The IP addresses used to access the management " +
					"server.",
				Type:     schema.TypeString,
				Optional: true,
			},
			"management_server_segment_id": {
				Description: "The segment ID used to access the management " +
					"server. This segment must be present in the list of " +
					"segments assigned to this Fortinet Firewall service.",
				Type:     schema.TypeString,
				Required: true,
			},
			"max_instance_count": {
				Description: "The maximum number of Fortinet Firewall instances " +
					"that should be deployed. `max_instance_count` must be " +
					"greater than or equal to `min_instance_count`.",
				Type:     schema.TypeInt,
				Required: true,
			},
			"min_instance_count": {
				Description: "The minimum number of Fortinet Firewall instances " +
					"that should be deployed.",
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
				Type:        schema.TypeSet,
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
							Type:        schema.TypeString,
							Required:    true,
						},
						"zone_name": {
							Description: "The name of the associated zone.",
							Type:        schema.TypeString,
							Required:    true,
						},
						"groups": {
							Description: "The list of groups associated with " +
								"the zone.",
							Type:     schema.TypeList,
							Required: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
					},
				},
			},
			"size": {
				Description: "The size of the service, one of `SMALL`, " +
					"`MEDIUM`, `LARGE`.",
				Type:     schema.TypeString,
				Required: true,
				ValidateFunc: validation.StringInSlice([]string{
					"SMALL", "MEDIUM", "LARGE"}, false),
			},
			"tunnel_protocol": {
				Description: "Tunnel Protocol. The default value is `IPSEC`. " +
					"it could be either `IPSEC` or `GRE`.",
				Type:     schema.TypeString,
				Optional: true,
				Default:  "IPSEC",
				ValidateFunc: validation.StringInSlice([]string{
					"IPSEC", "GRE"}, false),
			},
			"version": {
				Description: "The version of the Fortinet Firewall. Please " +
					"check Alkira Portal for all supported versions.",
				Type:     schema.TypeString,
				Required: true,
			},
		},
	}
}

func resourceFortinetCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {

	client := m.(*alkira.AlkiraClient)
	api := alkira.NewServiceFortinet(m.(*alkira.AlkiraClient))

	// Construct request
	request, err := generateFortinetRequest(d, m)

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

	return resourceFortinetRead(ctx, d, m)
}

func resourceFortinetRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {

	client := m.(*alkira.AlkiraClient)
	api := alkira.NewServiceFortinet(m.(*alkira.AlkiraClient))

	f, provState, err := api.GetById(d.Id())

	if err != nil {
		return diag.Diagnostics{{
			Severity: diag.Warning,
			Summary:  "FAILED TO GET RESOURCE",
			Detail:   fmt.Sprintf("%s", err),
		}}
	}

	d.Set("auto_scale", f.AutoScale)
	d.Set("billing_tag_ids", f.BillingTags)
	d.Set("credential_id", f.CredentialId)
	d.Set("cxp", f.Cxp)
	d.Set("license_type", f.LicenseType)
	d.Set("license_scheme", f.Scheme)
	d.Set("management_server_ip", f.ManagementServer.IpAddress)
	d.Set("max_instance_count", f.MaxInstanceCount)
	d.Set("min_instance_count", f.MinInstanceCount)
	d.Set("name", f.Name)
	d.Set("segment_options", deflateSegmentOptions(f.SegmentOptions))
	d.Set("size", f.Size)
	d.Set("tunnel_protocol", f.TunnelProtocol)
	d.Set("version", f.Version)

	// Set management server segment
	managementServerSegmentId, err := getSegmentIdByName(f.ManagementServer.Segment, m)

	if err != nil {
		return diag.FromErr(err)
	}
	d.Set("management_server_segment_id", managementServerSegmentId)

	// Set segments
	segments := make([]int, len(f.Segments))

	for _, seg := range f.Segments {
		seg, err := getSegmentIdByName(seg, m)

		if err != nil {
			return diag.FromErr(err)
		}
		segId, _ := strconv.Atoi(seg)
		segments = append(segments, segId)
	}
	d.Set("segment_ids", segments)

	// Set instances
	setInstance(d, f)

	// Set provision state
	if client.Provision == true && provState != "" {
		d.Set("provision_state", provState)
	}

	return nil
}

func resourceFortinetUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {

	client := m.(*alkira.AlkiraClient)
	api := alkira.NewServiceFortinet(m.(*alkira.AlkiraClient))

	// Construct request
	request, err := generateFortinetRequest(d, m)

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

	return nil
}

func resourceFortinetDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {

	client := m.(*alkira.AlkiraClient)
	api := alkira.NewServiceFortinet(m.(*alkira.AlkiraClient))

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
