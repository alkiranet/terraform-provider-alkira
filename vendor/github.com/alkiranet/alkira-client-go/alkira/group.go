package alkira

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

type Group struct {
	Id   int    `json:"id"`
	Name string `json:"name"`
}

// Get all groups from the given tenant network
func (ac *AlkiraClient) GetGroups() (string, error) {
	uri := fmt.Sprintf("%s/tenantnetworks/%s/groups", ac.URI, ac.TenantNetworkId)

	request, err := http.NewRequest("GET", uri, nil)
	request.Header.Set("Content-Type", "application/json")
	response, err := ac.Client.Do(request)

	if err != nil {
		return "", fmt.Errorf("failed to get groups: %v", err)
	}

	defer response.Body.Close()
	data, _ := ioutil.ReadAll(response.Body)

	if response.StatusCode != 200 {
		return "", fmt.Errorf("(%d) %s", response.StatusCode, string(data))
	}

	return string(data), nil
}

// Create a new Group
func (ac *AlkiraClient) CreateGroup(name string) (int, error) {
	uri := fmt.Sprintf("%s/tenantnetworks/%s/groups", ac.URI, ac.TenantNetworkId)

	body, err := json.Marshal(map[string]string{
		"name": name,
	})

	request, err := http.NewRequest("POST", uri, bytes.NewBuffer(body))
	request.Header.Set("Content-Type", "application/json")
	response, err := ac.Client.Do(request)

	if err != nil {
		return 0, fmt.Errorf("failed to create group: %v", err)
	}

	defer response.Body.Close()
	data, _ := ioutil.ReadAll(response.Body)

	var result Group
	json.Unmarshal([]byte(data), &result)

	if response.StatusCode != 201 {
		return 0, fmt.Errorf("(%d) %s", response.StatusCode, string(data))
	}

	return result.Id, nil
}

// Delete a group
func (ac *AlkiraClient) DeleteGroup(id int) error {
	uri := fmt.Sprintf("%s/tenantnetworks/%s/groups/%d", ac.URI, ac.TenantNetworkId, id)

	request, err := http.NewRequest("DELETE", uri, nil)
	request.Header.Set("Content-Type", "application/json")
	response, err := ac.Client.Do(request)

	if err != nil {
		return fmt.Errorf("failed to delete group: %v", err)
	}

	defer response.Body.Close()
	data, _ := ioutil.ReadAll(response.Body)

	if response.StatusCode != 200 {
		return fmt.Errorf("(%d) %s", response.StatusCode, string(data))
	}

	return nil
}
