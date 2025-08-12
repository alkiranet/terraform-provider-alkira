// Copyright (C) 2023-2025 Alkira Inc. All Rights Reserved.

package alkira

import (
	"encoding/json"
	"fmt"
)

type ConnectorJuniperSsrVrfMapping struct {
	AdvertiseOnPremRoutes bool   `json:"advertiseOnPremRoutes"`
	DisableInternetExit   bool   `json:"disableInternetExit"`
	JuniperSsrBgpAsn      int    `json:"juniperSsrBgpAsn"`
	SegmentId             int    `json:"segmentId"`
	JuniperSsrVrfName     string `json:"juniperSsrVrfName"`
}

type ConnectorJuniperSdwanInstance struct {
	CredentialId                string `json:"credentialId"`
	RegistrationKeyCredentialId string `json:"registrationKeyCredentialId"`
	HostName                    string `json:"hostName"`
	UserName                    string `json:"userName"`
	Id                          int    `json:"id,omitempty"`
}

type ConnectorJuniperSdwan struct {
	BillingTags           []int                           `json:"billingTags,omitempty"`
	Cxp                   string                          `json:"cxp"`
	Enabled               bool                            `json:"enabled"`
	JuniperSsrVrfMappings []ConnectorJuniperSsrVrfMapping `json:"juniperSsrVrfMappings"`
	Group                 string                          `json:"group,omitempty"`
	Id                    json.Number                     `json:"id,omitempty"`              // response only
	ImplicitGroupId       int                             `json:"implicitGroupId,omitempty"` // response only
	Instances             []ConnectorJuniperSdwanInstance `json:"instances"`
	Name                  string                          `json:"name"`
	Size                  string                          `json:"size"`
	TunnelProtocol        string                          `json:"tunnelProtocol"`
	Version               string                          `json:"version"`
	Description           string                          `json:"description,omitempty"`
}

// NewConnectorJuniperSdwan
func NewConnectorJuniperSdwan(ac *AlkiraClient) *AlkiraAPI[ConnectorJuniperSdwan] {
	uri := fmt.Sprintf("%s/tenantnetworks/%s/juniper-sdwan-connectors", ac.URI, ac.TenantNetworkId)
	api := &AlkiraAPI[ConnectorJuniperSdwan]{ac, uri, true}
	return api
}
