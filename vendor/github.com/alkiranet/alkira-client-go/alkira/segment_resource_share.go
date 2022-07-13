// Copyright (C) 2022 Alkira Inc. All Rights Reserved.

package alkira

import (
	"encoding/json"
	"fmt"
)

type SegmentResourceShare struct {
	Id                  json.Number `json:"id,omitempty"` // response only
	Name                string      `json:"name"`
	ServiceList         []int       `json:"serviceList"`
	DesignatedSegment   string      `json:"designatedSegment"`
	EndAResources       []int       `json:"endAResources"`
	EndBResources       []int       `json:"endBResources"`
	EndARouteLimit      int         `json:"endARouteLimit"`
	EndBRouteLimit      int         `json:"endBRouteLimit"`
	Direction           string      `json:"direction"`
	DesignatedSegmentId int         `json:"designatedSegmentId,omitempty"` // response only
	EndASegmentId       int         `json:"endASegmentId,omitempty"`       // response only
	EndBSegmentId       int         `json:"endBSegmentId,omitempty"`       // response only
}

// CreateSegmentResourceShare create a new segment resource
func (ac *AlkiraClient) CreateSegmentResourceShare(resource *SegmentResourceShare) (string, error) {
	uri := fmt.Sprintf("%s/tenantnetworks/%s/segment-resource-shares", ac.URI, ac.TenantNetworkId)

	body, err := json.Marshal(resource)

	if err != nil {
		return "", fmt.Errorf("CreateSegmentResourceShare: failed to marshal: %v", err)
	}

	data, err := ac.create(uri, body, true)

	if err != nil {
		return "", err
	}

	var result SegmentResourceShare
	err = json.Unmarshal([]byte(data), &result)

	if err != nil {
		return "", fmt.Errorf("CreateSegmentResourceShare: failed to unmarshal: %v", err)
	}

	return string(result.Id), nil
}

// DeleteSegmentResourceShare delete a segment resource share by ID
func (ac *AlkiraClient) DeleteSegmentResourceShare(id string) error {
	uri := fmt.Sprintf("%s/tenantnetworks/%s/segment-resource-shares/%s", ac.URI, ac.TenantNetworkId, id)
	return ac.delete(uri, true)
}

// UpdateSegmentResourceShare update a segment resource share by ID
func (ac *AlkiraClient) UpdateSegmentResourceShare(id string, resource *SegmentResourceShare) error {
	uri := fmt.Sprintf("%s/tenantnetworks/%s/segment-resource-shares/%s", ac.URI, ac.TenantNetworkId, id)

	body, err := json.Marshal(resource)

	if err != nil {
		return fmt.Errorf("UpdateSegmentResourceShare: failed to marshal: %v", err)
	}

	return ac.update(uri, body, true)
}

// GetSegmentResourceShares get all segment resource shares from the given tenant network
func (ac *AlkiraClient) GetSegmentResourceShares() (string, error) {
	uri := fmt.Sprintf("%s/tenantnetworks/%s/segment-resource-shares", ac.URI, ac.TenantNetworkId)
	data, err := ac.get(uri)
	return string(data), err
}

// GetSegmentResourceShareById get a single segment resource share by ID
func (ac *AlkiraClient) GetSegmentResourceShareById(id string) (*SegmentResourceShare, error) {
	uri := fmt.Sprintf("%s/tenantnetworks/%s/segment-resource-shares/%s", ac.URI, ac.TenantNetworkId, id)

	var share SegmentResourceShare

	data, err := ac.get(uri)

	if err != nil {
		return nil, err
	}

	err = json.Unmarshal([]byte(data), &share)

	if err != nil {
		return nil, fmt.Errorf("GetSegmentResourceShareById: failed to unmarshal: %v", err)
	}

	return &share, nil
}

// GetSegmentResourceShareByName get a segment resource share by its name
func (ac *AlkiraClient) GetSegmentResourceShareByName(name string) (*SegmentResourceShare, error) {

	if len(name) == 0 {
		return nil, fmt.Errorf("Invalid segmentResource name input")
	}

	shares, err := ac.GetSegmentResourceShares()

	if err != nil {
		return nil, err
	}

	var result []SegmentResourceShare
	json.Unmarshal([]byte(shares), &result)

	for _, t := range result {
		if t.Name == name {
			return &t, nil
		}
	}

	return nil, fmt.Errorf("failed to find the segmentResource by %s", name)
}
