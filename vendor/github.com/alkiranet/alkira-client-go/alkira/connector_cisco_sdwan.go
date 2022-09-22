// Copyright (C) 2021-2022 Alkira Inc. All Rights Reserved.

package alkira

import (
	"encoding/json"
	"fmt"
	"strconv"
)

type CiscoSdwanEdgeVrfMapping struct {
	AdvertiseOnPremRoutes bool `json:"advertiseOnPremRoutes"`
	CustomerAsn           int  `json:"customerAsn,omitempty"`
	DisableInternetExit   bool `json:"disableInternetExit"`
	SegmentId             int  `json:"segmentId"`
	Vrf                   int  `json:"vrf"`
}

type CiscoSdwanEdgeInfo struct {
	CloudInitFile          string `json:"cloudInitFile"`
	CredentialId           string `json:"credentialId"`
	HostName               string `json:"hostName"`
	Id                     int    `json:"id,omitempty"`
	Name                   string `json:"name,omitempty"`
	SshKeyPairCredentialId string `json:"sshKeyPairCredentialId,omitempty"`
}

type ConnectorCiscoSdwan struct {
	BillingTags          []int                      `json:"billingTags"`
	CiscoEdgeInfo        []CiscoSdwanEdgeInfo       `json:"ciscoEdgeInfo"`
	CiscoEdgeVrfMappings []CiscoSdwanEdgeVrfMapping `json:"ciscoEdgeVRFMappings"`
	Cxp                  string                     `json:"cxp"`
	Group                string                     `json:"group,omitempty"`
	Enabled              bool                       `json:"enabled"`
	Name                 string                     `json:"name"`
	Id                   int                        `json:"id,omitempty"`              // response only
	ImplicitGroupId      int                        `json:"implicitGroupId,omitempty"` // response only
	Size                 string                     `json:"size"`
	Type                 string                     `json:"type,omitempty"`
	Version              string                     `json:"version"`
}

// getCiscoSdwanConnecots get all Cisco Sdwan Connectors from the given tenant network
func (ac *AlkiraClient) getCiscoSdwanConnectors() (string, error) {
	uri := fmt.Sprintf("%s/tenantnetworks/%s/ciscosdwaningresses", ac.URI, ac.TenantNetworkId)

	data, err := ac.get(uri)
	return string(data), err
}

// CreateConnectorCiscoSdwan create a Cisco SDWAN connector
func (ac *AlkiraClient) CreateConnectorCiscoSdwan(connector *ConnectorCiscoSdwan) (string, error) {
	uri := fmt.Sprintf("%s/v1/tenantnetworks/%s/ciscosdwaningresses", ac.URI, ac.TenantNetworkId)

	// Construct the request
	body, err := json.Marshal(connector)

	if err != nil {
		return "", fmt.Errorf("CreateConnectorCiscoSdwan: failed to marshal: %v", err)
	}

	data, err := ac.create(uri, body, true)

	if err != nil {
		return "", err
	}

	var result ConnectorCiscoSdwan
	err = json.Unmarshal([]byte(data), &result)

	if err != nil {
		return "", fmt.Errorf("CreateConnectorCiscoSdwan: failed to unmarshal: %v", err)
	}

	return strconv.Itoa(result.Id), nil
}

// DeleteConnectorCiscoSdwan delete an existing Cisco SDWAN connector by its ID
func (ac *AlkiraClient) DeleteConnectorCiscoSdwan(id string) error {
	uri := fmt.Sprintf("%s/v1/tenantnetworks/%s/ciscosdwaningresses/%s", ac.URI, ac.TenantNetworkId, id)

	return ac.delete(uri, true)
}

// UpdateConnectorCiscoSdwan update an existing Cisco SDWAN connector by its ID
func (ac *AlkiraClient) UpdateConnectorCiscoSdwan(id string, connector *ConnectorCiscoSdwan) error {
	uri := fmt.Sprintf("%s/v1/tenantnetworks/%s/ciscosdwaningresses/%s", ac.URI, ac.TenantNetworkId, id)

	// Construct the request
	body, err := json.Marshal(connector)

	if err != nil {
		return fmt.Errorf("UpdateConnectorCiscoSdwan: failed to marshal: %v", err)
	}

	return ac.update(uri, body, true)
}

// GetConnectorCiscoSdwan get an existing Cisco SDWAN connector by its Id
func (ac *AlkiraClient) GetConnectorCiscoSdwan(id string) (*ConnectorCiscoSdwan, error) {
	uri := fmt.Sprintf("%s/tenantnetworks/%s/ciscosdwaningresses/%s", ac.URI, ac.TenantNetworkId, id)

	var connector ConnectorCiscoSdwan

	data, err := ac.get(uri)

	if err != nil {
		return nil, err
	}

	err = json.Unmarshal([]byte(data), &connector)

	if err != nil {
		return nil, fmt.Errorf("GetConnectorIPSec: failed to unmarshal: %v", err)
	}

	return &connector, nil
}

// GetConnectorCiscoSdwanByName get a Cisco Sdwan Connector by Name
func (ac *AlkiraClient) GetConnectorCiscoSdwanByName(name string) (ConnectorCiscoSdwan, error) {
	var ciscoSdwanConnector ConnectorCiscoSdwan

	if len(name) == 0 {
		return ciscoSdwanConnector, fmt.Errorf("GetConnectorCiscoSdwanByName: Invalid Connector name")
	}

	ciscoSdwanConnectors, err := ac.getCiscoSdwanConnectors()

	if err != nil {
		return ciscoSdwanConnector, err
	}

	var result []ConnectorCiscoSdwan
	json.Unmarshal([]byte(ciscoSdwanConnectors), &result)

	for _, l := range result {
		if l.Name == name {
			return l, nil
		}
	}

	return ciscoSdwanConnector, fmt.Errorf("GetConnectorCiscoSdwanByName: failed to find the connector by %s", name)
}
