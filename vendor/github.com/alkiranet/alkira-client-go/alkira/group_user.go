// Copyright (C) 2022 Alkira Inc. All Rights Reserved.

package alkira

import (
	"encoding/json"
	"fmt"
)

type UserGroup struct {
	Id          string `json:"id,omitempty"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

// GetUserGroups get all user groups from the given tenant network
func (ac *AlkiraClient) GetUserGroups() (string, error) {
	uri := fmt.Sprintf("%s/user-groups", ac.URI)

	data, err := ac.get(uri)
	return string(data), err
}

// CreateUserGroup create a new user group
func (ac *AlkiraClient) CreateUserGroup(name string, description string) (string, error) {
	uri := fmt.Sprintf("%s/user-groups", ac.URI)

	body, err := json.Marshal(map[string]string{
		"name":        name,
		"description": description,
	})

	if err != nil {
		return "", fmt.Errorf("CreateUserGroup: failed to marshal: %v", err)
	}

	data, err := ac.create(uri, body, true)

	if err != nil {
		return "", err
	}

	var result UserGroup
	err = json.Unmarshal([]byte(data), &result)

	if err != nil {
		return "", fmt.Errorf("CreateUserGroup: failed to unmarshal: %v", err)
	}

	return result.Id, nil
}

// GetUserGroup get an user group by its ID
func (ac *AlkiraClient) GetUserGroupById(id string) (UserGroup, error) {
	uri := fmt.Sprintf("%s/user-groups/%s", ac.URI, id)
	var group UserGroup

	data, err := ac.get(uri)

	if err != nil {
		return group, err
	}

	err = json.Unmarshal([]byte(data), &group)

	if err != nil {
		return group, fmt.Errorf("GetUserGroup: failed to unmarshal: %v", err)
	}

	return group, nil
}

// GetUserGroupByName get the user group by its name
func (ac *AlkiraClient) GetUserGroupByName(name string) (UserGroup, error) {
	var group UserGroup

	if len(name) == 0 {
		return group, fmt.Errorf("Invalid group name input")
	}

	groups, err := ac.GetUserGroups()

	if err != nil {
		return group, err
	}

	var result []UserGroup
	json.Unmarshal([]byte(groups), &result)

	for _, g := range result {
		if g.Name == name {
			return g, nil
		}
	}

	return group, fmt.Errorf("failed to find the group by %s", name)
}

// DeleteUserGroup delete an user group
func (ac *AlkiraClient) DeleteUserGroup(id string) error {
	uri := fmt.Sprintf("%s/user-groups/%s", ac.URI, id)
	return ac.delete(uri, true)
}

// UpdateUserGroup update an user group by its ID
func (ac *AlkiraClient) UpdateUserGroup(id string, name string, description string) error {
	uri := fmt.Sprintf("%s/user-groups/%s", ac.URI, id)

	body, err := json.Marshal(map[string]string{
		"name":        name,
		"description": description,
	})

	if err != nil {
		return fmt.Errorf("UpdateUserGroup: failed to marshal: %v", err)
	}

	return ac.update(uri, body, true)
}
