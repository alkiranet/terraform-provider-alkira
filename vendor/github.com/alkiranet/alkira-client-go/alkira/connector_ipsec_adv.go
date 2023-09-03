// Copyright (C) 2023 Alkira Inc. All Rights Reserved.

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
	OverlayIpReservationId string `json:"OverlayIpReservationId"`
}

type ConnectorAdvIPSecTunnelCxpEnd struct {
	PublicIpReservationId  string `json:"PublicIpReservationId"`
	OverlayIpReservationId string `json:"OverlayIpReservationId"`
}

type ConnectorAdvIPSecTunnel struct {
	Id           int                                `json:"id,omitempty"`
	TunnelNo     int                                `json:"tunnelNo"`
	PresharedKey string                             `json:"preSharedKey"`
	ProfileId    string                             `json:"profileId"`
	CustomerEnd  ConnectorAdvIPSecTunnelCustomerEnd `json:"customerEnd"`
	CxpEnd       ConnectorAdvIPSecTunnelCxpEnd      `json:"cxpEnd"`

	Advanced *ConnectorAdvIPSecAdvanced `json:"advanced,omitempty"`
}
type ConnectorAdvIPSecTunnelGateway struct {
	Id           int                        `json:"id,omitempty"`
	Name         string                     `json:"name"`
	CustomerGwIp string                     `json:"customerGwIp"`
	HaMode       string                     `json:"haMode,omitempty"`
	Tunnels      []*ConnectorAdvIPSecTunnel `json:"tunnels,omitempty"`
}

// gateway
type ConnectorAdvIPSecGateway struct {
	Advanced               *ConnectorAdvIPSecAdvanced `json:"advanced,omitempty"`
	BillingTags            []int                      `json:"billingTags,omitempty"`
	CustomerGwAsn          string                     `json:"customerGwAsn,omitempty"`
	CustomerGwIp           string                     `json:"customerGwIp"`
	EnableTunnelRedundancy bool                       `json:"enableTunnelRedundancy"`
	GatewayIpType          string                     `json:"gatewayIpType,omitempty"`
	HaMode                 string                     `json:"haMode,omitempty"`
	Id                     int                        `json:"id,omitempty"` // response only
	Name                   string                     `json:"name"`
	PresharedKeys          []string                   `json:"presharedKeys"`
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
	DynamicRouting *ConnectorIPSecDynamicRouting `json:"dynamicRouting"`
	StaticRouting  *ConnectorIPSecStaticRouting  `json:"staticRouting"`
}

// Top Level
type ConnectorAdvIPSec struct {
	AdvertiseDefaultRoute *bool                            `json:"advertiseDefaultRoute,omitempty"`
	AdvertiseOnPremRoutes *bool                            `json:"advertiseOnPremRoutes,omitempty"`
	BillingTags           []int                            `json:"billingTags"`
	CXP                   string                           `json:"cxp"`
	Enabled               bool                             `json:"enabled"`
	DestinationType       string                           `json:"destinationType"`
	Gateways              []ConnectorAdvIPSecGateway       `json:"gateways"`
	Group                 string                           `json:"group,omitempty"`
	Id                    json.Number                      `json:"id,omitempty"`              // response only
	ImplicitGroupId       int                              `json:"implicitGroupId,omitempty"` // response only
	Name                  string                           `json:"name"`
	PolicyOptions         *ConnectorAdvIPSecPolicyOptions  `json:"policyOptions"`
	RoutingOptions        *ConnectorAdvIPSecRoutingOptions `json:"routingOptions"`
	Segment               string                           `json:"segment"`
	TunnelsPerGateway     int                              `json:"tunnelsPerGateway"`
	Size                  string                           `json:"size"`
	VpnMode               string                           `json:"vpnMode"`
}

// NewConnectorAdvIPSec initialize a new connector
func NewConnectorAdvIPSec(ac *AlkiraClient) *AlkiraAPI[ConnectorAdvIPSec] {
	uri := fmt.Sprintf("%s/v1/tenantnetworks/%s/adv-ipsec-connectors", ac.URI, ac.TenantNetworkId)
	api := &AlkiraAPI[ConnectorAdvIPSec]{ac, uri, true}
	return api
}
