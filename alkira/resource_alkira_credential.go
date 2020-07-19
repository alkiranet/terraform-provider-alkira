package alkira

import (
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/alkiranet/terraform-provider-alkira/alkira/internal"
)

func resourceAlkiraCredentialAwsVpc() *schema.Resource {
	return &schema.Resource{
		Create: resourceCredential,
		Read:   resourceCredentialRead,
		Update: resourceCredentialUpdate,
		Delete: resourceCredentialDelete,

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

func resourceCredential(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*internal.AlkiraClient)

	name      := d.Get("name").(string)
	accessKey := d.Get("aws_access_key").(string)
	secretKey := d.Get("aws_secret_key").(string)
	authType  := d.Get("type").(string)

	log.Printf("[INFO] Createing credential-aws-vpc")
	credentialId, err := client.CreateCredentialAwsVpc(name, accessKey, secretKey, authType)

	if err != nil {
		return err
	}

	d.SetId(credentialId)
	return resourceCredentialRead(d, meta)
}

func resourceCredentialRead(d *schema.ResourceData, meta interface{}) error {
        return nil
}

func resourceCredentialUpdate(d *schema.ResourceData, meta interface{}) error {
        return resourceCredentialRead(d, meta)
}

func resourceCredentialDelete(d *schema.ResourceData, meta interface{}) error {
	client       := meta.(*internal.AlkiraClient)
	credentialId := d.Id()

	log.Printf("[INFO] Deleting credential-aws-vpc %s\n", credentialId)
	err := client.DeleteCredentialAwsVpc(credentialId)

	if err != nil {
		return err
	}

	return nil
}
