package alkira

import (
	"context"
	"fmt"

	"github.com/alkiranet/alkira-client-go/alkira"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func resourceAlkiraConnectorGcpInterconnect() *schema.Resource {
	return &schema.Resource{
		Description: "Manage GCP Interconnect.",

		CreateContext: resourceConnectorGcpInterconnectCreate,
		ReadContext:   resourceConnectorGcpInterconnectRead,
		UpdateContext: resourceConnectorGcpInterconnectUpdate,
		DeleteContext: resourceConnectorGcpInterconnectDelete,
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
				Description: "The name of the connector.",
				Type:        schema.TypeString,
				Required:    true,
			},
			"size": {
				Description: "The size of the connector, one of `5XSMALL`,`XSMALL`,`SMALL`, " +
					"`MEDIUM`, `LARGE`.",
				Type:     schema.TypeString,
				Required: true,
			},
			"description": {
				Description: "The description of the connector.",
				Type:        schema.TypeString,
				Optional:    true,
			},
			"cxp": {
				Description:  "The CXP where the connector should be provisioned.",
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringInSlice([]string{"US-EAST1", "US-WEST1", "US-CENTRAL1", "AP-NORTHEAST", "EU-WEST"}, false),
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
			"tunnel_protocol": {
				Description: "The tunnel protocol used by the connector." +
					"Can be one of `GRE`, `IPSEC`, `VXLAN`, `VXLAN_GPE`.",
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringInSlice([]string{"GRE", "IPSEC", "VXLAN", "VXLAN_GPE"}, false),
			},
			"billing_tag_ids": {
				Description: "Billing tags to be associated with " +
					"the resource. (see resource `alkira_billing_tag`).",
				Type:     schema.TypeSet,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeInt},
			},
			"loopback_prefixes": {
				Description: "A list of prefixes that should be " +
					"associated with the connector.",
				Type:     schema.TypeSet,
				Required: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"instance": {
				Description: "A list of instances of the Interconnect",
				Type:        schema.TypeSet,
				Required:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Description: "The ID of the instance.",
							Type:        schema.TypeInt,
							Optional:    true,
							Computed:    true,
						},
						"name": {
							Description: "The name of the instance.",
							Type:        schema.TypeString,
							Required:    true,
						},
						"edge_availibility_domain": {
							Description:  "The Availibility Domain of the instance.",
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: validation.StringInSlice([]string{"AVAILABILITY_DOMAIN_ANY", "AVAILIBILITY_DOMAIN_1", "AVAILBILITY_DOMAIN_2"}, false),
						},
						"candidate_subnets": {
							Description: "A list of candidate subnets to be used by the instance.",
							Type:        schema.TypeSet,
							Elem:        &schema.Schema{Type: schema.TypeString},
						},
						"customer_asn": {
							Description: "The customer ASN.",
							Type:        schema.TypeInt,
							Required:    true,
						},
						"bgp_auth_key": {
							Description: "The BGP MD5 authentication key to authenticate Alkira CXP.",
							Type:        schema.TypeString,
							Optional:    true,
						},
						"gateway_mac_address": {
							Description: "The MAC address of the gateway." +
								"Required if tunnel protocol is VXLAN.",
							Type:     schema.TypeString,
							Optional: true,
						},
						"vni_id": {
							Description: "The VXLAN Network Identifier." +
								"Required if tunnel protocol is VXLAN.",
							Type:     schema.TypeInt,
							Optional: true,
						},
					},
				},
			},

			"segment_options": {
				Description: "Options for each segment associated with the instance.",
				Type:        schema.TypeSet,
				Required:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"instance_name": {
							Description: "The name of the instance.",
							Type:        schema.TypeString,
							Required:    true,
						},
						"segment_id": {
							Description: "The ID of the segment.",
							Type:        schema.TypeString,
							Required:    true,
						},
						"advertise_on_prem_routes": {
							Description: "Advertise on-prem routes. Default is `false`.",
							Default:     false,
							Optional:    true,
							Type:        schema.TypeBool,
						},
						"disable_internet_exit": {
							Description: "Disable access to the internet. Default is `false`.",
							Default:     false,
							Optional:    true,
							Type:        schema.TypeBool,
						},
						"customer_gateways": {
							Description: "The customer gateway associated with the segment.",
							Type:        schema.TypeSet,
							Required:    true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"loopback_ip": {
										Description: "The customer gateway IP address " +
											"which is set as tunnel source",
										Type:     schema.TypeString,
										Required: true,
									},
									"tunnel_count": {
										Description: "Number of tunnels per customer gateway. " +
											"This must be a multiple of 2.",
										Type:     schema.TypeInt,
										Required: true,
									},
								},
							},
						},
					},
				},
			},
			"scale_group_id": {
				Description: "The ID of the scale group associated with the connector.",
				Type:        schema.TypeString,
				Optional:    true,
			},
			"implicit_group_id": {
				Description: "The ID of the implicit group associated with the connector.",
				Type:        schema.TypeInt,
				Optional:    true,
				Computed:    true,
			},
		},
	}
}

func resourceConnectorGcpInterconnectCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {

	// INIT
	client := m.(*alkira.AlkiraClient)
	api := alkira.NewConnectorGcpInterconnect(m.(*alkira.AlkiraClient))

	request, err := generateGcpInterconnectRequest(d, m)

	if err != nil {
		return diag.FromErr(err)
	}

	// CREATE
	response, provState, err, provErr := api.Create(request)

	if err != nil {
		return diag.FromErr(err)
	}

	// Set states
	d.SetId(string(response.Id))

	if client.Provision == true {
		d.Set("provision_state", provState)

		if provErr != nil {
			return diag.Diagnostics{{
				Severity: diag.Warning,
				Summary:  "PROVISION (CREATE) FAILED",
				Detail:   fmt.Sprintf("%s", provErr),
			}}
		}
	}

	return resourceConnectorGcpInterconnectRead(ctx, d, m)
}

func resourceConnectorGcpInterconnectRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {

	// INIT
	client := m.(*alkira.AlkiraClient)
	api := alkira.NewConnectorGcpInterconnect(m.(*alkira.AlkiraClient))

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
	d.Set("tunnel_protocol", connector.TunnelProtocol)
	d.Set("scale_group_id", connector.ScaleGroupId)
	d.Set("billing_tag_ids", connector.BillingTags)
	d.Set("loopback_prefixes", connector.LoopbackPrefixes)
	d.Set("enabled", connector.Enabled)
	d.Set("implicit_group_id", connector.ImplicitGroupId)
	setGcpInterconnectInstance(d, connector, m)

	// Set provision state
	if client.Provision == true && provState != "" {
		d.Set("provision_state", provState)
	}

	return nil
}

func resourceConnectorGcpInterconnectUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {

	// INIT
	client := m.(*alkira.AlkiraClient)
	api := alkira.NewConnectorGcpInterconnect(m.(*alkira.AlkiraClient))

	request, err := generateGcpInterconnectRequest(d, m)

	if err != nil {
		return diag.FromErr(err)
	}

	// UPDATE
	provState, err, provErr := api.Update(d.Id(), request)

	if err != nil {
		return diag.FromErr(err)
	}

	// Set provision state
	if client.Provision == true {
		d.Set("provision_state", provState)

		if provErr != nil {
			return diag.Diagnostics{{
				Severity: diag.Warning,
				Summary:  "PROVISION (UPDATE) FAILED",
				Detail:   fmt.Sprintf("%s", provErr),
			}}
		}
	}

	return resourceConnectorGcpInterconnectRead(ctx, d, m)
}

func resourceConnectorGcpInterconnectDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {

	// INIT
	client := m.(*alkira.AlkiraClient)
	api := alkira.NewConnectorGcpInterconnect(m.(*alkira.AlkiraClient))

	// DELETE
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
