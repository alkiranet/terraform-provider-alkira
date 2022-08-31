package alkira

import (
	"fmt"
	"log"

	"github.com/alkiranet/alkira-client-go/alkira"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// expandConnectorIPSecEndpointAdvanced
func expandConnectorIPSecEndpointAdvanced(in []interface{}) (*alkira.ConnectorIPSecSiteAdvanced, error) {

	if in == nil || len(in) == 0 {
		log.Printf("[DEBUG] empty IPSec endpoint advanced")
		return nil, nil
	}

	if in == nil || len(in) > 1 {
		log.Printf("[DEBUG] invalid IPSec endpoint advanced")
		return nil, nil
	}

	advanced := &alkira.ConnectorIPSecSiteAdvanced{}

	for _, input := range in {
		config := input.(map[string]interface{})

		if v, ok := config["dpd_delay"].(int); ok {
			advanced.DPDDelay = v
		}
		if v, ok := config["dpd_timeout"].(int); ok {
			advanced.DPDTimeout = v
		}
		if v, ok := config["esp_dh_group_numbers"].([]interface{}); ok {
			advanced.EspDHGroupNumbers = convertTypeListToStringList(v)
		}
		if v, ok := config["esp_encryption_algorithms"].([]interface{}); ok {
			advanced.EspEncryptionAlgorithms = convertTypeListToStringList(v)
		}
		if v, ok := config["esp_integrity_algorithms"].([]interface{}); ok {
			advanced.EspIntegrityAlgorithms = convertTypeListToStringList(v)
		}
		if v, ok := config["esp_life_time"].(int); ok {
			advanced.EspLifeTime = v
		}
		if v, ok := config["esp_random_time"].(int); ok {
			advanced.EspRandomTime = v
		}
		if v, ok := config["esp_rekey_time"].(int); ok {
			advanced.EspRekeyTime = v
		}
		if v, ok := config["ike_dh_group_numbers"].([]interface{}); ok {
			advanced.IkeDHGroupNumbers = convertTypeListToStringList(v)
		}
		if v, ok := config["ike_encryption_algorithms"].([]interface{}); ok {
			advanced.IkeEncryptionAlgorithms = convertTypeListToStringList(v)
		}
		if v, ok := config["ike_integrity_algorithms"].([]interface{}); ok {
			advanced.IkeIntegrityAlgorithms = convertTypeListToStringList(v)
		}
		if v, ok := config["ike_over_time"].(int); ok {
			advanced.IkeOverTime = v
		}
		if v, ok := config["ike_random_time"].(int); ok {
			advanced.IkeRandomTime = v
		}
		if v, ok := config["ike_rekey_time"].(int); ok {
			advanced.IkeRekeyTime = v
		}
		if v, ok := config["ike_version"].(string); ok {
			advanced.IkeVersion = v
		}
		if v, ok := config["initiator"].(bool); ok {
			advanced.Initiator = v
		}
		if v, ok := config["local_auth_type"].(string); ok {
			advanced.LocalAuthType = v
		}
		if v, ok := config["local_auth_value"].(string); ok {
			advanced.LocalAuthValue = v
		}
		if v, ok := config["remote_auth_type"].(string); ok {
			advanced.RemoteAuthType = v
		}
		if v, ok := config["remote_auth_value"].(string); ok {
			advanced.RemoteAuthValue = v
		}
		if v, ok := config["replay_window_size"].(int); ok {
			advanced.ReplayWindowSize = v
		}
	}

	return advanced, nil
}

// expandIPSecEndpoint expand IPSEC endpoint section
func expandConnectorIPSecEndpoint(in []interface{}) []*alkira.ConnectorIPSecSite {
	if in == nil || len(in) == 0 {
		log.Printf("[DEBUG] empty IPSec endpoint input")
		return nil
	}

	sites := make([]*alkira.ConnectorIPSecSite, len(in))

	for i, site := range in {
		siteConfig := site.(map[string]interface{})
		r := alkira.ConnectorIPSecSite{}

		r.Name = siteConfig["name"].(string)
		r.CustomerGwIp = siteConfig["customer_gateway_ip"].(string)
		r.HaMode = siteConfig["ha_mode"].(string)
		r.Id = siteConfig["id"].(int)

		if v, ok := siteConfig["billing_tag_ids"].([]interface{}); ok {
			r.BillingTags = convertTypeListToIntList(v)
		}

		if v, ok := siteConfig["preshared_keys"].([]interface{}); ok {
			r.PresharedKeys = convertTypeListToStringList(v)
		}

		if v, ok := siteConfig["advanced_options"].([]interface{}); ok {

			var err error
			r.Advanced, err = expandConnectorIPSecEndpointAdvanced(v)

			if err != nil {
				log.Printf("[ERROR] failed to parse advanced block of endpoint.")
				break
			}
		}

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

		if v, ok := segmentOptionsInput["allow_nat_exit"].(bool); ok {
			t := !v
			segmentOption.DisableInternetExit = &t
		}

		if v, ok := segmentOptionsInput["advertise_on_prem_routes"].(bool); ok {
			segmentOption.AdvertiseOnPremRoutes = &v
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
