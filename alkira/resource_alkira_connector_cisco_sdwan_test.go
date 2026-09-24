package alkira

import (
	"context"
	"testing"

	"github.com/alkiranet/alkira-client-go/alkira"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCiscoSdwanAllowListSchema(t *testing.T) {
	resourceSchema := resourceAlkiraConnectorCiscoSdwan().Schema

	t.Run("field exists and is Optional TypeSet of String", func(t *testing.T) {
		field, ok := resourceSchema["allow_list"]
		require.True(t, ok, "allow_list must be present in schema")
		assert.Equal(t, schema.TypeSet, field.Type)
		assert.True(t, field.Optional)
		assert.False(t, field.Required)

		elem, ok := field.Elem.(*schema.Schema)
		require.True(t, ok, "allow_list Elem must be a *schema.Schema")
		assert.Equal(t, schema.TypeString, elem.Type)
		assert.NotNil(t, elem.ValidateFunc, "allow_list elements must be validated")
	})
}

func TestCiscoSdwanAllowListValidator(t *testing.T) {
	validate := resourceAlkiraConnectorCiscoSdwan().Schema["allow_list"].Elem.(*schema.Schema).ValidateFunc

	valid := []string{"10.0.0.0/24", "192.168.1.0/24", "10.0.0.5", "10.0.0.5/32"}
	for _, v := range valid {
		_, errs := validate(v, "allow_list")
		assert.Emptyf(t, errs, "expected %q to be accepted (IPv4 CIDR or IPv4 IP)", v)
	}

	// Rejected at plan time: garbage, out-of-range masks, IPv6, and
	// host-bits-set CIDRs (not silently masked).
	invalid := []string{"not-an-ip", "10.0.0.0/33", "999.0.0.1", "2001:db8::/32", "::1", "::ffff:1.2.3.4", "10.0.0.5/24"}
	for _, v := range invalid {
		_, errs := validate(v, "allow_list")
		assert.NotEmptyf(t, errs, "expected %q to be rejected", v)
	}
}

// TestCiscoSdwanAllowListExpand documents the request-expansion contract:
// a TypeSet is converted with convertTypeSetToStringList (not the List helper),
// which de-duplicates entries.
func TestCiscoSdwanAllowListExpand(t *testing.T) {
	set := schema.NewSet(schema.HashString, []interface{}{
		"10.0.0.0/24", "10.0.0.5", "10.0.0.0/24", // duplicate collapses
	})

	got := convertTypeSetToStringList(set)

	assert.ElementsMatch(t, []string{"10.0.0.0/24", "10.0.0.5"}, got)
}

func TestCiscoSdwanAllowListCat8kGate(t *testing.T) {
	client := &alkira.AlkiraClient{}

	baseConfig := func(connType string, allowList []interface{}) map[string]interface{} {
		return map[string]interface{}{
			"name":       "test",
			"cxp":        "US-WEST",
			"size":       "SMALL",
			"version":    "20.9.1",
			"type":       connType,
			"allow_list": allowList,
		}
	}

	tests := []struct {
		name      string
		connType  string
		allowList []interface{}
		wantError bool
	}{
		{"cat8k with allow_list is allowed", "CAT8000V", []interface{}{"10.0.0.0/24"}, false},
		{"non-cat8k with allow_list is rejected", "CSR", []interface{}{"10.0.0.0/24"}, true},
		{"non-cat8k without allow_list is allowed", "CSR", []interface{}{}, false},
		{"cat8k without allow_list is allowed", "CAT8000V", []interface{}{}, false},
	}

	resource := resourceAlkiraConnectorCiscoSdwan()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg := terraform.NewResourceConfigRaw(baseConfig(tt.connType, tt.allowList))
			_, err := resource.Diff(context.Background(), nil, cfg, client)

			if tt.wantError {
				require.Error(t, err)
				assert.Contains(t, err.Error(), "CAT8000V")
			} else {
				require.NoError(t, err)
			}
		})
	}
}
