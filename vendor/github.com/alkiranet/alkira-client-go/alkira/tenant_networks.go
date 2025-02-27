// Copyright (C) 2020-2025 Alkira Inc. All Rights Reserved.

package alkira

import (
	"encoding/json"
	"fmt"
	"strconv"
)

type TenantNetworkId struct {
	Id int `json:"id"`
}

type TenantNetworkState struct {
	State string `json:"state"`
}

type TenantNetworkConnectorState struct {
	State    string `json:"state"`
	DocState string `json:"docState"`
}

type TenantNetworkServiceState struct {
	State    string `json:"state"`
	DocState string `json:"docState"`
}

type TenantNetworkProvisionRequest struct {
	Id    string `json:"id"`
	State string `json:"state"`
}

// GetTenantNetworks get the tenant networks of the current tenant
func (ac *AlkiraClient) GetTenantNetworks() (string, error) {
	uri := fmt.Sprintf("%s/tenantnetworks", ac.URI)

	data, _, err := ac.get(uri)

	if err != nil {
		return "", err
	}

	return string(data), nil
}

// GetTenantNetworkId get the tenant network Id of the current tenant
func (ac *AlkiraClient) GetTenantNetworkId() (string, error) {

	uri := fmt.Sprintf("%s/tenantnetworks", ac.URI)

	data, _, err := ac.get(uri)

	if err != nil {
		return "", err
	}

	var result []TenantNetworkId
	err = json.Unmarshal([]byte(data), &result)

	if err != nil {
		return "", fmt.Errorf("GetTenantNetworkId: failed to unmarshal: %v", err)
	}

	return strconv.Itoa(result[0].Id), nil
}

// GetTenantNetworkState get the tenant network state
func (ac *AlkiraClient) GetTenantNetworkState() (string, error) {
	uri := fmt.Sprintf("%s/tenantnetworks/%s", ac.URI, ac.TenantNetworkId)

	data, _, err := ac.get(uri)

	if err != nil {
		return "", err
	}

	var result TenantNetworkState
	err = json.Unmarshal([]byte(data), &result)

	if err != nil {
		return "", fmt.Errorf("GetTenantNetworkState: failed to unmarshal: %v", err)
	}

	return result.State, nil
}

// GetTenantNetworkConnectorState get the tenant network connector state by its Id
func (ac *AlkiraClient) GetTenantNetworkConnectorState(id string) (string, error) {
	uri := fmt.Sprintf("%s/tenantnetworks/%s/connectors/%s", ac.URI, ac.TenantNetworkId, id)

	data, _, err := ac.get(uri)

	if err != nil {
		return "", err
	}

	var result TenantNetworkConnectorState
	err = json.Unmarshal([]byte(data), &result)

	if err != nil {
		return "", fmt.Errorf("GetTenantNetworkConnectorState: failed to unmarshal: %v", err)
	}

	return result.State, nil
}

// GetTenantNetworkServiceState get the tenant network service state by its Id
func (ac *AlkiraClient) GetTenantNetworkServiceState(id string) (string, error) {
	uri := fmt.Sprintf("%s/tenantnetworks/%s/services/%s", ac.URI, ac.TenantNetworkId, id)

	data, _, err := ac.get(uri)

	if err != nil {
		return "", err
	}

	var result TenantNetworkServiceState
	err = json.Unmarshal([]byte(data), &result)

	if err != nil {
		return "", fmt.Errorf("GetTenantNetworkConnectorState: failed to unmarshal: %v", err)
	}

	return result.State, nil
}

// GetTenantNetworkProvisionRequest get the tenant network provision request
func (ac *AlkiraClient) GetTenantNetworkProvisionRequest(id string) (*TenantNetworkProvisionRequest, error) {
	uri := fmt.Sprintf("%s/tenantnetworks/%s/provision-requests/%s", ac.URI, ac.TenantNetworkId, id)

	data, _, err := ac.get(uri)

	if err != nil {
		return nil, err
	}

	var result TenantNetworkProvisionRequest
	err = json.Unmarshal([]byte(data), &result)

	if err != nil {
		return nil, fmt.Errorf("GetTenantNetworkProvisionRequest: failed to unmarshal: %v", err)
	}

	return &result, nil
}

// ProvisionTenantNetwork provisioning the current tenant network by its Id
func (ac *AlkiraClient) ProvisionTenantNetwork() (string, error) {
	uri := fmt.Sprintf("%s/tenantnetworks/%s/provision", ac.URI, ac.TenantNetworkId)

	data, _, err, _ := ac.create(uri, nil, false)

	if err != nil {
		return "", err
	}

	var result TenantNetworkState
	err = json.Unmarshal([]byte(data), &result)

	if err != nil {
		return "", fmt.Errorf("ProvisionTenantNetwork: failed to unmarshal: %v", err)
	}

	return result.State, nil
}
