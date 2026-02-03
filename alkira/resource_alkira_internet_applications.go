package alkira

import (
	"context"
	"fmt"

	"github.com/alkiranet/alkira-client-go/alkira"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func resourceAlkiraInternetApplication() *schema.Resource {
	return &schema.Resource{
		Description:   "Manage Internet Application.",
		CreateContext: resourceInternetApplicationCreate,
		ReadContext:   resourceInternetApplicationRead,
		UpdateContext: resourceInternetApplicationUpdate,
		DeleteContext: resourceInternetApplicationDelete,
		CustomizeDiff: func(ctx context.Context, d *schema.ResourceDiff, m interface{}) error {
			client := m.(*alkira.AlkiraClient)

			old, _ := d.GetChange("provision_state")

			if client.Provision && old == "FAILED" {
				d.SetNew("provision_state", "SUCCESS")
			}

			return nil
		},
		Importer: &schema.ResourceImporter{
			StateContext: importWithReadValidation(resourceInternetApplicationRead),
		},

		Schema: map[string]*schema.Schema{
			"billing_tag_ids": {
				Description: "Billing tags to be associated with " +
					"the resource. (see resource `alkira_billing_tag`).",
				Type:     schema.TypeSet,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeInt},
			},
			"bi_directional_az": {
				Description: "Bi-directional IFA AZ. The value could be " +
					"either `AZ0` or `AZ1`",
				Type:     schema.TypeString,
				Optional: true,
			},
			"byoip_id": {
				Description: "BYOIP ID.",
				Type:        schema.TypeInt,
				Optional:    true,
			},
			"connector_id": {
				Description: "Connector ID.",
				Type:        schema.TypeInt,
				Required:    true,
			},
			"connector_type": {
				Description: "Connector Type.The value could be `AWS_VPC`, " +
					"`AZURE_VNET`, `GCP_VPC`, `OCI_VCN`, `SD_WAN`, `IP_SEC` " +
					"`ARUBA_EDGE_CONNECT`, or `EXPRESS_ROUTE`.",
				Type:     schema.TypeString,
				Required: true,
			},
			"fqdn_prefix": {
				Description: "User provided FQDN prefix that will be " +
					"published on AWS Route 53.",
				Type:     schema.TypeString,
				Required: true,
			},
			"group_id": {
				Description: "ID of the auto generated system group.",
				Type:        schema.TypeInt,
				Computed:    true,
			},
			"inbound_connector_id": {
				Description: "Inbound connector ID.",
				Type:        schema.TypeString,
				Optional:    true,
			},
			"inbound_connector_type": {
				Description: "This field defines how the internet application " +
					"to be opened up to the public. Value `DEFAULT` means that " +
					"the native cloud internet connector is used. In this " +
					"case, Alkira takes care of creating this inbound internet " +
					"connector implicitly. When value `AKAMAI_PROLEXIC` is used " +
					"it means that the inbound traffic is through " +
					"`alkira_connector_akamai_prolexic`. You need to create " +
					"and configure that connector and use it with the " +
					"internet application.",
				Type:     schema.TypeString,
				Optional: true,
				Default:  "DEFAULT",
				ValidateFunc: validation.StringInSlice([]string{
					"DEFAULT", "AKAMAI_PROLEXIC"}, false),
			},
			"internet_protocol": {
				Description: "Protocol to be associated with the resource. " +
					"The value could be: `IPV4`, `IPV6` or `BOTH`. In order " +
					"to use the option `IPV6` or `BOTH`, field " +
					"`enable_ipv6_to_ipv4_translation` " +
					"should be enabled on the associated segment " +
					"and a valid IP pool range should be provided. " +
					"`IPV6` and `BOTH` options are only available to Internet " +
					"Applications on AWS CXPs. (**BETA**)",
				Type:     schema.TypeString,
				Optional: true,
				Default:  "IPV4",
				ValidateFunc: validation.StringInSlice([]string{
					"IPV4", "IPV6", "BOTH"}, false),
			},
			"name": {
				Description: "The name of the internet application.",
				Type:        schema.TypeString,
				Required:    true,
			},
			"provision_state": {
				Description: "The provision state of the resource.",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"public_ips": {
				Description: "This option pertains to the `AKAMAI_PROLEXIC` " +
					"`inbound_connector_type`. The public IPs are to be used " +
					"to access the internet application. These public IPs " +
					"must belong to one of the BYOIP ranges configured for " +
					"the connector-akamai-prolexic.",
				Type:     schema.TypeList,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"segment_id": {
				Description: "The ID of segment associated with the internet " +
					"application.",
				Type:     schema.TypeString,
				Required: true,
			},
			"size": {
				Description: "The size of the internet application, one of " +
					"`SMALL`, `MEDIUM` and `LARGE`.",
				Type:     schema.TypeString,
				Required: true,
			},
			"source_nat_ip_pool": {
				Description: "A IP range to be used for source NAT with this " +
					"internet application. It could be only one defined for " +
					"now. The endpoints of each range are inclusive. Source " +
					"NAT can only be used if `inbound_connector_type` is " +
					"`DEFAULT`.",
				Type: schema.TypeSet,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"start_ip": {
							Description: "The start IP of the range.",
							Type:        schema.TypeString,
							Required:    true,
						},
						"end_ip": {
							Description: "The end IP of the range.",
							Type:        schema.TypeString,
							Required:    true,
						},
					},
				},
				Optional: true,
			},
			"ilb_credential_id": {
				Description: "The credential ID of AWS account for `target` " +
					"This field can only be used when `connector_type` is " +
					"`AWS_VPC` and `target`'s `type` is `ILB_NAME`.",
				Type:     schema.TypeString,
				Optional: true,
			},
			"target": {
				Type: schema.TypeSet,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"type": {
							Description: "The type of the target, one of " +
								"`IP` or `ILB_NAME`.",
							Type:     schema.TypeString,
							Required: true,
							ValidateFunc: validation.StringInSlice(
								[]string{"IP", "ILB_NAME"}, false),
						},
						"value": {
							Description: "IFA ILB name or private IP.",
							Type:        schema.TypeString,
							Required:    true,
						},
						"port_ranges": {
							Description: "list of ports or port ranges. " +
								"Values can be mixed i.e. " +
								"`[\"20\", \"100-200\"]`. Value [\"-1\"] " +
								"means any port.",
							Type:     schema.TypeList,
							Elem:     &schema.Schema{Type: schema.TypeString},
							Required: true,
						},
					},
				},
				Required: true,
			},
		},
	}
}

func resourceInternetApplicationCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {

	// INIT
	client := m.(*alkira.AlkiraClient)
	api := alkira.NewInternetApplication(m.(*alkira.AlkiraClient))

	request, err := generateInternetApplicationRequest(d, m)

	if err != nil {
		return diag.FromErr(err)
	}

	// CREATE
	response, provState, err, valErr, provErr := api.Create(request)

	if err != nil {
		return diag.FromErr(err)
	}

	// Handle validation error
	if client.Validate && valErr != nil {
		var diags diag.Diagnostics
		readDiags := resourceInternetApplicationRead(ctx, d, m)
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

		if provErr != nil {
			return diag.Diagnostics{{
				Severity: diag.Warning,
				Summary:  "PROVISION (CREATE) FAILED",
				Detail:   fmt.Sprintf("%s", provErr),
			}}
		}
	}

	d.SetId(string(response.Id))
	return resourceInternetApplicationRead(ctx, d, m)
}

func resourceInternetApplicationRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {

	// INIT
	client := m.(*alkira.AlkiraClient)
	api := alkira.NewInternetApplication(m.(*alkira.AlkiraClient))

	// GET
	app, provState, err := api.GetById(d.Id())

	if err != nil {
		return diag.Diagnostics{{
			Severity: diag.Warning,
			Summary:  "FAILED TO GET RESOURCE",
			Detail:   fmt.Sprintf("%s", err),
		}}
	}

	d.Set("billing_tag_ids", app.BillingTags)
	d.Set("bi_directional_az", app.BiDirectionalAvailabilityZone)
	d.Set("byoip_id", app.ByoipId)
	d.Set("connector_id", app.ConnectorId)
	d.Set("connector_type", app.ConnectorType)
	d.Set("fqdn_prefix", app.FqdnPrefix)
	d.Set("name", app.Name)
	d.Set("internet_protocol", app.InternetProtocol)
	d.Set("public_ips", app.PublicIps)
	d.Set("size", app.Size)
	d.Set("ilb_credential_id", app.IlbCredentialId)

	// Segment
	segmentId, err := getSegmentIdByName(app.SegmentName, m)

	if err != nil {
		return diag.FromErr(err)
	}

	d.Set("segment_id", segmentId)

	// Source NAT IP pool
	var pool []map[string]interface{}

	for _, ipRange := range app.SnatIpv4Ranges {
		i := map[string]interface{}{
			"start_ip": ipRange.StartIp,
			"end_ip":   ipRange.EndIp,
		}
		pool = append(pool, i)
	}

	d.Set("source_nat_ip_pool", pool)

	// targets
	var targets []map[string]interface{}

	for _, target := range app.Targets {
		i := map[string]interface{}{
			"type":       target.Type,
			"value":      target.Value,
			"ports":      target.Ports,
			"portRanges": target.PortRanges,
		}
		targets = append(targets, i)
	}

	d.Set("targets", targets)

	// Set provision state
	if client.Provision && provState != "" {
		d.Set("provision_state", provState)
	}

	return nil
}

func resourceInternetApplicationUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {

	// INIT
	client := m.(*alkira.AlkiraClient)
	api := alkira.NewInternetApplication(m.(*alkira.AlkiraClient))

	request, err := generateInternetApplicationRequest(d, m)

	if err != nil {
		return diag.FromErr(err)
	}

	// UPDATE
	provState, err, valErr, provErr := api.Update(d.Id(), request)

	if err != nil {
		return diag.FromErr(err)
	}

	// Handle validation error
	if client.Validate && valErr != nil {
		var diags diag.Diagnostics
		readDiags := resourceInternetApplicationRead(ctx, d, m)
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

		if provErr != nil {
			return diag.Diagnostics{{
				Severity: diag.Warning,
				Summary:  "PROVISION FAILED",
				Detail:   fmt.Sprintf("%s", provErr),
			}}
		}
	}

	return resourceInternetApplicationRead(ctx, d, m)
}

func resourceInternetApplicationDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {

	// INIT
	client := m.(*alkira.AlkiraClient)
	api := alkira.NewInternetApplication(m.(*alkira.AlkiraClient))

	provState, err, valErr, provErr := api.Delete(d.Id())

	if err != nil {
		// Terraform may not print "with <resource address>" for destroys of objects
		// that are no longer in configuration, so include identifying context here.
		name, _ := d.GetOk("name")
		if nameStr, ok := name.(string); ok && nameStr != "" {
			return diag.FromErr(fmt.Errorf("%w alkira_internet_applications (name=%q id=%s)", err, nameStr, d.Id()))
		}
		return diag.FromErr(fmt.Errorf("%w alkira_internet_applications (id=%s)", err, d.Id()))
	}

	// Handle validation error
	if client.Validate && valErr != nil {
		return diag.Diagnostics{{
			Severity: diag.Error,
			Summary:  "VALIDATION (DELETE) FAILED",
			Detail:   fmt.Sprintf("%s", valErr),
		}}
	}

	d.SetId("")

	if client.Provision && provState != "SUCCESS" {
		return diag.Diagnostics{{
			Severity: diag.Warning,
			Summary:  "PROVISION (DELETE) FAILED",
			Detail:   fmt.Sprintf("%s", provErr),
		}}
	}

	return nil
}

func generateInternetApplicationRequest(d *schema.ResourceData, m interface{}) (*alkira.InternetApplication, error) {

	//
	// Segment
	//
	segmentName, err := getSegmentNameById(d.Get("segment_id").(string), m)

	if err != nil {
		return nil, err
	}

	//
	// Targets
	//
	pool := expandInternetApplicationSourceNatPool(d.Get("source_nat_ip_pool").(*schema.Set))

	//
	// Targets
	//
	targets := expandInternetApplicationTargets(d.Get("target").(*schema.Set))

	// Assemble request
	request := &alkira.InternetApplication{
		BillingTags:                   convertTypeSetToIntList(d.Get("billing_tag_ids").(*schema.Set)),
		BiDirectionalAvailabilityZone: d.Get("bi_directional_az").(string),
		ByoipId:                       d.Get("byoip_id").(int),
		ConnectorId:                   d.Get("connector_id").(int),
		ConnectorType:                 d.Get("connector_type").(string),
		FqdnPrefix:                    d.Get("fqdn_prefix").(string),
		InboundConnectorId:            d.Get("inbound_connector_id").(string),
		InboundConnectorType:          d.Get("inbound_connector_type").(string),
		InternetProtocol:              d.Get("internet_protocol").(string),
		Name:                          d.Get("name").(string),
		PublicIps:                     convertTypeListToStringList(d.Get("public_ips").([]interface{})),
		SegmentName:                   segmentName,
		SnatIpv4Ranges:                pool,
		Size:                          d.Get("size").(string),
		Targets:                       targets,
		IlbCredentialId:               d.Get("ilb_credential_id").(string),
	}

	return request, nil
}
