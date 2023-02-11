// Copyright (C) 2021-2023 Alkira Inc. All Rights Reserved.

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
	BillingTags          []int                      `json:"billingTags"`
	CiscoEdgeInfo        []CiscoSdwanEdgeInfo       `json:"ciscoEdgeInfo"`
	CiscoEdgeVrfMappings []CiscoSdwanEdgeVrfMapping `json:"ciscoEdgeVRFMappings"`
	Cxp                  string                     `json:"cxp"`
	Group                string                     `json:"group,omitempty"`
	Enabled              bool                       `json:"enabled"`
	Name                 string                     `json:"name"`
	Id                   json.Number                `json:"id,omitempty"`              // response only
	ImplicitGroupId      int                        `json:"implicitGroupId,omitempty"` // response only
	Size                 string                     `json:"size"`
	Type                 string                     `json:"type,omitempty"`
	Version              string                     `json:"version"`
}

// NewConnectorCiscoSdwan initialize a new connector
func NewConnectorCiscoSdwan(ac *AlkiraClient) *AlkiraAPI[ConnectorCiscoSdwan] {
	uri := fmt.Sprintf("%s/tenantnetworks/%s/ciscosdwaningresses", ac.URI, ac.TenantNetworkId)
	api := &AlkiraAPI[ConnectorCiscoSdwan]{ac, uri}
	return api
}
