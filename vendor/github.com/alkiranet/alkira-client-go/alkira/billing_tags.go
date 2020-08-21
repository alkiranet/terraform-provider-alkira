package alkira

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

type Billingtag struct {
	Id   int    `json:"id"`
	Name string `json:"name"`
}

// GetBillingTags get all billing tags from the given tenant network
func (ac *AlkiraClient) GetBillingTags() (string, error) {
	uri := fmt.Sprintf("%s/tags", ac.URI)

	request, err := http.NewRequest("GET", uri, nil)
	request.Header.Set("Content-Type", "application/json")
	response, err := ac.Client.Do(request)

	if err != nil {
		return "", fmt.Errorf("GetBillingTags: request failed, %v", err)
	}

	defer response.Body.Close()
	data, _ := ioutil.ReadAll(response.Body)

	if response.StatusCode != 200 {
		return "", fmt.Errorf("(%d) %s", response.StatusCode, string(data))
	}

	return string(data), nil
}

// CreateBillingTag create a new billing tag
func (ac *AlkiraClient) CreateBillingTag(name string) (int, error) {
	uri := fmt.Sprintf("%s/tags", ac.URI)

	body, err := json.Marshal(map[string]string{
		"name": name,
	})

	request, err := http.NewRequest("POST", uri, bytes.NewBuffer(body))
	request.Header.Set("Content-Type", "application/json")
	response, err := ac.Client.Do(request)

	if err != nil {
		return 0, fmt.Errorf("CreateBillingTag: request failed, %v", err)
	}

	defer response.Body.Close()
	data, _ := ioutil.ReadAll(response.Body)

	var result Billingtag
	json.Unmarshal([]byte(data), &result)

	if response.StatusCode != 201 {
		return 0, fmt.Errorf("(%d) %s", response.StatusCode, string(data))
	}

	return result.Id, nil
}

// GetBillingTag get single billing tag by Id
func (ac *AlkiraClient) GetBillingTag(id int) (string, error) {
	uri := fmt.Sprintf("%s/tags/%d", ac.URI, id)

	request, err := http.NewRequest("GET", uri, nil)
	request.Header.Set("Content-Type", "application/json")
	response, err := ac.Client.Do(request)

	if err != nil {
		return "", fmt.Errorf("GetBillingTag: request failed, %v", err)
	}

	defer response.Body.Close()
	data, _ := ioutil.ReadAll(response.Body)

	if response.StatusCode != 200 {
		return "", fmt.Errorf("(%d) %s", response.StatusCode, string(data))
	}

	return string(data), nil
}

// DeleteBillingTag delete a billing tag by Id
func (ac *AlkiraClient) DeleteBillingTag(id int) error {
	uri := fmt.Sprintf("%s/tags/%d", ac.URI, id)

	request, err := http.NewRequest("DELETE", uri, nil)
	request.Header.Set("Content-Type", "application/json")
	response, err := ac.Client.Do(request)

	if err != nil {
		return fmt.Errorf("DeleteBillingTag: request faile, %v", err)
	}

	defer response.Body.Close()
	data, _ := ioutil.ReadAll(response.Body)

	if response.StatusCode != 200 {
		return fmt.Errorf("(%d) %s", response.StatusCode, string(data))
	}

	return nil
}
