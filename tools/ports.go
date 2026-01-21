package tools

import (
	"context"

	"devspace-mcp/executor"

	"github.com/mark3labs/mcp-go/mcp"
)

// DevspaceListPortsTool returns the tool definition for listing port forwards
func DevspaceListPortsTool() mcp.Tool {
	return mcp.NewTool("devspace_list_ports",
		mcp.WithDescription("Lists configured port forwarding rules from devspace.yaml. Shows which local ports will be forwarded to which container ports when running devspace dev."),
		mcp.WithString("working_dir",
			mcp.Required(),
			mcp.Description("Working directory containing devspace.yaml"),
		),
		mcp.WithString("output",
			mcp.Description("Output format: table (default) or json"),
		),
	)
}

// DevspaceListPortsHandler handles the list ports command
func DevspaceListPortsHandler(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	workingDir := req.GetString("working_dir", "")
	if workingDir == "" {
		return mcp.NewToolResultError("working_dir parameter is required"), nil
	}

	// Build args
	args := []string{"list", "ports"}

	// Add output format if specified
	if output := req.GetString("output", ""); output == "json" {
		if err := ValidateStringParam("output", output); err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}
		args = append(args, "-o", "json")
	}

	// Execute command
	result := executor.ExecuteInDir(ctx, workingDir, args...)

	if !result.Success() {
		return mcp.NewToolResultError(EnhanceError(result)), nil
	}

	return mcp.NewToolResultText(result.FormatOutput()), nil
}
