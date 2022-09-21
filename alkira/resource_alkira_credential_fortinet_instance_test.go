package alkira

import (
	"testing"

	"github.com/stretchr/testify/require"
)

// Test invalid file path
func TestFortinetCredentialInvalidPath(t *testing.T) {
	_, err := extractLicenseKey("", "/no/file/at/this/path")
	require.Error(t, err)
}

func TestFortinetLicenseKeyPriority(t *testing.T) {
	expected := "literal_file_contents"

	actual, err := extractLicenseKey(expected, "/some/file/path")
	require.NoError(t, err)
	require.Equal(t, actual, expected)
}

func TestFortinetLicenseKeyAndPathAreEmptyStr(t *testing.T) {
	_, err := extractLicenseKey("", "")
	require.Error(t, err)
}
