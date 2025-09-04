// Copyright (C) 2025 Alkira Inc. All Rights Reserved.

package alkira

import (
	"fmt"
)

// GetHealthAll get all resources health status
func (ac *AlkiraClient) GetHealthAll() (string, error) {
	uri := fmt.Sprintf("%s/tenantnetworks/%s/health", ac.URI, ac.TenantNetworkId)
	data, _, err := ac.get(uri)

	return string(data), err
}

// GetHealthOfConnector get the health status by given connector ID
func (ac *AlkiraClient) GetHealthOfConnector(connectorId string) (string, error) {

	if connectorId == "" {
		return "", fmt.Errorf("Invalid connector ID %s.", connectorId)
	}

	uri := fmt.Sprintf("%s/tenantnetworks/%s/health/connector/%s", ac.URI, ac.TenantNetworkId, connectorId)
	data, _, err := ac.get(uri)

	return string(data), err
}

// GetHealthOfConnectorInstance get the health status by given
// connector instance ID
func (ac *AlkiraClient) GetHealthOfConnectorInstance(connectorId string, instanceId string) (string, error) {

	if connectorId == "" || instanceId == "" {
		return "", fmt.Errorf("Invalid connector ID %s or instance ID %s.", connectorId, instanceId)
	}

	uri := fmt.Sprintf("%s/tenantnetworks/%s/health/connector/%s/instance/%s", ac.URI, ac.TenantNetworkId, connectorId, instanceId)
	data, _, err := ac.get(uri)

	return string(data), err
}

// GetHealthOfService get the health status by given service ID
func (ac *AlkiraClient) GetHealthOfService(serviceId string) (string, error) {

	if serviceId == "" {
		return "", fmt.Errorf("Invalid service ID %s.", serviceId)
	}

	uri := fmt.Sprintf("%s/tenantnetworks/%s/health/service/%s", ac.URI, ac.TenantNetworkId, serviceId)
	data, _, err := ac.get(uri)

	return string(data), err
}

// GetHealthOfServiceInstance get the health status by given service
// instance ID
func (ac *AlkiraClient) GetHealthOfServiceInstance(serviceId string, instanceId string) (string, error) {

	if serviceId == "" || instanceId == "" {
		return "", fmt.Errorf("Invalid service ID %s or instance ID %s.", serviceId, instanceId)
	}

	uri := fmt.Sprintf("%s/tenantnetworks/%s/health/service/%s/instance/%s", ac.URI, ac.TenantNetworkId, serviceId, instanceId)
	data, _, err := ac.get(uri)

	return string(data), err
}
