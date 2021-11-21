// Copyright (C) 2021 Alkira Inc. All Rights Reserved.

package alkira

import (
	"encoding/json"
	"fmt"
)

type CommunityList struct {
	Name        string      `json:"name"`
	Description string      `json:"description"`
	Id          json.Number `json:"id,omitempty"`
	Values      []string    `json:"values"`
}

// GetCommunityLists get all community lists from the given tenant network
func (ac *AlkiraClient) GetCommunityLists() (string, error) {
	uri := fmt.Sprintf("%s/tenantnetworks/%s/community-lists", ac.URI, ac.TenantNetworkId)

	data, err := ac.get(uri)
	return string(data), err
}

// GetCommunityListById get single community list by Id
func (ac *AlkiraClient) GetCommunityListById(id string) (CommunityList, error) {
	uri := fmt.Sprintf("%s/tenantnetworks/%s/community-lists/%s", ac.URI, ac.TenantNetworkId, id)

	var list CommunityList

	data, err := ac.get(uri)

	if err != nil {
		return list, err
	}

	err = json.Unmarshal([]byte(data), &list)

	if err != nil {
		return list, fmt.Errorf("GetCommunityListById: failed to unmarshal: %v", err)
	}

	return list, nil
}

// GetCommunityListByName get the community list by its name
func (ac *AlkiraClient) GetCommunityListByName(name string) (CommunityList, error) {
	var list CommunityList

	if len(name) == 0 {
		return list, fmt.Errorf("GetCommunityListByName: Invalid community list name")
	}

	lists, err := ac.GetCommunityLists()

	if err != nil {
		return list, err
	}

	var result []CommunityList
	json.Unmarshal([]byte(lists), &result)

	for _, l := range result {
		if l.Name == name {
			return l, nil
		}
	}

	return list, fmt.Errorf("GetCommunityListByName: failed to find the community list by %s", name)
}

// CreateCommunityList create a community list
func (ac *AlkiraClient) CreateCommunityList(list *CommunityList) (string, error) {
	uri := fmt.Sprintf("%s/tenantnetworks/%s/community-lists", ac.URI, ac.TenantNetworkId)

	// Construct the request
	body, err := json.Marshal(list)

	data, err := ac.create(uri, body)

	if err != nil {
		return "", err
	}

	var result CommunityList
	err = json.Unmarshal([]byte(data), &result)

	if err != nil {
		return "", fmt.Errorf("CreateCommunityList: failed to unmarshal: %v", err)
	}

	return string(result.Id), nil
}

// DeleteCommunityList delete a community list by Id
func (ac *AlkiraClient) DeleteCommunityList(id string) error {
	uri := fmt.Sprintf("%s/tenantnetworks/%s/community-lists/%s", ac.URI, ac.TenantNetworkId, id)

	return ac.delete(uri)
}

// UpdateCommunityList update a community list by id
func (ac *AlkiraClient) UpdateCommunityList(id string, list *CommunityList) error {
	uri := fmt.Sprintf("%s/tenantnetworks/%s/community-lists/%s", ac.URI, ac.TenantNetworkId, id)

	body, err := json.Marshal(list)

	if err != nil {
		return fmt.Errorf("UpdateCommunityList: failed to marshal: %v", err)
	}

	return ac.update(uri, body)
}
