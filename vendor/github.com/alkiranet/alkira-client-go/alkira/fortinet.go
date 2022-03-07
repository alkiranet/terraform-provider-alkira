// Copyright (C) 2022 Alkira Inc. All Rights Reserved.

package alkira

import (
	"encoding/json"
	"fmt"
	"strconv"
)

type Fortinet struct {
	AutoScale        string                          `json:"autoScale,omitempty"`
	BillingTags      []int                           `json:"billingTags"`
	CredentialId     string                          `json:"credentialId"`
	Cxp              string                          `json:"cxp"`
	Id               int                             `json:"id"`
	Instances        []FortinetInstance              `json:"instances"`
	InternalName     string                          `json:"internalName"`
	LicenseType      string                          `json:"licenseType"`
	ManagementServer *FortinetManagmentServer        `json:"managementServer"`
	MaxInstanceCount int                             `json:"maxInstanceCount"`
	MinInstanceCount int                             `json:"minInstanceCount"`
	Name             string                          `json:"name"`
	Segments         []string                        `json:"segments"`
	SegmentOptions   map[string]*FortinetSegmentName `json:"segmentOptions"`
	Size             string                          `json:"size"`
	State            string                          `json:"state,omitempty"`
	TunnelProtocol   string                          `json:"tunnelProtocol"`
	Version          string                          `json:"version"`
}

type FortinetInstance struct {
	Name         string `json:"name"`
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

type FortinetSegmentName struct {
	SegmentId     int                 `json:"segmentId"`
	ZonesToGroups map[string][]string `json:"zonesToGroups"`
}

func (ac *AlkiraClient) CreateFortinet(f *Fortinet) (string, error) {
	uri := fmt.Sprintf("%s/tenantnetworks/%s/ftnt-fw-services", ac.URI, ac.TenantNetworkId)

	body, err := json.Marshal(f)

	if err != nil {
		return "", fmt.Errorf("CreateFortinet: marshal failed: %v", err)
	}

	data, err := ac.create(uri, body)

	if err != nil {
		return "", err
	}

	var result Fortinet
	err = json.Unmarshal([]byte(data), &result)

	if err != nil {
		return "", fmt.Errorf("CreateFortinet: failed to unmarshal: %v", err)
	}

	return strconv.Itoa(result.Id), nil
}

func (ac *AlkiraClient) GetFortinets() (string, error) {
	uri := fmt.Sprintf("%s/tenantnetworks/%s/ftnt-fw-services", ac.URI, ac.TenantNetworkId)
	data, err := ac.get(uri)

	if err != nil {
		return "", err
	}

	return string(data), nil
}

func (ac *AlkiraClient) GetFortinetById(id string) (*Fortinet, error) {

	uri := fmt.Sprintf("%s/tenantnetworks/%s/ftnt-fw-services/%s", ac.URI, ac.TenantNetworkId, id)

	var fortinet Fortinet

	data, err := ac.get(uri)

	if err != nil {
		return nil, err
	}

	err = json.Unmarshal([]byte(data), &fortinet)

	if err != nil {
		return nil, fmt.Errorf("GetFortinetById: failed to unmarshal: %v", err)
	}

	return &fortinet, nil
}

func (ac *AlkiraClient) GetFortinetInstanceConfig(serviceId string, instanceId string) (*FortinetInstanceConfig, error) {
	uri := fmt.Sprintf("%s/tenantnetworks/%s/ftnt-fw-services/%s/instances/%s/configuration", ac.URI, ac.TenantNetworkId, serviceId, instanceId)

	var fortinet FortinetInstanceConfig

	data, err := ac.get(uri)

	if err != nil {
		return nil, err
	}

	err = json.Unmarshal([]byte(data), &fortinet)

	if err != nil {
		return nil, fmt.Errorf("GetFortinetInstanceConfig: failed to unmarshal: %v", err)
	}

	return &fortinet, nil
}

func (ac *AlkiraClient) UpdateFortinet(id string, f *Fortinet) error {

	uri := fmt.Sprintf("%s/tenantnetworks/%s/ftnt-fw-services/%s", ac.URI, ac.TenantNetworkId, id)

	body, err := json.Marshal(f)

	if err != nil {
		return fmt.Errorf("UpdateFortinet: failed to marshal request: %v", err)
	}

	return ac.update(uri, body)
}

func (ac *AlkiraClient) DeleteFortinet(id string) error {
	uri := fmt.Sprintf("%s/tenantnetworks/%s/ftnt-fw-services/%s", ac.URI, ac.TenantNetworkId, id)

	return ac.delete(uri)
}
