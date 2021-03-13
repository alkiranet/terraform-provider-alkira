// Copyright (C) 2020-2021 Alkira Inc. All Rights Reserved.

package alkira

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

type TenantNetworkId struct {
	Id int `json:"id"`
}

type TenantNetworkState struct {
	State string `json:"state"`
}

// GetTenantNetworks get the tenant networks of the current tenant
func (ac *AlkiraClient) GetTenantNetworks() (string, error) {
	uri := fmt.Sprintf("%s/tenantnetworks", ac.URI)

	request, err := http.NewRequest("GET", uri, nil)
	request.Header.Set("Content-Type", "application/json")
	response, err := ac.Client.Do(request)

	if err != nil {
		return "", fmt.Errorf("GetTenantNetworks: request failed: %v", err)
	}

	defer response.Body.Close()
	data, _ := ioutil.ReadAll(response.Body)

	if response.StatusCode != 200 {
		return "", fmt.Errorf("(%d) %s", response.StatusCode, string(data))
	}

	return string(data), nil
}

// GetTenantNetworkId get the tenant network Id of the current tenant
func (ac *AlkiraClient) GetTenantNetworkId() (int, error) {

	uri := fmt.Sprintf("%s/tenantnetworks", ac.URI)
	id := 0

	request, err := http.NewRequest("GET", uri, nil)
	request.Header.Set("Content-Type", "application/json")
	response, err := ac.Client.Do(request)

	if err != nil {
		return id, fmt.Errorf("GetTenantNetworkId: request failed: %v", err)
	}

	defer response.Body.Close()
	data, _ := ioutil.ReadAll(response.Body)

	if response.StatusCode != 200 {
		return id, fmt.Errorf("(%d) %s", response.StatusCode, string(data))
	}

	var result []TenantNetworkId
	err = json.Unmarshal([]byte(data), &result)

	if err != nil {
		return id, fmt.Errorf("GetTenantNetworkId: parse failed: %v", err)
	}

	id = result[0].Id
	return result[0].Id, nil
}

// GetTenantNetworkState get the tenant network state
func (ac *AlkiraClient) GetTenantNetworkState() (string, error) {
	uri := fmt.Sprintf("%s/tenantnetworks/%s", ac.URI, ac.TenantNetworkId)
	state := ""

	request, err := http.NewRequest("GET", uri, nil)
	request.Header.Set("Content-Type", "application/json")
	response, err := ac.Client.Do(request)

	if err != nil {
		return state, fmt.Errorf("GetTenantNetworkState: request failed: %v", err)
	}

	defer response.Body.Close()
	data, _ := ioutil.ReadAll(response.Body)

	if response.StatusCode != 200 {
		return state, fmt.Errorf("(%d) %s", response.StatusCode, string(data))
	}

	var result TenantNetworkState
	err = json.Unmarshal([]byte(data), &result)

	if err != nil {
		return state, fmt.Errorf("GetTenantNetworkState: parse failed: %v", err)
	}

	state = result.State
	return state, nil
}

// ProvisionTenantNetwork provisioning the current tenant network
func (ac *AlkiraClient) ProvisionTenantNetwork() (string, error) {
	uri := fmt.Sprintf("%s/tenantnetworks/%s/provision", ac.URI, ac.TenantNetworkId)
	state := ""

	request, err := http.NewRequest("POST", uri, nil)
	request.Header.Set("Content-Type", "application/json")
	response, err := ac.Client.Do(request)

	if err != nil {
		return state, fmt.Errorf("ProvisionTenantNetwork: request failed: %v", err)
	}

	defer response.Body.Close()
	data, _ := ioutil.ReadAll(response.Body)

	if response.StatusCode != 200 && response.StatusCode != 202 {
		return state, fmt.Errorf("(%d) %s", response.StatusCode, string(data))
	}

	var result TenantNetworkState
	err = json.Unmarshal([]byte(data), &result)

	if err != nil {
		return state, fmt.Errorf("ProvisionTenantNetwork: parse failed: %v", err)
	}

	state = result.State
	return result.State, nil
}
