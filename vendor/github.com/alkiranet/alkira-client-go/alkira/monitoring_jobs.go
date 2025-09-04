// Copyright (C) 2025 Alkira Inc. All Rights Reserved.

package alkira

import (
	"fmt"
	"net/url"
)

// GetJobs get jobs with optional filter with status or type
func (ac *AlkiraClient) GetJobs(queryStatus string, queryType string) (string, error) {

	baseUri := fmt.Sprintf("%s/api/jobs", ac.URI)

	uri, err := url.Parse(baseUri)

	if err != nil {
		return "", fmt.Errorf("GetJobs: failed to parse URI: %v", err)
	}

	// Process optional query parameters
	q := uri.Query()

	if queryStatus != "" {
		q.Add("status", queryStatus)
	}

	if queryType != "" {
		q.Add("type", queryType)
	}

	// GET
	uri.RawQuery = q.Encode()
	data, _, err := ac.get(uri.String())

	return string(data), err
}
