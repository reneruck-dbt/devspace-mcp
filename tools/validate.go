package tools

import (
	"fmt"
	"strings"
)

// ValidateStringParam performs basic validation on string parameters
// to prevent flag injection (values starting with -)
func ValidateStringParam(name, value string) error {
	if value == "" {
		return nil
	}
	if strings.HasPrefix(value, "-") {
		return fmt.Errorf("invalid %s: value cannot start with '-'", name)
	}
	return nil
}

// ValidateCommandName validates that a command name contains only safe characters
func ValidateCommandName(name string) error {
	if name == "" {
		return fmt.Errorf("command name cannot be empty")
	}
	// Command name cannot start with a dash (would be interpreted as a flag)
	if strings.HasPrefix(name, "-") {
		return fmt.Errorf("invalid command name: cannot start with '-'")
	}
	for _, r := range name {
		if !isValidCommandChar(r) {
			return fmt.Errorf("invalid command name: contains invalid character '%c'", r)
		}
	}
	return nil
}

// isValidCommandChar returns true if the character is valid for a command name
// (alphanumeric, hyphen, underscore, colon)
func isValidCommandChar(r rune) bool {
	return (r >= 'a' && r <= 'z') ||
		(r >= 'A' && r <= 'Z') ||
		(r >= '0' && r <= '9') ||
		r == '-' || r == '_' || r == ':'
}
