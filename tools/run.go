package tools

import (
	"context"
	"strings"

	"devspace-mcp/executor"

	"github.com/mark3labs/mcp-go/mcp"
)

// DevspaceRunTool returns the tool definition for running commands
func DevspaceRunTool() mcp.Tool {
	return mcp.NewTool("devspace_run",
		mcp.WithDescription("Execute a predefined command from devspace.yaml"),
		mcp.WithString("command",
			mcp.Description("Name of the command to run (as defined in devspace.yaml)"),
			mcp.Required(),
		),
		mcp.WithString("args",
			mcp.Description("Arguments to pass to the command (space-separated)"),
		),
		mcp.WithString("working_dir",
			mcp.Description("Working directory containing devspace.yaml"),
		),
	)
}

// DevspaceRunHandler handles the run command
func DevspaceRunHandler(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	command := req.GetString("command", "")
	if command == "" {
		return mcp.NewToolResultError("command parameter is required"), nil
	}

	// Validate command name
	if err := ValidateCommandName(command); err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}

	args := []string{"run", command}

	if argsStr := req.GetString("args", ""); argsStr != "" {
		// Split args by space and append
		extraArgs := strings.Fields(argsStr)
		args = append(args, extraArgs...)
	}

	workingDir := req.GetString("working_dir", "")
	if workingDir != "" {
		if err := ValidateStringParam("working_dir", workingDir); err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}
	}

	result := executor.ExecuteInDir(ctx, workingDir, args...)

	if !result.Success() {
		return mcp.NewToolResultError(result.FormatOutput()), nil
	}

	return mcp.NewToolResultText(result.FormatOutput()), nil
}
