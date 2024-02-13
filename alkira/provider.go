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
				Description: "Your Tenant Username. If this is not provided then `api_key` must have a value.",
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: envDefaultFunc("ALKIRA_USERNAME"),
			},
			"password": {
				Description: "Your Tenant Password. If this is not provided then `api_key` must have a value.",
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: envDefaultFunc("ALKIRA_PASSWORD"),
			},
			"api_key": {
				Description: "Your Alkira API key. If thie is not provided then `username` and `password` must have a value.",
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: envDefaultFunc("ALKIRA_API_KEY"),
			},
			"provision": {
				Description: "With provision or not.",
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				DefaultFunc: envDefaultFunc("ALKIRA_PROVISION"),
			},
			"auth": {
				Description: "Authentication Method.",
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: envDefaultFunc("ALKIRA_AUTH"),
			},
		},

		ResourcesMap: map[string]*schema.Resource{
			"alkira_billing_tag":                    resourceAlkiraBillingTag(),
			"alkira_byoip_prefix":                   resourceAlkiraByoipPrefix(),
			"alkira_cloudvisor_account":             resourceAlkiraCloudVisorAccount(),
			"alkira_connector_akamai_prolexic":      resourceAlkiraConnectorAkamaiProlexic(),
			"alkira_connector_aruba_edge":           resourceAlkiraConnectorArubaEdge(),
			"alkira_connector_aws_vpc":              resourceAlkiraConnectorAwsVpc(),
			"alkira_connector_azure_vnet":           resourceAlkiraConnectorAzureVnet(),
			"alkira_connector_azure_expressroute":   resourceAlkiraConnectorAzureExpressRoute(),
			"alkira_connector_cisco_sdwan":          resourceAlkiraConnectorCiscoSdwan(),
			"alkira_connector_fortinet_sdwan":       resourceAlkiraConnectorFortinetSdwan(),
			"alkira_connector_gcp_vpc":              resourceAlkiraConnectorGcpVpc(),
			"alkira_connector_oci_vcn":              resourceAlkiraConnectorOciVcn(),
			"alkira_connector_internet_exit":        resourceAlkiraConnectorInternetExit(),
			"alkira_connector_ipsec":                resourceAlkiraConnectorIPSec(),
			"alkira_connector_ipsec_adv":            resourceAlkiraConnectorIPSecAdv(),
			"alkira_connector_ipsec_tunnel_profile": resourceAlkiraConnectorIpsecTunnelProfile(),
			"alkira_connector_remote_access":        resourceAlkiraConnectorRemoteAccess(),
			"alkira_connector_versa_sdwan":          resourceAlkiraConnectorVersaSdwan(),
			"alkira_connector_vmware_sdwan":         resourceAlkiraConnectorVmwareSdwan(),
			"alkira_credential_aws_vpc":             resourceAlkiraCredentialAwsVpc(),
			"alkira_credential_azure_vnet":          resourceAlkiraCredentialAzureVnet(),
			"alkira_credential_gcp_vpc":             resourceAlkiraCredentialGcpVpc(),
			"alkira_credential_oci_vcn":             resourceAlkiraCredentialOciVcn(),
			"alkira_credential_ssh_key_pair":        resourceAlkiraCredentialSshKeyPair(),
			"alkira_flow_collector":                 resourceAlkiraFlowCollector(),
			"alkira_group":                          resourceAlkiraGroup(),
			"alkira_group_user":                     resourceAlkiraGroupUser(),
			"alkira_group_direct_inter_connector":   resourceAlkiraDirectInterConnectorGroup(),
			"alkira_internet_application":           resourceAlkiraInternetApplication(),
			"alkira_ip_reservation":                 resourceAlkiraIpReservation(),
			"alkira_list_as_path":                   resourceAlkiraListAsPath(),
			"alkira_list_community":                 resourceAlkiraListCommunity(),
			"alkira_list_dns_server":                resourceAlkiraListDnsServer(),
			"alkira_list_extended_community":        resourceAlkiraListExtendedCommunity(),
			"alkira_list_global_cidr":               resourceAlkiraListGlobalCidr(),
			"alkira_list_policy_fqdn":               resourceAlkiraListPolicyFqdn(),
			"alkira_policy":                         resourceAlkiraPolicy(),
			"alkira_policy_nat":                     resourceAlkiraPolicyNat(),
			"alkira_policy_nat_rule":                resourceAlkiraPolicyNatRule(),
			"alkira_policy_prefix_list":             resourceAlkiraPolicyPrefixList(),
			"alkira_policy_routing":                 resourceAlkiraPolicyRouting(),
			"alkira_policy_rule":                    resourceAlkiraPolicyRule(),
			"alkira_policy_rule_list":               resourceAlkiraPolicyRuleList(),
			"alkira_segment":                        resourceAlkiraSegment(),
			"alkira_segment_resource":               resourceAlkiraSegmentResource(),
			"alkira_segment_resource_share":         resourceAlkiraSegmentResourceShare(),
			"alkira_service_checkpoint":             resourceAlkiraCheckpoint(),
			"alkira_service_cisco_ftdv":             resourceAlkiraServiceCiscoFTDv(),
			"alkira_service_fortinet":               resourceAlkiraServiceFortinet(),
			"alkira_service_infoblox":               resourceAlkiraInfoblox(),
			"alkira_service_zscaler":                resourceAlkiraServiceZscaler(),
			"alkira_service_pan":                    resourceAlkiraServicePan(),
		},
		DataSourcesMap: map[string]*schema.Resource{
			"alkira_billing_tag":                  dataSourceAlkiraBillingTag(),
			"alkira_byoip":                        dataSourceAlkiraByoip(),
			"alkira_credential":                   dataSourceAlkiraCredential(),
			"alkira_connector_akamai_prolexic":    dataSourceAlkiraConnectorAkamaiProlexic(),
			"alkira_connector_aruba_edge":         dataSourceAlkiraConnectorArubaEdge(),
			"alkira_connector_azure_vnet":         dataSourceAlkiraConnectorAzureVnet(),
			"alkira_connector_azure_expressroute": dataSourceAlkiraConnectorAzureExpressRoute(),
			"alkira_connector_aws_vpc":            dataSourceAlkiraConnectorAwsVpc(),
			"alkira_connector_ipsec":              dataSourceAlkiraConnectorIpsec(),
			"alkira_connector_ipsec_adv":          dataSourceAlkiraConnectorIpsecAdv(),
			"alkira_connector_cisco_sdwan":        dataSourceAlkiraConnectorCiscoSdwan(),
			"alkira_connector_internet_exit":      dataSourceAlkiraConnectorInternetExit(),
			"alkira_connector_oci_vcn":            dataSourceAlkiraConnectorOciVcn(),
			"alkira_connector_gcp_vpc":            dataSourceAlkiraConnectorGcpVpc(),
			"alkira_group":                        dataSourceAlkiraGroup(),
			"alkira_group_user":                   dataSourceAlkiraGroupUser(),
			"alkira_ip_reservation":               dataSourceAlkiraIpReservation(),
			"alkira_list_as_path":                 dataSourceAlkiraListAsPath(),
			"alkira_list_community":               dataSourceAlkiraListCommunity(),
			"alkira_list_extended_community":      dataSourceAlkiraListExtendedCommunity(),
			"alkira_list_global_cidr":             dataSourceAlkiraListGlobalCidr(),
			"alkira_policy":                       dataSourceAlkiraPolicy(),
			"alkira_policy_nat_rule":              dataSourceAlkiraPolicyNatRule(),
			"alkira_policy_prefix_list":           dataSourceAlkiraPolicyPrefixList(),
			"alkira_policy_rule":                  dataSourceAlkiraPolicyRule(),
			"alkira_policy_rule_list":             dataSourceAlkiraPolicyRuleList(),
			"alkira_segment":                      dataSourceAlkiraSegment(),
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
	alkiraClient, err := alkira.NewAlkiraClient(
		d.Get("portal").(string),
		d.Get("username").(string),
		d.Get("password").(string),
		d.Get("api_key").(string),
		d.Get("provision").(bool),
		d.Get("auth").(string))

	if err != nil {
		log.Printf("[ERROR] failed to initialize alkira provider, please check your credential and portal URI.")
		return nil, err
	}

	return alkiraClient, nil
}
