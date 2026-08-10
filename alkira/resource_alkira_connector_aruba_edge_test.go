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

func TestConnectorArubaEdgeScaleGroupId(t *testing.T) {
	t.Run("omitted when empty", func(t *testing.T) {
		connector := alkira.ConnectorArubaEdge{
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
		connector := alkira.ConnectorArubaEdge{
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
		original := alkira.ConnectorArubaEdge{
			Name:         "test-connector",
			Cxp:          "US-WEST-1",
			ScaleGroupId: "sg-xyz789",
		}

		data, err := json.Marshal(original)
		require.NoError(t, err)

		var decoded alkira.ConnectorArubaEdge
		require.NoError(t, json.Unmarshal(data, &decoded))
		require.Equal(t, original.ScaleGroupId, decoded.ScaleGroupId)
	})
}

func TestArubaEdgeDefalteInstance(t *testing.T) {
	expectedInstances := generateNumArubaEdgeInstance(4)
	m := deflateArubaEdgeInstances(expectedInstances)

	for _, v := range m {
		require.Contains(t, v, "account_name")
		require.NotZero(t, v["account_name"])
		require.Contains(t, v, "credential_id")
		require.NotZero(t, v["credential_id"])
		require.Contains(t, v, "host_name")
		require.NotZero(t, v["host_name"])
		require.Contains(t, v, "id")
		require.NotZero(t, v["id"])
		require.Contains(t, v, "name")
		require.NotZero(t, v["name"])
		require.Contains(t, v, "site_tag")
		require.NotZero(t, v["site_tag"])
	}
}

//
// HELPERS
//

func generateNumArubaEdgeInstance(num int) []alkira.ArubaEdgeInstance {
	if num <= 0 {
		return nil
	}

	var instances []alkira.ArubaEdgeInstance
	for i := 1; i <= num; i++ {
		instances = append(instances, alkira.ArubaEdgeInstance{
			Id:           json.Number(strconv.Itoa(i)),
			AccountName:  "accountName" + strconv.Itoa(i),
			CredentialId: "credentialId" + strconv.Itoa(i),
			HostName:     "hostName" + strconv.Itoa(i),
			Name:         "name" + strconv.Itoa(i),
			SiteTag:      "siteTag" + strconv.Itoa(i),
		})
	}

	return instances
}

func TestArubaEdgeAllowListSchema(t *testing.T) {
	resourceSchema := resourceAlkiraConnectorArubaEdge().Schema

	t.Run("field exists and is Optional TypeSet of String", func(t *testing.T) {
		field, ok := resourceSchema["allow_list"]
		require.True(t, ok, "allow_list must be present in schema")
		assert.Equal(t, schema.TypeSet, field.Type)
		assert.True(t, field.Optional)
		assert.False(t, field.Required)

		elem, ok := field.Elem.(*schema.Schema)
		require.True(t, ok, "allow_list Elem must be a *schema.Schema")
		assert.Equal(t, schema.TypeString, elem.Type)
		assert.NotNil(t, elem.ValidateFunc, "allow_list elements must be validated")
	})
}

func TestArubaEdgeAllowListValidator(t *testing.T) {
	validate := resourceAlkiraConnectorArubaEdge().Schema["allow_list"].Elem.(*schema.Schema).ValidateFunc

	valid := []string{"10.0.0.0/24", "192.168.1.0/24", "10.0.0.5", "10.0.0.5/32"}
	for _, v := range valid {
		_, errs := validate(v, "allow_list")
		assert.Emptyf(t, errs, "expected %q to be accepted (IPv4 CIDR or IPv4 IP)", v)
	}

	invalid := []string{"not-an-ip", "10.0.0.0/33", "999.0.0.1", "2001:db8::/32", "::1", "::ffff:1.2.3.4", "10.0.0.5/24"}
	for _, v := range invalid {
		_, errs := validate(v, "allow_list")
		assert.NotEmptyf(t, errs, "expected %q to be rejected", v)
	}
}

func TestArubaEdgeAllowListExpand(t *testing.T) {
	t.Run("populated set produces expected slice", func(t *testing.T) {
		set := schema.NewSet(schema.HashString, []interface{}{
			"10.0.0.0/24", "10.0.0.5", "10.0.0.0/24",
		})
		got := convertTypeSetToStringList(set)
		assert.ElementsMatch(t, []string{"10.0.0.0/24", "10.0.0.5"}, got)
	})

	t.Run("empty set produces empty slice", func(t *testing.T) {
		set := schema.NewSet(schema.HashString, []interface{}{})
		got := convertTypeSetToStringList(set)
		assert.Empty(t, got)
	})
}

func TestArubaEdgeScaleGroupIdSchema(t *testing.T) {
	resourceSchema := resourceAlkiraConnectorArubaEdge().Schema

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

func TestArubaEdgeScaleGroupIdRequest(t *testing.T) {
	t.Run("ScaleGroupId populated when set", func(t *testing.T) {
		connector := alkira.ConnectorArubaEdge{
			ScaleGroupId: "sg-123",
		}
		assert.Equal(t, "sg-123", connector.ScaleGroupId)
	})

	t.Run("ScaleGroupId empty when not set", func(t *testing.T) {
		connector := alkira.ConnectorArubaEdge{}
		assert.Empty(t, connector.ScaleGroupId)
	})
}
