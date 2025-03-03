// Copyright (C) 2020-2025 Alkira Inc. All Rights Reserved.

package alkira

import (
	"encoding/json"
	"fmt"
)

type TrafficPolicyRule struct {
	Description    string                   `json:"description"`
	Id             json.Number              `json:"id,omitempty"`
	MatchCondition PolicyRuleMatchCondition `json:"matchCondition"`
	Name           string                   `json:"name"`
	RuleAction     PolicyRuleAction         `json:"ruleAction"`
}

type PolicyRuleMatchCondition struct {
	ApplicationList       []int    `json:"applicationList"`
	Dscp                  string   `json:"dscp"`
	DstIp                 string   `json:"dstIp,omitempty"`
	DstPortList           []string `json:"dstPortList,omitempty"`
	DstPrefixListId       int      `json:"dstPrefixListId,omitempty"`
	InternetApplicationId int      `json:"internetApplicationId,omitempty"`
	Protocol              string   `json:"protocol"`
	SrcIp                 string   `json:"srcIp,omitempty"`
	SrcPortList           []string `json:"srcPortList,omitempty"`
	SrcPrefixListId       int      `json:"srcPrefixListId,omitempty"`
}

type PolicyRuleAction struct {
	Action          string   `json:"action"`
	ServiceTypeList []string `json:"serviceTypeList"`
	ServiceList     []int    `json:"serviceList"`
	FlowCollectors  []int    `json:"flowCollectors,omitempty"`
}

// NewTrafficPolicyRule new traffic policy rule
func NewTrafficPolicyRule(ac *AlkiraClient) *AlkiraAPI[TrafficPolicyRule] {
	uri := fmt.Sprintf("%s/tenantnetworks/%s/policy/rules", ac.URI, ac.TenantNetworkId)
	api := &AlkiraAPI[TrafficPolicyRule]{ac, uri, true}
	return api
}
