package alkira

import (
	"encoding/json"
	"fmt"
)

type AzureVnetThirdPartyConnectorAttachment struct {
	Name                string      `json:"name"`
	Description         string      `json:"description,omitempty"`
	CxpPeeringGatewayId int         `json:"cxpPeeringGatewayId"`
	VnetId              string      `json:"vnetId"`
	Id                  json.Number `json:"id,omitempty"`    // response only
	InternalName        string      `json:"internalName,omitempty"` // response only
	State               string      `json:"state,omitempty"` // response only
}

func NewAzureVnetThirdPartyConnectorAttachment(ac *AlkiraClient) *AlkiraAPI[AzureVnetThirdPartyConnectorAttachment] {
	uri := fmt.Sprintf("%s/tenantnetworks/%s/azure-vnet-third-party-connector-attachments", ac.URI, ac.TenantNetworkId)
	api := &AlkiraAPI[AzureVnetThirdPartyConnectorAttachment]{ac, uri, false}
	return api
}
