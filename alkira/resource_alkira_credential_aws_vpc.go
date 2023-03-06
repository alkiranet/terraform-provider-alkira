package alkira

import (
	"context"
	"errors"
	"log"

	"github.com/alkiranet/alkira-client-go/alkira"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceAlkiraCredentialAwsVpc() *schema.Resource {
	return &schema.Resource{
		Description: "Manage AWS credential for authentication.\n\n" +
			"The following methods are supported:\n\n" +
			" - Static credentials\n" +
			" - Environment variables\n\n" +
			"### Static Credentials:\n\n" +
			"Static credentials can be provided by adding an `aws_access_key`" +
			"and `aws_secret_key` in-line in the AWS provider block.\n\n" +
			"### Environment Variables:\n\n" +
			"You can provide your credentials via enviromental variables:\n\n " +
			"* AK_AWS_ACCESS_KEY_ID\n * AK_AWS_SECRET_ACCESS_KEY\n * AK_AWS_ROLE_ARN\n " +
			"* AK_AWS_ROLE_EXTERNAL_ID\n\n",
		CreateContext: resourceCredentialAwsVpc,
		ReadContext:   resourceCredentialAwsVpcRead,
		UpdateContext: resourceCredentialAwsVpcUpdate,
		DeleteContext: resourceCredentialAwsVpcDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			"name": &schema.Schema{
				Description: "Name of the credential.",
				Type:        schema.TypeString,
				Required:    true,
			},
			"aws_access_key": &schema.Schema{
				Description: "AWS access key.",
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc(
					"AK_AWS_ACCESS_KEY_ID",
					nil),
			},
			"aws_secret_key": &schema.Schema{
				Description: "AWS secret key.",
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc(
					"AK_AWS_SECRET_ACCESS_KEY",
					nil),
			},
			"aws_role_arn": &schema.Schema{
				Description: "AWS Role ARN.",
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc(
					"AK_AWS_ROLE_ARN",
					nil),
			},
			"aws_external_id": &schema.Schema{
				Description: "AWS Role External ID.",
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc(
					"AK_AWS_ROLE_EXTERNAL_ID",
					nil),
			},
			"type": &schema.Schema{
				Description: "Type of AWS-VPC credential.",
				Type:        schema.TypeString,
				Required:    true,
			},
		},
	}
}

func resourceCredentialAwsVpc(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*alkira.AlkiraClient)

	c, err := generateCredentialAwsVpc(d)

	if err != nil {
		return diag.FromErr(err)
	}

	id, err := client.CreateCredential(d.Get("name").(string), alkira.CredentialTypeAwsVpc, c, 0)

	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(id)
	return resourceCredentialAwsVpcRead(ctx, d, meta)
}

func resourceCredentialAwsVpcRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	return nil
}

func resourceCredentialAwsVpcUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*alkira.AlkiraClient)

	c, err := generateCredentialAwsVpc(d)

	if err != nil {
		return diag.FromErr(err)
	}

	log.Printf("[INFO] Updating credential (AWS-VPC) %s", d.Id())
	err = client.UpdateCredential(d.Id(), d.Get("name").(string), alkira.CredentialTypeAwsVpc, c, 0)

	if err != nil {
		return diag.FromErr(err)
	}

	return resourceCredentialAwsVpcRead(ctx, d, meta)
}

func resourceCredentialAwsVpcDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*alkira.AlkiraClient)
	credentialId := d.Id()

	log.Printf("[INFO] Deleting credential (AWS-VPC %s)\n", credentialId)
	err := client.DeleteCredential(credentialId, alkira.CredentialTypeAwsVpc)

	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId("")
	return nil
}

func generateCredentialAwsVpc(d *schema.ResourceData) (interface{}, error) {
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
		return nil, errors.New("Invalid AWS-VPC Credential Type")
	}

	log.Printf("[INFO] Creating credential (AWS-VPC) with type %s", credentialType)
	return c, nil
}
