package alkira

import (
	"context"
	"fmt"

	"github.com/alkiranet/alkira-client-go/alkira"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func resourceAlkiraConnectorAzureVnet() *schema.Resource {
	return &schema.Resource{
		Description:   "Manage Azure VNET Connector.",
		CreateContext: resourceConnectorAzureVnetCreate,
		ReadContext:   resourceConnectorAzureVnetRead,
		UpdateContext: resourceConnectorAzureVnetUpdate,
		DeleteContext: resourceConnectorAzureVnetDelete,
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
			"azure_vnet_id": {
				Description: "Azure Virtual Network Id.",
				Type:        schema.TypeString,
				Required:    true,
			},
			"billing_tag_ids": {
				Description: "Tags for billing.",
				Type:        schema.TypeSet,
				Optional:    true,
				Elem:        &schema.Schema{Type: schema.TypeInt},
			},
			"credential_id": {
				Description: "ID of credential managed by Credential Manager.",
				Type:        schema.TypeString,
				Required:    true,
			},
			"cxp": {
				Description: "The CXP where the connector should be provisioned.",
				Type:        schema.TypeString,
				Required:    true,
			},
			"connection_mode": {
				Description: "The mode that connector will use to connect to the " +
					"Alkira CXP. `VNET_GATEWAY` will connect with a Virtual " +
					"Gateway, `VNET_PEERING` will connect using an Alkira " +
					"Transit Hub (ATH).",
				Type:     schema.TypeString,
				Optional: true,
				Default:  "VNET_GATEWAY",
				ValidateFunc: validation.StringInSlice([]string{
					"VNET_GATEWAY", "VNET_PEERING"}, false),
			},
			"enabled": {
				Description: "Is the connector enabled. Default is `true`.",
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     true,
			},
			"failover_cxps": {
				Description: "A list of additional CXPs where the connector " +
					"should be provisioned for failover.",
				Type:     schema.TypeList,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"group": {
				Description: "The group of the connector.",
				Type:        schema.TypeString,
				Optional:    true,
			},
			"implicit_group_id": {
				Description: "The ID of implicit group automaticaly created " +
					"with the connector.",
				Type:     schema.TypeInt,
				Computed: true,
			},
			"name": {
				Description: "The name of the connector.",
				Type:        schema.TypeString,
				Required:    true,
			},
			"provision_state": {
				Description: "The provision state of the connector.",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"routing_options": {
				Description: " Routing options for the entire VNET, either " +
					"`ADVERTISE_DEFAULT_ROUTE` or `ADVERTISE_CUSTOM_PREFIX`. " +
					"Default value is `AVERTISE_DEFAULT_ROUTE`.",
				Type:     schema.TypeString,
				Optional: true,
				Default:  "ADVERTISE_DEFAULT_ROUTE",
				ValidateFunc: validation.StringInSlice([]string{
					"ADVERTISE_DEFAULT_ROUTE",
					"ADVERTISE_CUSTOM_PREFIX"}, false),
			},
			"routing_prefix_list_ids": {
				Description: "Prefix List IDs.",
				Type:        schema.TypeList,
				Optional:    true,
				Elem:        &schema.Schema{Type: schema.TypeInt},
			},
			"segment_id": {
				Description: "The ID of the segment assoicated with the connector.",
				Type:        schema.TypeString,
				Required:    true,
			},
			"vnet_cidr": &schema.Schema{
				Description: "Configure routing options on specified VNET CIDR.",
				Type:        schema.TypeSet,
				Optional:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"cidr": {
							Description: "VNET CIDR.",
							Type:        schema.TypeString,
							Required:    true,
						},
						"routing_options": {
							Description: "Routing options for the CIDR, either " +
								"`ADVERTISE_DEFAULT_ROUTE` or " +
								"`ADVERTISE_CUSTOM_PREFIX`.",
							Type:     schema.TypeString,
							Optional: true,
							ValidateFunc: validation.StringInSlice([]string{
								"ADVERTISE_DEFAULT_ROUTE",
								"ADVERTISE_CUSTOM_PREFIX"}, false),
						},
						"prefix_list_ids": {
							Description: "Prefix List IDs.",
							Type:        schema.TypeList,
							Optional:    true,
							Elem:        &schema.Schema{Type: schema.TypeInt},
						},
						"service_tags": {
							Description: "List of service tags provided by Azure.",
							Type:        schema.TypeSet,
							Optional:    true,
							Elem:        &schema.Schema{Type: schema.TypeString},
						},
						"udr_list_ids": {
							Description: "User defined routes list (`list_udr`).",
							Type:        schema.TypeSet,
							Optional:    true,
							Elem:        &schema.Schema{Type: schema.TypeInt},
						},
					},
				},
			},
			"vnet_subnet": &schema.Schema{
				Description: "Configure routing options on the specified VNET subnet.",
				Type:        schema.TypeSet,
				Optional:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"subnet_id": {
							Description: "VNET subnet ID.",
							Type:        schema.TypeString,
							Required:    true,
						},
						"subnet_cidr": {
							Description: "VNET subnet CIDR.",
							Type:        schema.TypeString,
							Optional:    true,
						},
						"routing_options": {
							Description: "Routing options for the subnet, " +
								"either `ADVERTISE_DEFAULT_ROUTE` " +
								"or `ADVERTISE_CUSTOM_PREFIX`.",
							Type:     schema.TypeString,
							Optional: true,
							ValidateFunc: validation.StringInSlice([]string{
								"ADVERTISE_DEFAULT_ROUTE",
								"ADVERTISE_CUSTOM_PREFIX"}, false),
						},
						"prefix_list_ids": {
							Description: "Prefix List IDs.",
							Type:        schema.TypeList,
							Optional:    true,
							Elem:        &schema.Schema{Type: schema.TypeInt},
						},
						"service_tags": {
							Description: "List of service tags provided by Azure.",
							Type:        schema.TypeSet,
							Optional:    true,
							Elem:        &schema.Schema{Type: schema.TypeString},
						},
						"udr_list_ids": {
							Description: "User defined routes list (`list_udr`).",
							Type:        schema.TypeSet,
							Optional:    true,
							Elem:        &schema.Schema{Type: schema.TypeInt},
						},
					},
				},
			},
			"service_tags": {
				Description: "list of service tags from Azure. Providing a service tag here " +
					"would result in service tag route configuration on VNET route table, so " +
					"that the traffic toward the service would directly steer towards those " +
					"services, and would not go via Alkira network.",
				Type:     schema.TypeList,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Optional: true,
			},
			"size": {
				Description: "The size of the connector, one of `SMALL`, `MEDIUM`, " +
					"`LARGE`, `2LARGE`, `4LARGE`, `5LARGE`, `10LARGE`, `20LARGE`.",
				Type:     schema.TypeString,
				Required: true,
				ValidateFunc: validation.StringInSlice([]string{
					"SMALL", "MEDIUM", "LARGE", `2LARGE`,
					`4LARGE`, `5LARGE`, `10LARGE`, `20LARGE`}, false),
			},
			"customer_asn": {
				Description: "A specific BGP ASN for the connector. This cannot be specified " +
					"when `connection_mode` is `VNET_PEERING`. This field cannot be updated " +
					"once the connector has been provisioned. The ASN cannot be value that " +
					"is [restricted by Azure]" +
					"(https://learn.microsoft.com/en-us/azure/vpn-gateway/vpn-gateway-vpn-faq#bgp).",
				Type:     schema.TypeInt,
				Optional: true,
			},
			"scale_group_id": {
				Description: "The ID of the scale group associated with the connector.",
				Type:        schema.TypeString,
				Optional:    true,
			},
		},
	}
}

func resourceConnectorAzureVnetCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {

	// INIT
	client := m.(*alkira.AlkiraClient)
	api := alkira.NewConnectorAzureVnet(m.(*alkira.AlkiraClient))

	request, err := generateConnectorAzureVnetRequest(d, m)

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

	return resourceConnectorAzureVnetRead(ctx, d, m)
}

func resourceConnectorAzureVnetRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {

	client := m.(*alkira.AlkiraClient)
	api := alkira.NewConnectorAzureVnet(m.(*alkira.AlkiraClient))

	// Get the resource
	connector, provState, err := api.GetById(d.Id())

	if err != nil {
		return diag.Diagnostics{{
			Severity: diag.Warning,
			Summary:  "FAILED TO GET RESOURCE",
			Detail:   fmt.Sprintf("%s", err),
		}}
	}

	d.Set("azure_vnet_id", connector.VnetId)
	d.Set("billing_tag_ids", connector.BillingTags)
	d.Set("credential_id", connector.CredentialId)
	d.Set("cxp", connector.CXP)
	d.Set("connection_mode", connector.ConnectionMode)
	d.Set("enabled", connector.Enabled)
	d.Set("failover_cxps", connector.SecondaryCXPs)
	d.Set("group", connector.Group)
	d.Set("implicit_group_id", connector.ImplicitGroupId)
	d.Set("name", connector.Name)
	d.Set("routing_options", connector.VnetRouting.ImportOptions.RouteImportMode)
	d.Set("routing_prefix_list_ids", connector.VnetRouting.ImportOptions.PrefixListIds)
	d.Set("size", connector.Size)
	d.Set("service_tags", connector.ServiceTags)
	d.Set("customer_asn", connector.CustomerASN)
	d.Set("scale_group_id", connector.ScaleGroupId)

	setVnetRouting(d, connector.VnetRouting)

	// Get segment
	numOfSegments := len(connector.Segments)
	if numOfSegments == 1 {
		segmentId, err := getSegmentIdByName(connector.Segments[0], m)

		if err != nil {
			return diag.FromErr(err)
		}
		d.Set("segment_id", segmentId)
	} else {
		return diag.FromErr(fmt.Errorf("the number of segments are invalid %n", numOfSegments))
	}

	// Set provision state
	if client.Provision == true && provState != "" {
		d.Set("provision_state", provState)
	}

	return nil
}

func resourceConnectorAzureVnetUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {

	// INIT
	client := m.(*alkira.AlkiraClient)
	api := alkira.NewConnectorAzureVnet(m.(*alkira.AlkiraClient))

	request, err := generateConnectorAzureVnetRequest(d, m)

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

	return resourceConnectorAzureVnetRead(ctx, d, m)
}

func resourceConnectorAzureVnetDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {

	// INIT
	client := m.(*alkira.AlkiraClient)
	api := alkira.NewConnectorAzureVnet(m.(*alkira.AlkiraClient))

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

// generateConnectorAzureVnetRequest generate request for connector-azure-vnet
func generateConnectorAzureVnetRequest(d *schema.ResourceData, m interface{}) (*alkira.ConnectorAzureVnet, error) {

	//
	// Segment
	//
	segmentName, err := getSegmentNameById(d.Get("segment_id").(string), m)

	if err != nil {
		return nil, err
	}

	//
	// Routing
	//
	routing, err := constructVnetRouting(d)

	if err != nil {
		return nil, err
	}

	// Assemble request
	request := &alkira.ConnectorAzureVnet{
		BillingTags:    convertTypeSetToIntList(d.Get("billing_tag_ids").(*schema.Set)),
		CXP:            d.Get("cxp").(string),
		ConnectionMode: d.Get("connection_mode").(string),
		CredentialId:   d.Get("credential_id").(string),
		Enabled:        d.Get("enabled").(bool),
		Group:          d.Get("group").(string),
		Name:           d.Get("name").(string),
		SecondaryCXPs:  convertTypeListToStringList(d.Get("failover_cxps").([]interface{})),
		Segments:       []string{segmentName},
		Size:           d.Get("size").(string),
		ServiceTags:    convertTypeListToStringList(d.Get("service_tags").([]interface{})),
		VnetId:         d.Get("azure_vnet_id").(string),
		VnetRouting:    routing,
		CustomerASN:    d.Get("customer_asn").(int),
		ScaleGroupId:   d.Get("scale_group_id").(string),
	}

	return request, nil
}
