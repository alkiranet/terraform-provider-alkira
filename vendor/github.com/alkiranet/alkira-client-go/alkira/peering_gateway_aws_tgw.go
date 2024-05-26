// Copyright (C) 2024 Alkira Inc. All Rights Reserved.

package alkira

import (
	"encoding/json"
	"fmt"
)

type PeeringGatewayAwsTgw struct {
	Name             string      `json:"name"`
	Description      string      `json:"description,omitempty"`
	Asn              int         `json:"asn"`
	AwsRegion        string      `json:"awsRegion"`
	Id               json.Number `json:"id,omitempty"` // response only
	State            string      `json:"state,omitempty"` // response only
}

// NewPeeringGatewayAwsTgw new peering gateway AWS-TGW
func NewPeeringGatewayAwsTgw(ac *AlkiraClient) *AlkiraAPI[PeeringGatewayAwsTgw] {
	uri := fmt.Sprintf("%s/tenantnetworks/%s/aws-tgws", ac.URI, ac.TenantNetworkId)
	api := &AlkiraAPI[PeeringGatewayAwsTgw]{ac, uri, false}
	return api
}
