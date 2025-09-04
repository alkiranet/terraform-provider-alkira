package alkira

import (
	"context"
	"fmt"

	"github.com/alkiranet/alkira-client-go/alkira"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceAlkiraConnectorJuniperSdwan() *schema.Resource {
	return &schema.Resource{
		Description:   "Manage Juniper SD-WAN Connector. (**BETA**)",
		CreateContext: resourceConnectorJuniperSdwanCreate,
		ReadContext:   resourceConnectorJuniperSdwanRead,
		UpdateContext: resourceConnectorJuniperSdwanUpdate,
		DeleteContext: resourceConnectorJuniperSdwanDelete,
		CustomizeDiff: func(ctx context.Context, d *schema.ResourceDiff, m any) error {
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
				Description: "The CXP where the connector should be " +
					"provisioned.",
				Type:     schema.TypeString,
				Required: true,
			},
			"juniper_ssr_version": {
				Description: "The Juniper SSR Version.",
				Type:        schema.TypeString,
				Required:    true,
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
				Description: "The ID of implicit group automaticaly created " +
					"with the connector.",
				Type:     schema.TypeInt,
				Computed: true,
			},
			"availability_zone": {
				Description: "Availability zone of the Juniper instance(s)",
				Type:        schema.TypeInt,
				Required:    true,
			},
			"size": &schema.Schema{
				Description: "The size of the connector, one of `SMALL`, " +
					"`MEDIUM`, `LARGE`, `2LARGE`, `4LARGE`, `5LARGE`.",
				Type:     schema.TypeString,
				Required: true,
			},
			"tunnel_protocol": {
				Description: "The tunnel protocol used by the connector.  Only accepted protocol is 'GRE'",
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "GRE",
				ValidateFunc: func(val any, key string) (warns []string, errs []error) {
					v := val.(string)
					if v != "GRE" {
						errs = append(errs, fmt.Errorf("%q must be GRE, got: %s", key, v))
					}
					return
				},
			},
			"instance": &schema.Schema{
				Description: "Juniper SSR Connector Instances",
				Type:        schema.TypeList,
				MinItems:    1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"registration_key": {
							Description: "The registration key of the Juniper instance.",
							Type:        schema.TypeString,
							Required:    true,
						},
						"registration_key_credential_id": {
							Description: "The generated registration key credential ID.",
							Type:        schema.TypeString,
							Computed:    true,
						},
						"hostname": {
							Description: "The hostname of the Juniper Instance.",
							Type:        schema.TypeString,
							Required:    true,
						},
						"id": {
							Description: "The ID of the Juniper instance.",
							Type:        schema.TypeInt,
							Computed:    true,
						},
					},
				},
				Required: true,
			},
			"juniper_ssr_vrf_mapping": {
				Description: "Juniper SSR Vrf Mapping.",
				Type:        schema.TypeSet,
				MaxItems:    1,
				MinItems:    1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"advertise_on_prem_routes": {
							Description: "Whether advertising On Prem Routes. " +
								"Default value is `false`.",
							Type:     schema.TypeBool,
							Optional: true,
							Default:  false,
						},
						"disable_internet_exit": {
							Description: "Enable or disable access to the " +
								"internet when traffic arrives via this " +
								"connector. Default value is `false`",
							Type:     schema.TypeBool,
							Optional: true,
							Default:  false,
						},
						"segment_id": {
							Description: "Alkira Segment ID.",
							Type:        schema.TypeInt,
							Required:    true,
						},
						"juniper_ssr_bgp_asn": {
							Description: "Gateway BGP ASN. Only accepts '65000'",
							Type:        schema.TypeInt,
							Optional:    true,
							Default:     65000,
							ValidateFunc: func(val any, key string) (warns []string, errs []error) {
								v := val.(int)
								if v != 65000 {
									errs = append(errs, fmt.Errorf("%q must be 65000, got: %d", key, v))
								}
								return
							},
						},
						"juniper_ssr_vrf_name": {
							Description: "Juniper VRF Name. Only accepts 'default'",
							Type:        schema.TypeString,
							Optional:    true,
							Default:     "default",
							ValidateFunc: func(val any, key string) (warns []string, errs []error) {
								v := val.(string)
								if v != "default" {
									errs = append(errs, fmt.Errorf("%q must be default, got: %s", key, v))
								}
								return
							},
						},
					},
				},
				Required: true,
			},
			"enabled": {
				Description: "Is the connector enabled. Default is `true`.",
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     true,
			},
		},
	}
}

func resourceConnectorJuniperSdwanCreate(ctx context.Context, d *schema.ResourceData, m any) diag.Diagnostics {
	// INIT
	client := m.(*alkira.AlkiraClient)
	api := alkira.NewConnectorJuniperSdwan(m.(*alkira.AlkiraClient))

	// Construct request
	request, err := generateConnectorJuniperSdwanRequest(d, m)

	if err != nil {
		return diag.FromErr(err)
	}

	// Send create request
	response, provState, err, valErr, provErr := api.Create(request)

	if err != nil {
		for _, instance := range request.Instances {
			client.DeleteCredential(instance.RegistrationKeyCredentialId, alkira.CredentialTypeApiKey)
		}
		return diag.FromErr(err)
	}

	// Set states
	d.SetId(string(response.Id))

	// Handle validation error
	if client.Validate && valErr != nil {
		var diags diag.Diagnostics
		readDiags := resourceConnectorIPSecRead(ctx, d, m)
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

	return resourceConnectorJuniperSdwanRead(ctx, d, m)
}

func resourceConnectorJuniperSdwanRead(ctx context.Context, d *schema.ResourceData, m any) diag.Diagnostics {
	// INIT
	client := m.(*alkira.AlkiraClient)
	api := alkira.NewConnectorJuniperSdwan(m.(*alkira.AlkiraClient))

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
	d.Set("group", connector.Group)
	d.Set("implicit_group_id", connector.ImplicitGroupId)
	d.Set("juniper_ssr_version", connector.Version)
	d.Set("name", connector.Name)
	d.Set("size", connector.Size)
	d.Set("tunnel_protocol", connector.TunnelProtocol)
	d.Set("enabled", connector.Enabled)
	d.Set("description", connector.Description)
	d.Set("availability_zone", connector.AvailabilityZone)

	// Set Juniper instances
	setJuniperInstances(d, connector)

	// Set VRF mapping
	var mappings []map[string]any

	for _, m := range connector.JuniperSsrVrfMappings {
		mapping := map[string]any{
			"advertise_on_prem_routes": m.AdvertiseOnPremRoutes,
			"disable_internet_exit":    m.DisableInternetExit,
			"juniper_ssr_bgp_asn":      m.JuniperSsrBgpAsn,
			"segment_id":               m.SegmentId,
			"juniper_ssr_vrf_name":     m.JuniperSsrVrfName,
		}
		mappings = append(mappings, mapping)
	}

	d.Set("juniper_ssr_vrf_mapping", mappings)

	// Set provision state
	if client.Provision == true && provState != "" {
		d.Set("provision_state", provState)
	}

	return nil
}

func resourceConnectorJuniperSdwanUpdate(ctx context.Context, d *schema.ResourceData, m any) diag.Diagnostics {

	// INIT
	client := m.(*alkira.AlkiraClient)
	api := alkira.NewConnectorJuniperSdwan(m.(*alkira.AlkiraClient))

	request, err := generateConnectorJuniperSdwanRequest(d, m)

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
		readDiags := resourceConnectorIPSecRead(ctx, d, m)
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

	return resourceConnectorJuniperSdwanRead(ctx, d, m)
}

func resourceConnectorJuniperSdwanDelete(ctx context.Context, d *schema.ResourceData, m any) diag.Diagnostics {

	// INIT
	client := m.(*alkira.AlkiraClient)
	api := alkira.NewConnectorJuniperSdwan(m.(*alkira.AlkiraClient))

	// DELETE
	provState, err, valErr, provErr := api.Delete(d.Id())

	if err != nil {
		return diag.FromErr(err)
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
	if client.Provision == true && provState != "SUCCESS" {
		return diag.Diagnostics{{
			Severity: diag.Warning,
			Summary:  "PROVISION (DELETE) FAILED",
			Detail:   fmt.Sprintf("%s", provErr),
		}}
	}
	return nil
}

// generateConnectorJuniperSdwanRequest generate request for Juniper SD-WAN connector
func generateConnectorJuniperSdwanRequest(d *schema.ResourceData, m any) (*alkira.ConnectorJuniperSdwan, error) {

	// Expand juniper instances
	instances, err := expandJuniperSdwanInstances(m.(*alkira.AlkiraClient), d.Get("instance").([]any))

	if err != nil {
		return nil, err
	}

	// Construct the request payload
	connector := &alkira.ConnectorJuniperSdwan{
		BillingTags:           convertTypeSetToIntList(d.Get("billing_tag_ids").(*schema.Set)),
		Instances:             instances,
		JuniperSsrVrfMappings: expandJuniperSdwanVrfMappings(d.Get("juniper_ssr_vrf_mapping").(*schema.Set)),
		Version:               d.Get("juniper_ssr_version").(string),
		Cxp:                   d.Get("cxp").(string),
		Group:                 d.Get("group").(string),
		Name:                  d.Get("name").(string),
		Size:                  d.Get("size").(string),
		TunnelProtocol:        d.Get("tunnel_protocol").(string),
		Enabled:               d.Get("enabled").(bool),
		Description:           d.Get("description").(string),
		AvailabilityZone:      d.Get("availability_zone").(int),
	}

	return connector, nil
}
