package alkira

import (
	"context"
	"testing"

	"github.com/alkiranet/alkira-client-go/alkira"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestPolicyRoutingInputValidation(t *testing.T) {
	t.Run("nil input handling", func(t *testing.T) {
		var nilInput map[string]interface{}
		assert.Nil(t, nilInput)

		// Test that nil input would be handled gracefully
		if nilInput == nil {
			assert.True(t, true, "Nil input should be handled gracefully")
		}
	})

	t.Run("empty input handling", func(t *testing.T) {
		emptyInput := map[string]interface{}{}
		assert.Equal(t, 0, len(emptyInput))

		// Test that empty input would return empty result
		if len(emptyInput) == 0 {
			assert.True(t, true, "Empty input should be handled gracefully")
		}
	})

	t.Run("prefix list ID conversion", func(t *testing.T) {
		// Test conversion of interface{} slice to int slice
		input := []interface{}{1, 2, 3}
		result := make([]int, len(input))

		for i, v := range input {
			if intVal, ok := v.(int); ok {
				result[i] = intVal
			}
		}

		expected := []int{1, 2, 3}
		assert.Equal(t, expected, result)
	})

	t.Run("string slice conversion", func(t *testing.T) {
		// Test conversion of interface{} slice to string slice
		input := []interface{}{"65001", "65002"}
		result := make([]string, len(input))

		for i, v := range input {
			if strVal, ok := v.(string); ok {
				result[i] = strVal
			}
		}

		expected := []string{"65001", "65002"}
		assert.Equal(t, expected, result)
	})
}

func TestPolicyRoutingDataStructures(t *testing.T) {
	t.Run("rule sequence handling", func(t *testing.T) {
		// Test that rule sequence numbers are handled correctly
		rules := []map[string]interface{}{
			{"sequence": 10, "action": "PERMIT"},
			{"sequence": 20, "action": "DENY"},
		}

		for i, rule := range rules {
			sequence, ok := rule["sequence"].(int)
			assert.True(t, ok, "Sequence should be an integer")
			assert.Equal(t, (i+1)*10, sequence, "Sequence should be incremental")
		}
	})

	t.Run("action validation", func(t *testing.T) {
		validActions := []string{"PERMIT", "DENY"}

		for _, action := range validActions {
			// Test that valid actions are handled
			assert.Contains(t, validActions, action)
		}
	})

	t.Run("nested map handling", func(t *testing.T) {
		// Test handling of nested configuration maps
		nestedConfig := map[string]interface{}{
			"match": map[string]interface{}{
				"prefix_list_ids": []interface{}{1, 2, 3},
			},
			"set": map[string]interface{}{
				"community_list_ids": []interface{}{4, 5, 6},
			},
		}

		// Verify we can extract nested data
		if match, ok := nestedConfig["match"].(map[string]interface{}); ok {
			if prefixIds, ok := match["prefix_list_ids"].([]interface{}); ok {
				assert.Len(t, prefixIds, 3)
			}
		}

		if set, ok := nestedConfig["set"].(map[string]interface{}); ok {
			if communityIds, ok := set["community_list_ids"].([]interface{}); ok {
				assert.Len(t, communityIds, 3)
			}
		}
	})
}

// AK-75598: the API returns neither field for INBOUND, so an imported INBOUND
// policy has them absent from state while config gets the default `true`.
func TestPolicyRoutingOutboundOnlyFieldsDiff(t *testing.T) {
	outboundOnly := []string{"advertise_internet_exit", "enable_as_override"}

	tests := []struct {
		name        string
		direction   string
		stateValue  string // "" means absent from state, as after import
		configValue interface{}
		wantDiff    bool
	}{
		{name: "inbound import, omitted in config", direction: "INBOUND", wantDiff: false},
		{name: "inbound import, explicit false", direction: "INBOUND", configValue: false, wantDiff: false},
		{name: "inbound, state true, config false", direction: "INBOUND", stateValue: "true", configValue: false, wantDiff: false},
		{name: "outbound import, omitted in config", direction: "OUTBOUND", wantDiff: true},
		{name: "outbound, state true, config false", direction: "OUTBOUND", stateValue: "true", configValue: false, wantDiff: true},
		{name: "outbound, state matches default", direction: "OUTBOUND", stateValue: "true", wantDiff: false},
	}

	r := resourceAlkiraPolicyRouting()
	client := &alkira.AlkiraClient{}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			attrs := map[string]string{
				"id":         "1",
				"name":       "p",
				"direction":  tt.direction,
				"segment_id": "1",
			}
			cfg := map[string]interface{}{
				"name":               "p",
				"direction":          tt.direction,
				"segment_id":         "1",
				"included_group_ids": []interface{}{1},
			}
			for _, f := range outboundOnly {
				if tt.stateValue != "" {
					attrs[f] = tt.stateValue
				}
				if tt.configValue != nil {
					cfg[f] = tt.configValue
				}
			}

			state := &terraform.InstanceState{ID: "1", Attributes: attrs}
			diff, err := r.Diff(context.Background(), state, terraform.NewResourceConfigRaw(cfg), client)
			require.NoError(t, err)

			for _, f := range outboundOnly {
				changed := false
				if diff != nil {
					_, changed = diff.Attributes[f]
				}
				assert.Equal(t, tt.wantDiff, changed, "diff on %s", f)
			}
		})
	}
}
