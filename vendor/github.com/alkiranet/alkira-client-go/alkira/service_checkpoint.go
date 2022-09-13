// Copyright (C) 2022 Alkira Inc. All Rights Reserved.

package alkira

import (
	"encoding/json"
	"fmt"
)

type Checkpoint struct {
	AutoScale        string                      `json:"autoScale"`
	BillingTags      []int                       `json:"billingTags"`
	Cxp              string                      `json:"cxp"`
	CredentialId     string                      `json:"credentialId"`
	Description      string                      `json:"description"`
	Id               json.Number                 `json:"id,omitempty"` //filled only on response
	Instances        []CheckpointInstance        `json:"instances"`
	InternalName     string                      `json:"internalName"`
	LicenseType      string                      `json:"licenseType"`
	ManagementServer *CheckpointManagementServer `json:"managementServer"`
	MaxInstanceCount int                         `json:"maxInstanceCount"`
	MinInstanceCount int                         `json:"minInstanceCount"`
	Name             string                      `json:"name"`
	PdpIps           []string                    `json:"pdpIps,omitempty"`
	Segments         []string                    `json:"segments"`
	SegmentOptions   SegmentNameToZone           `json:"segmentOptions"`
	Size             string                      `json:"size"`
	TunnelProtocol   string                      `json:"tunnelProtocol"`
	Version          string                      `json:"version"`
}

type CheckpointInstance struct {
	Id           int    `json:"id,omitempty"` //filled only on response
	Name         string `json:"name"`
	CredentialId string `json:"credentialId"`
	InternalName string `json:"internalName,omitempty"` //filled only on response
}

type CheckpointInstanceConfig struct {
	Data string //The response is string data the entire body of the response whould be interpreted together. There is no json structure.
}

type CheckpointManagementServer struct {
	ConfigurationMode string   `json:"configurationMode"`
	CredentialId      string   `json:"credentialId"`
	Domain            string   `json:"domain"`
	GlobalCidrListId  int      `json:"globalCidrListId"`
	Ips               []string `json:"ips"`
	Reachability      string   `json:"reachability"`
	Segment           string   `json:"segment"`
	SegmentId         int      `json:"segmentId"`
	Type              string   `json:"type"`
	UserName          string   `json:"userName"`
}

func (ac *AlkiraClient) CreateCheckpoint(c *Checkpoint) (string, error) {
	uri := fmt.Sprintf("%s/tenantnetworks/%s/chkp-fw-services", ac.URI, ac.TenantNetworkId)

	body, err := json.Marshal(c)

	if err != nil {
		return "", fmt.Errorf("CreateCheckpoint: marshal failed: %v", err)
	}

	data, err := ac.create(uri, body, true)

	if err != nil {
		return "", err
	}

	var result Checkpoint
	err = json.Unmarshal([]byte(data), &result)

	if err != nil {
		return "", fmt.Errorf("CreateCheckpoint: failed to unmarshal: %v", err)
	}

	return string(result.Id), nil
}

func (ac *AlkiraClient) GetCheckpoints() (string, error) {
	uri := fmt.Sprintf("%s/tenantnetworks/%s/chkp-fw-services", ac.URI, ac.TenantNetworkId)
	data, err := ac.get(uri)

	if err != nil {
		return "", err
	}

	return string(data), nil
}

func (ac *AlkiraClient) GetCheckpointById(id string) (*Checkpoint, error) {
	uri := fmt.Sprintf("%s/tenantnetworks/%s/chkp-fw-services/%s", ac.URI, ac.TenantNetworkId, id)

	var checkpoint Checkpoint

	data, err := ac.get(uri)

	if err != nil {
		return nil, err
	}

	err = json.Unmarshal([]byte(data), &checkpoint)

	if err != nil {
		return nil, fmt.Errorf("GetCheckpointById: failed to unmarshal: %v", err)
	}

	return &checkpoint, nil
}

func (ac *AlkiraClient) GetCheckpointInstanceConfig(serviceId string, instanceId string) (*CheckpointInstanceConfig, error) {
	uri := fmt.Sprintf("%s/tenantnetworks/%s/chkp-fw-services/%s/instances/%s/configuration", ac.URI, ac.TenantNetworkId, serviceId, instanceId)

	data, err := ac.get(uri)

	if err != nil {
		return nil, err
	}

	//The response is string data. The entire body of the response should be interpreted together.
	//There is no json structure for CheckpointInstanceConfig.
	return &CheckpointInstanceConfig{Data: string(data)}, nil
}

func (ac *AlkiraClient) UpdateCheckpoint(id string, c *Checkpoint) error {
	uri := fmt.Sprintf("%s/tenantnetworks/%s/chkp-fw-services/%s", ac.URI, ac.TenantNetworkId, id)

	body, err := json.Marshal(c)

	if err != nil {
		return fmt.Errorf("UpdateCheckpoint: failed to marshal: %v", err)
	}

	return ac.update(uri, body, true)
}

func (ac *AlkiraClient) DeleteCheckpoint(id string) error {
	uri := fmt.Sprintf("%s/tenantnetworks/%s/chkp-fw-services/%s", ac.URI, ac.TenantNetworkId, id)

	return ac.delete(uri, true)
}
