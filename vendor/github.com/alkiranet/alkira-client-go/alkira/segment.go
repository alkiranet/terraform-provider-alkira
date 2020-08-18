package alkira

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

type Segment struct {
	Id   int    `json:"id"`
	Name string `json:"name"`
}

// Get all segments from the given tenant network
func (ac *AlkiraClient) GetSegments() (string, error) {
	uri := fmt.Sprintf("%s/tenantnetworks/%s/segments", ac.URI, ac.TenantNetworkId)

	request, err := http.NewRequest("GET", uri, nil)
	request.Header.Set("Content-Type", "application/json")
	response, err := ac.Client.Do(request)

	if err != nil {
		return "", fmt.Errorf("failed to get segments: %v", err)
	}

	defer response.Body.Close()
	data, _ := ioutil.ReadAll(response.Body)

	if response.StatusCode != 200 {
		return "", fmt.Errorf("(%d) %s", response.StatusCode, string(data))
	}

	return string(data), nil
}

// Get single segment from the given tenant network by segment Id
func (ac *AlkiraClient) GetSegment(id int) (string, error) {
	uri := fmt.Sprintf("%s/tenantnetworks/%s/segments/%d", ac.URI, ac.TenantNetworkId, id)

	request, err := http.NewRequest("GET", uri, nil)
	request.Header.Set("Content-Type", "application/json")
	response, err := ac.Client.Do(request)

	if err != nil {
		return "", fmt.Errorf("failed to get segment %d: %v", id, err)
	}

	defer response.Body.Close()
	data, _ := ioutil.ReadAll(response.Body)

	if response.StatusCode != 200 {
		return "", fmt.Errorf("(%d) %s", response.StatusCode, string(data))
	}

	return string(data), nil
}

// Create a new Segment
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
		log.Printf("Error : %s", err)
		return 0, fmt.Errorf("failed to create segment: %v", err)
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

// Delete a segment by given segment Id
func (ac *AlkiraClient) DeleteSegment(id int) error {

	uri := fmt.Sprintf("%s/tenantnetworks/%s/segments/%d", ac.URI, ac.TenantNetworkId, id)

	request, err := http.NewRequest("DELETE", uri, nil)
	request.Header.Set("Content-Type", "application/json")
	response, err := ac.Client.Do(request)

	if err != nil {
		return fmt.Errorf("failed to delete segment: %v", err)
	}

	defer response.Body.Close()
	data, _ := ioutil.ReadAll(response.Body)

	if response.StatusCode != 202 && response.StatusCode != 204 {
		return fmt.Errorf("(%d) %s", response.StatusCode, string(data))
	}

	return nil
}
