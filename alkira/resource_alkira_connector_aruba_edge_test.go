package alkira

import (
	"sort"
	"testing"

	"github.com/alkiranet/alkira-client-go/alkira"
	"github.com/stretchr/testify/require"
)

func TestGenerateArubaEdgeResource(t *testing.T) {
	expectedConnector := defaultConnector()
	r := resourceAlkiraConnectorArubaEdge()
	d := r.TestResourceData()
	setArubaEdgeResourceFields(&expectedConnector, d)

	actualConnector, err := generateConnectorArubaEdgeRequest(d, nil)
	require.NoError(t, err)

	require.ElementsMatch(t, expectedConnector.ArubaEdgeVrfMapping, actualConnector.ArubaEdgeVrfMapping)
	require.ElementsMatch(t, expectedConnector.Segments, actualConnector.Segments)
	requireArubeEdgeInstancesMatch(t, expectedConnector.Instances, actualConnector.Instances)
}

func defaultConnector() alkira.ConnectorArubaEdge {
	expectedArubaEdgeVrfMapping := []alkira.ArubaEdgeVRFMapping{
		alkira.ArubaEdgeVRFMapping{false, 0, "DEFAULT", false, 0},
		alkira.ArubaEdgeVRFMapping{true, 1, "DEFAULT1", true, 1},
		alkira.ArubaEdgeVRFMapping{false, 2, "DEFAULT2", false, 2},
	}
	expectedSegments := []string{"1", "2", "3"}
	expectedInstances := []alkira.ArubaEdgeInstance{
		alkira.ArubaEdgeInstance{"0", "accountName0", "credentialId0", "hostName0", "name0", "siteTag0"},
		alkira.ArubaEdgeInstance{"1", "accountName1", "credentialId1", "hostName1", "name1", "siteTag1"},
		alkira.ArubaEdgeInstance{"2", "accountName2", "credentialId2", "hostName2", "name2", "siteTag2"},
	}

	return alkira.ConnectorArubaEdge{
		ArubaEdgeVrfMapping: expectedArubaEdgeVrfMapping,
		Segments:            expectedSegments,
		Instances:           expectedInstances,
	}
}

func requireArubeEdgeInstancesMatch(t *testing.T, expected, actual []alkira.ArubaEdgeInstance) {
	require.Equal(t, len(expected), len(actual))

	//sort slices and make sure that they are in the same order. Name is a better choice than Id
	//since we do not populate the Id field for a create request.
	sort.Slice(expected, func(i, j int) bool { return expected[i].Name < expected[j].Name })
	sort.Slice(actual, func(i, j int) bool { return actual[i].Name < actual[j].Name })

	//When generating a request we do not include the id field. The backend API is finnicky about
	//its inclusion. Here we are just making sure that expected IDs equal the actual IDs so that we
	//can continue to use testify's ElementsMatch function.
	for i, _ := range actual {
		actual[i].Id = expected[i].Id
	}

	require.ElementsMatch(t, expected, actual)
}
