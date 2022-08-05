package alkira

import (
	"fmt"

	"github.com/alkiranet/alkira-client-go/alkira"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func expandPolicyRoutingRuleMatch(in map[string]interface{}) (*alkira.RoutePolicyRulesMatch, error) {

	match := alkira.RoutePolicyRulesMatch{}

	if v, ok := in["all"].(bool); ok {
		match.All = v
	}
	if v, ok := in["as_path_list_ids"].([]interface{}); ok {
		match.AsPathListIds = convertTypeListToIntList(v)
	}
	if v, ok := in["community_list_ids"].([]interface{}); ok {
		match.CommunityListIds = convertTypeListToIntList(v)
	}
	if v, ok := in["extended_community_list_ids"].([]interface{}); ok {
		match.ExtendedCommunityListIds = convertTypeListToIntList(v)
	}
	if v, ok := in["prefix_list_ids"].([]interface{}); ok {
		match.PrefixListIds = convertTypeListToIntList(v)
	}
	if v, ok := in["cxps"].([]interface{}); ok {
		match.Cxps = convertTypeListToStringList(v)
	}
	if v, ok := in["group_ids"].([]interface{}); ok {
		match.ConnectorGroupIds = convertTypeListToIntList(v)
	}

	return &match, nil
}

// expandPolicyRoutingRuleSet expand the "set" section of the policy routing rule
func expandPolicyRoutingRuleSet(in map[string]interface{}) (*alkira.RoutePolicyRulesSet, error) {

	set := alkira.RoutePolicyRulesSet{}

	if v, ok := in["as_path_prepend"].(string); ok {
		set.AsPathPrepend = v
	}
	if v, ok := in["community"].(string); ok {
		set.Community = v
	}
	if v, ok := in["extended_community"].(string); ok {
		set.ExtendedCommunity = v
	}

	return &set, nil
}

// expandPolicyRoutingRuleInterCxpRoutesRedistribution expand the
//   "inter_cxp_routes_redistribution" section of the policy routing rule
func expandPolicyRoutingRuleInterCxpRoutesRedistribution(in map[string]interface{}) (*alkira.RoutePolicyRulesInterCxpRoutesRedistribution, error) {

	distrib := alkira.RoutePolicyRulesInterCxpRoutesRedistribution{}

	if v, ok := in["distribution_type"].(string); ok {
		distrib.DistributionType = v
	}
	if v, ok := in["redistribute_as_secondary"].(bool); ok {
		distrib.RedistributeAsSecondary = v
	}
	if v, ok := in["restricted_cxps"].([]interface{}); ok {
		distrib.RestrictedCxps = convertTypeListToStringList(v)

		if len(distrib.RestrictedCxps) == 0 {
			distrib.RestrictedCxps = nil
		}
	}

	return &distrib, nil
}

// expandPolicyRoutingRule expanding the "rule" sections of the routing policy
func expandPolicyRoutingRule(in *schema.Set) ([]alkira.RoutePolicyRules, error) {

	if in == nil || in.Len() == 0 {
		return nil, fmt.Errorf("[ERROR] Invalid route policy rule")
	}

	rules := make([]alkira.RoutePolicyRules, in.Len())

	for i, ruleInput := range in.List() {

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
