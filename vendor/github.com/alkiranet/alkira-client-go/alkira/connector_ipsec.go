// Copyright (C) 2020-2022 Alkira Inc. All Rights Reserved.

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

//
// From the current version of API, the routing option is enforced and
// you have to set it.
//
// routingOptions: {
//   staticRouting: {}
//   dynamicRouting: {}
// }
//
type ConnectorIPSecStaticRouting struct {
	Availability string `json:"availability"`
	PrefixListId int    `json:"prefixListId"`
}

type ConnectorIPSecDynamicRouting struct {
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
}

// getIpsecConnectors get all ipsec connectors from the given tenant network
func (ac *AlkiraClient) getIpsecConnectors() (string, error) {
	uri := fmt.Sprintf("%s/tenantnetworks/%s/ipsecconnectors", ac.URI, ac.TenantNetworkId)

	data, err := ac.get(uri)
	return string(data), err
}

// CreateConnectorIPSec create an IPSEC connector
func (ac *AlkiraClient) CreateConnectorIPSec(connector *ConnectorIPSec) (string, error) {
	uri := fmt.Sprintf("%s/v1/tenantnetworks/%s/ipsecconnectors", ac.URI, ac.TenantNetworkId)

	// Construct the request
	body, err := json.Marshal(connector)

	if err != nil {
		return "", fmt.Errorf("CreateConnectorIPSec: failed to marshal: %v", err)
	}

	data, err := ac.create(uri, body, true)

	if err != nil {
		return "", err
	}

	var result ConnectorIPSec
	err = json.Unmarshal([]byte(data), &result)

	if err != nil {
		return "", fmt.Errorf("CreateConnectorAwsVpc: failed to unmarshal: %v", err)
	}

	return string(result.Id), nil
}

// DeleteConnectorIPSec delete an IPSEC connector by Id
func (ac *AlkiraClient) DeleteConnectorIPSec(id string) error {
	uri := fmt.Sprintf("%s/v1/tenantnetworks/%s/ipsecconnectors/%s", ac.URI, ac.TenantNetworkId, id)

	return ac.delete(uri, true)
}

// UpdateConnectorIPSec update an IPSEC connector by Id
func (ac *AlkiraClient) UpdateConnectorIPSec(id string, connector *ConnectorIPSec) error {
	uri := fmt.Sprintf("%s/v1/tenantnetworks/%s/ipsecconnectors/%s", ac.URI, ac.TenantNetworkId, id)

	body, err := json.Marshal(connector)

	if err != nil {
		return fmt.Errorf("UpdateConnectorIPSec: failed to marshal request: %v", err)
	}

	return ac.update(uri, body, true)
}

// GetConnectorIPSec get an IPSEC connector by Id
func (ac *AlkiraClient) GetConnectorIPSec(id string) (*ConnectorIPSec, error) {
	uri := fmt.Sprintf("%s/tenantnetworks/%s/ipsecconnectors/%s", ac.URI, ac.TenantNetworkId, id)

	var connector ConnectorIPSec

	data, err := ac.get(uri)

	if err != nil {
		return nil, err
	}

	err = json.Unmarshal([]byte(data), &connector)

	if err != nil {
		return nil, fmt.Errorf("GetConnectorIPSec: failed to unmarshal: %v", err)
	}

	return &connector, nil
}

// GetConnectorIpsecByName get an internet connector by name
func (ac *AlkiraClient) GetConnectorIpsecByName(name string) (ConnectorIPSec, error) {
	var ipsecConnector ConnectorIPSec

	if len(name) == 0 {
		return ipsecConnector, fmt.Errorf("GetConnectorIpsecByName: Invalid Connector name")
	}

	ipsecConnectors, err := ac.getIpsecConnectors()

	if err != nil {
		return ipsecConnector, err
	}

	var result []ConnectorIPSec
	json.Unmarshal([]byte(ipsecConnectors), &result)

	for _, l := range result {
		if l.Name == name {
			return l, nil
		}
	}

	return ipsecConnector, fmt.Errorf("GetConnectorIpsecByName: failed to find the connector by %s", name)
}
