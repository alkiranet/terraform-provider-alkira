package alkira

import (
	"encoding/json"
	"fmt"
)

type AzureVnetThirdPartyConnector struct {
	Name                                     string      `json:"name"`
	Description                              string      `json:"description,omitempty"`
	CXP                                      string      `json:"cxp"`
	Enabled                                  bool        `json:"enabled"`
	Group                                    string      `json:"group,omitempty"`
	Segments                                 []string    `json:"segments"`
	Size                                     string      `json:"size"`
	AzureVnetThirdPartyConnectorAttachmentId int         `json:"azureVnetThirdPartyConnectorAttachmentId"`
	BillingTags                              []int       `json:"billingTags,omitempty"`
	StaticRoutes                             []int       `json:"staticRoutes,omitempty"`
	Id                                       json.Number `json:"id,omitempty"`              // response only
	ImplicitGroupId                          int         `json:"implicitGroupId,omitempty"` // response only
	CxpPeeringGatewayId                      int         `json:"cxpPeeringGatewayId,omitempty"` // response only
}

func NewAzureVnetThirdPartyConnector(ac *AlkiraClient) *AlkiraAPI[AzureVnetThirdPartyConnector] {
	uri := fmt.Sprintf("%s/tenantnetworks/%s/azure-vnet-third-party-connectors", ac.URI, ac.TenantNetworkId)
	api := &AlkiraAPI[AzureVnetThirdPartyConnector]{ac, uri, true}
	return api
}
