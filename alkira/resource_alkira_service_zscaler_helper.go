package alkira

import (
	"errors"

	"github.com/alkiranet/alkira-client-go/alkira"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func expandZscalerIpsecConfigurations(in *schema.Set) (*alkira.ZscalerIpSecConfig, error) {
	if in == nil || in.Len() <= 0 {
		return nil, errors.New("ZscalerIpSecConfig must be length 1")
	}

	if in.Len() > 1 {
		return nil, errors.New("expandIpsecConfigurations Set input can have only one entry")
	}

	ip := &alkira.ZscalerIpSecConfig{}
	for _, c := range in.List() {
		cfg := c.(map[string]interface{})
		if v, ok := cfg["esp_dh_group_number"].(string); ok {
			ip.EspDhGroupNumber = v
		}
		if v, ok := cfg["esp_encryption_algorithm"].(string); ok {
			ip.EspEncryptionAlgorithm = v
		}
		if v, ok := cfg["esp_integrity_algorithm"].(string); ok {
			ip.EspIntegrityAlgorithm = v
		}
		if v, ok := cfg["health_check_type"].(string); ok {
			ip.HealthCheckType = v
		}
		if v, ok := cfg["http_probe_url"].(string); ok {
			ip.HttpProbeUrl = v
		}
		if v, ok := cfg["ike_dh_group_number"].(string); ok {
			ip.IkeDhGroupNumber = v
		}
		if v, ok := cfg["ike_encryption_algorithm"].(string); ok {
			ip.IkeEncryptionAlgorithm = v
		}
		if v, ok := cfg["ike_integrity_algorithm"].(string); ok {
			ip.IkeIntegrityAlgorithm = v
		}
		if v, ok := cfg["local_fpdn_id"].(string); ok {
			ip.LocalFqdnId = v
		}
		if v, ok := cfg["pre_shared_key"].(string); ok {
			ip.PreSharedKey = v
		}
		if v, ok := cfg["ping_probe_ip"].(string); ok {
			ip.PingProbeIp = v
		}
	}

	return ip, nil
}

func deflateZscalerIpsecConfiguration(z *alkira.ZscalerIpSecConfig) []map[string]interface{} {
	cfg := make(map[string]interface{})
	cfg["esp_dh_group_number"] = z.EspDhGroupNumber
	cfg["esp_encryption_algorithm"] = z.EspEncryptionAlgorithm
	cfg["esp_integrity_algorithm"] = z.EspIntegrityAlgorithm
	cfg["health_check_type"] = z.HealthCheckType
	cfg["http_probe_url"] = z.HttpProbeUrl
	cfg["ike_dh_group_number"] = z.IkeDhGroupNumber
	cfg["ike_encryption_algorithm"] = z.IkeEncryptionAlgorithm
	cfg["ike_integrity_algorithm"] = z.IkeIntegrityAlgorithm
	cfg["local_fpdn_id"] = z.LocalFqdnId
	cfg["pre_shared_key"] = z.PreSharedKey
	cfg["ping_probe_ip"] = z.PingProbeIp

	return []map[string]interface{}{cfg}
}
