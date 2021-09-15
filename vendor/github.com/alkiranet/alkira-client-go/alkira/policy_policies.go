// Copyright (C) 2020-2021 Alkira Inc. All Rights Reserved.

package alkira

import (
	"encoding/json"
	"fmt"
)

type Policy struct {
	Description string      `json:"description"`
	Enabled     string      `json:"enabled"`
	FromGroups  []int       `json:"fromGroups"`
	Id          json.Number `json:"id,omitempty"`
	Name        string      `json:"name"`
	RuleListId  int         `json:"ruleListId"`
	SegmentIds  []int       `json:"segmentIds"`
	ToGroups    []int       `json:"toGroups"`
}

// CreatePolicy Create a policy
func (ac *AlkiraClient) CreatePolicy(p *Policy) (string, error) {
	uri := fmt.Sprintf("%s/tenantnetworks/%s/policy/policies", ac.URI, ac.TenantNetworkId)

	// Construct the request
	body, err := json.Marshal(p)

	if err != nil {
		return "", err
	}

	data, err := ac.create(uri, body)

	if err != nil {
		return "", fmt.Errorf("CreatePolicy: request failed: %v", err)
	}

	var result Policy
	json.Unmarshal([]byte(data), &result)

	if err != nil {
		return "", fmt.Errorf("CreatePolicy: request failed: %v", err)
	}

	return string(result.Id), nil
}

// DeletePolicy delete a policy by Id
func (ac *AlkiraClient) DeletePolicy(id string) error {
	uri := fmt.Sprintf("%s/tenantnetworks/%s/policy/policies/%s", ac.URI, ac.TenantNetworkId, id)

	return ac.delete(uri)
}

// UpdatePolicy update a policy by Id
func (ac *AlkiraClient) UpdatePolicy(id string, p *Policy) error {
	uri := fmt.Sprintf("%s/tenantnetworks/%s/policy/policies/%s", ac.URI, ac.TenantNetworkId, id)

	// Construct the request
	body, err := json.Marshal(p)

	if err != nil {
		return fmt.Errorf("UpdatePolicy: failed to marshal: %v", err)
	}

	return ac.update(uri, body)
}

// GetPolicy get a policy by Id
func (ac *AlkiraClient) GetPolicy(id string) (*Policy, error) {
	uri := fmt.Sprintf("%s/tenantnetworks/%s/policy/policies/%s", ac.URI, ac.TenantNetworkId, id)

	data, err := ac.get(uri)

	if err != nil {
		return nil, err
	}

	var result Policy
	err = json.Unmarshal([]byte(data), &result)

	if err != nil {
		return nil, fmt.Errorf("GetPolicy: failed to unmarshal: %v", err)
	}

	return &result, nil
}
