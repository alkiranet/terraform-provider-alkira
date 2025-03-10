// Copyright (C) 2023-2025 Alkira Inc. All Rights Reserved.

package alkira

import (
	"encoding/json"
	"fmt"
)

type VmwareSdwanVrfMapping struct {
	AdvertiseOnPremRoutes  bool   `json:"advertiseOnPremRoutes"`
	GatewayBgpAsn          int    `json:"gatewayBgpAsn,omitempty"`
	DisableInternetExit    bool   `json:"disableInternetExit"`
	SegmentId              int    `json:"segmentId"`
	VmWareSdWanSegmentName string `json:"vmWareSdWanSegmentName"`
}

type VmwareSdwanInstance struct {
	CredentialId string `json:"credentialId"`
	HostName     string `json:"hostName"`
	Id           int    `json:"id,omitempty"`
	Name         string `json:"name,omitempty"`
}

type ConnectorVmwareSdwan struct {
	BillingTags             []int                   `json:"billingTags"`
	Cxp                     string                  `json:"cxp"`
	Group                   string                  `json:"group,omitempty"`
	Id                      json.Number             `json:"id,omitempty"`              // response only
	ImplicitGroupId         int                     `json:"implicitGroupId,omitempty"` // response only
	Instances               []VmwareSdwanInstance   `json:"instances"`
	Name                    string                  `json:"name"`
	OrchestratorHostAddress string                  `json:"orchestratorHostAddress"`
	Size                    string                  `json:"size"`
	TunnelProtocol          string                  `json:"tunnelProtocol"`
	Version                 string                  `json:"version"`
	VmWareSdWanVRFMappings  []VmwareSdwanVrfMapping `json:"vmWareSdWanVRFMappings"`
	Enabled                 bool                    `json:"enabled"`
	Description             string                  `json:"description,omitempty"`
}

// NewConnectorVmwareSdwan new connector
func NewConnectorVmwareSdwan(ac *AlkiraClient) *AlkiraAPI[ConnectorVmwareSdwan] {
	uri := fmt.Sprintf("%s/tenantnetworks/%s/vmware-sdwan-connectors", ac.URI, ac.TenantNetworkId)
	api := &AlkiraAPI[ConnectorVmwareSdwan]{ac, uri, true}
	return api
}
