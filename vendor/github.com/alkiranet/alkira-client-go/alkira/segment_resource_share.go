// Copyright (C) 2022-2025 Alkira Inc. All Rights Reserved.

package alkira

import (
	"encoding/json"
	"fmt"
)

type SegmentResourceShare struct {
	Id                  json.Number `json:"id,omitempty"` // response only
	Name                string      `json:"name"`
	Description         string      `json:"description,omitempty"`
	ServiceList         []int       `json:"serviceList"`
	DesignatedSegment   string      `json:"designatedSegment"`
	EndAResources       []int       `json:"endAResources"`
	EndBResources       []int       `json:"endBResources"`
	EndARouteLimit      int         `json:"endARouteLimit"`
	EndBRouteLimit      int         `json:"endBRouteLimit"`
	FromEnd             string      `json:"fromEnd,omitempty"`
	Direction           string      `json:"direction"`
	DesignatedSegmentId int         `json:"designatedSegmentId,omitempty"` // response only
	EndASegmentId       int         `json:"endASegmentId,omitempty"`       // response only
	EndBSegmentId       int         `json:"endBSegmentId,omitempty"`       // response only
	RuleListId          int         `json:"ruleListId,omitempty"`
}

// NewSegmentResourceShare new segment resource share
func NewSegmentResourceShare(ac *AlkiraClient) *AlkiraAPI[SegmentResourceShare] {
	uri := fmt.Sprintf("%s/tenantnetworks/%s/segment-resource-shares", ac.URI, ac.TenantNetworkId)
	api := &AlkiraAPI[SegmentResourceShare]{ac, uri, true}
	return api
}
