// Copyright (C) 2020-2021 Alkira Inc. All Rights Reserved.

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

// CreatePolicyRuleList create a policy rule list
func (ac *AlkiraClient) CreatePolicyRuleList(p *PolicyRuleList) (string, error) {
	uri := fmt.Sprintf("%s/tenantnetworks/%s/policy/rulelists", ac.URI, ac.TenantNetworkId)

	// Construct the request
	body, err := json.Marshal(p)

	if err != nil {
		return "", err
	}

	data, err := ac.create(uri, body, true)

	if err != nil {
		return "", fmt.Errorf("CreatePolicyRuleList: request failed: %v", err)
	}

	var result PolicyRuleList
	err = json.Unmarshal([]byte(data), &result)

	if err != nil {
		return "", fmt.Errorf("CreatePolicyRuleList: request failed: %v", err)
	}

	return string(result.Id), nil
}

// DeletePolicyRuleList delete a policy rule list
func (ac *AlkiraClient) DeletePolicyRuleList(id string) error {
	uri := fmt.Sprintf("%s/tenantnetworks/%s/policy/rulelists/%s", ac.URI, ac.TenantNetworkId, id)

	return ac.delete(uri, true)
}

// UpdatePolicyRuleList update a policy rule list
func (ac *AlkiraClient) UpdatePolicyRuleList(id string, p *PolicyRuleList) error {
	uri := fmt.Sprintf("%s/tenantnetworks/%s/policy/rulelists/%s", ac.URI, ac.TenantNetworkId, id)

	// Construct the request
	body, err := json.Marshal(p)

	if err != nil {
		return fmt.Errorf("UpdatePolicyRuleList: failed to marshal: %v", err)
	}

	return ac.update(uri, body, true)
}

// GetPolicyRuleList get a policy rule list
func (ac *AlkiraClient) GetPolicyRuleList(id string) (*PolicyRuleList, error) {
	uri := fmt.Sprintf("%s/tenantnetworks/%s/policy/rulelists/%s", ac.URI, ac.TenantNetworkId, id)

	data, err := ac.get(uri)

	if err != nil {
		return nil, err
	}

	var result PolicyRuleList
	err = json.Unmarshal([]byte(data), &result)

	if err != nil {
		return nil, fmt.Errorf("GetPolicyRuleList: failed to unmarshal: %v", err)
	}

	return &result, nil
}
