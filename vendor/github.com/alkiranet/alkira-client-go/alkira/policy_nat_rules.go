// Copyright (C) 2021 Alkira Inc. All Rights Reserved.

package alkira

import (
	"encoding/json"
	"fmt"
)

type NatRule struct {
	Name        string        `json:"name"`
	Description string        `json:"description"`
	Id          json.Number   `json:"id,omitempty"`
	Enabled     bool          `json:"enabled"`
	Match       NatRuleMatch  `json:"match"`
	Action      NatRuleAction `json:"action"`
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

// CreateNatRule create a policy rule
func (ac *AlkiraClient) CreateNatRule(rule *NatRule) (string, error) {
	uri := fmt.Sprintf("%s/tenantnetworks/%s/nat-rules", ac.URI, ac.TenantNetworkId)

	// Construct the request
	body, err := json.Marshal(rule)

	if err != nil {
		return "", err
	}

	data, err := ac.create(uri, body, true)

	if err != nil {
		return "", fmt.Errorf("CreateNatRule: request failed: %v", err)
	}

	var result NatRule
	err = json.Unmarshal([]byte(data), &result)

	if err != nil {
		return "", fmt.Errorf("CreateNatRule: request failed: %v", err)
	}

	return string(result.Id), nil
}

// DeleteNatRule delete a policy rule
func (ac *AlkiraClient) DeleteNatRule(id string) error {
	uri := fmt.Sprintf("%s/tenantnetworks/%s/nat-rules/%s", ac.URI, ac.TenantNetworkId, id)
	return ac.delete(uri, true)
}

// UpdateNatRule update a policy rule list
func (ac *AlkiraClient) UpdateNatRule(id string, rule *NatRule) error {
	uri := fmt.Sprintf("%s/tenantnetworks/%s/nat-rules/%s", ac.URI, ac.TenantNetworkId, id)

	// Construct the request
	body, err := json.Marshal(rule)

	if err != nil {
		return fmt.Errorf("UpdateNatRule: failed to marshal: %v", err)
	}

	return ac.update(uri, body, true)
}

// GetNatRule get a policy rule list
func (ac *AlkiraClient) GetNatRule(id string) (*NatRule, error) {
	uri := fmt.Sprintf("%s/tenantnetworks/%s/nat-rules/%s", ac.URI, ac.TenantNetworkId, id)

	data, err := ac.get(uri)

	if err != nil {
		return nil, err
	}

	var result NatRule
	err = json.Unmarshal([]byte(data), &result)

	if err != nil {
		return nil, fmt.Errorf("GetNatRule: failed to unmarshal: %v", err)
	}

	return &result, nil
}
