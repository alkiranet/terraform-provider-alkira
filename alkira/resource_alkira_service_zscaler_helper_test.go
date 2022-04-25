package alkira

import (
	"testing"

	"github.com/alkiranet/alkira-client-go/alkira"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/stretchr/testify/require"
)

func TestExpandIpsecConfiguration(t *testing.T) {
	expectedZscalerIpsecConfig := defaultZscalerIpsecConfig("prefix")

	r := resourceAlkiraServiceZscaler()
	z := schema.HashResource(r)
	s := schema.NewSet(z, []interface{}{makeMapIpsecConfiguration(expectedZscalerIpsecConfig)})

	actual, err := expandZscalerIpsecConfigurations(s)
	require.NoError(t, err)
	require.Equal(t, expectedZscalerIpsecConfig, actual)
}

func makeMapIpsecConfiguration(z *alkira.ZscalerIpSecConfig) map[string]interface{} {
	m := make(map[string]interface{})
	m["esp_dh_group_number"] = z.EspDhGroupNumber
	m["esp_encryption_algorithm"] = z.EspEncryptionAlgorithm
	m["esp_integrity_algorithm"] = z.EspIntegrityAlgorithm
	m["health_check_type"] = z.HealthCheckType
	m["http_probe_url"] = z.HttpProbeUrl
	m["ike_dh_group_number"] = z.IkeDhGroupNumber
	m["ike_encryption_algorithm"] = z.IkeEncryptionAlgorithm
	m["ike_integrity_algorithm"] = z.IkeIntegrityAlgorithm
	m["local_fpdn_id"] = z.LocalFqdnId
	m["pre_shared_key"] = z.PreSharedKey
	m["ping_probe_ip"] = z.PingProbeIp

	return m
}

func defaultZscalerIpsecConfig(prefix string) *alkira.ZscalerIpSecConfig {
	return &alkira.ZscalerIpSecConfig{
		EspDhGroupNumber:       prefix + "EspDhGroupNumber",
		EspEncryptionAlgorithm: prefix + "EspEncryptionAlgorithm",
		EspIntegrityAlgorithm:  prefix + "EspIntegrityAlgorithm",
		HealthCheckType:        prefix + "HealthCheckType",
		HttpProbeUrl:           prefix + "HttpProbeUrl",
		IkeDhGroupNumber:       prefix + "IkeDhGroupNumber",
		IkeEncryptionAlgorithm: prefix + "IkeEncryptionAlgorithm",
		IkeIntegrityAlgorithm:  prefix + "IkeIntegrityAlgorithm",
		LocalFqdnId:            prefix + "LocalFqdnId",
		PreSharedKey:           prefix + "PreSharedKey",
		PingProbeIp:            prefix + "PingProbeIp",
	}
}
