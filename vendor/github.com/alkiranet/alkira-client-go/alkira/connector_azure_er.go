// Copyright (C) 2020-2022 Alkira Inc. All Rights Reserved.

package alkira

import (
	"encoding/json"
	"fmt"
	"strconv"
)

type ConnectorAzureErSegment struct {
	SegmentName           string `json:"segmentName"`
	SegmentId             int    `json:"segmentId:omitempty"`
	CustomerAsn           int    `json:"customerAsn"`
	DisableInternetExit   bool   `json:"disableInternetExit"`
	AdvertiseOnPremRoutes bool   `json:"advertiseOnPremRoutes"`
}

type ConnectorAzureErInstance struct {
	Name                  string   `json:"name"`
	Id                    int      `json:"id,omitempty"`
	ExpressRouteCircuitId string   `json:"expressRouteCircuitId"`
	RedundantRouter       bool     `json:"redundantRouter,omitempty"`
	LoopbackSubnet        string   `json:"loopbackSubnet,omitempty"`
	CredentialId          string   `json:"credentialId"`
	GatewayMacAddress     []string `json:"gatewayMacAddress,omitempty"`
	Vnis                  []int    `json:"vnis,omitempty"`
}

type ConnectorAzureEr struct {
	Name           string                     `json:"name"`
	Id             int                        `json:"id,omitempty"`
	CredentialId   string                     `json:"credentialId"`
	Size           string                     `json:"size"`
	Enabled        bool                       `json:"enabled"`
	VhubPrefix     string                     `json:"vhubPrefix"`
	TunnelProtocol string                     `json:"tunnelProtocol"`
	Cxp            string                     `json:"cxp"`
	Group          string                     `json:"group,omitempty"`
	Instances      []ConnectorAzureErInstance `json:"instances,omitempty"`
	SegmentOptions []ConnectorAzureErSegment  `json:"segmentOptions,omitempty"`
	BillingTags    []int                      `json:"billingTags"`
}

// getAzureErConnectors get all Azure Express Route connectors from the given tenant network
func (ac *AlkiraClient) getAzureErConnectors() (string, error) {
	uri := fmt.Sprintf("%s/tenantnetworks/%s/azure-express-route-connectors", ac.URI, ac.TenantNetworkId)

	data, err := ac.get(uri)
	return string(data), err
}

// CreateConnectorAzureEr create an Azure Express Route connector
func (ac *AlkiraClient) CreateConnectorAzureEr(connector *ConnectorAzureEr) (string, error) {
	uri := fmt.Sprintf("%s/tenantnetworks/%s/azure-express-route-connectors", ac.URI, ac.TenantNetworkId)

	// Construct the request
	body, err := json.Marshal(connector)

	if err != nil {
		return "", fmt.Errorf("CreateConnectorAzureEr: failed to marshal: %v", err)
	}

	data, err := ac.create(uri, body, true)

	if err != nil {
		return "", err
	}

	var result ConnectorAzureEr
	json.Unmarshal([]byte(data), &result)

	if err != nil {
		return "", fmt.Errorf("CreateConnectorAzureEr: failed to unmarshal: %v", err)
	}

	return strconv.Itoa(result.Id), nil
}

// GetConnectorAzureEr get one Azure Express Route connector by Id
func (ac *AlkiraClient) GetConnectorAzureEr(id string) (*ConnectorAzureEr, error) {
	uri := fmt.Sprintf("%s/tenantnetworks/%s/azure-express-route-connectors/%s", ac.URI, ac.TenantNetworkId, id)

	data, err := ac.get(uri)

	if err != nil {
		return nil, err
	}

	var result ConnectorAzureEr
	err = json.Unmarshal([]byte(data), &result)

	if err != nil {
		return nil, fmt.Errorf("GetConnectorAzureEr: failed to unmarshal: %v", err)
	}

	return &result, nil
}

// DeleteConnectorAzureEr delete the given Azure Express Route connector by Id
func (ac *AlkiraClient) DeleteConnectorAzureEr(id string) error {
	uri := fmt.Sprintf("%s/tenantnetworks/%s/azure-express-route-connectors/%s", ac.URI, ac.TenantNetworkId, id)

	return ac.delete(uri, true)
}

// UpdateConnectorAzureEr update an Azure Express Route connector by Id
func (ac *AlkiraClient) UpdateConnectorAzureEr(id string, connector *ConnectorAzureEr) error {
	uri := fmt.Sprintf("%s/v1/tenantnetworks/%s/azure-express-route-connectors/%s", ac.URI, ac.TenantNetworkId, id)

	// Construct the request
	body, err := json.Marshal(connector)

	if err != nil {
		return fmt.Errorf("UpdateConnectorAzureEr: failed to marshal: %v", err)
	}

	return ac.update(uri, body, true)
}

// GetConnectorAzureErByName get an Azure Express Route connector by name
func (ac *AlkiraClient) GetConnectorAzureErByName(name string) (ConnectorAzureEr, error) {
	var azureErConnector ConnectorAzureEr

	if len(name) == 0 {
		return azureErConnector, fmt.Errorf("GetConnectorAzureErByName: Invalid Connector name")
	}

	azureErConnectors, err := ac.getAzureErConnectors()

	if err != nil {
		return azureErConnector, err
	}

	var result []ConnectorAzureEr
	json.Unmarshal([]byte(azureErConnectors), &result)

	for _, l := range result {
		if l.Name == name {
			return l, nil
		}
	}

	return azureErConnector, fmt.Errorf("GetConnectorAzureErByName: failed to find the connector by %s", name)
}
