// Copyright (C) 2023-2025 Alkira Inc. All Rights Reserved.

package alkira

import (
	"encoding/json"
	"fmt"
)

type InterConnectorCommunicationGroup struct {
	Id                           json.Number `json:"id"`
	Name                         string      `json:"name"`
	Description                  string      `json:"description"`
	Segment                      string      `json:"segment"`
	Cxp                          string      `json:"cxp"`
	ConnectorProviderRegion      string      `json:"connectorProviderRegion"`
	ConnectorType                string      `json:"connectorType"`
	VirtualNetworkManagerAzureId int         `json:"azureVirtualNetworkManagerId"`
}

func NewInterConnectorCommunicationGroup(ac *AlkiraClient) *AlkiraAPI[InterConnectorCommunicationGroup] {
	uri := fmt.Sprintf("%s/tenantnetworks/%s/inter-connector-communication-groups", ac.URI, ac.TenantNetworkId)
	api := &AlkiraAPI[InterConnectorCommunicationGroup]{ac, uri, true}
	return api
}
