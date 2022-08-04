package alkira

import (
	"fmt"

	"github.com/alkiranet/alkira-client-go/alkira"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func expandPolicyRoutingRuleMatch(in *schema.Set) (*alkira.RoutePolicyRulesMatch, error) {

	if in == nil || in.Len() == 0 {
		return nil, nil
	}

	if in.Len() > 1 {
		return nil, fmt.Errorf("Only one match block should be defined in routing policy.")
	}

	match := alkira.RoutePolicyRulesMatch{}

	for _, matchInput := range in.List() {

		input := matchInput.(map[string]interface{})

		if v, ok := input["all"].(bool); ok {
			match.All = v
		}
		if v, ok := input["as_path_list_ids"].([]interface{}); ok {
			match.AsPathListIds = convertTypeListToIntList(v)
		}
		if v, ok := input["community_list_ids"].([]interface{}); ok {
			match.CommunityListIds = convertTypeListToIntList(v)
		}
		if v, ok := input["extended_community_list_ids"].([]interface{}); ok {
			match.ExtendedCommunityListIds = convertTypeListToIntList(v)
		}
		if v, ok := input["prefix_list_ids"].([]interface{}); ok {
			match.PrefixListIds = convertTypeListToIntList(v)
		}
		if v, ok := input["cxps"].([]interface{}); ok {
			match.Cxps = convertTypeListToStringList(v)
		}
		if v, ok := input["group_ids"].([]interface{}); ok {
			match.ConnectorGroupIds = convertTypeListToIntList(v)
		}
	}

	return &match, nil
}

// expandPolicyRoutingRuleSet expand the "set" section of the policy routing rule
func expandPolicyRoutingRuleSet(in *schema.Set) (*alkira.RoutePolicyRulesSet, error) {

	if in == nil || in.Len() == 0 {
		return nil, nil
	}

	if in.Len() > 1 {
		return nil, fmt.Errorf("Only one match block should be defined in routing policy.")
	}

	set := alkira.RoutePolicyRulesSet{}

	for _, setInput := range in.List() {

		input := setInput.(map[string]interface{})

		if v, ok := input["as_path_prepend"].(string); ok {
			set.AsPathPrepend = v
		}
		if v, ok := input["community"].(string); ok {
			set.Community = v
		}
		if v, ok := input["extended_community"].(string); ok {
			set.ExtendedCommunity = v
		}
	}

	return &set, nil
}

// expandPolicyRoutingRuleInterCxpRoutesRedistribution expand the
//   "inter_cxp_routes_redistribution" section of the policy routing rule
func expandPolicyRoutingRuleInterCxpRoutesRedistribution(in *schema.Set) (*alkira.RoutePolicyRulesInterCxpRoutesRedistribution, error) {

	if in == nil || in.Len() == 0 {
		return nil, nil
	}

	distrib := alkira.RoutePolicyRulesInterCxpRoutesRedistribution{}

	for _, disInput := range in.List() {

		input := disInput.(map[string]interface{})

		if v, ok := input["distribution_type"].(string); ok {
			distrib.DistributionType = v
		}
		if v, ok := input["redistribute_as_secondary"].(bool); ok {
			distrib.RedistributeAsSecondary = v
		}
		if v, ok := input["restricted_cxps"].([]interface{}); ok {
			distrib.RestrictedCxps = convertTypeListToStringList(v)

			if len(distrib.RestrictedCxps) == 0 {
				distrib.RestrictedCxps = nil
			}
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

	for i, rule := range in.List() {

		r := alkira.RoutePolicyRules{}
		ruleCfg := rule.(map[string]interface{})

		if v, ok := ruleCfg["action"].(string); ok {
			r.Action = v
		}
		if v, ok := ruleCfg["name"].(string); ok {
			r.Name = v
		}
		if v, ok := ruleCfg["match"].(*schema.Set); ok {
			match, err := expandPolicyRoutingRuleMatch(v)

			if err != nil {
				return nil, err
			}

			r.Match = *match
		}
		if v, ok := ruleCfg["set"].(*schema.Set); ok {
			set, err := expandPolicyRoutingRuleSet(v)
			if err != nil {
				return nil, err
			}

			r.Set = set
		}
		if v, ok := ruleCfg["inter_cxp_routes_redistribution"].(*schema.Set); ok {
			distribution, err := expandPolicyRoutingRuleInterCxpRoutesRedistribution(v)
			if err != nil {
				return nil, err
			}

			r.InterCxpRoutesRedistribution = distribution
		}
		rules[i] = r
	}

	return rules, nil
}
