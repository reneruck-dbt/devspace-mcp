package tools

import (
	"context"

	"devspace-mcp/executor"

	"github.com/mark3labs/mcp-go/mcp"
)

// DevspaceBuildTool returns the tool definition for building images
func DevspaceBuildTool() mcp.Tool {
	return mcp.NewTool("devspace_build",
		mcp.WithDescription("Build all images defined in devspace.yaml"),
		mcp.WithString("namespace",
			mcp.Description("Kubernetes namespace"),
		),
		mcp.WithString("kube_context",
			mcp.Description("Kubernetes context to use"),
		),
		mcp.WithString("profile",
			mcp.Description("Profile to use"),
		),
		mcp.WithBoolean("skip_push",
			mcp.Description("Skip pushing images to registry"),
		),
		mcp.WithString("tag",
			mcp.Description("Tag to use for built images"),
		),
		mcp.WithString("working_dir",
			mcp.Description("Working directory containing devspace.yaml"),
		),
	)
}

// DevspaceBuildHandler handles the build command
func DevspaceBuildHandler(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	args := []string{"build"}

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
	if req.GetBool("skip_push", false) {
		args = append(args, "--skip-push")
	}
	if tag := req.GetString("tag", ""); tag != "" {
		if err := ValidateStringParam("tag", tag); err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}
		args = append(args, "--tag", tag)
	}

	workingDir := req.GetString("working_dir", "")

	// Build can take a while, use long running timeout
	result := executor.ExecuteWithOptions(ctx, executor.LongRunningTimeout, workingDir, args...)

	if !result.Success() {
		return mcp.NewToolResultError(EnhanceError(result)), nil
	}

	return mcp.NewToolResultText(result.FormatOutput()), nil
}
