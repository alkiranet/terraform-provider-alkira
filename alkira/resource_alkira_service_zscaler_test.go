package alkira

import (
	"testing"

	"github.com/alkiranet/alkira-client-go/alkira"
	"github.com/stretchr/testify/require"
)

func TestGenerateRequesetServiceZscaler(t *testing.T) {
	expectedSegmentNames := []string{"default", "default1"}
	expectedTunnelType := "IPSEC"

	z := &alkira.Zscaler{
		Segments:   expectedSegmentNames,
		TunnelType: expectedTunnelType,
	}

	r := resourceAlkiraServiceZscaler()
	d := r.TestResourceData()
	d.Set("segment_names", z.Segments)
	d.Set("tunnel_protocol", z.TunnelType)

	actual, err := generateZscalerRequest(d, z)
	require.NoError(t, err)
	require.Equal(t, expectedSegmentNames, actual.Segments)
}
