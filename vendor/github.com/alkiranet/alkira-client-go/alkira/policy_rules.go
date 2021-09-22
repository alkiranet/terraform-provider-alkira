// Copyright (C) 2020-2021 Alkira Inc. All Rights Reserved.

package alkira

import (
	"encoding/json"
	"fmt"
)

type PolicyRule struct {
	Description    string                   `json:"description"`
	Id             json.Number              `json:"id,omitempty"`
	MatchCondition PolicyRuleMatchCondition `json:"matchCondition"`
	Name           string                   `json:"name"`
	RuleAction     PolicyRuleAction         `json:"ruleAction"`
}

type PolicyRuleMatchCondition struct {
	ApplicationFamilyList []string `json:"applicationFamilyList"`
	ApplicationList       []string `json:"applicationList"`
	Dscp                  string   `json:"dscp"`
	DstIp                 string   `json:"dstIp"`
	DstPortList           []string `json:"dstPortList,omitempty"`
	DstPrefixListId       int      `json:"dstPrefixListId,omitempty"`
	Protocol              string   `json:"protocol"`
	SrcIp                 string   `json:"srcIp"`
	SrcPortList           []string `json:"srcPortList,omitempty"`
	SrcPrefixListId       int      `json:"srcPrefixListId,omitempty"`
}

type PolicyRuleAction struct {
	Action          string   `json:"action"`
	ServiceTypeList []string `json:"serviceTypeList"`
	ServiceList     []int    `json:"serviceList"`
}

// CreatePolicyRule create a policy rule
func (ac *AlkiraClient) CreatePolicyRule(p *PolicyRule) (string, error) {
	uri := fmt.Sprintf("%s/tenantnetworks/%s/policy/rules", ac.URI, ac.TenantNetworkId)

	// Construct the request
	body, err := json.Marshal(p)

	if err != nil {
		return "", err
	}

	data, err := ac.create(uri, body)

	if err != nil {
		return "", fmt.Errorf("CreatePolicyRule: request failed: %v", err)
	}

	var result PolicyRule
	err = json.Unmarshal([]byte(data), &result)

	if err != nil {
		return "", fmt.Errorf("CreatePolicyRule: request failed: %v", err)
	}

	return string(result.Id), nil
}

// DeletePolicyRule delete a policy rule
func (ac *AlkiraClient) DeletePolicyRule(id string) error {
	uri := fmt.Sprintf("%s/tenantnetworks/%s/policy/rules/%s", ac.URI, ac.TenantNetworkId, id)
	return ac.delete(uri)
}

// UpdatePolicyRule update a policy rule list
func (ac *AlkiraClient) UpdatePolicyRule(id string, p *PolicyRule) error {
	uri := fmt.Sprintf("%s/tenantnetworks/%s/policy/rules/%s", ac.URI, ac.TenantNetworkId, id)

	// Construct the request
	body, err := json.Marshal(p)

	if err != nil {
		return fmt.Errorf("UpdatePolicyRule: failed to marshal: %v", err)
	}

	return ac.update(uri, body)
}

// GetPolicyRule get a policy rule list
func (ac *AlkiraClient) GetPolicyRule(id string) (*PolicyRule, error) {
	uri := fmt.Sprintf("%s/tenantnetworks/%s/policy/rules/%s", ac.URI, ac.TenantNetworkId, id)

	data, err := ac.get(uri)

	if err != nil {
		return nil, err
	}

	var result PolicyRule
	err = json.Unmarshal([]byte(data), &result)

	if err != nil {
		return nil, fmt.Errorf("GetPolicyRule: failed to unmarshal: %v", err)
	}

	return &result, nil
}
