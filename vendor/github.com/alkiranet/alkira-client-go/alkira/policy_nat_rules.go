// Copyright (C) 2021-2022 Alkira Inc. All Rights Reserved.

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

// CreateNatRule create a policy NAT rule
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

// DeleteNatRule delete a policy NAT rule
func (ac *AlkiraClient) DeleteNatRule(id string) error {
	uri := fmt.Sprintf("%s/tenantnetworks/%s/nat-rules/%s", ac.URI, ac.TenantNetworkId, id)
	return ac.delete(uri, true)
}

// UpdateNatRule update a policy NAT rule
func (ac *AlkiraClient) UpdateNatRule(id string, rule *NatRule) error {
	uri := fmt.Sprintf("%s/tenantnetworks/%s/nat-rules/%s", ac.URI, ac.TenantNetworkId, id)

	// Construct the request
	body, err := json.Marshal(rule)

	if err != nil {
		return fmt.Errorf("UpdateNatRule: failed to marshal: %v", err)
	}

	return ac.update(uri, body, true)
}

// GetNatRule get all policy NAT rules from the given tenant network
func (ac *AlkiraClient) GetNatRules() (string, error) {
	uri := fmt.Sprintf("%s/tenantnetworks/%s/nat-rules", ac.URI, ac.TenantNetworkId)

	data, err := ac.get(uri)
	return string(data), err
}

// GetNatRuleByName get the policy NAT rule by its name
func (ac *AlkiraClient) GetNatRuleByName(name string) (*NatRule, error) {

	if len(name) == 0 {
		return nil, fmt.Errorf("GetNatRuleByName: Invalid rule name")
	}

	rules, err := ac.GetNatRules()

	if err != nil {
		return nil, err
	}

	var result []NatRule
	json.Unmarshal([]byte(rules), &result)

	for _, l := range result {
		if l.Name == name {
			return &l, nil
		}
	}

	return nil, fmt.Errorf("GetNatRuleByName: failed to find the rule by name %s", name)
}

// GetNatRule get a policy NAT rule by ID
func (ac *AlkiraClient) GetNatRuleById(id string) (*NatRule, error) {
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
