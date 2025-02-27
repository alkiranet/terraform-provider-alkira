// Copyright (C) 2022-2025 Alkira Inc. All Rights Reserved.

package alkira

import (
	"encoding/json"
	"fmt"
)

type ServiceInfoblox struct {
	AnyCast                     InfobloxAnycast    `json:"anycast"`
	BillingTags                 []int              `json:"billingTags"`
	Cxp                         string             `json:"cxp"`
	Description                 string             `json:"description,omitempty"`
	GlobalCidrListId            int                `json:"globalCidrListId"`
	GridMaster                  InfobloxGridMaster `json:"gridMaster"`
	Id                          json.Number        `json:"id,omitempty"` // RESPONSE ONLY
	Instances                   []InfobloxInstance `json:"instances"`
	InternalName                string             `json:"internalName,omitempty"`
	LicenseType                 string             `json:"licenseType,omitempty"`
	Name                        string             `json:"name"`
	Segments                    []string           `json:"segments"`
	ServiceGroupId              int                `json:"serviceGroupId,omitempty"`              // RESPONSE ONLY
	ServiceGroupImplicitGroupId int                `json:"serviceGroupImplicitGroupId,omitempty"` // RESPONSE ONLY
	ServiceGroupName            string             `json:"serviceGroupName"`
	Size                        string             `json:"size,omitempty"`
	AllowListId                 int                `json:"allowListId,omitempty"`
}

type InfobloxAnycast struct {
	BackupCxps []string `json:"backupCxps,omitempty"`
	Enabled    bool     `json:"enabled"`
	Ips        []string `json:"ips,omitempty"`
}

type InfobloxGridMaster struct {
	External                 bool   `json:"external,omitempty"`
	GridMasterCredentialId   string `json:"gridMasterCredentialId"`
	Ip                       string `json:"ip,omitempty"`
	Name                     string `json:"name"`
	SharedSecretCredentialId string `json:"sharedSecretCredentialId"`
}

type InfobloxInstance struct {
	AnyCastEnabled     bool        `json:"anyCastEnabled"`
	ConfiguredMasterIp string      `json:"configuredMasterIp,omitempty"`
	CredentialId       string      `json:"credentialId"`
	HostName           string      `json:"hostName"`
	Id                 json.Number `json:"id,omitempty"`
	InternalName       string      `json:"internalName,omitempty"`
	LanPrefix          string      `json:"lanPrefix,omitempty"`
	ManagementPrefix   string      `json:"managementPrefix,omitempty"`
	Model              string      `json:"model"`
	Name               string      `json:"name,omitempty"`
	ProductId          string      `json:"productId,omitempty"`
	PublicIp           string      `json:"publicIp,omitempty"`
	Type               string      `json:"type"`
	Version            string      `json:"version,omitempty"`
}

// NewServiceInfoblox new service infoblox
func NewServiceInfoblox(ac *AlkiraClient) *AlkiraAPI[ServiceInfoblox] {
	uri := fmt.Sprintf("%s/tenantnetworks/%s/infoblox-services", ac.URI, ac.TenantNetworkId)
	api := &AlkiraAPI[ServiceInfoblox]{ac, uri, true}
	return api
}
