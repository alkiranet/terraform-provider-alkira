package alkira

import (
	"fmt"
	"log"

	"github.com/alkiranet/alkira-client-go/alkira"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// expandConnectorAdvIPSecAdvancedOptions
func expandConnectorAdvIPSecAdvancedOptions(in []interface{}) (*alkira.ConnectorAdvIPSecAdvanced, error) {

	if in == nil || len(in) == 0 {
		log.Printf("[DEBUG] empty IPSec endpoint advanced")
		return nil, nil
	}

	if in == nil || len(in) > 1 {
		log.Printf("[DEBUG] invalid IPSec endpoint advanced")
		return nil, nil
	}

	advanced := &alkira.ConnectorAdvIPSecAdvanced{}

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

// expandConnectorAdvIPSecTunnel expand IPSec gateway tunnels
func expandConnectorAdvIPSecTunnel(in []interface{}) []*alkira.ConnectorAdvIPSecTunnel {
	if in == nil || len(in) == 0 {
		log.Printf("[DEBUG] empty IPSec gateway tunnel")
		return nil
	}

	tunnels := make([]*alkira.ConnectorAdvIPSecTunnel, len(in))

	for i, t := range in {
		config := t.(map[string]interface{})
		r := alkira.ConnectorAdvIPSecTunnel{}

		r.CustomerEnd.OverlayIpReservationId = config["customer_end_overlay_ip_reservation_id"].(string)
		r.CxpEnd.OverlayIpReservationId = config["cxp_end_overlay_ip_reservation_id"].(string)
		r.CxpEnd.PublicIpReservationId = config["cxp_end_public_ip_reservation_id"].(string)
		r.Id = config["id"].(int)
		r.PresharedKey = config["preshared_key"].(string)
		r.ProfileId = config["profile_id"].(int)
		r.TunnelNo = config["number"].(int)

		if v, ok := config["advanced_options"].([]interface{}); ok {

			var err error
			r.Advanced, err = expandConnectorAdvIPSecAdvancedOptions(v)

			if err != nil {
				log.Printf("[ERROR] failed to parse advanced options.")
				break
			}
		}

		tunnels[i] = &r
	}
	return tunnels
}

// expandConnectorAdvIPSecGateway expand IPSEC gateway
func expandConnectorAdvIPSecGateway(in []interface{}) []*alkira.ConnectorAdvIPSecGateway {
	if in == nil || len(in) == 0 {
		log.Printf("[DEBUG] empty IPSec gateway input")
		return nil
	}

	gws := make([]*alkira.ConnectorAdvIPSecGateway, len(in))

	for i, gw := range in {
		gwConfig := gw.(map[string]interface{})
		r := alkira.ConnectorAdvIPSecGateway{}

		r.Name = gwConfig["name"].(string)
		r.CustomerGwIp = gwConfig["customer_gateway_ip"].(string)
		r.HaMode = gwConfig["ha_mode"].(string)
		r.Id = gwConfig["id"].(int)

		if v, ok := gwConfig["tunnel"].([]interface{}); ok {

			var err error
			r.Tunnels = expandConnectorAdvIPSecTunnel(v)

			if err != nil {
				log.Printf("[ERROR] failed to expand tunnels the of gateway.")
				break
			}
		}

		gws[i] = &r
	}
	return gws
}

// expandConnectorAdvIPSecPolicyOptions expand policy_options
func expandConnectorAdvIPSecPolicyOptions(in *schema.Set) (*alkira.ConnectorAdvIPSecPolicyOptions, error) {
	if in == nil || in.Len() == 0 {
		log.Printf("[DEBUG] Empty policy options of IPSec connector.")
		return &alkira.ConnectorAdvIPSecPolicyOptions{}, nil
	}

	if in.Len() > 1 {
		return nil, fmt.Errorf("ERROR: only one policy_options could be specified.")
	}

	policyOptions := alkira.ConnectorAdvIPSecPolicyOptions{}

	for _, input := range in.List() {
		policyOptionsInput := input.(map[string]interface{})

		policyOptions.BranchTSPrefixListIds = convertTypeListToIntList(policyOptionsInput["branch_ts_prefix_list_ids"].([]interface{}))
		policyOptions.CxpTSPrefixListIds = convertTypeListToIntList(policyOptionsInput["cxp_ts_prefix_list_ids"].([]interface{}))
	}

	return &policyOptions, nil
}

// expandConnectorAdvIPSecRoutingOptions expand routing_options
func expandConnectorAdvIPSecRoutingOptions(in *schema.Set) (*alkira.ConnectorAdvIPSecRoutingOptions, error) {
	if in == nil || in.Len() == 0 {
		log.Printf("[DEBUG] Empty routing options of IPSec connector.")
		return &alkira.ConnectorAdvIPSecRoutingOptions{}, nil
	}

	if in.Len() > 1 {
		return nil, fmt.Errorf("ERROR: only one routing_options could be specified.")
	}

	staticOption := alkira.ConnectorAdvIPSecStaticRouting{}
	dynamicOption := alkira.ConnectorAdvIPSecDynamicRouting{}
	routingOptions := alkira.ConnectorAdvIPSecRoutingOptions{}

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

				routingOptions = alkira.ConnectorAdvIPSecRoutingOptions{
					StaticRouting: &staticOption,
				}
			}
		case "DYNAMIC":
			{
				avail, availOk := routingOptionsInput["availability"].(string)

				if availOk {
					dynamicOption.Availability = avail
				}

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

				routingOptions = alkira.ConnectorAdvIPSecRoutingOptions{
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
					dynamicOption.Availability = avail
				}

				asn, asnOk := routingOptionsInput["customer_gateway_asn"].(string)

				if asnOk {
					dynamicOption.CustomerGwAsn = asn
				} else {
					return nil, fmt.Errorf("ERROR: if BOTH routing type is specified, customer_gateway_asn is required.")
				}

				routingOptions = alkira.ConnectorAdvIPSecRoutingOptions{
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

// deflateConnectorAdvIPSecPolicyOptions
func deflateConnectorAdvIPSecPolicyOptions(cfg *alkira.ConnectorAdvIPSecPolicyOptions) map[string]interface{} {
	option := map[string]interface{}{
		"branch_ts_prefix_list_ids": cfg.BranchTSPrefixListIds,
		"cxp_ts_prefix_list_ids":    cfg.CxpTSPrefixListIds,
	}

	return option
}

// deflateConnectorAdvIPSecTunnel
func deflateConnectorAdvIPSecTunnel(tunnelConfig *alkira.ConnectorAdvIPSecTunnel) map[string]interface{} {
	if tunnelConfig == nil {
		log.Printf("[ERROR] invalid IPSec tunnel")
		return nil
	}

	var advanced []map[string]interface{}

	if tunnelConfig.Advanced != nil {
		config := map[string]interface{}{
			"dpd_delay":                 tunnelConfig.Advanced.DPDDelay,
			"dpd_timeout":               tunnelConfig.Advanced.DPDTimeout,
			"esp_dh_group_numbers":      tunnelConfig.Advanced.EspDHGroupNumbers,
			"esp_encryption_algorithms": tunnelConfig.Advanced.EspEncryptionAlgorithms,
			"esp_integrity_algorithms":  tunnelConfig.Advanced.EspIntegrityAlgorithms,
			"esp_life_time":             tunnelConfig.Advanced.EspLifeTime,
			"esp_random_time":           tunnelConfig.Advanced.EspRandomTime,
			"esp_rekey_time":            tunnelConfig.Advanced.EspRekeyTime,
			"ike_dh_group_numbers":      tunnelConfig.Advanced.IkeDHGroupNumbers,
			"ike_encryption_algorithms": tunnelConfig.Advanced.IkeEncryptionAlgorithms,
			"ike_integrity_algorithms":  tunnelConfig.Advanced.IkeIntegrityAlgorithms,
			"ike_over_time":             tunnelConfig.Advanced.IkeOverTime,
			"ike_random_time":           tunnelConfig.Advanced.IkeRandomTime,
			"ike_rekey_time":            tunnelConfig.Advanced.IkeRekeyTime,
			"ike_version":               tunnelConfig.Advanced.IkeVersion,
			"initiator":                 tunnelConfig.Advanced.Initiator,
			"local_auth_type":           tunnelConfig.Advanced.LocalAuthType,
			"local_auth_value":          tunnelConfig.Advanced.LocalAuthValue,
			"remote_auth_type":          tunnelConfig.Advanced.RemoteAuthType,
			"remote_auth_value":         tunnelConfig.Advanced.RemoteAuthValue,
			"replay_window_size":        tunnelConfig.Advanced.ReplayWindowSize,
		}
		advanced = append(advanced, config)
	}

	tunnel := map[string]interface{}{
		"number":                                 tunnelConfig.TunnelNo,
		"preshared_key":                          tunnelConfig.PresharedKey,
		"profile_id":                             tunnelConfig.ProfileId,
		"id":                                     tunnelConfig.Id,
		"customer_end_overlay_ip_reservation_id": tunnelConfig.CustomerEnd.OverlayIpReservationId,
		"cxp_end_overlay_ip_reservation_id":      tunnelConfig.CxpEnd.OverlayIpReservationId,
		"cxp_end_public_ip_reservation_id":       tunnelConfig.CxpEnd.PublicIpReservationId,
		"advanced_options":                       advanced,
	}

	return tunnel
}

// deflateConnectorAdvIPSecGatewayInstance
func deflateConnectorAdvIPSecGatewayInstance(gatewayConfig *alkira.ConnectorAdvIPSecGateway) map[string]interface{} {
	if gatewayConfig == nil {
		log.Printf("[ERROR] invalid IPSec gateway")
		return nil
	}

	var tunnels []map[string]interface{}

	for _, t := range gatewayConfig.Tunnels {
		config := deflateConnectorAdvIPSecTunnel(t)
		tunnels = append(tunnels, config)
	}

	gateway := map[string]interface{}{
		"customer_gateway_ip": gatewayConfig.CustomerGwIp,
		"ha_mode":             gatewayConfig.HaMode,
		"id":                  gatewayConfig.Id,
		"name":                gatewayConfig.Name,
		"tunnels":             tunnels,
	}

	return gateway
}

// deflateConnectorAdvIPSecGateway
func deflateConnectorAdvIPSecGateway(connector *alkira.ConnectorAdvIPSec, d *schema.ResourceData) []map[string]interface{} {
	//
	// Go through all gateways from the config firstly to find a
	// match, either gateway's ID or gateway's name should be
	// uniquely identifying an gateway.
	//
	// On the first read call at the end of the create call, Terraform
	// didn't track any gateway IDs yet.
	//
	var gateways []map[string]interface{}

	for _, gateway := range d.Get("gateway").([]interface{}) {
		gatewayConfig := gateway.(map[string]interface{})

		for _, gw := range connector.Gateways {
			if gatewayConfig["id"].(int) == gw.Id || gatewayConfig["name"].(string) == gw.Name {
				gateway := deflateConnectorAdvIPSecGatewayInstance(gw)
				gateways = append(gateways, gateway)
				break
			}
		}
	}

	//
	// Go through all gateways from the API response one more time to
	// find any gateway that has not been tracked from Terraform
	// config.
	//
	for _, gw := range connector.Gateways {
		new := true

		// Check if the gateway already exists in the Terraform config
		for _, gateway := range d.Get("gateway").([]interface{}) {
			gatewayConfig := gateway.(map[string]interface{})

			if gatewayConfig["id"].(int) == gw.Id || gatewayConfig["name"].(string) == gw.Name {
				new = false
				break
			}
		}

		// If the gateway is new, add it to the tail of the list,
		// this will generate a diff
		if new {
			gateway := deflateConnectorAdvIPSecGatewayInstance(gw)
			gateways = append(gateways, gateway)
			break
		}
	}

	return gateways
}
