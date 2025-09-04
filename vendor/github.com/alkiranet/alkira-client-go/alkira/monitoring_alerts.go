// Copyright (C) 2025 Alkira Inc. All Rights Reserved.

package alkira

import (
	"fmt"
	"net/url"
)

// GetAlerts get alerts with optional query parameters
func (ac *AlkiraClient) GetAlerts(queryStatus string, queryType string, queryPriority string) (string, error) {

	baseUri := fmt.Sprintf("%s/api/alerts", ac.URI)

	uri, err := url.Parse(baseUri)

	if err != nil {
		return "", fmt.Errorf("GetAlerts: failed to parse URI: %v", err)
	}

	// Process optional query parameters
	q := uri.Query()

	if queryStatus != "" {
		q.Add("status", queryStatus)
	}

	if queryType != "" {
		q.Add("type", queryType)
	}

	if queryPriority != "" {
		q.Add("priority", queryPriority)
	}

	// GET
	uri.RawQuery = q.Encode()
	data, _, err := ac.get(uri.String())

	return string(data), err
}
