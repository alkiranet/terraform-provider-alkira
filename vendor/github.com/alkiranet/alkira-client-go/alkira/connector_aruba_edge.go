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
	Id                  json.Number           `json:"id,omitempty"`
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
