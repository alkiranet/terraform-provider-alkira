// Copyright (C) 2021-2025 Alkira Inc. All Rights Reserved.

package alkira

import (
	"encoding/json"
	"fmt"
)

type NatPolicyRule struct {
	Name        string        `json:"name"`
	Description string        `json:"description"`
	Id          json.Number   `json:"id,omitempty"`
	Enabled     bool          `json:"enabled"`
	Match       NatRuleMatch  `json:"match"`
	Action      NatRuleAction `json:"action"`
	Category    string        `json:"category"`
	Direction   string        `json:"direction,omitempty"`
}

type NatRuleMatch struct {
	SourcePrefixes      []string `json:"sourcePrefixes,omitempty"`
	SourcePrefixListIds []int    `json:"sourcePrefixListIds,omitempty"`
	DestPrefixes        []string `json:"destPrefixes,omitempty"`
	DestPrefixListIds   []int    `json:"destPrefixListIds,omitempty"`
	SourcePortList      []string `json:"sourcePortList,omitempty"`
	DestPortList        []string `json:"destPortList,omitempty"`
	Protocol            string   `json:"protocol"`
}

type NatRuleAction struct {
	SourceAddressTranslation      NatRuleActionSrcTranslation `json:"sourceAddressTranslation"`
	DestinationAddressTranslation NatRuleActionDstTranslation `json:"destinationAddressTranslation"`
	Egress                        EgressAction                `json:"egress"`
}

type NatRuleActionSrcTranslation struct {
	TranslationType         string                `json:"translationType"`
	TranslatedPrefixes      []string              `json:"translatedPrefixes,omitempty"`
	TranslatedPrefixListIds []int                 `json:"translatedPrefixListIds,omitempty"`
	Bidirectional           *bool                 `json:"bidirectional,omitempty"`
	MatchAndInvalidate      *bool                 `json:"matchAndInvalidate,omitempty"`
	RoutingOptions          NatRuleRoutingOptions `json:"routingOptions,omitempty"`
}

type NatRuleActionDstTranslation struct {
	TranslationType            string                `json:"translationType"`
	TranslatedPrefixes         []string              `json:"translatedPrefixes,omitempty"`
	TranslatedPrefixListIds    []int                 `json:"translatedPrefixListIds,omitempty"`
	TranslatedPortList         []string              `json:"translatedPortList,omitempty"`
	TranslatedPolicyFqdnListId int                   `json:"translatedPolicyFqdnListId,omitempty"`
	Bidirectional              *bool                 `json:"bidirectional,omitempty"`
	AdvertiseToConnector       *bool                 `json:"advertiseToConnector,omitempty"`
	RoutingOptions             NatRuleRoutingOptions `json:"routingOptions,omitempty"`
}

type NatRuleRoutingOptions struct {
	TrackPrefixes                  []string `json:"trackPrefixes,omitempty"`
	TrackPrefixListIds             []int    `json:"trackPrefixListIds,omitempty"`
	InvalidateRoutingTrackPrefixes *bool    `json:"invalidateRoutingTrackPrefixes,omitempty"`
}

type EgressAction struct {
	IpType string `json:"ipType"`
}

// NewNatPolicyRule new NAT policy rule
func NewNatRule(ac *AlkiraClient) *AlkiraAPI[NatPolicyRule] {
	uri := fmt.Sprintf("%s/tenantnetworks/%s/nat-rules", ac.URI, ac.TenantNetworkId)
	api := &AlkiraAPI[NatPolicyRule]{ac, uri, true}
	return api
}
