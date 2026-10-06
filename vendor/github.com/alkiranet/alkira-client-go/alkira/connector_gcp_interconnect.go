package alkira

import (
	"encoding/json"
	"fmt"
)

type ConnectorGcpInterconnect struct {
	ScaleGroupId     string                             `json:"scaleGroupId,omitempty"`
	Description      string                             `json:"description,omitempty"`
	Cxp              string                             `json:"cxp"`
	Group            string                             `json:"group,omitempty"`
	Size             string                             `json:"size"`
	TunnelProtocol   string                             `json:"tunnelProtocol"`
	Name             string                             `json:"name"`
	Id               json.Number                        `json:"id,omitempty"` // RESPONSE ONLY
	BillingTags      []int                              `json:"billingTags,omitempty"`
	LoopbackPrefixes []string                           `json:"loopbackPrefixes"`
	Instances        []ConnectorGcpInterconnectInstance `json:"instances"`
	ImplicitGroupId  int                                `json:"implicitGroupId,omitempty"` // RESPONSE ONLY
	Enabled          bool                               `json:"enabled"`
}

type ConnectorGcpInterconnectInstance struct {
	Name                      string                                  `json:"name"`
	GcpEdgeAvailabilityDomain string                                  `json:"gcpEdgeAvailabilityDomain"`
	BgpAuthKeyAlkira          string                                  `json:"bgpAuthKeyAlkira,omitempty"`
	GatewayMacAddress         string                                  `json:"gatewayMacAddress,omitempty"`
	SegmentOptions            []ConnectorGcpInterconnectSegmentOption `json:"segmentOptions"`
	Id                        int                                     `json:"id,omitempty"` // RESPONSE ONLY
	CustomerAsn               int                                     `json:"customerAsn"`
	Vni                       int                                     `json:"vni,omitempty"`
}

type ConnectorGcpInterconnectSegmentOption struct {
	SegmentName           string                                    `json:"segmentName"`
	CustomerGateways      []ConnectorGcpInterconnectCustomerGateway `json:"customerGateways"`
	AdvertiseOnPremRoutes bool                                      `json:"advertiseOnPremRoutes,omitempty"`
	DisableInternetExit   bool                                      `json:"disableInternetExit"`
}

type ConnectorGcpInterconnectCustomerGateway struct {
	LoopbackIp  string `json:"loopbackIp,omitempty"`
	TunnelCount int    `json:"tunnelCount"`
}

func NewConnectorGcpInterconnect(ac *AlkiraClient) *AlkiraAPI[ConnectorGcpInterconnect] {
	uri := fmt.Sprintf("%s/tenantnetworks/%s/gcp-interconnect-connectors", ac.URI, ac.TenantNetworkId)
	api := &AlkiraAPI[ConnectorGcpInterconnect]{ac, uri, true}
	return api
}
