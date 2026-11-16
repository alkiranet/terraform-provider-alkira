// Copyright (C) 2026 Alkira Inc. All Rights Reserved.

package alkira

import (
	"encoding/json"
	"fmt"
)

type ConnectorMerakiSDWANSegmentOption struct {
	AdvertiseOnPremRoutes bool  `json:"advertiseOnPremRoutes"`
	DisableInternetExit   bool  `json:"disableInternetExit"`
	MerakiBgpAsn          int64 `json:"merakiBgpAsn"` // up to 4294967295
	SegmentId             int   `json:"segmentId"`
}

type ConnectorMerakiSDWANInstance struct {
	Az           int         `json:"az,omitempty"` // response only
	CredentialId string      `json:"credentialId"`
	HostName     string      `json:"hostName"`
	Id           json.Number `json:"id,omitempty"`
	Name         string      `json:"name,omitempty"`  // response only
	WanIp        string      `json:"wanIp,omitempty"` // response only
}

type ConnectorMerakiSDWAN struct {
	BillingTags     []int                               `json:"billingTags,omitempty"`
	Cxp             string                              `json:"cxp"`
	Description     string                              `json:"description"`
	Enabled         bool                                `json:"enabled"`
	Group           string                              `json:"group,omitempty"`
	GroupId         int                                 `json:"groupId,omitempty"`         // response only
	Id              json.Number                         `json:"id,omitempty"`              // response only
	ImplicitGroupId int                                 `json:"implicitGroupId,omitempty"` // response only
	Instances       []ConnectorMerakiSDWANInstance      `json:"instances"`
	Name            string                              `json:"name"`
	SegmentOptions  []ConnectorMerakiSDWANSegmentOption `json:"segmentOptions"`
	Size            string                              `json:"size"`
	TunnelProtocol  string                              `json:"tunnelProtocol,omitempty"`
	Version         string                              `json:"version,omitempty"`
}

type ConnectorMerakiSDWANBgpTimers struct {
	Hold      int `json:"hold"`
	Keepalive int `json:"keepalive"`
}

type ConnectorMerakiSDWANIke struct {
	DhGroup    int    `json:"dhGroup"`
	Encryption string `json:"encryption"`
	Hash       string `json:"hash"`
	Lifetime   int    `json:"lifetime"`
	Prf        string `json:"prf"`
	Version    string `json:"version"`
}

type ConnectorMerakiSDWANEsp struct {
	Encryption string `json:"encryption"`
	Hash       string `json:"hash"`
	Lifetime   int    `json:"lifetime"`
	PfsGroup   int    `json:"pfsGroup"`
}

type ConnectorMerakiSDWANTunnelConfiguration struct {
	Esp           ConnectorMerakiSDWANEsp `json:"esp"`
	Ike           ConnectorMerakiSDWANIke `json:"ike"`
	IkeId         string                  `json:"ikeId"`
	InsideLocalIp string                  `json:"insideLocalIp"`
	InsidePeerIp  string                  `json:"insidePeerIp"`
	PeerIp        string                  `json:"peerIp"`
	PresharedKey  string                  `json:"presharedKey"`
	Primary       bool                    `json:"primary"`
	TunnelSubnet  string                  `json:"tunnelSubnet"`
}

// ConnectorMerakiSDWANInstanceConfiguration the values the customer enters in
// the Meraki Dashboard for one instance
type ConnectorMerakiSDWANInstanceConfiguration struct {
	AlkiraAsn   int64                                     `json:"alkiraAsn"`
	BgpTimers   ConnectorMerakiSDWANBgpTimers             `json:"bgpTimers"`
	MerakiAsn   int64                                     `json:"merakiAsn"`
	MinFirmware string                                    `json:"minFirmware"`
	Tunnels     []ConnectorMerakiSDWANTunnelConfiguration `json:"tunnels"`
	VmxWanIp    string                                    `json:"vmxWanIp"`
}

// NewConnectorMerakiSDWAN new connector-meraki-sdwan
func NewConnectorMerakiSDWAN(ac *AlkiraClient) *AlkiraAPI[ConnectorMerakiSDWAN] {
	uri := fmt.Sprintf("%s/tenantnetworks/%s/meraki-sdwan-connectors", ac.URI, ac.TenantNetworkId)
	api := &AlkiraAPI[ConnectorMerakiSDWAN]{ac, uri, true}
	return api
}

// GetMerakiSDWANInstanceConfiguration get the Dashboard configuration of a
// connector instance. TPS returns 400 until the instance is deployed. The
// response holds the tunnel pre-shared keys, and the client logs response
// bodies when TF_LOG=DEBUG.
func GetMerakiSDWANInstanceConfiguration(ac *AlkiraClient, connectorId string, instanceId string) (*ConnectorMerakiSDWANInstanceConfiguration, error) {
	uri := fmt.Sprintf("%s/tenantnetworks/%s/meraki-sdwan-connectors/%s/instances/%s/configuration",
		ac.URI, ac.TenantNetworkId, connectorId, instanceId)

	data, _, err := ac.get(uri)
	if err != nil {
		return nil, fmt.Errorf("GetMerakiSDWANInstanceConfiguration: %w", err)
	}

	var result ConnectorMerakiSDWANInstanceConfiguration
	if err := json.Unmarshal(data, &result); err != nil {
		return nil, fmt.Errorf("GetMerakiSDWANInstanceConfiguration: failed to unmarshal: %w", err)
	}

	return &result, nil
}
