// Copyright (C) 2020-2025 Alkira Inc. All Rights Reserved.

package alkira

import (
	"encoding/json"
	"fmt"
)

type UserInputPrefixes struct {
	Id    string `json:"id,omitempty"`
	FqId  string `json:"fqId,omitempty"`
	Type  string `json:"type"`
	Value string `json:"value"`
}

type ConnectorGcpVpcExportOptions struct {
	ExportAllSubnets bool                `json:"exportAllSubnets,omitempty"`
	Prefixes         []UserInputPrefixes `json:"userInputPrefixes,omitempty"`
}

type ConnectorGcpVpcImportOptions struct {
	RouteImportMode string `json:"routeImportMode"`
	PrefixListIds   []int  `json:"prefixListIds,omitempty"`
}

type ConnectorGcpVpcRouting struct {
	ExportOptions ConnectorGcpVpcExportOptions `json:"exportToCXPOptions"`
	ImportOptions ConnectorGcpVpcImportOptions `json:"importFromCXPOptions"`
}

type ConnectorGcpVpc struct {
	BillingTags     []int                   `json:"billingTags"`
	CXP             string                  `json:"cxp"`
	CredentialId    string                  `json:"credentialId"`
	CustomerRegion  string                  `json:"customerRegion"`
	Enabled         bool                    `json:"enabled"`
	GcpRouting      *ConnectorGcpVpcRouting `json:"gcpRouting,omitempty"`
	Group           string                  `json:"group,omitempty"`
	Id              json.Number             `json:"id,omitempty"`              // response only
	ImplicitGroupId int                     `json:"implicitGroupId,omitempty"` // response only
	Name            string                  `json:"name"`
	ProjectId       string                  `json:"projectId,omitempty"`
	SecondaryCXPs   []string                `json:"secondaryCXPs,omitempty"`
	Segments        []string                `json:"segments"`
	Size            string                  `json:"size"`
	VpcId           string                  `json:"vpcId"`
	VpcName         string                  `json:"vpcName"`
	CustomerASN     int                     `json:"customerAsn,omitempty"`
	ScaleGroupId    string                  `json:"scaleGroupId,omitempty"`
	Description     string                  `json:"description,omitempty"`
}

// NewConnectorGcpVpc initialize a new connector
func NewConnectorGcpVpc(ac *AlkiraClient) *AlkiraAPI[ConnectorGcpVpc] {
	uri := fmt.Sprintf("%s/tenantnetworks/%s/gcpvpcconnectors", ac.URI, ac.TenantNetworkId)
	api := &AlkiraAPI[ConnectorGcpVpc]{ac, uri, true}
	return api
}
