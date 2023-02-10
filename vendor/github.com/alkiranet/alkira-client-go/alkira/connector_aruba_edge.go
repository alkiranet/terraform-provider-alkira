// Copyright (C) 2022-2023 Alkira Inc. All Rights Reserved.

package alkira

import (
	"encoding/json"
	"fmt"
)

type ConnectorArubaEdge struct {
	ArubaEdgeVrfMapping []ArubaEdgeVRFMapping `json:"arubaEdgeVRFMapping,omitempty"`
	BillingTags         []int                 `json:"billingTags"`
	BoostMode           bool                  `json:"boostMode"`
	Cxp                 string                `json:"cxp"`
	GatewayBgpAsn       int                   `json:"gatewayBgpAsn"`
	Group               string                `json:"group,omitempty"`
	Id                  json.Number           `json:"id,omitempty"`              // response only
	ImplicitGroupId     int                   `json:"implicitGroupId,omitempty"` // response only
	Instances           []ArubaEdgeInstance   `json:"instances"`
	Name                string                `json:"name"`
	Segments            []string              `json:"segments"`
	Size                string                `json:"size"`
	TunnelProtocol      string                `json:"tunnelProtocol"`
	Version             string                `json:"version"`
}

type ArubaEdgeVRFMapping struct {
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
	api := &AlkiraAPI[ConnectorArubaEdge]{ac, uri}
	return api
}
