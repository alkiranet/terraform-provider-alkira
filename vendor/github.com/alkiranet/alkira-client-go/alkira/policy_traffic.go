// Copyright (C) 2020-2025 Alkira Inc. All Rights Reserved.

package alkira

import (
	"encoding/json"
	"fmt"
)

type TrafficPolicy struct {
	Description   string      `json:"description"`
	Enabled       bool        `json:"enabled"`
	FromGroups    []int       `json:"fromGroups"`
	Id            json.Number `json:"id,omitempty"`
	Name          string      `json:"name"`
	RuleListId    int         `json:"ruleListId"`
	SegmentIds    []int       `json:"segmentIds"`
	ToGroups      []int       `json:"toGroups"`
	ZTAProfileIds []string    `json:"ztaProfileIds"`
}

// NewTrafficPolicy new traffic policy
func NewTrafficPolicy(ac *AlkiraClient) *AlkiraAPI[TrafficPolicy] {
	uri := fmt.Sprintf("%s/tenantnetworks/%s/policy/policies", ac.URI, ac.TenantNetworkId)
	api := &AlkiraAPI[TrafficPolicy]{ac, uri, true}
	return api
}
