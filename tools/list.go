package tools

import (
	"context"

	"devspace-mcp/executor"

	"github.com/mark3labs/mcp-go/mcp"
)

// DevspaceListNamespacesTool returns the tool definition for listing namespaces
func DevspaceListNamespacesTool() mcp.Tool {
	return mcp.NewTool("devspace_list_namespaces",
		mcp.WithDescription("List Kubernetes namespaces"),
		mcp.WithString("kube_context",
			mcp.Description("Kubernetes context to use"),
		),
	)
}

// DevspaceListNamespacesHandler handles the list namespaces command
func DevspaceListNamespacesHandler(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	args := []string{"list", "namespaces"}

	if kubeContext := req.GetString("kube_context", ""); kubeContext != "" {
		if err := ValidateStringParam("kube_context", kubeContext); err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}
		args = append(args, "--kube-context", kubeContext)
	}

	result := executor.Execute(ctx, args...)

	if !result.Success() {
		return mcp.NewToolResultError(result.FormatOutput()), nil
	}

	return mcp.NewToolResultText(result.FormatOutput()), nil
}

// DevspaceListContextsTool returns the tool definition for listing contexts
func DevspaceListContextsTool() mcp.Tool {
	return mcp.NewTool("devspace_list_contexts",
		mcp.WithDescription("List available Kubernetes contexts"),
	)
}

// DevspaceListContextsHandler handles the list contexts command
func DevspaceListContextsHandler(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	result := executor.Execute(ctx, "list", "contexts")

	if !result.Success() {
		return mcp.NewToolResultError(result.FormatOutput()), nil
	}

	return mcp.NewToolResultText(result.FormatOutput()), nil
}

// DevspaceListDeploymentsTool returns the tool definition for listing deployments
func DevspaceListDeploymentsTool() mcp.Tool {
	return mcp.NewTool("devspace_list_deployments",
		mcp.WithDescription("List deployments and their status"),
		mcp.WithString("namespace",
			mcp.Description("Kubernetes namespace"),
		),
		mcp.WithString("kube_context",
			mcp.Description("Kubernetes context to use"),
		),
		mcp.WithString("working_dir",
			mcp.Description("Working directory containing devspace.yaml"),
		),
	)
}

// DevspaceListDeploymentsHandler handles the list deployments command
func DevspaceListDeploymentsHandler(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	args := []string{"list", "deployments"}

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

	workingDir := req.GetString("working_dir", "")

	result := executor.ExecuteInDir(ctx, workingDir, args...)

	if !result.Success() {
		return mcp.NewToolResultError(result.FormatOutput()), nil
	}

	return mcp.NewToolResultText(result.FormatOutput()), nil
}

// DevspaceListProfilesTool returns the tool definition for listing profiles
func DevspaceListProfilesTool() mcp.Tool {
	return mcp.NewTool("devspace_list_profiles",
		mcp.WithDescription("List available DevSpace profiles from devspace.yaml"),
		mcp.WithString("working_dir",
			mcp.Description("Working directory containing devspace.yaml"),
		),
	)
}

// DevspaceListProfilesHandler handles the list profiles command
func DevspaceListProfilesHandler(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	workingDir := req.GetString("working_dir", "")

	result := executor.ExecuteInDir(ctx, workingDir, "list", "profiles")

	if !result.Success() {
		return mcp.NewToolResultError(result.FormatOutput()), nil
	}

	return mcp.NewToolResultText(result.FormatOutput()), nil
}

// DevspaceListVarsTool returns the tool definition for listing variables
func DevspaceListVarsTool() mcp.Tool {
	return mcp.NewTool("devspace_list_vars",
		mcp.WithDescription("List variables defined in the active devspace configuration"),
		mcp.WithString("profile",
			mcp.Description("Profile to use when resolving variables"),
		),
		mcp.WithString("working_dir",
			mcp.Description("Working directory containing devspace.yaml"),
		),
	)
}

// DevspaceListVarsHandler handles the list vars command
func DevspaceListVarsHandler(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	args := []string{"list", "vars"}

	if profile := req.GetString("profile", ""); profile != "" {
		if err := ValidateStringParam("profile", profile); err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}
		args = append(args, "--profile", profile)
	}

	workingDir := req.GetString("working_dir", "")

	result := executor.ExecuteInDir(ctx, workingDir, args...)

	if !result.Success() {
		return mcp.NewToolResultError(result.FormatOutput()), nil
	}

	return mcp.NewToolResultText(result.FormatOutput()), nil
}
