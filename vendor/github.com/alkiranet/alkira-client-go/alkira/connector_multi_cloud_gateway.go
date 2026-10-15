// Copyright (C) 2026 Alkira Inc. All Rights Reserved.

package alkira

import (
	"encoding/json"
	"fmt"
)

type ConnectorMultiCloudGateway struct {
	Name                string                                    `json:"name"`
	Description         string                                    `json:"description,omitempty"`
	Cxp                 string                                    `json:"cxp"`
	Group               string                                    `json:"group,omitempty"`
	Size                string                                    `json:"size"`
	Enabled             bool                                      `json:"enabled"`
	BillingTags         []int                                     `json:"billingTags,omitempty"`
	Id                  json.Number                               `json:"id,omitempty"`              // RESPONSE ONLY
	ConfigStatus        string                                    `json:"configStatus,omitempty"`    // RESPONSE ONLY
	ImplicitGroupId     int                                       `json:"implicitGroupId,omitempty"` // RESPONSE ONLY
	MultiCloudGatewayId int                                       `json:"multiCloudGatewayId"`
	LoopbackPrefix      string                                    `json:"loopbackPrefix"`
	SegmentOptions      []ConnectorMultiCloudGatewaySegmentOption `json:"segmentOptions"`
	Instances           []ConnectorMultiCloudGatewayInstance      `json:"instances"`
}

type ConnectorMultiCloudGatewaySegmentOption struct {
	SegmentName           string `json:"segmentName"`
	AdvertiseOnPremRoutes bool   `json:"advertiseOnPremRoutes,omitempty"`
	DisableInternetExit   bool   `json:"disableInternetExit,omitempty"`
}

type ConnectorMultiCloudGatewayInstance struct {
	Id                int    `json:"id,omitempty"` // sent back on PUT to keep an existing instance
	Name              string `json:"name"`
	BandwidthInMbps   int    `json:"bandwidthInMbps,omitempty"`   // RESPONSE ONLY
	OnRamp            string `json:"onRamp,omitempty"`            // RESPONSE ONLY
	LumenConnectionId string `json:"lumenConnectionId,omitempty"` // RESPONSE ONLY
}

// NewConnectorMultiCloudGateway new Multi-Cloud Gateway connector
func NewConnectorMultiCloudGateway(ac *AlkiraClient) *AlkiraAPI[ConnectorMultiCloudGateway] {
	uri := fmt.Sprintf("%s/tenantnetworks/%s/multi-cloud-gateway-connectors", ac.URI, ac.TenantNetworkId)
	api := &AlkiraAPI[ConnectorMultiCloudGateway]{ac, uri, true}
	return api
}
