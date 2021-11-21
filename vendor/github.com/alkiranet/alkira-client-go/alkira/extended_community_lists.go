// Copyright (C) 2021 Alkira Inc. All Rights Reserved.

package alkira

import (
	"encoding/json"
	"fmt"
)

type ExtendedCommunityList struct {
	Name        string      `json:"name"`
	Description string      `json:"description"`
	Id          json.Number `json:"id,omitempty"`
	Values      []string    `json:"values"`
}

// GetExtendedCommunityLists get all extended community lists from the given tenant network
func (ac *AlkiraClient) GetExtendedCommunityLists() (string, error) {
	uri := fmt.Sprintf("%s/tenantnetworks/%s/extended-community-lists", ac.URI, ac.TenantNetworkId)

	data, err := ac.get(uri)
	return string(data), err
}

// GetExtendedCommunityList get single extended community list by Id
func (ac *AlkiraClient) GetExtendedCommunityListById(id string) (ExtendedCommunityList, error) {
	uri := fmt.Sprintf("%s/tenantnetworks/%s/extended-community-lists/%s", ac.URI, ac.TenantNetworkId, id)

	var list ExtendedCommunityList

	data, err := ac.get(uri)

	if err != nil {
		return list, err
	}

	err = json.Unmarshal([]byte(data), &list)

	if err != nil {
		return list, fmt.Errorf("GetExtendedCommunityListById: failed to unmarshal: %v", err)
	}

	return list, nil
}

// GetExtendedCommunityList get the extended community list by its name
func (ac *AlkiraClient) GetExtendedCommunityListByName(name string) (ExtendedCommunityList, error) {
	var list ExtendedCommunityList

	if len(name) == 0 {
		return list, fmt.Errorf("GetExtendedCommunityListByName: Invalid extended community list name")
	}

	lists, err := ac.GetExtendedCommunityLists()

	if err != nil {
		return list, err
	}

	var result []ExtendedCommunityList
	json.Unmarshal([]byte(lists), &result)

	for _, l := range result {
		if l.Name == name {
			return l, nil
		}
	}

	return list, fmt.Errorf("GetExtendedCommunityListByName: failed to find the extended community list by %s", name)
}

// CreateExtendedCommunityList create a extended community list
func (ac *AlkiraClient) CreateExtendedCommunityList(list *ExtendedCommunityList) (string, error) {
	uri := fmt.Sprintf("%s/tenantnetworks/%s/extended-community-lists", ac.URI, ac.TenantNetworkId)

	// Construct the request
	body, err := json.Marshal(list)

	data, err := ac.create(uri, body)

	if err != nil {
		return "", err
	}

	var result ExtendedCommunityList
	err = json.Unmarshal([]byte(data), &result)

	if err != nil {
		return "", fmt.Errorf("CreateExtendedCommunityList: failed to unmarshal: %v", err)
	}

	return string(result.Id), nil
}

// DeleteExtendedCommunityList delete an extended community list by Id
func (ac *AlkiraClient) DeleteExtendedCommunityList(id string) error {
	uri := fmt.Sprintf("%s/tenantnetworks/%s/extended-community-lists/%s", ac.URI, ac.TenantNetworkId, id)

	return ac.delete(uri)
}

// UpdateExtendedCommunityList update an extended community list by id
func (ac *AlkiraClient) UpdateExtendedCommunityList(id string, list *ExtendedCommunityList) error {
	uri := fmt.Sprintf("%s/tenantnetworks/%s/extended-community-lists/%s", ac.URI, ac.TenantNetworkId, id)

	body, err := json.Marshal(list)

	if err != nil {
		return fmt.Errorf("UpdateExtendedCommunityList: failed to marshal: %v", err)
	}

	return ac.update(uri, body)
}
