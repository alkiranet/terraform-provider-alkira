// Copyright (C) 2024-2025 Alkira Inc. All Rights Reserved.

package alkira

import (
	"encoding/json"
	"fmt"
)

type PeeringGatewayCxp struct {
	Name          string      `json:"name"`
	Description   string      `json:"description,omitempty"`
	Cxp           string      `json:"cxp"`
	CloudProvider string      `json:"cloudProvider"`
	CloudRegion   string      `json:"cloudRegion"`
	Segment       string      `json:"segment"`
	Id            json.Number `json:"id,omitempty"`    // response only
	State         string      `json:"state,omitempty"` // response only
}

func NewPeeringGatewayCxp(ac *AlkiraClient) *AlkiraAPI[PeeringGatewayCxp] {
	uri := fmt.Sprintf("%s/tenantnetworks/%s/cxp-peering-gateways", ac.URI, ac.TenantNetworkId)
	api := &AlkiraAPI[PeeringGatewayCxp]{ac, uri, false}
	return api
}
