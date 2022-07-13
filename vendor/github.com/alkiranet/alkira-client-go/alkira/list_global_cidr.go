// Copyright (C) 2021 Alkira Inc. All Rights Reserved.

package alkira

import (
	"encoding/json"
	"fmt"
)

type GlobalCidrList struct {
	Description string      `json:"description"`
	CXP         string      `json:"cxp"`
	Id          json.Number `json:"id,omitempty"`
	Name        string      `json:"name"`
	Values      []string    `json:"values"`
}

// GetGlobalCidrLists Get all global CIDR list from the given tenant network
func (ac *AlkiraClient) GetGlobalCidrLists() (string, error) {
	uri := fmt.Sprintf("%s/tenantnetworks/%s/global-cidr-lists", ac.URI, ac.TenantNetworkId)

	data, err := ac.get(uri)
	return string(data), err
}

// GetGlobalCidrListById get single Global CIDR list by Id
func (ac *AlkiraClient) GetGlobalCidrListById(id string) (GlobalCidrList, error) {
	uri := fmt.Sprintf("%s/tenantnetworks/%s/global-cidr-lists/%s", ac.URI, ac.TenantNetworkId, id)

	var list GlobalCidrList

	data, err := ac.get(uri)

	if err != nil {
		return list, err
	}

	err = json.Unmarshal([]byte(data), &list)

	if err != nil {
		return list, fmt.Errorf("GetGlobalCidrListById: failed to unmarshal: %v", err)
	}

	return list, nil
}

// GetGlobalCidrListByName get the global CIDR by its name
func (ac *AlkiraClient) GetGlobalCidrListByName(name string) (GlobalCidrList, error) {
	var list GlobalCidrList

	if len(name) == 0 {
		return list, fmt.Errorf("Invalid Global CIDR list name")
	}

	lists, err := ac.GetGlobalCidrLists()

	if err != nil {
		return list, err
	}

	var result []GlobalCidrList
	json.Unmarshal([]byte(lists), &result)

	for _, p := range result {
		if p.Name == name {
			return p, nil
		}
	}

	return list, fmt.Errorf("failed to find the global CIDR list with name %s", name)
}

// CreateGlobalCidrList create a global CIDR list
func (ac *AlkiraClient) CreateGlobalCidrList(p *GlobalCidrList) (string, error) {
	uri := fmt.Sprintf("%s/tenantnetworks/%s/global-cidr-lists", ac.URI, ac.TenantNetworkId)

	// Construct the request
	body, err := json.Marshal(p)

	data, err := ac.create(uri, body, true)

	if err != nil {
		return "", err
	}

	var result GlobalCidrList
	err = json.Unmarshal([]byte(data), &result)

	if err != nil {
		return "", fmt.Errorf("CreateGlobalCidrList: failed to unmarshal: %v", err)
	}

	return string(result.Id), nil
}

// DeleteGlobalCidrList delete a global CIDR list
func (ac *AlkiraClient) DeleteGlobalCidrList(id string) error {
	uri := fmt.Sprintf("%s/tenantnetworks/%s/global-cidr-lists/%s", ac.URI, ac.TenantNetworkId, id)

	return ac.delete(uri, true)
}

// UpdateGlobalCidrList update a global CIDR list by Id
func (ac *AlkiraClient) UpdateGlobalCidrList(id string, l *GlobalCidrList) error {

	uri := fmt.Sprintf("%s/tenantnetworks/%s/global-cidr-lists/%s", ac.URI, ac.TenantNetworkId, id)

	body, err := json.Marshal(l)

	if err != nil {
		return fmt.Errorf("UpdateGlobalCidrList: failed to marshal: %v", err)
	}

	return ac.update(uri, body, true)
}
