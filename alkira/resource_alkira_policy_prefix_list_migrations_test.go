package alkira

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func TestResourcePolicyPrefixListV0(t *testing.T) {
	v0Schema := resourcePolicyPrefixListV0()

	if v0Schema == nil {
		t.Fatal("Expected V0 schema to be non-nil")
	}

	// Verify V0 schema has TypeList for prefix and prefix_range
	prefixSchema, ok := v0Schema.Schema["prefix"]
	if !ok {
		t.Fatal("Expected 'prefix' field in V0 schema")
	}

	if prefixSchema.Type != schema.TypeList {
		t.Errorf("Expected prefix to be TypeList in V0, got %v", prefixSchema.Type)
	}

	prefixRangeSchema, ok := v0Schema.Schema["prefix_range"]
	if !ok {
		t.Fatal("Expected 'prefix_range' field in V0 schema")
	}

	if prefixRangeSchema.Type != schema.TypeList {
		t.Errorf("Expected prefix_range to be TypeList in V0, got %v", prefixRangeSchema.Type)
	}
}

func TestResourcePolicyPrefixListStateUpgradeV0(t *testing.T) {
	ctx := context.Background()

	tests := []struct {
		name     string
		input    map[string]interface{}
		validate func(t *testing.T, output map[string]interface{})
	}{
		{
			name: "migrate prefix from TypeList to TypeSet",
			input: map[string]interface{}{
				"name":        "test-list",
				"description": "test description",
				"prefix": []interface{}{
					map[string]interface{}{
						"cidr":        "10.0.0.0/24",
						"description": "Network 1",
					},
					map[string]interface{}{
						"cidr":        "10.0.1.0/24",
						"description": "Network 2",
					},
				},
			},
			validate: func(t *testing.T, output map[string]interface{}) {
				prefix, ok := output["prefix"]
				if !ok {
					t.Fatal("Expected 'prefix' in output")
				}

				prefixList, ok := prefix.([]interface{})
				if !ok {
					t.Fatalf("Expected prefix to be []interface{}, got %T", prefix)
				}

				if len(prefixList) != 2 {
					t.Errorf("Expected 2 prefixes, got %d", len(prefixList))
				}

				// Verify the content was preserved
				cidrs := make(map[string]string)
				for _, p := range prefixList {
					m, ok := p.(map[string]interface{})
					if !ok {
						t.Fatalf("Expected prefix entry to be map, got %T", p)
					}
					cidr, ok := m["cidr"].(string)
					if !ok {
						continue
					}
					desc, _ := m["description"].(string)
					cidrs[cidr] = desc
				}

				if cidrs["10.0.0.0/24"] != "Network 1" {
					t.Errorf("Expected description 'Network 1' for 10.0.0.0/24, got %q", cidrs["10.0.0.0/24"])
				}
				if cidrs["10.0.1.0/24"] != "Network 2" {
					t.Errorf("Expected description 'Network 2' for 10.0.1.0/24, got %q", cidrs["10.0.1.0/24"])
				}
			},
		},
		{
			name: "migrate prefix_range from TypeList to TypeSet",
			input: map[string]interface{}{
				"name": "test-range-list",
				"prefix_range": []interface{}{
					map[string]interface{}{
						"prefix":      "10.0.0.0/16",
						"description": "Range 1",
						"le":          24,
						"ge":          16,
					},
					map[string]interface{}{
						"prefix":      "192.168.0.0/16",
						"description": "Range 2",
						"le":          32,
						"ge":          16,
					},
				},
			},
			validate: func(t *testing.T, output map[string]interface{}) {
				prefixRange, ok := output["prefix_range"]
				if !ok {
					t.Fatal("Expected 'prefix_range' in output")
				}

				rangeList, ok := prefixRange.([]interface{})
				if !ok {
					t.Fatalf("Expected prefix_range to be []interface{}, got %T", prefixRange)
				}

				if len(rangeList) != 2 {
					t.Errorf("Expected 2 prefix_ranges, got %d", len(rangeList))
				}

				// Verify the content was preserved
				ranges := make(map[string]map[string]interface{})
				for _, r := range rangeList {
					m, ok := r.(map[string]interface{})
					if !ok {
						t.Fatalf("Expected prefix_range entry to be map, got %T", r)
					}
					prefix, ok := m["prefix"].(string)
					if !ok {
						continue
					}
					ranges[prefix] = m
				}

				r1 := ranges["10.0.0.0/16"]
				if r1 == nil {
					t.Fatal("Expected 10.0.0.0/16 in ranges")
				}
				if r1["description"] != "Range 1" {
					t.Errorf("Expected description 'Range 1', got %v", r1["description"])
				}
				if toInt(r1["le"]) != 24 {
					t.Errorf("Expected le=24, got %v", r1["le"])
				}

				r2 := ranges["192.168.0.0/16"]
				if r2 == nil {
					t.Fatal("Expected 192.168.0.0/16 in ranges")
				}
				if r2["description"] != "Range 2" {
					t.Errorf("Expected description 'Range 2', got %v", r2["description"])
				}
				if toInt(r2["le"]) != 32 {
					t.Errorf("Expected le=32, got %v", r2["le"])
				}
			},
		},
		{
			name: "handle missing prefix field",
			input: map[string]interface{}{
				"name":   "test-no-prefix",
				"prefix": nil,
				"prefix_range": []interface{}{
					map[string]interface{}{
						"prefix": "10.0.0.0/16",
						"le":     24,
						"ge":     16,
					},
				},
			},
			validate: func(t *testing.T, output map[string]interface{}) {
				// Should not error, just skip the missing prefix
				prefixRange, ok := output["prefix_range"]
				if !ok {
					t.Fatal("Expected 'prefix_range' in output")
				}
				rangeList, ok := prefixRange.([]interface{})
				if !ok || len(rangeList) != 1 {
					t.Errorf("Expected 1 prefix_range, got %d", len(rangeList))
				}
			},
		},
		{
			name: "handle string le/ge values",
			input: map[string]interface{}{
				"name": "test-string-nums",
				"prefix_range": []interface{}{
					map[string]interface{}{
						"prefix": "10.0.0.0/16",
						"le":     "24",
						"ge":     "16",
					},
				},
			},
			validate: func(t *testing.T, output map[string]interface{}) {
				prefixRange, ok := output["prefix_range"]
				if !ok {
					t.Fatal("Expected 'prefix_range' in output")
				}
				rangeList, ok := prefixRange.([]interface{})
				if !ok || len(rangeList) != 1 {
					t.Fatalf("Expected 1 prefix_range, got %d", len(rangeList))
				}
				r := rangeList[0].(map[string]interface{})
				if toInt(r["le"]) != 24 {
					t.Errorf("Expected le=24, got %v", r["le"])
				}
				if toInt(r["ge"]) != 16 {
					t.Errorf("Expected ge=16, got %v", r["ge"])
				}
			},
		},
		{
			name: "handle empty input",
			input: map[string]interface{}{
				"name":         "test-empty",
				"prefix":       []interface{}{},
				"prefix_range": []interface{}{},
			},
			validate: func(t *testing.T, output map[string]interface{}) {
				prefix, ok := output["prefix"].([]interface{})
				if !ok || len(prefix) != 0 {
					t.Errorf("Expected empty prefix list, got %v", prefix)
				}
				prefixRange, ok := output["prefix_range"].([]interface{})
				if !ok || len(prefixRange) != 0 {
					t.Errorf("Expected empty prefix_range list, got %v", prefixRange)
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			output, err := resourcePolicyPrefixListStateUpgradeV0(ctx, tt.input, nil)
			if err != nil {
				t.Fatalf("Unexpected error: %v", err)
			}
			tt.validate(t, output)
		})
	}
}

func TestToInt(t *testing.T) {
	tests := []struct {
		name     string
		input    interface{}
		expected int
	}{
		{"int value", 42, 42},
		{"float64 value", float64(42.7), 42},
		{"string number", "42", 42},
		{"string with spaces", " 42 ", 42},
		{"invalid string", "abc", 0},
		{"zero", 0, 0},
		{"negative int", -5, -5},
		{"negative string", "-5", -5},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := toInt(tt.input)
			if result != tt.expected {
				t.Errorf("toInt(%v) = %d, want %d", tt.input, result, tt.expected)
			}
		})
	}
}
