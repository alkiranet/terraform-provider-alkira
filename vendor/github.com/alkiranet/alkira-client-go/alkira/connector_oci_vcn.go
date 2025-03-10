// Copyright (C) 2021-2025 Alkira Inc. All Rights Reserved.

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
	BillingTags     []int       `json:"billingTags"`
	CXP             string      `json:"cxp"`
	CredentialId    string      `json:"credentialId"`
	CustomerRegion  string      `json:"customerRegion"`
	Enabled         bool        `json:"enabled"`
	Group           string      `json:"group,omitempty"`
	Id              json.Number `json:"id,omitempty"`              // response only
	ImplicitGroupId int         `json:"implicitGroupId,omitempty"` // response only
	Name            string      `json:"name"`
	Primary         bool        `json:"primary"`
	SecondaryCXPs   []string    `json:"secondaryCXPs,omitempty"`
	Segments        []string    `json:"segments"`
	Size            string      `json:"size"`
	VcnId           string      `json:"vcnId"`
	VcnRouting      interface{} `json:"vcnRouting,omitempty"`
	Description     string      `json:"description,omitempty"`
}

// NewConnectorOciVcn new connector-oci-vcn
func NewConnectorOciVcn(ac *AlkiraClient) *AlkiraAPI[ConnectorOciVcn] {
	uri := fmt.Sprintf("%s/tenantnetworks/%s/oci-vcn-connectors", ac.URI, ac.TenantNetworkId)
	api := &AlkiraAPI[ConnectorOciVcn]{ac, uri, true}
	return api
}
