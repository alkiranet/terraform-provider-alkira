// Copyright (C) 2021-2022 Alkira Inc. All Rights Reserved.

package alkira

import (
	"encoding/json"
	"fmt"
)

type ConnectorOciVcnInputPrefixes struct {
	Id    string `json:"id,omitempty"`
	Type  string `json:"type"`
	Value string `json:"value"`
}

type ConnectorOciVcnExportOptions struct {
	Mode     string                         `json:"routeExportMode,omitempty"`
	Prefixes []ConnectorOciVcnInputPrefixes `json:"userInputPrefixes,omitempty"`
}

type ConnectorOciVcnRouteTables struct {
	Id            string `json:"id"`
	PrefixListIds []int  `json:"prefixListIds"`
	Mode          string `json:"routeImportMode"`
}

type ConnectorOciVcnImportOptions struct {
	RouteTables []ConnectorOciVcnRouteTables `json:"routeTables"`
}

type ConnectorOciVcnRouting struct {
	Export interface{} `json:"exportToCXPOptions"`
	Import interface{} `json:"importFromCXPOptions"`
}

type ConnectorOciVcn struct {
	BillingTags    []int       `json:"billingTags"`
	CXP            string      `json:"cxp"`
	CredentialId   string      `json:"credentialId"`
	CustomerRegion string      `json:"customerRegion"`
	Enabled        bool        `json:"enabled"`
	Group          string      `json:"group"`
	Id             json.Number `json:"id,omitempty"`
	Name           string      `json:"name"`
	Primary        bool        `json:"primary"`
	Segments       []string    `json:"segments"`
	Size           string      `json:"size"`
	VcnId          string      `json:"vcnId"`
	VcnRouting     interface{} `json:"vcnRouting,omitempty"`
}

// CreateConnectorOciVcn create an OCI-VCN connector
func (ac *AlkiraClient) CreateConnectorOciVcn(connector *ConnectorOciVcn) (string, error) {
	uri := fmt.Sprintf("%s/v1/tenantnetworks/%s/oci-vcn-connectors", ac.URI, ac.TenantNetworkId)

	// Construct the request
	body, err := json.Marshal(connector)

	if err != nil {
		return "", fmt.Errorf("CreateConnectorOciVcn: failed to marshal: %v", err)
	}

	data, err := ac.create(uri, body)

	if err != nil {
		return "", err
	}

	var result ConnectorOciVcn
	err = json.Unmarshal([]byte(data), &result)

	if err != nil {
		return "", fmt.Errorf("CreateConnectorOciVcn: failed to unmarshal: %v", err)
	}

	return string(result.Id), nil
}

// DeleteConnectorOciVcn delete an OCI-VCN connector
func (ac *AlkiraClient) DeleteConnectorOciVcn(id string) error {
	uri := fmt.Sprintf("%s/v1/tenantnetworks/%s/oci-vcn-connectors/%s", ac.URI, ac.TenantNetworkId, id)

	return ac.delete(uri)
}

// UpdateConnectorOciVcn update an OCI-VCN connector
func (ac *AlkiraClient) UpdateConnectorOciVcn(id string, connector *ConnectorOciVcn) error {
	uri := fmt.Sprintf("%s/v1/tenantnetworks/%s/oci-vcn-connectors/%s", ac.URI, ac.TenantNetworkId, id)

	// Construct the request
	body, err := json.Marshal(connector)

	if err != nil {
		return fmt.Errorf("UpdateConnectorOciVcn: failed to marshal: %v", err)
	}

	return ac.update(uri, body)
}

// GetConnectorOciVcn get one OCI-VCN connector by Id
func (ac *AlkiraClient) GetConnectorOciVcn(id string) (*ConnectorOciVcn, error) {
	uri := fmt.Sprintf("%s/tenantnetworks/%s/oci-vcn-connectors/%s", ac.URI, ac.TenantNetworkId, id)

	data, err := ac.get(uri)

	if err != nil {
		return nil, err
	}

	var result ConnectorOciVcn
	err = json.Unmarshal([]byte(data), &result)

	if err != nil {
		return nil, fmt.Errorf("GetConnectorOciVcn: failed to unmarshal: %v", err)
	}

	return &result, nil
}
