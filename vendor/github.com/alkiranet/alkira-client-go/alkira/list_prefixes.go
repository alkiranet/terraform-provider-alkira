// Copyright (C) 2020-2025 Alkira Inc. All Rights Reserved.

package alkira

import (
	"encoding/json"
	"fmt"
)

type PolicyPrefixListRange struct {
	Prefix      string `json:"prefix"`
	Le          int    `json:"le,omitempty"`
	Ge          int    `json:"ge,omitempty"`
	Description string `json:"description,omitempty"`
}

type PolicyPrefixListDetails struct {
	Description string `json:"description,omitempty"`
}

type PolicyPrefixList struct {
	Description   string                              `json:"description"`
	Id            json.Number                         `json:"id,omitempty"`
	Name          string                              `json:"name"`
	Prefixes      []string                            `json:"prefixes"`
	PrefixDetails map[string]*PolicyPrefixListDetails `json:"prefixDetails,omitempty"`
	PrefixRanges  []PolicyPrefixListRange             `json:"prefixRanges,omitempty"`
	Type          string                              `json:"type,omitempty"`
}

// NewPolicyPrefixList new policy prefix list
func NewPolicyPrefixList(ac *AlkiraClient) *AlkiraAPI[PolicyPrefixList] {
	uri := fmt.Sprintf("%s/tenantnetworks/%s/policy/prefixlists", ac.URI, ac.TenantNetworkId)
	api := &AlkiraAPI[PolicyPrefixList]{ac, uri, true}
	return api
}
