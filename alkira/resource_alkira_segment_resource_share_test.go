package alkira

import (
	"context"
	"encoding/json"
	"net/http"
	"strings"
	"testing"

	"github.com/alkiranet/alkira-client-go/alkira"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
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

// TestSegmentResourceShareDesignatedSegmentIdSchema pins designated_segment_id
// to a numeric-ID-only value: getSegmentNameById hands it straight to
// GET /segments/<value>, and a segment name there is a 500 the client retries
// five times before failing.
func TestSegmentResourceShareDesignatedSegmentIdSchema(t *testing.T) {
	field, ok := resourceAlkiraSegmentResourceShare().Schema["designated_segment_id"]
	require.True(t, ok, "designated_segment_id must be present in schema")

	assert.Equal(t, schema.TypeString, field.Type)
	assert.True(t, field.Required)
	assert.False(t, field.Computed)
	require.NotNil(t, field.ValidateFunc, "designated_segment_id must be validated")

	for _, v := range []string{"1145", "690"} {
		_, errs := field.ValidateFunc(v, "designated_segment_id")
		assert.Emptyf(t, errs, "expected %q to be accepted", v)
	}

	for _, v := range []string{"ak74335-seg-a", "seg_1", "", "-1", "12ab", "0690", "007"} {
		_, errs := field.ValidateFunc(v, "designated_segment_id")
		assert.Lenf(t, errs, 1, "expected %q to be rejected with exactly one error", v)
	}

	_, errs := field.ValidateFunc("ak74335-seg-a", "designated_segment_id")
	require.Len(t, errs, 1)
	assert.Contains(t, errs[0].Error(), "designated_segment_id")
	assert.Contains(t, errs[0].Error(), "alkira_segment.example.id")
}

// TestSegmentResourceShareDesignatedSegmentIdRead pins designated_segment_id
// to the segment's ID even though the API's share response names the segment
// instead, and pins that Read still sets the attributes that follow that
// conversion in the function.
func TestSegmentResourceShareDesignatedSegmentIdRead(t *testing.T) {
	name := "ak74335-seg-a"

	client := createMockAlkiraClient(t, func(w http.ResponseWriter, req *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		// The segment lookup is a get-by-name, which returns a list.
		if req.URL.Query().Get("name") != "" {
			require.Contains(t, req.URL.Path, "/segments",
				"segment ID must be looked up against the segment endpoint")
			json.NewEncoder(w).Encode([]alkira.Segment{{
				Id:   json.Number("1145"),
				Name: name,
			}})
			return
		}

		json.NewEncoder(w).Encode(alkira.SegmentResourceShare{
			Id:                json.Number("112"),
			Name:              "ak74335-share",
			DesignatedSegment: name,
			ServiceList:       []int{0},
			EndAResources:     []int{1149},
			EndBResources:     []int{1150},
			Direction:         "BIDIRECTIONAL",
		})
	})

	r := resourceAlkiraSegmentResourceShare()
	d := r.TestResourceData()
	d.SetId("112")

	// A failed GET surfaces as a Warning here, which HasError() would not catch.
	diags := resourceSegmentResourceShareRead(context.Background(), d, client)
	require.Empty(t, diags, "resourceSegmentResourceShareRead should return no diagnostics")

	assert.Equal(t, "1145", d.Get("designated_segment_id"),
		"Read should store the segment's ID even though the API returns its name")
	assert.NotEmpty(t, d.Get("end_a_segment_resource_ids"),
		"Read should still set attributes that follow the segment conversion")
}

// TestSegmentResourceShareReadFailedSegmentLookupKeepsRefreshing pins a failed
// segment lookup to a warning that skips only designated_segment_id. Returning
// there would stop the attributes after it from refreshing while Terraform
// still counted the refresh a success, and it would abort the refresh that
// terraform destroy runs first whenever the segment is already marked for
// deletion, which the share's own GetById tolerates and the segment
// get-by-name does not.
func TestSegmentResourceShareReadFailedSegmentLookupKeepsRefreshing(t *testing.T) {
	client := createMockAlkiraClient(t, func(w http.ResponseWriter, req *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		// A get-by-name that matches nothing returns an empty list, which the
		// client reports as a failure to resolve the name.
		if req.URL.Query().Get("name") != "" {
			json.NewEncoder(w).Encode([]alkira.Segment{})
			return
		}

		json.NewEncoder(w).Encode(alkira.SegmentResourceShare{
			Id:                json.Number("112"),
			Name:              "ak74335-share",
			DesignatedSegment: "ak74335-seg-a",
			ServiceList:       []int{0},
			EndAResources:     []int{1149},
			EndBResources:     []int{1150},
			EndARouteLimit:    100,
			EndBRouteLimit:    100,
			Direction:         "BIDIRECTIONAL",
		})
	})

	r := resourceAlkiraSegmentResourceShare()
	d := r.TestResourceData()
	d.SetId("112")

	diags := resourceSegmentResourceShareRead(context.Background(), d, client)

	require.Len(t, diags, 1, "an unresolvable designated segment should warn once")
	assert.Equal(t, diag.Warning, diags[0].Severity,
		"the warning must not abort the refresh that terraform destroy runs first")
	assert.Contains(t, diags[0].Detail, "ak74335-seg-a",
		"the warning should name the segment it could not resolve")

	assert.Equal(t, []interface{}{1149}, d.Get("end_a_segment_resource_ids"),
		"attributes after the segment lookup must still refresh")
	assert.Equal(t, "BIDIRECTIONAL", d.Get("traffic_direction"))
	assert.Equal(t, 100, d.Get("end_b_route_limit"))
}
