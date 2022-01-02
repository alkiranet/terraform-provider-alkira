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
		Description: "Manage Azure Cloud Connector.",

		Create: resourceConnectorAzureVnetCreate,
		Read:   resourceConnectorAzureVnetRead,
		Update: resourceConnectorAzureVnetUpdate,
		Delete: resourceConnectorAzureVnetDelete,

		Schema: map[string]*schema.Schema{
			"azure_region": {
				Description: "Azure Region where VNET resides.",
				Type:        schema.TypeString,
				Required:    true,
			},
			"azure_vnet_id": {
				Description: "Azure Virtual Network Id.",
				Type:        schema.TypeString,
				Required:    true,
			},
			"azure_subscription_id": {
				Description: "The Azure subscription ID of the VNET. If the" +
					"`subscirption_id` was provided in the credential, the one" +
					"in the credential will be always used.",
				Type:     schema.TypeString,
				Optional: true,
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
			"group": {
				Description: "The group of the connector.",
				Type:        schema.TypeString,
				Optional:    true,
			},
			"name": {
				Description: "The name of the connector.",
				Type:        schema.TypeString,
				Required:    true,
			},
			"routing_options": {
				Description:  "Routing options, either `ADVERTISE_DEFAULT_ROUTE` or `ADVERTISE_CUSTOM_PREFIX`.",
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
			"size": {
				Description:  "The size of the connector, one of `SMALL`, `MEDIUM` or `LARGE`.",
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringInSlice([]string{"SMALL", "MEDIUM", "LARGE"}, false),
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

	d.Set("azure_region", connector.CustomerRegion)
	d.Set("azure_subscription_id", connector.SubscriptionId)
	d.Set("azure_vnet_id", connector.VnetId)
	d.Set("billing_tag_ids", connector.BillingTags)
	d.Set("credential_id", connector.CredentialId)
	d.Set("cxp", connector.CXP)
	d.Set("group", connector.Group)
	d.Set("name", connector.Name)
	d.Set("routing_options", connector.VnetRouting.ImportOptions.RouteImportMode)
	d.Set("routing_prefix_list_ids", connector.VnetRouting.ImportOptions.PrefixListIds)
	d.Set("size", connector.Size)

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
	routing := constructVnetRouting(d.Get("routing_options").(string), d.Get("routing_prefix_list_ids").([]interface{}))

	segment, err := client.GetSegmentById(strconv.Itoa(d.Get("segment_id").(int)))

	if err != nil {
		log.Printf("[ERROR] failed to get segment by Id: %d", d.Get("segment_id"))
		return nil, err
	}

	request := &alkira.ConnectorAzureVnet{
		BillingTags:    billingTags,
		CXP:            d.Get("cxp").(string),
		CredentialId:   d.Get("credential_id").(string),
		CustomerRegion: d.Get("azure_region").(string),
		Group:          d.Get("group").(string),
		Name:           d.Get("name").(string),
		Segments:       []string{segment.Name},
		Size:           d.Get("size").(string),
		SubscriptionId: d.Get("azure_subscription_id").(string),
		VnetId:         d.Get("azure_vnet_id").(string),
		VnetRouting:    routing,
	}

	return request, nil
}

// constructVnetRouting expand AZURE VNET routing options
func constructVnetRouting(option string, prefixList []interface{}) *alkira.ConnectorVnetRouting {

	routing := alkira.ConnectorVnetImportOptions{}

	routing.RouteImportMode = option
	routing.PrefixListIds = convertTypeListToIntList(prefixList)

	return &alkira.ConnectorVnetRouting{routing}
}
