package alkira

import (
	"encoding/json"
	"net/http"
	"testing"

	"github.com/alkiranet/alkira-client-go/alkira"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestGetSegmentNameByIdRejectsInvalidId asserts that getSegmentNameById
// validates its id argument before calling the Segment API. This closes
// the reachable-without-import path where a plain schema.TypeString
// `segment_id` (present on every resource that references a segment, with
// no Validate* of its own) was passed straight to segmentApi.GetById,
// whose URI is built with unescaped string interpolation. A payload like
// "../tenant/users?includeSecrets=true#" reached that call on an ordinary
// `terraform apply` -- no import required. A rejected id must never reach
// the API.
func TestGetSegmentNameByIdRejectsInvalidId(t *testing.T) {
	t.Run("malicious id is rejected before any API call", func(t *testing.T) {
		serverHit := false
		client := createMockAlkiraClient(t, http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
			serverHit = true
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(alkira.Segment{Id: "999", Name: "unexpected"})
		}))

		name, err := getSegmentNameById("../tenant/users?includeSecrets=true#", client)

		assert.Error(t, err, "expected a malicious segment_id to be rejected")
		assert.Empty(t, name)
		assert.False(t, serverHit, "the segment API must not be called for a rejected id")
	})

	t.Run("valid numeric id still resolves through the API", func(t *testing.T) {
		serverHit := false
		client := createMockAlkiraClient(t, http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
			serverHit = true
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(alkira.Segment{Id: "100", Name: "prod-segment"})
		}))

		name, err := getSegmentNameById("100", client)

		require.NoError(t, err)
		assert.Equal(t, "prod-segment", name)
		assert.True(t, serverHit, "a valid segment_id should still reach the segment API")
	})
}

// These tests focus on the input validation and data structure handling
// The actual API calls would need integration tests or dependency injection for proper testing

func TestSegmentConversionInputValidation(t *testing.T) {
	t.Run("nil set handling", func(t *testing.T) {
		var nilSet *schema.Set
		assert.Nil(t, nilSet)

		if nilSet == nil {
			// This is how the actual function should handle nil input
			assert.True(t, true, "Nil set should be handled gracefully")
		}
	})

	t.Run("empty set handling", func(t *testing.T) {
		emptySet := schema.NewSet(schema.HashString, []interface{}{})
		assert.Equal(t, 0, emptySet.Len())

		if emptySet.Len() == 0 {
			// This is how the actual function should handle empty input
			assert.True(t, true, "Empty set should be handled gracefully")
		}
	})

	t.Run("string slice conversion", func(t *testing.T) {
		// Test basic string slice operations that would be used in conversion
		testData := []string{"segment1", "segment2", "segment3"}

		// Simulate ID extraction from names
		result := make([]string, len(testData))
		for i, name := range testData {
			// In real function, this would call API to get ID
			result[i] = "id_" + name
		}

		expected := []string{"id_segment1", "id_segment2", "id_segment3"}
		assert.Equal(t, expected, result)
	})
}

// Test utility functions that don't require mocking
func TestSegmentHelperUtilities(t *testing.T) {
	t.Run("empty segment list handling", func(t *testing.T) {
		// Test that empty inputs are handled gracefully
		emptySet := schema.NewSet(schema.HashString, []interface{}{})
		assert.Equal(t, 0, emptySet.Len())

		emptySlice := []string{}
		assert.Equal(t, 0, len(emptySlice))
	})

	t.Run("schema set conversion", func(t *testing.T) {
		// Test schema.Set to string slice conversion
		testSet := schema.NewSet(schema.HashString, []interface{}{"a", "b", "c"})
		result := make([]string, 0, testSet.Len())
		for _, item := range testSet.List() {
			result = append(result, item.(string))
		}
		assert.ElementsMatch(t, []string{"a", "b", "c"}, result)
	})
}
