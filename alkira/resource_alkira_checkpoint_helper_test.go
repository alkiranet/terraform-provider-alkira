package alkira

import (
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"testing"

	"github.com/alkiranet/alkira-client-go/alkira"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/stretchr/testify/require"
)

func TestExpandCheckpointSegmentOptionsValid(t *testing.T) {
	//Setup
	expectedSegment := initCheckpointSegment()
	expectedZoneName := "testZoneName"
	expectedGroupName := "testGroup"

	m := makeMapCheckpointSegmentOptions(expectedSegment.Id, expectedZoneName, []interface{}{expectedGroupName})
	mArr := []interface{}{m}
	s := newSetFromCheckpointResource(mArr)

	mgs, err := expandCheckpointSegmentOptions(s, getCheckpointSegmentInTest)
	require.Nil(t, err)

	//Assertions
	expected := fmt.Sprintf(`{"%s":{"%s":["%s"]}}`, expectedSegment.Name, expectedZoneName, expectedGroupName)
	actual, err := json.Marshal(mgs)
	require.Nil(t, err)

	require.JSONEq(t, expected, string(actual))
}

func TestExpandCheckpointSegmentOptionsInvalid(t *testing.T) {
	//test nil set
	c, err := expandCheckpointSegmentOptions(nil, nil)
	require.Nil(t, c)
	require.Nil(t, err)

	//test empty set
	s := newSetFromCheckpointResource(nil)
	c, err = expandCheckpointSegmentOptions(s, nil)
	require.Nil(t, c)
	require.Nil(t, err)
}

func TestExpandCheckpointInstanceValid(t *testing.T) {
	instanceName := "testInstanceName"
	instanceCredentialId := "testInstanceCredentialId"

	mArr := makeNumCheckpointInstances(3, instanceName, instanceCredentialId)
	s := newSetFromCheckpointResource(mArr)

	instances := expandCheckpointInstances(s)
	require.NotZero(t, len(instances))
	require.Equal(t, len(instances), len(mArr))
	requireAllKeyValuesCheckpointInstances(t, instances, mArr)
}

func TestExpandCheckpointInstanceInvalid(t *testing.T) {
	//test nil Set
	c := expandCheckpointInstances(nil)
	require.Nil(t, c)

	//test empty Set
	s := newSetFromCheckpointResource(nil)
	c = expandCheckpointInstances(s)
	require.Nil(t, c)
}

func TestConvertCheckpointSegmentOptionsValid(t *testing.T) {
	expectedId := 11
	expectedZoneName := "testZoneName"
	expectedGroups := []string{"group1", "group2", "group3"}

	mArr := makeNumCheckpointSegmentOptions(3, expectedId, expectedZoneName, expectedGroups)
	s := newSetFromCheckpointResource(mArr)

	_, actualZoneName, actualGroups, err := convertCheckpointSegmentOptions(s, getCheckpointSegmentInTest)
	if err != nil {
		fmt.Println(err)
	}

	require.Equal(t, expectedZoneName, actualZoneName)
	require.Equal(t, expectedGroups, actualGroups)
}

func TestConvertCheckpointSegmentOptionsInvalid(t *testing.T) {
	//test nil Set
	actualSegmentName, actualZoneName, actualGroups, err := convertCheckpointSegmentOptions(nil, getCheckpointSegmentInTest)
	require.Error(t, err)
	require.Empty(t, actualSegmentName)
	require.Empty(t, actualZoneName)
	require.Empty(t, actualGroups)

	//test empty Set
	s := newSetFromCheckpointResource(nil)
	actualSegmentName, actualZoneName, actualGroups, err = convertCheckpointSegmentOptions(s, getCheckpointSegmentInTest)
	require.Nil(t, err)
	require.Empty(t, actualSegmentName)
	require.Empty(t, actualZoneName)
	require.Empty(t, actualGroups)

}

func TestConvertCheckpointSegmentOptionsGetError(t *testing.T) {
	//test get Func is nil
	mArr := makeNumCheckpointSegmentOptions(1, 0, "", []string{})
	s := newSetFromCheckpointResource(mArr)
	actualSegmentName, actualZoneName, actualGroups, err := convertCheckpointSegmentOptions(s, nil)
	require.Error(t, err)
	require.Empty(t, actualSegmentName)
	require.Empty(t, actualZoneName)
	require.Empty(t, actualGroups)

	//test get Func returns error
	s = newSetFromCheckpointResource(mArr)
	actualSegmentName, actualZoneName, actualGroups, err = convertCheckpointSegmentOptions(s, getCheckpointSegmentError)
	require.Error(t, err)
	require.Empty(t, actualSegmentName)
	require.Empty(t, actualZoneName)
	require.Empty(t, actualGroups)
}

func TestExpandCheckpointManagementServerValid(t *testing.T) {
	//Setup
	expectedManagementServer := initCheckpointTestManagementServer()
	s := newSetFromCheckpointResource([]interface{}{deflateCheckpointManagementServer(expectedManagementServer)})

	//Assertions
	actualManagementServer, err := expandCheckpointManagementServer(s, getCheckpointSegmentInTest)
	require.Nil(t, err)
	require.Equal(t, expectedManagementServer.Ips, actualManagementServer.Ips)
	require.Equal(t, expectedManagementServer.ConfigurationMode, actualManagementServer.ConfigurationMode)
}

func TestExpandCheckpointManagementServerGetError(t *testing.T) {
	m := make(map[string]interface{})
	m["segment_id"] = 0
	s := newSetFromCheckpointResource([]interface{}{m})

	_, err := expandCheckpointManagementServer(s, getCheckpointSegmentError)
	require.Error(t, err)
}

func TestDeflateManagementServerValid(t *testing.T) {
	expected := initCheckpointTestManagementServer()

	m := deflateCheckpointManagementServer(expected)
	require.Equal(t, m["configuration_mode"], expected.ConfigurationMode)
	require.Equal(t, m["credential_id"], expected.CredentialId)
	require.Equal(t, m["domain"], expected.Domain)
	require.Equal(t, m["global_cidr_list_id"], expected.GlobalCidrListId)
	require.Equal(t, m["ips"], expected.Ips)
	require.Equal(t, m["reachability"], expected.Reachability)
	require.Equal(t, m["segment"], expected.Segment)
	require.Equal(t, m["segment_id"], expected.SegmentId)
	require.Equal(t, m["type"], expected.Type)
	require.Equal(t, m["user_name"], expected.UserName)
}

func TestDeflateCheckpointSegmentOptionsValid(t *testing.T) {
	expectedSegmentName := "0"
	expectedZoneName := "zoneName"
	expectedGroups := []string{"1", "2", "3", "4"}

	nameToZone := make(alkira.CheckpointSegmentNameToZone)
	zoneToGroups := make(alkira.CheckpointZoneToGroups)

	nameToZone[expectedSegmentName] = zoneToGroups
	zoneToGroups[expectedZoneName] = expectedGroups

	actual, err := deflateCheckpointSegmentOptions(nameToZone, getCheckpointSegmentInTest)
	require.Nil(t, err)
	require.Len(t, actual, 1)
}

func TestDeflateCheckpointSegmentOptionsGetError(t *testing.T) {
	_, err := deflateCheckpointSegmentOptions(nil, getCheckpointSegmentError)
	require.Error(t, err)
}

func TestDeflateCheckpointInstances(t *testing.T) {
	numInstances := 9
	testName := "testName"
	testCredId := "testCredId"
	c := []alkira.CheckpointInstance{}

	for i := 0; i < numInstances; i++ {
		c = append(c, alkira.CheckpointInstance{
			Name:         testName + strconv.Itoa(i),
			CredentialId: testCredId + strconv.Itoa(i),
		})
	}

	m := deflateCheckpointInstances(c)
	require.Len(t, m, len(c))

	for i, _ := range m {
		require.Contains(t, m[i]["name"], testName+strconv.Itoa(i))
		require.Contains(t, m[i]["credential_id"], testCredId+strconv.Itoa(i))
	}
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

func makeNumCheckpointInstances(num int, name string, id string) []interface{} {
	mArr := []interface{}{}
	for i := 0; i < num; i++ {
		name := name + strconv.Itoa(i)
		credentialId := id + strconv.Itoa(i)

		mArr = append(mArr, makeMapCheckpointInstance(name, credentialId))
	}

	return mArr
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
		SegmentId:         0,
		Type:              "type",
		UserName:          "userName",
	}
}

func initCheckpointSegment() alkira.Segment {
	return alkira.Segment{
		Asn:     0,
		Id:      0,
		IpBlock: "10.0.1.1/23",
		Name:    "segmentName",
		ReservePublicIPsForUserAndSiteConnectivity: false,
	}
}

//actual := alkira.CheckpointManagementServer{
//	ConfigurationMode: "configurationMode",
//	CredentialId:      "credentialId",
//	Domain:            "domain",
//	GlobalCidrListId:  "globalCidrListId",
//	Ips:               []string{"1", "2", "3"},
//	Reachability:      "reachability",
//	Segment:           "segment",
//	SegmentId:         "segmentId",
//	Type:              "type",
//	UserName:          "userName",
//}
