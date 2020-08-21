package alkira

import (
	"log"
	"strconv"

	"github.com/alkiranet/alkira-client-go/alkira"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceAlkiraConnectorAwsVpc() *schema.Resource {
	return &schema.Resource{
		Create: resourceConnectorAwsVpcCreate,
		Read:   resourceConnectorAwsVpcRead,
		Update: resourceConnectorAwsVpcUpdate,
		Delete: resourceConnectorAwsVpcDelete,

		Schema: map[string]*schema.Schema{
			"aws_account_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"aws_region": {
				Type:     schema.TypeString,
				Required: true,
			},
			"billing_tags": {
				Type:     schema.TypeList,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"credential_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"connector_id": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"cxp": {
				Type:     schema.TypeString,
				Required: true,
			},
			"group":{
				Type:     schema.TypeString,
				Optional: true,
			},
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"segment": {
				Type:     schema.TypeString,
				Required: true,
			},
			"size": {
				Type:     schema.TypeString,
				Required: true,
			},
			"vpc_id": {
				Type:     schema.TypeString,
				Required: true,
			},
		},
	}
}

func resourceConnectorAwsVpcCreate(d *schema.ResourceData, m interface{}) error {
	client := m.(*alkira.AlkiraClient)

	billingTags := convertTypeListToStringList(d.Get("billing_tags").([]interface{}))
	segments    := []string{d.Get("segment").(string)}

	connectorAwsVpc := &alkira.ConnectorAwsVpcRequest{
		BillingTags:    billingTags,
		CXP:            d.Get("cxp").(string),
		CredentialId:   d.Get("credential_id").(string),
		CustomerName:   client.Username,
		CustomerRegion: d.Get("aws_region").(string),
		Group:          d.Get("group").(string),
		Name:           d.Get("name").(string),
		Segments:       segments,
		Size:           d.Get("size").(string),
		VpcId:          d.Get("vpc_id").(string),
		VpcOwnerId:     d.Get("aws_account_id").(string),
	}

	id, err := client.CreateConnectorAwsVpc(connectorAwsVpc)

	if err != nil {
		return err
	}

	d.SetId(strconv.Itoa(id))
	d.Set("connector_id", id)
	return resourceConnectorAwsVpcRead(d, m)
}

func resourceConnectorAwsVpcRead(d *schema.ResourceData, m interface{}) error {
	return nil
}

func resourceConnectorAwsVpcUpdate(d *schema.ResourceData, m interface{}) error {
	return resourceConnectorAwsVpcRead(d, m)
}

func resourceConnectorAwsVpcDelete(d *schema.ResourceData, m interface{}) error {
	client := m.(*alkira.AlkiraClient)

	log.Printf("[INFO] Deleting Connector (AWS-VPC) %s", d.Id())
	err := client.DeleteConnectorAwsVpc(d.Get("connector_id").(int))

	if err != nil {
		return err
	}

	return nil
}
