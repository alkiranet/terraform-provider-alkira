// Copyright (C) 2021-2025 Alkira Inc. All Rights Reserved.

package alkira

import (
	"encoding/json"
	"fmt"
)

type CiscoSdwanEdgeVrfMapping struct {
	AdvertiseOnPremRoutes bool `json:"advertiseOnPremRoutes"`
	CustomerAsn           int  `json:"customerAsn,omitempty"`
	DisableInternetExit   bool `json:"disableInternetExit"`
	SegmentId             int  `json:"segmentId"`
	Vrf                   int  `json:"vrf"`
}

type CiscoSdwanEdgeInfo struct {
	CloudInitFile          string `json:"cloudInitFile"`
	CredentialId           string `json:"credentialId"`
	HostName               string `json:"hostName"`
	Id                     int    `json:"id,omitempty"`
	Name                   string `json:"name,omitempty"`
	SshKeyPairCredentialId string `json:"sshKeyPairCredentialId,omitempty"`
}

type ConnectorCiscoSdwan struct {
	Name                 string                     `json:"name"`
	Cxp                  string                     `json:"cxp"`
	Group                string                     `json:"group,omitempty"`
	Id                   json.Number                `json:"id,omitempty"`
	Size                 string                     `json:"size"`
	Type                 string                     `json:"type,omitempty"`
	Version              string                     `json:"version"`
	TunnelProtocol       string                     `json:"tunnelProtocol,omitempty"`
	CiscoEdgeInfo        []CiscoSdwanEdgeInfo       `json:"ciscoEdgeInfo"`
	CiscoEdgeVrfMappings []CiscoSdwanEdgeVrfMapping `json:"ciscoEdgeVRFMappings"`
	BillingTags          []int                      `json:"billingTags"`
	ImplicitGroupId      int                        `json:"implicitGroupId,omitempty"`
	Enabled              bool                       `json:"enabled"`
	Description          string                     `json:"description,omitempty"`
}

// NewConnectorCiscoSdwan initialize a new connector
func NewConnectorCiscoSdwan(ac *AlkiraClient) *AlkiraAPI[ConnectorCiscoSdwan] {
	uri := fmt.Sprintf("%s/tenantnetworks/%s/ciscosdwaningresses", ac.URI, ac.TenantNetworkId)
	api := &AlkiraAPI[ConnectorCiscoSdwan]{ac, uri, true}
	return api
}
