package alkira

import (
	"errors"

	"github.com/alkiranet/alkira-client-go/alkira"
)

func expandAzureErInstances(in []interface{}, m interface{}) ([]alkira.ConnectorAzureErInstance, error) {
	if in == nil || len(in) == 0 {
		return nil, errors.New("Invalid Azure Er Instance input")
	}

	instances := make([]alkira.ConnectorAzureErInstance, len(in))
	for i, instance := range in {
		r := alkira.ConnectorAzureErInstance{}
		instanceCfg := instance.(map[string]interface{})

		if v, ok := instanceCfg["name"].(string); ok {
			r.Name = v
		}
		if v, ok := instanceCfg["express_route_circuit_id"].(string); ok {
			r.ExpressRouterCircuitId = v
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
		if v, ok := instanceCfg["vnis"].([]int); ok {
			r.Vnis = v
		}
		instances[i] = r
	}

	return instances, nil
}

func expandAzureErSegments(seg []interface{}, m interface{}) ([]alkira.ConnectorAzureErSegment, error) {
	if seg == nil || len(seg) == 0 {
		return nil, errors.New("Invalid Azure Er Segment Options input")
	}

	segments := make([]alkira.ConnectorAzureErSegment, len(seg))
	for i, segment := range seg {
		r := alkira.ConnectorAzureErSegment{}
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
			r.DisabledInternetExit = v
		}
		if v, ok := instanceCfg["advertise_on_prem_routes"].(bool); ok {
			r.AdvertiseOnPremRoutes = v
		}
		segments[i] = r
	}

	return segments, nil
}
