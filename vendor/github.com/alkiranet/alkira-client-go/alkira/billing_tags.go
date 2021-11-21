// Copyright (C) 2020-2021 Alkira Inc. All Rights Reserved.

package alkira

import (
	"encoding/json"
	"fmt"
	"strconv"
)

type BillingTag struct {
	Id          int    `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

// CreateBillingTag create a new billing tag
func (ac *AlkiraClient) CreateBillingTag(name string, description string) (string, error) {
	uri := fmt.Sprintf("%s/tags", ac.URI)

	body, err := json.Marshal(map[string]string{
		"name":        name,
		"description": description,
	})

	if err != nil {
		return "", fmt.Errorf("CreateBillingTag: failed to marshal: %v", err)
	}

	data, err := ac.create(uri, body)

	if err != nil {
		return "", err
	}

	var result BillingTag
	err = json.Unmarshal([]byte(data), &result)

	if err != nil {
		return "", fmt.Errorf("CreateBillingTag: failed to unmarshal: %v", err)
	}

	return strconv.Itoa(result.Id), nil
}

// DeleteBillingTag delete a billing tag by Id
func (ac *AlkiraClient) DeleteBillingTag(id string) error {
	uri := fmt.Sprintf("%s/tags/%s", ac.URI, id)
	return ac.delete(uri)
}

// UpdateBillingTag update a billing tag by Id
func (ac *AlkiraClient) UpdateBillingTag(id string, name string, description string) error {
	uri := fmt.Sprintf("%s/tags/%s", ac.URI, id)

	body, err := json.Marshal(map[string]string{
		"name":        name,
		"description": description,
	})

	if err != nil {
		return fmt.Errorf("UpdateBillingTag: failed to marshal: %v", err)
	}

	return ac.update(uri, body)
}

// GetBillingTags get all billing tags from the given tenant network
func (ac *AlkiraClient) GetBillingTags() (string, error) {
	uri := fmt.Sprintf("%s/tags", ac.URI)
	data, err := ac.get(uri)
	return string(data), err
}

// GetBillingTagById get a single billing tag by Id
func (ac *AlkiraClient) GetBillingTagById(id string) (BillingTag, error) {
	uri := fmt.Sprintf("%s/tags/%s", ac.URI, id)

	var billingTag BillingTag

	data, err := ac.get(uri)

	if err != nil {
		return billingTag, err
	}

	err = json.Unmarshal([]byte(data), &billingTag)

	if err != nil {
		return billingTag, fmt.Errorf("GetBillingTagById: failed to unmarshal: %v", err)
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
