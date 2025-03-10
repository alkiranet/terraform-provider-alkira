// Copyright (C) 2021-2025 Alkira Inc. All Rights Reserved.

package alkira

import (
	"encoding/json"
	"fmt"
)

type NatPolicy struct {
	Name                               string      `json:"name"`
	Description                        string      `json:"description"`
	Type                               string      `json:"type"`
	Segment                            string      `json:"segment"`
	IncludedGroups                     []int       `json:"includedGroups"`
	ExcludedGroups                     []int       `json:"excludedGroups"`
	Id                                 json.Number `json:"id,omitempty"`
	NatRuleIds                         []int       `json:"natRuleIds"`
	Category                           string      `json:"category"`
	AllowOverlappingTranslatedPrefixes *bool       `json:"allowOverlappingTranslatedPrefixes"`
}

// NewNatPolicy new nat policy
func NewNatPolicy(ac *AlkiraClient) *AlkiraAPI[NatPolicy] {
	uri := fmt.Sprintf("%s/tenantnetworks/%s/nat-policies", ac.URI, ac.TenantNetworkId)
	api := &AlkiraAPI[NatPolicy]{ac, uri, true}
	return api
}
