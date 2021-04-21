// Copyright (C) 2020-2021 Alkira Inc. All Rights Reserved.

package alkira

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
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

	request, err := http.NewRequest("GET", uri, nil)
	request.Header.Set("Content-Type", "application/json")
	response, err := ac.Client.Do(request)

	if err != nil {
		return "", fmt.Errorf("GetSegments: request failed: %v", err)
	}

	defer response.Body.Close()
	data, _ := ioutil.ReadAll(response.Body)

	if response.StatusCode != 200 {
		return "", fmt.Errorf("(%d) %s", response.StatusCode, string(data))
	}

	return string(data), nil
}

// GetSegment get single segment by Id
func (ac *AlkiraClient) GetSegmentById(id int) (Segment, error) {
	uri := fmt.Sprintf("%s/tenantnetworks/%s/segments/%d", ac.URI, ac.TenantNetworkId, id)

	var segment Segment

	request, err := http.NewRequest("GET", uri, nil)
	request.Header.Set("Content-Type", "application/json")
	response, err := ac.Client.Do(request)

	if err != nil {
		return segment, fmt.Errorf("GetSegmentById: request failed: %v", err)
	}

	defer response.Body.Close()
	data, _ := ioutil.ReadAll(response.Body)

	if response.StatusCode != 200 {
		return segment, fmt.Errorf("(%d) %s", response.StatusCode, string(data))
	}

	err = json.Unmarshal([]byte(data), &segment)

	if err != nil {
		return segment, fmt.Errorf("GetSegmentById: parse failed: %v", err)
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
func (ac *AlkiraClient) CreateSegment(name string, asn string, ipBlock string) (int, error) {

	uri := fmt.Sprintf("%s/tenantnetworks/%s/segments", ac.URI, ac.TenantNetworkId)

	body, err := json.Marshal(map[string]string{
		"name":    name,
		"asn":     asn,
		"ipBlock": ipBlock,
	})

	request, err := http.NewRequest("POST", uri, bytes.NewBuffer(body))
	request.Header.Set("Content-Type", "application/json")
	response, err := ac.Client.Do(request)

	if err != nil {
		return 0, fmt.Errorf("CreateSegment: request failed, %v", err)
	}

	defer response.Body.Close()
	data, _ := ioutil.ReadAll(response.Body)

	var result Segment
	json.Unmarshal([]byte(data), &result)

	if response.StatusCode != 201 {
		return 0, fmt.Errorf("(%d) %s", response.StatusCode, string(data))
	}

	return result.Id, nil
}

// DeleteSegment delete a segment by given segment Id
func (ac *AlkiraClient) DeleteSegment(id int) error {

	uri := fmt.Sprintf("%s/tenantnetworks/%s/segments/%d", ac.URI, ac.TenantNetworkId, id)

	request, err := http.NewRequest("DELETE", uri, nil)
	request.Header.Set("Content-Type", "application/json")
	response, err := ac.Client.Do(request)

	if err != nil {
		return fmt.Errorf("DeleteSegment: request failed, %v", err)
	}

	defer response.Body.Close()
	data, _ := ioutil.ReadAll(response.Body)

	if response.StatusCode != 200 {
		return fmt.Errorf("(%d) %s", response.StatusCode, string(data))
	}

	return nil
}

// UpdateSegment update a segment by segment Id
func (ac *AlkiraClient) UpdateSegment(id int, name string, asn string, ipBlock string) error {

	uri := fmt.Sprintf("%s/tenantnetworks/%s/segments/%d", ac.URI, ac.TenantNetworkId, id)

	body, err := json.Marshal(map[string]string{
		"name":    name,
		"asn":     asn,
		"ipBlock": ipBlock,
	})

	request, err := http.NewRequest("PUT", uri, bytes.NewBuffer(body))
	request.Header.Set("Content-Type", "application/json")
	response, err := ac.Client.Do(request)

	if err != nil {
		return fmt.Errorf("UpdateSegment: request failed, %v", err)
	}

	defer response.Body.Close()
	data, _ := ioutil.ReadAll(response.Body)

	if response.StatusCode != 200 {
		return fmt.Errorf("(%d) %s", response.StatusCode, string(data))
	}

	return nil
}
