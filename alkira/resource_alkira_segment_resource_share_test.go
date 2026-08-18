package alkira

import (
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// AK-73557: the accepted values for traffic_direction, traffic_from_end and the
// service_ids sentinels were not discoverable from the generated documentation. The
// docs are rendered from these Description strings, so assert the values are named
// there and that the validators agree with what is documented.

func TestSegmentResourceShareTrafficDirectionSchema(t *testing.T) {
	field, ok := resourceAlkiraSegmentResourceShare().Schema["traffic_direction"]
	require.True(t, ok, "traffic_direction must be present in schema")

	assert.Equal(t, schema.TypeString, field.Type)
	assert.True(t, field.Optional)
	assert.Equal(t, "BIDIRECTIONAL", field.Default)
	require.NotNil(t, field.ValidateFunc, "traffic_direction must be validated")

	for _, v := range []string{"BIDIRECTIONAL", "UNIDIRECTIONAL"} {
		assert.Containsf(t, field.Description, v, "traffic_direction docs must name %q", v)

		_, errs := field.ValidateFunc(v, "traffic_direction")
		assert.Emptyf(t, errs, "expected %q to be accepted", v)
	}

	for _, v := range []string{"bidirectional", "ANY", "A", ""} {
		_, errs := field.ValidateFunc(v, "traffic_direction")
		assert.NotEmptyf(t, errs, "expected %q to be rejected", v)
	}
}

func TestSegmentResourceShareTrafficFromEndSchema(t *testing.T) {
	field, ok := resourceAlkiraSegmentResourceShare().Schema["traffic_from_end"]
	require.True(t, ok, "traffic_from_end must be present in schema")

	assert.Equal(t, schema.TypeString, field.Type)
	assert.True(t, field.Optional)
	assert.Nil(t, field.Default, "traffic_from_end must stay unset by default so it can be omitted for BIDIRECTIONAL")
	require.NotNil(t, field.ValidateFunc, "traffic_from_end must be validated")

	// The server-side enum is SegmentResourceShare.End, which is exactly A or B.
	for _, v := range []string{"A", "B"} {
		assert.Containsf(t, field.Description, "`"+v+"`", "traffic_from_end docs must name %q", v)

		_, errs := field.ValidateFunc(v, "traffic_from_end")
		assert.Emptyf(t, errs, "expected %q to be accepted", v)
	}

	for _, v := range []string{"a", "b", "C", "END_A", ""} {
		_, errs := field.ValidateFunc(v, "traffic_from_end")
		assert.NotEmptyf(t, errs, "expected %q to be rejected", v)
	}

	assert.Contains(t, field.Description, "UNIDIRECTIONAL",
		"traffic_from_end docs must say when the field applies")
}

func TestSegmentResourceShareServiceIdsSchema(t *testing.T) {
	field, ok := resourceAlkiraSegmentResourceShare().Schema["service_ids"]
	require.True(t, ok, "service_ids must be present in schema")

	assert.Equal(t, schema.TypeList, field.Type)
	assert.True(t, field.Required)

	elem, ok := field.Elem.(*schema.Schema)
	require.True(t, ok, "service_ids Elem must be a *schema.Schema")
	assert.Equal(t, schema.TypeInt, elem.Type)

	// Both sentinels carry meaning that cannot be guessed from the type alone:
	// 0 is no service, -1 is any service (tenant-provisioning-service
	// SegmentResourceShare.serviceList, NO_FIREWALL / ANY_FIREWALL).
	for _, sentinel := range []string{"`0`", "`-1`"} {
		assert.Containsf(t, field.Description, sentinel,
			"service_ids docs must document the %s sentinel", sentinel)
	}
	assert.True(t,
		strings.Contains(field.Description, "no service") && strings.Contains(field.Description, "any service"),
		"service_ids docs must explain what the sentinels select, got: %s", field.Description)
}
