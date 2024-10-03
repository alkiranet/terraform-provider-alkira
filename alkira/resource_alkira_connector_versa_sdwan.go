package alkira

import (
	"context"
	"fmt"

	"github.com/alkiranet/alkira-client-go/alkira"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func resourceAlkiraConnectorVersaSdwan() *schema.Resource {
	return &schema.Resource{
		Description:   "Manage Versa SD-WAN Connector. (**BETA**)",
		CreateContext: resourceConnectorVersaSdwanCreate,
		ReadContext:   resourceConnectorVersaSdwanRead,
		UpdateContext: resourceConnectorVersaSdwanUpdate,
		DeleteContext: resourceConnectorVersaSdwanDelete,
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
				Description: "Billing tags to be associated with " +
					"the resource. (see resource `alkira_billing_tag`).",
				Type:     schema.TypeSet,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeInt},
			},
			"cxp": {
				Description: "The CXP where the connector should be provisioned.",
				Type:        schema.TypeString,
				Required:    true,
			},
			"enabled": {
				Description: "Is the connector enabled. Default value is `true`.",
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
				Description: "The ID of implicit group automaticaly " +
					"created with the connector.",
				Type:     schema.TypeInt,
				Computed: true,
			},
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
			"provision_state": {
				Description: "The provision state of the connector.",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"global_tenant_id": {
				Description: "The global tenant ID of Versa SD-WAN. Default " +
					"value is `1`.",
				Type:     schema.TypeInt,
				Optional: true,
				Default:  1,
			},
			"local_id": {
				Description: "The local ID.",
				Type:        schema.TypeString,
				Required:    true,
			},
			"local_public_shared_key": {
				Description: "The local public shared key. Default value is" +
					"`1234`.",
				Type:     schema.TypeString,
				Optional: true,
				Default:  "1234",
			},
			"remote_id": {
				Description: "The remote ID.",
				Type:        schema.TypeString,
				Required:    true,
			},
			"remote_public_shared_key": {
				Description: "The remote public shared key. Default value is" +
					"`1234`.",
				Type:     schema.TypeString,
				Optional: true,
				Default:  "1234",
			},
			"size": &schema.Schema{
				Description: "The size of the connector, one of `SMALL`, " +
					"`MEDIUM`, `LARGE`, `2LARGE`, `5LARGE`. ",
				Type:     schema.TypeString,
				Required: true,
				ValidateFunc: validation.StringInSlice([]string{
					"SMALL", "MEDIUM", "LARGE", "2LARGE",
					"5LARGE"}, false),
			},
			"tunnel_protocol": {
				Description: "The tunnel protocol of Versa SD-WAN.",
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "IPSEC",
			},
			"versa_controller_host": {
				Description: "The Versa controller IP/FQDN.",
				Type:        schema.TypeString,
				Required:    true,
			},
			"versa_vos_device": &schema.Schema{
				Description: "Versa VOS Device.",
				Type:        schema.TypeList,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"hostname": {
							Description: "The hostname of the VOS Device.",
							Type:        schema.TypeString,
							Required:    true,
						},
						"id": {
							Description: "The ID of the VOS device.",
							Type:        schema.TypeInt,
							Computed:    true,
						},
						"local_device_serial_number": &schema.Schema{
							Description: "Local device serial number.",
							Type:        schema.TypeString,
							Required:    true,
						},
						"version": &schema.Schema{
							Description: "Versa version.",
							Type:        schema.TypeString,
							Required:    true,
						},
					},
				},
				Required: true,
			},
			"vrf_segment_mapping": {
				Description: "Specify target segment for VRF.",
				Type:        schema.TypeSet,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"advertise_on_prem_routes": {
							Description: "Advertise On Prem Routes. Default value " +
								"is `false`.",
							Type:     schema.TypeBool,
							Optional: true,
							Default:  false,
						},
						"advertise_default_route": {
							Description: "Whether advertise default route of " +
								"internet connector. Default value is `false`.",
							Type:     schema.TypeBool,
							Optional: true,
							Default:  false,
						},
						"versa_bgp_asn": {
							Description: "BGP ASN on the Versa. A typical value " +
								"for 2 byte segment is `64523` and `4200064523` " +
								"for 4 byte segment.",
							Type:     schema.TypeInt,
							Required: true,
						},
						"segment_id": {
							Description: "Segment ID.",
							Type:        schema.TypeInt,
							Required:    true,
						},
						"vrf_name": {
							Description: "VRF Name.",
							Type:        schema.TypeString,
							Required:    true,
						},
					},
				},
				Required: true,
			},
		},
	}
}

func resourceConnectorVersaSdwanCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {

	// INIT
	client := m.(*alkira.AlkiraClient)
	api := alkira.NewConnectorVersaSdwan(m.(*alkira.AlkiraClient))

	request, err := generateConnectorVersaSdwanRequest(d, m)

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

	return resourceConnectorVersaSdwanRead(ctx, d, m)
}

func resourceConnectorVersaSdwanRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {

	// INIT
	client := m.(*alkira.AlkiraClient)
	api := alkira.NewConnectorVersaSdwan(m.(*alkira.AlkiraClient))

	// GET
	connector, provState, err := api.GetById(d.Id())

	if err != nil {
		return diag.Diagnostics{{
			Severity: diag.Warning,
			Summary:  "FAILED TO GET RESOURCE",
			Detail:   fmt.Sprintf("%s", err),
		}}
	}

	d.Set("billing_tag_ids", connector.BillingTags)
	d.Set("cxp", connector.Cxp)
	d.Set("enabled", connector.Enabled)
	d.Set("group", connector.Group)
	d.Set("implicit_group_id", connector.ImplicitGroupId)
	d.Set("local_id", connector.LocalId)
	d.Set("local_public_shared_key", connector.LocalPublicSharedKey)
	d.Set("name", connector.Name)
	d.Set("remote_id", connector.RemoteId)
	d.Set("remote_public_shared_key", connector.RemotePublicSharedKey)
	d.Set("size", connector.Size)
	d.Set("tunnel_protocol", connector.TunnelProtocol)
	d.Set("description", connector.Description)

	// Set Instances
	setVersaSdwanInstance(d, connector)

	// Set vrf_segment_mapping
	var mappings []map[string]interface{}

	for _, m := range connector.VersaSdWanVRFMappings {
		mapping := map[string]interface{}{
			"advertise_on_prem_routes": m.AdvertiseOnPremRoutes,
			"advertise_default_route":  !m.DisableInternetExit,
			"versa_bgp_asn":            m.GatewayBgpAsn,
			"segment_id":               m.SegmentId,
			"vrf_name":                 m.VrfName,
		}
		mappings = append(mappings, mapping)
	}

	d.Set("vrf_segment_mapping", mappings)

	// Set provision state
	if client.Provision == true && provState != "" {
		d.Set("provision_state", provState)
	}

	return nil
}

func resourceConnectorVersaSdwanUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {

	// INIT
	client := m.(*alkira.AlkiraClient)
	api := alkira.NewConnectorVersaSdwan(m.(*alkira.AlkiraClient))

	request, err := generateConnectorVersaSdwanRequest(d, m)

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

	return resourceConnectorVersaSdwanRead(ctx, d, m)
}

func resourceConnectorVersaSdwanDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {

	// INIT
	client := m.(*alkira.AlkiraClient)
	api := alkira.NewConnectorVersaSdwan(m.(*alkira.AlkiraClient))

	// DELETE
	provState, err, provErr := api.Delete(d.Id())

	if err != nil {
		return diag.FromErr(err)
	}

	if client.Provision == true && provState != "SUCCESS" {
		return diag.Diagnostics{{
			Severity: diag.Warning,
			Summary:  "PROVISION (DELETE) FAILED",
			Detail:   fmt.Sprintf("%s", provErr),
		}}
	}

	d.SetId("")

	return nil
}
