// Copyright (C) 2021 Alkira Inc. All Rights Reserved.

package alkira

import (
	"encoding/json"
	"fmt"
)

type CloudProviderAccount struct {
	Name          string `json:"name"`
	Id            string `json:"id,omitempty"`
	CredentialId  string `json:"credentialId"`
	CloudProvider string `json:"cloudProvider"`
	AutoSync      string `json:"autoSync"`
	NativeId      string `json:"nativeId"`
}

// GetCloudProviderAccounts get all cloud provider accounts
func (ac *AlkiraClient) GetCloudProviderAccounts() (string, error) {
	uri := fmt.Sprintf("%s/cloud-provider-accounts", ac.URI)

	data, err := ac.get(uri)
	return string(data), err
}

// GetCloudProviderAccountById get single cloud provider account by Id
func (ac *AlkiraClient) GetCloudProviderAccountById(id string) (CloudProviderAccount, error) {
	uri := fmt.Sprintf("%s/cloud-provider-accounts/%s", ac.URI, id)

	var account CloudProviderAccount

	data, err := ac.get(uri)

	if err != nil {
		return account, err
	}

	err = json.Unmarshal([]byte(data), &account)

	if err != nil {
		return account, fmt.Errorf("GetCloudProviderAccountById: failed to unmarshal: %v", err)
	}

	return account, nil
}

// GetCloudProviderAccountByName get cloud provider account by its name
func (ac *AlkiraClient) GetCloudProviderAccountByName(name string) (CloudProviderAccount, error) {
	var account CloudProviderAccount

	if len(name) == 0 {
		return account, fmt.Errorf("Invalid cloud provider account name input")
	}

	accounts, err := ac.GetCloudProviderAccounts()

	if err != nil {
		return account, err
	}

	var result []CloudProviderAccount
	json.Unmarshal([]byte(accounts), &result)

	for _, g := range result {
		if g.Name == name {
			return g, nil
		}
	}

	return account, fmt.Errorf("failed to find cloud provider account by %s", name)
}

// CreateCloudProviderAccount create a new cloud provider account
func (ac *AlkiraClient) CreateCloudProviderAccount(account CloudProviderAccount) (string, error) {
	uri := fmt.Sprintf("%s/cloud-provider-accounts", ac.URI)

	// Construct the request
	body, err := json.Marshal(account)

	if err != nil {
		return "", fmt.Errorf("CreateCloudProviderAccount: failed to marshal: %v", err)
	}

	data, err := ac.create(uri, body)

	if err != nil {
		return "", err
	}

	var result CloudProviderAccount
	err = json.Unmarshal([]byte(data), &result)

	if err != nil {
		return "", fmt.Errorf("CreateCloudProviderAccount: failed to unmarshal: %v", err)
	}

	return result.Id, nil
}

// DeleteCloudProviderAccount delete a cloud provider account by Id
func (ac *AlkiraClient) DeleteCloudProviderAccount(id string) error {
	uri := fmt.Sprintf("%s/cloud-provider-accounts/%s", ac.URI, id)
	return ac.delete(uri)
}

// UpdateCloudProviderAccount update a cloud provider account by Id
func (ac *AlkiraClient) UpdateCloudProviderAccount(id string, account CloudProviderAccount) error {
	uri := fmt.Sprintf("%s/cloud-provider-accounts/%s", ac.URI, id)

	// Construct the request
	body, err := json.Marshal(account)

	if err != nil {
		return fmt.Errorf("UpdateCloudProviderAccount: failed to marshal: %v", err)
	}

	return ac.update(uri, body)
}
