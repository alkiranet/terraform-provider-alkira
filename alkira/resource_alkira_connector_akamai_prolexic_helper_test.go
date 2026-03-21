package alkira

import (
	"testing"

	"github.com/alkiranet/alkira-client-go/alkira"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/stretchr/testify/assert"
)

func akamaiProlexicTestSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"tunnel_configuration": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"alkira_public_ip": {
							Type:     schema.TypeString,
							Required: true,
						},
						"tunnel_ips": {
							Type:     schema.TypeSet,
							Required: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"ran_tunnel_ip": {
										Type:     schema.TypeString,
										Required: true,
									},
									"alkira_overlay_tunnel_ip": {
										Type:     schema.TypeString,
										Required: true,
									},
									"akamai_overlay_tunnel_ip": {
										Type:     schema.TypeString,
										Required: true,
									},
								},
							},
						},
					},
				},
			},
		},
	}
}

func TestSetConnectorAkamaiTunnelConfiguration_Empty(t *testing.T) {
	d := akamaiProlexicTestSchema().TestResourceData()
	setConnectorAkamaiTunnelConfiguration(d, []alkira.ConnectorAkamaiProlexicOverlayConfiguration{})
	assert.Equal(t, 0, d.Get("tunnel_configuration").(*schema.Set).Len())
}

func TestSetConnectorAkamaiTunnelConfiguration_Nil(t *testing.T) {
	d := akamaiProlexicTestSchema().TestResourceData()
	setConnectorAkamaiTunnelConfiguration(d, nil)
	assert.Equal(t, 0, d.Get("tunnel_configuration").(*schema.Set).Len())
}

func TestSetConnectorAkamaiTunnelConfiguration_SingleConfig(t *testing.T) {
	d := akamaiProlexicTestSchema().TestResourceData()

	overlayConfig := []alkira.ConnectorAkamaiProlexicOverlayConfiguration{
		{
			AlkiraPublicIp: "203.0.113.1",
			TunnelIps: []alkira.ConnectorAkamaiProlexicTunnelIp{
				{
					RanTunnelDestinationIp: "10.0.0.1",
					AlkiraOverlayTunnelIp:  "10.1.0.1",
					AkamaiOverlayTunnelIp:  "10.2.0.1",
				},
			},
		},
	}

	setConnectorAkamaiTunnelConfiguration(d, overlayConfig)

	configs := d.Get("tunnel_configuration").(*schema.Set).List()
	assert.Len(t, configs, 1)

	cfg := configs[0].(map[string]interface{})
	assert.Equal(t, "203.0.113.1", cfg["alkira_public_ip"])

	ips := cfg["tunnel_ips"].(*schema.Set).List()
	assert.Len(t, ips, 1)
	ip := ips[0].(map[string]interface{})
	assert.Equal(t, "10.0.0.1", ip["ran_tunnel_ip"])
	assert.Equal(t, "10.1.0.1", ip["alkira_overlay_tunnel_ip"])
	assert.Equal(t, "10.2.0.1", ip["akamai_overlay_tunnel_ip"])
}

func TestSetConnectorAkamaiTunnelConfiguration_MultipleTunnelIps(t *testing.T) {
	d := akamaiProlexicTestSchema().TestResourceData()

	overlayConfig := []alkira.ConnectorAkamaiProlexicOverlayConfiguration{
		{
			AlkiraPublicIp: "203.0.113.1",
			TunnelIps: []alkira.ConnectorAkamaiProlexicTunnelIp{
				{
					RanTunnelDestinationIp: "10.0.0.1",
					AlkiraOverlayTunnelIp:  "10.1.0.1",
					AkamaiOverlayTunnelIp:  "10.2.0.1",
				},
				{
					RanTunnelDestinationIp: "10.0.0.2",
					AlkiraOverlayTunnelIp:  "10.1.0.2",
					AkamaiOverlayTunnelIp:  "10.2.0.2",
				},
				{
					RanTunnelDestinationIp: "10.0.0.3",
					AlkiraOverlayTunnelIp:  "10.1.0.3",
					AkamaiOverlayTunnelIp:  "10.2.0.3",
				},
			},
		},
	}

	setConnectorAkamaiTunnelConfiguration(d, overlayConfig)

	configs := d.Get("tunnel_configuration").(*schema.Set).List()
	assert.Len(t, configs, 1)

	cfg := configs[0].(map[string]interface{})
	assert.Equal(t, "203.0.113.1", cfg["alkira_public_ip"])

	ips := cfg["tunnel_ips"].(*schema.Set).List()
	assert.Len(t, ips, 3)

	// Collect all ran_tunnel_ips to verify all 3 are present (Set order not guaranteed)
	ranIps := make([]string, 0, 3)
	for _, entry := range ips {
		ip := entry.(map[string]interface{})
		ranIps = append(ranIps, ip["ran_tunnel_ip"].(string))
	}
	assert.ElementsMatch(t, []string{"10.0.0.1", "10.0.0.2", "10.0.0.3"}, ranIps)
}

func TestSetConnectorAkamaiTunnelConfiguration_MultipleConfigs(t *testing.T) {
	d := akamaiProlexicTestSchema().TestResourceData()

	overlayConfig := []alkira.ConnectorAkamaiProlexicOverlayConfiguration{
		{
			AlkiraPublicIp: "203.0.113.1",
			TunnelIps: []alkira.ConnectorAkamaiProlexicTunnelIp{
				{RanTunnelDestinationIp: "10.0.0.1", AlkiraOverlayTunnelIp: "10.1.0.1", AkamaiOverlayTunnelIp: "10.2.0.1"},
			},
		},
		{
			AlkiraPublicIp: "203.0.113.2",
			TunnelIps: []alkira.ConnectorAkamaiProlexicTunnelIp{
				{RanTunnelDestinationIp: "10.0.1.1", AlkiraOverlayTunnelIp: "10.1.1.1", AkamaiOverlayTunnelIp: "10.2.1.1"},
			},
		},
	}

	setConnectorAkamaiTunnelConfiguration(d, overlayConfig)

	configs := d.Get("tunnel_configuration").(*schema.Set).List()
	assert.Len(t, configs, 2)

	publicIps := make([]string, 0, 2)
	for _, c := range configs {
		cfg := c.(map[string]interface{})
		publicIps = append(publicIps, cfg["alkira_public_ip"].(string))
	}
	assert.ElementsMatch(t, []string{"203.0.113.1", "203.0.113.2"}, publicIps)
}
