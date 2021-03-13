// Copyright (C) 2020-2021 Alkira Inc. All Rights Reserved.

package alkira

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

type PolicyPrefixList struct {
	Description string   `json:"description"`
	Id          int      `json:"id"`
	Name        string   `json:"name"`
	Prefixes    []string `json:"prefixes"`
}

// GetPolicyPrefixLists Get all prefixes from the given tenant network
func (ac *AlkiraClient) GetPolicyPrefixLists() (string, error) {
	uri := fmt.Sprintf("%s/tenantnetworks/%s/policy/prefixlists", ac.URI, ac.TenantNetworkId)

	request, err := http.NewRequest("GET", uri, nil)
	request.Header.Set("Content-Type", "application/json")
	response, err := ac.Client.Do(request)

	if err != nil {
		return "", fmt.Errorf("GetPolicyPrefixLists: request failed: %v", err)
	}

	defer response.Body.Close()
	data, _ := ioutil.ReadAll(response.Body)

	if response.StatusCode != 200 {
		return "", fmt.Errorf("(%d) %s", response.StatusCode, string(data))
	}

	return string(data), nil
}

// GetPolicyPrefixListById get single prefix list by Id
func (ac *AlkiraClient) GetPolicyPrefixListById(id int) (PolicyPrefixList, error) {
	uri := fmt.Sprintf("%s/tenantnetworks/%s/policy/prefixlists/%d", ac.URI, ac.TenantNetworkId, id)

	var prefixList PolicyPrefixList

	request, err := http.NewRequest("GET", uri, nil)
	request.Header.Set("Content-Type", "application/json")
	response, err := ac.Client.Do(request)

	if err != nil {
		return prefixList, fmt.Errorf("GetPolicyPrefixListById: request failed: %v", err)
	}

	defer response.Body.Close()
	data, _ := ioutil.ReadAll(response.Body)

	if response.StatusCode != 200 {
		return prefixList, fmt.Errorf("(%d) %s", response.StatusCode, string(data))
	}

	err = json.Unmarshal([]byte(data), &prefixList)

	if err != nil {
		return prefixList, fmt.Errorf("GetPolicyPrefixListById: parse failed: %v", err)
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
func (ac *AlkiraClient) CreatePolicyPrefixList(p *PolicyPrefixList) (int, error) {
	uri := fmt.Sprintf("%s/tenantnetworks/%s/policy/prefixlists", ac.URI, ac.TenantNetworkId)
	id := 0

	// Construct the request
	body, err := json.Marshal(p)

	request, err := http.NewRequest("POST", uri, bytes.NewBuffer(body))
	request.Header.Set("Content-Type", "application/json")
	response, err := ac.Client.Do(request)

	if err != nil {
		return id, fmt.Errorf("CreatePolicyPrefixList: request failed: %v", err)
	}

	defer response.Body.Close()
	data, _ := ioutil.ReadAll(response.Body)

	var result PolicyPrefixList
	json.Unmarshal([]byte(data), &result)

	if response.StatusCode != 201 {
		return id, fmt.Errorf("(%d) %s", response.StatusCode, string(data))
	}

	id = result.Id
	return id, nil
}

// DeletePolicyPrefixList delete a policy prefix list
func (ac *AlkiraClient) DeletePolicyPrefixList(id int) error {
	uri := fmt.Sprintf("%s/tenantnetworks/%s/policy/prefixlists/%d", ac.URI, ac.TenantNetworkId, id)

	request, err := http.NewRequest("DELETE", uri, nil)

	request.Header.Set("Content-Type", "application/json")
	response, err := ac.Client.Do(request)

	if err != nil {
		return fmt.Errorf("DeletePolicyPrefixList: request failed: %v", err)
	}

	defer response.Body.Close()
	data, _ := ioutil.ReadAll(response.Body)

	if response.StatusCode != 200 && response.StatusCode != 202 {
		return fmt.Errorf("(%d) %s", response.StatusCode, string(data))
	}

	return nil
}
