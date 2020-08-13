package alkira

import (
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/alkiranet/alkira-client-go/alkira"
)

func resourceAlkiraCredentialAwsVpc() *schema.Resource {
	return &schema.Resource{
		Create: resourceCredentialAwsVpc,
		Read:   resourceCredentialAwsVpcRead,
		Update: resourceCredentialAwsVpcUpdate,
		Delete: resourceCredentialAwsVpcDelete,

		Schema: map[string]*schema.Schema{
			"name": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				Description: "The name of the credential",
			},
			"aws_access_key": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				Description: "The EC2 access key",

			},
			"aws_secret_key": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				Description: "The EC2 secret key",

			},
			"type": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				Description: "The Type of AWS-VPC credential",

			},
		},
	}
}

func resourceCredentialAwsVpc(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*alkira.AlkiraClient)

	c := alkira.CredentialAwsVpc{
		Ec2AccessKey: d.Get("aws_access_key").(string),
		Ec2SecretKey: d.Get("aws_secret_key").(string),
		Type:         d.Get("type").(string),
	}

	log.Printf("[INFO] Createing credential-aws-vpc")
	credentialId, err := client.CreateCredential(d.Get("name").(string), "awsvpc", c)

	if err != nil {
		return err
	}

	d.SetId(credentialId)
	return resourceCredentialAwsVpcRead(d, meta)
}

func resourceCredentialAwsVpcRead(d *schema.ResourceData, meta interface{}) error {
        return nil
}

func resourceCredentialAwsVpcUpdate(d *schema.ResourceData, meta interface{}) error {
        return resourceCredentialAwsVpcRead(d, meta)
}

func resourceCredentialAwsVpcDelete(d *schema.ResourceData, meta interface{}) error {
	client       := meta.(*alkira.AlkiraClient)
	credentialId := d.Id()

	log.Printf("[INFO] Deleting credential (AWS-VPC %s)\n", credentialId)
	err := client.DeleteCredential(credentialId, "awsvpc")

	if err != nil {
		return err
	}

	return nil
}
