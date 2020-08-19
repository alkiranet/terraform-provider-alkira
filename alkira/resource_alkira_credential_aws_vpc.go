package alkira

import (
	"errors"
	"log"

	"github.com/alkiranet/alkira-client-go/alkira"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
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
				Optional:    true,
				Description: "AWS access key",
			},
			"aws_secret_key": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Description: "AWS secret key",
			},
			"aws_role_arn": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The AWS Role Arn",
			},
			"aws_external_id": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The AWS Role External ID",
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

	credentialType := d.Get("type").(string)
	var c interface{}

	if credentialType == "ACCESS_KEY" {
		c = alkira.CredentialAwsVpcKey{
			Ec2AccessKey: d.Get("aws_access_key").(string),
			Ec2SecretKey: d.Get("aws_secret_key").(string),
			Type:         d.Get("type").(string),
		}
	} else if credentialType == "ROLE" {
		c = alkira.CredentialAwsVpcRole{
			Ec2RoleArn:    d.Get("aws_role_arn").(string),
			Ec2ExternalId: d.Get("aws_external_id").(string),
			Type:          d.Get("type").(string),
		}
	} else {
		return errors.New("Invalid AWS-VPC Credential Type")
	}

	log.Printf("[INFO] Createing credential (AWS-VPC) with type %s", credentialType)
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
	client := meta.(*alkira.AlkiraClient)
	credentialId := d.Id()

	log.Printf("[INFO] Deleting credential (AWS-VPC %s)\n", credentialId)
	err := client.DeleteCredential(credentialId, "awsvpc")

	if err != nil {
		return err
	}

	return nil
}
