// Copyright (C) 2020-2025 Alkira Inc. All Rights Reserved.

package alkira

import (
	"encoding/json"
	"fmt"
)

type InputPrefixes struct {
	Id    string `json:"id,omitempty"`
	Type  string `json:"type"`
	Value string `json:"value"`
}

type ExportOptions struct {
	Mode     string          `json:"routeExportMode,omitempty"`
	Prefixes []InputPrefixes `json:"userInputPrefixes,omitempty"`
}

type RouteTables struct {
	Id            string `json:"id"`
	PrefixListIds []int  `json:"prefixListIds"`
	Mode          string `json:"routeImportMode"`
}

type ImportOptions struct {
	RouteTables []RouteTables `json:"routeTables"`
}

type ConnectorAwsVpcRouting struct {
	Export interface{} `json:"exportToCXPOptions"`
	Import interface{} `json:"importFromCXPOptions"`
}

type TgwAttachment struct {
	SubnetId         string `json:"subnetId"`
	AvailabilityZone string `json:"availabilityZone"`
}

type ConnectorAwsVpc struct {
	BillingTags                        []int           `json:"billingTags"`
	CXP                                string          `json:"cxp"`
	CredentialId                       string          `json:"credentialId"`
	CustomerName                       string          `json:"customerName"`
	CustomerRegion                     string          `json:"customerRegion"`
	DirectInterVPCCommunicationEnabled bool            `json:"directInterVPCCommunicationEnabled"`
	DirectInterVPCCommunicationGroup   string          `json:"directInterVPCCommunicationGroup,omitempty"`
	Enabled                            bool            `json:"enabled"`
	Group                              string          `json:"group,omitempty"`
	Id                                 json.Number     `json:"id,omitempty"`              // response only
	ImplicitGroupId                    int             `json:"implicitGroupId,omitempty"` // response only
	Name                               string          `json:"name"`
	SecondaryCXPs                      []string        `json:"secondaryCXPs,omitempty"`
	Segments                           []string        `json:"segments"`
	Size                               string          `json:"size"`
	TgwAttachments                     []TgwAttachment `json:"tgwAttachments,omitempty"`
	VpcId                              string          `json:"vpcId"`
	VpcOwnerId                         string          `json:"vpcOwnerId"`
	VpcRouting                         interface{}     `json:"vpcRouting"`
	TgwConnectEnabled                  bool            `json:"tgwConnectEnabled"`
	ScaleGroupId                       string          `json:"scaleGroupId,omitempty"`
	Description                        string          `json:"description,omitempty"`
}

// NewConnectorAwsVpc new connector-aws-vpc
func NewConnectorAwsVpc(ac *AlkiraClient) *AlkiraAPI[ConnectorAwsVpc] {
	uri := fmt.Sprintf("%s/tenantnetworks/%s/awsvpcconnectors", ac.URI, ac.TenantNetworkId)
	api := &AlkiraAPI[ConnectorAwsVpc]{ac, uri, true}
	return api
}
