package alkira

import (
	"fmt"
)

type VirtualNetworkManagerAzure struct {
	Name                 string   `json:"name"`
	Region               string   `json:"region"`
	SubscriptionId       string   `json:"subscriptionId"`
	ResourceGroup        string   `json:"resourceGroup"`
	Description          string   `json:"description,omitempty"`
	CredentialsId        string   `json:"credentialsId"`
	SubscriptionsInScope []string `json:"subscriptionsInScope"`
}

func NewVritualNetworkManagerAzure(ac *AlkiraClient) *AlkiraAPI[VirtualNetworkManagerAzure] {
	uri := fmt.Sprintf("%s/tenantnetworks/%s/azure-virtual-network-managers", ac.URI, ac.TenantNetworkId)
	api := &AlkiraAPI[VirtualNetworkManagerAzure]{ac, uri, true}
	return api
}
