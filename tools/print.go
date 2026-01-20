package tools

import (
	"context"

	"devspace-mcp/executor"

	"github.com/mark3labs/mcp-go/mcp"
)

// DevspacePrintTool returns the tool definition for printing resolved config
func DevspacePrintTool() mcp.Tool {
	return mcp.NewTool("devspace_print",
		mcp.WithDescription("Print the resolved devspace configuration as YAML"),
		mcp.WithString("profile",
			mcp.Description("Profile to apply when resolving the configuration"),
		),
		mcp.WithBoolean("skip_info",
			mcp.Description("Only print the configuration without additional info"),
		),
		mcp.WithString("working_dir",
			mcp.Description("Working directory containing devspace.yaml"),
		),
	)
}

// DevspacePrintHandler handles the print command
func DevspacePrintHandler(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	args := []string{"print"}

	if profile := req.GetString("profile", ""); profile != "" {
		if err := ValidateStringParam("profile", profile); err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}
		args = append(args, "--profile", profile)
	}
	if req.GetBool("skip_info", false) {
		args = append(args, "--skip-info")
	}

	workingDir := req.GetString("working_dir", "")

	result := executor.ExecuteInDir(ctx, workingDir, args...)

	if !result.Success() {
		return mcp.NewToolResultError(result.FormatOutput()), nil
	}

	return mcp.NewToolResultText(result.FormatOutput()), nil
}
