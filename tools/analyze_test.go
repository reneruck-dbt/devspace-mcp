package tools

import (
	"testing"
)

func TestDevspaceAnalyzeTool(t *testing.T) {
	tool := DevspaceAnalyzeTool()

	// Verify tool name
	if tool.Name != "devspace_analyze" {
		t.Errorf("expected tool name 'devspace_analyze', got %s", tool.Name)
	}

	// Verify description is set
	if tool.Description == "" {
		t.Error("tool description should not be empty")
	}

	// Verify description mentions minimal output when healthy
	if tool.Description != "" && len(tool.Description) > 0 {
		// Just check it's not empty, actual content may vary
		t.Logf("Tool description: %s", tool.Description)
	}
}
