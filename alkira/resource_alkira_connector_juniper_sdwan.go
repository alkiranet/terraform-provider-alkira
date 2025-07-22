package alkira

import (
	"context"
	"fmt"

	"github.com/alkiranet/alkira-client-go/alkira"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func resourceAlkiraConnectorJuniperSdwan() *schema.Resource {
	return &schema.Resource{
		Description:   "Manage Juniper SD-WAN Connector. (**BETA**)",
		CreateContext: resourceConnectorJuniperSdwanCreate,
		ReadContext:   resourceConnectorJuniperSdwanRead,
		UpdateContext: resourceConnectorJuniperSdwanUpdate,
		DeleteContext: resourceConnectorJuniperSdwanDelete,
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
			"version": {
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
			"size": &schema.Schema{
				Description: "The size of the connector, one of `SMALL`, " +
					"`MEDIUM`, `LARGE`, `2LARGE`, `4LARGE`, `5LARGE`.",
				Type:     schema.TypeString,
				Required: true,
				ValidateFunc: validation.StringInSlice([]string{
					"SMALL", "MEDIUM", "LARGE", "2LARGE", "4LARGE", "5LARGE"}, false),
			},
			"instance": &schema.Schema{
				Description: "Juniper SSR Connector Instances",
				Type:        schema.TypeList,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"credential_id": {
							Description: "The generated username password credential ID.",
							Type:        schema.TypeString,
							Computed:    true,
						},
						"username": {
							Description: "The username of the Juniper instance.",
							Type:        schema.TypeString,
							Required:    true,
						},
						"password": {
							Description: "The password of the Juniper instance.",
							Type:        schema.TypeString,
							Required:    true,
						},
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

func resourceConnectorJuniperSdwanCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {

	// INIT
	client := m.(*alkira.AlkiraClient)
	api := alkira.NewConnectorJuniperSdwan(m.(*alkira.AlkiraClient))

	// Construct request
	request, err := generateConnectorJuniperSdwanRequest(d, m)

	if err != nil {
		return diag.FromErr(err)
	}

	// Send create request
	response, provState, err, provErr := api.Create(request)

	if err != nil {
		for _, instance := range request.Instances {
			client.DeleteCredential(instance.CredentialId, alkira.CredentialTypeUserNamePassword)
			client.DeleteCredential(instance.RegistrationKeyCredentialId, alkira.CredentialTypeApiKey)
		}
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

	return resourceConnectorJuniperSdwanRead(ctx, d, m)
}

func resourceConnectorJuniperSdwanRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {

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
	d.Set("version", connector.Version)
	d.Set("name", connector.Name)
	d.Set("size", connector.Size)
	d.Set("tunnel_protocol", connector.TunnelProtocol)
	d.Set("enabled", connector.Enabled)
	d.Set("description", connector.Description)

	// Set Juniper instances
	setJuniperInstances(d, connector)

	// Set VRF mapping
	var mappings []map[string]interface{}

	for _, m := range connector.JuniperSsrVrfMappings {
		mapping := map[string]interface{}{
			"advertise_on_prem_routes": m.AdvertiseOnPremRoutes,
			"disable_internet_exit":    m.DisableInternetExit,
			"juniper_ssr_bgp_asn":      m.JuniperSsrBgpAsn,
			"segment_id":               m.SegmentId,
			"juniper_ssr_vrf_name":     m.JuniperSsrVrfName,
		}
		mappings = append(mappings, mapping)
	}

	d.Set("target_segment", mappings)

	// Set provision state
	if client.Provision == true && provState != "" {
		d.Set("provision_state", provState)
	}

	return nil
}

func resourceConnectorJuniperSdwanUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {

	// INIT
	client := m.(*alkira.AlkiraClient)
	api := alkira.NewConnectorJuniperSdwan(m.(*alkira.AlkiraClient))

	request, err := generateConnectorJuniperSdwanRequest(d, m)

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

	return resourceConnectorJuniperSdwanRead(ctx, d, m)
}

func resourceConnectorJuniperSdwanDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {

	// INIT
	client := m.(*alkira.AlkiraClient)
	api := alkira.NewConnectorJuniperSdwan(m.(*alkira.AlkiraClient))

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

// generateConnectorJuniperSdwanRequest generate request for Juniper SD-WAN connector
func generateConnectorJuniperSdwanRequest(d *schema.ResourceData, m interface{}) (*alkira.ConnectorJuniperSdwan, error) {

	// Expand juniper instances
	instances, err := expandJuniperSdwanInstances(m.(*alkira.AlkiraClient), d.Get("instance").([]interface{}))

	if err != nil {
		return nil, err
	}

	// Construct the request payload
	connector := &alkira.ConnectorJuniperSdwan{
		BillingTags:           convertTypeSetToIntList(d.Get("billing_tag_ids").(*schema.Set)),
		Instances:             instances,
		JuniperSsrVrfMappings: expandJuniperSdwanVrfMappings(d.Get("juniper_ssr_vrf_mapping").(*schema.Set)),
		Version:               d.Get("version").(string),
		Cxp:                   d.Get("cxp").(string),
		Group:                 d.Get("group").(string),
		Name:                  d.Get("name").(string),
		Size:                  d.Get("size").(string),
		TunnelProtocol:        "GRE",
		Enabled:               d.Get("enabled").(bool),
		Description:           d.Get("description").(string),
	}

	return connector, nil
}
