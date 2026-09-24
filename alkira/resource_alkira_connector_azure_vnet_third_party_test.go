package alkira

import (
	"encoding/json"
	"testing"

	"github.com/alkiranet/alkira-client-go/alkira"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestConnectorAzureVnetThirdPartyScaleGroupId(t *testing.T) {
	t.Run("omitted when empty", func(t *testing.T) {
		connector := alkira.AzureVnetThirdPartyConnector{
			Name: "test-connector",
			CXP:  "US-EAST-2",
		}

		data, err := json.Marshal(connector)
		require.NoError(t, err)

		var m map[string]interface{}
		require.NoError(t, json.Unmarshal(data, &m))
		require.NotContains(t, m, "scaleGroupId", "scaleGroupId should be omitted when empty")
	})

	t.Run("included when set", func(t *testing.T) {
		connector := alkira.AzureVnetThirdPartyConnector{
			Name:         "test-connector",
			CXP:          "US-EAST-2",
			ScaleGroupId: "sg-abc123",
		}

		data, err := json.Marshal(connector)
		require.NoError(t, err)

		var m map[string]interface{}
		require.NoError(t, json.Unmarshal(data, &m))
		require.Contains(t, m, "scaleGroupId", "scaleGroupId should be present when set")
		require.Equal(t, "sg-abc123", m["scaleGroupId"])
	})

	t.Run("round-trip marshal/unmarshal", func(t *testing.T) {
		original := alkira.AzureVnetThirdPartyConnector{
			Name:         "test-connector",
			CXP:          "US-WEST-1",
			ScaleGroupId: "sg-xyz789",
		}

		data, err := json.Marshal(original)
		require.NoError(t, err)

		var decoded alkira.AzureVnetThirdPartyConnector
		require.NoError(t, json.Unmarshal(data, &decoded))
		require.Equal(t, original.ScaleGroupId, decoded.ScaleGroupId)
	})
}

func TestAzureVnetThirdPartyScaleGroupIdSchema(t *testing.T) {
	resourceSchema := resourceAlkiraConnectorAzureVnetThirdParty().Schema

	t.Run("field exists in schema", func(t *testing.T) {
		_, ok := resourceSchema["scale_group_id"]
		assert.True(t, ok, "scale_group_id must be present in schema")
	})

	t.Run("field is Optional TypeString", func(t *testing.T) {
		field := resourceSchema["scale_group_id"]
		assert.Equal(t, schema.TypeString, field.Type)
		assert.True(t, field.Optional)
		assert.False(t, field.Required)
		assert.False(t, field.Computed)
	})
}

func TestAzureVnetThirdPartyScaleGroupIdRequest(t *testing.T) {
	t.Run("ScaleGroupId populated when set", func(t *testing.T) {
		connector := alkira.AzureVnetThirdPartyConnector{
			ScaleGroupId: "sg-123",
		}
		assert.Equal(t, "sg-123", connector.ScaleGroupId)
	})

	t.Run("ScaleGroupId empty when not set", func(t *testing.T) {
		connector := alkira.AzureVnetThirdPartyConnector{}
		assert.Empty(t, connector.ScaleGroupId)
	})
}
