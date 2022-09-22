// Copyright (C) 2022 Alkira Inc. All Rights Reserved.

package alkira

import (
	"encoding/json"
	"fmt"
)

type ConnectorArubaEdge struct {
	ArubaEdgeVrfMapping []ArubaEdgeVRFMapping `json:"arubaEdgeVRFMapping,omitempty"`
	BillingTags         []int                 `json:"billingTags"`
	BoostMode           bool                  `json:"boostMode"`
	Cxp                 string                `json:"cxp"`
	GatewayBgpAsn       int                   `json:"gatewayBgpAsn"`
	Group               string                `json:"group,omitempty"`
	Id                  json.Number           `json:"id,omitempty"`              // response only
	ImplicitGroupId     int                   `json:"implicitGroupId,omitempty"` // response only
	Instances           []ArubaEdgeInstance   `json:"instances"`
	Name                string                `json:"name"`
	Segments            []string              `json:"segments"`
	Size                string                `json:"size"`
	TunnelProtocol      string                `json:"tunnelProtocol"`
	Version             string                `json:"version"`
}

type ArubaEdgeVRFMapping struct {
	AdvertiseOnPremRoutes       bool   `json:"advertiseOnPremRoutes"`
	AlkiraSegmentId             int    `json:"alkiraSegmentId"`
	ArubaEdgeConnectSegmentName string `json:"arubaEdgeConnectSegmentName"`
	DisableInternetExit         bool   `json:"disableInternetExit"`
	GatewayBgpAsn               int    `json:"gatewayBgpAsn"`
}

type ArubaEdgeInstance struct {
	Id           json.Number `json:"id,omitempty"`
	AccountName  string      `json:"accountName"`
	CredentialId string      `json:"credentialId"`
	HostName     string      `json:"hostName"`
	Name         string      `json:"name"`
	SiteTag      string      `json:"siteTag"`
}

type ArubaEdgeInstanceConfig struct {
	Data string //The response is string data. The entire body of the response should be interpreted together. There is no json structure.
}

// getArubaEdgeConnectors get all aruba edge connectors from the given tenant network
func (ac *AlkiraClient) getArubaEdgeConnectors() (string, error) {
	uri := fmt.Sprintf("%s/tenantnetworks/%s/aruba-edge-connectors", ac.URI, ac.TenantNetworkId)

	data, err := ac.get(uri)
	return string(data), err
}

func (ac *AlkiraClient) CreateConnectorArubaEdge(c *ConnectorArubaEdge) (string, error) {
	uri := fmt.Sprintf("%s/tenantnetworks/%s/aruba-edge-connectors", ac.URI, ac.TenantNetworkId)

	body, err := json.Marshal(c)

	if err != nil {
		return "", fmt.Errorf("CreateConnectorArubaEdge: marshal failed: %v", err)
	}

	data, err := ac.create(uri, body, true)

	if err != nil {
		return "", err
	}

	var result ConnectorArubaEdge
	err = json.Unmarshal([]byte(data), &result)

	if err != nil {
		return "", fmt.Errorf("CreateConnectorArubaEdge: failed to unmarshal: %v", err)
	}

	return string(result.Id), nil
}

func (ac *AlkiraClient) GetAllConnectorArubaEdge() (string, error) {
	uri := fmt.Sprintf("%s/tenantnetworks/%s/aruba-edge-connectors", ac.URI, ac.TenantNetworkId)
	data, err := ac.get(uri)

	if err != nil {
		return "", err
	}

	return string(data), nil
}

func (ac *AlkiraClient) GetConnectorArubaEdgeById(id string) (*ConnectorArubaEdge, error) {
	uri := fmt.Sprintf("%s/tenantnetworks/%s/aruba-edge-connectors/%s", ac.URI, ac.TenantNetworkId, id)

	var arubaEdge ConnectorArubaEdge

	data, err := ac.get(uri)

	if err != nil {
		return nil, err
	}

	err = json.Unmarshal([]byte(data), &arubaEdge)

	if err != nil {
		return nil, fmt.Errorf("GetConnectorArubaEdgeById: failed to unmarshal: %v", err)
	}

	return &arubaEdge, nil
}

func (ac *AlkiraClient) GetArubaInstanceConfig(serviceId string, instanceId string) (*ArubaEdgeInstanceConfig, error) {
	uri := fmt.Sprintf("%s/tenantnetworks/%s/aruba-edge-connectors/%s/instances/%s/configuration", ac.URI, ac.TenantNetworkId, serviceId, instanceId)

	data, err := ac.get(uri)

	if err != nil {
		return nil, err
	}

	//The response is string data. The entire body of the response should be interpreted together.
	//There is no json structure for ArubaEdgeInstanceConfig.
	return &ArubaEdgeInstanceConfig{Data: string(data)}, nil
}

func (ac *AlkiraClient) UpdateConnectorArubaEdge(id string, c *ConnectorArubaEdge) error {
	uri := fmt.Sprintf("%s/tenantnetworks/%s/aruba-edge-connectors/%s", ac.URI, ac.TenantNetworkId, id)

	body, err := json.Marshal(c)

	if err != nil {
		return fmt.Errorf("UpdateConnectorArubaEdge: failed to marshal request: %v", err)
	}

	return ac.update(uri, body, true)
}

func (ac *AlkiraClient) DeleteConnectorArubaEdge(id string) error {
	uri := fmt.Sprintf("%s/tenantnetworks/%s/aruba-edge-connectors/%s", ac.URI, ac.TenantNetworkId, id)

	return ac.delete(uri, true)
}

// GetConnectorArubaEdgeByName get an Aruba Edge connector by Name
func (ac *AlkiraClient) GetConnectorArubaEdgeByName(name string) (ConnectorArubaEdge, error) {
	var arubaConnector ConnectorArubaEdge

	if len(name) == 0 {
		return arubaConnector, fmt.Errorf("GetConnectorArubaEdgeByName: Invalid Connector name")
	}

	arubaConnectors, err := ac.getArubaEdgeConnectors()

	if err != nil {
		return arubaConnector, err
	}

	var result []ConnectorArubaEdge
	json.Unmarshal([]byte(arubaConnectors), &result)

	for _, l := range result {
		if l.Name == name {
			return l, nil
		}
	}

	return arubaConnector, fmt.Errorf("GetConnectorArubaEdgeByName: failed to find the connector by %s", name)
}
