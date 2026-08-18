package alkira

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// countSensitiveByName walks a resource schema (including nested TypeList/TypeSet
// blocks whose Elem is a *schema.Resource) and, for the given attribute name,
// returns how many occurrences it found and how many of those were marked
// Sensitive. A secret attribute is only masked in `terraform plan`/`apply`
// output and CI logs when Sensitive is true.
func countSensitiveByName(s map[string]*schema.Schema, name string) (found int, sensitive int) {
	for key, attr := range s {
		if key == name {
			found++
			if attr.Sensitive {
				sensitive++
			}
		}
		if attr.Elem != nil {
			if r, ok := attr.Elem.(*schema.Resource); ok {
				f, sen := countSensitiveByName(r.Schema, name)
				found += f
				sensitive += sen
			}
		}
	}
	return found, sensitive
}

// TestSecretResourceAttributesAreSensitive asserts that credential-bearing
// resource attributes are marked Sensitive so their values are not rendered in
// cleartext in the plan/apply diff shown on the operator terminal or captured in
// CI job logs. This fails against code where any listed attribute lacks the flag.
func TestSecretResourceAttributesAreSensitive(t *testing.T) {
	cases := []struct {
		resource string
		schema   map[string]*schema.Schema
		attrs    []string
	}{
		{"alkira_connector_cisco_sdwan", resourceAlkiraConnectorCiscoSdwan().Schema, []string{"password"}},
		{"alkira_connector_fortinet_sdwan", resourceAlkiraConnectorFortinetSdwan().Schema, []string{"password"}},
		{"alkira_connector_ipsec_adv", resourceAlkiraConnectorIPSecAdv().Schema, []string{"preshared_key"}},
		{"alkira_credential_azure_vnet", resourceAlkiraCredentialAzureVnet().Schema, []string{"secret_key"}},
		{"alkira_credential_gcp_vpc", resourceAlkiraCredentialGcpVpc().Schema, []string{"private_key"}},
		{"alkira_service_checkpoint", resourceAlkiraCheckpoint().Schema, []string{"password"}},
		{"alkira_service_cisco_ftdv", resourceAlkiraServiceCiscoFTDv().Schema, []string{"password", "admin_password"}},
		{"alkira_service_fortinet", resourceAlkiraServiceFortinet().Schema, []string{"password", "license_key"}},
		{"alkira_service_infoblox", resourceAlkiraInfoblox().Schema, []string{"password", "shared_secret"}},
		{"alkira_service_pan", resourceAlkiraServicePan().Schema, []string{"auth_key"}},
	}

	for _, c := range cases {
		for _, attr := range c.attrs {
			t.Run(c.resource+"/"+attr, func(t *testing.T) {
				found, sensitive := countSensitiveByName(c.schema, attr)
				if found == 0 {
					t.Fatalf("%s: attribute %q not found in schema", c.resource, attr)
				}
				if sensitive != found {
					t.Fatalf("%s: attribute %q present %d time(s) but only %d marked Sensitive; "+
						"an unmarked secret is printed in cleartext in plan/apply output and CI logs",
						c.resource, attr, found, sensitive)
				}
			})
		}
	}
}
