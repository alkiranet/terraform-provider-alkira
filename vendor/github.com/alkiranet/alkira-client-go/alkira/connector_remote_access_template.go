// Copyright (C) 2022-2023 Alkira Inc. All Rights Reserved.

package alkira

import (
	"encoding/json"
	"fmt"
)

type ConnectorRemoteAccessTemplate struct {
	AdvancedOptions       RemoteAccessConnectorTemplateAdvancedOptions  `json:"advancedOptions"`
	Arguments             []RemoteAccessConnectorTemplateArguments      `json:"arguments"`
	AuthenticationOptions RemoteAccessConnectorTemplateAuthOptions      `json:"authenticationOptions"`
	DocState              string                                        `json:"docState,omitempty"`
	Id                    json.Number                                   `json:"id"`
	InternalName          string                                        `json:"internalName,omitempty"`
	Name                  string                                        `json:"name"`
	SamlIDPMetadata       string                                        `json:"samlIDPMetadata"`
	SegmentIds            []int                                         `json:"segmentIds"`
	SegmentOptions        []RemoteAccessConnectorTemplateSegmentOptions `json:"segmentOptions"`
	Segments              []string                                      `json:"segments"`
	State                 string                                        `json:"state,omitempty"`
}

type RemoteAccessConnectorTemplateAdvancedOptions struct {
	EnableDynamicRegionMapping bool   `json:"enableDynamicRegionMapping"`
	MaxActiveUsersThreshold    int    `json:"maxActiveUsersThreshold"`
	NameServer                 string `json:"nameServer"`
}

type RemoteAccessConnectorTemplateArguments struct {
	BillingTags []int  `json:"billingTags,omitempty"`
	Cxp         string `json:"cxp"`
	Size        string `json:"size"`
}

type RemoteAccessConnectorTemplateAuthOptions struct {
	LdapSettings   *RemoteAccessTemplateLdapSettings `json:"ldapSettings,omitempty"`
	SupportedModes []string                          `json:"supportedModes"`
}

type RemoteAccessTemplateLdapSettings struct {
	BindUserDomain     string `json:"bindUserDomain,omitempty"`
	CredentialID       string `json:"credentialId,omitempty"`
	DestinationAddress string `json:"destinationAddress,omitempty"`
	LdapType           string `json:"ldapType,omitempty"`
	ManagementSegment  string `json:"managementSegment,omitempty"`
	SearchScopeDomain  string `json:"searchScopeDomain,omitempty"`
}

type RemoteAccessConnectorTemplateSegmentOptions struct {
	Name              string                                  `json:"name"`
	SegmentId         int                                     `json:"segmentId"`
	UserGroupMappings []RemoteAccessTemplateUserGroupMappings `json:"userGroupMappings"`
}

type RemoteAccessTemplateCxpToSubnetMappings struct {
	Cxp     string   `json:"cxp"`
	Subnets []string `json:"subnets"`
}

type RemoteAccessTemplateUserGroupMappings struct {
	BillingTag         int                                       `json:"billingTag,omitempty"`
	CxpToSubnetMapping []RemoteAccessTemplateCxpToSubnetMappings `json:"cxpToSubnetsMapping"`
	GroupID            int                                       `json:"groupId,omitempty"`
	Name               string                                    `json:"name"`
	PrefixListID       *int                                      `json:"prefixListId"`
	RoutingTagID       int                                       `json:"routingTagId,omitempty"`
	SplitTunneling     bool                                      `json:"splitTunneling"`
	UserGroupID        int                                       `json:"userGroupId,omitempty"`
}

// NewConnectorRemoteAccess new connector
func NewConnectorRemoteAccessTemplate(ac *AlkiraClient) *AlkiraAPI[ConnectorRemoteAccessTemplate] {
	uri := fmt.Sprintf("%s/tenantnetworks/%s/alkira-remote-access-connector-templates", ac.URI, ac.TenantNetworkId)
	api := &AlkiraAPI[ConnectorRemoteAccessTemplate]{ac, uri}
	return api
}
