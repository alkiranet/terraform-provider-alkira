// Copyright (C) 2020-2025 Alkira Inc. All Rights Reserved.

package alkira

import (
	"encoding/json"
	"fmt"
)

type ServicePan struct {
	BillingTagIds               []int                                `json:"billingTags"`
	Bundle                      string                               `json:"bundle,omitempty"`
	CXP                         string                               `json:"cxp"`
	CredentialId                string                               `json:"credentialId"`
	GlobalProtectEnabled        bool                                 `json:"globalProtectEnabled"`
	GlobalProtectSegmentOptions map[string]*GlobalProtectSegmentName `json:"globalProtectSegmentOptions,omitempty"`
	Id                          json.Number                          `json:"id,omitempty"`
	Instances                   []ServicePanInstance                 `json:"instances,omitempty"`
	LicenseType                 string                               `json:"licenseType"`
	LicenseKey                  string                               `json:"licenseKey"`
	ManagementSegmentId         int                                  `json:"managementSegment"`
	MaxInstanceCount            int                                  `json:"maxInstanceCount"`
	MinInstanceCount            int                                  `json:"minInstanceCount"`
	MasterKeyEnabled            bool                                 `json:"masterKeyEnabled,omitempty"`
	MasterKeyCredentialId       string                               `json:"masterKeyCredentialId,omitempty"`
	Name                        string                               `json:"name"`
	PanoramaEnabled             bool                                 `json:"panoramaEnabled"`
	PanoramaDeviceGroup         *string                              `json:"panoramaDeviceGroup,omitempty"`
	PanoramaIpAddress           *string                              `json:"panoramaIPAddress,omitempty"`
	PanoramaIpAddresses         []string                             `json:"panoramaIPAddresses,omitempty"`
	PanoramaTemplate            *string                              `json:"panoramaTemplate,omitempty"`
	PanWarmBootEnabled          bool                                 `json:"panWarmBootEnabled,omitempty"`
	RegistrationCredentialId    string                               `json:"registrationCredentialId,omitempty"`
	SegmentIds                  []int                                `json:"segments"`
	SegmentOptions              SegmentNameToZone                    `json:"segmentOptions,omitempty"`
	Size                        string                               `json:"size"`
	SubLicenseType              string                               `json:"subLicenseType,omitempty"`
	TunnelProtocol              string                               `json:"tunnelProtocol,omitempty"`
	Type                        string                               `json:"type"`
	Version                     string                               `json:"version"`
	Description                 string                               `json:"description,omitempty"`
}

type GlobalProtectSegmentOptions struct {
	SegmentName *GlobalProtectSegmentName `json:"segmentName"`
}

type GlobalProtectSegmentName struct {
	RemoteUserZoneName string `json:"remoteUserZoneName"`
	PortalFqdnPrefix   string `json:"portalFqdnPrefix"`
	ServiceGroupName   string `json:"serviceGroupName"`
}

type GlobalProtectSegmentOptionsInstance struct {
	SegmentName *GlobalProtectSegmentNameInstance `json:"segmentName"`
}

type GlobalProtectSegmentNameInstance struct {
	PortalEnabled  bool `json:"portalEnabled"`
	GatewayEnabled bool `json:"gatewayEnabled"`
	PrefixListId   int  `json:"prefixListId"`
}

type ServicePanInstance struct {
	CredentialId                string                                       `json:"credentialId"`
	GlobalProtectSegmentOptions map[string]*GlobalProtectSegmentNameInstance `json:"globalProtectSegmentOptions,omitempty"`
	Id                          int                                          `json:"id,omitempty"`
	MasterKeyEnabled            bool                                         `json:"masterKeyEnabled,omitempty"`
	Name                        string                                       `json:"name"`
	TrafficEnabled              bool                                         `json:"trafficEnabled"`
}

// NewServicePan new service pan
func NewServicePan(ac *AlkiraClient) *AlkiraAPI[ServicePan] {
	uri := fmt.Sprintf("%s/v1/tenantnetworks/%s/panfwservices", ac.URI, ac.TenantNetworkId)
	api := &AlkiraAPI[ServicePan]{ac, uri, true}
	return api
}
