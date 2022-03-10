package alkira

import (
	"log"
	"os"

	"github.com/alkiranet/alkira-client-go/alkira"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// Provider returns a schema.Provider for Alkira.
func Provider() *schema.Provider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"portal": {
				Description: "The URL for Alkira Custom Portal.",
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: envDefaultFunc("ALKIRA_PORTAL"),
			},
			"username": {
				Description: "Your Tenant Username.",
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: envDefaultFunc("ALKIRA_USERNAME"),
			},
			"password": {
				Description: "Your Tenant Password.",
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: envDefaultFunc("ALKIRA_PASSWORD"),
			},
		},

		ResourcesMap: map[string]*schema.Resource{
			"alkira_billing_tag":             resourceAlkiraBillingTag(),
			"alkira_connector_aws_vpc":       resourceAlkiraConnectorAwsVpc(),
			"alkira_connector_azure_vnet":    resourceAlkiraConnectorAzureVnet(),
			"alkira_connector_cisco_sdwan":   resourceAlkiraConnectorCiscoSdwan(),
			"alkira_connector_gcp_vpc":       resourceAlkiraConnectorGcpVpc(),
			"alkira_connector_oci_vcn":       resourceAlkiraConnectorOciVcn(),
			"alkira_connector_internet_exit": resourceAlkiraConnectorInternetExit(),
			"alkira_connector_ipsec":         resourceAlkiraConnectorIPSec(),
			"alkira_credential_aws_vpc":      resourceAlkiraCredentialAwsVpc(),
			"alkira_credential_azure_vnet":   resourceAlkiraCredentialAzureVnet(),
			"alkira_credential_cisco_sdwan":  resourceAlkiraCredentialCiscoSdwan(),
			"alkira_credential_fortinet":     resourceAlkiraCredentialFortinet(),
			"alkira_credential_gcp_vpc":      resourceAlkiraCredentialGcpVpc(),
			"alkira_credential_oci_vcn":      resourceAlkiraCredentialOciVcn(),
			"alkira_credential_pan":          resourceAlkiraCredentialPan(),
			"alkira_credential_pan_instance": resourceAlkiraCredentialPanInstance(),
			"alkira_cloudvisor_account":      resourceAlkiraCloudVisorAccount(),
			"alkira_group_connector":         resourceAlkiraGroupConnector(),
			"alkira_fortinet":                resourceAlkiraFortinet(),
			"alkira_group":                   resourceAlkiraGroup(),
			"alkira_group_user":              resourceAlkiraGroupUser(),
			"alkira_list_as_path":            resourceAlkiraListAsPath(),
			"alkira_list_community":          resourceAlkiraListCommunity(),
			"alkira_list_extended_community": resourceAlkiraListExtendedCommunity(),
			"alkira_list_global_cidr":        resourceAlkiraListGlobalCidr(),
			"alkira_internet_application":    resourceAlkiraInternetApplication(),
			"alkira_policy":                  resourceAlkiraPolicy(),
			"alkira_policy_nat":              resourceAlkiraPolicyNat(),
			"alkira_policy_nat_rule":         resourceAlkiraPolicyNatRule(),
			"alkira_policy_prefix_list":      resourceAlkiraPolicyPrefixList(),
			"alkira_policy_rule":             resourceAlkiraPolicyRule(),
			"alkira_policy_rule_list":        resourceAlkiraPolicyRuleList(),
			"alkira_segment":                 resourceAlkiraSegment(),
			"alkira_segment_resource":        resourceAlkiraSegmentResource(),
			"alkira_segment_resource_share":  resourceAlkiraSegmentResourceShare(),
      "alkira_service_pan":             resourceAlkiraServicePan(),
			"alkira_tenant_network":          resourceAlkiraTenantNetwork(),
		},
		DataSourcesMap: map[string]*schema.Resource{
			"alkira_billing_tag":        dataSourceAlkiraBillingTag(),
			"alkira_credential":         dataSourceAlkiraCredential(),
			"alkira_group":              dataSourceAlkiraGroup(),
			"alkira_group_connector":    dataSourceAlkiraGroupConnector(),
			"alkira_group_user":         dataSourceAlkiraGroupUser(),
			"alkira_policy_prefix_list": dataSourceAlkiraPolicyPrefixList(),
			"alkira_segment":            dataSourceAlkiraSegment(),
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
	alkiraClient, err := alkira.NewAlkiraClient(d.Get("portal").(string), d.Get("username").(string), d.Get("password").(string))

	if err != nil {
		log.Printf("[ERROR] failed to initialize alkira provider, please check your credential and portal URI.")
		return nil, err
	}

	return alkiraClient, nil
}
