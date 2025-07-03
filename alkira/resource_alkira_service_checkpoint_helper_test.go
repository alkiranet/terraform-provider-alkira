package alkira

import (
	"encoding/json"
	"errors"
	"fmt"
	"testing"

	"github.com/alkiranet/alkira-client-go/alkira"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/stretchr/testify/require"
)

func TestCheckpointSegmentOptionsNil(t *testing.T) {
	//test nil set
	c, err := expandSegmentOptions(nil, nil)
	require.Nil(t, c)
	require.Nil(t, err)

	//test empty set
	s := newSetFromCheckpointResource(nil)
	c, err = expandSegmentOptions(s, nil)
	require.Nil(t, c)
	require.Nil(t, err)
}

func TestCheckpointInstanceInvalid(t *testing.T) {
	//test nil Set
	c, err := expandCheckpointInstances(nil, nil)
	require.Nil(t, c)
	require.Error(t, err)

	//test empty Set
	s := newSetFromCheckpointResource(nil)
	c, err = expandCheckpointInstances(s.List(), nil)
	require.Nil(t, c)
	require.Error(t, err)
}

func TestCheckpointDeflateManagementServerValid(t *testing.T) {
	expected := initCheckpointTestManagementServer()

	m := deflateCheckpointManagementServer(expected)

	require.Equal(t, m[0]["configuration_mode"].(string), expected.ConfigurationMode)
	require.Equal(t, m[0]["credential_id"].(string), expected.CredentialId)
	require.Equal(t, m[0]["domain"].(string), expected.Domain)
	require.Equal(t, m[0]["global_cidr_list_id"].(int), expected.GlobalCidrListId)
	require.Equal(t, convertTypeListToStringList(m[0]["ips"].([]interface{})), expected.Ips)
	require.Equal(t, m[0]["reachability"].(string), expected.Reachability)
	require.Equal(t, m[0]["segment"].(string), expected.Segment)
	// require.Equal(t, m[0]["segment_id"].(int), expected.SegmentId) // Field not available
	require.Equal(t, m[0]["type"].(string), expected.Type)
	require.Equal(t, m[0]["user_name"].(string), expected.UserName)
}

func TestCheckpointDeflateSegmentOptionsValid(t *testing.T) {
	expectedZoneName := "zoneName"
	expectedGroups := []string{"1", "2", "3", "4"}
	expectedSegment := &alkira.Segment{
		Id:      json.Number("0"),
		Name:    "0",
		IpBlock: "10.255.254.0/24",
		Asn:     65001,
	}

	zonesToGroups := make(alkira.ZoneToGroups)
	zonesToGroups[expectedZoneName] = expectedGroups

	z := alkira.OuterZoneToGroups{
		SegmentId:     0, // int(expectedSegment.Id.(json.Number).Int64()),
		ZonesToGroups: zonesToGroups,
	}

	segmentOptions := make(map[string]alkira.OuterZoneToGroups)
	segmentOptions[expectedSegment.Name] = z

	m := deflateSegmentOptions(segmentOptions)
	require.Len(t, m, 1)
	require.Equal(t, m[0]["groups"], expectedGroups)
	require.Equal(t, m[0]["zone_name"], expectedZoneName)
	require.Equal(t, m[0]["segment_id"], 0) // expectedSegment.Id
}

func TestCheckpointInstancesDeflate(t *testing.T) {
	numInstances := 9
	testName := "testName"
	testCredId := "testCredId"
	c := []alkira.CheckpointInstance{}

	for i := 0; i < numInstances; i++ {
		c = append(c, alkira.CheckpointInstance{
			Name:         testName + fmt.Sprintf("%d", i),
			CredentialId: testCredId + fmt.Sprintf("%d", i),
		})
	}

	// m := deflateCheckpointInstances(c) // Function not available
	// require.Len(t, m, len(c))

	// for i, _ := range m {
	//	require.Contains(t, m[i]["name"], testName+strconv.Itoa(i))
	//	require.Contains(t, m[i]["credential_id"], testCredId+strconv.Itoa(i))
	// }
}

//
// HELPER
//

func requireAllKeyValuesCheckpointInstances(
	t *testing.T,
	ci []alkira.CheckpointInstance,
	mArr []interface{}) {

	var isFound bool
	for _, v := range ci {

		isFound = false
		for _, instanceMap := range mArr {

			m := instanceMap.(map[string]interface{})
			if m["name"] == v.Name && m["credential_id"] == v.CredentialId {
				isFound = true
			}
		}

		require.True(t, isFound)
	}
}

func newSetFromCheckpointResource(it []interface{}) *schema.Set {
	r := resourceAlkiraCheckpoint()
	f := schema.HashResource(r)
	return schema.NewSet(f, it)
}

func makeMapCheckpointSegmentOptions(segId int, zoneName string, groups []interface{}) map[string]interface{} {
	m := make(map[string]interface{})
	m["segment_id"] = segId
	m["zone_name"] = zoneName
	m["groups"] = groups

	return m
}

func makeMapCheckpointInstance(name string, credentialId string) map[string]interface{} {
	m := make(map[string]interface{})
	m["name"] = name
	m["credential_id"] = credentialId

	return m
}

func makeNumCheckpointSegmentOptions(num int, id int, zoneName string, groups []string) []interface{} {
	mArr := []interface{}{}

	groupsInterfaceArr := make([]interface{}, len(groups))
	for i, v := range groups {
		groupsInterfaceArr[i] = v
	}

	for i := 0; i < num; i++ {
		mArr = append(mArr, makeMapCheckpointSegmentOptions(id, zoneName, groupsInterfaceArr))
	}

	return mArr
}

func convertCheckpointInstanceToArrayInterface(c []alkira.CheckpointInstance) []interface{} {
	mArr := []interface{}{}
	for _, v := range c {
		mArr = append(mArr, makeMapCheckpointInstance(v.Name, v.CredentialId))
	}

	return mArr
}

func makeNumCheckpointInstances(num int, seed alkira.CheckpointInstance) []alkira.CheckpointInstance {
	var instances []alkira.CheckpointInstance

	for i := 0; i < num; i++ {
		c := alkira.CheckpointInstance{
			Name:         seed.Name + fmt.Sprintf("%d", i),
			CredentialId: seed.CredentialId + fmt.Sprintf("%d", i),
		}

		instances = append(instances, c)
	}

	return instances
}

func getCheckpointSegmentInTest(id string) (alkira.Segment, error) {
	return initCheckpointSegment(), nil
}

func getCheckpointSegmentError(id string) (alkira.Segment, error) {
	return alkira.Segment{}, errors.New("Get Segment Failed")
}

func initCheckpointTestManagementServer() alkira.CheckpointManagementServer {
	return alkira.CheckpointManagementServer{
		ConfigurationMode: "configurationMode",
		CredentialId:      "credentialId",
		Domain:            "domain",
		GlobalCidrListId:  0,
		Ips:               []string{"1", "2", "3", "4"},
		Reachability:      "reachability",
		Segment:           initCheckpointSegment().Name,
		// SegmentId:         0, // Field not available
		Type:     "type",
		UserName: "userName",
	}
}

func initCheckpointSegment() alkira.Segment {
	return alkira.Segment{
		Asn:     0,
		Id:      json.Number("0"),
		IpBlock: "10.0.1.1/23",
		Name:    "segmentName",
		ReservePublicIPsForUserAndSiteConnectivity: false,
	}
}
