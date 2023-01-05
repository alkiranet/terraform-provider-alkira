// Copyright (C) 2020-2023 Alkira Inc. All Rights Reserved.

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
	ApplicationFamilyList []int    `json:"applicationFamilyList"`
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
}

// CreatePolicyRule create a policy rule
func (ac *AlkiraClient) CreatePolicyRule(p *PolicyRule) (string, error) {
	uri := fmt.Sprintf("%s/tenantnetworks/%s/policy/rules", ac.URI, ac.TenantNetworkId)

	// Construct the request
	body, err := json.Marshal(p)

	if err != nil {
		return "", err
	}

	data, err := ac.create(uri, body, true)

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
	return ac.delete(uri, true)
}

// UpdatePolicyRule update a policy rule
func (ac *AlkiraClient) UpdatePolicyRule(id string, p *PolicyRule) error {
	uri := fmt.Sprintf("%s/tenantnetworks/%s/policy/rules/%s", ac.URI, ac.TenantNetworkId, id)

	// Construct the request
	body, err := json.Marshal(p)

	if err != nil {
		return fmt.Errorf("UpdatePolicyRule: failed to marshal: %v", err)
	}

	return ac.update(uri, body, true)
}

// GetPolicyRule get all policy rules from the given tenant network
func (ac *AlkiraClient) GetPolicyRules() (string, error) {
	uri := fmt.Sprintf("%s/tenantnetworks/%s/policy/rules", ac.URI, ac.TenantNetworkId)

	data, err := ac.get(uri)
	return string(data), err
}

// GetPolicyRuleByName get the policy rule by its name
func (ac *AlkiraClient) GetPolicyRuleByName(name string) (*PolicyRule, error) {

	if len(name) == 0 {
		return nil, fmt.Errorf("GetPolicyRuleByName: Invalid list name")
	}

	rules, err := ac.GetPolicyRules()

	if err != nil {
		return nil, err
	}

	var result []PolicyRule
	json.Unmarshal([]byte(rules), &result)

	for _, l := range result {
		if l.Name == name {
			return &l, nil
		}
	}

	return nil, fmt.Errorf("GetPolicyRuleByName: failed to find the rule by name %s", name)
}

// GetPolicyRule get a policy rule by ID
func (ac *AlkiraClient) GetPolicyRuleById(id string) (*PolicyRule, error) {
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
