package alkira

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
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
				DefaultFunc: schema.EnvDefaultFunc(
					"AK_AWS_ACCESS_KEY_ID",
					nil),
			},
			"aws_secret_key": {
				Description: "AWS secret key.",
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc(
					"AK_AWS_SECRET_ACCESS_KEY",
					nil),
			},
			"aws_role_arn": {
				Description: "AWS Role ARN.",
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc(
					"AK_AWS_ROLE_ARN",
					nil),
			},
			"aws_external_id": {
				Description: "AWS Role External ID.",
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc(
					"AK_AWS_ROLE_EXTERNAL_ID",
					nil),
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

	// Get all credentials and find the one matching the ID
	// Note: GetCredentialById() doesn't work in all environments (405 error),
	// so we use GetCredentials() which lists all credentials.
	credentialsJSON, err := client.GetCredentials()
	if err != nil {
		return diag.Diagnostics{{
			Severity: diag.Warning,
			Summary:  "FAILED TO GET RESOURCE",
			Detail:   fmt.Sprintf("Error reading credential %s: %s", d.Id(), err),
		}}
	}

	// Parse the JSON array of credentials
	var credentials []alkira.CredentialResponseDetail
	if err := json.Unmarshal([]byte(credentialsJSON), &credentials); err != nil {
		return diag.Diagnostics{{
			Severity: diag.Warning,
			Summary:  "FAILED TO GET RESOURCE",
			Detail:   fmt.Sprintf("Error parsing credentials for %s: %s", d.Id(), err),
		}}
	}

	// Find the credential matching the ID
	var credential *alkira.CredentialResponseDetail
	for i := range credentials {
		if credentials[i].Id == d.Id() {
			credential = &credentials[i]
			break
		}
	}

	// If not found, return error
	if credential == nil {
		return diag.Diagnostics{{
			Severity: diag.Warning,
			Summary:  "FAILED TO GET RESOURCE",
			Detail:   fmt.Sprintf("Credential %s not found", d.Id()),
		}}
	}

	// Set fields returned by API
	d.Set("name", credential.Name)
	d.Set("type", credential.SubType)

	// Note: Sensitive fields (access_key, secret_key, role_arn, external_id)
	// are NOT returned by the API for security reasons and must be
	// maintained in the user's HCL configuration.

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

func generateCredentialAwsVpc(d *schema.ResourceData) (interface{}, error) {
	credentialType := d.Get("type").(string)
	var c interface{}

	switch credentialType {
	case "ACCESS_KEY":
		c = alkira.CredentialAwsVpcKey{
			Ec2AccessKey: d.Get("aws_access_key").(string),
			Ec2SecretKey: d.Get("aws_secret_key").(string),
			Type:         d.Get("type").(string),
		}
	case "ROLE":
		c = alkira.CredentialAwsVpcRole{
			Ec2RoleArn:    d.Get("aws_role_arn").(string),
			Ec2ExternalId: d.Get("aws_external_id").(string),
			Type:          d.Get("type").(string),
		}
	default:
		return nil, errors.New("ERROR: Invalid AWS-VPC credential type")
	}

	log.Printf("[INFO] Creating credential (AWS-VPC) with type %s", credentialType)
	return c, nil
}
