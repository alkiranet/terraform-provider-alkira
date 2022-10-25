package alkira

import (
	"encoding/json"
	"fmt"
)

type RemoteAccessConnector struct {
	Id                    int                                   `json:"id,omitempty"`
	TemplateID            int                                   `json:"templateId"`
	Name                  string                                `json:"name"`
	InternalName          string                                `json:"internalName,omitempty"`
	Cxp                   string                                `json:"cxp"`
	TagID                 int                                   `json:"tagId"`
	Size                  string                                `json:"size"`
	BillingTags           []int                                 `json:"billingTags"`
	Segments              []string                              `json:"segments"`
	SegmentIds            []int                                 `json:"segmentIds"`
	SegmentOptions        []RemoteAccessConnectorSegmentOptions `json:"segmentOptions"`
	AuthenticationOptions RemoteAccessConnectorAuthOptions      `json:"authenticationOptions"`
	AdvancedOptions       RemoteAccessAdvancedOptions           `json:"advancedOptions"`
	DocState              string                                `json:"docState,omitempty"`
	State                 string                                `json:"state"`
	DhParamKeysID         int                                   `json:"dhParamKeysId"`
	MaxActiveUsers        int                                   `json:"maxActiveUsers"`
	ServerCertificates    []RemoteAccessServerCerts             `json:"serverCertificates"`
}

type RemoteAccessConnectorSegmentOptions struct {
	Name              string `json:"name"`
	UserGroupMappings []struct {
		BillingTag     int      `json:"billingTag"`
		GroupID        int      `json:"groupId"`
		Name           string   `json:"name"`
		PrefixListID   int      `json:"prefixListId"`
		RoutingTagID   int      `json:"routingTagId"`
		SplitTunneling bool     `json:"splitTunneling"`
		Subnets        []string `json:"subnets"`
		UserGroupID    string   `json:"userGroupId"`
	} `json:"userGroupMappings"`
}

type RemoteAccessConnectorAuthOptions struct {
	LdapSettings struct {
		BindUserDomain     string `json:"bindUserDomain"`
		CredentialID       string `json:"credentialId"`
		DestinationAddress string `json:"destinationAddress"`
		LdapType           string `json:"ldapType"`
		ManagementSegment  string `json:"managementSegment"`
		SearchScopeDomain  string `json:"searchScopeDomain"`
	} `json:"ldapSettings"`
	SupportedModes []string `json:"supportedModes"`
}

type RemoteAccessAdvancedOptions struct {
	EnableDynamicRegionMapping bool   `json:"enableDynamicRegionMapping"`
	MaxActiveUsersThreshold    int    `json:"maxActiveUsersThreshold"`
	NameServer                 string `json:"nameServer"`
}

type RemoteAccessServerCerts struct {
	ServerCertificateID     int `json:"serverCertificateId"`
	ServerRootCertificateID int `json:"serverRootCertificateId"`
}

func (ac *AlkiraClient) GetRemoteAccessConnectors() (string, error) {
	uri := fmt.Sprintf("%s/tenantnetworks/%s/alkira-remote-access-connectors", ac.URI, ac.TenantNetworkId)
	data, err := ac.get(uri)

	if err != nil {
		return "", err
	}

	return string(data), nil
}

func (ac *AlkiraClient) GetRemoteAccessConnectorById(id string) (*RemoteAccessConnector, error) {
	uri := fmt.Sprintf("%s/tenantnetworks/%s/alkira-remote-access-connectors/%s", ac.URI, ac.TenantNetworkId, id)

	var rac RemoteAccessConnector

	data, err := ac.get(uri)

	if err != nil {
		return nil, err
	}

	err = json.Unmarshal([]byte(data), &rac)

	if err != nil {
		return nil, fmt.Errorf("GetRemoteAccessConnectorById: failed to unmarshal: %v", err)
	}

	return &rac, nil
}

func (ac *AlkiraClient) GetRemoteAccessConnectorConfig(id string) (string, error) {
	uri := fmt.Sprintf("%s/tenantnetworks/%s/alkira-remote-access-connectors/%s", ac.URI, ac.TenantNetworkId, id)

	data, err := ac.get(uri)

	if err != nil {
		return "", err
	}

	return string(data), nil
}
