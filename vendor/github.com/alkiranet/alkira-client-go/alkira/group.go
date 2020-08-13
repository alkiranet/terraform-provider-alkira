package alkira

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
)

type Group struct {
	Id              int         `json:"id"`
	Name            string      `json:"name"`
}


// Get all groups from the given tenant network
func (ac *AlkiraClient) GetGroups() ([]byte, int) {
	groupEndpoint := ac.URI + "tenantnetworks/" + ac.TenantNetworkId + "/groups"

	request, err := http.NewRequest("GET", groupEndpoint, nil)
	request.Header.Set("Content-Type", "application/json")
	response, err := ac.Client.Do(request)

	if err != nil {
		log.Printf("Error : %s", err)
	}

	defer response.Body.Close()
	data, _ := ioutil.ReadAll(response.Body)

	//log.Println(response.StatusCode)
	log.Println(string(data))

	return data, response.StatusCode
}

// Create a new Group
func (ac *AlkiraClient) CreateGroup(name string) (int, int) {
	var result Group

	groupEndpoint := ac.URI + "tenantnetworks/" + ac.TenantNetworkId + "/groups"

	body, err := json.Marshal(map[string]string{
		"name":    name,
	})

	request, err := http.NewRequest("POST", groupEndpoint, bytes.NewBuffer(body))
	request.Header.Set("Content-Type", "application/json")
	response, err := ac.Client.Do(request)

	if err != nil {
		log.Printf("Error : %s", err)
		return 0, 0
	}

	defer response.Body.Close()
	data, _ := ioutil.ReadAll(response.Body)

	json.Unmarshal([]byte(data), &result)

	return result.Id, response.StatusCode
}

// Delete a group
func (ac *AlkiraClient) DeleteGroup(id string) (int) {
	groupEndpoint := ac.URI + "tenantnetworks/" + ac.TenantNetworkId + "/groups/" + id

	request, err := http.NewRequest("DELETE", groupEndpoint, nil)
	request.Header.Set("Content-Type", "application/json")
	response, err := ac.Client.Do(request)

	if err != nil {
		log.Printf("Error : %s", err)
		return 0
	}

	defer response.Body.Close()
	data, _ := ioutil.ReadAll(response.Body)

	//log.Println(response.StatusCode)
	log.Println(string(data))

	return response.StatusCode
}
