package alkira

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"
	"time"

	"github.com/alkiranet/alkira-client-go/alkira"
	"github.com/hashicorp/go-retryablehttp"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// Common mock server creator for Alkira resources
func createMockAlkiraClient(t *testing.T, handler http.HandlerFunc) *alkira.AlkiraClient {
	server := httptest.NewServer(handler)
	t.Cleanup(server.Close)

	retryClient := retryablehttp.NewClient()
	retryClient.HTTPClient.Timeout = time.Duration(1) * time.Second

	return &alkira.AlkiraClient{
		URI:             server.URL,
		TenantNetworkId: "0",
		Client:          retryClient,
		Provision:       false,
	}
}

// Common HTTP handler for mock responses
func createMockHandler(responseData interface{}, statusCode int) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(statusCode)

		if responseData != nil && statusCode != http.StatusNoContent {
			// This will work for any JSON-serializable type
			if jsonData, ok := responseData.([]byte); ok {
				w.Write(jsonData)
			} else {
				// For struct types, they will be JSON-encoded by the specific mock servers
				// This is a generic placeholder
			}
		}
	}
}

// Common validation function for names
func validateResourceName(v interface{}, k string) (warnings []string, errors []error) {
	value := v.(string)
	if value == "" {
		errors = append(errors, fmt.Errorf("%q cannot be empty", k))
	}
	return warnings, errors
}

// Common validation function for IDs
func validateResourceId(id string) bool {
	if id == "" {
		return false
	}
	_, err := strconv.Atoi(id)
	return err == nil
}

// Common string utility functions
func containsString(s, substr string) bool {
	return len(substr) <= len(s) && (len(substr) == 0 || indexOf(s, substr) >= 0)
}

func indexOf(s, substr string) int {
	if len(substr) == 0 {
		return 0
	}
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return i
		}
	}
	return -1
}

// Common helper function to safely get string from ResourceData
func getStringFromResourceData(d *schema.ResourceData, key string) string {
	if val := d.Get(key); val != nil {
		return val.(string)
	}
	return ""
}

// Common helper function to safely get int from ResourceData
func getIntFromResourceData(d *schema.ResourceData, key string) int {
	if val := d.Get(key); val != nil {
		return val.(int)
	}
	return 0
}

// Common helper function to safely get bool from ResourceData
func getBoolFromResourceData(d *schema.ResourceData, key string) bool {
	if val := d.Get(key); val != nil {
		return val.(bool)
	}
	return false
}

// Common helper function to safely get string slice from ResourceData
func getStringSliceFromResourceData(d *schema.ResourceData, key string) []string {
	if val := d.Get(key); val != nil {
		if list, ok := val.([]interface{}); ok {
			result := make([]string, len(list))
			for i, v := range list {
				if str, ok := v.(string); ok {
					result[i] = str
				}
			}
			return result
		}
	}
	return []string{}
}

// Common helper function to safely get int slice from ResourceData
func getIntSliceFromResourceData(d *schema.ResourceData, key string) []int {
	if val := d.Get(key); val != nil {
		if list, ok := val.([]interface{}); ok {
			result := make([]int, len(list))
			for i, v := range list {
				if num, ok := v.(int); ok {
					result[i] = num
				}
			}
			return result
		}
	}
	return []int{}
}

// Common test table structure for validation tests
type ValidationTestCase struct {
	Name      string
	Input     interface{}
	ExpectErr bool
	ErrCount  int
}

// Common test table structure for ID validation tests
type IdValidationTestCase struct {
	Name  string
	Id    string
	Valid bool
}

// Common ID validation test cases
func GetCommonIdValidationTestCases() []IdValidationTestCase {
	return []IdValidationTestCase{
		{
			Name:  "valid numeric ID",
			Id:    "123",
			Valid: true,
		},
		{
			Name:  "valid large numeric ID",
			Id:    "999999999999",
			Valid: true,
		},
		{
			Name:  "invalid empty ID",
			Id:    "",
			Valid: false,
		},
		{
			Name:  "invalid non-numeric ID",
			Id:    "abc",
			Valid: false,
		},
	}
}

// Common name validation test cases
func GetCommonNameValidationTestCases() []ValidationTestCase {
	return []ValidationTestCase{
		{
			Name:      "valid name",
			Input:     "test-resource",
			ExpectErr: false,
			ErrCount:  0,
		},
		{
			Name:      "valid name with numbers",
			Input:     "resource-123",
			ExpectErr: false,
			ErrCount:  0,
		},
		{
			Name:      "valid name with underscores",
			Input:     "test_resource_name",
			ExpectErr: false,
			ErrCount:  0,
		},
		{
			Name:      "empty name",
			Input:     "",
			ExpectErr: true,
			ErrCount:  1,
		},
		{
			Name:      "name with spaces",
			Input:     "test resource",
			ExpectErr: false,
			ErrCount:  0,
		},
	}
}
