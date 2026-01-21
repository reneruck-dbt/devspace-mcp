package tools

import (
	"testing"
)

func TestDevspaceExecTool(t *testing.T) {
	tool := DevspaceExecTool()

	// Verify tool name
	if tool.Name != "devspace_exec" {
		t.Errorf("expected tool name 'devspace_exec', got %s", tool.Name)
	}

	// Verify description is set
	if tool.Description == "" {
		t.Error("tool description should not be empty")
	}

	// Verify command is in required parameters
	isRequired := false
	for _, req := range tool.InputSchema.Required {
		if req == "command" {
			isRequired = true
			break
		}
	}
	if !isRequired {
		t.Error("command parameter should be required")
	}
}

func TestDevspaceExecValidation(t *testing.T) {
	tests := []struct {
		name    string
		param   string
		value   string
		wantErr bool
	}{
		{
			name:    "valid command",
			param:   "command",
			value:   "ls -la",
			wantErr: false,
		},
		{
			name:    "command with flag injection",
			param:   "command",
			value:   "-malicious",
			wantErr: true,
		},
		{
			name:    "valid namespace",
			param:   "namespace",
			value:   "default",
			wantErr: false,
		},
		{
			name:    "namespace with flag injection",
			param:   "namespace",
			value:   "--help",
			wantErr: true,
		},
		{
			name:    "valid pod name",
			param:   "pod",
			value:   "test-pod-123",
			wantErr: false,
		},
		{
			name:    "pod with flag injection",
			param:   "pod",
			value:   "-malicious",
			wantErr: true,
		},
		{
			name:    "valid container name",
			param:   "container",
			value:   "main",
			wantErr: false,
		},
		{
			name:    "valid label selector",
			param:   "label_selector",
			value:   "app=test",
			wantErr: false,
		},
		{
			name:    "valid image selector",
			param:   "image_selector",
			value:   "nginx:latest",
			wantErr: false,
		},
		{
			name:    "valid workdir",
			param:   "workdir",
			value:   "/app/src",
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateStringParam(tt.param, tt.value)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateStringParam(%q, %q) error = %v, wantErr %v", tt.param, tt.value, err, tt.wantErr)
			}
		})
	}
}
