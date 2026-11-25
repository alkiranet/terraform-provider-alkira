package alkira

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// scale_group_id is Optional only: not ForceNew, so a change is never a destroy and re-create
// (the API rejects a change after provisioning), and not Computed, so drift is shown.
func TestSdwanConnectorScaleGroupIdSchema(t *testing.T) {
	resources := map[string]*schema.Resource{
		"cisco":    resourceAlkiraConnectorCiscoSdwan(),
		"fortinet": resourceAlkiraConnectorFortinetSdwan(),
		"versa":    resourceAlkiraConnectorVersaSdwan(),
		"vmware":   resourceAlkiraConnectorVmwareSdwan(),
		"juniper":  resourceAlkiraConnectorJuniperSdwan(),
	}

	for name, r := range resources {
		t.Run(name, func(t *testing.T) {
			field, ok := r.Schema["scale_group_id"]
			require.True(t, ok, "scale_group_id must be present in schema")
			assert.Equal(t, schema.TypeString, field.Type)
			assert.True(t, field.Optional)
			assert.False(t, field.Required)
			assert.False(t, field.Computed)
			assert.False(t, field.ForceNew)
		})
	}
}

func TestVersaSdwanScaleGroupIdRequest(t *testing.T) {
	r := resourceAlkiraConnectorVersaSdwan()

	t.Run("set", func(t *testing.T) {
		d := schema.TestResourceDataRaw(t, r.Schema, map[string]interface{}{"scale_group_id": "sg-123"})
		connector, err := generateConnectorVersaSdwanRequest(d)
		require.NoError(t, err)
		assert.Equal(t, "sg-123", connector.ScaleGroupId)
	})

	t.Run("not set", func(t *testing.T) {
		d := schema.TestResourceDataRaw(t, r.Schema, map[string]interface{}{})
		connector, err := generateConnectorVersaSdwanRequest(d)
		require.NoError(t, err)
		assert.Empty(t, connector.ScaleGroupId)
	})
}
