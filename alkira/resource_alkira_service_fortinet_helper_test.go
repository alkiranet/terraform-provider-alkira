package alkira

import (
	"testing"

	"github.com/alkiranet/alkira-client-go/alkira"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/stretchr/testify/require"
)

func TestExpandGenerateFortinetInstances(t *testing.T) {

	expectedName := "testName"
	expectedSerialNumber := "serialNumber"
	expectedCredentialId := "credentialId"

	expectedFortinetInstance := alkira.FortinetInstance{expectedName, expectedName, expectedSerialNumber, expectedCredentialId}

	m := makeMapFortinetInstance(expectedName, expectedSerialNumber, expectedCredentialId)
	m1 := makeMapFortinetInstance(expectedName+"1", expectedSerialNumber+"1", expectedCredentialId+"1")
	m2 := makeMapFortinetInstance(expectedName+"2", expectedSerialNumber+"2", expectedCredentialId+"2")
	mArr := []interface{}{m, m1, m2}

	r := resourceAlkiraFortinet()
	f := schema.HashResource(r)
	s := schema.NewSet(f, mArr)

	actual := expandFortinetInstances(s)
	require.Equal(t, len(actual), len(mArr))

	//Sets are unordered. We need to find our comparable item
	mIndex := 0
	for i, v := range actual {
		if v.Name == expectedName {
			mIndex = i
			break
		}
	}

	require.Equal(t, expectedFortinetInstance, actual[mIndex])
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
