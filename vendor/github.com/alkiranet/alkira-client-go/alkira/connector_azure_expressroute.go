// Copyright (C) 2020-2022 Alkira Inc. All Rights Reserved.

package alkira

import (
	"encoding/json"
	"fmt"
	"strconv"
)

type ConnectorAzureExpressRouteSegment struct {
	SegmentName           string `json:"segmentName"`
	SegmentId             int    `json:"segmentId:omitempty"`
	CustomerAsn           int    `json:"customerAsn"`
	DisableInternetExit   bool   `json:"disableInternetExit"`
	AdvertiseOnPremRoutes bool   `json:"advertiseOnPremRoutes"`
}

type ConnectorAzureExpressRouteInstance struct {
	Name                  string   `json:"name"`
	Id                    int      `json:"id,omitempty"`
	ExpressRouteCircuitId string   `json:"expressRouteCircuitId"`
	RedundantRouter       bool     `json:"redundantRouter,omitempty"`
	LoopbackSubnet        string   `json:"loopbackSubnet,omitempty"`
	CredentialId          string   `json:"credentialId"`
	GatewayMacAddress     []string `json:"gatewayMacAddress,omitempty"`
	Vnis                  []int    `json:"vnis,omitempty"`
}

type ConnectorAzureExpressRoute struct {
	Name           string                               `json:"name"`
	Id             int                                  `json:"id,omitempty"`
	Size           string                               `json:"size"`
	Enabled        bool                                 `json:"enabled"`
	VhubPrefix     string                               `json:"vhubPrefix"`
	TunnelProtocol string                               `json:"tunnelProtocol"`
	Cxp            string                               `json:"cxp"`
	Group          string                               `json:"group,omitempty"`
	Instances      []ConnectorAzureExpressRouteInstance `json:"instances,omitempty"`
	SegmentOptions []ConnectorAzureExpressRouteSegment  `json:"segmentOptions,omitempty"`
	BillingTags    []int                                `json:"billingTags"`
}

// getAzureExpressRouteConnectors get all Azure ExpressRoute connectors from the given tenant network
func (ac *AlkiraClient) getAzureExpressRouteConnectors() (string, error) {
	uri := fmt.Sprintf("%s/tenantnetworks/%s/azure-express-route-connectors", ac.URI, ac.TenantNetworkId)

	data, err := ac.get(uri)
	return string(data), err
}

// CreateConnectorAzureExpressRoute create an Azure ExpressRoute connector
func (ac *AlkiraClient) CreateConnectorAzureExpressRoute(connector *ConnectorAzureExpressRoute) (string, error) {
	uri := fmt.Sprintf("%s/tenantnetworks/%s/azure-express-route-connectors", ac.URI, ac.TenantNetworkId)

	// Construct the request
	body, err := json.Marshal(connector)

	if err != nil {
		return "", fmt.Errorf("CreateConnectorAzureExpressRoute: failed to marshal: %v", err)
	}

	data, err := ac.create(uri, body, true)

	if err != nil {
		return "", err
	}

	var result ConnectorAzureExpressRoute
	json.Unmarshal([]byte(data), &result)

	if err != nil {
		return "", fmt.Errorf("CreateConnectorAzureExpressRoute: failed to unmarshal: %v", err)
	}

	return strconv.Itoa(result.Id), nil
}

// GetConnectorAzureExpressRoute get one Azure ExpressRoute connector by Id
func (ac *AlkiraClient) GetConnectorAzureExpressRoute(id string) (*ConnectorAzureExpressRoute, error) {
	uri := fmt.Sprintf("%s/tenantnetworks/%s/azure-express-route-connectors/%s", ac.URI, ac.TenantNetworkId, id)

	data, err := ac.get(uri)

	if err != nil {
		return nil, err
	}

	var result ConnectorAzureExpressRoute
	err = json.Unmarshal([]byte(data), &result)

	if err != nil {
		return nil, fmt.Errorf("GetConnectorAzureExpressRoute: failed to unmarshal: %v", err)
	}

	return &result, nil
}

// DeleteConnectorAzureExpressRoute delete the given Azure ExpressRoute connector by Id
func (ac *AlkiraClient) DeleteConnectorAzureExpressRoute(id string) error {
	uri := fmt.Sprintf("%s/tenantnetworks/%s/azure-express-route-connectors/%s", ac.URI, ac.TenantNetworkId, id)

	return ac.delete(uri, true)
}

// UpdateConnectorAzureExpressRoute update an Azure ExpressRoute connector by Id
func (ac *AlkiraClient) UpdateConnectorAzureExpressRoute(id string, connector *ConnectorAzureExpressRoute) error {
	uri := fmt.Sprintf("%s/v1/tenantnetworks/%s/azure-express-route-connectors/%s", ac.URI, ac.TenantNetworkId, id)

	// Construct the request
	body, err := json.Marshal(connector)

	if err != nil {
		return fmt.Errorf("UpdateConnectorAzureExpressRoute: failed to marshal: %v", err)
	}

	return ac.update(uri, body, true)
}

// GetConnectorAzureExpressRouteByName get an Azure ExpressRoute connector by name
func (ac *AlkiraClient) GetConnectorAzureExpressRouteByName(name string) (ConnectorAzureExpressRoute, error) {
	var azureErConnector ConnectorAzureExpressRoute

	if len(name) == 0 {
		return azureErConnector, fmt.Errorf("GetConnectorAzureExpressRouteByName: Invalid Connector name")
	}

	azureErConnectors, err := ac.getAzureExpressRouteConnectors()

	if err != nil {
		return azureErConnector, err
	}

	var result []ConnectorAzureExpressRoute
	json.Unmarshal([]byte(azureErConnectors), &result)

	for _, l := range result {
		if l.Name == name {
			return l, nil
		}
	}

	return azureErConnector, fmt.Errorf("GetConnectorAzureExpressRouteByName: failed to find the connector by %s", name)
}