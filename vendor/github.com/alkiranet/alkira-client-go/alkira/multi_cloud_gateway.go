// Copyright (C) 2026 Alkira Inc. All Rights Reserved.

package alkira

import (
	"encoding/json"
	"fmt"
)

type MultiCloudGateway struct {
	Name         string      `json:"name"`
	Tier         string      `json:"tier"`
	Description  string      `json:"description"`
	Id           json.Number `json:"id,omitempty"`           // RESPONSE ONLY
	NativeId     string      `json:"nativeId,omitempty"`     // RESPONSE ONLY
	State        string      `json:"state,omitempty"`        // RESPONSE ONLY
	External     bool        `json:"external,omitempty"`     // RESPONSE ONLY
	ConnectorIds []int       `json:"connectorIds,omitempty"` // RESPONSE ONLY
}

// NewMultiCloudGateway new Multi-Cloud Gateway, as peering_gateway_aws_tgw.go does
func NewMultiCloudGateway(ac *AlkiraClient) *AlkiraAPI[MultiCloudGateway] {
	uri := fmt.Sprintf("%s/tenantnetworks/%s/multi-cloud-gateways", ac.URI, ac.TenantNetworkId)
	api := &AlkiraAPI[MultiCloudGateway]{ac, uri, false}
	return api
}

// multiCloudGatewayListResponse is the MCGW list page shape: always
// paginated, with no name filter and no ?paginated=false switch, unlike the
// generic AlkiraAPI.GetByName's bare array.
type multiCloudGatewayListResponse struct {
	Data       []MultiCloudGateway `json:"data"`
	Pagination Pagination          `json:"pagination"`
}

const multiCloudGatewayPageLimit = 100

// GetMultiCloudGatewayByName walks the offset/limit pages of the MCGW list
// and matches name exactly. It does not use the generic AlkiraAPI.GetByName,
// which expects a bare array and a ?paginated=false switch this list doesn't
// support.
func GetMultiCloudGatewayByName(ac *AlkiraClient, name string) (*MultiCloudGateway, string, error) {
	if len(name) == 0 {
		return nil, "", fmt.Errorf("GetMultiCloudGatewayByName: invalid resource name")
	}

	baseUri := fmt.Sprintf("%s/tenantnetworks/%s/multi-cloud-gateways", ac.URI, ac.TenantNetworkId)

	var matches []MultiCloudGateway
	var state string
	offset := 0

	for {
		uri := fmt.Sprintf("%s?offset=%d&limit=%d", baseUri, offset, multiCloudGatewayPageLimit)

		data, pageState, err := ac.get(uri)
		if err != nil {
			return nil, pageState, fmt.Errorf("GetMultiCloudGatewayByName: failed to get page at offset %d: %w", offset, err)
		}
		state = pageState

		var page multiCloudGatewayListResponse
		if err := json.Unmarshal(data, &page); err != nil {
			return nil, state, fmt.Errorf("GetMultiCloudGatewayByName: failed to unmarshal page at offset %d: %w", offset, err)
		}

		for _, mcgw := range page.Data {
			if mcgw.Name == name {
				matches = append(matches, mcgw)
			}
		}

		offset += len(page.Data)

		if len(page.Data) == 0 || offset >= page.Pagination.Hits {
			break
		}
	}

	if len(matches) != 1 {
		return nil, state, fmt.Errorf("GetMultiCloudGatewayByName: found %d Multi-Cloud Gateways named %q, want exactly 1", len(matches), name)
	}

	return &matches[0], state, nil
}
