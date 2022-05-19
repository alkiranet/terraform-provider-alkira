package alkira

import (
	"encoding/json"
	"fmt"
	"sort"
	"strconv"
	"testing"

	"github.com/alkiranet/alkira-client-go/alkira"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/kr/pretty"
	"github.com/stretchr/testify/require"
)

func TestExpandDeflateInfoblox(t *testing.T) {
	expectedInstances := makeInfobloxInstances(4, true)
	//fmt.Println("expectedInstances: ", len(expectedInstances))
	expectedGridMaster := makeInfobloxGridMaster("gm", false)
	expectedAnycast := makeInfobloxAnycast(3, 5, "backupVal", "ipsVal", false)

	expectedInfoblox := &alkira.Infoblox{
		AnyCast:    expectedAnycast,
		Instances:  expectedInstances,
		GridMaster: expectedGridMaster,
	}
	expectedInfoblox = expectedInfoblox

	createCredentialFn := func(name string, ctype alkira.CredentialType, credential interface{}) (string, error) {
		return "credentialId", nil
	}

	d := resourceAlkiraInfoblox().TestResourceData()
	d.Set("instances", deflateInfobloxInstances(expectedInstances))
	d.Set("grid_master", deflateInfobloxGridMaster(expectedGridMaster))
	d.Set("anycast", deflateInfobloxAnycast(expectedAnycast))

	fmt.Println("SKUMPS SKUMPS A TOAST TO THIS NIGHT")
	pretty.Println(d.Get("instances").(*schema.Set))
	pretty.Println(d.Get("anycast").(*schema.Set))
	pretty.Println(d.Get("grid_master").(*schema.Set))

	actualInfoblox, err := generateInfobloxRequest(d, nil, createCredentialFn)
	fmt.Println("err: ", err)
	require.NoError(t, err)
	actualInfoblox = actualInfoblox

	//require.Equal(t, expectedGridMaster, actualInfoblox.GridMaster)
	//require.Equal(t, expectedAnycast, actualInfoblox.AnyCast)
	//requireInfobloxInstancesEqual(t, expectedInstances, actualInfoblox.Instances)
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
	return alkira.InfobloxInstance{
		AnyCastEnabled:     anycast,
		ConfiguredMasterIp: prefix + "ConfiguredMasterIp",
		CredentialId:       prefix + "CredentialId",
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
	return alkira.InfobloxGridMaster{
		External:                 external,
		GridMasterCredentialId:   prefix + "GridMasterCredentialId",
		Ip:                       prefix + "Ip",
		Name:                     prefix + "Name",
		SharedSecretCredentialId: prefix + "SharedSecretCredentialId",
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

func sortInfobloxInstancesByCredentialId(ins []alkira.InfobloxInstance) {
	sort.Slice(ins, func(i, j int) bool { return ins[i].CredentialId < ins[j].CredentialId })
}

//This is required because the infoblox client has quite a few fields that are not currently in use
//in the backend. Here we are making comparisons only to the fields that the terraform user might
//input.
func requireInfobloxInstancesEqual(t *testing.T, expected, actual []alkira.InfobloxInstance) {
	require.Equal(t, len(expected), len(actual))

	//It doesn't really matter how we sort as long as they are in the same order. I chose
	//credentialId because it has the most potential for uniqueness i.e. the model field could
	//conceivably all be the same.
	sortInfobloxInstancesByCredentialId(expected)
	sortInfobloxInstancesByCredentialId(actual)

	for i, v := range actual {
		require.Equal(t, expected[i].AnyCastEnabled, v.AnyCastEnabled)
		require.Equal(t, expected[i].CredentialId, v.CredentialId)
		require.Equal(t, expected[i].HostName, v.HostName)
		require.Equal(t, expected[i].Model, v.Model)
		require.Equal(t, expected[i].Type, v.Type)
		require.Equal(t, expected[i].Version, v.Version)
	}
}
