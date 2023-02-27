package alkira

import (
	"context"
	"fmt"

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
			"orchestrator_host": {
				Description: "VMWare (Velo) Orchestrator portal host address.",
				Type:        schema.TypeString,
				Required:    true,
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
			"tunnel_protocol": {
				Description: "Only supported tunnel protocol is `IPSEC` for now.",
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "IPSEC",
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
						"activation_code": &schema.Schema{
							Description: "Activation code generated in " +
								"VMWare orchestrator account.",
							Type:     schema.TypeString,
							Required: true,
							DefaultFunc: schema.EnvDefaultFunc(
								"AK_VMWARE_SDWAN_ACTIVATION_CODE",
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
			"target_segment": {
				Description: "Specify target segment.",
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
							Optional: true,
							Default:  65000,
						},
						"segment_id": {
							Description: "Alkira Segment ID.",
							Type:        schema.TypeInt,
							Required:    true,
						},
						"vmware_sdwan_segment_name": {
							Description: "VMWare SD-WAN Segment name for " +
								"correlating with Alkria segment.",
							Type:     schema.TypeString,
							Required: true,
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
	d.Set("orchestrator_host", connector.OrchestratorHostAddress)
	d.Set("size", connector.Size)
	d.Set("tunnel_protocol", connector.TunnelProtocol)

	// Set virtual edge
	setVirtualEdge(d, connector)

	// Set vrf_segment_mapping
	var mappings []map[string]interface{}

	for _, m := range connector.VmWareSdWanVRFMappings {
		mapping := map[string]interface{}{
			"advertise_on_prem_routes":   m.AdvertiseOnPremRoutes,
			"allow_nat_exit":             m.DisableInternetExit,
			"gateway_bgp_asn":            m.GatewayBgpAsn,
			"segment_id":                 m.SegmentId,
			"vmware_sdwang_segment_name": m.VmWareSdWanSegmentName,
		}
		mappings = append(mappings, mapping)
	}

	d.Set("target_segment", mappings)
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
		BillingTags:             convertTypeListToIntList(d.Get("billing_tag_ids").([]interface{})),
		Instances:               virtualEdges,
		VmWareSdWanVRFMappings:  expandVmwareSdwanVrfMappings(d.Get("target_segment").(*schema.Set)),
		Cxp:                     d.Get("cxp").(string),
		Group:                   d.Get("group").(string),
		OrchestratorHostAddress: d.Get("orchestrator_host").(string),
		Name:                    d.Get("name").(string),
		Size:                    d.Get("size").(string),
		TunnelProtocol:          d.Get("tunnel_protocol").(string),
		Version:                 d.Get("version").(string),
	}

	return connector, nil
}
