package alkira

import (
	"errors"
	"log"

	"github.com/alkiranet/alkira-client-go/alkira"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func expandGcpInterconnectCustomerGateways(d []interface{}) ([]alkira.ConnectorGcpInterconnectCustomerGateway, error) {
	if d == nil || len(d) == 0 {
		log.Printf("[ERROR] invalid GCP interconnect customer gateway input")
		return nil, errors.New("[ERROR] invalid GCP interconnect customer gateway input")
	}
	customerGateways := make([]alkira.ConnectorGcpInterconnectCustomerGateway, len(d))
	for i, c := range d {
		cfgCustomerGateway := c.(map[string]interface{})
		newCustomerGateway := alkira.ConnectorGcpInterconnectCustomerGateway{}
		if v, ok := cfgCustomerGateway["loopback_ip"].(string); ok {
			newCustomerGateway.LoopbackIp = v
		}
		if v, ok := cfgCustomerGateway["tunnel_count"].(int); ok {
			newCustomerGateway.TunnelCount = v
		}
		customerGateways[i] = newCustomerGateway
	}
	return customerGateways, nil
}

func expandGcpInterconnectSegmentOptions(d []interface{}, m interface{}) ([]alkira.ConnectorGcpInterconnectSegmentOption, error) {
	if d == nil || len(d) == 0 {
		log.Printf("[ERROR] invalid GCP interconnect segment option input")
		return nil, errors.New("[ERROR] invalid GCP interconnect segment option input")
	}

	segmentOptions := make([]alkira.ConnectorGcpInterconnectSegmentOption, len(d))
	for i, s := range d {
		cfgSegmentOption := s.(map[string]interface{})
		newSegmentOption := alkira.ConnectorGcpInterconnectSegmentOption{}
		if v, ok := cfgSegmentOption["segment_id"].(string); ok {
			segmentName, err := getSegmentNameById(v, m)
			if err != nil {
				return nil, err
			}
			newSegmentOption.SegmentName = segmentName
		}
		if v, ok := cfgSegmentOption["advertise_on_prem_routes"].(bool); ok {
			newSegmentOption.AdvertiseOnPremRoutes = v
		}
		if v, ok := cfgSegmentOption["advertise_default_route"].(bool); ok {
			// advertise_default_route is negation of disable_internet_exit
			newSegmentOption.DisableInternetExit = !v
		}

		if v, ok := cfgSegmentOption["customer_gateways"].([]interface{}); ok {
			customerGateways, err := expandGcpInterconnectCustomerGateways(v)
			if err != nil {
				return nil, err
			}
			newSegmentOption.CustomerGateways = customerGateways
		}
		segmentOptions[i] = newSegmentOption

	}
	return segmentOptions, nil
}

func expandGcpInterconnectInstances(in []interface{}, m interface{}) ([]alkira.ConnectorGcpInterconnectInstance, error) {
	if in == nil || len(in) == 0 {
		log.Printf("[ERROR] invalid GCP interconnect instance input")
		return nil, errors.New("[ERROR] invalid GCP interconnect instance input")
	}

	instances := make([]alkira.ConnectorGcpInterconnectInstance, len(in))

	// loop over the instances from the config and copy the values from the config to the struct
	// to create the API payload
	for i, instance := range in {
		newInstance := alkira.ConnectorGcpInterconnectInstance{}
		cfgInstance := instance.(map[string]interface{})

		if v, ok := cfgInstance["id"].(int); ok {
			newInstance.Id = v
		}
		if v, ok := cfgInstance["name"].(string); ok {
			newInstance.Name = v
		}
		if v, ok := cfgInstance["edge_availability_domain"].(string); ok {
			newInstance.GcpEdgeAvailabilityDomain = v
		}
		if v, ok := cfgInstance["bgp_auth_key"].(string); ok {
			newInstance.BgpAuthKeyAlkira = v
		}
		if v, ok := cfgInstance["gateway_mac_address"].(string); ok {
			newInstance.GatewayMacAddress = v
		}
		if v, ok := cfgInstance["customer_asn"].(int); ok {
			newInstance.CustomerAsn = v
		}
		if v, ok := cfgInstance["vni_id"].(int); ok {
			newInstance.Vni = v
		}
		if v, ok := cfgInstance["segment_options"].([]interface{}); ok {
			segmentOptions, err := expandGcpInterconnectSegmentOptions(v, m)
			if err != nil {
				return nil, err
			}
			newInstance.SegmentOptions = segmentOptions
		}

		instances[i] = newInstance
	}
	return instances, nil
}

func setGcpInterconnectSegmentOptions(instance alkira.ConnectorGcpInterconnectInstance, m interface{}) ([]map[string]interface{}, error) {
	var segmentOptions []map[string]interface{}
	sO := instance.SegmentOptions

	// loop over the segment options
	for _, aSegmentOption := range sO {
		segmentId, err := getSegmentIdByName(aSegmentOption.SegmentName, m)
		if err != nil {
			log.Printf("[ERROR] error getting segment ID for Segment Name %v", aSegmentOption.SegmentName)
			return nil, err
		}

		// create a list of map for customer gateways
		var customerGateways []map[string]interface{}
		for _, aCustomerGateway := range aSegmentOption.CustomerGateways {
			// add all gateways to the list
			customerGateway := map[string]interface{}{
				"loopback_ip":  aCustomerGateway.LoopbackIp,
				"tunnel_count": aCustomerGateway.TunnelCount,
			}
			customerGateways = append(customerGateways, customerGateway)
		}
		segmentOption := map[string]interface{}{
			"segment_id":               segmentId,
			"advertise_on_prem_routes": aSegmentOption.AdvertiseOnPremRoutes,
			"advertise_default_route":  !aSegmentOption.DisableInternetExit,
			"customer_gateways":        customerGateways,
		}
		segmentOptions = append(segmentOptions, segmentOption)
	}
	return segmentOptions, nil
}

func setGcpInterconnectInstance(d *schema.ResourceData, ins []alkira.ConnectorGcpInterconnectInstance, m interface{}) []map[string]interface{} {
	var instances []map[string]interface{}
	for _, in := range ins {
		instanceSegmentOptions, err := setGcpInterconnectSegmentOptions(in, m)
		if err != nil {
			log.Printf("[ERROR] error setting segment options")
			return nil
		}
		instance := map[string]interface{}{
			"id":                       in.Id,
			"name":                     in.Name,
			"edge_availability_domain": in.GcpEdgeAvailabilityDomain,
			"customer_asn":             in.CustomerAsn,
			"bgp_auth_key":             in.BgpAuthKeyAlkira,
			"gateway_mac_address":      in.GatewayMacAddress,
			"vni_id":                   in.Vni,
			"segment_options":          instanceSegmentOptions,
		}
		instances = append(instances, instance)

	}
	return instances

}

func generateGcpInterconnectRequest(d *schema.ResourceData, m interface{}) (*alkira.ConnectorGcpInterconnect, error) {
	instances, err := expandGcpInterconnectInstances(d.Get("instances").([]interface{}), m)
	if err != nil {
		return nil, err
	}

	// Assemble request
	connector := &alkira.ConnectorGcpInterconnect{
		Name:             d.Get("name").(string),
		Size:             d.Get("size").(string),
		Description:      d.Get("description").(string),
		Cxp:              d.Get("cxp").(string),
		Enabled:          d.Get("enabled").(bool),
		Group:            d.Get("group").(string),
		TunnelProtocol:   d.Get("tunnel_protocol").(string),
		BillingTags:      convertTypeSetToIntList(d.Get("billing_tag_ids").(*schema.Set)),
		LoopbackPrefixes: convertTypeSetToStringList(d.Get("loopback_prefixes").(*schema.Set)),
		Instances:        instances,
		ScaleGroupId:     d.Get("scale_group_id").(string),
		ImplicitGroupId:  d.Get("implicit_group_id").(int),
	}

	return connector, nil
}
