// Copyright (C) 2022 Alkira Inc. All Rights Reserved.

package alkira

import (
	"encoding/json"
	"fmt"
)

type ConnectorAkamaiProlexic struct {
	AkamaiBgpAsn         int                                           `json:"akamaiBgpAsn"`
	BillingTags          []int                                         `json:"billingTags,omitempty"`
	ByoipOptions         []ConnectorAkamaiProlexicByoipOption          `json:"byoipOptions,omitempty"`
	CXP                  string                                        `json:"cxp"`
	CredentialId         string                                        `json:"credentialId"`
	Enabled              bool                                          `json:"enabled"`
	Group                string                                        `json:"group"`
	Id                   json.Number                                   `json:"id,omitempty"`              // response only
	ImplicitGroupId      int                                           `json:"implicitGroupId,omitempty"` // response only
	Name                 string                                        `json:"name"`
	OverlayConfiguration []ConnectorAkamaiProlexicOverlayConfiguration `json:"overlayConfiguration,omitempty"`
	Segments             []string                                      `json:"segments"`
	Size                 string                                        `json:"size"`
}

type ConnectorAkamaiProlexicByoipOption struct {
	ByoipId                   int  `json:"byoipId"`
	RouteAdvertisementEnabled bool `json:"routeAdvertisementEnabled"`
}

type ConnectorAkamaiProlexicOverlayConfiguration struct {
	AlkiraPublicIp string                            `json:"alkiraPublicIp"`
	TunnelIps      []ConnectorAkamaiProlexicTunnelIp `json:"tunnelIps"`
}

type ConnectorAkamaiProlexicTunnelIp struct {
	RanTunnelDestinationIp string `json:"ranTunnelDestinationIp"`
	AlkiraOverlayTunnelIp  string `json:"alkiraOverlayTunnelIp"`
	AkamaiOverlayTunnelIp  string `json:"akamaiOverlayTunnelIp"`
}

// getAkamaiProlexicConnectors get all akamai prolexic connectors from the given tenant network
func (ac *AlkiraClient) getAkamaiProlexicConnectors() (string, error) {
	uri := fmt.Sprintf("%s/tenantnetworks/%s/akamai-prolexic-connectors", ac.URI, ac.TenantNetworkId)

	data, err := ac.get(uri)
	return string(data), err
}

// CreateConnectorAkamaiProlexic create a Akamai Prolexic connector
func (ac *AlkiraClient) CreateConnectorAkamaiProlexic(c *ConnectorAkamaiProlexic) (string, error) {
	uri := fmt.Sprintf("%s/tenantnetworks/%s/akamai-prolexic-connectors", ac.URI, ac.TenantNetworkId)

	// Construct the request
	body, err := json.Marshal(c)

	if err != nil {
		return "", fmt.Errorf("CreateConnectorAkamaiProlexic: failed to marshal: %v", err)
	}

	data, err := ac.create(uri, body, true)

	if err != nil {
		return "", err
	}

	var result ConnectorAkamaiProlexic
	err = json.Unmarshal([]byte(data), &result)

	if err != nil {
		return "", fmt.Errorf("CreateConnectorAkamaiProlexic: parse failed: %v", err)
	}

	return string(result.Id), nil
}

// DeleteConnectorAkamaiProlexic delete a Akamai Prolexic connector by ID
func (ac *AlkiraClient) DeleteConnectorAkamaiProlexic(id string) error {
	uri := fmt.Sprintf("%s/tenantnetworks/%s/akamai-prolexic-connectors/%s", ac.URI, ac.TenantNetworkId, id)

	return ac.delete(uri, true)
}

// UpdateConnectorAkamaiProlexic update an Akamai Prolexic connector
func (ac *AlkiraClient) UpdateConnectorAkamaiProlexic(id string, c *ConnectorAkamaiProlexic) error {
	uri := fmt.Sprintf("%s/tenantnetworks/%s/akamai-prolexic-connectors/%s", ac.URI, ac.TenantNetworkId, id)

	// Construct the request
	body, err := json.Marshal(c)

	if err != nil {
		return fmt.Errorf("UpdateConnectorAkamaiProlexic: failed to marshal: %v", err)
	}

	return ac.update(uri, body, true)
}

// GetConnectorAkamaiProlexic get a Akamai Prolexic connector by ID
func (ac *AlkiraClient) GetConnectorAkamaiProlexic(id string) (*ConnectorAkamaiProlexic, error) {
	uri := fmt.Sprintf("%s/tenantnetworks/%s/akamai-prolexic-connectors/%s", ac.URI, ac.TenantNetworkId, id)

	data, err := ac.get(uri)

	if err != nil {
		return nil, err
	}

	var result ConnectorAkamaiProlexic
	err = json.Unmarshal([]byte(data), &result)

	if err != nil {
		return nil, fmt.Errorf("GetConnectorAkamaiProlexic: failed to unmarshal: %v", err)
	}

	return &result, nil
}

// GetConnectorAkamaiProlexicByName get a Akamai Prolexic connector by name
func (ac *AlkiraClient) GetConnectorAkamaiProlexicByName(name string) (ConnectorAkamaiProlexic, error) {
	var akamaiConnector ConnectorAkamaiProlexic

	if len(name) == 0 {
		return akamaiConnector, fmt.Errorf("GetConnectorAkamaiProlexicByName: Invalid Connector name")
	}

	akamaiConnectors, err := ac.getAkamaiProlexicConnectors()

	if err != nil {
		return akamaiConnector, err
	}

	var result []ConnectorAkamaiProlexic
	json.Unmarshal([]byte(akamaiConnectors), &result)

	for _, l := range result {
		if l.Name == name {
			return l, nil
		}
	}

	return akamaiConnector, fmt.Errorf("GetConnectorAkamaiProlexicByName: failed to find the connector by %s", name)
}
