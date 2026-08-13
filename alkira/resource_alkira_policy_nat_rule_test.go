package alkira

import (
	"context"
	"net/http"
	"testing"

	"github.com/alkiranet/alkira-client-go/alkira"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// readNatRule drives the real Read function against a mock API returning the
// given rule and hands back the resulting resource data.
func readNatRule(t *testing.T, rule *alkira.NatPolicyRule) *schema.ResourceData {
	t.Helper()

	client := serveMockServer(t, rule, http.StatusOK)

	d := resourceAlkiraPolicyNatRule().TestResourceData()
	d.SetId("123")

	diags := resourcePolicyNatRuleRead(context.Background(), d, client)
	require.False(t, diags.HasError(), "Read should succeed: %v", diags)

	return d
}

// Read must persist "direction". The field is sent on create and update but
// was never read back, so "terraform import alkira_policy_nat_rule" left it
// null in state and the next plan showed a spurious
// "+ direction = INBOUND".
func TestPolicyNatRuleReadSetsDirection(t *testing.T) {
	newRule := func(direction string) *alkira.NatPolicyRule {
		return &alkira.NatPolicyRule{
			Name:        "nat-rule-1",
			Description: "nat rule under test",
			Enabled:     true,
			Category:    "DEFAULT",
			Direction:   direction,
			Match: alkira.NatRuleMatch{
				SourcePrefixes: []string{"10.0.0.0/8"},
				DestPrefixes:   []string{"192.168.0.0/16"},
				Protocol:       "any",
			},
		}
	}

	t.Run("INBOUND is written to state", func(t *testing.T) {
		d := readNatRule(t, newRule("INBOUND"))

		assert.Equal(t, "INBOUND", d.Get("direction"),
			"Read must persist direction, otherwise import leaves it null")
	})

	t.Run("OUTBOUND is written to state", func(t *testing.T) {
		d := readNatRule(t, newRule("OUTBOUND"))

		assert.Equal(t, "OUTBOUND", d.Get("direction"))
	})

	t.Run("direction omitted by the API stays empty", func(t *testing.T) {
		d := readNatRule(t, newRule(""))

		assert.Equal(t, "", d.Get("direction"),
			"an API response without direction must not invent a value")
	})

	t.Run("the other fields Read owns are unaffected", func(t *testing.T) {
		rule := newRule("INBOUND")

		d := readNatRule(t, rule)

		assert.Equal(t, rule.Name, d.Get("name"))
		assert.Equal(t, rule.Description, d.Get("description"))
		assert.Equal(t, rule.Enabled, d.Get("enabled"))
		assert.Equal(t, rule.Category, d.Get("category"))

		match := d.Get("match").(*schema.Set).List()
		require.Len(t, match, 1)
		matchMap := match[0].(map[string]interface{})
		assert.Equal(t, "any", matchMap["protocol"])
	})
}
