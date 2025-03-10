// Copyright (C) 2022-2025 Alkira Inc. All Rights Reserved.

package alkira

import (
	"encoding/json"
	"fmt"
)

type ConnectorRemoteAccessTemplate struct {
	AdvancedOptions       ConnectorRemoteAccessAdvancedOptions  `json:"advancedOptions"`
	Arguments             []ConnectorRemoteAccessArguments      `json:"arguments"`
	AuthenticationOptions ConnectorRemoteAccessAuthOptions      `json:"authenticationOptions"`
	Id                    json.Number                           `json:"id,omitempty"`
	Name                  string                                `json:"name"`
	SegmentOptions        []ConnectorRemoteAccessSegmentOptions `json:"segmentOptions"`
	Segments              []string                              `json:"segments,omitempty"`
	BannerText            string                                `json:"bannerText,omitempty"`
}

type ConnectorRemoteAccessAdvancedOptions struct {
	EnableDynamicRegionMapping bool   `json:"enableDynamicRegionMapping"`
	MaxActiveUsersThreshold    int    `json:"maxActiveUsersThreshold"`
	NameServer                 string `json:"nameServer"`
	FallbackToTcp              bool   `json:"fallbackToTcp"`
}

type ConnectorRemoteAccessArguments struct {
	BillingTags []int  `json:"billingTags,omitempty"`
	Cxp         string `json:"cxp"`
	Size        string `json:"size"`
}

type ConnectorRemoteAccessAuthOptions struct {
	LdapSettings   *ConnectorRemoteAccessLdapSettings `json:"ldapSettings,omitempty"`
	SupportedModes []string                           `json:"supportedModes"`
}

type ConnectorRemoteAccessLdapSettings struct {
	BindUserDomain      string `json:"bindUserDomain,omitempty"`
	DestinationAddress  string `json:"destinationAddress,omitempty"`
	LdapType            string `json:"ldapType,omitempty"`
	ManagementSegmentId int    `json:"managementSegmentId,omitempty"`
	SearchScopeDomain   string `json:"searchScopeDomain,omitempty"`
}

type ConnectorRemoteAccessSegmentOptions struct {
	SegmentId         int                                      `json:"segmentId"`
	UserGroupMappings []ConnectorRemoteAccessUserGroupMappings `json:"userGroupMappings"`
}

type ConnectorRemoteAccessUserGroupMappings struct {
	BillingTag          int                                       `json:"billingTag,omitempty"`
	Name                string                                    `json:"name"`
	PrefixListId        int                                       `json:"prefixListId"`
	SplitTunneling      bool                                      `json:"splitTunneling"`
	CxpToSubnetsMapping []ConnectorRemoteAccessCxpToSubnetMapping `json:"cxpToSubnetsMapping"`
}

type ConnectorRemoteAccessCxpToSubnetMapping struct {
	Cxp     string   `json:"cxp"`
	Subnets []string `json:"subnets"`
}

// NewConnectorRemoteAccessTemplate
func NewConnectorRemoteAccessTemplate(ac *AlkiraClient) *AlkiraAPI[ConnectorRemoteAccessTemplate] {
	uri := fmt.Sprintf("%s/tenantnetworks/%s/alkira-remote-access-connector-templates", ac.URI, ac.TenantNetworkId)
	api := &AlkiraAPI[ConnectorRemoteAccessTemplate]{ac, uri, true}
	return api
}
