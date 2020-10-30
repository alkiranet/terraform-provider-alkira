package alkira

import (
	"log"
	"strconv"

	"github.com/alkiranet/alkira-client-go/alkira"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceAlkiraConnectorAzureVnet() *schema.Resource {
	return &schema.Resource{
		Create: resourceConnectorAzureVnetCreate,
		Read:   resourceConnectorAzureVnetRead,
		Update: resourceConnectorAzureVnetUpdate,
		Delete: resourceConnectorAzureVnetDelete,

		Schema: map[string]*schema.Schema{
			"azure_region": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				Description: "Azure Region",
			},
			"azure_vnet_id": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				Description: "Azure Virutal Network Id",
			},
			"billing_tags": {
				Type:     schema.TypeList,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeInt},
			},
			"connector_id": &schema.Schema{
				Type:     schema.TypeInt,
				Computed: true,
			},
			"credential_id": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				Description: "The credentials for creating connector",
			},
			"cxp": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				Description: "The CXP to be used for the connector",
			},
			"group": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Description: "A user group that the connector belongs to",
			},
			"name": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				Description: "The name of the connector",
			},
			"segment": {
				Type:     schema.TypeString,
				Required: true,
			},
			"size": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
		},
	}
}

func resourceConnectorAzureVnetCreate(d *schema.ResourceData, m interface{}) error {
	client := m.(*alkira.AlkiraClient)

	billingTags := convertTypeListToIntList(d.Get("billing_tags").([]interface{}))
	segments := []string{d.Get("segment").(string)}

	connector := &alkira.ConnectorAzureVnetRequest{
		BillingTags:    billingTags,
		CXP:            d.Get("cxp").(string),
		CredentialId:   d.Get("credential_id").(string),
		CustomerRegion: d.Get("azure_region").(string),
		Group:          d.Get("group").(string),
		Name:           d.Get("name").(string),
		Segments:       segments,
		Size:           d.Get("size").(string),
		VnetId:         d.Get("azure_vnet_id").(string),
	}

	log.Printf("[INFO] Creating Connector (AZURE-VNET)")
	id, err := client.CreateConnectorAzureVnet(connector)

	if err != nil {
		return err
	}

	d.SetId(strconv.Itoa(id))
	d.Set("connector_id", id)

	return resourceConnectorAzureVnetRead(d, m)
}

func resourceConnectorAzureVnetRead(d *schema.ResourceData, m interface{}) error {
	return nil
}

func resourceConnectorAzureVnetUpdate(d *schema.ResourceData, m interface{}) error {
	return resourceConnectorAzureVnetRead(d, m)
}

func resourceConnectorAzureVnetDelete(d *schema.ResourceData, m interface{}) error {
	client := m.(*alkira.AlkiraClient)

	log.Printf("[INFO] Deleting Connector (AZURE-VNET) %s", d.Id())
	err := client.DeleteConnectorAzureVnet(d.Get("connector_id").(int))

	if err != nil {
		return err
	}

	return nil
}
