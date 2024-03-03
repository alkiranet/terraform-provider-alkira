package alkira

import (
	"context"
	"fmt"

	"github.com/alkiranet/alkira-client-go/alkira"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func resourceAlkiraServiceCiscoFTDv() *schema.Resource {
	return &schema.Resource{
		Description: "Manage Cisco FTDv Service. (**BETA**)",

		CreateContext: resourceServiceCiscoFTDvCreate,
		ReadContext:   resourceServiceCiscoFTDvRead,
		UpdateContext: resourceServiceCiscoFTDvUpdate,
		DeleteContext: resourceServiceCiscoFTDvDelete,
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
			"name": {
				Description: "The name of the service.",
				Type:        schema.TypeString,
				Required:    true,
			},
			"auto_scale": {
				Description: "Indicate if `auto_scale` should be enabled for your Cisco FTDv service." +
					" `ON` and `OFF` are accepted values. Default is `OFF`.",
				Type:         schema.TypeString,
				Default:      "OFF",
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"ON", "OFF"}, false),
			},
			"provision_state": {
				Description: "The provision state of the resource.",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"size": {
				Description: "The size of the service, one of `SMALL`, " +
					"`MEDIUM`, `LARGE`, `2LARGE`.",
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringInSlice([]string{"SMALL", `MEDIUM`, `LARGE`, `2LARGE`}, false),
			},
			"tunnel_protocol": {
				Description:  "The tunnel protocol. Default is `IPSEC`.",
				Type:         schema.TypeString,
				Optional:     true,
				Default:      "IPSEC",
				ValidateFunc: validation.StringInSlice([]string{"IPSEC"}, false),
			},
			"cxp": {
				Description: "The CXP where the service should be provisioned.",
				Type:        schema.TypeString,
				Required:    true,
			},
			"global_cidr_list_id": {
				Description: "The ID of the `alkira_list_global_cidr` to be " +
					"associated with the service. The list must be tagged " +
					"with `CISCO FTDV`. CIDR must be at least `/25`.",
				Type:     schema.TypeInt,
				Required: true,
			},
			"max_instance_count": {
				Description: "The maximum number of instances that should be deployed.",
				Type:        schema.TypeInt,
				Required:    true,
			},
			"min_instance_count": {
				Description: "The minimum number of instances that should be deployed.",
				Type:        schema.TypeInt,
				Optional:    true,
				Default:     0,
			},
			"billing_tag_ids": {
				Description: "IDs of billing tags.",
				Type:        schema.TypeSet,
				Optional:    true,
				Elem:        &schema.Schema{Type: schema.TypeInt},
			},
			"segment_ids": {
				Description: "IDs of segments associated with the service.",
				Type:        schema.TypeSet,
				Required:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
			"firepower_management_center": {
				Description: "The Firepower Management Center options.",
				Type:        schema.TypeSet,
				Required:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"server_ip": {
							Description: "IP address of the Firepower Management Center.",
							Type:        schema.TypeString,
							Required:    true,
						},
						"credential_id": {
							Description: "An opaque identifier generated when " +
								"storing firepower_management_center credentials.",
							Type:     schema.TypeString,
							Computed: true,
						},
						"username": {
							Description: "Firepower Management Center (FMC) username.",
							Type:        schema.TypeString,
							Required:    true,
						},
						"password": {
							Description: "Firepower Management Center (FMC) password.",
							Type:        schema.TypeString,
							Required:    true,
						},
						"segment_id": {
							Description: "ID of the segment accociated with the " +
								"Firepower Management Center.",
							Type:     schema.TypeString,
							Required: true,
						},
						"ip_allow_list": {
							Description: "List of IP addresses and CIDRs to access the " +
								"Firepower Management Center.",
							Type:     schema.TypeList,
							Optional: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
					},
				},
			},
			"instance": {
				Type:     schema.TypeList,
				Required: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Description: "ID of the Cisco Firepower Firewall instance.",
							Type:        schema.TypeInt,
							Computed:    true,
						},
						"credential_id": {
							Description: "An opaque identifier generated when " +
								"storing Cisco Firepower Firewall instance " +
								"credentials.",
							Type:     schema.TypeString,
							Computed: true,
						},
						"hostname": {
							Description: "Hostname of the Firepower Firewall.",
							Type:        schema.TypeString,
							Required:    true,
						},
						"version": {
							Description: "Cisco Firepower Firewall version.",
							Type:        schema.TypeString,
							Required:    true,
						},
						"license_type": {
							Description: "Cisco Firepower Firewall license " +
								"type, either `BRING_YOUR_OWN` or `PAY_AS_YOU_GO`.",
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: validation.StringInSlice([]string{"BRING_YOUR_OWN", "PAY_AS_YOU_GO"}, false),
						},
						"admin_password": {
							Description: "Firepower Firewall Admin Password.",
							Type:        schema.TypeString,
							Required:    true,
						},
						"fmc_registration_key": {
							Description: "FMC Registration Key.",
							Type:        schema.TypeString,
							Required:    true,
						},
						"ftdv_nat_id": {
							Description: "FTDv NAT ID.",
							Type:        schema.TypeString,
							Optional:    true,
						},
					},
				},
			},
			"segment_options": {
				Type:        schema.TypeSet,
				Optional:    true,
				Description: "The segment options used by the Cisco FTDv.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"segment_id": {
							Description: "ID of the segment.",
							Type:        schema.TypeString,
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
							Optional:    true,
							Elem:        &schema.Schema{Type: schema.TypeString},
						},
					},
				},
			},
		},
	}
}

// resourceServiceCiscoFTDvCreate create a Cisco FTDv service
func resourceServiceCiscoFTDvCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {

	client := m.(*alkira.AlkiraClient)
	api := alkira.NewServiceCiscoFTDv(m.(*alkira.AlkiraClient))

	// Construct request
	request, err := generateServiceCiscoFTDvRequest(d, m)

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

	return resourceServiceCiscoFTDvRead(ctx, d, m)
}

// resourceServiceCiscoFTDvRead get and save a Cisco FTDv services
func resourceServiceCiscoFTDvRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {

	client := m.(*alkira.AlkiraClient)
	api := alkira.NewServiceCiscoFTDv(m.(*alkira.AlkiraClient))

	service, provState, err := api.GetById(d.Id())

	if err != nil {
		return diag.Diagnostics{{
			Severity: diag.Warning,
			Summary:  "FAILED TO GET RESOURCE",
			Detail:   fmt.Sprintf("%s", err),
		}}
	}

	d.Set("auto_scale", service.AutoScale)
	d.Set("billing_tag_ids", service.BillingTags)
	d.Set("credential_id", service.CredentialId)
	d.Set("cxp", service.Cxp)
	d.Set("firepower_management_center", deflateCiscoFTDvManagementServer(service))
	d.Set("global_cidr_list_id", service.GlobalCidrListId)
	d.Set("instance", setCiscoFTDvInstances(d, service.Instances))
	d.Set("max_instance_count", service.MaxInstanceCount)
	d.Set("min_instance_count", service.MinInstanceCount)
	d.Set("name", service.Name)
	d.Set("segment_options", deflateSegmentOptions(service.SegmentOptions))
	d.Set("size", service.Size)
	d.Set("tunnel_protocol", service.TunnelProtocol)

	// Set provision state
	if client.Provision == true && provState != "" {
		d.Set("provision_state", provState)
	}

	return nil
}

// resourceServiceCiscoFTDvUpdate update a Cisco FTDv service
func resourceServiceCiscoFTDvUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {

	client := m.(*alkira.AlkiraClient)
	api := alkira.NewServiceCiscoFTDv(m.(*alkira.AlkiraClient))

	// Construct request
	request, err := generateServiceCiscoFTDvRequest(d, m)

	if err != nil {
		return diag.FromErr(fmt.Errorf("UpdateServiceCiscoFTDv: failed to marshal: %v", err))
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

	return resourceServiceCiscoFTDvRead(ctx, d, m)
}

// resourceServiceCiscoFTDvDelete delete
func resourceServiceCiscoFTDvDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {

	client := m.(*alkira.AlkiraClient)
	api := alkira.NewServiceCiscoFTDv(m.(*alkira.AlkiraClient))

	provState, err, provErr := api.Delete((d.Id()))

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

// generateServiceCiscoFTDvRequest generate a request
func generateServiceCiscoFTDvRequest(d *schema.ResourceData, m interface{}) (*alkira.ServiceCiscoFTDv, error) {

	// Segments
	segmentNames, err := convertSegmentIdsToSegmentNames(d.Get("segment_ids").(*schema.Set), m)

	if err != nil {
		return nil, err
	}

	// Segment Options
	segmentOptions, err := expandCiscoFtdvSegmentOptions(d.Get("segment_options").(*schema.Set), m)

	if err != nil {
		return nil, err
	}

	//
	// Management Server
	//
	// credential_id and ip_allow_list is on top level of the service,
	// but those fields should be part of the management_center.
	//
	credentialId, ipAllowList, managementServer, err := expandCiscoFtdvManagementServer(d.Get("firepower_management_center").(*schema.Set), m)

	if err != nil {
		return nil, err
	}

	//
	// Instances
	//
	instances, err := expandCiscoFTDvInstances(d.Get("instance").([]interface{}), m)

	if err != nil {
		return nil, err
	}

	//
	// Requests
	//
	request := &alkira.ServiceCiscoFTDv{
		Name:             d.Get("name").(string),
		GlobalCidrListId: d.Get("global_cidr_list_id").(int),
		Size:             d.Get("size").(string),
		CredentialId:     credentialId,
		Cxp:              d.Get("cxp").(string),
		ManagementServer: managementServer,
		IpAllowList:      ipAllowList,
		MaxInstanceCount: d.Get("max_instance_count").(int),
		MinInstanceCount: d.Get("min_instance_count").(int),
		Segments:         segmentNames,
		SegmentOptions:   segmentOptions,
		AutoScale:        d.Get("auto_scale").(string),
		TunnelProtocol:   d.Get("tunnel_protocol").(string),
		BillingTags:      convertTypeSetToIntList(d.Get("billing_tag_ids").(*schema.Set)),
		Instances:        instances,
	}

	return request, nil
}
