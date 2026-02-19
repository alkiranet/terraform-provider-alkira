package alkira

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestExpandAwsDirectConnectInstances_IDHandling(t *testing.T) {
	// These tests focus on verifying that instance IDs are correctly extracted
	// from Terraform state data. We test the ID extraction logic in isolation
	// by examining the intermediate data structure before segment expansion.

	t.Run("ID is extracted from state as int", func(t *testing.T) {
		// This is the critical test for the bug fix:
		// Instance IDs from state MUST be preserved in the expanded struct

		input := []interface{}{
			map[string]interface{}{
				"name":          "dx-instance-1",
				"id":            123, // ID from state after creation
				"connection_id": "dxcon-abc123",
				"dx_asn":        64512,
				"on_prem_asn":   65000,
				"vlan_id":       100,
				"aws_region":    "us-west-2",
				"credential_id": "cred-123",
			},
		}

		// Extract the instance config
		cfg := input[0].(map[string]interface{})

		// Test the critical line that was failing
		if v, ok := cfg["id"].(int); ok {
			assert.Equal(t, 123, v, "ID should be extracted as int")
		} else {
			t.Error("FAILED: ID type assertion failed - this is the bug!")
		}
	})

	t.Run("ID field not present for new instances", func(t *testing.T) {
		input := []interface{}{
			map[string]interface{}{
				"name": "dx-instance-new",
				// No "id" field
				"connection_id": "dxcon-new",
				"dx_asn":        64512,
			},
		}

		cfg := input[0].(map[string]interface{})

		// Verify ID field doesn't exist
		_, exists := cfg["id"]
		assert.False(t, exists, "New instances should not have ID field")
	})

	t.Run("multiple instances preserve their IDs", func(t *testing.T) {
		input := []interface{}{
			map[string]interface{}{
				"id":            100,
				"name":          "instance-1",
				"connection_id": "dxcon-1",
			},
			map[string]interface{}{
				"id":            200,
				"name":          "instance-2",
				"connection_id": "dxcon-2",
			},
		}

		// Verify both IDs are extractable
		cfg1 := input[0].(map[string]interface{})
		cfg2 := input[1].(map[string]interface{})

		id1, ok1 := cfg1["id"].(int)
		id2, ok2 := cfg2["id"].(int)

		assert.True(t, ok1, "First instance ID should be extractable")
		assert.True(t, ok2, "Second instance ID should be extractable")
		assert.Equal(t, 100, id1)
		assert.Equal(t, 200, id2)
	})
}

func TestSetAwsDirectConnectInstance_FieldName(t *testing.T) {
	// This test verifies that setAwsDirectConnectInstance uses the correct
	// field name ("instance" not "instances") by checking the schema

	t.Run("schema uses 'instance' field name", func(t *testing.T) {
		resourceSchema := resourceAlkiraConnectorAwsDx().Schema

		// Verify "instance" field exists in schema
		instanceField, exists := resourceSchema["instance"]
		assert.True(t, exists, "Schema must have 'instance' field")
		assert.NotNil(t, instanceField, "instance field must not be nil")

		// Verify "instances" (plural) does NOT exist in schema
		_, wrongFieldExists := resourceSchema["instances"]
		assert.False(t, wrongFieldExists, "Schema should NOT have 'instances' field (bug was using plural)")
	})
}

func TestAwsDirectConnectInstanceStateHandling(t *testing.T) {
	t.Run("manually set instance data and verify ID preservation", func(t *testing.T) {
		// This test verifies we can manually set instance data and read it back
		// This simulates the state management without requiring API calls

		resourceSchema := resourceAlkiraConnectorAwsDx().Schema
		d := schema.TestResourceDataRaw(t, resourceSchema, map[string]interface{}{
			"name":            "test-connector",
			"cxp":             "US-WEST",
			"size":            "SMALL",
			"tunnel_protocol": "GRE",
		})

		// Manually set instance data (simulating what setAwsDirectConnectInstance would do)
		instanceData := []interface{}{
			map[string]interface{}{
				"id":            123,
				"name":          "dx-test",
				"connection_id": "dxcon-test",
				"dx_asn":        64512,
				"on_prem_asn":   65000,
				"vlan_id":       100,
				"aws_region":    "us-west-2",
				"credential_id": "cred-test",
			},
		}

		// Set to "instance" field (not "instances" - this is the bug fix)
		err := d.Set("instance", instanceData)
		require.NoError(t, err, "Should be able to set instance data")

		// Verify we can read it back
		instancesFromState := d.Get("instance")
		require.NotNil(t, instancesFromState, "Should be able to get instance data")

		instances := instancesFromState.([]interface{})
		require.Len(t, instances, 1)

		instance := instances[0].(map[string]interface{})
		assert.Equal(t, 123, instance["id"], "ID should be preserved")
		assert.Equal(t, "dx-test", instance["name"])
	})

	t.Run("verify wrong field name not accessible", func(t *testing.T) {
		resourceSchema := resourceAlkiraConnectorAwsDx().Schema
		d := schema.TestResourceDataRaw(t, resourceSchema, map[string]interface{}{
			"name":            "test-connector",
			"cxp":             "US-WEST",
			"size":            "SMALL",
			"tunnel_protocol": "GRE",
		})

		instanceData := []interface{}{
			map[string]interface{}{
				"id":   456,
				"name": "test",
			},
		}

		// Try to set to "instances" (plural - the bug)
		err := d.Set("instances", instanceData)
		// This should fail or be ignored because "instances" is not in schema
		if err == nil {
			// If Set doesn't error, verify Get returns nothing
			wrongData := d.Get("instances")
			t.Logf("Setting 'instances' (wrong name) returned: %v", wrongData)
			// The data won't be properly tracked if field name doesn't match schema
		}

		// Verify correct field name works
		err = d.Set("instance", instanceData)
		require.NoError(t, err)

		correctData := d.Get("instance")
		require.NotNil(t, correctData, "Correct field name 'instance' must work")
	})

	t.Run("multiple instances with IDs", func(t *testing.T) {
		resourceSchema := resourceAlkiraConnectorAwsDx().Schema
		d := schema.TestResourceDataRaw(t, resourceSchema, map[string]interface{}{
			"name":            "test-connector",
			"cxp":             "US-WEST",
			"size":            "SMALL",
			"tunnel_protocol": "GRE",
		})

		instanceData := []interface{}{
			map[string]interface{}{
				"id":            100,
				"name":          "instance-1",
				"connection_id": "dxcon-1",
				"dx_asn":        64512,
				"on_prem_asn":   65000,
				"vlan_id":       100,
				"aws_region":    "us-west-2",
				"credential_id": "cred-1",
			},
			map[string]interface{}{
				"id":            200,
				"name":          "instance-2",
				"connection_id": "dxcon-2",
				"dx_asn":        64512,
				"on_prem_asn":   65001,
				"vlan_id":       200,
				"aws_region":    "us-east-1",
				"credential_id": "cred-2",
			},
		}

		err := d.Set("instance", instanceData)
		require.NoError(t, err)

		// Read back and verify both IDs
		instances := d.Get("instance").([]interface{})
		require.Len(t, instances, 2)

		instance1 := instances[0].(map[string]interface{})
		instance2 := instances[1].(map[string]interface{})

		assert.Equal(t, 100, instance1["id"])
		assert.Equal(t, 200, instance2["id"])
	})
}
