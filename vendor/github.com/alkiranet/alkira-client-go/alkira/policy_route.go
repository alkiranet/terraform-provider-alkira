// Copyright (C) 2021-2024 Alkira Inc. All Rights Reserved.

package alkira

import (
	"encoding/json"
	"fmt"
)

type RoutePolicy struct {
	Name                          string             `json:"name"`
	Description                   string             `json:"description"`
	Enabled                       bool               `json:"enabled"`
	Direction                     string             `json:"direction"`
	Segment                       string             `json:"segment"`
	IncludedGroups                []int              `json:"includedGroups"`
	ExcludedGroups                []int              `json:"excludedGroups,omitempty"`
	Id                            json.Number        `json:"id,omitempty"` // response only
	AdvertiseInternetExit         *bool              `json:"advertiseInternetExit"`
	AdvertiseOnPremRoutes         bool               `json:"advertiseOnPremRoutes,omitempty"`
	AdvertiseCustomRoutesPrefixId int                `json:"advertiseCustomRoutesPrefixId,omitempty"`
	Rules                         []RoutePolicyRules `json:"rules,omitempty"`
}

type RoutePolicyRules struct {
	Action                       string                                        `json:"action"`
	Name                         string                                        `json:"name"`
	Match                        RoutePolicyRulesMatch                         `json:"match"`
	SequenceNo                   int                                           `json:"sequenceNo,omitempty"` // response only
	Set                          *RoutePolicyRulesSet                          `json:"set,omitempty"`
	InterCxpRoutesRedistribution *RoutePolicyRulesInterCxpRoutesRedistribution `json:"interCxpRoutesRedistribution,omitempty"`
}

type RoutePolicyRulesMatch struct {
	All                      bool     `json:"all"`
	AsPathListIds            []int    `json:"asPathListIds,omitempty"`
	CommunityListIds         []int    `json:"communityListIds,omitempty"`
	ExtendedCommunityListIds []int    `json:"extendedCommunityListIds,omitempty"`
	PrefixListIds            []int    `json:"prefixListIds,omitempty"`
	Cxps                     []string `json:"cxps,omitempty"`
	ConnectorGroupIds        []int    `json:"connectorGroupIds,omitempty"`
}

type RoutePolicyRulesSet struct {
	AsPathPrepend     string `json:"asPathPrepend"`
	Community         string `json:"community"`
	ExtendedCommunity string `json:"extendedCommunity"`
}

type RoutePolicyRulesInterCxpRoutesRedistribution struct {
	DistributionType        string   `json:"distributionType,omitempty"`
	RedistributeAsSecondary bool     `json:"redistributeAsSecondary,omitempty"`
	RestrictedCxps          []string `json:"restrictedCxps,omitempty"`
}

// NewRoutePolicy new route policy
func NewRoutePolicy(ac *AlkiraClient) *AlkiraAPI[RoutePolicy] {
	uri := fmt.Sprintf("%s/tenantnetworks/%s/route-policies", ac.URI, ac.TenantNetworkId)
	api := &AlkiraAPI[RoutePolicy]{ac, uri, true}
	return api
}
