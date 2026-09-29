// Copyright (C) 2023-2025 Alkira Inc. All Rights Reserved.

package alkira

import (
	"encoding/json"
	"fmt"
)

type ConnectorPrismaSDWANVRFMapping struct {
	AdvertiseOnPremRoutes bool   `json:"advertiseOnPremRoutes"`
	DisableInternetExit   bool   `json:"disableInternetExit"`
	GatewayBgpAsn         int    `json:"gatewayBgpAsn,omitempty"`
	SegmentId             int    `json:"segmentId"`
	VrfName               string `json:"vrfName"`
}

type ConnectorPrismaSDWANInstance struct {
	CredentialId string      `json:"credentialId"`
	HostName     string      `json:"hostName"`
	Id           json.Number `json:"id,omitempty"`
	IonModel     string      `json:"ionModel"`
	Version      string      `json:"version"`
}

type ConnectorPrismaSDWAN struct {
	AllowList               []string                         `json:"allowList,omitempty"`
	BillingTags             []int                            `json:"billingTags,omitempty"`
	Cxp                     string                           `json:"cxp"`
	Description             string                           `json:"description,omitempty"`
	Enabled                 bool                             `json:"enabled"`
	Group                   string                           `json:"group,omitempty"`
	GroupId                 int                              `json:"groupId,omitempty"`         // response only
	Id                      json.Number                      `json:"id,omitempty"`              // response only
	ImplicitGroupId         int                              `json:"implicitGroupId,omitempty"` // response only
	Instances               []ConnectorPrismaSDWANInstance   `json:"instances"`
	Name                    string                           `json:"name"`
	PrismaSDWANVRFMappings  []ConnectorPrismaSDWANVRFMapping `json:"prismaSDWANVRFMappings"`
	ScaleGroupId            string                           `json:"scaleGroupId,omitempty"`
	Size                    string                           `json:"size"`
	TunnelProtocol          string                           `json:"tunnelProtocol"`
}

// NewConnectorPrismaSDWAN new connector-prisma-sdwan
func NewConnectorPrismaSDWAN(ac *AlkiraClient) *AlkiraAPI[ConnectorPrismaSDWAN] {
	uri := fmt.Sprintf("%s/tenantnetworks/%s/prisma-sdwan-connectors", ac.URI, ac.TenantNetworkId)
	api := &AlkiraAPI[ConnectorPrismaSDWAN]{ac, uri, true}
	return api
}
