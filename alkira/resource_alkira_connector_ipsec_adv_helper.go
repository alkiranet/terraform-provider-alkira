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
		log.Printf("[DEBUG] empty connector-ipsec-adv advanced options")
		return nil, nil
	}

	if in == nil || len(in) > 1 {
		log.Printf("[DEBUG] invalid connector-ipsec-adv endpoint advanced")
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
		r.Id = config["id"].(string)
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

// expandConnectorAdvIPSecGateway expand "gateway" block
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

// expandConnectorAdvIPSecPolicyOptions expand "policy_options" block
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

		policyOptions.BranchTSPrefixListIds = convertTypeSetToIntList(policyOptionsInput["on_prem_prefix_list_ids"].(*schema.Set))
		policyOptions.CxpTSPrefixListIds = convertTypeSetToIntList(policyOptionsInput["cxp_prefix_list_ids"].(*schema.Set))
	}

	return &policyOptions, nil
}

// expandConnectorAdvIPSecRoutingOptions expand "routing_options" block
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

// generateConnectorIPSecAdvRequest generate request for connector-ipsec
func generateConnectorIPSecAdvRequest(d *schema.ResourceData, m interface{}) (*alkira.ConnectorAdvIPSec, error) {

	gateways := expandConnectorAdvIPSecGateway(d.Get("gateway").([]interface{}))

	//
	// Segment
	//
	segmentName, err := getSegmentNameById(d.Get("segment_id").(string), m)

	if err != nil {
		return nil, err
	}

	//
	// Construct Policy Options and Routing Options
	//
	// Base on the vpn_mode, switch what options to use
	//
	vpnMode := d.Get("vpn_mode").(string)

	var policyOptions *alkira.ConnectorAdvIPSecPolicyOptions
	var routingOptions *alkira.ConnectorAdvIPSecRoutingOptions

	switch vpnMode := d.Get("vpn_mode").(string); vpnMode {
	case "ROUTE_BASED":
		{
			routingOptions, err = expandConnectorAdvIPSecRoutingOptions(
				d.Get("routing_options").(*schema.Set))

			if err != nil {
				return nil, err
			}
		}
	case "POLICY_BASED":
		{
			policyOptions, err = expandConnectorAdvIPSecPolicyOptions(
				d.Get("policy_options").(*schema.Set))

			if err != nil {
				return nil, err
			}
		}
	}

	// Construct the request
	connector := &alkira.ConnectorAdvIPSec{
		AdvertiseDefaultRoute: d.Get("advertise_default_route").(bool),
		AdvertiseOnPremRoutes: d.Get("advertise_on_prem_routes").(bool),
		BillingTags:           convertTypeSetToIntList(d.Get("billing_tag_ids").(*schema.Set)),
		CXP:                   d.Get("cxp").(string),
		Enabled:               d.Get("enabled").(bool),
		DestinationType:       d.Get("destination_type").(string),
		Group:                 d.Get("group").(string),
		Name:                  d.Get("name").(string),
		Segment:               segmentName,
		Size:                  d.Get("size").(string),
		TunnelsPerGateway:     d.Get("tunnels_per_gateway").(int),
		VpnMode:               vpnMode,
		Description:           d.Get("description").(string),
		Gateways:              gateways,
		PolicyOptions:         policyOptions,
		RoutingOptions:        routingOptions,
	}

	return connector, nil
}

// setConnectorAdvIPSecPolicyOptions
func setConnectorAdvIPSecPolicyOptions(policyOptions *alkira.ConnectorAdvIPSecPolicyOptions, d *schema.ResourceData) {

	if policyOptions == nil {
		log.Printf("[DEBUG] empty policy options")
		return
	}

	var options []map[string]interface{}

	option := map[string]interface{}{
		"on_prem_prefix_list_ids": policyOptions.BranchTSPrefixListIds,
		"cxp_prefix_list_ids":     policyOptions.CxpTSPrefixListIds,
	}

	options = append(options, option)
	d.Set("policy_options", options)
}

// setConnectorAdvIPSecRoutingOptions
func setConnectorAdvIPSecRoutingOptions(routingOptions *alkira.ConnectorAdvIPSecRoutingOptions, d *schema.ResourceData) {

	if routingOptions == nil {
		log.Printf("[DEBUG] empty routing options")
		return
	}

	var options []map[string]interface{}

	// If the routing type is "BOTH"
	if routingOptions.StaticRouting != nil && routingOptions.DynamicRouting != nil {
		log.Printf("[DEBUG] routing_type is BOTH")

		option := map[string]interface{}{
			"type":                 "BOTH",
			"prefix_list_id":       routingOptions.StaticRouting.PrefixListId,
			"availability":         routingOptions.StaticRouting.Availability,
			"customer_gateway_asn": routingOptions.DynamicRouting.CustomerGwAsn,
		}
		options = append(options, option)
	} else if routingOptions.DynamicRouting == nil {
		log.Printf("[DEBUG] routing_type is STATIC")

		option := map[string]interface{}{
			"type":           "STATIC",
			"prefix_list_id": routingOptions.StaticRouting.PrefixListId,
			"availability":   routingOptions.StaticRouting.Availability,
		}
		options = append(options, option)
	} else if routingOptions.StaticRouting == nil {
		log.Printf("[DEBUG] routing_type is DYNAMIC")

		option := map[string]interface{}{
			"type":                 "DYNAMIC",
			"availability":         routingOptions.DynamicRouting.Availability,
			"customer_gateway_asn": routingOptions.DynamicRouting.CustomerGwAsn,
			"bgp_auth_key":         routingOptions.DynamicRouting.BgpAuthKeyAlkira,
		}
		options = append(options, option)
	} else {
		log.Printf("[DEBUG] no routing options")
		return
	}

	d.Set("routing_options", options)
}

// deflateConnectorAdvIPSecTunnel
func deflateConnectorAdvIPSecTunnel(tunnelConfig *alkira.ConnectorAdvIPSecTunnel) map[string]interface{} {

	if tunnelConfig == nil {
		log.Printf("[ERROR] invalid IPSec tunnel")
		return nil
	}

	advancedConfig := make(map[string]interface{})

	advancedConfig["dpd_delay"] = tunnelConfig.Advanced.DPDDelay
	advancedConfig["dpd_timeout"] = tunnelConfig.Advanced.DPDTimeout
	advancedConfig["esp_dh_group_numbers"] = tunnelConfig.Advanced.EspDHGroupNumbers
	advancedConfig["esp_encryption_algorithms"] = tunnelConfig.Advanced.EspEncryptionAlgorithms
	advancedConfig["esp_integrity_algorithms"] = tunnelConfig.Advanced.EspIntegrityAlgorithms
	advancedConfig["esp_life_time"] = tunnelConfig.Advanced.EspLifeTime
	advancedConfig["esp_random_time"] = tunnelConfig.Advanced.EspRandomTime
	advancedConfig["esp_rekey_time"] = tunnelConfig.Advanced.EspRekeyTime
	advancedConfig["ike_dh_group_numbers"] = tunnelConfig.Advanced.IkeDHGroupNumbers
	advancedConfig["ike_encryption_algorithms"] = tunnelConfig.Advanced.IkeEncryptionAlgorithms
	advancedConfig["ike_integrity_algorithms"] = tunnelConfig.Advanced.IkeIntegrityAlgorithms
	advancedConfig["ike_over_time"] = tunnelConfig.Advanced.IkeOverTime
	advancedConfig["ike_random_time"] = tunnelConfig.Advanced.IkeRandomTime
	advancedConfig["ike_rekey_time"] = tunnelConfig.Advanced.IkeRekeyTime
	advancedConfig["ike_version"] = tunnelConfig.Advanced.IkeVersion
	advancedConfig["initiator"] = tunnelConfig.Advanced.Initiator
	advancedConfig["local_auth_type"] = tunnelConfig.Advanced.LocalAuthType
	advancedConfig["local_auth_value"] = tunnelConfig.Advanced.LocalAuthValue
	advancedConfig["remote_auth_type"] = tunnelConfig.Advanced.RemoteAuthType
	advancedConfig["remote_auth_value"] = tunnelConfig.Advanced.RemoteAuthValue
	advancedConfig["replay_window_size"] = tunnelConfig.Advanced.ReplayWindowSize

	tunnel = map[string]interface{}{
		"number":                                 tunnelConfig.TunnelNo,
		"preshared_key":                          tunnelConfig.PresharedKey,
		"profile_id":                             tunnelConfig.ProfileId,
		"id":                                     tunnelConfig.Id,
		"customer_end_overlay_ip":                tunnelConfig.CustomerEnd.OverlayIp,
		"customer_end_overlay_ip_reservation_id": tunnelConfig.CustomerEnd.OverlayIpReservationId,
		"cxp_end_overlay_ip_reservation_id":      tunnelConfig.CxpEnd.OverlayIpReservationId,
		"cxp_end_public_ip_reservation_id":       tunnelConfig.CxpEnd.PublicIpReservationId,
		"advanced_options":                       []interface{}{advancedConfig},
	}

	return tunnel
}

// deflateConnectorAdvIPSecGatewayInstance
func deflateConnectorAdvIPSecGatewayInstance(gatewayConfig *alkira.ConnectorAdvIPSecGateway) map[string]interface{} {
	if gatewayConfig == nil {
		log.Printf("[ERROR] invalid IPSec gateway")
		return nil
	}

	tunnels := make([]interface{}, len(gatewayConfig.Tunnels), len(gatewayConfig.Tunnels))

	for i, t := range gatewayConfig.Tunnels {
		config := deflateConnectorAdvIPSecTunnel(t)
		tunnels[i] = config
	}

	gateway := map[string]interface{}{
		"customer_gateway_ip": gatewayConfig.CustomerGwIp,
		"ha_mode":             gatewayConfig.HaMode,
		"id":                  gatewayConfig.Id,
		"name":                gatewayConfig.Name,
		"tunnel":              tunnels,
	}

	return gateway
}

// deflateConnectorAdvIPSecGateway
func deflateConnectorAdvIPSecGateway(connector *alkira.ConnectorAdvIPSec, d *schema.ResourceData) []interface{} {

	gateways := make([]interface{}, len(connector.Gateways), len(connector.Gateways))

	for i, gw := range connector.Gateways {
		gateway := deflateConnectorAdvIPSecGatewayInstance(gw)
		gateways[i] = gateway
	}

	return gateways
}

// setConnectorAdvIPSec
func setConnectorAdvIPSec(connector *alkira.ConnectorAdvIPSec, d *schema.ResourceData, m interface{}) error {

	d.Set("advertise_default_route", connector.AdvertiseDefaultRoute)
	d.Set("advertise_on_prem_routes", connector.AdvertiseOnPremRoutes)
	d.Set("billing_tag_ids", connector.BillingTags)
	d.Set("cxp", connector.CXP)
	d.Set("destination_type", connector.DestinationType)
	d.Set("enabled", connector.Enabled)
	d.Set("group", connector.Group)
	d.Set("implicit_group_id", connector.ImplicitGroupId)
	d.Set("name", connector.Name)
	d.Set("size", connector.Size)
	d.Set("tunnels_per_gateway", connector.TunnelsPerGateway)
	d.Set("vpn_mode", connector.VpnMode)
	d.Set("description", connector.Description)

	// gateway block
	d.Set("gateway", deflateConnectorAdvIPSecGateway(connector, d))

	// policy_options block
	setConnectorAdvIPSecPolicyOptions(connector.PolicyOptions, d)

	// routing_options block
	setConnectorAdvIPSecRoutingOptions(connector.RoutingOptions, d)

	// segment
	segmentId, err := getSegmentIdByName(connector.Segment, m)

	if err != nil {
		return err
	}
	d.Set("segment_id", segmentId)

	return nil
}
