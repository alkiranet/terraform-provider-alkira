// Copyright (C) 2024 Alkira Inc. All Rights Reserved.
package alkira

import (
	"encoding/json"
	"fmt"
)

type VirtualNetworkManagerAzure struct {
	Name                 string      `json:"name"`
	Region               string      `json:"region"`
	SubscriptionId       string      `json:"subscriptionId"`
	ResourceGroup        string      `json:"resourceGroup"`
	Description          string      `json:"description,omitempty"`
	CredentialId         string      `json:"credentialId"`
	Id                   json.Number `json:"id,omitempty"`    // RESPONSE ONLY.
	State                string      `json:"state,omitempty"` // RESPONSE ONLY.
	SubscriptionsInScope []string    `json:"subscriptionsInScope"`
}

func NewVirtualNetworkManagerAzure(ac *AlkiraClient) *AlkiraAPI[VirtualNetworkManagerAzure] {
	uri := fmt.Sprintf("%s/tenantnetworks/%s/azure-virtual-network-managers", ac.URI, ac.TenantNetworkId)
	api := &AlkiraAPI[VirtualNetworkManagerAzure]{ac, uri, true}
	return api
}
