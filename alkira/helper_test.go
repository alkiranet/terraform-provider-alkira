package alkira

import (
	"testing"
)

func TestExpandGlobalProtectSegmentOptions(t *testing.T) {

	type obj struct {
		segment_name          string
		remote_user_zone_name string
		portal_fdqn_prefix    string
		service_group_name    string
	}

	o := &obj{
		segment_name:          "test-seg-name",
		remote_user_zone_name: "remote_user_zone_name",
		portal_fdqn_prefix:    "portal_fdqn_prefix",
		service_group_name:    "service_group_name",
	}

	o1 := &obj{
		segment_name:          "1test-seg-name",
		remote_user_zone_name: "1remote_user_zone_name",
		portal_fdqn_prefix:    "1portal_fdqn_prefix",
		service_group_name:    "1service_group_name",
	}

	r := resourceAlkiraServicePan()
	rd := r.TestResourceData()
	rd.Set("global_protect_segment_options", []interface{}{o, o1})
	//setFunc := schema.HashResource(r)
	//s := schema.NewSet(setFunc, []interface{}{o, o1})

	//fmt.Println(s)
	//rd.Set("globalProtectSegmentOptions", s)

	//fmt.Println(rd)

	if "" != "" {
		t.Fatal("Whooops")
	}
}
