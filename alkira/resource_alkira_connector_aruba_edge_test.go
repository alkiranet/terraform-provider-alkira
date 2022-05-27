package alkira

import (
	"encoding/json"
	"sort"
	"strconv"
	"testing"

	"github.com/alkiranet/alkira-client-go/alkira"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
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

func TestExpandArubaEdgeVrfMapping(t *testing.T) {
	expectedArubaEdgeVrfMappings := generateNumArubaEdgeVrfMapping(3)

	r := resourceAlkiraConnectorArubaEdge()
	s := schema.NewSet(schema.HashResource(r), makeArrInterfaceFromArubaEdgeVrf(expectedArubaEdgeVrfMappings))

	getSegmentFn := func(id string) (alkira.Segment, error) {
		i, err := strconv.Atoi(id)
		require.Nil(t, err)

		return alkira.Segment{Id: i, Name: id}, nil
	}

	actualArubaEdgeVrfMappings, err := expandArubaEdgeVrfMappings(s, getSegmentFn)
	require.NoError(t, err)

	//Sets are not guaranteed an order
	sortArubaEdgeVrfMappingBySegmentName(expectedArubaEdgeVrfMappings)
	sortArubaEdgeVrfMappingBySegmentName(actualArubaEdgeVrfMappings)

	require.Equal(t, expectedArubaEdgeVrfMappings, actualArubaEdgeVrfMappings)
}

//
// HELPERS
//

func defaultConnector() alkira.ConnectorArubaEdge {
	return alkira.ConnectorArubaEdge{
		ArubaEdgeVrfMapping: generateNumArubaEdgeVrfMapping(3),
		Segments:            []string{"1", "2", "3"},
		Instances:           generateNumArubaEdgeInstance(4),
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

func generateNumArubaEdgeVrfMapping(num int) []alkira.ArubaEdgeVRFMapping {
	if num <= 0 {
		return nil
	}

	var instances []alkira.ArubaEdgeVRFMapping
	for i := 0; i < num; i++ {
		//None of these values should be there zero values to verify that we are setting them
		instances = append(instances, alkira.ArubaEdgeVRFMapping{
			AdvertiseOnPremRoutes:       true,
			AlkiraSegmentId:             i + 1,
			ArubaEdgeConnectSegmentName: strconv.Itoa(i + 1),
			DisableInternetExit:         true,
			GatewayBgpAsn:               i + 1,
		})
	}

	return instances
}

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

func makeMapArubaEdgeVrfMapping(ar alkira.ArubaEdgeVRFMapping, id int) map[string]interface{} {
	m := make(map[string]interface{})
	m["advertise_on_prem_routes"] = ar.AdvertiseOnPremRoutes
	m["segment_id"] = strconv.Itoa(ar.AlkiraSegmentId)
	m["aruba_edge_connect_segment_id"] = strconv.Itoa(id)
	m["disable_internet_exit"] = ar.DisableInternetExit
	m["gateway_gbp_asn"] = ar.GatewayBgpAsn

	return m
}

func makeArrInterfaceFromArubaEdgeVrf(ar []alkira.ArubaEdgeVRFMapping) []interface{} {
	var i []interface{}
	for _, a := range ar {
		m := makeMapArubaEdgeVrfMapping(a, a.AlkiraSegmentId)
		i = append(i, m)
	}

	return i
}

func sortArubaEdgeVrfMappingBySegmentName(ar []alkira.ArubaEdgeVRFMapping) {
	sort.Slice(ar, func(i, j int) bool {
		return ar[i].ArubaEdgeConnectSegmentName < ar[j].ArubaEdgeConnectSegmentName
	})
}
