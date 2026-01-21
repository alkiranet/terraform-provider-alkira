// Copyright (C) 2022-2025 Alkira Inc. All Rights Reserved.

package alkira

import (
	"encoding/json"
	"fmt"
)

type ServiceBluecat struct {
	BddsAnycast                 BluecatAnycast    `json:"bddsAnycast"`
	EdgeAnycast                 BluecatAnycast    `json:"edgeAnycast"`
	BillingTags                 []int             `json:"billingTags"`
	Cxp                         string            `json:"cxp"`
	Description                 string            `json:"description,omitempty"`
	GlobalCidrListId            int               `json:"globalCidrListId"`
	Id                          json.Number       `json:"id,omitempty"` // RESPONSE ONLY
	Instances                   []BluecatInstance `json:"instances"`
	InternalName                string            `json:"internalName,omitempty"`
	LicenseType                 string            `json:"licenseType,omitempty"`
	Name                        string            `json:"name"`
	Segments                    []string          `json:"segments"`
	ServiceGroupId              int               `json:"serviceGroupId,omitempty"`              // RESPONSE ONLY
	ServiceGroupImplicitGroupId int               `json:"serviceGroupImplicitGroupId,omitempty"` // RESPONSE ONLY
	ServiceGroupName            string            `json:"serviceGroupName"`
	Size                        string            `json:"size,omitempty"`
}

type BluecatAnycast struct {
	BackupCxps []string `json:"backupCxps,omitempty"`
	Ips        []string `json:"ips,omitempty"`
}

type BluecatInstance struct {
	BddsOptions      *BDDSOptions `json:"bddsOptions,omitempty"`
	EdgeOptions      *EdgeOptions `json:"edgeOptions,omitempty"`
	Id               json.Number  `json:"id,omitempty"`
	InternalName     string       `json:"internalName,omitempty"`
	Name             string       `json:"name"`
	Type             string       `json:"type"`                       // BDDS or EDGE
	CxpBgpIp         string       `json:"cxpBgpIp,omitempty"`         // RESPONSE ONLY
	CxpBgpAsn        string       `json:"cxpBgpAsn,omitempty"`        // RESPONSE ONLY
	InstanceBgpIp    string       `json:"instanceBgpIp,omitempty"`    // RESPONSE ONLY
	InstanceBgpAsn   string       `json:"instanceBgpAsn,omitempty"`   // RESPONSE ONLY
	ManagementPrefix string       `json:"managementPrefix,omitempty"` // RESPONSE ONLY
	ProductId        string       `json:"productId,omitempty"`        // RESPONSE ONLY
}

type BDDSOptions struct {
	LicenseCredentialId string `json:"licenseCredentialId"`
	HostName            string `json:"hostName"`
	Model               string `json:"model"`
	Version             string `json:"version"`
}

type EdgeOptions struct {
	CredentialId string `json:"credentialId"`
	HostName     string `json:"hostName"`
	Version      string `json:"version"`
}

// NewServiceBluecat new service bluecat
func NewServiceBluecat(ac *AlkiraClient) *AlkiraAPI[ServiceBluecat] {
	uri := fmt.Sprintf("%s/tenantnetworks/%s/bluecat-services", ac.URI, ac.TenantNetworkId)
	api := &AlkiraAPI[ServiceBluecat]{ac, uri, true}
	return api
}
