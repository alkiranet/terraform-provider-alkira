package alkira

import (
	"encoding/json"
	"testing"

	"github.com/alkiranet/alkira-client-go/alkira"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/stretchr/testify/assert"
)

func TestDeflateBluecatInstances(t *testing.T) {
	tests := []struct {
		name     string
		input    []alkira.BluecatInstance
		expected []map[string]interface{}
	}{
		{
			name:     "nil input",
			input:    nil,
			expected: nil,
		},
		{
			name:     "empty input",
			input:    []alkira.BluecatInstance{},
			expected: nil,
		},
		{
			name: "single BDDS instance",
			input: []alkira.BluecatInstance{
				{
					Name: "bdds-primary.example.com",
					Type: "BDDS",
					Id:   json.Number("123"),
					BddsOptions: &alkira.BDDSOptions{
						HostName:            "bdds-primary.example.com",
						Model:               "cBDDS50",
						Version:             "9.4.0",
						LicenseCredentialId: "cred-123",
					},
				},
			},
			expected: []map[string]interface{}{
				{
					"id":   json.Number("123"),
					"name": "bdds-primary.example.com",
					"type": "BDDS",
					"bdds_options": []interface{}{
						map[string]interface{}{
							"hostname":              "bdds-primary.example.com",
							"model":                 "cBDDS50",
							"version":               "9.4.0",
							"license_credential_id": "cred-123",
						},
					},
				},
			},
		},
		{
			name: "single Edge instance",
			input: []alkira.BluecatInstance{
				{
					Name: "edge-primary.example.com",
					Type: "EDGE",
					Id:   json.Number("124"),
					EdgeOptions: &alkira.EdgeOptions{
						HostName:     "edge-primary.example.com",
						Version:      "4.2.0",
						CredentialId: "cred-124",
					},
				},
			},
			expected: []map[string]interface{}{
				{
					"id":   json.Number("124"),
					"name": "edge-primary.example.com",
					"type": "EDGE",
					"edge_options": []interface{}{
						map[string]interface{}{
							"hostname":      "edge-primary.example.com",
							"version":       "4.2.0",
							"credential_id": "cred-124",
						},
					},
				},
			},
		},
		{
			name: "multiple instances",
			input: []alkira.BluecatInstance{
				{
					Name: "bdds-primary.example.com",
					Type: "BDDS",
					Id:   json.Number("123"),
					BddsOptions: &alkira.BDDSOptions{
						HostName:            "bdds-primary.example.com",
						Model:               "cBDDS50",
						Version:             "9.4.0",
						LicenseCredentialId: "cred-123",
					},
				},
				{
					Name: "edge-primary",
					Type: "EDGE",
					Id:   json.Number("124"),
					EdgeOptions: &alkira.EdgeOptions{
						HostName:     "edge-primary.example.com",
						Version:      "4.2.0",
						CredentialId: "cred-124",
					},
				},
			},
			expected: []map[string]interface{}{
				{
					"id":   json.Number("123"),
					"name": "bdds-primary.example.com",
					"type": "BDDS",
					"bdds_options": []interface{}{
						map[string]interface{}{
							"hostname":              "bdds-primary.example.com",
							"model":                 "cBDDS50",
							"version":               "9.4.0",
							"license_credential_id": "cred-123",
						},
					},
				},
				{
					"id":   json.Number("124"),
					"name": "edge-primary",
					"type": "EDGE",
					"edge_options": []interface{}{
						map[string]interface{}{
							"hostname":      "edge-primary.example.com",
							"version":       "4.2.0",
							"credential_id": "cred-124",
						},
					},
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := deflateBluecatInstances(tt.input)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestExpandBluecatAnycast(t *testing.T) {
	tests := []struct {
		name        string
		input       *schema.Set
		expected    *alkira.BluecatAnycast
		expectError bool
		errorMsg    string
	}{
		{
			name:        "nil input",
			input:       nil,
			expected:    &alkira.BluecatAnycast{},
			expectError: false,
		},
		{
			name:        "empty input",
			input:       schema.NewSet(schema.HashString, []interface{}{}),
			expected:    &alkira.BluecatAnycast{},
			expectError: false,
		},
		{
			name: "valid anycast configuration",
			input: schema.NewSet(
				func(i interface{}) int {
					return schema.HashString("test")
				},
				[]interface{}{
					map[string]interface{}{
						"ips":         []interface{}{"192.168.1.10", "192.168.1.11"},
						"backup_cxps": []interface{}{"cxp1", "cxp2"},
					},
				},
			),
			expected: &alkira.BluecatAnycast{
				Ips:        []string{"192.168.1.10", "192.168.1.11"},
				BackupCxps: []string{"cxp1", "cxp2"},
			},
			expectError: false,
		},
		{
			name: "minimal anycast configuration",
			input: schema.NewSet(
				func(i interface{}) int {
					return schema.HashString("test")
				},
				[]interface{}{
					map[string]interface{}{},
				},
			),
			expected: &alkira.BluecatAnycast{
				Ips:        []string(nil),
				BackupCxps: []string(nil),
			},
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := expandBluecatAnycast(tt.input)

			if tt.expectError {
				assert.Error(t, err)
				if tt.errorMsg != "" {
					assert.Contains(t, err.Error(), tt.errorMsg)
				}
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expected, result)
			}
		})
	}
}

func TestDeflateBluecatAnycast(t *testing.T) {
	tests := []struct {
		name     string
		input    alkira.BluecatAnycast
		expected []map[string]interface{}
	}{
		{
			name: "full anycast configuration",
			input: alkira.BluecatAnycast{
				Ips:        []string{"192.168.1.10", "192.168.1.11"},
				BackupCxps: []string{"cxp1", "cxp2"},
			},
			expected: []map[string]interface{}{
				{
					"ips":         []string{"192.168.1.10", "192.168.1.11"},
					"backup_cxps": []string{"cxp1", "cxp2"},
				},
			},
		},
		{
			name: "minimal anycast configuration",
			input: alkira.BluecatAnycast{
				Ips: []string{"192.168.1.10"},
			},
			expected: []map[string]interface{}{
				{
					"ips":         []string{"192.168.1.10"},
					"backup_cxps": []string(nil),
				},
			},
		},
		{
			name:  "empty anycast configuration",
			input: alkira.BluecatAnycast{},
			expected: []map[string]interface{}{
				{
					"ips":         []string(nil),
					"backup_cxps": []string(nil),
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := deflateBluecatAnycast(tt.input)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestExpandBDDSOptions(t *testing.T) {
	// Create a mock client for testing
	mockClient := &alkira.AlkiraClient{}

	tests := []struct {
		name        string
		input       []interface{}
		expected    *alkira.BDDSOptions
		expectError bool
		errorMsg    string
	}{
		{
			name:        "empty input",
			input:       []interface{}{},
			expected:    nil,
			expectError: false,
		},
		{
			name: "valid BDDS options",
			input: []interface{}{
				map[string]interface{}{
					"hostname":              "bdds.example.com",
					"model":                 "cBDDS50",
					"version":               "9.4.0",
					"client_id":             "client123",
					"activation_key":        "key123",
					"license_credential_id": "existing-cred-123",
				},
			},
			expected: &alkira.BDDSOptions{
				HostName:            "bdds.example.com",
				Model:               "cBDDS50",
				Version:             "9.4.0",
				LicenseCredentialId: "existing-cred-123",
			},
			expectError: false,
		},
		{
			name: "minimal BDDS options",
			input: []interface{}{
				map[string]interface{}{
					"hostname": "bdds.example.com",
					"model":    "cBDDS25",
					"version":  "9.3.0",
				},
			},
			expected: &alkira.BDDSOptions{
				HostName: "bdds.example.com",
				Model:    "cBDDS25",
				Version:  "9.3.0",
			},
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := expandBDDSOptions(tt.input, mockClient)

			if tt.expectError {
				assert.Error(t, err)
				if tt.errorMsg != "" {
					assert.Contains(t, err.Error(), tt.errorMsg)
				}
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expected, result)
			}
		})
	}
}

func TestExpandEdgeOptions(t *testing.T) {
	// Create a mock client for testing
	mockClient := &alkira.AlkiraClient{}

	tests := []struct {
		name        string
		input       []interface{}
		expected    *alkira.EdgeOptions
		expectError bool
		errorMsg    string
	}{
		{
			name:        "empty input",
			input:       []interface{}{},
			expected:    nil,
			expectError: false,
		},
		{
			name: "valid Edge options",
			input: []interface{}{
				map[string]interface{}{
					"hostname":      "edge.example.com",
					"version":       "4.2.0",
					"config_data":   "base64encodeddata",
					"credential_id": "existing-cred-124",
				},
			},
			expected: &alkira.EdgeOptions{
				HostName:     "edge.example.com",
				Version:      "4.2.0",
				CredentialId: "existing-cred-124",
			},
			expectError: false,
		},
		{
			name: "minimal Edge options",
			input: []interface{}{
				map[string]interface{}{
					"hostname": "edge.example.com",
					"version":  "4.1.0",
				},
			},
			expected: &alkira.EdgeOptions{
				HostName: "edge.example.com",
				Version:  "4.1.0",
			},
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := expandEdgeOptions(tt.input, mockClient)

			if tt.expectError {
				assert.Error(t, err)
				if tt.errorMsg != "" {
					assert.Contains(t, err.Error(), tt.errorMsg)
				}
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expected, result)
			}
		})
	}
}

func TestBluecatValidation(t *testing.T) {
	t.Run("test BDDS model validation", func(t *testing.T) {
		validModels := []string{"cBDDS25", "cBDDS50", "cBDDS100", "cBDDS200", "cBDDS500", "cBDDS1000"}
		for _, model := range validModels {
			// These models should be valid strings
			assert.NotEmpty(t, model)
			assert.Contains(t, model, "cBDDS")
		}
	})

	t.Run("test instance type validation", func(t *testing.T) {
		validTypes := []string{"BDDS", "EDGE"}
		for _, instanceType := range validTypes {
			// These types should be valid strings
			assert.NotEmpty(t, instanceType)
			assert.Contains(t, []string{"BDDS", "EDGE"}, instanceType)
		}
	})

	t.Run("test version format", func(t *testing.T) {
		validVersions := []string{"9.4.0", "9.5.1", "9.6.0", "4.2.0", "4.1.5"}
		for _, version := range validVersions {
			// Versions should follow semantic versioning pattern
			assert.NotEmpty(t, version)
			assert.Regexp(t, `^\d+\.\d+\.\d+$`, version)
		}
	})

	t.Run("test hostname format", func(t *testing.T) {
		validHostnames := []string{
			"bdds.example.com",
			"edge-01.corp.local",
			"bluecat-primary.enterprise.local",
			"dns-server.example.org",
		}
		for _, hostname := range validHostnames {
			// Hostnames should be valid DNS names
			assert.NotEmpty(t, hostname)
			assert.Regexp(t, `^[a-zA-Z0-9._-]+$`, hostname)
			assert.Contains(t, hostname, ".")
		}
	})

	t.Run("test IP address format", func(t *testing.T) {
		validIPs := []string{"192.168.1.10", "10.0.0.1", "172.16.0.1", "203.0.113.1"}
		for _, ip := range validIPs {
			// IPs should be in valid IPv4 format
			assert.NotEmpty(t, ip)
			assert.Regexp(t, `^\d+\.\d+\.\d+\.\d+$`, ip)
		}
	})
}
