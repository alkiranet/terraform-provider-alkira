// Copyright (C) 2021 Alkira Inc. All Rights Reserved.

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
	Id                            json.Number        `json:"id,omitempty"`
	AdvertiseInternetExit         bool               `json:"advertiseInternetExit,omitempty"`
	AdvertiseOnPremRoutes         bool               `json:"advertiseOnPremRoutes,omitempty"`
	AdvertiseCustomRoutesPrefixId int                `json:"advertiseCustomRoutesPrefixId,omitempty"`
	Rules                         []RoutePolicyRules `json:"rules"`
}

type RoutePolicyRules struct {
	Action                       string                                         `json:"action"`
	Name                         string                                         `json:"name"`
	Match                        RoutePolicyRulesMatch                          `json:"match"`
	Set                          *RoutePolicyRulesSet                           `json:"set,omitempty"`
	InterCxpRoutesRedistribution *RoutePolicyRulesInterCxpRoutesRedistribution   `json:"interCxpRoutesRedistribution"`
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
	DistributionType        string   `json:"distributionType"`
	RedistributeAsSecondary bool     `json:"redistributeAsSecondary,omitempty"`
	RestrictedCxps          []string `json:"restrictedCxps,omitempty"`
}

// CreateRoutePolicy create a route policy
func (ac *AlkiraClient) CreateRoutePolicy(p *RoutePolicy) (string, error) {
	uri := fmt.Sprintf("%s/tenantnetworks/%s/route-policies", ac.URI, ac.TenantNetworkId)

	// Construct the request
	body, err := json.Marshal(p)

	if err != nil {
		return "", err
	}

	data, err := ac.create(uri, body, true)

	if err != nil {
		return "", fmt.Errorf("CreateRoutePolicy: request failed: %v", err)
	}

	var result RoutePolicy
	json.Unmarshal([]byte(data), &result)

	if err != nil {
		return "", fmt.Errorf("CreateRoutePolicy: failed to unmarshal: %v", err)
	}

	return string(result.Id), nil
}

// DeleteRoutePolicy delete a route policy by Id
func (ac *AlkiraClient) DeleteRoutePolicy(id string) error {
	uri := fmt.Sprintf("%s/tenantnetworks/%s/route-policies/%s", ac.URI, ac.TenantNetworkId, id)

	return ac.delete(uri, true)
}

// UpdateRoutePolicy update a route policy by Id
func (ac *AlkiraClient) UpdateRoutePolicy(id string, p *RoutePolicy) error {
	uri := fmt.Sprintf("%s/tenantnetworks/%s/route-policies/%s", ac.URI, ac.TenantNetworkId, id)

	// Construct the request
	body, err := json.Marshal(p)

	if err != nil {
		return fmt.Errorf("UpdateRoutePolicy: failed to marshal: %v", err)
	}

	return ac.update(uri, body, true)
}

// GetRoutePolicy get a route policy by Id
func (ac *AlkiraClient) GetRoutePolicy(id string) (*RoutePolicy, error) {
	uri := fmt.Sprintf("%s/tenantnetworks/%s/route-policies/%s", ac.URI, ac.TenantNetworkId, id)

	data, err := ac.get(uri)

	if err != nil {
		return nil, err
	}

	var result RoutePolicy
	err = json.Unmarshal([]byte(data), &result)

	if err != nil {
		return nil, fmt.Errorf("GetRoutePolicy: failed to unmarshal: %v", err)
	}

	return &result, nil
}
