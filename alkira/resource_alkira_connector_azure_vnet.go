package alkira

import (
	"log"

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
			"billing_tags": {
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
			"segment": {
				Description: "The segment of the connector.",
				Type:        schema.TypeString,
				Required:    true,
			},
			"size": {
				Description: "The size of the connector, one of `SMALL`, `MEDIUM` or `LARGE`.",
				Type:        schema.TypeString,
				Required:    true,
			},
			"routing": {
				Description: "Routing Configurations.",
				Type:        schema.TypeSet,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"prefix_list_ids": {
							Description: "Prefix List Ids.",
							Type:        schema.TypeList,
							Optional:    true,
							Elem:        &schema.Schema{Type: schema.TypeInt},
						},
						"options": {
							Description:  "Routing options, either `ADVERTISE_DEFAULT_ROUTE` or `ADVERTISE_CUSTOM_PREFIX`.",
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: validation.StringInSlice([]string{"ADVERTISE_DEFAULT_ROUTE", "ADVERTISE_CUSTOM_PREFIX"}, false),
						},
					},
				},
				Optional: true,
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

	log.Printf("[INFO] Creating Connector (AZURE-VNET)")
	id, err := client.CreateConnectorAzureVnet(connector)

	if err != nil {
		return err
	}

	d.SetId(id)

	return resourceConnectorAzureVnetRead(d, m)
}

func resourceConnectorAzureVnetRead(d *schema.ResourceData, m interface{}) error {
	return nil
}

func resourceConnectorAzureVnetUpdate(d *schema.ResourceData, m interface{}) error {
	client := m.(*alkira.AlkiraClient)
	connector, err := generateConnectorAzureVnetRequest(d, m)

	if err != nil {
		return err
	}

	log.Printf("[INFO] Updating Connector (AZURE-VNET) %s", d.Id())
	err = client.UpdateConnectorAzureVnet(d.Id(), connector)

	if err != nil {
		return err
	}

	return resourceConnectorAzureVnetRead(d, m)
}

func resourceConnectorAzureVnetDelete(d *schema.ResourceData, m interface{}) error {
	client := m.(*alkira.AlkiraClient)

	log.Printf("[INFO] Deleting Connector (AZURE-VNET) %s", d.Id())
	err := client.DeleteConnectorAzureVnet(d.Id())

	return err
}

// generateConnectorAzureVnetRequest generate request for connector-azure-vnet
func generateConnectorAzureVnetRequest(d *schema.ResourceData, m interface{}) (*alkira.ConnectorAzureVnetRequest, error) {
	billingTags := convertTypeListToIntList(d.Get("billing_tags").([]interface{}))
	segments := []string{d.Get("segment").(string)}
	routing := expandVnetRouting(d.Get("routing").(*schema.Set))

	request := &alkira.ConnectorAzureVnetRequest{
		BillingTags:    billingTags,
		CXP:            d.Get("cxp").(string),
		CredentialId:   d.Get("credential_id").(string),
		CustomerRegion: d.Get("azure_region").(string),
		Group:          d.Get("group").(string),
		Name:           d.Get("name").(string),
		Segments:       segments,
		Size:           d.Get("size").(string),
		VnetId:         d.Get("azure_vnet_id").(string),
		VnetRouting:    routing,
	}

	return request, nil
}

// expandVnetRouting expand AZURE VNET routing options
func expandVnetRouting(in *schema.Set) *alkira.ConnectorVnetRouting {
	if in == nil || in.Len() == 0 {
		log.Printf("[DEBUG] Empty routing input.")
		return nil
	}

	if in.Len() > 1 {
		log.Printf("[ERROR] Only one routing section is allowed.")
		return nil
	}

	routing := alkira.ConnectorVnetImportOptions{}

	for _, r := range in.List() {
		t := r.(map[string]interface{})

		if v, ok := t["options"].(string); ok {
			routing.RouteImportMode = v
		}

		routing.PrefixListIds = convertTypeListToIntList(t["prefix_list_ids"].([]interface{}))
	}

	return &alkira.ConnectorVnetRouting{routing}
}
