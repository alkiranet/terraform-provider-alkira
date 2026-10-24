package alkira

import (
	"context"
	"fmt"

	"github.com/alkiranet/alkira-client-go/alkira"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func resourceAlkiraConnectorMerakiSDWAN() *schema.Resource {
	return &schema.Resource{
		Description:   "Manage Cisco Meraki SD-WAN Connector.",
		CreateContext: resourceConnectorMerakiSDWANCreate,
		ReadContext:   resourceConnectorMerakiSDWANRead,
		UpdateContext: warnOnFailedStateUpdate(resourceConnectorMerakiSDWANUpdate),
		DeleteContext: resourceConnectorMerakiSDWANDelete,
		CustomizeDiff: func(ctx context.Context, d *schema.ResourceDiff, m interface{}) error {
			client := m.(*alkira.AlkiraClient)

			old, _ := d.GetChange("provision_state")

			if client.Provision && old == "FAILED" {
				d.SetNew("provision_state", "SUCCESS")
			}

			return nil
		},
		Importer: &schema.ResourceImporter{
			StateContext: importWithReadValidation(resourceConnectorMerakiSDWANRead),
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
			"billing_tag_ids": {
				Description: "Billing tags to be associated with " +
					"the resource. (see resource `alkira_billing_tag`).",
				Type:     schema.TypeSet,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeInt},
			},
			"cxp": {
				Description: "The Azure CXP where the connector should be " +
					"provisioned.",
				Type:     schema.TypeString,
				Required: true,
			},
			"group": {
				Description: "The group of the connector.",
				Type:        schema.TypeString,
				Optional:    true,
			},
			"provision_state": {
				Description: "The provision state of the connector.",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"implicit_group_id": {
				Description: "The ID of the implicit group automatically " +
					"created with the connector.",
				Type:     schema.TypeInt,
				Computed: true,
			},
			"size": {
				Description: "The size of the connector, one of `SMALL`, " +
					"`MEDIUM` or `LARGE` (vMX-S, vMX-M or vMX-L). A new " +
					"size applies only to instances added afterwards.",
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringInSlice([]string{"SMALL", "MEDIUM", "LARGE"}, false),
			},
			"version": {
				Description: "The vMX image version. Default is the " +
					"version Alkira supports for the size.",
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"tunnel_protocol": {
				Description: "The tunnel protocol. Only `IPSEC` is supported. " +
					"Default value is `IPSEC`.",
				Type:         schema.TypeString,
				Optional:     true,
				Default:      "IPSEC",
				ValidateFunc: validation.StringInSlice([]string{"IPSEC"}, false),
			},
			"enabled": {
				Description: "Is the connector enabled. Default is `true`.",
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     true,
			},
			"instances": {
				Description: "Meraki vMX instances. Use 2 or more for HA. " +
					"An instance is identified by its `host_name`.",
				Type:     schema.TypeSet,
				Required: true,
				MinItems: 1,
				Set:      merakiSDWANInstanceHash,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Description: "The ID of the instance.",
							Type:        schema.TypeInt,
							Computed:    true,
						},
						"name": {
							Description: "The name of the instance.",
							Type:        schema.TypeString,
							Computed:    true,
						},
						"host_name": {
							Description: "The host name of the instance. It " +
								"must be unique in the connector and " +
								"identifies the instance: a new host name " +
								"replaces the vMX.",
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: validation.StringIsNotEmpty,
						},
						"auth_token": {
							Description: "The vMX authentication token " +
								"generated in the Meraki Dashboard. It is " +
								"stored as a credential of the instance. " +
								"A change replaces the token and redeploys " +
								"the vMX.",
							Type:         schema.TypeString,
							Required:     true,
							Sensitive:    true,
							ValidateFunc: validation.StringIsNotEmpty,
						},
						"credential_id": {
							Description: "The ID of the generated credential.",
							Type:        schema.TypeString,
							Computed:    true,
						},
						"az": {
							Description: "The availability zone Alkira placed " +
								"the instance in.",
							Type:     schema.TypeInt,
							Computed: true,
						},
						"wan_ip": {
							Description: "The vMX WAN private IP.",
							Type:        schema.TypeString,
							Computed:    true,
						},
						"public_ip": {
							Description: "The vMX public IP.",
							Type:        schema.TypeString,
							Computed:    true,
						},
					},
				},
			},
			"target_segment": {
				Description: "The segment of the connector.",
				Type:        schema.TypeList,
				Required:    true,
				MinItems:    1,
				MaxItems:    1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"segment_id": {
							Description: "Alkira Segment ID.",
							Type:        schema.TypeInt,
							Required:    true,
						},
						"meraki_bgp_asn": {
							Description: "BGP ASN of the vMX, from 1 to " +
								"4294967295 (2147483647 on 32-bit builds " +
								"of the provider). It must differ from the " +
								"segment ASN.",
							Type:         schema.TypeInt,
							Required:     true,
							ValidateFunc: validateMerakiSDWANBgpAsn,
						},
						"advertise_on_prem_routes": {
							Description: "Whether to advertise on-premises " +
								"routes. Default value is `false`.",
							Type:     schema.TypeBool,
							Optional: true,
							Default:  false,
						},
						"advertise_default_route": {
							Description: "Whether to advertise the default " +
								"route of the internet connector. Default " +
								"value is `false`.",
							Type:     schema.TypeBool,
							Optional: true,
							Default:  false,
						},
					},
				},
			},
		},
	}
}

func resourceConnectorMerakiSDWANCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {

	// INIT
	client := m.(*alkira.AlkiraClient)
	api := alkira.NewConnectorMerakiSDWAN(m.(*alkira.AlkiraClient))

	// Construct request
	request, created, _, err := generateConnectorMerakiSDWANRequest(d, m, nil)

	if err != nil {
		return diag.FromErr(err)
	}

	// Send create request
	response, provState, err, valErr, provErr := api.Create(request)

	if err != nil {
		deleteMerakiSDWANCredentials(client, created)
		return diag.FromErr(err)
	}

	// Set states
	d.SetId(string(response.Id))

	if client.Validate && valErr != nil {
		var diags diag.Diagnostics
		readDiags := resourceConnectorMerakiSDWANRead(ctx, d, m)
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

		if provState == "FAILED" {
			return diag.Diagnostics{{
				Severity: diag.Warning,
				Summary:  "PROVISION (CREATE) FAILED",
				Detail:   fmt.Sprintf("%s", provErr),
			}}
		}
	}

	return resourceConnectorMerakiSDWANRead(ctx, d, m)
}

func resourceConnectorMerakiSDWANRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {

	// INIT
	client := m.(*alkira.AlkiraClient)
	api := alkira.NewConnectorMerakiSDWAN(m.(*alkira.AlkiraClient))

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
	d.Set("description", connector.Description)
	d.Set("enabled", connector.Enabled)
	d.Set("group", connector.Group)
	d.Set("implicit_group_id", connector.ImplicitGroupId)
	d.Set("name", connector.Name)
	d.Set("size", connector.Size)
	d.Set("tunnel_protocol", connector.TunnelProtocol)
	d.Set("version", connector.Version)

	// Set instances
	setMerakiSDWANInstances(d, connector)

	// Set segment
	var segments []map[string]interface{}

	for _, o := range connector.SegmentOptions {
		segment := map[string]interface{}{
			"advertise_on_prem_routes": o.AdvertiseOnPremRoutes,
			"advertise_default_route":  !o.DisableInternetExit,
			"meraki_bgp_asn":           int(o.MerakiBgpAsn),
			"segment_id":               o.SegmentId,
		}
		segments = append(segments, segment)
	}

	d.Set("target_segment", segments)

	// Set provision state
	if client.Provision && provState != "" {
		d.Set("provision_state", provState)
	}

	return nil
}

func resourceConnectorMerakiSDWANUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {

	// INIT
	client := m.(*alkira.AlkiraClient)
	api := alkira.NewConnectorMerakiSDWAN(m.(*alkira.AlkiraClient))

	// Keep the old state if the update fails, so that the next plan
	// retries it. Otherwise a changed auth_token is saved and the token
	// replace is lost.
	d.Partial(true)

	// Changed tokens are replaced before the update, so that it
	// redeploys with them
	oldInstances, _ := d.GetChange("instances")
	request, created, replaced, err := generateConnectorMerakiSDWANRequest(d, m, oldInstances.(*schema.Set).List())

	if err != nil {
		return diag.FromErr(merakiSDWANReplacedTokensError(replaced, err))
	}

	// UPDATE
	provState, err, valErr, provErr := api.Update(d.Id(), request)

	if err != nil {
		deleteMerakiSDWANCredentials(client, created)
		return diag.FromErr(merakiSDWANReplacedTokensError(replaced, err))
	}

	// The connector is saved before it is validated, so keep the new
	// state from here on
	d.Partial(false)

	if client.Validate && valErr != nil {
		var diags diag.Diagnostics
		readDiags := resourceConnectorMerakiSDWANRead(ctx, d, m)
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

	return resourceConnectorMerakiSDWANRead(ctx, d, m)
}

func resourceConnectorMerakiSDWANDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {

	// INIT
	client := m.(*alkira.AlkiraClient)
	api := alkira.NewConnectorMerakiSDWAN(m.(*alkira.AlkiraClient))

	// DELETE
	_, err, valErr, provErr := api.Delete(d.Id())

	if err != nil {
		// Terraform may not print "with <resource address>" for destroys of objects
		// that are no longer in configuration, so include identifying context here.
		name, _ := d.GetOk("name")
		if nameStr, ok := name.(string); ok && nameStr != "" {
			return diag.FromErr(fmt.Errorf("%w alkira_connector_meraki_sdwan (name=%q id=%s)", err, nameStr, d.Id()))
		}
		return diag.FromErr(fmt.Errorf("%w alkira_connector_meraki_sdwan (id=%s)", err, d.Id()))
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
	if client.Provision && provErr != nil {
		return diag.Diagnostics{{
			Severity: diag.Warning,
			Summary:  "PROVISION (DELETE) FAILED",
			Detail:   fmt.Sprintf("%s", provErr),
		}}
	}

	return nil
}

// generateConnectorMerakiSDWANRequest generate request for Meraki SD-WAN
// connector. old is the instances block before the change, nil on create.
// It also returns the IDs of the credentials created for new instances and
// the instances whose token was replaced.
func generateConnectorMerakiSDWANRequest(d *schema.ResourceData, m interface{}, old []interface{}) (*alkira.ConnectorMerakiSDWAN, []string, []string, error) {

	instances, created, replaced, err := expandMerakiSDWANInstances(
		m.(*alkira.AlkiraClient),
		d.Get("instances").(*schema.Set).List(),
		old,
	)

	if err != nil {
		return nil, nil, replaced, err
	}

	// Construct the request payload
	connector := &alkira.ConnectorMerakiSDWAN{
		BillingTags:    convertTypeSetToIntList(d.Get("billing_tag_ids").(*schema.Set)),
		Cxp:            d.Get("cxp").(string),
		Description:    d.Get("description").(string),
		Enabled:        d.Get("enabled").(bool),
		Group:          d.Get("group").(string),
		Instances:      instances,
		Name:           d.Get("name").(string),
		SegmentOptions: expandMerakiSDWANSegmentOptions(d.Get("target_segment").([]interface{})),
		Size:           d.Get("size").(string),
		TunnelProtocol: d.Get("tunnel_protocol").(string),
		Version:        d.Get("version").(string),
	}

	return connector, created, replaced, nil
}
