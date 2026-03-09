package alkira

import (
	"encoding/json"
	"strconv"
	"testing"

	"github.com/alkiranet/alkira-client-go/alkira"
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
