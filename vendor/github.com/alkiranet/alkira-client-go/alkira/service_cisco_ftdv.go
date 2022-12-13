// Copyright (C) 2022 Alkira Inc. All Rights Reserved.

package alkira

import (
	"encoding/json"
	"fmt"
	"strconv"
)

type CiscoFTDvInstance struct {
	Id           int    `json:"id,omitempty"`           // filled in response
	CredentialId string `json:"credentialId,omitempty"` // filled in response
	InternalName string `json:"internalName,omitempty"` // filled in response
	State        string `json:"state,omitempty"`        // filled in response
	Hostname     string `json:"hostName"`
	LicenseType  string `json:"licenseType"`
	Version      string `json:"version"`
}

type CiscoFTDvManagementServer struct {
	IPAddress string `json:"ipAddress"`
	Segment   string `json:"segment"`
	SegmentId int    `json:"segmentId"`
}

type ServiceCiscoFTDv struct {
	Id               int                       `json:"id,omitempty"` // filled in response
	Name             string                    `json:"name"`
	GlobalCidrListId int                       `json:"globalCidrListId"`
	Size             string                    `json:"size"`
	CredentialId     string                    `json:"credentialId,omitempty"` // filled in response
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
	InternalName     string                    `json:"internalName,omitempty"` // filled in response
	State            string                    `json:"state,omitempty"`        // filled in response
}

// getCiscoFTDvServices get all Cisco FTDv services from the given tenant network
func (ac *AlkiraClient) getCiscoFTDvServices() (string, error) {
	uri := fmt.Sprintf("%s/tenantnetworks/%s/cisco-ftdv-fw-services", ac.URI, ac.TenantNetworkId)

	data, err := ac.get(uri)
	return string(data), err
}

// CreateServiceCiscoFTDv create a Cisco FTDv service
func (ac *AlkiraClient) CreateServiceCiscoFTDv(service *ServiceCiscoFTDv) (string, error) {
	uri := fmt.Sprintf("%s/tenantnetworks/%s/cisco-ftdv-fw-services", ac.URI, ac.TenantNetworkId)

	// Construct the request
	body, err := json.Marshal(service)

	if err != nil {
		return "", fmt.Errorf("CreateServiceCiscoFTDv: failed to marshal: %v", err)
	}

	data, err := ac.create(uri, body, true)

	if err != nil {
		return "", err
	}

	var result ServiceCiscoFTDv
	json.Unmarshal([]byte(data), &result)

	if err != nil {
		return "", fmt.Errorf("CreateServiceCiscoFTDv: failed to unmarshal: %v", err)
	}

	return strconv.Itoa(result.Id), nil
}

// GetServiceCiscoFTDv get one Cisco FTDv service by ID
func (ac *AlkiraClient) GetServiceCiscoFTDv(id string) (*ServiceCiscoFTDv, error) {
	uri := fmt.Sprintf("%s/tenantnetworks/%s/cisco-ftdv-fw-services/%s", ac.URI, ac.TenantNetworkId, id)

	data, err := ac.get(uri)

	if err != nil {
		return nil, err
	}

	var result ServiceCiscoFTDv
	err = json.Unmarshal([]byte(data), &result)

	if err != nil {
		return nil, fmt.Errorf("GetServiceCiscoFTDv: failed to unmarshal: %v", err)
	}

	return &result, nil
}

// DeleteServiceCiscoFTDv delete the given Cisco FTDv service by ID
func (ac *AlkiraClient) DeleteServiceCiscoFTDv(id string) error {
	uri := fmt.Sprintf("%s/tenantnetworks/%s/cisco-ftdv-fw-services/%s", ac.URI, ac.TenantNetworkId, id)

	return ac.delete(uri, true)
}

// UpdateServiceCiscoFTDv update a Cisco FTDv service by ID
func (ac *AlkiraClient) UpdateServiceCiscoFTDv(id string, service *ServiceCiscoFTDv) error {
	uri := fmt.Sprintf("%s/v1/tenantnetworks/%s/cisco-ftdv-fw-services/%s", ac.URI, ac.TenantNetworkId, id)

	// Construct the request
	body, err := json.Marshal(service)

	if err != nil {
		return fmt.Errorf("UpdateServiceCiscoFTDv: failed to marshal: %v", err)
	}

	return ac.update(uri, body, true)
}

// GetServiceCiscoFTDvByName get a Cisco FTDv service by name
func (ac *AlkiraClient) GetServiceCiscoFTDvByName(name string) (ServiceCiscoFTDv, error) {
	var ciscoFTDvService ServiceCiscoFTDv

	if len(name) == 0 {
		return ciscoFTDvService, fmt.Errorf("GetServiceCiscoFTDvByName: Invalid Service name")
	}

	ciscoFTDvServices, err := ac.getCiscoFTDvServices()

	if err != nil {
		return ciscoFTDvService, err
	}

	var result []ServiceCiscoFTDv
	json.Unmarshal([]byte(ciscoFTDvServices), &result)

	for _, l := range result {
		if l.Name == name {
			return l, nil
		}
	}

	return ciscoFTDvService, fmt.Errorf("GetServiceCiscoFTDvByName: failed to find the service by %s", name)
}
