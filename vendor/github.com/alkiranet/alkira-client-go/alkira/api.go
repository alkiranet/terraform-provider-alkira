// Copyright (C) 2023 Alkira Inc. All Rights Reserved.

package alkira

import (
	"encoding/json"
	"fmt"
)

// Generic struct to define a Alkira API
type AlkiraAPI[T any] struct {
	Client *AlkiraClient
	Uri    string
}

// Create create a resource
func (a *AlkiraAPI[T]) Create(resource *T) (*T, string, error) {

	// Construct the request
	body, err := json.Marshal(resource)

	if err != nil {
		return nil, "", fmt.Errorf("Create: failed to marshal: %v", err)
	}

	data, state, err := a.Client.create(a.Uri, body, true)

	if err != nil {
		return nil, state, err
	}

	var result T
	err = json.Unmarshal([]byte(data), &result)

	if err != nil {
		return nil, state, fmt.Errorf("Create: failed to unmarshal: %v", err)
	}

	return &result, state, nil
}

// Delete delete a resource by its ID
func (a *AlkiraAPI[T]) Delete(id string) (string, error) {

	// Construct single resource URI
	uri := fmt.Sprintf("%s/%s", a.Uri, id)

	return a.Client.delete(uri, true)
}

// Update update a resource by its ID
func (a *AlkiraAPI[T]) Update(id string, resource *T) (string, error) {

	// Construct single resource URI
	uri := fmt.Sprintf("%s/%s", a.Uri, id)

	// Construct the request
	body, err := json.Marshal(resource)

	if err != nil {
		return "", fmt.Errorf("Update: failed to marshal: %v", err)
	}

	return a.Client.update(uri, body, true)
}

// GetAll get all resources
func (a *AlkiraAPI[T]) GetAll() (string, error) {
	data, err := a.Client.get(a.Uri)
	return string(data), err
}

// GetById get a resource by its ID
func (a *AlkiraAPI[T]) GetById(id string) (*T, error) {

	// Construct single resource URI
	uri := fmt.Sprintf("%s/%s", a.Uri, id)

	data, err := a.Client.get(uri)

	if err != nil {
		return nil, err
	}

	var result T
	err = json.Unmarshal([]byte(data), &result)

	if err != nil {
		return nil, fmt.Errorf("Get: failed to unmarshal: %v", err)
	}

	return &result, nil
}

// GetByName get a resource by its name
func (a *AlkiraAPI[T]) GetByName(name string) (*T, string, error) {

	if len(name) == 0 {
		return nil, "", fmt.Errorf("Get: Invalid resource name")
	}

	// Construct single resource URI
	uri := fmt.Sprintf("%s?name=%s", a.Uri, name)

	data, state, err := a.Client.getByName(uri)

	if err != nil {
		return nil, "", err
	}

	var result T
	err = json.Unmarshal([]byte(data), &result)

	if err != nil {
		return nil, state, fmt.Errorf("Get: failed to unmarshal: %v", err)
	}

	return &result, state, nil
}
