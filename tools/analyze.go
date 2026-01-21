package tools

import (
	"context"
	"fmt"
	"time"

	"devspace-mcp/executor"

	"github.com/mark3labs/mcp-go/mcp"
)

// DevspaceAnalyzeTool returns the tool definition for analyzing a namespace
func DevspaceAnalyzeTool() mcp.Tool {
	return mcp.NewTool("devspace_analyze",
		mcp.WithDescription("Analyze a Kubernetes namespace for potential problems and issues"),
		mcp.WithString("namespace",
			mcp.Description("Kubernetes namespace to analyze"),
		),
		mcp.WithString("kube_context",
			mcp.Description("Kubernetes context to use"),
		),
		mcp.WithBoolean("wait",
			mcp.Description("Wait for pods to be ready before analyzing (default: true)"),
		),
		mcp.WithNumber("timeout",
			mcp.Description("Timeout in seconds (default: 120, max: 600)"),
		),
		mcp.WithString("working_dir",
			mcp.Description("Working directory containing devspace.yaml"),
		),
	)
}

// DevspaceAnalyzeHandler handles the analyze command
func DevspaceAnalyzeHandler(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	args := []string{"analyze"}

	if namespace := req.GetString("namespace", ""); namespace != "" {
		if err := ValidateStringParam("namespace", namespace); err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}
		args = append(args, "--namespace", namespace)
	}
	if kubeContext := req.GetString("kube_context", ""); kubeContext != "" {
		if err := ValidateStringParam("kube_context", kubeContext); err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}
		args = append(args, "--kube-context", kubeContext)
	}

	// Handle wait parameter - default is true, so only add flag if explicitly false
	allArgs := req.GetArguments()
	if _, exists := allArgs["wait"]; exists && !req.GetBool("wait", true) {
		args = append(args, "--wait=false")
	}

	timeout := executor.DefaultTimeout
	timeoutSec := req.GetInt("timeout", 0)
	if timeoutSec > 0 {
		// Cap timeout at 600 seconds (10 minutes)
		if timeoutSec > 600 {
			timeoutSec = 600
		}
		timeout = time.Duration(timeoutSec) * time.Second
		args = append(args, "--timeout", fmt.Sprintf("%d", timeoutSec))
	}

	workingDir := req.GetString("working_dir", "")

	result := executor.ExecuteWithOptions(ctx, timeout, workingDir, args...)

	if !result.Success() {
		return mcp.NewToolResultError(EnhanceError(result)), nil
	}

	return mcp.NewToolResultText(result.FormatOutput()), nil
}
