// Copyright (C) 2024-2025 Alkira Inc. All Rights Reserved.

package alkira

import (
	"encoding/json"
	"fmt"
)

type UdrListUdrs struct {
	Prefix       string `json:"prefix"`
	NextHopType  string `json:"nextHopType"`
	NextHopValue string `json:"nextHopValue"`
	Description  string `json:"description,omitempty"`
}

type UdrList struct {
	Name          string        `json:"name"`
	Description   string        `json:"description,omitempty"`
	CloudProvider string        `json:"cloudProvider"`
	Id            json.Number   `json:"id,omitempty"`
	Udrs          []UdrListUdrs `json:"udrs"`
}

// NewUdrList new UDR list
func NewUdrList(ac *AlkiraClient) *AlkiraAPI[UdrList] {
	uri := fmt.Sprintf("%s/tenantnetworks/%s/udr-lists", ac.URI, ac.TenantNetworkId)
	api := &AlkiraAPI[UdrList]{ac, uri, true}
	return api
}
