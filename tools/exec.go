package tools

import (
	"context"
	"time"

	"devspace-mcp/executor"

	"github.com/mark3labs/mcp-go/mcp"
)

// DevspaceExecTool returns the tool definition for executing commands in containers
func DevspaceExecTool() mcp.Tool {
	return mcp.NewTool("devspace_exec",
		mcp.WithDescription("Execute a command in a container using DevSpace. Uses 'devspace enter' with non-interactive mode. Useful for running debugging commands, checking file contents, or testing connectivity inside pods."),
		mcp.WithString("command",
			mcp.Required(),
			mcp.Description("Command to execute in the container (e.g., 'ls -la', 'curl localhost:8080', 'cat /etc/hosts')"),
		),
		mcp.WithString("working_dir",
			mcp.Description("Working directory containing devspace.yaml"),
		),
		mcp.WithString("namespace",
			mcp.Description("Kubernetes namespace"),
		),
		mcp.WithString("pod",
			mcp.Description("Specific pod name to execute command in"),
		),
		mcp.WithString("container",
			mcp.Description("Specific container name within the pod"),
		),
		mcp.WithString("label_selector",
			mcp.Description("Label selector to filter pods (e.g., 'app=myapp')"),
		),
		mcp.WithString("image_selector",
			mcp.Description("Image selector to filter by container image (e.g., 'nginx:latest')"),
		),
		mcp.WithString("workdir",
			mcp.Description("Working directory inside the container where command will be executed"),
		),
	)
}

// DevspaceExecHandler handles the exec command
func DevspaceExecHandler(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	command := req.GetString("command", "")
	if command == "" {
		return mcp.NewToolResultError("command parameter is required"), nil
	}

	// Validate command parameter
	if err := ValidateStringParam("command", command); err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}

	// Build args for devspace enter
	args := []string{"enter", "--tty=false", "--pick=false"}

	// Add namespace if specified
	if namespace := req.GetString("namespace", ""); namespace != "" {
		if err := ValidateStringParam("namespace", namespace); err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}
		args = append(args, "--namespace", namespace)
	}

	// Add pod if specified
	if pod := req.GetString("pod", ""); pod != "" {
		if err := ValidateStringParam("pod", pod); err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}
		args = append(args, "--pod", pod)
	}

	// Add container if specified
	if container := req.GetString("container", ""); container != "" {
		if err := ValidateStringParam("container", container); err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}
		args = append(args, "--container", container)
	}

	// Add label selector if specified
	if labelSelector := req.GetString("label_selector", ""); labelSelector != "" {
		if err := ValidateStringParam("label_selector", labelSelector); err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}
		args = append(args, "--label-selector", labelSelector)
	}

	// Add image selector if specified
	if imageSelector := req.GetString("image_selector", ""); imageSelector != "" {
		if err := ValidateStringParam("image_selector", imageSelector); err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}
		args = append(args, "--image-selector", imageSelector)
	}

	// Add workdir if specified
	if workdir := req.GetString("workdir", ""); workdir != "" {
		if err := ValidateStringParam("workdir", workdir); err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}
		args = append(args, "--workdir", workdir)
	}

	// Add the command after --
	args = append(args, "--", command)

	// Get working directory
	workingDir := req.GetString("working_dir", "")

	// Execute with extended timeout for exec commands
	result := executor.ExecuteWithOptions(ctx, 5*time.Minute, workingDir, args...)

	if !result.Success() {
		return mcp.NewToolResultError(result.FormatOutput()), nil
	}

	return mcp.NewToolResultText(result.FormatOutput()), nil
}
