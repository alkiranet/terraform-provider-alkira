// Copyright (C) 2022-2025 Alkira Inc. All Rights Reserved.

package alkira

import (
	"encoding/json"
	"fmt"
)

type ServiceFortinet struct {
	AutoScale        string                   `json:"autoScale,omitempty"`
	BillingTags      []int                    `json:"billingTags"`
	CredentialId     string                   `json:"credentialId"`
	Cxp              string                   `json:"cxp"`
	Id               json.Number              `json:"id"`
	Instances        []FortinetInstance       `json:"instances"`
	InternalName     string                   `json:"internalName"`
	LicenseType      string                   `json:"licenseType"`
	ManagementServer *FortinetManagmentServer `json:"managementServer"`
	MaxInstanceCount int                      `json:"maxInstanceCount"`
	MinInstanceCount int                      `json:"minInstanceCount"`
	Name             string                   `json:"name"`
	Scheme           string                   `json:"scheme,omitempty"`
	Segments         []string                 `json:"segments"`
	SegmentOptions   SegmentNameToZone        `json:"segmentOptions"`
	Size             string                   `json:"size"`
	State            string                   `json:"state,omitempty"`
	TunnelProtocol   string                   `json:"tunnelProtocol"`
	Version          string                   `json:"version"`
	Description      string                   `json:"description,omitempty"`
}

type FortinetInstance struct {
	Name         string `json:"name"`
	Id           int    `json:"id,omitempty"`
	HostName     string `json:"hostName"`
	SerialNumber string `json:"serialNumber"`
	CredentialId string `json:"credentialId"`
}

type FortinetManagmentServer struct {
	IpAddress string `json:"ipAddress"`
	Segment   string `json:"segment"`
}

type FortinetInstanceConfig struct {
	ManagementIp string `json:"managementIP"`
	SerialNumber string `json:"serialNumber"`
}

// NewServiceFortinet new service fortinet
func NewServiceFortinet(ac *AlkiraClient) *AlkiraAPI[ServiceFortinet] {
	uri := fmt.Sprintf("%s/tenantnetworks/%s/ftnt-fw-services", ac.URI, ac.TenantNetworkId)
	api := &AlkiraAPI[ServiceFortinet]{ac, uri, true}
	return api
}
