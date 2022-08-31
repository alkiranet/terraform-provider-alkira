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
	Enabled             bool                 `json:"enabled"`
	Id                  int                  `json:"id"`
	Name                string               `json:"name"`
	NumOfPublicIPs      int                  `json:"numOfPublicIPs,omitempty"`
	Segments            []string             `json:"segments"`
	Size                string               `json:"size"`
	TrafficDistribution *TrafficDistribution `json:"trafficDistribution,omitempty"`
}

// getInternetConnectors get all Internet Connectors from the given tenant network
func (ac *AlkiraClient) getInternetConnectors() (string, error) {
	uri := fmt.Sprintf("%s/tenantnetworks/%s/internetconnectors", ac.URI, ac.TenantNetworkId)

	data, err := ac.get(uri)
	return string(data), err
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

// GetConnectorInternetByName get internet connector by name
func (ac *AlkiraClient) GetConnectorInternetByName(name string) (ConnectorInternet, error) {
	var internetConnector ConnectorInternet

	if len(name) == 0 {
		return internetConnector, fmt.Errorf("GetConnectorInternetByName: Invalid Connector name")
	}

	internetConnectors, err := ac.getInternetConnectors()

	if err != nil {
		return internetConnector, err
	}

	var result []ConnectorInternet
	json.Unmarshal([]byte(internetConnectors), &result)

	for _, l := range result {
		if l.Name == name {
			return l, nil
		}
	}

	return internetConnector, fmt.Errorf("GetConnectorInternetByName: failed to find the connector by %s", name)
}
