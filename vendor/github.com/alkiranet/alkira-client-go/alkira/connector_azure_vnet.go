// Copyright (C) 2020-2022 Alkira Inc. All Rights Reserved.

package alkira

import (
	"encoding/json"
	"fmt"
)

type ConnectorVnetImportOptions struct {
	RouteImportMode string `json:"routeImportMode"`
	PrefixListIds   []int  `json:"prefixListIds,omitempty"`
}

type ConnectorVnetRouting struct {
	ImportOptions ConnectorVnetImportOptions `json:"importFromCXPOptions"`
}

type ConnectorAzureVnet struct {
	BillingTags       []int                 `json:"billingTags"`
	CXP               string                `json:"cxp"`
	CredentialId      string                `json:"credentialId"`
	Group             string                `json:"group"`
	Enabled           bool                  `json:"enabled,omitempty"`
	Id                json.Number           `json:"id,omitempty"`
	Name              string                `json:"name"`
	NativeServices    []string              `json:"nativeServices,omitempty"`
	ResourceGroupName string                `json:"resourceGroupName,omitempty"`
	Segments          []string              `json:"segments"`
	ServiceTags       []string              `json:"serviceTags,omitempty"`
	Size              string                `json:"size"`
	VnetId            string                `json:"vnetId"`
	VnetRouting       *ConnectorVnetRouting `json:"vnetRouting"`
}

// CreateConnectorAzureVnet create a AZURE-VNET connector
func (ac *AlkiraClient) CreateConnectorAzureVnet(connector *ConnectorAzureVnet) (string, error) {
	uri := fmt.Sprintf("%s/tenantnetworks/%s/azurevnetconnectors", ac.URI, ac.TenantNetworkId)

	// Construct the request
	body, err := json.Marshal(connector)

	if err != nil {
		return "", fmt.Errorf("CreateConnectorAzureVnet: failed to marshal: %v", err)
	}

	data, err := ac.create(uri, body)

	if err != nil {
		return "", err
	}

	var result ConnectorAzureVnet
	json.Unmarshal([]byte(data), &result)

	if err != nil {
		return "", fmt.Errorf("CreateConnectorAzureVnet: failed to unmarshal: %v", err)
	}

	return string(result.Id), nil
}

// GetConnectorAzureVnet get one AZURE-VNET connector by Id
func (ac *AlkiraClient) GetConnectorAzureVnet(id string) (*ConnectorAzureVnet, error) {
	uri := fmt.Sprintf("%s/tenantnetworks/%s/azurevnetconnectors/%s", ac.URI, ac.TenantNetworkId, id)

	data, err := ac.get(uri)

	if err != nil {
		return nil, err
	}

	var result ConnectorAzureVnet
	err = json.Unmarshal([]byte(data), &result)

	if err != nil {
		return nil, fmt.Errorf("GetConnectorAzureVnet: failed to unmarshal: %v", err)
	}

	return &result, nil
}

// DeleteConnectorAzureVnet delete the given AZURE-VNET connector by Id
func (ac *AlkiraClient) DeleteConnectorAzureVnet(id string) error {
	uri := fmt.Sprintf("%s/tenantnetworks/%s/azurevnetconnectors/%s", ac.URI, ac.TenantNetworkId, id)

	return ac.delete(uri)
}

// UpdateConnectorAzureVnet update an AZURE-VNET connector
func (ac *AlkiraClient) UpdateConnectorAzureVnet(id string, connector *ConnectorAzureVnet) error {
	uri := fmt.Sprintf("%s/v1/tenantnetworks/%s/azurevnetconnectors/%s", ac.URI, ac.TenantNetworkId, id)

	// Construct the request
	body, err := json.Marshal(connector)

	if err != nil {
		return fmt.Errorf("UpdateConnectorAzureVnet: failed to marshal: %v", err)
	}

	return ac.update(uri, body)
}
