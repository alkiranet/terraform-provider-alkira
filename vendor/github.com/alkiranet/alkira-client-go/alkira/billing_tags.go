// Copyright (C) 2020-2021 Alkira Inc. All Rights Reserved.

package alkira

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

type BillingTag struct {
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

	var result BillingTag
	json.Unmarshal([]byte(data), &result)

	if response.StatusCode != 201 {
		return 0, fmt.Errorf("(%d) %s", response.StatusCode, string(data))
	}

	return result.Id, nil
}

// GetBillingTag get a single billing tag by Id
func (ac *AlkiraClient) GetBillingTagById(id int) (BillingTag, error) {
	uri := fmt.Sprintf("%s/tags/%d", ac.URI, id)

	var billingTag BillingTag

	request, err := http.NewRequest("GET", uri, nil)
	request.Header.Set("Content-Type", "application/json")
	response, err := ac.Client.Do(request)

	if err != nil {
		return billingTag, fmt.Errorf("GetBillingTag: request failed, %v", err)
	}

	defer response.Body.Close()
	data, _ := ioutil.ReadAll(response.Body)

	if response.StatusCode != 200 {
		return billingTag, fmt.Errorf("(%d) %s", response.StatusCode, string(data))
	}

	err = json.Unmarshal([]byte(data), &billingTag)

	if err != nil {
		return billingTag, fmt.Errorf("GetBillingTagById: parse failed: %v", err)
	}

	return billingTag, nil
}

// GetBillingTagByName get a billing tag by its name
func (ac *AlkiraClient) GetBillingTagByName(name string) (BillingTag, error) {
	var billingTag BillingTag

	if len(name) == 0 {
		return billingTag, fmt.Errorf("Invalid billingTag name input")
	}

	billingTags, err := ac.GetBillingTags()

	if err != nil {
		return billingTag, err
	}

	var result []BillingTag
	json.Unmarshal([]byte(billingTags), &result)

	for _, t := range result {
		if t.Name == name {
			return t, nil
		}
	}

	return billingTag, fmt.Errorf("failed to find the billingTag by %s", name)
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

// UpdateBillingTag update a billing tag by Id
func (ac *AlkiraClient) UpdateBillingTag(id int, name string) error {
	uri := fmt.Sprintf("%s/tags/%d", ac.URI, id)

	body, err := json.Marshal(map[string]string{
		"name": name,
	})

	request, err := http.NewRequest("PUT", uri, bytes.NewBuffer(body))
	request.Header.Set("Content-Type", "application/json")
	response, err := ac.Client.Do(request)

	if err != nil {
		return fmt.Errorf("UpdateBillingTag: request faile, %v", err)
	}

	defer response.Body.Close()
	data, _ := ioutil.ReadAll(response.Body)

	if response.StatusCode != 200 {
		return fmt.Errorf("(%d) %s", response.StatusCode, string(data))
	}

	return nil
}
