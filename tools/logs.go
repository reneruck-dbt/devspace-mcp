package tools

import (
	"context"
	"fmt"

	"devspace-mcp/executor"

	"github.com/mark3labs/mcp-go/mcp"
)

// DevspaceLogsTool returns the tool definition for getting pod logs
func DevspaceLogsTool() mcp.Tool {
	return mcp.NewTool("devspace_logs",
		mcp.WithDescription("Get logs from a pod in the Kubernetes cluster"),
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
		return mcp.NewToolResultError(result.FormatOutput()), nil
	}

	return mcp.NewToolResultText(result.FormatOutput()), nil
}
