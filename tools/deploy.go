package tools

import (
	"context"

	"devspace-mcp/executor"

	"github.com/mark3labs/mcp-go/mcp"
)

// DevspaceDeployTool returns the tool definition for deploying
func DevspaceDeployTool() mcp.Tool {
	return mcp.NewTool("devspace_deploy",
		mcp.WithDescription("Deploy the project to Kubernetes using devspace"),
		mcp.WithString("namespace",
			mcp.Description("Kubernetes namespace to deploy to"),
		),
		mcp.WithString("kube_context",
			mcp.Description("Kubernetes context to use"),
		),
		mcp.WithString("profile",
			mcp.Description("Profile to use"),
		),
		mcp.WithBoolean("force_build",
			mcp.Description("Force rebuilding images even if not changed"),
		),
		mcp.WithBoolean("force_deploy",
			mcp.Description("Force redeployment even if not changed"),
		),
		mcp.WithBoolean("skip_build",
			mcp.Description("Skip building images"),
		),
		mcp.WithString("working_dir",
			mcp.Description("Working directory containing devspace.yaml"),
		),
	)
}

// DevspaceDeployHandler handles the deploy command
func DevspaceDeployHandler(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	args := []string{"deploy"}

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
	if req.GetBool("force_build", false) {
		args = append(args, "--force-build")
	}
	if req.GetBool("force_deploy", false) {
		args = append(args, "--force-deploy")
	}
	if req.GetBool("skip_build", false) {
		args = append(args, "--skip-build")
	}

	workingDir := req.GetString("working_dir", "")

	// Deploy can take a while, use long running timeout
	result := executor.ExecuteWithOptions(ctx, executor.LongRunningTimeout, workingDir, args...)

	if !result.Success() {
		return mcp.NewToolResultError(result.FormatOutput()), nil
	}

	return mcp.NewToolResultText(result.FormatOutput()), nil
}
