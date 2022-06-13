// Copyright (C) 2020-2022 Alkira Inc. All Rights Reserved.

package alkira

import (
	"encoding/json"
	"fmt"
	"strconv"
)

type Segment struct {
	Asn                                        int      `json:"asn"`
	EnterpriseDNSServerIP                      string   `json:"enterpriseDNSServerIP,omitempty"`
	Id                                         int      `json:"id,omitempty"` // only for response
	IpBlock                                    string   `json:"ipBlock"`
	IpBlocks                                   []string `json:"ipBlocks,omitempty"`
	Name                                       string   `json:"name"`
	ReservePublicIPsForUserAndSiteConnectivity bool     `json:"reservePublicIPsForUserAndSiteConnectivity,omitempty"`
}

// Get all segments from the given tenant network
func (ac *AlkiraClient) GetSegments() (string, error) {
	uri := fmt.Sprintf("%s/tenantnetworks/%s/segments", ac.URI, ac.TenantNetworkId)

	data, err := ac.get(uri)
	return string(data), err
}

// GetSegment get single segment by ID
func (ac *AlkiraClient) GetSegmentById(id string) (Segment, error) {
	uri := fmt.Sprintf("%s/tenantnetworks/%s/segments/%s", ac.URI, ac.TenantNetworkId, id)

	var segment Segment

	data, err := ac.get(uri)

	if err != nil {
		return segment, err
	}

	err = json.Unmarshal([]byte(data), &segment)

	if err != nil {
		return segment, fmt.Errorf("GetSegmentById: failed to unmarshal: %v", err)
	}

	return segment, nil
}

// GetSegmentByName get the segment by its name
func (ac *AlkiraClient) GetSegmentByName(name string) (Segment, error) {
	var segment Segment

	if len(name) == 0 {
		return segment, fmt.Errorf("Invalid segment name input")
	}

	segments, err := ac.GetSegments()

	if err != nil {
		return segment, err
	}

	var result []Segment
	json.Unmarshal([]byte(segments), &result)

	for _, g := range result {
		if g.Name == name {
			return g, nil
		}
	}

	return segment, fmt.Errorf("failed to find the segment by %s", name)
}

// CreateSegment create a new Segment
func (ac *AlkiraClient) CreateSegment(segment *Segment) (string, error) {

	uri := fmt.Sprintf("%s/tenantnetworks/%s/segments", ac.URI, ac.TenantNetworkId)

	// Construct the request
	body, err := json.Marshal(segment)

	if err != nil {
		return "", fmt.Errorf("CreateSegment: failed to marshal: %v", err)
	}

	data, err := ac.create(uri, body)

	if err != nil {
		return "", err
	}

	var result Segment
	err = json.Unmarshal([]byte(data), &result)

	if err != nil {
		return "", fmt.Errorf("CreateSegment: failed to unmarshal: %v", err)
	}

	return strconv.Itoa(result.Id), nil
}

// DeleteSegment delete a segment by given segment ID
func (ac *AlkiraClient) DeleteSegment(id string) error {
	uri := fmt.Sprintf("%s/tenantnetworks/%s/segments/%s", ac.URI, ac.TenantNetworkId, id)
	return ac.delete(uri)
}

// UpdateSegment update a segment by segment ID
func (ac *AlkiraClient) UpdateSegment(id string, segment *Segment) error {

	uri := fmt.Sprintf("%s/tenantnetworks/%s/segments/%s", ac.URI, ac.TenantNetworkId, id)

	// Construct the request
	body, err := json.Marshal(segment)

	if err != nil {
		return fmt.Errorf("UpdateSegment: failed to marshal: %v", err)
	}

	return ac.update(uri, body)
}
