// Copyright (C) 2020-2023 Alkira Inc. All Rights Reserved.

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
	Name                  string   `json:"name"`
	Id                    int      `json:"id,omitempty"`
	ExpressRouteCircuitId string   `json:"expressRouteCircuitId"`
	RedundantRouter       bool     `json:"redundantRouter,omitempty"`
	LoopbackSubnet        string   `json:"loopbackSubnet,omitempty"`
	CredentialId          string   `json:"credentialId"`
	GatewayMacAddress     []string `json:"gatewayMacAddresses,omitempty"`
	Vnis                  []int    `json:"vnis,omitempty"`
}

type ConnectorAzureExpressRoute struct {
	Name            string                               `json:"name"`
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
