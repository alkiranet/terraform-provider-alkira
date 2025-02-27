// Copyright (C) 2022-2025 Alkira Inc. All Rights Reserved.

package alkira

import (
	"encoding/json"
	"fmt"
)

type CiscoFTDvInstance struct {
	Id             int    `json:"id,omitempty"`           // response only
	CredentialId   string `json:"credentialId,omitempty"` // response only
	InternalName   string `json:"internalName,omitempty"` // response only
	State          string `json:"state,omitempty"`        // response only
	Hostname       string `json:"hostName"`
	LicenseType    string `json:"licenseType"`
	Version        string `json:"version"`
	TrafficEnabled bool   `json:"trafficEnabled"`
}

type CiscoFTDvManagementServer struct {
	IPAddress string `json:"ipAddress"`
	Segment   string `json:"segment"`
	SegmentId int    `json:"segmentId"`
}

type ServiceCiscoFTDv struct {
	Id               json.Number               `json:"id,omitempty"` // response only
	Name             string                    `json:"name"`
	GlobalCidrListId int                       `json:"globalCidrListId"`
	Size             string                    `json:"size"`
	CredentialId     string                    `json:"credentialId,omitempty"` // response only
	Cxp              string                    `json:"cxp"`
	ManagementServer CiscoFTDvManagementServer `json:"managementServer"`
	IpAllowList      []string                  `json:"servicesIpAllowList"`
	MaxInstanceCount int                       `json:"maxInstanceCount"`
	MinInstanceCount int                       `json:"minInstanceCount"`
	Segments         []string                  `json:"segments"`
	SegmentOptions   SegmentNameToZone         `json:"segmentOptions,omitempty"`
	Instances        []CiscoFTDvInstance       `json:"instances"`
	BillingTags      []int                     `json:"billingTags"`
	TunnelProtocol   string                    `json:"tunnelProtocol"`
	AutoScale        string                    `json:"autoScale"`
	InternalName     string                    `json:"internalName,omitempty"` // response only
	State            string                    `json:"state,omitempty"`        // response only
	Description      string                    `json:"description,omitempty"`
}

// NewServiceCiscoFTDv new service cisco FTDv
func NewServiceCiscoFTDv(ac *AlkiraClient) *AlkiraAPI[ServiceCiscoFTDv] {
	uri := fmt.Sprintf("%s/tenantnetworks/%s/cisco-ftdv-fw-services", ac.URI, ac.TenantNetworkId)
	api := &AlkiraAPI[ServiceCiscoFTDv]{ac, uri, true}
	return api
}
