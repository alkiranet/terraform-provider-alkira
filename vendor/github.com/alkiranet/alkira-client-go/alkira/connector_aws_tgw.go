// Copyright (C) 2024 Alkira Inc. All Rights Reserved.

package alkira

import (
	"encoding/json"
	"fmt"
)

type ConnectorAwsTgw struct {
	CXP                       string      `json:"cxp"`
	Name                      string      `json:"name"`
	Group                     string      `json:"group,omitempty"`
	Segments                  []string    `json:"segments"`
	Size                      string      `json:"size"`
	Enabled                   bool        `json:"enabled"`
	AwsTgwPeeringAttachmentId int         `json:"awsTgwPeeringAttachmentId"`
	BillingTags               []int       `json:"billingTags"`
	StaticRoutes              []int       `json:"staticRoutes"`
	Id                        json.Number `json:"id,omitempty"`              // response only
	ImplicitGroupId           int         `json:"implicitGroupId,omitempty"` // response only
}

// NewConnectorAwsTgw new connector-aws-tgw
func NewConnectorAwsTgw(ac *AlkiraClient) *AlkiraAPI[ConnectorAwsTgw] {
	uri := fmt.Sprintf("%s/tenantnetworks/%s/aws-tgw-connectors", ac.URI, ac.TenantNetworkId)
	api := &AlkiraAPI[ConnectorAwsTgw]{ac, uri, true}
	return api
}
