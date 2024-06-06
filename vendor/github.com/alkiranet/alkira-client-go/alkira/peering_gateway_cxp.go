// Copyright (C) 2024 Alkira Inc. All Rights Reserved.

package alkira

import (
	"encoding/json"
	"fmt"
)

type PeeringGatewayCxp struct {
	Name          string `json:"name"`
	Description   string `json:"description,omitempty"`
	Cxp           string `json:"cxp"`
	CloudProvider string `json:"cloudProvider"`
	CloudRegion   string `json:"cloudRegion"`
	Segment       string `json:"segment"`
	// SegmentId     string      `json:"segmentId"`
	Id    json.Number `json:"id,omitempty"`    // response only
	State string      `json:"state,omitempty"` // response only
}

// NewConnectorPeeringGatewayAwsTgwAttachment new peering gateway aws tgw attachment
func NewPeeringGatewayCxp(ac *AlkiraClient) *AlkiraAPI[PeeringGatewayCxp] {
	uri := fmt.Sprintf("%s/tenantnetworks/%s/cxp-peering-gateways", ac.URI, ac.TenantNetworkId)
	api := &AlkiraAPI[PeeringGatewayCxp]{ac, uri, false}
	return api
}
