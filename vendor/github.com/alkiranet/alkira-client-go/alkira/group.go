// Copyright (C) 2020-2021 Alkira Inc. All Rights Reserved.
package alkira

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

type Group struct {
	Id          int    `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

// GetGroups get all groups from the given tenant network
func (ac *AlkiraClient) GetGroups() (string, error) {
	uri := fmt.Sprintf("%s/tenantnetworks/%s/groups", ac.URI, ac.TenantNetworkId)

	request, err := http.NewRequest("GET", uri, nil)
	request.Header.Set("Content-Type", "application/json")
	response, err := ac.Client.Do(request)

	if err != nil {
		return "", fmt.Errorf("GetGroups: request failed, %v", err)
	}

	defer response.Body.Close()
	data, _ := ioutil.ReadAll(response.Body)

	if response.StatusCode != 200 {
		return "", fmt.Errorf("(%d) %s", response.StatusCode, string(data))
	}

	return string(data), nil
}

// CreateGroup create a new Group
func (ac *AlkiraClient) CreateGroup(name string, description string) (int, error) {
	uri := fmt.Sprintf("%s/tenantnetworks/%s/groups", ac.URI, ac.TenantNetworkId)

	body, err := json.Marshal(map[string]string{
		"name":        name,
		"description": description,
	})

	request, err := http.NewRequest("POST", uri, bytes.NewBuffer(body))
	request.Header.Set("Content-Type", "application/json")
	response, err := ac.Client.Do(request)

	if err != nil {
		return 0, fmt.Errorf("CreateGroup: request failed, %v", err)
	}

	defer response.Body.Close()
	data, _ := ioutil.ReadAll(response.Body)

	var result Group
	json.Unmarshal([]byte(data), &result)

	if response.StatusCode != 201 {
		return 0, fmt.Errorf("(%d) %s", response.StatusCode, string(data))
	}

	return result.Id, nil
}

// GetGroup get a group by its Id
func (ac *AlkiraClient) GetGroupById(id int) (Group, error) {
	uri := fmt.Sprintf("%s/tenantnetworks/%s/groups/%d", ac.URI, ac.TenantNetworkId, id)
	var group Group

	request, err := http.NewRequest("GET", uri, nil)
	request.Header.Set("Content-Type", "application/json")
	response, err := ac.Client.Do(request)

	if err != nil {
		return group, fmt.Errorf("GetGroup: request failed, %v", err)
	}

	defer response.Body.Close()
	data, _ := ioutil.ReadAll(response.Body)

	if response.StatusCode != 200 {
		return group, fmt.Errorf("(%d) %s", response.StatusCode, string(data))
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
func (ac *AlkiraClient) DeleteGroup(id int) error {
	uri := fmt.Sprintf("%s/tenantnetworks/%s/groups/%d", ac.URI, ac.TenantNetworkId, id)

	request, err := http.NewRequest("DELETE", uri, nil)
	request.Header.Set("Content-Type", "application/json")
	response, err := ac.Client.Do(request)

	if err != nil {
		return fmt.Errorf("DeleteGroup: request failed, %v", err)
	}

	defer response.Body.Close()
	data, _ := ioutil.ReadAll(response.Body)

	if response.StatusCode != 200 {
		return fmt.Errorf("(%d) %s", response.StatusCode, string(data))
	}

	return nil
}

// UpdateGroup update a group by its id
func (ac *AlkiraClient) UpdateGroup(id int, name string, description string) error {
	uri := fmt.Sprintf("%s/tenantnetworks/%s/groups/%d", ac.URI, ac.TenantNetworkId, id)

	body, err := json.Marshal(map[string]string{
		"name":        name,
		"description": description,
	})

	request, err := http.NewRequest("PUT", uri, bytes.NewBuffer(body))
	request.Header.Set("Content-Type", "application/json")
	response, err := ac.Client.Do(request)

	if err != nil {
		return fmt.Errorf("UpdateGroup: request failed, %v", err)
	}

	defer response.Body.Close()
	data, _ := ioutil.ReadAll(response.Body)

	var result Group
	json.Unmarshal([]byte(data), &result)

	if response.StatusCode != 200 {
		return fmt.Errorf("(%d) %s", response.StatusCode, string(data))
	}

	return nil
}
