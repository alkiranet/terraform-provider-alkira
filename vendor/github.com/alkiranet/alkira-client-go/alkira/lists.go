// Copyright (C) 2021 Alkira Inc. All Rights Reserved.

//
// Implementation of the generic lists with common structure. For
// special lists, it's implementated separately as its own.
//
package alkira

import (
	"encoding/json"
	"fmt"
)

type ListType string

const (
	ListTypeAsPath            ListType = "as-path-lists"
	ListTypeCommunity                  = "community-lists"
	ListTypeExtendedCommunity          = "extended-community-lists"
)

type List struct {
	Description string      `json:"description"`
	Id          json.Number `json:"id,omitempty"`
	Name        string      `json:"name"`
	Values      []string    `json:"values"`
}

// GetLists get all lists from the given tenant network
func (ac *AlkiraClient) GetLists(t ListType) (string, error) {
	uri := fmt.Sprintf("%s/tenantnetworks/%s/%s", ac.URI, ac.TenantNetworkId, t)

	data, err := ac.get(uri)
	return string(data), err
}

// GetListById get single list by Id
func (ac *AlkiraClient) GetListById(id string, t ListType) (List, error) {
	uri := fmt.Sprintf("%s/tenantnetworks/%s/%s/%s", ac.URI, ac.TenantNetworkId, t, id)

	var list List

	data, err := ac.get(uri)

	if err != nil {
		return list, err
	}

	err = json.Unmarshal([]byte(data), &list)

	if err != nil {
		return list, fmt.Errorf("GetListById: failed to unmarshal: %v", err)
	}

	return list, nil
}

// GetListByName get the list by its name
func (ac *AlkiraClient) GetListByName(name string, t ListType) (List, error) {
	var list List

	if len(name) == 0 {
		return list, fmt.Errorf("GetListByName: Invalid list name")
	}

	lists, err := ac.GetLists(t)

	if err != nil {
		return list, err
	}

	var result []List
	json.Unmarshal([]byte(lists), &result)

	for _, l := range result {
		if l.Name == name {
			return l, nil
		}
	}

	return list, fmt.Errorf("GetListByName: failed to find the list(%s) by %s", t, name)
}

// CreateList create a list
func (ac *AlkiraClient) CreateList(list *List, t ListType) (string, error) {
	uri := fmt.Sprintf("%s/tenantnetworks/%s/%s", ac.URI, ac.TenantNetworkId, t)

	// Construct the request
	body, err := json.Marshal(list)

	data, err := ac.create(uri, body, true)

	if err != nil {
		return "", err
	}

	var result List
	err = json.Unmarshal([]byte(data), &result)

	if err != nil {
		return "", fmt.Errorf("CreateList: failed to unmarshal: %v", err)
	}

	return string(result.Id), nil
}

// DeleteList delete a list by Id
func (ac *AlkiraClient) DeleteList(id string, t ListType) error {
	uri := fmt.Sprintf("%s/tenantnetworks/%s/%s/%s", ac.URI, ac.TenantNetworkId, t, id)

	return ac.delete(uri, true)
}

// UpdateList update a list by Id
func (ac *AlkiraClient) UpdateList(id string, list *List, t ListType) error {
	uri := fmt.Sprintf("%s/tenantnetworks/%s/%s/%s", ac.URI, ac.TenantNetworkId, t, id)

	body, err := json.Marshal(list)

	if err != nil {
		return fmt.Errorf("UpdateList: failed to marshal: %v", err)
	}

	return ac.update(uri, body, true)
}
