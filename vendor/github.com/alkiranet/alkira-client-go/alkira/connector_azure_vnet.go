// Copyright (C) 2020-2022 Alkira Inc. All Rights Reserved.

package alkira

import (
	"encoding/json"
	"fmt"
)

type ConnectorVnetServiceRoute struct {
	Id             string   `json:"id"`
	ServiceTags    []string `json:"serviceTags"`
	NativeServices []string `json:"nativeServices,omitempty"`
	Value          string   `json:"value"`
}

type ConnectorVnetServiceRoutes struct {
	Cidrs   []ConnectorVnetServiceRoute `json:"cidrs"`
	Subnets []ConnectorVnetServiceRoute `json:"subnets"`
}

type ConnectorVnetExportOptionUserInputPrefix struct {
	Id    string `json:"id"`
	Type  string `json:"type"`
	Value string `json:"value"`
}

type ConnectorVnetExportOptions struct {
	UserInputPrefixes []ConnectorVnetExportOptionUserInputPrefix `json:"userInputPrefixes"`
}

type ConnectorVnetImportOptionsCidr struct {
	RouteImportMode string `json:"routeImportMode"`
	PrefixListIds   []int  `json:"prefixListIds"`
	Value           string `json:"value"`
}

type ConnectorVnetImportOptionsSubnet struct {
	Id              string `json:"id"`
	RouteImportMode string `json:"routeImportMode"`
	PrefixListIds   []int  `json:"prefixListIds"`
	Value           string `json:"value"`
}

type ConnectorVnetImportOptions struct {
	Cidrs           []ConnectorVnetImportOptionsCidr   `json:"cidrs,omitempty"`
	PrefixListIds   []int                              `json:"prefixListIds,omitempty"`
	RouteImportMode string                             `json:"routeImportMode"`
	Subnets         []ConnectorVnetImportOptionsSubnet `json:"subnets,omitempty"`
}

type ConnectorVnetRouting struct {
	ExportOptions ConnectorVnetExportOptions `json:"exportToCXPOptions,omitempty"`
	ImportOptions ConnectorVnetImportOptions `json:"importFromCXPOptions"`
	ServiceRoutes ConnectorVnetServiceRoutes `json:"serviceRoutes,omitempty"`
}

type ConnectorAzureVnet struct {
	BillingTags       []int                 `json:"billingTags"`
	CXP               string                `json:"cxp"`
	CredentialId      string                `json:"credentialId"`
	Group             string                `json:"group"`
	Enabled           bool                  `json:"enabled"`
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

// getAzureVnetConnectors get all Azure Vnet connectors from the given tenant network
func (ac *AlkiraClient) getAzureVnetConnectors() (string, error) {
	uri := fmt.Sprintf("%s/tenantnetworks/%s/azurevnetconnectors", ac.URI, ac.TenantNetworkId)

	data, err := ac.get(uri)
	return string(data), err
}

// CreateConnectorAzureVnet create a AZURE-VNET connector
func (ac *AlkiraClient) CreateConnectorAzureVnet(connector *ConnectorAzureVnet) (string, error) {
	uri := fmt.Sprintf("%s/tenantnetworks/%s/azurevnetconnectors", ac.URI, ac.TenantNetworkId)

	// Construct the request
	body, err := json.Marshal(connector)

	if err != nil {
		return "", fmt.Errorf("CreateConnectorAzureVnet: failed to marshal: %v", err)
	}

	data, err := ac.create(uri, body, true)

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

	return ac.delete(uri, true)
}

// UpdateConnectorAzureVnet update an AZURE-VNET connector
func (ac *AlkiraClient) UpdateConnectorAzureVnet(id string, connector *ConnectorAzureVnet) error {
	uri := fmt.Sprintf("%s/v1/tenantnetworks/%s/azurevnetconnectors/%s", ac.URI, ac.TenantNetworkId, id)

	// Construct the request
	body, err := json.Marshal(connector)

	if err != nil {
		return fmt.Errorf("UpdateConnectorAzureVnet: failed to marshal: %v", err)
	}

	return ac.update(uri, body, true)
}

// GetConnectorAzureVnetByName get an Azure VNET connector by name
func (ac *AlkiraClient) GetConnectorAzureVnetByName(name string) (ConnectorAzureVnet, error) {
	var azureVnetConnector ConnectorAzureVnet

	if len(name) == 0 {
		return azureVnetConnector, fmt.Errorf("GetConnectorAzureVnetByName: Invalid Connector name")
	}

	azureVnetConnectors, err := ac.getAzureVnetConnectors()

	if err != nil {
		return azureVnetConnector, err
	}

	var result []ConnectorAzureVnet
	json.Unmarshal([]byte(azureVnetConnectors), &result)

	for _, l := range result {
		if l.Name == name {
			return l, nil
		}
	}

	return azureVnetConnector, fmt.Errorf("GetConnectorAzureVnetByName: failed to find the connector by %s", name)
}
