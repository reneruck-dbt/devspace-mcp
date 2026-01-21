package tools

import (
	"context"
	"fmt"
	"strings"

	"devspace-mcp/executor"

	"github.com/mark3labs/mcp-go/mcp"
)

// DevspaceLogsTool returns the tool definition for getting pod logs
func DevspaceLogsTool() mcp.Tool {
	return mcp.NewTool("devspace_logs",
		mcp.WithDescription("Get logs from a pod in the Kubernetes cluster with optional filtering by text or log level"),
		mcp.WithString("namespace",
			mcp.Description("Kubernetes namespace"),
		),
		mcp.WithString("pod",
			mcp.Description("Specific pod name to get logs from"),
		),
		mcp.WithString("container",
			mcp.Description("Container name within the pod"),
		),
		mcp.WithString("label_selector",
			mcp.Description("Label selector to filter pods (e.g., 'app=myapp')"),
		),
		mcp.WithNumber("lines",
			mcp.Description("Maximum number of lines to return (default: 200, max: 10000)"),
		),
		mcp.WithString("grep",
			mcp.Description("Filter logs to only show lines containing this text (case-insensitive)"),
		),
		mcp.WithString("grep_level",
			mcp.Description("Filter logs by level: error, warn, or info"),
		),
		mcp.WithString("working_dir",
			mcp.Description("Working directory containing devspace.yaml"),
		),
	)
}

// DevspaceLogsHandler handles the logs command
func DevspaceLogsHandler(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	args := []string{"logs"}

	if namespace := req.GetString("namespace", ""); namespace != "" {
		if err := ValidateStringParam("namespace", namespace); err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}
		args = append(args, "--namespace", namespace)
	}
	if pod := req.GetString("pod", ""); pod != "" {
		if err := ValidateStringParam("pod", pod); err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}
		args = append(args, "--pod", pod)
	}
	if container := req.GetString("container", ""); container != "" {
		if err := ValidateStringParam("container", container); err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}
		args = append(args, "--container", container)
	}
	if labelSelector := req.GetString("label_selector", ""); labelSelector != "" {
		if err := ValidateStringParam("label_selector", labelSelector); err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}
		args = append(args, "--label-selector", labelSelector)
	}

	lines := req.GetInt("lines", 200)
	// Ensure lines is positive and capped
	if lines < 1 {
		lines = 200
	} else if lines > 10000 {
		lines = 10000
	}
	args = append(args, "--lines", fmt.Sprintf("%d", lines))

	workingDir := req.GetString("working_dir", "")

	result := executor.ExecuteInDir(ctx, workingDir, args...)

	if !result.Success() {
		return mcp.NewToolResultError(EnhanceError(result)), nil
	}

	// Post-process output with filters if specified
	output := result.Stdout

	// Apply grep filter
	if grepPattern := req.GetString("grep", ""); grepPattern != "" {
		output = filterLines(output, grepPattern)
	}

	// Apply level filter
	if level := req.GetString("grep_level", ""); level != "" {
		output = filterByLevel(output, level)
	}

	return mcp.NewToolResultText(output), nil
}

// filterLines filters log lines that contain the pattern (case-insensitive)
func filterLines(input, pattern string) string {
	if input == "" || pattern == "" {
		return input
	}

	var filtered []string
	lines := strings.Split(input, "\n")
	for _, line := range lines {
		if containsIgnoreCase(line, pattern) {
			filtered = append(filtered, line)
		}
	}
	return strings.Join(filtered, "\n")
}

// filterByLevel filters log lines by log level
func filterByLevel(input, level string) string {
	if input == "" {
		return input
	}

	// Define patterns for each level
	patterns := map[string][]string{
		"error": {"error", "err", "fatal", "panic", "failed", "failure"},
		"warn":  {"warn", "warning"},
		"info":  {"info"},
	}

	levelPatterns, ok := patterns[strings.ToLower(level)]
	if !ok {
		// Unknown level, return input unchanged
		return input
	}

	var filtered []string
	lines := strings.Split(input, "\n")
	for _, line := range lines {
		lineLower := strings.ToLower(line)
		for _, pattern := range levelPatterns {
			if strings.Contains(lineLower, pattern) {
				filtered = append(filtered, line)
				break
			}
		}
	}
	return strings.Join(filtered, "\n")
}
