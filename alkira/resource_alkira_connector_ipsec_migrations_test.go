package alkira

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestConnectorRoutingOptionsSetUpgraderV0(t *testing.T) {
	for name, r := range map[string]*schema.Resource{
		"alkira_connector_ipsec":     resourceAlkiraConnectorIPSec(),
		"alkira_connector_ipsec_adv": resourceAlkiraConnectorIPSecAdv(),
	} {
		t.Run(name, func(t *testing.T) {
			assert.Equal(t, 1, r.SchemaVersion)
			require.Len(t, r.StateUpgraders, 1)

			upgrader := r.StateUpgraders[0]
			assert.Equal(t, 0, upgrader.Version)
			assert.True(t, upgrader.Type.AttributeType("routing_options").IsSetType(),
				"V0 routing_options must be a set")

			// Deriving V0 must not mutate the live schema.
			assert.Equal(t, schema.TypeList, r.Schema["routing_options"].Type)
			assert.Equal(t, 1, r.Schema["routing_options"].MaxItems)

			rawState := map[string]interface{}{
				"id": "17304",
				"routing_options": []interface{}{
					map[string]interface{}{
						"type":                 "DYNAMIC",
						"customer_gateway_asn": "65001",
						"local_asn":            "27205",
						"bgp_auth_key":         "secret",
					},
				},
			}
			out, err := upgrader.Upgrade(context.Background(), rawState, nil)
			require.NoError(t, err)
			assert.Equal(t, rawState, out)
		})
	}
}
