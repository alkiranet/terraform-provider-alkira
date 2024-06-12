package alkira

import (
	"encoding/json"
	"fmt"
)

type AzureVirtualNetworkManager struct {
	Name                 string      `json:"name"`
	Region               string      `json:"region"`
	SubscriptionId       string      `json:"subscriptionId"`
	ResourceGroup        string      `json:"resourceGroup"`
	Description          string      `json:"description,omitempty"`
	CredentialsId        string      `json:"credentialsId"`
	Id                   json.Number `json:"id,omitempty"`
	State                string      `json:"state,omitempty"`
	SubscriptionsInScope []string    `json:"subscriptionsInScope"`
}

func NewAzureVirtualNetworkManager(ac *AlkiraClient) *AlkiraAPI[AzureVirtualNetworkManager] {
	uri := fmt.Sprintf("%s/tenantnetworks/%s/azure-virtual-network-managers", ac.URI, ac.TenantNetworkId)
	api := &AlkiraAPI[VirtualNetworkManagerAzure]{ac, uri, true}
	return api
}
