package alkira

import (
	"encoding/json"
	"net/http"
	"testing"

	"github.com/alkiranet/alkira-client-go/alkira"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/stretchr/testify/require"
)

func TestGenerateRequesetServiceZscaler(t *testing.T) {
	expectedIpSecConfig := alkira.ZscalerIpSecConfig{
		EspDhGroupNumber:       "espDhGroupNumber",
		EspEncryptionAlgorithm: "espEncryptionAlgorithm",
		EspIntegrityAlgorithm:  "espIntegrityAlgorithm",
		HealthCheckType:        "healthCheckType",
		HttpProbeUrl:           "httpProbeUrl",
		IkeDhGroupNumber:       "ikeDhGroupNumber",
		IkeEncryptionAlgorithm: "ikeEncryptionAlgorithm",
		IkeIntegrityAlgorithm:  "ikeIntegrityAlgorithm",
		LocalFqdnId:            "localFqdnId",
		PreSharedKey:           "preSharedKey",
		PingProbeIp:            "pingProbeIp",
	}

	expectedTunnelType := "IPSEC"

	z := &alkira.ServiceZscaler{
		TunnelType:         expectedTunnelType,
		IpsecConfiguration: &expectedIpSecConfig,
	}

	deflatedIpSecCfg, err := deflateZscalerIpsecConfiguration(&expectedIpSecConfig)
	require.NoError(t, err)
	r := resourceAlkiraServiceZscaler()
	d := r.TestResourceData()

	d.Set("ipsec_configuration", schema.NewSet(schema.HashResource(r), []interface{}{deflatedIpSecCfg[0]})) //z.IpsecConfiguration)
	d.Set("tunnel_protocol", z.TunnelType)

	client := serveZscaler(t, z)
	actual, err := generateZscalerRequest(d, client)
	require.NoError(t, err)
	require.Equal(t, expectedIpSecConfig, *actual.IpsecConfiguration)
}

func serveZscaler(t *testing.T, z *alkira.ServiceZscaler) *alkira.AlkiraClient {
	return createMockAlkiraClient(t, func(w http.ResponseWriter, req *http.Request) {
		json.NewEncoder(w).Encode(z)
		w.Header().Set("Content-Type", "application/json")
	})
}
