package alkira

import (
	"encoding/json"
	"fmt"
)

type F5vServerEndpoint struct {
	F5ServiceId          int                            `json:"f5ServiceId"`
	Name                 string                         `json:"name"`
	Id                   json.Number                    `json:"id,omitempty"` // RESPONSE ONLY
	Type                 string                         `json:"type"`
	Segment              string                         `json:"segment"`
	FqdnPrefix           string                         `json:"fqdnPrefix,omitempty"`
	Protocol             string                         `json:"protocol"`
	Snat                 string                         `json:"snat"`
	F5ServiceInstanceIds []int                          `json:"f5ServiceInstanceIds,omitempty"`
	PortRanges           []string                       `json:"portRanges,omitempty"`
	DestinationEndpoints *F5VServerDestinationEndpoints `json:"destinationEndpoints,omitempty"`
	Metadata             *F5VServerEndpointMetadata     `json:"metadata,omitempty"` // RESPONSE ONLY
}

type F5VServerDestinationEndpoints struct {
	PortRanges  []string `json:"portRanges"`
	IpAddresses []string `json:"ipAddresses"`
}

type F5VServerEndpointInstanceMetadata struct {
	ElasticIp     string `json:"elasticIp,omitempty"`
	SecondaryIp   string `json:"secondaryIp,omitempty"`
	RouteDomainId int64  `json:"routeDomainId,omitempty"`
	SegmentId     int64  `json:"segmentId,omitempty"`
	Vlan          string `json:"vlan,omitempty"`
	EcmpPoolName  string `json:"ecmpPoolName,omitempty"`
	VirtualIp     string `json:"virtualIp,omitempty"`
}

type F5VServerEndpointMetadata struct {
	InstanceToMetadata map[string]F5VServerEndpointInstanceMetadata `json:"instanceToMetadata,omitempty"`
}

func NewF5vServerEndpoint(ac *AlkiraClient) *AlkiraAPI[F5vServerEndpoint] {
	uri := fmt.Sprintf("%s/tenantnetworks/%s/f5-vserver-endpoints", ac.URI, ac.TenantNetworkId)
	api := &AlkiraAPI[F5vServerEndpoint]{ac, uri, true}
	return api

}
