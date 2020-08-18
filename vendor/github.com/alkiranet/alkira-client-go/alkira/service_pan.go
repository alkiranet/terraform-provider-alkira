package alkira

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

type ServicePanRequest struct {
	CXP               string               `json:"cxp"`
	CredentialId      string               `json:"credentialId"`
	Instances         []ServicePanInstance `json:"instances"`
	LicenseType       string               `json:"licenseType"`
	ManagementSegment string               `json:"managementSegment"`
	MaxInstanceCount  int                  `json:"maxInstanceCount"`
	MinInstanceCount  int                  `json:"minInstanceCount"`
	Name              string               `json:"name"`
	PanoramaEnabled   string               `json:"panoramaEnabled"`
	PanoramaTemplate  string               `json:"panoramaTemplate"`
	Segments          []string             `json:"segments"`
	SegmentOptions    interface{}          `json:"segmentOptions"`
	Size              string               `json:"size"`
	Type              string               `json:"type"`
	Version           string               `json:"version"`
}

type ServicePanInstance struct {
	CredentialId string `json:"credentialId"`
	Name         string `json:"name"`
}

type ServicePanResponse struct {
	Id int `json:"id"`
}

// CreateServicePan create service PAN
func (ac *AlkiraClient) CreateServicePan(service *ServicePanRequest) (int, error) {
	uri := fmt.Sprintf("%s/v1/tenantnetworks/%s/panfwservices", ac.URI, ac.TenantNetworkId)
	id := 0

	// Construct the request
	body, err := json.Marshal(service)

	if err != nil {
		return id, fmt.Errorf("CreateServicePan: marshal failed: %v", err)
	}

	request, err := http.NewRequest("POST", uri, bytes.NewBuffer(body))
	request.Header.Set("Content-Type", "application/json")
	response, err := ac.Client.Do(request)

	if err != nil {
		return id, fmt.Errorf("CreateServicePan: request failed: %v", err)
	}

	defer response.Body.Close()
	data, _ := ioutil.ReadAll(response.Body)

	if response.StatusCode != 201 {
		return id, fmt.Errorf("(%d) %s", response.StatusCode, string(data))
	}

	var result ServicePanResponse
	err = json.Unmarshal([]byte(data), &result)

	if err != nil {
		return id, fmt.Errorf("CreateServicePan: parse failed: %v", err)
	}

	id = result.Id
	return id, nil
}

// DeleteServicePan delete a Service PAN
func (ac *AlkiraClient) DeleteServicePan(id int) error {
	uri := fmt.Sprintf("%s/v1/tenantnetworks/%s/panfwservices/%d", ac.URI, ac.TenantNetworkId, id)

	request, err := http.NewRequest("DELETE", uri, nil)
	request.Header.Set("Content-Type", "application/json")
	response, err := ac.Client.Do(request)

	if err != nil {
		return fmt.Errorf("DeleteServicePan: request failed: %v", err)
	}

	defer response.Body.Close()
	data, _ := ioutil.ReadAll(response.Body)

	if response.StatusCode != 200 {
		return fmt.Errorf("(%d) %s", response.StatusCode, string(data))
	}

	return nil
}
