// Copyright (C) 2023-2025 Alkira Inc. All Rights Reserved.

package alkira

import (
	"encoding/json"
	"fmt"
)

type AdditionalTunnelOptionsPerNode struct {
	Id      int    `json:"id"`
	Label   string `json:"label"`
	Enabled bool   `json:"enabled"`
}

type SegmentScaleOptions struct {
	AdditionalTunnelsPerNode       int                              `json:"additionalTunnelsPerNode"`
	AdditionalNodes                int                              `json:"additionalNodes"`
	SegmentId                      int                              `json:"segmentId"`
	ZoneName                       string                           `json:"zoneName"`
	AdditionalTunnelOptionsPerNode []AdditionalTunnelOptionsPerNode `json:"additionalTunnelOptionsPerNode"`
}

type ControllerScaleOptions struct {
	Description          string                `json:"description"`
	DocState             string                `json:"docState,omitempty"`
	EntityId             int                   `json:"entityId"`
	EntityType           string                `json:"entityType"`
	Id                   json.Number           `json:"id,omitempty"`                  // response only
	LastConfigUpdatedAt  int                   `json:"lastConfigUpdatedAt,omitempty"` // response only
	Name                 string                `json:"name"`
	NetworkEntityId      string                `json:"networkEntityId"`
	NetworkEntitySubType string                `json:"networkEntitySubType"`
	NetworkEntityType    string                `json:"networkEntityType"`
	SegmentScaleOptions  []SegmentScaleOptions `json:"segmentScaleOptions"`
	State                string                `json:"state,omitempty"` // response only
}

// NewControllerScaleOptions new controller scale options
func NewControllerScaleOptions(ac *AlkiraClient) *AlkiraAPI[ControllerScaleOptions] {
	uri := fmt.Sprintf("%s/tenantnetworks/%s/network-entity-scale-options", ac.URI, ac.TenantNetworkId)
	api := &AlkiraAPI[ControllerScaleOptions]{ac, uri, true}
	return api
}
