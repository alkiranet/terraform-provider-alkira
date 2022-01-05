package alkira

import (
	"log"

	"github.com/alkiranet/alkira-client-go/alkira"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func resourceAlkiraConnectorCiscoSdwan() *schema.Resource {
	return &schema.Resource{
		Description: "Manage Cisco SD-WAN Connector.",
		Create:      resourceConnectorCiscoSdwanCreate,
		Read:        resourceConnectorCiscoSdwanRead,
		Update:      resourceConnectorCiscoSdwanUpdate,
		Delete:      resourceConnectorCiscoSdwanDelete,

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
			"group": {
				Description: "The group of the connector.",
				Type:        schema.TypeString,
				Optional:    true,
			},
			"type": {
				Description: "The type of Cisco SD-WAN. Default value is `VEDGE`.",
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "VEDGE",
			},
			"size": &schema.Schema{
				Description:  "The size of the connector. one of `SMALL`, `MEDIUM` and `LARGE`.",
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringInSlice([]string{"SMALL", "MEDIUM", "LARGE"}, false),
			},
			"vedge": &schema.Schema{
				Description: "Cisco vEdge",
				Type:        schema.TypeSet,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"hostname": {
							Description: "The hostname of the vEdge.",
							Type:        schema.TypeString,
							Required:    true,
						},
						"cloud_init_file": {
							Description: "The cloud-init file for the vEdge.",
							Type:        schema.TypeString,
							Required:    true,
						},
						"credential_id": {
							Description: "The ID of the credential for Cisco SD-WAN.",
							Type:        schema.TypeString,
							Required:    true,
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

func resourceConnectorCiscoSdwanCreate(d *schema.ResourceData, m interface{}) error {
	client := m.(*alkira.AlkiraClient)
	connector, err := generateConnectorCiscoSdwanRequest(client, d, m)

	if err != nil {
		return err
	}

	log.Printf("[INFO] Creating Connector (Cisco SD-WAN)")
	id, err := client.CreateConnectorCiscoSdwan(connector)

	if err != nil {
		return err
	}

	d.SetId(id)

	return resourceConnectorCiscoSdwanRead(d, m)
}

func resourceConnectorCiscoSdwanRead(d *schema.ResourceData, m interface{}) error {
	client := m.(*alkira.AlkiraClient)

	connector, err := client.GetConnectorCiscoSdwan(d.Id())

	if err != nil {
		return err
	}

	d.Set("billing_tag_ids", connector.BillingTags)
	d.Set("cxp", connector.Cxp)
	d.Set("group", connector.Group)
	d.Set("name", connector.Name)
	d.Set("size", connector.Size)
	d.Set("type", connector.Type)

	// Set vedge
	var vedge []map[string]interface{}

	for _, info := range connector.CiscoEdgeInfo {
		edge := map[string]interface{}{
			"hostname":        info.HostName,
			"cloud_init_file": info.CloudInitFile,
			"credential_id":   info.CredentialId,
		}
		vedge = append(vedge, edge)
	}

	d.Set("vedge", vedge)

	// Set vrf_segment_mapping
	var mappings []map[string]interface{}

	for _, m := range connector.CiscoEdgeVrfMappings {
		mapping := map[string]interface{}{
			"advertise_on_prem_routes": m.AdvertiseOnPremRoutes,
			"allow_nat_exit":           m.DisableInternetExit,
			"segment_id":               m.SegmentId,
			"vrf_id":                   m.Vrf,
		}
		mappings = append(mappings, mapping)
	}

	d.Set("vrf_segment_mapping", vedge)
	d.Set("version", connector.Version)

	return nil
}

func resourceConnectorCiscoSdwanUpdate(d *schema.ResourceData, m interface{}) error {
	client := m.(*alkira.AlkiraClient)

	connector, err := generateConnectorCiscoSdwanRequest(client, d, m)

	if err != nil {
		return err
	}

	log.Printf("[INFO] Updating Connector (Cisco SD-WAN) %s", d.Id())
	err = client.UpdateConnectorCiscoSdwan(d.Id(), connector)

	if err != nil {
		return err
	}

	return resourceConnectorCiscoSdwanRead(d, m)
}

func resourceConnectorCiscoSdwanDelete(d *schema.ResourceData, m interface{}) error {
	client := m.(*alkira.AlkiraClient)
	id := d.Id()

	log.Printf("[INFO] Deleting Connector (Cisco SD-WAN) %s", id)
	err := client.DeleteConnectorCiscoSdwan(id)

	return err
}

// generateConnectorCiscoSdwanRequest generate request for Cisco SD-WAN connector
func generateConnectorCiscoSdwanRequest(ac *alkira.AlkiraClient, d *schema.ResourceData, m interface{}) (*alkira.ConnectorCiscoSdwan, error) {

	billingTags := convertTypeListToIntList(d.Get("billing_tag_ids").([]interface{}))
	mappings := expandCiscoSdwanVrfMappings(d.Get("vrf_segment_mapping").(*schema.Set))
	vedges := expandCiscoSdwanVedges(ac, d.Get("vedge").(*schema.Set))

	connector := &alkira.ConnectorCiscoSdwan{
		BillingTags:          billingTags,
		CiscoEdgeInfo:        vedges,
		CiscoEdgeVrfMappings: mappings,
		Cxp:                  d.Get("cxp").(string),
		Group:                d.Get("group").(string),
		Name:                 d.Get("name").(string),
		Size:                 d.Get("size").(string),
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
func expandCiscoSdwanVedges(ac *alkira.AlkiraClient, in *schema.Set) []alkira.CiscoSdwanEdgeInfo {
	if in == nil || in.Len() == 0 {
		log.Printf("[DEBUG] Empty vedges")
		return []alkira.CiscoSdwanEdgeInfo{}
	}

	mappings := make([]alkira.CiscoSdwanEdgeInfo, in.Len())

	for i, mapping := range in.List() {
		r := alkira.CiscoSdwanEdgeInfo{}
		t := mapping.(map[string]interface{})

		if v, ok := t["hostname"].(string); ok {
			r.HostName = v
		}
		if v, ok := t["cloud_init_file"].(string); ok {
			r.CloudInitFile = v
		}
		if v, ok := t["credential_id"].(string); ok {
			r.CredentialId = v
		}

		mappings[i] = r
	}

	return mappings
}
