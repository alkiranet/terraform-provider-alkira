package alkira

import (
	"errors"
	"fmt"
	"log"

	"github.com/alkiranet/alkira-client-go/alkira"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func expandGcpInterconnectCustomerGateways(cg []interface{}) ([]alkira.ConnectorGcpInterconnectCustomerGateway, error) {
	if cg == nil || len(cg) == 0 {
		log.Printf("[ERROR] invalid GCP interconnect customer gateway input")
		return nil, errors.New("[ERROR] invalid GCP interconnect customer gateway input")
	}
	customerGateways := make([]alkira.ConnectorGcpInterconnectCustomerGateway, len(cg))
	for i, c := range cg {
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

func expandGcpInterconnectSegmentOptions(so []interface{}, instanceName string, m interface{}) ([]alkira.ConnectorGcpInterconnectSegmentOption, error) {
	if so == nil || len(so) == 0 {
		log.Printf("[ERROR] invalid GCP interconnect segment option input")
		return nil, errors.New("[ERROR] invalid GCP interconnect segment option input")
	}

	segmentOptions := make([]alkira.ConnectorGcpInterconnectSegmentOption, 0)
	for _, s := range so {
		cfgSegmentOption := s.(map[string]interface{})
		if cfgSegmentOption["instance_name"] == instanceName {
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
			if v, ok := cfgSegmentOption["disable_internet_exit"].(bool); ok {
				newSegmentOption.DisableInternetExit = v
			}

			if v, ok := cfgSegmentOption["customer_gateways"].([]interface{}); ok {
				customerGateways, err := expandGcpInterconnectCustomerGateways(v)
				if err != nil {
					return nil, err
				}
				newSegmentOption.CustomerGateways = customerGateways
			}
			if newSegmentOption.SegmentName != "" ||
				newSegmentOption.AdvertiseOnPremRoutes ||
				newSegmentOption.DisableInternetExit ||
				len(newSegmentOption.CustomerGateways) > 0 {
				segmentOptions = append(segmentOptions, newSegmentOption) // append only if we have newSegmentOption
			}
		}
	}
	return segmentOptions, nil
}

func validateGcpInterconnectInstances(in []interface{}, so []interface{}) error {
	// Create a map to track instances with segment options
	instanceSegmentMap := make(map[string]bool)

	// Iterate over segment options to populate the map
	for _, segmentOption := range so {
		cfgSegmentOption := segmentOption.(map[string]interface{})
		if instanceName, ok := cfgSegmentOption["instance_name"].(string); ok {
			instanceSegmentMap[instanceName] = true
		}
	}

	// Check each instance to ensure it has a corresponding segment option
	for _, instance := range in {
		cfgInstance := instance.(map[string]interface{})
		instanceName, ok := cfgInstance["name"].(string)
		if !ok || instanceName == "" {
			return errors.New("[ERROR] GCP interconnect instance name is required")
		}

		if _, exists := instanceSegmentMap[instanceName]; !exists {
			return fmt.Errorf("[ERROR] No segment option found for instance '%s'", instanceName)
		}
	}
	return nil
}
func expandGcpInterconnectInstances(in []interface{}, so []interface{}, m interface{}) ([]alkira.ConnectorGcpInterconnectInstance, error) {
	if in == nil || len(in) == 0 {
		log.Printf("[ERROR] invalid GCP interconnect instance input")
		return nil, errors.New("[ERROR] invalid GCP interconnect instance input")
	}
	err := validateGcpInterconnectInstances(in, so)
	if err != nil {
		return nil, err
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
		if v, ok := cfgInstance["bgp_auth_key_alkira"].(string); ok {
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
		segmentOptions, err := expandGcpInterconnectSegmentOptions(so, cfgInstance["name"].(string), m)
		if err != nil {
			return nil, err
		}
		newInstance.SegmentOptions = segmentOptions
		instances[i] = newInstance
	}
	return instances, nil
}

func setGcpInterconnectSegmentOptions(d *schema.ResourceData, instance *alkira.ConnectorGcpInterconnectInstance, m interface{}) ([]map[string]interface{}, error) {
	var segmentOptions []map[string]interface{}

	// loop over the segment options
	for _, cSegmentOption := range d.Get("segment_options").([]interface{}) {
		cfgSegmentOption := cSegmentOption.(map[string]interface{})
		for _, aSegmentOption := range instance.SegmentOptions {
			// make segmentOptions map for each instance using the instance name
			if cfgSegmentOption["instance_name"].(string) == instance.Name {
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
					"instance_name":            instance.Name,
					"segment_id":               segmentId,
					"advertise_on_prem_routes": aSegmentOption.AdvertiseOnPremRoutes,
					"disable_internet_exit":    aSegmentOption.DisableInternetExit,
					"customer_gateways":        customerGateways,
				}
				segmentOptions = append(segmentOptions, segmentOption)
			}
		}
	}
	return segmentOptions, nil
}

func setGcpInterconnectInstance(d *schema.ResourceData, connector *alkira.ConnectorGcpInterconnect, m interface{}) {
	var instances []map[string]interface{}
	var segmentOptions []map[string]interface{}
	for _, cInstance := range d.Get("instances").([]interface{}) {
		configInstance := cInstance.(map[string]interface{})
		for _, aInstance := range connector.Instances {
			if configInstance["id"].(int) == aInstance.Id ||
				configInstance["name"].(string) == aInstance.Name {
				log.Printf("[DEBUG] instance found [%v]", aInstance.Name)
				instanceSegmentOptions, err := setGcpInterconnectSegmentOptions(d, &aInstance, m)
				if err != nil {
					log.Printf("[ERROR] error setting segment options")
					return
				}
				segmentOptions = append(segmentOptions, instanceSegmentOptions...)

				instance := map[string]interface{}{
					"id":                       aInstance.Id,
					"name":                     aInstance.Name,
					"edge_availability_domain": aInstance.GcpEdgeAvailabilityDomain,
					"customer_asn":             aInstance.CustomerAsn,
					"bgp_auth_key":             aInstance.BgpAuthKeyAlkira,
					"gateway_mac_address":      aInstance.GatewayMacAddress,
					"vni_id":                   aInstance.Vni,
				}
				instances = append(instances, instance)
			}
		}
	}

	for _, aInstance := range connector.Instances {
		new := true
		for _, cInstance := range d.Get("instances").([]interface{}) {
			instanceConfig := cInstance.(map[string]interface{})
			if instanceConfig["id"].(int) == aInstance.Id ||
				instanceConfig["name"].(string) == aInstance.Name {
				new = false
			}
		}

		if new {
			instanceSegmentOptions, err := setGcpInterconnectSegmentOptions(d, &aInstance, m)
			if err != nil {
				log.Printf("[DEBUG] error setting segment options")
				return
			}
			segmentOptions = append(segmentOptions, instanceSegmentOptions...)

			i := map[string]interface{}{
				"id":                       aInstance.Id,
				"name":                     aInstance.Name,
				"edge_availability_domain": aInstance.GcpEdgeAvailabilityDomain,
				"customer_asn":             aInstance.CustomerAsn,
				"bgp_auth_key":             aInstance.BgpAuthKeyAlkira,
				"gateway_mac_address":      aInstance.GatewayMacAddress,
				"vni_id":                   aInstance.Vni,
			}
			instances = append(instances, i)
		}
	}
	d.Set("instances", instances)
	d.Set("segment_options", segmentOptions)
}

func generateGcpInterconnectRequest(d *schema.ResourceData, m interface{}) (*alkira.ConnectorGcpInterconnect, error) {
	instances, err := expandGcpInterconnectInstances(d.Get("instances").([]interface{}), d.Get("segment_options").([]interface{}), m)
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
