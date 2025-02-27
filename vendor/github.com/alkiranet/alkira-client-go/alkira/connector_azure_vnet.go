// Copyright (C) 2020-2025 Alkira Inc. All Rights Reserved.

package alkira

import (
	"encoding/json"
	"fmt"
)

type ConnectorVnetUdrList struct {
	Id         string `json:"id"`
	UdrListIds []int  `json:"udrListIds"`
	Value      string `json:"value"`
}

type ConnectorVnetUdrLists struct {
	Cidrs   []ConnectorVnetUdrList `json:"cidrs"`
	Subnets []ConnectorVnetUdrList `json:"subnets"`
}

type ConnectorVnetServiceRoute struct {
	Id                 string   `json:"id"`
	ServiceTags        []string `json:"serviceTags"`
	NativeServices     []string `json:"nativeServices,omitempty"`
	NativeServiceNames []string `json:"nativeServiceNames,omitempty"`
	Value              string   `json:"value"`
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
	UdrLists      ConnectorVnetUdrLists      `json:"udrLists,omitempty"`
}

type ConnectorAzureVnet struct {
	BillingTags                       []int                 `json:"billingTags"`
	CXP                               string                `json:"cxp"`
	CredentialId                      string                `json:"credentialId"`
	Group                             string                `json:"group,omitempty"`
	Enabled                           bool                  `json:"enabled"`
	Id                                json.Number           `json:"id,omitempty"`              // RESPONSE ONLY
	ImplicitGroupId                   int                   `json:"implicitGroupId,omitempty"` // RESPONSE ONLY
	Name                              string                `json:"name"`
	NativeServiceNames                []string              `json:"nativeServiceNames,omitempty"`
	ResourceGroupName                 string                `json:"resourceGroupName,omitempty"`
	SecondaryCXPs                     []string              `json:"secondaryCXPs,omitempty"`
	Segments                          []string              `json:"segments"`
	ServiceTags                       []string              `json:"serviceTags,omitempty"`
	Size                              string                `json:"size"`
	VnetId                            string                `json:"vnetId"`
	ConnectionMode                    string                `json:"connectionMode,omitempty"`
	VnetRouting                       *ConnectorVnetRouting `json:"vnetRouting"`
	CustomerASN                       int                   `json:"customerAsn,omitempty"`
	ScaleGroupId                      string                `json:"scaleGroupId,omitempty"`
	PeeringGatewayCxpId               int                   `json:"cxpPeeringGatewayId,omitempty"`
	DirectInterVNETCommunicationGroup string                `json:"directInterVnetCommunicationGroup,omitempty"`
	UdrListIds                        []int                 `json:"udrListIds,omitempty"`
	Description                       string                `json:"description,omitempty"`
}

// NewConnectorAzureVnet initalize a new connector
func NewConnectorAzureVnet(ac *AlkiraClient) *AlkiraAPI[ConnectorAzureVnet] {
	uri := fmt.Sprintf("%s/tenantnetworks/%s/azurevnetconnectors", ac.URI, ac.TenantNetworkId)
	api := &AlkiraAPI[ConnectorAzureVnet]{ac, uri, true}
	return api
}
