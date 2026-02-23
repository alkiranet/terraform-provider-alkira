package alkira

import (
	"context"
	"fmt"
	"log"

	"github.com/alkiranet/alkira-client-go/alkira"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceAlkiraConnectorAkamaiProlexic() *schema.Resource {
	return &schema.Resource{
		Description: "Manage Connector for Akamai Prolexic.\n\n" +
			"This resource is still under active development. It may have" +
			" further changes in the near future. Today, to use this " +
			"connector, you will need to have onboarded a BYOIP with " +
			"Do Not Advertise set to `true`. Also, the " +
			"segment with public IPs needs to be reported to " +
			"Akamai Representative.",
		CreateContext: resourceConnectorAkamaiProlexicCreate,
		ReadContext:   resourceConnectorAkamaiProlexicRead,
		UpdateContext: resourceConnectorAkamaiProlexicUpdate,
		DeleteContext: resourceConnectorAkamaiProlexicDelete,
		CustomizeDiff: func(ctx context.Context, d *schema.ResourceDiff, m interface{}) error {
			client := m.(*alkira.AlkiraClient)

			old, _ := d.GetChange("provision_state")

			if client.Provision && old == "FAILED" {
				d.SetNew("provision_state", "SUCCESS")
			}

			return nil
		},
		Importer: &schema.ResourceImporter{
			StateContext: importWithReadValidation(resourceConnectorAkamaiProlexicRead),
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
			"akamai_bgp_asn": {
				Description: "The Akamai BGP ASN.",
				Type:        schema.TypeInt,
				Required:    true,
			},
			"akamai_bgp_authentication_key": {
				Description: "The Akamai BGP Authentication Key.",
				Type:        schema.TypeString,
				Required:    true,
			},
			"byoip_options": {
				Description: "BYOIP options.",
				Type:        schema.TypeSet,
				Required:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"byoip_prefix_id": {
							Description: "BYOIP prefix ID.",
							Type:        schema.TypeInt,
							Required:    true,
						},
						"enable_route_advertisement": {
							Description: "Whether enabling route advertisement.",
							Type:        schema.TypeBool,
							Required:    true,
						},
					},
				},
			},
			"billing_tag_ids": {
				Description: "Billing tags to be associated with " +
					"the resource. (see resource `alkira_billing_tag`).",
				Type:     schema.TypeSet,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeInt},
			},
			"credential_id": {
				Description: "The credential ID for storing Akamai BGP " +
					"authentication key.",
				Type:     schema.TypeString,
				Computed: true,
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
			"size": {
				Description: "The size of the connector, one of `SMALL`, `MEDIUM`, " +
					"`LARGE`, `2LARGE`, `5LARGE`.",
				Type:     schema.TypeString,
				Required: true,
			},
			"segment_id": {
				Description: "The ID of segments associated with the connector. " +
					"Currently, only `1` segment is allowed.",
				Type:     schema.TypeString,
				Required: true,
			},
			"tunnel_configuration": {
				Description: "Tunnel Configurations.",
				Type:        schema.TypeSet,
				Required:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"alkira_public_ip": {
							Description: "Alkira public IP.",
							Type:        schema.TypeString,
							Required:    true,
						},
						"tunnel_ips": {
							Description: "Tunnel IPs.",
							Type:        schema.TypeSet,
							Required:    true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"ran_tunnel_ip": {
										Description: "The underlay tunnel IP on the Akamai " +
											"side to be used to configure tunnels between " +
											"the Alkira CXP and the Akamai Prolexic service. " +
											"A RAN (Routed Access Network) is the unit of " +
											"availability for the Route GRE 3.0 service.",
										Type:     schema.TypeString,
										Required: true,
									},
									"alkira_overlay_tunnel_ip": {
										Description: "The overlay IP of the GRE " +
											"tunnel on the Alkira side.",
										Type:     schema.TypeString,
										Required: true,
									},
									"akamai_overlay_tunnel_ip": {
										Description: "The overlay IP of the GRE " +
											"tunnel on the Alkira side.",
										Type:     schema.TypeString,
										Required: true,
									},
								},
							},
						},
					},
				},
			},
		},
	}
}

func resourceConnectorAkamaiProlexicCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {

	client := m.(*alkira.AlkiraClient)
	api := alkira.NewConnectorAkamaiProlexic(m.(*alkira.AlkiraClient))

	// Construct request
	request, err := generateConnectorAkamaiProlexicRequest(d, m)

	if err != nil {
		return diag.FromErr(err)
	}

	// Send create request
	resource, provState, err, valErr, provErr := api.Create(request)

	if err != nil {
		return diag.FromErr(err)
	}

	// Set the state
	d.SetId(string(resource.Id))

	// Handle validation error
	if client.Validate && valErr != nil {
		var diags diag.Diagnostics
		readDiags := resourceConnectorAkamaiProlexicRead(ctx, d, m)
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

	return resourceConnectorAkamaiProlexicRead(ctx, d, m)
}

func resourceConnectorAkamaiProlexicRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {

	client := m.(*alkira.AlkiraClient)
	api := alkira.NewConnectorAkamaiProlexic(m.(*alkira.AlkiraClient))

	// Get resource
	connector, provState, err := api.GetById(d.Id())

	if err != nil {
		return diag.Diagnostics{{
			Severity: diag.Warning,
			Summary:  "FAILED TO GET RESOURCE",
			Detail:   fmt.Sprintf("%s", err),
		}}
	}

	d.Set("akamai_bgp_asn", connector.AkamaiBgpAsn)
	d.Set("billing_tag_ids", connector.BillingTags)
	d.Set("cxp", connector.CXP)
	d.Set("enabled", connector.Enabled)
	d.Set("group", connector.Group)
	d.Set("implicit_group_id", connector.ImplicitGroupId)
	d.Set("name", connector.Name)
	d.Set("size", connector.Size)
	d.Set("description", connector.Description)

	// Get segment
	numOfSegments := len(connector.Segments)
	if numOfSegments == 1 {
		segmentId, err := getSegmentIdByName(connector.Segments[0], m)

		if err != nil {
			return diag.FromErr(err)
		}
		d.Set("segment_id", segmentId)
	} else {
		return diag.FromErr(fmt.Errorf("failed to find segment"))
	}

	// byoip_options
	var options []map[string]interface{}

	for _, option := range connector.ByoipOptions {
		i := map[string]interface{}{
			"byoip_prefix_id":            option.ByoipId,
			"enable_route_advertisement": option.RouteAdvertisementEnabled,
		}
		options = append(options, i)
	}
	d.Set("byoip_options", options)
	// Set provision state
	if client.Provision && provState != "" {
		d.Set("provision_state", provState)
	}

	return nil
}

func resourceConnectorAkamaiProlexicUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {

	client := m.(*alkira.AlkiraClient)
	api := alkira.NewConnectorAkamaiProlexic(m.(*alkira.AlkiraClient))

	// Construct update request
	connector, err := generateConnectorAkamaiProlexicRequest(d, m)

	if err != nil {
		return diag.FromErr(err)
	}

	// Send update request
	provState, err, valErr, provisionErr := api.Update(d.Id(), connector)

	if err != nil {
		return diag.FromErr(err)
	}

	// Handle validation error
	if client.Validate && valErr != nil {
		var diags diag.Diagnostics
		readDiags := resourceConnectorAkamaiProlexicRead(ctx, d, m)
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

		if provisionErr != nil {
			return diag.Diagnostics{{
				Severity: diag.Warning,
				Summary:  "PROVISION (UPDATE) FAILED",
				Detail:   fmt.Sprintf("%s", provisionErr),
			}}
		}
	}

	return resourceConnectorAkamaiProlexicRead(ctx, d, m)
}

func resourceConnectorAkamaiProlexicDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {

	client := m.(*alkira.AlkiraClient)
	api := alkira.NewConnectorAkamaiProlexic(m.(*alkira.AlkiraClient))

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

	if client.Provision && provState != "SUCCESS" {
		return diag.Diagnostics{{
			Severity: diag.Warning,
			Summary:  "PROVISION (DELETE) FAILED",
			Detail:   fmt.Sprintf("%s", provErr),
		}}
	}

	return nil
}

// generateConnectorAkamaiProlexicRequest generate request for the connector
func generateConnectorAkamaiProlexicRequest(d *schema.ResourceData, m interface{}) (*alkira.ConnectorAkamaiProlexic, error) {

	byoipOptions := expandConnectorAkamaiByoipOptions(d.Get("byoip_options").(*schema.Set))
	tunnelConfigurations := expandConnectorAkamaiTunnelConfiguration(d.Get("tunnel_configuration").(*schema.Set))

	// Convert Segment
	segmentName, err := getSegmentNameById(d.Get("segment_id").(string), m)

	if err != nil {
		return nil, err
	}

	// Create implicit akamai-prolexic credential
	log.Printf("[INFO] Creating credential-akamai-prolexic")
	c := alkira.CredentialAkamaiProlexic{
		BgpAuthenticationKey: d.Get("akamai_bgp_authentication_key").(string),
	}

	client := m.(*alkira.AlkiraClient)
	credentialId, err := client.CreateCredential(d.Get("name").(string), alkira.CredentialTypeAkamaiProlexic, c, 0)

	if err != nil {
		return nil, err
	}

	d.Set("credential_id", credentialId)

	connector := &alkira.ConnectorAkamaiProlexic{
		AkamaiBgpAsn:         d.Get("akamai_bgp_asn").(int),
		BillingTags:          convertTypeSetToIntList(d.Get("billing_tag_ids").(*schema.Set)),
		ByoipOptions:         byoipOptions,
		CXP:                  d.Get("cxp").(string),
		CredentialId:         credentialId,
		Group:                d.Get("group").(string),
		Enabled:              d.Get("enabled").(bool),
		Name:                 d.Get("name").(string),
		Segments:             []string{segmentName},
		Size:                 d.Get("size").(string),
		OverlayConfiguration: tunnelConfigurations,
		Description:          d.Get("description").(string),
	}

	return connector, nil
}
