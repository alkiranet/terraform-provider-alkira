// Copyright (C) 2022-2025 Alkira Inc. All Rights Reserved.

package alkira

import (
	"encoding/json"
	"fmt"
)

type ConnectorAkamaiProlexic struct {
	AkamaiBgpAsn         int                                           `json:"akamaiBgpAsn"`
	BillingTags          []int                                         `json:"billingTags,omitempty"`
	ByoipOptions         []ConnectorAkamaiProlexicByoipOption          `json:"byoipOptions,omitempty"`
	CXP                  string                                        `json:"cxp"`
	CredentialId         string                                        `json:"credentialId"`
	Enabled              bool                                          `json:"enabled"`
	Group                string                                        `json:"group,omitempty"`
	Id                   json.Number                                   `json:"id,omitempty"`              // response only
	ImplicitGroupId      int                                           `json:"implicitGroupId,omitempty"` // response only
	Name                 string                                        `json:"name"`
	OverlayConfiguration []ConnectorAkamaiProlexicOverlayConfiguration `json:"overlayConfiguration,omitempty"`
	Segments             []string                                      `json:"segments"`
	Size                 string                                        `json:"size"`
	Description          string                                        `json:"description,omitempty"`
}

type ConnectorAkamaiProlexicByoipOption struct {
	ByoipId                   int  `json:"byoipId"`
	RouteAdvertisementEnabled bool `json:"routeAdvertisementEnabled"`
}

type ConnectorAkamaiProlexicOverlayConfiguration struct {
	AlkiraPublicIp string                            `json:"alkiraPublicIp"`
	TunnelIps      []ConnectorAkamaiProlexicTunnelIp `json:"tunnelIps"`
}

type ConnectorAkamaiProlexicTunnelIp struct {
	RanTunnelDestinationIp string `json:"ranTunnelDestinationIp"`
	AlkiraOverlayTunnelIp  string `json:"alkiraOverlayTunnelIp"`
	AkamaiOverlayTunnelIp  string `json:"akamaiOverlayTunnelIp"`
}

func NewConnectorAkamaiProlexic(ac *AlkiraClient) *AlkiraAPI[ConnectorAkamaiProlexic] {
	uri := fmt.Sprintf("%s/tenantnetworks/%s/akamai-prolexic-connectors", ac.URI, ac.TenantNetworkId)
	api := &AlkiraAPI[ConnectorAkamaiProlexic]{ac, uri, true}
	return api
}
