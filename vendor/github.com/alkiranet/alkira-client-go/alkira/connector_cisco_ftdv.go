// Copyright (C) 2020-2022 Alkira Inc. All Rights Reserved.

package alkira

import (
	"encoding/json"
	"fmt"
	"strconv"
)

type CiscoFTDvInstance struct {
	Id           int    `json:"id,omitempty"`           // filled in response
	CredentialId string `json:"credentialId,omitempty"` // filled in response
	InternalName string `json:"internalName,omitempty"` // filled in response
	State        string `json:"state,omitempty"`        // filled in response
	Hostname     string `json:"hostName"`
	LicenseType  string `json:"licenseType"`
	Version      string `json:"version"`
}

type CiscoFTDvManagementServer struct {
	IPAddress string `json:"ipAddress"`
	Segment   string `json:"segment"`
	SegmentId int    `json:"segmentId"`
}

type ConnectorCiscoFTDv struct {
	Id               int                       `json:"id,omitempty"` // filled in response
	Name             string                    `json:"name"`
	GlobalCidrListId int                       `json:"globalCidrListId"`
	Size             string                    `json:"size"`
	CredentialId     string                    `json:"credentialId,omitempty"` // filled in response
	Cxp              string                    `json:"cxp"`
	ManagementServer CiscoFTDvManagementServer `json:"managementServer"`
	IpAllowList      []string                  `json:"servicesIpAllowList"`
	MaxInstanceCount int                       `json:"maxInstanceCount"`
	MinInstanceCount int                       `json:"minInstanceCount"`
	Segments         []string                  `json:"segments"`
	SegmentOptions   SegmentNameToZone         `json:"segmentOptions,omitempty"`
	Instances        []CiscoFTDvInstance       `json:"instances"`
	BillingTags      []int                     `json:"billingTags"`
	TunnelProtocol   string                    `json:"tunnelProtocol"`
	AutoScale        string                    `json:"autoScale"`
	InternalName     string                    `json:"internalName,omitempty"` // filled in response
	State            string                    `json:"state,omitempty"`        // filled in response
}

// getCiscoFTDvConnectors get all Cisco FTDv connectors from the given tenant network
func (ac *AlkiraClient) getCiscoFTDvConnectors() (string, error) {
	uri := fmt.Sprintf("%s/tenantnetworks/%s/cisco-ftdv-fw-services", ac.URI, ac.TenantNetworkId)

	data, err := ac.get(uri)
	return string(data), err
}

// CreateConnectorCiscoFTDv create a Cisco FTDv connector
func (ac *AlkiraClient) CreateConnectorCiscoFTDv(connector *ConnectorCiscoFTDv) (string, error) {
	uri := fmt.Sprintf("%s/tenantnetworks/%s/cisco-ftdv-fw-services", ac.URI, ac.TenantNetworkId)

	// Construct the request
	body, err := json.Marshal(connector)

	if err != nil {
		return "", fmt.Errorf("CreateConnectorCiscoFTDv: failed to marshal: %v", err)
	}

	data, err := ac.create(uri, body, true)

	if err != nil {
		return "", err
	}

	var result ConnectorCiscoFTDv
	json.Unmarshal([]byte(data), &result)

	if err != nil {
		return "", fmt.Errorf("CreateConnectorCiscoFTDv: failed to unmarshal: %v", err)
	}

	return strconv.Itoa(result.Id), nil
}

// GetConnectorCiscoFTDv get one Cisco FTDv connector by Id
func (ac *AlkiraClient) GetConnectorCiscoFTDv(id string) (*ConnectorCiscoFTDv, error) {
	uri := fmt.Sprintf("%s/tenantnetworks/%s/cisco-ftdv-fw-services/%s", ac.URI, ac.TenantNetworkId, id)

	data, err := ac.get(uri)

	if err != nil {
		return nil, err
	}

	var result ConnectorCiscoFTDv
	err = json.Unmarshal([]byte(data), &result)

	if err != nil {
		return nil, fmt.Errorf("GetConnectorCiscoFTDv: failed to unmarshal: %v", err)
	}

	return &result, nil
}

// DeleteConnectorCiscoFTDv delete the given Cisco FTDv connector by Id
func (ac *AlkiraClient) DeleteConnectorCiscoFTDv(id string) error {
	uri := fmt.Sprintf("%s/tenantnetworks/%s/cisco-ftdv-fw-services/%s", ac.URI, ac.TenantNetworkId, id)

	return ac.delete(uri, true)
}

// UpdateConnectorCiscoFTDv update a Cisco FTDv connector by Id
func (ac *AlkiraClient) UpdateConnectorCiscoFTDv(id string, connector *ConnectorCiscoFTDv) error {
	uri := fmt.Sprintf("%s/v1/tenantnetworks/%s/cisco-ftdv-fw-services/%s", ac.URI, ac.TenantNetworkId, id)

	// Construct the request
	body, err := json.Marshal(connector)

	if err != nil {
		return fmt.Errorf("UpdateConnectorCiscoFTDv: failed to marshal: %v", err)
	}

	return ac.update(uri, body, true)
}

// GetConnectorCiscoFTDvByName get an Azure ExpressRoute connector by name
func (ac *AlkiraClient) GetConnectorCiscoFTDvByName(name string) (ConnectorCiscoFTDv, error) {
	var ciscoFTDvConnector ConnectorCiscoFTDv

	if len(name) == 0 {
		return ciscoFTDvConnector, fmt.Errorf("GetConnectorCiscoFTDvByName: Invalid Connector name")
	}

	ciscoFTDvConnectors, err := ac.getCiscoFTDvConnectors()

	if err != nil {
		return ciscoFTDvConnector, err
	}

	var result []ConnectorCiscoFTDv
	json.Unmarshal([]byte(ciscoFTDvConnectors), &result)

	for _, l := range result {
		if l.Name == name {
			return l, nil
		}
	}

	return ciscoFTDvConnector, fmt.Errorf("GetConnectorCiscoFTDvByName: failed to find the connector by %s", name)
}
