package alkira

import (
	"os"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
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
			"alkira_connector_aws_vpc":    resourceAlkiraConnectorAwsVpc(),
			"alkira_connector_azure_vnet": resourceAlkiraConnectorAzureVnet(),
			"alkira_connector_gcp_vpc":    resourceAlkiraConnectorGcpVpc(),
		},
		DataSourcesMap: map[string]*schema.Resource{
			"alkira_connector_aws_vpc":    dataSourceAlkiraConnectorAwsVpc(),
			"alkira_connector_azure_vnet": dataSourceAlkiraConnectorAzureVnet(),
			"alkira_connector_gcp_vpc":    dataSourceAlkiraConnectorGcpVpc(),
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
	config := Config{
		Portal:   d.Get("portal").(string),
		Username: d.Get("username").(string),
		Password: d.Get("password").(string),
	}

	skipVersionValidation := d.Get("skip_version_validation").(bool)
	if skipVersionValidation {
		return config.Client()
	}

	client, err := config.Client()
	if err != nil {
		return nil, err
	}

	return client, nil
}

func alkiraConfigureWithoutVersionValidation(d *schema.ResourceData) (interface{}, error) {
	config := Config{
		Portal:   d.Get("portal").(string),
		Username: d.Get("username").(string),
		Password: d.Get("password").(string),
	}

	return config.Client()
}
