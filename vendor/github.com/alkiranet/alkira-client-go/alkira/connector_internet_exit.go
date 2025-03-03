// Copyright (C) 2020-2025 Alkira Inc. All Rights Reserved.

package alkira

import (
	"encoding/json"
	"fmt"
)

// ENUM for keys: "DEFAULT", "SRC_IP"
type AlgorithmAttributes struct {
	Keys string `json:"keys"`
}

type TrafficDistribution struct {
	Algorithm           string              `json:"algorithm"`
	AlgorithmAttributes AlgorithmAttributes `json:"algorithmAttributes"`
}

type ConnectorInternet struct {
	BillingTags         []int                `json:"billingTags"`
	ByoipId             int                  `json:"byoipId,omitempty"`
	CXP                 string               `json:"cxp"`
	Description         string               `json:"description"`
	Group               string               `json:"group,omitempty"`
	Enabled             bool                 `json:"enabled"`
	Id                  json.Number          `json:"id,omitempty"`              // response only
	ImplicitGroupId     int                  `json:"implicitGroupId,omitempty"` // response only
	PublicIps           []string             `json:"publicIps,omitempty"`
	Name                string               `json:"name"`
	NumOfPublicIPs      int                  `json:"numOfPublicIPs,omitempty"`
	Segments            []string             `json:"segments"`
	TrafficDistribution *TrafficDistribution `json:"trafficDistribution,omitempty"`
	EgressIpTypes       []string             `json:"egressIPTypes"`
}

func NewConnectorInternet(ac *AlkiraClient) *AlkiraAPI[ConnectorInternet] {
	uri := fmt.Sprintf("%s/tenantnetworks/%s/internetconnectors", ac.URI, ac.TenantNetworkId)
	api := &AlkiraAPI[ConnectorInternet]{ac, uri, true}
	return api
}
