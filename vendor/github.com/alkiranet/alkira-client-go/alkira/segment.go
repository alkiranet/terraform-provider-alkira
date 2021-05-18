// Copyright (C) 2020-2021 Alkira Inc. All Rights Reserved.

package alkira

import (
	"encoding/json"
	"fmt"
	"strconv"
)

type Segment struct {
	Asn     int    `json:"asn"`
	Id      int    `json:"id"`
	IpBlock string `json:"ipBlock"`
	Name    string `json:"name"`
}

// Get all segments from the given tenant network
func (ac *AlkiraClient) GetSegments() (string, error) {
	uri := fmt.Sprintf("%s/tenantnetworks/%s/segments", ac.URI, ac.TenantNetworkId)

	data, err := ac.get(uri)
	return string(data), err
}

// GetSegment get single segment by Id
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
func (ac *AlkiraClient) CreateSegment(name string, asn string, ipBlock string) (string, error) {

	uri := fmt.Sprintf("%s/tenantnetworks/%s/segments", ac.URI, ac.TenantNetworkId)

	body, err := json.Marshal(map[string]string{
		"name":    name,
		"asn":     asn,
		"ipBlock": ipBlock,
	})

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

// DeleteSegment delete a segment by given segment Id
func (ac *AlkiraClient) DeleteSegment(id string) error {
	uri := fmt.Sprintf("%s/tenantnetworks/%s/segments/%s", ac.URI, ac.TenantNetworkId, id)
	return ac.delete(uri)
}

// UpdateSegment update a segment by segment Id
func (ac *AlkiraClient) UpdateSegment(id string, name string, asn string, ipBlock string) error {

	uri := fmt.Sprintf("%s/tenantnetworks/%s/segments/%s", ac.URI, ac.TenantNetworkId, id)

	body, err := json.Marshal(map[string]string{
		"name":    name,
		"asn":     asn,
		"ipBlock": ipBlock,
	})

	if err != nil {
		return fmt.Errorf("UpdateSegment: failed to marshal: %v", err)
	}

	return ac.update(uri, body)
}
