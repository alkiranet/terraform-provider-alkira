package alkira

import (
	"bytes"
	"encoding/json"
	"errors"
	"io/ioutil"
	"log"
	"net/http"
)

type PolicyRequest struct {
	Description    string   `json:"description"`
	Enabled        string   `json:"enabled"`
	FromGroups     []string `json:"fromGroups"`
	Name           string   `json:"name"`
	RuleListId     string   `json:"ruleListId"`
	SegmentIds     []string `json:"segmentIds"`
	ToGroups       []string `json:"toGroups"`
}

type policyResponse struct {
	Id int `json:"id"`
}

// Create a policy
func (ac *AlkiraClient) CreatePolicy(p *PolicyRequest) (int, error) {
	url := ac.URI + "v1/tenantnetworks/" + ac.TenantNetworkId + "/policy/policies"
	id  := 0

	// Construct the request
	body, err := json.Marshal(p)

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

	var result policyResponse
	json.Unmarshal([]byte(data), &result)

	if response.StatusCode != 201 {
		return id, errors.New("Failed to create policy")
	}

	id = result.Id

	return id, nil
}

// Delete a policy
func (ac *AlkiraClient) DeletePolicy(id string) (error) {
	url := ac.URI + "v1/tenantnetworks/" + ac.TenantNetworkId + "/policy/policies/" + id

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
		return errors.New("Failed to delete policy " + id)
	}

	return nil
}
