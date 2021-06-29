// Copyright (C) 2020-2021 Alkira Inc. All Rights Reserved.

package alkira

import (
	"encoding/json"
	"fmt"
	"strconv"
)

// Structs for "VPC Routing" options
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

// AWS-VPC connector
type ConnectorAwsVpcRequest struct {
	BillingTags    []int       `json:"billingTags"`
	CXP            string      `json:"cxp"`
	CredentialId   string      `json:"credentialId"`
	CustomerName   string      `json:"customerName"`
	CustomerRegion string      `json:"customerRegion"`
	Group          string      `json:"group"`
	Name           string      `json:"name"`
	Segments       []string    `json:"segments"`
	Size           string      `json:"size"`
	VpcId          string      `json:"vpcId"`
	VpcOwnerId     string      `json:"vpcOwnerId"`
	VpcRouting     interface{} `json:"vpcRouting"`
}

type ConnectorAwsVpcResponse struct {
	Id int `json:"id"`
}

// CreateConnectorAwsVPC create an AWS-VPC connector
func (ac *AlkiraClient) CreateConnectorAwsVpc(connector *ConnectorAwsVpcRequest) (string, error) {
	uri := fmt.Sprintf("%s/v1/tenantnetworks/%s/awsvpcconnectors", ac.URI, ac.TenantNetworkId)

	// Construct the request
	body, err := json.Marshal(connector)

	if err != nil {
		return "", fmt.Errorf("CreateConnectorAwsVpc: failed to marshal: %v", err)
	}

	data, err := ac.create(uri, body)

	if err != nil {
		return "", err
	}

	var result ConnectorAwsVpcResponse
	err = json.Unmarshal([]byte(data), &result)

	if err != nil {
		return "", fmt.Errorf("CreateConnectorAwsVpc: failed to unmarshal: %v", err)
	}

	return strconv.Itoa(result.Id), nil
}

// DeleteConnectorAwsVpc delete an AWS-VPC connector
func (ac *AlkiraClient) DeleteConnectorAwsVpc(id string) error {
	uri := fmt.Sprintf("%s/v1/tenantnetworks/%s/awsvpcconnectors/%s", ac.URI, ac.TenantNetworkId, id)

	return ac.delete(uri)
}

// UpdateConnectorAwsVPC update an AWS-VPC connector
func (ac *AlkiraClient) UpdateConnectorAwsVpc(id string, connector *ConnectorAwsVpcRequest) error {
	uri := fmt.Sprintf("%s/v1/tenantnetworks/%s/awsvpcconnectors/%s", ac.URI, ac.TenantNetworkId, id)

	// Construct the request
	body, err := json.Marshal(connector)

	if err != nil {
		return fmt.Errorf("UpdateConnectorAwsVpc: failed to marshal: %v", err)
	}

	return ac.update(uri, body)
}
