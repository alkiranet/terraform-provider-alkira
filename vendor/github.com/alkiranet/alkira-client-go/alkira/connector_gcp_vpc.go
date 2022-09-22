// Copyright (C) 2020-2022 Alkira Inc. All Rights Reserved.

package alkira

import (
	"encoding/json"
	"fmt"
)

type ConnectorGcpVpcImportOptions struct {
	RouteImportMode string `json:"routeImportMode"`
	PrefixListIds   []int  `json:"prefixListIds,omitempty"`
}

type ConnectorGcpVpcRouting struct {
	ImportOptions ConnectorGcpVpcImportOptions `json:"importFromCXPOptions"`
}

type ConnectorGcpVpc struct {
	BillingTags     []int                   `json:"billingTags"`
	CXP             string                  `json:"cxp"`
	CredentialId    string                  `json:"credentialId"`
	CustomerRegion  string                  `json:"customerRegion"`
	Enabled         bool                    `json:"enabled"`
	GcpRouting      *ConnectorGcpVpcRouting `json:"gcpRouting,omitempty"`
	Group           string                  `json:"group"`
	Id              json.Number             `json:"id,omitempty"`              // response only
	ImplicitGroupId int                     `json:"implicitGroupId,omitempty"` // response only
	Name            string                  `json:"name"`
	ProjectId       string                  `json:"projectId,omitempty"`
	Segments        []string                `json:"segments"`
	SecondaryCXPs   []string                `json:"secondaryCXPs,omitempty"`
	Size            string                  `json:"size"`
	VpcId           string                  `json:"vpcId"`
	VpcName         string                  `json:"vpcName"`
}

// getGcpVpcConnectors get all GPC VPC connectors from the given tenant network
func (ac *AlkiraClient) getGcpVpcConnectors() (string, error) {
	uri := fmt.Sprintf("%s/tenantnetworks/%s/gcpvpcconnectors", ac.URI, ac.TenantNetworkId)

	data, err := ac.get(uri)
	return string(data), err
}

// CreateConnectorGcpVpc create a GCP-VPC connector
func (ac *AlkiraClient) CreateConnectorGcpVpc(c *ConnectorGcpVpc) (string, error) {
	uri := fmt.Sprintf("%s/v1/tenantnetworks/%s/gcpvpcconnectors", ac.URI, ac.TenantNetworkId)

	// Construct the request
	body, err := json.Marshal(c)

	if err != nil {
		return "", fmt.Errorf("CreateConnectorGcpVpc: failed to marshal: %v", err)
	}

	data, err := ac.create(uri, body, true)

	if err != nil {
		return "", err
	}

	var result ConnectorGcpVpc
	err = json.Unmarshal([]byte(data), &result)

	if err != nil {
		return "", fmt.Errorf("CreateConnectorGcpVpc: parse failed: %v", err)
	}

	return string(result.Id), nil
}

// DeleteConnectorGcpVpc delete a GCP-VPC connector by Id
func (ac *AlkiraClient) DeleteConnectorGcpVpc(id string) error {
	uri := fmt.Sprintf("%s/v1/tenantnetworks/%s/gcpvpcconnectors/%s", ac.URI, ac.TenantNetworkId, id)

	return ac.delete(uri, true)
}

// UpdateConnectorGcpVpc update an GCP-VPC connector
func (ac *AlkiraClient) UpdateConnectorGcpVpc(id string, c *ConnectorGcpVpc) error {
	uri := fmt.Sprintf("%s/v1/tenantnetworks/%s/gcpvpcconnectors/%s", ac.URI, ac.TenantNetworkId, id)

	// Construct the request
	body, err := json.Marshal(c)

	if err != nil {
		return fmt.Errorf("UpdateConnectorGcpVpc: failed to marshal: %v", err)
	}

	return ac.update(uri, body, true)
}

// GetConnectorGcpVpc get a GCP-VPC connector by Id
func (ac *AlkiraClient) GetConnectorGcpVpc(id string) (*ConnectorGcpVpc, error) {
	uri := fmt.Sprintf("%s/v1/tenantnetworks/%s/gcpvpcconnectors/%s", ac.URI, ac.TenantNetworkId, id)

	data, err := ac.get(uri)

	if err != nil {
		return nil, err
	}

	var result ConnectorGcpVpc
	err = json.Unmarshal([]byte(data), &result)

	if err != nil {
		return nil, fmt.Errorf("GetConnectorGcpVpc: failed to unmarshal: %v", err)
	}

	return &result, nil
}

// GetConnectorGcpVpcByName get a GCP connector by name
func (ac *AlkiraClient) GetConnectorGcpVpcByName(name string) (ConnectorGcpVpc, error) {
	var gcpVpcConnector ConnectorGcpVpc

	if len(name) == 0 {
		return gcpVpcConnector, fmt.Errorf("GetConnectorGcpVpcByName: Invalid Connector name")
	}

	gcpVpcConnectors, err := ac.getGcpVpcConnectors()

	if err != nil {
		return gcpVpcConnector, err
	}

	var result []ConnectorGcpVpc
	json.Unmarshal([]byte(gcpVpcConnectors), &result)

	for _, l := range result {
		if l.Name == name {
			return l, nil
		}
	}

	return gcpVpcConnector, fmt.Errorf("GetConnectorGcpVpcByName: failed to find the connector by %s", name)
}
