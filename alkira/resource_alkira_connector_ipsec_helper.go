package alkira

import (
	"fmt"
	"log"

	"github.com/alkiranet/alkira-client-go/alkira"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// expandIPSecEndpoint expand IPSEC endpoint section
func expandConnectorIPSecEndpoint(in *schema.Set) []*alkira.ConnectorIPSecSite {
	if in == nil || in.Len() == 0 {
		log.Printf("[DEBUG] empty IPSec endpoint input")
		return nil
	}

	sites := make([]*alkira.ConnectorIPSecSite, in.Len())

	for i, site := range in.List() {

		r := alkira.ConnectorIPSecSite{}
		siteCfg := site.(map[string]interface{})

		if v, ok := siteCfg["name"].(string); ok {
			r.Name = v
		}
		if v, ok := siteCfg["customer_gateway_ip"].(string); ok {
			r.CustomerGwIp = v
		}

		r.BillingTags = convertTypeListToIntList(siteCfg["billing_tag_ids"].([]interface{}))
		r.PresharedKeys = convertTypeListToStringList(siteCfg["preshared_keys"].([]interface{}))

		sites[i] = &r
	}

	return sites
}

// expandConnectorIPSecSegmentOptions expand segment_options
func expandConnectorIPSecSegmentOptions(in *schema.Set) (interface{}, error) {
	if in == nil || in.Len() == 0 {
		log.Printf("[DEBUG] Empty segment options of IPSec connector.")
		return nil, nil
	}

	// Only one segment in IPSec connector is supported
	if in.Len() > 1 {
		return nil, fmt.Errorf("ERROR: only one segment_options could be specified.")
	}

	segmentOptions := make(map[string]alkira.ConnectorIPSecSegmentOptions)

	for _, input := range in.List() {
		segmentOptionsInput := input.(map[string]interface{})

		segmentName, _ := segmentOptionsInput["name"].(string)
		var segmentOption alkira.ConnectorIPSecSegmentOptions

		if v, ok := segmentOptionsInput["disable_internet_exit"].(bool); ok {
			segmentOption.DisableInternetExit = &v
		}

		if v, ok := segmentOptionsInput["disable_advertise_on_prem_routes"].(bool); ok {
			t := !v
			segmentOption.AdvertiseOnPremRoutes = &t
		}

		segmentOptions[segmentName] = segmentOption
	}

	return segmentOptions, nil
}

// expandConnectorIPSecPolicyOptions expand policy_options
func expandConnectorIPSecPolicyOptions(in *schema.Set) (*alkira.ConnectorIPSecPolicyOptions, error) {
	if in == nil || in.Len() == 0 {
		log.Printf("[DEBUG] Empty policy options of IPSec connector.")
		return &alkira.ConnectorIPSecPolicyOptions{}, nil
	}

	if in.Len() > 1 {
		return nil, fmt.Errorf("ERROR: only one policy_options could be specified.")
	}

	policyOptions := alkira.ConnectorIPSecPolicyOptions{}

	for _, input := range in.List() {
		policyOptionsInput := input.(map[string]interface{})

		policyOptions.BranchTSPrefixListIds = convertTypeListToIntList(policyOptionsInput["on_prem_prefix_list_ids"].([]interface{}))
		policyOptions.CxpTSPrefixListIds = convertTypeListToIntList(policyOptionsInput["cxp_prefix_list_ids"].([]interface{}))
	}

	return &policyOptions, nil
}

// expandConnectorIPSecRoutingOptions expand routing_options
func expandConnectorIPSecRoutingOptions(in *schema.Set) (*alkira.ConnectorIPSecRoutingOptions, error) {
	if in == nil || in.Len() == 0 {
		log.Printf("[DEBUG] Empty routing options of IPSec connector.")
		return &alkira.ConnectorIPSecRoutingOptions{}, nil
	}

	if in.Len() > 1 {
		return nil, fmt.Errorf("ERROR: only one routing_options could be specified.")
	}

	staticOption := alkira.ConnectorIPSecStaticRouting{}
	dynamicOption := alkira.ConnectorIPSecDynamicRouting{}
	routingOptions := alkira.ConnectorIPSecRoutingOptions{}

	for _, input := range in.List() {
		routingOptionsInput := input.(map[string]interface{})

		switch routingType := routingOptionsInput["type"].(string); routingType {
		case "STATIC":
			{
				v, ok := routingOptionsInput["prefix_list_id"].(int)

				if ok {
					staticOption.PrefixListId = v
				} else {
					return nil, fmt.Errorf("ERROR: if STATIC routing type is specified, prefix_list_id is required.")
				}

				avail, availOk := routingOptionsInput["availability"].(string)

				if availOk {
					staticOption.Availability = avail
				}

				routingOptions = alkira.ConnectorIPSecRoutingOptions{
					StaticRouting: &staticOption,
				}
			}
		case "DYNAMIC":
			{
				v, ok := routingOptionsInput["customer_gateway_asn"].(string)

				if ok {
					dynamicOption.CustomerGwAsn = v
				} else {
					return nil, fmt.Errorf("ERROR: if DYNAMIC routing type is specified, customer_gateway_asn is required.")
				}

				bgp, ok := routingOptionsInput["bgp_auth_key"].(string)

				if ok {
					dynamicOption.BgpAuthKeyAlkira = bgp
				}

				routingOptions = alkira.ConnectorIPSecRoutingOptions{
					DynamicRouting: &dynamicOption,
				}
			}
		case "BOTH":
			{
				id, idOk := routingOptionsInput["prefix_list_id"].(int)

				if idOk {
					staticOption.PrefixListId = id
				} else {
					return nil, fmt.Errorf("ERROR: if BOTH routing type is specified, prefix_list_id is required.")
				}

				avail, availOk := routingOptionsInput["availability"].(string)

				if availOk {
					staticOption.Availability = avail
				}

				asn, asnOk := routingOptionsInput["customer_gateway_asn"].(string)

				if asnOk {
					dynamicOption.CustomerGwAsn = asn
				} else {
					return nil, fmt.Errorf("ERROR: if BOTH routing type is specified, customer_gateway_asn is required.")
				}

				routingOptions = alkira.ConnectorIPSecRoutingOptions{
					StaticRouting:  &staticOption,
					DynamicRouting: &dynamicOption,
				}
			}
		default:
			return nil, fmt.Errorf("ERROR: invalid routing type")
		}
	}

	return &routingOptions, nil
}
