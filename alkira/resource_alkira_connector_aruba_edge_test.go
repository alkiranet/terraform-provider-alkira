package alkira

import (
	"encoding/json"
	"strconv"
	"testing"

	"github.com/alkiranet/alkira-client-go/alkira"
	"github.com/stretchr/testify/require"
)

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
	for i := 0; i < num; i++ {
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
