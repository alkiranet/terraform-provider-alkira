package alkira

import (
	"context"
	"fmt"
	"log"

	"github.com/alkiranet/alkira-client-go/alkira"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func resourceAlkiraConnectorCiscoSdwan() *schema.Resource {
	return &schema.Resource{
		Description:   "Manage Cisco SD-WAN Connector.",
		CreateContext: resourceConnectorCiscoSdwanCreate,
		ReadContext:   resourceConnectorCiscoSdwanRead,
		UpdateContext: resourceConnectorCiscoSdwanUpdate,
		DeleteContext: resourceConnectorCiscoSdwanDelete,
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
			"billing_tag_ids": {
				Description: "A list of Billing Tag by ID associated with the connector.",
				Type:        schema.TypeList,
				Optional:    true,
				Elem:        &schema.Schema{Type: schema.TypeInt},
			},
			"cxp": {
				Description: "The CXP where the connector should be provisioned.",
				Type:        schema.TypeString,
				Required:    true,
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
				Description: "The ID of implicit group automaticaly created with the connector.",
				Type:        schema.TypeInt,
				Computed:    true,
			},
			"type": {
				Description: "The type of Cisco SD-WAN. Default value is `VEDGE`.",
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "VEDGE",
			},
			"size": &schema.Schema{
				Description:  "The size of the connector, one of `SMALL`, `MEDIUM` and `LARGE`, `2LARGE`, `4LARGE`, `5LARGE`, `10LARGE` and `20LARGE`.",
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringInSlice([]string{"SMALL", "MEDIUM", "LARGE", "2LARGE", "4LARGE", "5LARGE", "10LARGE", "20LARGE"}, false),
			},
			"vedge": &schema.Schema{
				Description: "Cisco vEdge",
				Type:        schema.TypeList,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"cloud_init_file": {
							Description: "The cloud-init file for the vEdge.",
							Type:        schema.TypeString,
							Required:    true,
						},
						"credential_id": {
							Description: "The generated credential ID for Cisco SD-WAN.",
							Type:        schema.TypeString,
							Computed:    true,
						},
						"credential_ssh_key_pair_id": {
							Description: "The ID of the credential for SSH Key Pair.",
							Type:        schema.TypeString,
							Optional:    true,
						},
						"hostname": {
							Description: "The hostname of the vEdge.",
							Type:        schema.TypeString,
							Required:    true,
						},
						"id": {
							Description: "The ID of the vEdge instance.",
							Type:        schema.TypeInt,
							Computed:    true,
						},
						"username": &schema.Schema{
							Description: "Cisco SD-WAN username. It could be also " +
								"set by environment variable `AK_CISCO_SDWAN_USERNAME`.",
							Type:     schema.TypeString,
							Required: true,
							DefaultFunc: schema.EnvDefaultFunc(
								"AK_CISCO_SDWAN_USERNAME",
								nil),
						},
						"password": &schema.Schema{
							Description: "Cisco SD-WAN password. It could be also " +
								"set by environment variable `AK_CISCO_SDWAN_PASSWORD`.",
							Type:     schema.TypeString,
							Required: true,
							DefaultFunc: schema.EnvDefaultFunc(
								"AK_CISCO_SDWAN_PASSWORD",
								nil),
						},
					},
				},
				Required: true,
			},
			"version": {
				Description: "The version of Cisco SD-WAN.",
				Type:        schema.TypeString,
				Required:    true,
			},
			"vrf_segment_mapping": {
				Description: "Specify target segment for VRF.",
				Type:        schema.TypeSet,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"advertise_on_prem_routes": {
							Description: "Advertise On Prem Routes.",
							Type:        schema.TypeBool,
							Optional:    true,
							Default:     false,
						},
						"allow_nat_exit": {
							Description: "Allow NAT exit.",
							Type:        schema.TypeBool,
							Optional:    true,
							Default:     true,
						},
						"customer_asn": {
							Description: "BGP ASN on the customer premise side. A typical value for 2 byte segment " +
								"is `64523` and `4200064523` for 4 byte segment.",
							Type:     schema.TypeInt,
							Required: true,
						},
						"segment_id": {
							Description: "Segment ID.",
							Type:        schema.TypeInt,
							Required:    true,
						},
						"vrf_id": {
							Description: "VRF ID.",
							Type:        schema.TypeInt,
							Required:    true,
						},
					},
				},
				Required: true,
			},
		},
	}
}

func resourceConnectorCiscoSdwanCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {

	// INIT
	client := m.(*alkira.AlkiraClient)
	api := alkira.NewConnectorCiscoSdwan(m.(*alkira.AlkiraClient))

	request, err := generateConnectorCiscoSdwanRequest(d, m)

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

	return resourceConnectorCiscoSdwanRead(ctx, d, m)
}

func resourceConnectorCiscoSdwanRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {

	// INIT
	client := m.(*alkira.AlkiraClient)
	api := alkira.NewConnectorCiscoSdwan(m.(*alkira.AlkiraClient))

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
	d.Set("name", connector.Name)
	d.Set("size", connector.Size)
	d.Set("type", connector.Type)

	// Set vedge
	setVedge(d, connector)

	// Set vrf_segment_mapping
	var mappings []map[string]interface{}

	for _, m := range connector.CiscoEdgeVrfMappings {
		mapping := map[string]interface{}{
			"advertise_on_prem_routes": m.AdvertiseOnPremRoutes,
			"allow_nat_exit":           m.DisableInternetExit,
			"customer_asn":             m.CustomerAsn,
			"segment_id":               m.SegmentId,
			"vrf_id":                   m.Vrf,
		}
		mappings = append(mappings, mapping)
	}

	d.Set("vrf_segment_mapping", mappings)
	d.Set("version", connector.Version)

	// Set provision state
	if client.Provision == true && provState != "" {
		d.Set("provision_state", provState)
	}

	return nil
}

func resourceConnectorCiscoSdwanUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {

	// INIT
	client := m.(*alkira.AlkiraClient)
	api := alkira.NewConnectorCiscoSdwan(m.(*alkira.AlkiraClient))

	request, err := generateConnectorCiscoSdwanRequest(d, m)

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

	return resourceConnectorCiscoSdwanRead(ctx, d, m)
}

func resourceConnectorCiscoSdwanDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {

	// INIT
	client := m.(*alkira.AlkiraClient)
	api := alkira.NewConnectorCiscoSdwan(m.(*alkira.AlkiraClient))

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

// generateConnectorCiscoSdwanRequest generate request for Cisco SD-WAN connector
func generateConnectorCiscoSdwanRequest(d *schema.ResourceData, m interface{}) (*alkira.ConnectorCiscoSdwan, error) {

	// Expand Cisco SDWAN vEdge block
	vedges, err := expandCiscoSdwanVedges(m.(*alkira.AlkiraClient), d.Get("vedge").([]interface{}))

	if err != nil {
		return nil, err
	}

	// Construct the request payload
	connector := &alkira.ConnectorCiscoSdwan{
		BillingTags:          convertTypeListToIntList(d.Get("billing_tag_ids").([]interface{})),
		CiscoEdgeInfo:        vedges,
		CiscoEdgeVrfMappings: expandCiscoSdwanVrfMappings(d.Get("vrf_segment_mapping").(*schema.Set)),
		Cxp:                  d.Get("cxp").(string),
		Group:                d.Get("group").(string),
		Enabled:              d.Get("enabled").(bool),
		Name:                 d.Get("name").(string),
		Size:                 d.Get("size").(string),
		Type:                 d.Get("type").(string),
		Version:              d.Get("version").(string),
	}

	return connector, nil
}

// expandCiscoSdwanVrfMappings expand Cisco SD-WAN VRF segment mapping
func expandCiscoSdwanVrfMappings(in *schema.Set) []alkira.CiscoSdwanEdgeVrfMapping {

	if in == nil || in.Len() == 0 {
		log.Printf("[DEBUG] Empty vrf_segment_mapping")
		return []alkira.CiscoSdwanEdgeVrfMapping{}
	}

	mappings := make([]alkira.CiscoSdwanEdgeVrfMapping, in.Len())
	for i, mapping := range in.List() {
		r := alkira.CiscoSdwanEdgeVrfMapping{}
		t := mapping.(map[string]interface{})

		if v, ok := t["advertise_on_prem_routes"].(bool); ok {
			r.AdvertiseOnPremRoutes = v
		}
		if v, ok := t["allow_nat_exit"].(bool); ok {
			r.DisableInternetExit = v
		}
		if v, ok := t["customer_asn"].(int); ok {
			r.CustomerAsn = v
		}
		if v, ok := t["segment_id"].(int); ok {
			r.SegmentId = v
		}
		if v, ok := t["vrf_id"].(int); ok {
			r.Vrf = v
		}

		mappings[i] = r
	}

	return mappings
}

// expandCiscoSdwanVedges expand Cisco SD-WAN Edge
func expandCiscoSdwanVedges(ac *alkira.AlkiraClient, in []interface{}) ([]alkira.CiscoSdwanEdgeInfo, error) {

	if in == nil || len(in) == 0 {
		log.Printf("[DEBUG] Empty vedges")
		return []alkira.CiscoSdwanEdgeInfo{}, nil
	}

	mappings := make([]alkira.CiscoSdwanEdgeInfo, len(in))

	for i, mapping := range in {
		r := alkira.CiscoSdwanEdgeInfo{}
		t := mapping.(map[string]interface{})

		var username string
		var password string

		if v, ok := t["hostname"].(string); ok {
			r.HostName = v
		}
		if v, ok := t["cloud_init_file"].(string); ok {
			r.CloudInitFile = v
		}
		if v, ok := t["username"].(string); ok {
			username = v
		}
		if v, ok := t["password"].(string); ok {
			password = v
		}
		if v, ok := t["credential_ssh_key_pair_id"].(string); ok {
			r.SshKeyPairCredentialId = v
		}
		if v, ok := t["credential_id"].(string); ok {
			if v == "" {
				log.Printf("[INFO] Creating CISCO-SDWAN Credential")
				credentialName := r.HostName + randomNameSuffix()

				credential := alkira.CredentialCiscoSdwan{
					Username: username,
					Password: password,
				}

				credentialId, err := ac.CreateCredential(
					credentialName,
					alkira.CredentialTypeCiscoSdwan,
					credential,
					0,
				)

				if err != nil {
					return nil, err
				}

				r.CredentialId = credentialId
			} else {
				r.CredentialId = v
			}
		}
		if v, ok := t["id"].(int); ok {
			r.Id = v
		}

		mappings[i] = r
	}

	return mappings, nil
}
