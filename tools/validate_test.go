package tools

import (
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
