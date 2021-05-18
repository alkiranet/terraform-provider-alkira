// Copyright (C) 2020-2021 Alkira Inc. All Rights Reserved.

package alkira

import (
	"encoding/json"
	"fmt"
	"strconv"
)

type Group struct {
	Id          int    `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

// GetGroups get all groups from the given tenant network
func (ac *AlkiraClient) GetGroups() (string, error) {
	uri := fmt.Sprintf("%s/tenantnetworks/%s/groups", ac.URI, ac.TenantNetworkId)

	data, err := ac.get(uri)
	return string(data), err
}

// CreateGroup create a new Group
func (ac *AlkiraClient) CreateGroup(name string, description string) (string, error) {
	uri := fmt.Sprintf("%s/tenantnetworks/%s/groups", ac.URI, ac.TenantNetworkId)

	body, err := json.Marshal(map[string]string{
		"name":        name,
		"description": description,
	})

	if err != nil {
		return "", fmt.Errorf("CreateGroup: failed to marshal: %v", err)
	}

	data, err := ac.create(uri, body)

	if err != nil {
		return "", err
	}

	var result Group
	err = json.Unmarshal([]byte(data), &result)

	if err != nil {
		return "", fmt.Errorf("CreateGroup: failed to unmarshal: %v", err)
	}

	return strconv.Itoa(result.Id), nil
}

// GetGroup get a group by its Id
func (ac *AlkiraClient) GetGroupById(id string) (Group, error) {
	uri := fmt.Sprintf("%s/tenantnetworks/%s/groups/%s", ac.URI, ac.TenantNetworkId, id)
	var group Group

	data, err := ac.get(uri)

	if err != nil {
		return group, err
	}

	err = json.Unmarshal([]byte(data), &group)

	if err != nil {
		return group, fmt.Errorf("GetGroup: failed to unmarshal: %v", err)
	}

	return group, nil
}

// GetGroupByName get the group by its name
func (ac *AlkiraClient) GetGroupByName(name string) (Group, error) {
	var group Group

	if len(name) == 0 {
		return group, fmt.Errorf("Invalid group name input")
	}

	groups, err := ac.GetGroups()

	if err != nil {
		return group, err
	}

	var result []Group
	json.Unmarshal([]byte(groups), &result)

	for _, g := range result {
		if g.Name == name {
			return g, nil
		}
	}

	return group, fmt.Errorf("failed to find the group by %s", name)
}

// DeleteGroup delete a group
func (ac *AlkiraClient) DeleteGroup(id string) error {
	uri := fmt.Sprintf("%s/tenantnetworks/%s/groups/%s", ac.URI, ac.TenantNetworkId, id)
	return ac.delete(uri)
}

// UpdateGroup update a group by its id
func (ac *AlkiraClient) UpdateGroup(id string, name string, description string) error {
	uri := fmt.Sprintf("%s/tenantnetworks/%s/groups/%s", ac.URI, ac.TenantNetworkId, id)

	body, err := json.Marshal(map[string]string{
		"name":        name,
		"description": description,
	})

	if err != nil {
		return fmt.Errorf("UpdateGroup: failed to marshal: %v", err)
	}

	return ac.update(uri, body)
}
