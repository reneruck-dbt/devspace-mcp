package tools

import (
	"testing"
)

func TestDevspaceListPodsTool(t *testing.T) {
	tool := DevspaceListPodsTool()

	// Verify tool name
	if tool.Name != "devspace_list_pods" {
		t.Errorf("expected tool name 'devspace_list_pods', got %s", tool.Name)
	}

	// Verify description is set
	if tool.Description == "" {
		t.Error("tool description should not be empty")
	}

	// Verify namespace is in required parameters
	isRequired := false
	for _, req := range tool.InputSchema.Required {
		if req == "namespace" {
			isRequired = true
			break
		}
	}
	if !isRequired {
		t.Error("namespace parameter should be required")
	}
}

func TestDevspaceListPodsValidation(t *testing.T) {
	tests := []struct {
		name    string
		param   string
		value   string
		wantErr bool
	}{
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
			name:    "valid label selector",
			param:   "label_selector",
			value:   "app=myapp",
			wantErr: false,
		},
		{
			name:    "complex label selector",
			param:   "label_selector",
			value:   "app=myapp,tier=frontend",
			wantErr: false,
		},
		{
			name:    "label selector with flag injection",
			param:   "label_selector",
			value:   "-l malicious",
			wantErr: true,
		},
		{
			name:    "valid field selector",
			param:   "field_selector",
			value:   "status.phase=Running",
			wantErr: false,
		},
		{
			name:    "field selector with flag injection",
			param:   "field_selector",
			value:   "--field-selector=malicious",
			wantErr: true,
		},
		{
			name:    "valid output format - wide",
			param:   "output",
			value:   "wide",
			wantErr: false,
		},
		{
			name:    "valid output format - json",
			param:   "output",
			value:   "json",
			wantErr: false,
		},
		{
			name:    "valid output format - yaml",
			param:   "output",
			value:   "yaml",
			wantErr: false,
		},
		{
			name:    "valid output format - name",
			param:   "output",
			value:   "name",
			wantErr: false,
		},
		{
			name:    "output format with flag injection",
			param:   "output",
			value:   "--output=json",
			wantErr: true,
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
