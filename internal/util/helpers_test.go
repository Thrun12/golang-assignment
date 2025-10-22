package util

import (
	"database/sql"
	"testing"
)

func TestToNullString(t *testing.T) {
	tests := []struct {
		name     string
		input    *string
		expected sql.NullString
	}{
		{
			name:     "Nil pointer",
			input:    nil,
			expected: sql.NullString{Valid: false},
		},
		{
			name: "Empty string",
			input: func() *string {
				s := ""
				return &s
			}(),
			expected: sql.NullString{Valid: false},
		},
		{
			name: "Valid string",
			input: func() *string {
				s := "Hello World"
				return &s
			}(),
			expected: sql.NullString{String: "Hello World", Valid: true},
		},
		{
			name: "Whitespace only",
			input: func() *string {
				s := "   "
				return &s
			}(),
			expected: sql.NullString{String: "   ", Valid: true},
		},
		{
			name: "String with special characters",
			input: func() *string {
				s := "Test@123!#$%"
				return &s
			}(),
			expected: sql.NullString{String: "Test@123!#$%", Valid: true},
		},
		{
			name: "Unicode string",
			input: func() *string {
				s := "Søholm-Boesen"
				return &s
			}(),
			expected: sql.NullString{String: "Søholm-Boesen", Valid: true},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ToNullString(tt.input)
			if result.Valid != tt.expected.Valid {
				t.Errorf("Expected Valid=%v, got Valid=%v", tt.expected.Valid, result.Valid)
			}
			if result.Valid && result.String != tt.expected.String {
				t.Errorf("Expected String=%q, got String=%q", tt.expected.String, result.String)
			}
		})
	}
}

func TestNullStringToString(t *testing.T) {
	tests := []struct {
		name     string
		input    sql.NullString
		expected string
	}{
		{
			name:     "Invalid NullString",
			input:    sql.NullString{Valid: false},
			expected: "",
		},
		{
			name:     "Valid NullString with empty string",
			input:    sql.NullString{String: "", Valid: true},
			expected: "",
		},
		{
			name:     "Valid NullString with value",
			input:    sql.NullString{String: "Hello World", Valid: true},
			expected: "Hello World",
		},
		{
			name:     "Valid NullString with whitespace",
			input:    sql.NullString{String: "   ", Valid: true},
			expected: "   ",
		},
		{
			name:     "Valid NullString with special characters",
			input:    sql.NullString{String: "Test@123!#$%", Valid: true},
			expected: "Test@123!#$%",
		},
		{
			name:     "Valid NullString with Unicode",
			input:    sql.NullString{String: "Søholm-Boesen", Valid: true},
			expected: "Søholm-Boesen",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := NullStringToString(tt.input)
			if result != tt.expected {
				t.Errorf("Expected %q, got %q", tt.expected, result)
			}
		})
	}
}

func TestRoundToTwoDecimals(t *testing.T) {
	tests := []struct {
		name     string
		input    float64
		expected float64
	}{
		{
			name:     "Already two decimals",
			input:    10.25,
			expected: 10.25,
		},
		{
			name:     "Round down",
			input:    10.244,
			expected: 10.24,
		},
		{
			name:     "Round up",
			input:    10.246,
			expected: 10.25,
		},
		{
			name:     "Round up at boundary",
			input:    10.245,
			expected: 10.25,
		},
		{
			name:     "Zero",
			input:    0.0,
			expected: 0.0,
		},
		{
			name:     "Negative number",
			input:    -10.246,
			expected: -10.25,
		},
		{
			name:     "Very small number",
			input:    0.004,
			expected: 0.0,
		},
		{
			name:     "Very small number round up",
			input:    0.005,
			expected: 0.01,
		},
		{
			name:     "Large number",
			input:    999999.999,
			expected: 1000000.0,
		},
		{
			name:     "Integer",
			input:    42.0,
			expected: 42.0,
		},
		{
			name:     "Many decimals",
			input:    3.141592653589793,
			expected: 3.14,
		},
		{
			name:     "Negative with many decimals",
			input:    -3.141592653589793,
			expected: -3.14,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := RoundToTwoDecimals(tt.input)
			// Use a small epsilon for floating point comparison
			epsilon := 0.0001
			if result < tt.expected-epsilon || result > tt.expected+epsilon {
				t.Errorf("Expected %.2f, got %.2f", tt.expected, result)
			}
		})
	}
}

func TestContainsSkill(t *testing.T) {
	tests := []struct {
		name     string
		skills   []string
		skill    string
		expected bool
	}{
		{
			name:     "Skill exists - exact match",
			skills:   []string{"Go", "Python", "JavaScript"},
			skill:    "Go",
			expected: true,
		},
		{
			name:     "Skill exists - case insensitive",
			skills:   []string{"Go", "Python", "JavaScript"},
			skill:    "go",
			expected: true,
		},
		{
			name:     "Skill exists - uppercase search",
			skills:   []string{"Go", "Python", "JavaScript"},
			skill:    "PYTHON",
			expected: true,
		},
		{
			name:     "Skill exists - mixed case",
			skills:   []string{"Go", "Python", "JavaScript"},
			skill:    "jAvAsCrIpT",
			expected: true,
		},
		{
			name:     "Skill does not exist",
			skills:   []string{"Go", "Python", "JavaScript"},
			skill:    "Rust",
			expected: false,
		},
		{
			name:     "Empty skills array",
			skills:   []string{},
			skill:    "Go",
			expected: false,
		},
		{
			name:     "Nil skills array",
			skills:   nil,
			skill:    "Go",
			expected: false,
		},
		{
			name:     "Empty skill search",
			skills:   []string{"Go", "Python"},
			skill:    "",
			expected: false,
		},
		{
			name:     "Partial match should not work",
			skills:   []string{"JavaScript"},
			skill:    "Java",
			expected: false,
		},
		{
			name:     "Skill with spaces - exact match",
			skills:   []string{"C++", "Objective-C", "Visual Basic"},
			skill:    "Visual Basic",
			expected: true,
		},
		{
			name:     "Skill with spaces - case insensitive",
			skills:   []string{"C++", "Objective-C", "Visual Basic"},
			skill:    "visual basic",
			expected: true,
		},
		{
			name:     "Special characters",
			skills:   []string{"C++", "C#", "F#"},
			skill:    "C#",
			expected: true,
		},
		{
			name:     "Special characters - case insensitive",
			skills:   []string{"C++", "C#", "F#"},
			skill:    "c#",
			expected: true,
		},
		{
			name:     "Hyphenated skill",
			skills:   []string{"Objective-C", "TypeScript"},
			skill:    "objective-c",
			expected: true,
		},
		{
			name:     "Unicode skill name",
			skills:   []string{"Go", "中文", "日本語"},
			skill:    "中文",
			expected: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ContainsSkill(tt.skills, tt.skill)
			if result != tt.expected {
				t.Errorf("Expected %v, got %v for skills=%v, skill=%q", tt.expected, result, tt.skills, tt.skill)
			}
		})
	}
}

func TestToNullStringRoundTrip(t *testing.T) {
	// Test that converting to NullString and back preserves the value
	testStrings := []string{"", "Hello", "Test@123", "Søholm-Boesen", "   "}

	for _, testStr := range testStrings {
		t.Run("RoundTrip: "+testStr, func(t *testing.T) {
			// Empty string should become invalid
			if testStr == "" {
				ptr := &testStr
				nullStr := ToNullString(ptr)
				result := NullStringToString(nullStr)
				if result != "" {
					t.Errorf("Expected empty string after round trip, got %q", result)
				}
			} else {
				ptr := &testStr
				nullStr := ToNullString(ptr)
				result := NullStringToString(nullStr)
				if result != testStr {
					t.Errorf("Expected %q after round trip, got %q", testStr, result)
				}
			}
		})
	}
}
