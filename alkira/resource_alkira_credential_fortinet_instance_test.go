package alkira

import (
	"testing"

	"github.com/stretchr/testify/require"
)

// Test invalid file path
func TestFortinetCredentialInvalidPath(t *testing.T) {
	isPath := true
	_, err := setLicenseKey(isPath, "/no/file/at/this/path")
	require.Error(t, err)
}

// Test string is unchanged if not file path
func TestFortinetStringUnchangedPathFalse(t *testing.T) {
	isPath := false
	expected := "literal_file_contents"

	actual, err := setLicenseKey(isPath, expected)
	require.NoError(t, err)
	require.Equal(t, actual, expected)
}
