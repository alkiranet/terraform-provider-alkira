// Copyright (C) 2020-2025 Alkira Inc. All Rights Reserved.

package alkira

import (
	"encoding/json"
	"fmt"
)

type PolicyRuleList struct {
	Description string               `json:"description"`
	Id          json.Number          `json:"id,omitempty"`
	Name        string               `json:"name"`
	Rules       []PolicyRuleListRule `json:"rules"`
}

type PolicyRuleListRule struct {
	Priority int `json:"priority"`
	RuleId   int `json:"ruleId"`
}

// NewPolicyRuleList new policy rule list
func NewPolicyRuleList(ac *AlkiraClient) *AlkiraAPI[PolicyRuleList] {
	uri := fmt.Sprintf("%s/tenantnetworks/%s/policy/rulelists", ac.URI, ac.TenantNetworkId)
	api := &AlkiraAPI[PolicyRuleList]{ac, uri, true}
	return api
}
