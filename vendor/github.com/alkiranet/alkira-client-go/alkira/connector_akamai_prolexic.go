// Copyright (C) 2022 Alkira Inc. All Rights Reserved.

package alkira

import (
	"encoding/json"
	"fmt"
)

type ConnectorAkamaiProlexic struct {
	Size                 string                                        `json:"size"`
	Enabled              bool                                          `json:"enabled,omitempty"`
	CredentialId         string                                        `json:"credentialId"`
	Segments             []string                                      `json:"segments"`
	OverlayConfiguration []ConnectorAkamaiProlexicOverlayConfiguration `json:"overlayConfiguration,omitempty"`
	AkamaiBgpAsn         int                                           `json:"akamaiBgpAsn"`
	ByoipOptions         []ConnectorAkamaiProlexicByoipOption          `json:"byoipOptions,omitempty"`
	Name                 string                                        `json:"name"`
	CXP                  string                                        `json:"cxp"`
	Group                string                                        `json:"group"`
	BillingTags          []int                                         `json:"billingTags,omitempty"`
	Id                   json.Number                                   `json:"id,omitempty"` // response only
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

// CreateConnectorAkamaiProlexic create a Akamai Prolexic connector
func (ac *AlkiraClient) CreateConnectorAkamaiProlexic(c *ConnectorAkamaiProlexic) (string, error) {
	uri := fmt.Sprintf("%s/tenantnetworks/%s/akamai-prolexic-connectors", ac.URI, ac.TenantNetworkId)

	// Construct the request
	body, err := json.Marshal(c)

	if err != nil {
		return "", fmt.Errorf("CreateConnectorAkamaiProlexic: failed to marshal: %v", err)
	}

	data, err := ac.create(uri, body)

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

	return ac.delete(uri)
}

// UpdateConnectorAkamaiProlexic update an Akamai Prolexic connector
func (ac *AlkiraClient) UpdateConnectorAkamaiProlexic(id string, c *ConnectorAkamaiProlexic) error {
	uri := fmt.Sprintf("%s/tenantnetworks/%s/akamai-prolexic-connectors/%s", ac.URI, ac.TenantNetworkId, id)

	// Construct the request
	body, err := json.Marshal(c)

	if err != nil {
		return fmt.Errorf("UpdateConnectorAkamaiProlexic: failed to marshal: %v", err)
	}

	return ac.update(uri, body)
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
