package alkira

import (
	"log"
	"strconv"

	"github.com/alkiranet/alkira-client-go/alkira"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func resourceAlkiraConnectorAzureVnet() *schema.Resource {
	return &schema.Resource{
		Description: "Manage Azure VNET Connector.",
		Create:      resourceConnectorAzureVnetCreate,
		Read:        resourceConnectorAzureVnetRead,
		Update:      resourceConnectorAzureVnetUpdate,
		Delete:      resourceConnectorAzureVnetDelete,

		Schema: map[string]*schema.Schema{
			"azure_vnet_id": {
				Description: "Azure Virtual Network Id.",
				Type:        schema.TypeString,
				Required:    true,
			},
			"billing_tag_ids": {
				Description: "Tags for billing.",
				Type:        schema.TypeList,
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
			"enabled": {
				Description: "Is the connector enabled. Default is `true`.",
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     true,
			},
			"failover_cxps": {
				Description: "A list of additional CXPs where the connector should be provisioned for failover.",
				Type:        schema.TypeList,
				Optional:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
			"group": {
				Description: "The group of the connector.",
				Type:        schema.TypeString,
				Optional:    true,
			},
			"implicit_group_id": {
				Description: "The ID of implicit group automaticaly created with the connector.",
				Type:        schema.TypeInt,
				Computed:    true,
			},
			"name": {
				Description: "The name of the connector.",
				Type:        schema.TypeString,
				Required:    true,
			},
			"routing_options": {
				Description: " Routing options for the entire VNET, either `ADVERTISE_DEFAULT_ROUTE` " +
					"or `ADVERTISE_CUSTOM_PREFIX`. Default is `AVERTISE_DEFAULT_ROUTE`.",
				Type:         schema.TypeString,
				Optional:     true,
				Default:      "ADVERTISE_DEFAULT_ROUTE",
				ValidateFunc: validation.StringInSlice([]string{"ADVERTISE_DEFAULT_ROUTE", "ADVERTISE_CUSTOM_PREFIX"}, false),
			},
			"routing_prefix_list_ids": {
				Description: "Prefix List IDs.",
				Type:        schema.TypeList,
				Optional:    true,
				Elem:        &schema.Schema{Type: schema.TypeInt},
			},
			"segment_id": {
				Description: "The ID of the segment assoicated with the connector.",
				Type:        schema.TypeInt,
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
							Description:  "Routing options for the CIDR, either `ADVERTISE_DEFAULT_ROUTE` or `ADVERTISE_CUSTOM_PREFIX`.",
							Type:         schema.TypeString,
							Optional:     true,
							ValidateFunc: validation.StringInSlice([]string{"ADVERTISE_DEFAULT_ROUTE", "ADVERTISE_CUSTOM_PREFIX"}, false),
						},
						"prefix_list_ids": {
							Description: "Prefix List IDs.",
							Type:        schema.TypeList,
							Optional:    true,
							Elem:        &schema.Schema{Type: schema.TypeInt},
						},
						"service_tags": {
							Description: "List of service tags provided by Azure.",
							Type:        schema.TypeList,
							Optional:    true,
							Elem:        &schema.Schema{Type: schema.TypeString},
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
							Description:  "Routing options for the subnet, either `ADVERTISE_DEFAULT_ROUTE` or `ADVERTISE_CUSTOM_PREFIX`.",
							Type:         schema.TypeString,
							Optional:     true,
							ValidateFunc: validation.StringInSlice([]string{"ADVERTISE_DEFAULT_ROUTE", "ADVERTISE_CUSTOM_PREFIX"}, false),
						},
						"prefix_list_ids": {
							Description: "Prefix List IDs.",
							Type:        schema.TypeList,
							Optional:    true,
							Elem:        &schema.Schema{Type: schema.TypeInt},
						},
						"service_tags": {
							Description: "List of service tags provided by Azure.",
							Type:        schema.TypeList,
							Optional:    true,
							Elem:        &schema.Schema{Type: schema.TypeString},
						},
					},
				},
			},
			"service_tags": {
				Description: "list of service tags from Azure. Providing a service tag here, " +
					"would result in service tag route configuration on VNET route table, so " +
					"that the traffic toward the service would directly steer towards those " +
					"services, and would not go via Alkira network.",
				Type:     schema.TypeList,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Optional: true,
			},
			"size": {
				Description:  "The size of the connector, one of `SMALL`, `MEDIUM`, `LARGE`, `2LARGE`, `4LARGE`, `5LARGE`, `10LARGE`, `20LARGE`.",
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringInSlice([]string{"SMALL", "MEDIUM", "LARGE", `2LARGE`, `4LARGE`, `5LARGE`, `10LARGE`, `20LARGE`}, false),
			},
		},
	}
}

func resourceConnectorAzureVnetCreate(d *schema.ResourceData, m interface{}) error {
	client := m.(*alkira.AlkiraClient)
	connector, err := generateConnectorAzureVnetRequest(d, m)

	if err != nil {
		return err
	}

	id, err := client.CreateConnectorAzureVnet(connector)

	if err != nil {
		return err
	}

	d.SetId(id)
	return resourceConnectorAzureVnetRead(d, m)
}

func resourceConnectorAzureVnetRead(d *schema.ResourceData, m interface{}) error {
	client := m.(*alkira.AlkiraClient)

	connector, err := client.GetConnectorAzureVnet(d.Id())

	if err != nil {
		return err
	}

	d.Set("azure_vnet_id", connector.VnetId)
	d.Set("billing_tag_ids", connector.BillingTags)
	d.Set("credential_id", connector.CredentialId)
	d.Set("cxp", connector.CXP)
	d.Set("enabled", connector.Enabled)
	d.Set("failover_cxps", connector.SecondaryCXPs)
	d.Set("group", connector.Group)
	d.Set("implicit_group_id", connector.ImplicitGroupId)
	d.Set("name", connector.Name)
	d.Set("routing_options", connector.VnetRouting.ImportOptions.RouteImportMode)
	d.Set("routing_prefix_list_ids", connector.VnetRouting.ImportOptions.PrefixListIds)
	d.Set("size", connector.Size)
	d.Set("service_tags", connector.ServiceTags)

	if len(connector.Segments) > 0 {
		segment, err := client.GetSegmentByName(connector.Segments[0])

		if err != nil {
			return err
		}
		d.Set("segment_id", segment.Id)
	}

	return nil
}

func resourceConnectorAzureVnetUpdate(d *schema.ResourceData, m interface{}) error {
	client := m.(*alkira.AlkiraClient)
	connector, err := generateConnectorAzureVnetRequest(d, m)

	if err != nil {
		return err
	}

	err = client.UpdateConnectorAzureVnet(d.Id(), connector)

	if err != nil {
		return err
	}

	return resourceConnectorAzureVnetRead(d, m)
}

func resourceConnectorAzureVnetDelete(d *schema.ResourceData, m interface{}) error {
	client := m.(*alkira.AlkiraClient)

	err := client.DeleteConnectorAzureVnet(d.Id())

	return err
}

// generateConnectorAzureVnetRequest generate request for connector-azure-vnet
func generateConnectorAzureVnetRequest(d *schema.ResourceData, m interface{}) (*alkira.ConnectorAzureVnet, error) {
	client := m.(*alkira.AlkiraClient)

	billingTags := convertTypeListToIntList(d.Get("billing_tag_ids").([]interface{}))
	failoverCXPs := convertTypeListToStringList(d.Get("failover_cxps").([]interface{}))
	serviceTags := convertTypeListToStringList(d.Get("service_tags").([]interface{}))

	segment, err := client.GetSegmentById(strconv.Itoa(d.Get("segment_id").(int)))

	if err != nil {
		log.Printf("[ERROR] failed to get segment by Id: %d", d.Get("segment_id"))
		return nil, err
	}

	routing, err := constructVnetRouting(d)

	if err != nil {
		return nil, err
	}

	request := &alkira.ConnectorAzureVnet{
		BillingTags:   billingTags,
		CXP:           d.Get("cxp").(string),
		CredentialId:  d.Get("credential_id").(string),
		Enabled:       d.Get("enabled").(bool),
		Group:         d.Get("group").(string),
		Name:          d.Get("name").(string),
		SecondaryCXPs: failoverCXPs,
		Segments:      []string{segment.Name},
		Size:          d.Get("size").(string),
		ServiceTags:   serviceTags,
		VnetId:        d.Get("azure_vnet_id").(string),
		VnetRouting:   routing,
	}

	return request, nil
}
