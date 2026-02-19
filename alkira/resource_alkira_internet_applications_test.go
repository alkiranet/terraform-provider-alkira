package alkira

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestInternetApplicationsFieldNameMatch(t *testing.T) {
	t.Run("schema uses target field name", func(t *testing.T) {
		resourceSchema := resourceAlkiraInternetApplication().Schema

		// Verify "target" field exists in schema
		targetField, exists := resourceSchema["target"]
		assert.True(t, exists, "Schema must have 'target' field")
		assert.NotNil(t, targetField, "target field must not be nil")

		// Verify "targets" (plural - the bug) does NOT exist in schema
		_, wrongFieldExists := resourceSchema["targets"]
		assert.False(t, wrongFieldExists, "Schema should NOT have 'targets' field (bug was using plural)")
	})
}
