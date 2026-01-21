package tools

import (
	"context"
	"fmt"
	"strings"
	"time"

	"devspace-mcp/executor"

	"github.com/mark3labs/mcp-go/mcp"
)

// DevspaceStatusTool returns the tool definition for getting environment status
func DevspaceStatusTool() mcp.Tool {
	return mcp.NewTool("devspace_status",
		mcp.WithDescription("Get comprehensive DevSpace environment health status. Aggregates information from multiple sources including devspace.yaml validation, deployments, analysis, sync paths, and port forwards."),
		mcp.WithString("working_dir",
			mcp.Required(),
			mcp.Description("Working directory containing devspace.yaml"),
		),
		mcp.WithString("namespace",
			mcp.Description("Kubernetes namespace to check"),
		),
	)
}

// DevspaceStatusHandler handles the status command
func DevspaceStatusHandler(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	workingDir := req.GetString("working_dir", "")
	if workingDir == "" {
		return mcp.NewToolResultError("working_dir parameter is required"), nil
	}

	namespace := req.GetString("namespace", "")

	var status strings.Builder
	status.WriteString("# DevSpace Environment Status\n\n")

	// 1. Check devspace.yaml exists
	status.WriteString("## Configuration\n")
	if err := ValidateDevspaceYaml(workingDir); err != nil {
		status.WriteString("❌ " + err.Error() + "\n")
		// If no devspace.yaml, can't continue with other checks
		return mcp.NewToolResultText(status.String()), nil
	}
	status.WriteString("✅ devspace.yaml found\n\n")

	// 2. Get deployments status
	status.WriteString("## Deployments\n")
	args := []string{"list", "deployments"}
	if namespace != "" {
		args = append(args, "--namespace", namespace)
	}
	result := executor.ExecuteInDir(ctx, workingDir, args...)
	if result.Success() {
		output := strings.TrimSpace(result.Stdout)
		if output != "" {
			status.WriteString(output + "\n")
		} else {
			status.WriteString("No deployments found\n")
		}
	} else {
		status.WriteString(fmt.Sprintf("⚠️  Could not fetch deployments: %s\n", result.Stderr))
	}
	status.WriteString("\n")

	// 3. Run analyze
	status.WriteString("## Analysis\n")
	args = []string{"analyze", "--timeout=30"}
	if namespace != "" {
		args = append(args, "--namespace", namespace)
	}
	result = executor.ExecuteWithOptions(ctx, 35*time.Second, workingDir, args...)
	if result.Success() {
		output := strings.TrimSpace(result.Stdout)
		if output == "" {
			status.WriteString("✅ No issues detected\n")
		} else {
			status.WriteString(output + "\n")
		}
	} else {
		status.WriteString(fmt.Sprintf("⚠️  Analysis failed: %s\n", result.Stderr))
	}
	status.WriteString("\n")

	// 4. List configured sync paths
	status.WriteString("## Configured Sync Paths\n")
	result = executor.ExecuteInDir(ctx, workingDir, "list", "sync")
	if result.Success() {
		output := strings.TrimSpace(result.Stdout)
		if output != "" {
			status.WriteString(output + "\n")
		} else {
			status.WriteString("No sync paths configured\n")
		}
	} else {
		status.WriteString(fmt.Sprintf("⚠️  Could not fetch sync paths: %s\n", result.Stderr))
	}
	status.WriteString("\n")

	// 5. List configured ports
	status.WriteString("## Configured Port Forwards\n")
	result = executor.ExecuteInDir(ctx, workingDir, "list", "ports")
	if result.Success() {
		output := strings.TrimSpace(result.Stdout)
		if output != "" {
			status.WriteString(output + "\n")
		} else {
			status.WriteString("No port forwards configured\n")
		}
	} else {
		status.WriteString(fmt.Sprintf("⚠️  Could not fetch port forwards: %s\n", result.Stderr))
	}

	return mcp.NewToolResultText(status.String()), nil
}
