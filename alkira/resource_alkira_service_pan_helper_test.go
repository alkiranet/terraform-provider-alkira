package alkira

import (
	"net/http"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/stretchr/testify/require"
)

// TestExpandGlobalProtectSegmentOptionsPropagatesSegmentIdError pins that a
// segment name in segment_id fails the expand instead of being silently
// dropped, and that the failure happens before any API call.
func TestExpandGlobalProtectSegmentOptionsPropagatesSegmentIdError(t *testing.T) {
	client := createMockAlkiraClient(t, http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		t.Errorf("segment API must not be called for a rejected segment_id")
	}))

	in := schema.NewSet(schema.HashResource(&schema.Resource{
		Schema: map[string]*schema.Schema{
			"segment_id":            {Type: schema.TypeString},
			"remote_user_zone_name": {Type: schema.TypeString},
			"portal_fqdn_prefix":    {Type: schema.TypeString},
			"service_group_name":    {Type: schema.TypeString},
		},
	}), []interface{}{
		map[string]interface{}{
			"segment_id":            "ak74335-seg-a",
			"remote_user_zone_name": "zone-a",
			"portal_fqdn_prefix":    "abc",
			"service_group_name":    "group-a",
		},
	})

	_, err := expandGlobalProtectSegmentOptions(in, client)

	require.Error(t, err)
}

// TestExpandGlobalProtectSegmentOptionsInstancePropagatesSegmentIdError pins
// the same guarantee for the per-instance variant of the block.
func TestExpandGlobalProtectSegmentOptionsInstancePropagatesSegmentIdError(t *testing.T) {
	client := createMockAlkiraClient(t, http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		t.Errorf("segment API must not be called for a rejected segment_id")
	}))

	in := schema.NewSet(schema.HashResource(&schema.Resource{
		Schema: map[string]*schema.Schema{
			"segment_id":      {Type: schema.TypeString},
			"portal_enabled":  {Type: schema.TypeBool},
			"gateway_enabled": {Type: schema.TypeBool},
			"prefix_list_id":  {Type: schema.TypeInt},
		},
	}), []interface{}{
		map[string]interface{}{
			"segment_id":      "ak74335-seg-a",
			"portal_enabled":  true,
			"gateway_enabled": false,
			"prefix_list_id":  1,
		},
	})

	_, err := expandGlobalProtectSegmentOptionsInstance(in, client)

	require.Error(t, err)
}
