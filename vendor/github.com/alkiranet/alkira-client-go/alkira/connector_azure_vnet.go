// Copyright (C) 2020-2021 Alkira Inc. All Rights Reserved.

package alkira

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
)

type ConnectorVnetImportOptions struct {
	RouteImportMode string `json:"routeImportMode"`
	PrefixListIds   []int  `json:"prefixListIds,omitempty"`
}

type ConnectorVnetRouting struct {
	ImportOptions ConnectorVnetImportOptions `json:"importFromCXPOptions"`
}

type ConnectorAzureVnetRequest struct {
	BillingTags    []int                 `json:"billingTags"`
	CXP            string                `json:"cxp"`
	CredentialId   string                `json:"credentialId"`
	CustomerRegion string                `json:"customerRegion"`
	Group          string                `json:"group"`
	Name           string                `json:"name"`
	NativeServices []string              `json:"nativeServices,omitempty"`
	Segments       []string              `json:"segments"`
	Size           string                `json:"size"`
	VnetId         string                `json:"vnetId"`
	VnetRouting    *ConnectorVnetRouting `json:"vnetRouting"`
}

type ConnectorAzureVnetResponse struct {
	Id int `json:"id"`
}

// CreateConnectorAzureVnet create a AZURE-VNET connector
func (ac *AlkiraClient) CreateConnectorAzureVnet(connector *ConnectorAzureVnetRequest) (string, error) {
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

	var result ConnectorAzureVnetResponse
	json.Unmarshal([]byte(data), &result)

	if err != nil {
		return "", fmt.Errorf("CreateConnectorAzureVnet: failed to unmarshal: %v", err)
	}

	return strconv.Itoa(result.Id), nil
}

// GetConnectorAzureVnet get one AZURE-VNET connector by Id
func (ac *AlkiraClient) GetConnectorAzureVnet(id int) error {
	uri := fmt.Sprintf("%s/tenantnetworks/%s/azurevnetconnectors/%d", ac.URI, ac.TenantNetworkId, id)

	request, err := http.NewRequest("GET", uri, nil)
	request.Header.Set("Content-Type", "application/json")
	response, err := ac.Client.Do(request)

	if err != nil {
		return fmt.Errorf("GetConnectorAzureVnet: request failed: %v", err)
	}

	defer response.Body.Close()
	data, _ := ioutil.ReadAll(response.Body)

	if response.StatusCode != 200 {
		return fmt.Errorf("(%d) %s", response.StatusCode, string(data))
	}

	return nil
}

// DeleteConnectorAzureVnet delete the given AZURE-VNET connector by Id
func (ac *AlkiraClient) DeleteConnectorAzureVnet(id string) error {
	uri := fmt.Sprintf("%s/tenantnetworks/%s/azurevnetconnectors/%s", ac.URI, ac.TenantNetworkId, id)

	return ac.delete(uri)
}

// UpdateConnectorAzureVnet update an AZURE-VNET connector
func (ac *AlkiraClient) UpdateConnectorAzureVnet(id string, connector *ConnectorAzureVnetRequest) error {
	uri := fmt.Sprintf("%s/v1/tenantnetworks/%s/azurevnetconnectors/%s", ac.URI, ac.TenantNetworkId, id)

	// Construct the request
	body, err := json.Marshal(connector)

	if err != nil {
		return fmt.Errorf("UpdateConnectorAzureVnet: failed to marshal: %v", err)
	}

	return ac.update(uri, body)
}
