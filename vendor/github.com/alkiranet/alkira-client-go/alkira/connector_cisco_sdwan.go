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
	CloudInitFile string `json:"cloudInitFile"`
	CredentialId  string `json:"credentialId"`
	HostName      string `json:"hostName"`
	Id            int    `json:"id,omitempty"`
	Name          string `json:"name,omitempty"`
}

type ConnectorCiscoSdwan struct {
	BillingTags          []int                      `json:"billingTags"`
	CiscoEdgeInfo        []CiscoSdwanEdgeInfo       `json:"ciscoEdgeInfo"`
	CiscoEdgeVrfMappings []CiscoSdwanEdgeVrfMapping `json:"ciscoEdgeVRFMappings"`
	Cxp                  string                     `json:"cxp"`
	Group                string                     `json:"group,omitempty"`
	Enabled              bool                       `json:"enabled,omitempty"`
	Name                 string                     `json:"name"`
	Id                   int                        `json:"id,omitempty"`
	Size                 string                     `json:"size"`
	Type                 string                     `json:"type,omitempty"`
	Version              string                     `json:"version"`
}

// CreateConnectorCiscoSdwan create a Cisco SDWAN connector
func (ac *AlkiraClient) CreateConnectorCiscoSdwan(connector *ConnectorCiscoSdwan) (string, error) {
	uri := fmt.Sprintf("%s/v1/tenantnetworks/%s/ciscosdwaningresses", ac.URI, ac.TenantNetworkId)

	// Construct the request
	body, err := json.Marshal(connector)

	if err != nil {
		return "", fmt.Errorf("CreateConnectorCiscoSdwan: failed to marshal: %v", err)
	}

	data, err := ac.create(uri, body)

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

	return ac.delete(uri)
}

// UpdateConnectorCiscoSdwan update an existing Cisco SDWAN connector by its ID
func (ac *AlkiraClient) UpdateConnectorCiscoSdwan(id string, connector *ConnectorCiscoSdwan) error {
	uri := fmt.Sprintf("%s/v1/tenantnetworks/%s/ciscosdwaningresses/%s", ac.URI, ac.TenantNetworkId, id)

	// Construct the request
	body, err := json.Marshal(connector)

	if err != nil {
		return fmt.Errorf("UpdateConnectorCiscoSdwan: failed to marshal: %v", err)
	}

	return ac.update(uri, body)
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
