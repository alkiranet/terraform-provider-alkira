// Copyright (C) 2024-2025 Alkira Inc. All Rights Reserved.

package alkira

import (
	"encoding/json"
	"fmt"
)

type PeeringGatewayAwsTgwAttachment struct {
	Name                       string           `json:"name"`
	Description                string           `json:"description,omitempty"`
	Requestor                  string           `json:"requestor"`
	PeerAwsRegion              string           `json:"peerAwsRegion,omitempty"`
	PeerAwsTgwId               string           `json:"peerAwsTgwId,omitempty"`
	PeerAwsAccountId           string           `json:"peerAwsAccountId"`
	PeerDirectConnectGatewayId string           `json:"peerDirectConnectGatewayId,omitempty"`
	PeerAllowedPrefixes        []string         `json:"peerAllowedPrefixes,omitempty"`
	AwsTgwId                   int              `json:"awsTgwId"`
	Type                       string           `json:"type,omitempty"`
	Id                         json.Number      `json:"id,omitempty"`               // response only
	State                      string           `json:"state,omitempty"`            // response only
	TriggerProposal            bool             `json:"trigger_proposal,omitempty"` // response only
	FailureReason              string           `json:"failureReason,omitempty"`    // response only
	ProposalDetails            *ProposalDetails `json:"proposalDetails,omitempty"`  //response only
	ProposalStatus             string           `json:"proposalStatus,omitempty"`   // response only
}

type ProposalDetails struct {
	ProposalId    string `json:"proposalId,omitempty"`
	ProposalState string `json:"proposalState,omitempty"`
	CreatedAt     int    `json:"createdAt,omitempty"`
	UpdatedAt     int    `json:"updatedAt,omitempty"`
}

// NewConnectorPeeringGatewayAwsTgwAttachment new peering gateway aws tgw attachment
func NewPeeringGatewayAwsTgwAttachment(ac *AlkiraClient) *AlkiraAPI[PeeringGatewayAwsTgwAttachment] {
	uri := fmt.Sprintf("%s/tenantnetworks/%s/aws-tgw-peering-attachments", ac.URI, ac.TenantNetworkId)
	api := &AlkiraAPI[PeeringGatewayAwsTgwAttachment]{ac, uri, false}
	return api
}
