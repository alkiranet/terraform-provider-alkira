package alkira

import (
	"fmt"
	"reflect"

	"github.com/alkiranet/alkira-client-go/alkira"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// validateInterCxpRoutingRules validates rule-level constraints for
// inter-CXP routing policies. Called from CustomizeDiff at plan time.
func validateInterCxpRoutingRules(rules []interface{}) error {
	for i, r := range rules {
		rule := r.(map[string]interface{})

		// match_all is exclusive — cannot combine with other match conditions.
		if matchAll, ok := rule["match_all"].(bool); ok && matchAll {
			matchFields := []string{
				"match_prefix_list_ids",
				"match_community_list_ids",
				"match_extended_community_list_ids",
				"match_as_path_list_ids",
				"match_segment_resource_ids",
				"match_group_ids",
			}
			for _, field := range matchFields {
				if v, ok := rule[field].(*schema.Set); ok && v.Len() > 0 {
					return fmt.Errorf("rule[%d] %q: match_all cannot be combined with %s", i, rule["name"], field)
				}
			}
		}

		// match_group_ids and match_segment_resource_ids are mutually exclusive.
		groupIds, hasGroups := rule["match_group_ids"].(*schema.Set)
		segResIds, hasSegRes := rule["match_segment_resource_ids"].(*schema.Set)
		if hasGroups && hasSegRes && groupIds.Len() > 0 && segResIds.Len() > 0 {
			return fmt.Errorf("rule[%d] %q: match_group_ids and match_segment_resource_ids are mutually exclusive", i, rule["name"])
		}

		// match_group_ids and match_extended_community_list_ids are mutually exclusive.
		extCommIds, hasExtComm := rule["match_extended_community_list_ids"].(*schema.Set)
		if hasGroups && hasExtComm && groupIds.Len() > 0 && extCommIds.Len() > 0 {
			return fmt.Errorf("rule[%d] %q: match_group_ids and match_extended_community_list_ids are mutually exclusive", i, rule["name"])
		}

		// DENY action cannot have set parameters.
		if action, ok := rule["action"].(string); ok && action == "DENY" {
			setFields := []string{"set_as_path_prepend", "set_community", "set_extended_community"}
			for _, field := range setFields {
				if v, ok := rule[field].(string); ok && v != "" {
					return fmt.Errorf("rule[%d] %q: %s cannot be set when action is DENY", i, rule["name"], field)
				}
			}
		}
	}

	return nil
}

func expandPolicyInterCxpRoutingRuleMatch(in map[string]interface{}) (*alkira.InterCxpRoutePolicyRuleMatch, error) {

	match := alkira.InterCxpRoutePolicyRuleMatch{}

	if v, ok := in["match_all"].(bool); ok {
		match.All = v
	}
	if v, ok := in["match_prefix_list_ids"].(*schema.Set); ok {
		match.PrefixListIds = convertTypeSetToIntList(v)
		if len(match.PrefixListIds) == 0 {
			match.PrefixListIds = nil
		}
	}
	if v, ok := in["match_community_list_ids"].(*schema.Set); ok {
		match.CommunityListIds = convertTypeSetToIntList(v)
		if len(match.CommunityListIds) == 0 {
			match.CommunityListIds = nil
		}
	}
	if v, ok := in["match_extended_community_list_ids"].(*schema.Set); ok {
		match.ExtendedCommunityListIds = convertTypeSetToIntList(v)
		if len(match.ExtendedCommunityListIds) == 0 {
			match.ExtendedCommunityListIds = nil
		}
	}
	if v, ok := in["match_as_path_list_ids"].(*schema.Set); ok {
		match.AsPathListIds = convertTypeSetToIntList(v)
		if len(match.AsPathListIds) == 0 {
			match.AsPathListIds = nil
		}
	}
	if v, ok := in["match_segment_resource_ids"].(*schema.Set); ok {
		match.SegmentResourceIds = convertTypeSetToIntList(v)
		if len(match.SegmentResourceIds) == 0 {
			match.SegmentResourceIds = nil
		}
	}
	if v, ok := in["match_group_ids"].(*schema.Set); ok {
		match.ConnectorGroupIds = convertTypeSetToIntList(v)
		if len(match.ConnectorGroupIds) == 0 {
			match.ConnectorGroupIds = nil
		}
	}

	return &match, nil
}

func expandPolicyInterCxpRoutingRuleSet(in map[string]interface{}) (*alkira.InterCxpRoutePolicyRuleSet, error) {

	set := alkira.InterCxpRoutePolicyRuleSet{}

	if v, ok := in["set_as_path_prepend"].(string); ok {
		set.AsPathPrepend = v
	}
	if v, ok := in["set_community"].(string); ok {
		set.Community = v
	}
	if v, ok := in["set_extended_community"].(string); ok {
		set.ExtendedCommunity = v
	}

	if reflect.DeepEqual(set, alkira.InterCxpRoutePolicyRuleSet{}) {
		return nil, nil
	}

	return &set, nil
}

func expandPolicyInterCxpRoutingRule(in []interface{}) ([]alkira.InterCxpRoutePolicyRule, error) {

	if len(in) == 0 {
		return nil, nil
	}

	rules := make([]alkira.InterCxpRoutePolicyRule, len(in))

	for i, ruleInput := range in {

		rule := alkira.InterCxpRoutePolicyRule{}
		input := ruleInput.(map[string]interface{})

		if v, ok := input["action"].(string); ok {
			rule.Action = v
		}
		if v, ok := input["name"].(string); ok {
			rule.Name = v
		}

		match, err := expandPolicyInterCxpRoutingRuleMatch(input)
		if err != nil {
			return nil, err
		}
		rule.Match = *match

		set, err := expandPolicyInterCxpRoutingRuleSet(input)
		if err != nil {
			return nil, err
		}
		rule.Set = set

		rules[i] = rule
	}

	return rules, nil
}

func setPolicyInterCxpRoutingRules(in []alkira.InterCxpRoutePolicyRule, d *schema.ResourceData) error {

	if len(in) == 0 {
		return nil
	}

	rules := make([]map[string]interface{}, len(in))

	for i, rule := range in {
		r := map[string]interface{}{
			"name":        rule.Name,
			"sequence_no": rule.SequenceNo,
			"action":      rule.Action,
			"match_all":   rule.Match.All,
		}

		if rule.Match.PrefixListIds != nil {
			r["match_prefix_list_ids"] = rule.Match.PrefixListIds
		}
		if rule.Match.CommunityListIds != nil {
			r["match_community_list_ids"] = rule.Match.CommunityListIds
		}
		if rule.Match.ExtendedCommunityListIds != nil {
			r["match_extended_community_list_ids"] = rule.Match.ExtendedCommunityListIds
		}
		if rule.Match.AsPathListIds != nil {
			r["match_as_path_list_ids"] = rule.Match.AsPathListIds
		}
		if rule.Match.SegmentResourceIds != nil {
			r["match_segment_resource_ids"] = rule.Match.SegmentResourceIds
		}
		if rule.Match.ConnectorGroupIds != nil {
			r["match_group_ids"] = rule.Match.ConnectorGroupIds
		}

		if rule.Set != nil {
			r["set_as_path_prepend"] = rule.Set.AsPathPrepend
			r["set_community"] = rule.Set.Community
			r["set_extended_community"] = rule.Set.ExtendedCommunity
		}

		rules[i] = r
	}

	d.Set("rule", rules)
	return nil
}
