package alkira

import (
	"reflect"

	"github.com/alkiranet/alkira-client-go/alkira"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func expandPolicyRoutingRuleMatch(in map[string]interface{}) (*alkira.RoutePolicyRulesMatch, error) {

	match := alkira.RoutePolicyRulesMatch{}

	if v, ok := in["match_all"].(bool); ok {
		match.All = v
	}
	if v, ok := in["match_as_path_list_ids"].([]interface{}); ok {
		match.AsPathListIds = convertTypeListToIntList(v)
	}
	if v, ok := in["match_community_list_ids"].([]interface{}); ok {
		match.CommunityListIds = convertTypeListToIntList(v)
	}
	if v, ok := in["match_extended_community_list_ids"].([]interface{}); ok {
		match.ExtendedCommunityListIds = convertTypeListToIntList(v)

		if len(match.ExtendedCommunityListIds) == 0 {
			match.ExtendedCommunityListIds = nil
		}
	}
	if v, ok := in["match_prefix_list_ids"].([]interface{}); ok {
		match.PrefixListIds = convertTypeListToIntList(v)

		if len(match.PrefixListIds) == 0 {
			match.PrefixListIds = nil
		}
	}
	if v, ok := in["match_cxps"].([]interface{}); ok {
		match.Cxps = convertTypeListToStringList(v)
		if len(match.Cxps) == 0 {
			match.Cxps = nil
		}
	}
	if v, ok := in["match_segment_resource_ids"].(*schema.Set); ok {
		match.SegmentResourceIds = convertTypeSetToIntList(v)
		if len(match.SegmentResourceIds) == 0 {
			match.SegmentResourceIds = nil
		}
	}
	if v, ok := in["match_group_ids"].([]interface{}); ok {
		match.ConnectorGroupIds = convertTypeListToIntList(v)
		if len(match.ConnectorGroupIds) == 0 {
			match.ConnectorGroupIds = nil
		}
	}

	return &match, nil
}

// expandPolicyRoutingRuleSet expand the "set" section of the policy routing rule
func expandPolicyRoutingRuleSet(in map[string]interface{}) (*alkira.RoutePolicyRulesSet, error) {

	set := alkira.RoutePolicyRulesSet{}

	if v, ok := in["set_as_path_prepend"].(string); ok {
		set.AsPathPrepend = v
	}
	if v, ok := in["set_community"].(string); ok {
		set.Community = v
	}
	if v, ok := in["set_extended_community"].(string); ok {
		set.ExtendedCommunity = v
	}
	if v, ok := in["set_med"].(int); ok {
		set.Med = v
	}

	return &set, nil
}

// expandPolicyRoutingRuleInterCxpRoutesRedistribution expand the
//
//	"inter_cxp_routes_redistribution" section of the policy routing rule
func expandPolicyRoutingRuleInterCxpRoutesRedistribution(in map[string]interface{}) (*alkira.RoutePolicyRulesInterCxpRoutesRedistribution, error) {

	distrib := alkira.RoutePolicyRulesInterCxpRoutesRedistribution{}

	if v, ok := in["routes_distribution_type"].(string); ok {
		distrib.DistributionType = v
	}
	if v, ok := in["routes_distribution_as_secondary"].(bool); ok {
		distrib.RedistributeAsSecondary = v
	}
	if v, ok := in["routes_distribution_restricted_cxps"].([]interface{}); ok {
		distrib.RestrictedCxps = convertTypeListToStringList(v)

		if len(distrib.RestrictedCxps) == 0 {
			distrib.RestrictedCxps = nil
		}
	}

	if reflect.DeepEqual(distrib, alkira.RoutePolicyRulesInterCxpRoutesRedistribution{}) {
		return nil, nil
	}

	return &distrib, nil
}

// expandPolicyRoutingRule expanding the "rule" sections of the routing policy
func expandPolicyRoutingRule(in []interface{}) ([]alkira.RoutePolicyRules, error) {

	if in == nil || len(in) == 0 {
		return nil, nil
	}

	rules := make([]alkira.RoutePolicyRules, len(in))

	for i, ruleInput := range in {

		rule := alkira.RoutePolicyRules{}
		input := ruleInput.(map[string]interface{})

		if v, ok := input["action"].(string); ok {
			rule.Action = v
		}
		if v, ok := input["name"].(string); ok {
			rule.Name = v
		}

		match, err := expandPolicyRoutingRuleMatch(input)

		if err != nil {
			return nil, err
		}

		rule.Match = *match

		set, err := expandPolicyRoutingRuleSet(input)
		if err != nil {
			return nil, err
		}

		rule.Set = set

		distribution, err := expandPolicyRoutingRuleInterCxpRoutesRedistribution(input)

		if err != nil {
			return nil, err
		}

		rule.InterCxpRoutesRedistribution = distribution

		rules[i] = rule
	}

	return rules, nil
}

// setPolicyRoutingRules set rules
func setPolicyRoutingRules(in []alkira.RoutePolicyRules, d *schema.ResourceData) error {

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

		if rule.Match.AsPathListIds != nil {
			r["match_as_path_list_ids"] = rule.Match.AsPathListIds
		}
		if rule.Match.CommunityListIds != nil {
			r["match_community_list_ids"] = rule.Match.CommunityListIds
		}
		if rule.Match.ExtendedCommunityListIds != nil {
			r["match_extended_community_list_ids"] = rule.Match.ExtendedCommunityListIds
		}
		if rule.Match.PrefixListIds != nil {
			r["match_prefix_list_ids"] = rule.Match.PrefixListIds
		}
		if rule.Match.ConnectorGroupIds != nil {
			r["match_group_ids"] = rule.Match.ConnectorGroupIds
		}
		if rule.Match.Cxps != nil {
			r["match_cxps"] = rule.Match.Cxps
		}
		if rule.Match.SegmentResourceIds != nil {
			r["match_segment_resource_ids"] = rule.Match.SegmentResourceIds
		}

		if rule.Set != nil {
			r["set_as_path_prepend"] = rule.Set.AsPathPrepend
			r["set_community"] = rule.Set.Community
			r["set_extended_community"] = rule.Set.ExtendedCommunity
			r["set_med"] = rule.Set.Med
		}

		if rule.InterCxpRoutesRedistribution != nil {
			r["routes_distribution_restricted_cxps"] = rule.InterCxpRoutesRedistribution.RestrictedCxps
			r["routes_distribution_type"] = rule.InterCxpRoutesRedistribution.DistributionType
			r["routes_distribution_as_secondary"] = rule.InterCxpRoutesRedistribution.RedistributeAsSecondary
		}

		rules[i] = r
	}

	d.Set("rule", rules)
	return nil
}
