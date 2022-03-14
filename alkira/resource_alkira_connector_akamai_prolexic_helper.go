package alkira

import (
	"log"

	"github.com/alkiranet/alkira-client-go/alkira"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func expandConnectorAkamaiTunnelIps(in *schema.Set) []alkira.ConnectorAkamaiProlexicTunnelIp {
	if in == nil || in.Len() == 0 {
		log.Printf("[DEBUG] invalid input of connector_akamai_prolexic tunnel ips")
		return nil
	}

	tunnelIps := make([]alkira.ConnectorAkamaiProlexicTunnelIp, in.Len())

	for i, ip := range in.List() {
		r := alkira.ConnectorAkamaiProlexicTunnelIp{}
		content := ip.(map[string]interface{})

		if v, ok := content["ran_tunnel_ip"].(string); ok {
			r.RanTunnelDestinationIp = v
		}
		if v, ok := content["alkira_overlay_tunnel_ip"].(string); ok {
			r.AlkiraOverlayTunnelIp = v
		}
		if v, ok := content["akamai_overlay_tunnel_ip"].(string); ok {
			r.AkamaiOverlayTunnelIp = v
		}

		tunnelIps[i] = r
	}

	return tunnelIps
}

// expandConnectorAkamaiTunnelConfiguration
func expandConnectorAkamaiTunnelConfiguration(in *schema.Set) []alkira.ConnectorAkamaiProlexicOverlayConfiguration {
	if in == nil || in.Len() == 0 {
		log.Printf("[DEBUG] invalid input of connector_akamai_prolexic tunnel configuration")
		return nil
	}

	configurations := make([]alkira.ConnectorAkamaiProlexicOverlayConfiguration, in.Len())

	for i, config := range in.List() {
		r := alkira.ConnectorAkamaiProlexicOverlayConfiguration{}
		cfg := config.(map[string]interface{})
		if v, ok := cfg["alkira_public_ip"].(string); ok {
			r.AlkiraPublicIp = v
		}
		if v, ok := cfg["tunnel_ips"].(*schema.Set); ok {
			r.TunnelIps = expandConnectorAkamaiTunnelIps(v)
		}
		configurations[i] = r
	}

	return configurations
}

// expandConnectorAkamaiByoipOptions
func expandConnectorAkamaiByoipOptions(in *schema.Set) []alkira.ConnectorAkamaiProlexicByoipOption {

	options := make([]alkira.ConnectorAkamaiProlexicByoipOption, in.Len())

	for i, option := range in.List() {
		r := alkira.ConnectorAkamaiProlexicByoipOption{}
		opt := option.(map[string]interface{})
		if v, ok := opt["byoip_prefix_id"].(int); ok {
			r.ByoipId = v
		}
		if v, ok := opt["enable_route_advertisement"].(bool); ok {
			r.RouteAdvertisementEnabled = v
		}
		options[i] = r
	}

	return options
}
