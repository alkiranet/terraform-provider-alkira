// Copyright (C) 2021-2025 Alkira Inc. All Rights Reserved.
//
// Implementation of the generic lists with common structure. For
// special lists, it's implementated separately as its own.
package alkira

import (
	"encoding/json"
	"fmt"
)

type List struct {
	Description string      `json:"description"`
	Id          json.Number `json:"id,omitempty"`
	Name        string      `json:"name"`
	Values      []string    `json:"values"`
}

func NewListAsPath(ac *AlkiraClient) *AlkiraAPI[List] {
	uri := fmt.Sprintf("%s/tenantnetworks/%s/as-path-lists", ac.URI, ac.TenantNetworkId)
	api := &AlkiraAPI[List]{ac, uri, true}
	return api
}

func NewListCommunity(ac *AlkiraClient) *AlkiraAPI[List] {
	uri := fmt.Sprintf("%s/tenantnetworks/%s/community-lists", ac.URI, ac.TenantNetworkId)
	api := &AlkiraAPI[List]{ac, uri, true}
	return api
}

func NewListExtendedCommunity(ac *AlkiraClient) *AlkiraAPI[List] {
	uri := fmt.Sprintf("%s/tenantnetworks/%s/extended-community-lists", ac.URI, ac.TenantNetworkId)
	api := &AlkiraAPI[List]{ac, uri, true}
	return api
}
