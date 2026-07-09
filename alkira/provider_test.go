package alkira

import (
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

var testAccProviders map[string]*schema.Provider
var testAccProvider *schema.Provider

var testAccProvidersVersionValidation map[string]*schema.Provider
var testAccProviderVersionValidation *schema.Provider

func init() {
	testAccProvider = Provider()
	testAccProviders = map[string]*schema.Provider{
		"alkira": testAccProvider,
	}

	testAccProviderVersionValidation = Provider()
	// testAccProviderVersionValidation.ConfigureFunc = alkiraConfigureWithoutVersionValidation // Function not available
	testAccProvidersVersionValidation = map[string]*schema.Provider{
		"alkira": testAccProviderVersionValidation,
	}
}

func TestProvider(t *testing.T) {
	if err := Provider().InternalValidate(); err != nil {
		t.Fatalf("err: %s", err)
	}
}

func TestProvider_impl(t *testing.T) {
	var _ = Provider()
}

func TestProviderSerializationEnabledDefault(t *testing.T) {
	const key = "ALKIRA_API_SERIALIZATION_ENABLED"
	t.Setenv(key, "")
	os.Unsetenv(key)

	s := Provider().Schema["serialization_enabled"]
	if s == nil {
		t.Fatal("serialization_enabled schema not found")
	}

	def, err := s.DefaultValue()
	if err != nil {
		t.Fatalf("error loading default: %s", err)
	}

	if def != true {
		t.Fatalf("expected serialization_enabled default to be true, got %v", def)
	}
}

func TestBoolEnvDefaultFunc(t *testing.T) {
	const key = "ALKIRA_API_SERIALIZATION_ENABLED"
	f := boolEnvDefaultFunc(key, true)

	cases := []struct {
		name    string
		set     bool
		value   string
		want    interface{}
		wantErr bool
	}{
		{name: "unset falls back to default", set: false, want: true},
		{name: "explicit true", set: true, value: "true", want: true},
		{name: "explicit false disables", set: true, value: "false", want: false},
		{name: "numeric one enables", set: true, value: "1", want: true},
		{name: "numeric zero disables", set: true, value: "0", want: false},
		{name: "invalid value errors", set: true, value: "yes-please", wantErr: true},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			if tc.set {
				t.Setenv(key, tc.value)
			} else {
				t.Setenv(key, "")
				os.Unsetenv(key)
			}

			got, err := f()
			if tc.wantErr {
				if err == nil {
					t.Fatalf("expected error, got value %v", got)
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %s", err)
			}
			if got != tc.want {
				t.Fatalf("expected %v, got %v", tc.want, got)
			}
		})
	}
}

// UNUSED: Commented out to suppress linter warnings
// func testAccPreCheck(t *testing.T) {
// 	if v := os.Getenv("ALKIRA_PORTAL"); v == "" {
// 		t.Fatal("ALKIRA_PORTAL must be set for acceptance tests.")
// 	}
// 	if v := os.Getenv("ALKIRA_USERNAME"); v == "" {
// 		t.Fatal("ALKIRA_USERNAME must be set for acceptance tests.")
// 	}
// 	if v := os.Getenv("ALKIRA_PASSWORD"); v == "" {
// 		t.Fatal("ALKIRA_PASSWORD must be set for acceptance tests.")
// 	}
// }
