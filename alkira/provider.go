package alkira

import (
	"os"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/alkiranet/terraform-provider-alkira/alkira/internal"
)

// Provider returns a schema.Provider for Alkira.
func Provider() terraform.ResourceProvider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"portal": {
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: envDefaultFunc("ALKIRA_PORTAL"),
			},
			"username": {
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: envDefaultFunc("ALKIRA_USERNAME"),
			},
			"password": {
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: envDefaultFunc("ALKIRA_PASSWORD"),
			},
			"skip_version_validation": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
		},

		ResourcesMap: map[string]*schema.Resource{
			"alkira_connector_aws_vpc":      resourceAlkiraConnectorAwsVpc(),
			"alkira_connector_azure_vnet":   resourceAlkiraConnectorAzureVnet(),
			"alkira_connector_gcp_vpc":      resourceAlkiraConnectorGcpVpc(),
			"alkira_credential_aws_vpc":     resourceAlkiraCredentialAwsVpc(),
			"alkira_credential_azure_vnet":  resourceAlkiraCredentialAzureVnet(),
			"alkira_credential_gcp_vpc":     resourceAlkiraCredentialGcpVpc(),
			"alkira_segment":                resourceAlkiraSegment(),
			"alkira_tenant_network":         resourceAlkiraTenantNetwork(),
		},
		DataSourcesMap: map[string]*schema.Resource{
			"alkira_connector_aws_vpc":    dataSourceAlkiraConnectorAwsVpc(),
			"alkira_connector_azure_vnet": dataSourceAlkiraConnectorAzureVnet(),
			"alkira_connector_gcp_vpc":    dataSourceAlkiraConnectorGcpVpc(),
			"alkira_segment":              dataSourceAlkiraSegment(),
			"alkira_tenant_network":       dataSourceAlkiraTenantNetwork(),
		},
		ConfigureFunc: alkiraConfigure,
	}
}

func envDefaultFunc(k string) schema.SchemaDefaultFunc {
	return func() (interface{}, error) {
		if v := os.Getenv(k); v != "" {
			return v, nil
		}

		return nil, nil
	}
}

func alkiraConfigure(d *schema.ResourceData) (interface{}, error) {
	alkiraClient, err := internal.NewAlkiraClient(d.Get("portal").(string), d.Get("username").(string), d.Get("password").(string))

	if err != nil {
		return nil, nil
	}

	return alkiraClient, nil
}
