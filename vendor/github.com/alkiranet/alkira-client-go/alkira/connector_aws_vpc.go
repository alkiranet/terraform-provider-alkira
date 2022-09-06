// Copyright (C) 2020-2022 Alkira Inc. All Rights Reserved.

package alkira

import (
	"encoding/json"
	"fmt"
)

type InputPrefixes struct {
	Id    string `json:"id,omitempty"`
	Type  string `json:"type"`
	Value string `json:"value"`
}

type ExportOptions struct {
	Mode     string          `json:"routeExportMode,omitempty"`
	Prefixes []InputPrefixes `json:"userInputPrefixes,omitempty"`
}

type RouteTables struct {
	Id            string `json:"id"`
	PrefixListIds []int  `json:"prefixListIds"`
	Mode          string `json:"routeImportMode"`
}

type ImportOptions struct {
	RouteTables []RouteTables `json:"routeTables"`
}

type ConnectorAwsVpcRouting struct {
	Export interface{} `json:"exportToCXPOptions"`
	Import interface{} `json:"importFromCXPOptions"`
}

type TgwAttachment struct {
	SubnetId         string `json:"subnetId"`
	AvailabilityZone string `json:"availabilityZone"`
}

type ConnectorAwsVpc struct {
	BillingTags                        []int           `json:"billingTags"`
	CXP                                string          `json:"cxp"`
	CredentialId                       string          `json:"credentialId"`
	CustomerName                       string          `json:"customerName"`
	CustomerRegion                     string          `json:"customerRegion"`
	DirectInterVPCCommunicationEnabled bool            `json:"directInterVPCCommunicationEnabled"`
	Enabled                            bool            `json:"enabled"`
	Group                              string          `json:"group"`
	Id                                 json.Number     `json:"id,omitempty"`
	Name                               string          `json:"name"`
	Segments                           []string        `json:"segments"`
	Size                               string          `json:"size"`
	TgwAttachments                     []TgwAttachment `json:"tgwAttachments,omitempty"`
	VpcId                              string          `json:"vpcId"`
	VpcOwnerId                         string          `json:"vpcOwnerId"`
	VpcRouting                         interface{}     `json:"vpcRouting"`
}

// getAwsVpcConnectors get all aws vpc connectors from the given tenant network
func (ac *AlkiraClient) getAwsVpcConnectors() (string, error) {
	uri := fmt.Sprintf("%s/tenantnetworks/%s/awsvpcconnectors", ac.URI, ac.TenantNetworkId)

	data, err := ac.get(uri)
	return string(data), err
}

// CreateConnectorAwsVPC create an AWS-VPC connector
func (ac *AlkiraClient) CreateConnectorAwsVpc(connector *ConnectorAwsVpc) (string, error) {
	uri := fmt.Sprintf("%s/v1/tenantnetworks/%s/awsvpcconnectors", ac.URI, ac.TenantNetworkId)

	// Construct the request
	body, err := json.Marshal(connector)

	if err != nil {
		return "", fmt.Errorf("CreateConnectorAwsVpc: failed to marshal: %v", err)
	}

	data, err := ac.create(uri, body, true)

	if err != nil {
		return "", err
	}

	var result ConnectorAwsVpc
	err = json.Unmarshal([]byte(data), &result)

	if err != nil {
		return "", fmt.Errorf("CreateConnectorAwsVpc: failed to unmarshal: %v", err)
	}

	return string(result.Id), nil
}

// DeleteConnectorAwsVpc delete an AWS-VPC connector
func (ac *AlkiraClient) DeleteConnectorAwsVpc(id string) error {
	uri := fmt.Sprintf("%s/v1/tenantnetworks/%s/awsvpcconnectors/%s", ac.URI, ac.TenantNetworkId, id)

	return ac.delete(uri, true)
}

// UpdateConnectorAwsVPC update an AWS-VPC connector
func (ac *AlkiraClient) UpdateConnectorAwsVpc(id string, connector *ConnectorAwsVpc) error {
	uri := fmt.Sprintf("%s/v1/tenantnetworks/%s/awsvpcconnectors/%s", ac.URI, ac.TenantNetworkId, id)

	// Construct the request
	body, err := json.Marshal(connector)

	if err != nil {
		return fmt.Errorf("UpdateConnectorAwsVpc: failed to marshal: %v", err)
	}

	return ac.update(uri, body, true)
}

// GetConnectorAwsVpc get one AWS-VPC connector by Id
func (ac *AlkiraClient) GetConnectorAwsVpc(id string) (*ConnectorAwsVpc, error) {
	uri := fmt.Sprintf("%s/tenantnetworks/%s/awsvpcconnectors/%s", ac.URI, ac.TenantNetworkId, id)

	data, err := ac.get(uri)

	if err != nil {
		return nil, err
	}

	var result ConnectorAwsVpc
	err = json.Unmarshal([]byte(data), &result)

	if err != nil {
		return nil, fmt.Errorf("GetConnectorAwsVpc: failed to unmarshal: %v", err)
	}

	return &result, nil
}

// GetConnectorAwsVpcByName get an Aws Vpc connector by name
func (ac *AlkiraClient) GetConnectorAwsVpcByName(name string) (ConnectorAwsVpc, error) {
	var awsVpcConnector ConnectorAwsVpc

	if len(name) == 0 {
		return awsVpcConnector, fmt.Errorf("GetConnectorAwsVpcByName: Invalid Connector name")
	}

	awsVpcConnectors, err := ac.getAwsVpcConnectors()

	if err != nil {
		return awsVpcConnector, err
	}

	var result []ConnectorAwsVpc
	json.Unmarshal([]byte(awsVpcConnectors), &result)

	for _, l := range result {
		if l.Name == name {
			return l, nil
		}
	}

	return awsVpcConnector, fmt.Errorf("GetConnectorAwsVpcByName: failed to find the connector by %s", name)
}
