// Copyright (C) 2022-2023 Alkira Inc. All Rights Reserved.

package alkira

import (
	"encoding/json"
	"fmt"
)

type ServiceCheckpoint struct {
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
	Data string //The response is string data the entire body of the
	//response whould be interpreted together. There is no
	//json structure.
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

// NewServiceCheckpoint new service checkpoint
func NewServiceCheckpoint(ac *AlkiraClient) *AlkiraAPI[ServiceCheckpoint] {
	uri := fmt.Sprintf("%s/tenantnetworks/%s/chkp-fw-services", ac.URI, ac.TenantNetworkId)
	api := &AlkiraAPI[ServiceCheckpoint]{ac, uri, true}
	return api
}
