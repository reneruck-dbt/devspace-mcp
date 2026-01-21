package tools

import (
	"strings"
	"testing"

	"devspace-mcp/executor"
)

func TestEnhanceError(t *testing.T) {
	tests := []struct {
		name           string
		result         executor.Result
		wantContains   string
		wantSuggestion bool
	}{
		{
			name: "AWS SSO token expired",
			result: executor.Result{
				Stderr:   "error: token has expired, please login again",
				ExitCode: 1,
			},
			wantContains:   "aws sso login",
			wantSuggestion: true,
		},
		{
			name: "Unable to connect to server",
			result: executor.Result{
				Stderr:   "Unable to connect to the server: dial tcp: i/o timeout",
				ExitCode: 1,
			},
			wantContains:   "VPN connection",
			wantSuggestion: true,
		},
		{
			name: "Forbidden error",
			result: executor.Result{
				Stderr:   "Error from server (Forbidden): pods is forbidden",
				ExitCode: 1,
			},
			wantContains:   "RBAC permissions",
			wantSuggestion: true,
		},
		{
			name: "Resource not found",
			result: executor.Result{
				Stderr:   "Error: pod \"test-pod\" not found",
				ExitCode: 1,
			},
			wantContains:   "Verify the namespace",
			wantSuggestion: true,
		},
		{
			name: "Context deadline exceeded",
			result: executor.Result{
				Stderr:   "error: context deadline exceeded",
				ExitCode: 1,
			},
			wantContains:   "timed out",
			wantSuggestion: true,
		},
		{
			name: "No devspace.yaml",
			result: executor.Result{
				Stderr:   "Cannot find a devspace.yaml in the current directory",
				ExitCode: 1,
			},
			wantContains:   "working_dir parameter",
			wantSuggestion: true,
		},
		{
			name: "ImagePullBackOff",
			result: executor.Result{
				Stderr:   "Pod status: ImagePullBackOff",
				ExitCode: 1,
			},
			wantContains:   "pull container image",
			wantSuggestion: true,
		},
		{
			name: "CrashLoopBackOff",
			result: executor.Result{
				Stderr:   "Pod is in CrashLoopBackOff state",
				ExitCode: 1,
			},
			wantContains:   "Check pod logs",
			wantSuggestion: true,
		},
		{
			name: "OOM Killed",
			result: executor.Result{
				Stderr:   "Container was OOMKilled",
				ExitCode: 137,
			},
			wantContains:   "memory limits",
			wantSuggestion: true,
		},
		{
			name: "Certificate expired",
			result: executor.Result{
				Stderr:   "x509: certificate has expired or is not yet valid",
				ExitCode: 1,
			},
			wantContains:   "certificate",
			wantSuggestion: true,
		},
		{
			name: "Unknown error without pattern",
			result: executor.Result{
				Stderr:   "Some random error message",
				ExitCode: 1,
			},
			wantContains:   "Some random error message",
			wantSuggestion: false,
		},
		{
			name: "Success with no error",
			result: executor.Result{
				Stdout:   "Command executed successfully",
				ExitCode: 0,
			},
			wantContains:   "Command executed successfully",
			wantSuggestion: false,
		},
		{
			name: "Error in Error field",
			result: executor.Result{
				Error:    "token has expired",
				ExitCode: 1,
			},
			wantContains:   "aws sso login",
			wantSuggestion: true,
		},
		{
			name: "Case insensitive matching",
			result: executor.Result{
				Stderr:   "Error: TOKEN HAS EXPIRED",
				ExitCode: 1,
			},
			wantContains:   "aws sso login",
			wantSuggestion: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := EnhanceError(tt.result)

			// Check if output contains expected text
			if !strings.Contains(got, tt.wantContains) {
				t.Errorf("EnhanceError() output should contain %q, got:\n%s", tt.wantContains, got)
			}

			// Check if suggestion is present when expected
			hasSuggestion := strings.Contains(got, "💡 Suggestion:")
			if hasSuggestion != tt.wantSuggestion {
				t.Errorf("EnhanceError() suggestion present = %v, want %v\nOutput:\n%s", hasSuggestion, tt.wantSuggestion, got)
			}
		})
	}
}

func TestContainsIgnoreCase(t *testing.T) {
	tests := []struct {
		name    string
		text    string
		pattern string
		want    bool
	}{
		{
			name:    "exact match lowercase",
			text:    "token has expired",
			pattern: "token has expired",
			want:    true,
		},
		{
			name:    "case insensitive match",
			text:    "TOKEN HAS EXPIRED",
			pattern: "token has expired",
			want:    true,
		},
		{
			name:    "pattern in middle of text",
			text:    "Error: token has expired, please login",
			pattern: "token has expired",
			want:    true,
		},
		{
			name:    "pattern not found",
			text:    "some other error",
			pattern: "token has expired",
			want:    false,
		},
		{
			name:    "empty pattern",
			text:    "any text",
			pattern: "",
			want:    true, // strings.Contains returns true for empty pattern
		},
		{
			name:    "empty text",
			text:    "",
			pattern: "something",
			want:    false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := containsIgnoreCase(tt.text, tt.pattern)
			if got != tt.want {
				t.Errorf("containsIgnoreCase(%q, %q) = %v, want %v", tt.text, tt.pattern, got, tt.want)
			}
		})
	}
}
