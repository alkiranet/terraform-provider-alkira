// Copyright (C) 2023-2025 Alkira Inc. All Rights Reserved.

package alkira

import (
	"encoding/json"
	"fmt"
)

type ConnectorFortinetSdwanVrfMapping struct {
	AdvertiseOnPremRoutes bool `json:"advertiseOnPremRoutes"`
	DisableInternetExit   bool `json:"disableInternetExit"`
	GatewayBgpAsn         int  `json:"gatewayBgpAsn,omitempty"`
	SegmentId             int  `json:"segmentId"`
	Vrf                   int  `json:"vrf"`
}

type ConnectorFortinetSdwanInstance struct {
	CredentialId string `json:"credentialId"`
	HostName     string `json:"hostName"`
	Id           int    `json:"id,omitempty"`
	LicenseType  string `json:"licenseType"`
	SerialNumber string `json:"serialNumber,omitempty"`
	Version      string `json:"version"`
}

type ConnectorFortinetSdwan struct {
	AllowList            []string                           `json:"allowList,omitempty"`
	BillingTags          []int                              `json:"billingTags,omitempty"`
	Cxp                  string                             `json:"cxp"`
	Enabled              bool                               `json:"enabled"`
	FtntSdWanVRFMappings []ConnectorFortinetSdwanVrfMapping `json:"ftntSDWANVRFMappings"`
	Group                string                             `json:"group,omitempty"`
	Id                   json.Number                        `json:"id,omitempty"`              // response only
	ImplicitGroupId      int                                `json:"implicitGroupId,omitempty"` // response only
	Instances            []ConnectorFortinetSdwanInstance   `json:"instances"`
	Name                 string                             `json:"name"`
	Size                 string                             `json:"size"`
	TunnelProtocol       string                             `json:"tunnelProtocol"`
	Description          string                             `json:"description,omitempty"`
}

// NewConnectorFortinetSdwan new connector-fortinet-sdwan
func NewConnectorFortinetSdwan(ac *AlkiraClient) *AlkiraAPI[ConnectorFortinetSdwan] {
	uri := fmt.Sprintf("%s/tenantnetworks/%s/ftnt-sdwan-connectors", ac.URI, ac.TenantNetworkId)
	api := &AlkiraAPI[ConnectorFortinetSdwan]{ac, uri, true}
	return api
}
