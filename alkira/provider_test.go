package alkira

import (
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

var testAccProviders map[string]terraform.ResourceProvider
var testAccProvider *schema.Provider

var testAccProvidersVersionValidation map[string]terraform.ResourceProvider
var testAccProviderVersionValidation *schema.Provider

func init() {
	testAccProvider = Provider().(*schema.Provider)
	testAccProviders = map[string]terraform.ResourceProvider{
		"alkira": testAccProvider,
	}

	testAccProviderVersionValidation = Provider().(*schema.Provider)
	testAccProviderVersionValidation.ConfigureFunc = alkiraConfigureWithoutVersionValidation
	testAccProvidersVersionValidation = map[string]terraform.ResourceProvider{
		"alkira": testAccProviderVersionValidation,
	}
}

func TestProvider(t *testing.T) {
	if err := Provider().(*schema.Provider).InternalValidate(); err != nil {
		t.Fatalf("err: %s", err)
	}
}

func TestProvider_impl(t *testing.T) {
	var _ = Provider()
}

func testAccPreCheck(t *testing.T) {
	if v := os.Getenv("ALKIRA_PORTAL"); v == "" {
		t.Fatal("ALKIRA_PORTAL must be set for acceptance tests.")
	}
	if v := os.Getenv("ALKIRA_USERNAME"); v == "" {
		t.Fatal("ALKIRA_USERNAME must be set for acceptance tests.")
	}
	if v := os.Getenv("ALKIRA_PASSWORD"); v == "" {
		t.Fatal("ALKIRA_PASSWORD must be set for acceptance tests.")
	}
}
