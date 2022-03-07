// Copyright (C) 2020-2022 Alkira Inc. All Rights Reserved.

package alkira

import (
	"encoding/json"
	"fmt"
	"strconv"
)

type ConnectorGroup struct {
	Id          int    `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

// GetConnectorGroups get all connector groups from the given tenant network
func (ac *AlkiraClient) GetConnectorGroups() (string, error) {
	uri := fmt.Sprintf("%s/tenantnetworks/%s/groups", ac.URI, ac.TenantNetworkId)

	data, err := ac.get(uri)
	return string(data), err
}

// CreateConnectorGroup create a new connector group
func (ac *AlkiraClient) CreateConnectorGroup(name string, description string) (string, error) {
	uri := fmt.Sprintf("%s/tenantnetworks/%s/groups", ac.URI, ac.TenantNetworkId)

	body, err := json.Marshal(map[string]string{
		"name":        name,
		"description": description,
	})

	if err != nil {
		return "", fmt.Errorf("CreateConnectorGroup: failed to marshal: %v", err)
	}

	data, err := ac.create(uri, body)

	if err != nil {
		return "", err
	}

	var result ConnectorGroup
	err = json.Unmarshal([]byte(data), &result)

	if err != nil {
		return "", fmt.Errorf("CreateConnectorGroup: failed to unmarshal: %v", err)
	}

	return strconv.Itoa(result.Id), nil
}

// GetConnectorGroup get a connector group by its ID
func (ac *AlkiraClient) GetConnectorGroupById(id string) (ConnectorGroup, error) {
	uri := fmt.Sprintf("%s/tenantnetworks/%s/groups/%s", ac.URI, ac.TenantNetworkId, id)
	var group ConnectorGroup

	data, err := ac.get(uri)

	if err != nil {
		return group, err
	}

	err = json.Unmarshal([]byte(data), &group)

	if err != nil {
		return group, fmt.Errorf("GetConnectorGroup: failed to unmarshal: %v", err)
	}

	return group, nil
}

// GetConnectorGroupByName get the connector group by its name
func (ac *AlkiraClient) GetConnectorGroupByName(name string) (ConnectorGroup, error) {
	var group ConnectorGroup

	if len(name) == 0 {
		return group, fmt.Errorf("Invalid group name input")
	}

	groups, err := ac.GetConnectorGroups()

	if err != nil {
		return group, err
	}

	var result []ConnectorGroup
	json.Unmarshal([]byte(groups), &result)

	for _, g := range result {
		if g.Name == name {
			return g, nil
		}
	}

	return group, fmt.Errorf("failed to find the group by %s", name)
}

// DeleteConnectorGroup delete a connector group
func (ac *AlkiraClient) DeleteConnectorGroup(id string) error {
	uri := fmt.Sprintf("%s/tenantnetworks/%s/groups/%s", ac.URI, ac.TenantNetworkId, id)
	return ac.delete(uri)
}

// UpdateConnectorGroup update a connector group by its ID
func (ac *AlkiraClient) UpdateConnectorGroup(id string, name string, description string) error {
	uri := fmt.Sprintf("%s/tenantnetworks/%s/groups/%s", ac.URI, ac.TenantNetworkId, id)

	body, err := json.Marshal(map[string]string{
		"name":        name,
		"description": description,
	})

	if err != nil {
		return fmt.Errorf("UpdateConnectorGroup: failed to marshal: %v", err)
	}

	return ac.update(uri, body)
}
