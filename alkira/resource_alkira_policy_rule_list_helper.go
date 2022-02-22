package alkira

import (
	"log"

	"github.com/alkiranet/alkira-client-go/alkira"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func expandPolicyRuleListRules(in *schema.Set) []alkira.PolicyRuleListRule {
	if in == nil || in.Len() == 0 {
		log.Printf("[DEBUG] invalid policy rule")
		return nil
	}

	rules := make([]alkira.PolicyRuleListRule, in.Len())
	for i, rule := range in.List() {
		r := alkira.PolicyRuleListRule{}
		ruleCfg := rule.(map[string]interface{})
		if v, ok := ruleCfg["priority"].(int); ok {
			r.Priority = v
		}
		if v, ok := ruleCfg["rule_id"].(int); ok {
			r.RuleId = v
		}
		rules[i] = r
	}

	return rules
}
