package tools

import (
	"context"

	"devspace-mcp/executor"

	"github.com/mark3labs/mcp-go/mcp"
)

// DevspacePurgeTool returns the tool definition for purging deployments
func DevspacePurgeTool() mcp.Tool {
	return mcp.NewTool("devspace_purge",
		mcp.WithDescription("WARNING: Destructive operation. Delete all deployed Kubernetes resources for the project. This cannot be undone."),
		mcp.WithString("namespace",
			mcp.Description("Kubernetes namespace"),
		),
		mcp.WithString("kube_context",
			mcp.Description("Kubernetes context to use"),
		),
		mcp.WithString("profile",
			mcp.Description("Profile to use"),
		),
		mcp.WithBoolean("force_purge",
			mcp.Description("Force purge even if resources are in use"),
		),
		mcp.WithString("working_dir",
			mcp.Description("Working directory containing devspace.yaml"),
		),
	)
}

// DevspacePurgeHandler handles the purge command
func DevspacePurgeHandler(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	args := []string{"purge"}

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
	if profile := req.GetString("profile", ""); profile != "" {
		if err := ValidateStringParam("profile", profile); err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}
		args = append(args, "--profile", profile)
	}
	if req.GetBool("force_purge", false) {
		args = append(args, "--force-purge")
	}

	workingDir := req.GetString("working_dir", "")

	result := executor.ExecuteInDir(ctx, workingDir, args...)

	if !result.Success() {
		return mcp.NewToolResultError(result.FormatOutput()), nil
	}

	return mcp.NewToolResultText(result.FormatOutput()), nil
}
