package alkira

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestValidateSegmentId pins what reaches GET /segments/<value>.
func TestValidateSegmentId(t *testing.T) {
	accepted := []string{"1", "690", "1145", "2147483647", "9223372036854775807"}

	for _, id := range accepted {
		t.Run("accepts "+id, func(t *testing.T) {
			assert.NoError(t, validateSegmentId(id))
		})
	}

	rejected := []string{"", "0", "0690", "007", "ak74335-seg-a", "12ab", "-1", "1.5", " 1145", "1145 ", "seg_1", "9223372036854775808", "99999999999999999999"}

	for _, id := range rejected {
		t.Run("rejects "+id, func(t *testing.T) {
			assert.Error(t, validateSegmentId(id))
		})
	}

	t.Run("error text names the offending value and the fix", func(t *testing.T) {
		err := validateSegmentId("ak74335-seg-a")

		require.Error(t, err)
		assert.Contains(t, err.Error(), "ak74335-seg-a")
		assert.Contains(t, err.Error(), "alkira_segment.example.id")
	})

	t.Run("error text for an oversized value names the value", func(t *testing.T) {
		err := validateSegmentId("99999999999999999999")

		require.Error(t, err)
		assert.Contains(t, err.Error(), "99999999999999999999")
		assert.Contains(t, err.Error(), "too large")
	})
}

// TestGetSegmentNameByIdRejectsNameWithoutCallingApi pins that a segment name
// never reaches the API, which is what made this cost five retries.
func TestGetSegmentNameByIdRejectsNameWithoutCallingApi(t *testing.T) {
	cases := []string{"ak74335-seg-a", "0690", "99999999999999999999"}

	for _, id := range cases {
		t.Run(id, func(t *testing.T) {
			client := createMockAlkiraClient(t, http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
				t.Errorf("segment API must not be called for rejected id %q", id)
			}))

			name, err := getSegmentNameById(id, client)

			require.Error(t, err)
			assert.Empty(t, name)
		})
	}
}
