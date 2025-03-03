// Copyright (C) 2020-2025 Alkira Inc. All Rights Reserved.

package alkira

import (
	"encoding/json"
	"fmt"
)

type InternetApplicationSnatIpv4 struct {
	StartIp string `json:"startIp"`
	EndIp   string `json:"endIp"`
}

type InternetApplicationTargets struct {
	Type       string   `json:"type"`
	Value      string   `json:"value"`
	Ports      []int    `json:"ports"`
	PortRanges []string `json:"portRanges,omitempty"`
}

type InternetApplication struct {
	BillingTags                   []int                          `json:"billingTags"`
	BiDirectionalAvailabilityZone string                         `json:"biDirectionalAvailabilityZone,omitempty"`
	ByoipId                       int                            `json:"byoipId,omitempty"`
	ConnectorId                   int                            `json:"connectorId"`
	ConnectorType                 string                         `json:"connectorType"`
	Description                   string                         `json:"description,omitempty"`
	FqdnPrefix                    string                         `json:"fqdnPrefix"`
	Group                         string                         `json:"group,omitempty"`
	Id                            json.Number                    `json:"id,omitempty"` // response only
	InboundConnectorId            string                         `json:"inboundConnectorId,omitempty"`
	InboundConnectorType          string                         `json:"inboundConnectorType,omitempty"`
	InboundInternetGroupId        json.Number                    `json:"inboundInternetGroupId,omitempty"`
	InternetProtocol              string                         `json:"internetProtocol"`
	Name                          string                         `json:"name"`
	PublicIps                     []string                       `json:"publicIps"`
	SegmentName                   string                         `json:"segmentName"`
	SnatIpv4Ranges                []*InternetApplicationSnatIpv4 `json:"snatIpv4Ranges,omitempty"`
	Size                          string                         `json:"size"`
	Targets                       []InternetApplicationTargets   `json:"targets,omitempty"`
	IlbCredentialId               string                         `json:"ilbCredentialId,omitempty"`
}

// NewInternetApplication new internet application
func NewInternetApplication(ac *AlkiraClient) *AlkiraAPI[InternetApplication] {
	uri := fmt.Sprintf("%s/v1/tenantnetworks/%s/internet-applications", ac.URI, ac.TenantNetworkId)
	api := &AlkiraAPI[InternetApplication]{ac, uri, true}
	return api
}
