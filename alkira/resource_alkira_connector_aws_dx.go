package alkira

import (
	"context"
	"fmt"

	"github.com/alkiranet/alkira-client-go/alkira"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func resourceAlkiraConnectorAwsDx() *schema.Resource {
	return &schema.Resource{
		Description: "Manage AWS Direct Connect (DX) connector. (**BETA**)",

		CreateContext: resourceConnectorAwsDxCreate,
		ReadContext:   resourceConnectorAwsDxRead,
		UpdateContext: resourceConnectorAwsDxUpdate,
		DeleteContext: resourceConnectorAwsDxDelete,
		CustomizeDiff: func(ctx context.Context, d *schema.ResourceDiff, m interface{}) error {
			client := m.(*alkira.AlkiraClient)

			old, _ := d.GetChange("provision_state")

			if client.Provision && old == "FAILED" {
				d.SetNew("provision_state", "SUCCESS")
			}

			return nil
		},
		Importer: &schema.ResourceImporter{
			StateContext: importWithReadValidation(resourceConnectorAwsDxRead),
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Description: "The name of the connector.",
				Type:        schema.TypeString,
				Required:    true,
			},
			"description": {
				Description: "The description of the connector.",
				Type:        schema.TypeString,
				Optional:    true,
			},
			"cxp": {
				Description: "The CXP where the connector should be " +
					"provisioned.",
				Type:     schema.TypeString,
				Required: true,
			},
			"enabled": {
				Description: "Is the connector enabled. Default is `true`.",
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     true,
			},
			"group": {
				Description: "The group of the connector.",
				Type:        schema.TypeString,
				Optional:    true,
			},
			"implicit_group_id": {
				Description: "ID of implicit group created for the " +
					"connector.",
				Type:     schema.TypeInt,
				Computed: true,
			},
			"tunnel_protocol": {
				Description: "The tunnel protocol used by the connector." +
					"The value should be one of `GRE`, `IPSEC`, `VXLAN`, " +
					"`VXLAN_GPE`.",
				Type:     schema.TypeString,
				Required: true,
				ValidateFunc: validation.StringInSlice([]string{
					"GRE", "IPSEC", "VXLAN", "VXLAN_GPE"}, false),
			},
			"billing_tag_ids": {
				Description: "Billing tags to be associated with " +
					"the resource. (see resource `alkira_billing_tag`).",
				Type:     schema.TypeSet,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeInt},
			},
			"size": {
				Description: "The size of the connector, one of `SMALL`, " +
					"`MEDIUM`, `LARGE`, `2LARGE`, `5LARGE` or `10LARGE`.",
				Type:     schema.TypeString,
				Required: true,
			},
			"scale_group_id": {
				Description: "The ID of the scale group associated with " +
					"the connector.",
				Type:     schema.TypeString,
				Optional: true,
			},
			"instance": {
				Description: "AWS DirectConnect (DX) instance.",
				Type:        schema.TypeList,
				Required:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Description: "The name of the instance.",
							Type:        schema.TypeString,
							Required:    true,
						},
						"id": {
							Description: "ID of the instance.",
							Type:        schema.TypeInt,
							Computed:    true,
						},
						"connection_id": {
							Description: "AWS DirctConnect connection ID.",
							Type:        schema.TypeString,
							Required:    true,
						},
						"dx_asn": {
							Description: "The ASN of AWS side of the " +
								"connection.",
							Type:     schema.TypeInt,
							Required: true,
						},
						"dx_gateway_ip": {
							Description: "Valid IP from underlay_prefix " +
								"network used on AWS Direct Connect gateway.",
							Type:     schema.TypeString,
							Optional: true,
						},
						"on_prem_asn": {
							Description: "The customer underlay ASN.",
							Type:        schema.TypeInt,
							Required:    true,
						},
						"on_prem_gateway_ip": {
							Description: "Valid IP from customer gateway.",
							Type:        schema.TypeString,
							Optional:    true,
						},
						"underlay_prefix": {
							Description: "A `/30` IP prefix for on-premise " +
								"gateway and DirectConnect gateway.",
							Type:     schema.TypeString,
							Optional: true,
						},
						"bgp_auth_key": {
							Description: "The BGP MD5 authentication key for" +
								"Direct Connect Gateway to verify peer.",
							Type:     schema.TypeString,
							Optional: true,
						},
						"bgp_auth_key_alkira": {
							Description: "The BGP MD5 authentication key for" +
								"Alkira to authenticate CXP.",
							Type:     schema.TypeString,
							Optional: true,
						},
						"vlan_id": {
							Description: "ID of customer facing VLAN " +
								"provided by the co-location provider, " +
								"configured for the link between colo " +
								"provider and the customer router.",
							Type:     schema.TypeInt,
							Required: true,
						},
						"aws_region": {
							Description: "AWS region of the Direct Connect.",
							Type:        schema.TypeString,
							Required:    true,
						},
						"credential_id": {
							Description: "ID of AWS credential.",
							Type:        schema.TypeString,
							Required:    true,
						},
						"gateway_mac_address": {
							Description: "The MAC address of the gateway." +
								"It's required if the `tunnel_protocol` " +
								"is `VXLAN`.",
							Type:     schema.TypeString,
							Optional: true,
						},
						"vni": {
							Description: "Customer provided VXLAN Network " +
								"Identifier (VNI). This field is required " +
								"only when `tunnel_protocol` is `VXLAN`.",
							Type:     schema.TypeInt,
							Optional: true,
						},
						"segment_options": {
							Description: "Options for each segment " +
								"associated with the instance.",
							Type:     schema.TypeSet,
							Required: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"segment_id": {
										Description: "The ID of the segment.",
										Type:        schema.TypeString,
										Required:    true,
									},
									"on_prem_segment_asn": {
										Description: "The ASN of customer " +
											"on-prem side.",
										Type:     schema.TypeInt,
										Required: true,
									},
									"customer_loopback_ip": {
										Description: "Customer loopback IP " +
											"which is set as tunnel source. " +
											"The field is applicable only " +
											"when `tunnel_protocol` is not " +
											"`IPSEC`.",
										Type:     schema.TypeString,
										Optional: true,
									},
									"alkira_loopback_ip1": {
										Description: "Alkira loopback IP " +
											"which is set as tunnel 1. " +
											"The field is applicable only " +
											"when `tunnel_protocol` is not " +
											"`IPSEC`.",
										Type:     schema.TypeString,
										Optional: true,
									},
									"alkira_loopback_ip2": {
										Description: "Alkira loopback IP " +
											"which is set as tunnel 2. " +
											"The field is applicable only " +
											"when `tunnel_protocol` is not " +
											"`IPSEC`.",
										Type:     schema.TypeString,
										Optional: true,
									},
									"loopback_subnet": {
										Description: "Prefix of all loopback " +
											"IPs, helps to identify the block " +
											"to reserve IPs from.",
										Type:     schema.TypeString,
										Required: true,
									},
									"advertise_on_prem_routes": {
										Description: "Advertise on-prem routes. " +
											"Default value is `false`.",
										Type:     schema.TypeBool,
										Optional: true,
										Default:  false,
									},
									"advertise_default_route": {
										Description: "Enable or disable access " +
											"to the internet when traffic " +
											"arrives via this connector. " +
											"Default value is `false`.",
										Type:     schema.TypeBool,
										Optional: true,
										Default:  false,
									},
									"number_of_customer_loopback_ips": {
										Description: "The number of customer " +
											"loopback IPs needs to be generated " +
											"by Alkira from `loopback_subnet`." +
											"The field is only applicable " +
											"when `tunnel_protocol` is `IPSEC`.",
										Type:     schema.TypeInt,
										Optional: true,
									},
									"tunnel_count_per_customer_loopback_ip": {
										Description: "The number of tunnels " +
											"needs to be created for each " +
											"customer loopback IP. The value " +
											"must be multiple of `2` (one " +
											"tunnel per AZ). The field is only " +
											"applicable when `tunnel_protocol` " +
											"is `IPSEC`.",
										Type:     schema.TypeInt,
										Optional: true,
									},
								},
							},
						}, // segment_options
					},
				},
			}, // instances
		},
	}
}

func resourceConnectorAwsDxCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {

	// INIT
	client := m.(*alkira.AlkiraClient)
	api := alkira.NewConnectorAwsDirectConnect(m.(*alkira.AlkiraClient))

	request, err := generateAwsDirectConnectRequest(d, m)

	if err != nil {
		return diag.FromErr(err)
	}

	// CREATE
	response, provState, err, valErr, provErr := api.Create(request)

	if err != nil {
		return diag.FromErr(err)
	}

	// Set states
	d.SetId(string(response.Id))

	// Handle validation error
	if client.Validate && valErr != nil {
		var diags diag.Diagnostics
		readDiags := resourceConnectorAwsDxRead(ctx, d, m)
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

	return resourceConnectorAwsDxRead(ctx, d, m)
}

func resourceConnectorAwsDxRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {

	// INIT
	client := m.(*alkira.AlkiraClient)
	api := alkira.NewConnectorAwsDirectConnect(m.(*alkira.AlkiraClient))

	// GET
	connector, provState, err := api.GetById(d.Id())

	if err != nil {
		return diag.Diagnostics{{
			Severity: diag.Warning,
			Summary:  "FAILED TO GET RESOURCE",
			Detail:   fmt.Sprintf("%s", err),
		}}
	}

	d.Set("name", connector.Name)
	d.Set("description", connector.Description)
	d.Set("cxp", connector.Cxp)
	d.Set("group", connector.Group)
	d.Set("size", connector.Size)
	d.Set("scale_group_id", connector.ScaleGroupId)
	d.Set("tunnel_protocol", connector.TunnelProtocol)
	d.Set("billing_tag_ids", connector.BillingTags)
	d.Set("enabled", connector.Enabled)
	d.Set("implicit_group_id", connector.ImplicitGroupId)

	// Set instances
	err = setAwsDirectConnectInstance(d, m, connector)

	if err != nil {
		return diag.Diagnostics{{
			Severity: diag.Warning,
			Summary:  "FAILED TO GET RESOURCE",
			Detail:   fmt.Sprintf("%s", err),
		}}
	}

	// Set provision state
	if client.Provision && provState != "" {
		d.Set("provision_state", provState)
	}

	return nil
}

func resourceConnectorAwsDxUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {

	// INIT
	client := m.(*alkira.AlkiraClient)
	api := alkira.NewConnectorAwsDirectConnect(m.(*alkira.AlkiraClient))

	request, err := generateAwsDirectConnectRequest(d, m)

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
		readDiags := resourceConnectorAwsDxRead(ctx, d, m)
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
				Summary:  "PROVISION (UPDATE) FAILED",
				Detail:   fmt.Sprintf("%s", provErr),
			}}
		}
	}

	return resourceConnectorAwsDxRead(ctx, d, m)
}

func resourceConnectorAwsDxDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {

	// INIT
	client := m.(*alkira.AlkiraClient)
	api := alkira.NewConnectorAwsDirectConnect(m.(*alkira.AlkiraClient))

	// DELETE
	provState, err, valErr, provErr := api.Delete(d.Id())

	if err != nil {
		// Terraform may not print "with <resource address>" for destroys of objects
		// that are no longer in configuration, so include identifying context here.
		name, _ := d.GetOk("name")
		if nameStr, ok := name.(string); ok && nameStr != "" {
			return diag.FromErr(fmt.Errorf("%w alkira_connector_aws_dx (name=%q id=%s)", err, nameStr, d.Id()))
		}
		return diag.FromErr(fmt.Errorf("%w alkira_connector_aws_dx (id=%s)", err, d.Id()))
	}

	d.SetId("")

	// Handle validation error
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
