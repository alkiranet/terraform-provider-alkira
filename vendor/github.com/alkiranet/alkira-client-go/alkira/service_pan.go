package alkira

import (
	"bytes"
	"encoding/json"
	"errors"
	"io/ioutil"
	"log"
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
	Id              int         `json:"id"`
}

// Create service PAN
func (ac *AlkiraClient) CreateServicePan(service *ServicePanRequest) (int, error) {
	url := ac.URI + "v1/tenantnetworks/" + ac.TenantNetworkId + "/panfwservices"
	id  := 0

	// Construct the request
	body, err := json.Marshal(service)

	log.Println(bytes.NewBuffer(body))
	request, err := http.NewRequest("POST", url, bytes.NewBuffer(body))
	request.Header.Set("Content-Type", "application/json")
	response, err := ac.Client.Do(request)

	if err != nil {
		return id, err
	}

	defer response.Body.Close()
	data, _ := ioutil.ReadAll(response.Body)

	log.Println(response.StatusCode)
	log.Println(string(data))

	var result ServicePanResponse
	json.Unmarshal([]byte(data), &result)

	if response.StatusCode != 201 {
		return id, errors.New("Failed to create service PAN")
	}

	id = result.Id

	return id, nil
}

// Delete a Service PAN
func (ac *AlkiraClient) DeleteServicePan(id string) (error) {
	url := ac.URI + "v1/tenantnetworks/" + ac.TenantNetworkId + "/panfwservices/" + id

	request, err := http.NewRequest("DELETE", url, nil)
	request.Header.Set("Content-Type", "application/json")
	response, err := ac.Client.Do(request)

	if err != nil {
		return err
	}

	defer response.Body.Close()
	data, _ := ioutil.ReadAll(response.Body)

	log.Println(response.StatusCode)
	log.Println(string(data))

	if response.StatusCode != 200 {
		return errors.New("Failed to delete PAN service " + id)
	}

	return nil
}
