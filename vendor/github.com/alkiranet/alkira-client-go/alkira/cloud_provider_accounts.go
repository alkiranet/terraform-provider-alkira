// Copyright (C) 2021-2025 Alkira Inc. All Rights Reserved.

package alkira

import (
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

// NewCloudProviderAccounts
func NewCloudProviderAccounts(ac *AlkiraClient) *AlkiraAPI[CloudProviderAccount] {
	uri := fmt.Sprintf("%s/cloud-provider-accounts", ac.URI)
	api := &AlkiraAPI[CloudProviderAccount]{ac, uri, false}
	return api
}
