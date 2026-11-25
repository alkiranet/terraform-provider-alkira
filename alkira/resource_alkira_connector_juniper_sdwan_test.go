package alkira

import (
	"testing"

	"github.com/alkiranet/alkira-client-go/alkira"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestJuniperSdwanScaleGroupId(t *testing.T) {
	r := resourceAlkiraConnectorJuniperSdwan()
	client := &alkira.AlkiraClient{}

	t.Run("set", func(t *testing.T) {
		d := schema.TestResourceDataRaw(t, r.Schema, map[string]interface{}{"scale_group_id": "sg-123"})
		connector, err := generateConnectorJuniperSdwanRequest(d, client)
		require.NoError(t, err)
		assert.Equal(t, "sg-123", connector.ScaleGroupId)
	})

	t.Run("not set", func(t *testing.T) {
		d := schema.TestResourceDataRaw(t, r.Schema, map[string]interface{}{})
		connector, err := generateConnectorJuniperSdwanRequest(d, client)
		require.NoError(t, err)
		assert.Empty(t, connector.ScaleGroupId)
	})
}
