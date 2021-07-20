// Copyright (C) 2020-2021 Alkira Inc. All Rights Reserved.

package alkira

import (
	"encoding/json"
	"fmt"
	"strconv"
)

type ServicePan struct {
	BillingTagIds       []int                `json:"billingTags"`
	CXP                 string               `json:"cxp"`
	CredentialId        string               `json:"credentialId"`
	Id                  int                  `json:"id,omitempty"`
	Instances           []ServicePanInstance `json:"instances,omitempty"`
	LicenseType         string               `json:"licenseType"`
	ManagementSegmentId int                  `json:"managementSegment"`
	MaxInstanceCount    int                  `json:"maxInstanceCount"`
	MinInstanceCount    int                  `json:"minInstanceCount"`
	Name                string               `json:"name"`
	PanoramaEnabled     bool                 `json:"panoramaEnabled"`
	PanoramaDeviceGroup string               `json:"panoramaDeviceGroup"`
	PanoramaIpAddress   string               `json:"panoramaIPAddress"`
	PanoramaTemplate    string               `json:"panoramaTemplate"`
	PanWarmBootEnabled  bool                 `json:"panWarmBootEnabled,omitempty"`
	SegmentIds          []int                `json:"segments"`
	SegmentOptions      interface{}          `json:"segmentOptions,omitempty"`
	Size                string               `json:"size"`
	TunnelProtocol      string               `json:"tunnelProtocol,omitempty"`
	Type                string               `json:"type"`
	Version             string               `json:"version"`
}

type ServicePanInstance struct {
	CredentialId string `json:"credentialId"`
	Name         string `json:"name"`
}

// CreateServicePan create service PAN
func (ac *AlkiraClient) CreateServicePan(service *ServicePan) (string, error) {
	uri := fmt.Sprintf("%s/v1/tenantnetworks/%s/panfwservices", ac.URI, ac.TenantNetworkId)

	// Construct the request
	body, err := json.Marshal(service)

	if err != nil {
		return "", fmt.Errorf("CreateServicePan: marshal failed: %v", err)
	}

	data, err := ac.create(uri, body)

	if err != nil {
		return "", err
	}

	var result ServicePan
	err = json.Unmarshal([]byte(data), &result)

	if err != nil {
		return "", fmt.Errorf("CreateServicePan: failed to unmarshal: %v", err)
	}

	return strconv.Itoa(result.Id), nil
}

// DeleteServicePan delete a Service PAN by Id
func (ac *AlkiraClient) DeleteServicePan(id string) error {
	uri := fmt.Sprintf("%s/v1/tenantnetworks/%s/panfwservices/%s", ac.URI, ac.TenantNetworkId, id)

	return ac.delete(uri)
}

// UpdateServicePan Update a Service PAN by Id
func (ac *AlkiraClient) UpdateServicePan(id string, service *ServicePan) error {
	uri := fmt.Sprintf("%s/v1/tenantnetworks/%s/panfwservices/%s", ac.URI, ac.TenantNetworkId, id)

	body, err := json.Marshal(service)

	if err != nil {
		return fmt.Errorf("UpdateServicePan: failed to marshal request: %v", err)
	}

	return ac.update(uri, body)
}

// GetServicePanById get an service-pan by Id
func (ac *AlkiraClient) GetServicePanById(id string) (*ServicePan, error) {
	uri := fmt.Sprintf("%s/tenantnetworks/%s/panfwservices/%s", ac.URI, ac.TenantNetworkId, id)

	var service ServicePan

	data, err := ac.get(uri)

	if err != nil {
		return nil, err
	}

	err = json.Unmarshal([]byte(data), &service)

	if err != nil {
		return nil, fmt.Errorf("GetServicePan: failed to unmarshal: %v", err)
	}

	return &service, nil
}
