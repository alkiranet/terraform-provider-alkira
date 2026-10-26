package alkira

import (
	"reflect"
	"testing"

	"github.com/alkiranet/alkira-client-go/alkira"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func toInterfaceSlice(in []map[string]interface{}) []interface{} {
	out := make([]interface{}, len(in))
	for i, v := range in {
		out[i] = v
	}
	return out
}

// TestExpandFlattenInstances: a round-trip through flatten (API -> state) and back through
// instanceIdsByName/expand (state -> request) keeps each instance's id by name, carries
// on_ramp_region through (it is sent on every request, not read-only), and the computed
// bandwidth_in_mbps and lumen_connection_id survive being stored in state.
func TestExpandFlattenInstances(t *testing.T) {
	r := resourceAlkiraConnectorMultiCloudGateway()
	instances := []alkira.ConnectorMultiCloudGatewayInstance{
		{Id: 1, Name: "a", OnRampRegion: "us-east-1 - US East (N. Virginia) - Reston", BandwidthInMbps: 1000, LumenConnectionId: "lumen-1"},
		{Id: 2, Name: "b", OnRampRegion: "us-east-1 - US East (N. Virginia) - Ashburn 100G", BandwidthInMbps: 2000, LumenConnectionId: "lumen-2"},
	}
	flattened := flattenConnectorMultiCloudGatewayInstances(instances)

	d := schema.TestResourceDataRaw(t, r.Schema, nil)
	if err := d.Set("instance", flattened); err != nil {
		t.Fatal(err)
	}
	stored := d.Get("instance").([]interface{})

	for i, want := range instances {
		got := stored[i].(map[string]interface{})
		if got["bandwidth_in_mbps"] != want.BandwidthInMbps || got["on_ramp_region"] != want.OnRampRegion || got["lumen_connection_id"] != want.LumenConnectionId {
			t.Fatalf("fields lost in state round-trip: %#v", got)
		}
	}

	ids := instanceIdsByName(stored)
	out := expandConnectorMultiCloudGatewayInstances(stored, ids)
	want := []alkira.ConnectorMultiCloudGatewayInstance{
		{Name: "a", Id: 1, OnRampRegion: "us-east-1 - US East (N. Virginia) - Reston"},
		{Name: "b", Id: 2, OnRampRegion: "us-east-1 - US East (N. Virginia) - Ashburn 100G"},
	}
	if !reflect.DeepEqual(out, want) {
		t.Fatalf("got %#v, want %#v", out, want)
	}
}

// TestInstanceIdsMatchedByName: ids are keyed by name, not list position.
func TestInstanceIdsMatchedByName(t *testing.T) {
	stored := toInterfaceSlice([]map[string]interface{}{
		{"name": "a", "id": 1}, {"name": "b", "id": 2}, {"name": "c", "id": 3},
	})
	ids := instanceIdsByName(stored)
	if ids["a"] != 1 || ids["b"] != 2 || ids["c"] != 3 || len(ids) != 3 {
		t.Fatalf("unexpected id map: %v", ids)
	}
}

// TestRemovingMiddleBlockKeepsOtherIds: blocks [a, b, c] with stored ids [1, 2, 3]; removing b
// sends [a:1, c:3], and renaming c to d sends [a:1, b:2, d] with no id on d.
func TestRemovingMiddleBlockKeepsOtherIds(t *testing.T) {
	ids := instanceIdsByName(toInterfaceSlice([]map[string]interface{}{
		{"name": "a", "id": 1}, {"name": "b", "id": 2}, {"name": "c", "id": 3},
	}))

	t.Run("remove_middle_block", func(t *testing.T) {
		in := toInterfaceSlice([]map[string]interface{}{{"name": "a"}, {"name": "c"}})
		out := expandConnectorMultiCloudGatewayInstances(in, ids)
		want := []alkira.ConnectorMultiCloudGatewayInstance{{Name: "a", Id: 1}, {Name: "c", Id: 3}}
		if !reflect.DeepEqual(out, want) {
			t.Fatalf("got %#v, want %#v", out, want)
		}
	})

	t.Run("rename_block", func(t *testing.T) {
		in := toInterfaceSlice([]map[string]interface{}{{"name": "a"}, {"name": "b"}, {"name": "d"}})
		out := expandConnectorMultiCloudGatewayInstances(in, ids)
		want := []alkira.ConnectorMultiCloudGatewayInstance{{Name: "a", Id: 1}, {Name: "b", Id: 2}, {Name: "d", Id: 0}}
		if !reflect.DeepEqual(out, want) {
			t.Fatalf("got %#v, want %#v", out, want)
		}
	})
}

// TestOrderInstancesByPriorState: the API returning instances in a different order than the
// prior state/config must not show as a permanent diff. A reordered API response is put back
// into state order, and a new API-only instance is appended at the end.
func TestOrderInstancesByPriorState(t *testing.T) {
	prior := toInterfaceSlice([]map[string]interface{}{{"name": "a"}, {"name": "b"}})

	t.Run("api_reorders_known_instances", func(t *testing.T) {
		apiOrder := []map[string]interface{}{{"name": "b"}, {"name": "a"}}
		got := orderInstancesByPriorState(apiOrder, prior)
		if got[0]["name"] != "a" || got[1]["name"] != "b" {
			t.Fatalf("got %#v, want state order [a, b]", got)
		}
	})

	t.Run("new_api_only_instance_appended", func(t *testing.T) {
		apiOrder := []map[string]interface{}{{"name": "b"}, {"name": "c"}, {"name": "a"}}
		got := orderInstancesByPriorState(apiOrder, prior)
		want := []string{"a", "b", "c"}
		for i, w := range want {
			if got[i]["name"] != w {
				t.Fatalf("got %#v, want order %v", got, want)
			}
		}
	})
}

// TestUpdateSendsInstanceIds: an update sends each matched instance's stored id, carries the
// config's on_ramp_region (it is not read-only), and sends no other read-only field.
func TestUpdateSendsInstanceIds(t *testing.T) {
	ids := instanceIdsByName(toInterfaceSlice([]map[string]interface{}{
		{"name": "a", "id": 1, "bandwidth_in_mbps": 1000, "on_ramp_region": "us-east-1 - US East (N. Virginia) - Reston", "lumen_connection_id": "lumen-1"},
	}))
	in := toInterfaceSlice([]map[string]interface{}{{"name": "a", "on_ramp_region": "us-east-1 - US East (N. Virginia) - Ashburn 100G"}})
	out := expandConnectorMultiCloudGatewayInstances(in, ids)
	want := []alkira.ConnectorMultiCloudGatewayInstance{{Name: "a", Id: 1, OnRampRegion: "us-east-1 - US East (N. Virginia) - Ashburn 100G"}}
	if !reflect.DeepEqual(out, want) {
		t.Fatalf("got %#v, want %#v: only id and on_ramp_region may be sent for an existing instance", out, want)
	}
}

// TestValidateMultiCloudGatewayOnRampRegions: two instances sharing an on_ramp_region fail
// plan-time validation naming the duplicate; unique regions pass.
func TestValidateMultiCloudGatewayOnRampRegions(t *testing.T) {
	t.Run("duplicate_rejected", func(t *testing.T) {
		instances := toInterfaceSlice([]map[string]interface{}{
			{"name": "a", "on_ramp_region": "us-east-1 - US East (N. Virginia) - Reston"},
			{"name": "b", "on_ramp_region": "us-east-1 - US East (N. Virginia) - Reston"},
		})
		err := validateMultiCloudGatewayOnRampRegions(instances)
		if err == nil {
			t.Fatal("want an error for a duplicate on_ramp_region, got nil")
		}
	})

	t.Run("unique_accepted", func(t *testing.T) {
		instances := toInterfaceSlice([]map[string]interface{}{
			{"name": "a", "on_ramp_region": "us-east-1 - US East (N. Virginia) - Reston"},
			{"name": "b", "on_ramp_region": "us-east-1 - US East (N. Virginia) - Ashburn 100G"},
		})
		if err := validateMultiCloudGatewayOnRampRegions(instances); err != nil {
			t.Fatalf("want no error for unique on_ramp_region values, got %v", err)
		}
	})
}

// TestDuplicateInstanceNamesRejected: two instance blocks sharing a name fail plan-time
// validation naming the duplicate; unique names pass.
