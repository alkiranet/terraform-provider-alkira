package alkira

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestJuniperSdwanAllowListSchema(t *testing.T) {
	resourceSchema := resourceAlkiraConnectorJuniperSdwan().Schema

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

func TestJuniperSdwanAllowListValidator(t *testing.T) {
	validate := resourceAlkiraConnectorJuniperSdwan().Schema["allow_list"].Elem.(*schema.Schema).ValidateFunc

	valid := []string{"10.0.0.0/24", "192.168.1.0/24", "10.0.0.5", "10.0.0.5/32"}
	for _, v := range valid {
		_, errs := validate(v, "allow_list")
		assert.Emptyf(t, errs, "expected %q to be accepted (IPv4 CIDR or IPv4 IP)", v)
	}

	invalid := []string{"not-an-ip", "10.0.0.0/33", "999.0.0.1", "2001:db8::/32", "::1", "::ffff:1.2.3.4", "10.0.0.5/24"}
	for _, v := range invalid {
		_, errs := validate(v, "allow_list")
		assert.NotEmptyf(t, errs, "expected %q to be rejected", v)
	}
}

func TestJuniperSdwanAllowListExpand(t *testing.T) {
	t.Run("populated set produces expected slice", func(t *testing.T) {
		set := schema.NewSet(schema.HashString, []interface{}{
			"10.0.0.0/24", "10.0.0.5", "10.0.0.0/24", // duplicate collapses
		})

		got := convertTypeSetToStringList(set)
		assert.ElementsMatch(t, []string{"10.0.0.0/24", "10.0.0.5"}, got)
	})

	t.Run("empty set produces empty slice", func(t *testing.T) {
		set := schema.NewSet(schema.HashString, []interface{}{})
		got := convertTypeSetToStringList(set)
		assert.Empty(t, got)
	})
}
