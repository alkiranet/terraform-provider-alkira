package alkira

import (
	"encoding/json"
	"testing"

	"github.com/alkiranet/alkira-client-go/alkira"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestExpandPrismaSDWANVrfMappings_NilInput(t *testing.T) {
	result := expandPrismaSDWANVrfMappings(nil)
	assert.Equal(t, []alkira.ConnectorPrismaSDWANVRFMapping{}, result)
}

func TestExpandPrismaSDWANVrfMappings_EmptyInput(t *testing.T) {
	emptySet := schema.NewSet(schema.HashString, []interface{}{})
	result := expandPrismaSDWANVrfMappings(emptySet)
	assert.Equal(t, []alkira.ConnectorPrismaSDWANVRFMapping{}, result)
}

func TestExpandPrismaSDWANVrfMappings_ValidInput(t *testing.T) {
	hashFunc := func(i interface{}) int {
		m := i.(map[string]interface{})
		return schema.HashString(m["vrf_name"].(string))
	}

	input := schema.NewSet(hashFunc, []interface{}{
		map[string]interface{}{
			"advertise_on_prem_routes": true,
			"advertise_default_route":  false,
			"gateway_bgp_asn":          65000,
			"segment_id":               42,
			"vrf_name":                 "test-vrf-1",
		},
	})

	result := expandPrismaSDWANVrfMappings(input)

	require.Len(t, result, 1)
	assert.Equal(t, true, result[0].AdvertiseOnPremRoutes)
	assert.Equal(t, true, result[0].DisableInternetExit) // advertise_default_route=false -> DisableInternetExit=true
	assert.Equal(t, 65000, result[0].GatewayBgpAsn)
	assert.Equal(t, 42, result[0].SegmentId)
	assert.Equal(t, "test-vrf-1", result[0].VrfName)
}

func TestExpandPrismaSDWANVrfMappings_AdvertiseDefaultRouteInversion(t *testing.T) {
	hashFunc := func(i interface{}) int {
		m := i.(map[string]interface{})
		return schema.HashString(m["vrf_name"].(string))
	}

	t.Run("advertise_default_route true sets DisableInternetExit false", func(t *testing.T) {
		input := schema.NewSet(hashFunc, []interface{}{
			map[string]interface{}{
				"advertise_on_prem_routes": false,
				"advertise_default_route":  true,
				"gateway_bgp_asn":          65001,
				"segment_id":               1,
				"vrf_name":                 "vrf-a",
			},
		})

		result := expandPrismaSDWANVrfMappings(input)
		require.Len(t, result, 1)
		assert.False(t, result[0].DisableInternetExit)
	})

	t.Run("advertise_default_route false sets DisableInternetExit true", func(t *testing.T) {
		input := schema.NewSet(hashFunc, []interface{}{
			map[string]interface{}{
				"advertise_on_prem_routes": false,
				"advertise_default_route":  false,
				"gateway_bgp_asn":          65002,
				"segment_id":               2,
				"vrf_name":                 "vrf-b",
			},
		})

		result := expandPrismaSDWANVrfMappings(input)
		require.Len(t, result, 1)
		assert.True(t, result[0].DisableInternetExit)
	})
}

func TestExpandPrismaSDWANVrfMappings_MultipleEntries(t *testing.T) {
	hashFunc := func(i interface{}) int {
		m := i.(map[string]interface{})
		return schema.HashString(m["vrf_name"].(string))
	}

	input := schema.NewSet(hashFunc, []interface{}{
		map[string]interface{}{
			"advertise_on_prem_routes": true,
			"advertise_default_route":  true,
			"gateway_bgp_asn":          65100,
			"segment_id":               10,
			"vrf_name":                 "vrf-first",
		},
		map[string]interface{}{
			"advertise_on_prem_routes": false,
			"advertise_default_route":  false,
			"gateway_bgp_asn":          65200,
			"segment_id":               20,
			"vrf_name":                 "vrf-second",
		},
	})

	result := expandPrismaSDWANVrfMappings(input)
	require.Len(t, result, 2)

	// Since sets are unordered, find each mapping by vrf_name
	vrfMap := make(map[string]alkira.ConnectorPrismaSDWANVRFMapping)
	for _, r := range result {
		vrfMap[r.VrfName] = r
	}

	first := vrfMap["vrf-first"]
	assert.True(t, first.AdvertiseOnPremRoutes)
	assert.Equal(t, 65100, first.GatewayBgpAsn)
	assert.Equal(t, 10, first.SegmentId)

	second := vrfMap["vrf-second"]
	assert.False(t, second.AdvertiseOnPremRoutes)
	assert.Equal(t, 65200, second.GatewayBgpAsn)
	assert.Equal(t, 20, second.SegmentId)
}

func TestExpandPrismaSDWANInstances_NilInput(t *testing.T) {
	result := expandPrismaSDWANInstances(nil)
	assert.Equal(t, []alkira.ConnectorPrismaSDWANInstance{}, result)
}

func TestExpandPrismaSDWANInstances_EmptyInput(t *testing.T) {
	result := expandPrismaSDWANInstances([]interface{}{})
	assert.Equal(t, []alkira.ConnectorPrismaSDWANInstance{}, result)
}

func TestExpandPrismaSDWANInstances_ValidInput(t *testing.T) {
	input := []interface{}{
		map[string]interface{}{
			"host_name":     "prisma-ion-1",
			"credential_id": "cred-abc",
			"ion_model":     "ion-3102v",
			"version":       "6.4.1",
			"id":            100,
		},
	}

	result := expandPrismaSDWANInstances(input)

	require.Len(t, result, 1)
	assert.Equal(t, "prisma-ion-1", result[0].HostName)
	assert.Equal(t, "cred-abc", result[0].CredentialId)
	assert.Equal(t, "ion-3102v", result[0].IonModel)
	assert.Equal(t, "6.4.1", result[0].Version)
	assert.Equal(t, json.Number("100"), result[0].Id)
}

func TestExpandPrismaSDWANInstances_MultipleInstances(t *testing.T) {
	input := []interface{}{
		map[string]interface{}{
			"host_name":     "prisma-ion-1",
			"credential_id": "cred-1",
			"ion_model":     "ion-3102v",
			"version":       "6.4.1",
			"id":            1,
		},
		map[string]interface{}{
			"host_name":     "prisma-ion-2",
			"credential_id": "cred-2",
			"ion_model":     "ion-3104v",
			"version":       "6.4.2",
			"id":            2,
		},
		map[string]interface{}{
			"host_name":     "prisma-ion-3",
			"credential_id": "cred-3",
			"ion_model":     "ion-3108v",
			"version":       "6.5.0",
			"id":            3,
		},
	}

	result := expandPrismaSDWANInstances(input)

	require.Len(t, result, 3)

	assert.Equal(t, "prisma-ion-1", result[0].HostName)
	assert.Equal(t, "cred-1", result[0].CredentialId)
	assert.Equal(t, "ion-3102v", result[0].IonModel)

	assert.Equal(t, "prisma-ion-2", result[1].HostName)
	assert.Equal(t, "cred-2", result[1].CredentialId)
	assert.Equal(t, "ion-3104v", result[1].IonModel)

	assert.Equal(t, "prisma-ion-3", result[2].HostName)
	assert.Equal(t, "cred-3", result[2].CredentialId)
	assert.Equal(t, "ion-3108v", result[2].IonModel)
}

func TestExpandPrismaSDWANInstances_ZeroId(t *testing.T) {
	input := []interface{}{
		map[string]interface{}{
			"host_name":     "new-instance",
			"credential_id": "cred-new",
			"ion_model":     "ion-3102v",
			"version":       "6.4.1",
			"id":            0,
		},
	}

	result := expandPrismaSDWANInstances(input)

	require.Len(t, result, 1)
	assert.Equal(t, json.Number("0"), result[0].Id)
	assert.Equal(t, "new-instance", result[0].HostName)
}

func TestSetPrismaSDWANInstances_MatchByHostName(t *testing.T) {
	// Build a ResourceData with the connector schema
	resourceSchema := resourceAlkiraConnectorPrismaSDWAN().Schema
	d := schema.TestResourceDataRaw(t, resourceSchema, map[string]interface{}{
		"name":            "test-connector",
		"cxp":             "US-EAST-2",
		"size":            "SMALL",
		"tunnel_protocol": "IPSEC",
		"instances": []interface{}{
			map[string]interface{}{
				"host_name":     "prisma-ion-1",
				"credential_id": "cred-abc",
				"ion_model":     "ion-3102v",
				"version":       "6.4.1",
				"id":            0,
			},
		},
	})

	// Simulate API response with assigned ID
	connector := &alkira.ConnectorPrismaSDWAN{
		Instances: []alkira.ConnectorPrismaSDWANInstance{
			{
				CredentialId: "cred-abc",
				HostName:     "prisma-ion-1",
				Id:           json.Number("42"),
				IonModel:     "ion-3102v",
				Version:      "6.4.1",
			},
		},
	}

	setPrismaSDWANInstances(d, connector)

	instances := d.Get("instances").([]interface{})
	require.Len(t, instances, 1)

	inst := instances[0].(map[string]interface{})
	assert.Equal(t, "prisma-ion-1", inst["host_name"])
	assert.Equal(t, "cred-abc", inst["credential_id"])
	assert.Equal(t, 42, inst["id"])
	assert.Equal(t, "ion-3102v", inst["ion_model"])
	assert.Equal(t, "6.4.1", inst["version"])
}

func TestSetPrismaSDWANInstances_MatchById(t *testing.T) {
	resourceSchema := resourceAlkiraConnectorPrismaSDWAN().Schema
	d := schema.TestResourceDataRaw(t, resourceSchema, map[string]interface{}{
		"name":            "test-connector",
		"cxp":             "US-EAST-2",
		"size":            "SMALL",
		"tunnel_protocol": "IPSEC",
		"instances": []interface{}{
			map[string]interface{}{
				"host_name":     "prisma-ion-1",
				"credential_id": "cred-abc",
				"ion_model":     "ion-3102v",
				"version":       "6.4.1",
				"id":            42,
			},
		},
	})

	connector := &alkira.ConnectorPrismaSDWAN{
		Instances: []alkira.ConnectorPrismaSDWANInstance{
			{
				CredentialId: "cred-abc-updated",
				HostName:     "prisma-ion-1-renamed",
				Id:           json.Number("42"),
				IonModel:     "ion-3102v",
				Version:      "6.5.0",
			},
		},
	}

	setPrismaSDWANInstances(d, connector)

	instances := d.Get("instances").([]interface{})
	require.Len(t, instances, 1)

	inst := instances[0].(map[string]interface{})
	assert.Equal(t, 42, inst["id"])
	assert.Equal(t, "prisma-ion-1-renamed", inst["host_name"])
	assert.Equal(t, "cred-abc-updated", inst["credential_id"])
	assert.Equal(t, "6.5.0", inst["version"])
}

func TestSetPrismaSDWANInstances_NewInstanceFromAPI(t *testing.T) {
	resourceSchema := resourceAlkiraConnectorPrismaSDWAN().Schema
	d := schema.TestResourceDataRaw(t, resourceSchema, map[string]interface{}{
		"name":            "test-connector",
		"cxp":             "US-EAST-2",
		"size":            "SMALL",
		"tunnel_protocol": "IPSEC",
		"instances": []interface{}{
			map[string]interface{}{
				"host_name":     "prisma-ion-1",
				"credential_id": "cred-1",
				"ion_model":     "ion-3102v",
				"version":       "6.4.1",
				"id":            10,
			},
		},
	})

	// API returns the existing instance plus a new one
	connector := &alkira.ConnectorPrismaSDWAN{
		Instances: []alkira.ConnectorPrismaSDWANInstance{
			{
				CredentialId: "cred-1",
				HostName:     "prisma-ion-1",
				Id:           json.Number("10"),
				IonModel:     "ion-3102v",
				Version:      "6.4.1",
			},
			{
				CredentialId: "cred-2",
				HostName:     "prisma-ion-new",
				Id:           json.Number("20"),
				IonModel:     "ion-3104v",
				Version:      "6.5.0",
			},
		},
	}

	setPrismaSDWANInstances(d, connector)

	instances := d.Get("instances").([]interface{})
	require.Len(t, instances, 2)

	// First instance should be the matched one
	inst0 := instances[0].(map[string]interface{})
	assert.Equal(t, 10, inst0["id"])
	assert.Equal(t, "prisma-ion-1", inst0["host_name"])

	// Second instance is the new one from the API
	inst1 := instances[1].(map[string]interface{})
	assert.Equal(t, 20, inst1["id"])
	assert.Equal(t, "prisma-ion-new", inst1["host_name"])
}
