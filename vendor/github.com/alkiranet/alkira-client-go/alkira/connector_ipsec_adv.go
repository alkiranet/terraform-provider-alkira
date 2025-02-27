// Copyright (C) 2023-2025 Alkira Inc. All Rights Reserved.

package alkira

import (
	"encoding/json"
	"fmt"
)

type ConnectorAdvIPSecAdvanced struct {
	DPDDelay                int      `json:"dPDDelay,omitempty"`
	DPDTimeout              int      `json:"dPDTimeout,omitempty"`
	EspDHGroupNumbers       []string `json:"espDHGroupNumbers,omitempty"`
	EspEncryptionAlgorithms []string `json:"espEncryptionAlgorithms,omitempty"`
	EspIntegrityAlgorithms  []string `json:"espIntegrityAlgorithms,omitempty"`
	EspLifeTime             int      `json:"espLifeTime,omitempty"`
	EspRandomTime           int      `json:"espRandomTime,omitempty"`
	EspRekeyTime            int      `json:"espRekeyTime,omitempty"`
	IkeDHGroupNumbers       []string `json:"ikeDHGroupNumbers,omitempty"`
	IkeEncryptionAlgorithms []string `json:"ikeEncryptionAlgorithms,omitempty"`
	IkeIntegrityAlgorithms  []string `json:"ikeIntegrityAlgorithms,omitempty"`
	IkeOverTime             int      `json:"ikeOverTime,omitempty"`
	IkeRandomTime           int      `json:"ikeRandomTime,omitempty"`
	IkeRekeyTime            int      `json:"ikeRekeyTime,omitempty"`
	IkeVersion              string   `json:"ikeVersion,omitempty"`
	Initiator               bool     `json:"initiator,omitempty"`
	LocalAuthType           string   `json:"localAuthType,omitempty"`
	LocalAuthValue          string   `json:"localAuthValue,omitempty"`
	RemoteAuthType          string   `json:"remoteAuthType,omitempty"`
	RemoteAuthValue         string   `json:"remoteAuthValue,omitempty"`
	ReplayWindowSize        int      `json:"replayWindowSize,omitempty"`
}

type ConnectorAdvIPSecTunnelCustomerEnd struct {
	OverlayIp              string `json:"overlayIp,omitempty"`
	OverlayIpReservationId string `json:"overlayIpReservationId"`
}

type ConnectorAdvIPSecTunnelCxpEnd struct {
	OverlayIpReservationId string `json:"overlayIpReservationId"`
	PublicIpReservationId  string `json:"publicIpReservationId"`
}

type ConnectorAdvIPSecTunnel struct {
	Advanced     *ConnectorAdvIPSecAdvanced         `json:"advanced,omitempty"`
	CustomerEnd  ConnectorAdvIPSecTunnelCustomerEnd `json:"customerEnd"`
	CxpEnd       ConnectorAdvIPSecTunnelCxpEnd      `json:"cxpEnd"`
	Id           string                             `json:"id,omitempty"`
	PresharedKey string                             `json:"preSharedKey"`
	ProfileId    int                                `json:"profileId,omitempty"`
	TunnelNo     int                                `json:"tunnelNo"`
}

// Gateway
type ConnectorAdvIPSecGateway struct {
	CustomerGwIp string                     `json:"customerGwIp"`
	HaMode       string                     `json:"haMode,omitempty"`
	Id           int                        `json:"id,omitempty"`
	Name         string                     `json:"name"`
	Tunnels      []*ConnectorAdvIPSecTunnel `json:"tunnels,omitempty"`
}

// Policy Options
type ConnectorAdvIPSecPolicyOptions struct {
	BranchTSPrefixListIds []int `json:"branchTSPrefixListIds"`
	CxpTSPrefixListIds    []int `json:"cxpTSPrefixListIds"`
}

// Routing Options
type ConnectorAdvIPSecStaticRouting struct {
	Availability string `json:"availability"`
	PrefixListId int    `json:"prefixListId"`
}

type ConnectorAdvIPSecDynamicRouting struct {
	Availability     string `json:"availability,omitempty"`
	BgpAuthKeyAlkira string `json:"bgpAuthKeyAlkira"`
	CustomerGwAsn    string `json:"customerGwAsn"`
}

type ConnectorAdvIPSecRoutingOptions struct {
	DynamicRouting *ConnectorAdvIPSecDynamicRouting `json:"dynamicRouting"`
	StaticRouting  *ConnectorAdvIPSecStaticRouting  `json:"staticRouting"`
}

// Top Level
type ConnectorAdvIPSec struct {
	AdvertiseDefaultRoute bool                             `json:"advertiseDefaultRoute"`
	AdvertiseOnPremRoutes bool                             `json:"advertiseOnPremRoutes"`
	BillingTags           []int                            `json:"billingTags"`
	CXP                   string                           `json:"cxp"`
	Enabled               bool                             `json:"enabled"`
	DestinationType       string                           `json:"destinationType"`
	Gateways              []*ConnectorAdvIPSecGateway      `json:"gateways"`
	Group                 string                           `json:"group,omitempty"`
	Id                    json.Number                      `json:"id,omitempty"`
	ImplicitGroupId       int                              `json:"implicitGroupId,omitempty"`
	Name                  string                           `json:"name"`
	PolicyOptions         *ConnectorAdvIPSecPolicyOptions  `json:"policyOptions"`
	RoutingOptions        *ConnectorAdvIPSecRoutingOptions `json:"routingOptions"`
	Segment               string                           `json:"segment"`
	TunnelsPerGateway     int                              `json:"tunnelsPerGateway"`
	Size                  string                           `json:"size"`
	VpnMode               string                           `json:"vpnMode"`
	Description           string                           `json:"description,omitempty"`
}

// NewConnectorAdvIPSec initialize a new connector
func NewConnectorAdvIPSec(ac *AlkiraClient) *AlkiraAPI[ConnectorAdvIPSec] {
	uri := fmt.Sprintf("%s/v1/tenantnetworks/%s/adv-ipsec-connectors", ac.URI, ac.TenantNetworkId)
	api := &AlkiraAPI[ConnectorAdvIPSec]{ac, uri, true}
	return api
}
