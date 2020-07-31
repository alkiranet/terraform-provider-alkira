package alkira

import (
	"log"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/alkiranet/alkira-client-go/alkira"
)

func resourceAlkiraConnectorGcpVpc() *schema.Resource {
	return &schema.Resource{
		Create: resourceConnectorGcpVpcCreate,
		Read:   resourceConnectorGcpVpcRead,
		Update: resourceConnectorGcpVpcUpdate,
		Delete: resourceConnectorGcpVpcDelete,

		Schema: map[string]*schema.Schema{
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
			"gcp_region": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"gcp_vpc_id": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"gcp_vpc_name": &schema.Schema{
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

func resourceConnectorGcpVpcCreate(d *schema.ResourceData, m interface{}) error {
	client := m.(*alkira.AlkiraClient)

	segments := []string{d.Get("segment").(string)}

	connector := &alkira.ConnectorGcpVpcRequest{
		CXP:            d.Get("cxp").(string),
		CredentialId:   d.Get("credential_id").(string),
		CustomerRegion: d.Get("gcp_region").(string),
		Group:          d.Get("group").(string),
		Name:           d.Get("name").(string),
        Segments:       segments,
        Size:           d.Get("size").(string),
        VpcId:          d.Get("gcp_vpc_id").(string),
        VpcName:        d.Get("gcp_vpc_name").(string),
	}

	log.Printf("[INFO] Creating Connector (GCP-VPC)")
	id, err := client.CreateConnectorGcpVpc(connector)

	if err != nil {
		return err
	}

	d.SetId(strconv.Itoa(id))
	d.Set("connector_id", strconv.Itoa(id))

	return resourceConnectorGcpVpcRead(d, m)
}

func resourceConnectorGcpVpcRead(d *schema.ResourceData, m interface{}) error {
        return nil
}

func resourceConnectorGcpVpcUpdate(d *schema.ResourceData, m interface{}) error {
        return resourceConnectorGcpVpcRead(d, m)
}

func resourceConnectorGcpVpcDelete(d *schema.ResourceData, m interface{}) error {
	client := m.(*alkira.AlkiraClient)

	log.Printf("[INFO] Deleting Connector (GCP-VPC) %s", d.Id())
	err := client.DeleteConnectorGcpVpc(d.Id())

	if err != nil {
		return err
	}

	return nil
}
