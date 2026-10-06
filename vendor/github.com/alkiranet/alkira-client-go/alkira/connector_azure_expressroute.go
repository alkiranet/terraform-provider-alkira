// Copyright (C) 2020-2025 Alkira Inc. All Rights Reserved.

package alkira

import (
	"encoding/json"
	"fmt"
)

type ConnectorAzureExpressRouteSegment struct {
	SegmentName           string `json:"segmentName"`
	SegmentId             int    `json:"segmentId:omitempty"`
	CustomerAsn           int    `json:"customerAsn"`
	DisableInternetExit   bool   `json:"disableInternetExit"`
	AdvertiseOnPremRoutes bool   `json:"advertiseOnPremRoutes"`
}

type ConnectorAzureExpressRouteInstance struct {
	Name                  string                  `json:"name"`
	Id                    int                     `json:"id,omitempty"`
	ExpressRouteCircuitId string                  `json:"expressRouteCircuitId"`
	RedundantRouter       bool                    `json:"redundantRouter,omitempty"`
	LoopbackSubnet        string                  `json:"loopbackSubnet,omitempty"`
	CredentialId          string                  `json:"credentialId"`
	GatewayMacAddress     []string                `json:"gatewayMacAddresses,omitempty"`
	Vnis                  []int                   `json:"vnis,omitempty"`
	SegmentOptions        []InstanceSegmentOption `json:"segmentOptions"`
}

type InstanceSegmentOption struct {
	SegmentName      string            `json:"segmentName"`
	CustomerGateways []CustomerGateway `json:"customerGateways"`
}

type CustomerGateway struct {
	Name    string                  `json:"name"`
	Id      string                  `json:"id,omitempty"`
	Tunnels []CustomerGatewayTunnel `json:"tunnels"`
}
type CustomerGatewayTunnel struct {
	Name            string `json:"name"`
	Id              string `json:"id,omitempty"`
	Initiator       bool   `json:"initiator,omitempty"`
	ProfileId       int    `json:"profileId,omitempty"`
	IkeVersion      string `json:"ikeVersion,omitempty"`
	PreSharedKey    string `json:"preSharedKey,omitempty"`
	RemoteAuthType  string `json:"remoteAuthType,omitempty"`
	RemoteAuthValue string `json:"remoteAuthValue,omitempty"`
}

type ConnectorAzureExpressRoute struct {
	Name            string                               `json:"name"`
	Description     string                               `json:"description,omitempty"`
	Id              json.Number                          `json:"id,omitempty"`
	Size            string                               `json:"size"`
	Enabled         bool                                 `json:"enabled"`
	ImplicitGroupId int                                  `json:"implicitGroupId,omitempty"` // response only
	VhubPrefix      string                               `json:"vhubPrefix"`
	TunnelProtocol  string                               `json:"tunnelProtocol"`
	Cxp             string                               `json:"cxp"`
	Group           string                               `json:"group,omitempty"`
	Instances       []ConnectorAzureExpressRouteInstance `json:"instances,omitempty"`
	SegmentOptions  []ConnectorAzureExpressRouteSegment  `json:"segmentOptions,omitempty"`
	BillingTags     []int                                `json:"billingTags"`
}

func NewConnectorAzureExpressRoute(ac *AlkiraClient) *AlkiraAPI[ConnectorAzureExpressRoute] {
	uri := fmt.Sprintf("%s/v1/tenantnetworks/%s/azure-express-route-connectors", ac.URI, ac.TenantNetworkId)
	api := &AlkiraAPI[ConnectorAzureExpressRoute]{ac, uri, true}
	return api
}
