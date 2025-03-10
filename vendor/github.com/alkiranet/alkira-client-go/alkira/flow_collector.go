// Copyright (C) 2020-2025 Alkira Inc. All Rights Reserved.

package alkira

import (
	"encoding/json"
	"fmt"
)

type FlowCollector struct {
	Id                   json.Number `json:"id"`
	Name                 string      `json:"name"`
	Description          string      `json:"description,omitempty"`
	CollectorType        string      `json:"collectorType"`
	Enabled              bool        `json:"enabled"`
	Segment              string      `json:"segment,omitempty"`
	DestinationIp        string      `json:"destinationIp,omitempty"`
	DestinationFqdn      string      `json:"destinationFqdn,omitempty"`
	DestinationPort      int         `json:"destinationPort"`
	TransportProtocol    string      `json:"transportProtocol"`
	ExportType           string      `json:"exportType"`
	FlowRecordTemplateId int         `json:"flowRecordTemplateId"`
	Cxps                 []string    `json:"cxps"`
}

// NewFlowCollector new flow collector
func NewFlowCollector(ac *AlkiraClient) *AlkiraAPI[FlowCollector] {
	uri := fmt.Sprintf("%s/tenantnetworks/%s/flow-collectors", ac.URI, ac.TenantNetworkId)
	api := &AlkiraAPI[FlowCollector]{ac, uri, true}
	return api
}
