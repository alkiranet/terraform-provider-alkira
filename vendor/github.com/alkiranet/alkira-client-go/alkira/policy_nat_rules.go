// Copyright (C) 2021-2023 Alkira Inc. All Rights Reserved.

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
	TranslationType         string   `json:"translationType"`
	TranslatedPrefixes      []string `json:"translatedPrefixes,omitempty"`
	TranslatedPrefixListIds []int    `json:"translatedPrefixListIds,omitempty"`
	Bidirectional           bool     `json:"bidirectional,omitempty"`
	MatchAndInvalidate      bool     `json:"matchAndInvalidate,omitempty"`
}

type NatRuleActionDstTranslation struct {
	TranslationType         string   `json:"translationType"`
	TranslatedPrefixes      []string `json:"translatedPrefixes,omitempty"`
	TranslatedPrefixListIds []int    `json:"translatedPrefixListIds,omitempty"`
	TranslatedPortList      []string `json:"translatedPortList,omitempty"`
	Bidirectional           bool     `json:"bidirectional,omitempty"`
	AdvertiseToConnector    bool     `json:"advertiseToConnector,omitempty"`
}

type EgressAction struct {
	IpType string `json:"ipType"`
}

// NewNatPolicyRule new NAT policy rule
func NewNatRule(ac *AlkiraClient) *AlkiraAPI[NatPolicyRule] {
	uri := fmt.Sprintf("%s/tenantnetworks/%s/nat-rules", ac.URI, ac.TenantNetworkId)
	api := &AlkiraAPI[NatPolicyRule]{ac, uri}
	return api
}
