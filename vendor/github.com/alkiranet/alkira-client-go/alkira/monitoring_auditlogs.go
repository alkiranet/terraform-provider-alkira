// Copyright (C) 2025 Alkira Inc. All Rights Reserved.

package alkira

import (
	"fmt"
	"net/url"
)

type AuditLogResponse struct {
	Data       []AuditLogEntry `json:"data"`
	Pagination Pagination      `json:"pagination"`
}

type AuditLogEntry struct {
	CreatedAt   int64         `json:"createdAt"`
	Description string        `json:"description"`
	ID          string        `json:"id"`
	Initiator   string        `json:"initiator"`
	IPAddress   string        `json:"ipAddress"`
	Status      string        `json:"status"`
	Tags        []AuditLogTag `json:"tags"`
	TenantID    int           `json:"tenantId"`
	Type        string        `json:"type"`
}

type AuditLogTag map[string]string

type Pagination struct {
	Offset int `json:"offset"`
	Limit  int `json:"limit"`
	Hits   int `json:"hits"`
}

// GetAuditLogs get audit logs with optional filter with log status or
// log type
func (ac *AlkiraClient) GetAuditLogs(queryStatus string, queryType string) (string, error) {

	baseUri := fmt.Sprintf("%s/api/auditlogs", ac.URI)

	uri, err := url.Parse(baseUri)

	if err != nil {
		return "", fmt.Errorf("GetAuditLogs: failed to parse URI: %v", err)
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
