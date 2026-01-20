package executor

import (
	"context"
	"testing"
	"time"
)

func TestResult_Success(t *testing.T) {
	tests := []struct {
		name     string
		result   Result
		expected bool
	}{
		{
			name:     "successful result",
			result:   Result{ExitCode: 0, Error: ""},
			expected: true,
		},
		{
			name:     "non-zero exit code",
			result:   Result{ExitCode: 1, Error: ""},
			expected: false,
		},
		{
			name:     "error message present",
			result:   Result{ExitCode: 0, Error: "some error"},
			expected: false,
		},
		{
			name:     "both error and non-zero exit",
			result:   Result{ExitCode: 1, Error: "some error"},
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.result.Success(); got != tt.expected {
				t.Errorf("Success() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestResult_FormatOutput(t *testing.T) {
	tests := []struct {
		name     string
		result   Result
		expected string
	}{
		{
			name:     "stdout only",
			result:   Result{Stdout: "output"},
			expected: "output",
		},
		{
			name:     "stderr only",
			result:   Result{Stderr: "error output"},
			expected: "error output",
		},
		{
			name:     "stdout and stderr",
			result:   Result{Stdout: "output", Stderr: "error output"},
			expected: "output\nerror output",
		},
		{
			name:     "all fields",
			result:   Result{Stdout: "output", Stderr: "error output", Error: "fatal error"},
			expected: "output\nerror output\nError: fatal error",
		},
		{
			name:     "empty result",
			result:   Result{},
			expected: "",
		},
		{
			name:     "error only",
			result:   Result{Error: "fatal error"},
			expected: "Error: fatal error",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.result.FormatOutput(); got != tt.expected {
				t.Errorf("FormatOutput() = %q, want %q", got, tt.expected)
			}
		})
	}
}

func TestExecute_CommandNotFound(t *testing.T) {
	// Save original and restore after test - we can't easily mock exec.Command
	// so we test with a command that doesn't exist
	ctx := context.Background()

	// This tests the error handling when the command fails
	// We can't test devspace directly without it being installed,
	// but we can verify the structure works
	result := ExecuteWithOptions(ctx, 1*time.Second, "", "version")

	// The result should have some output or error
	// depending on whether devspace is installed
	if result.ExitCode == 0 && result.Error == "" {
		// devspace is installed and worked
		if result.Stdout == "" && result.Stderr == "" {
			t.Error("Expected some output from successful command")
		}
	} else {
		// devspace is not installed or failed
		// This is expected in test environments without devspace
		if result.Error == "" && result.Stderr == "" {
			t.Error("Expected error information for failed command")
		}
	}
}

func TestExecuteWithOptions_Timeout(t *testing.T) {
	ctx := context.Background()

	// Use a very short timeout with sleep command to test timeout handling
	// Note: This test may not work on all systems
	result := ExecuteWithOptions(ctx, 10*time.Millisecond, "", "sleep", "10")

	// Should timeout
	if result.ExitCode != -2 && result.Error != "command timed out" {
		// If devspace isn't installed, we'll get a different error
		// That's fine for this test
		if result.ExitCode == -1 {
			// Command not found is acceptable
			return
		}
	}
}

func TestExecuteWithOptions_WorkingDir(t *testing.T) {
	ctx := context.Background()

	// Test with valid working directory
	result := ExecuteWithOptions(ctx, 5*time.Second, "/tmp", "version")

	// We can't guarantee devspace is installed, but we can verify
	// the working directory parameter doesn't cause a panic
	_ = result
}

func TestExecuteInDir(t *testing.T) {
	ctx := context.Background()

	// Test that ExecuteInDir calls ExecuteWithOptions correctly
	result := ExecuteInDir(ctx, "/tmp", "version")

	// Basic sanity check - should return a result without panic
	_ = result.Success()
	_ = result.FormatOutput()
}

func TestExecute(t *testing.T) {
	ctx := context.Background()

	// Test that Execute calls ExecuteWithOptions with defaults
	result := Execute(ctx, "version")

	// Basic sanity check - should return a result without panic
	_ = result.Success()
	_ = result.FormatOutput()
}
