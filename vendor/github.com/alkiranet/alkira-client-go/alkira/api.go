// Copyright (C) 2023-2025 Alkira Inc. All Rights Reserved.

package alkira

import (
	"encoding/json"
	"fmt"
)

// Generic struct to define a Alkira API
type AlkiraAPI[T any] struct {
	Client    *AlkiraClient
	Uri       string
	Provision bool
}

// Create create a resource
func (a *AlkiraAPI[T]) Create(resource *T) (*T, string, error, error) {

	// Construct the request
	body, err := json.Marshal(resource)

	if err != nil {
		return nil, "", fmt.Errorf("api-create: failed to marshal: %v", err), nil
	}

	data, state, err, errProv := a.Client.create(a.Uri, body, a.Provision)

	if err != nil {
		return nil, state, err, errProv
	}

	var result T
	err = json.Unmarshal([]byte(data), &result)

	if err != nil {
		return nil, state, fmt.Errorf("api-create: failed to unmarshal: %v", err), errProv
	}

	return &result, state, nil, errProv
}

// Delete delete a resource by its ID
func (a *AlkiraAPI[T]) Delete(id string) (string, error, error) {

	// Construct single resource URI
	uri := fmt.Sprintf("%s/%s", a.Uri, id)

	return a.Client.delete(uri, a.Provision)
}

// Update update a resource by its ID
func (a *AlkiraAPI[T]) Update(id string, resource *T) (string, error, error) {

	// Construct single resource URI
	uri := fmt.Sprintf("%s/%s", a.Uri, id)

	// Construct the request
	body, err := json.Marshal(resource)

	if err != nil {
		return "", fmt.Errorf("api-update: failed to marshal: %v", err), nil
	}

	return a.Client.update(uri, body, a.Provision)
}

// GetAll get all resources
func (a *AlkiraAPI[T]) GetAll() (string, error) {
	data, _, err := a.Client.get(a.Uri)
	return string(data), err
}

// GetById get a resource by its ID
func (a *AlkiraAPI[T]) GetById(id string) (*T, string, error) {

	// Construct single resource URI
	uri := fmt.Sprintf("%s/%s?includeMarkedForDeletion=true", a.Uri, id)

	data, provState, err := a.Client.get(uri)

	if err != nil {
		return nil, provState, err
	}

	var result T
	err = json.Unmarshal([]byte(data), &result)

	if err != nil {
		return nil, provState, fmt.Errorf("api-get-all: failed to unmarshal: %v", err)
	}

	return &result, provState, nil
}

// GetByName get a resource by its name
func (a *AlkiraAPI[T]) GetByName(name string) (*T, string, error) {

	if len(name) == 0 {
		return nil, "", fmt.Errorf("api-get-by-name: Invalid resource name")
	}

	// Construct single resource URI
	uri := fmt.Sprintf("%s?name=%s&paginated=false", a.Uri, name)

	data, state, err := a.Client.getByName(uri)

	if err != nil {
		return nil, "", err
	}

	var result []T
	err = json.Unmarshal([]byte(data), &result)

	if err != nil {
		return nil, state, fmt.Errorf("api-get-by-name: failed to unmarshal: %v", err)
	}

	if len(result) != 1 {
		return nil, state, fmt.Errorf("api-get-by-name: failed to get resource by name: %s", name)
	}

	return &result[0], state, nil
}
