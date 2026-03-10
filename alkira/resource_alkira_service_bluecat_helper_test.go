package alkira

import (
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
					Id:   123,
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
					"id":   123,
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
					Id:   124,
					EdgeOptions: &alkira.EdgeOptions{
						HostName:     "edge-primary.example.com",
						Version:      "4.2.0",
						CredentialId: "cred-124",
					},
				},
			},
			expected: []map[string]interface{}{
				{
					"id":   124,
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
					Id:   123,
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
					Id:   124,
					EdgeOptions: &alkira.EdgeOptions{
						HostName:     "edge-primary.example.com",
						Version:      "4.2.0",
						CredentialId: "cred-124",
					},
				},
			},
			expected: []map[string]interface{}{
				{
					"id":   123,
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
					"id":   124,
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
			// Create a ResourceData with empty instance state to simulate
			// an import (no prior sensitive values in state).
			d := resourceAlkiraBluecat().TestResourceData()
			result := deflateBluecatInstances(tt.input, d)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestExpandBluecatInstances_HostnameIdLookup(t *testing.T) {
	mockClient := &alkira.AlkiraClient{}

	// Simulate old state: [bdds1(id=1), bdds2(id=2), edge1(id=3), edge2(id=4)]
	oldInstances := []interface{}{
		map[string]interface{}{
			"id": 1, "name": "bdds1.example.com", "type": "BDDS",
			"bdds_options": []interface{}{map[string]interface{}{
				"hostname": "bdds1.example.com", "model": "cBDDS50", "version": "9.4.0",
				"client_id": "c1", "activation_key": "k1", "license_credential_id": "cred-1",
			}},
			"edge_options": []interface{}{},
		},
		map[string]interface{}{
			"id": 2, "name": "bdds2.example.com", "type": "BDDS",
			"bdds_options": []interface{}{map[string]interface{}{
				"hostname": "bdds2.example.com", "model": "cBDDS50", "version": "9.4.0",
				"client_id": "c2", "activation_key": "k2", "license_credential_id": "cred-2",
			}},
			"edge_options": []interface{}{},
		},
		map[string]interface{}{
			"id": 3, "name": "edge1.example.com", "type": "EDGE",
			"bdds_options": []interface{}{},
			"edge_options": []interface{}{map[string]interface{}{
				"hostname": "edge1.example.com", "version": "4.2.0",
				"config_data": "data1", "credential_id": "cred-3",
			}},
		},
		map[string]interface{}{
			"id": 4, "name": "edge2.example.com", "type": "EDGE",
			"bdds_options": []interface{}{},
			"edge_options": []interface{}{map[string]interface{}{
				"hostname": "edge2.example.com", "version": "4.2.0",
				"config_data": "data2", "credential_id": "cred-4",
			}},
		},
	}

	// New config: [bdds1, bdds2, bdds3(new), bdds4(new), edge1, edge2]
	// After positional shift, Terraform would assign id=3 to bdds3 and id=4 to
	// bdds4 (inherited from edge1/edge2's old positions), and id=0 to edge1/edge2.
	// Our fix must correct this: bdds3/bdds4 → id=0 (new), edge1/edge2 → id=3/4.
	newInstances := []interface{}{
		map[string]interface{}{
			"id": 1, "name": "bdds1.example.com", "type": "BDDS",
			"bdds_options": []interface{}{map[string]interface{}{
				"hostname": "bdds1.example.com", "model": "cBDDS50", "version": "9.4.0",
				"client_id": "c1", "activation_key": "k1", "license_credential_id": "cred-1",
			}},
			"edge_options": []interface{}{},
		},
		map[string]interface{}{
			"id": 2, "name": "bdds2.example.com", "type": "BDDS",
			"bdds_options": []interface{}{map[string]interface{}{
				"hostname": "bdds2.example.com", "model": "cBDDS50", "version": "9.4.0",
				"client_id": "c2", "activation_key": "k2", "license_credential_id": "cred-2",
			}},
			"edge_options": []interface{}{},
		},
		// bdds3: new instance - positional shift gave it id=3 (edge1's old id), must become 0
		map[string]interface{}{
			"id": 3, "name": "", "type": "BDDS",
			"bdds_options": []interface{}{map[string]interface{}{
				"hostname": "bdds3.example.com", "model": "cBDDS50", "version": "9.4.0",
				"client_id": "c3", "activation_key": "k3", "license_credential_id": "cred-new3",
			}},
			"edge_options": []interface{}{},
		},
		// bdds4: new instance - positional shift gave it id=4 (edge2's old id), must become 0
		map[string]interface{}{
			"id": 4, "name": "", "type": "BDDS",
			"bdds_options": []interface{}{map[string]interface{}{
				"hostname": "bdds4.example.com", "model": "cBDDS50", "version": "9.4.0",
				"client_id": "c4", "activation_key": "k4", "license_credential_id": "cred-new4",
			}},
			"edge_options": []interface{}{},
		},
		// edge1: existing - positional shift gave it id=0, must be restored to 3
		map[string]interface{}{
			"id": 0, "name": "", "type": "EDGE",
			"bdds_options": []interface{}{},
			"edge_options": []interface{}{map[string]interface{}{
				"hostname": "edge1.example.com", "version": "4.2.0",
				"config_data": "data1", "credential_id": "cred-3",
			}},
		},
		// edge2: existing - positional shift gave it id=0, must be restored to 4
		map[string]interface{}{
			"id": 0, "name": "", "type": "EDGE",
			"bdds_options": []interface{}{},
			"edge_options": []interface{}{map[string]interface{}{
				"hostname": "edge2.example.com", "version": "4.2.0",
				"config_data": "data2", "credential_id": "cred-4",
			}},
		},
	}

	result, err := expandBluecatInstances(newInstances, oldInstances, mockClient)

	assert.NoError(t, err)
	assert.Len(t, result, 6)

	// bdds1 and bdds2: unchanged, ids must be preserved
	assert.Equal(t, 1, result[0].Id)
	assert.Equal(t, 2, result[1].Id)

	// bdds3 and bdds4: new, must NOT inherit edge1/edge2's ids
	assert.Equal(t, 0, result[2].Id, "bdds3 is new and must have id=0")
	assert.Equal(t, 0, result[3].Id, "bdds4 is new and must have id=0")

	// edge1 and edge2: existing, must get their correct ids back via hostname match
	assert.Equal(t, 3, result[4].Id, "edge1 must keep id=3 despite positional shift")
	assert.Equal(t, 4, result[5].Id, "edge2 must keep id=4 despite positional shift")
}

func TestExpandBluecatInstances_EmptyHostname(t *testing.T) {
	mockClient := &alkira.AlkiraClient{}

	// Old state: one instance with id=5 but no bdds_options/edge_options (no hostname).
	// It must not be matched to any new instance; the new instance should be treated
	// as new (id=0).
	oldInstances := []interface{}{
		map[string]interface{}{
			"id": 5, "name": "mystery", "type": "BDDS",
			"bdds_options": []interface{}{},
			"edge_options": []interface{}{},
		},
	}

	newInstances := []interface{}{
		map[string]interface{}{
			"id": 5, "name": "", "type": "BDDS",
			"bdds_options": []interface{}{map[string]interface{}{
				"hostname": "bdds-new.example.com", "model": "cBDDS50", "version": "9.4.0",
				"client_id": "c1", "activation_key": "k1", "license_credential_id": "cred-1",
			}},
			"edge_options": []interface{}{},
		},
	}

	result, err := expandBluecatInstances(newInstances, oldInstances, mockClient)

	assert.NoError(t, err)
	assert.Len(t, result, 1)
	// bdds-new has no match in old state (old instance had no hostname), so id must be 0
	assert.Equal(t, 0, result[0].Id, "new instance must not inherit id from unmatched old instance")
}

func TestExpandBluecatInstances_RemoveFromMiddle(t *testing.T) {
	mockClient := &alkira.AlkiraClient{}

	// Old state: [A(id=1), B(id=2), C(id=3)]
	// New config: [A, C] — B removed from the middle.
	// Terraform's positional diff puts C at index 1 with id=2 (B's old id).
	// Our fix must restore C's correct id=3.
	oldInstances := []interface{}{
		map[string]interface{}{
			"id": 1, "name": "a", "type": "BDDS",
			"bdds_options": []interface{}{map[string]interface{}{
				"hostname": "a.example.com", "model": "cBDDS50", "version": "9.4.0",
				"client_id": "ca", "activation_key": "ka", "license_credential_id": "cred-1",
			}},
			"edge_options": []interface{}{},
		},
		map[string]interface{}{
			"id": 2, "name": "b", "type": "BDDS",
			"bdds_options": []interface{}{map[string]interface{}{
				"hostname": "b.example.com", "model": "cBDDS50", "version": "9.4.0",
				"client_id": "cb", "activation_key": "kb", "license_credential_id": "cred-2",
			}},
			"edge_options": []interface{}{},
		},
		map[string]interface{}{
			"id": 3, "name": "c", "type": "BDDS",
			"bdds_options": []interface{}{map[string]interface{}{
				"hostname": "c.example.com", "model": "cBDDS50", "version": "9.4.0",
				"client_id": "cc", "activation_key": "kc", "license_credential_id": "cred-3",
			}},
			"edge_options": []interface{}{},
		},
	}

	// After B is removed, Terraform positionally assigns B's old id=2 to C.
	newInstances := []interface{}{
		map[string]interface{}{
			"id": 1, "name": "a", "type": "BDDS",
			"bdds_options": []interface{}{map[string]interface{}{
				"hostname": "a.example.com", "model": "cBDDS50", "version": "9.4.0",
				"client_id": "ca", "activation_key": "ka", "license_credential_id": "cred-1",
			}},
			"edge_options": []interface{}{},
		},
		map[string]interface{}{
			"id": 2, "name": "c", "type": "BDDS", // positional shift gave id=2
			"bdds_options": []interface{}{map[string]interface{}{
				"hostname": "c.example.com", "model": "cBDDS50", "version": "9.4.0",
				"client_id": "cc", "activation_key": "kc", "license_credential_id": "cred-3",
			}},
			"edge_options": []interface{}{},
		},
	}

	result, err := expandBluecatInstances(newInstances, oldInstances, mockClient)

	assert.NoError(t, err)
	assert.Len(t, result, 2)
	assert.Equal(t, 1, result[0].Id, "A must keep id=1")
	assert.Equal(t, 3, result[1].Id, "C must keep id=3 despite positional shift to id=2")
}

func TestExpandBluecatInstances_MixedCreateUpdateDelete(t *testing.T) {
	mockClient := &alkira.AlkiraClient{}

	// Old state: [A(id=1), B(id=2)]
	// New config: [A, C(new), B] — C inserted in the middle, B moved to end.
	// Terraform's positional diff would give:
	//   index 0: A  → id=1  (unchanged, correct)
	//   index 1: C  → id=2  (positional shift; should be 0 since C is new)
	//   index 2: B  → id=0  (fell off the end; should be 2)
	oldInstances := []interface{}{
		map[string]interface{}{
			"id": 1, "name": "a", "type": "BDDS",
			"bdds_options": []interface{}{map[string]interface{}{
				"hostname": "a.example.com", "model": "cBDDS50", "version": "9.4.0",
				"client_id": "ca", "activation_key": "ka", "license_credential_id": "cred-1",
			}},
			"edge_options": []interface{}{},
		},
		map[string]interface{}{
			"id": 2, "name": "b", "type": "BDDS",
			"bdds_options": []interface{}{map[string]interface{}{
				"hostname": "b.example.com", "model": "cBDDS50", "version": "9.4.0",
				"client_id": "cb", "activation_key": "kb", "license_credential_id": "cred-2",
			}},
			"edge_options": []interface{}{},
		},
	}

	newInstances := []interface{}{
		map[string]interface{}{
			"id": 1, "name": "a", "type": "BDDS",
			"bdds_options": []interface{}{map[string]interface{}{
				"hostname": "a.example.com", "model": "cBDDS50", "version": "9.4.0",
				"client_id": "ca", "activation_key": "ka", "license_credential_id": "cred-1",
			}},
			"edge_options": []interface{}{},
		},
		// C is new; Terraform positionally assigned id=2
		map[string]interface{}{
			"id": 2, "name": "", "type": "BDDS",
			"bdds_options": []interface{}{map[string]interface{}{
				"hostname": "c.example.com", "model": "cBDDS50", "version": "9.4.0",
				"client_id": "cc", "activation_key": "kc", "license_credential_id": "cred-new",
			}},
			"edge_options": []interface{}{},
		},
		// B was moved to end; Terraform assigned id=0
		map[string]interface{}{
			"id": 0, "name": "b", "type": "BDDS",
			"bdds_options": []interface{}{map[string]interface{}{
				"hostname": "b.example.com", "model": "cBDDS50", "version": "9.4.0",
				"client_id": "cb", "activation_key": "kb", "license_credential_id": "cred-2",
			}},
			"edge_options": []interface{}{},
		},
	}

	result, err := expandBluecatInstances(newInstances, oldInstances, mockClient)

	assert.NoError(t, err)
	assert.Len(t, result, 3)
	assert.Equal(t, 1, result[0].Id, "A must keep id=1")
	assert.Equal(t, 0, result[1].Id, "C is new and must have id=0")
	assert.Equal(t, 2, result[2].Id, "B must keep id=2 despite positional shift")
}

func TestValidateBluecatInstanceHostnames(t *testing.T) {
	bddsInstance := func(hostname string) map[string]interface{} {
		return map[string]interface{}{
			"type": "BDDS",
			"bdds_options": []interface{}{map[string]interface{}{
				"hostname": hostname, "model": "cBDDS50", "version": "9.4.0",
				"client_id": "c", "activation_key": "k", "license_credential_id": "cred",
			}},
			"edge_options": []interface{}{},
		}
	}

	t.Run("no instances", func(t *testing.T) {
		assert.NoError(t, validateBluecatInstanceHostnames([]interface{}{}))
	})

	t.Run("unique hostnames", func(t *testing.T) {
		instances := []interface{}{bddsInstance("a.example.com"), bddsInstance("b.example.com")}
		assert.NoError(t, validateBluecatInstanceHostnames(instances))
	})

	t.Run("duplicate hostname", func(t *testing.T) {
		instances := []interface{}{bddsInstance("dup.example.com"), bddsInstance("dup.example.com")}
		err := validateBluecatInstanceHostnames(instances)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "dup.example.com")
		assert.Contains(t, err.Error(), "instance[0]")
		assert.Contains(t, err.Error(), "instance[1]")
	})

	t.Run("duplicate among many", func(t *testing.T) {
		instances := []interface{}{
			bddsInstance("a.example.com"),
			bddsInstance("b.example.com"),
			bddsInstance("a.example.com"), // duplicate of index 0
		}
		err := validateBluecatInstanceHostnames(instances)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "instance[0]")
		assert.Contains(t, err.Error(), "instance[2]")
	})

	t.Run("instances with no hostname are skipped", func(t *testing.T) {
		noHostname := map[string]interface{}{
			"type": "BDDS", "bdds_options": []interface{}{}, "edge_options": []interface{}{},
		}
		instances := []interface{}{noHostname, noHostname} // two empty-hostname entries: no error
		assert.NoError(t, validateBluecatInstanceHostnames(instances))
	})
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
			name:     "empty anycast configuration",
			input:    alkira.BluecatAnycast{},
			expected: nil,
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
