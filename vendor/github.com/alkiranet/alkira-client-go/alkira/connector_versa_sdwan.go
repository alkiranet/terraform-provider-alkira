// Copyright (C) 2023-2025 Alkira Inc. All Rights Reserved.

package alkira

import (
	"encoding/json"
	"fmt"
)

type VersaSdwanVrfMapping struct {
	AdvertiseOnPremRoutes bool   `json:"advertiseOnPremRoutes"`
	GatewayBgpAsn         int    `json:"gatewayBgpAsn,omitempty"`
	DisableInternetExit   bool   `json:"disableInternetExit"`
	SegmentId             int    `json:"segmentId"`
	VrfName               string `json:"vrfName"`
}

type VersaSdwanInstance struct {
	HostName     string `json:"hostName"`
	Id           int    `json:"id,omitempty"`
	SerialNumber string `json:"serialNumber,omitempty"`
	Version      string `json:"version"`
}

type ConnectorVersaSdwan struct {
	BillingTags           []int                  `json:"billingTags"`
	Cxp                   string                 `json:"cxp"`
	Group                 string                 `json:"group,omitempty"`
	Id                    json.Number            `json:"id,omitempty"`              // response only
	ImplicitGroupId       int                    `json:"implicitGroupId,omitempty"` // response only
	Instances             []VersaSdwanInstance   `json:"instances"`
	Name                  string                 `json:"name"`
	GlobalTenantId        int                    `json:"globalTenantId"`
	LocalId               string                 `json:"localId"`
	LocalPublicSharedKey  string                 `json:"localPublicSharedKey"`
	RemoteId              string                 `json:"remoteId"`
	RemotePublicSharedKey string                 `json:"remotePublicSharedKey"`
	Size                  string                 `json:"size"`
	TunnelProtocol        string                 `json:"tunnelProtocol"`
	VersaControllerHost   string                 `json:"versaControllerHost"`
	VersaSdWanVRFMappings []VersaSdwanVrfMapping `json:"versaSDWANVRFMappings"`
	Enabled               bool                   `json:"enabled"`
	Description           string                 `json:"description,omitempty"`
}

// NewConnectorVersaSdwan new connector
func NewConnectorVersaSdwan(ac *AlkiraClient) *AlkiraAPI[ConnectorVersaSdwan] {
	uri := fmt.Sprintf("%s/tenantnetworks/%s/versa-sdwan-connectors", ac.URI, ac.TenantNetworkId)
	api := &AlkiraAPI[ConnectorVersaSdwan]{ac, uri, true}
	return api
}
