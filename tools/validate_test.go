package tools

import (
	"os"
	"path/filepath"
	"testing"
)

func TestValidateStringParam(t *testing.T) {
	tests := []struct {
		name      string
		paramName string
		value     string
		wantErr   bool
	}{
		{
			name:      "empty value is valid",
			paramName: "namespace",
			value:     "",
			wantErr:   false,
		},
		{
			name:      "normal value is valid",
			paramName: "namespace",
			value:     "my-namespace",
			wantErr:   false,
		},
		{
			name:      "value starting with dash is invalid",
			paramName: "namespace",
			value:     "-malicious",
			wantErr:   true,
		},
		{
			name:      "value starting with double dash is invalid",
			paramName: "kube_context",
			value:     "--help",
			wantErr:   true,
		},
		{
			name:      "value with dash in middle is valid",
			paramName: "profile",
			value:     "my-profile",
			wantErr:   false,
		},
		{
			name:      "value with numbers is valid",
			paramName: "tag",
			value:     "v1.2.3",
			wantErr:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateStringParam(tt.paramName, tt.value)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateStringParam() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestValidateCommandName(t *testing.T) {
	tests := []struct {
		name    string
		cmdName string
		wantErr bool
	}{
		{
			name:    "empty command name is invalid",
			cmdName: "",
			wantErr: true,
		},
		{
			name:    "simple command name is valid",
			cmdName: "migrate",
			wantErr: false,
		},
		{
			name:    "command with hyphen is valid",
			cmdName: "run-tests",
			wantErr: false,
		},
		{
			name:    "command with underscore is valid",
			cmdName: "run_tests",
			wantErr: false,
		},
		{
			name:    "command with colon is valid",
			cmdName: "db:migrate",
			wantErr: false,
		},
		{
			name:    "command with numbers is valid",
			cmdName: "test123",
			wantErr: false,
		},
		{
			name:    "command starting with dash is invalid",
			cmdName: "-help",
			wantErr: true,
		},
		{
			name:    "command with space is invalid",
			cmdName: "run test",
			wantErr: true,
		},
		{
			name:    "command with semicolon is invalid",
			cmdName: "run;ls",
			wantErr: true,
		},
		{
			name:    "command with ampersand is invalid",
			cmdName: "run&ls",
			wantErr: true,
		},
		{
			name:    "command with pipe is invalid",
			cmdName: "run|ls",
			wantErr: true,
		},
		{
			name:    "command with backtick is invalid",
			cmdName: "run`ls`",
			wantErr: true,
		},
		{
			name:    "command with dollar sign is invalid",
			cmdName: "$HOME",
			wantErr: true,
		},
		{
			name:    "command with parenthesis is invalid",
			cmdName: "run()",
			wantErr: true,
		},
		{
			name:    "uppercase command is valid",
			cmdName: "RunTests",
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateCommandName(tt.cmdName)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateCommandName(%q) error = %v, wantErr %v", tt.cmdName, err, tt.wantErr)
			}
		})
	}
}

func TestIsValidCommandChar(t *testing.T) {
	validChars := "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789-_:"
	invalidChars := " !@#$%^&*()+=[]{}|\\;'\"<>,./?"

	for _, c := range validChars {
		if !isValidCommandChar(c) {
			t.Errorf("isValidCommandChar(%q) = false, want true", c)
		}
	}

	for _, c := range invalidChars {
		if isValidCommandChar(c) {
			t.Errorf("isValidCommandChar(%q) = true, want false", c)
		}
	}
}

func TestValidateDevspaceYaml(t *testing.T) {
	// Create a temporary directory for testing
	tempDir := t.TempDir()

	tests := []struct {
		name       string
		workingDir string
		createFile bool
		wantErr    bool
	}{
		{
			name:       "devspace.yaml exists",
			workingDir: tempDir,
			createFile: true,
			wantErr:    false,
		},
		{
			name:       "devspace.yaml does not exist",
			workingDir: tempDir,
			createFile: false,
			wantErr:    true,
		},
		{
			name:       "invalid directory",
			workingDir: filepath.Join(tempDir, "nonexistent"),
			createFile: false,
			wantErr:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup: create devspace.yaml if needed
			if tt.createFile {
				configPath := filepath.Join(tt.workingDir, "devspace.yaml")
				if err := os.WriteFile(configPath, []byte("version: v2beta1"), 0644); err != nil {
					t.Fatalf("failed to create test file: %v", err)
				}
				defer os.Remove(configPath)
			}

			err := ValidateDevspaceYaml(tt.workingDir)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateDevspaceYaml() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestValidateDevspaceYamlWithCurrentDir(t *testing.T) {
	// Test with empty working_dir (should use current directory)
	// This will likely fail since we're not in a devspace project
	err := ValidateDevspaceYaml("")
	if err == nil {
		t.Log("ValidateDevspaceYaml with empty dir succeeded (devspace.yaml exists in test dir)")
	} else {
		// Expected to fail in most test environments
		if err.Error() == "" {
			t.Error("ValidateDevspaceYaml should return a descriptive error message")
		}
	}
}
