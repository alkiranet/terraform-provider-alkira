package alkira

import (
	"encoding/json"
	"fmt"
)

type InterCxpRoutePolicy struct {
	Name        string                    `json:"name"`
	Description string                    `json:"description,omitempty"`
	Enabled     bool                      `json:"enabled"`
	Direction   string                    `json:"direction"`
	Segment     string                    `json:"segment"`
	SourceCxps  []string                  `json:"sourceCxps"`
	DestCxps    []string                  `json:"destCxps"`
	Id          json.Number               `json:"id,omitempty"` // response only
	Rules       []InterCxpRoutePolicyRule `json:"rules,omitempty"`
}

type InterCxpRoutePolicyRule struct {
	Action     string                        `json:"action"`
	Name       string                        `json:"name"`
	Match      InterCxpRoutePolicyRuleMatch  `json:"match"`
	SequenceNo int                           `json:"sequenceNo,omitempty"` // response only
	Set        *InterCxpRoutePolicyRuleSet   `json:"set,omitempty"`
}

type InterCxpRoutePolicyRuleMatch struct {
	All                      bool  `json:"all"`
	PrefixListIds            []int `json:"prefixListIds,omitempty"`
	CommunityListIds         []int `json:"communityListIds,omitempty"`
	ExtendedCommunityListIds []int `json:"extendedCommunityListIds,omitempty"`
	AsPathListIds            []int `json:"asPathListIds,omitempty"`
	SegmentResourceIds       []int `json:"segmentResourceIds,omitempty"`
	ConnectorGroupIds        []int `json:"connectorGroupIds,omitempty"`
}

type InterCxpRoutePolicyRuleSet struct {
	AsPathPrepend     string `json:"asPathPrepend,omitempty"`
	Community         string `json:"community,omitempty"`
	ExtendedCommunity string `json:"extendedCommunity,omitempty"`
}

// NewInterCxpRoutePolicy creates a new inter-CXP route policy API client
func NewInterCxpRoutePolicy(ac *AlkiraClient) *AlkiraAPI[InterCxpRoutePolicy] {
	uri := fmt.Sprintf("%s/tenantnetworks/%s/intercxp-route-policies", ac.URI, ac.TenantNetworkId)
	api := &AlkiraAPI[InterCxpRoutePolicy]{ac, uri, true}
	return api
}
