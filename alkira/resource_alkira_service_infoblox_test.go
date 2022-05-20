package alkira

import (
	"encoding/json"
	"sort"
	"strconv"
	"testing"

	"github.com/alkiranet/alkira-client-go/alkira"
	"github.com/stretchr/testify/require"
)

func TestExpandDeflateInfoblox(t *testing.T) {
	expectedInstances := makeInfobloxInstances(4, true)
	expectedGridMaster := makeInfobloxGridMaster("gm", false)
	expectedAnycast := makeInfobloxAnycast(3, 5, "backupVal", "ipsVal", false)

	expectedInfoblox := &alkira.Infoblox{
		AnyCast:    expectedAnycast,
		Instances:  expectedInstances,
		GridMaster: expectedGridMaster,
	}

	createCredentialFn := func(name string, ctype alkira.CredentialType, credential interface{}) (string, error) {
		return "credentialId", nil
	}

	d := resourceAlkiraInfoblox().TestResourceData()
	setAllInfobloxResourceFields(d, expectedInfoblox)

	actualInfoblox, err := generateInfobloxRequest(d, nil, createCredentialFn)
	require.NoError(t, err)

	require.Equal(t, expectedGridMaster, actualInfoblox.GridMaster)
	require.Equal(t, expectedAnycast, actualInfoblox.AnyCast)
	requireInfobloxInstancesEqual(t, expectedInstances, actualInfoblox.Instances)
}

func makeInfobloxInstances(num int, anycast bool) []alkira.InfobloxInstance {
	var ins []alkira.InfobloxInstance
	for i := 0; i < num; i++ {
		in := makeInfobloxInstance(strconv.Itoa(i), i, anycast)
		ins = append(ins, in)
	}

	return ins
}

func makeInfobloxInstance(prefix string, id int, anycast bool) alkira.InfobloxInstance {
	credentialId, _ := createInfobloxInstanceCredentialInTest("", alkira.CredentialTypeAkamaiProlexic, nil)
	return alkira.InfobloxInstance{
		AnyCastEnabled:     anycast,
		ConfiguredMasterIp: prefix + "ConfiguredMasterIp",
		CredentialId:       credentialId,
		HostName:           prefix + "HostName",
		Id:                 json.Number(strconv.Itoa(id)),
		InternalName:       prefix + "InternalName",
		LanPrefix:          prefix + "LanPrefix",
		ManagementPrefix:   prefix + "ManagementPrefix",
		Model:              prefix + "Model",
		Name:               prefix + "Name",
		ProductId:          prefix + "ProductId",
		PublicIp:           prefix + "PublicIp",
		Type:               prefix + "Type",
		Version:            prefix + "Version",
	}
}

func makeInfobloxGridMaster(prefix string, external bool) alkira.InfobloxGridMaster {
	credentialId, _ := createInfobloxInstanceCredentialInTest("", alkira.CredentialTypeAkamaiProlexic, nil)
	return alkira.InfobloxGridMaster{
		External:                 external,
		GridMasterCredentialId:   credentialId,
		Ip:                       prefix + "Ip",
		Name:                     prefix + "Name",
		SharedSecretCredentialId: credentialId,
	}
}

func makeInfobloxAnycast(lenBackups, lenIps int, backupValue, ipsValue string, external bool) alkira.InfobloxAnycast {
	backups := make([]string, lenBackups)
	for i, _ := range backups {
		backups[i] = strconv.Itoa(i) + backupValue
	}

	ips := make([]string, lenIps)
	for i, _ := range ips {
		ips[i] = strconv.Itoa(i) + ipsValue
	}

	return alkira.InfobloxAnycast{
		BackupCxps: backups,
		Enabled:    external,
		Ips:        ips,
	}
}

func sortInfobloxInstancesByHost(ins []alkira.InfobloxInstance) {
	sort.Slice(ins, func(i, j int) bool { return ins[i].HostName < ins[j].HostName })
}

//This is required because the infoblox client has quite a few fields that are not currently in use
//in the backend. Here we are making comparisons only to the fields that the terraform user might
//input.
func requireInfobloxInstancesEqual(t *testing.T, expected, actual []alkira.InfobloxInstance) {
	require.Equal(t, len(expected), len(actual))

	//It doesn't really matter how we sort as long as they are in the same order.
	sortInfobloxInstancesByHost(expected)
	sortInfobloxInstancesByHost(actual)

	for i, v := range actual {
		require.Equal(t, expected[i].AnyCastEnabled, v.AnyCastEnabled)
		require.Equal(t, expected[i].CredentialId, v.CredentialId)
		require.Equal(t, expected[i].HostName, v.HostName)
		require.Equal(t, expected[i].Model, v.Model)
		require.Equal(t, expected[i].Type, v.Type)
		require.Equal(t, expected[i].Version, v.Version)
	}
}

func createInfobloxInstanceCredentialInTest(name string, ctype alkira.CredentialType, credential interface{}) (string, error) {
	return "credentialId", nil
}
