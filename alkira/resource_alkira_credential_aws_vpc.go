package alkira

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/alkiranet/alkira-client-go/alkira"

	"github.com/hashicorp/go-cty/cty"
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
			"You can provide your credentials via environmental variables:\n\n " +
			"* AK_AWS_ACCESS_KEY_ID\n * AK_AWS_SECRET_ACCESS_KEY\n * AK_AWS_ROLE_ARN\n " +
			"* AK_AWS_ROLE_EXTERNAL_ID\n\n",
		CreateContext: resourceCredentialAwsVpc,
		ReadContext:   resourceCredentialAwsVpcRead,
		UpdateContext: resourceCredentialAwsVpcUpdate,
		DeleteContext: resourceCredentialAwsVpcDelete,
		Importer: &schema.ResourceImporter{
			StateContext: importWithReadValidation(resourceCredentialAwsVpcRead),
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Description: "Name of the credential.",
				Type:        schema.TypeString,
				Required:    true,
			},
			"aws_access_key": {
				Description: "AWS access key.",
				Type:        schema.TypeString,
				Optional:    true,
				Sensitive:   true,
				WriteOnly:   true,
			},
			"aws_secret_key": {
				Description: "AWS secret key.",
				Type:        schema.TypeString,
				Optional:    true,
				Sensitive:   true,
				WriteOnly:   true,
			},
			"aws_role_arn": {
				Description: "AWS Role ARN.",
				Type:        schema.TypeString,
				Optional:    true,
				Sensitive:   true,
				WriteOnly:   true,
			},
			"aws_external_id": {
				Description: "AWS Role External ID.",
				Type:        schema.TypeString,
				Optional:    true,
				Sensitive:   true,
				WriteOnly:   true,
			},
			"type": {
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
	client := meta.(*alkira.AlkiraClient)

	credential, err := client.GetCredentialById(d.Id())
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			d.SetId("")
			return nil
		}
		return diag.FromErr(err)
	}

	d.Set("name", credential.Name)

	if credential.SubType != "" {
		d.Set("type", credential.SubType)
	}

	// Sensitive fields (access_key, secret_key, role_arn, external_id)
	// are NOT returned by the API for security reasons and are
	// maintained in the user's HCL configuration via WriteOnly.

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
		// Terraform may not print "with <resource address>" for destroys of objects
		// that are no longer in configuration, so include identifying context here.
		name, _ := d.GetOk("name")
		if nameStr, ok := name.(string); ok && nameStr != "" {
			return diag.FromErr(fmt.Errorf("%w alkira_credential_aws_vpc (name=%q id=%s)", err, nameStr, d.Id()))
		}
		return diag.FromErr(fmt.Errorf("%w alkira_credential_aws_vpc (id=%s)", err, d.Id()))
	}

	d.SetId("")
	return nil
}

// getAwsVpcCredentialValue gets a value from config or environment variable.
// For WriteOnly fields, reads from raw config since values are not stored in state.
func getAwsVpcCredentialValue(d *schema.ResourceData, field string, envVar string, required bool) (string, error) {
	// First try raw config (for WriteOnly fields)
	attrPath := cty.Path{cty.GetAttrStep{Name: field}}
	val, diags := d.GetRawConfigAt(attrPath)

	if !diags.HasError() && !val.IsNull() && val.IsKnown() && val.Type() == cty.String {
		strVal := val.AsString()
		if strVal != "" {
			return strVal, nil
		}
	}

	// Fall back to environment variable
	envValue := os.Getenv(envVar)
	if envValue != "" {
		return envValue, nil
	}

	if required {
		return "", fmt.Errorf("required field '%s' is not set in configuration and environment variable '%s' is not set", field, envVar)
	}
	return "", nil
}

func generateCredentialAwsVpc(d *schema.ResourceData) (interface{}, error) {
	credentialType := d.Get("type").(string)
	var c interface{}

	switch credentialType {
	case "ACCESS_KEY":
		accessKey, err := getAwsVpcCredentialValue(d, "aws_access_key", "AK_AWS_ACCESS_KEY_ID", true)
		if err != nil {
			return nil, err
		}

		secretKey, err := getAwsVpcCredentialValue(d, "aws_secret_key", "AK_AWS_SECRET_ACCESS_KEY", true)
		if err != nil {
			return nil, err
		}

		c = alkira.CredentialAwsVpcKey{
			Ec2AccessKey: accessKey,
			Ec2SecretKey: secretKey,
			Type:         credentialType,
		}
	case "ROLE":
		roleArn, err := getAwsVpcCredentialValue(d, "aws_role_arn", "AK_AWS_ROLE_ARN", true)
		if err != nil {
			return nil, err
		}

		externalId, _ := getAwsVpcCredentialValue(d, "aws_external_id", "AK_AWS_ROLE_EXTERNAL_ID", false)

		c = alkira.CredentialAwsVpcRole{
			Ec2RoleArn:    roleArn,
			Ec2ExternalId: externalId,
			Type:          credentialType,
		}
	default:
		return nil, errors.New("ERROR: Invalid AWS-VPC credential type")
	}

	log.Printf("[INFO] Creating credential (AWS-VPC) with type %s", credentialType)
	return c, nil
}
