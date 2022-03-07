// Copyright (C) 2022 Alkira Inc. All Rights Reserved.

package alkira

import (
	"encoding/json"
	"fmt"
	"strconv"
)

type SegmentResource struct {
	Id            int                          `json:"id"`
	Name          string                       `json:"name"`
	Segment       string                       `json:"segment"`
	GroupPrefixes []SegmentResourceGroupPrefix `json:"groupPrefixes"`
}

type SegmentResourceGroupPrefix struct {
	GroupId      int `json:"groupId"`
	PrefixListId int `json:"prefixListId"`
}

// CreateSegmentResource create a new segment resource
func (ac *AlkiraClient) CreateSegmentResource(resource *SegmentResource) (string, error) {
	uri := fmt.Sprintf("%s/tenantnetworks/%s/segment-resources", ac.URI, ac.TenantNetworkId)

	body, err := json.Marshal(resource)

	if err != nil {
		return "", fmt.Errorf("CreateSegmentResource: failed to marshal: %v", err)
	}

	data, err := ac.create(uri, body)

	if err != nil {
		return "", err
	}

	var result SegmentResource
	err = json.Unmarshal([]byte(data), &result)

	if err != nil {
		return "", fmt.Errorf("CreateSegmentResource: failed to unmarshal: %v", err)
	}

	return strconv.Itoa(result.Id), nil
}

// DeleteSegmentResource delete a segment resource by ID
func (ac *AlkiraClient) DeleteSegmentResource(id string) error {
	uri := fmt.Sprintf("%s/tenantnetworks/%s/segment-resources/%s", ac.URI, ac.TenantNetworkId, id)
	return ac.delete(uri)
}

// UpdateSegmentResource update a segment resource by ID
func (ac *AlkiraClient) UpdateSegmentResource(id string, resource *SegmentResource) error {
	uri := fmt.Sprintf("%s/tenantnetworks/%s/segment-resources/%s", ac.URI, ac.TenantNetworkId, id)

	body, err := json.Marshal(resource)

	if err != nil {
		return fmt.Errorf("UpdateSegmentResource: failed to marshal: %v", err)
	}

	return ac.update(uri, body)
}

// GetSegmentResources get all segment resources from the given tenant network
func (ac *AlkiraClient) GetSegmentResources() (string, error) {
	uri := fmt.Sprintf("%s/tenantnetworks/%s/segment-resources", ac.URI, ac.TenantNetworkId)
	data, err := ac.get(uri)
	return string(data), err
}

// GetSegmentResourceById get a single segment resource by ID
func (ac *AlkiraClient) GetSegmentResourceById(id string) (SegmentResource, error) {
	uri := fmt.Sprintf("%s/tenantnetworks/%s/segment-resources", ac.URI, ac.TenantNetworkId)

	var segmentResource SegmentResource

	data, err := ac.get(uri)

	if err != nil {
		return segmentResource, err
	}

	err = json.Unmarshal([]byte(data), &segmentResource)

	if err != nil {
		return segmentResource, fmt.Errorf("GetSegmentResourceById: failed to unmarshal: %v", err)
	}

	return segmentResource, nil
}

// GetSegmentResourceByName get a segment resource by its name
func (ac *AlkiraClient) GetSegmentResourceByName(name string) (SegmentResource, error) {
	var segmentResource SegmentResource

	if len(name) == 0 {
		return segmentResource, fmt.Errorf("Invalid segmentResource name input")
	}

	segmentResources, err := ac.GetSegmentResources()

	if err != nil {
		return segmentResource, err
	}

	var result []SegmentResource
	json.Unmarshal([]byte(segmentResources), &result)

	for _, t := range result {
		if t.Name == name {
			return t, nil
		}
	}

	return segmentResource, fmt.Errorf("failed to find the segmentResource by %s", name)
}
