// Copyright (C) 2021 Alkira Inc. All Rights Reserved.

package alkira

import (
	"encoding/json"
	"fmt"
)

type NatPolicy struct {
	Name           string      `json:"name"`
	Description    string      `json:"description"`
	Type           string      `json:"type"`
	Segment        string      `json:"segment"`
	IncludedGroups []int       `json:"includedGroups"`
	ExcludedGroups []int       `json:"excludedGroups"`
	Id             json.Number `json:"id,omitempty"`
	NatRuleIds     []int       `json:"natRuleIds"`
}

// CreateNatPolicy create a nat policy
func (ac *AlkiraClient) CreateNatPolicy(p *NatPolicy) (string, error) {
	uri := fmt.Sprintf("%s/tenantnetworks/%s/nat-policies", ac.URI, ac.TenantNetworkId)

	// Construct the request
	body, err := json.Marshal(p)

	if err != nil {
		return "", err
	}

	data, err := ac.create(uri, body, false)

	if err != nil {
		return "", fmt.Errorf("CreateNatPolicy: request failed: %v", err)
	}

	var result NatPolicy
	json.Unmarshal([]byte(data), &result)

	if err != nil {
		return "", fmt.Errorf("CreateNatPolicy: failed to unmarshal: %v", err)
	}

	return string(result.Id), nil
}

// DeleteNatPolicy delete a nat policy by Id
func (ac *AlkiraClient) DeleteNatPolicy(id string) error {
	uri := fmt.Sprintf("%s/tenantnetworks/%s/nat-policies/%s", ac.URI, ac.TenantNetworkId, id)

	return ac.delete(uri, false)
}

// UpdateNatPolicy update a nat policy by Id
func (ac *AlkiraClient) UpdateNatPolicy(id string, p *NatPolicy) error {
	uri := fmt.Sprintf("%s/tenantnetworks/%s/nat-policies/%s", ac.URI, ac.TenantNetworkId, id)

	// Construct the request
	body, err := json.Marshal(p)

	if err != nil {
		return fmt.Errorf("UpdateNatPolicy: failed to marshal: %v", err)
	}

	return ac.update(uri, body, false)
}

// GetNatPolicy get a nat policy by Id
func (ac *AlkiraClient) GetNatPolicy(id string) (*NatPolicy, error) {
	uri := fmt.Sprintf("%s/tenantnetworks/%s/nat-policies/%s", ac.URI, ac.TenantNetworkId, id)

	data, err := ac.get(uri)

	if err != nil {
		return nil, err
	}

	var result NatPolicy
	err = json.Unmarshal([]byte(data), &result)

	if err != nil {
		return nil, fmt.Errorf("GetNatPolicy: failed to unmarshal: %v", err)
	}

	return &result, nil
}
