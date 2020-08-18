package alkira

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

type InternetApplicationRequest struct {
	ConnectorId   string `json:"connectorId"`
	ConnectorType string `json:"connectorType"`
	FqdnPrefix    string `json:"fqdnPrefix"`
	Group         string `json:"group"`
	Name          string `json:"name"`
	PrivateIp     string `json:"privateIp"`
	PrivatePort   string `json:"privatePort"`
	SegmentName   string `json:"segmentName"`
	Size          string `json:"size"`
}

type InternetApplicationResponse struct {
	Id int `json:"id"`
}

// CreateInternetApplication create an internet application
func (ac *AlkiraClient) CreateInternetApplication(app *InternetApplicationRequest) (int, error) {
	uri := fmt.Sprintf("%s/v1/tenantnetworks/%s/internet-applications", ac.URI, ac.TenantNetworkId)
	id := 0

	// Construct the request
	body, err := json.Marshal(app)

	request, err := http.NewRequest("POST", uri, bytes.NewBuffer(body))
	request.Header.Set("Content-Type", "application/json")
	response, err := ac.Client.Do(request)

	if err != nil {
		return id, fmt.Errorf("CreateInternetApplication: request failed: %v", err)
	}

	defer response.Body.Close()
	data, _ := ioutil.ReadAll(response.Body)

	var result InternetApplicationResponse
	json.Unmarshal([]byte(data), &result)

	if response.StatusCode != 201 {
		return id, fmt.Errorf("(%d) %s", response.StatusCode, string(data))
	}

	id = result.Id
	return id, nil
}

// DeleteInternetApplication delete given internet application by id
func (ac *AlkiraClient) DeleteInternetApplication(id int) error {
	uri := fmt.Sprintf("%s/v1/tenantnetworks/%s/internet-applications/%d", ac.URI, ac.TenantNetworkId, id)

	request, err := http.NewRequest("DELETE", uri, nil)
	request.Header.Set("Content-Type", "application/json")
	response, err := ac.Client.Do(request)

	if err != nil {
		return fmt.Errorf("DeleteInternetApplication: request failed: %v", err)
	}

	defer response.Body.Close()
	data, _ := ioutil.ReadAll(response.Body)

	if response.StatusCode != 200 {
		return fmt.Errorf("(%d) %s", response.StatusCode, string(data))
	}

	return nil
}
