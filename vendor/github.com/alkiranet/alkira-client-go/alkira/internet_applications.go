// Copyright (C) 2020-2022 Alkira Inc. All Rights Reserved.

package alkira

import (
	"encoding/json"
	"fmt"
)

type InternetApplicationTargets struct {
	Type  string `json:"type"`
	Value string `json:"value"`
	Ports []int  `json:"ports"`
}

type InternetApplication struct {
	BillingTags            []int                        `json:"billingTags"`
	ConnectorId            int                          `json:"connectorId"`
	ConnectorType          string                       `json:"connectorType"`
	FqdnPrefix             string                       `json:"fqdnPrefix"`
	Group                  string                       `json:"group"`
	Id                     json.Number                  `json:"id,omitempty"` // response only
	InboundConnectorId     string                       `json:"inboundConnectorId,omitempty"`
	InboundConnectorType   string                       `json:"inboundConnectorType,omitempty"`
	InboundInternetGroupId json.Number                  `json:"inboundInternetGroupId,omitempty"`
	Name                   string                       `json:"name"`
	PublicIps              []string                     `json:"publicIps"`
	SegmentName            string                       `json:"segmentName"`
	Size                   string                       `json:"size"`
	Targets                []InternetApplicationTargets `json:"targets,omitempty"`
}

// CreateInternetApplication create an internet application
func (ac *AlkiraClient) CreateInternetApplication(app *InternetApplication) (string, string, error) {
	uri := fmt.Sprintf("%s/v1/tenantnetworks/%s/internet-applications", ac.URI, ac.TenantNetworkId)

	// Construct the request
	body, err := json.Marshal(app)

	if err != nil {
		return "", "", fmt.Errorf("CreateInternetApplication: failed to marshal: %v", err)
	}

	data, err := ac.create(uri, body, true)

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

// DeleteInternetApplication delete given internet application by ID
func (ac *AlkiraClient) DeleteInternetApplication(id string) error {
	uri := fmt.Sprintf("%s/v1/tenantnetworks/%s/internet-applications/%s", ac.URI, ac.TenantNetworkId, id)

	return ac.delete(uri, true)
}

// UpdateInternetApplication update a given internet application by ID
func (ac *AlkiraClient) UpdateInternetApplication(id string, app *InternetApplication) error {
	uri := fmt.Sprintf("%s/v1/tenantnetworks/%s/internet-applications/%s", ac.URI, ac.TenantNetworkId, id)

	// Construct the request
	body, err := json.Marshal(app)

	if err != nil {
		return fmt.Errorf("UpdateInternetApplication: failed to marshal: %v", err)
	}

	return ac.update(uri, body, true)
}

// GetInternetApplication get internet application by ID
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
