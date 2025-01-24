package alkira

import (
	"errors"

	"github.com/alkiranet/alkira-client-go/alkira"
)

// expandAzureExpressRouteInstances expands the instances for Azure ExpressRoute connector
func expandAzureExpressRouteInstances(in []interface{}, m interface{}) ([]alkira.ConnectorAzureExpressRouteInstance, error) {
	if in == nil || len(in) == 0 {
		return nil, errors.New("Invalid Azure ExpressRoute Instance input")
	}

	instances := make([]alkira.ConnectorAzureExpressRouteInstance, len(in))
	for i, instance := range in {
		r := alkira.ConnectorAzureExpressRouteInstance{}
		instanceCfg := instance.(map[string]interface{})

		// Basic fields
		if v, ok := instanceCfg["id"].(int); ok {
			r.Id = v
		}
		if v, ok := instanceCfg["name"].(string); ok {
			r.Name = v
		}
		if v, ok := instanceCfg["expressroute_circuit_id"].(string); ok {
			r.ExpressRouteCircuitId = v
		}
		if v, ok := instanceCfg["redundant_router"].(bool); ok {
			r.RedundantRouter = v
		}
		if v, ok := instanceCfg["loopback_subnet"].(string); ok {
			r.LoopbackSubnet = v
		}
		if v, ok := instanceCfg["credential_id"].(string); ok {
			r.CredentialId = v
		}

		// Gateway MAC Addresses
		gatewayMacAddresses := []string{}
		if v, ok := instanceCfg["gateway_mac_address"].([]interface{}); ok {
			for _, addr := range v {
				gatewayMacAddresses = append(gatewayMacAddresses, addr.(string))
			}
		}
		r.GatewayMacAddress = gatewayMacAddresses

		// Virtual Network Interfaces (VNIs)
		vnis := []int{}
		if v, ok := instanceCfg["virtual_network_interface"].([]interface{}); ok {
			for _, vni := range v {
				vnis = append(vnis, vni.(int))
			}
		}
		r.Vnis = vnis

		// Segment Options
		if v, ok := instanceCfg["segment_options"].([]interface{}); ok {
			segmentOptions, err := expandInstanceSegmentOptions(v)
			if err != nil {
				return nil, err
			}
			r.SegmentOptions = segmentOptions
		}

		instances[i] = r
	}

	return instances, nil
}

// expandInstanceSegmentOptions expands the segment options for an instance
func expandInstanceSegmentOptions(in []interface{}) ([]alkira.InstanceSegmentOption, error) {
	segmentOptions := make([]alkira.InstanceSegmentOption, 0)
	for _, segOpt := range in {
		segOptMap := segOpt.(map[string]interface{})
		segmentOption := alkira.InstanceSegmentOption{}

		if v, ok := segOptMap["segment_name"].(string); ok {
			segmentOption.SegmentName = v
		}

		// Customer Gateways
		customerGateways, err := expandCustomerGateways(segOptMap["customer_gateways"].([]interface{}))
		if err != nil {
			return nil, err
		}
		segmentOption.CustomerGateways = customerGateways

		segmentOptions = append(segmentOptions, segmentOption)
	}
	return segmentOptions, nil
}

// expandCustomerGateways expands the customer gateways for a segment option
func expandCustomerGateways(in []interface{}) ([]alkira.CustomerGateway, error) {
	customerGateways := make([]alkira.CustomerGateway, 0)
	for _, cg := range in {
		cgMap := cg.(map[string]interface{})
		customerGateway := alkira.CustomerGateway{}

		if v, ok := cgMap["name"].(string); ok {
			customerGateway.Name = v
		}

		// Tunnels
		tunnels, err := expandCustomerGatewayTunnels(cgMap["tunnels"].([]interface{}))
		if err != nil {
			return nil, err
		}
		customerGateway.Tunnels = tunnels

		customerGateways = append(customerGateways, customerGateway)
	}
	return customerGateways, nil
}

// expandCustomerGatewayTunnels expands the tunnels for a customer gateway
func expandCustomerGatewayTunnels(in []interface{}) ([]alkira.CustomerGatewayTunnel, error) {
	tunnels := make([]alkira.CustomerGatewayTunnel, 0)
	for _, t := range in {
		tMap := t.(map[string]interface{})
		tunnel := alkira.CustomerGatewayTunnel{}

		if v, ok := tMap["name"].(string); ok {
			tunnel.Name = v
		}
		if v, ok := tMap["initiator"].(bool); ok {
			tunnel.Initiator = v
		}
		if v, ok := tMap["profile_id"].(int); ok {
			tunnel.ProfileId = v
		}
		if v, ok := tMap["ike_version"].(string); ok {
			tunnel.IkeVersion = v
		}
		if v, ok := tMap["pre_shared_key"].(string); ok {
			tunnel.PreSharedKey = v
		}
		if v, ok := tMap["remote_auth_type"].(string); ok {
			tunnel.RemoteAuthType = v
		}
		if v, ok := tMap["remote_auth_value"].(string); ok {
			tunnel.RemoteAuthValue = v
		}

		tunnels = append(tunnels, tunnel)
	}
	return tunnels, nil
}

// expandAzureExpressRouteSegments expands the segment options for Azure ExpressRoute connector
func expandAzureExpressRouteSegments(seg []interface{}, m interface{}) ([]alkira.ConnectorAzureExpressRouteSegment, error) {
	if seg == nil || len(seg) == 0 {
		return nil, errors.New("Invalid Azure ExpressRoute Segment Options input")
	}

	segments := make([]alkira.ConnectorAzureExpressRouteSegment, len(seg))
	for i, segment := range seg {
		r := alkira.ConnectorAzureExpressRouteSegment{}
		instanceCfg := segment.(map[string]interface{})
		if v, ok := instanceCfg["segment_name"].(string); ok {
			r.SegmentName = v
		}
		if v, ok := instanceCfg["customer_asn"].(int); ok {
			r.CustomerAsn = v
		}
		if v, ok := instanceCfg["disable_internet_exit"].(bool); ok {
			r.DisableInternetExit = v
		}
		if v, ok := instanceCfg["advertise_on_prem_routes"].(bool); ok {
			r.AdvertiseOnPremRoutes = v
		}
		segments[i] = r
	}

	return segments, nil
}

// flattenInstance flattens an Azure ExpressRoute instance into a Terraform schema-compatible map
func flattenInstance(instance alkira.ConnectorAzureExpressRouteInstance) map[string]interface{} {
	result := map[string]interface{}{
		"credential_id":             instance.CredentialId,
		"expressroute_circuit_id":   instance.ExpressRouteCircuitId,
		"gateway_mac_address":       instance.GatewayMacAddress,
		"id":                        instance.Id,
		"loopback_subnet":           instance.LoopbackSubnet,
		"name":                      instance.Name,
		"redundant_router":          instance.RedundantRouter,
		"virtual_network_interface": instance.Vnis,
		"segment_options":           flattenInstanceSegmentOptions(instance.SegmentOptions),
	}
	return result
}

// flattenInstanceSegmentOptions flattens the segment options for an instance
func flattenInstanceSegmentOptions(segmentOptions []alkira.InstanceSegmentOption) []interface{} {
	if segmentOptions == nil {
		return nil
	}

	result := make([]interface{}, len(segmentOptions))
	for i, segOpt := range segmentOptions {
		s := map[string]interface{}{
			"segment_name":      segOpt.SegmentName,
			"customer_gateways": flattenCustomerGateways(segOpt.CustomerGateways),
		}
		result[i] = s
	}
	return result
}

// flattenCustomerGateways flattens the customer gateways for a segment option
func flattenCustomerGateways(customerGateways []alkira.CustomerGateway) []interface{} {
	if customerGateways == nil {
		return nil
	}

	result := make([]interface{}, len(customerGateways))
	for i, cg := range customerGateways {
		c := map[string]interface{}{
			"name":    cg.Name,
			"tunnels": flattenCustomerGatewayTunnels(cg.Tunnels),
		}
		result[i] = c
	}
	return result
}

// flattenCustomerGatewayTunnels flattens the tunnels for a customer gateway
func flattenCustomerGatewayTunnels(tunnels []alkira.CustomerGatewayTunnel) []interface{} {
	if tunnels == nil {
		return nil
	}

	result := make([]interface{}, len(tunnels))
	for i, t := range tunnels {
		tunnel := map[string]interface{}{
			"name":              t.Name,
			"initiator":         t.Initiator,
			"profile_id":        t.ProfileId,
			"ike_version":       t.IkeVersion,
			"pre_shared_key":    t.PreSharedKey,
			"remote_auth_type":  t.RemoteAuthType,
			"remote_auth_value": t.RemoteAuthValue,
		}
		result[i] = tunnel
	}
	return result
}
