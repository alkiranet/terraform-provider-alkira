package alkira

import (
	"context"
	"testing"

	ctymsgpack "github.com/hashicorp/go-cty/cty/msgpack"
	"github.com/hashicorp/terraform-plugin-go/tfprotov5"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestSegmentOptionsGroupsStateFromList asserts that state written while
// segment_options.groups was a list still loads now that it is a set. JSON
// state stores both as an array, so no state migration is needed.
func TestSegmentOptionsGroupsStateFromList(t *testing.T) {
	resources := map[string]*schema.Resource{
		"alkira_service_pan":        resourceAlkiraServicePan(),
		"alkira_service_fortinet":   resourceAlkiraServiceFortinet(),
		"alkira_service_checkpoint": resourceAlkiraCheckpoint(),
		"alkira_service_cisco_ftdv": resourceAlkiraServiceCiscoFTDv(),
	}

	oldState := []byte(`{"id":"1","segment_options":[{"segment_id":"100","zone_name":"zone1","groups":["tgw","egress","aws"]}]}`)

	for typeName, res := range resources {
		t.Run(typeName, func(t *testing.T) {
			server := schema.NewGRPCProviderServer(&schema.Provider{
				ResourcesMap: map[string]*schema.Resource{typeName: res},
			})

			resp, err := server.UpgradeResourceState(context.Background(), &tfprotov5.UpgradeResourceStateRequest{
				TypeName: typeName,
				Version:  int64(res.SchemaVersion),
				RawState: &tfprotov5.RawState{JSON: oldState},
			})
			require.NoError(t, err)
			require.Empty(t, resp.Diagnostics)

			state, err := ctymsgpack.Unmarshal(resp.UpgradedState.MsgPack, res.CoreConfigSchema().ImpliedType())
			require.NoError(t, err)

			options := state.GetAttr("segment_options").AsValueSlice()
			require.Len(t, options, 1)

			var groups []string
			for _, g := range options[0].GetAttr("groups").AsValueSlice() {
				groups = append(groups, g.AsString())
			}
			assert.ElementsMatch(t, []string{"tgw", "egress", "aws"}, groups)
		})
	}
}
