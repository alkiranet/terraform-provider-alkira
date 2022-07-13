// Copyright (C) 2020-2021 Alkira Inc. All Rights Reserved.

package alkira

import (
	"encoding/json"
	"fmt"
)

type PolicyPrefixList struct {
	Description string      `json:"description"`
	Id          json.Number `json:"id,omitempty"`
	Name        string      `json:"name"`
	Prefixes    []string    `json:"prefixes"`
}

// GetPolicyPrefixLists Get all prefixes from the given tenant network
func (ac *AlkiraClient) GetPolicyPrefixLists() (string, error) {
	uri := fmt.Sprintf("%s/tenantnetworks/%s/policy/prefixlists", ac.URI, ac.TenantNetworkId)

	data, err := ac.get(uri)
	return string(data), err
}

// GetPolicyPrefixListById get single prefix list by Id
func (ac *AlkiraClient) GetPolicyPrefixListById(id string) (PolicyPrefixList, error) {
	uri := fmt.Sprintf("%s/tenantnetworks/%s/policy/prefixlists/%s", ac.URI, ac.TenantNetworkId, id)

	var prefixList PolicyPrefixList

	data, err := ac.get(uri)

	if err != nil {
		return prefixList, err
	}

	err = json.Unmarshal([]byte(data), &prefixList)

	if err != nil {
		return prefixList, fmt.Errorf("GetPolicyPrefixListById: failed to unmarshal: %v", err)
	}

	return prefixList, nil
}

// GetPolicyPrefixListByName get the prefix list by its name
func (ac *AlkiraClient) GetPolicyPrefixListByName(name string) (PolicyPrefixList, error) {
	var prefixList PolicyPrefixList

	if len(name) == 0 {
		return prefixList, fmt.Errorf("Invalid prefix list name")
	}

	prefixLists, err := ac.GetPolicyPrefixLists()

	if err != nil {
		return prefixList, err
	}

	var result []PolicyPrefixList
	json.Unmarshal([]byte(prefixLists), &result)

	for _, p := range result {
		if p.Name == name {
			return p, nil
		}
	}

	return prefixList, fmt.Errorf("failed to find the prefix list by %s", name)
}

// CreatePolicyPrefixList create a policy prefix
func (ac *AlkiraClient) CreatePolicyPrefixList(p *PolicyPrefixList) (string, error) {
	uri := fmt.Sprintf("%s/tenantnetworks/%s/policy/prefixlists", ac.URI, ac.TenantNetworkId)

	// Construct the request
	body, err := json.Marshal(p)

	data, err := ac.create(uri, body, true)

	if err != nil {
		return "", err
	}

	var result PolicyPrefixList
	err = json.Unmarshal([]byte(data), &result)

	if err != nil {
		return "", fmt.Errorf("CreatePolicyPrefixList: failed to unmarshal: %v", err)
	}

	return string(result.Id), nil
}

// DeletePolicyPrefixList delete a policy prefix list
func (ac *AlkiraClient) DeletePolicyPrefixList(id string) error {
	uri := fmt.Sprintf("%s/tenantnetworks/%s/policy/prefixlists/%s", ac.URI, ac.TenantNetworkId, id)

	return ac.delete(uri, true)
}

// UpdatePolicyPrefixList update a PolicyPrefixList by id
func (ac *AlkiraClient) UpdatePolicyPrefixList(id string, p *PolicyPrefixList) error {

	uri := fmt.Sprintf("%s/tenantnetworks/%s/prefixlists/%s", ac.URI, ac.TenantNetworkId, id)

	body, err := json.Marshal(p)

	if err != nil {
		return fmt.Errorf("UpdatePolicyPrefixList: failed to marshal: %v", err)
	}

	return ac.update(uri, body, true)
}
