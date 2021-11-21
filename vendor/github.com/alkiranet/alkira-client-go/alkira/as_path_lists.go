// Copyright (C) 2021 Alkira Inc. All Rights Reserved.

package alkira

import (
	"encoding/json"
	"fmt"
)

type AsPathList struct {
	Name        string      `json:"name"`
	Description string      `json:"description"`
	Id          json.Number `json:"id,omitempty"`
	Values      []string    `json:"values"`
}

// GetAsPathLists get all AS path lists from the given tenant network
func (ac *AlkiraClient) GetAsPathLists() (string, error) {
	uri := fmt.Sprintf("%s/tenantnetworks/%s/as-path-lists", ac.URI, ac.TenantNetworkId)

	data, err := ac.get(uri)
	return string(data), err
}

// GetAsPathListById get single AS path list by Id
func (ac *AlkiraClient) GetAsPathListById(id string) (AsPathList, error) {
	uri := fmt.Sprintf("%s/tenantnetworks/%s/as-path-lists/%s", ac.URI, ac.TenantNetworkId, id)

	var list AsPathList

	data, err := ac.get(uri)

	if err != nil {
		return list, err
	}

	err = json.Unmarshal([]byte(data), &list)

	if err != nil {
		return list, fmt.Errorf("GetAsPathListById: failed to unmarshal: %v", err)
	}

	return list, nil
}

// GetAsPathListByName get the AS path list by name
func (ac *AlkiraClient) GetAsPathListByName(name string) (AsPathList, error) {
	var list AsPathList

	if len(name) == 0 {
		return list, fmt.Errorf("GetAsPathListByName: Invalid AS path list name")
	}

	lists, err := ac.GetAsPathLists()

	if err != nil {
		return list, err
	}

	var result []AsPathList
	json.Unmarshal([]byte(lists), &result)

	for _, l := range result {
		if l.Name == name {
			return l, nil
		}
	}

	return list, fmt.Errorf("GetAsPathListByName: failed to find the AS path list by %s", name)
}

// CreateAsPathList create an AS path list
func (ac *AlkiraClient) CreateAsPathList(list *AsPathList) (string, error) {
	uri := fmt.Sprintf("%s/tenantnetworks/%s/as-path-lists", ac.URI, ac.TenantNetworkId)

	// Construct the request
	body, err := json.Marshal(list)

	data, err := ac.create(uri, body)

	if err != nil {
		return "", err
	}

	var result AsPathList
	err = json.Unmarshal([]byte(data), &result)

	if err != nil {
		return "", fmt.Errorf("CreateAsPathList: failed to unmarshal: %v", err)
	}

	return string(result.Id), nil
}

// DeleteAsPathList delete an AS path list by Id
func (ac *AlkiraClient) DeleteAsPathList(id string) error {
	uri := fmt.Sprintf("%s/tenantnetworks/%s/as-path-lists/%s", ac.URI, ac.TenantNetworkId, id)

	return ac.delete(uri)
}

// UpdateAsPathList update an AS path list by Id
func (ac *AlkiraClient) UpdateAsPathList(id string, list *AsPathList) error {
	uri := fmt.Sprintf("%s/tenantnetworks/%s/as-path-lists/%s", ac.URI, ac.TenantNetworkId, id)

	body, err := json.Marshal(list)

	if err != nil {
		return fmt.Errorf("UpdateAsPathList: failed to marshal: %v", err)
	}

	return ac.update(uri, body)
}
