package alkira

import (
	"testing"

	"github.com/alkiranet/alkira-client-go/alkira"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestPolicyNatFieldNameMatch(t *testing.T) {
	t.Run("schema uses allow_overlapping_translated_source_addresses field name", func(t *testing.T) {
		resourceSchema := resourceAlkiraPolicyNat().Schema

		// Verify "allow_overlapping_translated_source_addresses" field exists in schema (plural)
		field, exists := resourceSchema["allow_overlapping_translated_source_addresses"]
		assert.True(t, exists, "Schema must have 'allow_overlapping_translated_source_addresses' field (plural)")
		assert.NotNil(t, field, "allow_overlapping_translated_source_addresses field must not be nil")

		// Verify "allow_overlapping_translated_source_address" (singular - the bug) does NOT exist
		_, wrongFieldExists := resourceSchema["allow_overlapping_translated_source_address"]
		assert.False(t, wrongFieldExists, "Schema should NOT have 'allow_overlapping_translated_source_address' field (bug was using singular)")
	})
}

func TestPolicyNatStateHandling(t *testing.T) {
	t.Run("allow_overlapping_translated_source_addresses data can be saved and retrieved", func(t *testing.T) {
		resourceSchema := resourceAlkiraPolicyNat().Schema
		d := schema.TestResourceDataRaw(t, resourceSchema, map[string]interface{}{
			"name":       "test-policy",
			"segment_id": "seg-123",
		})

		// Set with correct field name (plural)
		testValue := true
		err := d.Set("allow_overlapping_translated_source_addresses", testValue)
		require.NoError(t, err, "Should be able to set allow_overlapping_translated_source_addresses data")

		// Verify we can read it back
		valueFromState := d.Get("allow_overlapping_translated_source_addresses")
		require.NotNil(t, valueFromState, "Should be able to get allow_overlapping_translated_source_addresses data")
		assert.Equal(t, testValue, valueFromState.(bool), "Value should match")
	})

	t.Run("verify wrong field name not accessible", func(t *testing.T) {
		resourceSchema := resourceAlkiraPolicyNat().Schema
		d := schema.TestResourceDataRaw(t, resourceSchema, map[string]interface{}{
			"name":       "test-policy",
			"segment_id": "seg-123",
		})

		// Try to set with singular form (the bug)
		err := d.Set("allow_overlapping_translated_source_address", true)
		// This should fail or be ignored because singular form is not in schema
		if err == nil {
			wrongData := d.Get("allow_overlapping_translated_source_address")
			t.Logf("Setting 'allow_overlapping_translated_source_address' (singular - wrong) returned: %v", wrongData)
		}

		// Verify correct field name (plural) works
		err = d.Set("allow_overlapping_translated_source_addresses", true)
		require.NoError(t, err)

		correctData := d.Get("allow_overlapping_translated_source_addresses")
		require.NotNil(t, correctData, "Correct field name 'allow_overlapping_translated_source_addresses' (plural) must work")
		assert.Equal(t, true, correctData.(bool))
	})
}

func TestPolicyNatFieldNaming(t *testing.T) {
	t.Run("document singular vs plural to prevent regression", func(t *testing.T) {
		// This test documents the singular/plural distinction to prevent reintroduction
		correctPlural := "allow_overlapping_translated_source_addresses"
		incorrectSingular := "allow_overlapping_translated_source_address"

		assert.NotEqual(t, correctPlural, incorrectSingular,
			"Singular/plural documented: 'addresses' (plural) != 'address' (singular)")

		// Verify the schema uses the plural form
		resourceSchema := resourceAlkiraPolicyNat().Schema
		_, hasPlural := resourceSchema[correctPlural]
		_, hasSingular := resourceSchema[incorrectSingular]

		assert.True(t, hasPlural, "Schema must use plural form 'addresses'")
		assert.False(t, hasSingular, "Schema must NOT have singular form 'address'")
	})
}

// TPS stores a NULL category for any NAT policy or rule created without an
// explicit one, and omits the key from the GET body, so the client-go value
// arrives as "". Read must translate that back to the schema default.
func TestNatCategoryOrDefault(t *testing.T) {
	cases := []struct {
		name     string
		apiValue string
		expected string
	}{
		{"omitted by the API becomes the schema default", "", "DEFAULT"},
		{"explicit DEFAULT is preserved", "DEFAULT", "DEFAULT"},
		{"INTERNET_CONNECTOR is preserved", "INTERNET_CONNECTOR", "INTERNET_CONNECTOR"},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			assert.Equal(t, tc.expected, natCategoryOrDefault(tc.apiValue))
		})
	}
}

// The normalization constant is only correct while it matches the schema
// default -- if they diverge, an imported resource diffs forever.
func TestNatCategoryOrDefaultMatchesSchemaDefault(t *testing.T) {
	for name, resourceSchema := range map[string]map[string]*schema.Schema{
		"alkira_policy_nat":      resourceAlkiraPolicyNat().Schema,
		"alkira_policy_nat_rule": resourceAlkiraPolicyNatRule().Schema,
	} {
		t.Run(name, func(t *testing.T) {
			category, exists := resourceSchema["category"]
			require.True(t, exists, "schema must declare category")
			assert.Equal(t, natCategoryOrDefault(""), category.Default,
				"normalization constant must equal the schema default")
		})
	}
}

// direction is Optional with no Default, so it must be written through
// unnormalized: Terraform's unset value and the API's "" already coincide.
func TestPolicyNatRuleDirectionHasNoSchemaDefault(t *testing.T) {
	direction, exists := resourceAlkiraPolicyNatRule().Schema["direction"]
	require.True(t, exists, "schema must declare direction")
	assert.Nil(t, direction.Default, "direction must not gain a default")
}

// Characterizes the import round-trip: every attribute the API returns has to
// land in state, or the plan taken right after `terraform import` is dirty.
func TestSetNatRuleFieldsPopulatesStateFromAPI(t *testing.T) {
	t.Run("category omitted by the API lands as DEFAULT", func(t *testing.T) {
		d := schema.TestResourceDataRaw(t, resourceAlkiraPolicyNatRule().Schema, map[string]interface{}{})

		setNatRuleFields(d, &alkira.NatPolicyRule{Name: "rule-1", Category: ""})

		assert.Equal(t, "DEFAULT", d.Get("category"))
	})

	t.Run("direction is read back into state", func(t *testing.T) {
		d := schema.TestResourceDataRaw(t, resourceAlkiraPolicyNatRule().Schema, map[string]interface{}{})

		setNatRuleFields(d, &alkira.NatPolicyRule{Name: "rule-1", Direction: "INBOUND"})

		assert.Equal(t, "INBOUND", d.Get("direction"))
	})

	t.Run("an explicit category survives the round-trip", func(t *testing.T) {
		d := schema.TestResourceDataRaw(t, resourceAlkiraPolicyNatRule().Schema, map[string]interface{}{})

		setNatRuleFields(d, &alkira.NatPolicyRule{Name: "rule-1", Category: "INTERNET_CONNECTOR"})

		assert.Equal(t, "INTERNET_CONNECTOR", d.Get("category"))
	})
}
