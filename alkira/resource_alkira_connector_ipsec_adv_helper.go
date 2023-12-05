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

		if v, ok := config["ike_version"].(string); ok {
			advanced.IkeVersion = v
		}
		if v, ok := config["initiator"].(bool); ok {
			advanced.Initiator = v
		}
		if v, ok := config["remote_auth_type"].(string); ok {
			advanced.RemoteAuthType = v
		}
		if v, ok := config["remote_auth_value"].(string); ok {
			advanced.RemoteAuthValue = v
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

		r.CustomerEnd.OverlayIp = config["customer_end_overlay_ip"].(string)
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

		policyOptions.BranchTSPrefixListIds = convertTypeSetToIntList(policyOptionsInput["on_prem_prefix_list_ids"].(*schema.Set))
		policyOptions.CxpTSPrefixListIds = convertTypeSetToIntList(policyOptionsInput["cxp_prefix_list_ids"].(*schema.Set))
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
		"on_prem_prefix_list_ids": cfg.BranchTSPrefixListIds,
		"cxp_prefix_list_ids":     cfg.CxpTSPrefixListIds,
	}

	return option
}

// deflateConnectorAdvIPSecTunnel
func deflateConnectorAdvIPSecTunnel(tunnelConfig *alkira.ConnectorAdvIPSecTunnel) map[string]interface{} {
	if tunnelConfig == nil {
		log.Printf("[ERROR] invalid IPSec tunnel")
		return nil
	}

	advancedConfig := make(map[string]interface{})

	if tunnelConfig.Advanced != nil {
		advancedConfig["ike_version"] = tunnelConfig.Advanced.IkeVersion
		advancedConfig["initiator"] = tunnelConfig.Advanced.Initiator
		advancedConfig["remote_auth_type"] = tunnelConfig.Advanced.RemoteAuthType
		advancedConfig["remote_auth_value"] = tunnelConfig.Advanced.RemoteAuthValue
	}

	tunnel := map[string]interface{}{
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
