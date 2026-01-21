package tools

import (
	"testing"
)

func TestDevspaceListPortsTool(t *testing.T) {
	tool := DevspaceListPortsTool()

	// Verify tool name
	if tool.Name != "devspace_list_ports" {
		t.Errorf("expected tool name 'devspace_list_ports', got %s", tool.Name)
	}

	// Verify description is set
	if tool.Description == "" {
		t.Error("tool description should not be empty")
	}

	// Verify working_dir is in required parameters
	isRequired := false
	for _, req := range tool.InputSchema.Required {
		if req == "working_dir" {
			isRequired = true
			break
		}
	}
	if !isRequired {
		t.Error("working_dir parameter should be required")
	}
}
