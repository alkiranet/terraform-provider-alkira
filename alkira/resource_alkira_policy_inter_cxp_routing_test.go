package alkira

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestInterCxpRouting_DirectionImmutable(t *testing.T) {
	r := resourceAlkiraPolicyInterCxpRouting()
	s := r.Schema["direction"]

	// direction is enforced as immutable via CustomizeDiff, not ForceNew
	assert.False(t, s.ForceNew, "direction should not use ForceNew — immutability is enforced in CustomizeDiff")
	assert.NotNil(t, r.CustomizeDiff, "CustomizeDiff must be set to enforce direction immutability")
}

// TestInterCxpRouting_ValidateMatchAllExclusivity uses TestResourceData to
// populate rule fields so that the *schema.Set type assertions in
// validateInterCxpRoutingRules are exercised with real SDK types.
func TestInterCxpRouting_ValidateMatchAllExclusivity(t *testing.T) {
	r := resourceAlkiraPolicyInterCxpRouting()
	d := r.TestResourceData()

	// Set a rule with match_all = true AND match_prefix_list_ids populated.
	d.Set("rule", []interface{}{
		map[string]interface{}{
			"name":                              "bad-rule",
			"action":                            "ALLOW",
			"match_all":                         true,
			"match_prefix_list_ids":             []interface{}{42},
			"match_community_list_ids":          []interface{}{},
			"match_extended_community_list_ids": []interface{}{},
			"match_as_path_list_ids":            []interface{}{},
			"match_segment_resource_ids":        []interface{}{},
			"match_group_ids":                   []interface{}{},
		},
	})

	rules := d.Get("rule").([]interface{})
	err := validateInterCxpRoutingRules(rules)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "match_all cannot be combined with match_prefix_list_ids")
}

func TestInterCxpRouting_ValidateMatchAllAlone(t *testing.T) {
	r := resourceAlkiraPolicyInterCxpRouting()
	d := r.TestResourceData()

	d.Set("rule", []interface{}{
		map[string]interface{}{
			"name":      "allow-all",
			"action":    "ALLOW",
			"match_all": true,
		},
	})

	rules := d.Get("rule").([]interface{})
	err := validateInterCxpRoutingRules(rules)
	assert.NoError(t, err)
}

func TestInterCxpRouting_ValidateMutuallyExclusiveGroupsAndSegmentResources(t *testing.T) {
	r := resourceAlkiraPolicyInterCxpRouting()
	d := r.TestResourceData()

	d.Set("rule", []interface{}{
		map[string]interface{}{
			"name":                       "conflict-rule",
			"action":                     "ALLOW",
			"match_group_ids":            []interface{}{1},
			"match_segment_resource_ids": []interface{}{2},
		},
	})

	rules := d.Get("rule").([]interface{})
	err := validateInterCxpRoutingRules(rules)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "match_group_ids and match_segment_resource_ids are mutually exclusive")
}

func TestInterCxpRouting_ValidateMutuallyExclusiveGroupsAndExtCommunity(t *testing.T) {
	r := resourceAlkiraPolicyInterCxpRouting()
	d := r.TestResourceData()

	d.Set("rule", []interface{}{
		map[string]interface{}{
			"name":                              "conflict-rule",
			"action":                            "ALLOW",
			"match_group_ids":                   []interface{}{1},
			"match_extended_community_list_ids": []interface{}{2},
		},
	})

	rules := d.Get("rule").([]interface{})
	err := validateInterCxpRoutingRules(rules)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "match_group_ids and match_extended_community_list_ids are mutually exclusive")
}

func TestInterCxpRouting_ValidateDenyWithSetParams(t *testing.T) {
	r := resourceAlkiraPolicyInterCxpRouting()
	d := r.TestResourceData()

	d.Set("rule", []interface{}{
		map[string]interface{}{
			"name":                "deny-with-set",
			"action":              "DENY",
			"match_all":           true,
			"set_as_path_prepend": "100 100",
		},
	})

	rules := d.Get("rule").([]interface{})
	err := validateInterCxpRoutingRules(rules)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "set_as_path_prepend cannot be set when action is DENY")
}

func TestInterCxpRouting_ValidateValidRule(t *testing.T) {
	r := resourceAlkiraPolicyInterCxpRouting()
	d := r.TestResourceData()

	d.Set("rule", []interface{}{
		map[string]interface{}{
			"name":                  "allow-prefixes",
			"action":                "ALLOW",
			"match_prefix_list_ids": []interface{}{10},
			"set_community":         "65512:20",
		},
	})

	rules := d.Get("rule").([]interface{})
	err := validateInterCxpRoutingRules(rules)
	assert.NoError(t, err)
}

func TestInterCxpRouting_SchemaFields(t *testing.T) {
	r := resourceAlkiraPolicyInterCxpRouting()

	requiredFields := []string{
		"name", "enabled", "direction", "segment_id",
		"source_cxps", "dest_cxps", "rule",
	}
	for _, field := range requiredFields {
		t.Run(field+"_required", func(t *testing.T) {
			s, exists := r.Schema[field]
			assert.True(t, exists, "field %q must exist", field)
			assert.True(t, s.Required, "field %q must be required", field)
		})
	}

	ruleSchema := r.Schema["rule"].Elem.(*schema.Resource).Schema
	ruleRequiredFields := []string{"name", "action"}
	for _, field := range ruleRequiredFields {
		t.Run("rule."+field+"_required", func(t *testing.T) {
			s, exists := ruleSchema[field]
			assert.True(t, exists, "rule field %q must exist", field)
			assert.True(t, s.Required, "rule field %q must be required", field)
		})
	}
}
