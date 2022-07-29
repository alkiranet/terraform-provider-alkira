package alkira

import (
	"strconv"
	"testing"

	"github.com/alkiranet/alkira-client-go/alkira"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/stretchr/testify/require"
)

func TestExpandGenerateFortinetInstances(t *testing.T) {

	//TODO(mac): uncomment and fix test
	expectedName := "testName"
	expectedSerialNumber := "serialNumber"
	expectedCredentialId := "credentialId"

	expectedFortinetInstance := alkira.FortinetInstance{
		Name:         expectedName,
		HostName:     expectedName,
		SerialNumber: expectedSerialNumber,
		CredentialId: expectedCredentialId,
	}

	m := makeMapFortinetInstance(expectedName, expectedSerialNumber, expectedCredentialId)
	m1 := makeMapFortinetInstance(expectedName+"1", expectedSerialNumber+"1", expectedCredentialId+"1")
	m2 := makeMapFortinetInstance(expectedName+"2", expectedSerialNumber+"2", expectedCredentialId+"2")
	mArr := []interface{}{m, m1, m2}

	r := resourceAlkiraServiceFortinet()
	f := schema.HashResource(r)
	s := schema.NewSet(f, mArr)

	actual := expandFortinetInstances(s)
	require.Equal(t, len(actual), len(mArr))

	//Sets are unordered. We need to find our comparable item
	//mIndex := new(int)
	var mIndex int
	for i, v := range actual {
		if v.Name == expectedName {
			mIndex = i
			break
		}
	}

	require.Equal(t, expectedFortinetInstance, actual[mIndex])
}

func TestExpandFortinetZone(t *testing.T) {
	expectedName := "ZONE_NAME"
	expectedGroupName := "GROUP_NAME"

	ifc := makeNumMapFortinetZone(9, expectedName, expectedGroupName)

	names, groups := getNamesAndGroups(ifc)

	s := schema.NewSet(schema.HashResource(zoneResourceFromFortinet()), ifc)

	// expand fortinet zone
	fz := expandFortinetZone(s)

	// test for name and group inclusion
	for k, v := range fz {
		require.Contains(t, names, k)
		require.Contains(t, groups, v)
	}
}

//
// TEST HELPERS
//

func makeMapFortinetInstance(name string, serialNumber string, credentialId string) map[string]interface{} {
	m := make(map[string]interface{})
	m["name"] = name
	m["serial_number"] = serialNumber
	m["credential_id"] = credentialId

	return m
}

func makeMapFortinetZone(name string, groups []interface{}) map[string]interface{} {
	m := make(map[string]interface{})
	m["name"] = name
	m["groups"] = groups
	return m
}

// makeNumMapFortinetZone could be adjusted for variability in the number of groups included
// in each map. For now it is fixed at one.
func makeNumMapFortinetZone(num int, baseName string, baseGroupsName string) []interface{} {
	var ifc []interface{}

	for i := 0; i < num; i++ {
		postfixedName := baseName + strconv.Itoa(i)
		postFixedGroupsName := []interface{}{baseGroupsName + strconv.Itoa(i)}
		m := makeMapFortinetZone(postfixedName, postFixedGroupsName)
		ifc = append(ifc, m)
	}

	return ifc
}

func getNamesAndGroups(i []interface{}) ([]string, [][]string) {
	var nameVals []string
	var groupVals [][]string

	for _, v := range i {
		name := v.(map[string]interface{})["name"]
		groups := v.(map[string]interface{})["groups"]

		nameVals = append(nameVals, name.(string))

		var s []string
		for _, p := range groups.([]interface{}) {
			s = append(s, p.(string))
		}
		groupVals = append(groupVals, s)
	}

	return nameVals, groupVals
}

// if tests break because of zoneResourceFromFortinet function it means the schema for our
// fortinet resource has changed. In that instance we would need to adjust tests and make
// sure that we haven't broken backward compatability.
func zoneResourceFromFortinet() *schema.Resource {
	r := resourceAlkiraServiceFortinet()
	return r.Schema["segment_options"].Elem.(*schema.Resource).Schema["zone"].Elem.(*schema.Resource)
}
