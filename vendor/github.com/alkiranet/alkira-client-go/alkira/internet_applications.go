// Copyright (C) 2020-2021 Alkira Inc. All Rights Reserved.

package alkira

import (
	"encoding/json"
	"fmt"
)

type InternetApplication struct {
	BillingTags            []int       `json:"billingTags"`
	ConnectorId            string      `json:"connectorId"`
	ConnectorType          string      `json:"connectorType"`
	FqdnPrefix             string      `json:"fqdnPrefix"`
	Id                     json.Number `json:"id,omitempty"`
	InboundInternetGroupId json.Number `json:"inboundInternetGroupId,omitempty"`
	Group                  string      `json:"group"`
	Name                   string      `json:"name"`
	PrivateIp              string      `json:"privateIp"`
	PrivatePort            string      `json:"privatePort"`
	SegmentName            string      `json:"segmentName"`
	Size                   string      `json:"size"`
}

// CreateInternetApplication create an internet application
func (ac *AlkiraClient) CreateInternetApplication(app *InternetApplication) (string, string, error) {
	uri := fmt.Sprintf("%s/v1/tenantnetworks/%s/internet-applications", ac.URI, ac.TenantNetworkId)

	// Construct the request
	body, err := json.Marshal(app)

	if err != nil {
		return "", "", fmt.Errorf("CreateInternetApplication: failed to marshal: %v", err)
	}

	data, err := ac.create(uri, body)

	if err != nil {
		return "", "", err
	}

	var result InternetApplication
	json.Unmarshal([]byte(data), &result)

	if err != nil {
		return "", "", fmt.Errorf("CreateInternetApplication: failed to unmarshal: %v", err)
	}

	return string(result.Id), string(result.InboundInternetGroupId), nil
}

// DeleteInternetApplication delete given internet application by Id
func (ac *AlkiraClient) DeleteInternetApplication(id string) error {
	uri := fmt.Sprintf("%s/v1/tenantnetworks/%s/internet-applications/%s", ac.URI, ac.TenantNetworkId, id)

	return ac.delete(uri)
}

// UpdateInternetApplication update a given internet application by Id
func (ac *AlkiraClient) UpdateInternetApplication(id string, app *InternetApplication) error {
	uri := fmt.Sprintf("%s/v1/tenantnetworks/%s/internet-applications/%s", ac.URI, ac.TenantNetworkId, id)

	// Construct the request
	body, err := json.Marshal(app)

	if err != nil {
		return fmt.Errorf("UpdateInternetApplication: failed to marshal: %v", err)
	}

	return ac.update(uri, body)
}

// GetInternetApplication get internet application by Id
func (ac *AlkiraClient) GetInternetApplication(id string) (*InternetApplication, error) {
	uri := fmt.Sprintf("%s/v1/tenantnetworks/%s/internet-applications/%s", ac.URI, ac.TenantNetworkId, id)

	data, err := ac.get(uri)

	if err != nil {
		return nil, err
	}

	var result InternetApplication
	err = json.Unmarshal([]byte(data), &result)

	if err != nil {
		return nil, fmt.Errorf("GetInternetApplication: failed to unmarshal: %v", err)
	}

	return &result, nil
}
