package alkira

import (
	"fmt"
	"testing"
)

func assertStrEquals(t *testing.T, str1, str2 string) {
	if str1 != str2 {
		t.Fatalf(fmt.Sprintf("failed asserting that %s is equal to %s", str1, str2))
	}

}

func assertTrue(t *testing.T, b bool, fieldName string) {
	if !b {
		t.Fatalf(fmt.Sprintf("failed asserting %s is true", fieldName))
	}
}
