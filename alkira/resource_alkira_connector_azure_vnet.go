package alkira

import (
	"log"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/alkiranet/alkira-client-go/alkira"
)

func resourceAlkiraConnectorAzureVnet() *schema.Resource {
	return &schema.Resource{
		Create: resourceConnectorAzureVnetCreate,
		Read:   resourceConnectorAzureVnetRead,
		Update: resourceConnectorAzureVnetUpdate,
		Delete: resourceConnectorAzureVnetDelete,

		Schema: map[string]*schema.Schema{
			"azure_region": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"azure_vnet_id": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"connector_id": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"credential_id": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"cxp": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"group": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"segment": {
				Type: schema.TypeString,
				Required: true,
				Description: "A segment associated with the connector AWS-VPC",
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

	segments := []string{d.Get("segment").(string)}

	connector := &alkira.ConnectorAzureVnetRequest{
		CXP:            d.Get("cxp").(string),
		CredentialId:   d.Get("credential_id").(string),
		CustomerRegion: d.Get("azure_region").(string),
		Group:          d.Get("group").(string),
		Name:           d.Get("name").(string),
        Segments:       segments,
        Size:           d.Get("size").(string),
        VnetId:         d.Get("azure_vnet_id").(string),
	}

	id, err := client.CreateConnectorAzureVnet(connector)

	if err != nil {
		return err
	}

	d.SetId(strconv.Itoa(id))
	d.Set("connector_id", strconv.Itoa(id))

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
	err := client.DeleteConnectorAzureVnet(d.Id())

	if err != nil {
		return err
	}

	return nil
}
