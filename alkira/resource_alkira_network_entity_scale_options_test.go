// Package alkira - Copyright (C) 2023-2025 Alkira Inc. All Rights Reserved.
package alkira

import (
	"context"
	"testing"

	"github.com/alkiranet/alkira-client-go/alkira"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestWarnOnFailedScaleOptionsUpdate mirrors TestWarnOnFailedStateUpdate but
// exercises the scale-options one-off wrapper, which keys off the "state"
// field rather than "provision_state".
func TestWarnOnFailedScaleOptionsUpdate(t *testing.T) {
	testSchema := map[string]*schema.Schema{
		"name": {
			Type:     schema.TypeString,
			Optional: true,
		},
		"state": {
			Type:     schema.TypeString,
			Computed: true,
		},
	}

	// newResourceData builds a ResourceData with prior state and a diff,
	// mirroring what the SDK hands to UpdateContext during apply. When the
	// prior state is FAILED, the diff carries the FAILED->SUCCESS change
	// forced by the resource's CustomizeDiff.
	newResourceData := func(t *testing.T, state string, configChanged bool) *schema.ResourceData {
		instanceState := &terraform.InstanceState{
			ID: "1",
			Attributes: map[string]string{
				"name":  "old-name",
				"state": state,
			},
		}

		diffAttrs := map[string]*terraform.ResourceAttrDiff{}
		if configChanged {
			diffAttrs["name"] = &terraform.ResourceAttrDiff{
				Old: "old-name",
				New: "new-name",
			}
		}
		if state == "FAILED" {
			diffAttrs["state"] = &terraform.ResourceAttrDiff{
				Old: "FAILED",
				New: "SUCCESS",
			}
		}

		d, err := schema.InternalMap(testSchema).Data(instanceState,
			&terraform.InstanceDiff{Attributes: diffAttrs})
		require.NoError(t, err)
		return d
	}

	tests := []struct {
		name          string
		provision     bool
		state         string
		configChanged bool
		updateDiags   diag.Diagnostics
		expectWarning bool
	}{
		{
			name:          "warns on FAILED state with config changes",
			provision:     true,
			state:         "FAILED",
			configChanged: true,
			expectWarning: true,
		},
		{
			name:          "no warning on retry-only re-apply (no config changes)",
			provision:     true,
			state:         "FAILED",
			configChanged: false,
			expectWarning: false,
		},
		{
			name:          "no warning on healthy resource",
			provision:     true,
			state:         "SUCCESS",
			configChanged: true,
			expectWarning: false,
		},
		{
			name:          "no warning when provision mode is off",
			provision:     false,
			state:         "FAILED",
			configChanged: true,
			expectWarning: false,
		},
		{
			name:          "warning suppressed when update errors",
			provision:     true,
			state:         "FAILED",
			configChanged: true,
			updateDiags: diag.Diagnostics{{
				Severity: diag.Error,
				Summary:  "UPDATE FAILED",
			}},
			expectWarning: false,
		},
		{
			name:          "warning appended to update's own warnings",
			provision:     true,
			state:         "FAILED",
			configChanged: true,
			updateDiags: diag.Diagnostics{{
				Severity: diag.Warning,
				Summary:  "PROVISION (UPDATE) FAILED",
			}},
			expectWarning: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := newResourceData(t, tt.state, tt.configChanged)
			client := &alkira.AlkiraClient{Provision: tt.provision}

			updateCalled := false
			update := func(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
				updateCalled = true
				return tt.updateDiags
			}

			diags := warnOnFailedScaleOptionsUpdate(update)(context.Background(), d, client)

			// The wrapped update must always run - the warning never
			// blocks the request.
			assert.True(t, updateCalled, "wrapped update should always be invoked")

			warningCount := 0
			for _, diagnostic := range diags {
				if diagnostic.Summary == "CONFIGURATION CHANGES SKIPPED" {
					warningCount++
					assert.Equal(t, diag.Warning, diagnostic.Severity)
				}
			}

			if tt.expectWarning {
				assert.Equal(t, 1, warningCount, "expected the skip warning to be emitted")
				// The update's own diagnostics must be preserved.
				assert.Len(t, diags, len(tt.updateDiags)+1)
			} else {
				assert.Zero(t, warningCount, "expected no skip warning")
				assert.Len(t, diags, len(tt.updateDiags))
			}
		})
	}
}

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
