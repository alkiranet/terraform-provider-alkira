// Package alkira - Copyright (C) 2023-2025 Alkira Inc. All Rights Reserved.
package alkira

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// segmentScaleOption builds a segment_scale_options block map as the SDK would
// hand it to CustomizeDiff: additional_tunnels_per_node as int (0 when omitted)
// and additional_tunnel_options_per_node as a []interface{} of labelCount labels.
func segmentScaleOption(tunnelsPerNode, labelCount int) map[string]interface{} {
	labels := make([]interface{}, labelCount)
	for i := 0; i < labelCount; i++ {
		labels[i] = map[string]interface{}{"id": i + 1, "label": "label", "enabled": true}
	}
	return map[string]interface{}{
		"segment_id":                         1,
		"additional_tunnels_per_node":        tunnelsPerNode,
		"additional_tunnel_options_per_node": labels,
	}
}

func TestValidateAdditionalTunnelsPerNodeMatchesTunnelOptions(t *testing.T) {
	tests := []struct {
		name      string
		options   []interface{}
		wantError bool
	}{
		{
			name:      "matching count is valid",
			options:   []interface{}{segmentScaleOption(2, 2)},
			wantError: false,
		},
		{
			name:      "mismatched count is rejected",
			options:   []interface{}{segmentScaleOption(2, 1)},
			wantError: true,
		},
		{
			name:      "omitted tunnels per node with labels is rejected",
			options:   []interface{}{segmentScaleOption(0, 1)},
			wantError: true,
		},
		{
			name:      "no labels leaves tunnels per node untouched",
			options:   []interface{}{segmentScaleOption(3, 0)},
			wantError: false,
		},
		{
			name:      "empty options is valid",
			options:   []interface{}{},
			wantError: false,
		},
		{
			name:      "one bad segment among many is rejected",
			options:   []interface{}{segmentScaleOption(2, 2), segmentScaleOption(1, 3)},
			wantError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateAdditionalTunnelsPerNodeMatchesTunnelOptions(tt.options)
			if tt.wantError {
				require.Error(t, err)
				assert.Contains(t, err.Error(), "additional_tunnels_per_node")
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
