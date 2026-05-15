package alkira

import (
	"encoding/json"
	"strconv"
	"testing"

	"github.com/alkiranet/alkira-client-go/alkira"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestConnectorPrismaSDWANScaleGroupId(t *testing.T) {
	t.Run("omitted when empty", func(t *testing.T) {
		connector := alkira.ConnectorPrismaSDWAN{
			Name: "test-connector",
			Cxp:  "US-EAST-2",
		}

		data, err := json.Marshal(connector)
		require.NoError(t, err)

		var m map[string]interface{}
		require.NoError(t, json.Unmarshal(data, &m))
		require.NotContains(t, m, "scaleGroupId", "scaleGroupId should be omitted when empty")
	})

	t.Run("included when set", func(t *testing.T) {
		connector := alkira.ConnectorPrismaSDWAN{
			Name:         "test-connector",
			Cxp:          "US-EAST-2",
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
		original := alkira.ConnectorPrismaSDWAN{
			Name:         "test-connector",
			Cxp:          "US-WEST-1",
			ScaleGroupId: "sg-xyz789",
		}

		data, err := json.Marshal(original)
		require.NoError(t, err)

		var decoded alkira.ConnectorPrismaSDWAN
		require.NoError(t, json.Unmarshal(data, &decoded))
		require.Equal(t, original.ScaleGroupId, decoded.ScaleGroupId)
	})
}

func TestPrismaSDWANDeflateInstances(t *testing.T) {
	expectedInstances := generateNumPrismaSDWANInstances(4)
	result := expandPrismaSDWANInstances(deflateToInterfaceSlice(expectedInstances))

	require.Len(t, result, 4)
	for i, inst := range result {
		assert.Equal(t, "hostName"+strconv.Itoa(i+1), inst.HostName)
		assert.Equal(t, "credentialId"+strconv.Itoa(i+1), inst.CredentialId)
		assert.Equal(t, "ion-3102v", inst.IonModel)
		assert.Equal(t, "6.4.1", inst.Version)
		assert.Equal(t, i+1, inst.Id)
	}
}

func TestPrismaSDWANScaleGroupIdSchema(t *testing.T) {
	resourceSchema := resourceAlkiraConnectorPrismaSDWAN().Schema

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

func TestPrismaSDWANScaleGroupIdRequest(t *testing.T) {
	t.Run("ScaleGroupId populated when set", func(t *testing.T) {
		connector := alkira.ConnectorPrismaSDWAN{
			ScaleGroupId: "sg-123",
		}
		assert.Equal(t, "sg-123", connector.ScaleGroupId)
	})

	t.Run("ScaleGroupId empty when not set", func(t *testing.T) {
		connector := alkira.ConnectorPrismaSDWAN{}
		assert.Empty(t, connector.ScaleGroupId)
	})
}

func TestPrismaSDWANInstancesSchema(t *testing.T) {
	resourceSchema := resourceAlkiraConnectorPrismaSDWAN().Schema

	t.Run("instances field exists", func(t *testing.T) {
		_, ok := resourceSchema["instances"]
		assert.True(t, ok, "instances must be present in schema")
	})

	t.Run("instances is Required TypeList", func(t *testing.T) {
		field := resourceSchema["instances"]
		assert.Equal(t, schema.TypeList, field.Type)
		assert.True(t, field.Required)
	})

	t.Run("instance sub-fields", func(t *testing.T) {
		instanceSchema := resourceSchema["instances"].Elem.(*schema.Resource).Schema

		assert.Contains(t, instanceSchema, "id")
		assert.True(t, instanceSchema["id"].Computed)

		assert.Contains(t, instanceSchema, "host_name")
		assert.True(t, instanceSchema["host_name"].Required)

		assert.Contains(t, instanceSchema, "credential_id")
		assert.True(t, instanceSchema["credential_id"].Required)

		assert.Contains(t, instanceSchema, "ion_model")
		assert.True(t, instanceSchema["ion_model"].Required)

		assert.Contains(t, instanceSchema, "version")
		assert.True(t, instanceSchema["version"].Required)
	})
}

func TestPrismaSDWANTargetSegmentSchema(t *testing.T) {
	resourceSchema := resourceAlkiraConnectorPrismaSDWAN().Schema

	t.Run("target_segment field exists", func(t *testing.T) {
		_, ok := resourceSchema["target_segment"]
		assert.True(t, ok, "target_segment must be present in schema")
	})

	t.Run("target_segment is Required TypeSet", func(t *testing.T) {
		field := resourceSchema["target_segment"]
		assert.Equal(t, schema.TypeSet, field.Type)
		assert.True(t, field.Required)
	})

	t.Run("target_segment sub-fields", func(t *testing.T) {
		segSchema := resourceSchema["target_segment"].Elem.(*schema.Resource).Schema

		assert.Contains(t, segSchema, "segment_id")
		assert.True(t, segSchema["segment_id"].Required)

		assert.Contains(t, segSchema, "gateway_bgp_asn")
		assert.True(t, segSchema["gateway_bgp_asn"].Required)

		assert.Contains(t, segSchema, "vrf_name")
		assert.True(t, segSchema["vrf_name"].Required)

		assert.Contains(t, segSchema, "advertise_on_prem_routes")
		assert.True(t, segSchema["advertise_on_prem_routes"].Optional)

		assert.Contains(t, segSchema, "advertise_default_route")
		assert.True(t, segSchema["advertise_default_route"].Optional)
	})
}

//
// HELPERS
//

func generateNumPrismaSDWANInstances(num int) []alkira.ConnectorPrismaSDWANInstance {
	if num <= 0 {
		return nil
	}

	var instances []alkira.ConnectorPrismaSDWANInstance
	for i := 1; i <= num; i++ {
		instances = append(instances, alkira.ConnectorPrismaSDWANInstance{
			Id:           i,
			CredentialId: "credentialId" + strconv.Itoa(i),
			HostName:     "hostName" + strconv.Itoa(i),
			IonModel:     "ion-3102v",
			Version:      "6.4.1",
		})
	}

	return instances
}

func deflateToInterfaceSlice(instances []alkira.ConnectorPrismaSDWANInstance) []interface{} {
	var result []interface{}
	for _, inst := range instances {
		result = append(result, map[string]interface{}{
			"host_name":     inst.HostName,
			"credential_id": inst.CredentialId,
			"ion_model":     inst.IonModel,
			"version":       inst.Version,
			"id":            inst.Id,
		})
	}
	return result
}
