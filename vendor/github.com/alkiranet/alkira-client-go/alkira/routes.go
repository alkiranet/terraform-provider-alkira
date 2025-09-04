// Copyright (C) 2025 Alkira Inc. All Rights Reserved.

package alkira

import (
	"encoding/json"
	"fmt"
	"net/url"
	"strconv"
)

// Route query parameters for filtering routes
type RouteQueryParams struct {
	Type                 string `json:"type"`                           // required: received, advertised, overlap
	SegmentName          string `json:"segmentName,omitempty"`          // segment name filter
	SegmentNames         string `json:"segmentNames,omitempty"`         // multiple segment names (comma-separated)
	CXP                  string `json:"cxp,omitempty"`                  // Cloud Exchange Point filter
	ConnectorID          string `json:"connectorId,omitempty"`          // connector ID filter
	Offset               int    `json:"offset,omitempty"`               // pagination offset
	Limit                int    `json:"limit,omitempty"`                // pagination limit
	Search               string `json:"search,omitempty"`               // search filter
	PrefixType           string `json:"prefixType,omitempty"`           // prefix type filter
	RouteType            string `json:"routeType,omitempty"`            // route type filter
	SegmentID            string `json:"segmentId,omitempty"`            // segment ID filter
	OverlapType          string `json:"overlapType,omitempty"`          // overlap type filter
	SourceCXP            string `json:"sourceCXP,omitempty"`            // source CXP filter
	EntityInstance       string `json:"entityInstance,omitempty"`       // entity instance filter
	EntityType           string `json:"entityType,omitempty"`           // entity type filter
	Group                string `json:"group,omitempty"`                // group filter
	SegmentResourceShare string `json:"segmentResourceShare,omitempty"` // segment resource share filter
	RouteRecvType        string `json:"routeRecvType,omitempty"`        // route receive type filter
	Prefix               string `json:"prefix,omitempty"`               // prefix filter
	LPMPrefix            string `json:"lpmPrefix,omitempty"`            // LPM prefix filter
}

// Route count query parameters
type RouteCountQueryParams struct {
	Type                 string `json:"type"`                           // route type filter
	SegmentName          string `json:"segmentName,omitempty"`          // segment name filter
	SegmentNames         string `json:"segmentNames,omitempty"`         // multiple segment names (comma-separated)
	CXP                  string `json:"cxp,omitempty"`                  // Cloud Exchange Point filter
	ConnectorID          string `json:"connectorId,omitempty"`          // connector ID filter
	SegmentID            string `json:"segmentId,omitempty"`            // segment ID filter
	RouteRecvType        string `json:"routeRecvType,omitempty"`        // route receive type filter
	OverlapType          string `json:"overlapType,omitempty"`          // overlap type filter
	SourceCXP            string `json:"sourceCXP,omitempty"`            // source CXP filter
	EntityInstance       string `json:"entityInstance,omitempty"`       // entity instance filter
	EntityType           string `json:"entityType,omitempty"`           // entity type filter
	Group                string `json:"group,omitempty"`                // group filter
	SegmentResourceShare string `json:"segmentResourceShare,omitempty"` // segment resource share filter
	Search               string `json:"search,omitempty"`               // search filter
	Prefix               string `json:"prefix,omitempty"`               // prefix filter
	LPMPrefix            string `json:"lpmPrefix,omitempty"`            // LPM prefix filter
	PrefixType           string `json:"prefixType,omitempty"`           // prefix type filter
}

// Route UI connector represents a connector in a route
type RouteUIConnector struct {
	Connector struct {
		ConnectorName         string `json:"connectorName"`
		ConnectorInstanceName string `json:"connectorInstanceName"`
		ConnectorType         string `json:"connectorType"`
		ConnectorID           int    `json:"connectorId"`
		ConnectorTag          string `json:"connectorTag"`
		ConnectorCXPID        string `json:"connectorCxpId"`
		ConnectorCXPName      string `json:"connectorCxpName"`
		ConnectorRouteCXP     string `json:"connectorRouteCxp"`
		ConnectorGroup        string `json:"connectorGroup"`
	} `json:"connector"`
	PrefixType          string      `json:"prefixType"`
	ConnRouteType       string      `json:"connRouteType"`
	ConnOriginalPrefix  string      `json:"connOriginalPrefix"`
	ResourceShareName   string      `json:"resourceShareName"`
	RouteSuppressed     bool        `json:"routeSuppressed"`
	OverlappingEntities interface{} `json:"overlappingEntities"`
}

// Route UI result represents a single route result
type RouteUIResult struct {
	TenantNetworkID        int                `json:"tenantNetworkId"`
	Prefix                 string             `json:"prefix"`
	SegmentName            string             `json:"segmentName"`
	VrfName                string             `json:"vrfName"`
	CxpName                string             `json:"cxpName"`
	RouteType              string             `json:"routeType"`
	OriginalPrefix         string             `json:"originalPrefix"`
	NatMetadata            interface{}        `json:"natMetadata"`
	OverlapCxps            []string           `json:"overlapCxps"`
	OverlapInvalidateNodes interface{}        `json:"overlapInvalidateNodes"`
	OverlappedRSMetadata   interface{}        `json:"overlappedRsMetadata"`
	CXPToOverlapEntities   interface{}        `json:"cxpToOverlapEntities"`
	OverlapEntityToCXPs    interface{}        `json:"overlapEntityToCXPs"`
	Connectors             []RouteUIConnector `json:"connectors"`
	SegShareOverlapping    interface{}        `json:"segShareOverlapping"`
}

// Routes response from API
type RoutesUIResponse struct {
	Data                 []RouteUIResult `json:"data"`
	Pagination           PaginationData  `json:"pagination"`
	LatestRouteTimestamp int64           `json:"latestRouteTimestamp"`
}

// Pagination data for routes response
type PaginationData struct {
	Offset int `json:"offset"`
	Limit  int `json:"limit"`
	Hits   int `json:"hits"`
}

// Route count response from API
type RouteCountUIResult struct {
	TenantNetworkID int `json:"tenantNetworkId"`
	Count           int `json:"count"`
}

// GetRoutes retrieves routes for a tenant network with optional filtering
func (ac *AlkiraClient) GetRoutes(params RouteQueryParams) (*RoutesUIResponse, error) {
	// Construct the URI
	uri := fmt.Sprintf("%s/tenantnetworks/%s/routes", ac.URI, ac.TenantNetworkId)

	// Build query parameters
	queryParams := url.Values{}

	// Add all parameters if they are not empty
	if params.Type != "" {
		queryParams.Set("type", params.Type)
	}
	if params.SegmentName != "" {
		queryParams.Set("segmentName", params.SegmentName)
	}
	if params.SegmentNames != "" {
		queryParams.Set("segmentNames", params.SegmentNames)
	}
	if params.CXP != "" {
		queryParams.Set("cxp", params.CXP)
	}
	if params.ConnectorID != "" {
		queryParams.Set("connectorId", params.ConnectorID)
	}
	if params.Offset != 0 {
		queryParams.Set("offset", strconv.Itoa(params.Offset))
	}
	if params.Limit != 0 {
		queryParams.Set("limit", strconv.Itoa(params.Limit))
	}
	if params.Search != "" {
		queryParams.Set("search", params.Search)
	}
	if params.PrefixType != "" {
		queryParams.Set("prefixType", params.PrefixType)
	}
	if params.RouteType != "" {
		queryParams.Set("routeType", params.RouteType)
	}
	if params.SegmentID != "" {
		queryParams.Set("segmentId", params.SegmentID)
	}
	if params.OverlapType != "" {
		queryParams.Set("overlapType", params.OverlapType)
	}
	if params.SourceCXP != "" {
		queryParams.Set("sourceCXP", params.SourceCXP)
	}
	if params.EntityInstance != "" {
		queryParams.Set("entityInstance", params.EntityInstance)
	}
	if params.EntityType != "" {
		queryParams.Set("entityType", params.EntityType)
	}
	if params.Group != "" {
		queryParams.Set("group", params.Group)
	}
	if params.SegmentResourceShare != "" {
		queryParams.Set("segmentResourceShare", params.SegmentResourceShare)
	}
	if params.RouteRecvType != "" {
		queryParams.Set("routeRecvType", params.RouteRecvType)
	}
	if params.Prefix != "" {
		queryParams.Set("prefix", params.Prefix)
	}
	if params.LPMPrefix != "" {
		queryParams.Set("lpmPrefix", params.LPMPrefix)
	}

	// Append query parameters to URI if any
	if len(queryParams) > 0 {
		uri += "?" + queryParams.Encode()
	}

	// Make the request
	data, _, err := ac.get(uri)
	if err != nil {
		return nil, fmt.Errorf("failed to get routes: %v", err)
	}

	// Parse response
	var routes RoutesUIResponse
	err = json.Unmarshal([]byte(data), &routes)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal routes response: %v", err)
	}

	return &routes, nil
}

// GetRouteCount retrieves route count for a tenant network with optional filtering
func (ac *AlkiraClient) GetRouteCount(params RouteCountQueryParams) (*RouteCountUIResult, error) {
	// Construct the URI
	uri := fmt.Sprintf("%s/tenantnetworks/%s/route-count", ac.URI, ac.TenantNetworkId)

	// Build query parameters
	queryParams := url.Values{}

	// Add all parameters if they are not empty
	if params.Type != "" {
		queryParams.Set("type", params.Type)
	}
	if params.SegmentName != "" {
		queryParams.Set("segmentName", params.SegmentName)
	}
	if params.SegmentNames != "" {
		queryParams.Set("segmentNames", params.SegmentNames)
	}
	if params.CXP != "" {
		queryParams.Set("cxp", params.CXP)
	}
	if params.ConnectorID != "" {
		queryParams.Set("connectorId", params.ConnectorID)
	}
	if params.SegmentID != "" {
		queryParams.Set("segmentId", params.SegmentID)
	}
	if params.RouteRecvType != "" {
		queryParams.Set("routeRecvType", params.RouteRecvType)
	}
	if params.OverlapType != "" {
		queryParams.Set("overlapType", params.OverlapType)
	}
	if params.SourceCXP != "" {
		queryParams.Set("sourceCXP", params.SourceCXP)
	}
	if params.EntityInstance != "" {
		queryParams.Set("entityInstance", params.EntityInstance)
	}
	if params.EntityType != "" {
		queryParams.Set("entityType", params.EntityType)
	}
	if params.Group != "" {
		queryParams.Set("group", params.Group)
	}
	if params.SegmentResourceShare != "" {
		queryParams.Set("segmentResourceShare", params.SegmentResourceShare)
	}
	if params.Search != "" {
		queryParams.Set("search", params.Search)
	}
	if params.Prefix != "" {
		queryParams.Set("prefix", params.Prefix)
	}
	if params.LPMPrefix != "" {
		queryParams.Set("lpmPrefix", params.LPMPrefix)
	}
	if params.PrefixType != "" {
		queryParams.Set("prefixType", params.PrefixType)
	}

	// Append query parameters to URI if any
	if len(queryParams) > 0 {
		uri += "?" + queryParams.Encode()
	}

	// Make the request
	data, _, err := ac.get(uri)
	if err != nil {
		return nil, fmt.Errorf("failed to get route count: %v", err)
	}

	// Parse response
	var routeCount RouteCountUIResult
	err = json.Unmarshal([]byte(data), &routeCount)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal route count response: %v", err)
	}

	return &routeCount, nil
}
