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

// TestSegmentResourceSegmentIdSchema pins segment_id to a numeric-ID-only
// value: getSegmentNameById hands it straight to GET /segments/<value>, and a
// segment name there is a 500 the client retries five times before failing.
func TestSegmentResourceSegmentIdSchema(t *testing.T) {
	field, ok := resourceAlkiraSegmentResource().Schema["segment_id"]
	require.True(t, ok, "segment_id must be present in schema")

	assert.Equal(t, schema.TypeString, field.Type)
	assert.True(t, field.Required)
	assert.False(t, field.Computed)
	require.NotNil(t, field.ValidateFunc, "segment_id must be validated")

	for _, v := range []string{"1145", "690"} {
		_, errs := field.ValidateFunc(v, "segment_id")
		assert.Emptyf(t, errs, "expected %q to be accepted", v)
	}

	for _, v := range []string{"ak74335-seg-a", "seg_1", "", "-1", "12ab", "0690", "007"} {
		_, errs := field.ValidateFunc(v, "segment_id")
		assert.Lenf(t, errs, 1, "expected %q to be rejected with exactly one error", v)
	}

	_, errs := field.ValidateFunc("ak74335-seg-a", "segment_id")
	require.Len(t, errs, 1)
	assert.Contains(t, errs[0].Error(), "segment_id")
	assert.Contains(t, errs[0].Error(), "alkira_segment.example.id")
}

// TestSegmentResourceSegmentIdRead pins segment_id to the segment's ID even
// though the API's segment-resource response names the segment instead.
func TestSegmentResourceSegmentIdRead(t *testing.T) {
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

		json.NewEncoder(w).Encode(alkira.SegmentResource{
			Id:      json.Number("1149"),
			Name:    "ak74335-res-a",
			Segment: name,
		})
	})

	r := resourceAlkiraSegmentResource()
	d := r.TestResourceData()
	d.SetId("1149")

	// A failed GET surfaces as a Warning here, which HasError() would not catch.
	diags := resourceSegmentResourceRead(context.Background(), d, client)
	require.Empty(t, diags, "resourceSegmentResourceRead should return no diagnostics")

	assert.Equal(t, "1145", d.Get("segment_id"),
		"Read should store the segment's ID even though the API returns its name")
}

// TestSegmentResourceReadFailedSegmentLookupIsFatal pins this resource's
// existing behavior on a failed segment lookup, which differs from
// alkira_segment_resource_share: segment_id is the only place the segment
// appears in state, so Read has nothing to fall back on and fails.
func TestSegmentResourceReadFailedSegmentLookupIsFatal(t *testing.T) {
	client := createMockAlkiraClient(t, func(w http.ResponseWriter, req *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		if req.URL.Query().Get("name") != "" {
			json.NewEncoder(w).Encode([]alkira.Segment{})
			return
		}

		json.NewEncoder(w).Encode(alkira.SegmentResource{
			Id:      json.Number("1149"),
			Name:    "ak74335-res-a",
			Segment: "ak74335-seg-a",
		})
	})

	r := resourceAlkiraSegmentResource()
	d := r.TestResourceData()
	d.SetId("1149")

	diags := resourceSegmentResourceRead(context.Background(), d, client)

	require.True(t, diags.HasError(),
		"an unresolvable segment must fail the refresh")
}
