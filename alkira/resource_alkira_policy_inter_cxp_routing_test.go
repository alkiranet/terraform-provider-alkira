package alkira

import (
	"encoding/json"
	"testing"

	"github.com/alkiranet/alkira-client-go/alkira"
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

	t.Run("rule.match_service_group_ids", func(t *testing.T) {
		s, exists := ruleSchema["match_service_group_ids"]
		require.True(t, exists)
		assert.Equal(t, schema.TypeSet, s.Type)
		assert.Equal(t, schema.TypeInt, s.Elem.(*schema.Schema).Type)
		assert.True(t, s.Optional)
		assert.Zero(t, s.MaxItems)
	})
	t.Run("rule.set_reset_local_pref_to_default", func(t *testing.T) {
		s, exists := ruleSchema["set_reset_local_pref_to_default"]
		require.True(t, exists)
		assert.Equal(t, schema.TypeString, s.Type)
		assert.True(t, s.Optional)
		require.NotNil(t, s.ValidateFunc)
		for _, v := range []string{"true", "false"} {
			_, errs := s.ValidateFunc(v, "set_reset_local_pref_to_default")
			assert.Empty(t, errs, "value %q must be valid", v)
		}
		_, errs := s.ValidateFunc("yes", "set_reset_local_pref_to_default")
		assert.NotEmpty(t, errs)
	})

	ruleRequiredFields := []string{"name", "action"}
	for _, field := range ruleRequiredFields {
		t.Run("rule."+field+"_required", func(t *testing.T) {
			s, exists := ruleSchema[field]
			assert.True(t, exists, "rule field %q must exist", field)
			assert.True(t, s.Required, "rule field %q must be required", field)
		})
	}
}

// interCxpRoutingRuleFromConfig builds one rule through TestResourceData so
// that set fields are real *schema.Set values, as at plan time.
func interCxpRoutingRuleFromConfig(t *testing.T, rule map[string]interface{}) map[string]interface{} {
	t.Helper()
	d := resourceAlkiraPolicyInterCxpRouting().TestResourceData()
	require.NoError(t, d.Set("rule", []interface{}{rule}))
	return d.Get("rule").([]interface{})[0].(map[string]interface{})
}

func TestInterCxpRouting_ValidateMatchAllWithServiceGroups(t *testing.T) {
	rule := interCxpRoutingRuleFromConfig(t, map[string]interface{}{
		"name":                    "bad-rule",
		"action":                  "ALLOW",
		"match_all":               true,
		"match_service_group_ids": []interface{}{7},
	})

	err := validateInterCxpRoutingRules([]interface{}{rule})
	require.Error(t, err)
	assert.Contains(t, err.Error(), "match_all cannot be combined with match_service_group_ids")
}

func TestInterCxpRouting_ValidateDenyWithResetLocalPref(t *testing.T) {
	rule := interCxpRoutingRuleFromConfig(t, map[string]interface{}{
		"name":                            "deny-with-reset",
		"action":                          "DENY",
		"match_service_group_ids":         []interface{}{7},
		"set_reset_local_pref_to_default": "false",
	})

	err := validateInterCxpRoutingRules([]interface{}{rule})
	require.Error(t, err)
	assert.Contains(t, err.Error(), "set_reset_local_pref_to_default cannot be set when action is DENY")
}

func TestInterCxpRouting_ValidateServiceGroupRuleWithReset(t *testing.T) {
	rule := interCxpRoutingRuleFromConfig(t, map[string]interface{}{
		"name":                            "service-group-prepend",
		"action":                          "ALLOW",
		"match_service_group_ids":         []interface{}{7},
		"set_as_path_prepend":             "65100 65100",
		"set_reset_local_pref_to_default": "true",
	})

	assert.NoError(t, validateInterCxpRoutingRules([]interface{}{rule}))
}

func TestInterCxpRouting_ExpandMatchServiceGroupIds(t *testing.T) {
	t.Run("populated", func(t *testing.T) {
		rule := interCxpRoutingRuleFromConfig(t, map[string]interface{}{
			"name":                    "r",
			"action":                  "ALLOW",
			"match_service_group_ids": []interface{}{1, 2},
		})

		match, err := expandPolicyInterCxpRoutingRuleMatch(rule)
		require.NoError(t, err)
		assert.ElementsMatch(t, []int{1, 2}, match.ServiceGroupIds)
	})

	t.Run("empty", func(t *testing.T) {
		rule := interCxpRoutingRuleFromConfig(t, map[string]interface{}{
			"name":                    "r",
			"action":                  "ALLOW",
			"match_all":               true,
			"match_service_group_ids": []interface{}{},
		})

		match, err := expandPolicyInterCxpRoutingRuleMatch(rule)
		require.NoError(t, err)
		assert.Nil(t, match.ServiceGroupIds)

		b, err := json.Marshal(match)
		require.NoError(t, err)
		assert.NotContains(t, string(b), "serviceGroupIds")
	})
}

func TestInterCxpRouting_ExpandSetResetLocalPrefToDefault(t *testing.T) {
	tests := []struct {
		name     string
		value    string
		expected *bool
		jsonPart string
	}{
		{name: "unset", value: "", expected: nil},
		{name: "true", value: "true", expected: boolPtr(true), jsonPart: `"resetLocalPrefToDefault":true`},
		{name: "false", value: "false", expected: boolPtr(false), jsonPart: `"resetLocalPrefToDefault":false`},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			rule := interCxpRoutingRuleFromConfig(t, map[string]interface{}{
				"name":                            "r",
				"action":                          "ALLOW",
				"match_prefix_list_ids":           []interface{}{10},
				"set_as_path_prepend":             "65100",
				"set_reset_local_pref_to_default": tc.value,
			})

			set, err := expandPolicyInterCxpRoutingRuleSet(rule)
			require.NoError(t, err)
			require.NotNil(t, set)
			assert.Equal(t, tc.expected, set.ResetLocalPrefToDefault)

			b, err := json.Marshal(set)
			require.NoError(t, err)
			if tc.jsonPart == "" {
				assert.NotContains(t, string(b), "resetLocalPrefToDefault")
			} else {
				assert.Contains(t, string(b), tc.jsonPart)
			}
		})
	}
}

func TestInterCxpRouting_ExpandSetOnlyResetLocalPrefIsNotEmpty(t *testing.T) {
	rule := interCxpRoutingRuleFromConfig(t, map[string]interface{}{
		"name":                            "r",
		"action":                          "ALLOW",
		"match_prefix_list_ids":           []interface{}{10},
		"set_reset_local_pref_to_default": "false",
	})

	set, err := expandPolicyInterCxpRoutingRuleSet(rule)
	require.NoError(t, err)
	require.NotNil(t, set)
	assert.Equal(t, boolPtr(false), set.ResetLocalPrefToDefault)
}

func TestInterCxpRouting_FlattenServiceGroupsAndResetLocalPref(t *testing.T) {
	tests := []struct {
		name     string
		value    *bool
		expected string
	}{
		{name: "nil", value: nil, expected: ""},
		{name: "true", value: boolPtr(true), expected: "true"},
		{name: "false", value: boolPtr(false), expected: "false"},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			d := resourceAlkiraPolicyInterCxpRouting().TestResourceData()
			in := []alkira.InterCxpRoutePolicyRule{{
				Name:   "r",
				Action: "ALLOW",
				Match:  alkira.InterCxpRoutePolicyRuleMatch{ServiceGroupIds: []int{1, 2}},
				Set: &alkira.InterCxpRoutePolicyRuleSet{
					AsPathPrepend:           "65100",
					ResetLocalPrefToDefault: tc.value,
				},
			}}

			require.NoError(t, setPolicyInterCxpRoutingRules(in, d))

			assert.Equal(t, tc.expected, d.Get("rule.0.set_reset_local_pref_to_default"))
			ids := d.Get("rule.0.match_service_group_ids").(*schema.Set)
			assert.ElementsMatch(t, []interface{}{1, 2}, ids.List())
		})
	}
}
