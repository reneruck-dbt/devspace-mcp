package executor

import (
	"bytes"
	"context"
	"os/exec"
	"time"
)

// DefaultTimeout is the default command execution timeout
const DefaultTimeout = 120 * time.Second

// LongRunningTimeout is used for build/deploy operations
const LongRunningTimeout = 10 * time.Minute

// Result contains the output from command execution
type Result struct {
	Stdout   string `json:"stdout"`
	Stderr   string `json:"stderr"`
	ExitCode int    `json:"exit_code"`
	Error    string `json:"error,omitempty"`
}

// Execute runs a devspace command with the given arguments
func Execute(ctx context.Context, args ...string) Result {
	return ExecuteWithOptions(ctx, DefaultTimeout, "", args...)
}

// ExecuteWithTimeout runs a devspace command with a custom timeout
func ExecuteWithTimeout(ctx context.Context, timeout time.Duration, args ...string) Result {
	return ExecuteWithOptions(ctx, timeout, "", args...)
}

// ExecuteInDir runs a devspace command in a specific directory
func ExecuteInDir(ctx context.Context, workingDir string, args ...string) Result {
	return ExecuteWithOptions(ctx, DefaultTimeout, workingDir, args...)
}

// ExecuteWithOptions runs a devspace command with custom timeout and working directory
func ExecuteWithOptions(ctx context.Context, timeout time.Duration, workingDir string, args ...string) Result {
	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	cmd := exec.CommandContext(ctx, "devspace", args...)

	if workingDir != "" {
		cmd.Dir = workingDir
	}

	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	err := cmd.Run()

	result := Result{
		Stdout:   stdout.String(),
		Stderr:   stderr.String(),
		ExitCode: 0,
	}

	if err != nil {
		// Check for context cancellation/timeout first
		if ctx.Err() == context.DeadlineExceeded {
			result.ExitCode = -2
			result.Error = "command timed out"
		} else if ctx.Err() == context.Canceled {
			result.ExitCode = -3
			result.Error = "command was cancelled"
		} else if exitErr, ok := err.(*exec.ExitError); ok {
			result.ExitCode = exitErr.ExitCode()
		} else {
			result.ExitCode = -1
			result.Error = err.Error()
		}
	}

	return result
}

// FormatOutput returns a formatted string combining stdout and stderr
func (r Result) FormatOutput() string {
	output := r.Stdout
	if r.Stderr != "" {
		if output != "" {
			output += "\n"
		}
		output += r.Stderr
	}
	if r.Error != "" {
		if output != "" {
			output += "\n"
		}
		output += "Error: " + r.Error
	}
	return output
}

// Success returns true if the command executed successfully
func (r Result) Success() bool {
	return r.ExitCode == 0 && r.Error == ""
}
