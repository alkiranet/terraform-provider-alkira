package alkira

import (
	"bytes"
	"encoding/json"
	"errors"
	"io/ioutil"
	"log"
	"net/http"
)

type InternetApplicationRequest struct {
	ConnectorId    string   `json:"connectorId"`
	ConnectorType  string   `json:"connectorType"`
	FqdnPrefix     string   `json:"fqdnPrefix"`
	Group          string   `json:"group"`
	Name           string   `json:"name"`
	PrivateIp      string   `json:"privateIp"`
	PrivatePort    string   `json:"privatePort"`
	SegmentName    string   `json:"segmentName"`
	Size           string   `json:"size"`
}

type InternetApplicationResponse struct {
	Id int `json:"id"`
}

// Create an internet application
func (ac *AlkiraClient) CreateInternetApplication(app *InternetApplicationRequest) (int, error) {
	url := ac.URI + "v1/tenantnetworks/" + ac.TenantNetworkId + "/internet-applications"
	id  := 0

	// Construct the request
	body, err := json.Marshal(app)

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

	var result InternetApplicationResponse
	json.Unmarshal([]byte(data), &result)

	if response.StatusCode != 201 {
		return id, errors.New("Failed to create internet application")
	}

	id = result.Id

	return id, nil
}

// Delete one internet application
func (ac *AlkiraClient) DeleteInternetApplication(id string) (error) {
	url := ac.URI + "v1/tenantnetworks/" + ac.TenantNetworkId + "/internet-applications/" + id

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
		return errors.New("Failed to delete internet application " + id)
	}

	return nil
}
