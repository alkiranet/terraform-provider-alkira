// Copyright (C) 2024-2025 Alkira Inc. All Rights Reserved.

package alkira

import (
	"fmt"
)

type IPReservation struct {
	Id                string `json:"id,omitempty"`
	Name              string `json:"name"`
	Type              string `json:"type"`
	Prefix            string `json:"prefix"`
	PrefixLen         int    `json:"prefixLen"`
	PrefixType        string `json:"prefixType"`
	FirstIpAssignedTo string `json:"firstIpAssignedTo"`
	NodeId            string `json:"nodeId"`
	Cxp               string `json:"cxp"`
	ScaleGroupId      string `json:"scaleGroupId"`
	Segment           string `json:"segment"`
}

func NewIPReservation(ac *AlkiraClient) *AlkiraAPI[IPReservation] {
	uri := fmt.Sprintf("%s/tenantnetworks/%s/ip-reservations", ac.URI, ac.TenantNetworkId)
	api := &AlkiraAPI[IPReservation]{ac, uri, false}
	return api
}
