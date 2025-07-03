package alkira

import (
	"testing"

	"github.com/alkiranet/alkira-client-go/alkira"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDeflateZscalerIpsecConfiguration(t *testing.T) {
	tests := []struct {
		name        string
		input       *alkira.ZscalerIpSecConfig
		expected    []map[string]interface{}
		expectError bool
	}{
		{
			name:        "nil input",
			input:       nil,
			expected:    nil,
			expectError: true,
		},
		{
			name: "valid ipsec configuration",
			input: &alkira.ZscalerIpSecConfig{
				EspDhGroupNumber:       "14",
				EspEncryptionAlgorithm: "AES256",
				EspIntegrityAlgorithm:  "SHA256",
				HealthCheckType:        "PING",
				HttpProbeUrl:           "http://example.com/health",
				IkeDhGroupNumber:       "14",
				IkeEncryptionAlgorithm: "AES256",
				IkeIntegrityAlgorithm:  "SHA256",
				LocalFqdnId:            "local.example.com",
				PreSharedKey:           "secret123",
				PingProbeIp:            "8.8.8.8",
			},
			expected: []map[string]interface{}{
				{
					"esp_dh_group_number":      "14",
					"esp_encryption_algorithm": "AES256",
					"esp_integrity_algorithm":  "SHA256",
					"health_check_type":        "PING",
					"http_probe_url":           "http://example.com/health",
					"ike_dh_group_number":      "14",
					"ike_encryption_algorithm": "AES256",
					"ike_integrity_algorithm":  "SHA256",
					"local_fpdn_id":            "local.example.com",
					"pre_shared_key":           "secret123",
					"ping_probe_ip":            "8.8.8.8",
				},
			},
			expectError: false,
		},
		{
			name: "minimal ipsec configuration",
			input: &alkira.ZscalerIpSecConfig{
				EspDhGroupNumber:       "5",
				EspEncryptionAlgorithm: "AES128",
				EspIntegrityAlgorithm:  "SHA1",
				IkeDhGroupNumber:       "5",
				IkeEncryptionAlgorithm: "AES128",
				IkeIntegrityAlgorithm:  "SHA1",
				PreSharedKey:           "secret",
			},
			expected: []map[string]interface{}{
				{
					"esp_dh_group_number":      "5",
					"esp_encryption_algorithm": "AES128",
					"esp_integrity_algorithm":  "SHA1",
					"health_check_type":        "",
					"http_probe_url":           "",
					"ike_dh_group_number":      "5",
					"ike_encryption_algorithm": "AES128",
					"ike_integrity_algorithm":  "SHA1",
					"local_fpdn_id":            "",
					"pre_shared_key":           "secret",
					"ping_probe_ip":            "",
				},
			},
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := deflateZscalerIpsecConfiguration(tt.input)

			if tt.expectError {
				assert.Error(t, err)
				assert.Nil(t, result)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expected, result)
			}
		})
	}
}

func TestZscalerConfigurationValidation(t *testing.T) {
	t.Run("valid encryption algorithms", func(t *testing.T) {
		validAlgorithms := []string{"AES128", "AES256", "3DES"}
		for _, alg := range validAlgorithms {
			config := &alkira.ZscalerIpSecConfig{
				EspEncryptionAlgorithm: alg,
				IkeEncryptionAlgorithm: alg,
			}

			result, err := deflateZscalerIpsecConfiguration(config)
			require.NoError(t, err)
			assert.Equal(t, alg, result[0]["esp_encryption_algorithm"])
			assert.Equal(t, alg, result[0]["ike_encryption_algorithm"])
		}
	})

	t.Run("valid integrity algorithms", func(t *testing.T) {
		validAlgorithms := []string{"SHA1", "SHA256", "SHA384", "SHA512"}
		for _, alg := range validAlgorithms {
			config := &alkira.ZscalerIpSecConfig{
				EspIntegrityAlgorithm: alg,
				IkeIntegrityAlgorithm: alg,
			}

			result, err := deflateZscalerIpsecConfiguration(config)
			require.NoError(t, err)
			assert.Equal(t, alg, result[0]["esp_integrity_algorithm"])
			assert.Equal(t, alg, result[0]["ike_integrity_algorithm"])
		}
	})

	t.Run("valid DH group numbers", func(t *testing.T) {
		validGroups := []string{"1", "2", "5", "14", "15", "16", "17", "18", "19", "20", "21"}
		for _, group := range validGroups {
			config := &alkira.ZscalerIpSecConfig{
				EspDhGroupNumber: group,
				IkeDhGroupNumber: group,
			}

			result, err := deflateZscalerIpsecConfiguration(config)
			require.NoError(t, err)
			assert.Equal(t, group, result[0]["esp_dh_group_number"])
			assert.Equal(t, group, result[0]["ike_dh_group_number"])
		}
	})

	t.Run("valid health check types", func(t *testing.T) {
		validTypes := []string{"PING", "HTTP", "HTTPS"}
		for _, hcType := range validTypes {
			config := &alkira.ZscalerIpSecConfig{
				HealthCheckType: hcType,
			}

			result, err := deflateZscalerIpsecConfiguration(config)
			require.NoError(t, err)
			assert.Equal(t, hcType, result[0]["health_check_type"])
		}
	})
}
