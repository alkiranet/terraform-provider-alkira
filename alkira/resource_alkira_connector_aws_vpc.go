package alkira

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceAlkiraConnectorAwsVpc() *schema.Resource {
	return &schema.Resource{
		Create: resourceConnectorAwsVpcCreate,
		Read:   resourceConnectorAwsVpcRead,
		Update: resourceConnectorAwsVpcUpdate,
		Delete: resourceConnectorAwsVpcDelete,

		Schema: map[string]*schema.Schema{
			"vpc_1_id": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"vpc_1_name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"vpc_1_owner_id": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"vpc_2_id": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"vpc_2_name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"vpc_2_owner_id": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"size": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"segments": {
				Type: schema.TypeList,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Required: true,
				Description: "A list of segments associated with the connector",
			},
		},
	}
}

func resourceConnectorAwsVpcCreate(d *schema.ResourceData, m interface{}) error {
        return resourceConnectorAwsVpcRead(d, m)
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
