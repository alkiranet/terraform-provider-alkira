// Copyright (C) 2020-2025 Alkira Inc. All Rights Reserved.

package alkira

import (
	"encoding/json"
	"fmt"
)

type ConnectorIPSecSiteAdvanced struct {
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

type ConnectorIPSecSite struct {
	Advanced               *ConnectorIPSecSiteAdvanced `json:"advanced,omitempty"`
	BillingTags            []int                       `json:"billingTags,omitempty"`
	CustomerGwAsn          string                      `json:"customerGwAsn,omitempty"`
	CustomerGwIp           string                      `json:"customerGwIp"`
	EnableTunnelRedundancy bool                        `json:"enableTunnelRedundancy"`
	GatewayIpType          string                      `json:"gatewayIpType,omitempty"`
	HaMode                 string                      `json:"haMode,omitempty"`
	Id                     int                         `json:"id,omitempty"` // response only
	Name                   string                      `json:"name"`
	PresharedKeys          []string                    `json:"presharedKeys"`
}

type ConnectorIPSecPolicyOptions struct {
	BranchTSPrefixListIds []int `json:"branchTSPrefixListIds"`
	CxpTSPrefixListIds    []int `json:"cxpTSPrefixListIds"`
}

// From the current version of API, the routing option is enforced and
// you have to set it.
//
//	routingOptions: {
//	  staticRouting: {}
//	  dynamicRouting: {}
//	}
type ConnectorIPSecStaticRouting struct {
	Availability string `json:"availability"`
	PrefixListId int    `json:"prefixListId"`
}

type ConnectorIPSecDynamicRouting struct {
	Availability     string `json:"availability,omitempty"`
	BgpAuthKeyAlkira string `json:"bgpAuthKeyAlkira"`
	CustomerGwAsn    string `json:"customerGwAsn"`
}

type ConnectorIPSecRoutingOptions struct {
	DynamicRouting *ConnectorIPSecDynamicRouting `json:"dynamicRouting"`
	StaticRouting  *ConnectorIPSecStaticRouting  `json:"staticRouting"`
}

// SegmentOptions block is dynamic and so can't be put into a
// structure. It needs to be marshalled from
// map[string]ConnectorIPSecSegmentOptions
type ConnectorIPSecSegmentOptions struct {
	AdvertiseOnPremRoutes *bool `json:"advertiseOnPremRoutes,omitempty"`
	DisableInternetExit   *bool `json:"disableInternetExit,omitempty"`
}

type ConnectorIPSec struct {
	BillingTags     []int                         `json:"billingTags"`
	CXP             string                        `json:"cxp"`
	Group           string                        `json:"group,omitempty"`
	Enabled         bool                          `json:"enabled"`
	Id              json.Number                   `json:"id,omitempty"`              // response only
	ImplicitGroupId int                           `json:"implicitGroupId,omitempty"` // response only
	Name            string                        `json:"name"`
	PolicyOptions   *ConnectorIPSecPolicyOptions  `json:"policyOptions"`
	RoutingOptions  *ConnectorIPSecRoutingOptions `json:"routingOptions"`
	SegmentOptions  interface{}                   `json:"segmentOptions"`
	Segments        []string                      `json:"segments"` // Only one segment is supported for now
	Sites           []*ConnectorIPSecSite         `json:"sites,omitempty"`
	Size            string                        `json:"size"`
	VpnMode         string                        `json:"vpnMode"`
	ScaleGroupId    string                        `json:"scaleGroupId,omitempty"`
	Description     string                        `json:"description,omitempty"`
}

// NewConnectorIPSec initialize a new connector
func NewConnectorIPSec(ac *AlkiraClient) *AlkiraAPI[ConnectorIPSec] {
	uri := fmt.Sprintf("%s/v1/tenantnetworks/%s/ipsecconnectors", ac.URI, ac.TenantNetworkId)
	api := &AlkiraAPI[ConnectorIPSec]{ac, uri, true}
	return api
}
