// Copyright (C) 2024 Alkira Inc. All Rights Reserved.

package alkira

import (
	"encoding/json"
	"fmt"
)

type ConnectorAwsDirectConnectSegmentOption struct {
	SegmentName                      string  `json:"segmentName"`
	CustomerAsn                      int     `json:"customerAsn"`
	CustomerLoopbackIp               string  `json:"customerLoopbackIp,omitempty"`
	AlkLoopbackIp1                   string  `json:"alkLoopbackIp1,omitempty"`
	AlkLoopbackIp2                   string  `json:"alkLoopbackIp2,omitempty"`
	LoopbackSubnet                   string  `json:"loopbackSubnet"`
	AdvertiseOnPremRoutes            bool `json:"advertiseOnPremRoutes"`
	DisableInternetExit              bool `json:"disableInternetExit"`
	NumOfCustomerLoopbackIps         int     `json:"numOfCustomerLoopbackIps"`
	TunnelCountPerCustomerLoopbackIp int     `json:"tunnelCountPerCustomerLoopbackIp,omitempty"`
}

type ConnectorAwsDirectConnectInstance struct {
	Name              string           `json:"name"`
	ConnectionId      string           `json:"connectionId"`
	DcGatewayAsn      int              `json:"dcGatewayAsn"`
	UnderlayAsn       int              `json:"underlayAsn"`
	UnderlayPrefix    string           `json:"underlayPrefix"`
	AwsUnderlayIp     string           `json:"awsUnderlayIp,omitempty"`
	OnPremUnderlayIp  string           `json:"onPremUnderlayIp,omitempty"`
	BgpAuthKey        string           `json:"bgpAuthKey,omitempty"`
	BgpAuthKeyAlkira  string           `json:"bgpAuthKeyAlkira,omitempty"`
	Vlan              int              `json:"vlan"`
	CustomerRegion    string           `json:"customerRegion"`
	CredentialId      string           `json:"credentialId"`
	GatewayMacAddress string           `json:"gatewayMacAddress,omitempty"`
	Vni               int              `json:"vni"`
	SegmentOptions    []ConnectorAwsDirectConnectSegmentOption `json:"segmentOptions"`
}

type ConnectorAwsDirectConnect struct {
	BillingTags    []int       `json:"billingTags"`
	Cxp            string      `json:"cxp"`
	Description    string      `json:"description,omitempty"`
	Enabled        bool        `json:"enabled"`
	Group          string      `json:"group,omitempty"`
	Id             json.Number `json:"id,omitempty"` // response only
	Name           string      `json:"name"`
	Size           string      `json:"size"`
	TunnelProtocol string      `json:"tunnelProtocol"`
	Instances      []ConnectorAwsDirectConnectInstance `json:"instances"`
}

// NewConnectorAwsDirectConnect new connector-aws-direct-connect
func NewConnectorAwsDirectConnect(ac *AlkiraClient) *AlkiraAPI[ConnectorAwsDirectConnect] {
	uri := fmt.Sprintf("%s/tenantnetworks/%s/directconnectconnectors", ac.URI, ac.TenantNetworkId)
	api := &AlkiraAPI[ConnectorAwsDirectConnect]{ac, uri, true}
	return api
}
