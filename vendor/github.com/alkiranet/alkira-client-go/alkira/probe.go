// Copyright (C) 2021-2025 Alkira Inc. All Rights Reserved.

package alkira

import (
	"encoding/json"
	"fmt"
)

// Probe represents any type of probe
type Probe struct {
	ID               string             `json:"id,omitempty"`
	Name             string             `json:"name"`
	Type             string             `json:"type"`
	Parameters       json.RawMessage    `json:"parameters"`
	Description      string             `json:"description"`
	Enabled          bool               `json:"enabled"`
	NetworkEntity    ProbeNetworkEntity `json:"networkEntity"`
	FailureThreshold int                `json:"failureThreshold,omitempty"`
	SuccessThreshold int                `json:"successThreshold,omitempty"`
	PeriodSeconds    int                `json:"periodSeconds,omitempty"`
	TimeoutSeconds   int                `json:"timeoutSeconds,omitempty"`
	LastUpdated      int64              `json:"lastUpdated,omitempty"`
}

// ProbeNetworkEntity represents a network entity
type ProbeNetworkEntity struct {
	Type string `json:"type"`
	ID   string `json:"id"`
}

// TcpProbe defines TCP probe parameters
type TcpProbe struct {
	Port int `json:"port"`
}

// HttpProbe defines HTTP probe parameters
type HttpProbe struct {
	URI        string           `json:"uri"`
	Validators []ProbeValidator `json:"validators,omitempty"`
	Headers    map[string]any   `json:"headers,omitempty"`
}

// HttpsProbe defines HTTPS probe parameters
type HttpsProbe struct {
	ServerName            string           `json:"serverName,omitempty"`
	DisableCertValidation bool             `json:"disableCertValidation,omitempty"`
	CaCertificate         string           `json:"caCertificate,omitempty"`
	URI                   string           `json:"uri"`
	Validators            []ProbeValidator `json:"validators,omitempty"`
	Headers               map[string]any   `json:"headers,omitempty"`
}

// ProbeValidator represents a validator for HTTP/HTTPS probes
type ProbeValidator struct {
	Type       string          `json:"type"`
	Parameters json.RawMessage `json:"parameters"`
}

// ProbeStatusCodeValidator defines HTTP status code validator
type ProbeStatusCodeValidator struct {
	StatusCode string `json:"statusCode"`
}

// ProbeResponseBodyValidator defines HTTP response body validator
type ProbeResponseBodyValidator struct {
	Regex string `json:"regex"`
}

// NewProbe creates a new Probe API client
func NewProbe(ac *AlkiraClient) *AlkiraAPI[Probe] {
	uri := fmt.Sprintf("%s/tenantnetworks/%s/probes", ac.URI, ac.TenantNetworkId)
	api := &AlkiraAPI[Probe]{ac, uri, true}
	return api
}
