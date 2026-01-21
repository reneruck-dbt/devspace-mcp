package tools

import (
	"testing"
)

func TestDevspaceLogsTool(t *testing.T) {
	tool := DevspaceLogsTool()

	// Verify tool name
	if tool.Name != "devspace_logs" {
		t.Errorf("expected tool name 'devspace_logs', got %s", tool.Name)
	}

	// Verify description is set
	if tool.Description == "" {
		t.Error("tool description should not be empty")
	}
}

func TestFilterLines(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		pattern  string
		expected string
	}{
		{
			name:     "simple match",
			input:    "line1\nerror occurred\nline3",
			pattern:  "error",
			expected: "error occurred",
		},
		{
			name:     "case insensitive",
			input:    "line1\nERROR occurred\nline3",
			pattern:  "error",
			expected: "ERROR occurred",
		},
		{
			name:     "multiple matches",
			input:    "error 1\ninfo\nerror 2\nwarn",
			pattern:  "error",
			expected: "error 1\nerror 2",
		},
		{
			name:     "no match",
			input:    "line1\nline2\nline3",
			pattern:  "error",
			expected: "",
		},
		{
			name:     "empty input",
			input:    "",
			pattern:  "error",
			expected: "",
		},
		{
			name:     "empty pattern",
			input:    "line1\nline2",
			pattern:  "",
			expected: "line1\nline2",
		},
		{
			name:     "partial match",
			input:    "line1\nthis is an error message\nline3",
			pattern:  "error",
			expected: "this is an error message",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := filterLines(tt.input, tt.pattern)
			if result != tt.expected {
				t.Errorf("filterLines() = %q, want %q", result, tt.expected)
			}
		})
	}
}

func TestFilterByLevel(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		level    string
		expected string
	}{
		{
			name:     "filter error level",
			input:    "INFO: started\nERROR: failed\nWARN: slow",
			level:    "error",
			expected: "ERROR: failed",
		},
		{
			name:     "filter warn level",
			input:    "INFO: started\nERROR: failed\nWARN: slow\nWARNING: deprecated",
			level:    "warn",
			expected: "WARN: slow\nWARNING: deprecated",
		},
		{
			name:     "filter info level",
			input:    "INFO: started\nERROR: failed\ninfo: completed",
			level:    "info",
			expected: "INFO: started\ninfo: completed",
		},
		{
			name:     "error variants",
			input:    "error occurred\nfailed to connect\npanic: nil pointer\nfatal error",
			level:    "error",
			expected: "error occurred\nfailed to connect\npanic: nil pointer\nfatal error",
		},
		{
			name:     "no match",
			input:    "INFO: line1\nDEBUG: line2",
			level:    "error",
			expected: "",
		},
		{
			name:     "unknown level",
			input:    "ERROR: failed\nINFO: ok",
			level:    "debug",
			expected: "ERROR: failed\nINFO: ok",
		},
		{
			name:     "empty input",
			input:    "",
			level:    "error",
			expected: "",
		},
		{
			name:     "case insensitive",
			input:    "Error: Something went wrong\nINFO: all good",
			level:    "error",
			expected: "Error: Something went wrong",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := filterByLevel(tt.input, tt.level)
			if result != tt.expected {
				t.Errorf("filterByLevel() = %q, want %q", result, tt.expected)
			}
		})
	}
}

func TestFilterByLevelErrorPatterns(t *testing.T) {
	// Test that all error-related keywords are detected
	errorKeywords := []string{"error", "err", "fatal", "panic", "failed", "failure"}
	for _, keyword := range errorKeywords {
		t.Run("detects_"+keyword, func(t *testing.T) {
			input := "This line contains " + keyword
			result := filterByLevel(input, "error")
			if result == "" {
				t.Errorf("filterByLevel should detect %q as error, but got empty result", keyword)
			}
		})
	}
}
