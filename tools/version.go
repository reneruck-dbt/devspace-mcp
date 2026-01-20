package tools

import (
	"context"

	"devspace-mcp/executor"

	"github.com/mark3labs/mcp-go/mcp"
)

// DevspaceVersionTool returns the tool definition for devspace version
func DevspaceVersionTool() mcp.Tool {
	return mcp.NewTool("devspace_version",
		mcp.WithDescription("Get the devspace CLI version"),
	)
}

// DevspaceVersionHandler handles the devspace version command
func DevspaceVersionHandler(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	result := executor.Execute(ctx, "version")

	if !result.Success() {
		return mcp.NewToolResultError(result.FormatOutput()), nil
	}

	return mcp.NewToolResultText(result.FormatOutput()), nil
}
