// Copyright (C) 2020-2022 Alkira Inc. All Rights Reserved.

package alkira

import (
	"encoding/json"
	"fmt"
	"strconv"
)

type ConnectorAzureErVnetSegment struct {
	SegmentName           string `json:"segmentName"`
	SegmentId             int    `json:"segmentId:omitempty"`
	CustomerAsn           int    `json:"customerAsn"`
	DisableInternetExit   bool   `json:"disableInternetExit"`
	AdvertiseOnPremRoutes bool   `json:"advertiseOnPremRoutes"`
}

type ConnectorAzureErVnetInstance struct {
	Name                  string   `json:"name"`
	Id                    int      `json:"id,omitempty"`
	ExpressRouteCircuitId string   `json:"expressRouteCircuitId"`
	RedundantRouter       bool     `json:"redundantRouter,omitempty"`
	LoopbackSubnet        string   `json:"loopbackSubnet,omitempty"`
	CredentialId          string   `json:"credentialId"`
	GatewayMacAddress     []string `json:"gatewayMacAddress,omitempty"`
	Vnis                  []int    `json:"vnis,omitempty"`
}

type ConnectorAzureErVnet struct {
	Name           string                         `json:"name"`
	Id             int                            `json:"id,omitempty"`
	CredentialId   string                         `json:"credentialId"`
	Size           string                         `json:"size"`
	Enabled        bool                           `json:"enabled"`
	VhubPrefix     string                         `json:"vhubPrefix"`
	TunnelProtocol string                         `json:"tunnelProtocol"`
	Cxp            string                         `json:"cxp"`
	Instances      []ConnectorAzureErVnetInstance `json:"instances,omitempty"`
	SegmentOptions []ConnectorAzureErVnetSegment  `json:"segmentOptions,omitempty"`
	BillingTags    []int                          `json:"billingTags"`
}

// getAzureErVnetConnectors get all Azure Express Route Vnet connectors from the given tenant network
func (ac *AlkiraClient) getAzureErVnetConnectors() (string, error) {
	uri := fmt.Sprintf("%s/tenantnetworks/%s/azure-express-route-connectors", ac.URI, ac.TenantNetworkId)

	data, err := ac.get(uri)
	return string(data), err
}

// CreateConnectorAzureErVnet create a Azure Express Route Vnet connector
func (ac *AlkiraClient) CreateConnectorAzureErVnet(connector *ConnectorAzureErVnet) (string, error) {
	uri := fmt.Sprintf("%s/tenantnetworks/%s/azure-express-route-connectors", ac.URI, ac.TenantNetworkId)

	// Construct the request
	body, err := json.Marshal(connector)

	if err != nil {
		return "", fmt.Errorf("CreateConnectorAzureErVnet: failed to marshal: %v", err)
	}

	data, err := ac.create(uri, body, true)

	if err != nil {
		return "", err
	}

	var result ConnectorAzureErVnet
	json.Unmarshal([]byte(data), &result)

	if err != nil {
		return "", fmt.Errorf("CreateConnectorAzureErVnet: failed to unmarshal: %v", err)
	}

	return strconv.Itoa(result.Id), nil
}

// GetConnectorAzureErVnet get one Azure Express Route Vnet connector by Id
func (ac *AlkiraClient) GetConnectorAzureErVnet(id string) (*ConnectorAzureErVnet, error) {
	uri := fmt.Sprintf("%s/tenantnetworks/%s/azure-express-route-connectors/%s", ac.URI, ac.TenantNetworkId, id)

	data, err := ac.get(uri)

	if err != nil {
		return nil, err
	}

	var result ConnectorAzureErVnet
	err = json.Unmarshal([]byte(data), &result)

	if err != nil {
		return nil, fmt.Errorf("GetConnectorAzureErVnet: failed to unmarshal: %v", err)
	}

	return &result, nil
}

// DeleteConnectorAzureErVnet delete the given Azure Express Route Vnet connector by Id
func (ac *AlkiraClient) DeleteConnectorAzureErVnet(id string) error {
	uri := fmt.Sprintf("%s/tenantnetworks/%s/azure-express-route-connectors/%s", ac.URI, ac.TenantNetworkId, id)

	return ac.delete(uri, true)
}

// UpdateConnectorAzureErVnet update an Azure Express Route Vnet connector by Id
func (ac *AlkiraClient) UpdateConnectorAzureErVnet(id string, connector *ConnectorAzureErVnet) error {
	uri := fmt.Sprintf("%s/v1/tenantnetworks/%s/azure-express-route-connectors/%s", ac.URI, ac.TenantNetworkId, id)

	// Construct the request
	body, err := json.Marshal(connector)

	if err != nil {
		return fmt.Errorf("UpdateConnectorAzureErVnet: failed to marshal: %v", err)
	}

	return ac.update(uri, body, true)
}

// GetConnectorAzureErVnetByName get an Azure Express Route Vnet connector by name
func (ac *AlkiraClient) GetConnectorAzureErVnetByName(name string) (ConnectorAzureErVnet, error) {
	var azureErVnetConnector ConnectorAzureErVnet

	if len(name) == 0 {
		return azureErVnetConnector, fmt.Errorf("GetConnectorAzureErVnetByName: Invalid Connector name")
	}

	azureErVnetConnectors, err := ac.getAzureErVnetConnectors()

	if err != nil {
		return azureErVnetConnector, err
	}

	var result []ConnectorAzureErVnet
	json.Unmarshal([]byte(azureErVnetConnectors), &result)

	for _, l := range result {
		if l.Name == name {
			return l, nil
		}
	}

	return azureErVnetConnector, fmt.Errorf("GetConnectorAzureErVnetByName: failed to find the connector by %s", name)
}
