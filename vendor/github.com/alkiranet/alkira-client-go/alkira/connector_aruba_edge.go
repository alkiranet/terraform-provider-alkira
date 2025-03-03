// Copyright (C) 2022-2025 Alkira Inc. All Rights Reserved.

package alkira

import (
	"encoding/json"
	"fmt"
)

type ConnectorArubaEdge struct {
	ArubaEdgeVrfMappings []ArubaEdgeVRFMappings `json:"arubaEdgeVRFMappings"`
	BillingTags          []int                  `json:"billingTags"`
	BoostMode            bool                   `json:"boostMode"`
	Cxp                  string                 `json:"cxp"`
	Group                string                 `json:"group,omitempty"`
	Id                   json.Number            `json:"id,omitempty"`              // response only
	ImplicitGroupId      int                    `json:"implicitGroupId,omitempty"` // response only
	Instances            []ArubaEdgeInstance    `json:"instances"`
	Name                 string                 `json:"name"`
	Size                 string                 `json:"size"`
	TunnelProtocol       string                 `json:"tunnelProtocol"`
	Version              string                 `json:"version"`
	Enabled              bool                   `json:"enabled"`
	Description          string                 `json:"description,omitempty"`
}

type ArubaEdgeVRFMappings struct {
	AdvertiseOnPremRoutes       bool   `json:"advertiseOnPremRoutes"`
	AlkiraSegmentId             int    `json:"alkiraSegmentId"`
	ArubaEdgeConnectSegmentName string `json:"arubaEdgeConnectSegmentName"`
	DisableInternetExit         bool   `json:"disableInternetExit"`
	GatewayBgpAsn               int    `json:"gatewayBgpAsn"`
}

type ArubaEdgeInstance struct {
	Id           json.Number `json:"id,omitempty"`
	AccountName  string      `json:"accountName"`
	CredentialId string      `json:"credentialId"`
	HostName     string      `json:"hostName"`
	Name         string      `json:"name"`
	SiteTag      string      `json:"siteTag"`
}

type ArubaEdgeInstanceConfig struct {
	Data string //The response is string data. The entire body of the
	//response should be interpreted together. There is no
	//json structure.
}

// NewConnectorArubaEdge initalize a new connector
func NewConnectorArubaEdge(ac *AlkiraClient) *AlkiraAPI[ConnectorArubaEdge] {
	uri := fmt.Sprintf("%s/tenantnetworks/%s/aruba-edge-connectors", ac.URI, ac.TenantNetworkId)
	api := &AlkiraAPI[ConnectorArubaEdge]{ac, uri, true}
	return api
}
