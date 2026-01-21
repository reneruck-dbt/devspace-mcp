package tools

import (
	"context"
	"os/exec"

	"devspace-mcp/executor"

	"github.com/mark3labs/mcp-go/mcp"
)

// DevspaceListPodsTool returns the tool definition for listing pods
func DevspaceListPodsTool() mcp.Tool {
	return mcp.NewTool("devspace_list_pods",
		mcp.WithDescription("Lists pods in a Kubernetes namespace using kubectl. Useful for inspecting running pods, checking their status, and identifying pod names for use with other tools."),
		mcp.WithString("namespace",
			mcp.Required(),
			mcp.Description("Kubernetes namespace to list pods from"),
		),
		mcp.WithString("label_selector",
			mcp.Description("Filter pods by labels (e.g., 'app=myapp,tier=frontend')"),
		),
		mcp.WithString("field_selector",
			mcp.Description("Filter pods by fields (e.g., 'status.phase=Running')"),
		),
		mcp.WithString("output",
			mcp.Description("Output format: wide, json, yaml, or name (default: wide)"),
		),
		mcp.WithBoolean("all_namespaces",
			mcp.Description("List pods from all namespaces (overrides namespace parameter)"),
		),
	)
}

// DevspaceListPodsHandler handles the list pods command using kubectl
func DevspaceListPodsHandler(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	// Build kubectl command
	args := []string{"get", "pods"}

	// Check if all namespaces flag is set
	allNamespaces := req.GetBool("all_namespaces", false)
	if allNamespaces {
		args = append(args, "--all-namespaces")
	} else {
		// Use specific namespace
		namespace := req.GetString("namespace", "")
		if namespace == "" {
			return mcp.NewToolResultError("namespace parameter is required when all_namespaces is false"), nil
		}
		if err := ValidateStringParam("namespace", namespace); err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}
		args = append(args, "-n", namespace)
	}

	// Add label selector if specified
	if labelSelector := req.GetString("label_selector", ""); labelSelector != "" {
		if err := ValidateStringParam("label_selector", labelSelector); err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}
		args = append(args, "-l", labelSelector)
	}

	// Add field selector if specified
	if fieldSelector := req.GetString("field_selector", ""); fieldSelector != "" {
		if err := ValidateStringParam("field_selector", fieldSelector); err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}
		args = append(args, "--field-selector", fieldSelector)
	}

	// Add output format
	output := req.GetString("output", "wide")
	if err := ValidateStringParam("output", output); err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}
	args = append(args, "-o", output)

	// Execute kubectl command
	result := executeKubectl(ctx, args...)

	if !result.Success() {
		return mcp.NewToolResultError(result.FormatOutput()), nil
	}

	return mcp.NewToolResultText(result.FormatOutput()), nil
}

// executeKubectl runs a kubectl command with the given arguments
func executeKubectl(ctx context.Context, args ...string) executor.Result {
	ctx, cancel := context.WithTimeout(ctx, executor.DefaultTimeout)
	defer cancel()

	cmd := exec.CommandContext(ctx, "kubectl", args...)

	var stdout, stderr string
	if output, err := cmd.CombinedOutput(); err != nil {
		stderr = string(output)
		exitCode := -1
		if exitErr, ok := err.(*exec.ExitError); ok {
			exitCode = exitErr.ExitCode()
		}
		return executor.Result{
			Stdout:   "",
			Stderr:   stderr,
			ExitCode: exitCode,
			Error:    err.Error(),
		}
	} else {
		stdout = string(output)
	}

	return executor.Result{
		Stdout:   stdout,
		Stderr:   stderr,
		ExitCode: 0,
	}
}
