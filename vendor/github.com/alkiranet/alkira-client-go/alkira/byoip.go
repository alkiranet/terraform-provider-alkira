// Copyright (C) 2022 Alkira Inc. All Rights Reserved.

package alkira

import (
	"encoding/json"
	"fmt"
)

type ByoipExtraAttributes struct {
	Message   string `json:"message"`
	Signature string `json:"signature"`
	PublicKey string `json:"publicKey"`
}

type Byoip struct {
	Prefix          string               `json:"prefix"`
	Cxp             string               `json:"cxp"`
	Description     string               `json:"description"`
	ExtraAttributes ByoipExtraAttributes `json:"extraAttributes"`
	DoNotAdvertise  bool                 `json:"doNotAdvertise"`
	Id              json.Number          `json:"id,omitempty"` // response only
}

// CreateByoip create a new BYOIP
func (ac *AlkiraClient) CreateByoip(byoip *Byoip) (string, error) {
	uri := fmt.Sprintf("%s/tenantnetworks/%s/byoips", ac.URI, ac.TenantNetworkId)

	body, err := json.Marshal(byoip)

	if err != nil {
		return "", fmt.Errorf("CreateByoip: failed to marshal: %v", err)
	}

	data, err := ac.create(uri, body, true)

	if err != nil {
		return "", err
	}

	var result Byoip
	err = json.Unmarshal([]byte(data), &result)

	if err != nil {
		return "", fmt.Errorf("CreateByoip: failed to unmarshal: %v", err)
	}

	return string(result.Id), nil
}

// DeleteByoip delete a BYOIP by ID
func (ac *AlkiraClient) DeleteByoip(id string) error {
	uri := fmt.Sprintf("%s/tenantnetworks/%s/byoips/%s", ac.URI, ac.TenantNetworkId, id)
	return ac.delete(uri, true)
}

// UpdateByoip update a BYOIP by ID
func (ac *AlkiraClient) UpdateByoip(id string, byoip *Byoip) error {
	uri := fmt.Sprintf("%s/tenantnetworks/%s/byoips/%s", ac.URI, ac.TenantNetworkId, id)

	body, err := json.Marshal(byoip)

	if err != nil {
		return fmt.Errorf("UpdateByoip: failed to marshal: %v", err)
	}

	return ac.update(uri, body, true)
}

// GetByoips get all BYOIP from the given tenant network
func (ac *AlkiraClient) GetByoips() (string, error) {
	uri := fmt.Sprintf("%s/tenantnetworks/%s/byoips", ac.URI, ac.TenantNetworkId)
	data, err := ac.get(uri)
	return string(data), err
}

// GetByoipById get a single BYOIP by ID
func (ac *AlkiraClient) GetByoipById(id string) (*Byoip, error) {
	uri := fmt.Sprintf("%s/tenantnetworks/%s/byoips/%s", ac.URI, ac.TenantNetworkId, id)

	var byoip Byoip

	data, err := ac.get(uri)

	if err != nil {
		return nil, err
	}

	err = json.Unmarshal([]byte(data), &byoip)

	if err != nil {
		return nil, fmt.Errorf("GetByoipById: failed to unmarshal: %v", err)
	}

	return &byoip, nil
}

// GetByoipByName get a single BYOIP by name
func (ac *AlkiraClient) GetByoipPrefixByName(name string) (Byoip, error) {
	var byoip Byoip

	if len(name) == 0 {
		return byoip, fmt.Errorf("Invalid BYOIP name input")
	}

	byoips, err := ac.GetByoips()

	if err != nil {
		return byoip, err
	}

	var result []Byoip
	err = json.Unmarshal([]byte(byoips), &result)

	if err != nil {
		return byoip, fmt.Errorf("GetByoipByName: failed to unmarshal: %v", err)
	}

	for _, t := range result {
		if t.Prefix == name {
			return t, nil
		}
	}

	return byoip, fmt.Errorf("failed to find the BYOIP by %s", name)
}
