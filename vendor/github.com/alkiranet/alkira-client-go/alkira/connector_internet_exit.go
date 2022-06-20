// Copyright (C) 2020-2021 Alkira Inc. All Rights Reserved.

package alkira

import (
	"encoding/json"
	"fmt"
	"strconv"
)

// ENUM for keys: "DEFAULT", "SRC_IP"
type AlgorithmAttributes struct {
	Keys string `json:"keys"`
}

type TrafficDistribution struct {
	Algorithm           string              `json:"algorithm"`
	AlgorithmAttributes AlgorithmAttributes `json:"algorithmAttributes"`
}

type ConnectorInternet struct {
	BillingTags         []int                `json:"billingTags"`
	CXP                 string               `json:"cxp"`
	Description         string               `json:"description"`
	Group               string               `json:"group,omitempty"`
	Enabled             bool                 `json:"enabled,omitempty"`
	Id                  int                  `json:"id"`
	Name                string               `json:"name"`
	NumOfPublicIPs      int                  `json:"numOfPublicIPs,omitempty"`
	Segments            []string             `json:"segments"`
	Size                string               `json:"size"`
	TrafficDistribution *TrafficDistribution `json:"trafficDistribution,omitempty"`
}

// CreateConnectorInternetExit create an internet connector
func (ac *AlkiraClient) CreateConnectorInternetExit(connector *ConnectorInternet) (string, error) {
	uri := fmt.Sprintf("%s/v1/tenantnetworks/%s/internetconnectors", ac.URI, ac.TenantNetworkId)

	// Construct the request
	body, err := json.Marshal(connector)
	if err != nil {
		return "", fmt.Errorf("CreateConnectorInternetExit: failed to marshal: %v", err)
	}

	data, err := ac.create(uri, body, true)

	if err != nil {
		return "", err
	}

	var result ConnectorInternet
	err = json.Unmarshal([]byte(data), &result)

	if err != nil {
		return "", fmt.Errorf("CreateConnectorInternetExit: failed to unmarshal: %v", err)
	}

	return strconv.Itoa(result.Id), nil
}

// GetConnectorInternetExitById get an internet exit connector by id
func (ac *AlkiraClient) GetConnectorInternetExitById(id string) (*ConnectorInternet, error) {
	uri := fmt.Sprintf("%s/v1/tenantnetworks/%s/internetconnectors/%s", ac.URI, ac.TenantNetworkId, id)
	var result ConnectorInternet

	data, err := ac.get(uri)

	if err != nil {
		return nil, err
	}

	err = json.Unmarshal([]byte(data), &result)

	if err != nil {
		return nil, fmt.Errorf("GetConnectorInternetExit: failed to unmarshal: %v", err)
	}

	return &result, nil
}

// DeleteConnectorInternetExit delete an internet exit connector
func (ac *AlkiraClient) DeleteConnectorInternetExit(id string) error {
	uri := fmt.Sprintf("%s/v1/tenantnetworks/%s/internetconnectors/%s", ac.URI, ac.TenantNetworkId, id)

	return ac.delete(uri, true)
}

// UpdateConnectorInternetExit update an internet exit connector by its Id
func (ac *AlkiraClient) UpdateConnectorInternetExit(id string, connector *ConnectorInternet) error {
	uri := fmt.Sprintf("%s/v1/tenantnetworks/%s/internetconnectors/%s", ac.URI, ac.TenantNetworkId, id)

	body, err := json.Marshal(connector)

	if err != nil {
		return fmt.Errorf("UpdateConnectorInternetExit: failed to marshal request: %v", err)
	}

	return ac.update(uri, body, true)
}
