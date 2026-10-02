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

// TestListUdrReadSetsRoutes covers AK-75539: route was never written to state
// by Read, so an imported list always planned its routes as additions.
func TestListUdrReadSetsRoutes(t *testing.T) {
	client := createMockAlkiraClient(t, func(w http.ResponseWriter, req *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(alkira.UdrList{
			Id:            json.Number("10"),
			Name:          "udr-list",
			CloudProvider: "AZURE",
			Udrs: []alkira.UdrListUdrs{
				{Prefix: "10.0.0.0/8", Description: "rfc1918", NextHopType: "INTERNET"},
				{Prefix: "192.168.0.0/16", NextHopType: "INTERNET"},
			},
		})
	})

	d := resourceAlkiraListUdr().TestResourceData()
	d.SetId("10")

	diags := resourceListUdrRead(context.Background(), d, client)
	require.Empty(t, diags)

	assert.ElementsMatch(t,
		[]alkira.UdrListUdrs{
			{Prefix: "10.0.0.0/8", Description: "rfc1918", NextHopType: "INTERNET"},
			{Prefix: "192.168.0.0/16", NextHopType: "INTERNET"},
		},
		expandListUdrRoutes(d.Get("route").(*schema.Set)),
		"routes read from the API must round-trip through state")
}

func TestFlattenListUdrRoutes(t *testing.T) {
	tests := []struct {
		name string
		in   []alkira.UdrListUdrs
		want []map[string]interface{}
	}{
		{name: "nil", in: nil, want: []map[string]interface{}{}},
		{
			// next hop fields are not in the schema, so they must be dropped.
			name: "drops next hop fields",
			in:   []alkira.UdrListUdrs{{Prefix: "10.0.0.0/8", Description: "d", NextHopType: "INTERNET", NextHopValue: "x"}},
			want: []map[string]interface{}{{"prefix": "10.0.0.0/8", "description": "d"}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, flattenListUdrRoutes(tt.in))
		})
	}
}
