package tools

import (
	"github.com/mark3labs/mcp-go/server"
)

// RegisterAll registers all devspace tools with the MCP server
func RegisterAll(s *server.MCPServer) {
	// Version tool
	s.AddTool(DevspaceVersionTool(), DevspaceVersionHandler)

	// List tools
	s.AddTool(DevspaceListNamespacesTool(), DevspaceListNamespacesHandler)
	s.AddTool(DevspaceListContextsTool(), DevspaceListContextsHandler)
	s.AddTool(DevspaceListDeploymentsTool(), DevspaceListDeploymentsHandler)
	s.AddTool(DevspaceListProfilesTool(), DevspaceListProfilesHandler)
	s.AddTool(DevspaceListVarsTool(), DevspaceListVarsHandler)

	// Print tool
	s.AddTool(DevspacePrintTool(), DevspacePrintHandler)

	// Analyze tool
	s.AddTool(DevspaceAnalyzeTool(), DevspaceAnalyzeHandler)

	// Logs tool
	s.AddTool(DevspaceLogsTool(), DevspaceLogsHandler)

	// Build tool
	s.AddTool(DevspaceBuildTool(), DevspaceBuildHandler)

	// Deploy tool
	s.AddTool(DevspaceDeployTool(), DevspaceDeployHandler)

	// Purge tool
	s.AddTool(DevspacePurgeTool(), DevspacePurgeHandler)

	// Run tool
	s.AddTool(DevspaceRunTool(), DevspaceRunHandler)
}
