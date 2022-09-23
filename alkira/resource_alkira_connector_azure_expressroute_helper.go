package alkira

import (
	"errors"

	"github.com/alkiranet/alkira-client-go/alkira"
)

func expandAzureExpressRouteInstances(in []interface{}, m interface{}) ([]alkira.ConnectorAzureExpressRouteInstance, error) {
	if in == nil || len(in) == 0 {
		return nil, errors.New("Invalid Azure ExpressRoute Instance input")
	}

	instances := make([]alkira.ConnectorAzureExpressRouteInstance, len(in))
	for i, instance := range in {
		r := alkira.ConnectorAzureExpressRouteInstance{}
		instanceCfg := instance.(map[string]interface{})
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
		if v, ok := instanceCfg["gateway_mac_address"].([]string); ok {
			r.GatewayMacAddress = v
		}
		if v, ok := instanceCfg["virtual_network_interface"].([]int); ok {
			r.Vnis = v
		}
		instances[i] = r
	}

	return instances, nil
}

func expandAzureExpressRouteSegments(seg []interface{}, m interface{}) ([]alkira.ConnectorAzureExpressRouteSegment, error) {
	if seg == nil || len(seg) == 0 {
		return nil, errors.New("Invalid Azure ExpresRoute Segment Options input")
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
		if v, ok := instanceCfg["customer_asn"].(int); ok {
			r.CustomerAsn = v
		}
		if v, ok := instanceCfg["disabled_internet_exit"].(bool); ok {
			r.DisableInternetExit = v
		}
		if v, ok := instanceCfg["advertise_on_prem_routes"].(bool); ok {
			r.AdvertiseOnPremRoutes = v
		}
		segments[i] = r
	}

	return segments, nil
}
