package alkira

import (
	"context"
	"fmt"
	"log"

	"github.com/alkiranet/alkira-client-go/alkira"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func resourceAlkiraConnectorVmwareSdwan() *schema.Resource {
	return &schema.Resource{
		Description: "Manage VMWARE SD-WAN Connector.",
		Create:      resourceConnectorVmwareSdwanCreate,
		Read:        resourceConnectorVmwareSdwanRead,
		Update:      resourceConnectorVmwareSdwanUpdate,
		Delete:      resourceConnectorVmwareSdwanDelete,
		CustomizeDiff: func(ctx context.Context, d *schema.ResourceDiff, m interface{}) error {
			client := m.(*alkira.AlkiraClient)

			old, _ := d.GetChange("provision_state")

			if client.Provision == true && old == "FAILED" {
				d.SetNew("provision_state", "SUCCESS")
			}

			return nil
		},
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Description: "The name of the connector.",
				Type:        schema.TypeString,
				Required:    true,
			},
			"billing_tag_ids": {
				Description: "A list of Billing Tag by ID associated " +
					"with the connector.",
				Type:     schema.TypeList,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeInt},
			},
			"cxp": {
				Description: "The CXP where the connector should be " +
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
				Description: "The ID of implicit group automaticaly created " +
					"with the connector.",
				Type:     schema.TypeInt,
				Computed: true,
			},
			"size": &schema.Schema{
				Description: "The size of the connector, one of `SMALL`, " +
					"`MEDIUM` and `LARGE`, `2LARGE`, `4LARGE`, `5LARGE`, " +
					"`10LARGE` and `20LARGE`.",
				Type:     schema.TypeString,
				Required: true,
				ValidateFunc: validation.StringInSlice([]string{
					"SMALL", "MEDIUM",
					"LARGE", "2LARGE",
					"4LARGE", "5LARGE",
					"10LARGE", "20LARGE"}, false),
			},
			"virtual_edge": &schema.Schema{
				Description: "Virtual Edge",
				Type:        schema.TypeList,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"credential_id": {
							Description: "The generated credential ID.",
							Type:        schema.TypeString,
							Computed:    true,
						},
						"hostname": {
							Description: "The hostname of the virtual edge.",
							Type:        schema.TypeString,
							Required:    true,
						},
						"id": {
							Description: "The ID of the virtual edge.",
							Type:        schema.TypeInt,
							Computed:    true,
						},
						"username": &schema.Schema{
							Description: "VMWARE SD-WAN username. It could be also " +
								"set by environment variable `AK_VMWARE_SDWAN_USERNAME`.",
							Type:     schema.TypeString,
							Required: true,
							DefaultFunc: schema.EnvDefaultFunc(
								"AK_VMWARE_SDWAN_USERNAME",
								nil),
						},
						"password": &schema.Schema{
							Description: "VMWARE SD-WAN password. It could be also " +
								"set by environment variable `AK_VMWARE_SDWAN_PASSWORD`.",
							Type:     schema.TypeString,
							Required: true,
							DefaultFunc: schema.EnvDefaultFunc(
								"AK_VMWARE_SDWAN_PASSWORD",
								nil),
						},
					},
				},
				Required: true,
			},
			"version": {
				Description: "The version of VMWARE SD-WAN.",
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
						"gateway_bgp_asn": {
							Description: "BGP ASN on the customer premise side. " +
								"A typical value for 2 byte segment " +
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

func resourceConnectorVmwareSdwanCreate(d *schema.ResourceData, m interface{}) error {

	client := m.(*alkira.AlkiraClient)
	api := alkira.NewConnectorVmwareSdwan(m.(*alkira.AlkiraClient))

	// Construct request
	request, err := generateConnectorVmwareSdwanRequest(d, m)

	if err != nil {
		return err
	}

	// Send create request
	response, provisionState, err := api.Create(request)

	if err != nil {
		return err
	}

	// Set states
	d.SetId(string(response.Id))

	if client.Provision == true {
		d.Set("provision_state", provisionState)
	}

	return resourceConnectorVmwareSdwanRead(d, m)
}

func resourceConnectorVmwareSdwanRead(d *schema.ResourceData, m interface{}) error {

	client := m.(*alkira.AlkiraClient)
	api := alkira.NewConnectorVmwareSdwan(m.(*alkira.AlkiraClient))

	connector, err := api.GetById(d.Id())

	if err != nil {
		return err
	}

	d.Set("billing_tag_ids", connector.BillingTags)
	d.Set("cxp", connector.Cxp)
	d.Set("group", connector.Group)
	d.Set("implicit_group_id", connector.ImplicitGroupId)
	d.Set("name", connector.Name)
	d.Set("size", connector.Size)

	// Set virtual edge
	setVirtualEdge(d, connector)

	// Set vrf_segment_mapping
	var mappings []map[string]interface{}

	for _, m := range connector.VmWareSdWanVRFMappings {
		mapping := map[string]interface{}{
			"advertise_on_prem_routes": m.AdvertiseOnPremRoutes,
			"allow_nat_exit":           m.DisableInternetExit,
			"gateway_bgp_asn":          m.GatewayBgpAsn,
			"segment_id":               m.SegmentId,
		}
		mappings = append(mappings, mapping)
	}

	d.Set("vrf_segment_mapping", mappings)
	d.Set("version", connector.Version)

	// Set provision state
	_, provisionState, err := api.GetByName(d.Get("name").(string))

	if client.Provision == true && provisionState != "" {
		d.Set("provision_state", provisionState)
	}

	return nil
}

func resourceConnectorVmwareSdwanUpdate(d *schema.ResourceData, m interface{}) error {

	client := m.(*alkira.AlkiraClient)
	api := alkira.NewConnectorVmwareSdwan(m.(*alkira.AlkiraClient))

	// Construct update request
	request, err := generateConnectorVmwareSdwanRequest(d, m)

	if err != nil {
		return err
	}

	// Send update request
	provisionState, err := api.Update(d.Id(), request)

	if err != nil {
		return err
	}

	// Set provision state
	if client.Provision == true {
		d.Set("provision_state", provisionState)
	}

	return resourceConnectorVmwareSdwanRead(d, m)
}

func resourceConnectorVmwareSdwanDelete(d *schema.ResourceData, m interface{}) error {

	client := m.(*alkira.AlkiraClient)
	api := alkira.NewConnectorVmwareSdwan(m.(*alkira.AlkiraClient))

	provisionState, err := api.Delete(d.Id())

	if err != nil {
		return err
	}

	if client.Provision == true && provisionState != "SUCCESS" {
		return fmt.Errorf("failed to delete connector_vmware_sdwan %s, provision failed", d.Id())
	}

	d.SetId("")

	return nil
}

// generateConnectorVmwareSdwanRequest generate request for VMWARE SD-WAN connector
func generateConnectorVmwareSdwanRequest(d *schema.ResourceData, m interface{}) (*alkira.ConnectorVmwareSdwan, error) {

	//
	// Expand virtual_edge block
	//
	virtualEdges, err := expandVmwareSdwanVirtualEdges(m.(*alkira.AlkiraClient), d.Get("virtual_edge").([]interface{}))

	if err != nil {
		return nil, err
	}

	// Construct the request payload
	connector := &alkira.ConnectorVmwareSdwan{
		BillingTags:            convertTypeListToIntList(d.Get("billing_tag_ids").([]interface{})),
		Instances:              virtualEdges,
		VmWareSdWanVRFMappings: expandVmwareSdwanVrfMappings(d.Get("vrf_segment_mapping").(*schema.Set)),
		Cxp:                    d.Get("cxp").(string),
		Group:                  d.Get("group").(string),
		Name:                   d.Get("name").(string),
		Size:                   d.Get("size").(string),
		Version:                d.Get("version").(string),
	}

	return connector, nil
}

// expandVmwareSdwanVrfMappings expand VMWARE SD-WAN VRF segment mapping
func expandVmwareSdwanVrfMappings(in *schema.Set) []alkira.VmwareSdwanEdgeVrfMapping {

	if in == nil || in.Len() == 0 {
		log.Printf("[DEBUG] Empty vrf_segment_mapping")
		return []alkira.VmwareSdwanVrfMapping{}
	}

	mappings := make([]alkira.VmwareSdwanVrfMapping, in.Len())
	for i, mapping := range in.List() {
		r := alkira.VmwareSdwanVrfMapping{}
		t := mapping.(map[string]interface{})

		if v, ok := t["advertise_on_prem_routes"].(bool); ok {
			r.AdvertiseOnPremRoutes = v
		}
		if v, ok := t["allow_nat_exit"].(bool); ok {
			r.DisableInternetExit = !v
		}
		if v, ok := t["gateway_bgp_asn"].(int); ok {
			r.GatewayBgpAsn = v
		}
		if v, ok := t["segment_id"].(int); ok {
			r.SegmentId = v
		}

		mappings[i] = r
	}

	return mappings
}

// expandVmwareSdwanVedges expand virtual edges
func expandVmwareSdwanVirtualEdges(ac *alkira.AlkiraClient, in []interface{}) ([]alkira.VmwareSdwanInstance, error) {

	if in == nil || len(in) == 0 {
		log.Printf("[DEBUG] Empty vedges")
		return []alkira.VmwareSdwanInstance{}, nil
	}

	mappings := make([]alkira.VmwareSdwanInstance, len(in))

	for i, mapping := range in {
		r := alkira.VmwareSdwanInstance{}
		t := mapping.(map[string]interface{})

		var username string
		var password string

		if v, ok := t["hostname"].(string); ok {
			r.HostName = v
		}
		if v, ok := t["username"].(string); ok {
			username = v
		}
		if v, ok := t["password"].(string); ok {
			password = v
		}
		if v, ok := t["credential_id"].(string); ok {
			if v == "" {
				log.Printf("[INFO] Creating VMWARE-SDWAN Credential")
				credentialName := r.HostName + randomNameSuffix()

				credential := alkira.CredentialVmwareSdwan{
					Username: username,
					Password: password,
				}

				credentialId, err := ac.CreateCredential(
					credentialName,
					alkira.CredentialTypeVmwareSdwan,
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
