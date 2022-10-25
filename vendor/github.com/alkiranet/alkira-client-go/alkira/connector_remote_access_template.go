package alkira

import (
	"encoding/json"
	"fmt"
)

// POST
type RemoteAccessConnectorTemplate struct {
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
	//Subnets            []string                                  `json:"subnets"`
	UserGroupID int `json:"userGroupId,omitempty"`
}

func (ac *AlkiraClient) CreateRemoteAccessConnectorTemplate(r *RemoteAccessConnectorTemplate) (string, error) {
	uri := fmt.Sprintf("%s/tenantnetworks/%s/alkira-remote-access-connector-templates", ac.URI, ac.TenantNetworkId)
	body, err := json.Marshal(r)

	if err != nil {
		return "", fmt.Errorf("CreateRemoteAccessConnectorTemplate: marshal failed: %v", err)
	}

	data, err := ac.create(uri, body, true)

	if err != nil {
		return "", err
	}

	var result RemoteAccessConnectorTemplate
	err = json.Unmarshal([]byte(data), &result)

	if err != nil {
		return "", fmt.Errorf("CreateRemoteAccessConnectorTemplate: failed to unmarshal: %v", err)
	}

	return string(result.Id), nil
}

func (ac *AlkiraClient) GetRemoteAccessConnectorTemplates() (string, error) {
	uri := fmt.Sprintf("%s/tenantnetworks/%s/alkira-remote-access-connector-templates", ac.URI, ac.TenantNetworkId)
	data, err := ac.get(uri)

	if err != nil {
		return "", err
	}

	return string(data), nil
}
func (ac *AlkiraClient) GetRemoteAccessConnectorTemplateById(id string) (*RemoteAccessConnectorTemplate, error) {
	uri := fmt.Sprintf("%s/tenantnetworks/%s/alkira-remote-access-connector-templates/%s", ac.URI, ac.TenantNetworkId, id)

	var ract RemoteAccessConnectorTemplate

	data, err := ac.get(uri)

	if err != nil {
		return nil, err
	}

	err = json.Unmarshal([]byte(data), &ract)

	if err != nil {
		return nil, fmt.Errorf("GetRemoteAccessConnectorTemplateById: failed to unmarshal: %v", err)
	}

	return &ract, nil
}

func (ac *AlkiraClient) UpdateRemoteAccessConnectorTemplate(id string, r *RemoteAccessConnectorTemplate) error {
	uri := fmt.Sprintf("%s/tenantnetworks/%s/alkira-remote-access-connector-templates/%s", ac.URI, ac.TenantNetworkId, id)

	body, err := json.Marshal(r)

	if err != nil {
		return fmt.Errorf("UpdateRemoteAccessConnectorTemplate: failed to marshal: %v", err)
	}

	return ac.update(uri, body, true)
}

func (ac *AlkiraClient) DeleteRemoteAccessConnectorTemplate(id string) error {
	uri := fmt.Sprintf("%s/tenantnetworks/%s/alkira-remote-access-connector-templates/%s", ac.URI, ac.TenantNetworkId, id)

	return ac.delete(uri, true)
}
