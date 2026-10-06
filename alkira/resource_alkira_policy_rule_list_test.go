package alkira

import (
	"context"
	"encoding/json"
	"net/http"
	"testing"

	"github.com/alkiranet/alkira-client-go/alkira"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestPolicyRuleListReadSetsRules covers AK-75539: import seeds state only
// from Read, so rules missing here showed up as drift on the next plan.
func TestPolicyRuleListReadSetsRules(t *testing.T) {
	client := createMockAlkiraClient(t, func(w http.ResponseWriter, req *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(alkira.PolicyRuleList{
			Id:          json.Number("2002"),
			Name:        "traffic-rule-list-tf-import-test",
			Description: "Traffic rule list for TF import validation",
			Rules: []alkira.PolicyRuleListRule{
				{RuleId: 6182, Priority: 1},
				{RuleId: 6183, Priority: 2},
			},
		})
	})

	d := resourceAlkiraPolicyRuleList().TestResourceData()
	d.SetId("2002")

	diags := resourcePolicyRuleListRead(context.Background(), d, client)
	require.Empty(t, diags)

	assert.Equal(t, "traffic-rule-list-tf-import-test", d.Get("name"))
	assert.ElementsMatch(t,
		[]alkira.PolicyRuleListRule{
			{RuleId: 6182, Priority: 1},
			{RuleId: 6183, Priority: 2},
		},
		expandPolicyRuleListRules(d.Get("rules").(*schema.Set)),
		"rules read from the API must round-trip through state")
}

func TestFlattenPolicyRuleListRules(t *testing.T) {
	tests := []struct {
		name string
		in   []alkira.PolicyRuleListRule
		want []map[string]interface{}
	}{
		{name: "nil", in: nil, want: []map[string]interface{}{}},
		{
			name: "two rules",
			in:   []alkira.PolicyRuleListRule{{RuleId: 6182, Priority: 1}, {RuleId: 6183, Priority: 2}},
			want: []map[string]interface{}{
				{"priority": 1, "rule_id": 6182},
				{"priority": 2, "rule_id": 6183},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, flattenPolicyRuleListRules(tt.in))
		})
	}
}
