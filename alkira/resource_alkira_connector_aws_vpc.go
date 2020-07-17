package alkira

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/alkiranet/terraform-provider-alkira/alkira/internal"
)

func resourceAlkiraConnectorAwsVpc() *schema.Resource {
	return &schema.Resource{
		Create: resourceConnectorAwsVpcCreate,
		Read:   resourceConnectorAwsVpcRead,
		Update: resourceConnectorAwsVpcUpdate,
		Delete: resourceConnectorAwsVpcDelete,

		Schema: map[string]*schema.Schema{
			"name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"vpc_id": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"aws_account_id": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"aws_region": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"aws_access_key": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"aws_secret_key": &schema.Schema{
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
			"id": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
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

func resourceConnectorAwsVpcCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*internal.AlkiraClient)

	name      := d.Get("name").(string)
	accessKey := d.Get("aws_access_key").(string)
	secretKey := d.Get("aws_secret_key").(string)

	credentialId, statusCode := client.CreateCredentialAwsVpc(name, accessKey, secretKey)

	if statusCode != 201 {
		return fmt.Errorf("failed to save credential (%d)", statusCode)
	}

	segments := []string{d.Get("segment").(string)}

	connectorAwsVpc := &internal.ConnectorAwsVpcRequest{
		CXP:            d.Get("cxp").(string),
		CredentialId:   credentialId,
		CustomerName:   client.Username,
		CustomerRegion: d.Get("aws_region").(string),
		Group:          d.Get("group").(string),
		Name:           d.Get("name").(string),
        Segments:       segments,
        Size:           d.Get("size").(string),
        VpcId:          d.Get("vpc_id").(string),
        VpcOwnerId:     d.Get("aws_account_id").(string),
	}

	client.CreateConnectorAwsVpc(connectorAwsVpc)
	return resourceConnectorAwsVpcRead(d, meta)
}

func resourceConnectorAwsVpcRead(d *schema.ResourceData, m interface{}) error {
        return nil
}

func resourceConnectorAwsVpcUpdate(d *schema.ResourceData, m interface{}) error {
        return resourceConnectorAwsVpcRead(d, m)
}

func resourceConnectorAwsVpcDelete(d *schema.ResourceData, m interface{}) error {
        return nil
}
