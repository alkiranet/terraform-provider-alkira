// Copyright (C) 2020-2025 Alkira Inc. All Rights Reserved.

package alkira

import (
	"encoding/json"
	"fmt"
)

type ConnectorVhubImportOptions struct {
	RouteImportMode string `json:"routeImportMode,omitempty"`
	PrefixListIds   []int  `json:"prefixListIds,omitempty"`
}

type ConnectorVhubRouting struct {
	ImportOptions *ConnectorVhubImportOptions `json:"importFromCXPOptions,omitempty"`
}

type ConnectorAzureVHub struct {
	// Request fields
	Name         string                `json:"name"`
	Description  string                `json:"description,omitempty"`
	Enabled      bool                  `json:"enabled"`
	CXP          string                `json:"cxp"`
	Segments     []string              `json:"segments"`
	Group        string                `json:"group,omitempty"`
	Size         string                `json:"size"`
	BillingTags  []int                 `json:"billingTags,omitempty"`
	CredentialId string                `json:"credentialId"`
	VirtualHubId string                `json:"virtualHubId"`
	ScaleGroupId string                `json:"scaleGroupId"`
	VhubRouting  *ConnectorVhubRouting `json:"vhubRouting"`

	// Response-only fields
	Id                json.Number `json:"id,omitempty"`
	ImplicitGroupId   int         `json:"implicitGroupId,omitempty"`
	GroupId           int         `json:"groupId,omitempty"`
	SegmentIds        []int       `json:"segmentIds,omitempty"`
	SubscriptionId    string      `json:"subscriptionId,omitempty"`
	ResourceGroupName string      `json:"resourceGroupName,omitempty"`
	VirtualWanId      string      `json:"virtualWanId,omitempty"`
	VirtualHubName    string      `json:"virtualHubName,omitempty"`
	VpnGatewayId      string      `json:"vpnGatewayId,omitempty"`
	ASN               int         `json:"asn,omitempty"`
}

func NewConnectorAzureVHub(ac *AlkiraClient) *AlkiraAPI[ConnectorAzureVHub] {
	uri := fmt.Sprintf("%s/tenantnetworks/%s/azure-vhub-connectors", ac.URI, ac.TenantNetworkId)
	api := &AlkiraAPI[ConnectorAzureVHub]{ac, uri, true}
	return api
}
