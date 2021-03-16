package alkira

import (
	"log"
	"strconv"

	"github.com/alkiranet/alkira-client-go/alkira"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceAlkiraConnectorGcpVpc() *schema.Resource {
	return &schema.Resource{
		Description: "Manage GCP Cloud Connector.",

		Create: resourceConnectorGcpVpcCreate,
		Read:   resourceConnectorGcpVpcRead,
		Update: resourceConnectorGcpVpcUpdate,
		Delete: resourceConnectorGcpVpcDelete,

		Schema: map[string]*schema.Schema{
			"billing_tags": {
				Description: "Tags for billing.",
				Type:        schema.TypeList,
				Optional:    true,
				Elem:        &schema.Schema{Type: schema.TypeInt},
			},
			"connector_id": {
				Type:     schema.TypeInt,
				Computed: true,
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
			"gcp_region": {
				Description: "GCP region where VPC resides.",
				Type:        schema.TypeString,
				Required:    true,
			},
			"gcp_vpc_id": {
				Description: "GCP VPC ID.",
				Type:        schema.TypeString,
				Required:    true,
			},
			"gcp_vpc_name": {
				Description: "GCP VPC name.",
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
		},
	}
}

func resourceConnectorGcpVpcCreate(d *schema.ResourceData, m interface{}) error {
	client := m.(*alkira.AlkiraClient)

	billingTags := convertTypeListToIntList(d.Get("billing_tags").([]interface{}))
	segments := []string{d.Get("segment").(string)}

	connector := &alkira.ConnectorGcpVpc{
		BillingTags:    billingTags,
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
	d.Set("connector_id", id)

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
	err := client.DeleteConnectorGcpVpc(d.Get("connector_id").(int))

	if err != nil {
		return err
	}

	return nil
}
