// Copyright (C) 2020-2023 Alkira Inc. All Rights Reserved.

package alkira

import (
	"encoding/json"
	"fmt"
)

type SegmentIpBlocks struct {
	Values []string `json:"values"`
}

type SegmentSrcIpv4PoolList struct {
	StartIp string `json:"startIp,omitempty"`
	EndIp   string `json:"endIp,omitempty"`
}

type Segment struct {
	Asn                                        int                      `json:"asn"`
	Description                                string                   `json:"description,omitempty"`
	EnableIpv6ToIpv4Translation                bool                     `json:"enableIpv6ToIpv4Translation"`
	EnterpriseDNSServerIP                      string                   `json:"enterpriseDNSServerIP,omitempty"`
	Id                                         json.Number              `json:"id,omitempty"` // only for response
	IpBlock                                    string                   `json:"ipBlock,omitempty"`
	IpBlocks                                   SegmentIpBlocks          `json:"ipBlocks,omitempty"`
	Name                                       string                   `json:"name"`
	ReservePublicIPsForUserAndSiteConnectivity bool                     `json:"reservePublicIPsForUserAndSiteConnectivity,omitempty"`
	SrcIpv4PoolList                            []SegmentSrcIpv4PoolList `json:"srcIpv4PoolList,omitempty"`
}

func NewSegment(ac *AlkiraClient) *AlkiraAPI[Segment] {
	uri := fmt.Sprintf("%s/tenantnetworks/%s/segments", ac.URI, ac.TenantNetworkId)
	api := &AlkiraAPI[Segment]{ac, uri}
	return api
}
